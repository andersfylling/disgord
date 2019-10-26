package disgord

import (
	"sync"
	"time"

	"github.com/andersfylling/djp"

	"github.com/andersfylling/disgord/internal/constant"
	"github.com/andersfylling/disgord/internal/crs"
	jp "github.com/buger/jsonparser"
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

	GuildMembersCache
	GuildRolesCache // warning: deletes previous roles
	GuildRoleCache  // updates or adds a new role
)

// gatewayCacher allows cache repositories to handle event content.
//
// handleGatewayEvent should only be called for situations where a object
// will be created, or updated. Not deleted. For deleting an object use
// the Delete method of a BasicCacheRepo implemented cache repository.
type gatewayCacher interface {
	handleGatewayEvent(evt string, data []byte, flags Flag) (updated interface{}, err error)
}
type restCacher interface {
	handleRESTResponse(obj interface{}) error
}

type evtCacheFunc func(data []byte, flags Flag) (updated interface{}, err error)

type BasicCacheRepo interface {
	Size() uint
	Cap() uint
	ListIDs() []Snowflake
	// Get returns nil when no match was found
	Get(id Snowflake) interface{}
	Del(id Snowflake)
}

// CacheConfig allows for tweaking the cacheLink system on a personal need
type CacheConfig struct {
	Mutable bool // Must be immutable to support concurrent access and long-running tasks(!)

	DisableUserCaching  bool
	UserCacheMaxEntries uint
	UserCacheLifetime   time.Duration

	DisableVoiceStateCaching  bool
	VoiceStateCacheMaxEntries uint
	VoiceStateCacheLifetime   time.Duration

	DisableChannelCaching  bool
	ChannelCacheMaxEntries uint
	ChannelCacheLifetime   time.Duration

	DisableGuildCaching  bool
	GuildCacheMaxEntries uint
	GuildCacheLifetime   time.Duration

	// Deprecated
	UserCacheAlgorithm string
	// Deprecated
	VoiceStateCacheAlgorithm string
	// Deprecated
	ChannelCacheAlgorithm string
	// Deprecated
	GuildCacheAlgorithm string
}

func newCache(websocketShardCount uint, conf *CacheConfig) (c *cache, err error) {
	if websocketShardCount == 0 {
		websocketShardCount = 1
	}

	c = &cache{
		conf: conf,
	}

	// pools
	poolUsers := &bottomlessPool{
		New: func() Reseter {
			return &User{}
		},
	}
	poolChannels := &bottomlessPool{
		New: func() Reseter {
			return &Channel{}
		},
	}

	userShardCount := 12 * websocketShardCount
	if userShardCount > 256 {
		userShardCount = 256
	}
	for i := uint(0); i < userShardCount; i++ {
		c.userRepos = append(c.userRepos, &usersCache{
			conf,
			crs.New(conf.UserCacheMaxEntries),
			poolUsers,
		})
	}

	c.channelRepos = append(c.channelRepos, &channelsCache{
		conf,
		crs.New(conf.ChannelCacheMaxEntries),
		poolChannels,
		c.userRepos[0],
	})
	c.guildRepos = append(c.guildRepos, &guildsCache{
		conf:  conf,
		items: crs.New(conf.ChannelCacheMaxEntries),
		pool: &bottomlessPool{
			New: func() Reseter {
				return &Guild{}
			},
		},
		users:    c.userRepos[0],
		channels: c.channelRepos[0],
	})
	c.presenceRepos = append(c.presenceRepos, &presencesCache{
		conf:  conf,
		items: crs.New(conf.ChannelCacheMaxEntries),
		pool: &bottomlessPool{
			New: func() Reseter {
				return &UserPresence{}
			},
		},
		users: c.userRepos[0],
	})

	return c, nil
}

type cache struct {
	conf          *CacheConfig
	bot           *User
	userRepos     []*usersCache
	channelRepos  []*channelsCache
	guildRepos    []*guildsCache
	presenceRepos []*presencesCache
}

func (c *cache) resultRef(x DeepCopier) interface{} {
	if c.conf.Mutable {
		return x
	}

	return x.DeepCopy()
}

func (c *cache) work(job func()) {
	// makes it easier to do a goroutine pool if needed
	go job()
}

//////////////////////////////////////////////////////
//
// sharding
//
//////////////////////////////////////////////////////

func (c *cache) shardID(id Snowflake, nrOfRepos int) int {
	return int(uint64(id) % uint64(nrOfRepos))
}

func (c *cache) users(id Snowflake) *usersCache {
	return c.userRepos[c.shardID(id, len(c.userRepos))]
}

func (c *cache) channels(id Snowflake) *channelsCache {
	return c.channelRepos[c.shardID(id, len(c.channelRepos))]
}

func (c *cache) guilds(id Snowflake) *guildsCache {
	return c.guildRepos[c.shardID(id, len(c.guildRepos))]
}

func (c *cache) presences(id Snowflake) *presencesCache {
	return c.presenceRepos[c.shardID(id, len(c.presenceRepos))]
}

//////////////////////////////////////////////////////
//
// websocket events
//
//////////////////////////////////////////////////////

func (c *cache) onReady(data []byte, flags Flag) (updated interface{}, err error) {
	r := &Ready{User: c.bot}
	err = Unmarshal(data, &r)

	if c.conf.DisableGuildCaching {
		return r, err
	}

	for i := range r.Guilds {
		id := r.Guilds[i].ID
		c.guilds(id).Prepare(id)
	}

	return r, nil
}

func (c *cache) onResumed(data []byte, flags Flag) (updated interface{}, err error) {
	return &Resumed{}, nil
}

func (c *cache) onChannelCreate(data []byte, flags Flag) (updated interface{}, err error) {
	if c.conf.DisableChannelCaching {
		var cc *ChannelCreate
		err = Unmarshal(data, &cc)
		return cc, err
	}

	id, err := djp.GetSnowflake(data, "id")
	if err != nil {
		return nil, err
	}

	var channel *Channel
	var channelErr error
	var recipients []*User
	var usersErr error

	wg := sync.WaitGroup{}
	wg.Add(2)
	c.work(func() {
		defer wg.Done()

		var channelI interface{}
		channelI, channelErr = c.channels(id).onChannelCreate(data, flags)
		if channelErr != nil {
			return
		}
		channel = channelI.(*Channel)
	})
	c.work(func() {
		defer wg.Done()

		var recipientsI interface{}
		recipientsI, usersErr = c.users(id).onChannelCreate(data, flags)
		if usersErr != nil {
			return
		}
		recipients = recipientsI.([]*User)
	})
	wg.Wait()

	if channelErr != nil {
		return nil, channelErr
	}
	// no need to worry about this. At this stage the json should have been valid
	// so the error is more likely to be related to missing recipients due to the
	// channel type not being group or DM.
	// TODO: check if error is only "missing users"
	//if usersErr != nil {
	//	return nil, usersErr
	//}

	channel.Recipients = recipients
	return channel, nil
}

// onChannelUpdate
//
// Discord: Sent when a channel is updated. The inner payload is a channel object. This is not
// sent when the field last_message_id is altered. To keep track of the last_message_id changes,
// you must listen for Message Create events.
func (c *cache) onChannelUpdate(data []byte, flags Flag) (updated interface{}, err error) {
	id, err := djp.GetSnowflake(data, "id")
	if err != nil {
		return nil, err
	}

	lock := func(m *Member, cb func(*Member)) {
		if constant.LockedMethods {
			m.Lock()
		}
		cb(m)
		if constant.LockedMethods {
			m.Unlock()
		}
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()
	var newMembers []*Member
	for i := range members {
		var updated bool
		for j := range guild.Members {
			if guild.Members[j].userID != 0 && guild.Members[j].userID == members[i].userID {
				var tmp *User
				lock(members[i], func(m *Member) {
					tmp = members[i].User
					members[i].User = nil
				})
				_ = members[i].CopyOverTo(guild.Members[j])
				lock(members[i], func(_ *Member) {
					members[i].User = tmp
				})
				updated = true
				break
			}
		}

		if !updated {
			var tmp *User
			lock(members[i], func(m *Member) {
				tmp = members[i].User
				members[i].User = nil
			})
			member := members[i].DeepCopy().(*Member)
			lock(members[i], func(_ *Member) {
				members[i].User = tmp
			})

			newMembers = append(newMembers, member)
		}
	}

	guild.Members = append(guild.Members, newMembers...)

	if guild.MemberCount < uint(len(guild.Members)) {
		guild.MemberCount = uint(len(guild.Members))
	}
}

// GetGuildMember ...
func (c *Cache) GetGuildMember(guildID, userID Snowflake) (member *Member, err error) {
	if c.guilds == nil {
		err = newErrorUsingDeactivatedCache("guilds")
		return
	}

	c.guilds.RLock()

	result, exists := c.guilds.Get(guildID)
	if !exists {
		err = newErrorCacheItemNotFound(guildID)
		return
	}

	guild := result.Val.(*guildCacheItem).guild
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

	// add user object if it exists
	member.User, _ = c.GetUser(userID)
	return
}

// GetGuildMembersAfter ...
func (c *Cache) GetGuildMembersAfter(guildID, after Snowflake, limit int) (members []*Member, err error) {
	if c.guilds == nil {
		err = newErrorUsingDeactivatedCache("guilds")
		return
	}

	var recipients []*User

	wg := sync.WaitGroup{}
	wg.Add(2)
	c.work(func() {
		defer wg.Done()

		recipientsI, err := c.users(id).onChannelUpdate(data, flags)
		if err != nil {
			return
		}
		recipients = recipientsI.([]*User)
	})
	c.work(func() {
		defer wg.Done()

		guildID, err := djp.GetSnowflake(data, "guild_id")
		if err != nil || guildID.IsZero() {
			return
		}

		_, _ = c.guilds(guildID).onChannelUpdate(data, flags)
	})
	wg.Wait()

	updated, err = c.channels(id).onChannelUpdate(data, flags)
	if err != nil {
		return nil, err
	}
	// for usersErr see onChannelCreate

	updated.(*Channel).Recipients = recipients
	return updated, nil
}

func (c *cache) onChannelDelete(data []byte, flags Flag) (updated interface{}, err error) {
	id, err := djp.GetSnowflake(data, "id")
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	c.work(func() {
		defer wg.Done()

		guildID, err := djp.GetSnowflake(data, "guild_id")
		if err != nil || guildID.IsZero() {
			return
		}

		_, _ = c.guilds(guildID).onChannelDelete(data, flags)
	})

	updated, err = c.channels(id).onChannelDelete(data, flags)
	if err != nil {
		return nil, err
	}
	wg.Wait()

	return updated, nil
}

func (c *cache) onChannelPinsUpdate(data []byte, flags Flag) (updated interface{}, err error) {
	if c.conf.DisableChannelCaching {
		var cpu *ChannelPinsUpdate
		err = Unmarshal(data, &cpu)
		return cpu, err
	}
}

func (c *cache) onGuildCreate(data []byte, flags Flag) (updated interface{}, err error) {

}

func (c *cache) onGuildUpdate(data []byte, flags Flag) (updated interface{}, err error)             {}
func (c *cache) onGuildDelete(data []byte, flags Flag) (updated interface{}, err error)             {}
func (c *cache) onGuildBanAdd(data []byte, flags Flag) (updated interface{}, err error)             {}
func (c *cache) onGuildBanRemove(data []byte, flags Flag) (updated interface{}, err error)          {}
func (c *cache) onGuildEmojisUpdate(data []byte, flags Flag) (updated interface{}, err error)       {}
func (c *cache) onGuildIntegrationsUpdate(data []byte, flags Flag) (updated interface{}, err error) {}
func (c *cache) onGuildMemberAdd(data []byte, flags Flag) (updated interface{}, err error)          {}
func (c *cache) onGuildMemberRemove(data []byte, flags Flag) (updated interface{}, err error)       {}
func (c *cache) onGuildMemberUpdate(data []byte, flags Flag) (updated interface{}, err error)       {}
func (c *cache) onGuildMembersChunk(data []byte, flags Flag) (updated interface{}, err error)       {}
func (c *cache) onGuildRoleCreate(data []byte, flags Flag) (updated interface{}, err error)         {}
func (c *cache) onGuildRoleUpdate(data []byte, flags Flag) (updated interface{}, err error)         {}
func (c *cache) onGuildRoleDelete(data []byte, flags Flag) (updated interface{}, err error)         {}

func (c *cache) onMessageCreate(data []byte, flags Flag) (updated interface{}, err error) {
	channelID, err := djp.GetSnowflake(data, "channel_id")
	if err == nil && !channelID.IsZero() {
		_, _ = c.channels(channelID).onMessageCreate(data, flags)
	}

	var msg *MessageCreate
	err = Unmarshal(data, &msg)
	return msg, err
}

func (c *cache) onMessageUpdate(data []byte, flags Flag) (updated interface{}, err error) {
	var msg *MessageUpdate
	err = Unmarshal(data, &msg)
	return msg, err
}

func (c *cache) onMessageDelete(data []byte, flags Flag) (updated interface{}, err error) {
	var msg *MessageDelete
	err = Unmarshal(data, &msg)
	return msg, err
}

func (c *cache) onMessageDeleteBulk(data []byte, flags Flag) (updated interface{}, err error) {
	var msg *MessageDeleteBulk
	err = Unmarshal(data, &msg)
	return msg, err
}

func (c *cache) onMessageReactionAdd(data []byte, flags Flag) (updated interface{}, err error)       {}
func (c *cache) onMessageReactionRemove(data []byte, flags Flag) (updated interface{}, err error)    {}
func (c *cache) onMessageReactionRemoveAll(data []byte, flags Flag) (updated interface{}, err error) {}
func (c *cache) onPresenceUpdate(data []byte, flags Flag) (updated interface{}, err error)           {}

func (c *cache) onTypingStart(data []byte, flags Flag) (updated interface{}, err error) {
	var ts *TypingStart
	err = Unmarshal(data, &ts)
	return ts, err
}

func (c *cache) onUserUpdate(data []byte, flags Flag) (updated interface{}, err error)        {}
func (c *cache) onVoiceStateUpdate(data []byte, flags Flag) (updated interface{}, err error)  {}
func (c *cache) onVoiceServerUpdate(data []byte, flags Flag) (updated interface{}, err error) {}

func (c *cache) onWebhooksUpdate(data []byte, flags Flag) (updated interface{}, err error) {
	var wu *WebhooksUpdate
	err = Unmarshal(data, &wu)
	return wu, err
}

//////////////////////////////////////////////////////
//
// JSON HELPERS
// All helpers must start with a lowercase "json".
//
//////////////////////////////////////////////////////

// jsonNumberOfKeys returns the number of json keys at depth 1.
func jsonNumberOfKeys(data []byte) (nrOfKeys uint) {
	jp.EachKey(data, func(i int, bytes []byte, valueType jp.ValueType, e error) {
		nrOfKeys++
	})
	return
}

func jsonArrayLen(data []byte, keys ...string) (len int) {
	_, _ = jp.ArrayEach(data, func(b []byte, _ jp.ValueType, _ int, _ error) {
		len++
	}, keys...)

	return
}
