package disgord

import (
	"context"
	"errors"
	"fmt"
	"github.com/andersfylling/discordgateway"
	"github.com/andersfylling/discordgateway/closecode"
	"github.com/andersfylling/discordgateway/command"
	"github.com/andersfylling/discordgateway/event"
	"github.com/andersfylling/discordgateway/gatewayshard"
	"github.com/andersfylling/discordgateway/intent"
	discordgatewaylog "github.com/andersfylling/discordgateway/log"
	"github.com/andersfylling/discordgateway/opcode"
	"github.com/andersfylling/disgord/internal/constant"
	"github.com/andersfylling/disgord/internal/gateway"
	"github.com/andersfylling/disgord/internal/logger"
	"github.com/andersfylling/disgord/json"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

type ShardManager interface {
	Connect(ctx context.Context) error
	Disconnect()
	SendCommand(ctx context.Context, cmd GatewayCommand) error
}

type GatewayCommand interface {
	GuildID() Snowflake
	OperationCode() uint
	CommandCode() uint
}

func runShard(shard *gatewayshard.Shard, workChan chan<- *gatewayshard.Shard) error {
	websocketUrl := fmt.Sprintf("wss://gateway.discord.gg/?v=%d&encoding=%s", constant.DiscordVersion, strings.ToLower(constant.Encoding))
	if _, err := shard.Dial(context.Background(), websocketUrl); err != nil {
		return fmt.Errorf("failed to open websocket connection. %w", err)
	}

	// process websocket messages as they arrive and trigger the handler whenever relevant
	if err := shard.EventLoop(context.Background()); err != nil {
		reconnect := true

		var discordErr *discordgateway.DiscordError
		if errors.As(err, &discordErr) {
			reconnect = discordErr.CanReconnect()
		}

		if reconnect {
			log.Info(fmt.Errorf("reconnecting: %w", err))
			if err := shard.PrepareForReconnect(); err != nil {
				return fmt.Errorf("failed to prepare for reconnect: %w", err)
			}
			workChan <- shard
			return nil
		}

		return err
	}

	return errors.New("unexpected error from shards event loop, no error was returned")
}

type BasicShardManager struct {
	ShardIDs   []uint
	ShardCount uint
	Intents    Intent
	BotToken   string
	Log        Logger
	Client     Session

	shardingConfiguredByUser bool
	setupOnce                sync.Once
	shards                   []*gatewayshard.Shard
	identityProperties       *discordgateway.IdentifyConnectionProperties
}

func (sm *BasicShardManager) setupShardDetails() error {
	gatewayBotInfo, err := sm.getGatewayBotInfo()
	if err != nil {
		return err
	}

	for i := 0; i < int(gatewayBotInfo.Shards); i++ {
		sm.ShardIDs = append(sm.ShardIDs, uint(i))
	}
	sm.ShardCount = gatewayBotInfo.Shards
	// TODO: rate limits
	return nil
}

func (sm *BasicShardManager) setup() error {
	discordgatewaylog.LogInstance = sm.Log

	if int(sm.ShardCount) < len(sm.ShardIDs) {
		return errors.New("shard count is less than the number of specified shard ids")
	}
	if sm.ShardCount > 0 && len(sm.ShardIDs) == 0 {
		return errors.New("shard ids must be specified when using setting shard count")
	}
	if sm.ShardIDs != nil {
		sm.shardingConfiguredByUser = true
	} else {
		if err := sm.setupShardDetails(); err != nil {
			return err
		}
	}

	sm.identityProperties = &discordgateway.IdentifyConnectionProperties{
		OS:      runtime.GOOS,
		Browser: "github.com/andersfylling/discordgateway v0",
		Device:  "github.com/andersfylling/disgord v0",
	}

	return sm.setupShards()
}

func (sm *BasicShardManager) setupShards() error {
	intents := intent.Type(sm.Intents)
	for i := range sm.ShardIDs {
		shardID := discordgateway.ShardID(sm.ShardIDs[i])
		shard, err := gatewayshard.NewShard(shardID, sm.BotToken, nil,
			discordgateway.WithIntents(intents),
			discordgateway.WithIdentifyConnectionProperties(sm.identityProperties),
			discordgateway.WithShardCount(sm.ShardCount),
		)
		if err != nil {
			return err
		}

		sm.shards = append(sm.shards, shard)
	}

	return nil
}

func (sm *BasicShardManager) getGatewayBotInfo() (*gateway.GatewayBot, error) {
	return sm.Client.Gateway().GetBot()
}

func (sm *BasicShardManager) Connect(ctx context.Context) error {
	for {
		err := sm.connect(ctx)
		var discordErr *discordgateway.DiscordError
		if errors.As(err, &discordErr) && discordErr.CloseCode == closecode.ShardingRequired {
			// figure out if the sharding information was specified by user or discord
			if sm.shardingConfiguredByUser {
				// user has to fix this
				return err
			}

			// otherwise we automatically increment the number of shards
			if err = sm.setupShardDetails(); err != nil {
				return err
			}
			if err = sm.setupShards(); err != nil {
				return err
			}
		} else {
			return err
		}
	}
}

func (sm *BasicShardManager) connect(parentCtx context.Context) (err error) {
	sm.setupOnce.Do(func() {
		err = sm.setup()
	})
	if err != nil {
		return fmt.Errorf("failed to setup shard manager: %w", err)
	}

	g, ctx := errgroup.WithContext(parentCtx)
	workChan := make(chan *gatewayshard.Shard, len(sm.shards))
	for _, shard := range sm.shards {
		workChan <- shard
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case shard, ok := <-workChan:
					if !ok {
						return errors.New("work channel for shard manager unexpectedly closed")
					}
					if err := runShard(shard, workChan); err != nil {
						return err
					}
				}
			}
		})
	}

	return g.Wait()
}

func (sm *BasicShardManager) Disconnect() {
	for i := range sm.shards {
		_ = sm.shards[i].Close()
	}
}

func (sm *BasicShardManager) SendCommand(parentCtx context.Context, cmd GatewayCommand) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	if !cmd.GuildID().IsZero() {
		shardID := discordgateway.DeriveShardID(uint64(cmd.GuildID()), uint(len(sm.shards)))
		shard := sm.shards[shardID]
		return shard.Write(command.Type(cmd.CommandCode()), data)
	}

	g, ctx := errgroup.WithContext(parentCtx)
	workChan := make(chan *gatewayshard.Shard, len(sm.shards))
	for i := range sm.shards {
		workChan <- sm.shards[i]
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case shard, ok := <-workChan:
				if !ok {
					return errors.New("work chan for sending command suddenly closed")
				}
				return shard.Write(command.Type(cmd.OperationCode()), data)
			}
		})
	}

	return g.Wait()
}
