package disgord

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/discordws"
	"github.com/andersfylling/disgord/guild"
	"github.com/andersfylling/disgord/request"
	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/snowflake"
	"github.com/sirupsen/logrus"
)

const (
	// APIComEncoding data format used when communicating with the discord API
	APIComEncoding string = "json"

	// APIVersion desired API version to use
	APIVersion int = 6 // February 5, 2018

	GitHubURL string = "https://github.com/andersfylling/disgord"

	Version string = "v0.2.0" // todo: eh?..
)

// Session the discord api is split in two. socket for keeping the client up to date, and http api for requests.
type Session interface {
	// main modules
	//

	// Request For interacting with Discord. Sending messages, creating channels, guilds, etc.
	// To read object state such as guilds, State() should be used in stead. However some data
	// might not exist in the state. If so it should be requested.
	Req() request.Requester

	// Event let's developers listen for specific events, event groups, or every event as one listener.
	// Supports both channels and callbacks
	Evt() EvtDispatcher

	// State reflects the latest changes received from Discord gateway.
	// Should be used instead of requesting objects.
	State() StateCacher

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
	Guild(guildID snowflake.ID) <-chan *guild.Guild
	Channel(channelID snowflake.ID) <-chan *channel.Channel
	Channels(guildID snowflake.ID) <-chan map[snowflake.ID]*channel.Channel
	Msg(msgID snowflake.ID) <-chan *channel.Message
	User(userID snowflake.ID) <-chan *user.User
	Member(guildID, userID snowflake.ID) <-chan *guild.Member
	Members(guildID snowflake.ID) <-chan map[snowflake.ID]*guild.Member
}

type Config struct {
	Token            string
	HTTPClient       *http.Client
	LoadAllMembers   bool
	LoadAllChannels  bool
	LoadAllRoles     bool
	LoadAllPresences bool
	Debug            bool
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
		APIVersion:         APIVersion,
		BotToken:           conf.Token,
		UserAgentSourceURL: GitHubURL,
		UserAgentVersion:   Version,
		HTTPClient:         conf.HTTPClient,
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
		state:         NewStateCache(evtDispatcher),
		req:           reqClient,
	}

	return c, nil
}

func NewClientMustConnect(conf *Config) *Client {
	client, err := NewClient(conf)
	if err != nil {
		panic(err)
	}

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	return client
}

func NewSession(conf *Config) (Session, error) {
	return NewClient(conf)
}

func NewSessionMustConnect(conf *Config) Session {
	return NewClientMustConnect(conf)
}

type Client struct {
	sync.RWMutex

	token string

	ws            *discordws.Client
	socketEvtChan <-chan discordws.EventInterface

	// register listeners for events
	evtDispatch *Dispatch

	// discord http api
	req *request.Client

	httpClient *http.Client

	// cache
	state *StateCache
}

func (c *Client) eventObserver() {
	for {
		select {
		case evt, alive := <-c.socketEvtChan:
			if !alive {
				logrus.Error("Event channel is dead!")
				break
			}

			ctx := context.Background()

			// TODO: parsing JSON uses panic and not logging on issues..

			eventName := evt.Name()
			data := evt.Data()

			// fan out to specific channel types
			go c.evtDispatch.trigger(eventName, c, ctx, data)
		}
	}
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
	go c.eventObserver()

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

func (c *Client) Req() request.Requester {
	return c.req
}

func (c *Client) Evt() EvtDispatcher {
	return c.evtDispatch
}

func (c *Client) State() StateCacher {
	return c.state
}

func (c *Client) AddListener(evtName string, listener interface{}) {
	c.evtDispatch.AddHandler(evtName, listener)
}

// AddListenerOnce not implemented. Do not use.
func (c *Client) AddListenerOnce(evtName string, listener interface{}) {
	c.evtDispatch.AddHandlerOnce(evtName, listener)
}

func (c *Client) Channel(channelID snowflake.ID) <-chan *channel.Channel {
	ch := make(chan *channel.Channel)

	go func(receiver chan<- *channel.Channel, storage StateCacher) {
		result := &channel.Channel{}
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

func (c *Client) Channels(GuildID snowflake.ID) <-chan map[snowflake.ID]*channel.Channel {
	ch := make(chan map[snowflake.ID]*channel.Channel)

	go func(receiver chan<- map[snowflake.ID]*channel.Channel, storage StateCacher) {
		result := make(map[snowflake.ID]*channel.Channel)
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
func (c *Client) Guild(guildID snowflake.ID) <-chan *guild.Guild {
	ch := make(chan *guild.Guild)

	go func(receiver chan<- *guild.Guild, storage StateCacher) {
		result := &guild.Guild{}
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
func (c *Client) Msg(msgID snowflake.ID) <-chan *channel.Message {
	ch := make(chan *channel.Message)

	go func(receiver chan<- *channel.Message, storage StateCacher) {
		result := &channel.Message{}
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
func (c *Client) User(userID snowflake.ID) <-chan *user.User {
	ch := make(chan *user.User)

	go func(userID snowflake.ID, receiver chan<- *user.User, storage StateCacher) {
		var result *user.User
		cached := true

		// check cache
		result, err := storage.User(userID)
		if err != nil {
			// log
			fmt.Printf("User not in cache: id: %s\n", userID.String())
		}

		// TODO: cache dead objects, to avoid http requesting the same none existance object?
		// will this ever be a problem

		// do http request if none found
		if result == nil {
			cached = false
			result = user.NewUser()
			err = c.req.Get("/users/"+userID.String(), result)
			if err != nil {
				fmt.Println("User does not exist in discord..")
				receiver <- nil
				close(receiver)
				return
			}
		}

		// return result
		receiver <- result

		// update cache with new result, if not found
		if !cached {
			storage.UserChan() <- result
		}

		// kill the channel
		close(receiver)
	}(userID, ch, c.state)

	return ch
}
func (c *Client) Member(guildID, userID snowflake.ID) <-chan *guild.Member {
	ch := make(chan *guild.Member)

	go func(receiver chan<- *guild.Member, storage StateCacher) {
		result := &guild.Member{}
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
func (c *Client) Members(guildID snowflake.ID) <-chan map[snowflake.ID]*guild.Member {
	ch := make(chan map[snowflake.ID]*guild.Member)

	go func(receiver chan<- map[snowflake.ID]*guild.Member, storage StateCacher) {
		result := make(map[snowflake.ID]*guild.Member)
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
