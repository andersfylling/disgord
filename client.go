package disgord

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/andersfylling/disgord/logger"
	"github.com/andersfylling/disgord/websocket"

	"github.com/andersfylling/disgord/constant"
	"golang.org/x/net/proxy"

	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/httd"
)

// NewRESTClient creates a Client for sending and handling Discord protocols such as rate limiting
func NewRESTClient(conf *Config) (*httd.Client, error) {
	return httd.NewClient(&httd.Config{
		APIVersion:                   constant.DiscordVersion,
		BotToken:                     conf.BotToken,
		UserAgentSourceURL:           constant.GitHubURL,
		UserAgentVersion:             constant.Version,
		UserAgentExtra:               conf.ProjectName,
		HTTPClient:                   conf.HTTPClient,
		CancelRequestWhenRateLimited: conf.CancelRequestWhenRateLimited,
	})
}

// New create a Client. But panics on configuration/setup errors.
func New(conf *Config) (c *Client) {
	var err error
	if c, err = NewClient(conf); err != nil {
		panic(err)
	}

	return c

}

// NewClient creates a new DisGord Client and returns an error on configuration issues
func NewClient(conf *Config) (c *Client, err error) {
	if conf.HTTPClient == nil {
		conf.HTTPClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}
	if conf.Proxy != nil {
		conf.HTTPClient.Transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return conf.Proxy.Dial(network, addr)
			},
		}
	}

	if conf.ProjectName == "" {
		conf.ProjectName = LibraryInfo()
	}

	conf.shutdownChan = make(chan interface{})

	if conf.Logger == nil {
		conf.Logger = logger.Empty{}
	}

	// request Client for REST requests
	reqClient, err := NewRESTClient(conf)
	if err != nil {
		return nil, err
	}

	eventTracker := &websocket.UniqueStringSlice{}

	// caching
	var cacher *Cache
	if !conf.DisableCache {
		if conf.CacheConfig == nil {
			conf.CacheConfig = DefaultCacheConfig()
		} else {
			ensureBasicCacheConfig(conf.CacheConfig)
		}
		conf.CacheConfig.Log = conf.Logger
		cacher, err = newCache(conf.CacheConfig)
		if err != nil {
			return nil, err
		}

		// register for events for activate caches
		if !conf.CacheConfig.DisableUserCaching {
			eventTracker.Add(event.Ready)
			eventTracker.Add(event.UserUpdate)
		}
		if !conf.CacheConfig.DisableChannelCaching {
			eventTracker.Add(event.ChannelCreate)
			eventTracker.Add(event.ChannelUpdate)
			eventTracker.Add(event.ChannelPinsUpdate)
			eventTracker.Add(event.ChannelDelete)
		}
		if !conf.CacheConfig.DisableGuildCaching {
			eventTracker.Add(event.GuildCreate)
			eventTracker.Add(event.GuildDelete)
			eventTracker.Add(event.GuildUpdate)
			eventTracker.Add(event.GuildEmojisUpdate)
			eventTracker.Add(event.GuildMemberAdd)
			eventTracker.Add(event.GuildMemberRemove)
			eventTracker.Add(event.GuildMembersChunk)
			eventTracker.Add(event.GuildMemberUpdate)
			eventTracker.Add(event.GuildRoleCreate)
			eventTracker.Add(event.GuildRoleDelete)
			eventTracker.Add(event.GuildRoleUpdate)
			eventTracker.Add(event.GuildIntegrationsUpdate)
		}
	} else {
		// create an empty cache to avoid nil panics
		cacher, err = newCache(&CacheConfig{
			DisableUserCaching:       true,
			DisableChannelCaching:    true,
			DisableGuildCaching:      true,
			DisableVoiceStateCaching: true,
		})
		if err != nil {
			return nil, err
		}
	}
	cacher.conf.clientConf = conf

	// Required for voice operation
	eventTracker.Add(event.VoiceStateUpdate)
	eventTracker.Add(event.VoiceServerUpdate)

	// websocket sharding
	evtChan := make(chan *websocket.Event, 2) // TODO: higher value when more shards?

	// event dispatcher
	dispatch := newDispatcher()

	// create a disgord Client/instance/session
	c = &Client{
		shutdownChan: conf.shutdownChan,
		config:       conf,
		httpClient:   conf.HTTPClient,
		proxy:        conf.Proxy,
		botToken:     conf.BotToken,
		dispatcher:   dispatch,
		req:          reqClient,
		cache:        cacher,
		log:          conf.Logger,
		pool:         newPools(),
		eventChan:    evtChan,
		eventTracker: eventTracker,
	}
	c.dispatcher.addSessionInstance(c)
	c.voiceRepository = newVoiceRepository(c)

	return c, err
}

type ShardConfig = websocket.ShardConfig

// Config Configuration for the DisGord Client
type Config struct {
	BotToken   string
	HTTPClient *http.Client
	Proxy      proxy.Dialer

	CancelRequestWhenRateLimited bool

	DisableCache bool
	CacheConfig  *CacheConfig
	ShardConfig  ShardConfig
	Presence     *UpdateStatusCommand

	// DisGord triggers custom events to handle caching in a simpler manner.
	// eg. on guild delete, all the channels should be deleted too. to make this simple
	// disgord dispatches a channel delete event for all channels owned by the guild.
	//
	// Setting AcceptCustomEvents to true, will allow you to receive these. Note that
	// they can be detected with evt.ShardID == disgord.MockedShardID or evt.ShardID < 0.
	AcceptCustomEvents bool

	//ImmutableCache bool

	//LoadAllMembers   bool
	//LoadAllChannels  bool
	//LoadAllRoles     bool
	//LoadAllPresences bool

	// for cancellation
	shutdownChan chan interface{}

	// your project name, name of bot, or application
	ProjectName string

	// Logger is a dependency that must be injected to support logging.
	// disgord.DefaultLogger() can be used
	Logger Logger
}

// Client is the main disgord Client to hold your state and data. You must always initiate it using the constructor
// methods (eg. New(..) or NewClient(..)).
//
// Note that this Client holds all the REST methods, and is split across files, into whatever category
// the REST methods regards.
type Client struct {
	sync.RWMutex

	shutdownChan chan interface{}
	config       *Config
	botToken     string

	myID        Snowflake
	permissions PermissionBits

	// reactor demultiplexer for events
	dispatcher *dispatcher

	// cancelRequestWhenRateLimited by default the Client waits until either the HTTPClient.timeout or
	// the rate limit ends before closing a request channel. If activated, in stead, requests will
	// instantly be denied, and the process ended with a rate limited error.
	cancelRequestWhenRateLimited bool

	// req holds the rate limiting logic and error parsing unique for Discord
	req *httd.Client

	// http Client used for connections
	httpClient *http.Client
	proxy      proxy.Dialer

	shardManager websocket.ShardManager
	eventChan    chan *websocket.Event
	eventTracker *websocket.UniqueStringSlice

	connectedGuilds      []Snowflake
	connectedGuildsMutex sync.RWMutex

	cache *Cache

	log Logger

	// voice
	*voiceRepository

	// pools
	pool *pools
}

//////////////////////////////////////////////////////
//
// COMPLIANCE'S / IMPLEMENTATIONS
//
//////////////////////////////////////////////////////
var _ fmt.Stringer = (*Client)(nil)
var _ Session = (*Client)(nil)
var _ Link = (*Client)(nil)

//////////////////////////////////////////////////////
//
// METHODS
//
//////////////////////////////////////////////////////

func (c *Client) Pool() *pools {
	return c.pool
}

// AddPermission adds a minimum required permission to the bot. If the permission is negative, it is overwritten to 0.
// This is useful for creating the bot URL.
//
// At the moment, this holds no other effect than aesthetics.
func (c *Client) AddPermission(permission PermissionBit) (updatedPermissions PermissionBits) {
	if permission < 0 {
		permission = 0
	}

	c.permissions |= permission
	return c.GetPermissions()
}

// GetPermissions returns the minimum bot requirements.
func (c *Client) GetPermissions() (permissions PermissionBits) {
	return c.permissions
}

// CreateBotURL creates a URL that can be used to invite this bot to a guild/server.
// Note that it depends on the bot ID to be after the Discord update where the Client ID
// is the same as the Bot ID.
//
// By default the permissions will be 0, as in none. If you want to add/set the minimum required permissions
// for your bot to run successfully, you should utilise
//  Client.
func (c *Client) CreateBotURL() (u string, err error) {
	_, _ = c.GetCurrentUser() // update c.myID

	if c.myID.IsZero() {
		err = errors.New("unable to get bot id")
		return "", err
	}

	// make sure the snowflake is new enough to be used as a Client ID
	t, err := time.Parse("2006-01-02 15:04:05", "2016-08-07 05:39:21.906")
	if err != nil {
		return "", err
	}

	loc, _ := time.LoadLocation("America/Los_Angeles")
	t = t.In(loc)

	if !c.myID.Date().After(t) {
		err = errors.New("the bot was not created after " + t.String() + " and can therefore not use the bot ID to generate a invite link")
		return "", err
	}

	format := "https://discordapp.com/oauth2/authorize?scope=bot&client_id=%s&permissions=%d"
	u = fmt.Sprintf(format, c.myID.String(), c.permissions)
	return u, nil
}

// AvgHeartbeatLatency checks the duration of waiting before receiving a response from Discord when a
// heartbeat packet was sent. Note that heartbeats are usually sent around once a minute and is not a accurate
// way to measure delay between the Client and Discord server
func (c *Client) AvgHeartbeatLatency() (duration time.Duration, err error) {
	latencies, err := c.shardManager.HeartbeatLatencies()
	if err != nil {
		return 0, err
	}

	var average int64
	for _, v := range latencies {
		average += v.Nanoseconds()
	}
	average /= int64(len(latencies))

	return time.Duration(average) * time.Nanosecond, nil
}

// HeartbeatLatencies returns latencies mapped to each shard, by their respective ID. shardID => latency.
func (c *Client) HeartbeatLatencies() (latencies map[uint]time.Duration, err error) {
	return c.shardManager.HeartbeatLatencies()
}

// Myself get the current user / connected user
// Deprecated: use GetCurrentUser instead
func (c *Client) Myself() (user *User, err error) {
	return c.GetCurrentUser()
}

// GetConnectedGuilds get a list over guild IDs that this Client is "connected to"; or have joined through the ws connection. This will always hold the different Guild IDs, while the GetGuilds or GetCurrentUserGuilds might be affected by cache configuration.
func (c *Client) GetConnectedGuilds() []Snowflake {
	c.connectedGuildsMutex.RLock()
	defer c.connectedGuildsMutex.RUnlock()
	return c.connectedGuilds
}

// Logger returns the log instance of DisGord.
// Note that this instance is never nil. When the conf.Logger is not assigned
// an empty struct is used instead. Such that all calls are simply discarded at compile time
// removing the need for nil checks.
func (c *Client) Logger() logger.Logger {
	return c.log
}

func (c *Client) String() string {
	return LibraryInfo()
}

// RateLimiter return the rate limiter object
func (c *Client) RateLimiter() httd.RateLimiter {
	return c.req.RateLimiter()
}

// Req return the request object. Used in REST requests to handle rate limits,
// wrong http responses, etc.
func (c *Client) Req() httd.Requester {
	return c.req
}

// Cache returns the cacheLink manager for the session
func (c *Client) Cache() Cacher {
	return c.cache
}

//////////////////////////////////////////////////////
//
// Socket connection
//
//////////////////////////////////////////////////////

func (c *Client) setupConnectEnv() {
	// set the user ID upon connection
	// only works with socket logic
	c.On(event.UserUpdate, c.handlerUpdateSelfBot)
	c.On(event.GuildCreate, c.handlerAddToConnectedGuilds)
	c.On(event.GuildDelete, c.handlerRemoveFromConnectedGuilds)

	// start demultiplexer which also trigger dispatching
	var cache *Cache
	if !c.config.DisableCache {
		cache = c.cache
	}
	go demultiplexer(c.dispatcher, c.eventChan, cache)
}

// Connect establishes a websocket connection to the discord API
func (c *Client) Connect() (err error) {
	// set the user ID upon connection
	// only works for socketing
	//
	// also verifies that the correct credentials were supplied
	var me *User
	if me, err = c.GetCurrentUser(); err != nil {
		return err
	}
	c.myID = me.ID

	if err = websocket.ConfigureShardConfig(c, &c.config.ShardConfig); err != nil {
		return err
	}

	sharding := websocket.NewShardMngr(websocket.ShardManagerConfig{
		ShardConfig:        c.config.ShardConfig,
		Logger:             c.config.Logger,
		ShutdownChan:       c.config.shutdownChan,
		DefaultBotPresence: c.config.Presence,
		TrackedEvents:      c.eventTracker,
		EventChan:          c.eventChan,
		DisgordInfo:        LibraryInfo(),
		ProjectName:        c.config.ProjectName,
		BotToken:           c.config.BotToken,
	})

	c.setupConnectEnv()

	c.log.Info("Connecting to discord Gateway")
	if err = sharding.Connect(); err != nil {
		c.log.Info(err)
		return err
	}

	c.log.Info("Connected")
	c.shardManager = sharding
	return nil
}

// Disconnect closes the discord websocket connection
func (c *Client) Disconnect() (err error) {
	fmt.Println() // to keep ^C on it's own line
	c.log.Info("Closing Discord gateway connection")
	close(c.dispatcher.shutdown)
	if err = c.shardManager.Disconnect(); err != nil {
		c.log.Error(err)
		return err
	}
	close(c.shutdownChan)
	c.log.Info("Disconnected")

	return nil
}

// Suspend in case you want to temporary disconnect from the Gateway. But plan on
// connecting again without restarting your software/application, this should be used.
func (c *Client) Suspend() (err error) {
	c.log.Info("Closing Discord gateway connection")
	if err = c.shardManager.Disconnect(); err != nil {
		return err
	}
	c.log.Info("Suspended")

	return nil
}

// DisconnectOnInterrupt wait until a termination signal is detected
func (c *Client) DisconnectOnInterrupt() (err error) {
	// create a channel to listen for termination signals (graceful shutdown)
	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-termSignal

	return c.Disconnect()
}

// StayConnectedUntilInterrupted is a simple wrapper for connect, and disconnect that listens for system interrupts.
// When a error happens you can terminate the application without worries.
func (c *Client) StayConnectedUntilInterrupted() (err error) {
	if err = c.Connect(); err != nil {
		c.log.Error(err)
		return err
	}

	if err = c.DisconnectOnInterrupt(); err != nil {
		c.log.Error(err)
		return err
	}

	return nil
}

//////////////////////////////////////////////////////
//
// Internal event handlers
//
//////////////////////////////////////////////////////

// handlerAddToConnectedGuilds update internal state when joining or creating a guild
func (c *Client) handlerAddToConnectedGuilds(_ Session, evt *GuildCreate) {
	c.connectedGuildsMutex.Lock()
	defer c.connectedGuildsMutex.Unlock()

	// don't add an entry if there already is one
	for i := range c.connectedGuilds {
		if c.connectedGuilds[i] == evt.Guild.ID {
			return
		}
	}

	c.connectedGuilds = append(c.connectedGuilds, evt.Guild.ID)
}

// handlerRemoveFromConnectedGuilds update internal state when deleting or leaving a guild
func (c *Client) handlerRemoveFromConnectedGuilds(_ Session, evt *GuildDelete) {
	c.connectedGuildsMutex.Lock()
	defer c.connectedGuildsMutex.Unlock()

	guilds := c.connectedGuilds
	for i := range guilds {
		if guilds[i] != evt.UnavailableGuild.ID {
			continue
		}
		guilds[i] = guilds[len(guilds)-1]
		guilds = guilds[:len(guilds)-1]
		break
	}

	c.connectedGuilds = guilds
}

func (c *Client) handlerUpdateSelfBot(_ Session, update *UserUpdate) {
	_ = c.cache.Update(UserCache, update.User)
}

//////////////////////////////////////////////////////
//
// Socket utilities
//
//////////////////////////////////////////////////////

// Ready triggers a given callback when all shards has gotten their first Ready event
// Warning: Do not call Client.Connect before this.
func (c *Client) Ready(cb func()) {
	ctrl := &rdyCtrl{
		cb: cb,
	}

	c.On(EvtReady, func(s Session, evt *Ready) {
		ctrl.Lock()
		defer ctrl.Unlock()

		l := c.shardManager.NrOfShards()
		if l != uint(len(ctrl.shardReady)) {
			ctrl.shardReady = make([]bool, l)
		}

		ctrl.shardReady[evt.ShardID] = true
	}, ctrl)
}

// On creates a specification to be executed on the given event. The specification
// consists of, in order, 0 or more middlewares, 1 or more handlers, 0 or 1 controller.
// On incorrect ordering, or types, the method will panic. See reactor.go for types.
//
// Each of the three sub-types of a specification is run in sequence, as well as the specifications
// registered for a event. However, the slice of specifications are executed in a goroutine to avoid
// blocking future events. The middlewares allows manipulating the event data before it reaches the
// handlers. The handlers executes short-running logic based on the event data (use go routine if
// you need a long running task). The controller dictates lifetime of the specification.
//
//  // a handler that is executed on every Ready event
//  Client.On(EvtReady, onReady)
//
//  // a handler that runs only the first three times a READY event is fired
//  Client.On(EvtReady, onReady, &Ctrl{Runs: 3})
//
//  // a handler that only runs for events within the first 10 minutes
//  Client.On(EvtReady, onReady, &Ctrl{Duration: 10*time.Minute})
//
// Another example is to create a voting system where you specify a deadline instead of a Runs counter:
//  On("MESSAGE_CREATE", mdlwHasMentions, handleMsgsWithMentions, saveVoteToDB, &Ctrl{Until:time.Now().Add(time.Hour)})
//
// You can use your own Ctrl struct, as long as it implements disgord.HandlerCtrl. Do not execute long running tasks
// in the methods. Use a go routine instead.
//
// If the HandlerCtrl.OnInsert returns an error, the related handlers are still added to the dispatcher.
// But the error is logged to the injected logger instance (log.Error).
//
// This ctrl feature was inspired by https://github.com/discordjs/discord.js
func (c *Client) On(event string, inputs ...interface{}) {
	if err := ValidateHandlerInputs(inputs...); err != nil {
		panic(err)
	}
	c.eventTracker.Add(event)

	if err := c.dispatcher.register(event, inputs...); err != nil {
		panic(err)
	}
}

// Emit sends a socket command directly to Discord.
func (c *Client) Emit(command SocketCommand, data interface{}) error {
	switch command {
	case CommandUpdateStatus, CommandUpdateVoiceState, CommandRequestGuildMembers:
	default:
		return errors.New("command is not supported")
	}

	if g, ok := data.(guilder); ok {
		// if this is guild specific, then only send data through the related shard
		guildID := g.getGuildID()
		shardID := GetShardForGuildID(guildID, c.shardManager.NrOfShards())
		shard, err := c.shardManager.GetShard(shardID)
		if err != nil {
			return err
		}
		return shard.Emit(command, data)
	}

	// otherwise it is sent through every shard
	return c.shardManager.Emit(command, data)
}

// AcceptEvent only events registered using this method is accepted from the Discord socket API. The rest is discarded
// to reduce unnecessary marshalling and controls.
func (c *Client) AcceptEvent(events ...string) {
	for _, evt := range events {
		c.eventTracker.Add(evt)
	}
}

//////////////////////////////////////////////////////
//
// Abstract CRUD operations
//
//////////////////////////////////////////////////////

// DeleteFromDiscord if the given object has implemented the private interface discordDeleter this method can
// be used to delete said object.
func (c *Client) DeleteFromDiscord(obj discordDeleter, flags ...Flag) (err error) {
	if obj == nil {
		return errors.New("object to save can not be nil")
	}

	err = obj.deleteFromDiscord(c, flags...)
	return
}

//////////////////////////////////////////////////////
//
// REST Methods
// customs
//
//////////////////////////////////////////////////////

func (c *Client) GetGuilds(params *GetCurrentUserGuildsParams, flags ...Flag) ([]*Guild, error) {
	// TODO: populate these partial guild objects
	return c.GetCurrentUserGuilds(params)
}

func (c *Client) KickVoiceParticipant(guildID, userID Snowflake) error {
	return c.UpdateGuildMember(guildID, userID).KickFromVoice().Execute()
}

// SendMsg Input anything and it will be converted to a message and sent. If you
// supply it with multiple data's, it will simply merge them. Even if they are multiple Message objects.
// However, if you supply multiple CreateMessageParams objects, you will face issues. But at this point
// you really need to reconsider your own code.
//
// Note that sending a &Message will simply refer to it, and not copy over the contents into
// the reply. example output: message{6434732342356}
//
// If you want to affect the actual message data besides .Content; provide a
// MessageCreateParams. The reply message will be updated by the last one provided.
func (c *Client) SendMsg(channelID Snowflake, data ...interface{}) (msg *Message, err error) {

	var flags []Flag
	params := &CreateMessageParams{}
	for i := range data {
		if data[i] == nil {
			continue
		}

		var s string
		switch t := data[i].(type) {
		case *CreateMessageParams:
			*params = *t
		case CreateMessageParams:
			*params = t
		case string:
			s = t
		case *Flag:
			flags = append(flags, *t)
		case Flag:
			flags = append(flags, t)
		default:
			if str, ok := t.(fmt.Stringer); ok {
				s = str.String()
			} else {
				s = fmt.Sprint(t)
			}
		}

		if s != "" {
			params.Content += " " + s
		}
	}

	// wtf?
	if data == nil {
		if mergeFlags(flags).IgnoreEmptyParams() {
			params.Content = ""
		} else {
			return nil, errors.New("params were nil")
		}
	}

	return c.CreateMessage(channelID, params, flags...)
}

/* status updates */

// UpdateStatus updates the Client's game status
// note: for simple games, check out UpdateStatusString
func (c *Client) UpdateStatus(s *UpdateStatusCommand) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return c.shardManager.Emit(CommandUpdateStatus, s)
}

// UpdateStatusString sets the Client's game activity to the provided string, status to online
// and type to Playing
func (c *Client) UpdateStatusString(s string) error {
	updateData := &UpdateStatusCommand{
		Since: nil,
		Game: &Activity{
			Name: s,
			Type: 0,
		},
		Status: StatusOnline,
		AFK:    false,
	}
	return c.UpdateStatus(updateData)
}

func (c *Client) newRESTRequest(conf *httd.Request, flags []Flag) *rest {
	r := &rest{
		c:    c,
		conf: conf,
	}
	r.init()
	r.flags = mergeFlags(flags)

	return r
}
