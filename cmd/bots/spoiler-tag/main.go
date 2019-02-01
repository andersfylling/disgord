package main

import (
	"os"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v3"
)

func main() {
	c, err := disgord.NewClient(&disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
		Logger:   disgord.DefaultLogger(true),
	})
	if err != nil {
		panic(err)
	}

	chanID := snowflake.ID(540519296640614416)
	_, err = c.CreateChannelMessage(chanID, &disgord.CreateChannelMessageParams{
		Content:           "testing",
		SpoilerTagContent: true,
	})

	if err != nil {
		panic(err)
	}
}
