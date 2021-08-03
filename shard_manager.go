package disgord

import (
	"github.com/andersfylling/discordgateway"
)

type ShardManager interface {
	Start(Session) error
	Stop()
	WriteMessage(message interface{}) error
}

var DefaultShardManager ShardManager = &SimpleShardManager{}

// SimpleShardManager is for single instance use only.
// Features:
//  - autoscaling
//    When discord closes the connection due to too few shards, the recommended shard number is requested from Discord.
//  - seamless restarts
//    Boots up a new shard manager before the current one is closed. During this time, events must be hashed and checked
//    to avoid duplicate events in the reactor.
type SimpleShardManager struct {
	session Session
	shards  []discordgateway.Shard
}

var _ ShardManager = &SimpleShardManager{}

func (s *SimpleShardManager) Start(session Session) error {
	s.session = session

	gatewayInfo, err := session.Gateway().GetBot()
	if err != nil {
		return err
	}

	s.shards = make([]discordgateway.Shard, gatewayInfo.Shards)

	panic("implement me")
}

func (s *SimpleShardManager) Stop() {
	for i := range s.shards {
		s.shards[i].Close()
	}
	s.shards = nil
	panic("implement me")
}

func (s *SimpleShardManager) WriteMessage(message interface{}) error {
	panic("implement me")
}
