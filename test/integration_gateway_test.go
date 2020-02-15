// +build integration

package test

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/andersfylling/disgord"
)

func TestConnect(t *testing.T) {
	<-time.After(6 * time.Second) // avoid identify abuse
	c := disgord.New(disgord.Config{
		BotToken:     token,
		DisableCache: true,
		Logger:       disgord.DefaultLogger(true),
	})
	defer c.Disconnect()
	if err := c.Connect(context.Background()); err != nil {
		t.Fatal(err)
	}

	done := make(chan interface{}, 2)
	c.Ready(func() {
		done <- true
	})

	select {
	case <-time.After(10 * time.Second):
		t.Fatal("unable to connect within time frame of 10s")
	case <-done:
	}
}

func TestConnectWithShards(t *testing.T) {
	<-time.After(6 * time.Second) // avoid identify abuse
	c := disgord.New(disgord.Config{
		BotToken:     token,
		DisableCache: true,
		Logger:       disgord.DefaultLogger(true),
		ShardConfig: disgord.ShardConfig{
			ShardIDs: []uint{0, 1},
		},
	})
	defer c.Disconnect()
	if err := c.Connect(context.Background()); err != nil {
		t.Fatal(err)
	}

	done := make(chan interface{}, 2)
	c.Ready(func() {
		done <- true
	})

	select {
	case <-time.After(15 * time.Second):
		t.Fatal("unable to connect within time frame of 10s")
	case <-done:
	}
}

func TestConnectWithSeveralInstances(t *testing.T) {
	<-time.After(6 * time.Second) // avoid identify abuse
	createInstance := func(shardIDs []uint, shardCount uint) *disgord.Client {
		return disgord.New(disgord.Config{
			BotToken:     token,
			DisableCache: true,
			Logger:       disgord.DefaultLogger(true),
			ShardConfig: disgord.ShardConfig{
				ShardIDs:   shardIDs,
				ShardCount: shardCount,
			},
		})
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(20*time.Second))
	done := make(chan interface{}, 2)
	instanceReady := make(chan interface{}, 3)
	go func() {
		untilZero := 2
		for {
			select {
			case <-instanceReady:
				untilZero--
			case <-ctx.Done():
				return
			}

			if untilZero == 0 {
				done <- true
				return
			}
		}
	}()

	shardCount := uint(2)
	var instances []*disgord.Client
	for i := uint(0); i < shardCount; i++ {
		instance := createInstance([]uint{i}, shardCount)
		instances = append(instances, instance)

		instance.Ready(func() {
			instanceReady <- true
		})
		if err := instance.Connect(context.Background()); err != nil {
			cancel()
			t.Error(err)
			return
		}
		<-time.After(5 * time.Second)
	}

	defer func() {
		for i := range instances {
			_ = instances[i].Disconnect()
		}
	}()
	select {
	case <-ctx.Done():
		t.Error("unable to connect within time frame")
	case <-done:
	}
}

func TestFailOnPrematureEmit(t *testing.T) {
	// TODO: update when a queue is added/or whatever, for c.Emit before c.Connect takes place
	c := disgord.New(disgord.Config{
		BotToken:     "dkjfhslkjfhksf",
		DisableCache: true,
		ShardConfig: disgord.ShardConfig{
			ShardIDs: []uint{0, 1},
		},
	})
	_, err := c.Emit(disgord.UpdateStatus, &disgord.UpdateStatusPayload{
		Status: "hello",
	})
	if err == nil {
		t.Fatal("Emit should have failed as no shards have been connected (initialised)")
	}
}

func TestDefaultStatus(t *testing.T) {
	<-time.After(6 * time.Second) // avoid identify abuse
	c := disgord.New(disgord.Config{
		BotToken:     token,
		DisableCache: true,
		Logger:       disgord.DefaultLogger(true),
		Presence: &disgord.UpdateStatusPayload{
			Status: disgord.StatusIdle,
			Game: &disgord.Activity{
				Name: "hello",
			},
		},
	})
	defer c.Disconnect()

	done := make(chan bool, 2)
	c.On(disgord.EvtPresenceUpdate, func(_ disgord.Session, evt *disgord.PresenceUpdate) {
		if !evt.User.Bot {
			return
		}
		usr, err := c.GetCurrentUser(context.Background())
		if err != nil {
			done <- false
			return
		}
		if evt.User.ID != usr.ID {
			return
		}

		if evt.Status != disgord.StatusIdle {
			done <- false
			return
		}
		if evt.Game == nil {
			done <- false
			return
		}
		if evt.Game.Name != "hello" {
			done <- false
			return
		}

		done <- true
	})
	_ = c.Connect(context.Background())

	select {
	case <-time.After(20 * time.Second):
		t.Fatal("unable to connect within time frame of 20s")
	case success := <-done:
		if !success {
			t.Fatal("was unable to set bot presence")
		}
	}
}
