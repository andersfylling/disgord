package std

import "github.com/andersfylling/disgord"

func mentionString(id disgord.Snowflake) string {
	return "<@" + id.String() + ">"
}

func getMsg(evt interface{}) (msg *disgord.Message) {
	switch t := evt.(type) {
	case *disgord.MessageCreate:
		msg = t.Message
	case *disgord.MessageUpdate:
		msg = t.Message
	default:
		msg = nil
	}

	return msg
}
