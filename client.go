package disgord

import (
	"context"
	"fmt"
	"net/http"
	"os/user"
	"sync"
	"time"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/discordws"
	"github.com/andersfylling/disgord/disgordctx"
	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/guild"
	"github.com/andersfylling/disgord/request"
	"github.com/andersfylling/snowflake"
	"github.com/sirupsen/logrus"
)

const (
	// APIComEncoding data encoding when communicating with the discord API
	APIComEncoding string = "json"

	// APIVersion desired API version to use
	APIVersion int = 6
)

// Session the discord api is split in two. socket for keeping the client up to date, and http api for requests.
type Session interface {
	// main modules
	//

	// Request For interacting with Discord. Sending messages, creating channels, guilds, etc.
	// To read object state such as guilds, State() should be used in stead. However some data
	// might not exist in the state. If so it should be requested.
	Request() request.Client

	// Event let's developers listen for specific events, event groups, or every event as one listener.
	// Supports both channels and callbacks
	Event() event.Dispatcher

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

	// event channels
	EvtChan(evt event.KeyType) <-chan interface{}

	// event callbacks
	//EvtAddListener(evt event.KeyType, callback interface{}) // use reflection based on keytype and cb params

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

	// create a disgord instance
	d := &Client{
		HTTPClient:    conf.HTTPClient,
		ws:            dws,
		socketEvtChan: dws.GetEventChannel(),
		Token:         conf.Token,
		Event:         event.NewDispatch(),
		State:         NewStateCache(),
	}

	return d, nil
}

// NewRequiredClient same as NewDisgord, but exits program if an error occours
func NewRequiredClient(conf *Config) *Client {
	dg, err := NewClient(conf)
	if err != nil {
		logrus.Fatal(err)
	}

	return dg
}

type Client struct {
	sync.RWMutex

	Token string

	ws *discordws.Client

	HTTPClient *http.Client

	socketEvtChan <-chan discordws.EventInterface

	// register listeners for events
	Event event.Dispatcher

	// cache
	State StateCacher
}

func (c *Client) eventObserver() {
	for {
		select {
		case evt, alive := <-c.socketEvtChan:
			if !alive {
				logrus.Error("Event channel is dead!")
				break
			}

			session := &disgordctx.Session{} //disgord context
			ctx := context.Background()

			// TODO: parsing JSON uses panic and not logging on issues..

			eventName := evt.Name()
			data := evt.Data()

			// fan out to specific channel types
			go c.Event.Trigger(eventName, session, ctx, data)
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
