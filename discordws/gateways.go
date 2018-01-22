package discordws

import (
	"bytes"
	"encoding/gob"
)

type GatewayPayload struct {
	Op             uint        `json:"op"`
	Data           interface{} `json:"d"`
	SequenceNumber uint        `json:"s,omitempty"`
	EventName      string      `json:"t,omitempty"`
}

func (gp *GatewayPayload) DataToByteArr() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(gp.Data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type GetGatewayResponse struct {
	URL    string `json:"url"`
	Shards uint   `json:"shards,omitempty"`
}
