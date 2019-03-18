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

	// Large bots only. If Discord did not give you a custom rate limit, do not touch this.
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
		Presence:          conf.Presence,
		TrackEvent:        &websocket.UniqueStringSlice{},
		discordPktPool: &sync.Pool{
			New: func() interface{} {
				return &websocket.DiscordPacket{}
			},
		},
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

	// Presence represents the desired bot status at any given time
	Presence *UpdateStatusCommand

	prepared bool
	log      Logger

	discordPktPool *sync.Pool

	client *client // hacky - used to register handlers
}

func (s *WSShardManager) GetConnectionDetails(c httd.Getter) (url string, shardCount uint, err error) {
	var d *GatewayBot
	if d, err = GetGatewayBot(c); err != nil {
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
		err = s.shards[i].Prepare(conf, s.discordPktPool, s.evtChan, s.conRequestChan, s.TrackEvent, conf.WSShardManagerConfig.FirstID+uint(i))
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

	id := GetShardForGuildID(guildID, uint(len(s.shards)))
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
		tmpErr := s.shards[i].Connect()
		if tmpErr != nil {
			err = tmpErr
			s.log.Error(err)
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
	if guild, ok := data.(guilder); ok {
		guildID := guild.getGuildID()
		shard, err := s.GetShard(guildID)
		if err != nil {
			return err
		}

		return shard.Emit(cmd, data)
	}

	s.RLock()
	defer s.RUnlock()
	for i := 0; i < len(s.shards); i++ {
		err = s.shards[i].ws.Emit(cmd, data)
	}

	return
}

//
//func (s *WSShardManager) UpdatePresence(command *UpdateStatusCommand) error {
//	s.Presence.mu.Lock()
//	command.mu = s.Presence.mu
//	*s.Presence = *command // don't change the pointer as each shard has a copy
//	s.Presence.mu.Unlock()
//
//	s.RLock()
//	for i := range s.shards {
//		s.shards[i].ws.SetPresence(command)
//	}
//	s.RUnlock()
//
//	go func() {
//		// this is just to ensure that the bot is updated on all shards
//
//		s.RLock()
//		updated := make([]bool, len(s.shards))
//		// TODO: before auto-scaling of shards is added, a token or hash must be added to check if the same shards
//		//  are used later on in the update confirmations
//		s.RUnlock()
//
//		for {
//			timeout := time.Now().Add(10 * time.Second)
//			done := make(chan interface{})
//
//			_ = s.client.On(
//				EventPresenceUpdate,
//				// middleware
//				func(evt interface{}) interface{} {
//					e := evt.(*PresenceUpdate)
//					if e.User.ID != s.client.myID {
//						return nil // don't proceed if this regards another user
//					}
//					return evt
//				},
//				// handler
//				func(s Session, evt *PresenceUpdate) {
//					if int(evt.ShardID) >= len(updated) {
//						s.Logger().Error("the ShardID does not exists. Got", evt.ShardID, "wants below", len(updated))
//						return
//					}
//
//					updated[evt.ShardID] = true
//					for i := range updated {
//						if updated[i] == false {
//							return
//						}
//					}
//
//					close(done)
//				},
//				// ctrl
//				&timeoutHandlerCtrl{timeout},
//			)
//
//			// either continue when all the shards completed, or wait for a system interrupt or timeout
//			select {
//			case <-s.client.shutdownChan:
//				return
//			case <-time.After(time.Now().Sub(timeout)):
//			case <-done:
//			}
//
//			// TODO: check shards hash here to make sure that these are the same shards, in the same order
//			//  as earlier.
//			for i := range updated {
//				if !updated[i] {
//					// retry updating the shard
//					// TODO: if the shard disconnected.. it would send a identify with the correct presence info -
//					//   how can this be checked, to avoid redundant traffic?
//					_ = s.shards[i].Emit(CommandUpdateStatus, command)
//				}
//			}
//		}
//	}()
//
//	_ = s.Emit(CommandUpdateStatus, command)
//
//	return nil
//}

var _ Emitter = (*WSShardManager)(nil)
var _ Link = (*WSShardManager)(nil)

type WSShard struct {
	sync.RWMutex
	id    uint
	total uint

	ws     *websocket.EvtClient
	guilds []snowflake.ID
}

func (s *WSShard) Prepare(conf *Config, discordPktPool *sync.Pool, evtChan chan *websocket.Event, conRequestChan websocket.A, trackEvents *websocket.UniqueStringSlice, id uint) (err error) {
	s.id = id
	s.total = conf.WSShardManagerConfig.ShardLimit

	s.ws, err = websocket.NewEventClient(&websocket.EvtConfig{
		// identity
		Browser:             LibraryInfo(),
		Device:              conf.ProjectName,
		GuildLargeThreshold: 250,
		ShardCount:          s.total,
		Presence:            conf.Presence,

		// lib specific
		Version:        constant.DiscordVersion,
		Encoding:       constant.JSONEncoding,
		ChannelBuffer:  3,
		Endpoint:       conf.WSShardManagerConfig.URL,
		EventChan:      evtChan,
		TrackedEvents:  trackEvents,
		Logger:         conf.Logger,
		A:              conRequestChan,
		DiscordPktPool: discordPktPool,

		// user settings
		BotToken: conf.BotToken,
		Proxy:    conf.Proxy,

		SystemShutdown: conf.shutdownChan,
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
