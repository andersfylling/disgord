// +build disgordperf

package gateway

import (
	"strconv"
	"strings"

	"github.com/andersfylling/disgord/internal/gateway/opcode"
	"github.com/andersfylling/disgord/json"
)

//UnmarshalJSON see interface json.Unmarshaler
//TODO: benchmark json.RawMessage and this for both voice and event!
//there haven't been found any bugs in this code, but just to be
//sure the voice works correctly; json.RawMessage is preferred. For now..
func (p *DiscordPacket) UnmarshalJSON(data []byte) (err error) {
	var i int

	// t
	t := []byte{
		'{', '"', 't', '"', ':',
	}
	for i = range t {
		if t[i] != data[i] {
			evt := discordPacketJSON{}
			err = json.Unmarshal(data, &evt)
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
		err = json.Unmarshal(data, &evt)
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
		tmp, err = strconv.ParseUint(val.String(), 10, 32)
		if err != nil {
			evt := discordPacketJSON{}
			err = json.Unmarshal(data, &evt)
			evt.CopyOverTo(p)
			return
		}
		p.SequenceNumber = uint32(tmp)
	}

	// op
	i += 2              // skip `,"`
	if data[i] != 'o' { // o as in op
		evt := discordPacketJSON{}
		err = json.Unmarshal(data, &evt)
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
			err = json.Unmarshal(data, &evt)
			evt.CopyOverTo(p)
			return
		}
		p.Op = opcode.OpCode(tmp)
	}

	// data
	i += 2 // skip `,"`
	if data[i] != 'd' {
		evt := discordPacketJSON{}
		err = json.Unmarshal(data, &evt)
		evt.CopyOverTo(p)
		return
	}
	i += 3 // skip `d":`
	p.Data = data[i : len(data)-1]

	return
}
