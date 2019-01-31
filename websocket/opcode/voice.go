package opcode

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
