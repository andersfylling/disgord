package disgord

import (
	"errors"
	"sort"
	"sync"

	"github.com/andersfylling/disgord/json"
)

var CacheMissErr = errors.New("no matching entry found in cache")
var CacheEntryAlreadyExistsErr = errors.New("cache entry already exists")

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

func NewBasicCache() *BasicCache {
	cache := &BasicCache{
		CurrentUser: &User{},
	}
	cache.Users.Store = make(map[Snowflake]*User)
	cache.Channels.Store = make(map[Snowflake]*Channel)
	cache.Guilds.Store = make(map[Snowflake]*guildCacheContainer)
	cache.VoiceStates.Store = make(map[Snowflake]*voiceStateCacheEntry)

	return cache
}

type voiceStateCache struct {
	sync.Mutex
	Store map[Snowflake]*voiceStateCacheEntry
}

type voiceStateCacheEntry struct {
	sync.Mutex
	GuildID Snowflake
	Store   map[Snowflake]*VoiceState
}

type channelsCache struct {
	sync.Mutex
	Store map[Snowflake]*Channel
}

type guildsCache struct {
	sync.Mutex
	Store map[Snowflake]*guildCacheContainer
}

type usersCache struct {
	sync.Mutex
	Store map[Snowflake]*User
}

type guildCacheContainer struct {
	Guild      *Guild
	ChannelIDs []Snowflake
	Members    map[Snowflake]*Member
}

func retrieveChannels(ids []Snowflake, repo *channelsCache) []*Channel {
	channels := make([]*Channel, 0, len(ids))

	repo.Lock()
	for i := range ids {
		channel, ok := repo.Store[ids[i]]
		if !ok {
			continue
		}

		channels = append(channels, DeepCopy(channel).(*Channel))
	}
	repo.Unlock()

	return channels
}

func constructMemberList(membersMap map[Snowflake]*Member) (members []*Member) {
	for _, member := range membersMap {
		if member == nil {
			continue
		}

		members = append(members, DeepCopy(member).(*Member))
	}
	return members
}

func buildGuildFromCacheContainer(guildCopy *Guild, ChannelIDs []Snowflake, members []*Member, users *usersCache, channels *channelsCache) *Guild {
	guildCopy.Channels = retrieveChannels(ChannelIDs, channels)
	guildCopy.Members = members

	users.Lock()
	for i := range guildCopy.Members {
		member := guildCopy.Members[i]
		user, ok := users.Store[member.UserID]
		if !ok {
			continue
		}

		member.User = DeepCopy(user).(*User)
	}
	users.Unlock()

	return guildCopy
}

// BasicCache cache with CRS support for Users and voice states
// use NewCacheLFUImmutable to instantiate it!
type BasicCache struct {
	CacheNop

	// set via disgord.createClient
	// must never be overwritten
	currentUserID Snowflake // dangerous: no verification that id is set

	CurrentUserMu sync.Mutex
	CurrentUser   *User

	Users       usersCache
	VoiceStates voiceStateCache
	Channels    channelsCache
	Guilds      guildsCache
}

var _ Cache = (*BasicCache)(nil)

func (c *BasicCache) createDMChannel(msg *Message) {
	channelID := msg.ChannelID

	c.Channels.Lock()
	defer c.Channels.Unlock()
	if _, exists := c.Channels.Store[channelID]; !exists {
		channel := &Channel{
			ID:            channelID,
			LastMessageID: msg.ID,
			Type:          ChannelTypeDM,
		}

		c.CurrentUserMu.Lock()
		if c.CurrentUser != nil {
			channel.Recipients = append(channel.Recipients, DeepCopy(c.CurrentUser).(*User))
		}
		c.CurrentUserMu.Unlock()

		if msg.Author != nil {
			channel.Recipients = append(channel.Recipients, DeepCopy(msg.Author).(*User))
		}

		c.Patch(channel)

		c.Channels.Store[channelID] = channel
	}
}

func (c *BasicCache) Ready(data []byte) (*Ready, error) {
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
func (c *BasicCache) MessageCreate(data []byte) (*MessageCreate, error) {
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

func (c *BasicCache) ChannelCreate(data []byte) (*ChannelCreate, error) {
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

func (c *BasicCache) saveChannel(channel *Channel) error {
	if _, exists := c.Channels.Store[channel.ID]; exists {
		return CacheEntryAlreadyExistsErr
	}

	c.Channels.Store[channel.ID] = channel
	return nil
}

func (c *BasicCache) ChannelUpdate(data []byte) (*ChannelUpdate, error) {
	// assumption#1: Create may not take place before an update event
	// assumption#2: The set of fields in both ChannelCreate and ChannelUpdate are the same
	// assumption#3: a channel can not change from one type to another (text => news, text => voice)

	updateChannel := func(channelID Snowflake, channel *Channel) (*Channel, error) {
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

	c.Channels.Lock()
	defer c.Channels.Unlock()

	var channel *Channel
	var err error
	if channelI, ok := c.Channels.Store[channelID]; ok {
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

		if storedChannel, exists := c.Channels.Store[channelID]; !exists {
			c.Channels.Store[channelID] = tmp
		} else if channel, err = updateChannel(channelID, storedChannel); err != nil { // double lock
			return nil, err
		}
	}

	return &ChannelUpdate{Channel: channel}, nil
}

func (c *BasicCache) ChannelDelete(data []byte) (*ChannelDelete, error) {
	cd := &ChannelDelete{}
	if err := json.Unmarshal(data, cd); err != nil {
		return nil, err
	}
	c.Patch(cd)

	c.Channels.Lock()
	defer c.Channels.Unlock()
	delete(c.Channels.Store, cd.Channel.ID)

	return cd, nil
}

func (c *BasicCache) ChannelPinsUpdate(data []byte) (*ChannelPinsUpdate, error) {
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
	if channel, exists := c.Channels.Store[cpu.ChannelID]; exists {
		if cpu.LastPinTimestamp.After(channel.LastPinTimestamp.Time) {
			channel.LastPinTimestamp = cpu.LastPinTimestamp
		}
	}

	return cpu, nil
}

//func (c *BasicCache) VoiceStateUpdate(data []byte) (*VoiceStateUpdate, error) {
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

func (c *BasicCache) UserUpdate(data []byte) (*UserUpdate, error) {
	// assumption#1: this user does not exist in users repo

	c.CurrentUserMu.Lock()
	defer c.CurrentUserMu.Unlock()
	if err := json.Unmarshal(data, c.CurrentUser); err != nil {
		return nil, err
	}
	c.Patch(c.CurrentUser)

	user := DeepCopy(c.CurrentUser).(*User)
	return &UserUpdate{User: user}, nil
}

func (c *BasicCache) saveUsers(users []*User) {
	c.Users.Lock()
	defer c.Users.Unlock()

	for i := range users {
		id := users[i].ID
		if _, ok := c.Users.Store[id]; ok {
			continue
		}

		c.Users.Store[id] = users[i]
	}
}

func (c *BasicCache) VoiceServerUpdate(data []byte) (*VoiceServerUpdate, error) {
	vsu := &VoiceServerUpdate{}
	if err := json.Unmarshal(data, vsu); err != nil {
		return nil, err
	}
	c.Patch(vsu)

	return vsu, nil
}

func (c *BasicCache) GuildMembersChunk(data []byte) (evt *GuildMembersChunk, err error) {
	if evt, err = c.CacheNop.GuildMembersChunk(data); err != nil {
		return nil, err
	}

	users := make([]*User, 0, len(evt.Members))
	for i := range evt.Members {
		user := DeepCopy(evt.Members[i].User).(*User)
		users = append(users, user)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		c.saveUsers(users)
		wg.Done()
	}()

	// stayed here for historical reasons, investigate if this can be removed
	sort.Slice(evt.Members, func(i, j int) bool {
		return evt.Members[i].UserID < evt.Members[j].UserID
	})

	c.Guilds.Lock()
	defer c.Guilds.Unlock()

	if container, ok := c.Guilds.Store[evt.GuildID]; ok {
		// this is fresh data so we just overwrite existing content
		for i := range evt.Members {
			member := DeepCopy(evt.Members[i]).(*Member)
			member.User = nil
			container.Members[member.UserID] = member
		}
	}
	wg.Wait()

	return evt, nil
}

func (c *BasicCache) GuildMemberRemove(data []byte) (*GuildMemberRemove, error) {
	gmr, err := c.CacheNop.GuildMemberRemove(data)
	if err != nil {
		return nil, err
	}

	c.Guilds.Lock()
	defer c.Guilds.Unlock()

	if container, ok := c.Guilds.Store[gmr.GuildID]; ok {
		if _, ok := container.Members[gmr.User.ID]; ok {
			container.Guild.MemberCount--
		}
		delete(container.Members, gmr.User.ID)
	}

	return gmr, nil
}

func (c *BasicCache) GuildMemberUpdate(data []byte) (evt *GuildMemberUpdate, err error) {
	if evt, err = c.CacheNop.GuildMemberUpdate(data); err != nil {
		return nil, err
	}

	c.Guilds.Lock()
	defer c.Guilds.Unlock()

	if container, ok := c.Guilds.Store[evt.GuildID]; ok {
		if member, ok := container.Members[evt.User.ID]; ok {
			if err = json.Unmarshal(data, member); err != nil {
				return nil, err
			}
			c.Patch(evt)
		} else {
			container.Guild.MemberCount++
			container.Members[evt.User.ID] = DeepCopy(evt.Member).(*Member)
		}
		container.Members[evt.User.ID].User = nil
	}

	return evt, nil
}

func (c *BasicCache) GuildMemberAdd(data []byte) (*GuildMemberAdd, error) {
	evt, err := c.CacheNop.GuildMemberAdd(data)
	if err != nil {
		return nil, err
	}

	// save user
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		user := DeepCopy(evt.Member.User).(*User)
		c.saveUsers([]*User{user})
		wg.Done()
	}()

	c.Guilds.Lock()
	defer c.Guilds.Unlock()

	if container, ok := c.Guilds.Store[evt.Member.GuildID]; ok {
		if _, ok := container.Members[evt.Member.UserID]; !ok {
			container.Guild.MemberCount++
		}
		container.Members[evt.Member.UserID] = DeepCopy(evt.Member).(*Member)
	}

	wg.Wait()
	return evt, nil
}

func (c *BasicCache) deconstructGuild(guild *Guild) (*Guild, []Snowflake, map[Snowflake]*Member) {
	channelIDs := make([]Snowflake, 0, len(guild.Channels))
	membersMap := make(map[Snowflake]*Member, len(guild.Members))
	if !guild.Unavailable {
		// cache channels
		c.Channels.Lock()
		for i := range guild.Channels {
			channel := DeepCopy(guild.Channels[i]).(*Channel)
			_ = c.saveChannel(channel)
			channelIDs = append(channelIDs, channel.ID)
		}
		c.Channels.Unlock()
		guild.Channels = nil

		// cache users
		users := make([]*User, 0, len(guild.Members))
		for i := range guild.Members {
			member := guild.Members[i]
			users = append(users, member.User)
			member.User = nil
		}
		c.saveUsers(users)

		// move members
		for i := range guild.Members {
			member := guild.Members[i]
			membersMap[member.UserID] = member
		}
		guild.Members = nil
	}

	return guild, channelIDs, membersMap
}

func (c *BasicCache) GuildCreate(data []byte) (*GuildCreate, error) {
	evt, err := c.CacheNop.GuildCreate(data)
	if err != nil {
		return nil, err
	}

	guild := DeepCopy(evt.Guild).(*Guild)
	_, channelIDs, membersMap := c.deconstructGuild(guild)

	c.Guilds.Lock()
	defer c.Guilds.Unlock()

	c.Guilds.Store[guild.ID] = &guildCacheContainer{
		Guild:      guild,
		ChannelIDs: channelIDs,
		Members:    membersMap,
	} // discard any previous data

	return evt, nil
}

func (c *BasicCache) GuildUpdate(data []byte) (*GuildUpdate, error) {
	evt, err := c.CacheNop.GuildUpdate(data)
	if err != nil {
		return nil, err
	}

	c.Guilds.Lock()
	defer c.Guilds.Unlock()

	container, ok := c.Guilds.Store[evt.Guild.ID]
	if !ok {
		// unlikely - slow case
		guild := DeepCopy(evt.Guild).(*Guild)
		_, channelIDs, membersMap := c.deconstructGuild(guild)

		c.Guilds.Store[guild.ID] = &guildCacheContainer{
			Guild:      guild,
			ChannelIDs: channelIDs,
			Members:    membersMap,
		}
		return evt, nil
	}

	// channels and members should not have been affected by this, so that's a lot of garbage.
	if err = json.Unmarshal(data, container.Guild); err != nil {
		return nil, err
	}
	container.Guild.Members = nil
	container.Guild.Channels = nil
	c.Patch(evt)

	return evt, nil
}

func (c *BasicCache) GuildDelete(data []byte) (*GuildDelete, error) {
	guildEvt := &GuildDelete{}
	if err := json.Unmarshal(data, guildEvt); err != nil {
		return nil, err
	}
	c.Patch(guildEvt)

	c.Guilds.Lock()
	defer c.Guilds.Unlock()
	delete(c.Guilds.Store, guildEvt.UnavailableGuild.ID)

	return guildEvt, nil
}

func (c *BasicCache) GuildRoleCreate(data []byte) (evt *GuildRoleCreate, err error) {
	// since guild create events have to destroy old data to make sure nothing is outdated
	// we do a nop if the guild doesn't exist
	if evt, err = c.CacheNop.GuildRoleCreate(data); err != nil {
		return nil, err
	}
	role := DeepCopy(evt.Role).(*Role)

	c.Guilds.Lock()
	defer c.Guilds.Unlock()

	if container, ok := c.Guilds.Store[evt.GuildID]; ok {
		guild := container.Guild

		var saved bool
		for i := range guild.Roles {
			if role.ID == guild.Roles[i].ID {
				guild.Roles[i] = role
				saved = true
				break
			}
		}

		if !saved {
			guild.Roles = append(guild.Roles, role)
		}
	}

	return evt, nil
}

func (c *BasicCache) GuildRoleUpdate(data []byte) (evt *GuildRoleUpdate, err error) {
	// since guild create events have to destroy old data to make sure nothing is outdated
	// we do a nop if the guild doesn't exist
	if evt, err = c.CacheNop.GuildRoleUpdate(data); err != nil {
		return nil, err
	}

	c.Guilds.Lock()
	defer c.Guilds.Unlock()

	if container, ok := c.Guilds.Store[evt.GuildID]; ok {
		guild := container.Guild
		var storedRole *Role

		for i := range guild.Roles {
			if evt.Role.ID == guild.Roles[i].ID {
				storedRole = guild.Roles[i]
				break
			}
		}

		if storedRole == nil {
			role := DeepCopy(evt.Role).(*Role)
			guild.Roles = append(guild.Roles, role)
		} else if err = json.Unmarshal(data, storedRole); err != nil {
			return nil, err
		} else {
			c.Patch(storedRole)
		}
	}

	return evt, nil
}

func (c *BasicCache) GuildRoleDelete(data []byte) (evt *GuildRoleDelete, err error) {
	if evt, err = c.CacheNop.GuildRoleDelete(data); err != nil {
		return nil, err
	}

	c.Guilds.Lock()
	defer c.Guilds.Unlock()

	if container, ok := c.Guilds.Store[evt.GuildID]; ok {
		container.Guild.DeleteRoleByID(evt.RoleID)
	}

	return evt, nil
}

// REST lookup
// func (c *BasicCache) GetMessage(channelID, messageID Snowflake) (*Message, error) {
// 	return nil, nil
// }
// func (c *BasicCache) GetCurrentUserGuilds(p *GetCurrentUserGuildsParams) ([]*PartialGuild, error) {
// 	return nil, nil
// }
// func (c *BasicCache) GetMessages(channel Snowflake, p *GetMessagesParams) ([]*Message, error) {
// 	return nil, nil
// }
// func (c *BasicCache) GetMembers(guildID Snowflake, p *GetMembersParams) ([]*Member, error) {
// 	return nil, nil
// }

func (c *BasicCache) GetChannel(id Snowflake) (*Channel, error) {
	c.Channels.Lock()
	defer c.Channels.Unlock()

	if channel, ok := c.Channels.Store[id]; ok {
		return DeepCopy(channel).(*Channel), nil
	}
	return nil, CacheMissErr
}

func (c *BasicCache) GetGuildEmoji(guildID, emojiID Snowflake) (*Emoji, error) {
	c.Guilds.Lock()
	defer c.Guilds.Unlock()

	if container, ok := c.Guilds.Store[guildID]; ok {
		if emoji, err := container.Guild.Emoji(emojiID); emoji != nil && err == nil {
			return DeepCopy(emoji).(*Emoji), nil
		}
	}
	return nil, CacheMissErr
}

func (c *BasicCache) GetGuildEmojis(id Snowflake) ([]*Emoji, error) {
	c.Guilds.Lock()
	defer c.Guilds.Unlock()

	if container, ok := c.Guilds.Store[id]; ok {
		emojis := make([]*Emoji, 0, len(container.Guild.Emojis))
		for _, emoji := range emojis {
			if emoji == nil { // shouldn't happen, but let's just be certain
				continue
			}
			emojis = append(emojis, DeepCopy(emoji).(*Emoji))
		}
		return emojis, nil
	}
	return nil, CacheMissErr
}

func (c *BasicCache) GetGuild(id Snowflake) (*Guild, error) {
	var guildCopy *Guild
	var channelIDs []Snowflake
	var members []*Member

	c.Guilds.Lock()
	if container, ok := c.Guilds.Store[id]; ok {
		guildCopy = DeepCopy(container.Guild).(*Guild)
		members = constructMemberList(container.Members)
		channelIDs = make([]Snowflake, len(container.ChannelIDs))
		copy(channelIDs, container.ChannelIDs)
	}
	c.Guilds.Unlock()

	if guildCopy == nil {
		return nil, CacheMissErr
	}

	return buildGuildFromCacheContainer(guildCopy, channelIDs, members, &c.Users, &c.Channels), nil
}

func (c *BasicCache) GetGuildChannels(id Snowflake) ([]*Channel, error) {
	var channelIDs []Snowflake
	var guildFound bool

	c.Guilds.Lock()
	if container, ok := c.Guilds.Store[id]; ok {
		channelIDs = make([]Snowflake, len(container.ChannelIDs))
		copy(channelIDs, container.ChannelIDs)
		guildFound = true
	}
	c.Guilds.Unlock()

	if !guildFound {
		return nil, CacheMissErr
	}
	return retrieveChannels(channelIDs, &c.Channels), nil
}

// GetMember fetches member and related user data from cache. User is not guaranteed to be populated.
// Tip: use Member.GetUser(..) instead of Member.User
func (c *BasicCache) GetMember(guildID, userID Snowflake) (*Member, error) {
	var user *User
	var member *Member

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		user, _ = c.GetUser(userID)
		wg.Done()
	}()

	c.Guilds.Lock()
	defer c.Guilds.Unlock()

	if container, ok := c.Guilds.Store[guildID]; ok {
		if member, _ = container.Members[userID]; member != nil {
			member = DeepCopy(member).(*Member)
		}
	}

	wg.Wait()
	if member != nil {
		member.User = user
		return member, nil
	}

	return nil, CacheMissErr
}
func (c *BasicCache) GetGuildRoles(id Snowflake) ([]*Role, error) {
	c.Guilds.Lock()
	defer c.Guilds.Unlock()

	if container, ok := c.Guilds.Store[id]; ok {
		roles := make([]*Role, 0, len(container.Guild.Channels))
		for _, role := range roles {
			if role == nil { // shouldn't happen, but let's just be certain
				continue
			}
			roles = append(roles, DeepCopy(role).(*Role))
		}
		return roles, nil
	}
	return nil, CacheMissErr
}
func (c *BasicCache) GetCurrentUser() (*User, error) {
	c.CurrentUserMu.Lock()
	defer c.CurrentUserMu.Unlock()
	if c.CurrentUser == nil {
		return nil, CacheMissErr
	}

	return DeepCopy(c.CurrentUser).(*User), nil
}
func (c *BasicCache) GetUser(id Snowflake) (*User, error) {
	if id == c.currentUserID {
		return c.GetCurrentUser()
	}

	c.Users.Lock()
	defer c.Users.Unlock()
	if user, ok := c.Users.Store[id]; ok {
		return DeepCopy(user).(*User), nil
	}
	return nil, CacheMissErr
}
