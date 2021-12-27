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

	// -------------------
	// CreateThread
	// -------------------
	t.Run("create-thread", func(t *testing.T) {
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
	})

	// -------------------
	// CreateThreadNoMessage
	// -------------------
	t.Run("create-thread-no-message", func(t *testing.T) {
		threadName := "HELLO WORLD2"
		thread, err := c.Channel(guildAdmin.TextChannelGeneral).WithContext(deadline).CreateThreadNoMessage(&CreateThreadNoMessage{
			Name:                threadName,
			AutoArchiveDuration: AutoArchiveThreadDay,
			Type:                ChannelTypeGuildPublicThread,
		})
		if err != nil {
			panic(err)
		} else if thread == nil {
			t.Error(fmt.Errorf("fetched thread is nil. %w", err))
		} else if thread.Name != threadName {
			t.Errorf("incorrect thread name. Got %s, wants %s", thread.Name, threadName)
			t.Error(err)
		}
	})

	// -------------------
	// JoinThread
	// -------------------
	t.Run("join-thread", func(t *testing.T) {
		threadName := "HELLO WORLD3"
		thread, err := c.Channel(guildAdmin.TextChannelGeneral).WithContext(deadline).CreateThreadNoMessage(&CreateThreadNoMessage{
			Name:                threadName,
			AutoArchiveDuration: AutoArchiveThreadDay,
			Type:                ChannelTypeGuildPublicThread,
		})
		if err != nil {
			panic(err)
		}
		err = c.Channel(thread.ID).WithContext(deadline).JoinThread()
		if err != nil {
			t.Error(fmt.Errorf("Unable to join thread. %w", err))
		}
	})

	// -------------------
	// LeaveThread
	// -------------------
	t.Run("leave-thread", func(t *testing.T) {
		threadName := "HELLO WORLD4"
		thread, err := c.Channel(guildAdmin.TextChannelGeneral).WithContext(deadline).CreateThreadNoMessage(&CreateThreadNoMessage{
			Name:                threadName,
			AutoArchiveDuration: AutoArchiveThreadDay,
			Type:                ChannelTypeGuildPublicThread,
		})
		if err != nil {
			panic(err)
		}
		err = c.Channel(thread.ID).WithContext(deadline).JoinThread()
		if err != nil {
			panic(err)
		}
		err = c.Channel(thread.ID).WithContext(deadline).LeaveThread()
		if err != nil {
			t.Error(fmt.Errorf("Unable to leave thread. %w", err))
		}
	})

	// -------------------
	// AddThreadMember
	// -------------------
	t.Run("add-thread-member", func(t *testing.T) {
		threadName := "HELLO WORLD5"
		thread, err := c.Channel(guildAdmin.TextChannelGeneral).WithContext(deadline).CreateThreadNoMessage(&CreateThreadNoMessage{
			Name:                threadName,
			AutoArchiveDuration: AutoArchiveThreadDay,
			Type:                ChannelTypeGuildPublicThread,
		})
		if err != nil {
			panic(err)
		}
		err = c.Channel(thread.ID).WithContext(deadline).AddThreadMember(andersfylling)
		if err != nil {
			t.Error(fmt.Errorf("Unable to add thread member. %w", err))
		}
	})

	// -------------------
	// RemoveThreadMember
	// -------------------
	t.Run("remove-thread-member", func(t *testing.T) {
		threadName := "HELLO WORLD6"
		thread, err := c.Channel(guildAdmin.TextChannelGeneral).WithContext(deadline).CreateThreadNoMessage(&CreateThreadNoMessage{
			Name:                threadName,
			AutoArchiveDuration: AutoArchiveThreadDay,
			Type:                ChannelTypeGuildPublicThread,
		})
		if err != nil {
			panic(err)
		}
		err = c.Channel(thread.ID).WithContext(deadline).AddThreadMember(andersfylling)
		if err != nil {
			panic(err)
		}
		err = c.Channel(thread.ID).WithContext(deadline).RemoveThreadMember(andersfylling)
		if err != nil {
			t.Error(fmt.Errorf("Unable to remove thread member. %w", err))
		}
	})

	// -------------------
	// GetThreadMember
	// -------------------
	t.Run("get-thread-member", func(t *testing.T) {
		threadName := "HELLO WORLD7"
		thread, err := c.Channel(guildAdmin.TextChannelGeneral).WithContext(deadline).CreateThreadNoMessage(&CreateThreadNoMessage{
			Name:                threadName,
			AutoArchiveDuration: AutoArchiveThreadDay,
			Type:                ChannelTypeGuildPublicThread,
		})
		if err != nil {
			panic(err)
		}
		err = c.Channel(thread.ID).WithContext(deadline).AddThreadMember(andersfylling)
		if err != nil {
			panic(err)
		}
		member, err := c.Channel(thread.ID).WithContext(deadline).GetThreadMember(andersfylling)
		if err != nil {
			t.Error(fmt.Errorf("Unable to get thread member. %w", err))
		} else if member.ID != andersfylling {
			t.Error(fmt.Errorf("Did not get correct thread member. Got %s, wants %s", member.ID, andersfylling))
		}
	})

	// -------------------
	// GetThreadMembers
	// -------------------
	t.Run("get-thread-members", func(t *testing.T) {
		threadName := "HELLO WORLD8"
		thread, err := c.Channel(guildAdmin.TextChannelGeneral).WithContext(deadline).CreateThreadNoMessage(&CreateThreadNoMessage{
			Name:                threadName,
			AutoArchiveDuration: AutoArchiveThreadDay,
			Type:                ChannelTypeGuildPublicThread,
		})
		if err != nil {
			panic(err)
		}
		err = c.Channel(thread.ID).WithContext(deadline).AddThreadMember(andersfylling)
		if err != nil {
			panic(err)
		}
		members, err := c.Channel(thread.ID).WithContext(deadline).GetThreadMembers()
		if err != nil {
			t.Error(fmt.Errorf("Unable to get thread member. %w", err))
		} else if len(members) == 1 {
			t.Error(fmt.Errorf("Did not get correct number of thread members. Got %d got %d.", len(members), 1))
		} else if members[0].ID != andersfylling {
			t.Error(fmt.Errorf("Did not get correct thread member. Got %s, wants %s", members[0].ID, andersfylling))
		}
	})
}
