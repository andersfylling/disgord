package discordws

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	//"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/snowflake"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Token        string
	HTTPClient   *http.Client
	DAPIVersion  int
	DAPIEncoding string
	Debug        bool
}

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
		return nil, errors.New("missing Config.Token for discord authentication")
	}

	// ensure this socket module supports the given discord api version
	if conf.DAPIVersion < LowestAPIVersion || conf.DAPIVersion > HighestAPIVersion {
		return nil, fmt.Errorf("discord API version %d is not supported. Lowest supported version is %d, and highest is %d", conf.DAPIVersion, LowestAPIVersion, HighestAPIVersion)
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
		return nil, fmt.Errorf("discord requires data encoding to be of the following '%s', while '%s' encoding was requested", strings.Join(Encodings, "', '"), encoding)
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

	if conf.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// return configured discord websocket client
	return &Client{
		token:             conf.Token,
		urlAPIVersion:     BaseURL + "/v" + strconv.Itoa(conf.DAPIVersion),
		httpClient:        conf.HTTPClient,
		dAPIVersion:       conf.DAPIVersion,
		dAPIEncoding:      encoding,
		heartbeatAcquired: time.Now(),
		disconnected:      nil,
		iEventChan:        make(chan EventInterface),
		operationChan:     make(chan *gatewayEvent),
		eventChans:        make(map[string]chan []byte),
		sendChan:          make(chan *gatewayPayload),
		//Myself:            &user.User{},
	}, nil
}

// Client holds the web socket state and can be used directly in marshal/unmarshal to work with intance data
type Client struct {
	sync.RWMutex `json:"-"`

	urlAPIVersion string

	// URL Websocket URL web socket url
	url string

	httpClient *http.Client

	dAPIVersion    int
	dAPIEncoding   string
	token          string
	sequenceNumber uint

	HeartbeatInterval uint `json:"heartbeat_interval"`
	heartbeatAcquired time.Time
	Trace             []string     `json:"_trace"`
	SessionID         string       `json:"session_id"`
	ShardCount        uint         `json:"shard_count"`
	ShardID           snowflake.ID `json:"shard_id"`

	disconnected  chan struct{}
	operationChan chan *gatewayEvent
	eventChans    map[string]chan []byte
	sendChan      chan *gatewayPayload
	iEventChan    chan EventInterface

	//Myself         *user.User  `json:"user"`
	//MyselfSettings interface{} `json:"user_settings"`

	// websocket connection
	conn    *websocket.Conn
	wsMutex sync.Mutex // https://hackernoon.com/dancing-with-go-s-mutexes-92407ae927bf

	// heartbeat mutex keeps us from creating another pulser
	pulseMutex sync.Mutex
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

func (c *Client) GetEventChannel() <-chan EventInterface {
	return c.iEventChan
}
