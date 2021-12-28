package disgord

import (
	"errors"
	"fmt"
	"github.com/andersfylling/disgord/internal/disgorderr"
)

// TODO: go generate from internal/errors/*

type Err = disgorderr.Err
type CloseConnectionErr = disgorderr.ClosedConnectionErr
type HandlerSpecErr = disgorderr.HandlerSpecErr

//////////////////////////////////////////////////////
//
// REST errors
//
//////////////////////////////////////////////////////

var MissingRequiredFieldErr = errors.New("missing required field")

var MissingIDErr = fmt.Errorf("id: %w", MissingRequiredFieldErr)
var MissingGuildIDErr = fmt.Errorf("guild: %w", MissingIDErr)
var MissingChannelIDErr = fmt.Errorf("channel: %w", MissingIDErr)
var MissingUserIDErr = fmt.Errorf("user: %w", MissingIDErr)
var MissingMessageIDErr = fmt.Errorf("message: %w", MissingIDErr)
var MissingEmojiIDErr = fmt.Errorf("emoji: %w", MissingIDErr)
var MissingRoleIDErr = fmt.Errorf("role: %w", MissingIDErr)
var MissingWebhookIDErr = fmt.Errorf("webhook: %w", MissingIDErr)
var MissingPermissionOverwriteIDErr = fmt.Errorf("channel permission overwrite: %w", MissingIDErr)

var MissingNameErr = fmt.Errorf("name: %w", MissingRequiredFieldErr)
var MissingGuildNameErr = fmt.Errorf("guild: %w", MissingNameErr)
var MissingChannelNameErr = fmt.Errorf("channel: %w", MissingNameErr)
var MissingWebhookNameErr = fmt.Errorf("webhook: %w", MissingNameErr)
var MissingThreadNameErr = fmt.Errorf("thread: %w", MissingNameErr)

var MissingWebhookTokenErr = errors.New("webhook token was not set")
