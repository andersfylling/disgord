package gateway

import (
	"testing"
)

func TestVoiceClient(t *testing.T) {
	conn := &testWS{
		closing: make(chan interface{}),
		opening: make(chan interface{}),
		writing: make(chan interface{}),
		reading: make(chan []byte),
	}
	conn.isConnected.Store(false)


	c, err := NewVoiceClient(&VoiceConfig{
		conn: conn,
		GuildID:           1,
		UserID:            0,
		SessionID:         "",
		Token:             "",
		HTTPClient:        nil,
		Endpoint:          "",
		MessageQueueLimit: 0,
		Logger:            nil,
		SystemShutdown:    nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	
	c.Connect()
}