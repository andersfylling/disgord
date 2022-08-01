package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

var commands = []*disgord.CreateApplicationCommand{
	{
		Name:        "test_command",
		Description: "just testing",
		Options: []*disgord.ApplicationCommandOption{
			{
				Name:        "test_option",
				Type:        disgord.OptionTypeString,
				Description: "testing options",
				Choices: []*disgord.ApplicationCommandOptionChoice{
					{
						Name:  "test_choice",
						Value: "test_val",
					},
				},
			},
		},
	},
}

func main() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISCORD_BOT_TOKEN"),
		Logger:   log,
	})
	defer client.Gateway().StayConnectedUntilInterrupted()

	// note the permission and scope are the minimum requirements for slash commands to operate
	u, err := client.BotAuthorizeURL(disgord.PermissionUseSlashCommands, []string{
		"bot",
		"applications.commands",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(u)

	// register commands once the bot is ready
	client.Gateway().BotReady(func() {
		for i := range commands {
			// application command id is 0 here
			// on a ready event, the client is updated to store the application id
			// you can fetch the application id using the bot id (current user id) or copy it from
			// the discord page.
			if err = client.ApplicationCommand(0).Guild(486833611564253184).Create(commands[i]); err != nil {
				log.Fatal(err)
			}
		}
	})

	// Respond hello any related discord slash command
	client.Gateway().InteractionCreate(func(s disgord.Session, h *disgord.InteractionCreate) {
		fmt.Printf("%+v", *h)
		err := s.SendInteractionResponse(context.Background(), h, &disgord.CreateInteractionResponse{
			Type: 4,
			Data: &disgord.MessageInteractionResponseData{
				Content:    "hello",
				Components: []*disgord.MessageComponent{},
			},
		})
		if err != nil {
			log.Error(err)
		}
	})

}
