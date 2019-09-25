package disgord

import (
	"sync"

	"github.com/andersfylling/disgord/websocket/cmd"
)

// SocketCommand represents the type used to emit commands to Discord
// over the socket connection
type SocketCommand = string

// CommandRequestGuildMembers Used to request offline members for a guild or
// a list of guilds. When initially connecting, the gateway will only send
// offline members if a guild has less than the large_threshold members
// (value in the Gateway Identify). If a Client wishes to receive additional
// members, they need to explicitly request them via this operation. The
// server will send Guild Members Chunk events in response with up to 1000
// members per chunk until all members that match the request have been sent.
const CommandRequestGuildMembers SocketCommand = cmd.RequestGuildMembers

// RequestGuildMembersCommand payload for socket command REQUEST_GUILD_MEMBERS.
// See CommandRequestGuildMembers
type RequestGuildMembersCommand struct {
	// GuildID	id of the guild(s) to get offline members for
	GuildID []Snowflake `json:"guild_id"`

	// Query string that username starts with, or an empty string to return all members
	Query string `json:"query"`

	// Limit maximum number of members to send or 0 to request all members matched
	Limit uint `json:"limit"`
}

func (u *RequestGuildMembersCommand) getGuildID() Snowflake {
	return u.GuildID
}

var _ guilder = (*RequestGuildMembersCommand)(nil)

// CommandUpdateVoiceState Sent when a Client wants to join, move, or
// disconnect from a voice channel.
const CommandUpdateVoiceState SocketCommand = cmd.UpdateVoiceState

// UpdateVoiceStateCommand payload for socket command UPDATE_VOICE_STATE.
// see CommandUpdateVoiceState
type UpdateVoiceStateCommand struct {
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

func (u *UpdateVoiceStateCommand) getGuildID() Snowflake {
	return u.GuildID
}

var _ guilder = (*UpdateVoiceStateCommand)(nil)

// CommandUpdateStatus Sent by the Client to indicate a presence or status
// update.
const CommandUpdateStatus SocketCommand = cmd.UpdateStatus

// UpdateStatusCommand payload for socket command UPDATE_STATUS.
// see CommandUpdateStatus
type UpdateStatusCommand struct {
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

type Presence = UpdateStatusCommand
