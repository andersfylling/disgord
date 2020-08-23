// Package disgord provides Go bindings for the documented Discord API, and allows for a stateful Client using the Session interface, with the option of a configurable caching system or bypass the built-in caching logic all together.
//
// Getting started
//
// Create a Disgord session to get access to the REST API and socket functionality. In the following example, we listen for new messages and write a "hello" message when our handler function gets fired.
//
// Session interface: https://pkg.go.dev/github.com/andersfylling/disgord?tab=doc#Session
//  discord := disgord.New(&disgord.Config{
//    BotToken: "my-secret-bot-token",
//  })
//  defer discord.StayConnectedUntilInterrupted()
//
//  // listen for incoming messages and reply with a "hello"
//  discord.On(event.MessageCreate, func(s disgord.Session, evt *disgord.MessageCreate) {
//      msg := evt.Message
//      msg.Reply(s, "hello")
//  })
//
//  // If you want some logic to fire when the bot is ready
//  // (all shards has received their ready event), please use the Ready method.
//  discord.Ready(func() {
//  	fmt.Println("READY NOW!")
//  })
//
//
//
// Listen for events using Channels
//
// Disgord also provides the option to listen for events using a channel. The setup is exactly the same as registering a function.
// Simply define your channel, add buffering if you need it, and register it as a handler in the `.On` method.
//
//  msgCreateChan := make(chan *disgord.MessageCreate, 10)
//  session.On(disgord.EvtMessageCreate, msgCreateChan)
//
// Never close a channel without removing the handler from Disgord. You can't directly call Remove, instead you
// inject a controller to dictate the handler's lifetime. Since you are the owner of the channel, disgord will never
// close it for you.
//
//  ctrl := &disgord.Ctrl{Channel: msgCreateChan}
//  session.On(disgord.EvtMessageCreate, msgCreateChan, ctrl)
//  go func() {
//    // close the channel after 20 seconds and safely remove it from Disgord
//    // without Disgord trying to send data through it after it has closed
//    <- time.After(20 * time.Second)
//    ctrl.CloseChannel()
//  }
//
// Here is what it would look like to use the channel for handling events. Please run this in a go routine unless
// you know what you are doing.
//
//  for {
//      var message *disgord.Message
//      var status string
//      select {
//      case evt, alive := <- msgCreateChan
//          if !alive {
//              return
//          }
//          msg = evt.Message
//          status = "created"
//      }
//
//      fmt.Printf("A message from %s was %s\n", msg.Author.Mention(), status)
//      // output: "A message from @Anders was created"
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
//  client := disgord.New(&disgord.Config{
//    ShardConfig: disgord.ShardConfig{
//      // this is a copy so u can't manipulate the config later on
//      ShardIDs: []uint{0,1,2,3,4},
//    },
//  })
//
// Running multiple instances each with 1 shard (note each instance must use unique shard ids)
//  client := disgord.New(&disgord.Config{
//    ShardConfig: disgord.ShardConfig{
//      // this is a copy so u can't manipulate the config later on
//      ShardIDs: []uint{0}, // this number must change for each instance. Try to automate this.
//      ShardCount: 5, // total of 5 shards, but this disgord instance only has one. AutoScaling is disabled - use OnScalingRequired.
//    },
//  })
//
// Handle scaling options yourself
//  client := disgord.New(&disgord.Config{
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
// > Note: if you create a CacheConfig you don't have to set every field.
//
// > Note: Only LFU is supported.
//
// > Note: Lifetime options does not currently work/do anything (yet).
//
// A part of Disgord is the control you have; while this can be a good detail for advanced Users, we recommend beginners to utilise the default configurations (by simply not editing the configuration).
// Example of configuring the cache:
//  discord, err := disgord.NewClient(&disgord.Config{
//    BotToken: "my-secret-bot-token",
//    CacheDefault: &disgord.CacheConfig{
//              Mutable: false, // everything going in and out of the cache is deep copied
//				// setting Mutable to true, might break your program as this is experimental and not supported.
//
//              DisableUserCaching: false, // activates caching for Users
//              UserCacheLifetime: time.Duration(4) * time.Hour, // removed from cache after 9 hours, unless updated
//
//              DisableVoiceStateCaching: true, // don't cache voice states
//
//              DisableChannelCaching: false,
//              ChannelCacheLifetime: 0, // lives forever unless cache replacement strategy kicks in
//           },
//  })
//
// If you just want to change a specific field, you can do so. The fields are always default values.
//
// > Note: Disabling caching for some types while activating it for others (eg. disabling Channels, but activating guild caching), can cause items extracted from the cache to not reflect the true discord state.
//
// Example, activated guild but disabled channel caching: The guild is stored to the cache, but it's Channels are discarded. Guild Channels are dismantled from the guild object and otherwise stored in the channel cache to improve performance and reduce memory use. So when you extract the cached guild object, all of the channel will only hold their channel ID, and nothing more.
//
//
// Immutable cache
//
// To keep it safe and reliable, you can not directly affect the contents of the cache. Unlike discordgo where everything is mutable, the caching in disgord is immutable. This does reduce performance as a copy must be made (only on new cache entries), but as a performance freak, I can tell you right now that a simple struct copy is not that expensive. This also means that, as long as discord sends their events properly, the caching will always reflect the true state of discord.
//
// If there is a bug in the cache and you keep getting the incorrect data, please file an issue at github.com/andersfylling/disgord so it can quickly be resolved(!)
//
// Bypass the built-in REST cache
//
// Whenever you call a REST method from the Session interface; the cache is always checked first. Upon a cache hit, no REST request is executed and you get the data from the cache in return. However, if this is problematic for you or there exist a bug which gives you bad/outdated data, you can bypass it by using Disgord flags.
//  // get a user using the Session implementation (checks cache, and updates the cache on cache miss)
//  user, err := session.GetUser(userID)
//
//  // bypass the cache checking. Same as before, but we insert a disgord.Flag type.
//  user, err := session.GetUser(userID, disgord.IgnoreCache)
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
// `json_std` switches out jsoniter with the json package from the std libs.
//
// `disgord_removeDiscordMutex` replaces mutexes in discord structures with a empty mutex; removes locking behaviour and any mutex code when compiled.
//
// `disgord_parallelism` activates built-in locking in discord structure methods. Eg. Guild.AddChannel(*Channel) does not do locking by default. But if you find yourself using these discord data structures in parallel environment, you can activate the internal locking to reduce race conditions. Note that activating `disgord_parallelism` and `disgord_removeDiscordMutex` at the same time, will cause you to have no locking as `disgord_removeDiscordMutex` affects the same mutexes.
//
// `disgord_legacy` adds wrapper methods with the original discord naming. eg. For REST requests you will notice Disgord uses a consistency between update/create/get/delete/set while discord uses edit/update/modify/close/delete/remove/etc. So if you struggle find a REST method, you can enable this build tag to gain access to mentioned wrappers.
//
// `disgordperf` does some low level tweaking that can help boost json unmarshalling and drops json validation from Discord responses/events. Other optimizations might take place as well.
//
// `disgord_websocket_gorilla` replaces nhooyr/websocket dependency with gorilla/websocket for gateway communication.
//
//
// Deleting Discord data
//
// In addition to the typical REST endpoints for deleting data, you can also use Client/Session.DeleteFromDiscord(...) for basic deletions. If you need to delete a specific range of messages, or anything complex as that; you can't use .DeleteFromDiscord(...). Not every struct has implemented the interface that allows you to call DeleteFromDiscord. Do not fret, if you try to pass a type that doesn't qualify, you get a compile error.
//
package disgord

//go:generate go run generate/intents/main.go

import (
	"fmt"
	"github.com/andersfylling/disgord/json"

	"github.com/andersfylling/disgord/internal/util"

	"github.com/andersfylling/disgord/internal/constant"
)

const Name = constant.Name
const Version = constant.Version

// LibraryInfo returns name + version
func LibraryInfo() string {
	return fmt.Sprintf("%s %s", constant.Name, constant.Version)
}

var defaultUnmarshaler = json.Unmarshal
var defaultMarshaler = json.Marshal

// Wrapper for github.com/andersfylling/snowflake
// ------------------

// Snowflake twitter snowflake identification for Discord
type Snowflake = util.Snowflake

// GetSnowflake see snowflake.GetSnowflake
func GetSnowflake(v interface{}) (Snowflake, error) {
	return util.GetSnowflake(v)
}

// NewSnowflake see snowflake.NewSnowflake
func NewSnowflake(id uint64) Snowflake {
	return util.NewSnowflake(id)
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
