package disgord

import (
	"github.com/andersfylling/disgord/cache/interfaces"
)

func createChannelCacher(conf *CacheConfig) (cacher interfaces.CacheAlger, err error) {
	if !conf.ChannelCaching {
		return nil, nil
	}

	const channelWeight = 1 // MiB. TODO: what is the actual max size?
	limit := conf.ChannelCacheLimitMiB / channelWeight

	cacher, err = constructSpecificCacher(conf.ChannelCacheAlgorithm, limit, conf.ChannelCacheLifetime)
	return
}

type channelCacheItem struct {
	channel    *Channel
	recipients []Snowflake
}

func (c *channelCacheItem) process(channel *Channel, immutable bool) {
	if !immutable {
		c.channel = channel
	}


	c.channel = channel.DeepCopy().(*Channel)
	for i := range channel.Recipients {
		c.recipients = append(c.recipients, channel.Recipients[i].ID)
	}
	c.channel.Recipients = []*User{} // clear
}

func (c *channelCacheItem) build(cache *Cache) (channel *Channel) {
	if !cache.conf.Immutable {
		channel = c.channel
		return
	}

	channel = c.channel.DeepCopy().(*Channel)
	for i := range c.recipients {
		usr, err := cache.GetUser(c.recipients[i]) // handles immutability on it's own
		if err != nil || usr == nil {
			usr = NewUser()
			usr.ID = c.recipients[i]
			// TODO: should this be loaded by REST request?...
			// TODO-2: maybe it can be a cache option to load dead members on read?
		}
		channel.Recipients = append(channel.Recipients, usr)
	}

	return
}

func (c *channelCacheItem) update(fresh *Channel, immutable bool) {
	if !immutable {
		c.channel = fresh
		return
	}

	fresh.copyOverToCache(c.channel)
	if len(fresh.Recipients) == 0 {
		return
	}

	c.recipients = []Snowflake{}
	for i := range fresh.Recipients {
		// TODO: compare user ID/hash and delete/append accordingly.
		c.recipients = append(c.recipients, fresh.Recipients[i].ID)
	}
}

func (c *Cache) SetChannel(new *Channel) {
	if c.channels == nil || new == nil {
		return
	}

	c.channels.Lock()
	defer c.channels.Unlock()
	if item, exists := c.channels.Get(new.ID); exists {
		item.Object().(*channelCacheItem).update(new, c.conf.Immutable)
		c.channels.RefreshAfterDiscordUpdate(item)
	} else {
		content := &channelCacheItem{}
		content.process(new, c.conf.Immutable)
		c.channels.Set(new.ID, c.channels.CreateCacheableItem(content))
	}
}

func (c *Cache) UpdateChannelPin(id Snowflake, timestamp Timestamp) {
	if c.channels == nil || id.Empty() {
		return
	}

	c.channels.Lock()
	defer c.channels.Unlock()
	if item, exists := c.channels.Get(id); exists {
		item.Object().(*channelCacheItem).channel.LastPinTimestamp = timestamp
		c.channels.RefreshAfterDiscordUpdate(item)
	} else {
		// channel does not exist in cache, create a partial channel
		partial := &PartialChannel{ID: id, LastPinTimestamp: timestamp}
		content := &channelCacheItem{}
		content.process(partial, c.conf.Immutable)
		c.channels.Set(id, c.channels.CreateCacheableItem(content))
	}
}

func (c *Cache) GetChannel(id Snowflake) (channel *Channel, err error) {
	if c.channels == nil {
		err = NewErrorUsingDeactivatedCache("channels")
		return
	}

	c.channels.RLock()
	defer c.channels.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.channels.Get(id); !exists {
		err = NewErrorCacheItemNotFound(id)
		return
	}

	channel = result.Object().(*channelCacheItem).build(c)
	return
}

func (c *Cache) DeleteChannel(id Snowflake) {
	c.channels.Lock()
	defer c.channels.Unlock()

	c.channels.Delete(id)
}
