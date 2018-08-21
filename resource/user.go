package resource

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/andersfylling/disgord/rest/httd"
	"github.com/andersfylling/snowflake"
)

type UserInterface interface {
	Mention() string
	MentionNickname() string
	String() string
}

// TODO: https://discordapp.com/developers/docs/topics/gateway#activity-object-activity-structure
type UserActivity struct{}

// ---------

// TODO: should a user object always have a ID?
func NewUser() *User {
	return &User{}
}

type User struct {
	ID            snowflake.ID `json:"id,omitempty"`
	Username      string       `json:"username,omitempty"`
	Discriminator string       `json:"discriminator,omitempty"`
	Email         string       `json:"email,omitempty"`
	Avatar        *string      `json:"avatar"` // data:image/jpeg;base64,BASE64_ENCODED_JPEG_IMAGE_DATA
	Token         string       `json:"token,omitempty"`
	Verified      bool         `json:"verified,omitempty"`
	MFAEnabled    bool         `json:"mfa_enabled,omitempty"`
	Bot           bool         `json:"bot,omitempty"`

	sync.RWMutex `json:"-"`
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

// Partial check if this is not a complete user object
// Assumption: has a snowflake.
func (u *User) Partial() bool {
	return (u.Username + u.Discriminator) == ""
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

func (u *User) SendMsg(requester httd.Requester, msg *Message) (channelID snowflake.ID, messageID snowflake.ID, err error) {
	return snowflake.NewID(0), snowflake.NewID(0), errors.New("not implemented")
}

func (u *User) SendMsgString(requester httd.Requester, msg string) (channelID snowflake.ID, messageID snowflake.ID, err error) {
	return snowflake.NewID(0), snowflake.NewID(0), errors.New("not implemented")
}

func (u *User) DeepCopy() *User {
	user := NewUser()

	u.RLock()

	user.ID = u.ID
	user.Username = u.Username
	user.Discriminator = u.Discriminator
	user.Email = u.Email
	user.Token = u.Token
	user.Verified = u.Verified
	user.MFAEnabled = u.MFAEnabled
	user.Bot = u.Bot

	if u.Avatar != nil {
		avatar := *u.Avatar
		user.Avatar = &avatar
	}

	u.RUnlock()

	return user
}

func (u *User) Valid() bool {
	return u.ID > 0
}

// -------

type UserPresence struct {
	User    *User          `json:"user"`
	Roles   []snowflake.ID `json:"roles"`
	Game    *UserActivity  `json:"activity"`
	GuildID snowflake.ID   `json:"guild_id"`
	Nick    string         `json:"nick"`
	Status  string         `json:"status"`
}

func NewUserPresence() *UserPresence {
	return &UserPresence{}
}

func (p *UserPresence) Update(status string) {
	// Update the presence.
	// talk to the discord api
}

func (p *UserPresence) String() string {
	return p.Status
}

func (p *UserPresence) Clear() {
	p.Game = nil
}

type UserConnection struct {
	ID           string                `json:"id"`           // id of the connection account
	Name         string                `json:"name"`         // the username of the connection account
	Type         string                `json:"type"`         // the service of the connection (twitch, youtube)
	Revoked      bool                  `json:"revoked"`      // whether the connection is revoked
	Integrations []*IntegrationAccount `json:"integrations"` // an array of partial server integrations
}
