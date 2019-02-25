package disgord

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/andersfylling/snowflake/v3"

	"github.com/andersfylling/disgord/constant"
)

// common functionality/types used by struct_*.go files goes here
//go:generate go run generate/interfaces/main.go

// Copier holds the CopyOverTo method which copies all it's content from one
// struct to another. Note that this requires a deep copy.
// useful when overwriting already existing content in the cacheLink to reduce GC.
type Copier interface {
	CopyOverTo(other interface{}) error
}

// cacheCopier is similar to Copier interface. Except that it only copies over fields which has a value, unlike Copier
// that creates an exact copy of everything. This will also ignore arrays that can be simplified to a snowflake array.
// An example of said simplification is Guild.Channels, as there will already exist a channel cacheLink.
//
// It is important to know that this should only be called by the cacheLink. The cacheLink must also make sure that the type
// given as an argument for `other` is correct. Failure to do so results in a panic.
type cacheCopier interface {
	copyOverToCache(other interface{}) error
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

// DiscordUpdater holds the Update method for updating any given Discord struct
// (fetch the latest content). If you only want to keep up to date with the
// cacheLink use the UpdateFromCache method.
// TODO: change param type for UpdateFromCache once caching is implemented
//type DiscordUpdater interface {
//	Update(session Session)
//	UpdateFromCache(session Session)
//}

// DiscordSaver holds the SaveToDiscord method for sending changes to the
// Discord API over REST.
// If you change any of the values and want to notify Discord about your change,
// use the Save method to send a REST request (assuming that the struct values
// can be updated).
//
// NOTE! if the struct has an snowflake/ID, it will update content. But if the
// snowflake is missing/not set, it will create content (if possible,
// otherwise you will get an error)
type discordSaver interface {
	saveToDiscord(session Session, changes discordSaver) error
}

// DiscordDeleter holds the DeleteFromDiscord method which deletes a given
// object from the Discord servers.
type discordDeleter interface {
	deleteFromDiscord(session Session) error
}

// DeepCopier holds the DeepCopy method which creates and returns a deep copy of
// any struct.
type DeepCopier interface {
	DeepCopy() interface{}
}

// hasher creates a hash for comparing objects. This excludes the identifier and object type as those are expected
// to be the same during a comparison.
type hasher interface {
	hash() string
}

type guilder interface {
	getGuildID() snowflake.ID
}

// zeroInitialiser zero initializes a struct by setting all the values to the default initialization values.
// Used in the flyweight pattern.
type zeroInitialiser interface {
	zeroInitialize()
}

// Reseter Reset() zero initialises or empties a struct instance
type Reseter interface {
	Reset()
}

// internalUpdater is called whenever a socket event or a REST response is created.
type internalUpdater interface {
	updateInternals()
}

type internalClientUpdater interface {
	updateInternalsWithClient(*client)
}

// Discord types

// helperTypes: timestamp, levels, etc.

// discordTimeFormat to be able to correctly convert timestamps back into json,
// we need the micro timestamp with an addition at the ending.
// time.RFC3331 does not yield an output similar to the discord timestamp input, the date is however correct.
const timestampFormat = "2006-01-02T15:04:05.000000+00:00"

// Timestamp handles Discord timestamps
type Timestamp time.Time

// MarshalJSON see json.Marshaler
// error: https://stackoverflow.com/questions/28464711/go-strange-json-hyphen-unmarshall-error
func (t Timestamp) MarshalJSON() ([]byte, error) {
	// wrap in double qoutes for valid json parsing
	jsonReady := fmt.Sprintf("\"%s\"", t.String())

	return []byte(jsonReady), nil
}

// UnmarshalJSON see json.Unmarshaler
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var ts time.Time
	err := unmarshal(data, &ts)
	if err != nil {
		return err
	}

	*t = Timestamp(ts)
	return nil
}

// String converts the timestamp into a discord formatted timestamp. time.RFC3331 does not suffice
func (t Timestamp) String() string {
	return t.Time().Format(timestampFormat)
}

// Time converts the DiscordTimestamp into a time.Time type.......
func (t Timestamp) Time() time.Time {
	return time.Time(t)
}

// Empty check if the timestamp holds no value / not set
func (t Timestamp) Empty() bool {
	return time.Time(t).UnixNano() == 0
}

// -----------
// levels

// ExplicitContentFilterLvl ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-explicit-content-filter-level
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
// https://discordapp.com/developers/docs/resources/guild#guild-object-mfa-level
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
// https://discordapp.com/developers/docs/resources/guild#guild-object-verification-level
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

// DefaultMessageNotificationLvl ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-default-message-notification-level
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
func (d *Discriminator) UnmarshalJSON(data []byte) (err error) {
	*d = 0
	length := len(data) - 1
	for i := 1; i < length; i++ {
		*d = *d*10 + Discriminator(data[i]-'0')
	}
	return
}

// MarshalJSON see interface json.Marshaler
func (d Discriminator) MarshalJSON() (data []byte, err error) {
	return []byte("\"" + d.String() + "\""), nil
}

// Gateway is for parsing the Gateway endpoint response
type Gateway struct {
	URL string `json:"url"`
}

// GatewayBot is for parsing the Gateway Bot endpoint response
type GatewayBot struct {
	Gateway
	Shards            uint `json:"shards"`
	SessionStartLimit struct {
		Total      uint `json:"total"`
		Remaining  uint `json:"remaining"`
		ResetAfter uint `json:"reset_after"`
	} `json:"session_start_limit"`
}

// extractAttribute extracts the snowflake value from a JSON string given a attribute filter. For extracting the root ID of an JSON byte array,
// set filter to `"id":"` and scope to `0`. Note that the filter holds the last character before the value starts.
func extractAttribute(filter []byte, scope int, data []byte) (id Snowflake, err error) {
	//filter := []byte(`"id":"`)
	filterLen := len(filter) - 1
	//scope := 0

	var start uint
	lastPos := len(data) - 1
	for i := 1; i <= lastPos-filterLen; i++ {
		if data[i] == '{' {
			scope++
		} else if data[i] == '}' {
			scope--
		}

		if scope != 0 {
			continue
		}

		for j := filterLen; j >= 0; j-- {
			if filter[j] != data[i+j] {
				break
			}

			if j == 0 {
				start = uint(i + len(filter))
			}
		}

		if start != 0 {
			break
		}
	}

	if start == 0 {
		err = errors.New("unable to locate ID")
		return
	}

	i := start
	//E:
	for {
		if data[i] >= '0' && data[i] <= '9' {
			i++
		} else {
			break
		}
		//
		//switch data[i] {
		//case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		//	i++
		//default:
		//	break E
		//}
	}

	if i > start {
		id = Snowflake(0)
		err = id.UnmarshalJSON(data[start-1 : i+1])
	} else {
		err = errors.New("id was empty")
	}
	return
}

func handleRWLocking(read, write *sync.RWMutex) {
	if constant.LockedMethods {
		read.RLock()
		write.Lock()
	}
}

func handleRWUnlocking(read, write *sync.RWMutex) {
	if constant.LockedMethods {
		read.RUnlock()
		write.Unlock()
	}
}
