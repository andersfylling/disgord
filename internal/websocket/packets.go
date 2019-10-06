package websocket

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"

	opcode2 "github.com/andersfylling/disgord/internal/websocket/opcode"
)

//////////////////////////////////////////////////////
//
// HELPER FUNC(TION)S
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
	GuildID   Snowflake `json:"server_id"` // Yay for eventual consistency
	UserID    Snowflake `json:"user_id"`
	SessionID string    `json:"session_id"`
	Token     string    `json:"token"`
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
	Token              string          `json:"token"`
	Properties         interface{}     `json:"properties"`
	Compress           bool            `json:"compress"`
	LargeThreshold     uint            `json:"large_threshold"`
	Shard              *[2]uint        `json:"shard,omitempty"`
	Presence           json.RawMessage `json:"presence,omitempty"`
	GuildSubscriptions bool            `json:"guild_subscriptions"` // most ambiguous naming ever but ok.
}

type evtResume struct {
	Token      string `json:"token"`
	SessionID  string `json:"session_id"`
	SequenceNr uint   `json:"seq"`
}

type RequestGuildMembersPayload struct {
	// GuildID	id of the guild(s) to get offline members for
	GuildIDs []Snowflake `json:"guild_id"`

	// Query string that username starts with, or an empty string to return all members
	Query string `json:"query"`

	// Limit maximum number of members to send or 0 to request all members matched
	Limit uint `json:"limit"`

	// UserIDs used to specify which users you wish to fetch
	UserIDs []Snowflake `json:"user_ids,omitempty"`
}

var _ CmdPayload = (*RequestGuildMembersPayload)(nil)

func (u *RequestGuildMembersPayload) isCmdPayload() bool { return true }

type UpdateVoiceStatePayload struct {
	// GuildID id of the guild
	GuildID Snowflake `json:"guild_id"`

	// ChannelID id of the voice channel Client wants to join
	// (set to 0 if disconnecting)
	ChannelID Snowflake `json:"channel_id"`

	// SelfMute is the Client mute
	SelfMute bool `json:"self_mute"`

	// SelfDeaf is the Client deafened
	SelfDeaf bool `json:"self_deaf"`
}

var _ CmdPayload = (*UpdateVoiceStatePayload)(nil)

func (u *UpdateVoiceStatePayload) isCmdPayload() bool { return true }

type UpdateStatusPayload struct {
	// Since unix time (in milliseconds) of when the Client went idle, or null if the Client is not idle
	Since *uint `json:"since"`

	// Game null, or the user's new activity
	Game interface{} `json:"game"`

	// Status the user's new status
	Status string `json:"status"`

	// AFK whether or not the Client is afk
	AFK bool `json:"afk"`
}

var _ CmdPayload = (*UpdateStatusPayload)(nil)

func (u *UpdateStatusPayload) isCmdPayload() bool { return true }

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
	Op      opcode2.OpCode `json:"op"`
	Data    interface{}    `json:"d"`
	CmdName string         `json:"-"`
}

type helloPacket struct {
	HeartbeatInterval uint `json:"heartbeat_interval"`
}

// discordPacketJSON is used when we need to fall back on the unmarshaler logic
type discordPacketJSON struct {
	Op             opcode2.OpCode `json:"op"`
	Data           []byte         `json:"d"`
	SequenceNumber uint           `json:"s"`
	EventName      string         `json:"t"`
}

func (p *discordPacketJSON) CopyOverTo(packet *DiscordPacket) {
	packet.Op = p.Op
	packet.Data = p.Data
	packet.SequenceNumber = p.SequenceNumber
	packet.EventName = p.EventName
}

// DiscordPacket is packets sent by Discord over the socket connection
type DiscordPacket struct {
	Op             opcode2.OpCode  `json:"op"`
	Data           json.RawMessage `json:"d"`
	SequenceNumber uint            `json:"s,omitempty"`
	EventName      string          `json:"t,omitempty"`
}

func (p *DiscordPacket) reset() {
	p.Op = 0
	p.SequenceNumber = 0
	// TODO: re-use data slice in unmarshal ?
	p.Data = nil
	p.EventName = ""
}
