package websocket

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"

	"github.com/andersfylling/snowflake/v3"
)

//////////////////////////////////////////////////////
//
// HELPER FUNCS
//
//////////////////////////////////////////////////////

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

type GatewayBotGetter interface {
	GetGatewayBot() (gateway *GatewayBot, err error)
}

func ConfigureShardConfig(client GatewayBotGetter, conf *ShardConfig) error {
	data, err := client.GetGatewayBot()
	if err != nil {
		return err
	}

	if conf.URL == "" {
		conf.URL = data.URL
	}
	if len(conf.ShardIDs) == 0 {
		for i := uint(0); i < data.Shards; i++ {
			conf.ShardIDs = append(conf.ShardIDs, i)
		}
	}
	if conf.ShardRateLimit == 0 {
		conf.ShardRateLimit = defaultShardRateLimit
	}

	return nil
}

//////////////////////////////////////////////////////
//
// VOICE SPECIFIC
//
//////////////////////////////////////////////////////

type voicePacket struct {
	Op   uint            `json:"op"`
	Data json.RawMessage `json:"d"`
}

type VoiceReady struct {
	SSRC  uint32   `json:"ssrc"`
	IP    string   `json:"ip"`
	Port  int      `json:"port"`
	Modes []string `json:"modes"`

	// From: https://discordapp.com/developers/docs/topics/voice-connections#establishing-a-voice-websocket-connection
	// `heartbeat_interval` here is an erroneous field and should be ignored.
	// The correct heartbeat_interval value comes from the Hello payload.

	// HeartbeatInterval uint `json:"heartbeat_interval"`
}

type voiceSelectProtocol struct {
	Protocol string                     `json:"protocol"`
	Data     *VoiceSelectProtocolParams `json:"data"`
}

type VoiceSelectProtocolParams struct {
	Address string `json:"address"`
	Port    uint16 `json:"port"`
	Mode    string `json:"mode"`
}

type VoiceSessionDescription struct {
	Mode      string   `json:"mode"`
	SecretKey [32]byte `json:"secret_key"`
}

type voiceIdentify struct {
	GuildID   snowflake.ID `json:"server_id"` // Yay for eventual consistency
	UserID    snowflake.ID `json:"user_id"`
	SessionID string       `json:"session_id"`
	Token     string       `json:"token"`
}

//////////////////////////////////////////////////////
//
// EVENT SPECIFIC
//
//////////////////////////////////////////////////////

type evtReadyPacket struct {
	SessionID string `json:"session_id"`
}

type evtIdentity struct {
	Token          string          `json:"token"`
	Properties     interface{}     `json:"properties"`
	Compress       bool            `json:"compress"`
	LargeThreshold uint            `json:"large_threshold"`
	Shard          *[2]uint        `json:"shard,omitempty"`
	Presence       json.RawMessage `json:"presence,omitempty"`
}

//////////////////////////////////////////////////////
//
// GENERAL PURPOSE
//
//////////////////////////////////////////////////////

// Gateway is for parsing the Gateway endpoint response
type Gateway struct {
	URL string `json:"url"`
}

// GatewayBot is for parsing the Gateway Bot endpoint response
type GatewayBot struct {
	Gateway
	Shards            uint `json:"shards"`
	SessionStartLimit struct {
		Total      uint `json:"total"`
		Remaining  uint `json:"remaining"`
		ResetAfter uint `json:"reset_after"`
	} `json:"session_start_limit"`
}

// clientPacket is outgoing packets by the client
type clientPacket struct {
	Op   uint        `json:"op"`
	Data interface{} `json:"d"`
}

type helloPacket struct {
	HeartbeatInterval uint `json:"heartbeat_interval"`
}

// discordPacketJSON is used when we need to fall back on the unmarshaler logic
type discordPacketJSON struct {
	Op             uint   `json:"op"`
	Data           []byte `json:"d"`
	SequenceNumber uint   `json:"s"`
	EventName      string `json:"t"`
}

func (p *discordPacketJSON) CopyOverTo(packet *DiscordPacket) {
	packet.Op = p.Op
	packet.Data = p.Data
	packet.SequenceNumber = p.SequenceNumber
	packet.EventName = p.EventName
}

// DiscordPacket is packets sent by Discord over the socket connection
type DiscordPacket struct {
	Op             uint            `json:"op"`
	Data           json.RawMessage `json:"d"`
	SequenceNumber uint            `json:"s,omitempty"`
	EventName      string          `json:"t,omitempty"`
}

// UnmarshalJSON see interface json.Unmarshaler
// TODO: benchmark json.RawMessage and this for both voice and event!
// there haven't been found any bugs in this code, but just to be
// sure the voice works correctly; json.RawMessage is preferred. For now..
//func (p *DiscordPacket) UnmarshalJSON(data []byte) (err error) {
//	var i int
//
//	// t
//	t := []byte{
//		'{', '"', 't', '"', ':',
//	}
//	for i = range t {
//		if t[i] != data[i] {
//			evt := discordPacketJSON{}
//			err = httd.Unmarshal(data, &evt)
//			evt.CopyOverTo(p)
//			return
//		}
//	}
//	i++                 // jump to next char
//	if data[i] == 'n' { // null
//		i += 4 // skip `null`
//	} else {
//		// extract the t value
//		var val strings.Builder
//		i++ // skip `"`
//		for ; data[i] != '"'; i++ {
//			val.WriteByte(data[i])
//		}
//		i++ // skip `"`
//		p.EventName = val.String()
//	}
//
//	// s
//	i += 2 // skip `,"`
//	if data[i] != 's' {
//		evt := discordPacketJSON{}
//		err = httd.Unmarshal(data, &evt)
//		evt.CopyOverTo(p)
//		return
//	}
//	i += 3              // skip `s":`
//	if data[i] == 'n' { // null value
//		i += 4 // skip `null`
//	} else {
//		// extract the s value
//		var val strings.Builder
//		for ; data[i] != ','; i++ {
//			val.WriteByte(data[i])
//		}
//		var tmp uint64
//		tmp, err = strconv.ParseUint(val.String(), 10, 64)
//		if err != nil {
//			evt := discordPacketJSON{}
//			err = httd.Unmarshal(data, &evt)
//			evt.CopyOverTo(p)
//			return
//		}
//		p.SequenceNumber = uint(tmp)
//	}
//
//	// op
//	i += 2              // skip `,"`
//	if data[i] != 'o' { // o as in op
//		evt := discordPacketJSON{}
//		err = httd.Unmarshal(data, &evt)
//		evt.CopyOverTo(p)
//		return
//	}
//	i += 4              // skip `op":`
//	if data[i] == 'n' { // null value
//		i += 4 // skip `null`
//	} else if data[i] == '0' {
//		p.Op = 0
//		i++ // skip 0
//	} else {
//		// extract the op value
//		var val strings.Builder
//		for ; data[i] != ','; i++ {
//			val.WriteByte(data[i])
//		}
//		var tmp uint64
//		tmp, err = strconv.ParseUint(val.String(), 10, 64)
//		if err != nil {
//			evt := discordPacketJSON{}
//			err = httd.Unmarshal(data, &evt)
//			evt.CopyOverTo(p)
//			return
//		}
//		p.Op = uint(tmp)
//	}
//
//	// data
//	i += 2 // skip `,"`
//	if data[i] != 'd' {
//		evt := discordPacketJSON{}
//		err = httd.Unmarshal(data, &evt)
//		evt.CopyOverTo(p)
//		return
//	}
//	i += 3 // skip `d":`
//	p.Data = data[i : len(data)-1]
//
//	return
//}

func (p *DiscordPacket) reset() {
	p.Op = 0
	p.SequenceNumber = 0
	// TODO: re-use data slice in unmarshal ?
	p.Data = nil
	p.EventName = ""
}
