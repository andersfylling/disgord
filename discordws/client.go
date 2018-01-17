package discordws

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// NewDefaultClient uses the latest discord API (2018-01-17), JSON encoding and a http.Client with a N second timeout. See discordws.DefaultHTTPTimeout (most likely 10seconds)
func NewDefaultClient() (*Client, error) {
	// http client configuration
	httpClient := &http.Client{
		Timeout: time.Second * DefaultHTTPTimeout,
	}

	return NewClient(httpClient, HighestAPIVersion, EncodingJSON)
}

// NewRequiredDefaultClient same as NewDefaultClient(), but program exits on failure.
func NewRequiredDefaultClient() *Client {
	c, err := NewDefaultClient()
	if err != nil {
		logrus.Fatal(err)
	}

	return c
}

// NewRequiredClient same as NewClient(...), but program exits on failure.
func NewRequiredClient(httpClient *http.Client, discordAPIVersion int, discordAPIEncoding string) *Client {
	c, err := NewClient(httpClient, discordAPIVersion, discordAPIEncoding)
	if err != nil {
		logrus.Fatal(err)
	}

	return c
}

// NewClient Creates a new discord websocket client
func NewClient(httpClient *http.Client, discordAPIVersion int, discordAPIEncoding string) (*Client, error) {
	if discordAPIVersion < LowestAPIVersion || discordAPIVersion > HighestAPIVersion {
		return nil, fmt.Errorf("Discord API version %d is not supported. Lowest supported version is %d, and highest is %d", discordAPIVersion, LowestAPIVersion, HighestAPIVersion)
	}

	encoding := strings.ToLower(discordAPIEncoding)
	var acceptedEncoding bool
	for _, required := range Encodings {
		if encoding == required {
			acceptedEncoding = true
			break
		}
	}
	if !acceptedEncoding {
		return nil, fmt.Errorf("Discord requires data encoding to be of the following %swhile %s encoding was requested", strings.Join(Encodings, ", "), encoding)
	}

	// configure logrus output
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// return configured discord websocket client
	return &Client{
		URLAPIVersion:      BaseURL + "/v" + strconv.Itoa(discordAPIVersion),
		HTTPClient:         httpClient,
		DiscordAPIVersion:  discordAPIVersion,
		DiscordAPIEncoding: encoding,
	}, nil
}

// Client holds the web socket state
type Client struct {
	sync.RWMutex

	URLAPIVersion string `json:"-"`

	// URL Websocket URL web socket url
	URL string `json:"-"`

	HTTPClient *http.Client

	DiscordAPIVersion  int    `json:"v"`
	DiscordAPIEncoding string `json:"encoding"`

	HeartbeatInterval uint     `json:"heartbeat_interval"`
	Trace             []string `json:"_trace"`

	// websocket connection
	*websocket.Conn `json:"-"`
}

// Dead check if the websocket connection isn't established AKA "dead"
func (dws *Client) Dead() bool {
	return dws.Conn == nil
}

// Routed checks if the client has recieved the root endpoint for discord API communication
func (dws *Client) Routed() bool {
	return dws.URL != ""
}

// RemoveRoute deletes cached discord wss endpoint
func (dws *Client) RemoveRoute() {
	dws.URL = ""
}

// Connect establishes a websocket connection to the discord API
func (dws *Client) Connect() error {
	dws.Lock()
	defer dws.Unlock()
	var err error // creates issue with ws connection later when created

	// check if web socket connection is already open
	if !dws.Dead() {
		return errors.New("websocket connection already established, cannot open a new one")
	}

	// discord API sends a web socket url which should be used.
	// It's required to be cached, and only re-requested whenever disgord is unable
	// to reconnect to the API..
	if !dws.Routed() {
		var resp *http.Response
		resp, err = dws.HTTPClient.Get(dws.URLAPIVersion + "/gateway")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		gatewayResponse := &GetGatewayResponse{}
		err = json.Unmarshal([]byte(body), gatewayResponse)
		if err != nil {
			return err
		}

		dws.URL = gatewayResponse.URL + "?v=" + strconv.Itoa(dws.DiscordAPIVersion) + "&encoding=" + dws.DiscordAPIEncoding
	}

	// establish ws connection
	dws.Conn, _, err = websocket.DefaultDialer.Dial(dws.URL, nil)
	if err != nil {
		return err
	}

	// NOTE. I completely stole this defer func from discordgo. It was a too nice not to take advantage of.
	//  Will this return any potential errors?
	defer func() error {
		if err != nil {
			return dws.Disconnect()
		}

		return nil
	}()

	// Once connected, the client should immediately receive an Opcode 10 Hello payload, with information on the
	// connection's heartbeat interval:
	// {
	//   "heartbeat_interval": 45000,
	//   "_trace": ["discord-gateway-prd-1-99"]
	// }
	messageType, p, err := dws.ReadMessage()
	if err != nil {
		return err
	}
	// Only accept text messages atm
	if messageType != websocket.TextMessage {
		err = errors.New("disgord only accepts text message{type:1} as the initial web socket contact")
		return err
	}

	// handle the incoming hello data and store it in our ws configuration
	payload := GatewayPayload{
		Data: dws,
	}
	err = json.Unmarshal(p, &payload)
	if err != nil {
		return err
	}

	return nil
}

// Disconnect closes the discord websocket connection
func (dws *Client) Disconnect() error {
	dws.Lock()
	defer dws.Unlock()
	err := dws.Close()
	if err != nil {
		return err
	}
	dws.Conn = nil

	return nil
}

// Reconnect to discord endpoint
func (dws *Client) Reconnect() error {
	dws.Lock()
	defer dws.Unlock()

	for try := 0; try < MaxReconnectTries; try++ {

	}

	return nil
}
