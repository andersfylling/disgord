// Package disgord provides Go bindings for the documented Discord API, and allows for a stateful Client using the Session interface, with the option of a configurable caching system or bypass the built-in caching logic all together.
//
// # Getting started
//
// Create a Disgord client to get access to the REST API and gateway functionality. In the following example, we listen for new messages and respond with "hello".
//
// Session interface: https://pkg.go.dev/github.com/andersfylling/disgord?tab=doc#Session
//
//	client := disgord.New(disgord.Config{
//	  BotToken: "my-secret-bot-token",
//	})
//	defer client.Gateway().StayConnectedUntilInterrupted()
//
//	client.Gateway().MessageCreate(func(s disgord.Session, evt *disgord.MessageCreate) {
//	  evt.Message.Reply(context.Background(), s, "hello")
//	})
//
// # Listen for events using Channels
//
// You don't have to use a callback function, channels are supported too!
//
//	msgChan := make(chan *disgord.MessageCreate, 10)
//	client.Gateway().MessageCreateChan(msgChan)
//
// Never close a channel without removing the handler from Disgord, as it will cause a panic. You can control the
// lifetime of a handler or injected channel by in injecting a controller: disgord.HandlerCtrl. Since you are the
// owner of the channel, disgord will never close it for you.
//
//	ctrl := &disgord.Ctrl{Channel: msgCreateChan}
//	client.Gateway().WithCtrl(ctrl).MessageCreateChan(msgChan)
//	go func() {
//	  // close the channel after 20 seconds and safely remove it from the Disgord reactor
//	  <- time.After(20 * time.Second)
//	  ctrl.CloseChannel()
//	}
//
// # WebSockets and Sharding
//
// Disgord handles sharding for you automatically; when starting the bot, when discord demands you to scale up your shards (during runtime), etc. It also gives you control over the shard setup in case you want to run multiple instances of Disgord (in these cases you must handle scaling yourself as Disgord can not).
//
// Sharding is done behind the scenes, so you do not need to worry about any settings. Disgord will simply ask Discord for the recommended amount of shards for your bot on startup. However, to set specific amount of shards you can use the `disgord.ShardConfig` to specify a range of valid shard IDs (starts from 0).
//
// starting a bot with exactly 5 shards
//
//	client := disgord.New(disgord.Config{
//	  ShardConfig: disgord.ShardConfig{
//	    // this is a copy so u can't manipulate the config later on
//	    ShardIDs: []uint{0,1,2,3,4},
//	  },
//	})
//
// Running multiple instances each with 1 shard (note each instance must use unique shard ids)
//
//	client := disgord.New(disgord.Config{
//	  ShardConfig: disgord.ShardConfig{
//	    // this is a copy so u can't manipulate the config later on
//	    ShardIDs: []uint{0}, // this number must change for each instance. Try to automate this.
//	    ShardCount: 5, // total of 5 shards, but this disgord instance only has one. AutoScaling is disabled - use OnScalingRequired.
//	  },
//	})
//
// Handle scaling options yourself
//
//	client := disgord.New(disgord.Config{
//	  ShardConfig: disgord.ShardConfig{
//	    // this is a copy so u can't manipulate it later on
//	    DisableAutoScaling: true,
//	    OnScalingRequired: func(shardIDs []uint) (TotalNrOfShards uint, AdditionalShardIDs []uint) {
//	      // instead of asking discord for exact number of shards recommended
//	      // this is increased by 50% every time discord complains you don't have enough shards
//	      // to reduce the number of times you have to scale
//	      TotalNrOfShards := uint(len(shardIDs) * 1.5)
//	      for i := len(shardIDs) - 1; i < TotalNrOfShards; i++ {
//	        AdditionalShardIDs = append(AdditionalShardIDs, i)
//	      }
//	      return
//	    }, // end OnScalingRequired
//	  }, // end ShardConfig
//	})
//
// # Caching
//
// You can inject your own cache implementation. By default a read only LFU implementation is used, this should be
// sufficient for the average user. But you can overwrite certain methods as well!
//
// Say you dislike the implementation for MESSAGE_CREATE events, you can embed the default cache and define your own
// logic:
//
//	 type MyCoolCache struct {
//	   disgord.BasicCache
//	   msgCache map[Snowflake]*Message // channelID => Message
//	 }
//	 func (c *BasicCache) MessageCreate(data []byte) (*MessageCreate, error) {
//		  // some smart implementation here
//	 }
//
// > Note: if you inject your own cache, remember that the cache is also responsible for initiating the objects.
// > See disgord.CacheNop
//
// # Bypass the built-in REST cache
//
// Whenever you call a REST method from the Session interface; the cache is always checked first. Upon a cache hit, no REST request is executed and you get the data from the cache in return. However, if this is problematic for you or there exist a bug which gives you bad/outdated data, you can bypass it by using Disgord flags.
//
//	// get a user using the Session implementation (checks cache, and updates the cache on cache miss)
//	user, err := client.User(userID).Get()
//
//	// bypass the cache checking. Same as before, but we insert a disgord.Flag type.
//	user, err := client.User(userID).Get(disgord.IgnoreCache)
//
// # Disgord Flags
//
// In addition to disgord.IgnoreCache, as shown above, you can pass in other flags such as: disgord.SortByID, disgord.OrderAscending, etc. You can find these flags in the flag.go file.
//
// # Build tags
//
// `disgord_diagnosews` will store all the incoming and outgoing JSON data as files in the directory "diagnose-report/packets". The file format is as follows: unix_clientType_direction_shardID_operationCode_sequenceNumber[_eventName].json
package disgord

//go:generate go run internal/generate/intents/main.go
//go:generate go run internal/generate/interfaces/main.go
//go:generate go run internal/generate/inter/main.go
//go:generate go run internal/generate/sorters/main.go
//go:generate go run internal/generate/querybuilders/main.go

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/andersfylling/disgord/internal/gateway"
	"github.com/andersfylling/disgord/json"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

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

func newErrorUnsupportedType(message string) *ErrorUnsupportedType {
	return &ErrorUnsupportedType{
		info: message,
	}
}

// ErrorUnsupportedType used when the given param type is not supported
type ErrorUnsupportedType struct {
	info string
}

func (e *ErrorUnsupportedType) Error() string {
	return e.info
}

// hasher creates a hash for comparing objects. This excludes the identifier and object type as those are expected
// to be the same during a comparison.
type hasher interface {
	hash() string
}

type guilder interface {
	getGuildIDs() []Snowflake
}

// Mentioner can be implemented by any type that is mentionable.
// https://discord.com/developers/docs/reference#message-formatting-formats
type Mentioner interface {
	Mention() string
}

// zeroInitialiser zero initializes a struct by setting all the values to the default initialization values.
// Used in the flyweight pattern.
type zeroInitialiser interface {
	zeroInitialize()
}

// internalUpdater is called whenever a socket event or a REST response is created.
type internalUpdater interface {
	updateInternals()
}

type internalClientUpdater interface {
	updateInternalsWithClient(*Client)
}

// Discord types

// helperTypes: timestamp, levels, etc.

// discordTimeFormat to be able to correctly convert timestamps back into json,
// we need the micro timestamp with an addition at the ending.
// time.RFC3331 does not yield an output similar to the discord timestamp input, the date is however correct.
const timestampFormat = "2006-01-02T15:04:05.000000+00:00"

// Time handles Discord timestamps
type Time struct {
	time.Time
}

var _ json.Marshaler = (*Time)(nil)
var _ json.Unmarshaler = (*Time)(nil)

// MarshalJSON implements json.Marshaler.
// error: https://stackoverflow.com/questions/28464711/go-strange-json-hyphen-unmarshall-error
func (t Time) MarshalJSON() ([]byte, error) {
	var ts string
	if !t.IsZero() {
		ts = t.String()
	}

	// wrap in double quotes for valid json parsing
	return []byte(`"` + ts + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *Time) UnmarshalJSON(data []byte) error {
	var ts time.Time

	// Don't try to unmarshal empty strings.
	if bytes.Equal([]byte("\"\""), data) {
		return nil
	}

	if err := json.Unmarshal(data, &ts); err != nil {
		return err
	}

	t.Time = ts
	return nil
}

// String returns the timestamp as a Discord formatted timestamp. Formatting
// with time.RFC3331 does not suffice.
func (t Time) String() string {
	return t.Format(timestampFormat)
}

// -----------
// levels

// ExplicitContentFilterLvl ...
// https://discord.com/developers/docs/resources/guild#guild-object-explicit-content-filter-level
type ExplicitContentFilterLvl uint

// Explicit content filter levels
const (
	ExplicitContentFilterLvlDisabled ExplicitContentFilterLvl = iota
	ExplicitContentFilterLvlMembersWithoutRoles
	ExplicitContentFilterLvlAllMembers
)

// Disabled if the content filter is disabled
func (ecfl *ExplicitContentFilterLvl) Disabled() bool {
	return *ecfl == ExplicitContentFilterLvlDisabled
}

// MembersWithoutRoles if the filter only applies for members without a role
func (ecfl *ExplicitContentFilterLvl) MembersWithoutRoles() bool {
	return *ecfl == ExplicitContentFilterLvlMembersWithoutRoles
}

// AllMembers if the filter applies for all members regardles of them having a role or not
func (ecfl *ExplicitContentFilterLvl) AllMembers() bool {
	return *ecfl == ExplicitContentFilterLvlAllMembers
}

// MFALvl ...
// https://discord.com/developers/docs/resources/guild#guild-object-mfa-level
type MFALvl uint

// Different MFA levels
const (
	MFALvlNone MFALvl = iota
	MFALvlElevated
)

// None ...
func (mfal *MFALvl) None() bool {
	return *mfal == MFALvlNone
}

// Elevated ...
func (mfal *MFALvl) Elevated() bool {
	return *mfal == MFALvlElevated
}

// VerificationLvl ...
// https://discord.com/developers/docs/resources/guild#guild-object-verification-level
type VerificationLvl uint

// the different verification levels
const (
	VerificationLvlNone VerificationLvl = iota
	VerificationLvlLow
	VerificationLvlMedium
	VerificationLvlHigh
	VerificationLvlVeryHigh
)

// None unrestricted
func (vl *VerificationLvl) None() bool {
	return *vl == VerificationLvlNone
}

// Low must have verified email on account
func (vl *VerificationLvl) Low() bool {
	return *vl == VerificationLvlLow
}

// Medium must be registered on Discord for longer than 5 minutes
func (vl *VerificationLvl) Medium() bool {
	return *vl == VerificationLvlMedium
}

// High (╯°□°）╯︵ ┻━┻ - must be a member of the server for longer than 10 minutes
func (vl *VerificationLvl) High() bool {
	return *vl == VerificationLvlHigh
}

// VeryHigh ┻━┻ミヽ(ಠ益ಠ)ﾉ彡┻━┻ - must have a verified phone number
func (vl *VerificationLvl) VeryHigh() bool {
	return *vl == VerificationLvlVeryHigh
}

// GuildScheduledEventPrivacyLevel ...
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-privacy-level
type GuildScheduledEventPrivacyLevel uint

// the different scheduled event privacy level
const (
	GuildScheduledEventPrivacyLevelGuildOnly GuildScheduledEventPrivacyLevel = iota + 2
)

// GuildScheduledEventEntityTypes ...
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-types
type GuildScheduledEventEntityTypes uint

// the different scheduled event entity types
const (
	GuildScheduledEventEntityTypesStageInstance GuildScheduledEventEntityTypes = iota + 1
	GuildScheduledEventEntityTypesVoice
	GuildScheduledEventEntityTypesExternal
)

type GuildScheduledEventStatus uint

const (
	GuildScheduledEventStatusScheduled GuildScheduledEventStatus = iota + 1
	GuildScheduledEventStatusActive
	GuildScheduledEventStatusCompleted
	GuildScheduledEventStatusCancelled
)

// PremiumTier ...
// https://discord.com/developers/docs/resources/guild#guild-object-premium-tier
type PremiumTier uint

// the different premium tier levels
const (
	PremiumTierNone PremiumTier = iota
	PremiumTier1
	PremiumTier2
	PremiumTier3
)

// DefaultMessageNotificationLvl ...
// https://discord.com/developers/docs/resources/guild#guild-object-default-message-notification-level
type DefaultMessageNotificationLvl uint

// different notification levels on new messages
const (
	DefaultMessageNotificationLvlAllMessages DefaultMessageNotificationLvl = iota
	DefaultMessageNotificationLvlOnlyMentions
)

// AllMessages ...
func (dmnl *DefaultMessageNotificationLvl) AllMessages() bool {
	return *dmnl == DefaultMessageNotificationLvlAllMessages
}

// OnlyMentions ...
func (dmnl *DefaultMessageNotificationLvl) OnlyMentions() bool {
	return *dmnl == DefaultMessageNotificationLvlOnlyMentions
}

// NewDiscriminator Discord user discriminator hashtag
func NewDiscriminator(d string) (discriminator Discriminator, err error) {
	var tmp uint64
	tmp, err = strconv.ParseUint(d, 10, 16)
	if err == nil {
		discriminator = Discriminator(tmp)
	}

	return
}

// Discriminator value
type Discriminator uint16

var _ json.Unmarshaler = (*Discriminator)(nil)
var _ json.Marshaler = (*Discriminator)(nil)

func (d Discriminator) String() (str string) {
	if d == 0 {
		str = ""
		return
	}
	if d == 1 {
		str = "0001"
		return
	}

	str = strconv.Itoa(int(d))
	if d < 1000 {
		shift := 4 - len(str)
		for i := 0; i < shift; i++ {
			str = "0" + str
		}
	}

	return
}

// NotSet checks if the discriminator is not set
func (d Discriminator) NotSet() bool {
	return d == 0
}

// UnmarshalJSON see interface json.Unmarshaler
func (d *Discriminator) UnmarshalJSON(data []byte) error {
	*d = 0
	length := len(data) - 1
	for i := 1; i < length; i++ {
		*d = *d*10 + Discriminator(data[i]-'0')
	}
	return nil
}

// MarshalJSON see interface json.Marshaler
func (d Discriminator) MarshalJSON() (data []byte, err error) {
	return []byte("\"" + d.String() + "\""), nil
}

// ShardID calculate the shard id for a given guild.
// https://discord.com/developers/docs/topics/gateway#sharding-sharding-formula
func ShardID(guildID Snowflake, nrOfShards uint) uint {
	return gateway.GetShardForGuildID(guildID, nrOfShards)
}

//////////////////////////////////////////////////////
//
// Validators
//
//////////////////////////////////////////////////////

// https://discord.com/developers/docs/resources/user#avatar-data
func validAvatarPrefix(avatar string) (valid bool) {
	if avatar == "" {
		return false
	}

	construct := func(encoding string) string {
		return "data:image/" + encoding + ";base64,"
	}

	if len(avatar) < len(construct("X")) {
		return false // missing base64 declaration
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

	return true
}

// ValidateUsername uses Discords rule-set to verify user-names and nicknames
// https://discord.com/developers/docs/resources/user#usernames-and-nicknames
//
// Note that not all the rules are listed in the docs:
//
//	There are other rules and restrictions not shared here for the sake of spam and abuse mitigation, but the
//	majority of Users won't encounter them. It's important to properly handle all error messages returned by
//	Discord when editing or updating names.
func ValidateUsername(name string) (err error) {
	if name == "" {
		return errors.New("empty")
	}

	// attributes
	length := len(name)

	// Names must be between 2 and 32 characters long.
	if length < 2 {
		err = fmt.Errorf("name is too short: %w", ErrIllegalValue)
	} else if length > 32 {
		err = fmt.Errorf("name is too long: %w", ErrIllegalValue)
	}
	if err != nil {
		return err
	}

	// Names are sanitized and trimmed of leading, trailing, and excessive internal whitespace.
	if name[0] == ' ' {
		err = fmt.Errorf("contains whitespace prefix: %w", ErrIllegalValue)
	} else if name[length-1] == ' ' {
		err = fmt.Errorf("contains whitespace suffix: %w", ErrIllegalValue)
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
		return err
	}

	// Names cannot contain the following substrings: '@', '#', ':', '```'
	illegalChars := []string{
		"@", "#", ":", "```",
	}
	for _, illegalChar := range illegalChars {
		if strings.Contains(name, illegalChar) {
			err = errors.New("can not contain the character " + illegalChar)
			return err
		}
	}

	// Names cannot be: 'discordtag', 'everyone', 'here'
	illegalNames := []string{
		"discordtag", "everyone", "here",
	}
	for _, illegalName := range illegalNames {
		if name == illegalName {
			err = fmt.Errorf("the given username is illegal: %w", ErrIllegalValue)
			return err
		}
	}

	return nil
}

func validateChannelName(name string) (err error) {
	if name == "" {
		return ErrMissingChannelName
	}

	// attributes
	length := len(name)

	// Names must be of length of minimum 2 and maximum 100 characters long.
	if length < 2 {
		err = fmt.Errorf("name is too short: %w", ErrIllegalValue)
	} else if length > 100 {
		err = fmt.Errorf("name is too long: %w", ErrIllegalValue)
	}
	if err != nil {
		return err
	}

	return nil
}

// CreateTermSigListener create a channel to listen for termination signals (graceful shutdown)
func CreateTermSigListener() <-chan os.Signal {
	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	return termSignal
}
