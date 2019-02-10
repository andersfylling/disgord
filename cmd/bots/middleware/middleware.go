package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/event"
)

func filterTestPrefix(evt interface{}) (ret interface{}) {
	msg := (evt.(*disgord.MessageCreate)).Message

	if strings.HasPrefix(msg.Content, "test") {
		// returning evt also allow us to make a copy that can be manipulated, and sent through the chain
		return evt
	}

	return nil // stops
}

func main() {
	client := disgord.New(&disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
		Logger:   disgord.DefaultLogger(true),
	})

	client.On(event.MessageCreate, func() {
		fmt.Println("this should fire on every event")
	})

	client.On(event.MessageCreate, filterTestPrefix, func() {
		fmt.Println("this should fire on every event")
	})

}
