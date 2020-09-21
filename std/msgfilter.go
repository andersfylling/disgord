package std

import (
	"context"
	"strings"

	"github.com/andersfylling/disgord"
)

func NewMsgFilter(ctx context.Context, client disgord.Session) (filter *msgFilter, err error) {
	if filter, err = newMsgFilter(ctx, client); err != nil {
		return nil, err
	}
	filter.s = client

	return filter, nil
}

type CurrentUserRESTResource interface {
	CurrentUser() disgord.CurrentUserQueryBuilder
}

func newMsgFilter(ctx context.Context, client CurrentUserRESTResource) (filter *msgFilter, err error) {
	usr, err := client.CurrentUser().WithContext(ctx).Get()
	if err != nil {
		return nil, err
	}

	filter = &msgFilter{}
	filter.botID = usr.ID

	return filter, nil
}

type msgFilter struct {
	s      disgord.Session
	botID  disgord.Snowflake
	prefix string

	permissions       disgord.PermissionBit
	eitherPermissions disgord.PermissionBit
}

// SetPrefix set the prefix attribute which is used in StripPrefix, HasPrefix.
// Do not set the prefix to be a space, " ", as this will be trimmed when checking later.
func (f *msgFilter) SetPrefix(p string) {
	f.prefix = p
}

func (f *msgFilter) ContainsBotMention(evt interface{}) interface{} {
	var msg *disgord.Message
	if msg = getMsg(evt); msg == nil {
		return nil
	}

	if !strings.Contains(msg.Content, mentionString(f.botID)) {
		return nil
	}

	return evt
}

func (f *msgFilter) NotByBot(evt interface{}) interface{} {
	return messageIsBot(evt, false)
}

func (f *msgFilter) IsByBot(evt interface{}) interface{} {
	return messageIsBot(evt, true)
}

func (f *msgFilter) HasBotMentionPrefix(evt interface{}) interface{} {
	return messageHasPrefix(evt, mentionString(f.botID))
}

func (f *msgFilter) HasPrefix(evt interface{}) interface{} {
	if f.prefix == "" {
		return evt
	}

	return messageHasPrefix(evt, f.prefix)
}

func (f *msgFilter) StripPrefix(evt interface{}) interface{} {
	if f.prefix == "" {
		return evt
	}

	if content := messageHasPrefix(evt, f.prefix); content == nil {
		return nil
	}

	msg := getMsg(evt)
	msg.Content = msg.Content[len(f.prefix):]
	return evt
}

func (f *msgFilter) HasPermissions(evt interface{}) interface{} {
	msg := getMsg(evt)
	uID := msg.Author.ID
	if uID.IsZero() {
		return nil
	}

	p, err := f.s.Guild(msg.GuildID).GetMemberPermissions(uID)
	if err != nil {
		return nil
	}

	if (p & f.permissions) != f.permissions {
		return nil
	}

	if f.eitherPermissions > 0 && (p&f.eitherPermissions) == 0 {
		return nil
	}

	return msg
}

// SetMinPermissions enforces message authors to have at least the given permission flags
// for the HasPermissions method to succeed
func (f *msgFilter) SetMinPermissions(min disgord.PermissionBit) {
	f.permissions = min
}

// SetAltPermissions enforces message authors to have at least one of the given permission flags for the
// HasPermissions method to succeed
func (f *msgFilter) SetAltPermissions(bits ...disgord.PermissionBit) {
	var permissions disgord.PermissionBit
	for i := range bits {
		permissions |= bits[i]
	}
	f.eitherPermissions = permissions
}
