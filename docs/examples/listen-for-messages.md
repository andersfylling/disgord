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
session.AddListener(disgord.EventMessageCreate, func(session Session, data *disgord.MessageCreate) {
    fmt.Println(data.Message.Content)
})

// connect to the discord gateway to receive events
err = session.Connect()
if err != nil {
    panic(err)
}

// Keep the socket connection alive, until you terminate the application
<-termSignal
session.Disconnect()
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
                break
            }
            msg = data.Message
        }

        // print the message
        fmt.Println(msg.Content)
    }
}()

// connect to the discord gateway to receive events
err = session.Connect()
if err != nil {
    panic(err)
}

// Keep the socket connection alive, until you terminate the application
<-termSignal
session.Disconnect()
```

> **Note:** That the use of the "termSignal" is a variable you must create yourself. It is highly recommended to handle terminal signals within your own system to do graceful shutdowns. Disgord does not have an internal mechanic for graceful shutdown so you must call the Disconnect function as shown above.
