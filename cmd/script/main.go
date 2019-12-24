package main

import (
	"context"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
)

// replyPongToPing is a handler that replies pong to ping messages
func replyPongToPing(s disgord.Session, data *disgord.MessageCreate) {
	msg := data.Message

	// whenever the message written is "ping", the bot replies "pong"
	if msg.Content == "ping" {
		_, _ = msg.Reply(context.Background(), s, "pong")
	}
}

func main() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
		Logger:   disgord.DefaultLogger(false), // debug=false
	})
	defer client.StayConnectedUntilInterrupted(context.Background())

	log, _ := std.NewLogFilter(client)
	filter, _ := std.NewMsgFilter(context.Background(), client)
	filter.SetPrefix("REPLACE_ME")

	// create a handler and bind it to new message events
	// tip: read the documentation for std.CopyMsgEvt and understand why it is used here.
	client.On(disgord.EvtMessageCreate,
		// middleware
		filter.NotByBot,    // ignore bot messages
		filter.HasPrefix,   // read original
		log.LogMsg,         // log command message
		std.CopyMsgEvt,     // read & copy original
		filter.StripPrefix, // write copy
		// handler
		replyPongToPing) // handles copy
}
