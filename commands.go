package disgord

// SocketCommand represents the type used to emit commands to Discord
// over the socket connection
type SocketCommand = string

// CommandRequestGuildMembers Used to request offline members for a guild or
// a list of guilds. When initially connecting, the gateway will only send
// offline members if a guild has less than the large_threshold members
// (value in the Gateway Identify). If a client wishes to receive additional
// members, they need to explicitly request them via this operation. The
// server will send Guild Members Chunk events in response with up to 1000
// members per chunk until all members that match the request have been sent.
const CommandRequestGuildMembers SocketCommand = "REQUEST_GUILD_MEMBERS"

// RequestGuildMembersCommand payload for socket command REQUEST_GUILD_MEMBERS.
// See CommandRequestGuildMembers
type RequestGuildMembersCommand struct {
	// GuildID	id of the guild(s) to get offline members for
	GuildID Snowflake `json:"guild_id"`

	// Query string that username starts with, or an empty string to return all members
	Query string `json:"query"`

	// Limit maximum number of members to send or 0 to request all members matched
	Limit uint `json:"limit"`
}

// CommandUpdateVoiceState Sent when a client wants to join, move, or
// disconnect from a voice channel.
const CommandUpdateVoiceState SocketCommand = "UPDATE_VOICE_STATE"

// UpdateVoiceStateCommand payload for socket command UPDATE_VOICE_STATE.
// see CommandUpdateVoiceState
type UpdateVoiceStateCommand struct {
	// GuildID id of the guild
	GuildID Snowflake `json:"guild_id"`

	// ChannelID id of the voice channel client wants to join
	// (null if disconnecting)
	ChannelID *Snowflake `json:"channel_id"`

	// SelfMute is the client mute
	SelfMute bool `json:"self_mute"`

	// SelfDeaf is the client deafened
	SelfDeaf bool `json:"self_deaf"`
}

// CommandUpdateStatus Sent by the client to indicate a presence or status
// update.
const CommandUpdateStatus SocketCommand = "UPDATE_STATUS"

// UpdateStatusCommand payload for socket command UPDATE_STATUS.
// see CommandUpdateStatus
type UpdateStatusCommand struct {
	// Since unix time (in milliseconds) of when the client went idle, or null if the client is not idle
	Since *uint `json:"since"`

	// Game null, or the user's new activity
	Game *Activity `json:"activity"`

	// Status the user's new status
	Status string `json:"status"`

	// AFK whether or not the client is afk
	AFK bool `json:"afk"`
}
