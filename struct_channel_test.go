package disgord

import "testing"

func TestChannel_InterfaceImplementations(t *testing.T) {
	var c interface{} = &Channel{}

	t.Run("DeepCopier", func(t *testing.T) {
		if _, ok := c.(DeepCopier); !ok {
			t.Error("Channel does not implement DeepCopier")
		}

		test := NewChannel()
		test.Icon = "sdljfdsjf"
		test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
			Type: "first",
		})
		test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
			Type: "second",
		})

		copy := test.DeepCopy().(*Channel)
		test.Icon = "sfkjdsf"
		if copy.Icon != "sdljfdsjf" {
			t.Error("deep copy failed")
		}

		test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
			Type: "third",
		})
		if len(copy.PermissionOverwrites) != 2 {
			t.Error("deep copy failed")
		}
	})

	t.Run("Copier", func(t *testing.T) {
		if _, ok := c.(Copier); !ok {
			t.Error("Channel does not implement Copier")
		}
	})

	t.Run("DiscordSaver", func(t *testing.T) {
		if _, ok := c.(DiscordSaver); !ok {
			t.Error("Channel does not implement DiscordSaver")
		}
	})

	t.Run("DiscordDeleter", func(t *testing.T) {
		if _, ok := c.(DiscordDeleter); !ok {
			t.Error("Channel does not implement DiscordDeleter")
		}
	})
}
