package disgord

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"errors"
	"github.com/andersfylling/disgordws"
	"github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/disgord/rest"
	"github.com/andersfylling/disgord/rest/httd"
	"github.com/andersfylling/disgord/state"
	. "github.com/andersfylling/snowflake"
	"github.com/sirupsen/logrus"
)

// Session the discord api is split in two. socket for keeping the client up to date, and http api for requests.
type Session interface {
	// main modules
	//

	// Request For interacting with Discord. Sending messages, creating channels, guilds, etc.
	// To read object state such as guilds, State() should be used in stead. However some data
	// might not exist in the state. If so it should be requested. Note that this only holds http
	// CRUD operation and not the actual rest endpoints for discord (See Rest()).
	Req() httd.Requester

	// todo
	//Rest()

	// Event let's developers listen for specific events, event groups, or every event as one listener.
	// Supports both channels and callbacks
	Evt() EvtDispatcher

	// State reflects the latest changes received from Discord gateway.
	// Should be used instead of requesting objects.
	State() state.Cacher

	// RateLimiter the ratelimiter for the discord REST API
	RateLimiter() httd.RateLimiter

	// Discord Gateway, web socket
	//
	Connect() error
	Disconnect() error

	// module wrappers
	//

	// event callbacks
	AddListener(evtName string, callback interface{})
	AddListenerOnce(evtName string, callback interface{})

	// state/caching module
	// checks the cache first, otherwise do a http request
	Guild(guildID Snowflake) <-chan *resource.Guild
	Channel(channelID Snowflake) <-chan *resource.Channel
	Channels(guildID Snowflake) <-chan map[Snowflake]*resource.Channel
	Msg(msgID Snowflake) <-chan *resource.Message
	User(userID Snowflake) <-chan *UserChan
	Member(guildID, userID Snowflake) <-chan *resource.Member
	Members(guildID Snowflake) <-chan map[Snowflake]*resource.Member
}

type Config struct {
	Token      string
	HTTPClient *http.Client

	APIVersion  int    // eg. version 6. 0 defaults to lowest supported api version
	APIEncoding string // eg. json, use const. defaults to json

	CancelRequestWhenRateLimited bool

	LoadAllMembers   bool
	LoadAllChannels  bool
	LoadAllRoles     bool
	LoadAllPresences bool

	Debug bool
}

// NewClient creates a new default disgord instance
func NewClient(conf *Config) (*Client, error) {

	// ensure valid api version
	if conf.APIVersion == 0 {
		conf.APIVersion = 6 // the current discord API, for now v6
	}
	switch conf.APIVersion { // todo: simplify
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		fallthrough
	case 4:
		fallthrough
	case 5:
		return nil, errors.New("outdated API version")
	case 6: // supported
	default:
		return nil, errors.New("Discord API version is not yet supported")
	}

	if conf.HTTPClient == nil {
		// http client configuration
		conf.HTTPClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	// Use disgordws to keep the socket connection going
	// default communication encoding to json
	if conf.APIEncoding == "" {
		conf.APIEncoding = JSONEncoding
	}
	dws, err := disgordws.NewClient(&disgordws.Config{
		// user settings
		Token:      conf.Token,
		HTTPClient: conf.HTTPClient,
		Debug:      conf.Debug,

		// lib specific
		DAPIVersion:  conf.APIVersion,
		DAPIEncoding: conf.APIEncoding,
	})
	if err != nil {
		return nil, err
	}

	// request client
	reqConf := &httd.Config{
		APIVersion:                   conf.APIVersion,
		BotToken:                     conf.Token,
		UserAgentSourceURL:           GitHubURL,
		UserAgentVersion:             Version,
		HTTPClient:                   conf.HTTPClient,
		CancelRequestWhenRateLimited: conf.CancelRequestWhenRateLimited,
	}
	reqClient := httd.NewClient(reqConf)

	// event dispatcher
	evtDispatcher := NewDispatch()

	// create a disgord client/instance/session
	c := &Client{
		httpClient:    conf.HTTPClient,
		ws:            dws,
		socketEvtChan: dws.GetEventChannel(),
		token:         conf.Token,
		evtDispatch:   evtDispatcher,
		state:         state.NewCache(),
		req:           reqClient,
	}

	return c, nil
}

func NewClientMustCompile(conf *Config) *Client {
	client, err := NewClient(conf)
	if err != nil {
		panic(err)
	}

	return client
}

func NewSession(conf *Config) (Session, error) {
	return NewClient(conf)
}

func NewSessionMustCompile(conf *Config) Session {
	return NewClientMustCompile(conf)
}

type Client struct {
	sync.RWMutex

	token string

	ws            *disgordws.Client
	socketEvtChan <-chan disgordws.EventInterface

	// register listeners for events
	evtDispatch *Dispatch

	// cancelRequestWhenRateLimited by default the client waits until either the HTTPClient.timeout or
	// the rate limit ends before closing a request channel. If activated, in stead, requests will
	// instantly be denied, and the channel closed.
	cancelRequestWhenRateLimited bool

	// discord http api
	req *httd.Client

	httpClient *http.Client

	// cache
	state *state.Cache
}

func (c *Client) logInfo(msg string) {
	logrus.WithFields(logrus.Fields{
		"lib": c.ws.String(),
	}).Info(msg)
}

func (c *Client) logErr(msg string) {
	logrus.WithFields(logrus.Fields{
		"lib": c.ws.String(),
	}).Error(msg)
}

func (c *Client) String() string {
	return c.ws.String()
}

func (c *Client) RateLimiter() httd.RateLimiter {
	return c.req.RateLimiter()
}

// Connect establishes a websocket connection to the discord API
func (c *Client) Connect() (err error) {
	c.logInfo("Connecting to discord Gateway")
	err = c.ws.Connect()
	if err != nil {
		c.logErr(err.Error())
		return
	}
	c.logInfo("Connected")

	// setup event observer
	go c.eventHandler()

	return nil
}

// Disconnect closes the discord websocket connection
func (c *Client) Disconnect() (err error) {
	fmt.Println()
	c.logInfo("Closing Discord gateway connection")
	err = c.ws.Disconnect()
	if err != nil {
		c.logErr(err.Error())
		return
	}
	c.logInfo("Disconnected")

	return nil
}

func (c *Client) Req() httd.Requester {
	return c.req
}

func (c *Client) Evt() EvtDispatcher {
	return c.evtDispatch
}

func (c *Client) State() state.Cacher {
	return c.state
}

func (c *Client) AddListener(evtName string, listener interface{}) {
	c.evtDispatch.AddHandler(evtName, listener)
}

// AddListenerOnce not implemented. Do not use.
func (c *Client) AddListenerOnce(evtName string, listener interface{}) {
	c.evtDispatch.AddHandlerOnce(evtName, listener)
}

func (c *Client) Channel(channelID Snowflake) <-chan *resource.Channel {
	ch := make(chan *resource.Channel)

	go func(receiver chan<- *resource.Channel, storage *state.Cache) {
		result := &resource.Channel{}
		cached := true

		// check cache

		// do http request if none found
		if result == nil {
			cached = false
		}

		// return result
		receiver <- result

		// update cache with new result, if not found
		if !cached {
			//storage.MemberChan <- result
		}

		// kill the channel
		close(ch)
	}(ch, c.state)

	return ch
}

func (c *Client) Channels(GuildID Snowflake) <-chan map[Snowflake]*resource.Channel {
	ch := make(chan map[Snowflake]*resource.Channel)

	go func(receiver chan<- map[Snowflake]*resource.Channel, storage *state.Cache) {
		result := make(map[Snowflake]*resource.Channel)
		cached := true

		// check cache

		// do http request if none found
		if result == nil {
			cached = false
		}

		// return result
		receiver <- result

		// update cache with new result, if not found
		if !cached {
			//storage.MemberChan <- result
		}

		// kill the channel
		close(ch)
	}(ch, c.state)

	return ch
}

// state/caching module
func (c *Client) Guild(guildID Snowflake) <-chan *resource.Guild {
	ch := make(chan *resource.Guild)

	go func(receiver chan<- *resource.Guild, storage *state.Cache) {
		result := &resource.Guild{}
		cached := true

		// check cache

		// do http request if none found
		if result == nil {
			cached = false
		}

		// return result
		receiver <- result

		// update cache with new result, if not found
		if !cached {
			//storage.MemberChan <- result
		}

		// kill the channel
		close(ch)
	}(ch, c.state)

	return ch
}
func (c *Client) Msg(msgID Snowflake) <-chan *resource.Message {
	ch := make(chan *resource.Message)

	go func(receiver chan<- *resource.Message, storage *state.Cache) {
		result := &resource.Message{}
		cached := true

		// check cache

		// do http request if none found
		if result == nil {
			cached = false
		}

		// return result
		receiver <- result

		// update cache with new result, if not found
		if !cached {
			//storage.MemberChan <- result
		}

		// kill the channel
		close(ch)
	}(ch, c.state)

	return ch
}

type UserChan struct {
	User *resource.User
	Err error
	Cache bool
}
func (c *Client) User(userID Snowflake) <-chan *UserChan {
	ch := make(chan *UserChan)

	go func(userID Snowflake, receiver chan<- *UserChan, storage *state.Cache) {
		response := &UserChan{
			Cache: true,
		}

		// check cache
		response.User, response.Err = storage.User(userID)
		if response.Err != nil {
			response.Cache = false
			response.Err = nil
			response.User, response.Err = rest.GetUser(c.req, userID)
		}

		// TODO: cache dead objects, to avoid http requesting the same none existent object?
		// will this ever be a problem

		// return result
		receiver <- response

		// update cache with new result, if not found
		if !response.Cache && response.User != nil {
			storage.ProcessUser(&state.UserDetail{
				User:  response.User,
				Dirty: false,
			})
		}

		// kill the channel
		close(receiver)
	}(userID, ch, c.state)

	return ch
}
func (c *Client) Member(guildID, userID Snowflake) <-chan *resource.Member {
	ch := make(chan *resource.Member)

	go func(receiver chan<- *resource.Member, storage *state.Cache) {
		result := &resource.Member{}
		cached := true

		// check cache

		// do http request if none found
		if result == nil {
			cached = false
		}

		// return result
		receiver <- result

		// update cache with new result, if not found
		if !cached {
			//storage.MemberChan <- result
		}

		// kill the channel
		close(ch)
	}(ch, c.state)

	return ch
}
func (c *Client) Members(guildID Snowflake) <-chan map[Snowflake]*resource.Member {
	ch := make(chan map[Snowflake]*resource.Member)

	go func(receiver chan<- map[Snowflake]*resource.Member, storage *state.Cache) {
		result := make(map[Snowflake]*resource.Member)
		cached := true

		// check cache

		// do http request if none found
		if result == nil {
			cached = false
		}

		// return result
		receiver <- result

		// update cache with new result, if not found
		if !cached {
			//storage.MemberChan <- result
		}

		// kill the channel
		close(ch)
	}(ch, c.state)

	return ch
}
