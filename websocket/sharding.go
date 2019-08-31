package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/logger"

	"golang.org/x/net/proxy"
)

const defaultShardRateLimit float64 = 5.1 // seconds
type shardID = uint

func ConfigureShardConfig(client GatewayBotGetter, conf *ShardConfig) error {
	data, err := client.GetGatewayBot()
	if err != nil {
		return err
	}

	if conf.URL == "" {
		conf.URL = data.URL
	}

	if len(conf.ShardIDs) == 0 {
		conf.TotalNumberOfShards = data.Shards
		for i := uint(0); i < data.Shards; i++ {
			conf.ShardIDs = append(conf.ShardIDs, i)
		}
	} else {
		conf.TotalNumberOfShards = uint(len(conf.ShardIDs))
	}

	if conf.ShardRateLimit == 0 {
		conf.ShardRateLimit = defaultShardRateLimit
	}

	return nil
}

func NewShardMngr(conf ShardManagerConfig) *shardMngr {
	mngr := &shardMngr{
		shards: map[shardID]*EvtClient{},
		conf:   conf,
		DiscordPktPool: &sync.Pool{
			New: func() interface{} {
				return &DiscordPacket{}
			},
		},
	}
	mngr.sync.logger = conf.Logger
	mngr.sync.timeoutMs = time.Duration(int64(conf.ShardRateLimit*1000)) * time.Millisecond

	return mngr
}

// ShardManager regards websocket shards.
type ShardManager interface {
	Connect() error
	Disconnect() error
	Emit(string, interface{}) error
	NrOfShards() uint
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

	// TotalNumberOfShards should reflect the "total number of shards" across all
	// instances for your bot. If you run 3 containers with 2 shards each, then
	// the TotalNumberOfShards should be 6, while the length of ShardIDs would be
	// two on each container.
	//
	// defaults to len(ShardIDs) if 0
	TotalNumberOfShards uint

	// Large bots only. If Discord did not give you a custom rate limit, do not touch this.
	ShardRateLimit float64

	// URL is fetched from the gateway before initialising a connection
	URL string
}

// ShardManagerConfig all fields, except proxy.Dialer, is required
type ShardManagerConfig struct {
	ShardConfig
	DisgordInfo  string
	BotToken     string
	Proxy        proxy.Dialer
	HTTPClient   *http.Client
	Logger       logger.Logger
	ShutdownChan chan interface{}

	// ...
	TrackedEvents *UniqueStringSlice

	// sync ---
	EventChan chan<- *Event

	// user specific
	DefaultBotPresence interface{}
	ProjectName        string
}

type shardMngr struct {
	mu             sync.RWMutex
	conf           ShardManagerConfig
	shards         map[shardID]*EvtClient
	DiscordPktPool *sync.Pool

	sync shardSync
}

var _ ShardManager = (*shardMngr)(nil)

func (s *shardMngr) initializeShards() error {
	baseConfig := EvtConfig{ // TIP: not nicely grouped, feel free to adjust
		// identity
		Browser:             s.conf.DisgordInfo,
		Device:              s.conf.ProjectName,
		GuildLargeThreshold: 0, // let's not sometimes load partial guilds info. Either load everything or nothing.
		ShardCount:          uint(len(s.conf.ShardIDs)),
		Presence:            s.conf.DefaultBotPresence,

		// lib specific
		Version:        constant.DiscordVersion,
		Encoding:       constant.JSONEncoding,
		Endpoint:       s.conf.URL,
		Logger:         s.conf.Logger,
		TrackedEvents:  s.conf.TrackedEvents,
		DiscordPktPool: s.DiscordPktPool,

		// synchronization
		EventChan:    s.conf.EventChan,
		connectQueue: s.sync.queueShard,

		// user settings
		BotToken:   s.conf.BotToken,
		Proxy:      s.conf.Proxy,
		HTTPClient: s.conf.HTTPClient,

		// other
		SystemShutdown: s.conf.ShutdownChan,
	}

	for _, id := range s.conf.ShardIDs {
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
		if err = s.initializeShards(); err != nil {
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
	}
	return nil
}

func (s *shardMngr) NrOfShards() uint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.conf.TotalNumberOfShards
}

func (s *shardMngr) Emit(cmd string, data interface{}) (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, shard := range s.shards {
		err = shard.Emit(cmd, data)
	}
	return err
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
