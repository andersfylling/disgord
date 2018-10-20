package event

// Different socket commands that can be sent to Discord
const (
	Heartbeat           = "HEARTBEAT"
	Ready               = "READY"
	Resume              = "RESUME"
	Resumed             = "RESUMED"
	Identify            = "IDENTIFY"
	StatusUpdate        = "STATUS_UPDATE"
	VoiceStateUpdate    = "VOICE_STATE_UPDATE"
	RequestGuildMembers = "REQUEST_GUILD_MEMBERS"
)

// custom events for Disgord. Don't use these.
const (
	Shutdown = "_"
	Close    = "-"
)
