package main

import (
	"context"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
)

// So the time has come where you want to be a bot engineer huh?
// In this article you are introduced to creating the common
// ping-pong bot. This snippet will contain the main
// function's body.
func main() {
	// configure a Disgord client
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISCORD_TOKEN"),
	})
	defer client.Gateway().StayConnectedUntilInterrupted()

	// create a mdlw that only accepts messages with a "ping" prefix
	// tip: use this to identify bot commands
	content, _ := std.NewMsgFilter(context.Background(), client)
	content.SetPrefix("ping")

	client.Gateway().
		WithMiddleware(content.HasPrefix).
		MessageCreate(func(s disgord.Session, evt *disgord.MessageCreate) {
			_, _ = evt.Message.Reply(context.Background(), s, "pong")
		})
}
