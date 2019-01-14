package disgord

import (
	"github.com/andersfylling/disgord/constant"
	"go.uber.org/zap"
)

// Logger super basic logging interface
type Logger = constant.Logger

// DefaultLogger create a new logger instance for DisGord with the option to activate debugging.
func DefaultLogger(debug bool) *LoggerZap {
	conf := zap.NewProductionConfig()
	if debug {
		conf.Development = true
		conf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	log, err := conf.Build()
	if err != nil {
		panic(err)
	}

	return &LoggerZap{
		instance: log.With(
			zap.String("lib", constant.Name),
			zap.String("ver", constant.Version)),
	}
}

type LoggerZap struct {
	instance *zap.Logger
}

func (log *LoggerZap) Debug(msg string) {
	log.instance.Debug(msg)
	_ = log.instance.Sync()
}
func (log *LoggerZap) Info(msg string) {
	log.instance.Info(msg)
	_ = log.instance.Sync()
}
func (log *LoggerZap) Error(msg string) {
	log.instance.Error(msg)
	_ = log.instance.Sync()
}
