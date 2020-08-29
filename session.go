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
	RESTVoice
	RESTWebhook

	Invite(code string) InviteQueryBuilder

	Channel(cid Snowflake) ChannelQueryBuilder

	User(uid Snowflake) UserQueryBuilder

	CurrentUser() CurrentUserQueryBuilder

	// CreateGuild Create a new guild. Returns a guild object on success. Fires a Guild Create Gateway event.
	CreateGuild(ctx context.Context, guildName string, params *CreateGuildParams, flags ...Flag) (*Guild, error)

	// Guild is used to create a guild query builder.
	Guild(id Snowflake) GuildQueryBuilder
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

	// Custom REST functions
	SendMsg(ctx context.Context, channelID Snowflake, data ...interface{}) (*Message, error)

	// Status update functions
	UpdateStatus(s *UpdateStatusPayload) error
	UpdateStatusString(s string) error

	GetGuilds(ctx context.Context, params *GetCurrentUserGuildsParams, flags ...Flag) ([]*Guild, error)
	GetConnectedGuilds() []Snowflake
}
