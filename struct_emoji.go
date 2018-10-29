package disgord

import (
	"github.com/andersfylling/disgord/constant"
)

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
	mu Lockable
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

	return "<" + prefix + e.Name + ":" + e.ID.String() + ">"
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

	if constant.LockedMethods {
		e.mu.RLock()
		emoji.mu.Lock()
	}

	emoji.ID = e.ID
	emoji.Name = e.Name
	emoji.Roles = e.Roles
	emoji.RequireColons = e.RequireColons
	emoji.Managed = e.Managed
	emoji.Animated = e.Animated
	emoji.mu = Lockable{}

	if e.User != nil {
		emoji.User = e.User.DeepCopy().(*User)
	}

	if constant.LockedMethods {
		e.mu.RUnlock()
		emoji.mu.Unlock()
	}

	return
}

// Missing GuildID...
//func (e *Emoji) saveToDiscord(session Session) (err error) {
//	session.Emoji
//}
//func (e *Emoji) deleteFromDiscord(session Session) (err error) {
//	session.DeleteGuildEmoji(guildID, emojiID)
//}

//func (e *Emoji) createGuildEmoji(session Session) (err error) {
//	params := &CreateGuildEmojiParams{
//		Name:  e.Name,
//		Image: e.image,
//		Roles: e.Roles,
//	}
//
//	var creation *Emoji
//	creation, err = session.CreateGuildEmoji(, params)
//	if err != nil {
//		return
//	}
//
//	creation.CopyOverTo(e)
//	return
//}
//func (e *Emoji) modifyGuildEmoji(session Session) (err error) {
//	params := &ModifyGuildEmojiParams{
//		Name:  e.Name,
//		Roles: e.Roles,
//	}
//}

// func (e *Emoji) Clear() {
// 	// obviously don't delete the user ...
// }
