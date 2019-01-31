package opcode

// operation codes sent by Discord over the socket connection
const (
	DiscordEvent uint = iota
	Heartbeat
	Identify
	StatusUpdate
	VoiceStateUpdate
	VoiceServerPing
	Resume
	Reconnect
	RequestGuildMembers
	InvalidSession
	Hello
	HeartbeatAck
)
