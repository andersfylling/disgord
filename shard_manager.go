package disgord

import (
	"context"
	"errors"
	"fmt"
	"github.com/andersfylling/discordgateway"
	"github.com/andersfylling/discordgateway/event"
	"github.com/andersfylling/discordgateway/gatewayshard"
	"github.com/andersfylling/discordgateway/intent"
	discordgatewaylog "github.com/andersfylling/discordgateway/log"
	"github.com/andersfylling/disgord/internal/constant"
	"github.com/andersfylling/disgord/internal/logger"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

type ShardManager interface {
	Connect() error
	Disconnect()
	SendCommand(cmd GatewayCommand) error
}

type GatewayCommand interface {
	GuildID() Snowflake
	OperationCode() uint
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
	ShardIDs []uint
	ShardCount uint
	Intents Intent
	BotToken string
	Log Logger

	setupOnce sync.Once
	shards []*gatewayshard.Shard
}

func (sm *BasicShardManager) setup() error {
	discordgatewaylog.LogInstance = sm.Log

	id := &discordgateway.IdentifyConnectionProperties{
		OS:      runtime.GOOS,
		Browser: "github.com/andersfylling/discordgateway v0",
		Device:  "github.com/andersfylling/disgord v0",
	}
	intents := intent.Type(sm.Intents)

	for i := range sm.ShardIDs {
		shardID := discordgateway.ShardID(sm.ShardIDs[i])
		shard, err := gatewayshard.NewShard(shardID, sm.BotToken, nil,
			discordgateway.WithIntents(intents),
			discordgateway.WithIdentifyConnectionProperties(id),
		)
		if err != nil {
			log.Fatal(err)
		}

		sm.shards = append(sm.shards, shard)
	}

	return nil
}

func (sm *BasicShardManager) Connect(ctx context.Context) error{
	for {
		err := sm.connect(ctx)
		if err
	}

}

func (sm *BasicShardManager) connect(parentCtx context.Context) (err error) {
	sm.setupOnce.Do(func () {
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
					if err := RunShard(shard, workChan); err != nil {
						return err
					}
				}
			}
		})
	}

	return g.Wait()
}

func (sm *BasicShardManager) Disconnect() error {

}

func (sm *BasicShardManager) SendCommand(cmd GatewayCommand) error {

}