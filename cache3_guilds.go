package disgord

import (
	"github.com/andersfylling/disgord/cache/interfaces"
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/websocket"
	"github.com/pkg/errors"
)

type guildsCache struct {
	items    interfaces.CacheAlger
	users    *usersCache
	channels *channelsCache
	config   *CacheConfig
	pool     Pool
	evt      chan<- *websocket.Event
}

type cachedGuild struct {
	guild    *Guild
	channels []Snowflake
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
		_ = item.Object().(*Guild).copyOverToCache(g)
	}
	c.items.RUnlock()

	return g
}
func (c *guildsCache) Channels(guildID Snowflake) (channels []Snowflake) {
	c.items.RLock()
	if item, exists := c.items.Get(guildID); exists {
		tmp := item.Object().(*cachedGuild).channels
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
		ShardID: FakeShardID,
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
	case EvtGuildDelete:
		f = c.onGuildDelete
	}
	if f == nil {
		return nil, nil
	}

	return f(data, flags)
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
