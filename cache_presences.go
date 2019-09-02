package disgord

import (
	"sync"

	"github.com/andersfylling/disgord/crs"

	jp "github.com/buger/jsonparser"
)

// the guild should nil their presence field, and fetch them from here on build
type presencesCache struct {
	items  *crs.LFU
	users  *usersCache
	config *CacheConfig
	pool   Pool // must never be nil !
}

type cachedGuildPresences struct {
	sync.RWMutex
	presences []*UserPresence
}

func (c *cachedGuildPresences) DeepCopy() (list []*UserPresence) {
	c.RLock()
	list = make([]*UserPresence, 0, len(c.presences))
	for _, p := range c.presences {
		list = append(list, p.DeepCopy().(*UserPresence))
	}
	c.RUnlock()

	return list
}
func (c *cachedGuildPresences) del(userID Snowflake) {
	// keep deleting nil, incorrect entries
	// TODO: move all entries to the end instead and then resize - check with benchmark
	for {
		i := -1
		c.Lock()
		for j := range c.presences {
			if c.presences[j] == nil || c.presences[j].User == nil || c.presences[j].userID == userID {
				i = j
				break
			}
		}
		if i == -1 {
			c.Unlock()
			break
		}

		// remove entry
		c.presences[i] = c.presences[len(c.presences)-1]
		c.presences[len(c.presences)-1] = nil
		c.presences = c.presences[:len(c.presences)-1]
		c.Unlock()
	}
}
func (c *cachedGuildPresences) update(data []byte, flags Flag) (updated *UserPresence, err error) {
	var userID Snowflake
	if userID, err = jsonGetSnowflake(data, "user", "id"); err != nil {
		return nil, err
	}
	data = jp.Delete(data, "user", "id")

	var p *UserPresence
	c.Lock()
	for j := range c.presences {
		if c.presences[j].userID == userID {
			p = c.presences[j]
			break
		}
	}
	// if the presence/user has not yet been added, we skip the update step
	if p == nil {
		p = &UserPresence{userID: userID}
		c.presences = append(c.presences, p)
	}

	// TODO: presence.Activities, handle each uniquely to avoid incorrect overwrites.
	_ = Unmarshal(data, p)
	if !flags.Ignorecache() {
		updated = p.DeepCopy().(*UserPresence)
	}
	c.Unlock()

	return updated, nil
}

func (c *presencesCache) handleRESTResponse(obj interface{}) (err error) {
	return nil
}
func (c *presencesCache) Del(guildID Snowflake) {
	c.items.Lock()
	c.items.Delete(guildID)
	c.items.Unlock()
}
func (c *presencesCache) Get(guildID Snowflake) (presences interface{}) {
	c.items.RLock()
	if item, exists := c.items.Get(guildID); exists {
		presences = item.Val.(*cachedGuildPresences)
	}
	c.items.RUnlock()
	if presences == nil {
		return nil
	}

	return presences.(*cachedGuildPresences).DeepCopy()
}
func (c *presencesCache) Size() (size uint) {
	c.items.RLock()
	size = c.items.Size()
	c.items.RUnlock()

	return
}
func (c *presencesCache) Cap() (cap uint) {
	c.items.RLock()
	cap = c.items.Cap()
	c.items.RUnlock()

	return
}
func (c *presencesCache) ListIDs() (list []Snowflake) {
	c.items.RLock()
	list = c.items.ListIDs()
	c.items.RUnlock()

	return
}

// var _ gatewayCacher = (*presencesCache)(nil)
// var _ restCacher = (*presencesCache)(nil)
var _ BasicCacheRepo = (*presencesCache)(nil)

//////////////////////////////////////////////////////
//
// Event handlers
//
//////////////////////////////////////////////////////

func (c *presencesCache) evtDemultiplexer(evt string, data []byte, flags Flag) (updated interface{}, err error) {
	var f func(data []byte, flag Flag) (interface{}, error)
	switch evt {
	case EvtGuildCreate:
		f = c.onGuildCreate
	case EvtGuildDelete:
		f = c.onGuildDelete
	case EvtGuildMemberRemove:
		f = c.onGuildMemberRemove
	case EvtPresenceUpdate:
		f = c.onPresenceUpdate
	}
	if f == nil {
		return nil, nil
	}

	return f(data, flags)
}

func (c *presencesCache) onPresenceUpdate(data []byte, flags Flag) (updated interface{}, err error) {
	var guildID Snowflake
	if guildID, err = jsonGetSnowflake(data, "guild_id"); err != nil {
		return nil, nil
	}

	var ps *cachedGuildPresences
	c.items.RLock()
	if item, exists := c.items.Get(guildID); exists {
		ps = item.Val.(*cachedGuildPresences)
	}
	c.items.RUnlock()
	if ps == nil {
		ps = &cachedGuildPresences{
			presences: make([]*UserPresence, 0, 1),
		}
		if updated, err = ps.update(data, flags); err != nil {
			return nil, err
		}
		container := c.items.CreateCacheableItem(ps)

		c.items.Lock()
		c.items.Set(guildID, container)
		c.items.Unlock()
	}

	return ps.update(data, flags)
}

func (c *presencesCache) onGuildMemberRemove(data []byte, flags Flag) (updated interface{}, err error) {
	var guildID Snowflake
	if guildID, err = jsonGetSnowflake(data, "guild_id"); err != nil {
		return nil, nil
	}

	var ps *cachedGuildPresences
	c.items.RLock()
	if item, exists := c.items.Get(guildID); exists {
		ps = item.Val.(*cachedGuildPresences)
	}
	c.items.RUnlock()
	if ps == nil {
		return nil, nil
	}

	var userID Snowflake
	if userID, err = jsonGetSnowflake(data, "user", "id"); err != nil {
		return nil, nil
	}
	ps.del(userID)
	return nil, nil
}

func (c *presencesCache) onGuildDelete(data []byte, flags Flag) (updated interface{}, err error) {
	var guildID Snowflake
	if guildID, err = jsonGetSnowflake(data, "id"); err != nil {
		return nil, nil
	}

	c.Del(guildID)
	return nil, nil
}

func (c *presencesCache) onGuildCreate(data []byte, flags Flag) (updated interface{}, err error) {
	var guildID Snowflake
	if guildID, err = jsonGetSnowflake(data, "id"); err != nil {
		return nil, nil
	}

	var ps *cachedGuildPresences
	c.items.RLock()
	if item, exists := c.items.Get(guildID); exists {
		// it's discord after all, a presence update might be sent before
		// a guild create event. I don't know if this is true. But I see
		// no reason to risk it.
		ps = item.Val.(*cachedGuildPresences)
	}
	c.items.RUnlock()
	if ps == nil {
		ps = &cachedGuildPresences{}
		if updated, err = ps.update(data, flags); err != nil {
			return nil, err
		}
		container := c.items.CreateCacheableItem(ps)

		c.items.Lock()
		c.items.Set(guildID, container)
		c.items.Unlock()
	}

	_, err = jp.ArrayEach(data, func(presenceData []byte, _ jp.ValueType, _ int, err error) {
		_, _ = ps.update(presenceData, IgnoreCache)
	}, "presences")
	if err != nil {
		return nil, err
	}

	return ps.DeepCopy(), nil
}
