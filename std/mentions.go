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
	botID disgord.Snowflake
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
	var msg *disgord.Message
	if msg = getMsg(evt); msg == nil {
		return nil
	}

	content := strings.TrimSpace(msg.Content)
	if !strings.HasPrefix(content, mentionString(f.botID)) {
		return nil
	}

	return evt
}
