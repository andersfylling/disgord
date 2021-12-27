## Listen for messages

In Disgord it is required that you specify the event you are listening to using an event key (see package event). Unlike discordgo this library does not figure out the event types by checking your handler signature, this is cause I believe anonymous functions should not be used. Hopefully this will make for cleaner code. 

> The somewhat annoying part is that you specify the event name as a parameter anyways.
```go
func printMessage(session disgord.Session, data *disgord.MessageCreate) {
    fmt.Println(data.Message.Content)
}

func main() {
    client := disgord.New(disgord.Config{
        BotToken: os.Getenv("DISGORD_TOKEN"),
    })
    // connect, and stay connected until a system interrupt takes place
    defer client.StayConnectedUntilInterrupted(context.Background())
    
    // create a handler and bind it to new message events
    // handlers/listener are run in sequence if you register more than one
    // so you should not need to worry about locking your objects unless you do any
    // parallel computing with said objects
    client.On(disgord.EvtMessageCreate, printMessage)
}
```

In addition, Disgord also supports the use of channels for handling events. They work just like registering handlers.
However, be careful if you plan on closing channels, as you might put disgord internals into a deadlock. Instead, use the disgord.Ctrl method CloseChannel to handle this.
```go
client := disgord.New(disgord.Config{
    BotToken: os.Getenv("DISGORD_TOKEN"),
})

// or use a channel to listen for events
msgCreateChan := make(chan *disgord.MessageCreate)
client.On(disgord.EvtMessageCreate, msgCreateChan)

go func() {
    for {
        var msg *disgord.Message

        // wait for a new message
        select {
        case data, alive := <- msgCreateChan:
            if !alive {
                fmt.Println("channel is dead")
                return
            }
            msg = data.Message
        }

        // print the message
        // note that, since you're using channels, you might need to lock your disgord objects.
        msg.RLock()
        fmt.Println(msg.Content)
        msg.RUnlock()
    }
}()
```

The nice part with go channels is that they can work as a load balancer. Thanks to disgord middlewares you can use the worker pattern for, say... playing computationally expensive games.
```go
func main() {
	client := disgord.New(&disgord.Config{
        BotToken: os.Getenv("DISGORD_TOKEN"),
    })
    defer client.StayConnectedUntilInterrupted(context.Background())
	
	maxNumberOfGames := 5
	workChan := make(chan *MessageCreate, maxNumberOfGames)
	
	// setup channel lifetime to stop being used after the five games are finished
	// we simply increment the Runs in each worker, if the game did not finish.
	// this will keep it alive until all the games are done.
	//
	// Note that disgord never closes channels for you since you create the channels
	// you have the ownership of them, and must handle closing on your own.
	// I recommend using Ctrl.OnRemove() and simply call ctrl.CloseChannel() in this situation.
	ctrl := &disgord.Ctrl{Channel: workChan, Runs: maxNumberOfGames}
	
	for i := 0; i < maxNumberOfGames; i++ {
		go chessWorker(s, workChan, ctrl)
	}
	
	// filters
    filter, _ := std.NewMsgFilter(context.Background(), client)
    filter.SetPrefix("!")
    
    chessFilter, _ := std.NewMsgFilter(context.Background(), client)
    chessFilter.SetPrefix("chess")
    
    moveFilter, _ := std.NewMsgFilter(context.Background(), client)
    moveFilter.SetPrefix("move")
	
	// command: !chess move e2e4
	s.On(disgord.EvtMessageCreate,
		// identify a command, copy the msg content and remove the prefix
        filter.NotByBot,
        filter.HasPrefix,
        log.LogMsg,
        std.CopyMsgEvt,  // message is reused, so create a copy to avoid issues
        filter.StripPrefix,
        
        // identify a chess command, remove the "chess" prefix
        chessFilter.HasPrefix,
        chessFilter.StripPrefix,
        
        // identify the sub command "move" and remove the "move" prefix
        moveFilter.HasPrefix,
        moveFilter.StripPrefix,
        
        // send the message to one of the workers
		workChan, ctrl)
}

func chessWorker(s disgord.Session, workChan chan *disgord.MessageCreate, ctrl *disgord.Ctrl) {
	for {
		var msg *disgord.Message
		select {
		case evt, open := <- workChan:
		    if !open {
		        s.Logger().Info("worker channel closed")
		        return
		    }
		    msg = evt.Message
		}
		
		// this section is more psuedo code, won't actually work unless you have a 
		// chess engine
		move := msg.Content 
		game := getChessGame(msg.Author.ID)
		game.ApplyMove(move)
		gameStatus := game.CalculateBestMoveAndExecute()
		
		if !gameStatus.Finished() {
			// yes, you need locking.
			// increment Runs such that the player can do another move
            // if the ctrl hits 0 or lower, the handler is removed. 
			ctrl.Runs++
		}
		
		m.Reply(s, gameStatus.String())
	}
}
```
There are libs for writing simpler commands. So there is no need for you to actually use middlewares for this purpose. This is just an example using pure disgord functionality.

> **Note:** That you might experience parallel handling of event objects if you choose to use channels. However, this will only happen if you use the same channel in two or more of your own goroutines.
