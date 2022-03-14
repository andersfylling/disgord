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
var ErrMissingScheduledEventName = fmt.Errorf("scheduled event name: %w", ErrMissingName)

var ErrMissingWebhookToken = errors.New("webhook token was not set")

var ErrIllegalValue = errors.New("illegal value")
var ErrIllegalScheduledEventPrivacyLevelValue = fmt.Errorf("scheduled event privacy level: %w", ErrIllegalValue)

var ErrMissingTime = fmt.Errorf("time: %w", ErrMissingRequiredField)
var ErrMissingScheduledEventStartTime = fmt.Errorf("scheduled event start: %w", ErrMissingTime)
var ErrMissingScheduledEventEndTime = fmt.Errorf("scheduled event end: %w", ErrMissingTime)

var ErrMissingScheduledEventLocation = fmt.Errorf("scheduled event: %w", ErrMissingRequiredField)

var ErrMissingType = fmt.Errorf("type: %w", ErrMissingRequiredField)
var ErrMissingScheduledEventEntityType = fmt.Errorf("scheduled event entity: %w", ErrMissingType)
