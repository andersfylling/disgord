package disgord

import (
	"sync"

	"github.com/andersfylling/disgord/websocket"

	"github.com/andersfylling/disgord/websocket/cmd"
)

// gatewayCmdName is the gateway command name for the payload to be sent to Discord over a websocket connection.
type gatewayCmdName string

type gatewayCmdPayload interface {
	guilder
	distribute(shardCount uint, filter func(guildID websocket.Snowflake) (shardID uint)) (pkts []websocket.ShardDistributer)
}

// proxy pattern(?)
type gatewayCommand struct {
	Data gatewayCmdPayload
	cmd  gatewayCmdName
}

var _ websocket.GatewayCommandPayload = (*gatewayCommand)(nil)
var _ websocket.ShardDistributer = (*gatewayCommand)(nil)

func (g *gatewayCommand) Distribute(shardCount uint, filter func(guildID Snowflake) (shardID uint)) (pkts []websocket.ShardDistributer) {
	pkts = g.Data.distribute(shardCount, filter)
	if len(pkts) == 0 {
		pkts = append(pkts, g) // haxor
	}

	return
}

func (g *gatewayCommand) CmdName() string {
	return string(g.cmd)
}

func (g *gatewayCommand) GetGuildIDs() []websocket.Snowflake {
	return g.Data.getGuildIDs()
}

// #################################################################
// GatewayCmdRequestGuildMembers Used to request offline members for a guild or
// a list of guilds. When initially connecting, the gateway will only send
// offline members if a guild has less than the large_threshold members
// (value in the Gateway Identify). If a Client wishes to receive additional
// members, they need to explicitly request them via this operation. The
// server will send Guild Members Chunk events in response with up to 1000
// members per chunk until all members that match the request have been sent.
const RequestGuildMembers gatewayCmdName = cmd.RequestGuildMembers

// RequestGuildMembersPayload payload for socket command REQUEST_GUILD_MEMBERS.
// See GatewayCmdRequestGuildMembers
//
// WARNING: If this request is in queue while a auto-scaling is forced, it will be removed from the queue
// and not re-inserted like the other commands. This is due to the guild id slice, which is a bit trickier
// to handle.
type RequestGuildMembersPayload struct {
	// GuildID	id of the guild(s) to get offline members for
	GuildIDs []Snowflake `json:"guild_id"`

	// Query string that username starts with, or an empty string to return all members
	Query string `json:"query"`

	// Limit maximum number of members to send or 0 to request all members matched
	Limit uint `json:"limit"`

	// UserIDs used to specify which users you wish to fetch
	UserIDs []Snowflake `json:"user_ids,omitempty"`
}

var _ gatewayCmdPayload = (*RequestGuildMembersPayload)(nil)

func (u *RequestGuildMembersPayload) distribute(shardCount uint, filter func(guildID websocket.Snowflake) (shardID uint)) (pkts []websocket.ShardDistributer) {
	messages := make(map[uint]*RequestGuildMembersPayload)
	for _, guildID := range u.GuildIDs {
		shardID := filter(guildID)
		if _, ok := messages[shardID]; !ok {
			messages[shardID] = &RequestGuildMembersPayload{
				Query: u.Query,
				Limit: u.Limit,
			}
		}
		messages[shardID].GuildIDs = append(messages[shardID].GuildIDs, guildID)
	}

	for _, msg := range messages {
		pkts = append(pkts, msg)
	}
	return pkts
}

func (u *RequestGuildMembersPayload) getGuildIDs() []Snowflake {
	return u.GuildIDs
}

// #################################################################
// UpdateVoiceState Sent when a Client wants to join, move, or
// disconnect from a voice channel.
const UpdateVoiceState gatewayCmdName = cmd.UpdateVoiceState

// UpdateVoiceStatePayload payload for socket command UPDATE_VOICE_STATE.
// see UpdateVoiceState
type UpdateVoiceStatePayload struct {
	// GuildID id of the guild
	GuildID Snowflake `json:"guild_id"`

	// ChannelID id of the voice channel Client wants to join
	// (null if disconnecting)
	ChannelID *Snowflake `json:"channel_id"`

	// SelfMute is the Client mute
	SelfMute bool `json:"self_mute"`

	// SelfDeaf is the Client deafened
	SelfDeaf bool `json:"self_deaf"`
}

var _ gatewayCmdPayload = (*UpdateVoiceStatePayload)(nil)

func (u *UpdateVoiceStatePayload) distribute(nr uint, _ func(websocket.Snowflake) uint) (pkts []websocket.ShardDistributer) {
	return nil
}

func (u *UpdateVoiceStatePayload) getGuildIDs() []Snowflake {
	return []Snowflake{u.GuildID}
}

// #################################################################
// UpdateStatus Sent by the Client to indicate a presence or status
// update.
const UpdateStatus gatewayCmdName = cmd.UpdateStatus

// UpdateStatusPayload payload for socket command UPDATE_STATUS.
// see UpdateStatus
type UpdateStatusPayload struct {
	mu sync.RWMutex
	// Since unix time (in milliseconds) of when the Client went idle, or null if the Client is not idle
	Since *uint `json:"since"`

	// Game null, or the user's new activity
	Game *Activity `json:"game"`

	// Status the user's new status
	Status string `json:"status"`

	// AFK whether or not the Client is afk
	AFK bool `json:"afk"`
}

var _ gatewayCmdPayload = (*UpdateStatusPayload)(nil)

func (u *UpdateStatusPayload) distribute(nr uint, _ func(websocket.Snowflake) uint) (pkts []websocket.ShardDistributer) {
	return nil
}

func (u *UpdateStatusPayload) getGuildIDs() []Snowflake {
	return nil
}
