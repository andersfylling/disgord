// +build integration

package disgord

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/andersfylling/disgord/internal/logger"
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

var guildAdmin = struct {
	ID Snowflake
}{
	ID: ParseSnowflakeString(os.Getenv("TEST_GUILD_ADMIN_ID")),
}

func validSnowflakes() {
	if guildAdmin.ID.IsZero() {
		panic("missing id for admin guild")
	}
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
		c, err = NewClient(context.Background(), Config{
			BotToken: token,
			Logger:   &logger.FmtPrinter{},
			Presence: status,
		})
		if err != nil {
			t.Fatal("failed to initiate a client")
		}
	})
	wg.Wait()

	wg.Add(1)
	t.Run("premature-emit", func(t *testing.T) {
		defer wg.Done()
		if _, err := c.Gateway().Dispatch(UpdateStatus, &UpdateStatusPayload{}); err == nil {
			t.Fatal("Emit should have failed as no shards have been connected (initialised)")
		}
	})
	wg.Wait()

	// We need this for later.
	guildCreateEvent := make(chan *GuildCreate, 2)
	c.Gateway().WithCtrl(&Ctrl{Runs: 1}).GuildCreate(func(_ Session, evt *GuildCreate) {
		guildCreateEvent <- evt
	})

	defer c.Gateway().Disconnect()
	wg.Add(1)
	t.Run("connect", func(t *testing.T) {
		defer wg.Done()
		if err := c.Gateway().Connect(); err != nil {
			t.Fatal(err)
		}
	})
	wg.Wait()

	wg.Add(1)
	t.Run("ready", func(t *testing.T) {
		defer wg.Done()
		ready := make(chan interface{}, 2)
		c.Gateway().BotReady(func() {
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
	t.Run("role", func(t *testing.T) {
		defer wg.Done()
		roleName := "t-" + strconv.Itoa(rand.Int())
		if len(roleName) > 10 {
			roleName = roleName[:10]
		}

		created := make(chan interface{}, 2)
		c.Gateway().GuildRoleCreate(func(s Session, h *GuildRoleCreate) {
			created <- 0
		})

		deleted := make(chan interface{}, 2)
		c.Gateway().GuildRoleDelete(func(s Session, h *GuildRoleDelete) {
			deleted <- 0
		})

		var roleID Snowflake
		t.Run("create", func(t *testing.T) {
			createdRole, err := c.Guild(guildAdmin.ID).CreateRole(&CreateGuildRoleParams{
				Name:   roleName,
				Reason: "integration test",
			})
			if err != nil {
				t.Error(fmt.Errorf("unable to create role. %w", err))
				return
			}

			select {
			case <-time.After(10 * time.Second):
				t.Error("failed to get role create event within time frame of 10s")
			case <-created:
				close(created)
			}

			guild, err := c.Cache().GetGuild(guildAdmin.ID)
			if err != nil || guild == nil {
				t.Error("somehow the admin guild is not in the cache")
				return
			}

			r, err := guild.Role(createdRole.ID)
			if err != nil {
				t.Fatal("role does not exist in cache")
			}

			if r.Name != roleName {
				t.Errorf("role name differs. Got %s, wants %s", r.Name, roleName)
			}

			roleID = createdRole.ID
		})

		t.Run("delete", func(t *testing.T) {
			if roleID.IsZero() {
				t.Fatal("unable to test role delete, as role create failed")
			}

			if err := c.Guild(guildAdmin.ID).Role(roleID).Delete(); err != nil {
				t.Error(fmt.Errorf("unable to delete role. %w", err))
				return
			}

			select {
			case <-time.After(10 * time.Second):
				t.Error("failed to get role deleted event within time frame of 10s")
			case <-deleted:
				close(deleted)
			}

			guild, err := c.Cache().GetGuild(guildAdmin.ID)
			if err != nil || guild == nil {
				t.Error("somehow the admin guild is not in the cache")
				return
			}

			if r, _ := guild.Role(roleID); r != nil {
				t.Fatal("role exist in cache")
			}
		})
	})
	wg.Wait()

	wg.Add(1)
	t.Run("default-presence", func(t *testing.T) {
		defer wg.Done()
		done := make(chan bool, 2)
		c.Gateway().PresenceUpdate(func(_ Session, evt *PresenceUpdate) {
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

			game, err := evt.Game()
			if err != nil {
				done <- false
				return
			}
			if game.Name != "hello" {
				done <- false
				return
			}

			done <- true
		})
		if _, err := c.Gateway().Dispatch(UpdateStatus, status); err != nil {
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

		c.Gateway().VoiceStateUpdate(func(_ Session, evt *VoiceStateUpdate) {
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
			v, err := c.Guild(guildTypical.ID).VoiceChannel(oldChannelID).Connect(false, true)
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
	t.Run("middleware", func(t *testing.T) {
		defer wg.Done()
		deadline, _ := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))

		const prefix = "test"
		content := prefix + " sads sdfjsd fkjsdf"
		channelID := guildTypical.TextChannelGeneral

		gotMessage := make(chan *MessageCreate)
		defer close(gotMessage)

		filterTestPrefix := func(evt interface{}) (ret interface{}) {
			msg := (evt.(*MessageCreate)).Message
			if strings.HasPrefix(msg.Content, prefix) {
				return evt
			}
			return nil
		}
		filterChannel := func(evt interface{}) (ret interface{}) {
			msg := (evt.(*MessageCreate)).Message
			if msg.ChannelID == channelID {
				return evt
			}
			return nil
		}

		c.Gateway().WithMiddleware(filterChannel, filterTestPrefix).MessageCreateChan(gotMessage)
		_, err := c.Channel(channelID).WithContext(deadline).CreateMessage(&CreateMessageParams{Content: content})
		if err != nil {
			panic(fmt.Errorf("unable to send message. %w", err))
		}

		select {
		case msg := <-gotMessage:
			if msg.Message.Content != content {
				panic("unexpected message content")
			}
		case <-deadline.Done():
			panic("message create event did not trigger within the deadline")
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
		c.Gateway().WithCtrl(&Ctrl{Runs: 1}).MessageCreate(func(_ Session, evt *MessageCreate) {
			if evt.Message.Author.Bot && evt.Message.Member != nil {
				snowflakeChan <- evt.Message.Member.GuildID
				snowflakeChan <- evt.Message.Member.UserID
			}
		})
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
		c.Gateway().WithCtrl(&Ctrl{Runs: 1}).MessageUpdate(func(_ Session, evt *MessageUpdate) {
			if evt.Message.Author.Bot {
				snowflakeChan <- evt.Message.Member.GuildID
				snowflakeChan <- evt.Message.Member.UserID
			}
		})
		_, err = c.Channel(guildTypical.TextChannelGeneral).Message(msg.ID).WithContext(deadline).UpdateBuilder().SetContent("world").Execute()
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
		_ = c.Channel(guildTypical.TextChannelGeneral).Message(msg.ID).WithContext(deadline).Delete()

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
		member, err := c.Guild(guildTypical.ID).Member(c.botID).WithContext(deadline).Get(IgnoreCache)
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
	defer c.Gateway().Disconnect()
	if err := c.Gateway().Connect(); err != nil {
		t.Fatal(err)
	}

	done := make(chan interface{}, 2)
	c.Gateway().BotReady(func() {
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

		instance.Gateway().BotReady(func() {
			instanceReady <- true
		})
		if err := instance.Gateway().Connect(); err != nil {
			cancel()
			t.Error(err)
			return
		}
		<-time.After(5 * time.Second)
	}

	defer func() {
		for i := range instances {
			_ = instances[i].Gateway().Disconnect()
		}
	}()
	select {
	case <-ctx.Done():
		t.Error("unable to connect within time frame")
	case <-done:
	}
}

func TestREST(t *testing.T) {
	const andersfylling = Snowflake(769640669135896586)
	validSnowflakes()

	c := New(Config{
		BotToken:     token,
		DisableCache: true,
		Logger:       &logger.FmtPrinter{},
	})

	deadline, _ := context.WithDeadline(context.Background(), time.Now().Add(25*time.Second))

	// -------------------
	// CHANNELS
	// -------------------
	t.Run("channel", func(t *testing.T) {
		func() {
			channel, err := c.Channel(guildTypical.TextChannelGeneral).WithContext(deadline).Get()
			if err != nil {
				panic(err)
			} else if channel == nil {
				t.Error(fmt.Errorf("fetched channel is nil. %w", err))
			} else if channel.ID != guildTypical.TextChannelGeneral {
				t.Errorf("incorrect channel id. Got %s, wants %s", channel.ID.String(), guildTypical.TextChannelGeneral.String())
			}
		}()

		// create DM & send a message
		func() {
			channel, err := c.User(andersfylling).WithContext(deadline).CreateDM()
			if err != nil {
				t.Error(fmt.Errorf("unable to create DM with user. %w", err))
			} else if channel == nil {
				t.Error(fmt.Errorf("returned DM channel is nil. %w", err))
			}

			content := "hi"
			msg, err := c.Channel(channel.ID).WithContext(deadline).CreateMessage(&CreateMessageParams{Content: content})
			if err != nil {
				t.Error(fmt.Errorf("unable to create message in DM channel. %w", err))
			}
			if msg == nil {
				t.Error("returned message was nil")
			} else if msg.Content != content {
				t.Errorf("unexpected message content from DM. Got %s, wants %s", msg.Content, content)
			}
		}()
	})

	// -------------------
	// Current User
	// -------------------
	t.Run("current-user", func(t *testing.T) {
		if _, err := c.CurrentUser().Get(IgnoreCache); err != nil {
			t.Error(fmt.Errorf("unable to fetch current user. %w", err))
		}
	})

	// -------------------
	// User
	// -------------------
	t.Run("user", func(t *testing.T) {
		const userID = andersfylling
		user, err := c.User(userID).WithContext(deadline).Get(IgnoreCache)
		if err != nil {
			t.Error(fmt.Errorf("unable to fetch user. %w", err))
		} else if user == nil {
			t.Error("fetched user was nil")
		} else if user.ID != userID {
			t.Errorf("unexpected user id. Got %s, wants %s", user.ID.String(), userID.String())
		}
	})

	// -------------------
	// Voice Region
	// -------------------
	t.Run("voice-region", func(t *testing.T) {
		regions, err := c.WithContext(deadline).GetVoiceRegions(IgnoreCache)
		if err != nil {
			t.Error(fmt.Errorf("unable to fetch voice regions. %w", err))
		}
		if len(regions) < 1 {
			t.Error("expected at least one voice region")
		}
	})

	// -------------------
	// Members
	// -------------------
	t.Run("members", func(t *testing.T) {
		members, err := c.Guild(guildTypical.ID).GetMembers(nil, IgnoreCache)
		if err != nil {
			t.Error("failed to fetched members over REST, ", err)
		}

		if len(members) == 0 {
			t.Error("expected there to be members. None found.")
		}
	})

	// -------------------
	// Audit Logs
	// -------------------
	// t.Run("audit-logs", func(t *testing.T) {
	// })
}
