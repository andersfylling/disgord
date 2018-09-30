package disgord

import (
	"github.com/andersfylling/disgord/cache/interfaces"
)

func (c *Cache) SetUser(new *User) {
	if c.users == nil || new == nil {
		return
	}

	c.users.Lock()
	defer c.users.Unlock()
	if item, exists := c.users.Get(new.ID); exists {
		if c.conf.Immutable {
			new.copyOverToCache(item.Object())
		} else {
			item.Set(new)
		}
		c.users.RefreshAfterDiscordUpdate(item)
	} else {
		var content interface{}
		if c.conf.Immutable {
			content = new.DeepCopy()
		} else {
			content = new
		}
		c.users.Set(new.ID, c.users.CreateCacheableItem(content))
	}
}

func (c *Cache) GetUser(id Snowflake) (user *User, err error) {
	if c.users == nil {
		err = NewErrorUsingDeactivatedCache("users")
		return
	}

	c.users.RLock()
	defer c.users.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.users.Get(id); !exists {
		err = NewErrorCacheItemNotFound(id)
		return
	}

	if c.conf.Immutable {
		user = result.Object().(*User).DeepCopy().(*User)
	} else {
		user = result.Object().(*User)
	}

	return
}
