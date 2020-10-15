## Ping-Pong bot
So the time has come where you want to be a bot engineer huh? In this article you are introduced to creating the common ping-pong bot. This snippet will contain the main function's body.

```go
package main

import (
	"context"
	"os"

	"github.com/andersfylling/disgord"
)

func main() {
    // configure a Disgord client
    client := disgord.New(disgord.Config{
        BotToken: os.Getenv("DISGORD_TOKEN"),
    })
    defer client.StayConnectedUntilInterrupted(context.Background())
    
    // create a handler and bind it to new message events
    client.On(disgord.EvtMessageCreate, func(s disgord.Session, evt *disgord.MessageCreate) {
        msg := evt.Message
        if msg.Content == "ping" {
            _, _ = msg.Reply(context.Background(), s, "pong")
        }
    })
}
```


Disgord also offers middlewares and a std package to checking the msg content

```go
package main

import (
	"context"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
)

const NewMessage = disgord.EvtMessageCreate

func main() {
    // configure a Disgord client
    client := disgord.New(disgord.Config{
        BotToken: os.Getenv("DISGORD_TOKEN"),
    })
    defer client.StayConnectedUntilInterrupted(context.Background())
    
    // create a mdlw that only accepts messages with a "ping" prefix
    // tip: use this to identify bot commands
    content, _ := std.NewMsgFilter(context.Background(), client)
    content.SetPrefix("ping")
    
    // create a handler and bind it to new message events
    // we add a middleware/filter to ensure that the message content 
    // starts with "ping".
    client.On(NewMessage, content.HasPrefix, func(s disgord.Session, evt *disgord.MessageCreate) {
        _, _ = evt.Message.Reply(context.Background(), s, "pong")
    })
}
```
