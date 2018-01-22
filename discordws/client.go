package discordws

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/andersfylling/snowflake"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// NewRequiredClient same as NewClient(...), but program exits on failure.
func NewRequiredClient(config *Config) *Client {
	c, err := NewClient(config)
	if err != nil {
		logrus.Fatal(err)
	}

	return c
}

// NewClient Creates a new discord websocket client
func NewClient(conf *Config) (*Client, error) {
	if conf == nil {
		return nil, errors.New("Missing Config.Token for discord authentication")
	}

	if conf.DAPIVersion < LowestAPIVersion || conf.DAPIVersion > HighestAPIVersion {
		return nil, fmt.Errorf("Discord API version %d is not supported. Lowest supported version is %d, and highest is %d", conf.DAPIVersion, LowestAPIVersion, HighestAPIVersion)
	}

	encoding := strings.ToLower(conf.DAPIEncoding)
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

	// check the http client exists. Otherwise create one.
	if conf.HTTPClient == nil {
		conf.HTTPClient = &http.Client{
			Timeout: time.Second * DefaultHTTPTimeout,
		}
	}

	// configure logrus output
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// return configured discord websocket client
	return &Client{
		token:              conf.Token,
		URLAPIVersion:      BaseURL + "/v" + strconv.Itoa(conf.DAPIVersion),
		HTTPClient:         conf.HTTPClient,
		DiscordAPIVersion:  conf.DAPIVersion,
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

	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	sequenceNumber    uint          `json:"-"`
	Trace             []string      `json:"_trace"`
	SessionID         string        `json:"session_id"`
	token             string        `json:"-"`
	ShardCount        uint          `json:"shard_count"`
	ShardID           snowflake.ID  `json:"shard_id"`

	connected chan struct{}

	// websocket connection
	*websocket.Conn `json:"-"`
	wsMutex         sync.Mutex // https://hackernoon.com/dancing-with-go-s-mutexes-92407ae927bf
}

func (dws *Client) String() string {
	return fmt.Sprintf("%s v%d.%d.%d", LibName, LibVersionMajor, LibVersionMinor, LibVersionPatch)
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
