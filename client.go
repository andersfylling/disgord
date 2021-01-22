package disgord

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/proxy"

	"github.com/andersfylling/disgord/internal/gateway"
	"github.com/andersfylling/disgord/internal/logger"

	"github.com/andersfylling/disgord/internal/constant"

	"github.com/andersfylling/disgord/internal/httd"
)

// New create a Client. But panics on configuration/setup errors.
func New(conf Config) *Client {
	client, err := NewClient(context.Background(), conf)
	if err != nil {
		panic(err)
	}
	return client
}

// NewClient creates a new Disgord Client and returns an error on configuration issues
// context is required since a single external request is made to verify bot details
func NewClient(ctx context.Context, conf Config) (*Client, error) {
	return createClient(ctx, &conf)
}

func verifyClientProduction(ctx context.Context, client *Client) (Snowflake, error) {
	usr, err := client.CurrentUser().WithContext(ctx).Get(IgnoreCache)
	if err != nil {
		return 0, err
	}
	if usr == nil {
		return 0, fmt.Errorf("unable to gather bot information")
	}
	if usr.ID.IsZero() {
		return 0, fmt.Errorf("for some reason the bot ID is unknown")
	}

	return usr.ID, nil
}

var verifyClient func(ctx context.Context, client *Client) (Snowflake, error) = verifyClientProduction

// NewClient creates a new Disgord Client and returns an error on configuration issues
func createClient(ctx context.Context, conf *Config) (c *Client, err error) {
	if conf.Logger == nil {
		conf.Logger = logger.Empty{}
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

	if conf.Intents > 0 {
		conf.DMIntents |= conf.Intents
	}

	const DMIntents = IntentDirectMessageReactions | IntentDirectMessages | IntentDirectMessageTyping
	if validRange := conf.DMIntents & DMIntents; (conf.DMIntents ^ validRange) > 0 {
		return nil, errors.New("you have specified intents that are not for DM usage. See documentation")
	}

	if conf.IgnoreEvents != nil {
		conf.Logger.Info("Config.IgnoreEvents has been deprecated. Use Config.RejectEvents instead")
	}
	conf.RejectEvents = append(conf.RejectEvents, conf.IgnoreEvents...)

	// remove extra/duplicates events
	uniqueEventNames := make(map[string]bool)
	for _, eventName := range conf.RejectEvents {
		uniqueEventNames[eventName] = false
	}
	// if _, ok := uniqueEventNames[EvtUserUpdate]; ok {
	// 	return nil, errors.New("you can not reject the event USER_UPDATE")
	// }
	if _, ok := uniqueEventNames["PRESENCES_REPLACE"]; !ok {
		// https://github.com/discord/discord-api-docs/issues/683
		uniqueEventNames["PRESENCES_REPLACE"] = false
	}
	if _, ok := uniqueEventNames[EvtReady]; ok && conf.LoadMembersQuietly {
		return nil, fmt.Errorf("you can not reject the READY event when LoadMembersQuietly is set to true")
	}
	conf.RejectEvents = make([]string, 0, len(uniqueEventNames))
	for eventName, _ := range uniqueEventNames {
		conf.RejectEvents = append(conf.RejectEvents, eventName)
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

	// caching
	var cache Cache
	if conf.DisableCache {
		if _, ok := conf.Cache.(*CacheNop); !ok {
			cache = &CacheNop{}
		} else {
			cache = conf.Cache
		}
	} else if conf.Cache == nil {
		// don't specify any limits, this should be done by the user instead
		cache = NewCacheLFUImmutable(0, 0, 0, 0)
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

	// this external requests ensures two things:
	//  - the bot token is valid (a disgord instance is locked to a bot token)
	//  - that the bot id is always known
	if c.botID, err = verifyClient(ctx, c); err != nil {
		return nil, err
	}

	// TODO: this is just waiting to fail
	if internalCache, ok := c.cache.(*CacheLFUImmutable); ok {
		internalCache.currentUserID = c.botID
	}

	return c, nil
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

	// Deprecated: use DMIntents (values here are copied to DMIntents for now)
	// For direct communication with you bot you must specify intents
	Intents Intent

	// DMIntents specify intents related to direct message capabilities. Guild related intents are derived
	// from the RejectEvents config option (I hope that one day Intents can be removed all together, and
	// such optimizations can be handled in the background).
	//
	// You can see sent intents by enabling debug logging. Remember that derived from RejectIntents are appended.
	DMIntents Intent

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
	// Deprecated: use RejectEvents instead (nothing changed, just better naming)
	IgnoreEvents []string

	RejectEvents []string
}

// Client is the main disgord Client to hold your state and data. You must always initiate it using the constructor
// methods (eg. New(..) or NewClient(..)).
//
// Note that this Client holds all the REST methods, and is split across files, into whatever category
// the REST methods regards.
type Client struct {
	mu sync.RWMutex

	// current bot id
	botID Snowflake

	clientQueryBuilder

	shutdownChan chan interface{}
	config       *Config
	botToken     string

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
		c.Gateway().Ready(c.handlers.loadMembers)
	}
	c.Gateway().GuildCreate(c.handlers.saveGuildID)
	c.Gateway().GuildDelete(c.handlers.deleteGuildID)

	// start demultiplexer which also trigger dispatching
	go c.demultiplexer(c.dispatcher, c.eventChan)
}

type helperGatewayBotGetter struct {
	c *Client
}

var _ gateway.GatewayBotGetter = (*helperGatewayBotGetter)(nil)

func (h helperGatewayBotGetter) GetGatewayBot(ctx context.Context) (gateway *gateway.GatewayBot, err error) {
	return h.c.Gateway().WithContext(ctx).GetBot()
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

	_, _ = client.Gateway().Dispatch(RequestGuildMembers, &RequestGuildMembersPayload{
		GuildIDs: guildIDs,
	})
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
	_, err := c.Gateway().Dispatch(UpdateStatus, s)
	return err
}

// UpdateStatusString sets the Client's game activity to the provided string, status to online
// and type to Playing
func (c *Client) UpdateStatusString(s string) error {
	updateData := &UpdateStatusPayload{
		Since: nil,
		Game: []*Activity{
			{
				Name: s,
				Type: ActivityTypeGame,
			},
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
