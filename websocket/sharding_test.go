package websocket

import (
	"testing"

	"github.com/andersfylling/disgord/websocket/cmd"

	"github.com/andersfylling/disgord/event"
)

type GatewayBotGetterMock struct {
	get func() (gateway *GatewayBot, err error)
}

func (g GatewayBotGetterMock) GetGatewayBot() (gateway *GatewayBot, err error) {
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
	if err := ConfigureShardConfig(mock, &conf); err != nil {
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
	if err := ConfigureShardConfig(mock, &conf); err != nil {
		t.Error(err)
	}
	if !conf.DisableAutoScaling {
		t.Error("DisableAutoScaling should be true")
	}

	conf = ShardConfig{
		ShardIDs:   []uint{34, 7, 2},
		ShardCount: 34,
	}
	if err := ConfigureShardConfig(mock, &conf); err != nil {
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
		ShutdownChan: make(chan interface{}),
		EventChan:    make(chan *Event),
	}
	defer func() {
		close(config.ShutdownChan)
		close(config.EventChan)
	}()

	if err := ConfigureShardConfig(mock, &config.ShardConfig); err != nil {
		t.Fatal(err)
	}

	mngr := NewShardMngr(config)
	if err := mngr.initShards(); err != nil {
		t.Fatal(err)
	}

	// trick shards into thinking they have connected so we can emit msgs
	connect := func() {
		for _, shard := range mngr.shards {
			shard.haveConnectedOnce = true
		}
	}
	connect()

	for i := 1; i <= int(mngr.conf.ShardCount*14); i++ {
		if err := mngr.Emit(cmd.UpdateVoiceState, true, Snowflake(i<<22)); err != nil {
			t.Fatal(err)
		}
	}

	verifyDistribution := func() {
		for id, shard := range mngr.shards {
			for i := range shard.messageQueue.messages {
				m := shard.messageQueue.messages[i]
				if GetShardForGuildID(m.guildID, mngr.conf.ShardCount) != id {
					t.Error("incorrect distribution")
				}
			}
			if len(shard.messageQueue.messages) == 0 {
				t.Error("there should be at least one message in shard", id)
			}
		}
	}

	verifyDistribution()
	mngr.redistributeMsgs(func() {})
	verifyDistribution()

	mngr.redistributeMsgs(func() {
		mngr.conf.ShardIDs = append(mngr.conf.ShardIDs, uint(len(mngr.conf.ShardIDs)))
		mngr.conf.ShardCount++
		if err := mngr.initShards(); err != nil {
			t.Fatal(err)
		}
		connect()
	})
	verifyDistribution()
}

//
//func TestShardAutoScalingFailsafe(t *testing.T) {
//	// when discord disconnects one or more shards with the websocket
//	// error 4011: require shard scaling
//
//	eChan := make(chan *Event)
//	shutdown := make(chan interface{})
//	done := make(chan interface{})
//	deadline := 1 * time.Second
//	nrOfShards := uint(4)
//	conn := &testWS{
//		closing:      make(chan interface{}),
//		opening:      make(chan interface{}),
//		writing:      make(chan interface{}),
//		reading:      make(chan []byte),
//		disconnected: true,
//	}
//
//	mngr := NewShardMngr(ShardManagerConfig{
//		ShardConfig: ShardConfig{
//			shardIDs: []uint{0, 1},
//		},
//		DisgordInfo:   "",
//		BotToken:      "",
//		Proxy:         nil,
//		HTTPClient:    nil,
//		Logger:        logger.DefaultLogger(true),
//		ShutdownChan:  shutdown,
//		conn:          conn,
//		TrackedEvents: nil,
//		EventChan:     eChan,
//		RESTClient: &GatewayBotGetterMock{
//			get: func() (gateway *GatewayBot, err error) {
//				return &GatewayBot{
//					Shards: nrOfShards,
//				}, nil
//			},
//		},
//		DefaultBotPresence: nil,
//		ProjectName:        "",
//	})
//}
