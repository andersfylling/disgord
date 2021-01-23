// Package disgord provides Go bindings for the documented Discord API, and allows for a stateful Client using the Session interface, with the option of a configurable caching system or bypass the built-in caching logic all together.
//
// Getting started
//
// Create a Disgord client to get access to the REST API and gateway functionality. In the following example, we listen for new messages and respond with "hello".
//
// Session interface: https://pkg.go.dev/github.com/andersfylling/disgord?tab=doc#Session
//  client := disgord.New(disgord.Config{
//    BotToken: "my-secret-bot-token",
//  })
//  defer client.Gateway().StayConnectedUntilInterrupted()
//
//  client.Gateway().MessageCreate(func(s disgord.Session, evt *disgord.MessageCreate) {
//    evt.Message.Reply(context.Background(), s, "hello")
//  })
//
//
// Listen for events using Channels
//
// You don't have to use a callback function, channels are supported too!
//
//  msgChan := make(chan *disgord.MessageCreate, 10)
//  client.Gateway().MessageCreateChan(msgChan)
//
// Never close a channel without removing the handler from Disgord, as it will cause a panic. You can control the
// lifetime of a handler or injected channel by in injecting a controller: disgord.HandlerCtrl. Since you are the
// owner of the channel, disgord will never close it for you.
//
//  ctrl := &disgord.Ctrl{Channel: msgCreateChan}
//  client.Gateway().WithCtrl(ctrl).MessageCreateChan(msgChan)
//  go func() {
//    // close the channel after 20 seconds and safely remove it from the Disgord reactor
//    <- time.After(20 * time.Second)
//    ctrl.CloseChannel()
//  }
//
//
// WebSockets and Sharding
//
// Disgord handles sharding for you automatically; when starting the bot, when discord demands you to scale up your shards (during runtime), etc. It also gives you control over the shard setup in case you want to run multiple instances of Disgord (in these cases you must handle scaling yourself as Disgord can not).
//
// Sharding is done behind the scenes, so you do not need to worry about any settings. Disgord will simply ask Discord for the recommended amount of shards for your bot on startup. However, to set specific amount of shards you can use the `disgord.ShardConfig` to specify a range of valid shard IDs (starts from 0).
//
// starting a bot with exactly 5 shards
//  client := disgord.New(disgord.Config{
//    ShardConfig: disgord.ShardConfig{
//      // this is a copy so u can't manipulate the config later on
//      ShardIDs: []uint{0,1,2,3,4},
//    },
//  })
//
// Running multiple instances each with 1 shard (note each instance must use unique shard ids)
//  client := disgord.New(disgord.Config{
//    ShardConfig: disgord.ShardConfig{
//      // this is a copy so u can't manipulate the config later on
//      ShardIDs: []uint{0}, // this number must change for each instance. Try to automate this.
//      ShardCount: 5, // total of 5 shards, but this disgord instance only has one. AutoScaling is disabled - use OnScalingRequired.
//    },
//  })
//
// Handle scaling options yourself
//  client := disgord.New(disgord.Config{
//    ShardConfig: disgord.ShardConfig{
//      // this is a copy so u can't manipulate it later on
//      DisableAutoScaling: true,
//      OnScalingRequired: func(shardIDs []uint) (TotalNrOfShards uint, AdditionalShardIDs []uint) {
//        // instead of asking discord for exact number of shards recommended
//        // this is increased by 50% every time discord complains you don't have enough shards
//        // to reduce the number of times you have to scale
//        TotalNrOfShards := uint(len(shardIDs) * 1.5)
//        for i := len(shardIDs) - 1; i < TotalNrOfShards; i++ {
//          AdditionalShardIDs = append(AdditionalShardIDs, i)
//        }
//        return
//      }, // end OnScalingRequired
//    }, // end ShardConfig
//  })
//
//
// Caching
//
// You can inject your own cache implementation. By default a read only LFU implementation is used, this should be
// sufficient for the average user. But you can overwrite certain methods as well!
//
// Say you dislike the implementation for MESSAGE_CREATE events, you can embed the default cache and define your own
// logic:
//
//  type MyCoolCache struct {
//    disgord.CacheLFUImmutable
//    msgCache map[Snowflake]*Message // channelID => Message
//  }
//  func (c *CacheLFUImmutable) MessageCreate(data []byte) (*MessageCreate, error) {
//	  // some smart implementation here
//  }
//
// > Note: if you inject your own cache, remember that the cache is also responsible for initiating the objects.
// > See disgord.CacheNop
//
//
// Bypass the built-in REST cache
//
// Whenever you call a REST method from the Session interface; the cache is always checked first. Upon a cache hit, no REST request is executed and you get the data from the cache in return. However, if this is problematic for you or there exist a bug which gives you bad/outdated data, you can bypass it by using Disgord flags.
//  // get a user using the Session implementation (checks cache, and updates the cache on cache miss)
//  user, err := client.User(userID).Get()
//
//  // bypass the cache checking. Same as before, but we insert a disgord.Flag type.
//  user, err := client.User(userID).Get(disgord.IgnoreCache)
//
// Disgord Flags
//
// In addition to disgord.IgnoreCache, as shown above, you can pass in other flags such as: disgord.SortByID, disgord.OrderAscending, etc. You can find these flags in the flag.go file.
//
//
// Build tags
//
// `disgord_diagnosews` will store all the incoming and outgoing JSON data as files in the directory "diagnose-report/packets". The file format is as follows: unix_clientType_direction_shardID_operationCode_sequenceNumber[_eventName].json
//
// `disgordperf` does some low level tweaking that can help boost json unmarshalling and drops json validation from Discord responses/events. Other optimizations might take place as well.
//
//
package disgord

//go:generate go run internal/generate/intents/main.go

import (
	"fmt"

	"github.com/andersfylling/disgord/internal/util"

	"github.com/andersfylling/disgord/internal/constant"
)

const Name = constant.Name
const Version = constant.Version

// LibraryInfo returns name + version
func LibraryInfo() string {
	return fmt.Sprintf("%s %s", constant.Name, constant.Version)
}

// DeepCopier holds the DeepCopy method which creates and returns a deep copy of
// any struct.
type DeepCopier interface {
	deepCopy() interface{}
}

func DeepCopy(cp DeepCopier) interface{} {
	return cp.deepCopy()
}

// Copier holds the CopyOverTo method which copies all it's content from one
// struct to another. Note that this requires a deep copy.
// useful when overwriting already existing content in the cacheLink to reduce GC.
type Copier interface {
	copyOverTo(other interface{}) error
}

func DeepCopyOver(dst Copier, src Copier) error {
	// TODO: make sure dst and src are of the same type!
	return src.copyOverTo(dst)
}

// Reseter Reset() zero initialises or empties a struct instance
type Reseter interface {
	reset()
}

func Reset(r Reseter) {
	r.reset()
}

// Wrapper for snowflake
// ------------------

// Snowflake twitter snowflake identification for Discord
type Snowflake = util.Snowflake

// GetSnowflake see snowflake.GetSnowflake
func GetSnowflake(v interface{}) (Snowflake, error) {
	return util.GetSnowflake(v)
}

// ParseSnowflakeString see snowflake.ParseSnowflakeString
func ParseSnowflakeString(v string) Snowflake {
	return util.ParseSnowflakeString(v)
}

func newErrorMissingSnowflake(message string) *ErrorMissingSnowflake {
	return &ErrorMissingSnowflake{
		info: message,
	}
}

// ErrorMissingSnowflake used by methods about to communicate with the Discord API. If a snowflake value is required
// this is used to identify that you must set the value before being able to interact with the Discord API
type ErrorMissingSnowflake struct {
	info string
}

func (e *ErrorMissingSnowflake) Error() string {
	return e.info
}

func newErrorEmptyValue(message string) *ErrorEmptyValue {
	return &ErrorEmptyValue{
		info: message,
	}
}

// ErrorEmptyValue when a required value was set as empty
type ErrorEmptyValue struct {
	info string
}

func (e *ErrorEmptyValue) Error() string {
	return e.info
}
