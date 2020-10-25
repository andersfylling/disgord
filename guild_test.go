// +build !integration

package disgord

import (
	"strconv"
	"testing"

	"github.com/andersfylling/disgord/json"
)

// NewGuild ...
func NewGuild() *Guild {
	return &Guild{
		Roles:       []*Role{},
		Emojis:      []*Emoji{},
		Features:    []string{},
		VoiceStates: []*VoiceState{},
		Members:     []*Member{},
		Channels:    []*Channel{},
		Presences:   []*UserPresence{},
	}
}

func TestGuild_ChannelSorting(t *testing.T) {
	g := &Guild{}
	total := 1000
	for i := total; i > 0; i-- {
		s := NewSnowflake(uint64(i))
		c := &Channel{ID: s}
		g.AddChannel(c)
	}

	chans := g.Channels
	for i := 1; i <= total; i++ {
		if chans[i-1].ID != NewSnowflake(uint64(i)) {
			t.Error("wrong order")
			break
		}
	}
}

// --------
func TestGuildEmbed(t *testing.T) {
	res := []byte("{\"enabled\":true,\"channel_id\":\"41771983444115456\"}")
	expects := []byte("{\"enabled\":true,\"channel_id\":41771983444115456}")

	// convert to struct
	guildEmbed := GuildEmbed{}
	if err := json.Unmarshal(res, &guildEmbed); err != nil {
		t.Error(err)
	}

	// back to json
	data, err := json.Marshal(&guildEmbed)
	if err != nil {
		t.Error(err)
	}

	// match
	if string(expects) != string(data) {
		t.Errorf("json data differs. Got %s, wants %s", string(data), string(expects))
	}
}

// -------------

func TestGuild_sortChannels(t *testing.T) {
	snowflakes := []Snowflake{
		NewSnowflake(6),
		NewSnowflake(65),
		NewSnowflake(324),
		NewSnowflake(5435),
		NewSnowflake(63453),
		NewSnowflake(111111111),
	}

	guild := NewGuild()

	for i := range snowflakes {
		channel := NewChannel()
		channel.ID = snowflakes[len(snowflakes)-1-i] // reverse

		guild.Channels = append(guild.Channels, channel)
	}

	guild.sortChannels()
	for i, c := range guild.Channels {
		if snowflakes[i] != c.ID {
			t.Error("Channels in guild did not sort correctly")
		}
	}
}

func TestGuild_AddChannel(t *testing.T) {
	snowflakes := []Snowflake{
		NewSnowflake(6),
		NewSnowflake(65),
		NewSnowflake(324),
		NewSnowflake(5435),
		NewSnowflake(63453),
		NewSnowflake(111111111),
	}

	guild := NewGuild()

	for i := range snowflakes {
		channel := NewChannel()
		channel.ID = snowflakes[len(snowflakes)-1-i] // reverse

		guild.AddChannel(channel)
	}

	for i, c := range guild.Channels {
		if snowflakes[i] != c.ID {
			t.Error("Channels in guild did not sort correctly")
		}
	}
}

func TestGuild_DeleteChannel(t *testing.T) {
	snowflakes := []Snowflake{
		NewSnowflake(6),
		NewSnowflake(65),
		NewSnowflake(324),
		NewSnowflake(5435),
		NewSnowflake(63453),
		NewSnowflake(111111111),
	}

	guild := NewGuild()

	for i := range snowflakes {
		channel := NewChannel()
		channel.ID = snowflakes[len(snowflakes)-1-i] // reverse

		guild.AddChannel(channel)
	}

	id := snowflakes[3]
	channel := NewChannel()
	channel.ID = id
	guild.DeleteChannel(channel)
	_, err := guild.Channel(id)
	if err == nil {
		t.Error("no error given when requesting a deleted channel")
	}
}

func TestPermissionBit(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		testBits := PermissionSendMessages | PermissionReadMessages
		if testBits.Contains(PermissionAdministrator) {
			t.Fatal("does not have administrator")
		}
		if !testBits.Contains(PermissionSendMessages) {
			t.Fatal("does have send messages")
		}
		if !testBits.Contains(PermissionReadMessages) {
			t.Fatal("does have read messages")
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		t.Run("single", func(t *testing.T) {
			container := struct {
				Permission PermissionBit `json:"permission"`
			}{PermissionSendMessages | PermissionReadMessages}

			b, err := json.Marshal(&container)
			if err != nil {
				t.Fatal(err)
			}

			tmp := container
			tmp.Permission = 0
			if err := json.Unmarshal(b, &tmp); err != nil {
				t.Fatal(err)
			}

			if tmp.Permission != container.Permission {
				t.Fatalf("unmarshaled value was unexpected. Got %d, wants %d", tmp.Permission, container.Permission)
			}
		})
		t.Run("array", func(t *testing.T) {
			perms := []PermissionBit{
				PermissionSendMessages | PermissionReadMessages,
				PermissionAddReactions,
				PermissionBanMembers,
			}
			container := struct {
				Permissions []PermissionBit `json:"permissions"`
			}{perms}

			contains := func(v PermissionBit) bool {
				for _, p := range container.Permissions {
					if p == v {
						return true
					}
				}
				return false
			}

			b, err := json.Marshal(&container)
			if err != nil {
				t.Fatal(err)
			}

			tmp := container
			tmp.Permissions = nil
			if err := json.Unmarshal(b, &tmp); err != nil {
				t.Fatal(err)
			}

			for i := range perms {
				if !contains(tmp.Permissions[i]) {
					t.Errorf("unmarshaled value was not found in original. Got %d", tmp.Permissions[i])
				}
			}
		})
		t.Run("array-extra", func(t *testing.T) {
			b := []byte(`{"permissions":["123", "4567"]}`)
			container := struct {
				Permissions []PermissionBit `json:"permissions"`
			}{}

			contains := func(v PermissionBit) bool {
				for _, p := range container.Permissions {
					if p == v {
						return true
					}
				}
				return false
			}

			if err := json.Unmarshal(b, &container); err != nil {
				t.Fatal(err)
			}

			if !contains(123) {
				t.Error("missing permission value 123")
			}
			if !contains(4567) {
				t.Error("missing permission value 4567")
			}
		})
	})

	t.Run("marshal", func(t *testing.T) {
		expects := PermissionBit(123456789)
		data := []byte(`{"permission":"` + strconv.FormatUint(uint64(expects), 10) + `"}`)
		container := struct {
			Permission PermissionBit `json:"permission"`
		}{}

		if container.Permission != 0 {
			t.Fatal("expected 0")
		}

		if err := json.Unmarshal(data, &container); err != nil {
			t.Fatal(err)
		}

		if container.Permission != expects {
			t.Fatalf("unmarshaled value was unexpected. Got %d, wants %d", container.Permission, expects)
		}
	})
}
