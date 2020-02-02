package opcode

import "math"

type OpCode uint

// custom op codes used by Disgord internally
const (
	EventReadyResumed OpCode = 102 // Discord use 0 here, but that is a shared op code
	None              OpCode = math.MaxUint16
)

// operation codes for the event client
const (
	EventDiscordEvent OpCode = iota
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
	VoiceIdentify OpCode = iota
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
	GetOperationCode() OpCode
}

// ExtractFrom extract the operation code
func ExtractFrom(holder OperationCodeHolder) OpCode {
	return holder.GetOperationCode()
}
