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

	// URL is fetched from the gateway before initialising a connection
	URL string
}

func NewShardManager(conf *WSShardManagerConfig) *WSShardManager {
	if conf == nil {
		conf = &WSShardManagerConfig{}
	}
	return &WSShardManager{
		conf:       conf,
		TrackEvent: &websocket.UniqueStringSlice{},
	}
}

type WSShardManager struct {
	sync.RWMutex

	evtChan chan *websocket.Event

	shards     []*WSShard
	conf       *WSShardManagerConfig
	TrackEvent *websocket.UniqueStringSlice

	prepared bool
}

func (s *WSShardManager) NrOfShards() uint {
	s.RLock()
	defer s.RUnlock()

	return uint(len(s.shards))
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

	s.evtChan = make(chan *websocket.Event, 1+conf.WSShardManagerConfig.ShardLimit*2)
	s.shards = make([]*WSShard, conf.WSShardManagerConfig.ShardLimit)

	var err error
	for i, e := range s.shards {
		e = &WSShard{}
		err = e.Prepare(conf, conf.WSShardManagerConfig.FirstID+uint(i))
		if err != nil {
			break
		}
	}

	return err
}

func (s *WSShardManager) GetShard(guildID snowflake.ID) *WSShard {
	s.RLock()
	defer s.RUnlock()

	if s.NrOfShards() == 0 {
		return nil
	}

	id := (uint64(guildID) >> 22) % uint64(len(s.shards))
	return s.shards[id]
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
	if latency > 0 {
		latency = latency / time.Duration(len(s.shards))
	}

	return
}

func (s *WSShardManager) ConnectAllShards() (err error) {
	s.RLock()
	defer s.RUnlock()

	if len(s.shards) == 0 {
		return errors.New("no shards exists")
	}

	if err = s.shards[0].Connect(); err != nil {
		return err
	}

	for i := 1; i < len(s.shards); i++ {
		// ratelimit: 1/5s
		<-time.After(5 * time.Second)
		err = s.shards[i].Connect()
		if err != nil {
			break
		}
	}

	return
}

func (s *WSShardManager) DisconnectAllShards() (err error) {
	s.RLock()
	defer s.RUnlock()

	for i := 0; i < len(s.shards); i++ {
		err = s.shards[i].Disconnect()
	}

	return
}

func (s *WSShardManager) EmitThroughAllShards(cmd SocketCommand, data interface{}) (err error) {
	s.RLock()
	defer s.RUnlock()
	for i := 0; i < len(s.shards); i++ {
		err = s.shards[i].ws.Emit(cmd, data)
	}

	return
}

type WSShard struct {
	sync.RWMutex
	id    uint
	total uint

	ws      *websocket.Client
	evtChan <-chan *websocket.Event
	guilds  []snowflake.ID
}

func (s *WSShard) Prepare(conf *Config, id uint) (err error) {
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

		// user settings
		BotToken:   conf.BotToken,
		HTTPClient: conf.HTTPClient,
	}, s.id)
	if err != nil {
		return err
	}

	return nil
}

func (s *WSShard) Connect() error {
	return s.ws.Connect()
}

func (s *WSShard) Disconnect() error {
	return s.ws.Disconnect()
}
