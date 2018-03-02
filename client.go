package disgord

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/andersfylling/disgord/discordws"
	"github.com/sirupsen/logrus"
	"github.com/andersfylling/disgord/request"
	"github.com/andersfylling/disgord/guild"
	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/snowflake"
	"github.com/andersfylling/disgord/user"
)

const (
	// APIComEncoding data format used when communicating with the discord API
	APIComEncoding string = "json"

	// APIVersion desired API version to use
	APIVersion int = 6 // February 5, 2018
)


// Session the discord api is split in two. socket for keeping the client up to date, and http api for requests.
type Session interface {
	// main modules
	//

	// Request For interacting with Discord. Sending messages, creating channels, guilds, etc.
	// To read object state such as guilds, State() should be used in stead. However some data
	// might not exist in the state. If so it should be requested.
	Req() *request.Client

	// Event let's developers listen for specific events, event groups, or every event as one listener.
	// Supports both channels and callbacks
	Evt() EvtDispatcher

	// State reflects the latest changes received from Discord gateway.
	// Should be used instead of requesting objects.
	State() StateCacher

	// module wrappers
	//

	// requests
	ReqGuild(guildID snowflake.ID) *guild.Guild
	ReqChannel(channelID snowflake.ID) *channel.Channel
	ReqChannels(guildID snowflake.ID) map[snowflake.ID]*channel.Channel
	ReqMsg(msgID snowflake.ID) *channel.Message
	ReqUser(userID snowflake.ID) *user.User
	ReqMember(guildID, userID snowflake.ID) *guild.Member
	ReqMembers(guildID snowflake.ID) map[snowflake.ID]*guild.Member

	// event callbacks
	EvtAddHandler(evtName string, callback interface{}) // use reflection based on keytype and cb params

	// state/caching module
	Guild(guildID snowflake.ID) *guild.Guild
	Channel(channelID snowflake.ID) *channel.Channel
	Channels(guildID snowflake.ID) map[snowflake.ID]*channel.Channel
	Msg(msgID snowflake.ID) *channel.Message
	User(userID snowflake.ID) *user.User
	Member(guildID, userID snowflake.ID) *guild.Member
	Members(guildID snowflake.ID) map[snowflake.ID]*guild.Member
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
func NewClient(conf *Config) *Client {

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
		logrus.Fatal(err)
	}

	// create a disgord instance
	c := &Client{
		httpClient:    conf.HTTPClient,
		ws:            dws,
		socketEvtChan: dws.GetEventChannel(),
		token:         conf.Token,
		evtDispatch:   NewDispatch(),
		state:         NewStateCache(),
	}

	return c
}


func NewSession(conf *Config) Session {
	return NewClient(conf)
}

type Client struct {
	sync.RWMutex

	token string

	ws *discordws.Client
	socketEvtChan <-chan discordws.EventInterface

	// register listeners for events
	evtDispatch *Dispatch

	// discord http api
	req *request.Client

	httpClient *http.Client

	// cache
	state StateCacher
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


func (c *Client) Req() *request.Client {
	return c.req
}

func (c *Client) Evt() EvtDispatcher {
	return c.evtDispatch
}

func (c *Client) State() StateCacher {
	return c.state
}


func (c *Client) ReqGuild(guildID snowflake.ID) *guild.Guild {
	return nil
}
func (c *Client) ReqChannel(channelID snowflake.ID) *channel.Channel {
	return nil
}
func (c *Client) ReqChannels(guildID snowflake.ID) map[snowflake.ID]*channel.Channel {
	result := make(map[snowflake.ID]*channel.Channel)
	return result
}
func (c *Client) ReqMsg(msgID snowflake.ID) *channel.Message {
	return nil
}
func (c *Client) ReqUser(userID snowflake.ID) *user.User {
	return nil
}
func (c *Client) ReqMember(guildID, userID snowflake.ID) *guild.Member {
	return nil
}
func (c *Client) ReqMembers(guildID snowflake.ID) map[snowflake.ID]*guild.Member {
	result := make(map[snowflake.ID]*guild.Member)
	return result
}

func (c *Client) EvtAddHandler(evtName string, listener interface{}) {
	c.evtDispatch.AddHandler(evtName, listener)
}

func (c *Client) Channel(channelID snowflake.ID) *channel.Channel {
	return nil
}

func (c *Client) Channels(GuildID snowflake.ID) map[snowflake.ID]*channel.Channel {
	result := make(map[snowflake.ID]*channel.Channel)
	return result
}


// state/caching module
func (c *Client) Guild(guildID snowflake.ID) *guild.Guild {
	return nil
}
func (c *Client) Msg(msgID snowflake.ID) *channel.Message {
	return nil
}
func (c *Client) User(userID snowflake.ID) *user.User {
	return nil
}
func (c *Client) Member(guildID, userID snowflake.ID) *guild.Member {
	return nil
}
func (c *Client) Members(guildID snowflake.ID) map[snowflake.ID]*guild.Member {
	result := make(map[snowflake.ID]*guild.Member)
	return result
}