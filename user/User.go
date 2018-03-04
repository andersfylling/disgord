package user

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/andersfylling/snowflake"
)

// UserMessager Methods required to create a new DM (or use an existing one) and send a DM.
type UserMessager interface {
	//CreateAndSendDM(recipientID snowflake.ID, msg *Message) error // hmmm...
}

type UserInterface interface {
	Mention() string
	MentionNickname() string

	// Update internal structure
	Update(*User) error
	Clear()

	// Send a direct message to this user
	SendMessage(UserMessager, string) (snowflake.ID, snowflake.ID, error)
}

type User struct {
	ID            snowflake.ID `json:"id,omitempty"`
	Username      string       `json:"username,omitempty"`
	Discriminator string       `json:"discriminator,omitempty"`
	Email         string       `json:"email,omitempty"`
	Avatar        *string      `json:"avatar"`
	Token         string       `json:"token,omitempty"`
	Verified      bool         `json:"verified,omitempty"`
	MFAEnabled    bool         `json:"mfa_enabled,omitempty"`
	Bot           bool         `json:"bot,omitempty"`

	sync.RWMutex `json:"-"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) Mention() string {
	return "<@" + u.ID.String() + ">"
}

func (u *User) MentionNickname() string {
	return "<@!" + u.ID.String() + ">"
}

func (u *User) String() string {
	return u.Username + "#" + u.Discriminator + "{" + u.ID.String() + "}"
}

func (u *User) MarshalJSON() ([]byte, error) {
	if u.ID.Empty() {
		return []byte("{}"), nil
	}

	return json.Marshal(User(*u))
}

// func (u *User) UnmarshalJSON(data []byte) error {
// 	return json.Unmarshal(data, &u.userJSON)
// }

func (u *User) Clear() {
	//u.d.Avatar = nil
}

func (u *User) SendMessage(client UserMessager, msg string) (channelID snowflake.ID, messageID snowflake.ID, err error) {
	return snowflake.NewID(0), snowflake.NewID(0), errors.New("not implemented")
}

func (u *User) Replicate(user *User) error {
	if user == nil {
		return errors.New("cannot copy nil object")
	}

	user.RLock()
	u.Lock()

	// deep copy, without changing the avatar pointer address
	var avatarAddress *string
	if u.Avatar != nil {
		avatarAddress = u.Avatar
	} else {
		avatar := ""
		avatarAddress = &avatar
	}

	*u = *user                   // copy all the fields
	u.Avatar = avatarAddress     // point to a different location than user.Avatar
	*(u.Avatar) = *(user.Avatar) // copy the Base64 image

	u.Unlock()
	user.RUnlock()

	return nil
}
