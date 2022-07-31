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
	Level:     logrus.DebugLevel,
}

func main() {
	slashCmds := []*disgord.CreateApplicationCommand{
		{
			Name:        "example_buttons",
			Description: "makes example button components. these components don't do anything.",
			Type:        disgord.ApplicationCommandChatInput,
		},
		{
			Name:        "example_select_menu",
			Description: "makes example select menu. these components don't do anything.",
			Type:        disgord.ApplicationCommandChatInput,
		},
		{
			Name:        "example_modal",
			Description: "makes example modal. these components don't do anything.",
			Type:        disgord.ApplicationCommandChatInput,
		},
	}
	client = disgord.New(disgord.Config{
		BotToken: os.Getenv("DISCORD_BOT_TOKEN"),
		Logger:   log,
		Intents:  disgord.IntentGuildMessages,
	})
	defer client.Gateway().StayConnectedUntilInterrupted()
	client.Gateway().BotReady(func() {
		for i := 0; i < len(slashCmds); i++ {
			err := client.ApplicationCommand(appID).Guild(486833611564253184).Create(slashCmds[i])
			if err != nil {
				log.Error(err)
			}
		}
	})
	client.Gateway().InteractionCreateChan(
		exampleButtonsIntCreateHandler,
		exampleSelectMenuIntCreateHandler,
		exampleModalIntCreateHandler,
	)
	go exampleButtonsHandler()
	go exampleSelectMenuHandler()
	go exampleModal()
}
