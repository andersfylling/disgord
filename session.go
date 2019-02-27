package disgord

import (
	"errors"
	"time"

	"github.com/andersfylling/disgord/logger"

	"github.com/andersfylling/snowflake/v3"

	"github.com/andersfylling/disgord/httd"
)

// NewSessionMock returns a session interface that triggers random events allows for fake rest requests.
// Ideal to test the behaviour of your new bot.
// Not implemented!
// TODO: what about a terminal interface for triggering specific events?
func NewSessionMock(conf *Config) (SessionMock, error) {
	return nil, errors.New("not implemented")
}

// EventChannels all methods for retrieving event channels
type EventChannels interface {
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

// Emitter for emitting data from A to B. Used in websocket connection
type Emitter interface {
	Emit(command SocketCommand, dataPointer interface{}) error
}

// Link is used to establish basic commands to create and destroy a link.
// See client.Disconnect() and client.Connect() for linking to the Discord servers
type Link interface {
	Connect() error
	Disconnect() error
}

// SocketHandler all socket related
type SocketHandler interface {
	// Link
	Disconnect() error

	// event handlers
	// inputs are in the following order: middlewares, handlers, controller
	On(event string, inputs ...interface{})
	Emitter

	// event channels
	EventChan(event string) (channel interface{}, err error)
	EventChannels() EventChannels

	// event register (which events to accept)
	// events which are not registered are discarded at socket level
	// to increase performance
	AcceptEvent(events ...string)
}

// AuditLogsRESTer REST interface for all audit-logs endpoints
type RESTAuditLogs interface {
	GetGuildAuditLogs(guildID Snowflake, flags ...Flag) *guildAuditLogsBuilder
}

type RESTMessage interface {
	GetMessages(channelID Snowflake, params URLQueryStringer, flags ...Flag) ([]*Message, error)
	GetMessage(channelID, messageID Snowflake, flags ...Flag) (ret *Message, err error)
	CreateMessage(channelID Snowflake, params *CreateMessageParams, flags ...Flag) (ret *Message, err error)
	UpdateMessage(chanID, msgID Snowflake, params *UpdateMessageParams, flags ...Flag) (ret *Message, err error)
	DeleteMessage(channelID, msgID Snowflake, flags ...Flag) (err error)
	DeleteMessages(chanID Snowflake, params *DeleteMessagesParams, flags ...Flag) (err error)
}

type RESTReaction interface {
	CreateReaction(channelID, messageID Snowflake, emoji interface{}, flags ...Flag) (err error)
	DeleteOwnReaction(channelID, messageID Snowflake, emoji interface{}, flags ...Flag) (err error)
	DeleteUserReaction(channelID, messageID, userID Snowflake, emoji interface{}, flags ...Flag) (err error)
	GetReaction(channelID, messageID Snowflake, emoji interface{}, params URLQueryStringer, flags ...Flag) (ret []*User, err error)
	DeleteAllReactions(channelID, messageID Snowflake, flags ...Flag) (err error)
}

// RESTChannel REST interface for all Channel endpoints
type RESTChannel interface {
	RESTMessage
	RESTReaction
	TriggerTypingIndicator(channelID Snowflake, flags ...Flag) (err error)
	GetPinnedMessages(channelID Snowflake, flags ...Flag) (ret []*Message, err error)
	PinMessage(msg *Message, flags ...Flag) (err error)
	PinMessageID(channelID, msgID Snowflake, flags ...Flag) (err error)
	UnpinMessage(msg *Message, flags ...Flag) (err error)
	UnpinMessageID(channelID, msgID Snowflake, flags ...Flag) (err error)
	GetChannel(id Snowflake, flags ...Flag) (ret *Channel, err error)
	UpdateChannel(id Snowflake, flags ...Flag) (builder *updateChannelBuilder)
	DeleteChannel(id Snowflake, flags ...Flag) (channel *Channel, err error)
	UpdateChannelPermissions(chanID, overwriteID Snowflake, params *UpdateChannelPermissionsParams, flags ...Flag) (err error)
	GetChannelInvites(id Snowflake, flags ...Flag) (ret []*Invite, err error)
	CreateChannelInvites(id Snowflake, params *CreateChannelInvitesParams, flags ...Flag) (ret *Invite, err error)
	DeleteChannelPermission(channelID, overwriteID Snowflake, flags ...Flag) (err error)
	AddDMParticipant(channelID Snowflake, participant *GroupDMParticipant, flags ...Flag) (err error)
	KickParticipant(channelID, userID Snowflake, flags ...Flag) (err error)
}

// RESTEmoji REST interface for all emoji endpoints
type RESTEmoji interface {
	GetGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) (*Emoji, error)
	GetGuildEmojis(id Snowflake, flags ...Flag) ([]*Emoji, error)
	CreateGuildEmoji(guildID Snowflake, params *CreateGuildEmojiParams, flags ...Flag) (*Emoji, error)
	UpdateGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) *updateGuildEmojiBuilder
	DeleteGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) error
}

// RESTGuild REST interface for all guild endpoints
type RESTGuild interface {
	CreateGuild(params *CreateGuildParams, flags ...Flag) (ret *Guild, err error)
	GetGuild(id Snowflake, flags ...Flag) (ret *Guild, err error)
	UpdateGuild(id Snowflake, params *UpdateGuildParams, flags ...Flag) (ret *Guild, err error)
	DeleteGuild(id Snowflake, flags ...Flag) (err error)
	GetGuildChannels(id Snowflake, flags ...Flag) (ret []*Channel, err error)
	CreateGuildChannel(id Snowflake, params *CreateGuildChannelParams, flags ...Flag) (ret *Channel, err error)
	UpdateGuildChannelPositions(id Snowflake, params []UpdateGuildChannelPositionsParams, flags ...Flag) (ret *Guild, err error)
	GetGuildMember(guildID, userID Snowflake, flags ...Flag) (ret *Member, err error)
	GetGuildMembers(guildID, after Snowflake, limit int, flags ...Flag) (ret []*Member, err error)
	AddGuildMember(guildID, userID Snowflake, params *AddGuildMemberParams, flags ...Flag) (ret *Member, err error)
	UpdateGuildMember(guildID, userID Snowflake, params *UpdateGuildMemberParams, flags ...Flag) (err error)
	SetCurrentUserNick(id Snowflake, nick string, flags ...Flag) (newNick string, err error)
	AddGuildMemberRole(guildID, userID, roleID Snowflake, flags ...Flag) (err error)
	RemoveGuildMemberRole(guildID, userID, roleID Snowflake, flags ...Flag) (err error)
	KickMember(guildID, userID Snowflake, flags ...Flag) (err error)
	GetGuildBans(id Snowflake, flags ...Flag) (ret []*Ban, err error)
	GetGuildBan(guildID, userID Snowflake, flags ...Flag) (ret *Ban, err error)
	BanMember(guildID, userID Snowflake, params *BanMemberParams, flags ...Flag) (err error)
	UnbanMember(guildID, userID Snowflake, flags ...Flag) (err error)
	GetGuildRoles(guildID Snowflake, flags ...Flag) (ret []*Role, err error)
	CreateGuildRole(id Snowflake, params *CreateGuildRoleParams, flags ...Flag) (ret *Role, err error)
	UpdateGuildRolePositions(guildID Snowflake, params []UpdateGuildRolePositionsParams, flags ...Flag) (ret []*Role, err error)
	UpdateGuildRole(guildID, roleID Snowflake, flags ...Flag) (builder *modifyGuildRoleBuilder)
	DeleteGuildRole(guildID, roleID Snowflake, flags ...Flag) (err error)
	EstimatePruneMembersCount(id Snowflake, days int, flags ...Flag) (estimate int, err error)
	PruneMembers(id Snowflake, days int, flags ...Flag) error
	GetGuildVoiceRegions(id Snowflake, flags ...Flag) (ret []*VoiceRegion, err error)
	GetGuildInvites(id Snowflake, flags ...Flag) (ret []*Invite, err error)
	GetGuildIntegrations(id Snowflake, flags ...Flag) (ret []*Integration, err error)
	CreateGuildIntegration(guildID Snowflake, params *CreateGuildIntegrationParams, flags ...Flag) (err error)
	UpdateGuildIntegration(guildID, integrationID Snowflake, params *UpdateGuildIntegrationParams, flags ...Flag) (err error)
	DeleteGuildIntegration(guildID, integrationID Snowflake, flags ...Flag) (err error)
	SyncGuildIntegration(guildID, integrationID Snowflake, flags ...Flag) (err error)
	GetGuildEmbed(guildID Snowflake, flags ...Flag) (ret *GuildEmbed, err error)
	UpdateGuildEmbed(guildID Snowflake, params *GuildEmbed, flags ...Flag) (ret *GuildEmbed, err error)
	GetGuildVanityURL(guildID Snowflake, flags ...Flag) (ret *PartialInvite, err error)
}

// RESTInvite REST interface for all invite endpoints
type RESTInvite interface {
	GetInvite(inviteCode string, params URLQueryStringer, flags ...Flag) (*Invite, error)
	DeleteInvite(inviteCode string, flags ...Flag) (deleted *Invite, err error)
}

// RESTUser REST interface for all user endpoints
type RESTUser interface {
	GetCurrentUser(flags ...Flag) (*User, error)
	GetUser(id Snowflake, flags ...Flag) (*User, error)
	UpdateCurrentUser(flags ...Flag) (builder *updateCurrentUserBuilder)
	GetCurrentUserGuilds(params *GetCurrentUserGuildsParams, flags ...Flag) (ret []*PartialGuild, err error)
	LeaveGuild(id Snowflake, flags ...Flag) (err error)
	GetUserDMs(flags ...Flag) (ret []*Channel, err error)
	CreateDM(recipientID Snowflake, flags ...Flag) (ret *Channel, err error)
	CreateGroupDM(params *CreateGroupDMParams, flags ...Flag) (ret *Channel, err error)
	GetUserConnections(flags ...Flag) (ret []*UserConnection, err error)
}

// RESTVoice REST interface for all voice endpoints
type RESTVoice interface {
	GetVoiceRegions(flags ...Flag) ([]*VoiceRegion, error)
}

// RESTWebhook REST interface for all Webhook endpoints
type RESTWebhook interface {
	CreateWebhook(channelID Snowflake, params *CreateWebhookParams, flags ...Flag) (ret *Webhook, err error)
	GetChannelWebhooks(channelID Snowflake, flags ...Flag) (ret []*Webhook, err error)
	GetGuildWebhooks(guildID Snowflake, flags ...Flag) (ret []*Webhook, err error)
	GetWebhook(id Snowflake, flags ...Flag) (ret *Webhook, err error)
	GetWebhookWithToken(id Snowflake, token string, flags ...Flag) (ret *Webhook, err error)
	UpdateWebhook(id Snowflake, flags ...Flag) (builder *updateWebhookBuilder)
	UpdateWebhookWithToken(id Snowflake, token string, flags ...Flag) (builder *updateWebhookBuilder)
	DeleteWebhook(webhookID Snowflake, flags ...Flag) (err error)
	DeleteWebhookWithToken(id Snowflake, token string, flags ...Flag) (err error)
	ExecuteWebhook(params *ExecuteWebhookParams, wait bool, URLSuffix string, flags ...Flag) (err error)
	ExecuteSlackWebhook(params *ExecuteWebhookParams, wait bool, flags ...Flag) (err error)
	ExecuteGitHubWebhook(params *ExecuteWebhookParams, wait bool, flags ...Flag) (err error)
}

// RESTer holds all the sub REST interfaces
type RESTMethods interface {
	RESTAuditLogs
	RESTChannel
	RESTEmoji
	RESTGuild
	RESTInvite
	RESTUser
	RESTVoice
	RESTWebhook
}

// VoiceHandler holds all the voice connection related methods
type VoiceHandler interface {
	VoiceConnect(guildID, channelID Snowflake) (ret VoiceConnection, err error)
}

// Session Is the runtime interface for DisGord. It allows you to interact with a live session (using sockets or not).
// Note that this interface is used after you've configured DisGord, and therefore won't allow you to configure it
// further.
type Session interface {
	// give information about the bot/connected user
	Myself() (*User, error)

	// Request For interacting with Discord. Sending messages, creating channels, guilds, etc.
	// To read object state such as guilds, State() should be used in stead. However some data
	// might not exist in the state. If so it should be requested. Note that this only holds http
	// CRUD operation and not the actual rest endpoints for discord (See Rest()).
	// Deprecated: will be unexported in next breaking release
	Req() httd.Requester

	// Cache reflects the latest changes received from Discord gateway.
	// Should be used instead of requesting objects.
	// Deprecated: will be unexported in next breaking release
	Cache() Cacher

	Logger() logger.Logger

	// RateLimiter the rate limiter for the discord REST API
	// Deprecated: will be unexported in next breaking release
	RateLimiter() httd.RateLimiter

	// Discord Gateway, web socket
	SocketHandler
	HeartbeatLatency() (duration time.Duration, err error)

	// Generic CRUD operations for Discord interaction
	DeleteFromDiscord(obj discordDeleter) error
	SaveToDiscord(original discordSaver, changes ...discordSaver) error

	AddPermission(permission int) (updatedPermissions int)
	GetPermissions() (permissions int)
	CreateBotURL() (u string, err error)

	Pool() *pools

	// state/caching module
	// checks the cacheLink first, otherwise do a http request
	RESTMethods

	// Custom REST functions
	SendMsg(channelID Snowflake, data ...interface{}) (msg *Message, err error)

	// Status update functions
	UpdateStatus(s *UpdateStatusCommand) (err error)
	UpdateStatusString(s string) (err error)

	GetGuilds(params *GetCurrentUserGuildsParams, flags ...Flag) ([]*Guild, error)
	GetConnectedGuilds() []snowflake.ID

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

	// Voice handler, responsible for opening up new voice channel connections
	VoiceHandler
}

type SessionMock interface {
	Session
	// TODO: methods for triggering certain events and controlling states/tracking
}
