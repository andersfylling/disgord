package disgord

import (
	"encoding/json"
	"errors"
	"github.com/andersfylling/disgord/constant"
)

const (
	// StatusIdle presence status for idle
	StatusIdle = "idle"
	// StatusDnd presence status for dnd
	StatusDnd = "dnd"
	// StatusOnline presence status for online
	StatusOnline = "online"
	// StatusOffline presence status for offline
	StatusOffline = "offline"
)

// flags for the Activity object to signify the type of action taken place
const (
	ActivityFlagInstance    = 1 << 0
	ActivityFlagJoin        = 1 << 1
	ActivityFlagSpectate    = 1 << 2
	ActivityFlagJoinRequest = 1 << 3
	ActivityFlagSync        = 1 << 4
	ActivityFlagPlay        = 1 << 5
)

//type UserInterface interface {
//	Mention() string
//	MentionNickname() string
//	String() string
//}

// ActivityParty ...
type ActivityParty struct {
	Lockable `json:"-"`

	ID   string `json:"id,omitempty"`   // the id of the party
	Size []int  `json:"size,omitempty"` // used to show the party's current and maximum size
}

// Limit shows the maximum number of guests/people allowed
func (ap *ActivityParty) Limit() int {
	if len(ap.Size) != 2 {
		return 0
	}

	return ap.Size[1]
}

// NumberOfPeople shows the current number of people attending the Party
func (ap *ActivityParty) NumberOfPeople() int {
	if len(ap.Size) != 1 {
		return 0
	}

	return ap.Size[0]
}

// DeepCopy see interface at struct.go#DeepCopier
func (ap *ActivityParty) DeepCopy() (copy interface{}) {
	copy = &ActivityParty{}
	ap.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (ap *ActivityParty) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var activity *ActivityParty
	if activity, ok = other.(*ActivityParty); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *ActivityParty")
		return
	}

	if constant.LockedMethods {
		ap.RLock()
		activity.Lock()
	}

	activity.ID = ap.ID
	activity.Size = ap.Size

	if constant.LockedMethods {
		ap.RUnlock()
		activity.Unlock()
	}

	return
}

// ActivityAssets ...
type ActivityAssets struct {
	Lockable `json:"-"`

	LargeImage string `json:"large_image,omitempty"` // the id for a large asset of the activity, usually a snowflake
	LargeText  string `json:"large_text,omitempty"`  //text displayed when hovering over the large image of the activity
	SmallImage string `json:"small_image,omitempty"` // the id for a small asset of the activity, usually a snowflake
	SmallText  string `json:"small_text,omitempty"`  //	text displayed when hovering over the small image of the activity
}

// DeepCopy see interface at struct.go#DeepCopier
func (a *ActivityAssets) DeepCopy() (copy interface{}) {
	copy = &ActivityAssets{}
	a.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (a *ActivityAssets) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var activity *ActivityAssets
	if activity, ok = other.(*ActivityAssets); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *ActivityAssets")
		return
	}

	if constant.LockedMethods {
		a.RLock()
		activity.Lock()
	}

	activity.LargeImage = a.LargeImage
	activity.LargeText = a.LargeText
	activity.SmallImage = a.SmallImage
	activity.SmallText = a.SmallText

	if constant.LockedMethods {
		a.RUnlock()
		activity.Unlock()
	}

	return
}

// ActivitySecrets ...
type ActivitySecrets struct {
	Lockable `json:"-"`

	Join     string `json:"join,omitempty"`     // the secret for joining a party
	Spectate string `json:"spectate,omitempty"` // the secret for spectating a game
	Match    string `json:"match,omitempty"`    // the secret for a specific instanced match
}

// DeepCopy see interface at struct.go#DeepCopier
func (a *ActivitySecrets) DeepCopy() (copy interface{}) {
	copy = &ActivitySecrets{}
	a.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (a *ActivitySecrets) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var activity *ActivitySecrets
	if activity, ok = other.(*ActivitySecrets); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *ActivitySecrets")
		return
	}

	if constant.LockedMethods {
		a.RLock()
		activity.Lock()
	}

	activity.Join = a.Join
	activity.Spectate = a.Spectate
	activity.Match = a.Match

	if constant.LockedMethods {
		a.RUnlock()
		activity.Unlock()
	}

	return
}

// ActivityTimestamp ...
type ActivityTimestamp struct {
	Lockable `json:"-"`

	Start int `json:"start,omitempty"` // unix time (in milliseconds) of when the activity started
	End   int `json:"end,omitempty"`   // unix time (in milliseconds) of when the activity ends
}

// DeepCopy see interface at struct.go#DeepCopier
func (a *ActivityTimestamp) DeepCopy() (copy interface{}) {
	copy = &ActivityTimestamp{}
	a.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (a *ActivityTimestamp) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var activity *ActivityTimestamp
	if activity, ok = other.(*ActivityTimestamp); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *ActivityTimestamp")
		return
	}

	if constant.LockedMethods {
		a.RLock()
		activity.Lock()
	}

	activity.Start = a.Start
	activity.End = a.End

	if constant.LockedMethods {
		a.RUnlock()
		activity.Unlock()
	}

	return
}

// NewActivity ...
func NewActivity() (activity *Activity) {
	return &Activity{
		Timestamps: []*ActivityTimestamp{},
	}
}

// Activity https://discordapp.com/developers/docs/topics/gateway#activity-object-activity-structure
type Activity struct {
	Lockable `json:"-"`

	Name          string               `json:"name"`                     // the activity's name
	Type          int                  `json:"type"`                     // activity type
	URL           *string              `json:"url,omitempty"`            //stream url, is validated when type is 1
	Timestamps    []*ActivityTimestamp `json:"timestamps,omitempty"`     // timestamps object	unix timestamps for start and/or end of the game
	ApplicationID Snowflake            `json:"application_id,omitempty"` //?	snowflake	application id for the game
	Details       *string              `json:"details,omitempty"`        //?	?string	what the player is currently doing
	State         *string              `json:"state,omitempty"`          //state?	?string	the user's current party status
	Party         *ActivityParty       `json:"party"`                    //party?	party object	information for the current party of the player
	Assets        *ActivityAssets      `json:"assets,omitempty"`         // assets?	assets object	images for the presence and their hover texts
	Secrets       *ActivitySecrets     `json:"secrets,omitempty"`        // secrets?	secrets object	secrets for Rich Presence joining and spectating
	Instance      bool                 `json:"instance,omitempty"`       // instance?	boolean	whether or not the activity is an instanced game session
	Flags         int                  `json:"flags,omitempty"`          // flags?	int	activity flags ORd together, describes what the payload includes
}

// DeepCopy see interface at struct.go#DeepCopier
func (a *Activity) DeepCopy() (copy interface{}) {
	copy = &Activity{}
	a.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (a *Activity) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var activity *Activity
	if activity, ok = other.(*Activity); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *Activity")
		return
	}

	if constant.LockedMethods {
		a.RLock()
		activity.Lock()
	}

	activity.Name = a.Name
	activity.Type = a.Type
	activity.ApplicationID = a.ApplicationID
	activity.Instance = a.Instance
	activity.Flags = a.Flags

	if a.URL != nil {
		url := *a.URL
		activity.URL = &url
	}
	if a.Timestamps != nil {
		if activity.Timestamps == nil {
			activity.Timestamps = make([]*ActivityTimestamp, len(a.Timestamps))
		}
		for i, timestampP := range a.Timestamps {
			if timestampP == nil {
				continue
			}
			activity.Timestamps[i] = timestampP.DeepCopy().(*ActivityTimestamp)
		}
	}
	if a.Details != nil {
		details := *a.Details
		activity.Details = &details
	}
	if a.State != nil {
		state := *a.State
		activity.State = &state
	}
	if a.Party != nil {
		activity.Party = a.Party.DeepCopy().(*ActivityParty)
	}
	if a.Assets != nil {
		activity.Assets = a.Assets.DeepCopy().(*ActivityAssets)
	}
	if a.Secrets != nil {
		activity.Secrets = a.Secrets.DeepCopy().(*ActivitySecrets)
	}

	if constant.LockedMethods {
		a.RUnlock()
		activity.Unlock()
	}

	return
}

// ---------

const (
	userOEmail      = 0x1 << iota
	userOAvatar     = 0x1 << iota
	userOToken      = 0x1 << iota
	userOVerified   = 0x1 << iota
	userOMFAEnabled = 0x1 << iota
	userOBot        = 0x1 << iota
)

// NewUser creates a new, empty user object
func NewUser() *User {
	return &User{}
}

func newUserJSON() *userJSON {
	d := "-"
	return &userJSON{
		Avatar: &d,
	}
}

type userJSON struct {
	/*-*/ ID Snowflake `json:"id,omitempty"`
	/*-*/ Username string `json:"username,omitempty"`
	/*-*/ Discriminator Discriminator `json:"discriminator,omitempty"`
	/*1*/ Email *string `json:"email"`
	/*2*/ Avatar *string `json:"avatar"`
	/*3*/ Token *string `json:"token"`
	/*4*/ Verified *bool `json:"verified"`
	/*5*/ MFAEnabled *bool `json:"mfa_enabled"`
	/*6*/ Bot *bool `json:"bot"`
}

func (u *userJSON) extractMap() uint8 {
	var overwritten uint8
	if u.Email != nil {
		overwritten |= userOEmail
	}
	if u.Avatar == nil || *u.Avatar != "-" {
		overwritten |= userOAvatar
	}
	if u.Token != nil {
		overwritten |= userOToken
	}
	if u.Verified != nil {
		overwritten |= userOVerified
	}
	if u.MFAEnabled != nil {
		overwritten |= userOMFAEnabled
	}
	if u.Bot != nil {
		overwritten |= userOBot
	}

	return overwritten
}

// User the Discord user object which is reused in most other data structures.
type User struct {
	Lockable `json:"-"`

	ID            Snowflake     `json:"id,omitempty"`
	Username      string        `json:"username,omitempty"`
	Discriminator Discriminator `json:"discriminator,omitempty"`
	Email         string        `json:"email,omitempty"`
	Avatar        *string       `json:"avatar"` // data:image/jpeg;base64,BASE64_ENCODED_JPEG_IMAGE_DATA //TODO: pointer?
	Token         string        `json:"token,omitempty"`
	Verified      bool          `json:"verified,omitempty"`
	MFAEnabled    bool          `json:"mfa_enabled,omitempty"`
	Bot           bool          `json:"bot,omitempty"`

	// Used to identify which fields are set by Discord in partial JSON objects. Yep.
	overwritten uint8 // map. see number left of field in userJSON struct.
}

// Mention returns the a string that Discord clients can format into a valid Discord mention
func (u *User) Mention() string {
	return "<@" + u.ID.String() + ">"
}

// MentionNickname same as Mention, but shows nicknames
// TODO: move to member object(?)
func (u *User) MentionNickname() string {
	return "<@!" + u.ID.String() + ">"
}

func (u *User) String() string {
	return u.Username + "#" + u.Discriminator.String() + "{" + u.ID.String() + "}"
}

// MarshalJSON see interface json.Marshaler
func (u *User) MarshalJSON() ([]byte, error) {
	if u.ID.Empty() {
		return []byte("{}"), nil
	}

	return json.Marshal(User(*u))
}

// UnmarshalJSON see interface json.Unmarshaler
func (u *User) UnmarshalJSON(data []byte) (err error) {
	j := userJSON{}
	err = json.Unmarshal(data, &j)
	if err != nil {
		return
	}

	changes := j.extractMap()
	u.ID = j.ID
	if j.Username != "" {
		u.Username = j.Username
	}
	if j.Discriminator != 0 {
		u.Discriminator = j.Discriminator
	}
	if (changes & userOEmail) > 0 {
		u.Email = *j.Email
	}
	if (changes & userOAvatar) > 0 {
		u.Avatar = j.Avatar
	}
	if (changes & userOToken) > 0 {
		u.Token = *j.Token
	}
	if (changes & userOVerified) > 0 {
		u.Verified = *j.Verified
	}
	if (changes & userOMFAEnabled) > 0 {
		u.MFAEnabled = *j.MFAEnabled
	}
	if (changes & userOBot) > 0 {
		u.Bot = *j.Bot
	}
	u.overwritten |= changes

	return
}

// SendMsg send a message to a user where you utilize a Message object instead of a string
func (u *User) SendMsg(session Session, message *Message) (channel *Channel, msg *Message, err error) {
	channel, err = session.CreateDM(u.ID)
	if err != nil {
		return
	}

	msg, err = session.SendMsg(channel.ID, message)
	return
}

// SendMsgString send a message to given user where the message is in the form of a string.
func (u *User) SendMsgString(session Session, content string) (channel *Channel, msg *Message, err error) {
	channel, msg, err = u.SendMsg(session, &Message{
		Content: content,
	})
	return
}

// DeepCopy see interface at struct.go#DeepCopier
// CopyOverTo see interface at struct.go#Copier
func (u *User) DeepCopy() (copy interface{}) {
	copy = NewUser()
	u.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (u *User) CopyOverTo(other interface{}) (err error) {
	var user *User
	var valid bool
	if user, valid = other.(*User); !valid {
		err = newErrorUnsupportedType("argument given is not a *User type")
		return
	}

	if constant.LockedMethods {
		u.RLock()
		user.Lock()
	}

	user.ID = u.ID
	user.Username = u.Username
	user.Discriminator = u.Discriminator
	user.Email = u.Email
	user.Token = u.Token
	user.Verified = u.Verified
	user.MFAEnabled = u.MFAEnabled
	user.Bot = u.Bot
	user.overwritten = u.overwritten

	if u.Avatar != nil {
		avatar := *u.Avatar
		user.Avatar = &avatar
	}

	if constant.LockedMethods {
		u.RUnlock()
		user.Unlock()
	}

	return
}

// copyOverToCache see interface at struct.go#CacheCopier
func (u *User) copyOverToCache(other interface{}) (err error) {
	user := other.(*User)

	if constant.LockedMethods {
		u.RLock()
		user.Lock()
	}

	if !u.ID.Empty() {
		user.ID = u.ID
	}
	if u.Username != "" {
		user.Username = u.Username
	}
	if u.Discriminator != 0 {
		user.Discriminator = u.Discriminator
	}
	if (u.overwritten & userOEmail) > 0 {
		user.Email = u.Email
	}
	if (u.overwritten & userOAvatar) > 0 {
		user.Avatar = u.Avatar
	}
	if (u.overwritten & userOToken) > 0 {
		user.Token = u.Token
	}
	if (u.overwritten & userOVerified) > 0 {
		user.Verified = u.Verified
	}
	if (u.overwritten & userOMFAEnabled) > 0 {
		user.MFAEnabled = u.MFAEnabled
	}
	if (u.overwritten & userOBot) > 0 {
		user.Bot = u.Bot
	}
	user.overwritten = u.overwritten

	if constant.LockedMethods {
		u.RUnlock()
		user.Unlock()
	}

	return
}

func (u *User) saveToDiscord(session Session) (err error) {
	var myself *User
	myself, err = session.Myself()
	if err != nil {
		return
	}
	if myself == nil {
		err = errors.New("can't get information about current user")
		return
	}

	if myself.ID != u.ID {
		err = errors.New("can only update current user")
		return
	}

	params := &ModifyCurrentUserParams{}
	if u.Username != "" {
		params.SetUsername(u.Username)
	}
	if u.Avatar != nil && u.Avatar != myself.Avatar {
		// TODO: allow resetting the avatar, somehow
		params.SetAvatar(*u.Avatar)
	}

	var updated *User
	updated, err = session.ModifyCurrentUser(params)
	if err != nil {
		return
	}

	*u = *updated
	return
}

// Valid ensure the user object has enough required information to be used in Discord interactions
func (u *User) Valid() bool {
	return u.ID > 0
}

// -------

// NewUserPresence creates a new user presence instance
func NewUserPresence() *UserPresence {
	return &UserPresence{
		Roles: []Snowflake{},
	}
}

// UserPresence presence info for a guild member or friend/user in a DM
type UserPresence struct {
	Lockable `json:"-"`

	User    *User       `json:"user"`
	Roles   []Snowflake `json:"roles"`
	Game    *Activity   `json:"activity"`
	GuildID Snowflake   `json:"guild_id"`
	Nick    string      `json:"nick"`
	Status  string      `json:"status"`
}

func (p *UserPresence) String() string {
	return p.Status
}

// DeepCopy see interface at struct.go#DeepCopier
func (p *UserPresence) DeepCopy() (copy interface{}) {
	copy = NewUserPresence()
	p.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (p *UserPresence) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var presence *UserPresence
	if presence, ok = other.(*UserPresence); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *UserPresence")
		return
	}

	if constant.LockedMethods {
		p.RLock()
		presence.Lock()
	}

	presence.User = p.User.DeepCopy().(*User)
	presence.Roles = p.Roles
	presence.Game = p.Game.DeepCopy().(*Activity)
	presence.GuildID = p.GuildID
	presence.Nick = p.Nick
	presence.Status = p.Status

	if constant.LockedMethods {
		p.RUnlock()
		presence.Unlock()
	}

	return
}

// UserConnection ...
type UserConnection struct {
	Lockable `json:"-"`

	ID           string                `json:"id"`           // id of the connection account
	Name         string                `json:"name"`         // the username of the connection account
	Type         string                `json:"type"`         // the service of the connection (twitch, youtube)
	Revoked      bool                  `json:"revoked"`      // whether the connection is revoked
	Integrations []*IntegrationAccount `json:"integrations"` // an array of partial server integrations
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *UserConnection) DeepCopy() (copy interface{}) {
	copy = &UserConnection{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *UserConnection) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var con *UserConnection
	if con, ok = other.(*UserConnection); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *UserConnection")
		return
	}

	if constant.LockedMethods {
		c.RLock()
		con.Lock()
	}

	con.ID = c.ID
	con.Name = c.Name
	con.Type = c.Type
	con.Revoked = c.Revoked

	con.Integrations = make([]*IntegrationAccount, len(c.Integrations))
	for i, account := range c.Integrations {
		con.Integrations[i] = account.DeepCopy().(*IntegrationAccount)
	}

	if constant.LockedMethods {
		c.RUnlock()
		con.Unlock()
	}

	return
}
