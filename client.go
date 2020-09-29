package disgord

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/andersfylling/disgord/internal/gateway"
	"github.com/andersfylling/disgord/internal/logger"

	"golang.org/x/net/proxy"

	"github.com/andersfylling/disgord/internal/constant"

	"github.com/andersfylling/disgord/internal/httd"
)

// New create a Client. But panics on configuration/setup errors.
func New(conf Config) *Client {
	client, err := createClient(&conf)
	if err != nil {
		panic(err)
	}
	return client
}

// NewClient creates a new Disgord Client and returns an error on configuration issues
func NewClient(conf Config) (*Client, error) {
	return createClient(&conf)
}

// NewClient creates a new Disgord Client and returns an error on configuration issues
func createClient(conf *Config) (c *Client, err error) {
	if conf.Presence != nil {
		if _, err := gateway.StringToStatusType(conf.Presence.Status); err != nil {
			return nil, fmt.Errorf("use a disgord value eg. disgord.StatusOnline: %w", err)
		}
	}
	if conf.HTTPClient == nil {
		// WARNING: do not set http.Client.Timeout (!)
		conf.HTTPClient = &http.Client{}
	} else if conf.HTTPClient.Timeout > 0 {
		// https://github.com/nhooyr/websocket/issues/67
		return nil, errors.New("do not set timeout in the http.Client, use context.Context instead")
	}
	if conf.Proxy != nil {
		conf.HTTPClient.Transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return conf.Proxy.Dial(network, addr)
			},
		}
	}

	httdClient, err := httd.NewClient(&httd.Config{
		APIVersion:                   constant.DiscordVersion,
		BotToken:                     conf.BotToken,
		UserAgentSourceURL:           constant.GitHubURL,
		UserAgentVersion:             constant.Version,
		UserAgentExtra:               conf.ProjectName,
		HTTPClient:                   conf.HTTPClient,
		CancelRequestWhenRateLimited: conf.CancelRequestWhenRateLimited,
		RESTBucketManager:            conf.RESTBucketManager,
	})
	if err != nil {
		return nil, err
	}

	if conf.ProjectName == "" {
		conf.ProjectName = LibraryInfo()
	}

	conf.shutdownChan = make(chan interface{})

	if conf.Logger == nil {
		conf.Logger = logger.Empty{}
	}

	// ignore PRESENCES_REPLACE: https://github.com/discord/discord-api-docs/issues/683
	conf.IgnoreEvents = append(conf.IgnoreEvents, "PRESENCES_REPLACE")

	// caching
	var cache Cache
	if conf.DisableCache {
		if _, ok := conf.Cache.(*CacheNop); !ok {
			cache = &CacheNop{}
		} else {
			cache = conf.Cache
		}
	} else if conf.Cache == nil {
		cache = NewCacheLFUImmutable(5000000, 10000, 0, 0)
	} else {
		cache = conf.Cache
	}

	// websocket sharding
	evtChan := make(chan *gateway.Event, 2) // TODO: higher value when more shards?

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
		req:          httdClient,
		cache:        cache,
		log:          conf.Logger,
		pool:         newPools(),
		eventChan:    evtChan,
	}
	c.handlers.c = c // parent reference
	c.dispatcher.addSessionInstance(c)
	c.clientQueryBuilder.client = c
	c.voiceRepository = newVoiceRepository(c)

	return c, err
}

type ShardConfig = gateway.ShardConfig

// Config Configuration for the Disgord Client
type Config struct {
	// ################################################
	// ##
	// ## Basic bot configuration.
	// ## This section is for everyone. And beginners
	// ## should stick to this section unless they know
	// ## what they are doing.
	// ##
	// ################################################
	BotToken   string
	HTTPClient *http.Client
	Proxy      proxy.Dialer

	// your project name, name of bot, or application
	ProjectName string

	// AlwaysParseChannelMentions will ensure that every message populates the
	// Message.ChannelsMentions, regardless of the Discord conditions.
	// AlwaysParseChannelMentions bool
	// TODO

	CancelRequestWhenRateLimited bool

	// LoadMembersQuietly will start fetching members for all Guilds in the background.
	// There is currently no proper way to detect when the loading is done nor if it
	// finished successfully.
	LoadMembersQuietly bool

	// Presence will automatically be emitted to discord on start up
	Presence *UpdateStatusPayload

	// for cancellation
	shutdownChan chan interface{}

	// Logger is a dependency that must be injected to support logging.
	// disgord.DefaultLogger() can be used
	Logger Logger

	// ################################################
	// ##
	// ## WARNING! For advanced Users only.
	// ## This section of options might break the bot,
	// ## make it incoherent to the Discord API requirements,
	// ## potentially causing your bot to be banned.
	// ## You use these features on your own risk.
	// ##
	// ################################################
	RESTBucketManager httd.RESTBucketManager

	DisableCache bool
	Cache        Cache
	ShardConfig  ShardConfig

	// IgnoreEvents will skip events that matches the given event names.
	// WARNING! This can break your caching, so be careful about what you want to ignore.
	//
	// Note this also triggers discord optimizations behind the scenes, such that disgord_diagnosews might
	// seem to be missing some events. But actually the lack of certain events will mean Discord aren't sending
	// them at all due to how the identify command was defined. eg. guildS_subscriptions
	IgnoreEvents []string

	Intents gateway.Intent
}

// Client is the main disgord Client to hold your state and data. You must always initiate it using the constructor
// methods (eg. New(..) or NewClient(..)).
//
// Note that this Client holds all the REST methods, and is split across files, into whatever category
// the REST methods regards.
type Client struct {
	sync.RWMutex

	clientQueryBuilder

	shutdownChan chan interface{}
	config       *Config
	botToken     string

	currentUser User
	myID        Snowflake
	permissions PermissionBit

	handlers internalHandlers

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

	shardManager gateway.ShardManager
	eventChan    chan *gateway.Event

	connectedGuilds      []Snowflake
	connectedGuildsMutex sync.RWMutex

	cache Cache

	log Logger

	// voice
	*voiceRepository

	// pools
	pool *pools
}

//////////////////////////////////////////////////////
//
// IMPLEMENTED INTERFACES
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
func (c *Client) AddPermission(permission PermissionBit) (updatedPermissions PermissionBit) {
	if permission < 0 {
		permission = 0
	}

	c.permissions |= permission
	return c.GetPermissions()
}

// GetPermissions returns the minimum bot requirements.
func (c *Client) GetPermissions() (permissions PermissionBit) {
	return c.permissions
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

// GetConnectedGuilds get a list over guild IDs that this Client is "connected to"; or have joined through the ws connection. This will always hold the different Guild IDs, while the GetGuilds or GetCurrentUserGuilds might be affected by cache configuration.
func (c *Client) GetConnectedGuilds() []Snowflake {
	c.connectedGuildsMutex.RLock()
	defer c.connectedGuildsMutex.RUnlock()

	guildIDs := make([]Snowflake, len(c.connectedGuilds))
	copy(guildIDs, c.connectedGuilds)

	return guildIDs
}

// Logger returns the log instance of Disgord.
// Note that this instance is never nil. When the conf.Logger is not assigned
// an empty struct is used instead. Such that all calls are simply discarded at compile time
// removing the need for nil checks.
func (c *Client) Logger() logger.Logger {
	return c.log
}

func (c *Client) String() string {
	return LibraryInfo()
}

// RESTBucketGrouping shows which hashed endpoints belong to which bucket hash for the REST API.
// Note that these bucket hashes are eventual consistent.
func (c *Client) RESTRatelimitBuckets() (group map[string][]string) {
	return c.req.BucketGrouping()
}

// Req return the request object. Used in REST requests to handle rate limits,
// wrong http responses, etc.
func (c *Client) Req() httd.Requester {
	return c.req
}

// Cache returns the cacheLink manager for the session
func (c *Client) Cache() Cache {
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
	if c.config.LoadMembersQuietly {
		c.On(EvtReady, c.handlers.loadMembers)
	}
	c.On(EvtGuildCreate, c.handlers.saveGuildID)
	c.On(EvtGuildDelete, c.handlers.deleteGuildID)

	// start demultiplexer which also trigger dispatching
	go c.demultiplexer(c.dispatcher, c.eventChan)
}

type helperGatewayBotGetter struct {
	c *Client
}

var _ gateway.GatewayBotGetter = (*helperGatewayBotGetter)(nil)

func (h helperGatewayBotGetter) GetGatewayBot(ctx context.Context) (gateway *gateway.GatewayBot, err error) {
	return h.c.WithContext(ctx).GetGatewayBot()
}

// Connect establishes a websocket connection to the discord API
func (c *Client) Connect(ctx context.Context) (err error) {
	// set the user ID upon connection
	// only works for socketing
	//
	// also verifies that the correct credentials were supplied

	// Avoid races during connection setup
	c.Lock()
	defer c.Unlock()

	var me *User
	if me, err = c.CurrentUser().WithContext(ctx).Get(); err != nil {
		return err
	}
	c.myID = me.ID

	if err = gateway.ConfigureShardConfig(ctx, helperGatewayBotGetter{c}, &c.config.ShardConfig); err != nil {
		return err
	}

	shardMngrConf := gateway.ShardManagerConfig{
		ShardConfig:  c.config.ShardConfig,
		Logger:       c.config.Logger,
		ShutdownChan: c.config.shutdownChan,
		IgnoreEvents: c.config.IgnoreEvents,
		Intents:      c.config.Intents,
		EventChan:    c.eventChan,
		DisgordInfo:  LibraryInfo(),
		ProjectName:  c.config.ProjectName,
		BotToken:     c.config.BotToken,
	}

	if c.config.Presence != nil {
		// assumption: error is handled when creating a new client
		status, _ := gateway.StringToStatusType(c.config.Presence.Status)
		shardMngrConf.DefaultBotPresence = &gateway.UpdateStatusPayload{
			Since:  c.config.Presence.Since,
			Game:   c.config.Presence.Game,
			Status: status,
			AFK:    c.config.Presence.AFK,
		}
	}

	sharding := gateway.NewShardMngr(shardMngrConf)

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
	// catches panic when being called as a deferred function
	if r := recover(); r != nil {
		panic("unable to connect due to above error")
	}

	<-CreateTermSigListener()
	return c.Disconnect()
}

// StayConnectedUntilInterrupted is a simple wrapper for connect, and disconnect that listens for system interrupts.
// When a error happens you can terminate the application without worries.
func (c *Client) StayConnectedUntilInterrupted(ctx context.Context) (err error) {
	// catches panic when being called as a deferred function
	if r := recover(); r != nil {
		panic("unable to connect due to above error")
	}

	if err = c.Connect(ctx); err != nil {
		c.log.Error(err)
		return err
	}

	select {
	case <-CreateTermSigListener():
	case <-ctx.Done():
	}

	return c.Disconnect()
}

//////////////////////////////////////////////////////
//
// Internal event handlers
//
//////////////////////////////////////////////////////

type internalHandlers struct {
	c *Client
}

// saveGuildID update internal state when joining or creating a guild
func (ih *internalHandlers) saveGuildID(_ Session, evt *GuildCreate) {
	client := ih.c
	client.connectedGuildsMutex.Lock()
	defer client.connectedGuildsMutex.Unlock()

	// don't add an entry if there already is one
	for i := range client.connectedGuilds {
		if client.connectedGuilds[i] == evt.Guild.ID {
			return
		}
	}

	client.connectedGuilds = append(client.connectedGuilds, evt.Guild.ID)
}

// deleteGuildID update internal state when deleting or leaving a guild
func (ih *internalHandlers) deleteGuildID(_ Session, evt *GuildDelete) {
	client := ih.c
	client.connectedGuildsMutex.Lock()
	defer client.connectedGuildsMutex.Unlock()

	guilds := client.connectedGuilds
	for i := range guilds {
		if guilds[i] != evt.UnavailableGuild.ID {
			continue
		}
		guilds[i] = guilds[len(guilds)-1]
		guilds = guilds[:len(guilds)-1]
		break
	}

	client.connectedGuilds = guilds
}

func (ih *internalHandlers) loadMembers(_ Session, evt *Ready) {
	client := ih.c
	guildIDs := make([]Snowflake, len(evt.Guilds))
	for i := range evt.Guilds {
		guildIDs[i] = evt.Guilds[i].ID
	}

	_, _ = client.Emit(RequestGuildMembers, &RequestGuildMembersPayload{
		GuildIDs: guildIDs,
	})
}

//////////////////////////////////////////////////////
//
// Socket utilities
//
//////////////////////////////////////////////////////

func (c *Client) Event() SocketHandlerRegistrator {
	return &socketHandlerRegister{
		register: c,
	}
}

// Ready triggers a given callback when all shards has gotten their first Ready event
// Warning: Do not call Client.Connect before this.
func (c *Client) Ready(cb func()) {
	ctrl := &rdyCtrl{
		cb: cb,
	}

	c.On(EvtReady, func(_ Session, evt *Ready) {
		ctrl.Lock()
		defer ctrl.Unlock()

		l := c.shardManager.ShardCount()
		if l != uint(len(ctrl.shardReady)) {
			ctrl.shardReady = make([]bool, l)
			ctrl.localShardIDs = c.shardManager.ShardIDs()
		}

		ctrl.shardReady[evt.ShardID] = true
	}, ctrl)
}

// GuildsReady is triggered once all unavailable Guilds given in the READY event has loaded from their respective GUILD_CREATE events.
func (c *Client) GuildsReady(cb func()) {
	ctrl := &guildsRdyCtrl{
		status: make(map[Snowflake]bool),
	}
	ctrl.cb = cb
	ctrl.status[0] = false

	c.On(EvtReady, func(_ Session, evt *Ready) {
		ctrl.Lock()
		defer ctrl.Unlock()

		for _, g := range evt.Guilds {
			if _, ok := ctrl.status[g.ID]; !ok {
				ctrl.status[g.ID] = false
			}
		}

		delete(ctrl.status, 0)
	}, ctrl)

	c.On(EvtGuildCreate, func(_ Session, evt *GuildCreate) {
		ctrl.Lock()
		defer ctrl.Unlock()
		ctrl.status[evt.Guild.ID] = true
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

	if err := c.dispatcher.register(event, inputs...); err != nil {
		panic(err)
	}
}

// Emit sends a socket command directly to Discord.
func (c *Client) Emit(name gatewayCmdName, payload gatewayCmdPayload) (unchandledGuildIDs []Snowflake, err error) {
	c.RLock()
	defer c.RUnlock()
	if c.shardManager == nil {
		return nil, errors.New("you must connect before you can Emit")
	}

	p, err := prepareGatewayCommand(payload)
	if err != nil {
		return nil, err
	}
	return c.shardManager.Emit(string(name), p)
}

//////////////////////////////////////////////////////
//
// REST Methods
// customs
//
//////////////////////////////////////////////////////

/* status updates */

// UpdateStatus updates the Client's game status
// note: for simple games, check out UpdateStatusString
func (c *Client) UpdateStatus(s *UpdateStatusPayload) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, err := c.Emit(UpdateStatus, s)
	return err
}

// UpdateStatusString sets the Client's game activity to the provided string, status to online
// and type to Playing
func (c *Client) UpdateStatusString(s string) error {
	updateData := &UpdateStatusPayload{
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
