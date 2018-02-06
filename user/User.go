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
	ID() snowflake.ID
	Username() string
	Discriminator() string
	Email() string
	Avatar() string
	Token() string
	Verified() bool
	MFAEnabled() bool
	Bot() bool

	// Update internal structure
	Update(UserInterface) error
	Clear()

	// Send a direct message to this user
	SendMessage(DisgordDMInterface, string) (snowflake.ID, snowflake.ID, error)
}

type userJSON struct {
	ID            snowflake.ID `json:"id,omitempty"`
	Username      string       `json:"username,omitempty"`
	Discriminator string       `json:"discriminator,omitempty"`
	Email         string       `json:"email,omitempty"`
	Avatar        string       `json:"avatar"`
	Token         string       `json:"token,omitempty"`
	Verified      bool         `json:"verified,omitempty"`
	MFAEnabled    bool         `json:"mfa_enabled,omitempty"`
	Bot           bool         `json:"bot,omitempty"`
}

type User struct {
	d userJSON
	sync.RWMutex
}

func NewUser() *User {
	return &User{}
}

func (u *User) Mention() string {
	return "<@" + u.d.ID.String() + ">"
}

func (u *User) MentionNickname() string {
	return "<@!" + u.d.ID.String() + ">"
}

func (u *User) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &u.d)
}

func (u *User) MarshalJSON() ([]byte, error) {
	if u.d.ID.Empty() {
		return []byte("{}"), nil
	}

	return json.Marshal(&u.d)
}

func (u *User) Clear() {
	//u.d.Avatar = nil
}

func (u *User) EqualID(ID snowflake.ID) bool {
	u.Lock()
	defer u.Unlock()
	return ID == u.d.ID
}

func (u *User) ID() snowflake.ID {
	u.Lock()
	defer u.Unlock()
	return u.d.ID
}
func (u *User) Username() string {
	u.Lock()
	defer u.Unlock()
	return u.d.Username
}
func (u *User) Discriminator() string {
	u.Lock()
	defer u.Unlock()
	return u.d.Discriminator
}
func (u *User) Email() string {
	u.Lock()
	defer u.Unlock()
	return u.d.Email
}
func (u *User) Avatar() string {
	u.Lock()
	defer u.Unlock()
	return u.d.Avatar
}
func (u *User) Token() string {
	u.Lock()
	defer u.Unlock()
	return u.d.Token
}
func (u *User) Verified() bool {
	u.Lock()
	defer u.Unlock()
	return u.d.Verified
}
func (u *User) MFAEnabled() bool {
	u.Lock()
	defer u.Unlock()
	return u.d.MFAEnabled
}
func (u *User) Bot() bool {
	u.Lock()
	defer u.Unlock()
	return u.d.Bot
}

func (u *User) Update(new UserInterface) (err error) {
	if !u.EqualID(new.ID()) {
		err = errors.New("cannot update user when the new struct has a different ID")
		return
	}
	// make sure that new is not the same pointer!
	if u == new.(*User) {
		err = errors.New("cannot update user when the new struct points to the same memory space")
		return
	}

	u.Lock()
	u.d.Username = new.Username()
	u.d.Discriminator = new.Discriminator()
	u.d.Email = new.Email()
	u.d.Avatar = new.Avatar()
	u.d.Token = new.Token()
	u.d.Verified = new.Verified()
	u.d.Verified = new.Verified()
	u.d.MFAEnabled = new.MFAEnabled()
	u.d.Bot = new.Bot()
	u.Unlock()

	return
}

func (u *User) SendMessage(client DisgordDMInterface, msg string) (channelID snowflake.ID, messageID snowflake.ID, err error) {
	return snowflake.NewID(0), snowflake.NewID(0), errors.New("not implemented")
}
