package test

import (
	"os"

	"github.com/andersfylling/disgord"
)

var token = os.Getenv("DISGORD_TOKEN_INTEGRATION_TEST")

var guildTypical = struct {
	ID                  disgord.Snowflake
	VoiceChannelGeneral disgord.Snowflake
	VoiceChannelOther1  disgord.Snowflake
	VoiceChannelOther2  disgord.Snowflake
}{
	ID:                  486833611564253184,
	VoiceChannelGeneral: 486833611564253188,
	VoiceChannelOther1:  673893473409171477,
	VoiceChannelOther2:  673893496356339724,
}
