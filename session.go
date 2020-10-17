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

	// AddPermission is to store the permissions required by the bot to function as intended.
	AddPermission(permission PermissionBit) (updatedPermissions PermissionBit)
	GetPermissions() (permissions PermissionBit)

	Pool() *pools

	ClientQueryBuilder

	// Status update functions
	UpdateStatus(s *UpdateStatusPayload) error
	UpdateStatusString(s string) error

	GetGuilds(ctx context.Context, params *GetCurrentUserGuildsParams, flags ...Flag) ([]*Guild, error)
	GetConnectedGuilds() []Snowflake
}
