package websocket

import (
	"github.com/andersfylling/disgord/constant"
	"net/http"
	"testing"
	"time"
)

func newTestClient() (*testClient, error) {
	return &testClient{
		receiveChan: make(chan *discordPacket),
		emitChan:    make(chan *clientPacket),
	}, nil
}

type testClient struct {
	receiveChan  chan *discordPacket
	emitChan     chan *clientPacket
	connected    int
	disconnected int
}

func (c *testClient) Connect() error {
	c.connected++
	return nil
}
func (c *testClient) Disconnect() error {
	c.disconnected++
	return nil
}
func (c *testClient) Emit(cmd string, data interface{}) error {
	return nil
}
func (c *testClient) Receive() <-chan *discordPacket {
	return nil
}

func TestManager_RegisterEvent(t *testing.T) {
	m := Manager{}
	t1 := "test"
	m.RegisterEvent(t1)

	if len(m.trackedEvents) == 0 {
		t.Error("expected length to be 1, got 0")
	}

	m.RegisterEvent(t1)
	if len(m.trackedEvents) == 2 {
		t.Error("expected length to be 1, got 2")
	}
}

func TestManager_RemoveEvent(t *testing.T) {
	m := Manager{}
	t1 := "test"
	m.RegisterEvent(t1)

	if len(m.trackedEvents) == 0 {
		t.Error("expected length to be 1, got 0")
	}

	m.RemoveEvent("sdfsdf")
	if len(m.trackedEvents) == 0 {
		t.Error("expected length to be 1, got 0")
	}

	m.RemoveEvent(t1)
	if len(m.trackedEvents) == 1 {
		t.Error("expected length to be 0, got 1")
	}
}

func TestManager_reconnect(t *testing.T) {
	client, _ := newTestClient()
	m := &Manager{
		conf: &ManagerConfig{
			// identity
			Browser:             "disgord",
			Device:              "disgord",
			GuildLargeThreshold: 250,
			ShardID:             0,
			ShardCount:          0,

			DefaultClientConfig: DefaultClientConfig{
				// lib specific
				Version:       constant.DiscordVersion,
				Encoding:      constant.JSONEncoding,
				ChannelBuffer: 1,
				Endpoint:      "sfkjsdlfsf",

				// user settings
				Token: "sifhsdoifhsdifhsdf",
				HTTPClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
		},
		shutdown:  make(chan interface{}),
		restart:   make(chan interface{}),
		Client:    client,
		eventChan: make(chan *Event),
	}
	//go m.operationHandlers()

	m.Connect()
	if client.connected != 1 {
		t.Error("expected 1 connect, got 0")
	}

}
