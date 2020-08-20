// +build !integration

package gateway

import (
	"github.com/andersfylling/disgord/json"
	"io/ioutil"
	"strconv"
	"sync"
	"testing"
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
		err := json.Unmarshal(file, &evt)
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
		err = json.Unmarshal(data, &evt)
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
			err := json.Unmarshal(file, evt)
			pool.Put(evt)
			if err != nil {
				t.Error(err)
			}
		}
	})
}
