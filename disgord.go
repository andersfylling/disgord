package disgord

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/andersfylling/disgord/endpoint"
	"github.com/gorilla/websocket"
)

func NewDisgord() *Disgord {
	d := &Disgord{
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	return d
}

// EventHook is an application-level type for handling discord requests.
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

	// WS web socket connection
	WS *websocket.Conn

	// WSSURL web socket url
	WSSURL string

	HTTPClient *http.Client

	// register listeners for events
	//*EventObserver
}

// Open establishes a websocket connection to the discord API
func (d *Disgord) Connect() error {
	d.Lock()
	defer d.Unlock()
	var err error // creates issue with ws connection later when created

	// check if web socket connection is already open
	if d.WS != nil {
		return errors.New("websocket connection already established. Cannot open a new connection.")
	}

	// discord API sends a web socket url which should be used.
	// It's required to be cached, and only re-requested whenever disgord is unable
	// to reconnect to the API..
	if d.WSSURL == "" {
		resp, err := d.HTTPClient.Get(endpoint.Gateway)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		gatewayResponse := &GetGatewayResponse{}
		err = json.Unmarshal([]byte(body), gatewayResponse)
		if err != nil {
			return err
		}

		d.WSSURL = gatewayResponse.URL + "?v=" + endpoint.APIVersion + "&encoding=" + endpoint.APIComEncoding
	}

	// establish ws connection
	d.WS, _, err = websocket.DefaultDialer.Dial(d.WSSURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	// NOTE. I completely stole this defer func from discordgo. It was a too nice not to.
	defer func() {
		// because of this, all code below must set err to the error
		// when exiting with an error :)  Maybe someone has a better
		// way :)
		if err != nil {
			d.Disconnect()
		}
	}()

	return nil
}

// Disconnect closes the discord websocket connection
func (d *Disgord) Disconnect() error {
	err := d.WS.Close()
	if err != nil {
		return err
	}
	d.WS = nil

	return nil
}
