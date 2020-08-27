package disgord

import (
	"context"
	"time"

	"github.com/andersfylling/disgord/internal/logger"
)

// Emitter for emitting data from A to B. Used in websocket connection
type Emitter interface {
	Emit(name gatewayCmdName, data gatewayCmdPayload) (unhandledGuildIDs []Snowflake, err error)
}

// Link allows basic Discord connection control. Affects all shards
type Link interface {
	// Connect establishes a websocket connection to the discord API
	Connect(ctx context.Context) error

	// Disconnect closes the discord websocket connection
	Disconnect() error
}

type OnSocketEventer interface {
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

	OnSocketEventer

	// Event gives access to type safe event handler registration using the builder pattern
	Event() SocketHandlerRegistrator

	Emitter
}

type RESTMessage interface {
	// GetMessages Returns the messages for a channel. If operating on a guild channel, this endpoint requires
	// the 'VIEW_CHANNEL' permission to be present on the current user. If the current user is missing
	// the 'READ_MESSAGE_HISTORY' permission in the channel then this will return no messages
	// (since they cannot read the message history). Returns an array of message objects on success.
	GetMessages(ctx context.Context, channelID Snowflake, params *GetMessagesParams, flags ...Flag) ([]*Message, error)

	// GetMessage Returns a specific message in the channel. If operating on a guild channel, this endpoints
	// requires the 'READ_MESSAGE_HISTORY' permission to be present on the current user.
	// Returns a message object on success.
	GetMessage(ctx context.Context, channelID, messageID Snowflake, flags ...Flag) (ret *Message, err error)

	// CreateMessage Post a message to a guild text or DM channel. If operating on a guild channel, this
	// endpoint requires the 'SEND_MESSAGES' permission to be present on the current user. If the tts field is set to true,
	// the SEND_TTS_MESSAGES permission is required for the message to be spoken. Returns a message object. Fires a
	// Message Create Gateway event. See message formatting for more information on how to properly format messages.
	// The maximum request size when sending a message is 8MB.
	CreateMessage(ctx context.Context, channelID Snowflake, params *CreateMessageParams, flags ...Flag) (ret *Message, err error)

	// UpdateMessage Edit a previously sent message. You can only edit messages that have been sent by the
	// current user. Returns a message object. Fires a Message Update Gateway event.
	UpdateMessage(ctx context.Context, chanID, msgID Snowflake, flags ...Flag) *updateMessageBuilder
	SetMsgContent(ctx context.Context, chanID, msgID Snowflake, content string) (*Message, error)
	SetMsgEmbed(ctx context.Context, chanID, msgID Snowflake, embed *Embed) (*Message, error)

	// DeleteMessage Delete a message. If operating on a guild channel and trying to delete a message that was not
	// sent by the current user, this endpoint requires the 'MANAGE_MESSAGES' permission. Returns a 204 empty response
	// on success. Fires a Message Delete Gateway event.
	DeleteMessage(ctx context.Context, channelID, msgID Snowflake, flags ...Flag) (err error)

	// DeleteMessages Delete multiple messages in a single request. This endpoint can only be used on guild
	// Channels and requires the 'MANAGE_MESSAGES' permission. Returns a 204 empty response on success. Fires multiple
	// Message Delete Gateway events.Any message IDs given that do not exist or are invalid will count towards
	// the minimum and maximum message count (currently 2 and 100 respectively). Additionally, duplicated IDs
	// will only be counted once.
	DeleteMessages(ctx context.Context, chanID Snowflake, params *DeleteMessagesParams, flags ...Flag) (err error)
}

type RESTReaction interface {
	// CreateReaction Create a reaction for the message. This endpoint requires the 'READ_MESSAGE_HISTORY'
	// permission to be present on the current user. Additionally, if nobody else has reacted to the message using this
	// emoji, this endpoint requires the 'ADD_REACTIONS' permission to be present on the current user. Returns a 204
	// empty response on success. The maximum request size when sending a message is 8MB.
	CreateReaction(ctx context.Context, channelID, messageID Snowflake, emoji interface{}, flags ...Flag) (err error)

	// DeleteOwnReaction Delete a reaction the current user has made for the message.
	// Returns a 204 empty response on success.
	DeleteOwnReaction(ctx context.Context, channelID, messageID Snowflake, emoji interface{}, flags ...Flag) (err error)

	// DeleteUserReaction Deletes another user's reaction. This endpoint requires the 'MANAGE_MESSAGES' permission
	// to be present on the current user. Returns a 204 empty response on success.
	DeleteUserReaction(ctx context.Context, channelID, messageID, userID Snowflake, emoji interface{}, flags ...Flag) (err error)

	// GetReaction Get a list of Users that reacted with this emoji. Returns an array of user objects on success.
	GetReaction(ctx context.Context, channelID, messageID Snowflake, emoji interface{}, params URLQueryStringer, flags ...Flag) (reactors []*User, err error)

	// DeleteAllReactions Deletes all reactions on a message. This endpoint requires the 'MANAGE_MESSAGES'
	// permission to be present on the current user.
	DeleteAllReactions(ctx context.Context, channelID, messageID Snowflake, flags ...Flag) (err error)
}

// RESTChannel REST interface for all Channel endpoints
type RESTChannel interface {
	RESTMessage
	RESTReaction

	// TriggerTypingIndicator Post a typing indicator for the specified channel. Generally bots should not implement
	// this route. However, if a bot is responding to a command and expects the computation to take a few seconds, this
	// endpoint may be called to let the user know that the bot is processing their message. Returns a 204 empty response
	// on success. Fires a Typing Start Gateway event.
	TriggerTypingIndicator(ctx context.Context, channelID Snowflake, flags ...Flag) (err error)

	// GetPinnedMessages Returns all pinned messages in the channel as an array of message objects.
	GetPinnedMessages(ctx context.Context, channelID Snowflake, flags ...Flag) (ret []*Message, err error)

	// PinMessage same as PinMessageID
	PinMessage(ctx context.Context, msg *Message, flags ...Flag) (err error)

	// PinMessageID Pin a message by its ID and channel ID. Requires the 'MANAGE_MESSAGES' permission.
	// Returns a 204 empty response on success.
	PinMessageID(ctx context.Context, channelID, msgID Snowflake, flags ...Flag) (err error)

	// UnpinMessage same as UnpinMessageID
	UnpinMessage(ctx context.Context, msg *Message, flags ...Flag) (err error)

	// UnpinMessageID Delete a pinned message in a channel. Requires the 'MANAGE_MESSAGES' permission.
	// Returns a 204 empty response on success. Returns a 204 empty response on success.
	UnpinMessageID(ctx context.Context, channelID, msgID Snowflake, flags ...Flag) (err error)

	// GetChannel Get a channel by Snowflake. Returns a channel object.
	GetChannel(ctx context.Context, id Snowflake, flags ...Flag) (ret *Channel, err error)

	// UpdateChannel Update a Channels settings. Requires the 'MANAGE_CHANNELS' permission for the guild. Returns
	// a channel on success, and a 400 BAD REQUEST on invalid parameters. Fires a Channel Update Gateway event. If
	// modifying a category, individual Channel Update events will fire for each child channel that also changes.
	// For the PATCH method, all the JSON Params are optional.
	UpdateChannel(ctx context.Context, id Snowflake, flags ...Flag) (builder *updateChannelBuilder)

	// DeleteChannel Delete a channel, or close a private message. Requires the 'MANAGE_CHANNELS' permission for
	// the guild. Deleting a category does not delete its child Channels; they will have their parent_id removed and a
	// Channel Update Gateway event will fire for each of them. Returns a channel object on success.
	// Fires a Channel Delete Gateway event.
	DeleteChannel(ctx context.Context, id Snowflake, flags ...Flag) (channel *Channel, err error)

	// EditChannelPermissions Edit the channel permission overwrites for a user or role in a channel. Only usable
	// for guild Channels. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty response on success.
	// For more information about permissions, see permissions.
	UpdateChannelPermissions(ctx context.Context, chanID, overwriteID Snowflake, params *UpdateChannelPermissionsParams, flags ...Flag) (err error)

	// GetChannelInvites Returns a list of invite objects (with invite metadata) for the channel. Only usable for
	// guild Channels. Requires the 'MANAGE_CHANNELS' permission.
	GetChannelInvites(ctx context.Context, id Snowflake, flags ...Flag) (ret []*Invite, err error)

	// CreateChannelInvite Create a new invite object for the channel. Only usable for guild Channels. Requires
	// the CREATE_INSTANT_INVITE permission. All JSON parameters for this route are optional, however the request
	// body is not. If you are not sending any fields, you still have to send an empty JSON object ({}).
	// Returns an invite object.
	CreateChannelInvite(ctx context.Context, id Snowflake, flags ...Flag) *createChannelInviteBuilder

	// DeleteChannelPermission Delete a channel permission overwrite for a user or role in a channel. Only usable
	// for guild Channels. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty response on success. For more
	// information about permissions,
	// see permissions: https://discord.com/developers/docs/topics/permissions#permissions
	DeleteChannelPermission(ctx context.Context, channelID, overwriteID Snowflake, flags ...Flag) (err error)

	// AddDMParticipant Adds a recipient to a Group DM using their access token. Returns a 204 empty response
	// on success.
	AddDMParticipant(ctx context.Context, channelID Snowflake, participant *GroupDMParticipant, flags ...Flag) (err error)

	// KickParticipant Removes a recipient from a Group DM. Returns a 204 empty response on success.
	KickParticipant(ctx context.Context, channelID, userID Snowflake, flags ...Flag) (err error)
}

// RESTInvite REST interface for all invite endpoints
type RESTInvite interface {
	// GetInvite Returns an invite object for the given code.
	GetInvite(ctx context.Context, inviteCode string, params URLQueryStringer, flags ...Flag) (*Invite, error)

	// DeleteInvite Delete an invite. Requires the MANAGE_CHANNELS permission. Returns an invite object on success.
	DeleteInvite(ctx context.Context, inviteCode string, flags ...Flag) (deleted *Invite, err error)
}

// RESTUser REST interface for all user endpoints
type RESTUser interface {
	// GetCurrentUser Returns the user object of the requester's account. For OAuth2, this requires the identify
	// scope, which will return the object without an email, and optionally the email scope, which returns the object
	// with an email.
	GetCurrentUser(ctx context.Context, flags ...Flag) (*User, error)

	// GetUser Returns a user object for a given user Snowflake.
	GetUser(ctx context.Context, id Snowflake, flags ...Flag) (*User, error)

	// UpdateCurrentUser Modify the requester's user account settings. Returns a user object on success.
	UpdateCurrentUser(ctx context.Context, flags ...Flag) (builder *updateCurrentUserBuilder)

	// GetCurrentUserGuilds Returns a list of partial guild objects the current user is a member of.
	// Requires the Guilds OAuth2 scope.
	GetCurrentUserGuilds(ctx context.Context, params *GetCurrentUserGuildsParams, flags ...Flag) (ret []*PartialGuild, err error)

	// LeaveGuild Leave a guild. Returns a 204 empty response on success.
	LeaveGuild(ctx context.Context, id Snowflake, flags ...Flag) (err error)

	// GetUserDMs Returns a list of DM channel objects.
	GetUserDMs(ctx context.Context, flags ...Flag) (ret []*Channel, err error)

	// CreateDM Create a new DM channel with a user. Returns a DM channel object.
	CreateDM(ctx context.Context, recipientID Snowflake, flags ...Flag) (ret *Channel, err error)

	// CreateGroupDM Create a new group DM channel with multiple Users. Returns a DM channel object.
	// This endpoint was intended to be used with the now-deprecated GameBridge SDK. DMs created with this
	// endpoint will not be shown in the Discord Client
	CreateGroupDM(ctx context.Context, params *CreateGroupDMParams, flags ...Flag) (ret *Channel, err error)

	// GetUserConnections Returns a list of connection objects. Requires the connections OAuth2 scope.
	GetUserConnections(ctx context.Context, flags ...Flag) (ret []*UserConnection, err error)
}

// RESTVoice REST interface for all voice endpoints
type RESTVoice interface {
	// GetVoiceRegionsBuilder Returns an array of voice region objects that can be used when creating servers.
	GetVoiceRegions(ctx context.Context, flags ...Flag) ([]*VoiceRegion, error)
}

// RESTWebhook REST interface for all Webhook endpoints
type RESTWebhook interface {
	// CreateWebhook Create a new webhook. Requires the 'MANAGE_WEBHOOKS' permission.
	// Returns a webhook object on success.
	CreateWebhook(ctx context.Context, channelID Snowflake, params *CreateWebhookParams, flags ...Flag) (ret *Webhook, err error)

	// GetChannelWebhooks Returns a list of channel webhook objects. Requires the 'MANAGE_WEBHOOKS' permission.
	GetChannelWebhooks(ctx context.Context, channelID Snowflake, flags ...Flag) (ret []*Webhook, err error)

	// GetWebhook Returns the new webhook object for the given id.
	GetWebhook(ctx context.Context, id Snowflake, flags ...Flag) (ret *Webhook, err error)

	// GetWebhookWithToken Same as GetWebhook, except this call does not require authentication and
	// returns no user in the webhook object.
	GetWebhookWithToken(ctx context.Context, id Snowflake, token string, flags ...Flag) (ret *Webhook, err error)

	// UpdateWebhook Modify a webhook. Requires the 'MANAGE_WEBHOOKS' permission.
	// Returns the updated webhook object on success.
	UpdateWebhook(ctx context.Context, id Snowflake, flags ...Flag) (builder *updateWebhookBuilder)

	// UpdateWebhookWithToken Same as UpdateWebhook, except this call does not require authentication,
	// does _not_ accept a channel_id parameter in the body, and does not return a user in the webhook object.
	UpdateWebhookWithToken(ctx context.Context, id Snowflake, token string, flags ...Flag) (builder *updateWebhookBuilder)

	// DeleteWebhook Delete a webhook permanently. User must be owner. Returns a 204 NO CONTENT response on success.
	DeleteWebhook(ctx context.Context, webhookID Snowflake, flags ...Flag) error

	// DeleteWebhookWithToken Same as DeleteWebhook, except this call does not require authentication.
	DeleteWebhookWithToken(ctx context.Context, id Snowflake, token string, flags ...Flag) error

	// ExecuteWebhook Trigger a webhook in Discord.
	ExecuteWebhook(ctx context.Context, params *ExecuteWebhookParams, wait bool, URLSuffix string, flags ...Flag) (*Message, error)

	// ExecuteSlackWebhook Trigger a webhook in Discord from the Slack app.
	ExecuteSlackWebhook(ctx context.Context, params *ExecuteWebhookParams, wait bool, flags ...Flag) (*Message, error)

	// ExecuteGitHubWebhook Trigger a webhook in Discord from the GitHub app.
	ExecuteGitHubWebhook(ctx context.Context, params *ExecuteWebhookParams, wait bool, flags ...Flag) (*Message, error)
}

// RESTer holds all the sub REST interfaces
type RESTMethods interface {
	RESTChannel
	RESTInvite
	RESTUser
	RESTVoice
	RESTWebhook
}

// Session Is the runtime interface for Disgord. It allows you to interact with a live session (using sockets or not).
// Note that this interface is used after you've configured Disgord, and therefore won't allow you to configure it
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

	RESTRatelimitBuckets() (group map[string][]string)

	// Abstract REST methods for Discord structs
	DeleteFromDiscord(ctx context.Context, obj discordDeleter, flags ...Flag) error

	// AddPermission is to store the permissions required by the bot to function as intended.
	AddPermission(permission PermissionBit) (updatedPermissions PermissionBit)
	GetPermissions() (permissions PermissionBit)

	// CreateBotURL
	InviteURL(ctx context.Context) (url string, err error)

	Pool() *pools

	RESTMethods

	// CreateGuild Create a new guild. Returns a guild object on success. Fires a Guild Create Gateway event.
	CreateGuild(ctx context.Context, guildName string, params *CreateGuildParams, flags ...Flag) (*Guild, error)

	// Guild is used to create a guild query builder.
	Guild(id Snowflake) GuildQueryBuilder

	// Custom REST functions
	SendMsg(ctx context.Context, channelID Snowflake, data ...interface{}) (*Message, error)

	// Status update functions
	UpdateStatus(s *UpdateStatusPayload) error
	UpdateStatusString(s string) error

	GetGuilds(ctx context.Context, params *GetCurrentUserGuildsParams, flags ...Flag) ([]*Guild, error)
	GetConnectedGuilds() []Snowflake
}
