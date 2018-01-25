package disgord

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/andersfylling/disgord/discordws"
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/guild"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Token      string
	HTTPClient *http.Client
	Debug      bool
}

// NewDisgord creates a new default disgord instance
func NewDisgord(conf *Config) (*Disgord, error) {

	if conf.HTTPClient == nil {
		// http client configuration
		conf.HTTPClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	// Use discordws to keep the socket connection going
	dws, err := discordws.NewClient(&discordws.Config{
		// user settings
		Token:      conf.Token,
		HTTPClient: conf.HTTPClient,
		Debug:      conf.Debug,

		// lib specific
		DAPIVersion:  endpoint.APIVersion,
		DAPIEncoding: endpoint.APIComEncoding,
	})
	if err != nil {
		return nil, err
	}

	// create a disgord instance
	d := &Disgord{
		HTTPClient: conf.HTTPClient,
		ws:         dws,
		EventChan:  dws.GetEventChannel(),
		Token:      conf.Token,
	}

	return d, nil
}

// NewRequiredDisgord same as NewDisgord, but exits program if an error occours
func NewRequiredDisgord(conf *Config) *Disgord {
	dg, err := NewDisgord(conf)
	if err != nil {
		logrus.Fatal(err)
	}

	return dg
}

// EventObserver is an application-level type for handling discord requests.
// All callbacks are optional, and whether they are defined or not
// is used to determine whether the EventDispatcher will send events to them.
type EventDispatcher struct {
	// current EventHook fields here

	// OnEvent is called for all events.
	// Handlers must typecast the event type manually, and ensure
	// that it can handle receiving the same event twice if a type-specific
	// callback also exists.
	//OnEvent func(ctx *Context, ev event.DiscordEvent) error

	// OnMessageEvent is called for every message-related event.
	//OnMessageEvent func(ctx *Context, ev event.MessageEvent) error

	// OnConnectionEvent ...
	// OnUserEvent ...
	// OnChannelEvent ...
	// OnGuildEvent ...
}

type Disgord struct {
	sync.RWMutex

	Token string

	ws *discordws.Client

	HTTPClient *http.Client

	EventChan <-chan discordws.EventInterface

	// register listeners for events
	EventDispatcher

	// Guilds all them guild objects
	Guilds []*guild.Guild
}

func (d *Disgord) String() string {
	return d.ws.String()
}

func (d *Disgord) logInfo(msg string) {
	logrus.WithFields(logrus.Fields{
		"lib": d.ws.String(),
	}).Info(msg)
}

func (d *Disgord) logErr(msg string) {
	logrus.WithFields(logrus.Fields{
		"lib": d.ws.String(),
	}).Error(msg)
}

func (d *Disgord) eventObserver() {
	for {
		select {
		case evt, ok := <-d.EventChan:
			if !ok {
				logrus.Error("Event channel is dead!")
				break
			}
			logrus.Infof("Event{%s}\n%+v\n", evt.Name(), string(evt.Data()))
		}
	}
}

// Connect establishes a websocket connection to the discord API
func (d *Disgord) Connect() (err error) {
	d.logInfo("Connecting to discord Gateway")
	err = d.ws.Connect()
	if err != nil {
		d.logErr(err.Error())
		return
	}
	d.logInfo("Connected")

	// setup event observer
	go d.eventObserver()

	return nil
}

// Disconnect closes the discord websocket connection
func (d *Disgord) Disconnect() (err error) {
	fmt.Println()
	d.logInfo("Closing Discord gateway connection")
	err = d.ws.Disconnect()
	if err != nil {
		d.logErr(err.Error())
		return
	}
	d.logInfo("Disconnected")

	return nil
}
