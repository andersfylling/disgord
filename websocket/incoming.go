package websocket

import (
	"strconv"
	"strings"
)

type payloadData []byte

func (pd *payloadData) UnmarshalJSON(data []byte) error {
	*pd = payloadData(data)
	return nil
}

func (pd *payloadData) ByteArr() []byte {
	return []byte(*pd)
}

type gatewayEventJSON struct {
	Op             uint        `json:"op"`
	Data           payloadData `json:"d"`
	SequenceNumber uint        `json:"s"`
	EventName      string      `json:"t"`
}

func (ge *gatewayEventJSON) gatewayEvent() *gatewayEvent {
	return &gatewayEvent{
		Op:             ge.Op,
		Data:           ge.Data,
		SequenceNumber: ge.SequenceNumber,
		EventName:      ge.EventName,
	}
}

type gatewayEvent struct {
	Op             uint        `json:"op"`
	Data           payloadData `json:"d"`
	SequenceNumber uint        `json:"s"`
	EventName      string      `json:"t"`
}

func (ge *gatewayEvent) GetOperationCode() uint {
	return ge.Op
}

func (ge *gatewayEvent) UnmarshalJSON(data []byte) (err error) {
	var i int

	// t
	t := []byte{
		'{', '"', 't', '"', ':',
	}
	for i = range t {
		if t[i] != data[i] {
			evt := &gatewayEventJSON{}
			err = unmarshal(data, evt)
			*ge = *evt.gatewayEvent()
			return
		}
	}
	i += 1              // jump to next char
	if data[i] == 'n' { // null
		i += 4 // skip `null`
	} else {
		// extract the t value
		var val strings.Builder
		i += 1 // skip `"`
		for ; data[i] != '"'; i++ {
			val.WriteByte(data[i])
		}
		i += 1 // skip `"`
		ge.EventName = val.String()
	}

	// s
	i += 2 // skip `,"`
	if data[i] != 's' {
		evt := &gatewayEventJSON{}
		err = unmarshal(data, evt)
		*ge = *evt.gatewayEvent()
		return
	}
	i += 3              // skip `s":`
	if data[i] == 'n' { // null value
		i += 4 // skip `null`
	} else {
		// extract the s value
		var val strings.Builder
		for ; data[i] != ','; i++ {
			val.WriteByte(data[i])
		}
		var tmp uint64
		tmp, err = strconv.ParseUint(val.String(), 10, 64)
		if err != nil {
			evt := &gatewayEventJSON{}
			err = unmarshal(data, evt)
			*ge = *evt.gatewayEvent()
			return
		}
		ge.SequenceNumber = uint(tmp)
	}

	// op
	i += 2              // skip `,"`
	if data[i] != 'o' { // o as in op
		evt := gatewayEventJSON{}
		err = unmarshal(data, &evt)
		*ge = *evt.gatewayEvent()
		return
	}
	i += 4              // skip `op":`
	if data[i] == 'n' { // null value
		i += 4 // skip `null`
	} else if data[i] == '0' {
		ge.Op = 0
		i += 1 // skip 0
	} else {
		// extract the op value
		var val strings.Builder
		for ; data[i] != ','; i++ {
			val.WriteByte(data[i])
		}
		var tmp uint64
		tmp, err = strconv.ParseUint(val.String(), 10, 64)
		if err != nil {
			evt := &gatewayEventJSON{}
			err = unmarshal(data, evt)
			*ge = *evt.gatewayEvent()
			return
		}
		ge.Op = uint(tmp)
	}

	// data
	i += 2 // skip `,"`
	if data[i] != 'd' {
		evt := &gatewayEventJSON{}
		err = unmarshal(data, evt)
		*ge = *evt.gatewayEvent()
		return
	}
	i += 3 // skip `d":`
	ge.Data = data[i : len(data)-1]

	//fmt.Println("asdas")
	return
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

type DiscordWSEvent interface {
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
	return unmarshal(evt.Data(), v)
}
