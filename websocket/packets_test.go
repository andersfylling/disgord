package websocket

import (
	"io/ioutil"
	"strconv"
	"sync"
	"testing"

	"github.com/andersfylling/disgord/httd"
)

func getAllJSONFiles(t *testing.T) (files [][]byte) {
	for _, i := range []int{1, 2, 3, 4} {
		data, err := ioutil.ReadFile("testdata/" + strconv.Itoa(i) + ".json")
		if err != nil {
			t.Error(err)
			break
		}
		files = append(files, data)
	}

	data, err := ioutil.ReadFile("testdata/large.json")
	if err != nil {
		t.Error(err)
		return
	}
	files = append(files, data)

	return
}

func TestDiscordPacket_UnmarshalJSON(t *testing.T) {
	files := getAllJSONFiles(t)
	for _, file := range files {
		evt := DiscordPacket{}
		err := httd.Unmarshal(file, &evt)
		if err != nil {
			t.Error(err)
		}
	}

	t.Run("std as fallback", func(t *testing.T) {
		data, err := ioutil.ReadFile("testdata/diff-structure.json")
		if err != nil {
			t.Skip(err)
			return
		}

		evt := DiscordPacket{}
		err = httd.Unmarshal(data, &evt)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("reusing data", func(t *testing.T) {
		pool := sync.Pool{
			New: func() interface{} {
				return &DiscordPacket{}
			},
		}
		files := getAllJSONFiles(t)
		for _, file := range files {
			evt := pool.Get().(*DiscordPacket)
			evt.reset()
			err := httd.Unmarshal(file, evt)
			pool.Put(evt)
			if err != nil {
				t.Error(err)
			}
		}
	})
}

func BenchmarkEvent_CustomUnmarshal_smallJSON(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/small.json")
	if err != nil {
		return
	}
	for n := 0; n < b.N; n++ {
		evt := DiscordPacket{}
		// evt.UnmarshalJSON(data) custom unmarshal
		_ = httd.Unmarshal(data, evt) // json.RawMessage
	}
}

func BenchmarkEvent_Unmarshal_smallJSON(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/small.json")
	if err != nil {
		return
	}
	for n := 0; n < b.N; n++ {
		evt := DiscordPacket{}
		// evt.UnmarshalJSON(data) custom unmarshal
		_ = httd.Unmarshal(data, &evt) // json.RawMessage
	}
}

func BenchmarkEvent_CustomUnmarshal_largeJSON(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/large.json")
	if err != nil {
		return
	}
	for n := 0; n < b.N; n++ {
		evt := DiscordPacket{}
		// evt.UnmarshalJSON(data) custom unmarshal
		_ = httd.Unmarshal(data, &evt) // json.RawMessage
	}
}

func BenchmarkEvent_Unmarshal_largeJSON(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/large.json")
	if err != nil {
		return
	}
	for n := 0; n < b.N; n++ {
		evt := DiscordPacket{}
		// evt.UnmarshalJSON(data) custom unmarshal
		_ = httd.Unmarshal(data, &evt) // json.RawMessage
	}
}
