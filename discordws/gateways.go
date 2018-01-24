package discordws

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

type EventInterface interface {
	Name() string
	Data() []byte
	Unmarshal(interface{}) error
}

type Event struct {
	content *gatewayEvent
}

func (evt *Event) Name() string {
	return evt.content.EventName
}

func (evt *Event) Data() []byte {
	return evt.content.Data.ByteArr()
}

func (evt *Event) Unmarshal(v interface{}) error {
	return json.Unmarshal(evt.Data(), v)
}

type gatewayPayload struct {
	Op             uint        `json:"op"`
	Data           interface{} `json:"d"`
	SequenceNumber uint        `json:"s,omitempty"`
	EventName      string      `json:"t,omitempty"`
}

// GatewayEvent used for incoming events from the gateway..
type gatewayEvent struct {
	Op             uint        `json:"op"`
	Data           payloadData `json:"d"`
	SequenceNumber uint        `json:"s"`
	EventName      string      `json:"t"`
}

type payloadData []byte

func (pd *payloadData) UnmarshalJSON(data []byte) error {
	*pd = payloadData(data)
	return nil
}

func (pd *payloadData) ByteArr() []byte {
	return []byte(*pd)
}

func (gp *gatewayPayload) DataToByteArr() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(gp.Data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type getGatewayResponse struct {
	URL    string `json:"url"`
	Shards uint   `json:"shards,omitempty"`
}

type helloPacket struct {
	HeartbeatInterval uint     `json:"heartbeat_interval"`
	Trace             []string `json:"_trace"`
}

type readyPacket struct {
	SessionID string   `json:"session_id"`
	Trace     []string `json:"_trace"`
}
