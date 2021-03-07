package disgord

import (
	"errors"
	"sort"
	"sync"

	"github.com/andersfylling/disgord/internal/crs"
	"github.com/andersfylling/disgord/json"
)

var ErrCacheMiss = errors.New("no matching entry found in cache")

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

type userHolder struct {
	User *User `json:"user"`
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

	// set via disgord.createClient
	// must never be overwritten
	currentUserID Snowflake // dangerous: no verification that id is set

	CurrentUserMu sync.Mutex
	CurrentUser   *User

	Users       crs.LFU
	VoiceStates crs.LFU
	Channels    crs.LFU
	Guilds      crs.LFU
}

var _ Cache = (*CacheLFUImmutable)(nil)

func (c *CacheLFUImmutable) getGuild(id Snowflake) (*crs.LFUItem, bool) {
	c.Guilds.Lock()
	defer c.Guilds.Unlock()
	return c.Guilds.Get(id)
}

func (c *CacheLFUImmutable) createDMChannel(msg *Message) {
	channelID := msg.ChannelID

	c.Channels.Lock()
	defer c.Channels.Unlock()
	if _, exists := c.Channels.Get(channelID); !exists {
		channel := &Channel{
			ID: channelID,
			Recipients: []*User{
				DeepCopy(c.CurrentUser).(*User),
				DeepCopy(msg.Author).(*User),
			},
			LastMessageID: msg.ID,
			Type:          ChannelTypeDM,
		}
		c.Patch(channel)

		item := c.Channels.CreateCacheableItem(channel)
		c.Channels.Set(channel.ID, item)
	}
}

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
	rdy.User = DeepCopy(c.CurrentUser).(*User)
	c.Patch(rdy)
	return rdy, err
}
func (c *CacheLFUImmutable) MessageCreate(data []byte) (*MessageCreate, error) {
	// assumption#1: Bots don't receive Channel Create Gateway Event for DMs

	msg, err := c.CacheNop.MessageCreate(data)
	if err != nil {
		return msg, err
	}

	if msg.Message.IsDirectMessage() {
		c.createDMChannel(msg.Message)
	}

	return msg, nil
}

func (c *CacheLFUImmutable) ChannelCreate(data []byte) (*ChannelCreate, error) {
	// assumption#1: Create may take place after an update to the channel
	// assumption#2: The set of fields in both ChannelCreate and ChannelUpdate are the same
	// assumption#3: a channel can not change from one type to another (text => news, text => voice)
	channel := &Channel{}
	if err := json.Unmarshal(data, channel); err != nil {
		return nil, err
	}
	c.Patch(channel)

	c.Channels.Lock()
	defer c.Channels.Unlock()
	_ = c.saveChannel(channel)
	return &ChannelCreate{Channel: DeepCopy(channel).(*Channel)}, nil
}

func (c *CacheLFUImmutable) saveChannel(channel *Channel) error {
	item := c.Channels.CreateCacheableItem(channel)
	if _, exists := c.Channels.Get(channel.ID); exists {
		return errors.New("already exists")
	}

	c.Channels.Set(channel.ID, item)
	return nil
}

func (c *CacheLFUImmutable) ChannelUpdate(data []byte) (*ChannelUpdate, error) {
	// assumption#1: Create may not take place before an update event
	// assumption#2: The set of fields in both ChannelCreate and ChannelUpdate are the same
	// assumption#3: a channel can not change from one type to another (text => news, text => voice)

	updateChannel := func(channelID Snowflake, item interface{}) (*Channel, error) {
		mutex := c.Mutex(&c.Channels, channelID)
		mutex.Lock()
		defer mutex.Unlock()

		channel := item.(*Channel)
		if err := json.Unmarshal(data, channel); err != nil {
			return nil, err
		}
		c.Patch(channel)

		channel = DeepCopy(channel).(*Channel)
		return channel, nil
	}

	var metadata *idHolder
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}
	channelID := metadata.ID

	var channel *Channel
	var err error
	if channelI, mu := c.get(&c.Channels, channelID); channelI != nil {
		mu.Lock()
		defer mu.Unlock()

		if channel, err = updateChannel(channelID, channelI); err != nil {
			return nil, err
		}
	} else {
		// unlikely
		tmp := &Channel{}
		if err := json.Unmarshal(data, tmp); err != nil {
			return nil, err
		}
		c.Patch(tmp)
		channel = DeepCopy(tmp).(*Channel)
		freshItem := c.Channels.CreateCacheableItem(tmp)

		c.Channels.Lock()
		if existingItem, exists := c.Channels.Get(channelID); !exists {
			c.Channels.Set(channelID, freshItem)
		} else if channel, err = updateChannel(channelID, existingItem.Val); err != nil { // double lock
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

	update.User = DeepCopy(c.CurrentUser).(*User)
	c.Patch(update)

	return update, nil
}

func (c *CacheLFUImmutable) saveUsers(users []*User) error {
	c.Users.Lock()
	defer c.Users.Unlock()

	// as slow as it gets
	for i := range users {
		if _, exists := c.Users.Get(users[i].ID); exists {
			continue
		}

		c.Users.Set(users[i].ID, c.Users.CreateCacheableItem(users[i]))
	}
	return nil
}

func (c *CacheLFUImmutable) VoiceServerUpdate(data []byte) (*VoiceServerUpdate, error) {
	vsu := &VoiceServerUpdate{}
	if err := json.Unmarshal(data, vsu); err != nil {
		return nil, err
	}
	c.Patch(vsu)

	return vsu, nil
}

func (c *CacheLFUImmutable) GuildMembersChunk(data []byte) (evt *GuildMembersChunk, err error) {
	if evt, err = c.CacheNop.GuildMembersChunk(data); err != nil {
		return nil, err
	}

	users := make([]*User, 0, len(evt.Members))
	for i := range evt.Members {
		users = append(users, evt.Members[i].User)
		evt.Members[i].User = nil
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		_ = c.saveUsers(users)
		wg.Done()
	}()

	sort.Slice(evt.Members, func(i, j int) bool {
		return evt.Members[i].UserID < evt.Members[j].UserID
	})

	c.Guilds.Lock()
	cachedGuild, exists := c.Guilds.Get(evt.GuildID)
	if !exists || cachedGuild == nil {
		cachedGuild = c.Guilds.CreateCacheableItem(&Guild{
			ID:          evt.GuildID,
			Unavailable: true,
		})
		c.Guilds.Set(evt.GuildID, cachedGuild)
	}
	guild := cachedGuild.Val.(*Guild)
	c.Guilds.Unlock()

	mutex := c.Mutex(&c.Guilds, evt.GuildID)
	mutex.Lock()

	// TODO: replace instead of re-allocating?
	//  should be designed for large guilds
	members := make([]*Member, 0, len(guild.Members)+len(evt.Members))
	for i := range guild.Members {
		pos := sort.Search(len(evt.Members), func(si int) bool {
			return evt.Members[si].UserID >= guild.Members[i].UserID
		})
		if pos == len(evt.Members) {
			members = append(members, guild.Members[i])
		}
	}
	members = append(members, evt.Members...)
	guild.Members = members

	mutex.Unlock()

	wg.Wait()
	return c.CacheNop.GuildMembersChunk(data)
}

func (c *CacheLFUImmutable) GuildMemberRemove(data []byte) (*GuildMemberRemove, error) {
	gmr := &GuildMemberRemove{}
	if err := json.Unmarshal(data, gmr); err != nil {
		return nil, err
	}
	c.Patch(gmr)

	if guildI, mu := c.get(&c.Guilds, gmr.GuildID); guildI != nil {
		mu.Lock()
		defer mu.Unlock()

		guild := guildI.(*Guild)
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

func (c *CacheLFUImmutable) GuildMemberUpdate(data []byte) (evt *GuildMemberUpdate, err error) {
	if evt, err = c.CacheNop.GuildMemberUpdate(data); err != nil {
		return nil, err
	}

	uid := evt.User.ID
	gid := evt.GuildID

	userwrap := &userHolder{}

	c.Users.Lock()
	cachedUser, userExists := c.Users.Get(uid)
	c.Users.Unlock()

	if userExists {
		mutex := c.Mutex(&c.Users, uid)
		mutex.Lock()
		userwrap.User = cachedUser.Val.(*User)
		if err := json.Unmarshal(data, userwrap); err == nil {
			c.Patch(userwrap.User)
		}
		mutex.Unlock()
	} else {
		userwrap.User = &User{}
		if err := json.Unmarshal(data, userwrap); err == nil {
			c.Patch(userwrap.User)
			usr := c.Users.CreateCacheableItem(userwrap.User)

			c.Users.Lock()
			if _, exists := c.Users.Get(uid); !exists {
				c.Users.Set(uid, usr)
			}
			c.Users.Unlock()
		}
	}
	userwrap = nil

	item, exists := c.getGuild(gid)
	if exists {
		mutex := c.Mutex(&c.Guilds, gid)
		mutex.Lock()
		defer mutex.Unlock()

		guild := item.Val.(*Guild)

		var member *Member
		for i := range guild.Members { // slow... map instead?
			if guild.Members[i].UserID == uid {
				member = guild.Members[i]
				break
			}
		}
		if member == nil {
			member = &Member{}

			guild.Members = append(guild.Members, member)
			guild.MemberCount++
		}

		if err := json.Unmarshal(data, member); err != nil {
			return nil, err
		}
		c.Patch(member)
		member.User = nil
	}

	return evt, nil
}

func (c *CacheLFUImmutable) GuildMemberAdd(data []byte) (*GuildMemberAdd, error) {
	gmr := &GuildMemberAdd{}
	if err := json.Unmarshal(data, gmr); err != nil {
		return nil, err
	}
	c.Patch(gmr)

	userID := gmr.Member.User.ID
	guildID := gmr.Member.GuildID

	// upsert user
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if user, mu := c.get(&c.Users, userID); user != nil {
			mu.Lock()
			defer mu.Unlock()

			gmr.Member.User = user.(*User)
			_ = json.Unmarshal(data, gmr)
			c.Patch(gmr)
		} else {
			usr := c.Users.CreateCacheableItem(DeepCopy(gmr.Member.User))

			c.Users.Lock()
			if _, exists := c.Users.Get(userID); !exists {
				c.Users.Set(userID, usr)
			}
			c.Users.Unlock()
		}
		wg.Done()
	}()

	// upsert member
	if guildI, mu := c.get(&c.Guilds, guildID); guildI != nil {
		mu.Lock()
		defer mu.Unlock()

		guild := guildI.(*Guild)

		var member *Member
		for i := range guild.Members { // TODO-slow:
			if guild.Members[i].UserID == gmr.Member.User.ID {
				member = guild.Members[i]
				if err := json.Unmarshal(data, member); err != nil {
					return nil, err
				}
				c.Patch(member)
				member.User = nil
				break
			}
		}
		if member == nil {
			member = DeepCopy(gmr.Member).(*Member)
			member.User = nil

			guild.Members = append(guild.Members, member)
			guild.MemberCount++
		}
	}

	wg.Wait()
	return gmr, nil
}

func (c *CacheLFUImmutable) GuildCreate(data []byte) (*GuildCreate, error) {
	var metadata *idHolder
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}
	guildID := metadata.ID

	item, exists := c.getGuild(guildID)
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

		guild = DeepCopy(guild).(*Guild)
	} else if !exists {
		// must create it
		if err := json.Unmarshal(data, &guild); err != nil {
			return nil, err
		}
		c.Patch(guild)

		e := c.Guilds.CreateCacheableItem(guild)
		guild = DeepCopy(guild).(*Guild)

		c.Guilds.Lock()
		if _, exists := c.Guilds.Get(guildID); !exists {
			c.Guilds.Set(guildID, e)
		}
		c.Guilds.Unlock()
	} else {
		// derp - this is really.. not supposed to happen but just in case
		evt, err := c.CacheNop.GuildCreate(data)
		if err != nil {
			return nil, err
		}
		guild = evt.Guild
	}

	// cache channels
	c.Channels.Lock()
	for i := range guild.Channels {
		channel := DeepCopy(guild.Channels[i]).(*Channel)
		_ = c.saveChannel(channel)
	}
	c.Channels.Unlock()

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

		return DeepCopy(guild).(*Guild), nil
	}

	var metadata *idHolder
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}
	guildID := metadata.ID

	item, exists := c.getGuild(guildID)
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
			guild = DeepCopy(guild).(*Guild)
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

func (c *CacheLFUImmutable) GuildRoleCreate(data []byte) (evt *GuildRoleCreate, err error) {
	if evt, err = c.CacheNop.GuildRoleCreate(data); err != nil {
		return nil, err
	}

	item, exists := c.getGuild(evt.GuildID)
	if exists {
		mutex := c.Mutex(&c.Guilds, evt.GuildID)
		mutex.Lock()
		defer mutex.Unlock()

		guild := item.Val.(*Guild)
		_ = guild.AddRole(evt.Role) // TODO: how do i handle this?
	}

	return evt, nil
}

func (c *CacheLFUImmutable) GuildRoleUpdate(data []byte) (evt *GuildRoleUpdate, err error) {
	if evt, err = c.CacheNop.GuildRoleUpdate(data); err != nil {
		return nil, err
	}

	item, exists := c.getGuild(evt.GuildID)
	if exists {
		mutex := c.Mutex(&c.Guilds, evt.GuildID)
		mutex.Lock()
		defer mutex.Unlock()

		guild := item.Val.(*Guild)
		role, err := guild.Role(evt.Role.ID)
		if err != nil {
			// role does not exist
			_ = guild.AddRole(DeepCopy(evt.Role).(*Role)) // TODO: how do i handle this?
		} else {
			tmp := &GuildRoleUpdate{Role: role}
			if err = json.Unmarshal(data, tmp); err != nil {
				return nil, err
			}
			c.Patch(evt)
		}
	}

	return evt, nil
}

func (c *CacheLFUImmutable) GuildRoleDelete(data []byte) (evt *GuildRoleDelete, err error) {
	if evt, err = c.CacheNop.GuildRoleDelete(data); err != nil {
		return nil, err
	}

	item, exists := c.getGuild(evt.GuildID)
	if exists {
		mutex := c.Mutex(&c.Guilds, evt.GuildID)
		mutex.Lock()
		defer mutex.Unlock()

		guild := item.Val.(*Guild)
		guild.DeleteRoleByID(evt.RoleID)
	}

	return evt, nil
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
func (c *CacheLFUImmutable) get(set *crs.LFU, id Snowflake) (interface{}, *sync.Mutex) {
	set.Lock()
	cachedItem, exists := set.Get(id)
	var val interface{}
	if exists {
		val = cachedItem.Val
	}
	set.Unlock()

	if exists {
		return val, c.Mutex(set, id)
	}
	return nil, nil
}

func (c *CacheLFUImmutable) GetChannel(id Snowflake) (*Channel, error) {
	if channel, mu := c.get(&c.Channels, id); channel != nil {
		mu.Lock()
		defer mu.Unlock()

		return DeepCopy(channel.(*Channel)).(*Channel), nil
	}
	return nil, ErrCacheMiss
}

func (c *CacheLFUImmutable) GetGuildEmoji(guildID, emojiID Snowflake) (*Emoji, error) {
	if guild, mu := c.get(&c.Guilds, guildID); guild != nil {
		mu.Lock()
		defer mu.Unlock()

		g := guild.(*Guild)
		emoji, err := g.Emoji(emojiID)
		if err != nil || emoji == nil {
			return nil, ErrCacheMiss
		}

		return DeepCopy(emoji).(*Emoji), nil
	}
	return nil, ErrCacheMiss
}
func (c *CacheLFUImmutable) GetGuildEmojis(id Snowflake) ([]*Emoji, error) {
	if guild, mu := c.get(&c.Guilds, id); guild != nil {
		mu.Lock()
		defer mu.Unlock()

		g := guild.(*Guild)
		emojis := make([]*Emoji, 0, len(g.Emojis))
		for _, emoji := range emojis {
			if emoji == nil { // shouldn't happen, but let's just be safe
				continue
			}
			emojis = append(emojis, DeepCopy(emoji).(*Emoji))
		}

		return emojis, nil
	}
	return nil, ErrCacheMiss
}
func (c *CacheLFUImmutable) GetGuild(id Snowflake) (*Guild, error) {
	if guild, mu := c.get(&c.Guilds, id); guild != nil {
		mu.Lock()
		defer mu.Unlock()

		return DeepCopy(guild.(*Guild)).(*Guild), nil
	}
	return nil, ErrCacheMiss
}
func (c *CacheLFUImmutable) GetGuildChannels(id Snowflake) ([]*Channel, error) {
	if guild, mu := c.get(&c.Guilds, id); guild != nil {
		mu.Lock()
		defer mu.Unlock()

		g := guild.(*Guild)
		channels := make([]*Channel, 0, len(g.Channels))
		for _, channel := range channels {
			if channel == nil { // shouldn't happen, but let's just be safe
				continue
			}
			channels = append(channels, DeepCopy(channel).(*Channel))
		}

		return channels, nil
	}
	return nil, ErrCacheMiss
}

// GetMember fetches member and related user data from cache. User is not guaranteed to be populated.
// Tip: use Member.GetUser(..) instead of Member.User
func (c *CacheLFUImmutable) GetMember(guildID, userID Snowflake) (*Member, error) {
	var user *User
	var member *Member

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		user, _ = c.GetUser(userID)
		wg.Done()
	}()

	if guild, mu := c.get(&c.Guilds, guildID); guild != nil {
		mu.Lock()
		defer mu.Unlock()

		g := guild.(*Guild)
		if member, _ = g.Member(userID); member != nil {
			member = DeepCopy(member).(*Member)
		}
	}
	wg.Wait()

	if member != nil {
		member.User = user
		return member, nil
	} else {
		return nil, ErrCacheMiss
	}
}
func (c *CacheLFUImmutable) GetGuildRoles(guildID Snowflake) ([]*Role, error) {
	if guild, mu := c.get(&c.Guilds, guildID); guild != nil {
		mu.Lock()
		defer mu.Unlock()

		g := guild.(*Guild)
		roles := make([]*Role, 0, len(g.Roles))
		for _, role := range roles {
			if role == nil { // shouldn't happen, but let's just be safe
				continue
			}
			roles = append(roles, DeepCopy(role).(*Role))
		}

		return roles, nil
	}
	return nil, ErrCacheMiss
}
func (c *CacheLFUImmutable) GetCurrentUser() (*User, error) {
	c.CurrentUserMu.Lock()
	defer c.CurrentUserMu.Unlock()
	if c.CurrentUser == nil {
		return nil, ErrCacheMiss
	}

	return DeepCopy(c.CurrentUser).(*User), nil
}
func (c *CacheLFUImmutable) GetUser(id Snowflake) (*User, error) {
	if id == c.currentUserID {
		return c.GetCurrentUser()
	}

	if user, mu := c.get(&c.Users, id); user != nil {
		mu.Lock()
		defer mu.Unlock()

		return DeepCopy(user.(*User)).(*User), nil
	}
	return nil, ErrCacheMiss
}
