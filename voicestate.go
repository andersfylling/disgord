package disgord

import "github.com/andersfylling/snowflake"

type VoiceState struct {
	UserID    snowflake.ID `json:"user_id"`
	SessionID snowflake.ID `json:"session_id"`
	ChannelID snowflake.ID `json:"channel_id"`
	GuildID   snowflake.ID `json:"guild_id"`
	Suppress  bool         `json:"suppress"`
	SelfMute  bool         `json:"self_mute"`
	SelfDeaf  bool         `json:"self_deaf"`
	Mute      bool         `json:"mute"`
	Deaf      bool         `json:"deaf"`
}
