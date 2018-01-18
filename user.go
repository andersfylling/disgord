package disgord

import (
	"fmt"

	"github.com/andersfylling/snowflake"
)

type User struct {
	Username      string       `json:"username,omitempty"`
	Discriminator string       `json:"discriminator,omitempty"`
	ID            snowflake.ID `json:"id,string,omitempty"`
	Email         string       `json:"email,omitempty"`
	Avatar        string       `json:"avatar,omitempty"`
	Token         string       `json:"token,omitempty"`
	Verified      bool         `json:"verified,omitempty"`
	MFAEnabled    bool         `json:"mfa_enabled,omitempty"`
	Bot           bool         `json:"bot,omitempty"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) Mention() string {
	return fmt.Sprintf("<@%d>", u.ID)
}

func (u *User) MentionNickname() string {
	return fmt.Sprintf("<@!%d>", u.ID)
}
