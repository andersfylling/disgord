package disgord

import (
	"time"

	"github.com/andersfylling/disgord/logger"
)

// Emitter for emitting data from A to B. Used in websocket connection
type Emitter interface {
	Emit(command SocketCommand, dataPointer interface{}) error
}

// Link allows basic Discord connection control. Affects all shards
type Link interface {
	// Connect establishes a websocket connection to the discord API
	Connect() error

	// Disconnect closes the discord websocket connection
	Disconnect() error
}

// SocketHandler all socket related logic
type SocketHandler interface {
	// Link controls the connection to the Discord API. Affects all shards.
	// Link

	// Disconnect closes the discord websocket connection
	Disconnect() error

	// Suspend temporary closes the socket connection, allowing resources to be
	// reused on reconnect
	Suspend() error

	// On creates a specification to be executed on the given event. The specification
	// consists of, in order, 0 or more middlewares, 1 or more handlers, 0 or 1 controller.
	// On incorrect ordering, or types, the method will panic. See reactor.go for types.
	//
	// Each of the three sub-types of a specification is run in sequence, as well as the specifications
	// registered for a event. However, the slice of specifications are executed in a goroutine to avoid
	// blocking future events. The middlewares allows manipulating the event data before it reaches the
	// handlers. The handlers executes short-running logic based on the event data (use go routine if
	// you need a long running task). The controller dictates lifetime of the specification.
	//
	//  // a handler that is executed on every Ready event
	//  Client.On(EvtReady, onReady)
	//
	//  // a handler that runs only the first three times a READY event is fired
	//  Client.On(EvtReady, onReady, &Ctrl{Runs: 3})
	//
	//  // a handler that only runs for events within the first 10 minutes
	//  Client.On(EvtReady, onReady, &Ctrl{Duration: 10*time.Minute})
	On(event string, inputs ...interface{})

	Emitter
}

// AuditLogsRESTer REST interface for all audit-logs endpoints
type RESTAuditLogs interface {
	// GetGuildAuditLogs Returns an audit log object for the guild. Requires the 'VIEW_AUDIT_LOG' permission.
	// Note that this request will _always_ send a REST request, regardless of you calling IgnoreCache or not.
	GetGuildAuditLogs(guildID Snowflake, flags ...Flag) *guildAuditLogsBuilder
}

type RESTMessage interface {
	// GetMessages Returns the messages for a channel. If operating on a guild channel, this endpoint requires
	// the 'VIEW_CHANNEL' permission to be present on the current user. If the current user is missing
	// the 'READ_MESSAGE_HISTORY' permission in the channel then this will return no messages
	// (since they cannot read the message history). Returns an array of message objects on success.
	GetMessages(channelID Snowflake, params *GetMessagesParams, flags ...Flag) ([]*Message, error)

	// GetMessage Returns a specific message in the channel. If operating on a guild channel, this endpoints
	// requires the 'READ_MESSAGE_HISTORY' permission to be present on the current user.
	// Returns a message object on success.
	GetMessage(channelID, messageID Snowflake, flags ...Flag) (ret *Message, err error)

	// CreateMessage Post a message to a guild text or DM channel. If operating on a guild channel, this
	// endpoint requires the 'SEND_MESSAGES' permission to be present on the current user. If the tts field is set to true,
	// the SEND_TTS_MESSAGES permission is required for the message to be spoken. Returns a message object. Fires a
	// Message Create Gateway event. See message formatting for more information on how to properly format messages.
	// The maximum request size when sending a message is 8MB.
	CreateMessage(channelID Snowflake, params *CreateMessageParams, flags ...Flag) (ret *Message, err error)

	// UpdateMessage Edit a previously sent message. You can only edit messages that have been sent by the
	// current user. Returns a message object. Fires a Message Update Gateway event.
	UpdateMessage(chanID, msgID Snowflake, flags ...Flag) *updateMessageBuilder
	SetMsgContent(chanID, msgID Snowflake, content string) (*Message, error)
	SetMsgEmbed(chanID, msgID Snowflake, embed *Embed) (*Message, error)

	// DeleteMessage Delete a message. If operating on a guild channel and trying to delete a message that was not
	// sent by the current user, this endpoint requires the 'MANAGE_MESSAGES' permission. Returns a 204 empty response
	// on success. Fires a Message Delete Gateway event.
	DeleteMessage(channelID, msgID Snowflake, flags ...Flag) (err error)

	// DeleteMessages Delete multiple messages in a single request. This endpoint can only be used on guild
	// channels and requires the 'MANAGE_MESSAGES' permission. Returns a 204 empty response on success. Fires multiple
	// Message Delete Gateway events.Any message IDs given that do not exist or are invalid will count towards
	// the minimum and maximum message count (currently 2 and 100 respectively). Additionally, duplicated IDs
	// will only be counted once.
	DeleteMessages(chanID Snowflake, params *DeleteMessagesParams, flags ...Flag) (err error)
}

type RESTReaction interface {
	// CreateReaction Create a reaction for the message. This endpoint requires the 'READ_MESSAGE_HISTORY'
	// permission to be present on the current user. Additionally, if nobody else has reacted to the message using this
	// emoji, this endpoint requires the 'ADD_REACTIONS' permission to be present on the current user. Returns a 204
	// empty response on success. The maximum request size when sending a message is 8MB.
	CreateReaction(channelID, messageID Snowflake, emoji interface{}, flags ...Flag) (err error)

	// DeleteOwnReaction Delete a reaction the current user has made for the message.
	// Returns a 204 empty response on success.
	DeleteOwnReaction(channelID, messageID Snowflake, emoji interface{}, flags ...Flag) (err error)

	// DeleteUserReaction Deletes another user's reaction. This endpoint requires the 'MANAGE_MESSAGES' permission
	// to be present on the current user. Returns a 204 empty response on success.
	DeleteUserReaction(channelID, messageID, userID Snowflake, emoji interface{}, flags ...Flag) (err error)

	// GetReaction Get a list of users that reacted with this emoji. Returns an array of user objects on success.
	GetReaction(channelID, messageID Snowflake, emoji interface{}, params URLQueryStringer, flags ...Flag) (reactors []*User, err error)

	// DeleteAllReactions Deletes all reactions on a message. This endpoint requires the 'MANAGE_MESSAGES'
	// permission to be present on the current user.
	DeleteAllReactions(channelID, messageID Snowflake, flags ...Flag) (err error)
}

// RESTChannel REST interface for all Channel endpoints
type RESTChannel interface {
	RESTMessage
	RESTReaction

	// TriggerTypingIndicator Post a typing indicator for the specified channel. Generally bots should not implement
	// this route. However, if a bot is responding to a command and expects the computation to take a few seconds, this
	// endpoint may be called to let the user know that the bot is processing their message. Returns a 204 empty response
	// on success. Fires a Typing Start Gateway event.
	TriggerTypingIndicator(channelID Snowflake, flags ...Flag) (err error)

	// GetPinnedMessages Returns all pinned messages in the channel as an array of message objects.
	GetPinnedMessages(channelID Snowflake, flags ...Flag) (ret []*Message, err error)

	// PinMessage same as PinMessageID
	PinMessage(msg *Message, flags ...Flag) (err error)

	// PinMessageID Pin a message by its ID and channel ID. Requires the 'MANAGE_MESSAGES' permission.
	// Returns a 204 empty response on success.
	PinMessageID(channelID, msgID Snowflake, flags ...Flag) (err error)

	// UnpinMessage same as UnpinMessageID
	UnpinMessage(msg *Message, flags ...Flag) (err error)

	// UnpinMessageID Delete a pinned message in a channel. Requires the 'MANAGE_MESSAGES' permission.
	// Returns a 204 empty response on success. Returns a 204 empty response on success.
	UnpinMessageID(channelID, msgID Snowflake, flags ...Flag) (err error)

	// GetChannel Get a channel by Snowflake. Returns a channel object.
	GetChannel(id Snowflake, flags ...Flag) (ret *Channel, err error)

	// UpdateChannel Update a channels settings. Requires the 'MANAGE_CHANNELS' permission for the guild. Returns
	// a channel on success, and a 400 BAD REQUEST on invalid parameters. Fires a Channel Update Gateway event. If
	// modifying a category, individual Channel Update events will fire for each child channel that also changes.
	// For the PATCH method, all the JSON Params are optional.
	UpdateChannel(id Snowflake, flags ...Flag) (builder *updateChannelBuilder)

	// DeleteChannel Delete a channel, or close a private message. Requires the 'MANAGE_CHANNELS' permission for
	// the guild. Deleting a category does not delete its child channels; they will have their parent_id removed and a
	// Channel Update Gateway event will fire for each of them. Returns a channel object on success.
	// Fires a Channel Delete Gateway event.
	DeleteChannel(id Snowflake, flags ...Flag) (channel *Channel, err error)

	// EditChannelPermissions Edit the channel permission overwrites for a user or role in a channel. Only usable
	// for guild channels. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty response on success.
	// For more information about permissions, see permissions.
	UpdateChannelPermissions(chanID, overwriteID Snowflake, params *UpdateChannelPermissionsParams, flags ...Flag) (err error)

	// GetChannelInvites Returns a list of invite objects (with invite metadata) for the channel. Only usable for
	// guild channels. Requires the 'MANAGE_CHANNELS' permission.
	GetChannelInvites(id Snowflake, flags ...Flag) (ret []*Invite, err error)

	// CreateChannelInvites Create a new invite object for the channel. Only usable for guild channels. Requires
	// the CREATE_INSTANT_INVITE permission. All JSON parameters for this route are optional, however the request
	// body is not. If you are not sending any fields, you still have to send an empty JSON object ({}).
	// Returns an invite object.
	CreateChannelInvites(id Snowflake, params *CreateChannelInvitesParams, flags ...Flag) (ret *Invite, err error)

	// DeleteChannelPermission Delete a channel permission overwrite for a user or role in a channel. Only usable
	// for guild channels. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty response on success. For more
	// information about permissions,
	// see permissions: https://discordapp.com/developers/docs/topics/permissions#permissions
	DeleteChannelPermission(channelID, overwriteID Snowflake, flags ...Flag) (err error)

	// AddDMParticipant Adds a recipient to a Group DM using their access token. Returns a 204 empty response
	// on success.
	AddDMParticipant(channelID Snowflake, participant *GroupDMParticipant, flags ...Flag) (err error)

	// KickParticipant Removes a recipient from a Group DM. Returns a 204 empty response on success.
	KickParticipant(channelID, userID Snowflake, flags ...Flag) (err error)
}

// RESTEmoji REST interface for all emoji endpoints
type RESTEmoji interface {
	// GetGuildEmoji Returns an emoji object for the given guild and emoji IDs.
	GetGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) (*Emoji, error)

	// GetGuildEmojis Returns a list of emoji objects for the given guild.
	GetGuildEmojis(id Snowflake, flags ...Flag) ([]*Emoji, error)

	// CreateGuildEmoji Create a new emoji for the guild. Requires the 'MANAGE_EMOJIS' permission.
	// Returns the new emoji object on success. Fires a Guild Emojis Update Gateway event.
	CreateGuildEmoji(guildID Snowflake, params *CreateGuildEmojiParams, flags ...Flag) (*Emoji, error)

	// UpdateGuildEmoji Modify the given emoji. Requires the 'MANAGE_EMOJIS' permission.
	// Returns the updated emoji object on success. Fires a Guild Emojis Update Gateway event.
	UpdateGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) *updateGuildEmojiBuilder

	// DeleteGuildEmoji Delete the given emoji. Requires the 'MANAGE_EMOJIS' permission. Returns 204 No Content on
	// success. Fires a Guild Emojis Update Gateway event.
	DeleteGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) error
}

// RESTGuild REST interface for all guild endpoints
type RESTGuild interface {
	// CreateGuild Create a new guild. Returns a guild object on success. Fires a Guild Create Gateway event.
	CreateGuild(guildName string, params *CreateGuildParams, flags ...Flag) (*Guild, error)

	// GetGuild Returns the guild object for the given id.
	GetGuild(id Snowflake, flags ...Flag) (*Guild, error)

	// ModifyGuild Modify a guild's settings. Requires the 'MANAGE_GUILD' permission. Returns the updated guild
	// object on success. Fires a Guild Update Gateway event.
	UpdateGuild(id Snowflake, flags ...Flag) *updateGuildBuilder

	// DeleteGuild Delete a guild permanently. User must be owner. Returns 204 No Content on success.
	// Fires a Guild Delete Gateway event.
	DeleteGuild(id Snowflake, flags ...Flag) error

	// GetGuildChannels Returns a list of guild channel objects.
	GetGuildChannels(id Snowflake, flags ...Flag) ([]*Channel, error)

	// CreateGuildChannel Create a new channel object for the guild. Requires the 'MANAGE_CHANNELS' permission.
	// Returns the new channel object on success. Fires a Channel Create Gateway event.
	CreateGuildChannel(id Snowflake, name string, params *CreateGuildChannelParams, flags ...Flag) (*Channel, error)

	// UpdateGuildChannelPositions Modify the positions of a set of channel objects for the guild.
	// Requires 'MANAGE_CHANNELS' permission. Returns a 204 empty response on success. Fires multiple Channel Update
	// Gateway events.
	UpdateGuildChannelPositions(id Snowflake, params []UpdateGuildChannelPositionsParams, flags ...Flag) error

	// GetMember Returns a guild member object for the specified user.
	GetMember(guildID, userID Snowflake, flags ...Flag) (*Member, error)

	// GetMembers uses the GetGuildMembers endpoint iteratively until your query params are met.
	GetMembers(guildID Snowflake, params *GetMembersParams, flags ...Flag) ([]*Member, error)

	// AddGuildMember Adds a user to the guild, provided you have a valid oauth2 access token for the user with
	// the guilds.join scope. Returns a 201 Created with the guild member as the body, or 204 No Content if the user is
	// already a member of the guild. Fires a Guild Member Add Gateway event. Requires the bot to have the
	// CREATE_INSTANT_INVITE permission.
	AddGuildMember(guildID, userID Snowflake, accessToken string, params *AddGuildMemberParams, flags ...Flag) (*Member, error)

	// ModifyGuildMember Modify attributes of a guild member. Returns a 204 empty response on success.
	// Fires a Guild Member Update Gateway event.
	UpdateGuildMember(guildID, userID Snowflake, flags ...Flag) *updateGuildMemberBuilder

	// SetCurrentUserNick Modifies the nickname of the current user in a guild. Returns a 200
	// with the nickname on success. Fires a Guild Member Update Gateway event.
	SetCurrentUserNick(id Snowflake, nick string, flags ...Flag) (newNick string, err error)

	// AddGuildMemberRole Adds a role to a guild member. Requires the 'MANAGE_ROLES' permission.
	// Returns a 204 empty response on success. Fires a Guild Member Update Gateway event.
	AddGuildMemberRole(guildID, userID, roleID Snowflake, flags ...Flag) error

	// RemoveGuildMemberRole Removes a role from a guild member. Requires the 'MANAGE_ROLES' permission.
	// Returns a 204 empty response on success. Fires a Guild Member Update Gateway event.
	RemoveGuildMemberRole(guildID, userID, roleID Snowflake, flags ...Flag) error

	// RemoveGuildMember Remove a member from a guild. Requires 'KICK_MEMBERS' permission.
	// Returns a 204 empty response on success. Fires a Guild Member Remove Gateway event.
	KickMember(guildID, userID Snowflake, flags ...Flag) error

	// GetGuildBans Returns a list of ban objects for the users banned from this guild. Requires the 'BAN_MEMBERS' permission.
	GetGuildBans(id Snowflake, flags ...Flag) ([]*Ban, error)

	// GetGuildBan Returns a ban object for the given user or a 404 not found if the ban cannot be found.
	// Requires the 'BAN_MEMBERS' permission.
	GetGuildBan(guildID, userID Snowflake, flags ...Flag) (*Ban, error)

	// BanMember Create a guild ban, and optionally delete previous messages sent by the banned user. Requires
	// the 'BAN_MEMBERS' permission. Returns a 204 empty response on success. Fires a Guild Ban Add Gateway event.
	BanMember(guildID, userID Snowflake, params *BanMemberParams, flags ...Flag) error

	// UnbanMember Remove the ban for a user. Requires the 'BAN_MEMBERS' permissions.
	// Returns a 204 empty response on success. Fires a Guild Ban Remove Gateway event.
	UnbanMember(guildID, userID Snowflake, flags ...Flag) error

	// GetGuildRoles Returns a list of role objects for the guild.
	GetGuildRoles(guildID Snowflake, flags ...Flag) ([]*Role, error)

	GetMemberPermissions(guildID, userID Snowflake, flags ...Flag) (permissions PermissionBits, err error)

	// CreateGuildRole Create a new role for the guild. Requires the 'MANAGE_ROLES' permission.
	// Returns the new role object on success. Fires a Guild Role Create Gateway event.
	CreateGuildRole(id Snowflake, params *CreateGuildRoleParams, flags ...Flag) (*Role, error)

	// UpdateGuildRolePositions Modify the positions of a set of role objects for the guild.
	// Requires the 'MANAGE_ROLES' permission. Returns a list of all of the guild's role objects on success.
	// Fires multiple Guild Role Update Gateway events.
	UpdateGuildRolePositions(guildID Snowflake, params []UpdateGuildRolePositionsParams, flags ...Flag) ([]*Role, error)

	// ModifyGuildRole Modify a guild role. Requires the 'MANAGE_ROLES' permission.
	// Returns the updated role on success. Fires a Guild Role Update Gateway event.
	UpdateGuildRole(guildID, roleID Snowflake, flags ...Flag) (builder *updateGuildRoleBuilder)

	// DeleteGuildRole Delete a guild role. Requires the 'MANAGE_ROLES' permission.
	// Returns a 204 empty response on success. Fires a Guild Role Delete Gateway event.
	DeleteGuildRole(guildID, roleID Snowflake, flags ...Flag) error

	// EstimatePruneMembersCount Returns an object with one 'pruned' key indicating the number of members that would be
	// removed in a prune operation. Requires the 'KICK_MEMBERS' permission.
	EstimatePruneMembersCount(id Snowflake, days int, flags ...Flag) (estimate int, err error)

	// PruneMembers Kicks members from N day back. Requires the 'KICK_MEMBERS' permission.
	// The estimate of kicked people is not returned. Use EstimatePruneMembersCount before calling PruneMembers
	// if you need it. Fires multiple Guild Member Remove Gateway events.
	PruneMembers(id Snowflake, days int, flags ...Flag) error

	// GetGuildVoiceRegions Returns a list of voice region objects for the guild. Unlike the similar /voice route,
	// this returns VIP servers when the guild is VIP-enabled.
	GetGuildVoiceRegions(id Snowflake, flags ...Flag) ([]*VoiceRegion, error)

	// GetGuildInvites Returns a list of invite objects (with invite metadata) for the guild.
	// Requires the 'MANAGE_GUILD' permission.
	GetGuildInvites(id Snowflake, flags ...Flag) ([]*Invite, error)

	// GetGuildIntegrations Returns a list of integration objects for the guild.
	// Requires the 'MANAGE_GUILD' permission.
	GetGuildIntegrations(id Snowflake, flags ...Flag) ([]*Integration, error)

	// CreateGuildIntegration Attach an integration object from the current user to the guild.
	// Requires the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
	// Fires a Guild Integrations Update Gateway event.
	CreateGuildIntegration(guildID Snowflake, params *CreateGuildIntegrationParams, flags ...Flag) error

	// UpdateGuildIntegration Modify the behavior and settings of a integration object for the guild.
	// Requires the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
	// Fires a Guild Integrations Update Gateway event.
	UpdateGuildIntegration(guildID, integrationID Snowflake, params *UpdateGuildIntegrationParams, flags ...Flag) error

	// DeleteGuildIntegration Delete the attached integration object for the guild.
	// Requires the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
	// Fires a Guild Integrations Update Gateway event.
	DeleteGuildIntegration(guildID, integrationID Snowflake, flags ...Flag) error

	// SyncGuildIntegration Sync an integration. Requires the 'MANAGE_GUILD' permission.
	// Returns a 204 empty response on success.
	SyncGuildIntegration(guildID, integrationID Snowflake, flags ...Flag) error

	// GetGuildEmbed Returns the guild embed object. Requires the 'MANAGE_GUILD' permission.
	GetGuildEmbed(guildID Snowflake, flags ...Flag) (*GuildEmbed, error)

	// UpdateGuildEmbed Modify a guild embed object for the guild. All attributes may be passed in with JSON and
	// modified. Requires the 'MANAGE_GUILD' permission. Returns the updated guild embed object.
	UpdateGuildEmbed(guildID Snowflake, flags ...Flag) *updateGuildEmbedBuilder

	// GetGuildVanityURL Returns a partial invite object for guilds with that feature enabled.
	// Requires the 'MANAGE_GUILD' permission.
	GetGuildVanityURL(guildID Snowflake, flags ...Flag) (*PartialInvite, error)
}

// RESTInvite REST interface for all invite endpoints
type RESTInvite interface {
	// GetInvite Returns an invite object for the given code.
	GetInvite(inviteCode string, params URLQueryStringer, flags ...Flag) (*Invite, error)

	// DeleteInvite Delete an invite. Requires the MANAGE_CHANNELS permission. Returns an invite object on success.
	DeleteInvite(inviteCode string, flags ...Flag) (deleted *Invite, err error)
}

// RESTUser REST interface for all user endpoints
type RESTUser interface {
	// GetCurrentUser Returns the user object of the requester's account. For OAuth2, this requires the identify
	// scope, which will return the object without an email, and optionally the email scope, which returns the object
	// with an email.
	GetCurrentUser(flags ...Flag) (*User, error)

	// GetUser Returns a user object for a given user Snowflake.
	GetUser(id Snowflake, flags ...Flag) (*User, error)

	// UpdateCurrentUser Modify the requester's user account settings. Returns a user object on success.
	UpdateCurrentUser(flags ...Flag) (builder *updateCurrentUserBuilder)

	// GetCurrentUserGuilds Returns a list of partial guild objects the current user is a member of.
	// Requires the guilds OAuth2 scope.
	GetCurrentUserGuilds(params *GetCurrentUserGuildsParams, flags ...Flag) (ret []*PartialGuild, err error)

	// LeaveGuild Leave a guild. Returns a 204 empty response on success.
	LeaveGuild(id Snowflake, flags ...Flag) (err error)

	// GetUserDMs Returns a list of DM channel objects.
	GetUserDMs(flags ...Flag) (ret []*Channel, err error)

	// CreateDM Create a new DM channel with a user. Returns a DM channel object.
	CreateDM(recipientID Snowflake, flags ...Flag) (ret *Channel, err error)

	// CreateGroupDM Create a new group DM channel with multiple users. Returns a DM channel object.
	// This endpoint was intended to be used with the now-deprecated GameBridge SDK. DMs created with this
	// endpoint will not be shown in the Discord Client
	CreateGroupDM(params *CreateGroupDMParams, flags ...Flag) (ret *Channel, err error)

	// GetUserConnections Returns a list of connection objects. Requires the connections OAuth2 scope.
	GetUserConnections(flags ...Flag) (ret []*UserConnection, err error)
}

// RESTVoice REST interface for all voice endpoints
type RESTVoice interface {
	// GetVoiceRegionsBuilder Returns an array of voice region objects that can be used when creating servers.
	GetVoiceRegions(flags ...Flag) ([]*VoiceRegion, error)
}

// RESTWebhook REST interface for all Webhook endpoints
type RESTWebhook interface {
	// CreateWebhook Create a new webhook. Requires the 'MANAGE_WEBHOOKS' permission.
	// Returns a webhook object on success.
	CreateWebhook(channelID Snowflake, params *CreateWebhookParams, flags ...Flag) (ret *Webhook, err error)

	// GetChannelWebhooks Returns a list of channel webhook objects. Requires the 'MANAGE_WEBHOOKS' permission.
	GetChannelWebhooks(channelID Snowflake, flags ...Flag) (ret []*Webhook, err error)

	// GetGuildWebhooks Returns a list of guild webhook objects. Requires the 'MANAGE_WEBHOOKS' permission.
	GetGuildWebhooks(guildID Snowflake, flags ...Flag) (ret []*Webhook, err error)

	// GetWebhook Returns the new webhook object for the given id.
	GetWebhook(id Snowflake, flags ...Flag) (ret *Webhook, err error)

	// GetWebhookWithToken Same as GetWebhook, except this call does not require authentication and
	// returns no user in the webhook object.
	GetWebhookWithToken(id Snowflake, token string, flags ...Flag) (ret *Webhook, err error)

	// UpdateWebhook Modify a webhook. Requires the 'MANAGE_WEBHOOKS' permission.
	// Returns the updated webhook object on success.
	UpdateWebhook(id Snowflake, flags ...Flag) (builder *updateWebhookBuilder)

	// UpdateWebhookWithToken Same as UpdateWebhook, except this call does not require authentication,
	// does _not_ accept a channel_id parameter in the body, and does not return a user in the webhook object.
	UpdateWebhookWithToken(id Snowflake, token string, flags ...Flag) (builder *updateWebhookBuilder)

	// DeleteWebhook Delete a webhook permanently. User must be owner. Returns a 204 NO CONTENT response on success.
	DeleteWebhook(webhookID Snowflake, flags ...Flag) error

	// DeleteWebhookWithToken Same as DeleteWebhook, except this call does not require authentication.
	DeleteWebhookWithToken(id Snowflake, token string, flags ...Flag) error

	// ExecuteWebhook Trigger a webhook in Discord.
	ExecuteWebhook(params *ExecuteWebhookParams, wait bool, URLSuffix string, flags ...Flag) error

	// ExecuteSlackWebhook Trigger a webhook in Discord from the Slack app.
	ExecuteSlackWebhook(params *ExecuteWebhookParams, wait bool, flags ...Flag) error

	// ExecuteGitHubWebhook Trigger a webhook in Discord from the GitHub app.
	ExecuteGitHubWebhook(params *ExecuteWebhookParams, wait bool, flags ...Flag) error
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
	// Logger returns the injected logger instance. If nothing was injected, a empty wrapper is returned
	// to avoid nil panics.
	Logger() logger.Logger

	// Discord Gateway, web socket
	SocketHandler

	// HeartbeatLatency returns the avg. ish time used to send and receive a heartbeat signal.
	// The latency is calculated as such:
	// 0. start timer (start)
	// 1. send heartbeat signal
	// 2. wait until a heartbeat ack is sent by Discord
	// 3. latency = time.Now().Sub(start)
	// 4. avg = (avg + latency) / 2
	//
	// This feature was requested. But should never be used as a proof for delay between client and Discord.
	AvgHeartbeatLatency() (duration time.Duration, err error)
	// returns the latency for each given shard id. shardID => latency
	HeartbeatLatencies() (latencies map[uint]time.Duration, err error)

	// Abstract REST methods for Discord structs
	DeleteFromDiscord(obj discordDeleter, flags ...Flag) error

	// AddPermission is to store the permissions required by the bot to function as intended.
	AddPermission(permission PermissionBit) (updatedPermissions PermissionBits)
	GetPermissions() (permissions PermissionBits)

	// CreateBotURL
	CreateBotURL() (url string, err error)

	Pool() *pools

	RESTMethods

	// Custom REST functions
	SendMsg(channelID Snowflake, data ...interface{}) (*Message, error)

	KickVoiceParticipant(guildID, userID Snowflake) error

	// Status update functions
	UpdateStatus(s *UpdateStatusCommand) error
	UpdateStatusString(s string) error

	GetGuilds(params *GetCurrentUserGuildsParams, flags ...Flag) ([]*Guild, error)
	GetConnectedGuilds() []Snowflake

	// Voice handler, responsible for opening up new voice channel connections
	VoiceHandler
}
