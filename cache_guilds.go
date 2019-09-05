package disgord

import (
	"sync"

	"github.com/buger/jsonparser"

	"github.com/andersfylling/disgord/crs"
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/websocket"
	"github.com/andersfylling/djp"
	"github.com/pkg/errors"
)

type cachedGuild struct {
	guild    *Guild
	channels []Snowflake
}

func (c *cachedGuild) transform() {
	// channels to ids
	if len(c.guild.Channels) > 0 {
		channelIDs := make([]Snowflake, 0, len(c.guild.Channels))
		for i := range c.guild.Channels {
			channelIDs = append(channelIDs, c.guild.Channels[i].ID)
		}
		c.guild.Channels = nil

		var unique bool
		for i := range channelIDs {
			unique = true
			for j := range c.channels {
				if channelIDs[i] == c.channels[j] {
					unique = false
					break
				}
			}

			if unique {
				c.channels = append(c.channels, channelIDs[i])
			}
		}
	}
}

type guildsCache struct {
	sync.RWMutex

	items    *crs.LFU
	users    *usersCache
	channels *channelsCache
	config   *CacheConfig
	pool     Pool
	evt      chan<- *websocket.Event
}

func (c *guildsCache) Del(guildID Snowflake) {
	c.items.Lock()
	c.items.Delete(guildID)
	c.items.Unlock()
}
func (c *guildsCache) Get(guildID Snowflake) (guild interface{}) {
	g := c.pool.Get().(*Guild)
	c.items.RLock()
	if item, exists := c.items.Get(guildID); exists {
		_ = item.Val.(*Guild).copyOverToCache(g)
	}
	c.items.RUnlock()

	return g
}
func (c *guildsCache) Channels(guildID Snowflake) (channels []Snowflake) {
	c.items.RLock()
	if item, exists := c.items.Get(guildID); exists {
		tmp := item.Val.(*cachedGuild).channels
		channels = make([]Snowflake, 0, len(tmp))
		copy(channels, tmp)
	}
	c.items.RUnlock()

	return channels
}
func (c *guildsCache) Size() (size uint) {
	c.items.RLock()
	size = c.items.Size()
	c.items.RUnlock()

	return
}
func (c *guildsCache) Cap() (cap uint) {
	c.items.RLock()
	cap = c.items.Cap()
	c.items.RUnlock()

	return
}
func (c *guildsCache) ListIDs() (list []Snowflake) {
	c.items.RLock()
	list = c.items.ListIDs()
	c.items.RUnlock()

	return
}

// var _ gatewayCacher = (*usersCache)(nil)
// var _ restCacher = (*usersCache)(nil)
var _ BasicCacheRepo = (*usersCache)(nil)

//////////////////////////////////////////////////////
//
// Event creators
// custom events to ensure proper caching.
// trigger these in a go routine
//
//////////////////////////////////////////////////////

func (c *guildsCache) triggerChannelDelete(channelID Snowflake) {
	// need a way to ensure that onChannelDelete methods only cares about ID
	info := Channel{
		ID: channelID,
	}

	data, err := httd.Marshal(&info)
	if err != nil {
		return
	}

	c.evt <- &websocket.Event{
		Name:    EvtChannelDelete,
		Data:    data,
		ShardID: MockedShardID,
	}
}

//////////////////////////////////////////////////////
//
// Event handlers
//
//////////////////////////////////////////////////////

func (c *guildsCache) evtDemultiplexer(evt string, data []byte, flags Flag) (updated interface{}, err error) {
	var f func(data []byte, flag Flag) (interface{}, error)
	switch evt {
	case EvtGuildCreate:
		f = c.onGuildCreate
	case EvtGuildDelete:
		f = c.onGuildDelete
	}
	if f == nil {
		return nil, nil
	}

	return f(data, flags)
}

func (c *guildsCache) onGuildCreate(data []byte, flags Flag) (updated interface{}, err error) {
	guildID, err := jsonGetSnowflake(data, "id")
	if err != nil {
		return nil, errors.New("missing guild id")
	}

	c.Lock()
	defer c.Unlock()

	var cg *cachedGuild
	// check if it already exists
	// it should _not_. But that's not an excuse in the discord realm.
	if item, exists := c.items.Get(guildID); exists {
		cg = item.Val.(*cachedGuild)
	} else {
		cg = &cachedGuild{
			guild: NewGuild(),
		}
	}

	// extract channel ids
	var channelIDs []Snowflake
	_, _ = jsonparser.ArrayEach(data, func(d []byte, _ jsonparser.ValueType, _ int, _ error) {
		id, err := jsonparser.GetUnsafeString(d, "id")
		if err != nil {
			return
		}
		var s Snowflake
		if err = s.UnmarshalJSON(jsonparser.StringToBytes(id)); err != nil {
			return
		}
		channelIDs = append(channelIDs, s)
	}, "channels")
	data = jsonparser.Delete(data, "channels")

	// avoid allocating N redundant user objects to the heap
	data = djp.MemberReplaceUserWithID(data, "members")

	if err := Unmarshal(data, cg.guild); err != nil {
		return nil, err
	}
	updated = cg.guild.DeepCopy()
	cg.transform()

	return updated, nil
}

func (c *guildsCache) onGuildUpdate(data []byte, flags Flag) (updated interface{}, err error) {
	guildID, err := jsonGetSnowflake(data, "id")
	if err != nil {
		return nil, errors.New("missing guild id")
	}

	// check if it already exists
	if entry, exists := c.items.Get(guildID); exists {
		return entry.update(data, flags)
	} else {

	}
	if _, exists := c.items.Get(guildID); !exists {
		return c.onGuildUpdate(data, flags)
	}

	// notify that channels were deleted
	channels := c.Channels(guildID)
	if len(channels) == 0 {
		return nil, nil
	}
	for _, channelID := range channels {
		go c.triggerChannelDelete(channelID)
	}

	return nil, nil
}

func (c *guildsCache) onGuildDelete(data []byte, flags Flag) (updated interface{}, err error) {
	// get user data
	guildID, err := jsonGetSnowflake(data, "id")
	if err != nil {
		return nil, errors.New("missing guild id")
	}

	// notify that channels were deleted
	channels := c.Channels(guildID)
	if len(channels) == 0 {
		return nil, nil
	}
	for _, channelID := range channels {
		go c.triggerChannelDelete(channelID)
	}

	return nil, nil
}
