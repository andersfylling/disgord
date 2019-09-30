package disgord

import (
	"github.com/andersfylling/disgord/crs"
	jp "github.com/buger/jsonparser"
)

type channelsCache struct {
	config *CacheConfig
	items  *crs.LFU
	pool   Pool

	users *usersCache
}

func (c *channelsCache) Del(channelID Snowflake) {
	c.items.Lock()
	c.items.Delete(channelID)
	c.items.Unlock()
}

// copy copies content from a cached channel into a freshly initialised channel
func (c *channelsCache) copy(dst, src *Channel) {
	_ = src.CopyOverTo(dst)
	dst.recipientsIDs = nil

	dst.Recipients = make([]*User, 0, len(src.recipientsIDs))
	for i := range src.recipientsIDs {
		userID := src.recipientsIDs[i]
		if user := c.users.Get(userID); user != nil {
			dst.Recipients = append(dst.Recipients, user.(*User))
		}
		// TODO: what if the user is missing due to a cache being disabled or
		// a cache replacement algorithm have removed it.
		// TODO: implement a weighted cache to reduce the chance
		//  Any way to ensure it? or add a helper method to check?
	}
}
func (c *channelsCache) Get(channelID Snowflake) (channel interface{}) {
	channel = c.pool.Get().(*Channel)
	c.Peek(channelID, func(cached *Channel) {
		c.copy(channel.(*Channel), cached)
	})

	return
}
func (c *channelsCache) Peek(channelID Snowflake, cb func(*Channel)) (exists bool) {
	if cb == nil {
		panic("callback can not be nil")
	}
	c.items.RLock()
	if item, exists := c.items.Get(channelID); exists {
		cb(item.Val.(*Channel))
		exists = true
	}
	c.items.RUnlock()

	return
}
func (c *channelsCache) Edit(channelID Snowflake, cb func(*Channel)) (exists bool) {
	if cb == nil {
		panic("callback can not be nil")
	}
	c.items.Lock()
	if item, exists := c.items.Get(channelID); exists {
		cb(item.Val.(*Channel))
		exists = true
	}
	c.items.Unlock()

	return
}
func (c *channelsCache) Size() (size uint) {
	c.items.RLock()
	size = c.items.Size()
	c.items.RUnlock()

	return
}
func (c *channelsCache) Cap() (cap uint) {
	c.items.RLock()
	cap = c.items.Cap()
	c.items.RUnlock()

	return
}
func (c *channelsCache) ListIDs() (list []Snowflake) {
	c.items.RLock()
	list = c.items.ListIDs()
	c.items.RUnlock()

	return
}

// var _ gatewayCacher = (*channelsCache)(nil)
// var _ restCacher = (*channelsCache)(nil)
var _ BasicCacheRepo = (*channelsCache)(nil)

//////////////////////////////////////////////////////
//
// Event handlers
//
//////////////////////////////////////////////////////

func (c *channelsCache) evtDemultiplexer(evt string, data []byte, flags Flag) (updated interface{}, err error) {
	var f func(data []byte, flag Flag) (interface{}, error)
	switch evt {
	case EvtChannelCreate:
		f = c.onChannelCreate
	case EvtChannelUpdate:
		f = c.onChannelUpdate
	case EvtChannelDelete:
		f = c.onChannelDelete
	case EvtChannelPinsUpdate:
		f = c.onChannelPinsUpdate
	case EvtMessageCreate:
		f = c.onMessageCreate
	//case EvtMessageDelete:
	// Channel.LastMessageID will be incorrect if the last
	// message is deleted. TODO: add a msg ID history?
	//	f = c.onMessageDelete
	case EvtGuildCreate:
		f = c.onGuildCreate
		//case EvtGuildDelete:
		//	f = c.onGuildDelete
		// when a guild is deleted, it should trigger fake ChannelDelete events for each
		// of its channels
	}
	if f == nil {
		return nil, nil
	}

	return f(data, flags)
}

func (c *channelsCache) onChannelCreate(data []byte, flags Flag) (updated interface{}, err error) {
	return c.onChannelUpdate(data, flags)
}

func (c *channelsCache) onGuildCreate(data []byte, flags Flag) (updated interface{}, err error) {
	const key = "channels"
	var channels []*Channel
	if !flags.Ignorecache() {
		channels = make([]*Channel, 0, jsonArrayLen(data, key))
	}
	_, _ = jp.ArrayEach(data, func(value []byte, dataType jp.ValueType, offset int, err error) {
		channel, err := c.onChannelCreate(value, flags)
		if err != nil {
			return
		}

		if channel != nil {
			channels = append(channels, channel.(*Channel))
		}
	}, key)

	return channels, nil
}

func (c *channelsCache) onChannelUpdate(data []byte, flags Flag) (updated interface{}, err error) {
	id, err := jsonGetSnowflake(data, "id")
	if err != nil {
		return nil, nil
	}

	var createChannel = func(channel *Channel) (interface{}, error) {
		if channel == nil {
			channel = &Channel{}
		}

		recipients, _, _, err := jp.Get(data, "recipients")
		if err != nil {
			recipients = nil
		}
		data = jp.Delete(data, "recipients")

		if err = Unmarshal(data, channel); err != nil {
			return nil, err
		}

		_, _ = jp.ArrayEach(recipients, func(value []byte, dataType jp.ValueType, offset int, err error) {
			id, err := jsonGetSnowflake(data, "id")
			if err != nil {
				return
			}

			for i := range channel.recipientsIDs {
				if channel.recipientsIDs[i] == id {
					return
				}
			}

			channel.recipientsIDs = append(channel.recipientsIDs, id)
		})

		if !flags.Ignorecache() {
			updated = &Channel{}
			c.copy(updated.(*Channel), channel)
		}

		return updated, err
	}

	if ok := c.Edit(id, func(channel *Channel) {
		updated, err = createChannel(channel)
	}); !ok {
		updated, err = createChannel(nil)
	}
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (c *channelsCache) onChannelDelete(data []byte, flags Flag) (updated interface{}, err error) {
	if id, err := jsonGetSnowflake(data, "id"); err == nil {
		c.Del(id)
	}
	return // don't really care about errors here
}

func (c *channelsCache) onChannelPinsUpdate(data []byte, flags Flag) (updated interface{}, err error) {
	info := &ChannelPinsUpdate{}
	if err = Unmarshal(data, updated); err != nil {
		return nil, err
	}

	c.Edit(info.ChannelID, func(channel *Channel) {
		channel.LastPinTimestamp = info.LastPinTimestamp
	})

	return info, nil
}

func (c *channelsCache) onMessageCreate(data []byte, flags Flag) (updated interface{}, err error) {
	channelID, err := jsonGetSnowflake(data, "channel_id")
	if err != nil {
		return nil, nil
	}
	msgID, err := jsonGetSnowflake(data, "id")
	if err != nil {
		return nil, nil
	}

	c.Edit(channelID, func(channel *Channel) {
		channel.LastMessageID = msgID
	})

	return nil, nil
}
