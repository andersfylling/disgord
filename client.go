package disgord

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/andersfylling/disgord/discordws"
	"github.com/andersfylling/disgord/request"
	"github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/disgord/state"
	"github.com/andersfylling/snowflake"
	"github.com/sirupsen/logrus"
)

// Session the discord api is split in two. socket for keeping the client up to date, and http api for requests.
type Session interface {
	// main modules
	//

	// Request For interacting with Discord. Sending messages, creating channels, guilds, etc.
	// To read object state such as guilds, State() should be used in stead. However some data
	// might not exist in the state. If so it should be requested.
	Req() request.DiscordRequester

	// Event let's developers listen for specific events, event groups, or every event as one listener.
	// Supports both channels and callbacks
	Evt() EvtDispatcher

	// State reflects the latest changes received from Discord gateway.
	// Should be used instead of requesting objects.
	State() state.Cacher

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
	Guild(guildID snowflake.ID) <-chan *resource.Guild
	Channel(channelID snowflake.ID) <-chan *resource.Channel
	Channels(guildID snowflake.ID) <-chan map[snowflake.ID]*resource.Channel
	Msg(msgID snowflake.ID) <-chan *resource.Message
	User(userID snowflake.ID) <-chan *resource.User
	Member(guildID, userID snowflake.ID) <-chan *resource.Member
	Members(guildID snowflake.ID) <-chan map[snowflake.ID]*resource.Member
}

type Config struct {
	Token      string
	HTTPClient *http.Client

	CancelRequestWhenRateLimited bool

	LoadAllMembers   bool
	LoadAllChannels  bool
	LoadAllRoles     bool
	LoadAllPresences bool

	Debug bool
}

// NewClient creates a new default disgord instance
func NewClient(conf *Config) (*Client, error) {

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
		DAPIVersion:  APIVersion,
		DAPIEncoding: APIComEncoding,
	})
	if err != nil {
		return nil, err
	}

	// request client
	reqConf := &request.Config{
		APIVersion:                   APIVersion,
		BotToken:                     conf.Token,
		UserAgentSourceURL:           GitHubURL,
		UserAgentVersion:             Version,
		HTTPClient:                   conf.HTTPClient,
		CancelRequestWhenRateLimited: conf.CancelRequestWhenRateLimited,
	}
	reqClient := request.NewClient(reqConf)

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

	ws            *discordws.Client
	socketEvtChan <-chan discordws.EventInterface

	// register listeners for events
	evtDispatch *Dispatch

	// cancelRequestWhenRateLimited by default the client waits until either the HTTPClient.timeout or
	// the rate limit ends before closing a request channel. If activated, in stead, requests will
	// instantly be denied, and the channel closed.
	cancelRequestWhenRateLimited bool

	// discord http api
	req *request.Client

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

func (c *Client) Req() request.DiscordRequester {
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

func (c *Client) Channel(channelID snowflake.ID) <-chan *resource.Channel {
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

func (c *Client) Channels(GuildID snowflake.ID) <-chan map[snowflake.ID]*resource.Channel {
	ch := make(chan map[snowflake.ID]*resource.Channel)

	go func(receiver chan<- map[snowflake.ID]*resource.Channel, storage *state.Cache) {
		result := make(map[snowflake.ID]*resource.Channel)
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
func (c *Client) Guild(guildID snowflake.ID) <-chan *resource.Guild {
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
func (c *Client) Msg(msgID snowflake.ID) <-chan *resource.Message {
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
func (c *Client) User(userID snowflake.ID) <-chan *resource.User {
	ch := make(chan *resource.User)

	go func(userID snowflake.ID, receiver chan<- *resource.User, storage *state.Cache) {
		var result *resource.User
		var err error
		cached := true

		// check cache
		result, err = storage.User(userID)
		if err != nil {
			// log
			fmt.Printf("User not in cache: id: %s\n", userID.String())
		}

		// TODO: cache dead objects, to avoid http requesting the same none existent object?
		// will this ever be a problem

		// do http request if none found
		if result == nil {
			cached = false
			result, err = resource.ReqUser(c.req, userID)
			if err != nil {
				// TODO: handle error
				// issue: devs might either be rate limited or user not found, how would they know tho?
				receiver <- nil
				close(receiver)
				return
			}
		}

		// return result
		receiver <- result

		// update cache with new result, if not found
		if !cached {
			storage.ProcessUser(&state.UserDetail{
				User:  result,
				Dirty: false,
			})
		}

		// kill the channel
		close(receiver)
	}(userID, ch, c.state)

	return ch
}
func (c *Client) Member(guildID, userID snowflake.ID) <-chan *resource.Member {
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
func (c *Client) Members(guildID snowflake.ID) <-chan map[snowflake.ID]*resource.Member {
	ch := make(chan map[snowflake.ID]*resource.Member)

	go func(receiver chan<- map[snowflake.ID]*resource.Member, storage *state.Cache) {
		result := make(map[snowflake.ID]*resource.Member)
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
