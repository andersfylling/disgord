package websocket

import (
	"time"
)

const defaultShardRateLimit float64 = 5.1 // seconds
type shardID = uint

func NewShardMngr(conf *ShardConfig) *shardMngr {
	return &shardMngr{
		conf: conf,
	}
}

// ShardManager regards websocket shards.
type ShardManager interface {
	Connect() error
	Disconnect() error
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

	// Large bots only. If Discord did not give you a custom rate limit, do not touch this.
	ShardRateLimit float64

	// URL is fetched from the gateway before initialising a connection
	URL string
}

type shardMngr struct {
	conf   *ShardConfig
	shards map[shardID]*EvtClient
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
