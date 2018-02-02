package disgord

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/andersfylling/disgord/discordws"
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/event"
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
		Dispatcher: event.NewDispatcher(),
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

type Disgord struct {
	sync.RWMutex

	Token string

	ws *discordws.Client

	HTTPClient *http.Client

	EventChan <-chan discordws.EventInterface

	// register listeners for events
	*event.Dispatcher

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

func (d *Disgord) eventObserver() {
	for {
		select {
		case evt, ok := <-d.EventChan:
			if !ok {
				logrus.Error("Event channel is dead!")
				break
			}

			switch evt.Name() {
			case event.GuildCreate:
				guild := &guild.Guild{}
				err := json.Unmarshal(evt.Data(), guild)
				if err != nil {
					panic(err)
				}
				ctx, _ := context.WithTimeout(context.Background(), time.Duration(5))
				d.GuildCreateEvent.Trigger(ctx, guild)
			default:
				fmt.Printf("------\nTODO\nImplement event handler for `%s`, data: \n%+v\n------\n\n", evt.Name(), string(evt.Data()))
			}
		}
	}
}
