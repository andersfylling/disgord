package disgord

import (
	"net/http"
	"time"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/websocket"
)

// NewSession create a client and return the Session interface
func NewSession(conf *Config) (Session, error) {
	if conf.HTTPClient == nil {
		// http client configuration
		conf.HTTPClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	if conf.ProjectName == "" {
		conf.ProjectName = LibraryInfo()
	}
	dws, err := websocket.NewClient(&websocket.Config{
		// user settings
		Token:      conf.Token,
		HTTPClient: conf.HTTPClient,
		Debug:      conf.Debug,

		// identity
		Browser:             LibraryInfo(),
		Device:              conf.ProjectName,
		GuildLargeThreshold: 250, // TODO: config

		// lib specific
		DAPIVersion:   constant.DiscordVersion,
		DAPIEncoding:  JSONEncoding,
		ChannelBuffer: 1,
	})
	if err != nil {
		return nil, err
	}

	// request client
	reqConf := &httd.Config{
		APIVersion:                   constant.DiscordVersion,
		BotToken:                     conf.Token,
		UserAgentSourceURL:           GitHubURL,
		UserAgentVersion:             constant.Version,
		HTTPClient:                   conf.HTTPClient,
		CancelRequestWhenRateLimited: conf.CancelRequestWhenRateLimited,
	}
	reqClient := httd.NewClient(reqConf)

	// event dispatcher
	evtDispatcher := NewDispatch(dws)

	// caching
	cacheConfig := &CacheConfig{
		Immutable: conf.ImmutableCache,

		UserCacheAlgorithm: CacheAlg_TLRU,
		UserCacheLifetime:  time.Duration(9) * time.Hour,
		UserCacheLimitMiB:  500,

		VoiceStateCacheAlgorithm: CacheAlg_LRU,
	}
	cacher, err := NewCache(cacheConfig)
	if err != nil {
		return nil, err
	}

	// register for events for activated caches
	if cacheConfig.UserCaching {
		dws.RegisterEvent(event.Ready)
		dws.RegisterEvent(event.UserUpdate)
	}
	if cacheConfig.VoiceStateCaching {
		dws.RegisterEvent(event.VoiceStateUpdate)
	}
	if cacheConfig.ChannelCaching {
		dws.RegisterEvent(event.ChannelCreate)
		dws.RegisterEvent(event.ChannelUpdate)
		dws.RegisterEvent(event.ChannelPinsUpdate)
		dws.RegisterEvent(event.ChannelDelete)
	}

	// create a disgord client/instance/session
	c := &Client{
		config:        conf,
		httpClient:    conf.HTTPClient,
		ws:            dws,
		socketEvtChan: dws.DiscordWSEventChan(),
		token:         conf.Token,
		evtDispatch:   evtDispatcher,
		cache:         cacher,
		req:           reqClient,
	}

	return c, nil
}

// NewSessionMustCompile same as NewClientMustCompile, but with the Session
// interface
func NewSessionMustCompile(conf *Config) (session Session) {
	var err error
	session, err = NewSession(conf)
	if err != nil {
		panic(err)
	}

	return
}

type EventChannels interface {
	All() <-chan interface{} // any event
	Ready() <-chan *Ready
	Resumed() <-chan *Resumed
	ChannelCreate() <-chan *ChannelCreate
	ChannelUpdate() <-chan *ChannelUpdate
	ChannelDelete() <-chan *ChannelDelete
	ChannelPinsUpdate() <-chan *ChannelPinsUpdate
	GuildCreate() <-chan *GuildCreate
	GuildUpdate() <-chan *GuildUpdate
	GuildDelete() <-chan *GuildDelete
	GuildBanAdd() <-chan *GuildBanAdd
	GuildBanRemove() <-chan *GuildBanRemove
	GuildEmojisUpdate() <-chan *GuildEmojisUpdate
	GuildIntegrationsUpdate() <-chan *GuildIntegrationsUpdate
	GuildMemberAdd() <-chan *GuildMemberAdd
	GuildMemberRemove() <-chan *GuildMemberRemove
	GuildMemberUpdate() <-chan *GuildMemberUpdate
	GuildMembersChunk() <-chan *GuildMembersChunk
	GuildRoleUpdate() <-chan *GuildRoleUpdate
	GuildRoleCreate() <-chan *GuildRoleCreate
	GuildRoleDelete() <-chan *GuildRoleDelete
	MessageCreate() <-chan *MessageCreate
	MessageUpdate() <-chan *MessageUpdate
	MessageDelete() <-chan *MessageDelete
	MessageDeleteBulk() <-chan *MessageDeleteBulk
	MessageReactionAdd() <-chan *MessageReactionAdd
	MessageReactionRemove() <-chan *MessageReactionRemove
	MessageReactionRemoveAll() <-chan *MessageReactionRemoveAll
	PresenceUpdate() <-chan *PresenceUpdate
	PresencesReplace() <-chan *PresencesReplace
	TypingStart() <-chan *TypingStart
	UserUpdate() <-chan *UserUpdate
	VoiceStateUpdate() <-chan *VoiceStateUpdate
	VoiceServerUpdate() <-chan *VoiceServerUpdate
	WebhooksUpdate() <-chan *WebhooksUpdate
}

type SocketHandler interface {
	Connect() error
	Disconnect() error
	DisconnectOnInterrupt() error

	// event handlers
	On(event string, handler ...interface{})
	Emit(command SocketCommand, dataPointer interface{})
	//Use(middleware ...interface{}) // TODO: is this useful?

	// event channels
	EventChan(event string) (channel interface{}, err error)
	EventChannels() EventChannels

	// event register (which events to accept)
	// events which are not registered are discarded at socket level
	// to increase performance
	AcceptEvent(events ...string)
}

type AuditLogsRESTer interface {
	GetGuildAuditLogs(guildID Snowflake, params *GuildAuditLogsParams) (log *AuditLog, err error)
}

type ChannelRESTer interface {
	GetChannel(id Snowflake) (ret *Channel, err error)
	ModifyChannel(changes *ModifyChannelParams) (ret *Channel, err error)
	DeleteChannel(id Snowflake) (err error)
	EditChannelPermissions(chanID, overwriteID Snowflake, params *EditChannelPermissionsParams) (err error)
	GetChannelInvites(id Snowflake) (ret []*Invite, err error)
	CreateChannelInvites(id Snowflake, params *CreateChannelInvitesParams) (ret *Invite, err error)
	DeleteChannelPermission(channelID, overwriteID Snowflake) (err error)
	TriggerTypingIndicator(channelID Snowflake) (err error)
	GetPinnedMessages(channelID Snowflake) (ret []*Message, err error)
	AddPinnedChannelMessage(channelID, msgID Snowflake) (err error)
	DeletePinnedChannelMessage(channelID, msgID Snowflake) (err error)
	GroupDMAddRecipient(channelID, userID Snowflake, params *GroupDMAddRecipientParams) (err error)
	GroupDMRemoveRecipient(channelID, userID Snowflake) (err error)
	GetChannelMessages(channelID Snowflake, params URLParameters) (ret []*Message, err error)
	GetChannelMessage(channelID, messageID Snowflake) (ret *Message, err error)
	CreateChannelMessage(channelID Snowflake, params *CreateChannelMessageParams) (ret *Message, err error)
	EditMessage(chanID, msgID Snowflake, params *EditMessageParams) (ret *Message, err error)
	DeleteMessage(channelID, msgID Snowflake) (err error)
	BulkDeleteMessages(chanID Snowflake, params *BulkDeleteMessagesParams) (err error)
	CreateReaction(channelID, messageID Snowflake, emoji interface{}) (ret *Reaction, err error)
	DeleteOwnReaction(channelID, messageID Snowflake, emoji interface{}) (err error)
	DeleteUserReaction(channelID, messageID, userID Snowflake, emoji interface{}) (err error)
	GetReaction(channelID, messageID Snowflake, emoji interface{}, params URLParameters) (ret []*User, err error)
	DeleteAllReactions(channelID, messageID Snowflake) (err error)
}

type EmojiRESTer interface {
	GetGuildEmojis(id Snowflake) (ret []*Emoji, err error)
	GetGuildEmoji(guildID, emojiID Snowflake) (ret *Emoji, err error)
	CreateGuildEmoji(guildID Snowflake, params *CreateGuildEmojiParams) (ret *Emoji, err error)
	ModifyGuildEmoji(guildID, emojiID Snowflake, params *ModifyGuildEmojiParams) (ret *Emoji, err error)
	DeleteGuildEmoji(guildID, emojiID Snowflake) (err error)
}

type GuildRESTer interface {
	CreateGuild(params *CreateGuildParams) (ret *Guild, err error)
	GetGuild(id Snowflake) (ret *Guild, err error)
	ModifyGuild(id Snowflake, params *ModifyGuildParams) (ret *Guild, err error)
	DeleteGuild(id Snowflake) (err error)
	GetGuildChannels(id Snowflake) (ret []*Channel, err error)
	CreateGuildChannel(id Snowflake, params *CreateGuildChannelParams) (ret *Channel, err error)
	GetGuildMember(guildID, userID Snowflake) (ret *Member, err error)
	GetGuildMembers(guildID, after Snowflake, limit int) (ret []*Member, err error)
	AddGuildMember(guildID, userID Snowflake, params *AddGuildMemberParams) (ret *Member, err error)
	ModifyGuildMember(guildID, userID Snowflake, params *ModifyGuildMemberParams) (err error)
	ModifyCurrentUserNick(id Snowflake, params *ModifyCurrentUserNickParams) (nick string, err error)
	AddGuildMemberRole(guildID, userID, roleID Snowflake) (err error)
	RemoveGuildMemberRole(guildID, userID, roleID Snowflake) (err error)
	RemoveGuildMember(guildID, userID Snowflake) (err error)
	GetGuildBans(id Snowflake) (ret []*Ban, err error)
	GetGuildBan(guildID, userID Snowflake) (ret *Ban, err error)
	CreateGuildBan(guildID, userID Snowflake, params *CreateGuildBanParams) (err error)
	RemoveGuildBan(guildID, userID Snowflake) (err error)
	GetGuildRoles(guildID Snowflake) (ret []*Role, err error)
	CreateGuildRole(id Snowflake, params *CreateGuildRoleParams) (ret *Role, err error)
	ModifyGuildRolePositions(guildID Snowflake, params *ModifyGuildRolePositionsParams) (ret []*Role, err error)
	ModifyGuildRole(guildID, roleID Snowflake, params *ModifyGuildRoleParams) (ret *Role, err error)
	DeleteGuildRole(guildID, roleID Snowflake) (err error)
	GetGuildPruneCount(id Snowflake, params *GuildPruneParams) (ret *GuildPruneCount, err error)
	BeginGuildPrune(id Snowflake, params *GuildPruneParams) (ret *GuildPruneCount, err error)
	GetGuildVoiceRegions(id Snowflake) (ret []*VoiceRegion, err error)
	GetGuildInvites(id Snowflake) (ret []*Invite, err error)
	GetGuildIntegrations(id Snowflake) (ret []*Integration, err error)
	CreateGuildIntegration(guildID Snowflake, params *CreateGuildIntegrationParams) (err error)
	ModifyGuildIntegration(guildID, integrationID Snowflake, params *ModifyGuildIntegrationParams) (err error)
	DeleteGuildIntegration(guildID, integrationID Snowflake) (err error)
	SyncGuildIntegration(guildID, integrationID Snowflake) (err error)
	GetGuildEmbed(guildID Snowflake) (ret *GuildEmbed, err error)
	ModifyGuildEmbed(guildID Snowflake, params *GuildEmbed) (ret *GuildEmbed, err error)
	GetGuildVanityURL(guildID Snowflake) (ret *PartialInvite, err error)
}

type InviteRESTer interface {
	GetInvite(inviteCode string, withCounts bool) (invite *Invite, err error)
	DeleteInvite(inviteCode string) (invite *Invite, err error)
}

type UserRESTer interface {
	GetCurrentUser() (ret *User, err error)
	GetUser(id Snowflake) (ret *User, err error)
	ModifyCurrentUser(params *ModifyCurrentUserParams) (ret *User, err error)
	GetCurrentUserGuilds(params *GetCurrentUserGuildsParams) (ret []*Guild, err error)
	LeaveGuild(id Snowflake) (err error)
	GetUserDMs() (ret []*Channel, err error)
	CreateDM(recipientID Snowflake) (ret *Channel, err error)
	CreateGroupDM(params *CreateGroupDMParams) (ret *Channel, err error)
	GetUserConnections() (ret []*UserConnection, err error)
}

type VoiceRESTer interface {
	GetVoiceRegions() (ret []*VoiceRegion, err error)
}

type WebhookRESTer interface {
	CreateWebhook(channelID Snowflake, params *CreateWebhookParams) (ret *Webhook, err error)
	GetChannelWebhooks(channelID Snowflake) (ret []*Webhook, err error)
	GetGuildWebhooks(guildID Snowflake) (ret []*Webhook, err error)
	GetWebhook(id Snowflake) (ret *Webhook, err error)
	GetWebhookWithToken(id Snowflake, token string) (ret *Webhook, err error)
	ModifyWebhook(newWebhook *Webhook) (ret *Webhook, err error)
	ModifyWebhookWithToken(newWebhook *Webhook) (ret *Webhook, err error)
	DeleteWebhook(webhookID Snowflake) (err error)
	DeleteWebhookWithToken(id Snowflake, token string) (err error)
	ExecuteWebhook(params *ExecuteWebhookParams, wait bool, URLSuffix string) (err error)
	ExecuteSlackWebhook(params *ExecuteWebhookParams, wait bool) (err error)
	ExecuteGitHubWebhook(params *ExecuteWebhookParams, wait bool) (err error)
}

type RESTer interface {
	AuditLogsRESTer
	ChannelRESTer
	EmojiRESTer
	GuildRESTer
	InviteRESTer
	UserRESTer
	VoiceRESTer
	WebhookRESTer
}

// The main interface for Disgord
type Session interface {
	// give information about the bot/connected user
	Myself() (*User, error)

	// Request For interacting with Discord. Sending messages, creating channels, guilds, etc.
	// To read object state such as guilds, State() should be used in stead. However some data
	// might not exist in the state. If so it should be requested. Note that this only holds http
	// CRUD operation and not the actual rest endpoints for discord (See Rest()).
	Req() httd.Requester

	// Cache reflects the latest changes received from Discord gateway.
	// Should be used instead of requesting objects.
	Cache() Cacher

	// RateLimiter the rate limiter for the discord REST API
	RateLimiter() httd.RateLimiter

	// Discord Gateway, web socket
	SocketHandler
	HeartbeatLatency() (duration time.Duration, err error)

	// Generic CRUD operations for Discord interaction
	DeleteFromDiscord(obj discordDeleter) error
	SaveToDiscord(obj discordSaver) error

	// state/caching module
	// checks the cache first, otherwise do a http request
	RESTer

	// Custom REST functions
	SendMsg(channelID Snowflake, message *Message) (msg *Message, err error)
	SendMsgString(channelID Snowflake, content string) (msg *Message, err error)
	UpdateMessage(message *Message) (msg *Message, err error)
	UpdateChannel(channel *Channel) (err error)

	// same as above. Except these returns a channel
	// WARNING: none below should be assumed to be working.
	// TODO: implement in the future!
	//GuildChan(guildID Snowflake) <-chan *Guild
	//ChannelChan(channelID Snowflake) <-chan *Channel
	//ChannelsChan(guildID Snowflake) <-chan map[Snowflake]*Channel
	//MsgChan(msgID Snowflake) <-chan *Message
	//UserChan(userID Snowflake) <-chan *UserChan
	//MemberChan(guildID, userID Snowflake) <-chan *Member
	//MembersChan(guildID Snowflake) <-chan map[Snowflake]*Member
}
