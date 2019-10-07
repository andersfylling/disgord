package std

import (
	"strconv"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/internal/logger"
)

type logFilterdg interface {
	Logger() logger.Logger
}

func NewLogFilter(client logFilterdg) (filter *logFilter, err error) {
	filter = &logFilter{}
	filter.log = client.Logger()

	return filter, nil
}

type logFilter struct {
	log logger.Logger
}

// LogMsg logs messages in the format:
//  - MESSAGE_CREATE: $user{usrid} created message{msgid} $content
//  - MESSAGE_UPDATE: $user{usrid} updated message{msgid} $newContent
//  - MESSAGE_DELETE: messages{msgid} was deleted
//  - MESSAGE_DELETE_BULK: messages{len: 2, ids: [1234,34567]} was deleted
// if unable to log the event, it does not exit the chain.
func (f *logFilter) LogMsg(evt interface{}) interface{} {
	var change string
	switch t := evt.(type) {
	case *disgord.MessageCreate:
		change = "created"
	case *disgord.MessageUpdate:
		change = "updated"
	case *disgord.MessageDelete:
		msgStr := "message{" + t.MessageID.String() + "}"
		f.log.Info(msgStr, "was deleted")
		return evt
	case *disgord.MessageDeleteBulk:
		if len(t.MessageIDs) == 0 {
			f.log.Info("0 messages was deleted")
			return evt
		}

		msgIds := "len: " + strconv.Itoa(len(t.MessageIDs)) + ", ids: ["
		for i := range t.MessageIDs {
			msgIds += t.MessageIDs[i].String() + ","
		}
		msgIds = msgIds[:len(msgIds)-1] + "]"
		msgStr := "messages{" + msgIds + "}"
		f.log.Info(msgStr, "was deleted")
		return evt
	default:
		f.log.Error("unable to log msg event", t)
		return evt
	}

	msg := getMsg(evt)
	user := msg.Author.String()
	content := msg.Content
	msgStr := "message{" + msg.ID.String() + "}"
	f.log.Info(user, change, msgStr, content)

	return evt
}
