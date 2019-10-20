Ever wanted a bot that delete their own messages after N seconds? Well here you go.
With this approach you register a listener for bot messages which you then call delete on.

You can also do a goroutine to delete a message right after you send it, this example allows you to handle multiple bots if desired. So instead we add a listener to delete your own messages.


```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/andersfylling/disgord"
)

const MessageLifeTime = 5 // seconds

func deleteMessageAfterTimeout(session disgord.Session, message *disgord.Message, lifetime time.Duration) {
	<-time.After(lifetime)
	err := session.DeleteFromDiscord(message)
	if err != nil {
		fmt.Println(err)
	}
}

func autoDeleteNewMessages(session disgord.Session, evt *disgord.MessageCreate) {
	// ignore humans
	if !evt.Message.Author.Bot {
		return
	}

	// ignore other bots
	// remove this check if you want to delete all bot messages after N seconds
	myself, err := session.Myself()
	if err != nil || evt.Message.Author.ID != myself.ID {
		return
	}

	// delete message after N seconds
	lifetime := time.Duration(MessageLifeTime) * time.Second
	go deleteMessageAfterTimeout(session, evt.Message, lifetime)
}

func main() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
		Logger: disgord.DefaultLogger(false), // optional logging, debug=false
	})

	client.On(disgord.EvtMessageCreate, autoDeleteNewMessages)

	// connect to the discord gateway to receive events
	if err := client.Connect(); err != nil {
		panic(err)
	}

	// graceful shutdown
	client.DisconnectOnInterrupt()
}
```
