package main

import (
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"go.uber.org/zap"
	"os"
)

// DefaultLogger create a new logger instance for Disgord with the option to activate debugging.
func MyInjectableLogger(conf zap.Config) *LoggerZap {
	log, err := conf.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	return &LoggerZap{
		instance: log.With(
			zap.String("lib", disgord.LibraryInfo()),
		),
	}
}

type LoggerZap struct {
	instance *zap.Logger
}

var _ disgord.Logger = (*LoggerZap)(nil)

func (log *LoggerZap) Debug(v ...interface{}) {
	log.instance.Debug(fmt.Sprint(v...))
	_ = log.instance.Sync()
}

func (log *LoggerZap) Info(v ...interface{}) {
	log.instance.Info(fmt.Sprint(v...))
	_ = log.instance.Sync()
}

func (log *LoggerZap) Error(v ...interface{}) {
	log.instance.Error(fmt.Sprint(v...))
	_ = log.instance.Sync()
}

func main() {
	logConf := zap.NewProductionConfig()
	logConf.Level.SetLevel(zap.DebugLevel)

	// Set up a new Disgord client
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
		Logger:   MyInjectableLogger(logConf),
	})
	defer client.StayConnectedUntilInterrupted(context.Background())
}
