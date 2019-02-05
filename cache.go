package disgord

import (
	"errors"
	"time"

	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/snowflake/v3"

	"github.com/andersfylling/disgord/cache/interfaces"
	"github.com/andersfylling/disgord/cache/lfu"
	"github.com/andersfylling/disgord/cache/lru"
)

type cacheRegistry uint

// cacheLink keys to redirect to the related cacheLink system
const (
	NoCacheSpecified cacheRegistry = iota
	UserCache
	ChannelCache
	GuildCache
	GuildEmojiCache
	VoiceStateCache
)

// the different cacheLink replacement algorithms
const (
	CacheAlgLRU  = "lru"
	CacheAlgLFU  = "lfu"
	CacheAlgTLRU = "tlru"
)

// Cacher gives basic cacheLink interaction options, and won't require changes when adding more cacheLink systems
type Cacher interface {
	Update(key cacheRegistry, v interface{}) (err error)
	Get(key cacheRegistry, id Snowflake, args ...interface{}) (v interface{}, err error)
	DeleteChannel(channelID snowflake.ID)
	DeleteGuildChannel(guildID snowflake.ID, channelID snowflake.ID)
	UpdateChannelPin(channelID snowflake.ID, lastPinTimestamp Timestamp)
	DeleteGuild(guildID snowflake.ID)
	DeleteGuildRole(guildID snowflake.ID, roleID snowflake.ID)
	UpdateChannelLastMessageID(channelID snowflake.ID, messageID snowflake.ID)
	SetGuildEmojis(guildID Snowflake, emojis []*Emoji)
	Updates(key cacheRegistry, vs []interface{}) error
}

func newErrorCacheItemNotFound(id Snowflake) *ErrorCacheItemNotFound {
	return &ErrorCacheItemNotFound{
		info: "item with id{" + id.String() + "} was not found in cacheLink",
	}
}

// ErrorCacheItemNotFound requested item was not found in cacheLink
type ErrorCacheItemNotFound struct {
	info string
}

// Error ...
func (e *ErrorCacheItemNotFound) Error() string {
	return e.info
}

func newErrorUsingDeactivatedCache(cacheName string) *ErrorUsingDeactivatedCache {
	return &ErrorUsingDeactivatedCache{
		info: "cannot use deactivated cacheLink: " + cacheName,
	}
}

// ErrorUsingDeactivatedCache the cacheLink system you are trying to access has been disabled in the CacheConfig
type ErrorUsingDeactivatedCache struct {
	info string
}

func (e *ErrorUsingDeactivatedCache) Error() string {
	return e.info
}

func constructSpecificCacher(alg string, limit uint, lifetime time.Duration) (cacher interfaces.CacheAlger, err error) {
	switch alg {
	case CacheAlgTLRU:
		//cacher = tlru.NewCacheList(limit, lifetime)
		err = errors.New("TLRU is missing schedulerer for clearing dead/timed out objects and is therefore deactivated")
	case CacheAlgLRU:
		cacher = lru.NewCacheList(limit)
	case CacheAlgLFU:
		cacher = lfu.NewCacheList(limit)
	default:
		err = errors.New("unsupported caching algorithm")
	}

	return
}

func newCache(conf *CacheConfig) (*Cache, error) {
	userCacher, err := createUserCacher(conf)
	if err != nil {
		return nil, err
	}

	voiceStateCacher, err := createVoiceStateCacher(conf)
	if err != nil {
		return nil, err
	}

	channelCacher, err := createChannelCacher(conf)
	if err != nil {
		return nil, err
	}

	guildCacher, err := createGuildCacher(conf)
	if err != nil {
		return nil, err
	}

	return &Cache{
		immutable:   !conf.Mutable,
		conf:        conf,
		users:       userCacher,
		voiceStates: voiceStateCacher,
		channels:    channelCacher,
		guilds:      guildCacher,
	}, nil
}

func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		UserCacheAlgorithm:       CacheAlgLFU,
		VoiceStateCacheAlgorithm: CacheAlgLFU,
		ChannelCacheAlgorithm:    CacheAlgLFU,
		GuildCacheAlgorithm:      CacheAlgLFU,
	}
}

func ensureBasicCacheConfig(conf *CacheConfig) {
	if conf.UserCacheAlgorithm == "" {
		conf.UserCacheAlgorithm = CacheAlgLFU
	}
	if conf.VoiceStateCacheAlgorithm == "" {
		conf.VoiceStateCacheAlgorithm = CacheAlgLFU
	}
	if conf.ChannelCacheAlgorithm == "" {
		conf.ChannelCacheAlgorithm = CacheAlgLFU
	}
	if conf.GuildCacheAlgorithm == "" {
		conf.GuildCacheAlgorithm = CacheAlgLFU
	}
}

// CacheConfig allows for tweaking the cacheLink system on a personal need
type CacheConfig struct {
	// may be false, the new Mutable makes it immutable by default
	// Deprecated
	Immutable bool
	Mutable   bool // Must be immutable to support concurrent access and long-running tasks(!)

	DisableUserCaching bool
	// Deprecated
	UserCacheLimitMiB   uint
	UserCacheMaxEntries uint
	UserCacheLifetime   time.Duration
	UserCacheAlgorithm  string

	DisableVoiceStateCaching  bool
	VoiceStateCacheMaxEntries uint
	VoiceStateCacheLifetime   time.Duration
	VoiceStateCacheAlgorithm  string

	DisableChannelCaching bool
	// Deprecated
	ChannelCacheLimitMiB   uint
	ChannelCacheMaxEntries uint
	ChannelCacheLifetime   time.Duration
	ChannelCacheAlgorithm  string

	DisableGuildCaching bool
	// Deprecated
	GuildCacheLimitMiB   uint
	GuildCacheMaxEntries uint
	GuildCacheLifetime   time.Duration
	GuildCacheAlgorithm  string
}

// Cache is the actual cacheLink. It holds the different systems which can be tweaked using the CacheConfig.
type Cache struct {
	conf        *CacheConfig
	immutable   bool
	users       interfaces.CacheAlger
	voiceStates interfaces.CacheAlger
	channels    interfaces.CacheAlger
	guilds      interfaces.CacheAlger
}

// Updates does the same as Update. But allows for a slice of entries instead.
func (c *Cache) Updates(key cacheRegistry, vs []interface{}) (err error) {
	for _, v := range vs {
		err = c.Update(key, v)
		if err != nil {
			return
		}
	}

	return
}

// Update updates a item in the cacheLink given the key identifier and the new content.
// It also checks if the given structs implements the required interfaces (See below).
func (c *Cache) Update(key cacheRegistry, v interface{}) (err error) {
	if v == nil {
		err = errors.New("object was nil")
		return
	}

	// Does not allow for bulk updates
	//_, implementsDeepCopier := v.(DeepCopier)
	//_, implementsCacheCopier := v.(cacheCopier)
	//if !implementsCacheCopier && !implementsDeepCopier && c.immutable {
	//	err = errors.New("object does not implement DeepCopier & cacheCopier and must do so when config.Mutable is set")
	//	return
	//}

	switch key {
	case UserCache:
		if user, isUser := v.(*User); isUser {
			c.SetUser(user)
		} else {
			err = errors.New("can only save *User structures to user cacheLink")
		}
	case VoiceStateCache:
		if state, isVS := v.(*VoiceState); isVS {
			c.SetVoiceState(state)
		} else {
			err = errors.New("can only save *VoiceState structures to voice state cacheLink")
		}
	case ChannelCache:
		if channel, isChannel := v.(*Channel); isChannel {
			c.SetChannel(channel)
		} else {
			err = errors.New("can only save *Channel structures to channel cacheLink")
		}
	case GuildEmojiCache:
		emojis := v.([]*Emoji)
		if len(emojis) == 0 {
			return
		}
		err = cacheEmoji_SetAll(c, emojis[0].guildID, emojis)
	case GuildCache:
		if guild, ok := v.(*Guild); ok {
			c.SetGuild(guild)
		} else {
			err = errors.New("can only save *Guild structures to guild cacheLink")
		}
	default:
		err = errors.New("caching for given type is not yet implemented")
	}

	return
}

// DirectUpdate is used for socket events to only update provided fields. Will peek into the cacheLink for a matching entry
// if found it updates it, otherwise a not found error is returned. May return a unmarshal error.
//
//  // user update
//  id := extractAttribute([]byte(`"id":"`), 0, jsonData)
//  err := cacheLink.DirectUpdate(UserCache, id, jsonData)
//  if err != nil {
//  	// most likely the user does not exist or it could not be updated
//  	// add the new user. See Cache.Update
//  }
//
// TODO-optimize: for bulk changes
func (c *Cache) DirectUpdate(registry cacheRegistry, id snowflake.ID, changes []byte) error {
	switch registry {
	case UserCache:
		usr, err := c.PeekUser(id)
		if err != nil {
			return err
		}

		err = httd.Unmarshal(changes, usr)
		return err
	}

	return errors.New("could not do a direct update for registry, most likely missing implementation")
}

// Get retrieve a item in the cacheLink, or get an error when not found or if the cacheLink system is disabled
// in your CacheConfig configuration.
func (c *Cache) Get(key cacheRegistry, id Snowflake, args ...interface{}) (v interface{}, err error) {
	switch key {
	case UserCache:
		v, err = c.GetUser(id)
	case VoiceStateCache:
		if len(args) > 0 {
			if params, ok := args[0].(*guildVoiceStateCacheParams); ok {
				v, err = c.GetVoiceState(id, params)
			} else {
				err = errors.New("voice state cacheLink extraction requires an addition argument of type *guildVoiceStateCacheParams")
			}
		} else {
			err = errors.New("voice state cacheLink extraction requires an addition argument for filtering")
		}
	case ChannelCache:
		v, err = c.GetChannel(id)
	default:
		err = errors.New("caching for given type is not yet implemented")
	}
	//
	//// TODO: do deep copying here to speed up the code
	//if copyable, implements := v.(DeepCopier); implements && c.immutable {
	//	v = copyable.DeepCopy()
	//}

	return
}

// --------------------------------------------------------
// Guild

func createGuildCacher(conf *CacheConfig) (cacher interfaces.CacheAlger, err error) {
	if conf.DisableGuildCaching {
		return nil, nil
	}

	var limit uint = conf.GuildCacheMaxEntries
	if limit == 0 && conf.GuildCacheLimitMiB > 0 {
		const guildWeight = 1 // MiB. TODO: what is the actual max size?
		limit = conf.GuildCacheLimitMiB / guildWeight
	}

	cacher, err = constructSpecificCacher(conf.ChannelCacheAlgorithm, limit, conf.ChannelCacheLifetime)
	return
}

type guildCacheItem struct {
	guild    *Guild
	channels []Snowflake
}

func (g *guildCacheItem) process(guild *Guild, immutable bool) {
	if immutable {
		g.guild = guild.DeepCopy().(*Guild)

		for _, member := range g.guild.Members {
			member.userID = member.User.ID
			member.User = nil
		}

		g.channels = make([]snowflake.ID, len(g.guild.Channels))
		for i, channel := range g.guild.Channels {
			g.channels[i] = channel.ID
		}
		g.guild.Channels = nil
	} else {
		g.guild = guild

		for _, member := range g.guild.Members {
			member.userID = member.User.ID
		}
	}

}

func (g *guildCacheItem) build(cache *Cache) (guild *Guild) {
	var err error

	if cache.immutable {
		guild = g.guild.DeepCopy().(*Guild)
		guild.Channels = make([]*Channel, len(g.channels))
		for i := range g.channels {
			guild.Channels[i], err = cache.GetChannel(g.channels[i])
			if err != nil {
				guild.Channels[i] = &Channel{
					ID: g.channels[i],
				}
			}
		}
		for _, member := range guild.Members {
			member.User, err = cache.GetUser(member.userID)
			// member has a GetUser method to handle nil users
		}

		// TODO: voice state
	} else {
		guild = g.guild

		channels := make([]*Channel, len(g.channels))
		for i := range g.channels {
			channels[i], err = cache.GetChannel(g.channels[i])
			if err != nil {
				channels[i] = &Channel{
					ID: g.channels[i],
				}
			}
		}
		guild.Channels = channels
	}

	return
}

func (g *guildCacheItem) update(fresh *Guild, immutable bool) {
	if immutable {
		fresh.copyOverToCache(g.guild)
		// roles
		if len(fresh.Roles) > 0 {
			g.guild.Roles = make([]*Role, len(fresh.Roles))
		}
		for i, role := range fresh.Roles {
			if role == nil {
				continue
			}
			g.guild.Roles[i] = role.DeepCopy().(*Role)
		}
		// emojis
		if len(fresh.Emojis) > 0 {
			g.guild.Emojis = make([]*Emoji, len(fresh.Emojis))
		}
		for i, emoji := range fresh.Emojis {
			if emoji == nil {
				continue
			}
			g.guild.Emojis[i] = emoji.DeepCopy().(*Emoji)
		}
		// voice states
		if len(fresh.VoiceStates) > 0 {
			g.guild.VoiceStates = make([]*VoiceState, len(fresh.VoiceStates))
		}
		for i, state := range fresh.VoiceStates {
			if state == nil {
				continue
			}
			g.guild.VoiceStates[i] = state.DeepCopy().(*VoiceState)
		}
		// members
		if len(fresh.Members) > 0 {
			g.guild.Members = make([]*Member, len(fresh.Members))
		}
		for i, m := range fresh.Members {
			if m == nil {
				continue
			}
			m.userID = m.User.ID
			m.User = nil
			g.guild.Members[i] = m.DeepCopy().(*Member)
		}
		// presences
		if len(fresh.Presences) > 0 {
			g.guild.Presences = make([]*UserPresence, len(fresh.Presences))
		}
		for i, p := range fresh.Presences {
			if p == nil {
				continue
			}
			g.guild.Presences[i] = p.DeepCopy().(*UserPresence)
		}
		// channels
		if len(fresh.Channels) > 0 {
			g.channels = make([]Snowflake, len(fresh.Channels))
		}
		for i, c := range fresh.Channels {
			if c == nil {
				continue
			}
			g.channels[i] = c.ID
		}
	} else {
		if len(fresh.Roles) == 0 && len(g.guild.Roles) > 0 {
			fresh.Roles = g.guild.Roles
		}
		if len(fresh.Emojis) == 0 && len(g.guild.Emojis) > 0 {
			fresh.Emojis = g.guild.Emojis
		}
		if len(fresh.VoiceStates) == 0 && len(g.guild.VoiceStates) > 0 {
			fresh.VoiceStates = g.guild.VoiceStates
		}
		if len(fresh.Members) == 0 && len(g.guild.Members) > 0 {
			fresh.Members = g.guild.Members
		}
		if len(fresh.Channels) == 0 && len(g.guild.Channels) > 0 {
			fresh.Channels = g.guild.Channels
		}
		if len(fresh.Presences) == 0 && len(g.guild.Presences) > 0 {
			fresh.Presences = g.guild.Presences
		}
		g.guild = fresh
	}
}

func (g *guildCacheItem) updateMembers(members []*Member, immutable bool) {
	newMembers := []*Member{}

	g.guild.Lock()
	defer g.guild.Unlock()

	var userID Snowflake
	for i := range members {
		userID = members[i].User.ID
		for j := range g.guild.Members {
			if g.guild.Members[j].userID == userID {
				userID = 0
				*g.guild.Members[j] = *members[i]
				g.guild.Members[j].User = nil
				break
			}
		}

		if !userID.Empty() {
			newMembers = append(newMembers, members[i])
		}
	}

	var member *Member
	for i := range newMembers {
		member = newMembers[i]
		member.userID = member.User.ID
		member.User = nil
		g.guild.Members = append(g.guild.Members, member)
	}
}

func (g *guildCacheItem) deleteChannel(id Snowflake) {
	for i := range g.channels {
		if g.channels[i] != id {
			continue
		}

		g.channels[i] = g.channels[len(g.channels)-1]
		g.channels = g.channels[:len(g.channels)-1]
		g.guild.DeleteChannelByID(id)
		return
	}
}

// SetGuild adds a new guild to cacheLink or updates an existing one
func (c *Cache) SetGuild(guild *Guild) {
	if c.guilds == nil || guild == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()
	if item, exists := c.guilds.Get(guild.ID); exists {
		item.Object().(*guildCacheItem).update(guild, c.immutable)
		c.guilds.RefreshAfterDiscordUpdate(item)
	} else {
		content := &guildCacheItem{}
		content.process(guild, c.immutable)
		c.guilds.Set(guild.ID, c.guilds.CreateCacheableItem(content))
	}
}

// SetGuildEmojis adds a new guild to cacheLink if no guild exist for the emojis or updates an existing guild with the new emojis
func (c *Cache) SetGuildEmojis(guildID Snowflake, emojis []*Emoji) {
	if c.guilds == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()
	if item, exists := c.guilds.Get(guildID); exists {
		guild := item.Object().(*guildCacheItem).guild
		if c.immutable {
			emojisCopy := make([]*Emoji, len(emojis))
			for i := range emojis {
				emojisCopy[i] = emojis[i].DeepCopy().(*Emoji)
			}
			// TODO-racecondition
			guild.Emojis = emojisCopy
		} else {
			// code smell => try to update only affected emoji's
			guild.Emojis = emojis
		}
		c.guilds.RefreshAfterDiscordUpdate(item)
	} else {
		content := &guildCacheItem{}
		content.process(&Guild{
			ID:     guildID,
			Emojis: emojis,
		}, c.immutable)
		c.guilds.Set(guildID, c.guilds.CreateCacheableItem(content))
	}
}

// SetGuildMember calls SetGuildMembers
func (c *Cache) SetGuildMember(guildID Snowflake, member *Member) {
	if c.guilds == nil || member == nil {
		return
	}

	c.SetGuildMembers(guildID, []*Member{member})
}

// SetGuildMembers adds the members to a guild or updates an existing guild
func (c *Cache) SetGuildMembers(guildID Snowflake, members []*Member) {
	if c.guilds == nil || members == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()
	if item, exists := c.guilds.Get(guildID); exists {
		item.Object().(*guildCacheItem).updateMembers(members, c.immutable)
		c.guilds.RefreshAfterDiscordUpdate(item)
	} else {
		content := &guildCacheItem{}
		content.process(&Guild{
			ID:      guildID,
			Members: members,
		}, c.immutable)
		c.guilds.Set(guildID, c.guilds.CreateCacheableItem(content))
	}
}

// SetGuildRoles creates a new guild if none is found and updates the roles for a given guild
func (c *Cache) SetGuildRoles(guildID Snowflake, roles []*Role) {
	if c.guilds == nil || roles == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()
	if item, exists := c.guilds.Get(guildID); exists {
		guild := item.Object().(*guildCacheItem).guild
		var newRoles []*Role
		if c.immutable {
			newRoles = make([]*Role, len(roles))
			for i := range roles {
				newRoles[i] = roles[i].DeepCopy().(*Role)
			}
		} else {
			newRoles = roles
		}
		guild.Roles = newRoles
		c.guilds.RefreshAfterDiscordUpdate(item)
	} else {
		content := &guildCacheItem{}
		content.process(&Guild{
			ID:    guildID,
			Roles: roles,
		}, c.immutable)
		c.guilds.Set(guildID, c.guilds.CreateCacheableItem(content))
	}
}

// GetGuild ...
func (c *Cache) GetGuild(id Snowflake) (guild *Guild, err error) {
	if c.guilds == nil {
		err = newErrorUsingDeactivatedCache("guilds")
		return
	}

	c.guilds.RLock()
	defer c.guilds.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.guilds.Get(id); !exists {
		err = newErrorCacheItemNotFound(id)
		return
	}

	guild = result.Object().(*guildCacheItem).build(c)
	return
}

func (c *Cache) PeekGuild(id snowflake.ID) (guild *Guild, err error) {
	if c.guilds == nil {
		err = newErrorUsingDeactivatedCache("guilds")
		return
	}

	c.guilds.RLock()
	defer c.guilds.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.guilds.Get(id); !exists {
		err = newErrorCacheItemNotFound(id)
		return
	}

	guild = result.Object().(*guildCacheItem).guild
	return
}

// GetGuildRoles ...
func (c *Cache) GetGuildRoles(id Snowflake) (roles []*Role, err error) {
	if c.guilds == nil {
		err = newErrorUsingDeactivatedCache("guilds")
		return
	}

	c.guilds.RLock()
	defer c.guilds.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.guilds.Get(id); !exists {
		err = newErrorCacheItemNotFound(id)
		return
	}

	rolePs := result.Object().(*guildCacheItem).guild.Roles
	if c.immutable {
		roles = make([]*Role, len(rolePs))
		for i := range rolePs {
			roles[i] = rolePs[i].DeepCopy().(*Role)
		}
	} else {
		roles = rolePs
	}

	return
}

// GetGuildRoles ...
func (c *Cache) GetGuildEmojis(id Snowflake) (emojis []*Emoji, err error) {
	if c.guilds == nil {
		err = newErrorUsingDeactivatedCache("guilds")
		return
	}

	c.guilds.RLock()
	defer c.guilds.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.guilds.Get(id); !exists {
		err = newErrorCacheItemNotFound(id)
		return
	}

	emojiPs := result.Object().(*guildCacheItem).guild.Emojis
	if c.immutable {
		emojis = make([]*Emoji, len(emojiPs))
		for i := range emojiPs {
			emojis[i] = emojiPs[i].DeepCopy().(*Emoji)
		}
	} else {
		emojis = emojiPs
	}

	return
}

// GetGuildMember ...
func (c *Cache) GetGuildMember(guildID, userID Snowflake) (member *Member, err error) {
	if c.guilds == nil {
		err = newErrorUsingDeactivatedCache("guilds")
		return
	}

	c.guilds.RLock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.guilds.Get(guildID); !exists {
		err = newErrorCacheItemNotFound(guildID)
		return
	}

	guild := result.Object().(*guildCacheItem).guild
	for i := range guild.Members {
		if guild.Members[i].userID == userID {
			member = guild.Members[i]
			if c.immutable {
				member = member.DeepCopy().(*Member)
			}
			break
		}
	}
	c.guilds.RUnlock()

	if member == nil {
		err = newErrorCacheItemNotFound(userID)
		return
	}

	// add user object
	member.User, err = c.GetUser(userID)
	if err != nil {
		member.User = &User{
			ID: userID,
		}
	}
	return
}

// GetGuildMembersAfter ...
func (c *Cache) GetGuildMembersAfter(guildID, after Snowflake, limit int) (members []*Member, err error) {
	if c.guilds == nil {
		err = newErrorUsingDeactivatedCache("guilds")
		return
	}

	c.guilds.RLock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.guilds.Get(guildID); !exists {
		err = newErrorCacheItemNotFound(guildID)
		return
	}

	guild := result.Object().(*guildCacheItem).guild
	for i := range guild.Members {
		if guild.Members[i].userID > after && len(members) <= limit {
			member := guild.Members[i]
			if c.immutable {
				member = member.DeepCopy().(*Member)
			}
			members = append(members, member)
		} else if len(members) > limit {
			break
		}
	}
	c.guilds.RUnlock()

	for i := range members {
		// add user object
		members[i].User, err = c.GetUser(members[i].userID)
		if err != nil {
			members[i].User = &User{
				ID: members[i].userID,
			}
		}
	}
	return
}

// DeleteGuild ...
func (c *Cache) DeleteGuild(id Snowflake) {
	if c.guilds == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()

	c.guilds.Delete(id)
}

// DeleteGuildChannel removes a channel from a cached guild object without removing the guild
func (c *Cache) DeleteGuildChannel(guildID, channelID Snowflake) {
	if c.guilds == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()
	if item, exists := c.guilds.Get(guildID); exists {
		item.Object().(*guildCacheItem).deleteChannel(channelID)
		c.guilds.RefreshAfterDiscordUpdate(item)
	}
}

// DeleteGuildRole removes a role from a cached guild object without removing the guild
func (c *Cache) DeleteGuildRole(guildID, roleID Snowflake) {
	if c.guilds == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()
	if item, exists := c.guilds.Get(guildID); exists {
		item.Object().(*guildCacheItem).guild.DeleteRoleByID(roleID)
		c.guilds.RefreshAfterDiscordUpdate(item)
	}
}

// --------------------------------------------------------
// Users

func createUserCacher(conf *CacheConfig) (cacher interfaces.CacheAlger, err error) {
	if conf.DisableUserCaching {
		return nil, nil
	}

	var limit uint = conf.UserCacheMaxEntries
	if limit == 0 && conf.UserCacheLimitMiB > 0 {
		const userWeight = 1 // MiB. TODO: what is the actual max size?
		limit = conf.UserCacheLimitMiB / userWeight
	}

	cacher, err = constructSpecificCacher(conf.UserCacheAlgorithm, limit, conf.UserCacheLifetime)
	return
}

// SetUser updates an existing user or adds a new one to the cacheLink
func (c *Cache) SetUser(new *User) {
	if c.users == nil || new == nil {
		return
	}

	c.users.Lock()
	defer c.users.Unlock()
	if item, exists := c.users.Get(new.ID); exists {
		if c.immutable {
			new.copyOverToCache(item.Object())
		} else {
			item.Set(new)
		}
		c.users.RefreshAfterDiscordUpdate(item)
	} else {
		var content interface{}
		if c.immutable {
			content = new.DeepCopy()
		} else {
			content = new
		}
		c.users.Set(new.ID, c.users.CreateCacheableItem(content))
	}
}

// GetUser ...
func (c *Cache) GetUser(id Snowflake) (user *User, err error) {
	if c.users == nil {
		err = newErrorUsingDeactivatedCache("users")
		return
	}

	c.users.RLock()
	defer c.users.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.users.Get(id); !exists {
		err = newErrorCacheItemNotFound(id)
		return
	}

	if c.immutable {
		user = result.Object().(*User).DeepCopy().(*User)
	} else {
		user = result.Object().(*User)
	}

	return
}

func (c *Cache) PeekUser(id snowflake.ID) (*User, error) {
	if c.users == nil {
		return nil, newErrorUsingDeactivatedCache("users")
	}

	c.users.RLock()
	defer c.users.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.users.Get(id); !exists {
		return nil, newErrorCacheItemNotFound(id)
	}

	return result.Object().(*User), nil
}

// --------------------------------------------------------
// Voice States

func createVoiceStateCacher(conf *CacheConfig) (cacher interfaces.CacheAlger, err error) {
	if conf.DisableVoiceStateCaching {
		return nil, nil
	}

	cacher, err = constructSpecificCacher(conf.VoiceStateCacheAlgorithm, conf.VoiceStateCacheMaxEntries, conf.VoiceStateCacheLifetime)
	return
}

type guildVoiceStatesCache struct {
	sessions []*VoiceState
}

func (g *guildVoiceStatesCache) existingSession(state *VoiceState) bool {
	return g.sessionPosition(state) > -1
}

func (g *guildVoiceStatesCache) sessionPosition(state *VoiceState) int {
	for i := range g.sessions {
		// if a channel is moved, the channelID should change(?)
		//if g.sessions[i].ChannelID != state.ChannelID {
		//	continue
		//}

		if g.sessions[i].UserID != state.UserID {
			continue
		}

		if g.sessions[i].SessionID == state.SessionID {
			return i
		}
	}

	return -1
}

func (g *guildVoiceStatesCache) update(state *VoiceState, copyOnWrite bool) {
	pos := g.sessionPosition(state)
	if state.ChannelID.Empty() {
		// someone left
		if pos > -1 {
			g.sessions[pos] = g.sessions[len(g.sessions)-1]
			g.sessions[len(g.sessions)-1] = nil
			g.sessions = g.sessions[:len(g.sessions)-1]
		}
		return
	}

	// someone joined / moved channel. But there exist no information about the session
	// so we register the information
	if pos < 0 {
		var data *VoiceState
		if copyOnWrite {
			data = state.DeepCopy().(*VoiceState) // TODO: handle member
		} else {
			data = state
		}
		g.sessions = append(g.sessions, data)
		return
	}

	// someone moved an existing channel
	if g.sessions[pos].ChannelID != state.ChannelID {
		g.sessions[pos].ChannelID = state.ChannelID
		return
	}

	// TODO: this point should not be reached, unless the above checks are incomplete
}

// SetVoiceState adds a new voice state to cacheLink or updates an existing one
func (c *Cache) SetVoiceState(state *VoiceState) {
	if c.voiceStates == nil || state == nil {
		return
	}

	c.voiceStates.Lock()
	defer c.voiceStates.Unlock()

	id := state.GuildID
	if item, exists := c.voiceStates.Get(id); exists {
		states := item.Object().(*guildVoiceStatesCache)
		states.update(state, c.immutable)
		c.users.RefreshAfterDiscordUpdate(item)
	} else {
		states := &guildVoiceStatesCache{}
		states.update(state, c.immutable)
		c.voiceStates.Set(id, c.voiceStates.CreateCacheableItem(states))
	}
}

type guildVoiceStateCacheParams struct {
	userID    Snowflake
	channelID Snowflake
	sessionID string
}

// GetVoiceState ...
func (c *Cache) GetVoiceState(guildID Snowflake, params *guildVoiceStateCacheParams) (state *VoiceState, err error) {
	if c.voiceStates == nil {
		err = newErrorUsingDeactivatedCache("voice-states")
		return
	}

	c.voiceStates.RLock()
	defer c.voiceStates.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.voiceStates.Get(guildID); !exists {
		err = newErrorCacheItemNotFound(guildID)
		return
	}

	states := result.Object().(*guildVoiceStatesCache)
	filter := &VoiceState{
		ChannelID: params.channelID,
		UserID:    params.userID,
		SessionID: params.sessionID,
	}
	pos := states.sessionPosition(filter)
	if pos < 0 {
		err = errors.New("unable to find state with given params filter")
		return
	}

	match := states.sessions[pos]
	if c.immutable {
		state = match.DeepCopy().(*VoiceState)
	} else {
		state = match
	}

	return
}

// --------------------------------------------------------
// Channels

func createChannelCacher(conf *CacheConfig) (cacher interfaces.CacheAlger, err error) {
	if conf.DisableChannelCaching {
		return nil, nil
	}
	var limit uint = conf.ChannelCacheMaxEntries
	if limit == 0 && conf.ChannelCacheLimitMiB > 0 {
		const channelWeight = 1 // MiB. TODO: what is the actual max size?
		limit = conf.ChannelCacheLimitMiB / channelWeight
	}

	cacher, err = constructSpecificCacher(conf.ChannelCacheAlgorithm, limit, conf.ChannelCacheLifetime)
	return
}

type channelCacheItem struct {
	channel *Channel
}

func (c *channelCacheItem) process(channel *Channel, immutable bool) {
	if immutable {
		c.channel = channel.DeepCopy().(*Channel)
		c.channel.Recipients = []*User{} // clear
	} else {
		c.channel = channel
	}

	c.channel.recipientsIDs = make([]Snowflake, len(channel.Recipients))
	for i := range channel.Recipients {
		c.channel.recipientsIDs = append(c.channel.recipientsIDs, channel.Recipients[i].ID)
	}
}

func (c *channelCacheItem) build(cache *Cache) (channel *Channel) {
	if cache.immutable {
		channel = c.channel.DeepCopy().(*Channel)
	} else {
		channel = c.channel
	}

	if channel.Type != ChannelTypeDM && channel.Type != ChannelTypeGroupDM {
		return
	}

	recipients := make([]*User, len(channel.recipientsIDs))
	for i := range c.channel.recipientsIDs {
		usr, err := cache.GetUser(c.channel.recipientsIDs[i]) // handles immutability on it's own
		if err != nil || usr == nil {
			usr = NewUser()
			usr.ID = c.channel.recipientsIDs[i]
			// TODO: should this be loaded by REST request?...
			// TODO-2: maybe it can be a cacheLink option to load dead members on read?
		}
		recipients[i] = usr
	}

	// TODO-racecondition: when !immutable
	channel.Recipients = recipients
	return
}

func (c *channelCacheItem) update(fresh *Channel, immutable bool) {
	if !immutable {
		c.channel = fresh
		return
	}

	fresh.copyOverToCache(c.channel)
}

// SetChannel adds a new channel to cacheLink or updates an existing one
func (c *Cache) SetChannel(new *Channel) {
	if c.channels == nil || new == nil {
		return
	}

	c.channels.Lock()
	defer c.channels.Unlock()
	if item, exists := c.channels.Get(new.ID); exists {
		item.Object().(*channelCacheItem).update(new, c.immutable)
		c.channels.RefreshAfterDiscordUpdate(item)
	} else {
		content := &channelCacheItem{}
		content.process(new, c.immutable)
		c.channels.Set(new.ID, c.channels.CreateCacheableItem(content))
	}
}

// UpdateChannelPin ...
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
		// channel does not exist in cacheLink, create a partial channel
		partial := &Channel{ID: id, LastPinTimestamp: timestamp}
		content := &channelCacheItem{}
		content.process(partial, c.immutable)
		c.channels.Set(id, c.channels.CreateCacheableItem(content))
	}
}

// UpdateChannelLastMessageID ...
func (c *Cache) UpdateChannelLastMessageID(channelID Snowflake, messageID Snowflake) {
	if c.channels == nil || channelID.Empty() || messageID.Empty() {
		return
	}

	c.channels.Lock()
	defer c.channels.Unlock()
	if item, exists := c.channels.Get(channelID); exists {
		item.Object().(*channelCacheItem).channel.LastMessageID = messageID
		c.channels.RefreshAfterDiscordUpdate(item)
	} else {
		// channel does not exist in cacheLink, create a partial channel
		// this is an indirect channel update..
		//partial := &PartialChannel{ID: channelID, LastMessageID: messageID}
		//content := &channelCacheItem{}
		//content.process(partial, c.immutable)
		//c.channels.Set(channelID, c.channels.CreateCacheableItem(content))
	}
}

// GetChannel ...
func (c *Cache) GetChannel(id Snowflake) (channel *Channel, err error) {
	if c.channels == nil {
		err = newErrorUsingDeactivatedCache("channels")
		return
	}

	c.channels.RLock()
	defer c.channels.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.channels.Get(id); !exists {
		err = newErrorCacheItemNotFound(id)
		return
	}

	channel = result.Object().(*channelCacheItem).build(c)
	return
}

// DeleteChannel ...
func (c *Cache) DeleteChannel(id Snowflake) {
	c.channels.Lock()
	defer c.channels.Unlock()

	c.channels.Delete(id)
}

// --------------------------------------------------------
// Guild

// --------------------------------------------------------
// Guild

var _ Cacher = (*Cache)(nil)
