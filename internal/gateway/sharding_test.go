// +build !integration

package gateway

import (
	"context"
	"testing"
	"time"

	"github.com/andersfylling/disgord/internal/event"
	"github.com/andersfylling/disgord/internal/gateway/cmd"
	"github.com/andersfylling/disgord/internal/logger"
)

type GatewayBotGetterMock struct {
	get func() (gateway *GatewayBot, err error)
}

func (g GatewayBotGetterMock) GetGatewayBot(_ context.Context) (gateway *GatewayBot, err error) {
	return g.get()
}

var _ GatewayBotGetter = (*GatewayBotGetterMock)(nil)

func TestConfigureShardConfig(t *testing.T) {
	nrOfShards := uint(4)
	u := "localhost:6060"
	mock := &GatewayBotGetterMock{
		get: func() (gateway *GatewayBot, err error) {
			return &GatewayBot{
				Shards:  nrOfShards,
				Gateway: Gateway{u},
			}, nil
		},
	}

	conf := ShardConfig{}
	if err := ConfigureShardConfig(context.Background(), mock, &conf); err != nil {
		t.Error(err)
	}
	if conf.URL != u {
		t.Error("url was not set")
	}
	if len(conf.ShardIDs) != int(conf.ShardCount) && conf.ShardCount != nrOfShards {
		t.Error("incorrectly set number of shards")
	}
	if conf.DisableAutoScaling {
		t.Error("DisableAutoScaling should not be true")
	}

	conf = ShardConfig{
		ShardIDs: []uint{34, 7, 2},
	}
	if err := ConfigureShardConfig(context.Background(), mock, &conf); err != nil {
		t.Error(err)
	}
	if !conf.DisableAutoScaling {
		t.Error("DisableAutoScaling should be true")
	}

	conf = ShardConfig{
		ShardIDs:   []uint{34, 7, 2},
		ShardCount: 34,
	}
	if err := ConfigureShardConfig(context.Background(), mock, &conf); err != nil {
		t.Error(err)
	}
	if !conf.DisableAutoScaling {
		t.Error("DisableAutoScaling should be true")
	}
}

func TestEnableGuildSubscriptions(t *testing.T) {
	ignore := []string{
		event.TypingStart, event.PresenceUpdate,
	}
	if _, ok := enableGuildSubscriptions(ignore); ok {
		t.Error("guild sub should be disabled")
	}

	ignore = []string{
		event.TypingStart, event.PresenceUpdate, event.Ready,
	}
	if _, ok := enableGuildSubscriptions(ignore); ok {
		t.Error("guild sub should be disabled")
	}

	ignore = []string{
		event.TypingStart, event.Ready,
	}
	if _, ok := enableGuildSubscriptions(ignore); !ok {
		t.Error("guild sub should be enabled")
	}

	ignore = []string{}
	if _, ok := enableGuildSubscriptions(ignore); !ok {
		t.Error("guild sub should be enabled")
	}
}

func TestRedistributeShardMessages(t *testing.T) {
	u := "localhost:6060"
	mock := &GatewayBotGetterMock{
		get: func() (gateway *GatewayBot, err error) {
			return &GatewayBot{
				Shards:  4,
				Gateway: Gateway{u},
			}, nil
		},
	}
	config := ShardManagerConfig{
		BotToken:     "test",
		ShutdownChan: make(chan interface{}),
		EventChan:    make(chan *Event),
		Logger:       &logger.Empty{},
	}
	defer func() {
		close(config.ShutdownChan)
		close(config.EventChan)
	}()

	if err := ConfigureShardConfig(context.Background(), mock, &config.ShardConfig); err != nil {
		t.Fatal(err)
	}

	mngr := NewShardMngr(config)
	if err := mngr.initShards(); err != nil {
		t.Fatal(err)
	}

	// trick shards into thinking they have connected so we can emit msgs
	connect := func() {
		for _, shard := range mngr.shards {
			shard.haveConnectedOnce.Store(true)
		}
	}
	connect()

	for i := 1; i <= int(mngr.conf.ShardCount*14); i++ {
		p := &RequestGuildMembersPayload{GuildIDs: []Snowflake{Snowflake(i << 22)}}
		if unhandledGuilds, err := mngr.Emit(cmd.RequestGuildMembers, p); err != nil || len(unhandledGuilds) != 0 {
			t.Error(err)
			t.Fatalf("%+v", unhandledGuilds)
		}
	}

	verifyDistribution := func(sid string) {
		for id, shard := range mngr.shards {
			for i := range shard.messageQueue.messages {
				m := shard.messageQueue.messages[i]
				if g, ok := m.Data.(*RequestGuildMembersPayload); ok {
					if GetShardForGuildID(g.GuildIDs[0], mngr.conf.ShardCount) != id {
						t.Error(sid, "incorrect distribution")
					}
				} else {
					panic(sid + "not *RequestGuildMembersPayload")
				}
			}
			if len(shard.messageQueue.messages) == 0 {
				t.Error(sid, "there should be at least one message in shard", id)
			}
		}
	}

	verifyDistribution("1")
	mngr.redistributeMsgs(func() {})
	verifyDistribution("2")

	mngr.redistributeMsgs(func() {
		mngr.conf.ShardIDs = append(mngr.conf.ShardIDs, uint(len(mngr.conf.ShardIDs)))
		mngr.conf.ShardCount++
		if err := mngr.initShards(); err != nil {
			t.Fatal(err)
		}
		connect()
	})
	verifyDistribution("3")
}

func TestIdentifyRateLimiting(t *testing.T) {
	u := "localhost:6060"
	mock := &GatewayBotGetterMock{
		get: func() (gateway *GatewayBot, err error) {
			return &GatewayBot{
				Shards:  1,
				Gateway: Gateway{u},
			}, nil
		},
	}
	config := ShardManagerConfig{
		BotToken:     "test",
		ShutdownChan: make(chan interface{}),
		EventChan:    make(chan *Event),
		Logger:       &logger.Empty{},
	}
	defer func() {
		close(config.EventChan)
		close(config.ShutdownChan)
	}()

	if err := ConfigureShardConfig(context.Background(), mock, &config.ShardConfig); err != nil {
		t.Fatal(err)
	}

	mngr := NewShardMngr(config)
	if err := mngr.initShards(); err != nil {
		t.Fatal(err)
	}

	ts := time.Now().Add(20 * time.Hour)
	reconnects := make([]time.Time, 0, DefaultIdentifyRateLimit+1)
	for i := 1; i <= DefaultIdentifyRateLimit-2; i++ {
		reconnects = append(reconnects, ts)
	}

	mngr.sync.metric.Lock()
	mngr.sync.metric.Reconnects = reconnects
	mngr.sync.metric.Unlock()

	nrOfTimestamps := mngr.sync.metric.ReconnectsSince(24 * time.Hour)
	if nrOfTimestamps != DefaultIdentifyRateLimit-2 {
		t.Fatalf("should be 998 reconnect time stamps, got %d", nrOfTimestamps)
	}

	// the timeout is after a run execution, so we add a entry before the test case
	mngr.connectQueue(0, func() error {
		return nil
	})
	connected := make(chan interface{})
	go mngr.connectQueue(0, func() error {
		connected <- true
		return nil
	})

	select {
	case <-connected:
		t.Fatal("should not be able to connect")
	case <-time.After(100 * time.Millisecond): // TODO: remove timeout, just don't know how yet
		select {
		case item, ok := <-mngr.sync.queue:
			if !ok {
				t.Fatal("queue was closed somehow")
			}
			if item == nil {
				t.Fatal("expected item to not be nil")
			}
		default:
		}
	}
}
