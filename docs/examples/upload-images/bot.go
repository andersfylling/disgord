package main

import (
	"os"

	"github.com/andersfylling/disgord"
)

const (
	ChannelID = disgord.Snowflake(93284097324)
)

// Uploading images as attachments or even to supply them in Embeds is quite easy, images not used
// in the same post will be added as attachments and images accessed using the `attachment://`
// scheme will be used on their respective locations.
func main() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISCORD_TOKEN"),
	})

	f1, err := os.Open("myfavouriteimage.jpg")
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	f2, err := os.Open("another.jpg")
	if err != nil {
		panic(err)
	}
	defer f2.Close()

	_, errUpload := client.Channel(ChannelID).CreateMessage(&disgord.CreateMessageParams{
		Content: "This is my favourite image, and another in an embed!",
		Files: []disgord.CreateMessageFileParams{
			{f1, "myfavouriteimage.jpg", false},
			{f2, "another.jpg", false},
		},
		Embed: &disgord.Embed{
			Description: "Look here!",
			Image: &disgord.EmbedImage{
				URL: "attachment://another.jpg",
			},
		},
	})
	if errUpload != nil {
		client.Logger().Error("unable to upload images.", errUpload)
	}
}
