package websocket

import (
	"time"
)

const defaultShardRateLimit float64 = 5.1 // seconds

type shard = EvtClient
type shardID = uint

func NewShardMngr(conf *ShardConfig) *shardMngr {
	return &shardMngr{
		conf: conf,
	}
}

type ShardManager interface {
	Connect() error
	Disconnect() error
	HeatbeatLatencies() (latencies map[shardID]time.Duration, err error)
}

type ShardConfig struct {
	// FirstID and ShardLimit creates the shard id range for this Client.
	// this can be useful if you have multiple clients and don't want to
	// duplicate the sharded connections. But have unique ones on each machine.
	//
	// NrOfShards overrides the recommended shards sent by Discord if specified.
	// If you do not understand sharding, and your bot is not considered "large" according
	// to the documentation, then just don't touch these and let DisGord configure them.
	FirstID    uint
	NrOfShards uint

	// Large bots only. If Discord did not give you a custom rate limit, do not touch this.
	ShardRateLimit float64

	// URL is fetched from the gateway before initialising a connection
	URL string
}

type shardMngr struct {
	conf   *ShardConfig
	shards map[shardID]*shard
}

var _ ShardManager = (*shardMngr)(nil)

func (s *shardMngr) initializeShards() error {

	panic("implement me")
}

func (s *shardMngr) Connect() (err error) {
	if len(s.shards) == 0 {
		if err = s.initializeShards(); err != nil {
			return err
		}
	}

	panic("implement me")
}
func (s *shardMngr) Disconnect() error {
	panic("implement me")
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
