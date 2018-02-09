package discord

import (
	"sync"

	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/snowflake"
)

func NewGuildUnavailable(ID snowflake.ID) *GuildUnavailable {
	gu := &GuildUnavailable{
		ID:          ID,
		Unavailable: true,
	}

	return gu
}

type GuildUnavailable struct {
	ID           snowflake.ID `json:"id"`
	Unavailable  bool         `json:"unavailable"` // ?*|
	sync.RWMutex `json:"-"`
}

// func (gu *GuildUnavailable) UnmarshalJSON(data []byte) error {
// 	return json.Unmarshal(data, GuildUnavailable(*gu))
// }
// func (gu *GuildUnavailable) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(GuildUnavailable(*gu))
// }

type Ready struct {
	APIVersion int                 `json:"v"`
	User       *user.User          `json:"user"`
	Guilds     []*GuildUnavailable `json:"guilds"`

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
}

// func (r *Ready) UnmarshalJSON(data []byte) error {
// 	return json.Unmarshal(data, Ready(*r))
// }
// func (r *Ready) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(Ready(*r))
// }

type Resumed struct {
	Trace []string `json:"_trace"`
}
