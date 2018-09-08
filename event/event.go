// Package description
//
package event

import (
	"context"
	"sync"

	"github.com/andersfylling/disgord/resource"
)

// event keys that does not fit within one of the existing files goes here
const KeyAllEvents = "GOD_DAMN_EVERYTHING"

// Gateway events

// Hello defines the heartbeat interval
type Hello struct {
	HeartbeatInterval uint            `json:"heartbeat_interval"`
	Trace             []string        `json:"_trace"`
	Ctx               context.Context `json:"-"`
}

// KeyReady The ready event is dispatched when a client has completed the
//       initial handshake with the gateway (for new sessions). The ready
//       event can be the largest and most complex event the gateway will
//       send, as it contains all the state required for a client to begin
//       interacting with the rest of the platform.
//       Fields:
//       * V uint8
//       * User *discord.user.User
//       * PrivateChannels []*discord.channel.Private
//       * Guilds []*discord.guild.Unavailable
//       * SessionID string
//       * Trace []string
const KeyReady = "READY"

// Ready	contains the initial state information
type Ready struct {
	APIVersion int                          `json:"v"`
	User       *resource.User               `json:"user"`
	Guilds     []*resource.GuildUnavailable `json:"guilds"`

	// not really needed, as it is handled on the socket layer.
	SessionID string   `json:"session_id"`
	Trace     []string `json:"_trace"`

	// private_channels will be an empty array. As bots receive private messages,
	// they will be notified via Channel Create events.
	//PrivateChannels []*channel.Channel `json:"private_channels"`

	// bot can't have presences
	//Presences []*Presence         `json:"presences"`

	// bot cant have relationships
	//Relationships []interface{} `son:"relationships"`

	// bot can't have user settings
	// UserSettings interface{}        `json:"user_settings"`

	sync.RWMutex `json:"-"`
	Ctx          context.Context `json:"-"`
}

// KeyResumed The resumed event is dispatched when a client has sent a resume
//         payload to the gateway (for resuming existing sessions).
//         Fields:
//         * Trace []string
const KeyResumed = "RESUMED"

// Resumed	response to Resume
type Resumed struct {
	Trace []string        `json:"_trace"`
	Ctx   context.Context `json:"-"`
}

// InvalidSession	failure response to Identify or Resume or invalid active session
type InvalidSession struct {
	Ctx context.Context `json:"-"`
}
