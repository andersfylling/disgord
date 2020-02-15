// +build integration

package test

import (
	"sync"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/andersfylling/disgord"
)

func TestVoice_ChangeChannel(t *testing.T) {
	<-time.After(6 * time.Second) // avoid identify abuse
	deadline, _ := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
	c := disgord.New(disgord.Config{
		BotToken:     token,
		DisableCache: true,
		Logger:       disgord.DefaultLogger(false),
	})
	defer c.Disconnect()
	if err := c.Connect(context.Background()); err != nil {
		t.Fatal(err)
	}

	oldChannelID := guildTypical.VoiceChannelGeneral
	newChannelID := guildTypical.VoiceChannelOther1
	connectedToVoiceChannel := make(chan bool)
	successfullyMoved := make(chan bool, 2)
	done := make(chan bool)
	defer close(successfullyMoved)

	c.On(disgord.EvtVoiceStateUpdate, func(_ disgord.Session, evt *disgord.VoiceStateUpdate) {
		myself, err := c.GetCurrentUser(context.Background())
		if err != nil {
			panic(err)
		}
		if evt.UserID != myself.ID {
			return
		}
		if evt.ChannelID == oldChannelID {
			connectedToVoiceChannel <- true
			return
		}
		if evt.ChannelID == newChannelID {
			successfullyMoved <- true
			successfullyMoved <- true
		} else {
			successfullyMoved <- false
			successfullyMoved <- false
		}
	})

	c.Ready(func() {
		v, err := c.VoiceConnect(guildTypical.ID, oldChannelID)
		if err != nil {
			t.Fatal(err)
		}

		select {
		case <-connectedToVoiceChannel:
		case <-deadline.Done():
			panic("connectedToVoiceChannel did not emit")
		}
		if err = v.MoveTo(newChannelID); err != nil {
			t.Fatal(err)
		}

		select {
		case <-successfullyMoved:
		case <-deadline.Done():
			panic("successfullyMoved did not emit")
		}
		defer func() {
			close(done)
		}()
		if err = v.Close(); err != nil {
			t.Fatal(err)
		}
		<-time.After(50 * time.Millisecond)
	})

	testFinished := sync.WaitGroup{}
	testFinished.Add(1)
	go func() {
		select {
		case <-time.After(10 * time.Second):
			t.Fatal("switching to a different voice channel failed")
		case success, ok := <-successfullyMoved:
			if !ok {
				t.Fatal("unexpected close of channel")
			}
			if !success {
				t.Fatal("did not go to a different voice channel")
			}
		}
		testFinished.Done()
	}()
	testFinished.Wait()

	select {
	case <-done:
	case <-deadline.Done():
		panic("done did not emit")
	}
}
