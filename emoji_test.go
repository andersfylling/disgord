package disgord

import (
	"testing"
)

func TestEmoji_InterfaceImplementations(t *testing.T) {
	var c interface{} = &Emoji{}

	t.Run("DeepCopier", func(t *testing.T) {
		if _, ok := c.(DeepCopier); !ok {
			t.Error("Emoji does not implement DeepCopier")
		}
	})

	t.Run("Copier", func(t *testing.T) {
		if _, ok := c.(Copier); !ok {
			t.Error("Emoji does not implement Copier")
		}
	})
	//
	// t.Run("DiscordSaver", func(t *testing.T) {
	// 	if _, ok := c.(discordSaver); !ok {
	// 		t.Error("Emoji does not implement DiscordSaver")
	// 	}
	// })
	//
	// t.Run("discordDeleter", func(t *testing.T) {
	// 	if _, ok := c.(discordDeleter); !ok {
	// 		t.Error("Emoji does not implement discordDeleter")
	// 	}
	// })
}
