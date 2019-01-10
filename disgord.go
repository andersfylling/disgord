// Package disgord provides Go bindings for the documented Discord API. And allows for a stateful client using the Session interface, with the option of a configurable caching system or bypass the built-in caching logic all together.
//
// Getting started
//
// Create a Disgord session to get access to the REST API and socket functionality. In the following example, we listen for new messages and write a "hello" message when our handler function gets fired.
//
// Session interface: https://godoc.org/github.com/andersfylling/disgord/#Session
//  discord, err := disgord.NewSession(&disgord.Config{
//    BotToken: "my-secret-bot-token",
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
// Optimizing your cacheLink logic
//
// > Note: if you create a CacheConfig you don't have to set every field. All the CacheAlgorithms are default to LFU when left blank.
//
// A part of Disgord is the control you have; while this can be a good detail for advanced users, we recommend beginners to utilise the default configurations (by simply not editing the configuration).
// Here we pass the cacheLink config when creating the session to access to the different cacheLink replacement algorithms, lifetime settings, and the option to disable different cacheLink systems.
//  discord, err := disgord.NewSession(&disgord.Config{
//    BotToken: "my-secret-bot-token",
//    Cache: &disgord.CacheConfig{
//              Mutable: false, // everything going in and out of the cacheLink is deep copied
//				// setting Mutable to true, might break your program as this is experimental.
//
//              DisableUserCaching: false, // activates caching for users
//              UserCacheLimitMiB: 500, // don't use more than ~500MiB of memory space for caching of users
//              UserCacheLifetime: time.Duration(4) * time.Hour, // removed from cacheLink after 9 hours, unless updated
//              UserCacheAlgorithm: disgord.CacheAlgLFU,
//
//              DisableVoiceStateCaching: true, // don't cache voice states
//              // VoiceStateCacheLifetime  time.Duration
//              // VoiceStateCacheAlgorithm string
//
//              DisableChannelCaching: false,
//              ChannelCacheLimitMiB: 300,
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
// > Note: Disabling caching for some types while activating it for others (eg. disabling channels, but activating guild caching), can cause items extracted from the cacheLink to not reflect the true discord state.
//
// Example, activated guild but disabled channel caching: The guild is stored to the cacheLink, but it's channels are discarded. Guild channels are dismantled from the guild object and otherwise stored in the channel cacheLink to improve performance and reduce memory use. So when you extract the cached guild object, all of the channel will only hold their channel ID, and nothing more.
//
//
// Immutable and concurrent accessible cacheLink
//
// The option CacheConfig.Immutable can greatly improve performance or break your system. If you utilize channels or you need concurrent access, the safest bet is to set immutable to `true`. While this is slower (as you create deep copies and don't share the same memory space with variables outside the cacheLink), it increases reliability that the cacheLink always reflects the last known Discord state.
// If you are uncertain, just set it to `true`. The default setting is `true` if `disgord.Cache.CacheConfig` is `nil`.
//
//
// Bypass the built-in REST cacheLink
//
// Whenever you call a REST method from the Session interface; the cacheLink is always checked first. Upon a cacheLink hit, no REST request is executed and you get the data from the cacheLink in return. However, if this is problematic for you or there exist a bug which gives you bad/outdated data, you can bypass it by using the REST functions directly. Remember that this will not update the cacheLink for you, and this needs to be done manually if you depend on the cacheLink.
//  // get a user using the Session implementation (checks cacheLink, and updates the cacheLink on cacheLink miss)
//  user, err := session.GetUser(userID)
//
//  // bypass the cacheLink checking. Same function name, but is found in the disgord package, not the session interface.
//  user, err := disgord.GetUser(userID)
//
// Manually updating the cacheLink
//
// If required, you can access the cacheLink and update it by hand. Note that this should not be required when you use the Session interface.
//  user, err := disgord.GetUser(userID)
//  if err != nil {
//      return err
//  }
//
//  // update the cacheLink
//  cacheLink := session.Cache()
//  err = cacheLink.Update(disgord.UserCache, user)
//  if err != nil {
//      return err
//  }
//
// Saving and Deleting Discord data
//
// > Note: when using SaveToDiscord(...) make sure the object reflects the Discord state. Calling Save on default values might overwrite or reset the object at Discord, causing literally.. Hell.
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
	"errors"
	"strings"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/snowflake/v3"
)

// LibraryInfo returns name + version
func LibraryInfo() string {
	return "Disgord " + constant.Version
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

// GetShardForGuildID converts a GuildID into a ShardID for correct retrieval of guild information
func GetShardForGuildID(guildID Snowflake, shardCount uint) (shardID uint) {
	return uint(guildID>>22) % shardCount
}

// https://discordapp.com/developers/docs/resources/user#avatar-data
func validAvatarPrefix(avatar string) (valid bool) {
	if avatar == "" {
		return
	}

	construct := func(encoding string) string {
		return "data:image/" + encoding + ";base64,"
	}

	if len(avatar) < len(construct("X")) {
		return
	}

	encodings := []string{
		"jpeg", "png", "gif",
	}
	for _, encoding := range encodings {
		prefix := construct(encoding)
		if strings.HasPrefix(avatar, prefix) {
			valid = len(avatar)-len(prefix) > 0 // it has content
			break
		}
	}

	return
}

// ValidateUsername uses Discords rule-set to verify user-names and nicknames
// https://discordapp.com/developers/docs/resources/user#usernames-and-nicknames
//
// Note that not all the rules are listed in the docs:
//  There are other rules and restrictions not shared here for the sake of spam and abuse mitigation, but the
//  majority of users won't encounter them. It's important to properly handle all error messages returned by
//  Discord when editing or updating names.
func ValidateUsername(name string) (err error) {
	if name == "" {
		err = errors.New("empty")
		return
	}

	// attributes
	length := len(name)

	// Names must be between 2 and 32 characters long.
	if length < 2 {
		err = errors.New("name is too short")
	} else if length > 32 {
		err = errors.New("name is too long")
	}
	if err != nil {
		return
	}

	// Names are sanitized and trimmed of leading, trailing, and excessive internal whitespace.
	if name[0] == ' ' {
		err = errors.New("contains whitespace prefix")
	} else if name[length-1] == ' ' {
		err = errors.New("contains whitespace suffix")
	} else {
		last := name[1]
		for i := 2; i < length-1; i++ {
			if name[i] == ' ' && last == name[i] {
				err = errors.New("contains excessive internal whitespace")
				break
			}
			last = name[i]
		}
	}
	if err != nil {
		return
	}

	// Names cannot contain the following substrings: '@', '#', ':', '```'
	illegalChars := []string{
		"@", "#", ":", "```",
	}
	for _, illegalChar := range illegalChars {
		if strings.Contains(name, illegalChar) {
			err = errors.New("can not contain the character " + illegalChar)
			return
		}
	}

	// Names cannot be: 'discordtag', 'everyone', 'here'
	illegalNames := []string{
		"discordtag", "everyone", "here",
	}
	for _, illegalName := range illegalNames {
		if name == illegalName {
			err = errors.New("the given username is illegal")
			break
		}
	}

	return
}

func validateChannelName(name string) (err error) {
	if name == "" {
		err = errors.New("empty")
		return
	}

	// attributes
	length := len(name)

	// Names must be of length of minimum 2 and maximum 100 characters long.
	if length < 2 {
		err = errors.New("name is too short")
	} else if length > 100 {
		err = errors.New("name is too long")
	}
	if err != nil {
		return
	}

	return
}
