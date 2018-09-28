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
				user.Discriminator = v.(Discriminator)
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
				d, _ := NewDiscriminator(v)
				user.Discriminator = d
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

func TestDiscriminator(t *testing.T) {
	t.Run("String()", func(t *testing.T) {
		var d Discriminator
		d = Discriminator(0)
		if d.String() != "" {
			t.Errorf("got %s, wants \"\"", d.String())
		}

		d = Discriminator(1)
		if d.String() != "0001" {
			t.Errorf("got %s, wants \"0001\"", d.String())
		}

		d = Discriminator(4)
		if d.String() != "0004" {
			t.Errorf("got %s, wants \"0004\"", d.String())
		}

		d = Discriminator(12)
		if d.String() != "0012" {
			t.Errorf("got %s, wants \"0012\"", d.String())
		}

		d = Discriminator(120)
		if d.String() != "0120" {
			t.Errorf("got %s, wants \"0120\"", d.String())
		}

		d = Discriminator(1201)
		if d.String() != "1201" {
			t.Errorf("got %s, wants \"1201\"", d.String())
		}
	})
	t.Run("UnmarshalJSON(...)", func(t *testing.T) {
		var data []byte
		var d Discriminator
		var err error

		data = []byte{
			'"', '0', '0', '0', '1', '"',
		}
		err = unmarshal(data, &d)
		if err != nil {
			t.Error(err)
		}
		if d.String() != "0001" {
			t.Errorf("got %s, wants \"0001\"", d.String())
		}

		data = []byte{
			'"', '0', '2', '0', '1', '"',
		}
		err = unmarshal(data, &d)
		if err != nil {
			t.Error(err)
		}
		if d.String() != "0201" {
			t.Errorf("got %s, wants \"0201\"", d.String())
		}

		data = []byte{
			'"', '"',
		}
		err = unmarshal(data, &d)
		if err != nil {
			t.Error(err)
		}
		if d.String() != "" {
			t.Errorf("got %s, wants \"\"", d.String())
		}
		if !d.NotSet() {
			t.Error("expected Discriminator to be NotSet")
		}

	})
	t.Run("MarshalJSON()", func(t *testing.T) {
		d := Discriminator(34)
		data, err := json.Marshal(&d)
		if err != nil {
			t.Error(err)
		}
		if string(data) != string([]byte("\"0034\"")) {
			t.Errorf("wrong. got %s, expects \"0034\"", string(data))
		}

		d = Discriminator(0)
		data, err = json.Marshal(&d)
		if err != nil {
			t.Error(err)
		}
		if string(data) != string([]byte("\"\"")) {
			t.Errorf("wrong. got %s, expects \"\"", string(data))
		}
	})
	t.Run("NotSet()", func(t *testing.T) {
		d := Discriminator(34)
		if d.NotSet() {
			t.Error("expected Discriminator.NotSet to be false, got true")
		}

		d = Discriminator(0)
		if !d.NotSet() {
			t.Error("expected Discriminator.NotSet to be true, got false")
		}
	})
}
