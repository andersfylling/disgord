// +build integration

package test

import (
	"testing"
	"time"

	"github.com/andersfylling/disgord"
)

func TestOn(t *testing.T) {
	c := disgord.New(disgord.Config{
		BotToken:     "sdkjfhdksfhskdjfhdkfjsd",
		DisableCache: true,
	})

	t.Run("normal Session", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("should not have triggered a panic")
			}
		}()

		c.On(disgord.EvtChannelCreate, func(s disgord.Session, e *disgord.ChannelCreate) {})
	})

	t.Run("normal Session with ctrl", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("should not have triggered a panic")
			}
		}()

		c.On(disgord.EvtChannelCreate, func(s disgord.Session, e *disgord.ChannelCreate) {}, &disgord.Ctrl{Runs: 1})
	})

	t.Run("normal Session with multiple ctrl's", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("multiple controllers should trigger a panic")
			}
		}()

		c.On(disgord.EvtChannelCreate,
			func(s disgord.Session, e *disgord.ChannelCreate) {},
			&disgord.Ctrl{Runs: 1},
			&disgord.Ctrl{Until: time.Now().Add(1 * time.Minute)})
	})

	t.Run("Session pointer", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Did not panic on incorrect handler signature")
			}
		}()

		c.On(disgord.EvtChannelCreate, func(s *disgord.Session, e *disgord.ChannelCreate) {})
	})
}
