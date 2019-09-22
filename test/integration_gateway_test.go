// +build integration

package test

import (
	"os"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/andersfylling/disgord"
)

var token = os.Getenv("DISGORD_TOKEN_INTEGRATION_TEST")

func TestConnect(t *testing.T) {
	c := disgord.New(&disgord.Config{
		BotToken:     token,
		DisableCache: true,
	})
	defer c.Disconnect()
	if err := c.Connect(); err != nil {
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
	c := disgord.New(&disgord.Config{
		BotToken:     token,
		DisableCache: true,
		ShardConfig: disgord.ShardConfig{
			ShardIDs: []uint{0, 1},
		},
	})
	defer c.Disconnect()
	if err := c.Connect(); err != nil {
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
	createInstance := func(shardIDs []uint, shardCount uint) *disgord.Client {
		return disgord.New(&disgord.Config{
			BotToken:     token,
			DisableCache: true,
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
		if err := instance.Connect(); err != nil {
			cancel()
			t.Fatal(err)
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
		t.Fatal("unable to connect within time frame")
	case <-done:
	}
}
