package std

import (
	"github.com/andersfylling/disgord"
)

const (
	thinking = "🤔"
	failure  = "👎"
	success  = "👌"
)

func WithProgressReactions(cmd func(s disgord.Session, evt *disgord.MessageCreate) error) func(s disgord.Session, evt *disgord.MessageCreate) {
	return func(s disgord.Session, evt *disgord.MessageCreate) {
		if evt.Message.Author != nil && evt.Message.Author.Bot {
			return
		}
		s.CreateReaction(evt.Message.ChannelID, evt.Message.ID, thinking)
		defer func() {
			s.DeleteOwnReaction(evt.Message.ChannelID, evt.Message.ID, thinking)
		}()
		if err := cmd(s, evt); err != nil {
			s.CreateReaction(evt.Message.ChannelID, evt.Message.ID, failure)
			return
		}
		s.CreateReaction(evt.Message.ChannelID, evt.Message.ID, success)
	}
}
