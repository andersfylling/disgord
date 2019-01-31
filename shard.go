package disgord

import (
	"errors"
	"sync"
	"time"

	"github.com/andersfylling/disgord/httd"

	"github.com/andersfylling/disgord/constant"

	"github.com/andersfylling/disgord/websocket"
	"github.com/andersfylling/snowflake/v3"
)

const DefaultShardRateLimit float64 = 5.5 // seconds

type WSShardManagerConfig struct {
	// FirstID and ShardLimit creates the shard id range for this client.
	// this can be useful if you have multiple clients and don't want to
	// duplicate the sharded connections. But have unique ones on each machine.
	//
	// ShardLimit overrides the recommended shards sent by Discord if specified.
	// If you do not understand sharding, and your bot is not considered "large" according
	// to the documentation, then just don't touch these and let DisGord configure them.
	FirstID    uint
	ShardLimit uint

	ShardRateLimit float64

	// URL is fetched from the gateway before initialising a connection
	URL string
}

func NewShardManager(conf *Config) *WSShardManager {
	if conf == nil || conf.WSShardManagerConfig == nil {
		panic("missing shard config")
	}

	return &WSShardManager{
		conf:              conf.WSShardManagerConfig,
		log:               conf.Logger,
		identifyRatelimit: 5,
		shutdownChan:      conf.shutdownChan,
		TrackEvent:        &websocket.UniqueStringSlice{},
	}
}

type WSShardManager struct {
	sync.RWMutex

	evtChan chan *websocket.Event

	shards     []*WSShard
	conf       *WSShardManagerConfig
	TrackEvent *websocket.UniqueStringSlice

	identifyRatelimit float64 // seconds
	previousIdentify  time.Time
	idMutex           sync.RWMutex

	conRequestChan websocket.A
	shutdownChan   <-chan interface{}

	prepared bool
	log      Logger
}

func (s *WSShardManager) GetConnectionDetails(c httd.Getter) (url string, shardCount uint, err error) {
	var d *GatewayBot
	d, err = GetGatewayBot(c)
	if err != nil {
		return
	}

	url = d.URL
	shardCount = d.Shards
	return
}

func (s *WSShardManager) Prepare(conf *Config) error {
	s.Lock()
	defer s.Unlock()
	if s.prepared {
		return errors.New("already prepared")
	}
	s.prepared = true

	s.evtChan = make(chan *websocket.Event, 1+1+conf.WSShardManagerConfig.ShardLimit*2)
	s.shards = make([]*WSShard, conf.WSShardManagerConfig.ShardLimit)

	s.conRequestChan = make(websocket.A, conf.WSShardManagerConfig.ShardLimit+1)

	// handle shards asking for permission to connect (rate limiting)
	go func(a websocket.A) {
		for {
			select {
			case <-s.shutdownChan:
				return
			case b, ok := <-a:
				if !ok {
					s.log.Error("b is closed")
					continue
				}

				releaser := make(websocket.B)
				b <- &websocket.K{
					Release: releaser,
					Key:     412, // random
					// TODO: store shard info for better error handling and potential metrics
				}
				select {
				case <-releaser:
					// apply rate limit
					<-time.After(time.Duration(s.conf.ShardRateLimit) * time.Second)
				case <-s.shutdownChan:
					return
				}
			}
		}
	}(s.conRequestChan)

	var err error
	for i := range s.shards {
		s.shards[i] = &WSShard{}
		err = s.shards[i].Prepare(conf, s.evtChan, s.conRequestChan, s.TrackEvent, conf.WSShardManagerConfig.FirstID+uint(i))
		if err != nil {
			break
		}
	}

	return err
}

func (s *WSShardManager) GetShard(guildID snowflake.ID) (*WSShard, error) {
	s.RLock()
	defer s.RUnlock()

	if len(s.shards) == 0 {
		return nil, errors.New("no shards exist")
	}

	id := (uint64(guildID) >> 22) % uint64(len(s.shards))
	return s.shards[id], nil
}

// GetAvgHeartbeatLatency can be 0 if no heartbeat has been measured yet
func (s *WSShardManager) GetAvgHeartbeatLatency() (latency time.Duration, err error) {
	s.RLock()
	defer s.RUnlock()

	var tmp time.Duration
	for i := range s.shards {
		tmp, err = s.shards[i].ws.HeartbeatLatency()
		if err != nil {
			break
		}

		latency += tmp
	}
	if latency > 0 && len(s.shards) > 0 {
		latency = latency / time.Duration(len(s.shards))
	}

	return
}

func (s *WSShardManager) Connect() (err error) {
	s.RLock()
	defer s.RUnlock()

	if len(s.shards) == 0 {
		return errors.New("no shards exists")
	}

	for i := 0; i < len(s.shards); i++ {
		err = s.shards[i].Connect()
		if err != nil {
			s.log.Error(err.Error())
			// break
		}
	}

	return
}

func (s *WSShardManager) Disconnect() (err error) {
	s.RLock()
	defer s.RUnlock()

	for i := 0; i < len(s.shards); i++ {
		err = s.shards[i].Disconnect()
	}

	return
}

func (s *WSShardManager) Emit(cmd SocketCommand, data interface{}) (err error) {
	s.RLock()
	defer s.RUnlock()
	for i := 0; i < len(s.shards); i++ {
		err = s.shards[i].ws.Emit(cmd, data)
	}

	return
}

// InitialReadyReceived checks if each shard has gotten at least one Ready event
func (s *WSShardManager) InitialReadyReceived() bool {
	s.RLock()
	defer s.RUnlock()

	for i := 0; i < len(s.shards); i++ {
		if s.shards[i].ws.ReadyCounter == 0 {
			return false
		}
	}

	return true
}

var _ Emitter = (*WSShardManager)(nil)
var _ Link = (*WSShardManager)(nil)

type WSShard struct {
	sync.RWMutex
	id    uint
	total uint

	ws     *websocket.Client
	guilds []snowflake.ID
}

func (s *WSShard) Prepare(conf *Config, evtChan chan *websocket.Event, conRequestChan websocket.A, trackEvents *websocket.UniqueStringSlice, id uint) (err error) {
	s.id = id
	s.total = conf.WSShardManagerConfig.ShardLimit

	s.ws, err = websocket.NewClient(&websocket.Config{
		// identity
		Browser:             LibraryInfo(),
		Device:              conf.ProjectName,
		GuildLargeThreshold: 250, // TODO: config
		ShardCount:          s.total,

		// lib specific
		Version:       constant.DiscordVersion,
		Encoding:      constant.JSONEncoding,
		ChannelBuffer: 3,
		Endpoint:      conf.WSShardManagerConfig.URL,
		EventChan:     evtChan,
		TrackedEvents: trackEvents,
		Logger:        conf.Logger,
		A:             conRequestChan,

		// user settings
		BotToken:   conf.BotToken,
		HTTPClient: conf.HTTPClient,
	}, s.id)
	if err != nil {
		return err
	}

	return nil
}

func (s *WSShard) Emit(cmd SocketCommand, data interface{}) (err error) {
	return s.ws.Emit(cmd, data)
}

func (s *WSShard) Connect() error {
	return s.ws.Connect()
}

func (s *WSShard) Disconnect() error {
	return s.ws.Disconnect()
}

var _ Emitter = (*WSShard)(nil)
var _ Link = (*WSShard)(nil)
