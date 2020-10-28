// +build !integration

package disgord

import (
	"testing"

	"github.com/andersfylling/disgord/json"
)

func TestChannel_DeepCopy(t *testing.T) {
	test := &Channel{}
	icon1 := "sdljfdsjf"
	test.Icon = icon1
	test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
		Type: 0,
	})
	test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
		Type: 1,
	})

	cp := DeepCopy(test).(*Channel)
	icon2 := "sfkjdsf"
	test.Icon = icon2
	if cp.Icon != icon1 {
		t.Error("deep copy failed")
	}

	test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
		Type: 2,
	})
	if len(cp.PermissionOverwrites) != 2 {
		t.Error("deep copy failed")
	}
}

func TestChannel_JSONIconNull(t *testing.T) {
	// check if null's in json are parsed as an empty string
	data := []byte(`{"id":"324234235","type":1,"icon":null}`)
	var c *struct {
		ID   Snowflake `json:"id"`
		Type int       `json:"type"`
		Icon string    `json:"icon"`
	}
	if err := json.Unmarshal(data, &c); err != nil {
		t.Fatal(err)
	}
	executeInternalUpdater(c)

	if c.Icon != "" {
		t.Error(c.Icon, "was not empty")
	}
}
