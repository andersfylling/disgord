// Disgord provides Go bindings for the documented Discord API. And allows for a stateful client using the Session interface, with the option of a configurable caching system or bypass the built-in caching logic all together.
//
// Getting started
//
// Create a Disgord session to get access to the REST API and socket functionality. In the following example, we listen for new messages and write a "hello" message when our handler function gets fired.
//
// Session interface: https://godoc.org/github.com/andersfylling/disgord/#Session
//  discord, err := disgord.NewSession(&disgord.Config{
//    Token: "my-secret-bot-token",
//  })
//  if err != nil {
//    panic(err)
//  }
//
//  // listen for incoming messages and reply with a "hello"
//  discord.On(event.MessageCreate, func(session disgord.Session, evt *disgord.MessageCreate) {
//      evt.Message.RespondString("hello")
//  })
//
//  // connect to the socket API to receive events
//  err = discord.Connect()
//  if err != nil {
//      panic(err)
//  }
//  discord.DisconnectOnInterrupt()
//
//
// Listen for events using channels
//
// Disgord also provides the option to listen for events using a channel, instead of registering a handler. However, before using the event channel, you must notify disgord that you care about the event (this is done automatically in the event handler registration).
//  session.AcceptEvent(event.MessageCreate) // alternative: disgord.EventMessageCreate
//  session.AcceptEvent(event.MessageUpdate)
//  for {
//      var message *disgord.Message
//      var status string
//      select {
//      case evt, alive := <- session.EventChannels().MessageCreate()
//          if !alive {
//              return
//          }
//          message = evt.Message
//          status = "created"
//      case evt, alive := <- session.EventChannels().MessageUpdate()
//          if !alive {
//              return
//          }
//          message = evt.Message
//          status = "updated"
//      }
//
//      fmt.Printf("A message from %s was %s\n", message.Author.Mention(), status)
//      // output example: "A message from @Anders was created"
//  }
//
// Optimizing your cache logic
//
// > Note: if you create a CacheConfig you must set the cache replacement algorithm for each cache (user, guild, etc.).
//
// A part of Disgord is the control you have; while this can be a good detail for advanced users, we recommend beginners to utilise the default configurations (by simply not editing the configuration).
// Here we pass the cache config when creating the session to access to the different cache replacement algorithms, lifetime settings, and the option to disable different cache systems.
//  discord, err := disgord.NewSession(&disgord.Config{
//    Token: "my-secret-bot-token",
//    Cache: &disgord.CacheConfig{
//              Immutable: true, // everything going in and out of the cache is deep copied
//
//              DisableUserCaching: false, // activates caching for users
//              UserCacheLimitMiB: 500, // don't use more than ~500MiB of memory space for caching of users
//              UserCacheLifetime: time.Duration(4) * time.Hour, // removed from cache after 9 hours, unless updated
//              UserCacheAlgorithm: disgord.CacheAlg_TLRU, // uses TLRU (Time aware Least Recently Used) for caching of users
//
//              DisableVoiceStateCaching: true, // don't cache voice states
//              // VoiceStateCacheLifetime  time.Duration
//              // VoiceStateCacheAlgorithm string
//
//              DisableChannelCaching: false,
//              ChannelCacheLimitMiB: 300,
//              ChannelCacheLifetime: 0, // lives forever
//              ChannelCacheAlgorithm: disgord.CacheAlg_LFU, // lfu (Least Frequently Used)
//           },
//  })
//
// > Note: Disabling caching for some types while activating it for others (eg. disabling channels, but activating guild caching), can cause items extracted from the cache to not reflect the true discord state.
//
// Example, activated guild but disabled channel caching: The guild is stored to the cache, but it's channels are discarded. Guild channels are dismantled from the guild object and otherwise stored in the channel cache to improve performance and reduce memory use. So when you extract the cached guild object, all of the channel will only hold their channel ID, and nothing more.
//
//
// Immutable and concurrent accessible cache
//
// The option CacheConfig.Immutable can greatly improve performance or break your system. If you utilize channels or you need concurrent access, the safest bet is to set immutable to `true`. While this is slower (as you create deep copies and don't share the same memory space with variables outside the cache), it increases reliability that the cache always reflects the last known Discord state.
// If you are uncertain, just set it to `true`. The default setting is `true` if `disgord.Cache.CacheConfig` is `nil`.
//
//
// Bypass the built-in REST cache
//
// Whenever you call a REST method from the Session interface; the cache is always checked first. Upon a cache hit, no REST request is executed and you get the data from the cache in return. However, if this is problematic for you or there exist a bug which gives you bad/outdated data, you can bypass it by using the REST functions directly. Remember that this will not update the cache for you, and this needs to be done manually if you depend on the cache.
//  // get a user using the Session implementation (checks cache, and updates the cache on cache miss)
//  user, err := session.GetUser(userID)
//
//  // bypass the cache checking. Same function name, but is found in the disgord package, not the session interface.
//  user, err := disgord.GetUser(userID)
//
// Manually updating the cache
//
// If required, you can access the cache and update it by hand. Note that this should not be required when you use the Session interface.
//  user, err := disgord.GetUser(userID)
//  if err != nil {
//      return err
//  }
//
//  // update the cache
//  cache := session.Cache()
//  err = cache.Update(disgord.UserCache, user)
//  if err != nil {
//      return err
//  }
//
// Saving and Deleting Discord data
//
// You might have seen the two methods in the session interface: SaveToDiscord(...) and DeleteFromDiscord(...).
// This are as straight forward as they sound. Passing a discord data structure into one of them executes their obvious behavior; to either save it to Discord, or delete it.
//  // create a new role and give it certain permissions
//  role := disgord.Role{}
//  role.Name = "Giraffes"
//  role.GuildID = guild.ID // required, for an obvious reason
//  role.Permissions = disgord.ManageChannelsPermission | disgord.ViewAuditLogsPermission
//  err := session.SaveToDiscord(&role)
//
// hang on! I don't want them to have the ManageChannel permission anyway
//  role.Permissions ^= disgord.ManageChannelsPermission // remove MANAGE_CHANNEL permission
//  err := session.SaveToDiscord(&role) // yes, it also updates existing objects
//
// You know what.. Let's just remove the role
//  err := session.DeleteFromDiscord(&role)
//
package disgord

import (
	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/websocket"
	"github.com/andersfylling/snowflake/v2"
)

const (
	// JSONEncoding const for JSON encoding type
	JSONEncoding = "json"

	// APIVersion desired API version to use
	APIVersion = 6 // February 5, 2018
	// DefaultAPIVersion the default Discord API version
	DefaultAPIVersion = 6

	// GitHubURL complete url for this project
	GitHubURL = "https://github.com/andersfylling/disgord"
)

// LibraryInfo returns name + version
func LibraryInfo() string {
	return "Disgord " + constant.Version
}

type DiscordWSEvent = websocket.DiscordWSEvent
type DiscordWebsocket = websocket.DiscordWebsocket

// Wrapper for github.com/andersfylling/snowflake
// ------------------

// Snowflake twitter snowflake identification for Discord
type Snowflake = snowflake.Snowflake

func GetSnowflake(v interface{}) (Snowflake, error) {
	s, err := snowflake.GetSnowflake(v)
	return Snowflake(s), err
}

func NewSnowflake(id uint64) Snowflake {
	return Snowflake(snowflake.NewSnowflake(id))
}

func ParseSnowflakeString(v string) Snowflake {
	return Snowflake(snowflake.ParseSnowflakeString(v))
}

func NewErrorMissingSnowflake(message string) *ErrorMissingSnowflake {
	return &ErrorMissingSnowflake{
		info: message,
	}
}

type ErrorMissingSnowflake struct {
	info string
}

func (e *ErrorMissingSnowflake) Error() string {
	return e.info
}

func NewErrorEmptyValue(message string) *ErrorEmptyValue {
	return &ErrorEmptyValue{
		info: message,
	}
}

type ErrorEmptyValue struct {
	info string
}

func (e *ErrorEmptyValue) Error() string {
	return e.info
}
