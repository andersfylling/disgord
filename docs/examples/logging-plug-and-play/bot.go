package main

import (
	"context"
	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
	"os"
)

var log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.DebugLevel,
}

func main() {
	// Set up a new Disgord client
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
		Logger:   log,
	})
	defer client.StayConnectedUntilInterrupted(context.Background())
}
