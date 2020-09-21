package disgord

import (
	"errors"
	"sync"

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

type gatewayCmdPayload interface { // TODO: go generate...
	isGatewayCmdPayload() bool
}

func prepareGatewayCommand(payload gatewayCmdPayload) (x gateway.CmdPayload, err error) {
	switch t := payload.(type) {
	case *RequestGuildMembersPayload:
		x = &gateway.RequestGuildMembersPayload{
			GuildIDs: t.GuildIDs,
			Query:    t.Query,
			Limit:    t.Limit,
			UserIDs:  t.UserIDs,
		}
	case *UpdateVoiceStatePayload:
		x = &gateway.UpdateVoiceStatePayload{
			GuildID:   t.GuildID,
			ChannelID: t.ChannelID,
			SelfMute:  t.SelfMute,
			SelfDeaf:  t.SelfDeaf,
		}
	case *UpdateStatusPayload:
		status, err := gateway.StringToStatusType(t.Status)
		if err != nil {
			return nil, err
		}
		x = &gateway.UpdateStatusPayload{
			Since:  t.Since,
			Game:   t.Game,
			Status: status,
			AFK:    t.AFK,
		}
	default:
		return nil, errors.New("missing support for payload")
	}

	return x, nil
}

// #################################################################
// RequestGuildMembersPayload payload for socket command REQUEST_GUILD_MEMBERS.
// See RequestGuildMembers
//
// WARNING: If this request is in queue while a auto-scaling is forced, it will be removed from the queue
// and not re-inserted like the other commands. This is due to the guild id slice, which is a bit trickier
// to handle.
//
// Wrapper for websocket.RequestGuildMembersPayload
type RequestGuildMembersPayload struct {
	// GuildID	id of the guild(s) to get the members for
	GuildIDs []Snowflake

	// Query string that username starts with, or an empty string to return all members
	Query string

	// Limit maximum number of members to send or 0 to request all members matched
	Limit uint

	// UserIDs used to specify which Users you wish to fetch
	UserIDs []Snowflake
}

var _ gatewayCmdPayload = (*RequestGuildMembersPayload)(nil)

func (r *RequestGuildMembersPayload) isGatewayCmdPayload() bool { return true }

// #################################################################

// UpdateVoiceStatePayload payload for socket command UPDATE_VOICE_STATE.
// see UpdateVoiceState
//
// Wrapper for websocket.UpdateVoiceStatePayload
type UpdateVoiceStatePayload struct {
	// GuildID id of the guild
	GuildID Snowflake

	// ChannelID id of the voice channel Client wants to join
	// (0 if disconnecting)
	ChannelID Snowflake

	// SelfMute is the Client mute
	SelfMute bool

	// SelfDeaf is the Client deafened
	// Currently, it's n ot documented how to receive sound from discord
	// therefore it's not supported and the bot is always going to be deaf.
	//
	// deprecated, this is always true
	SelfDeaf bool
}

var _ gatewayCmdPayload = (*UpdateVoiceStatePayload)(nil)

func (u *UpdateVoiceStatePayload) isGatewayCmdPayload() bool { return true }

// #################################################################

const (
	StatusOnline  = string(gateway.StatusOnline)
	StatusOffline = string(gateway.StatusOffline)
	StatusDnd     = string(gateway.StatusDND)
	StatusIdle    = string(gateway.StatusIdle)
)

// UpdateStatusPayload payload for socket command UPDATE_STATUS.
// see UpdateStatus
//
// Wrapper for websocket.UpdateStatusPayload
type UpdateStatusPayload struct {
	mu sync.RWMutex
	// Since unix time (in milliseconds) of when the Client went idle, or null if the Client is not idle
	Since *uint

	// Game null, or the user's new activity
	Game *Activity

	// Status the user's new status
	Status string

	// AFK whether or not the Client is afk
	AFK bool
}

var _ gatewayCmdPayload = (*UpdateStatusPayload)(nil)

func (u *UpdateStatusPayload) isGatewayCmdPayload() bool { return true }
