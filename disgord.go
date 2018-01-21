package disgord

import (
	"net/http"
	"sync"
	"time"

	"github.com/andersfylling/disgord/discordws"
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/guild"
	"github.com/sirupsen/logrus"
)

// NewDisgord creates a new default disgord instance
func NewDisgord() (*Disgord, error) {
	// http client configuration
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	return NewDisgordWithHTTPClient(httpClient)
}

// NewRequiredDisgord same as NewDisgord, but exits program if an error occours
func NewRequiredDisgord() *Disgord {
	dg, err := NewDisgord()
	if err != nil {
		logrus.Fatal(err)
	}

	return dg
}

// NewDisgordWithHTTPClient specify http configuration for the discord connection. Affects REST and pre- websocket handshake, for wss endpoint request
func NewDisgordWithHTTPClient(httpClient *http.Client) (*Disgord, error) {
	// Use discordws to keep the socket connection going
	dws, err := discordws.NewClient(httpClient, endpoint.APIVersion, endpoint.APIComEncoding)
	if err != nil {
		return nil, err
	}

	// create a disgord instance
	d := &Disgord{
		HTTPClient: httpClient,
		ws:         dws,
	}

	return d, nil
}

// NewRequiredDisgordWithHTTPClient same as NewDisgordWithHTTPClient, but exits program if an error occours
func NewRequiredDisgordWithHTTPClient(httpClient *http.Client) *Disgord {
	dg, err := NewDisgordWithHTTPClient(httpClient)
	if err != nil {
		logrus.Fatal(err)
	}

	return dg
}

// EventObserver is an application-level type for handling discord requests.
// All callbacks are optional, and whether they are defined or not
// is used to determine whether the EventDispatcher will send events to them.
type EventObserver struct {
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

	ws *discordws.Client

	HTTPClient *http.Client

	// register listeners for events
	//*EventObserver

	// Guilds all them guild objects
	Guilds []*guild.Guild
}

// Connect establishes a websocket connection to the discord API
func (d *Disgord) Connect() error {
	d.Lock()
	defer d.Unlock()
	return d.ws.Connect()
}

// Disconnect closes the discord websocket connection
func (d *Disgord) Disconnect() error {
	d.Lock()
	defer d.Unlock()
	return d.ws.Disconnect()
}
