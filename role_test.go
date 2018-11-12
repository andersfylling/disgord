package disgord

import "testing"

func TestRole_InterfaceImplementations(t *testing.T) {
	var c interface{} = &Role{}

	t.Run("DeepCopier", func(t *testing.T) {
		if _, ok := c.(DeepCopier); !ok {
			t.Error("Role does not implement DeepCopier")
		}
	})

	t.Run("Copier", func(t *testing.T) {
		if _, ok := c.(Copier); !ok {
			t.Error("Role does not implement Copier")
		}
	})

	t.Run("discordSaver", func(t *testing.T) {
		if _, ok := c.(discordSaver); !ok {
			t.Error("Role does not implement discordSaver")
		}
	})

	t.Run("discordDeleter", func(t *testing.T) {
		if _, ok := c.(discordDeleter); !ok {
			t.Error("Role does not implement discordDeleter")
		}
	})
}
