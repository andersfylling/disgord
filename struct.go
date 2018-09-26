package disgord

import (
	"fmt"
	"time"
)

// common functionality/types used by struct_*.go files goes here

// Copier holds the CopyOverTo method which copies all it's content from one
// struct to another. Note that this requires a deep copy.
// useful when overwriting already existing content in the cache to reduce GC.
type Copier interface {
	CopyOverTo(other interface{}) error
}

func NewErrorUnsupportedType(message string) *ErrorUnsupportedType {
	return &ErrorUnsupportedType{
		info: message,
	}
}

type ErrorUnsupportedType struct {
	info string
}

func (eut *ErrorUnsupportedType) Error() string {
	return eut.info
}

// DiscordUpdater holds the Update method for updating any given Discord struct
// (fetch the latest content). If you only want to keep up to date with the
// cache use the UpdateFromCache method.
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
	saveToDiscord(session Session) error
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

// Discord types

// helperTypes: timestamp, levels, etc.

// discordTimeFormat to be able to correctly convert timestamps back into json,
// we need the micro timestamp with an addition at the ending.
// time.RFC3331 does not yield an output similar to the discord timestamp input, the date is however correct.
const timestampFormat = "2006-01-02T15:04:05.000000+00:00"

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalJSON(data []byte) error
}

type Timestamp time.Time

// error: https://stackoverflow.com/questions/28464711/go-strange-json-hyphen-unmarshall-error
func (t Timestamp) MarshalJSON() ([]byte, error) {
	// wrap in double qoutes for valid json parsing
	jsonReady := fmt.Sprintf("\"%s\"", t.String())

	return []byte(jsonReady), nil
}

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

// -----------
// levels

// ExplicitContentFilterLvl ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-explicit-content-filter-level
type ExplicitContentFilterLvl uint

func (ecfl *ExplicitContentFilterLvl) Disabled() bool {
	return *ecfl == 0
}
func (ecfl *ExplicitContentFilterLvl) MembersWithoutRoles() bool {
	return *ecfl == 1
}
func (ecfl *ExplicitContentFilterLvl) AllMembers() bool {
	return *ecfl == 2
}

// MFA ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-mfa-level
type MFALvl uint

func (mfal *MFALvl) None() bool {
	return *mfal == 0
}
func (mfal *MFALvl) Elevated() bool {
	return *mfal == 1
}

// Verification ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-verification-level
type VerificationLvl uint

// None unrestricted
func (vl *VerificationLvl) None() bool {
	return *vl == 0
}

// Low must have verified email on account
func (vl *VerificationLvl) Low() bool {
	return *vl == 1
}

// Medium must be registered on Discord for longer than 5 minutes
func (vl *VerificationLvl) Medium() bool {
	return *vl == 2
}

// High (╯°□°）╯︵ ┻━┻ - must be a member of the server for longer than 10 minutes
func (vl *VerificationLvl) High() bool {
	return *vl == 3
}

// VeryHigh ┻━┻ミヽ(ಠ益ಠ)ﾉ彡┻━┻ - must have a verified phone number
func (vl *VerificationLvl) VeryHigh() bool {
	return *vl == 4
}

// DefaultMessageNotification ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-default-message-notification-level
type DefaultMessageNotificationLvl uint

func (dmnl *DefaultMessageNotificationLvl) AllMessages() bool {
	return *dmnl == 0
}
func (dmnl *DefaultMessageNotificationLvl) OnlyMentions() bool {
	return *dmnl == 1
}
func (dmnl *DefaultMessageNotificationLvl) Equals(v uint) bool {
	return uint(*dmnl) == v
}
