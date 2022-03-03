package main

import (
	"context"
	"os"

	"github.com/andersfylling/disgord"
)

func interaction(session disgord.Session, evt *disgord.InteractionCreate) {
	if evt.Type == disgord.InteractionApplicationCommand {
		f1, err := os.Open("myfavouriteimage.jpg")
		if err != nil {
			panic(err)
		}
		defer f1.Close()
		session.SendInteractionResponse(context.Background(), evt, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "This is my favourite image",
				Files: []disgord.CreateMessageFile{
					{f1, "myfavouriteimage.jpg", false},
				},
				Embeds: []*disgord.Embed{{
					Description: "Look here!",
					Image: &disgord.EmbedImage{
						URL: "attachment://myfavouriteimage.jpg",
					},
				}},
			},
		})
	}
}
func main() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
	})
	defer client.Gateway().StayConnectedUntilInterrupted()
	client.Gateway().InteractionCreate(interaction)
}
