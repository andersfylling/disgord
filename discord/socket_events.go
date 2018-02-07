package discord

import (
	"encoding/json"
	"sync"

	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/snowflake"
)

func NewGuildUnavailable(ID snowflake.ID) *GuildUnavailable {
	gu := &GuildUnavailable{}
	gu.ID = ID
	gu.Unavailable = true

	return gu
}

type UserSettings struct{}

type guildUnavailableJSON struct {
	ID          snowflake.ID `json:"id"`
	Unavailable bool         `json:"unavailable"` // ?*|
}

type GuildUnavailable struct {
	guildUnavailableJSON
	sync.RWMutex
}

func (gu *GuildUnavailable) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &gu.guildUnavailableJSON)
}
func (gu *GuildUnavailable) MarshalJSON() ([]byte, error) {
	return json.Marshal(&gu.guildUnavailableJSON)
}

type readyJSON struct {
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
	// UserSettings UserSettings        `json:"user_settings"`

}

type Ready struct {
	readyJSON
	sync.RWMutex
}

func (r *Ready) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.readyJSON)
}
func (r *Ready) MarshalJSON() ([]byte, error) {
	return json.Marshal(&r.readyJSON)
}

type Resumed struct{}
