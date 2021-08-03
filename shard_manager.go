package disgord

import (
	"fmt"
	"github.com/andersfylling/discordgateway"
	"github.com/andersfylling/discordgateway/opcode"
	"github.com/andersfylling/disgord/json"
)

type ShardManager interface {
	Start(Session) error
	Stop()
	SendCMD(GatewayCommand) error
}

type GatewayCommand interface {
	GuildID() Snowflake
	OperationCode() uint
}

type ErrorShardManager struct {
	Errors []error
}

func (e *ErrorShardManager) ErrorCount() (counter int) {
	for i := range e.Errors {
		if e.Errors[i] != nil {
			counter++
		}
	}
	return counter
}

func (e *ErrorShardManager) Error() string {
	return fmt.Sprintf("%d shard interactions failed: %+v", e.ErrorCount(), e.Errors)
}

var _ error = &ErrorShardManager{}

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

func (s *SimpleShardManager) SendCMD(cmd GatewayCommand) error {
	if cmd.GuildID().IsZero() {
		return s.sendCMDToAllShards(cmd)
	}
	return s.sendCMDToShard(cmd)
}

func (s *SimpleShardManager) sendCMDToAllShards(cmd GatewayCommand) error {
	code := opcode.Type(cmd.OperationCode())
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	shardErr := &ErrorShardManager{
		Errors: make([]error, len(s.shards)),
	}
	for i := range s.shards {
		shardErr.Errors[i] = s.shards[i].Write(code, message)
	}
	return shardErr
}

func (s *SimpleShardManager) sendCMDToShard(cmd GatewayCommand) error {
	code := opcode.Type(cmd.OperationCode())
	guildID := cmd.GuildID()
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	shardID := discordgateway.DeriveShardID(uint64(guildID), uint(len(s.shards)))
	return s.shards[shardID].Write(code, message)
}
