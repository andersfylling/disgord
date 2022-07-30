package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/andersfylling/disgord"
)

var (
	appID  disgord.Snowflake = disgord.Snowflake(0)
	client disgord.Session
)

var log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.InfoLevel,
}

func main() {
	slashCmd := &disgord.CreateApplicationCommand{
		Name:        "testcmd",
		Description: "makes a bunch of example message components in the channel the command was invoked",
	}
	client = disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
		Logger:   log,
		Intents:  disgord.IntentGuildMessages,
	})
	defer client.Gateway().StayConnectedUntilInterrupted()
	client.Gateway().BotReady(func() {
		err := client.ApplicationCommand(appID).Global().Create(slashCmd)
		if err != nil {
			log.Error(err)
		}
	})
	client.Gateway().InteractionCreateChan(testcmdIntCreateHandler)
	go testcmdHandler()
}
