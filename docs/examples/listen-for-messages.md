## Listen for messages

In Disgord it is required that you specify the event you are listening to using an event key (see package event). Unlike discordgo this library does not use reflection to understand which event type your function reacts to.
```go

client := disgord.New(&disgord.Config{
    BotToken: os.Getenv("DISGORD_TOKEN"),
    Logger: disgord.DefaultLogger(false), // optional logging, debug=false
})

// create a handler and bind it to new message events
// handlers/listener are run in sequence if you register more than one
// so you should not need to worry about locking your objects unless you do any
// parallel computing with said objects
client.On(disgord.EventMessageCreate, func(session disgord.Session, data *disgord.MessageCreate) {
    fmt.Println(data.Message.Content)
})

// connect to the discord gateway to receive events
err = client.Connect()
if err != nil {
    panic(err)
}

// Keep the socket connection alive, until you terminate the application
client.DisconnectOnInterrupt()
```

Note that if you dislike the long `disgord.EventMessageCreate` name. You can use the sub package `event`. However, the `event` package will only hold valid Discord events.
```go 
session.On(event.MessageCreate, func(session disgord.Session, data *disgord.MessageCreate) {
    fmt.Println(data.Message.Content)
})
```

In addition, Disgord also supports the use of channels for handling events. It is extremely important that you remember to tell disgord which events you want, before you start using channels. Failure to do so may result in the desired events being discarded at the socket layer, as you did not notify you wanted them. See the code example below.
```go
client := disgord.New(&disgord.Config{
    BotToken: os.Getenv("DISGORD_TOKEN"),
    Logger: disgord.DefaultLogger(false), // optional logging, debug=false
})

// or use a channel to listen for events
go func() {
    // channels are more advanced and requires you to register which event-channels
    // you will be using, in advanced. Otherwise you may not receive an event on the given channel.
    // See disgord.Session.AcceptEvent
    client.AcceptEvent(disgord.EventGuildCreate) // IMPORTANT!
    var messageCreateChan = client.EventChannels().MessageCreate()
    for {
        var msg *disgord.Message

        // wait for a new message
        select {
        case data, alive := <- messageCreateChan:
            if !alive {
                fmt.Println("channel is dead")
                return
            }
            msg = data.Message
        }

        // print the message
        // locking in case you are using the same channel somewhere else as well
        msg.RLock()
        fmt.Println(msg.Content)
        msg.RUnlock()
    }
}()

// connect to the discord gateway to receive events
err = client.Connect()
if err != nil {
    panic(err)
}

// Keep the socket connection alive, until you terminate the application
client.DisconnectOnInterrupt()
```

> **Note:** That you might experience parallel handling of event objects if you choose to use channels. However, this will only happen if you use the same channel in two or more of your own goroutines.
