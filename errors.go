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

var ErrMissingRequiredField = errors.New("missing required field")

var ErrMissingID = fmt.Errorf("id: %w", ErrMissingRequiredField)
var ErrMissingGuildID = fmt.Errorf("guild: %w", ErrMissingID)
var ErrMissingChannelID = fmt.Errorf("channel: %w", ErrMissingID)
var ErrMissingUserID = fmt.Errorf("user: %w", ErrMissingID)
var ErrMissingMessageID = fmt.Errorf("message: %w", ErrMissingID)
var ErrMissingEmojiID = fmt.Errorf("emoji: %w", ErrMissingID)
var ErrMissingRoleID = fmt.Errorf("role: %w", ErrMissingID)
var ErrMissingWebhookID = fmt.Errorf("webhook: %w", ErrMissingID)
var ErrMissingPermissionOverwriteID = fmt.Errorf("channel permission overwrite: %w", ErrMissingID)

var ErrMissingName = fmt.Errorf("name: %w", ErrMissingRequiredField)
var ErrMissingGuildName = fmt.Errorf("guild: %w", ErrMissingName)
var ErrMissingChannelName = fmt.Errorf("channel: %w", ErrMissingName)
var ErrMissingWebhookName = fmt.Errorf("webhook: %w", ErrMissingName)
var ErrMissingThreadName = fmt.Errorf("thread: %w", ErrMissingName)

var ErrMissingWebhookToken = errors.New("webhook token was not set")

var ErrIllegalValue = errors.New("illegal value")
