## Listen for messages

In Disgord it is required that you specify the event you are listening to using an event key (see package event). Unlike discordgo this library does not use reflection to understand which event type your function reacts to.
```GoLang
session, err := disgord.NewSession(&disgord.Config{
    Token: os.Getenv("DISGORD_TOKEN"),
})
if err != nil {
    panic(err)
}

// create a handler and bind it to new message events
// handlers/listener are run in sequence if you register more than one
// so you should not need to worry about locking your objects unless you do any
// parallel computing with said objects
session.On(disgord.EventMessageCreate, func(session disgord.Session, data *disgord.MessageCreate) {
    fmt.Println(data.Message.Content)
})

// connect to the discord gateway to receive events
err = session.Connect()
if err != nil {
    panic(err)
}

// Keep the socket connection alive, until you terminate the application
session.DisconnectOnInterrupt()
```

In addition, Disgord also supports the use of channels for handling events.
```GoLang
session, err := disgord.NewSession(&disgord.Config{
    Token: os.Getenv("DISGORD_TOKEN"),
})
if err != nil {
    panic(err)
}

// or use a channel to listen for events
go func() {
    for {
        var msg *disgord.Message

        // wait for a new message
        select {
        case data, alive := <- session.Evt().MessageCreateChan():
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
err = session.Connect()
if err != nil {
    panic(err)
}

// Keep the socket connection alive, until you terminate the application
session.DisconnectOnInterrupt()
```

> **Note:** That you might experience parallel handling of event objects if you choose to use channels. However, this will only happen if you use the same channel in two or more of your own goroutines.
