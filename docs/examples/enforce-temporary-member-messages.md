If you ever want to create a channel where the messages are deleted after N seconds, kinda like snapchat, see the code below.


```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/andersfylling/disgord"
)

const MessageLifeTime = 5 // seconds

func deleteDeadMessage(session disgord.Session, message *disgord.Message, lifetime time.Duration) {
	<-time.After(lifetime)
	err := session.DeleteFromDiscord(message)
	if err != nil {
		fmt.Println(err)
	}
}

func autoDeleteNewMessages(session disgord.Session, evt *disgord.MessageCreate) {
	lifetime := time.Duration(MessageLifeTime) * time.Second
	go deleteDeadMessage(session, evt.Message, lifetime)
}

func main() {
	client := disgord.New(&disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
	})

	client.On(disgord.EventMessageCreate, autoDeleteNewMessages)

	// connect to the discord gateway to receive events
	err = client.Connect()
	if err != nil {
		panic(err)
	}

	// graceful shutdown
	client.DisconnectOnInterrupt()
}
```
