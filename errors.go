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

var MissingGuildIDErr = fmt.Errorf("guild id: %w", MissingRequiredFieldErr)
var MissingChannelIDErr = fmt.Errorf("channel id: %w", MissingRequiredFieldErr)
var MissingUserIDErr = fmt.Errorf("user id: %w", MissingRequiredFieldErr)
var MissingMessageIDErr = fmt.Errorf("message id: %w", MissingRequiredFieldErr)
var MissingEmojiIDErr = fmt.Errorf("emoji id: %w", MissingRequiredFieldErr)
var MissingRoleIDErr = fmt.Errorf("role id: %w", MissingRequiredFieldErr)
var MissingWebhookIDErr = fmt.Errorf("webhook id: %w", MissingRequiredFieldErr)
var MissingPermissionOverwriteIDErr = fmt.Errorf("channel permission overwrite id: %w", MissingRequiredFieldErr)

var MissingGuildNameErr = fmt.Errorf("guild name: %w", MissingRequiredFieldErr)
var MissingChannelNameErr = fmt.Errorf("channel name: %w", MissingRequiredFieldErr)
var MissingWebhookNameErr = fmt.Errorf("webhook name: %w", MissingRequiredFieldErr)
var MissingThreadNameErr = fmt.Errorf("thread name: %w", MissingRequiredFieldErr)

var MissingWebhookTokenErr = errors.New("webhook token was not set")
