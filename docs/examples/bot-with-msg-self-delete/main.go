package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
)

var log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.DebugLevel,
}

const MessageLifeTime = 5 * time.Second

func main() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
		Logger:   log, // optional logging
		Cache:    &disgord.CacheNop{},
	})
	run(client)
}

func run(client *disgord.Client) {
	mdlw, err := NewMiddlewareHolder(client, context.Background())
	if err != nil {
		panic(err)
	}

	// listen for messages
	client.On(disgord.EvtMessageCreate,
		mdlw.filterOutHumans,
		mdlw.filterOutOthersMsgs,
		autoDeleteNewMessages)

	// connect now, and disconnect on system interrupt
	client.StayConnectedUntilInterrupted(context.Background())
}

//////////////////////////////////////////////////////
//
// HANDLERS
//
//////////////////////////////////////////////////////
func autoDeleteNewMessages(s disgord.Session, evt *disgord.MessageCreate) {
	// delete message after N seconds
	<-time.After(MessageLifeTime)

	err := s.DeleteFromDiscord(context.Background(), evt.Message)
	if err != nil {
		log.Error(err)
	}
}

//////////////////////////////////////////////////////
//
// MIDDLEWARES
//
//////////////////////////////////////////////////////
func NewMiddlewareHolder(s disgord.Session, ctx context.Context) (m *MiddlewareHolder, err error) {
	m = &MiddlewareHolder{session: s}
	if m.myself, err = s.GetCurrentUser(ctx); err != nil {
		return nil, errors.New("unable to fetch info about the bot instance")
	}

	return m, nil
}

// instead of storing the instances in global variables. Middlewares can ofcourse be standalone functions too.
type MiddlewareHolder struct {
	session disgord.Session
	myself  *disgord.User
}

func (m *MiddlewareHolder) filterOutHumans(evt interface{}) interface{} {
	if e, ok := evt.(*disgord.MessageCreate); ok {
		// ignore humans
		if !e.Message.Author.Bot {
			return nil
		}
	}

	return evt
}

func (m *MiddlewareHolder) filterOutOthersMsgs(evt interface{}) interface{} {
	// ignore other bots
	// remove this check if you want to delete all bot messages after N seconds
	if e, ok := evt.(*disgord.MessageCreate); ok {
		if e.Message.Author.ID != m.myself.ID {
			return nil
		}
	}

	return evt
}
