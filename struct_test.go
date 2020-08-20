// +build !integration

package disgord

import (
	"github.com/andersfylling/disgord/json"
	"io/ioutil"
	"strconv"
	"testing"
)

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

func TestTime(t *testing.T) {
	t.Run("omitempty", func(t *testing.T) {
		b := struct {
			T Time `json:"time,omitempty"`
		}{}

		bBytes, err := defaultMarshaler(b)
		if err != nil {
			t.Fatal(err)
		}

		if string(bBytes) != `{"time":""}` {
			t.Errorf("did not get an 'omitted' field. Got %s", string(bBytes))
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
		if err = json.Unmarshal(data, &d); err != nil {
			t.Error(err)
		}
		executeInternalUpdater(d)

		if d.String() != "0001" {
			t.Errorf("got %s, wants \"0001\"", d.String())
		}

		data = []byte{
			'"', '0', '2', '0', '1', '"',
		}
		err = json.Unmarshal(data, &d)
		if err != nil {
			t.Error(err)
		}
		if d.String() != "0201" {
			t.Errorf("got %s, wants \"0201\"", d.String())
		}

		data = []byte{
			'"', '"',
		}
		if err = json.Unmarshal(data, &d); err != nil {
			t.Error(err)
		}
		executeInternalUpdater(d)

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

func BenchmarkDiscriminator(b *testing.B) {
	b.Run("comparison", func(b *testing.B) {
		b.Run("string", func(b *testing.B) {
			val := "0401"
			vals := []string{
				"0401", "0001", "0400", "0011", "0101", "5435", "0010",
			}
			var i int
			var match bool
			for n := 0; n < b.N; n++ {
				for i = range vals {
					match = val == vals[i]
				}
			}
			if match {
			}
		})
		b.Run("uint16", func(b *testing.B) {
			val := Discriminator(0401)
			vals := []Discriminator{
				0401, 0001, 0400, 0011, 0101, 5435, 0010,
			}
			var i int
			var match bool
			for n := 0; n < b.N; n++ {
				for i = range vals {
					match = val == vals[i]
				}
			}
			if match {
			}
		})
	})
	b.Run("Unmarshal", func(b *testing.B) {
		dataSets := [][]byte{
			[]byte("\"0001\""),
			[]byte("\"0452\""),
			[]byte("\"4342\""),
			[]byte("\"0100\""),
			[]byte("\"5100\""),
			[]byte("\"1000\""),
			[]byte("\"5129\""),
			[]byte("\"3020\""),
			[]byte("\"0010\""),
		}
		b.Run("string", func(b *testing.B) {
			var result string
			var i int
			length := len(dataSets)
			for n := 0; n < b.N; n++ {
				result = string(dataSets[i])
				if i == length {
					i = 0
				}
			}
			if result == "" {
			}
		})
		b.Run("uint16-a", func(b *testing.B) {
			var result uint16
			var i int
			lengthi := len(dataSets)
			for n := 0; n < b.N; n++ {
				data := dataSets[i]
				result = 0
				length := len(data) - 1
				for j := 1; j < length; j++ {
					result = result*10 + uint16(data[j]-'0')
				}
				if i == lengthi {
					i = 0
				}
			}
			if result == 0 {
			}
		})
		b.Run("uint16-b", func(b *testing.B) {
			var result uint16
			var i int
			length := len(dataSets)
			for n := 0; n < b.N; n++ {
				data := dataSets[i]
				var tmp uint64
				tmp, _ = strconv.ParseUint(string(data), 10, 16)
				result = uint16(tmp)
				if i == length {
					i = 0
				}
			}
			if result == 0 {
			}
		})
		type fooOld struct {
			Foo string `json:"discriminator,omitempty"`
		}
		type fooNew struct {
			Foo Discriminator `json:"discriminator,omitempty"`
		}
		b.Run("string-struct", func(b *testing.B) {
			foo := &fooOld{}
			var i int
			length := len(dataSets)
			for n := 0; n < b.N; n++ {
				_ = json.Unmarshal(dataSets[i], foo)
				executeInternalUpdater(foo)
				if i == length {
					i = 0
				}
			}
		})
		b.Run("uint16-struct", func(b *testing.B) {
			foo := &fooNew{}
			var i int
			length := len(dataSets)
			for n := 0; n < b.N; n++ {
				_ = json.Unmarshal(dataSets[i], foo)
				executeInternalUpdater(foo)
				if i == length {
					i = 0
				}
			}
		})
	})
}

func TestIDExtraction(t *testing.T) {
	data := []byte(`{"id":"80351110224678912","test":{},username":"Nelly","discriminator":"1337","email":"nelly@discordapp.com","avatar":"8342729096ea3675442027381ff50dfe","verified":true}`)
	id, err := extractAttribute([]byte(`"id":"`), 0, data)
	if err != nil {
		t.Error(err)
	}

	if id != Snowflake(80351110224678912) {
		t.Error("incorrect snowflake id")
	}

	data, err = ioutil.ReadFile("testdata/guild/complete-guild.json")
	if err != nil {
		t.Error(err)
	}
	id, err = extractAttribute([]byte(`"id":"`), 0, data)
	if err != nil {
		t.Error(err)
	}

	if id != Snowflake(244200618854580224) {
		t.Error("incorrect snowflake id")
	}

}
