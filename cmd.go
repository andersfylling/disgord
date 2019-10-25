package disgord

// CmdSettings holds all customizable command settings
type CmdSettings struct {
	Prefix     string
	IgnoreCase bool
}

// Cmd holds all important information about registering a new command
type Cmd struct {
	Name,
	Description string
	Handler func()
}

// CmdCtx holds all important information about a triggered command
type CmdCtx struct {
}

// activeCmdSettings holds the current command settings
var activeCmdSettings *CmdSettings

// registeredCommands holds all active commands
var registeredCommands = map[string]*Cmd{}
