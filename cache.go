package disgord

import (
	"errors"
	"sync"

	"github.com/andersfylling/disgord/internal/crs"
	"github.com/andersfylling/disgord/json"
)

type idHolder struct {
	ID      Snowflake `json:"id"`
	Channel struct {
		ID Snowflake `json:"id"`
	} `json:"channel"`
	Guild struct {
		ID Snowflake `json:"id"`
	} `json:"guild"`
	User struct {
		ID Snowflake `json:"id"`
	} `json:"user"`
	UserID    Snowflake `json:"user_id"`
	GuildID   Snowflake `json:"guild_id"`
	ChannelID Snowflake `json:"channel_id"`
}

func NewCacheLFUImmutable(limitUsers, limitVoiceStates, limitChannels, limitGuilds uint) Cache {
	lfus := &CacheLFUImmutable{
		CurrentUser: &User{},
	}
	crs.SetLimit(&lfus.Users, limitUsers)
	crs.SetLimit(&lfus.VoiceStates, limitVoiceStates)
	crs.SetLimit(&lfus.Channels, limitChannels)
	crs.SetLimit(&lfus.Guilds, limitGuilds)

	return lfus
}

// CacheLFUImmutable cache with CRS support for Users and voice states
// use NewCacheLFUImmutable to instantiate it!
type CacheLFUImmutable struct {
	CacheNop

	shardedMutex struct {
		Guilds      [4]sync.Mutex
		Users       [10]sync.Mutex
		Channels    [5]sync.Mutex
		VoiceStates [12]sync.Mutex
	}

	CurrentUserMu sync.Mutex
	CurrentUser   *User

	Users       crs.LFU
	VoiceStates crs.LFU
	Channels    crs.LFU
	Guilds      crs.LFU
}

var _ Cache = (*CacheLFUImmutable)(nil)

func (c *CacheLFUImmutable) Mutex(repo *crs.LFU, id Snowflake) *sync.Mutex {
	switch repo {
	case &c.Users:
		return &c.shardedMutex.Users[int(id)%len(c.shardedMutex.Users)]
	case &c.Channels:
		return &c.shardedMutex.Channels[int(id)%len(c.shardedMutex.Channels)]
	case &c.Guilds:
		return &c.shardedMutex.Guilds[int(id)%len(c.shardedMutex.Guilds)]
	case &c.VoiceStates:
		return &c.shardedMutex.VoiceStates[int(id)%len(c.shardedMutex.VoiceStates)]
	}
	panic("unknown cache repo")
}

func (c *CacheLFUImmutable) Ready(data []byte) (*Ready, error) {
	c.CurrentUserMu.Lock()
	defer c.CurrentUserMu.Unlock()

	rdy := &Ready{
		User: c.CurrentUser,
	}

	err := json.Unmarshal(data, rdy)
	rdy.User = c.CurrentUser.DeepCopy().(*User)
	c.Patch(rdy)
	return rdy, err
}

func (c *CacheLFUImmutable) ChannelCreate(data []byte) (*ChannelCreate, error) {
	// assumption#1: Create may take place after an update to the channel
	// assumption#2: The set of fields in both ChannelCreate and ChannelUpdate are the same
	// assumption#3: a channel can not change from one type to another (text => news, text => voice)

	wrap := func(c *Channel) *ChannelCreate {
		return &ChannelCreate{Channel: c.DeepCopy().(*Channel)}
	}

	channel := &Channel{}
	if err := json.Unmarshal(data, channel); err != nil {
		return nil, err
	}
	c.Patch(channel)

	c.Channels.Lock()
	defer c.Channels.Unlock()
	if wrapper, exists := c.Channels.Get(channel.ID); exists {
		// don't update it. It might overwrite a update event(!)
		// TODO: timestamps would be helpful here(?)
		//  or some queue of updates
		err := json.Unmarshal(data, wrapper.Val)
		return wrap(channel), err
	}

	item := c.Channels.CreateCacheableItem(channel)
	c.Channels.Set(channel.ID, item)

	return wrap(channel), nil
}

func (c *CacheLFUImmutable) ChannelUpdate(data []byte) (*ChannelUpdate, error) {
	// assumption#1: Create may not take place before an update event
	// assumption#2: The set of fields in both ChannelCreate and ChannelUpdate are the same
	// assumption#3: a channel can not change from one type to another (text => news, text => voice)

	updateChannel := func(channelID Snowflake, item *crs.LFUItem) (*Channel, error) {
		mutex := c.Mutex(&c.Channels, channelID)
		mutex.Lock()
		channel := item.Val.(*Channel)
		if err := json.Unmarshal(data, channel); err != nil {
			return nil, err
		}
		c.Patch(channel)

		channel = channel.DeepCopy().(*Channel)
		mutex.Unlock()

		return channel, nil
	}

	var metadata *idHolder
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}
	channelID := metadata.ID

	c.Channels.RLock()
	item, exists := c.Channels.Get(channelID)
	c.Channels.RUnlock()

	var channel *Channel
	var err error
	if exists {
		if channel, err = updateChannel(channelID, item); err != nil {
			return nil, err
		}
	} else {
		// unlikely
		tmp := &Channel{}
		if err := json.Unmarshal(data, tmp); err != nil {
			return nil, err
		}
		c.Patch(tmp)
		channel = tmp.DeepCopy().(*Channel)
		freshItem := c.Channels.CreateCacheableItem(tmp)

		c.Channels.Lock()
		if existingItem, exists := c.Channels.Get(channelID); !exists {
			c.Channels.Set(channelID, freshItem)
		} else if channel, err = updateChannel(channelID, existingItem); err != nil { // double lock
			return nil, err
		}
		c.Channels.Unlock()
	}

	return &ChannelUpdate{Channel: channel}, nil
}

func (c *CacheLFUImmutable) ChannelDelete(data []byte) (*ChannelDelete, error) {
	cd := &ChannelDelete{}
	if err := json.Unmarshal(data, cd); err != nil {
		return nil, err
	}
	c.Patch(cd)

	c.Channels.Lock()
	defer c.Channels.Unlock()
	c.Channels.Delete(cd.Channel.ID)

	return cd, nil
}

func (c *CacheLFUImmutable) ChannelPinsUpdate(data []byte) (*ChannelPinsUpdate, error) {
	// assumption#1: not sent on deleted pins

	cpu := &ChannelPinsUpdate{}
	if err := json.Unmarshal(data, cpu); err != nil {
		return nil, err
	}
	c.Patch(cpu)

	if cpu.LastPinTimestamp.IsZero() {
		return cpu, nil
	}

	c.Channels.Lock()
	defer c.Channels.Unlock()
	if item, exists := c.Channels.Get(cpu.ChannelID); exists {
		channel := item.Val.(*Channel)
		channel.LastPinTimestamp = cpu.LastPinTimestamp
	}

	return cpu, nil
}

//func (c *CacheLFUImmutable) VoiceStateUpdate(data []byte) (*VoiceStateUpdate, error) {
//	// assumption#1: not sent on deleted pins
//
//	type voiceStateUpdateHolder struct {
//
//	}
//
//	var vsu *VoiceStateUpdate
//	if err := json.UnmarshalUpdate(data, &vsu); err != nil {
//		return nil, err
//	}
//
//	c.Channels.Lock()
//	defer c.Channels.Unlock()
//	if item, exists := c.Channels.Get(cpu.ChannelID); exists {
//		channel := item.Val.(*Channel)
//		channel.LastPinTimestamp = cpu.LastPinTimestamp
//	}
//
//	return cpu, nil
//}

func (c *CacheLFUImmutable) UserUpdate(data []byte) (*UserUpdate, error) {
	update := &UserUpdate{User: c.CurrentUser}

	c.CurrentUserMu.Lock()
	defer c.CurrentUserMu.Unlock()
	if err := json.Unmarshal(data, update); err != nil {
		return nil, err
	}

	update.User = c.CurrentUser.DeepCopy().(*User)
	c.Patch(update)

	return update, nil
}

func (c *CacheLFUImmutable) VoiceServerUpdate(data []byte) (*VoiceServerUpdate, error) {
	vsu := &VoiceServerUpdate{}
	if err := json.Unmarshal(data, vsu); err != nil {
		return nil, err
	}
	c.Patch(vsu)

	return vsu, nil
}

func (c *CacheLFUImmutable) GuildMemberRemove(data []byte) (*GuildMemberRemove, error) {
	gmr := &GuildMemberRemove{}
	if err := json.Unmarshal(data, gmr); err != nil {
		return nil, err
	}
	c.Patch(gmr)

	c.Guilds.Lock()
	defer c.Guilds.Unlock()

	if item, exists := c.Guilds.Get(gmr.GuildID); exists {
		guild := item.Val.(*Guild)

		for i := range guild.Members {
			if guild.Members[i].UserID == gmr.User.ID {
				guild.MemberCount--
				guild.Members[i] = guild.Members[len(guild.Members)-1]
				guild.Members = guild.Members[:len(guild.Members)-1]
				break
			}
		}
	}

	return gmr, nil
}

func (c *CacheLFUImmutable) GuildMemberAdd(data []byte) (*GuildMemberAdd, error) {
	gmr := &GuildMemberAdd{}
	if err := json.Unmarshal(data, gmr); err != nil {
		return nil, err
	}
	c.Patch(gmr)

	userID := gmr.Member.User.ID
	guildID := gmr.Member.GuildID

	c.Users.RLock()
	cachedUser, userExists := c.Users.Get(userID)
	c.Users.RUnlock()

	if userExists {
		// TODO: i assume the user is partial and doesn't hold any real updates
		usr := cachedUser.Val.(*User)
		// if err := json.Unmarshal(data, &Member{User:usr}); err == nil {
		// 	gmr.Member.User = usr.DeepCopy().(*User)
		// }
		gmr.Member.User = usr.DeepCopy().(*User)
	} else {
		c.Users.Lock()
		defer c.Users.Unlock()

		if _, exists := c.Users.Get(userID); !exists {
			usr := c.Users.CreateCacheableItem(gmr.Member.User.DeepCopy().(*User))
			c.Users.Set(userID, usr)
		}
	}

	c.Guilds.RLock()
	item, exists := c.Guilds.Get(guildID)
	c.Guilds.RUnlock()

	if exists {
		mutex := c.Mutex(&c.Guilds, guildID)
		mutex.Lock()
		defer mutex.Unlock()

		guild := item.Val.(*Guild)

		var member *Member
		for i := range guild.Members { // slow... map instead?
			if guild.Members[i].UserID == gmr.Member.User.ID {
				member = guild.Members[i]
				if err := json.Unmarshal(data, member); err != nil {
					return nil, err
				}
				c.Patch(member)
				break
			}
		}
		if member == nil {
			member = gmr.Member.DeepCopy().(*Member)

			guild.Members = append(guild.Members, member)
			guild.MemberCount++
		}
		member.User = nil
	}

	return gmr, nil
}

func (c *CacheLFUImmutable) GuildCreate(data []byte) (*GuildCreate, error) {
	var metadata *idHolder
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}
	guildID := metadata.ID

	c.Guilds.RLock()
	item, exists := c.Guilds.Get(guildID)
	c.Guilds.RUnlock()

	var guild *Guild
	if exists && item.Val.(*Guild).Unavailable {
		// pre-loaded from ready event
		mutex := c.Mutex(&c.Guilds, guildID)
		mutex.Lock()
		defer mutex.Unlock()

		guild = item.Val.(*Guild)
		if err := json.Unmarshal(data, guild); err != nil {
			return nil, err
		}
		guild.Unavailable = false
		c.Patch(guild)

		guild = guild.DeepCopy().(*Guild)
	} else if exists {
		// not pre-loaded from ready event
		if err := json.Unmarshal(data, &guild); err != nil {
			return nil, err
		}
		c.Patch(guild)
	} else {
		// must create it
		if err := json.Unmarshal(data, &guild); err != nil {
			return nil, err
		}
		c.Patch(guild)

		e := c.Guilds.CreateCacheableItem(guild)
		guild = guild.DeepCopy().(*Guild)

		c.Guilds.Lock()
		if _, exists := c.Guilds.Get(guildID); !exists {
			c.Guilds.Set(guildID, e)
		} // TODO: unmarshal if unavailable
		c.Guilds.Unlock()
	}

	return &GuildCreate{Guild: guild}, nil
}

func (c *CacheLFUImmutable) GuildUpdate(data []byte) (*GuildUpdate, error) {
	updateGuild := func(guildID Snowflake, item *crs.LFUItem) (*Guild, error) {
		mutex := c.Mutex(&c.Guilds, guildID)
		mutex.Lock()
		defer mutex.Unlock()

		guild := item.Val.(*Guild)
		if guild.Unavailable {
			guild.Unavailable = false
		}
		if err := json.Unmarshal(data, guild); err != nil {
			return nil, err
		}
		c.Patch(guild)

		return guild.DeepCopy().(*Guild), nil
	}

	var metadata *idHolder
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}
	guildID := metadata.ID

	c.Guilds.RLock()
	item, exists := c.Guilds.Get(guildID)
	c.Guilds.RUnlock()

	var guild *Guild
	var err error
	if exists {
		guild, err = updateGuild(guildID, item)
	} else {
		if err := json.Unmarshal(data, &guild); err != nil {
			return nil, err
		}
		e := c.Guilds.CreateCacheableItem(guild)

		c.Guilds.Lock()
		defer c.Guilds.Unlock()

		if oldItem, exists := c.Guilds.Get(guildID); exists {
			guild, err = updateGuild(guildID, oldItem) // fallback
		} else {
			c.Guilds.Set(guildID, e)
			guild = guild.DeepCopy().(*Guild)
		}
	}

	return &GuildUpdate{Guild: guild}, err
}

func (c *CacheLFUImmutable) GuildDelete(data []byte) (*GuildDelete, error) {
	guildEvt := &GuildDelete{}
	if err := json.Unmarshal(data, guildEvt); err != nil {
		return nil, err
	}
	c.Patch(guildEvt)

	c.Guilds.Lock()
	defer c.Guilds.Unlock()
	c.Guilds.Delete(guildEvt.UnavailableGuild.ID)

	return guildEvt, nil
}

// REST lookup
// func (c *CacheLFUImmutable) GetMessage(channelID, messageID Snowflake) (*Message, error) {
// 	return nil, nil
// }
// func (c *CacheLFUImmutable) GetCurrentUserGuilds(p *GetCurrentUserGuildsParams) ([]*PartialGuild, error) {
// 	return nil, nil
// }
// func (c *CacheLFUImmutable) GetMessages(channel Snowflake, p *GetMessagesParams) ([]*Message, error) {
// 	return nil, nil
// }
// func (c *CacheLFUImmutable) GetMembers(guildID Snowflake, p *GetMembersParams) ([]*Member, error) {
// 	return nil, nil
// }
func (c *CacheLFUImmutable) GetChannel(id Snowflake) (*Channel, error) {
	c.Channels.RLock()
	cachedItem, exists := c.Channels.Get(id)
	c.Channels.RUnlock()

	if exists {
		mutex := c.Mutex(&c.Channels, id)
		mutex.Lock()
		defer mutex.Lock()

		channel := cachedItem.Val.(*Channel)
		return channel.DeepCopy().(*Channel), nil
	}
	return nil, nil
}
func (c *CacheLFUImmutable) GetGuildEmoji(guildID, emojiID Snowflake) (*Emoji, error) {
	c.Guilds.RLock()
	cachedItem, exists := c.Guilds.Get(guildID)
	c.Guilds.RUnlock()

	if exists {
		mutex := c.Mutex(&c.Guilds, guildID)
		mutex.Lock()
		defer mutex.Lock()

		guild := cachedItem.Val.(*Guild)
		emoji, _ := guild.Emoji(emojiID)
		return emoji.DeepCopy().(*Emoji), nil
	}
	return nil, errors.New("guild does not exist")
}
func (c *CacheLFUImmutable) GetGuildEmojis(id Snowflake) ([]*Emoji, error) {
	c.Guilds.RLock()
	cachedItem, exists := c.Guilds.Get(id)
	c.Guilds.RUnlock()

	if exists {
		mutex := c.Mutex(&c.Guilds, id)
		mutex.Lock()
		defer mutex.Lock()

		guild := cachedItem.Val.(*Guild)
		emojis := make([]*Emoji, len(guild.Emojis))
		for i, emoji := range emojis {
			emojis[i] = emoji.DeepCopy().(*Emoji)
		}

		return emojis, nil
	}
	return nil, errors.New("guild does not exist")
}
func (c *CacheLFUImmutable) GetGuild(id Snowflake) (*Guild, error) {
	c.Guilds.RLock()
	cachedItem, exists := c.Guilds.Get(id)
	c.Guilds.RUnlock()

	var guild *Guild
	if exists {
		mutex := c.Mutex(&c.Guilds, id)
		mutex.Lock()
		defer mutex.Lock()

		guild = cachedItem.Val.(*Guild).DeepCopy().(*Guild)
	}

	return guild, nil
}
func (c *CacheLFUImmutable) GetGuildChannels(id Snowflake) ([]*Channel, error) {
	c.Guilds.RLock()
	cachedItem, exists := c.Guilds.Get(id)
	c.Guilds.RUnlock()

	if exists {
		mutex := c.Mutex(&c.Guilds, id)
		mutex.Lock()
		defer mutex.Lock()

		guild := cachedItem.Val.(*Guild)

		channels := make([]*Channel, len(guild.Channels))
		for i, channel := range guild.Channels {
			channels[i] = channel.DeepCopy().(*Channel)
		}

		return channels, nil
	}
	return nil, errors.New("guild does not exist")
}
func (c *CacheLFUImmutable) GetMember(guildID, userID Snowflake) (*Member, error) {
	var user *User
	var member *Member

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		c.Users.RLock()
		user, _ = c.GetUser(userID)
		c.Users.RUnlock()

		if user != nil {
			mutex := c.Mutex(&c.Users, userID)
			mutex.Lock()
			user = user.DeepCopy().(*User)
			mutex.Unlock()
		}
		wg.Done()
	}()

	c.Guilds.RLock()
	cachedItem, exists := c.Guilds.Get(guildID)
	c.Guilds.RUnlock()

	if exists {
		mutex := c.Mutex(&c.Users, userID)
		mutex.Lock()
		defer mutex.Unlock()

		guild := cachedItem.Val.(*Guild)
		member, _ = guild.Member(userID)
		if member != nil {
			member = member.DeepCopy().(*Member)
		}
	}

	wg.Wait()

	if member != nil {
		member.User = user
		return member, nil
	} else {
		return nil, nil
	}
}
func (c *CacheLFUImmutable) GetGuildRoles(guildID Snowflake) ([]*Role, error) {
	c.Guilds.RLock()
	cachedItem, exists := c.Guilds.Get(guildID)
	c.Guilds.RUnlock()

	if exists {
		mutex := c.Mutex(&c.Guilds, guildID)
		mutex.Lock()
		defer mutex.Unlock()

		guild := cachedItem.Val.(*Guild)
		roles := make([]*Role, len(guild.Roles))
		for i, role := range guild.Roles {
			roles[i] = role.DeepCopy().(*Role)
		}

		return roles, nil
	}
	return nil, errors.New("guild does not exist")
}
func (c *CacheLFUImmutable) GetCurrentUser() (*User, error) {
	c.CurrentUserMu.Lock()
	defer c.CurrentUserMu.Unlock()

	return c.CurrentUser.DeepCopy().(*User), nil
}
func (c *CacheLFUImmutable) GetUser(id Snowflake) (*User, error) {
	currentUser := func() *User {
		c.CurrentUserMu.Lock()
		defer c.CurrentUserMu.Unlock()

		if id == c.CurrentUser.ID {
			return c.CurrentUser
		}
		return nil
	}
	// hmmm.. ugly
	if match := currentUser(); match != nil {
		return match.DeepCopy().(*User), nil
	}

	c.Users.RLock()
	item, exists := c.Users.Get(id)
	c.Users.RUnlock()

	var user *User
	if exists {
		mutex := c.Mutex(&c.Users, id)
		mutex.Lock()
		defer mutex.Unlock()

		user = item.Val.(*User).DeepCopy().(*User)
	}

	return user, nil
}
