package user

import (
	"encoding/json"
	"fmt"

	"github.com/andersfylling/snowflake"
)

type User struct {
	ID            snowflake.ID `json:"id,omitempty,string"`
	Username      string       `json:"username,omitempty"`
	Discriminator string       `json:"discriminator,omitempty"`
	Email         string       `json:"email,omitempty"`
	Avatar        *string      `json:"avatar"`
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

func (u *User) MarshalJSON() ([]byte, error) {
	if u.ID.Empty() {
		return []byte("{}"), nil
	}

	// use an alias to avoid stack overflow by recursion
	type Alias User
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	})
}
