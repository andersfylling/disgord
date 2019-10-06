// +build disgordperf

package websocket

import (
	"io/ioutil"
	"testing"

	httd2 "github.com/andersfylling/disgord/internal/httd"
)

func BenchmarkEvent_CustomUnmarshal_smallJSON(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/small.json")
	if err != nil {
		return
	}
	for n := 0; n < b.N; n++ {
		evt := DiscordPacket{}
		if err := evt.UnmarshalJSON(data); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEvent_Unmarshal_smallJSON(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/small.json")
	if err != nil {
		return
	}
	for n := 0; n < b.N; n++ {
		evt := DiscordPacket{}
		if err := httd2.Unmarshal(data, &evt); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEvent_CustomUnmarshal_largeJSON(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/large.json")
	if err != nil {
		return
	}
	for n := 0; n < b.N; n++ {
		evt := DiscordPacket{}
		if err := evt.UnmarshalJSON(data); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEvent_Unmarshal_largeJSON(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/large.json")
	if err != nil {
		return
	}
	for n := 0; n < b.N; n++ {
		evt := DiscordPacket{}
		if err := httd2.Unmarshal(data, &evt); err != nil {
			b.Fatal(err)
		}
	}
}
