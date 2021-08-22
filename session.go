package disgord

import (
	"context"
	"time"

	"github.com/andersfylling/disgord/internal/logger"
)

// Session Is the runtime interface for Disgord. It allows you to interact with a live session (using sockets or not).
// Note that this interface is used after you've configured Disgord, and therefore won't allow you to configure it
// further.
type Session interface {
	// Logger returns the injected logger instance. If nothing was injected, a empty wrapper is returned
	// to avoid nil panics.
	Logger() logger.Logger

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
	EditInteractionResponse(ctx context.Context, interaction *InteractionCreate, message *Message) error
	SendInteractionResponse(context context.Context, interaction *InteractionCreate, data *InteractionResponse) error

	CreateGlobalCommand(ctx context.Context, applicationID Snowflake, command *ApplicationCommand) error
	CreateGuildCommand(ctx context.Context, applicationID Snowflake, guildID Snowflake, command *ApplicationCommand) error
	UpdateCommandGlobal(ctx context.Context, applicationID, commandID Snowflake, command *ApplicationCommand) error
	UpdateCommandGuild(ctx context.Context, applicationID, guildID, commandID Snowflake, command *ApplicationCommand) error
	DeleteCommandGlobal(ctx context.Context, applicationID, commandID Snowflake) error
	DeleteCommandGuild(ctx context.Context, applicationID, guildID, commandID Snowflake) error

	// Status update functions
	UpdateStatus(s *UpdateStatusPayload) error
	UpdateStatusString(s string) error

	GetConnectedGuilds() []Snowflake
}
