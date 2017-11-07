package presence

const (
	// Idle presence status for idle
	Idle = "idle"
	// Dnd presence status for dnd
	Dnd = "dnd"
	// Online presence status for online
	Online = "online"
	// Offline presence status for offline
	Offline = "offline"
)

// Status can either be "idle", "dnd", "online", or "offline"
type Status string
