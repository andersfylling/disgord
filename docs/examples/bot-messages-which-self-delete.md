Ever wanted a bot that delete their own messages after N seconds? Well here you go.
With this approach you register a listener for bot messages which you then call delete on.

You can also do a goroutine to delete a message right after you send it, but this can be tedious. So instead we add a listener to delete your own messages.


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
	// you don't have to do this though
	if evt.Message.Author.ID != session.Myself().ID {
		return
	}

	// delete message after N seconds
	lifetime := time.Duration(MessageLifeTime) * time.Second
	go deleteMessageAfterTimeout(session, evt.Message, lifetime)
}

func main() {
	session, err := disgord.NewSession(&disgord.Config{
		Token: os.Getenv("DISGORD_TOKEN"),
	})
	if err != nil {
		panic(err)
	}

	session.On(disgord.EventMessageCreate, autoDeleteNewMessages)

	// connect to the discord gateway to receive events
	err = session.Connect()
	if err != nil {
		panic(err)
	}

	// graceful shutdown
	session.DisconnectOnInterrupt()
}
```
