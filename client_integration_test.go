// +build integration

package disgord

import (
	"context"
	"github.com/andersfylling/disgord/internal/logger"
	"os"
	"sync"
	"testing"
	"time"
)

var token = os.Getenv("DISGORD_TOKEN_INTEGRATION_TEST")

var guildTypical = struct {
	ID                  Snowflake
	TextChannelGeneral  Snowflake
	VoiceChannelGeneral Snowflake
	VoiceChannelOther1  Snowflake
	VoiceChannelOther2  Snowflake
}{
	ID:                  ParseSnowflakeString(os.Getenv("TEST_GUILD_TYPICAL_ID")),
	TextChannelGeneral:  ParseSnowflakeString(os.Getenv("TEST_GUILD_TYPICAL_TEXT_GENERAL")),
	VoiceChannelGeneral: ParseSnowflakeString(os.Getenv("TEST_GUILD_TYPICAL_VOICE_GENERAL")),
	VoiceChannelOther1:  ParseSnowflakeString(os.Getenv("TEST_GUILD_TYPICAL_VOICE_1")),
	VoiceChannelOther2:  ParseSnowflakeString(os.Getenv("TEST_GUILD_TYPICAL_VOICE_2")),
}

func validSnowflakes() {
	if guildTypical.ID.IsZero() {
		panic("missing id for typical guild")
	}
	if guildTypical.TextChannelGeneral.IsZero() {
		panic("missing id for typical guild TextChannelGeneral")
	}
	if guildTypical.VoiceChannelGeneral.IsZero() {
		panic("missing id for typical guild VoiceChannelGeneral")
	}
	if guildTypical.VoiceChannelOther1.IsZero() {
		panic("missing id for typical guild VoiceChannelOther1")
	}
	if guildTypical.VoiceChannelOther2.IsZero() {
		panic("missing id for typical guild VoiceChannelOther2")
	}
}

func TestClient(t *testing.T) {
	validSnowflakes()

	wg := &sync.WaitGroup{}

	status := &UpdateStatusPayload{
		Status: StatusIdle,
		Game: &Activity{
			Name: "hello",
		},
	}

	var c *Client
	wg.Add(1)
	t.Run("New", func(t *testing.T) {
		defer wg.Done()
		var err error
		c, err = NewClient(Config{
			BotToken:     token,
			DisableCache: true,
			Logger:       &logger.FmtPrinter{},
			Presence:     status,
		})
		if err != nil {
			t.Fatal("failed to initiate a client")
		}
	})
	wg.Wait()

	wg.Add(1)
	t.Run("premature-emit", func(t *testing.T) {
		defer wg.Done()
		if _, err := c.Emit(UpdateStatus, &UpdateStatusPayload{}); err == nil {
			t.Fatal("Emit should have failed as no shards have been connected (initialised)")
		}
	})
	wg.Wait()

	// We need this for later.
	guildCreateEvent := make(chan *GuildCreate, 2)
	c.On(EvtGuildCreate, func(_ Session, evt *GuildCreate) {
		guildCreateEvent <- evt
	}, &Ctrl{Runs: 1})

	defer c.Disconnect()
	wg.Add(1)
	t.Run("connect", func(t *testing.T) {
		defer wg.Done()
		if err := c.Connect(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
	wg.Wait()

	wg.Add(1)
	t.Run("ready", func(t *testing.T) {
		defer wg.Done()
		ready := make(chan interface{}, 2)
		c.Ready(func() {
			ready <- true
		})
		select {
		case <-time.After(10 * time.Second):
			t.Fatal("unable to connect within time frame of 10s")
		case <-ready:
		}
	})
	wg.Wait()

	wg.Add(1)
	t.Run("default-presence", func(t *testing.T) {
		defer wg.Done()
		done := make(chan bool, 2)
		c.On(EvtPresenceUpdate, func(_ Session, evt *PresenceUpdate) {
			if !evt.User.Bot {
				c.Logger().Info("was not bot")
				return
			}
			usr, err := c.CurrentUser().Get()
			if err != nil {
				done <- false
				return
			}
			if evt.User.ID != usr.ID {
				return
			}

			if evt.Status != StatusIdle {
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
		if _, err := c.Emit(UpdateStatus, status); err != nil {
			t.Fatal(err)
		}

		select {
		case <-time.After(10 * time.Second):
			// yay
			// if no presence update is fired after calling emit,
			// that means that no change took place.
			// TODO: this test is fragile
		case success := <-done:
			if success {
				t.Fatal("unable to set presence at boot")
			}
		}
	})
	wg.Wait()

	// Add the voice state channel for later.
	voiceStateChan := make(chan *VoiceStateUpdate)

	wg.Add(1)
	t.Run("voice/MoveTo", func(t *testing.T) {
		defer wg.Done()
		deadline, _ := context.WithDeadline(context.Background(), time.Now().Add(25*time.Second))

		oldChannelID := guildTypical.VoiceChannelGeneral
		newChannelID := guildTypical.VoiceChannelOther1
		connectedToVoiceChannel := make(chan bool)
		successfullyMoved := make(chan bool, 2)
		done := make(chan bool)

		c.On(EvtVoiceStateUpdate, func(_ Session, evt *VoiceStateUpdate) {
			myself, err := c.CurrentUser().Get()
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
			voiceStateChan <- evt
		})

		go func() {
			v, err := c.Guild(guildTypical.ID).VoiceConnect(oldChannelID)
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
		}()

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
	})
	wg.Wait()

	wg.Add(1)
	t.Run("test-member-guild-user-id-non-zero", func(t *testing.T) {
		defer wg.Done()
		deadline, _ := context.WithDeadline(context.Background(), time.Now().Add(25*time.Second))

		// Test guild create event
		select {
		case x := <-guildCreateEvent:
			firstMember := x.Guild.Members[0]
			if firstMember.GuildID == 0 {
				panic("GuildID is zero")
			} else if firstMember.UserID == 0 {
				panic("UserID is zero")
			}
		case <-deadline.Done():
			panic("guildCreateEvent did not emit")
		}

		// Test message create event
		snowflakeChan := make(chan Snowflake, 2)
		c.On(EvtMessageCreate, func(_ Session, evt *MessageCreate) {
			if evt.Message.Author.Bot && evt.Message.Member != nil {
				snowflakeChan <- evt.Message.Member.GuildID
				snowflakeChan <- evt.Message.Member.UserID
			}
		}, &Ctrl{Runs: 1})
		msg, err := c.WithContext(deadline).SendMsg(guildTypical.TextChannelGeneral, "Hello World!")
		if err != nil {
			panic(err)
		}
		select {
		case x := <-snowflakeChan:
			if x == 0 {
				panic("GuildID is zero")
			}
		case <-deadline.Done():
			panic("snowflakeChan did not emit")
		}
		if <-snowflakeChan == 0 {
			panic("UserID is zero")
		}

		// Test message update event
		snowflakeChan = make(chan Snowflake, 2)
		c.On(EvtMessageUpdate, func(_ Session, evt *MessageUpdate) {
			if evt.Message.Author.Bot {
				snowflakeChan <- evt.Message.Member.GuildID
				snowflakeChan <- evt.Message.Member.UserID
			}
		}, &Ctrl{Runs: 1})
		_, err = c.Channel(guildTypical.TextChannelGeneral).Message(msg.ID).Update(deadline).SetContent("world").Execute()
		if err != nil {
			panic(err)
		}
		select {
		case x := <-snowflakeChan:
			if x == 0 {
				panic("GuildID is zero")
			}
		case <-deadline.Done():
			panic("snowflakeChan did not emit")
		}
		if <-snowflakeChan == 0 {
			panic("UserID is zero")
		}

		// GC the message
		_ = c.Channel(guildTypical.TextChannelGeneral).Message(msg.ID).Delete(deadline)

		// Handle voice state update
		select {
		case x := <-voiceStateChan:
			if x.Member.GuildID == 0 {
				panic("GuildID is zero")
			} else if x.Member.UserID == 0 {
				panic("UserID is zero")
			}
		case <-deadline.Done():
			panic("voiceStateChan did not emit")
		}

		// Test getting a member
		member, err := c.Guild(guildTypical.ID).Member(c.myID).WithContext(deadline).Get(IgnoreCache)
		if err != nil {
			panic(err)
		}
		if member.GuildID == 0 {
			panic("GuildID is zero")
		} else if member.UserID == 0 {
			panic("UserID is zero")
		}
	})
	wg.Wait()
}

func TestConnectWithShards(t *testing.T) {
	validSnowflakes()

	<-time.After(6 * time.Second) // avoid identify abuse
	c := New(Config{
		BotToken:     token,
		DisableCache: true,
		Logger:       &logger.FmtPrinter{},
		ShardConfig: ShardConfig{
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
	validSnowflakes()

	<-time.After(6 * time.Second) // avoid identify abuse
	createInstance := func(shardIDs []uint, shardCount uint) *Client {
		return New(Config{
			BotToken:     token,
			DisableCache: true,
			Logger:       &logger.FmtPrinter{},
			ShardConfig: ShardConfig{
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
	var instances []*Client
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
