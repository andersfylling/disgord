package disgord

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/websocket"

	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/httd"
	"github.com/sirupsen/logrus"
)

// Config Configuration for the Disgord client
type Config struct {
	Token      string
	HTTPClient *http.Client

	CancelRequestWhenRateLimited bool

	CacheConfig *CacheConfig

	ShardID      uint
	TotalShards  uint
	WebsocketURL string

	//ImmutableCache bool

	//LoadAllMembers   bool
	//LoadAllChannels  bool
	//LoadAllRoles     bool
	//LoadAllPresences bool

	Debug bool

	// your project name, name of bot, or application
	ProjectName string

	// ActivateEventChannels signifies that the developer will use channels to handle incoming events. May it be
	// in addition to handlers or not. This forces the use of a scheduler to empty the buffered channels when they
	// reach their capacity. Since it requires extra resources, others who have no interest in utilizing channels
	// should not experience any performance penalty (even though it might be unnoticeable).
	ActivateEventChannels bool

	//Logger logger.Logrus
}

// Client is the main disgord client to hold your state and data
type Client struct {
	sync.RWMutex

	config *Config
	token  string

	connected     sync.Mutex
	ws            *websocket.Client
	socketEvtChan <-chan *websocket.Event

	myID Snowflake

	// register listeners for events
	evtDispatch *Dispatch

	// cancelRequestWhenRateLimited by default the client waits until either the HTTPClient.timeout or
	// the rate limit ends before closing a request channel. If activated, in stead, requests will
	// instantly be denied, and the channel closed.
	cancelRequestWhenRateLimited bool

	// discord http api
	req *httd.Client

	httpClient *http.Client

	// cache
	cache *Cache
}

// HeartbeatLatency checks the duration of waiting before receiving a response from Discord when a
// heartbeat packet was sent. Note that heartbeats are usually sent around once a minute and is not a accurate
// way to measure delay between the client and Discord server
func (c *Client) HeartbeatLatency() (duration time.Duration, err error) {
	return c.ws.HeartbeatLatency()
}

// ShardID ...
func (c *Client) ShardID() uint {
	return c.config.ShardID
}

// ShardIDString convert the shard ID to a string
func (c *Client) ShardIDString() string {
	return strconv.Itoa(int(c.ShardID()))
}

// Myself get the current user / connected user
func (c *Client) Myself() (user *User, err error) {
	if c.myID.Empty() {
		user, err = c.GetCurrentUser()
		if err == nil {
			c.myID = user.ID
		}
		return
	}

	var usr interface{}
	usr, err = c.cache.Get(UserCache, c.myID)
	if err == nil {
		user = usr.(*User)
	}
	return
}

func (c *Client) logInfo(msg string) {
	logrus.WithFields(logrus.Fields{
		"lib": LibraryInfo(),
	}).Info(msg)
}

func (c *Client) logErr(msg string) {
	logrus.WithFields(logrus.Fields{
		"lib": LibraryInfo(),
	}).Error(msg)
}

func (c *Client) String() string {
	return LibraryInfo()
}

// RateLimiter return the rate limiter object
func (c *Client) RateLimiter() httd.RateLimiter {
	return c.req.RateLimiter()
}

// Connect establishes a websocket connection to the discord API
func (c *Client) Connect() (err error) {
	// set the user ID upon connection
	// only works for socketing
	c.Once(event.Ready, func(session Session, rdy *Ready) {
		c.myID = rdy.User.ID
	})
	c.On(event.UserUpdate, func(session Session, update *UserUpdate) {
		session.Cache().Update(UserCache, update.User)
	})

	c.logInfo("Connecting to discord Gateway")
	c.evtDispatch.start()
	err = c.ws.Connect()
	if err != nil {
		c.logErr(err.Error())
		return
	}
	c.logInfo("Connected")

	// setup event observer
	go c.eventHandler()

	return nil
}

// Disconnect closes the discord websocket connection
func (c *Client) Disconnect() (err error) {
	fmt.Println() // to keep ^C on it's own line
	c.logInfo("Closing Discord gateway connection")
	c.evtDispatch.stop()
	err = c.ws.Disconnect()
	if err != nil {
		c.logErr(err.Error())
		return
	}
	c.logInfo("Disconnected")

	return nil
}

// DisconnectOnInterrupt wait until a termination signal is detected
func (c *Client) DisconnectOnInterrupt() (err error) {
	// create a channel to listen for termination signals (graceful shutdown)
	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-termSignal

	return c.Disconnect()
}

// Req return the request object. Used in REST requests to handle rate limits,
// wrong http responses, etc.
func (c *Client) Req() httd.Requester {
	return c.req
}

// Cache returns the cache manager for the session
func (c *Client) Cache() Cacher {
	return c.cache
}

// On adds a event handler on the given event.
// On => event => handle the content like this
func (c *Client) On(event string, handlers ...interface{}) {
	c.evtDispatch.On(event, handlers...)
}

// Once same as `On`, however, once the handler is triggered, it is removed. In other words, it is only triggered once.
func (c *Client) Once(event string, handlers ...interface{}) {
	c.evtDispatch.Once(event, handlers...)
}

// Emit sends a socket command directly to Discord.
func (c *Client) Emit(command SocketCommand, data interface{}) error {
	switch command {
	case CommandUpdateStatus, CommandUpdateVoiceState, CommandRequestGuildMembers:
	default:
		return errors.New("command is not supported")
	}
	return c.ws.Emit(command, data)
}

// EventChan get a event channel using the event name
func (c *Client) EventChan(event string) (channel interface{}, err error) {
	return c.evtDispatch.EventChan(event)
}

// EventChannels get access to all the event channels
func (c *Client) EventChannels() (channels EventChannels) {
	return c.evtDispatch
}

// AcceptEvent only events registered using this method is accepted from the Discord socket API. The rest is discarded
// to improve performance.
func (c *Client) AcceptEvent(events ...string) {
	for _, evt := range events {
		c.ws.RegisterEvent(evt)
	}
}

// Generic CRUDS

// DeleteFromDiscord if the given object has implemented the private interface discordDeleter this method can
// be used to delete said object.
func (c *Client) DeleteFromDiscord(obj discordDeleter) (err error) {
	err = obj.deleteFromDiscord(c)
	return
}

// SaveToDiscord saves an object to the Discord servers. This supports creation of new objects or udpating/modifying
// existing objects. It really depends on how the object has implemented the private interface discordSaver.
func (c *Client) SaveToDiscord(obj discordSaver) (err error) {
	err = obj.saveToDiscord(c)
	return
}

// REST
// Audit-log

// GetGuildAuditLogs ...
func (c *Client) GetGuildAuditLogs(guildID Snowflake, params *GuildAuditLogsParams) (log *AuditLog, err error) {
	log, err = GuildAuditLogs(c.req, guildID, params)
	return
}

// Channel

// GetChannel ...
func (c *Client) GetChannel(id Snowflake) (ret *Channel, err error) {
	if ret, err = c.cache.GetChannel(id); err != nil {
		ret, err = GetChannel(c.req, id)
		if err != nil {
			return
		}
		_ = c.cache.Update(ChannelCache, ret)
	}
	return
}

// ModifyChannel ...
func (c *Client) ModifyChannel(id Snowflake, changes *ModifyChannelParams) (ret *Channel, err error) {
	ret, err = ModifyChannel(c.req, id, changes) // should trigger a socket event, no need to update cache
	return
}

// DeleteChannel ...
func (c *Client) DeleteChannel(id Snowflake) (channel *Channel, err error) {
	channel, err = DeleteChannel(c.req, id) // should trigger a socket event, no need to update cache
	return
}

// EditChannelPermissions ...
func (c *Client) EditChannelPermissions(chanID, overwriteID Snowflake, params *EditChannelPermissionsParams) (err error) {
	err = EditChannelPermissions(c.req, chanID, overwriteID, params)
	return
}

// GetChannelInvites ...
func (c *Client) GetChannelInvites(id Snowflake) (ret []*Invite, err error) {
	ret, err = GetChannelInvites(c.req, id)
	return
}

// CreateChannelInvites ...
func (c *Client) CreateChannelInvites(id Snowflake, params *CreateChannelInvitesParams) (ret *Invite, err error) {
	ret, err = CreateChannelInvites(c.req, id, params)
	return
}

// DeleteChannelPermission .
func (c *Client) DeleteChannelPermission(channelID, overwriteID Snowflake) (err error) {
	err = DeleteChannelPermission(c.req, channelID, overwriteID)
	return
}

// TriggerTypingIndicator .
func (c *Client) TriggerTypingIndicator(channelID Snowflake) (err error) {
	err = TriggerTypingIndicator(c.req, channelID)
	return
}

// GetPinnedMessages .
func (c *Client) GetPinnedMessages(channelID Snowflake) (ret []*Message, err error) {
	ret, err = GetPinnedMessages(c.req, channelID)
	return
}

// AddPinnedChannelMessage .
func (c *Client) AddPinnedChannelMessage(channelID, msgID Snowflake) (err error) {
	err = AddPinnedChannelMessage(c.req, channelID, msgID)
	return
}

// DeletePinnedChannelMessage .
func (c *Client) DeletePinnedChannelMessage(channelID, msgID Snowflake) (err error) {
	err = DeletePinnedChannelMessage(c.req, channelID, msgID)
	return
}

// GroupDMAddRecipient .
func (c *Client) GroupDMAddRecipient(channelID, userID Snowflake, params *GroupDMAddRecipientParams) (err error) {
	err = GroupDMAddRecipient(c.req, channelID, userID, params)
	return
}

// GroupDMRemoveRecipient .
func (c *Client) GroupDMRemoveRecipient(channelID, userID Snowflake) (err error) {
	err = GroupDMRemoveRecipient(c.req, channelID, userID)
	return
}

// GetChannelMessages .
func (c *Client) GetChannelMessages(channelID Snowflake, params URLParameters) (ret []*Message, err error) {
	ret, err = GetChannelMessages(c.req, channelID, params)
	return
}

// GetChannelMessage .
func (c *Client) GetChannelMessage(channelID, messageID Snowflake) (ret *Message, err error) {
	ret, err = GetChannelMessage(c.req, channelID, messageID)
	return
}

// CreateChannelMessage .
func (c *Client) CreateChannelMessage(channelID Snowflake, params *CreateChannelMessageParams) (ret *Message, err error) {
	ret, err = CreateChannelMessage(c.req, channelID, params)
	return
}

// EditMessage .
func (c *Client) EditMessage(chanID, msgID Snowflake, params *EditMessageParams) (ret *Message, err error) {
	ret, err = EditMessage(c.req, chanID, msgID, params)
	return
}

// DeleteMessage .
func (c *Client) DeleteMessage(channelID, msgID Snowflake) (err error) {
	err = DeleteMessage(c.req, channelID, msgID)
	return
}

// BulkDeleteMessages .
func (c *Client) BulkDeleteMessages(chanID Snowflake, params *BulkDeleteMessagesParams) (err error) {
	err = BulkDeleteMessages(c.req, chanID, params)
	return
}

// CreateReaction .
func (c *Client) CreateReaction(channelID, messageID Snowflake, emoji interface{}) (ret *Reaction, err error) {
	err = CreateReaction(c.req, channelID, messageID, emoji)
	return
}

// DeleteOwnReaction .
func (c *Client) DeleteOwnReaction(channelID, messageID Snowflake, emoji interface{}) (err error) {
	err = DeleteOwnReaction(c.req, channelID, messageID, emoji)
	return
}

// DeleteUserReaction .
func (c *Client) DeleteUserReaction(channelID, messageID, userID Snowflake, emoji interface{}) (err error) {
	err = DeleteUserReaction(c.req, channelID, messageID, userID, emoji)
	return
}

// GetReaction .
func (c *Client) GetReaction(channelID, messageID Snowflake, emoji interface{}, params URLParameters) (ret []*User, err error) {
	ret, err = GetReaction(c.req, channelID, messageID, emoji, params)
	return
}

// DeleteAllReactions .
func (c *Client) DeleteAllReactions(channelID, messageID Snowflake) (err error) {
	err = DeleteAllReactions(c.req, channelID, messageID)
	return
}

// Emoji

// GetGuildEmojis .
func (c *Client) GetGuildEmojis(id Snowflake) (ret []*Emoji, err error) {
	var guild *Guild
	guild, err = c.cache.GetGuild(id)
	if err != nil {
		ret, err = ListGuildEmojis(c.req, id)
		if err == nil {
			c.cache.SetGuildEmojis(id, ret)
		}
		return
	}

	ret = guild.Emojis
	return
}

// GetGuildEmoji .
func (c *Client) GetGuildEmoji(guildID, emojiID Snowflake) (ret *Emoji, err error) {
	var guild *Guild
	guild, err = c.cache.GetGuild(guildID)
	if err != nil {
		ret, err = GetGuildEmoji(c.req, guildID, emojiID)
		// TODO: cache
		return
	}
	ret, err = guild.Emoji(emojiID)
	if err != nil {
		ret, err = GetGuildEmoji(c.req, guildID, emojiID)
		// TODO: cache
		return
	}
	return
}

// CreateGuildEmoji .
func (c *Client) CreateGuildEmoji(guildID Snowflake, params *CreateGuildEmojiParams) (ret *Emoji, err error) {
	ret, err = CreateGuildEmoji(c.req, guildID, params)
	return
}

// ModifyGuildEmoji .
func (c *Client) ModifyGuildEmoji(guildID, emojiID Snowflake, params *ModifyGuildEmojiParams) (ret *Emoji, err error) {
	ret, err = ModifyGuildEmoji(c.req, guildID, emojiID, params)
	return
}

// DeleteGuildEmoji .
func (c *Client) DeleteGuildEmoji(guildID, emojiID Snowflake) (err error) {
	err = DeleteGuildEmoji(c.req, guildID, emojiID)
	return
}

// Guild

// CreateGuild .
func (c *Client) CreateGuild(params *CreateGuildParams) (ret *Guild, err error) {
	ret, err = CreateGuild(c.req, params)
	return
}

// GetGuild .
func (c *Client) GetGuild(id Snowflake) (ret *Guild, err error) {
	ret, err = c.cache.GetGuild(id)
	if err != nil {
		ret, err = GetGuild(c.req, id)
		if err != nil {
			return
		}
		c.cache.SetGuild(ret)
	}
	return
}

// ModifyGuild .
func (c *Client) ModifyGuild(id Snowflake, params *ModifyGuildParams) (ret *Guild, err error) {
	ret, err = ModifyGuild(c.req, id, params)
	return
}

// DeleteGuild .
func (c *Client) DeleteGuild(id Snowflake) (err error) {
	err = DeleteGuild(c.req, id)
	return
}

// GetGuildChannels .
func (c *Client) GetGuildChannels(id Snowflake) (ret []*Channel, err error) {
	var guild *Guild
	guild, err = c.cache.GetGuild(id)
	if err != nil {
		ret, err = GetGuildChannels(c.req, id)
		if err != nil {
			return
		}
		c.cache.SetGuild(&Guild{
			ID:       id,
			Channels: ret,
		})
	} else {
		ret = guild.Channels
	}
	return
}

// CreateGuildChannel .
func (c *Client) CreateGuildChannel(id Snowflake, params *CreateGuildChannelParams) (ret *Channel, err error) {
	ret, err = CreateGuildChannel(c.req, id, params)
	return
}

// GetGuildMember .
func (c *Client) GetGuildMember(guildID, userID Snowflake) (ret *Member, err error) {
	ret, err = c.cache.GetGuildMember(guildID, userID)
	if err != nil {
		ret, err = GetGuildMember(c.req, guildID, userID)
		if err != nil {
			return
		}
		c.cache.SetGuildMember(guildID, ret)
		c.cache.Update(UserCache, ret.User)
	}
	return
}

// GetGuildMembers .
func (c *Client) GetGuildMembers(guildID, after Snowflake, limit int) (ret []*Member, err error) {
	ret, err = c.cache.GetGuildMembersAfter(guildID, after, limit)
	if err != nil {
		ret, err = GetGuildMembers(c.req, guildID, after, limit)
		if err != nil {
			return
		}
		c.cache.SetGuildMembers(guildID, ret)
		//c.cache.Update(UserCache, ret.User)
		// TODO: update users
	}
	return
}

// AddGuildMember .
func (c *Client) AddGuildMember(guildID, userID Snowflake, params *AddGuildMemberParams) (ret *Member, err error) {
	ret, err = AddGuildMember(c.req, guildID, userID, params)
	return
}

// ModifyGuildMember .
func (c *Client) ModifyGuildMember(guildID, userID Snowflake, params *ModifyGuildMemberParams) (err error) {
	err = ModifyGuildMember(c.req, guildID, userID, params)
	return
}

// ModifyCurrentUserNick .
func (c *Client) ModifyCurrentUserNick(id Snowflake, params *ModifyCurrentUserNickParams) (nick string, err error) {
	nick, err = ModifyCurrentUserNick(c.req, id, params)
	return
}

// AddGuildMemberRole .
func (c *Client) AddGuildMemberRole(guildID, userID, roleID Snowflake) (err error) {
	err = AddGuildMemberRole(c.req, guildID, userID, roleID)
	return
}

// RemoveGuildMemberRole .
func (c *Client) RemoveGuildMemberRole(guildID, userID, roleID Snowflake) (err error) {
	err = RemoveGuildMemberRole(c.req, guildID, userID, roleID)
	return
}

// RemoveGuildMember .
func (c *Client) RemoveGuildMember(guildID, userID Snowflake) (err error) {
	err = RemoveGuildMember(c.req, guildID, userID)
	return
}

// GetGuildBans .
func (c *Client) GetGuildBans(id Snowflake) (ret []*Ban, err error) {
	ret, err = GetGuildBans(c.req, id)
	return
}

// GetGuildBan .
func (c *Client) GetGuildBan(guildID, userID Snowflake) (ret *Ban, err error) {
	ret, err = GetGuildBan(c.req, guildID, userID)
	return
}

// CreateGuildBan .
func (c *Client) CreateGuildBan(guildID, userID Snowflake, params *CreateGuildBanParams) (err error) {
	err = CreateGuildBan(c.req, guildID, userID, params)
	return
}

// RemoveGuildBan .
func (c *Client) RemoveGuildBan(guildID, userID Snowflake) (err error) {
	err = RemoveGuildBan(c.req, guildID, userID)
	return
}

// GetGuildRoles .
func (c *Client) GetGuildRoles(guildID Snowflake) (ret []*Role, err error) {
	ret, err = c.cache.GetGuildRoles(guildID)
	if err != nil {
		ret, err = GetGuildRoles(c.req, guildID)
		if err != nil {
			return
		}
		c.cache.SetGuildRoles(guildID, ret)
	}

	return
}

// CreateGuildRole .
func (c *Client) CreateGuildRole(id Snowflake, params *CreateGuildRoleParams) (ret *Role, err error) {
	ret, err = CreateGuildRole(c.req, id, params)
	return
}

// ModifyGuildRolePositions .
func (c *Client) ModifyGuildRolePositions(guildID Snowflake, params *ModifyGuildRolePositionsParams) (ret []*Role, err error) {
	ret, err = ModifyGuildRolePositions(c.req, guildID, params)
	return
}

// ModifyGuildRole .
func (c *Client) ModifyGuildRole(guildID, roleID Snowflake, params *ModifyGuildRoleParams) (ret *Role, err error) {
	ret, err = ModifyGuildRole(c.req, guildID, roleID, params)
	return
}

// DeleteGuildRole .
func (c *Client) DeleteGuildRole(guildID, roleID Snowflake) (err error) {
	err = DeleteGuildRole(c.req, guildID, roleID)
	return
}

// GetGuildPruneCount .
func (c *Client) GetGuildPruneCount(id Snowflake, params *GuildPruneParams) (ret *GuildPruneCount, err error) {
	ret, err = GetGuildPruneCount(c.req, id, params)
	return
}

// BeginGuildPrune .
func (c *Client) BeginGuildPrune(id Snowflake, params *GuildPruneParams) (ret *GuildPruneCount, err error) {
	ret, err = BeginGuildPrune(c.req, id, params)
	return
}

// GetGuildVoiceRegions .
func (c *Client) GetGuildVoiceRegions(id Snowflake) (ret []*VoiceRegion, err error) {
	ret, err = GetGuildVoiceRegions(c.req, id)
	return
}

// GetGuildInvites .
func (c *Client) GetGuildInvites(id Snowflake) (ret []*Invite, err error) {
	ret, err = GetGuildInvites(c.req, id)
	return
}

// GetGuildIntegrations .
func (c *Client) GetGuildIntegrations(id Snowflake) (ret []*Integration, err error) {
	ret, err = GetGuildIntegrations(c.req, id)
	return
}

// CreateGuildIntegration .
func (c *Client) CreateGuildIntegration(guildID Snowflake, params *CreateGuildIntegrationParams) (err error) {
	err = CreateGuildIntegration(c.req, guildID, params)
	return
}

// ModifyGuildIntegration .
func (c *Client) ModifyGuildIntegration(guildID, integrationID Snowflake, params *ModifyGuildIntegrationParams) (err error) {
	err = ModifyGuildIntegration(c.req, guildID, integrationID, params)
	return
}

// DeleteGuildIntegration .
func (c *Client) DeleteGuildIntegration(guildID, integrationID Snowflake) (err error) {
	err = DeleteGuildIntegration(c.req, guildID, integrationID)
	return
}

// SyncGuildIntegration .
func (c *Client) SyncGuildIntegration(guildID, integrationID Snowflake) (err error) {
	err = SyncGuildIntegration(c.req, guildID, integrationID)
	return
}

// GetGuildEmbed .
func (c *Client) GetGuildEmbed(guildID Snowflake) (ret *GuildEmbed, err error) {
	ret, err = GetGuildEmbed(c.req, guildID)
	return
}

// ModifyGuildEmbed .
func (c *Client) ModifyGuildEmbed(guildID Snowflake, params *GuildEmbed) (ret *GuildEmbed, err error) {
	ret, err = ModifyGuildEmbed(c.req, guildID, params)
	return
}

// GetGuildVanityURL .
func (c *Client) GetGuildVanityURL(guildID Snowflake) (ret *PartialInvite, err error) {
	ret, err = GetGuildVanityURL(c.req, guildID)
	return
}

// Invite

// GetInvite .
func (c *Client) GetInvite(inviteCode string, withCounts bool) (invite *Invite, err error) {
	invite, err = GetInvite(c.req, inviteCode, withCounts)
	return
}

// DeleteInvite .
func (c *Client) DeleteInvite(inviteCode string) (invite *Invite, err error) {
	invite, err = DeleteInvite(c.req, inviteCode)
	return
}

// User

// GetCurrentUser .
func (c *Client) GetCurrentUser() (ret *User, err error) {
	ret, err = GetCurrentUser(c.req)
	return
}

// GetUser .
func (c *Client) GetUser(id Snowflake) (ret *User, err error) {
	ret, err = c.cache.GetUser(id)
	if err != nil {
		ret, err = GetUser(c.req, id)
		if err != nil {
			return
		}
		c.cache.Update(UserCache, ret)
	}
	return
}

// ModifyCurrentUser .
func (c *Client) ModifyCurrentUser(params *ModifyCurrentUserParams) (ret *User, err error) {
	ret, err = ModifyCurrentUser(c.req, params)
	return
}

// GetCurrentUserGuilds .
func (c *Client) GetCurrentUserGuilds(params *GetCurrentUserGuildsParams) (ret []*Guild, err error) {
	ret, err = GetCurrentUserGuilds(c.req, params)
	return
}

// LeaveGuild .
func (c *Client) LeaveGuild(id Snowflake) (err error) {
	err = LeaveGuild(c.req, id)
	return
}

// GetUserDMs .
func (c *Client) GetUserDMs() (ret []*Channel, err error) {
	ret, err = GetUserDMs(c.req)
	return
}

// CreateDM .
func (c *Client) CreateDM(recipientID Snowflake) (ret *Channel, err error) {
	ret, err = CreateDM(c.req, recipientID)
	return
}

// CreateGroupDM .
func (c *Client) CreateGroupDM(params *CreateGroupDMParams) (ret *Channel, err error) {
	ret, err = CreateGroupDM(c.req, params)
	return
}

// GetUserConnections .
func (c *Client) GetUserConnections() (ret []*UserConnection, err error) {
	ret, err = GetUserConnections(c.req)
	return
}

// Voice

// GetVoiceRegions .
func (c *Client) GetVoiceRegions() (ret []*VoiceRegion, err error) {
	ret, err = ListVoiceRegions(c.req)
	return
}

// Webhook

// CreateWebhook .
func (c *Client) CreateWebhook(channelID Snowflake, params *CreateWebhookParams) (ret *Webhook, err error) {
	ret, err = CreateWebhook(c.req, channelID, params)
	return
}

// GetChannelWebhooks .
func (c *Client) GetChannelWebhooks(channelID Snowflake) (ret []*Webhook, err error) {
	ret, err = GetChannelWebhooks(c.req, channelID)
	return
}

// GetGuildWebhooks .
func (c *Client) GetGuildWebhooks(guildID Snowflake) (ret []*Webhook, err error) {
	ret, err = GetGuildWebhooks(c.req, guildID)
	return
}

// GetWebhook .
func (c *Client) GetWebhook(id Snowflake) (ret *Webhook, err error) {
	ret, err = GetWebhook(c.req, id)
	return
}

// GetWebhookWithToken .
func (c *Client) GetWebhookWithToken(id Snowflake, token string) (ret *Webhook, err error) {
	ret, err = GetWebhookWithToken(c.req, id, token)
	return
}

// ModifyWebhook .
func (c *Client) ModifyWebhook(id Snowflake, params *ModifyWebhookParams) (ret *Webhook, err error) {
	if id.Empty() {
		err = errors.New("given webhook ID was not set, there is nothing to modify")
		return
	}
	if params == nil {
		err = errors.New("given param object was nil, there is nothing to modify")
		return
	}
	if params.Empty() {
		err = errors.New("given param object was empty, there is nothing to modify")
		return
	}

	// verify avatar string prefix
	if params.avatarIsSet && params.avatar != "" && !validAvatarPrefix(params.avatar) {
		err = errors.New("given avatar string is invalid. Must specify data encoding. Eg. data:image/jpeg;base64,")
		return
	}

	// TODO: check if user has permission to modify webhook
	ret, err = ModifyWebhook(c.req, id, params)
	return
}

// ModifyWebhookWithToken .
func (c *Client) ModifyWebhookWithToken(newWebhook *Webhook) (ret *Webhook, err error) {
	ret, err = ModifyWebhookWithToken(c.req, newWebhook)
	return
}

// DeleteWebhook .
func (c *Client) DeleteWebhook(webhookID Snowflake) (err error) {
	err = DeleteWebhook(c.req, webhookID)
	return
}

// DeleteWebhookWithToken .
func (c *Client) DeleteWebhookWithToken(id Snowflake, token string) (err error) {
	err = DeleteWebhookWithToken(c.req, id, token)
	return
}

// ExecuteWebhook .
func (c *Client) ExecuteWebhook(params *ExecuteWebhookParams, wait bool, URLSuffix string) (err error) {
	err = ExecuteWebhook(c.req, params, wait, URLSuffix)
	return
}

// ExecuteSlackWebhook .
func (c *Client) ExecuteSlackWebhook(params *ExecuteWebhookParams, wait bool) (err error) {
	err = ExecuteSlackWebhook(c.req, params, wait)
	return
}

// ExecuteGitHubWebhook .
func (c *Client) ExecuteGitHubWebhook(params *ExecuteWebhookParams, wait bool) (err error) {
	err = ExecuteGitHubWebhook(c.req, params, wait)
	return
}

// Custom methods are usually reused by the resource package for readability
// -----

// SendMsg .
func (c *Client) SendMsg(channelID Snowflake, message *Message) (msg *Message, err error) {
	if constant.LockedMethods {
		message.RLock()
	}
	params := &CreateChannelMessageParams{
		Content: message.Content,
		Tts:     message.Tts,
		// File: ...
		// Embed: ...
	}
	if !message.Nonce.Empty() {
		params.Nonce = message.Nonce
	}
	if len(message.Embeds) > 0 {
		params.Embed = message.Embeds[0]
	}

	if constant.LockedMethods {
		message.RUnlock()
	}

	return c.CreateChannelMessage(channelID, params)
}

// SendMsgString .
func (c *Client) SendMsgString(channelID Snowflake, content string) (msg *Message, err error) {
	params := &CreateChannelMessageParams{
		Content: content,
	}

	msg, err = c.CreateChannelMessage(channelID, params)
	return
}

// UpdateMessage .
func (c *Client) UpdateMessage(message *Message) (msg *Message, err error) {
	if constant.LockedMethods {
		message.RLock()
		defer message.RUnlock()
	}

	params := &EditMessageParams{
		Content: message.Content,
	}
	if len(message.Embeds) > 0 {
		params.Embed = message.Embeds[0]
	}

	msg, err = c.EditMessage(message.ChannelID, message.ID, params)
	return
}

// UpdateChannel Not implemented yet
func (c *Client) UpdateChannel(channel *Channel) (err error) {
	// there are several different REST calls that needs to be made in order
	// to update the channel. But how exactly do we know what has changed?
	return errors.New("not implemented")
}

func waitForEvent(eventEmitter <-chan *websocket.Event) (event *websocket.Event, err error) {
	var alive bool
	event, alive = <-eventEmitter
	if !alive {
		err = errors.New("event emitter (channel) is dead")
	}

	return
}

// eventHandler Takes a incoming event from the websocket package, parses it, and sends
// trigger requests to the event dispatcher and state cacher.
func (c *Client) eventHandler() {
	for {
		var err error
		var evt *websocket.Event

		evt, err = waitForEvent(c.socketEvtChan)
		if err != nil {
			return
		}

		var box eventBox

		switch evt.Name {
		case EventReady:
			box = &Ready{}
		case EventResumed:
			box = &Resumed{}
		case EventChannelCreate:
			box = &ChannelCreate{}
		case EventChannelUpdate:
			box = &ChannelUpdate{}
		case EventChannelDelete:
			box = &ChannelDelete{}
		case EventChannelPinsUpdate:
			box = &ChannelPinsUpdate{}
		case EventGuildCreate:
			box = &GuildCreate{}
		case EventGuildUpdate:
			box = &GuildUpdate{}
		case EventGuildDelete:
			box = &GuildDelete{}
		case EventGuildBanAdd:
			box = &GuildBanAdd{}
		case EventGuildBanRemove:
			box = &GuildBanRemove{}
		case EventGuildEmojisUpdate:
			box = &GuildEmojisUpdate{}
		case EventGuildIntegrationsUpdate:
			box = &GuildIntegrationsUpdate{}
		case EventGuildMemberAdd:
			box = &GuildMemberAdd{}
		case EventGuildMemberRemove:
			box = &GuildMemberRemove{}
		case EventGuildMemberUpdate:
			box = &GuildMemberUpdate{}
		case EventGuildMembersChunk:
			box = &GuildMembersChunk{}
		case EventGuildRoleCreate:
			box = &GuildRoleCreate{}
		case EventGuildRoleUpdate:
			box = &GuildRoleUpdate{}
		case EventGuildRoleDelete:
			box = &GuildRoleDelete{}
		case EventMessageCreate:
			box = &MessageCreate{}
		case EventMessageUpdate:
			box = &MessageUpdate{}
		case EventMessageDelete:
			box = &MessageDelete{}
		case EventMessageDeleteBulk:
			box = &MessageDeleteBulk{}
		case EventMessageReactionAdd:
			box = &MessageReactionAdd{}
		case EventMessageReactionRemove:
			box = &MessageReactionRemove{}
		case EventMessageReactionRemoveAll:
			box = &MessageReactionRemoveAll{}
		case EventPresenceUpdate:
			box = &PresenceUpdate{}
		case EventTypingStart:
			box = &TypingStart{}
		case EventUserUpdate:
			box = &UserUpdate{}
		case EventVoiceStateUpdate:
			box = &VoiceStateUpdate{}
		case EventVoiceServerUpdate:
			box = &VoiceServerUpdate{}
		case EventWebhooksUpdate:
			box = &WebhooksUpdate{}
		default:
			fmt.Printf("------\nTODO\nImplement event handler for `%s`, data: \n%+v\n------\n\n", evt.Name, string(evt.Data))
			continue // move on to next event
		}

		// populate box
		ctx := context.Background()
		box.registerContext(ctx)

		// first unmarshal to get identifiers
		//tmp := *box

		// unmarshal into cache
		//err := c.cacheEvent2(evtName, box)

		err = unmarshal(evt.Data, box)
		if err != nil {
			logrus.Error(err)
			continue // ignore event
			// TODO: if an event is ignored, should it not at least send a signal for listeners with no parameters?
		}

		// cache
		c.cacheEvent(evt.Name, box)

		// trigger listeners
		prepareBox(evt.Name, box)
		c.evtDispatch.triggerChan(ctx, evt.Name, c, box)
		go c.evtDispatch.triggerCallbacks(ctx, evt.Name, c, box)
	}
}

func (c *Client) cacheEvent(event string, v interface{}) (err error) {
	// updates holds key and object to be cached
	updates := map[int]([]interface{}){}

	switch event {
	case EventReady:
		ready := v.(*Ready)
		updates[UserCache] = append(updates[UserCache], ready.User)

		for _, guild := range ready.Guilds {
			updates[GuildCache] = append(updates[GuildCache], guild)
		}
	case EventVoiceStateUpdate:
		update := v.(*VoiceStateUpdate)
		updates[VoiceStateCache] = append(updates[VoiceStateCache], update.VoiceState)
	case EventChannelCreate, EventChannelUpdate:
		var channel *Channel
		if event == EventChannelCreate {
			channel = (v.(*ChannelCreate)).Channel
		} else if event == EventChannelUpdate {
			channel = (v.(*ChannelUpdate)).Channel
		}
		if len(channel.Recipients) > 0 {
			for i := range channel.Recipients {
				updates[UserCache] = append(updates[UserCache], channel.Recipients[i])
			}
		}

		updates[ChannelCache] = append(updates[ChannelCache], channel)
	case EventChannelDelete:
		channel := (v.(*ChannelDelete)).Channel
		c.cache.DeleteChannel(channel.ID)
		c.cache.DeleteGuildChannel(channel.GuildID, channel.ID)
	case EventChannelPinsUpdate:
		evt := v.(*ChannelPinsUpdate)
		c.cache.UpdateChannelPin(evt.ChannelID, evt.LastPinTimestamp)
	case EventGuildCreate, EventGuildUpdate:
		var guild *Guild
		if event == EventGuildCreate {
			guild = (v.(*GuildCreate)).Guild
		} else if event == EventGuildUpdate {
			guild = (v.(*GuildUpdate)).Guild
		}
		updates[GuildCache] = append(updates[GuildCache], guild)

		// update all users
		if len(guild.Members) > 0 {
			updates[UserCache] = make([]interface{}, len(guild.Members))
			for i := range guild.Members {
				updates[UserCache][i] = guild.Members[i].User
			}
		}
	case EventGuildDelete:
		uguild := (v.(*GuildDelete)).UnavailableGuild
		c.cache.DeleteGuild(uguild.ID)
	case EventGuildRoleDelete:
		evt := v.(*GuildRoleDelete)
		c.cache.DeleteGuildRole(evt.GuildID, evt.RoleID)
	case EventGuildEmojisUpdate:
		evt := v.(*GuildEmojisUpdate)
		c.cache.SetGuildEmojis(evt.GuildID, evt.Emojis)
	case EventUserUpdate:
		usr := v.(*UserUpdate).User
		updates[UserCache] = append(updates[UserCache], usr)
	case EventMessageCreate:
		// TODO: performance issues?
		msg := (v.(*MessageCreate)).Message
		c.cache.UpdateChannelLastMessageID(msg.ChannelID, msg.ID)
	default:
		err = errors.New("unsupported event for caching")
		//case EventResumed:
		//case EventGuildBanAdd:
		//case EventGuildBanRemove:
		//case EventGuildIntegrationsUpdate:
		//case EventGuildMemberAdd:
		//case EventGuildMemberRemove:
		//case EventGuildMemberUpdate:
		//case EventGuildMembersChunk:
		//case EventGuildRoleCreate:
		//case EventGuildRoleUpdate:
		//case EventMessageUpdate:
		//case EventMessageDelete:
		//case EventMessageDeleteBulk:
		//case EventMessageReactionAdd:
		//case EventMessageReactionRemove:
		//case EventMessageReactionRemoveAll:
		//case EventPresenceUpdate:
		//case EventTypingStart:
		//case EventVoiceServerUpdate:
		//case EventWebhooksUpdate:
	}

	for key, structs := range updates {
		err = c.cache.Updates(key, structs)
	}
	return
}
