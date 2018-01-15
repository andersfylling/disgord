package disgord

import (
	"fmt"

	"github.com/andersfylling/snowflake"
)

type User struct {
	ID            snowflake.ID `json:"id"`
	Email         string       `json:"email"`
	Username      string       `json:"username"`
	Avatar        string       `json:"avatar"`
	Discriminator string       `json:"discriminator"`
	Token         string       `json:"token"`
	Verified      bool         `json:"verified"`
	MFAEnabled    bool         `json:"mfa_enabled"`
	Bot           bool         `json:"bot"`
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
