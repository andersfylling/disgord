package guild

import "github.com/andersfylling/disgord/user"

type Ban struct {
	Reason *string    `json:"reason"`
	User   *user.User `json:"user"`
}
