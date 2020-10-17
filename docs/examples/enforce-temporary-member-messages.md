If you ever want to create a channel where the messages are deleted after N seconds, kinda like snapchat, see the code below.


```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
)

const MessageLifeTime = 5 // seconds

func deleteDeadMessage(session disgord.Session, message *disgord.Message, lifetime time.Duration) {
	<-time.After(lifetime)
	if err := session.DeleteFromDiscord(context.Background(), message); err != nil {
		fmt.Println(err)
	}
}

// please consider using a queue instead
func autoDeleteNewMessages(session disgord.Session, evt *disgord.MessageCreate) {
	lifetime := time.Duration(MessageLifeTime) * time.Second
	go deleteDeadMessage(session, evt.Message, lifetime)
}

func main() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
	})
    defer client.StayConnectedUntilInterrupted(context.Background())
	
	filter, _ := std.NewMsgFilter(context.Background(), client)
	filter.SetMinPermissions(disgord.PermissionManageMessages) // make sure u can actually delete messages

	client.On(disgord.EvtMessageCreate, filter.HasPermissions, autoDeleteNewMessages)
}
```
