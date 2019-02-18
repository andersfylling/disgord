package std

import (
	"strings"

	"github.com/andersfylling/disgord"
)

type msgFilterdg interface {
	Myself() (*disgord.User, error)
}

func NewMsgFilter(client msgFilterdg) (filter *msgFilter, err error) {
	usr, err := client.Myself()
	if err != nil {
		return nil, err
	}

	filter = &msgFilter{}
	filter.botID = usr.ID

	return filter, nil
}

type msgFilter struct {
	botID  disgord.Snowflake
	prefix string
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
