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

// OperationCodeHolder Used on objects that holds a operation code
type OperationCodeHolder interface {
	GetOperationCode() uint
}

// ExtractFrom extract the operation code
func ExtractFrom(holder OperationCodeHolder) uint {
	return holder.GetOperationCode()
}
