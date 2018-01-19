package guild

import "github.com/andersfylling/snowflake"

type Embed struct {
	Enabled   bool         `json:"enabled"`
	ChannelID snowflake.ID `json:"channel_id,string"`
}
