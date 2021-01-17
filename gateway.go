package disgord

import (
	"context"
	"errors"
	"fmt"

	"github.com/andersfylling/disgord/internal/gateway"
	"github.com/andersfylling/disgord/internal/gateway/cmd"
	"github.com/andersfylling/disgord/internal/httd"
	"github.com/andersfylling/disgord/json"
)

func (c *Client) Gateway() GatewayQueryBuilder {
	return &gatewayQueryBuilder{client: c, ctx: context.Background(), socketHandlerRegister: socketHandlerRegister{reactor: c.dispatcher}}
}

type GatewayQueryBuilder interface {
	WithContext(ctx context.Context) GatewayQueryBuilder

	Get() (gateway *gateway.Gateway, err error)
	GetBot() (gateway *gateway.GatewayBot, err error)

	BotReady(func())
	BotGuildsReady(func())

	Dispatch(name gatewayCmdName, payload gateway.CmdPayload) (unchandledGuildIDs []Snowflake, err error)

	// Connect establishes a websocket connection to the discord API
	Connect() error
	StayConnectedUntilInterrupted() error

	// Disconnect closes the discord websocket connection
	Disconnect() error
	DisconnectOnInterrupt() error

	SocketHandlerRegistrator
}

type gatewayQueryBuilder struct {
	ctx    context.Context
	client *Client
	socketHandlerRegister
}

func (g gatewayQueryBuilder) WithContext(ctx context.Context) GatewayQueryBuilder {
	g.ctx = ctx
	return &g
}

// Connect establishes a websocket connection to the discord API
func (g gatewayQueryBuilder) Connect() (err error) {
	// set the user ID upon connection
	// only works for socketing
	//
	// also verifies that the correct credentials were supplied

	// Avoid races during connection setup
	g.client.mu.Lock()
	defer g.client.mu.Unlock()

	if err = gateway.ConfigureShardConfig(g.ctx, helperGatewayBotGetter{g.client}, &g.client.config.ShardConfig); err != nil {
		return err
	}

	shardMngrConf := gateway.ShardManagerConfig{
		ShardConfig:  g.client.config.ShardConfig,
		Logger:       g.client.config.Logger,
		ShutdownChan: g.client.config.shutdownChan,
		IgnoreEvents: g.client.config.RejectEvents,
		Intents:      g.client.config.DMIntents,
		EventChan:    g.client.eventChan,
		DisgordInfo:  LibraryInfo(),
		ProjectName:  g.client.config.ProjectName,
		BotToken:     g.client.config.BotToken,
	}

	if g.client.config.Presence != nil {
		if g.client.config.Presence.Status == "" {
			g.client.config.Presence.Status = StatusOnline // default
		}
		shardMngrConf.DefaultBotPresence = g.client.config.Presence
	}

	sharding := gateway.NewShardMngr(shardMngrConf)

	g.client.setupConnectEnv()

	g.client.log.Info("Connecting to discord Gateway")
	if err = sharding.Connect(); err != nil {
		g.client.log.Info(err)
		return err
	}

	g.client.log.Info("Connected")
	g.client.shardManager = sharding
	return nil
}

// Disconnect closes the discord websocket connection
func (g gatewayQueryBuilder) Disconnect() (err error) {
	fmt.Println() // to keep ^C on it's own line
	g.client.log.Info("Closing Discord gateway connection")
	close(g.client.dispatcher.shutdown)
	if err = g.client.shardManager.Disconnect(); err != nil {
		g.client.log.Error(err)
		return err
	}
	close(g.client.shutdownChan)
	g.client.log.Info("Disconnected")

	return nil
}

// DisconnectOnInterrupt wait until a termination signal is detected
func (g gatewayQueryBuilder) DisconnectOnInterrupt() (err error) {
	// catches panic when being called as a deferred function
	if r := recover(); r != nil {
		panic("unable to connect due to above error")
	}

	<-CreateTermSigListener()
	return g.Disconnect()
}

// StayConnectedUntilInterrupted is a simple wrapper for connect, and disconnect that listens for system interrupts.
// When a error happens you can terminate the application without worries.
func (g gatewayQueryBuilder) StayConnectedUntilInterrupted() (err error) {
	// catches panic when being called as a deferred function
	if r := recover(); r != nil {
		panic("unable to connect due to above error")
	}

	if err = g.Connect(); err != nil {
		g.client.log.Error(err)
		return err
	}

	ctx := g.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	select {
	case <-CreateTermSigListener():
	case <-ctx.Done():
	}

	return g.Disconnect()
}

// BotReady triggers a given callback when all shards has gotten their first Ready event
// Warning: Do not call Client.Connect before this.
func (g gatewayQueryBuilder) BotReady(cb func()) {
	ctrl := &rdyCtrl{
		cb: cb,
	}

	g.WithCtrl(ctrl).Ready(func(_ Session, evt *Ready) {
		ctrl.Lock()
		defer ctrl.Unlock()

		l := g.client.shardManager.ShardCount()
		if l != uint(len(ctrl.shardReady)) {
			ctrl.shardReady = make([]bool, l)
			ctrl.localShardIDs = g.client.shardManager.ShardIDs()
		}

		ctrl.shardReady[evt.ShardID] = true
	})
}

// BotGuildsReady is triggered once all unavailable Guilds given in the READY event has loaded from their respective GUILD_CREATE events.
func (g gatewayQueryBuilder) BotGuildsReady(cb func()) {
	ctrl := &guildsRdyCtrl{
		status: make(map[Snowflake]bool),
	}
	ctrl.cb = cb
	ctrl.status[0] = false

	g.WithCtrl(ctrl).Ready(func(_ Session, evt *Ready) {
		ctrl.Lock()
		defer ctrl.Unlock()

		for _, g := range evt.Guilds {
			if _, ok := ctrl.status[g.ID]; !ok {
				ctrl.status[g.ID] = false
			}
		}

		delete(ctrl.status, 0)
	})

	g.WithCtrl(ctrl).GuildCreate(func(_ Session, evt *GuildCreate) {
		ctrl.Lock()
		defer ctrl.Unlock()
		ctrl.status[evt.Guild.ID] = true
	})
}

// Emit sends a socket command directly to Discord.
func (g gatewayQueryBuilder) Dispatch(name gatewayCmdName, payload gateway.CmdPayload) (unchandledGuildIDs []Snowflake, err error) {
	g.client.mu.RLock()
	defer g.client.mu.RUnlock()
	if g.client.shardManager == nil {
		return nil, errors.New("you must connect before you can Dispatch requests")
	}

	return g.client.shardManager.Emit(string(name), payload)
}

// Get Returns an object with a single valid WSS URL, which the Client can use for Connecting.
// Clients should cacheLink this value and only call this endpoint to retrieve a new URL if they are unable to
// properly establish a connection using the cached version of the URL.
//  Method                  GET
//  Endpoint                /gateway
//  Discord documentation   https://discord.com/developers/docs/topics/gateway#get-gateway
//  Reviewed                2018-10-12
//  Comment                 This endpoint does not require authentication.
func (g gatewayQueryBuilder) Get() (gateway *gateway.Gateway, err error) {
	var body []byte
	_, body, err = g.client.req.Do(g.ctx, &httd.Request{
		Method:   httd.MethodGet,
		Endpoint: "/gateway",
	})
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &gateway)

	if gateway.URL, err = ensureDiscordGatewayURLHasQueryParams(gateway.URL); err != nil {
		return gateway, err
	}

	return
}

// GetBot Returns an object based on the information in Get Gateway, plus additional metadata
// that can help during the operation of large or sharded bots. Unlike the Get Gateway, this route should not
// be cached for extended periods of time as the value is not guaranteed to be the same per-call, and
// changes as the bot joins/leaves Guilds.
//  Method                  GET
//  Endpoint                /gateway/bot
//  Discord documentation   https://discord.com/developers/docs/topics/gateway#get-gateway-bot
//  Reviewed                2018-10-12
//  Comment                 This endpoint requires authentication using a valid bot token.
func (g gatewayQueryBuilder) GetBot() (gateway *gateway.GatewayBot, err error) {
	var body []byte
	_, body, err = g.client.req.Do(g.ctx, &httd.Request{
		Method:   httd.MethodGet,
		Endpoint: "/gateway/bot",
	})
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &gateway); err != nil {
		return nil, err
	}

	if gateway.URL, err = ensureDiscordGatewayURLHasQueryParams(gateway.URL); err != nil {
		return gateway, err
	}

	return gateway, nil
}

// ##############################################################################################################

// gatewayCmdName is the gateway command name for the payload to be sent to Discord over a websocket connection.
type gatewayCmdName string

const (
	// GatewayCmdRequestGuildMembers Used to request offline members for a guild or
	// a list of Guilds. When initially connecting, the gateway will only send
	// offline members if a guild has less than the large_threshold members
	// (value in the Gateway Identify). If a Client wishes to receive additional
	// members, they need to explicitly request them via this operation. The
	// server will send Guild Members Chunk events in response with up to 1000
	// members per chunk until all members that match the request have been sent.
	RequestGuildMembers gatewayCmdName = cmd.RequestGuildMembers

	// UpdateVoiceState Sent when a Client wants to join, move, or
	// disconnect from a voice channel.
	UpdateVoiceState gatewayCmdName = cmd.UpdateVoiceState

	// UpdateStatus Sent by the Client to indicate a presence or status
	// update.
	UpdateStatus gatewayCmdName = cmd.UpdateStatus
)

// #################################################################
// RequestGuildMembersPayload payload for socket command REQUEST_GUILD_MEMBERS.
// See RequestGuildMembers
//
// WARNING: If this request is in queue while a auto-scaling is forced, it will be removed from the queue
// and not re-inserted like the other commands. This is due to the guild id slice, which is a bit trickier
// to handle.
//
// Wrapper for websocket.RequestGuildMembersPayload
type RequestGuildMembersPayload = gateway.RequestGuildMembersPayload

var _ gateway.CmdPayload = (*RequestGuildMembersPayload)(nil)

// UpdateVoiceStatePayload payload for socket command UPDATE_VOICE_STATE.
// see UpdateVoiceState
//
// Wrapper for websocket.UpdateVoiceStatePayload
type UpdateVoiceStatePayload = gateway.UpdateVoiceStatePayload

var _ gateway.CmdPayload = (*UpdateVoiceStatePayload)(nil)

const (
	StatusOnline  = gateway.StatusOnline
	StatusOffline = gateway.StatusOffline
	StatusDnd     = gateway.StatusDND
	StatusIdle    = gateway.StatusIdle
)

// UpdateStatusPayload payload for socket command UPDATE_STATUS.
// see UpdateStatus
//
// Wrapper for websocket.UpdateStatusPayload
type UpdateStatusPayload = gateway.UpdateStatusPayload

var _ gateway.CmdPayload = (*UpdateStatusPayload)(nil)
