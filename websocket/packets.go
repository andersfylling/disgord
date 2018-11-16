package websocket

import (
	"bytes"
	"compress/zlib"
	"io"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord/httd"
)

// discordPacketJSON is used when we need to fall back on the unmarshaler logic
type discordPacketJSON struct {
	Op             uint   `json:"op"`
	Data           []byte `json:"d"`
	SequenceNumber uint   `json:"s"`
	EventName      string `json:"t"`
}

func (p *discordPacketJSON) CopyOverTo(packet *discordPacket) {
	packet.Op = p.Op
	packet.Data = p.Data
	packet.SequenceNumber = p.SequenceNumber
	packet.EventName = p.EventName
}

// discordPacket is packets sent by Discord over the socket connection
type discordPacket struct {
	Op             uint   `json:"op"`
	Data           []byte `json:"d"`
	SequenceNumber uint   `json:"s"`
	EventName      string `json:"t"`
}

// UnmarshalJSON see interface json.Unmarshaler
func (p *discordPacket) UnmarshalJSON(data []byte) (err error) {
	var i int

	// t
	t := []byte{
		'{', '"', 't', '"', ':',
	}
	for i = range t {
		if t[i] != data[i] {
			evt := discordPacketJSON{}
			err = httd.Unmarshal(data, &evt)
			evt.CopyOverTo(p)
			return
		}
	}
	i++                 // jump to next char
	if data[i] == 'n' { // null
		i += 4 // skip `null`
	} else {
		// extract the t value
		var val strings.Builder
		i++ // skip `"`
		for ; data[i] != '"'; i++ {
			val.WriteByte(data[i])
		}
		i++ // skip `"`
		p.EventName = val.String()
	}

	// s
	i += 2 // skip `,"`
	if data[i] != 's' {
		evt := discordPacketJSON{}
		err = httd.Unmarshal(data, &evt)
		evt.CopyOverTo(p)
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
			evt := discordPacketJSON{}
			err = httd.Unmarshal(data, &evt)
			evt.CopyOverTo(p)
			return
		}
		p.SequenceNumber = uint(tmp)
	}

	// op
	i += 2              // skip `,"`
	if data[i] != 'o' { // o as in op
		evt := discordPacketJSON{}
		err = httd.Unmarshal(data, &evt)
		evt.CopyOverTo(p)
		return
	}
	i += 4              // skip `op":`
	if data[i] == 'n' { // null value
		i += 4 // skip `null`
	} else if data[i] == '0' {
		p.Op = 0
		i++ // skip 0
	} else {
		// extract the op value
		var val strings.Builder
		for ; data[i] != ','; i++ {
			val.WriteByte(data[i])
		}
		var tmp uint64
		tmp, err = strconv.ParseUint(val.String(), 10, 64)
		if err != nil {
			evt := discordPacketJSON{}
			err = httd.Unmarshal(data, &evt)
			evt.CopyOverTo(p)
			return
		}
		p.Op = uint(tmp)
	}

	// data
	i += 2 // skip `,"`
	if data[i] != 'd' {
		evt := discordPacketJSON{}
		err = httd.Unmarshal(data, &evt)
		evt.CopyOverTo(p)
		return
	}
	i += 3 // skip `d":`
	p.Data = data[i : len(data)-1]

	//fmt.Println("asdas")
	return
}

// clientPacket is outgoing packets by the client
type clientPacket struct {
	Op   uint        `json:"op"`
	Data interface{} `json:"d"`
}

type traceData struct {
	Trace []string `json:"_trace"`
}

type helloPacket struct {
	HeartbeatInterval uint `json:"heartbeat_interval"`
	traceData
}

type readyPacket struct {
	SessionID string `json:"session_id"`
	traceData
}

// decompressBytes decompresses a binary message
func decompressBytes(input []byte) (output []byte, err error) {
	b := bytes.NewReader(input)
	var r io.ReadCloser

	r, err = zlib.NewReader(b)
	if err != nil {
		return
	}
	defer r.Close()

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(r)
	if err != nil {
		return
	}

	output = buffer.Bytes()
	return
}
