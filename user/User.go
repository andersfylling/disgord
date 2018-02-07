package user

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/andersfylling/snowflake"
)

// DisgordDMInterface Methods required to create a new DM (or use an existing one) and send a DM.
type DisgordDMInterface interface {
}

type UserInterface interface {
	Mention() string
	MentionNickname() string
	String() string
	UnmarshalJSON([]byte) error
	MarshalJSON() ([]byte, error)

	// Update internal structure
	Update(*User) error
	Clear()

	// Send a direct message to this user
	SendMessage(DisgordDMInterface, string) (snowflake.ID, snowflake.ID, error)
}

type userJSON struct {
	ID            snowflake.ID `json:"id,omitempty"`
	Username      string       `json:"username,omitempty"`
	Discriminator string       `json:"discriminator,omitempty"`
	Email         string       `json:"email,omitempty"`
	Avatar        *string      `json:"avatar"`
	Token         string       `json:"token,omitempty"`
	Verified      bool         `json:"verified,omitempty"`
	MFAEnabled    bool         `json:"mfa_enabled,omitempty"`
	Bot           bool         `json:"bot,omitempty"`
}

type User struct {
	userJSON // simplifies marshalling and `userJSON` doesn't appear, but exported fields can still be accessed
	sync.RWMutex
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
	return u.Username + "#" + u.Discriminator
}

func (u *User) MarshalJSON() ([]byte, error) {
	if u.ID.Empty() {
		return []byte("{}"), nil
	}

	return json.Marshal(&u.userJSON)
}

func (u *User) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &u.userJSON)
}

func (u *User) Clear() {
	//u.d.Avatar = nil
}

func (u *User) Update(new *User) (err error) {
	if u.ID != new.ID {
		err = errors.New("cannot update user when the new struct has a different ID")
		return
	}
	// make sure that new is not the same pointer!
	if u == new {
		err = errors.New("cannot update user when the new struct points to the same memory space")
		return
	}

	u.Lock()
	new.Lock()
	u.Username = new.Username
	u.Discriminator = new.Discriminator
	u.Email = new.Email
	u.Avatar = new.Avatar
	u.Token = new.Token
	u.Verified = new.Verified
	u.MFAEnabled = new.MFAEnabled
	u.Bot = new.Bot
	new.Unlock()
	u.Unlock()

	return
}

func (u *User) SendMessage(client DisgordDMInterface, msg string) (channelID snowflake.ID, messageID snowflake.ID, err error) {
	return snowflake.NewID(0), snowflake.NewID(0), errors.New("not implemented")
}
