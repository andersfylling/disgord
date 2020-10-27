package disgord

import (
	"github.com/andersfylling/disgord/internal/gateway"
	"github.com/andersfylling/disgord/internal/gateway/cmd"
)

// gatewayCmdName is the gateway command name for the payload to be sent to Discord over a websocket connection.
type gatewayCmdName string

const (
	// GatewayCmdRequestGuildMembers Used to request offline members for a guild or
	// a list of Guilds. When initially connecting, the gateway will only send
	// offline members if a guild has less than the large_threshold members
	// (value in the Gateway Identify). If a Client wishes to receive additional
	// members, they need to explicitly request them via this operation. The
	// server will send Guild Members Chunk events in response with up to 1000
	// members per chunk until all members that match the request have been sent.
	RequestGuildMembers gatewayCmdName = cmd.RequestGuildMembers

	// UpdateVoiceState Sent when a Client wants to join, move, or
	// disconnect from a voice channel.
	UpdateVoiceState gatewayCmdName = cmd.UpdateVoiceState

	// UpdateStatus Sent by the Client to indicate a presence or status
	// update.
	UpdateStatus gatewayCmdName = cmd.UpdateStatus
)

// #################################################################
// RequestGuildMembersPayload payload for socket command REQUEST_GUILD_MEMBERS.
// See RequestGuildMembers
//
// WARNING: If this request is in queue while a auto-scaling is forced, it will be removed from the queue
// and not re-inserted like the other commands. This is due to the guild id slice, which is a bit trickier
// to handle.
//
// Wrapper for websocket.RequestGuildMembersPayload
type RequestGuildMembersPayload = gateway.RequestGuildMembersPayload

var _ gateway.CmdPayload = (*RequestGuildMembersPayload)(nil)

// UpdateVoiceStatePayload payload for socket command UPDATE_VOICE_STATE.
// see UpdateVoiceState
//
// Wrapper for websocket.UpdateVoiceStatePayload
type UpdateVoiceStatePayload = gateway.UpdateVoiceStatePayload

var _ gateway.CmdPayload = (*UpdateVoiceStatePayload)(nil)

const (
	StatusOnline  = gateway.StatusOnline
	StatusOffline = gateway.StatusOffline
	StatusDnd     = gateway.StatusDND
	StatusIdle    = gateway.StatusIdle
)

// UpdateStatusPayload payload for socket command UPDATE_STATUS.
// see UpdateStatus
//
// Wrapper for websocket.UpdateStatusPayload
type UpdateStatusPayload = gateway.UpdateStatusPayload

var _ gateway.CmdPayload = (*UpdateStatusPayload)(nil)
