package websocket

import "testing"
import "os"
import "net/http"
import (
	"fmt"
)

const DISCORD_TOKEN = "DISCORD_DISGORDWS_TOKEN_TEST"
const DISCORD_TEST_MANUALLY = "DISCORD_DISGORDWS_MANUAL_TEST"

func createClient(t *testing.T) (DiscordWebsocket, error) {
	id := "DisgordWS v2.0.0"
	conf := &Config{
		Token:         os.Getenv(DISCORD_TOKEN),
		HTTPClient:    &http.Client{},
		DAPIVersion:   6,
		DAPIEncoding:  "json",
		Browser:       id,
		Device:        id,
		ChannelBuffer: 1,
		Debug:         true,
	}
	if conf.Token == "" {
		fmt.Printf("missing environment token '%s', skipping real connection tests.\n", DISCORD_TOKEN)
		t.Skip()
	}

	return NewClient(conf)
}

func TestClientConnection(t *testing.T) {
	if os.Getenv(DISCORD_TEST_MANUALLY) != "true" {
		fmt.Printf("missing environment token '%s', skipping manual connection test.\n", DISCORD_TEST_MANUALLY)
		t.Skip()
	}
	var err error

	client, err := createClient(t)
	if err != nil {
		t.Fatal(err)
	}
	err = client.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer func(client DiscordWebsocket) {
		select {
		case <-client.DiscordWSEventChan():
			client.Disconnect()
		}
	}(client)
}
