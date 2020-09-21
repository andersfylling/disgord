package main

import (
	"os"

	"github.com/andersfylling/disgord"
)

func main() {
	c := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
		Logger:   disgord.DefaultLogger(true),
	})

	chanID := disgord.Snowflake(540519296640614416)
	_, err := c.Channel(chanID).CreateMessage(&disgord.CreateMessageParams{
		Content:           "testing",
		SpoilerTagContent: true,
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

	_, _ = c.Channel(chanID).CreateMessage(&disgord.CreateMessageParams{
		Content: "with embed",
		Files: []disgord.CreateMessageFileParams{
			{Reader: f1, FileName: "myfavouriteimage.jpg", SpoilerTag: true},
			{Reader: f2, FileName: "another.jpg"},
		},
		Embed: &disgord.Embed{
			Description: "Look here!",
			Image: &disgord.EmbedImage{
				URL: "attachment://another.jpg",
			},
		},
	})

	_, _ = c.Channel(chanID).CreateMessage(&disgord.CreateMessageParams{
		Content: "This is my favourite image, and another in an embed!",
		Files: []disgord.CreateMessageFileParams{
			{Reader: f1, FileName: "myfavouriteimage.jpg"},
			{Reader: f2, FileName: "another.jpg", SpoilerTag: true},
		},
		Embed: &disgord.Embed{
			Description: "Look here!",
			Image: &disgord.EmbedImage{
				URL: "attachment://another.jpg",
			},
		},
	})

	if err != nil {
		panic(err)
	}

	<-make(chan interface{})
}
