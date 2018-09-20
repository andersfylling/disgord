package opcode

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

type OperationCodeHolder interface {
	GetOperationCode() uint
}

func ExtractFrom(holder OperationCodeHolder) uint {
	return holder.GetOperationCode()
}
