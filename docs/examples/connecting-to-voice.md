In this example code example, you can see how to let the bot send an audio fragment after receiving a command.

**Please note that for the sake of brevity error handling in this example has been omitted,
you should never do this in actual code**

```go
package main

import (
	"context"
	"os"

	"github.com/andersfylling/disgord"
)

func main() {
	// Set up a new Disgord client
	discord := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
	})

	var voice disgord.VoiceConnection
	discord.Ready(func() {
		// Once the bot has connected to the websocket, also connect to the voice channel
		voice, _ = discord.VoiceConnect(myGuildID, myChannelID)
	})
	discord.On(disgord.EvtMessageCreate, func(_ disgord.Session, m *disgord.MessageCreate) {
		// Upon receiving a message with content !airhorn, play a sound to the connection made earlier
		if m.Message.Content == "!airhorn" {
			f, _ := os.Open("airhorn.dca")
			defer f.Close()

			_ = voice.StartSpeaking() // Sending a speaking signal is mandatory before sending voice data
			_ = voice.SendDCA(f) // Or use voice.SendOpusFrame, this blocks until done sending (realtime audio duration)
			_ = voice.StopSpeaking() // Tell Discord we are done sending data.
		}
	})

	_ = discord.Connect(context.Background())
	_ = discord.DisconnectOnInterrupt()
}

```
