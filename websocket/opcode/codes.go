package opcode

const (
	DiscordEvent   = 0
	Ping           = 1
	Identity       = 2
	Resume         = 6
	Reconnect      = 7
	InvalidSession = 9
	Hello          = 10
	Heartbeat      = 11
)

type OperationCodeHolder interface {
	GetOperationCode() uint
}

func ExtractFrom(holder OperationCodeHolder) uint {
	return holder.GetOperationCode()
}
