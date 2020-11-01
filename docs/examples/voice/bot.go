package main

import (
	"os"

	"github.com/andersfylling/disgord"
)

const (
	MyGuildID = disgord.Snowflake(26854385)
	MyChannelID = disgord.Snowflake(93284097324)
)

// In this example code example, you can see how to let the bot send an audio fragment after receiving a command.
func main() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISCORD_TOKEN"),
	})
	gateway := client.Gateway()
	defer gateway.StayConnectedUntilInterrupted()

	var voice disgord.VoiceConnection
	gateway.BotReady(func() {
		// Once the bot has connected to the websocket, also connect to the voice channel
		voice, _ = client.VoiceConnectOptions(MyGuildID, MyChannelID, false, true)
	})

	gateway.MessageCreate(func(_ disgord.Session, m *disgord.MessageCreate) {
		// Upon receiving a message with content !airhorn, play a sound to the connection made earlier
		if m.Message.Content == "!airhorn" {
			f, _ := os.Open("airhorn.dca")
			defer f.Close()

			_ = voice.StartSpeaking() // Sending a speaking signal is mandatory before sending voice data
			_ = voice.SendDCA(f) // Or use voice.SendOpusFrame, this blocks until done sending (realtime audio duration)
			_ = voice.StopSpeaking() // Tell Discord we are done sending data.
		}
	})
}