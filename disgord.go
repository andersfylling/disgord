// Package disgord provides Go bindings for the documented Discord API. And allows for a stateful Client using the Session interface, with the option of a configurable caching system or bypass the built-in caching logic all together.
//
// Getting started
//
// Create a DisGord session to get access to the REST API and socket functionality. In the following example, we listen for new messages and write a "hello" message when our handler function gets fired.
//
// Session interface: https://godoc.org/github.com/andersfylling/disgord/#Session
//  discord := disgord.New(&disgord.Config{
//    Token: "my-secret-bot-token",
//  })
//  defer discord.StayConnectedUntilInterrupt()
//
//  // listen for incoming messages and reply with a "hello"
//  discord.On(event.MessageCreate, func(s disgord.Session, evt *disgord.MessageCreate) {
//      msg := evt.Message
//      msg.Reply(s, "hello")
//  })
//
// // If you want some logic to fire when the bot is ready
// // (all shards has received their ready event), please use the Ready method.
//  discord.Ready(func() {
//  	fmt.Println("READY NOW!")
//  })
//
//
//
// Listen for events using channels
//
// Disgord also provides the option to listen for events using a channel. The setup is exactly the same as registering a handler.
// Simply define your channel, add buffering if you need it, and register it as a handler in the .On method.
//
//  msgCreateChan := make(chan *disgord.MessageCreate, 10)
//  session.On(disgord.EvtMessageCreate, msgCreateChan)
//
// Never close a channel without removing the handler from disgord. You can't directly call Remove, instead you can
// inject a controller to dictate the handler's lifetime. Since you are the owner of the channel, disgord will not
// close it for you.
//
//  ctrl := &disgord.Ctrl{Channel: msgCreateChan}
//  go func() {
//    // close the channel after 20 seconds and safely remove it from disgord
//    // without disgord trying to send data through it after it has closed
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
//          message = evt.Message
//          status = "created"
//      }
//
//      fmt.Printf("A message from %s was %s\n", message.Author.Mention(), status)
//      // output example: "A message from @Anders was created"
//  }
//
// Optimizing your cache logic
//
// > Note: if you create a CacheConfig you don't have to set every field. All the CacheAlgorithms are default to LFU when left blank.
//
// A part of Disgord is the control you have; while this can be a good detail for advanced users, we recommend beginners to utilise the default configurations (by simply not editing the configuration).
// Here we pass the cache config when creating the session to access to the different cache replacement algorithms, lifetime settings, and the option to disable different cache systems.
//  discord, err := disgord.NewClient(&disgord.Config{
//    BotToken: "my-secret-bot-token",
//    Cache: &disgord.CacheConfig{
//              Mutable: false, // everything going in and out of the cache is deep copied
//				// setting Mutable to true, might break your program as this is experimental.
//
//              DisableUserCaching: false, // activates caching for users
//              UserCacheLifetime: time.Duration(4) * time.Hour, // removed from cache after 9 hours, unless updated
//              UserCacheAlgorithm: disgord.CacheAlgLFU,
//
//              DisableVoiceStateCaching: true, // don't cache voice states
//              // VoiceStateCacheLifetime  time.Duration
//              // VoiceStateCacheAlgorithm string
//
//              DisableChannelCaching: false,
//              ChannelCacheLifetime: 0, // lives forever
//              ChannelCacheAlgorithm: disgord.CacheAlgLFU, // lfu (Least Frequently Used)
//
//				GuildCacheAlgorithm: disgord.CacheAlgLFU, // no limit set, so the strategy to replace entries is not used
//           },
//  })
//
// If you just want to change a specific field, you can do so. By either calling the disgord.DefaultCacheConfig which gives you a Cache configuration designed by DisGord. Or you can set specific fields in a new CacheConfig since the different Cache Strategies are automatically set to LFU if missing.
// 	&disgord.Config{}
// Will automatically become
//  &disgord.Config{
//  	UserCacheAlgorithm: disgord.CacheAlgLFU,
//		VoiceStateCacheAlgorithm disgord.CacheAlgLFU,
//		ChannelCacheAlgorithm: disgord.CacheAlgLFU,
//		GuildCacheAlgorithm: disgord.CacheAlgLFU,
//  }
//
// And writing
//  &disgord.Config{
//  	UserCacheAlgorithm: disgord.CacheAlgLRU,
//		VoiceStateCacheAlgorithm disgord.CacheAlgLRU,
//  }
// Becomes
//  &disgord.Config{
//  	UserCacheAlgorithm: disgord.CacheAlgLRU, // unchanged
//		VoiceStateCacheAlgorithm disgord.CacheAlgLRU,  // unchanged
//		ChannelCacheAlgorithm: disgord.CacheAlgLFU,
//		GuildCacheAlgorithm: disgord.CacheAlgLFU,
//  }
//
// > Note: Disabling caching for some types while activating it for others (eg. disabling channels, but activating guild caching), can cause items extracted from the cache to not reflect the true discord state.
//
// Example, activated guild but disabled channel caching: The guild is stored to the cache, but it's channels are discarded. Guild channels are dismantled from the guild object and otherwise stored in the channel cache to improve performance and reduce memory use. So when you extract the cached guild object, all of the channel will only hold their channel ID, and nothing more.
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
// Whenever you call a REST method from the Session interface; the cache is always checked first. Upon a cache hit, no REST request is executed and you get the data from the cache in return. However, if this is problematic for you or there exist a bug which gives you bad/outdated data, you can bypass it by using disgord flags.
//  // get a user using the Session implementation (checks cache, and updates the cache on cache miss)
//  user, err := session.GetUser(userID)
//
//  // bypass the cache checking. Same as before, but we insert a disgord.Flag type.
//  user, err := session.GetUser(userID, disgord.IgnoreCache)
//
// DisGord Flags
//
// In addition to disgord.IgnoreCache, as shown above, you can pass in other flags such as: disgord.SortByID, disgord.OrderAscending, etc. You can find these flags in the flag.go file.
//
// Manually updating the cache
//
// Currently not supported. Should it ever be?
//
//
// Build tags
//
// `disgord_diagnosews` will store all the incoming and outgoing json data as files in the directory "diagnose-report/packets". The file format is as follows: unix_clientType_direction_shardID_operationCode_sequenceNumber[_eventName].json
//
// `json-std` switches out jsoniter with the json package from the std libs.
//
// `disgord_removeDiscordMutex` replaces mutexes in discord structures with a empty mutex; removes locking behaviour and any mutex code when compiled.
//
// `disgord_parallelism` activates built-in locking in discord structure methods. Eg. Guild.AddChannel(*Channel) does not do locking by default. But if you find yourself using these discord data structures in parallel environment, you can activate the internal locking to reduce race conditions. Note that activating `disgord_parallelism` and `disgord_removeDiscordMutex` at the same time, will cause you to have no locking as `disgord_removeDiscordMutex` affects the same mutexes.
//
//
// Deleting Discord data
//
// In addition to the typical REST endpoints for deleting data, you can also use Client/Session.DeleteFromDiscord(...) for basic deletions. If you need to delete a specific range of messages, or anything complex as that; you can't use .DeleteFromDiscord(...). Not every struct has implemented the interface that allows you to call DeleteFromDiscord. Do not fret, if you try to pass a type that doesn't qualify, you get a compile error.
//
package disgord

import (
	"fmt"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/snowflake/v3"
)

// LibraryInfo returns name + version
func LibraryInfo() string {
	return fmt.Sprint(constant.Name, constant.Version)
}

// Wrapper for github.com/andersfylling/snowflake
// ------------------

// Snowflake twitter snowflake identification for Discord
type Snowflake = snowflake.Snowflake

// GetSnowflake see snowflake.GetSnowflake
func GetSnowflake(v interface{}) (Snowflake, error) {
	s, err := snowflake.GetSnowflake(v)
	return Snowflake(s), err
}

// NewSnowflake see snowflake.NewSnowflake
func NewSnowflake(id uint64) Snowflake {
	return Snowflake(snowflake.NewSnowflake(id))
}

// ParseSnowflakeString see snowflake.ParseSnowflakeString
func ParseSnowflakeString(v string) Snowflake {
	return Snowflake(snowflake.ParseSnowflakeString(v))
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
