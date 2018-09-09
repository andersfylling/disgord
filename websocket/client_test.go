package websocket

import "testing"
import "os"
import "net/http"
import (
	"fmt"
	"os/signal"
	"syscall"
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

func TestClientConnectionManually(t *testing.T) {
	if os.Getenv(DISCORD_TEST_MANUALLY) != "true" {
		fmt.Printf("missing environment token '%s', skipping manual connection test.\n", DISCORD_TEST_MANUALLY)
		t.Skip()
	}
	var err error

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	client, err := createClient(t)
	if err != nil {
		t.Fatal(err)
	}
	client.MockEventChanReciever()
	err = client.Connect()
	if err != nil {
		t.Fatal(err)
	}

	<-termSignal
	client.Disconnect()
}
