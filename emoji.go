package disgord

import "github.com/andersfylling/snowflake"

type Emoji struct {
	ID            snowflake.ID   `json:"id,string"`
	Name          string         `json:"name"`
	User          *User          `json:"user"` // the user who created the emoji
	Roles         []snowflake.ID `json:"roles,string"`
	RequireColons bool           `json:"require_colons"`
	Managed       bool           `json:"managed"`
}
