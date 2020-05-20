package disgord

import (
	"github.com/andersfylling/disgord/internal/logger"
)

// Logger super basic logging interface
type Logger = logger.Logger

// Deprecated
func DefaultLogger(debug bool) logger.Logger {
	panic("this has been removed, please see examples/docs/logging-* for more information")
}

// Deprecated
func DefaultLoggerWithInstance(log logger.Logger) logger.Logger {
	return DefaultLogger(true)
}
