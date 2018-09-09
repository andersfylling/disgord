package resource

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	. "github.com/andersfylling/snowflake"
)

func BenchmarkUnmarshalReflection(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/user/user1.json")
	if err != nil {
		b.Skip("missing file for benchmarking unmarshal")
		return
	}

	b.Run("using reflection", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var user *User
			json.Unmarshal(data, user)
		}
	})

	b.Run("using interface wiring", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			user := &User{}
			var m = make(map[string]interface{})
			json.Unmarshal(data, m)

			var v interface{}
			var ok bool

			if v, ok = m["id"]; ok {
				user.ID = v.(Snowflake)
			}
			if v, ok = m["username"]; ok {
				user.Username = v.(string)
			}
			if v, ok = m["discriminator"]; ok {
				user.Discriminator = v.(string)
			}
			if v, ok = m["email"]; ok {
				user.Email = v.(string)
			}
			if v, ok = m["avatar"]; ok {
				user.Avatar = v.(*string)
			}
			if v, ok = m["token"]; ok {
				user.Token = v.(string)
			}
			if v, ok = m["verified"]; ok {
				user.Verified = v.(bool)
			}
			if v, ok = m["mfa_enabled"]; ok {
				user.MFAEnabled = v.(bool)
			}
			if v, ok = m["bot"]; ok {
				user.Bot = v.(bool)
			}

		}
	})

	b.Run("using string wiring", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			user := &User{}
			var m = make(map[string]string)
			json.Unmarshal(data, m)

			var v string
			var ok bool

			if v, ok = m["id"]; ok {
				user.ID = ParseSnowflakeString(v)
			}
			if v, ok = m["username"]; ok {
				user.Username = v
			}
			if v, ok = m["discriminator"]; ok {
				user.Discriminator = v
			}
			if v, ok = m["email"]; ok {
				user.Email = v
			}
			if vv, ok := m["avatar"]; ok {
				user.Avatar = &vv
			}
			if v, ok = m["token"]; ok {
				user.Token = v
			}
			if v, ok = m["verified"]; ok {
				user.Verified = v == "true"
			}
			if v, ok = m["mfa_enabled"]; ok {
				user.MFAEnabled = v == "true"
			}
			if v, ok = m["bot"]; ok {
				user.Bot = v == "true"
			}

		}
	})
}
