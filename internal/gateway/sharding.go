package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/andersfylling/disgord/internal/constant"
	"github.com/andersfylling/disgord/internal/event"
	"github.com/andersfylling/disgord/internal/gateway/cmd"
	"github.com/andersfylling/disgord/internal/logger"
)

const defaultShardRateLimit time.Duration = 5*time.Second + 100*time.Millisecond
const discordErrShardScalingRequired = 4011

type shardID = uint

type CmdPayload interface {
	isCmdPayload() bool
}

// GetShardForGuildID converts a GuildID into a ShardID for correct retrieval of guild information
func GetShardForGuildID(guildID Snowflake, shardCount uint) (shardID uint) {
	return uint(guildID>>22) % shardCount
}

func ConfigureShardConfig(ctx context.Context, client GatewayBotGetter, conf *ShardConfig) error {
	if len(conf.ShardIDs) == 0 && conf.ShardCount != 0 {
		return errors.New("ShardCount should only be set when you use distributed bots and have set the ShardIDs field - ShardCount is an optional field")
	}

	data, err := client.GetGatewayBot(ctx)
	if err != nil {
		return err
	}

	if len(conf.ShardIDs) > 0 || conf.ShardCount > 0 {
		conf.DisableAutoScaling = true
	}

	if conf.URL == "" {
		conf.URL = data.URL
	}

	if conf.IdentifiesPer24H == 0 {
		conf.IdentifiesPer24H = DefaultIdentifyRateLimit
	}

	if len(conf.ShardIDs) == 0 {
		conf.ShardCount = data.Shards
		for i := uint(0); i < data.Shards; i++ {
			conf.ShardIDs = append(conf.ShardIDs, i)
		}
	} else if conf.ShardCount == 0 {
		conf.ShardCount = uint(len(conf.ShardIDs))
	}

	if conf.ShardRateLimit == 0 {
		conf.ShardRateLimit = defaultShardRateLimit
	}

	return nil
}

// enableGuildSubscriptions if both typing event and presence event are to be ignore, we can disable GuildSubscription
// https://discord.com/developers/docs/topics/gateway#guild-subscriptions
func enableGuildSubscriptions(ignore []string) (updatedIgnores []string, ok bool) {
	requires := []string{
		event.TypingStart, event.PresenceUpdate,
	}
	for i := range ignore {
		for j := range requires {
			if ignore[i] == requires[j] {
				// remove matched requirements
				requires = append(requires[:j], requires[j+1:]...)
				break
			}
		}
		if len(requires) == 0 {
			break
		}
	}
	ok = len(requires) > 0
	// TODO: remove unnecessary events from the ignore slice

	return ignore, ok
}

func NewShardMngr(conf ShardManagerConfig) *shardMngr {
	conf.IgnoreEvents, conf.GuildSubscriptions = enableGuildSubscriptions(conf.IgnoreEvents)

	mngr := &shardMngr{
		conf:   conf,
		shards: map[shardID]*EvtClient{},
		DiscordPktPool: &sync.Pool{
			New: func() interface{} {
				return &DiscordPacket{}
			},
		},
	}
	if conf.ConnectQueue == nil {
		mngr.sync = newShardSync(&conf.ShardConfig, conf.Logger, "[shardSync]", conf.ShutdownChan)
		mngr.connectQueue = mngr.sync.queueShard

		go mngr.sync.process() // handle requests
	} else {
		mngr.connectQueue = conf.ConnectQueue
	}

	return mngr
}

// ShardManager regards websocket shards.
type ShardManager interface {
	Connect() error
	Disconnect() error
	Emit(string, CmdPayload) (unhandledGuildIDs []Snowflake, err error)
	LocalShardCount() uint
	ShardCount() uint
	ShardIDs() (shardIDs []uint)
	GetShard(shardID shardID) (shard *EvtClient, err error)
	HeartbeatLatencies() (latencies map[shardID]time.Duration, err error)
}

type ShardConfig struct {
	// Specify the shard ids that can be used by this instance.
	//  eg. ShardIds = []uint{0,1,2,3,11,12,13,14,32}
	//
	// This control is only useful if you have more than once instance of your bot duo to
	// high traffic or whatever reason you might possess.
	//
	// This also allows you to manually specify the number of shards, you just have to
	// specify their ID as well. You start from 0 until the number of shards you desire.
	//
	// Default value is populated by discord if this slice is nil.
	ShardIDs []uint

	// ShardCount should reflect the "total number of shards" across all
	// instances for your bot. If you run 3 containers with 2 shards each, then
	// the ShardCount should be 6, while the length of shardIDs would be
	// two on each container.
	//
	// defaults to len(shardIDs) if 0
	ShardCount uint

	// Large bots only. If Discord did not give you a custom rate limit, do not touch this.
	ShardRateLimit time.Duration

	// ConnectQueue is used to control how often shards can connect by sending an identify command.
	// For distributed systems, this must be overwritten as, by default, you can only send one identify
	// every five seconds. The default implementation can be found in shard_sync.go.
	ConnectQueue connectQueue

	// DisableAutoScaling is triggered when at least one shard gets a 4011 websocket
	// error from Discord. This causes all the shards to disconnect and new ones are created.
	//
	// default value is false unless shardIDs or ShardCount is set.
	DisableAutoScaling bool

	// OnScalingRequired is triggered when Discord closes the websocket connection
	// with a 4011 websocket error. It may run multiple times per session. You should
	// immediately call disconnect and scale your shards, unless you know what you're doing.
	//
	// This is triggered when DisableAutoScaling is true. If DisableAutoScaling is true and
	// OnScalingRequired is nil, this is considered an user error and will panic.
	//
	// You must return the new number of total shards and additional shard ids this instance
	// should setup. If you do not want this instance to gain extra shards, set AdditionalShardIDs
	// to nil.
	OnScalingRequired func(shardIDs []uint) (TotalNrOfShards uint, AdditionalShardIDs []uint)

	// OnScalingDiscardedRequests When scaling is triggered, some of the guilds might have moved to other shards
	// that do not exist on this disgord instance. This callback will return a list of guild ID that exists in
	// outgoing requests that were discarded due to no local shard match.
	//
	// Note: only regards systems with multiple disgord instances
	// TODO: return a list of outgoing requests instead such that people can re-trigger these on other instances.
	OnScalingDiscardedRequests func(unhandledGuildIDs []Snowflake)

	// IdentifiesPer24H regards how many identify packets a bot can send per a 24h period. Normally this
	// is 1000, but in some cases discord might allow you to increase it.
	//
	// Setting it to 0 will default it to 1000.
	IdentifiesPer24H uint

	// URL is fetched from the gateway before initialising a connection
	URL string
}

// ShardManagerConfig all fields, except proxy.Dialer, is required
type ShardManagerConfig struct {
	ShardConfig
	DisgordInfo  string
	BotToken     string
	HTTPClient   *http.Client
	Logger       logger.Logger
	ShutdownChan chan interface{}
	conn         Conn

	// ...
	IgnoreEvents []string
	Intents      Intent

	// sync ---
	EventChan chan<- *Event

	RESTClient GatewayBotGetter

	// user specific
	DefaultBotPresence *UpdateStatusPayload
	ProjectName        string
	GuildSubscriptions bool
}

type shardMngr struct {
	mu             sync.RWMutex
	conf           ShardManagerConfig
	shards         map[shardID]*EvtClient
	DiscordPktPool *sync.Pool

	sync         *shardSync
	connectQueue connectQueue
}

var _ ShardManager = (*shardMngr)(nil)

func (s *shardMngr) initShards() error {
	baseConfig := EvtConfig{ // TODO: not nicely grouped, feel free to adjust
		// identity
		Browser:             s.conf.DisgordInfo,
		Device:              s.conf.ProjectName,
		GuildLargeThreshold: 0, // let's not sometimes load partial guilds info. Either load everything or nothing.
		ShardCount:          s.conf.ShardCount,
		Presence:            s.conf.DefaultBotPresence,
		GuildSubscriptions:  s.conf.GuildSubscriptions,

		// lib specific
		Version:        constant.DiscordVersion,
		Encoding:       constant.JSONEncoding,
		Endpoint:       s.conf.URL,
		Logger:         s.conf.Logger,
		IgnoreEvents:   s.conf.IgnoreEvents,
		Intents:        s.conf.Intents,
		DiscordPktPool: s.DiscordPktPool,

		// synchronization
		EventChan:    s.conf.EventChan,
		connectQueue: s.connectQueue,

		// user settings
		BotToken:   s.conf.BotToken,
		HTTPClient: s.conf.HTTPClient,

		// other
		SystemShutdown: s.conf.ShutdownChan,
		discordErrListener: func(code int, reason string) {
			if code != discordErrShardScalingRequired {
				return
			}
			s.mu.Lock()
			defer s.mu.Unlock()
			s.conf.Logger.Info("scaling")

			if !s.conf.DisableAutoScaling {
				s.scale(code, reason)
			} else {
				if s.conf.OnScalingRequired == nil {
					panic("ShardConfig.OnScalingRequired must be set")
				}
				var newShards []uint
				s.conf.ShardCount, newShards = s.conf.OnScalingRequired(s.ShardIDs())
				s.conf.ShardIDs = append(s.conf.ShardIDs, newShards...)

				_ = s.Disconnect()
				if err := s.initShards(); err != nil {
					s.conf.Logger.Error("scaling", "init-shards", err)
					return
				}
				s.conf.Logger.Info("scaling", "connecting shards")
				if err := s.Connect(); err != nil {
					s.conf.Logger.Error("scaling", "connect", err)
				}
				s.conf.Logger.Info("scaling", "connected")
			}
		},
		conn: s.conf.conn,
	}

	for _, id := range s.conf.ShardIDs {
		if shard, alreadyConfigured := s.shards[id]; alreadyConfigured {
			shard.evtConf.ShardCount = s.conf.ShardCount
			continue
		}

		uniqueConfig := baseConfig // create copy, review requirement
		shard, err := NewEventClient(id, &uniqueConfig)
		if err != nil {
			return err
		}

		s.shards[id] = shard
	}
	return nil
}

func (s *shardMngr) Connect() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.conf.ShardIDs) == 0 {
		return errors.New("no shard ids has been registered")
	}

	if len(s.shards) == 0 {
		if err = s.initShards(); err != nil {
			return err
		}
	}

	for _, shard := range s.shards {
		err := shard.reconnectLoop()
		if err != nil {
			s.conf.Logger.Error(err)
		}
	}
	return nil
}
func (s *shardMngr) Disconnect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, shard := range s.shards {
		err := shard.Disconnect()
		if err != nil {
			s.conf.Logger.Error("Disconnect error (trivial):", err)
		}
		// possible connect/disconnect race..
		shard.sessionID = ""
		shard.sequenceNumber.Store(0)

		shard.haveConnectedOnce.Store(false)
	}
	return nil
}

func (s *shardMngr) LocalShardCount() uint {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// ShardIDs will always reflect the number of shards for this instance
	return uint(len(s.conf.ShardIDs))
}

func (s *shardMngr) ShardCount() uint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.conf.ShardCount
}

func (s *shardMngr) ShardIDs() (shardIDs []uint) {
	for id := range s.shards {
		shardIDs = append(shardIDs, id)
	}
	return shardIDs
}

// Emit splits up and dispatches the payload into the correct shards
// returns the guild ids it can not support and a error message
func (s *shardMngr) Emit(cmd string, payload CmdPayload) (guildIDs []Snowflake, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.shards) == 0 {
		return nil, errors.New("can not use Emit before Connected")
	}

	switch t := payload.(type) {
	case *RequestGuildMembersPayload:
		if len(s.shards) == 1 {
			for _, shard := range s.shards {
				return t.GuildIDs, shard.Emit(cmd, payload)
			}
		}

		requests := make(map[uint][]Snowflake)
		for i := range t.GuildIDs {
			shardID := GetShardForGuildID(t.GuildIDs[i], s.ShardCount())
			requests[shardID] = append(requests[shardID], t.GuildIDs[i])
		}

		for shardID := range requests {
			r := *t
			r.GuildIDs = requests[shardID]
			if shard, ok := s.shards[shardID]; ok {
				err = shard.Emit(cmd, &r)
				if err != nil {
					guildIDs = append(guildIDs, r.GuildIDs...)
				}
			} else {
				guildIDs = append(guildIDs, r.GuildIDs...)
			}
		}
	case *UpdateVoiceStatePayload:
		shardID := GetShardForGuildID(t.GuildID, s.ShardCount())
		if shard, ok := s.shards[shardID]; ok {
			err = shard.Emit(cmd, payload)
		} else {
			guildIDs = append(guildIDs, t.GuildID)
			err = errors.New("this guild is not handled by this shard")
		}
	case *UpdateStatusPayload:
		for _, shard := range s.shards {
			err = shard.Emit(cmd, payload)
		}
	default:
		err = errors.New("missing support for payload type")
	}

	return guildIDs, err
}

func (s *shardMngr) GetShard(shardID shardID) (shard *EvtClient, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if shard, ok := s.shards[shardID]; ok {
		return shard, nil
	}

	return nil, errors.New("no shard with given id " + fmt.Sprint(shardID))
}

func (s *shardMngr) HeartbeatLatencies() (latencies map[shardID]time.Duration, err error) {
	latencies = make(map[shardID]time.Duration)
	for id := range s.shards {
		latencies[id], err = s.shards[id].HeartbeatLatency()
		if err != nil {
			break
		}
	}
	return
}

func (s *shardMngr) scale(code int, reason string) {
	if s.conf.DisableAutoScaling {
		s.conf.Logger.Debug("discord require websocket shards to scale up but auto scaling is disabled - did not handle scaling internally")
		return
	}

	s.conf.Logger.Error("discord require websocket shards to scale up - starting auto scaling:", reason)

	unchandledGuilds := s.redistributeMsgs(func() {
		data, err := s.conf.RESTClient.GetGatewayBot(context.Background())
		if err != nil {
			s.conf.Logger.Error("autoscaling", err)
			return
		}

		_ = s.Disconnect()

		s.conf.URL = data.URL
		for i := uint(len(s.conf.ShardIDs) - 1); i < data.Shards; i++ {
			s.conf.ShardIDs = append(s.conf.ShardIDs, i)
			s.conf.ShardCount++
		}
		if err := s.initShards(); err != nil {
			s.conf.Logger.Error("autoscaling", "init-shards", err)
			return
		}
		if err := s.Connect(); err != nil {
			s.conf.Logger.Error("autoscaling", "connect", err)
		}
	})

	if s.conf.OnScalingDiscardedRequests != nil {
		s.conf.OnScalingDiscardedRequests(unchandledGuilds)
	}
}

func (s *shardMngr) redistributeMsgs(scaleShards func()) (unhandledGuildIDs []Snowflake) {
	var messages []*clientPacket
	for _, shard := range s.shards {
		messages = append(messages, shard.messageQueue.Steal()...)
	}

	scaleShards()

	// merge similar requests that only differs by guild IDs
	opToMerge := CmdNameToOpCode(cmd.RequestGuildMembers, clientTypeEvent)
	for i := range messages {
		m1 := messages[i]
		if m1 == nil || m1.Op != opToMerge {
			continue
		}

		for j := i + 1; j < len(messages); j++ {
			m2 := messages[j]
			if m2 == nil || m1.Op != m2.Op {
				continue
			}

			var rgm1 *RequestGuildMembersPayload
			var rgm2 *RequestGuildMembersPayload
			var ok bool
			if rgm1, ok = m1.Data.(*RequestGuildMembersPayload); !ok {
				continue
			}
			if rgm2, ok = m2.Data.(*RequestGuildMembersPayload); !ok {
				continue
			}
			rgm1.GuildIDs = append(rgm1.GuildIDs, rgm2.GuildIDs...)
			messages[j] = nil
		}
	}

	// reverse such that injected order stays the same
	// and merge similar requests that only differs by guild IDs
	for i := len(messages) - 1; i >= 0; i-- {
		m := messages[i]
		if m == nil {
			continue
		}

		if payload, ok := m.Data.(CmdPayload); ok {
			gIDs, _ := s.Emit(m.CmdName, payload)
			unhandledGuildIDs = append(unhandledGuildIDs, gIDs...)
			messages[i] = nil
		}
	}

	return unhandledGuildIDs
}
