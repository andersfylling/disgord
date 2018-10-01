package disgord

//
// import (
// 	//"encoding/json"
// 	"errors"
//
// 	"io/ioutil"
// 	"testing"
// )
//
// func extractRootID(data []byte) (id Snowflake, err error) {
// 	filter := []byte(`"id":"`)
// 	filterLen := len(filter) - 1
// 	scope := 0
//
// 	var start uint
// 	lastPos := len(data) - 1
// 	for i := 1; i <= lastPos-filterLen; i++ {
// 		if data[i] == '{' {
// 			scope++
// 		} else if data[i] == '}' {
// 			scope--
// 		}
//
// 		if scope != 0 {
// 			continue
// 		}
//
// 		for j := filterLen; j >= 0; j-- {
// 			if filter[j] != data[i+j] {
// 				break
// 			}
//
// 			if j == 0 {
// 				start = uint(i + len(filter))
// 			}
// 		}
//
// 		if start != 0 {
// 			break
// 		}
// 	}
//
// 	if start == 0 {
// 		err = errors.New("unable to locate ID")
// 		return
// 	}
//
// 	i := start
// E:
// 	for {
// 		switch data[i] {
// 		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
// 			i++
// 		default:
// 			break E
// 		}
// 	}
//
// 	if i > start {
// 		id = Snowflake(0)
// 		err = id.UnmarshalJSON(data[start-1 : i+1])
// 	} else {
// 		err = errors.New("id was empty")
// 	}
// 	return
// }
//
// func allocateGuild() *Guild {
// 	return new(Guild)
// }
//
// func TestIDExtraction(t *testing.T) {
// 	data := []byte(`{"id":"80351110224678912","test":{},username":"Nelly","discriminator":"1337","email":"nelly@discordapp.com","avatar":"8342729096ea3675442027381ff50dfe","verified":true}`)
// 	id, err := extractRootID(data)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	if id != Snowflake(80351110224678912) {
// 		t.Error("incorrect snowflake id")
// 	}
//
// 	data, err = ioutil.ReadFile("testdata/guild/complete-guild.json")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	id, err = extractRootID(data)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	if id != Snowflake(244200618854580224) {
// 		t.Error("incorrect snowflake id")
// 	}
//
// }
//
// var sink *Guild
//
// func BenchmarkUnmarshal(b *testing.B) {
// 	data, err := ioutil.ReadFile("testdata/guild/complete-guild.json")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	b.Run("cache", func(b *testing.B) {
// 		var cache *Guild = allocateGuild()
// 		for i := 0; i < b.N; i++ {
// 			id, err := extractRootID(data)
// 			if err != nil {
// 				panic(err)
// 			}
//
// 			if id == Snowflake(244200618854580224) {
// 				err = unmarshal(data, cache)
// 				sink = cache
// 				if err != nil {
// 					panic(err)
// 				}
// 			} else {
// 				panic("wrong id")
// 			}
// 		}
// 	})
// 	b.Run("new obj", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			var guild *Guild = allocateGuild()
// 			err = unmarshal(data, guild)
// 			if err != nil {
// 				panic(err)
// 			}
// 			sink = guild
// 			guild = nil
// 		}
// 	})
// }
//
// func BenchmarkUnmarshal_A(b *testing.B) {
// 	data, err := ioutil.ReadFile("testdata/guild/complete-guild.json")
// 	if err != nil {
// 		panic(err)
// 	}
// 	var cache *Guild = allocateGuild()
// 	for i := 0; i < b.N; i++ {
// 		A(data, cache)
// 	}
// }
// func BenchmarkUnmarshal_B(b *testing.B) {
// 	data, err := ioutil.ReadFile("testdata/guild/complete-guild.json")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	for i := 0; i < b.N; i++ {
// 		sink = B(data)
// 	}
// }
//
// func A(data []byte, cache *Guild) *Guild {
// 	id, err := extractRootID(data)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	if id == Snowflake(244200618854580224) {
// 		err = unmarshal(data, cache)
// 		sink = cache
// 		if err != nil {
// 			panic(err)
// 		}
// 	} else {
// 		panic("wrong id")
// 	}
//
// 	return cache
// }
//
// func B(data []byte) *Guild {
// 	var guild *Guild = allocateGuild()
// 	err := unmarshal(data, guild)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return guild
// }
