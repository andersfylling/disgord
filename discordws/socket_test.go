package discordws

import (
	"fmt"
	"os"
	"testing"
)

func TestConnect(t *testing.T) {
	token := os.Getenv("DISGORD_TOKEN")
	if token == "" {
		token = "NDA0NzY4MzUxMjgyMzk3MTg1.DUapbA.9uF6sXXIiOs7NzWC-nYdBz6Oaos"
		// panic("Missing disgord token in env var: DISGORD_TOKEN")
	}
	d := NewRequiredClient(&Config{
		Token:        token,
		DAPIVersion:  6,
		DAPIEncoding: EncodingJSON,
	})
	fmt.Println(d.token)
	//err := d.Connect()
	//if err != nil {
	//	t.Error(err)
	//}
	//<-d.Connected
	//d.Disconnect()
	//d.Kill()
}
