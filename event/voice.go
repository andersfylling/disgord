package event

import (
	"context"

	"github.com/andersfylling/disgord/resource"
	. "github.com/andersfylling/snowflake"
)

// KeyVoiceStateUpdate Sent when someone joins/leaves/moves voice channels.
//                  Inner payload is a voice state object.
const KeyVoiceStateUpdate = "VOICE_STATE_UPDATE"

// VoiceStateUpdate	someone joined, left, or moved a voice channel
type VoiceStateUpdate struct {
	VoiceState *resource.VoiceState `json:"voice_state"`
	Ctx        context.Context      `json:"-"`
}

// KeyVoiceServerUpdate Sent when a guild's voice server is updated. This is
//                   sent when initially connecting to voice, and when the
//                   current voice instance fails over to a new server.
//                   Fields:
//                   * Token     string or discord.Token
//                   * ChannelID int64 or discord.Snowflake
//                   * Endpoint  string or discord.Endpoint
const KeyVoiceServerUpdate = "VOICE_SERVER_UPDATE"

// VoiceServerUpdate	guild's voice server was updated
// Sent when a guild's voice server is updated.
// This is sent when initially connecting to voice,
// and when the current voice instance fails over to a new server.
type VoiceServerUpdate struct {
	Token    string          `json:"token"`
	GuildID  Snowflake       `json:"guild_id"`
	Endpoint string          `json:"endpoint"`
	Ctx      context.Context `json:"-"`
}
