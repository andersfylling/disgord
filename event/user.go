package event

import (
	"context"

	"github.com/andersfylling/disgord/resource"
)

// KeyUserUpdate Sent when properties about the user change. Inner payload is a
//            user object.
const KeyUserUpdate = "USER_UPDATE"

// UserUpdate	properties about a user changed
type UserUpdate struct {
	User *resource.User  `json:"user"`
	Ctx  context.Context `json:"-"`
}
