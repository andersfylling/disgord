package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
)

var log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.InfoLevel,
}

const MessageLifeTime = 5 // seconds

func deleteDeadMessage(session disgord.Session, msg *disgord.Message, lifetime time.Duration) {
	<-time.After(lifetime)
	if err := session.Channel(msg.ChannelID).Message(msg.ID).Delete(); err != nil {
		log.Error(fmt.Errorf("failed to delete message: %w", err))
	}
}

// please consider using a queue instead
func autoDeleteNewMessages(session disgord.Session, evt *disgord.MessageCreate) {
	lifetime := time.Duration(MessageLifeTime) * time.Second
	go deleteDeadMessage(session, evt.Message, lifetime)
}

func main() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
		Intents:  disgord.IntentGuildMessages,
	})
	defer client.Gateway().StayConnectedUntilInterrupted()

	filter, _ := std.NewMsgFilter(context.Background(), client)
	filter.SetMinPermissions(disgord.PermissionManageMessages) // make sure u can actually delete messages

	client.Gateway().
		WithMiddleware(filter.HasPermissions).
		MessageCreate(autoDeleteNewMessages)
}
