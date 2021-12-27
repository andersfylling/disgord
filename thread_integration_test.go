//go:build integration
// +build integration

package disgord

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/andersfylling/disgord/internal/logger"
)

func TestThreadEndpoints(t *testing.T) {
	const andersfylling = Snowflake(769640669135896586)
	validSnowflakes()

	c := New(Config{
		BotToken:     token,
		DisableCache: true,
		Logger:       &logger.FmtPrinter{},
	})

	deadline, _ := context.WithDeadline(context.Background(), time.Now().Add(25*time.Second))

	t.Run("create", func(t *testing.T) {
		threadName := "HELLO WORLD1"
		msg, err := c.Channel(guildAdmin.TextChannelGeneral).WithContext(deadline).CreateMessage(&CreateMessage{Content: threadName})
		if err != nil {
			panic(err)
		}

		thread, err := c.Channel(guildAdmin.TextChannelGeneral).WithContext(deadline).CreateThread(msg.ID, &CreateThread{
			Name:                threadName,
			AutoArchiveDuration: AutoArchiveThreadDay,
		})
		if err != nil {
			panic(err)
		} else if thread == nil {
			t.Error(fmt.Errorf("fetched thread is nil. %w", err))
		} else if thread.Name != threadName {
			t.Errorf("incorrect thread name. Got %s, wants %s", thread.Name, threadName)
		}

		_, _ = c.Channel(thread.ID).Delete()
	})

	var thread *Channel
	t.Run("create-thread-no-message", func(t *testing.T) {
		threadName := "Some Thread"
		var err error
		thread, err = c.Channel(guildAdmin.TextChannelGeneral).WithContext(deadline).CreateThreadNoMessage(&CreateThreadNoMessage{
			Name:                threadName,
			AutoArchiveDuration: AutoArchiveThreadDay,
			Type:                ChannelTypeGuildPublicThread,
		})
		if err != nil {
			panic(err)
		} else if thread == nil {
			t.Error(fmt.Errorf("fetched thread is nil: %w", err))
		} else if thread.Name != threadName {
			t.Errorf("incorrect thread name. Got %s, wants %s", thread.Name, threadName)
			t.Error(err)
		}
	})

	if thread == nil {
		t.Fatal("unable to create thread, cannot continue test suite")
	}

	t.Run("join", func(t *testing.T) {
		if err := c.Channel(thread.ID).WithContext(deadline).JoinThread(); err != nil {
			t.Fatal(fmt.Errorf("unable to join thread: %w", err))
		}
	})

	t.Run("leave", func(t *testing.T) {
		if err := c.Channel(thread.ID).WithContext(deadline).LeaveThread(); err != nil {
			t.Error(fmt.Errorf("unable to leave thread: %w", err))
		}
	})

	t.Run("add-member", func(t *testing.T) {
		if err := c.Channel(thread.ID).WithContext(deadline).AddThreadMember(andersfylling); err != nil {
			t.Fatal(fmt.Errorf("unable to add thread member: %w", err))
		}
	})

	t.Run("get-member", func(t *testing.T) {
		member, err := c.Channel(thread.ID).WithContext(deadline).GetThreadMember(andersfylling)
		if err != nil {
			t.Error(fmt.Errorf("unable to get thread member: %w", err))
		} else if member.ID != andersfylling {
			t.Error(fmt.Errorf("did not get correct thread member. Got %s, wants %s", member.ID, andersfylling))
		}
	})
	t.Run("get-members", func(t *testing.T) {
		members, err := c.Channel(thread.ID).WithContext(deadline).GetThreadMembers()
		if err != nil {
			t.Error(fmt.Errorf("unable to get thread member: %w", err))
		} else if len(members) == 1 {
			t.Error(fmt.Errorf("did not get correct number of thread members. Got %d got %d", len(members), 1))
		} else if members[0].ID != andersfylling {
			t.Error(fmt.Errorf("did not get correct thread member. Got %s, wants %s", members[0].ID, andersfylling))
		}
	})

	t.Run("remove-member", func(t *testing.T) {
		if err := c.Channel(thread.ID).WithContext(deadline).RemoveThreadMember(andersfylling); err != nil {
			t.Error(fmt.Errorf("unable to remove thread member. %w", err))
		}
	})

	t.Run("delete", func(t *testing.T) {
		if _, err := c.Channel(thread.ID).Delete(); err != nil {
			t.Error(fmt.Errorf("unable to delete thread: %w", err))
		}
	})
}
