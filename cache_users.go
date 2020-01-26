package disgord

import (
	"github.com/andersfylling/disgord/internal/crs"
	"github.com/andersfylling/disgord/internal/util"

	"github.com/andersfylling/djp"
	jp "github.com/buger/jsonparser"
)

type usersCache struct {
	config *CacheConfig
	items  *crs.LFU
	pool   Pool // must never be nil !
}

func (c *usersCache) Del(userID Snowflake) {
	c.items.Lock()
	c.items.Delete(userID)
	c.items.Unlock()
}
func (c *usersCache) Get(userID Snowflake) (user interface{}) {
	usr := c.pool.Get().(*User)
	c.items.RLock()
	if item, exists := c.items.Get(userID); exists {
		_ = item.Val.(*User).copyOverToCache(usr)
	}
	c.items.RUnlock()

	if usr.IsEmpty() || usr.Partial() {
		c.pool.Put(usr)
		return nil
	}
	return usr
}
func (c *usersCache) Size() (size uint) {
	c.items.RLock()
	size = c.items.Size()
	c.items.RUnlock()

	return
}
func (c *usersCache) Cap() (cap uint) {
	c.items.RLock()
	cap = c.items.Cap()
	c.items.RUnlock()

	return
}
func (c *usersCache) ListIDs() (list []Snowflake) {
	c.items.RLock()
	list = c.items.ListIDs()
	c.items.RUnlock()

	return
}

// Foreach allows you iterate over the users. This is not blocking for the rest of the system
// as it blocks only when it copies or extract data from one user at the time.
// This is faster when you make the cache mutable, but then again that introduces higher
// risk are then involved (race conditions, incorrect cache, etc).
func (c *usersCache) Foreach(cb func(*User)) {
	ids := c.ListIDs()

	for i := range ids {
		if user := c.Get(ids[i]); user != nil {
			cb(user.(*User))
		}
	}
}

// handleUserData takes a byte slice of a user object only
func (c *usersCache) handleUserData(data []byte, flags Flag) (updated *User, err error) {
	// get-user-id
	var id Snowflake
	if id, err = djp.GetSnowflake(data, "id"); err != nil {
		return nil, err
	}
	// end get-user-id

	// update-user
	var usr *User
	if !flags.Ignorecache() {
		// create a copy if we want to extract the user data
		updated = c.pool.Get().(*User)
	}
	c.items.Lock()
	defer c.items.Unlock()
	if item, exists := c.items.Get(id); exists {
		usr = item.Val.(*User)
		err = util.Unmarshal(data, usr)
	} else { // create new entry
		// TODO: check creation/update time on create events - might be able to skip the need to update the user obj
		usr = c.pool.Get().(*User)
		if err = util.Unmarshal(data, usr); err == nil {
			c.items.Set(id, c.items.CreateCacheableItem(usr))
		}
	}

	if err != nil {
		return nil, err
	}
	if !flags.Ignorecache() {
		_ = usr.copyOverToCache(updated)
	}
	// end update-user

	return updated, nil
}

// var _ gatewayCacher = (*usersCache)(nil)
// var _ restCacher = (*usersCache)(nil)
var _ BasicCacheRepo = (*usersCache)(nil)

//////////////////////////////////////////////////////
//
// Event handlers
//
//////////////////////////////////////////////////////

func (c *usersCache) evtDemultiplexer(evt string, data []byte, flags Flag) (updated interface{}, err error) {
	var f func(data []byte, flag Flag) (interface{}, error)
	switch evt {
	// Note! I don't know wtf Discord is on about by sending some user info
	//  in some of these. Some only have ID which is fine, while others have usernames, avatar, etc.
	//  Am I supposed to update the cache with this?
	// TODO: revisit the idea of updating the cache for every damn event type
	case EvtReady:
		f = c.onReady
	case EvtUserUpdate:
		f = c.onUserUpdate
	case EvtGuildMemberAdd:
		f = c.onGuildMemberAdd
	case EvtGuildMemberRemove:
		f = c.onGuildMemberRemove
	case EvtGuildMemberUpdate:
		f = c.onGuildMemberUpdate
	case EvtGuildMembersChunk:
		f = c.onGuildMembersChunk
	case EvtChannelCreate: // when DM only
		f = c.onChannelCreate
	case EvtChannelUpdate:
		f = c.onChannelUpdate
	case EvtGuildBanAdd:
		f = c.onGuildBanAdd
	case EvtGuildBanRemove:
		f = c.onGuildBanRemove

	// Don't see the need to handle user changes for this. Already have Presence update
	//case EvtVoiceStateUpdate:
	//	f = c.onVoiceStateUpdate

	// These are events that triggers often.
	case EvtMessageCreate: // TODO: consider dropping cache updates for users
		f = c.onMessageCreate
	//case EvtMessageUpdate: // does this even hold the user object?
	//	f = c.onMessageUpdate

	// This triggers a few times a second depending on the guild
	case EvtPresenceUpdate:
		f = c.onPresenceUpdate
	}
	if f == nil {
		return nil, nil
	}

	return f(data, flags)
}

func (c *usersCache) onReady(data []byte, flags Flag) (updated interface{}, err error) {
	// get user data
	value, _, _, err := jp.Get(data, "user")
	if err != nil {
		return nil, err
	}

	return c.handleUserData(value, flags)
}

func (c *usersCache) onUserUpdate(data []byte, flags Flag) (updated interface{}, err error) {
	return c.handleUserData(data, flags)
}

func (c *usersCache) onPresenceUpdate(data []byte, flags Flag) (updated interface{}, err error) {
	// get user data
	value, _, _, err := jp.Get(data, "user")
	if err != nil {
		return nil, err
	}

	// if the user data only has 1 field, it holds no change. Just an ID
	keysLen := jsonNumberOfKeys(value)
	if keysLen == 1 {
		var id Snowflake
		if id, err = jsonGetSnowflake(data, "id"); err != nil {
			return nil, err
		}

		if updated = c.Get(id); updated == nil {
			// DO NOT ADD USER TO MEMORY.
			// we do not want partial objects in memory, ever.
			usr := c.pool.Get().(*User)
			usr.ID = id
			return usr, nil
		}
		return updated, nil
	} else {
		return c.handleUserData(value, flags)
	}
}

func (c *usersCache) onGuildMemberUpdate(data []byte, flags Flag) (updated interface{}, err error) {
	// get user data
	value, _, _, err := jp.Get(data, "user")
	if err != nil {
		return nil, err
	}

	// if the user data only has 1 field, it holds no change. Just an ID
	keysLen := jsonNumberOfKeys(value)
	if keysLen == 0 {
		return nil, nil
	} else if keysLen == 1 {
		var id Snowflake
		if id, err = jsonGetSnowflake(data, "id"); err != nil {
			return nil, err
		}

		if updated = c.Get(id); updated == nil {
			// DO NOT ADD USER TO MEMORY.
			// we do not want partial objects in memory, ever.
			usr := c.pool.Get().(*User)
			usr.ID = id
			return usr, nil
		}
		return updated, nil
	} else {
		return c.handleUserData(value, flags)
	}
}

func (c *usersCache) onGuildMemberAdd(data []byte, flags Flag) (updated interface{}, err error) {
	return c.onGuildMemberUpdate(data, flags)
}

func (c *usersCache) onGuildMemberRemove(data []byte, flags Flag) (updated interface{}, err error) {
	return c.onGuildMemberUpdate(data, flags)
}

func (c *usersCache) onGuildMembersChunk(data []byte, flags Flag) (updated interface{}, err error) {
	// get user data
	membersBytes, _, _, err := jp.Get(data, "members")
	if err != nil {
		return nil, err
	}

	var users []*User
	var index int
	if !flags.Ignorecache() {
		users = make([]*User, 0, jsonArrayLen(membersBytes))
		for range users {
			users = append(users, c.pool.Get().(*User))
		}
	}
	c.items.Lock()
	_, err = jp.ArrayEach(membersBytes, func(memberData []byte, dataType jp.ValueType, offset int, err error) {
		usrData, _, _, err := jp.Get(memberData, "user")
		if err != nil || len(usrData) < 3 {
			return
		}

		// get-user-id
		var id Snowflake
		if id, err = jsonGetSnowflake(usrData, "id"); err != nil {
			return
		}
		// end get-user-id

		// update-user
		var usr *User
		if item, exists := c.items.Get(id); exists {
			usr = item.Val.(*User)
			err = Unmarshal(data, usr)
		} else {
			usr = c.pool.Get().(*User)
			if err = Unmarshal(data, usr); err == nil {
				c.items.Set(id, c.items.CreateCacheableItem(usr))
			}
		}

		if err != nil {
			return
		}
		if !flags.Ignorecache() {
			_ = usr.copyOverToCache(users[index])
			index++
		}
	})
	c.items.Unlock()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (c *usersCache) onChannelCreate(data []byte, flags Flag) (updated interface{}, err error) {
	return c.onChannelUpdate(data, flags)
}

func (c *usersCache) onChannelUpdate(data []byte, flags Flag) (updated interface{}, err error) {
	// if the channel is a DM/GroupDM the recipients field is populated
	usersData, _, _, err := jp.Get(data, "recipients")
	if err != nil {
		return nil, err
	}

	var users []*User
	_, err = jp.ArrayEach(usersData, func(usrData []byte, dataType jp.ValueType, offset int, err error) {
		usr, err := c.handleUserData(usrData, flags)
		if err == nil && usr != nil && !flags.Ignorecache() {
			users = append(users, usr)
		}
	})
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (c *usersCache) onGuildBanAdd(data []byte, flags Flag) (updated interface{}, err error) {
	usrData, _, _, err := jp.Get(data, "user")
	if err != nil {
		return nil, err
	}

	return c.handleUserData(usrData, flags)
}

func (c *usersCache) onGuildBanRemove(data []byte, flags Flag) (updated interface{}, err error) {
	usrData, _, _, err := jp.Get(data, "user")
	if err != nil {
		return nil, err
	}

	return c.handleUserData(usrData, flags)
}

func (c *usersCache) onMessageCreate(data []byte, flags Flag) (updated interface{}, err error) {
	usrData, _, _, err := jp.Get(data, "author")
	if err != nil || usrData == nil || len(usrData) < 3 {
		return nil, nil
	}

	// check for webhook message
	hasID, _, _, err := jp.Get(usrData, "id")
	if err != nil || hasID == nil || len(hasID) < 3 {
		return nil, nil
	}

	return c.handleUserData(usrData, flags)
}
