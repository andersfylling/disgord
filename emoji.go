package disgord

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

func validEmojiName(name string) bool {
	if name == "" {
		return false
	}
	// TODO: what is the allowed format?
	// a test showed that using "-" caused regex issues
	return !strings.Contains(name, "-")
}

// Emoji ...
type Emoji struct {
	ID            Snowflake   `json:"id"`
	Name          string      `json:"name"`
	Roles         []Snowflake `json:"roles,omitempty"`
	User          *User       `json:"user,omitempty"` // the user who created the emoji
	RequireColons bool        `json:"require_colons,omitempty"`
	Managed       bool        `json:"managed,omitempty"`
	Animated      bool        `json:"animated,omitempty"`

	//	image string // base 64 string, with prefix and everything
	guildID Snowflake
}

var _ Reseter = (*Emoji)(nil)
var _ DeepCopier = (*Emoji)(nil)
var _ Copier = (*Emoji)(nil)
var _ discordDeleter = (*Emoji)(nil)
var _ Mentioner = (*Emoji)(nil)

// var _ discordSaver = (*Emoji)(nil) // TODO
var _ fmt.Stringer = (*Emoji)(nil)

func (e *Emoji) String() string {
	return "emoji{name:" + e.Name + ", id:" + e.ID.String() + "}"
}

// PartialEmoji see Emoji
type PartialEmoji = Emoji

// SetBase64Image use this before creating the emoji for the first time
//func (e *Emoji) SetBase64Image(img string) {
//	e.image = img
//}

// Mention mentions an emoji. Adds the animation prefix, if animated
func (e *Emoji) Mention() string {
	prefix := ""
	if e.Animated {
		prefix = "a:"
	}

	return "<:" + prefix + e.Name + ":" + e.ID.String() + ">"
}

func (e *Emoji) LinkToGuild(guildID Snowflake) {
	e.guildID = guildID
}

// DeepCopy see interface at struct.go#DeepCopier
func (e *Emoji) DeepCopy() (copy interface{}) {
	copy = &Emoji{}
	e.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (e *Emoji) CopyOverTo(other interface{}) (err error) {
	var emoji *Emoji
	var ok bool
	if emoji, ok = other.(*Emoji); !ok {
		err = newErrorUnsupportedType("given type is not *Emoji")
		return
	}

	emoji.ID = e.ID
	emoji.Name = e.Name
	emoji.Roles = e.Roles
	emoji.RequireColons = e.RequireColons
	emoji.Managed = e.Managed
	emoji.Animated = e.Animated
	emoji.guildID = e.guildID

	if e.User != nil {
		emoji.User = e.User.DeepCopy().(*User)
	}
	return
}

// Missing GuildID...
//func (e *Emoji) saveToDiscord(s Session) (err error) {
//	session.Emoji
//}
func (e *Emoji) deleteFromDiscord(ctx context.Context, s Session, flags ...Flag) (err error) {
	if e.guildID.IsZero() {
		err = errors.New("missing guild ID, call Emoji.LinkToGuild")
		return
	}
	if e.ID.IsZero() {
		err = errors.New("missing emoji ID, cannot delete a not identified emoji")
		return
	}

	return s.Guild(e.guildID).DeleteEmoji(ctx, e.ID, flags...)
}

//////////////////////////////////////////////////////
//
// REST Methods
//
// https://discord.com/developers/docs/resources/emoji#emoji-resource
// Routes for controlling emojis do not follow the normal rate limit conventions.
// These routes are specifically limited on a per-guild basis to prevent abuse.
// This means that the quota returned by our APIs may be inaccurate,
// and you may encounter 429s.
//
//////////////////////////////////////////////////////

// CreateGuildEmojiParams JSON params for func CreateGuildEmoji
type CreateGuildEmojiParams struct {
	Name  string      `json:"name"`  // required
	Image string      `json:"image"` // required
	Roles []Snowflake `json:"roles"` // optional

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string `json:"-"`
}

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

//generate-rest-params: name:string, roles:[]Snowflake,
//generate-rest-basic-execute: emoji:*Emoji,
type updateGuildEmojiBuilder struct {
	r RESTBuilder
}

//generate-rest-params: roles:[]Snowflake,
//generate-rest-basic-execute: emoji:*Emoji,
type createGuildEmojiBuilder struct {
	r RESTBuilder
}
