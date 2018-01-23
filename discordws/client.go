package discordws

import (
	"errors"
	"fmt"
	"net/http"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/andersfylling/snowflake"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// NewRequiredClient same as NewClient(...), but program exits on failure.
func NewRequiredClient(conf *Config) *Client {
	c, err := NewClient(conf)
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
		return nil, fmt.Errorf("Discord requires data encoding to be of the following '%s', while '%s' encoding was requested", strings.Join(Encodings, "', '"), encoding)
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
		token:             conf.Token,
		urlAPIVersion:     BaseURL + "/v" + strconv.Itoa(conf.DAPIVersion),
		httpClient:        conf.HTTPClient,
		dAPIVersion:       conf.DAPIVersion,
		dAPIEncoding:      encoding,
		heartbeatAcquired: time.Now(),
		disconnected:      nil,
		operationChan:     make(chan *GatewayPayload),
		eventChans:        make(map[string](chan []byte)),
		sendChan:          make(chan *GatewayPayload),
		Myself:            &user.User{},
	}, nil
}

// Client holds the web socket state and can be used directly in marshal/unmarshal to work with intance data
type Client struct {
	sync.RWMutex `json:"-"`

	urlAPIVersion string `json:"-"`

	// URL Websocket URL web socket url
	url string `json:"-"`

	httpClient *http.Client `json:"-"`

	dAPIVersion    int    `json:"-"`
	dAPIEncoding   string `json:"-"`
	token          string `json:"-"`
	sequenceNumber uint   `json:"-"`

	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	heartbeatAcquired time.Time     `json:"-"`
	Trace             []string      `json:"_trace"`
	SessionID         string        `json:"session_id"`
	ShardCount        uint          `json:"shard_count"`
	ShardID           snowflake.ID  `json:"shard_id"`

	disconnected  chan struct{}            `json:"-"`
	operationChan chan *GatewayPayload     `json:"-"`
	eventChans    map[string](chan []byte) `json:"-"`
	sendChan      chan *GatewayPayload     `json:"-"`

	Myself *user.User `json:"user"`

	// websocket connection
	conn    *websocket.Conn `json:"-"`
	wsMutex sync.Mutex      `json:"-"` // https://hackernoon.com/dancing-with-go-s-mutexes-92407ae927bf

	// heartbeat mutex keeps us from creating another pulser
	pulseMutex sync.Mutex `json:"-"`
}

func (c *Client) String() string {
	return fmt.Sprintf("%s v%d.%d.%d", LibName, LibVersionMajor, LibVersionMinor, LibVersionPatch)
}

// Dead check if the websocket connection isn't established AKA "dead"
func (c *Client) Dead() bool {
	return c.conn == nil
}

// Routed checks if the client has recieved the root endpoint for discord API communication
func (c *Client) Routed() bool {
	return c.url != ""
}

// RemoveRoute deletes cached discord wss endpoint
func (c *Client) RemoveRoute() {
	c.url = ""
}
