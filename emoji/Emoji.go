package emoji

import (
	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/snowflake"
)

type Emoji struct {
	ID            snowflake.ID   `json:"id"`
	Name          string         `json:"name"`
	User          *user.User     `json:"user"` // the user who created the emoji
	Roles         []snowflake.ID `json:"roles"`
	RequireColons bool           `json:"require_colons"`
	Managed       bool           `json:"managed"`
}

func (e *Emoji) Clear() {
	// obviusly don't delete the user ...
}
