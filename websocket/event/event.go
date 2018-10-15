package event

// Different socket commands that can be sent to Discord
const (
	Heartbeat           = "HEARTBEAT"
	Ready               = "READY"
	Resume              = "RESUME"
	Identify            = "IDENTIFY"
	StatusUpdate        = "STATUS_UPDATE"
	VoiceStateUpdate    = "VOICE_STATE_UPDATE"
	RequestGuildMembers = "REQUEST_GUILD_MEMBERS"
)

// custom events for Disgord. Don't use these.
const (
	Shutdown = "_"
)
