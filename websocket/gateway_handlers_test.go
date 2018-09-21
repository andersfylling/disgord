package websocket

import (
	"io/ioutil"
	"strconv"
	"testing"
)

func getAllJSONFiles(t *testing.T) (files [][]byte) {
	for _, i := range []int{1,2,3,4} {
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



func TestDiscordEvent_CustomUnmarshaller(t *testing.T) {
	files := getAllJSONFiles(t)
	for _, file := range files {
		_, err, optimized := discordSocketUnmarshaller(file)
		if err != nil {
			t.Error(err)
		}
		if !optimized {
			t.Error("not opimized")
		}
	}

	t.Run("std as fallback", func(t *testing.T) {
		data, err := ioutil.ReadFile("testdata/diff-structure.json")
		if err != nil {
			t.Skip(err)
			return
		}

		evt, err, optimized := discordSocketUnmarshaller(data)
		if err != nil {
			t.Error(err)
		}
		if optimized || evt == nil {
			t.Error("was optimized")
			if evt == nil {
				t.Error("evt was nil")
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
		discordSocketUnmarshaller(data)
	}
}

func BenchmarkEvent_Unmarshal_smallJSON(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/small.json")
	if err != nil {
		return
	}
	for n := 0; n < b.N; n++ {
		evt := &gatewayEvent{}
		unmarshal(data, evt)
	}
}

func BenchmarkEvent_CustomUnmarshal_largeJSON(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/large.json")
	if err != nil {
		return
	}
	for n := 0; n < b.N; n++ {
		discordSocketUnmarshaller(data)
	}
}

func BenchmarkEvent_Unmarshal_largeJSON(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/large.json")
	if err != nil {
		return
	}
	for n := 0; n < b.N; n++ {
		evt := &gatewayEvent{}
		unmarshal(data, evt)
	}
}