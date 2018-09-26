package disgord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// ValidateJSONMarshalling
func validateJSONMarshalling(b []byte, v interface{}) error {
	var err error

	// convert to struct
	err = unmarshal(b, &v)
	if err != nil {
		return err
	}

	// back to json
	prettyJSON, err := json.MarshalIndent(&v, "", "    ")
	if err != nil {
		return err
	}

	// sort the data by keys
	// omg im getting lost in my own train of thought
	omg := make(map[string]interface{})
	err = unmarshal(prettyJSON, &omg)
	if err != nil {
		return err
	}

	omgAgain := make(map[string]interface{})
	err = unmarshal(b, &omgAgain)
	if err != nil {
		return err
	}

	// back to json
	prettyJSON, err = json.MarshalIndent(&omg, "", "    ")
	if err != nil {
		return err
	}

	b, err = json.MarshalIndent(&omgAgain, "", "    ")
	if err != nil {
		return err
	}

	// minify for comparison
	dst1 := bytes.Buffer{}
	err = json.Compact(&dst1, b)
	if err != nil {
		return err
	}
	dst2 := bytes.Buffer{}
	err = json.Compact(&dst2, prettyJSON)
	if err != nil {
		return err
	}

	// compare
	if dst2.String() != dst1.String() {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(b), string(prettyJSON), false)
		return fmt.Errorf("json data differs. \nDifference \n%s", dmp.DiffPrettyText(diffs))
	}

	return nil
}

func check(err error, t *testing.T) {
	// Hide function from stacktrace, PR#3
	t.Helper()

	if err != nil {
		t.Error(err)
	}
}

func TestError_InterfaceImplementations(t *testing.T) {
	var u interface{} = &ErrorUnsupportedType{}

	t.Run("error", func(t *testing.T) {
		if _, ok := u.(error); !ok {
			t.Error("ErrorUnsupportedType does not implement error")
		}
	})
}

// unmarshalling

func BenchmarkUnmarshalReflection(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/user/user1.json")
	if err != nil {
		b.Skip("missing file for benchmarking unmarshal")
		return
	}

	b.Run("using reflection", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var user *User
			unmarshal(data, user)
		}
	})

	b.Run("using interface wiring", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			user := &User{}
			var m = make(map[string]interface{})
			unmarshal(data, m)

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
			unmarshal(data, m)

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
