package disgord

// Emoji
type Emoji struct {
	ID            Snowflake   `json:"id"`
	Name          string      `json:"name"`
	Roles         []Snowflake `json:"roles,omitempty"`
	User          *User       `json:"user,omitempty"` // the user who created the emoji
	RequireColons bool        `json:"require_colons,omitempty"`
	Managed       bool        `json:"managed,omitempty"`
	Animated      bool        `json:"animated,omitempty"`
}
type PartialEmoji = Emoji

// Mention
// TODO: review
func (e *Emoji) Mention() string {
	return "<" + e.Name + ":" + e.ID.String() + ">"
}

// MentionAnimated add the animation prefix if a animated emoji
// TODO: review
func (e *Emoji) MentionAnimated() string {
	prefix := ""
	if e.Animated {
		prefix = "a:"
	}

	return "<" + prefix + e.Name + ":" + e.ID.String() + ">"
}

func (e *Emoji) Clear() {
	// obviously don't delete the user ...
}
