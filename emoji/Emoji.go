package emoji

import (
	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/snowflake"
)

type Emoji struct {
	ID            snowflake.ID   `json:"id"`
	Name          string         `json:"name"`
	Roles         []snowflake.ID `json:"roles,omitempty"`
	User          *user.User     `json:"user,omitempty"` // the user who created the emoji
	RequireColons bool           `json:"require_colons,omitempty"`
	Managed       bool           `json:"managed,omitempty"`
	Animated      bool           `json:"animated,omitempty"`
}

func (e *Emoji) Mention() string {
	animated := ""
	if e.Animated {
		animated = "a:"
	}

	return "<" + animated + e.Name + ":" + e.ID.String() + ">"
}

func (e *Emoji) Clear() {
	// obviusly don't delete the user ...
}
