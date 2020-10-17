package disgord

import (
	"bytes"
	"errors"
	"strconv"
	"time"

	"github.com/andersfylling/disgord/json"
)

// common functionality/types used by struct_*.go files goes here
//go:generate go run generate/interfaces/main.go
//go:generate go run generate/sorters/main.go
//go:generate go run generate/json/main.go

// Copier holds the CopyOverTo method which copies all it's content from one
// struct to another. Note that this requires a deep copy.
// useful when overwriting already existing content in the cacheLink to reduce GC.
type Copier interface {
	CopyOverTo(other interface{}) error
}

// DeepCopier holds the DeepCopy method which creates and returns a deep copy of
// any struct.
type DeepCopier interface {
	DeepCopy() interface{}
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

// Reseter Reset() zero initialises or empties a struct instance
type Reseter interface {
	Reset()
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
