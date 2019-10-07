package logger

import (
	"fmt"

	"github.com/andersfylling/disgord/internal/constant"
	"go.uber.org/zap"
)

// DefaultLogger create a new logger instance for DisGord with the option to activate debugging.
func DefaultLogger(debug bool) *LoggerZap {
	conf := zap.NewProductionConfig()
	if debug {
		conf.Development = true
		conf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	log, err := conf.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	return &LoggerZap{
		instance: log.With(
			zap.String("lib", constant.Name),
			zap.String("ver", constant.Version)),
	}
}

func DefaultLoggerWithInstance(log *zap.Logger) *LoggerZap {
	return &LoggerZap{
		instance: log,
	}
}

type LoggerZap struct {
	instance *zap.Logger
}

func (log *LoggerZap) getMessage(v ...interface{}) (sentence string) {
	for i := range v {
		var msg string
		if str, ok := v[i].(fmt.Stringer); ok {
			msg = str.String()
		} else {
			switch t := v[i].(type) {
			case string:
				msg = t
			case error:
				msg = t.Error()
			default:
				// TODO
				msg = fmt.Sprint(v[i])
			}
		}

		if sentence != "" {
			sentence += " " + msg
		} else {
			sentence = msg
		}
	}

	return sentence
}

func (log *LoggerZap) Debug(v ...interface{}) {
	log.instance.Debug(log.getMessage(v))
	_ = log.instance.Sync()
}

func (log *LoggerZap) Info(v ...interface{}) {
	log.instance.Info(log.getMessage(v))
	_ = log.instance.Sync()
}

func (log *LoggerZap) Error(v ...interface{}) {
	log.instance.Error(log.getMessage(v))
	_ = log.instance.Sync()
}

var _ Logger = (*LoggerZap)(nil)
