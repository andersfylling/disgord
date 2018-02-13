package discord

import (
	"sync"

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
