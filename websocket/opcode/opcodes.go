package opcode

import "math"

// custom op codes used by DisGord internally
const (
	Shutdown          uint = 100
	Close             uint = 101
	EventReadyResumed uint = 102 // Discord use 0 here, but that is a shared op code
	NoOPCode          uint = math.MaxUint16
)

// operation codes for the event client
const (
	EventDiscordEvent uint = iota
	EventHeartbeat
	EventIdentify
	EventStatusUpdate
	EventVoiceStateUpdate
	EventVoiceServerPing
	EventResume
	EventReconnect
	EventRequestGuildMembers
	EventInvalidSession
	EventHello
	EventHeartbeatAck
)

// operation codes for the voice client
const (
	VoiceIdentify uint = iota
	VoiceSelectProtocol
	VoiceReady
	VoiceHeartbeat
	VoiceSessionDescription
	VoiceSpeaking
	VoiceHeartbeatAck
	VoiceResume
	VoiceHello
	VoiceResumed
	_ // unused
	_ // unused
	_ // unused
	VoiceClientDisconnect
)

// OperationCodeHolder Used on objects that holds a operation code
type OperationCodeHolder interface {
	GetOperationCode() uint
}

// ExtractFrom extract the operation code
func ExtractFrom(holder OperationCodeHolder) uint {
	return holder.GetOperationCode()
}
