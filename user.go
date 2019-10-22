package disgord

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/andersfylling/disgord/internal/constant"
	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
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

// Activity types https://discordapp.com/developers/docs/topics/gateway#activity-object-activity-types
const (
	ActivityTypeGame = iota
	ActivityTypeStreaming
	ActivityTypeListening
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

var _ DeepCopier = (*ActivityParty)(nil)

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

var _ DeepCopier = (*ActivitySecrets)(nil)

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
		Timestamps: &ActivityTimestamp{},
	}
}

// Activity https://discordapp.com/developers/docs/topics/gateway#activity-object-activity-structure
type Activity struct {
	Lockable `json:"-"`

	Name          string             `json:"name"`                     // the activity's name
	Type          int                `json:"type"`                     // activity type
	URL           string             `json:"url,omitempty"`            //stream url, is validated when type is 1
	Timestamps    *ActivityTimestamp `json:"timestamps,omitempty"`     // timestamps object	unix timestamps for start and/or end of the game
	ApplicationID Snowflake          `json:"application_id,omitempty"` //?	snowflake	application id for the game
	Details       string             `json:"details,omitempty"`        //?	?string	what the player is currently doing
	State         string             `json:"state,omitempty"`          //state?	?string	the user's current party status
	Party         *ActivityParty     `json:"party,omitempty"`          //party?	party object	information for the current party of the player
	Assets        *ActivityAssets    `json:"assets,omitempty"`         // assets?	assets object	images for the presence and their hover texts
	Secrets       *ActivitySecrets   `json:"secrets,omitempty"`        // secrets?	secrets object	secrets for Rich Presence joining and spectating
	Instance      bool               `json:"instance,omitempty"`       // instance?	boolean	whether or not the activity is an instanced game session
	Flags         int                `json:"flags,omitempty"`          // flags?	int	activity flags ORd together, describes what the payload includes
}

var _ Reseter = (*Activity)(nil)
var _ DeepCopier = (*Activity)(nil)

// ---------

const (
	userOEmail = 0x1 << iota
	userOAvatar
	userOToken
	userOVerified
	userOMFAEnabled
	userOBot
	userOPremiumType
)

type PremiumType int

func (p PremiumType) String() (t string) {
	switch p {
	case PremiumTypeNitroClassic:
		t = "Nitro Classic"
	case PremiumTypeNitro:
		t = "Nitro"
	default:
		t = ""
	}

	return t
}

var _ fmt.Stringer = (*PremiumType)(nil)

const (
	// PremiumTypeNitroClassic includes app perks like animated emojis and avatars, but not games
	PremiumTypeNitroClassic PremiumType = 1

	// PremiumTypeNitro includes app perks as well as the games subscription service
	PremiumTypeNitro PremiumType = 2
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
	/*7*/ PremiumType *PremiumType `json:"premium_type,omitempty"`
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
	if u.PremiumType != nil {
		overwritten |= userOPremiumType
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
	Avatar        string        `json:"avatar"` // data:image/jpeg;base64,BASE64_ENCODED_JPEG_IMAGE_DATA //TODO: pointer?
	Token         string        `json:"token,omitempty"`
	Verified      bool          `json:"verified,omitempty"`
	MFAEnabled    bool          `json:"mfa_enabled,omitempty"`
	Bot           bool          `json:"bot,omitempty"`
	PremiumType   PremiumType   `json:"premium_type,omitempty"`

	// Used to identify which fields are set by Discord in partial JSON objects. Yep.
	overwritten uint8 // map. see number left of field in userJSON struct.
}

var _ Reseter = (*User)(nil)
var _ DeepCopier = (*User)(nil)

// Mention returns the a string that Discord clients can format into a valid Discord mention
func (u *User) Mention() string {
	return "<@" + u.ID.String() + ">"
}

// AvatarURL returns a link to the users avatar with the given size.
func (u *User) AvatarURL(size int, preferGIF bool) (url string, err error) {
	if size > 2048 || size < 16 || (size&(size-1)) > 0 {
		return "", errors.New("image size can be any power of two between 16 and 2048")
	}

	if u.Avatar == "" {
		url = fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%d.webp?size=%d", u.Discriminator%5, size)
	} else if strings.HasPrefix(u.Avatar, "a_") && preferGIF {
		url = fmt.Sprintf("https://cdn.discordapp.com/avatars/%d/%s.gif?size=%d", u.ID, u.Avatar, size)
	} else {
		url = fmt.Sprintf("https://cdn.discordapp.com/avatars/%d/%s.webp?size=%d", u.ID, u.Avatar, size)
	}

	return
}

// Tag formats the user to Anders#1234
func (u *User) Tag() string {
	return u.Username + "#" + u.Discriminator.String()
}

// String formats the user to Anders#1234{1234567890}
func (u *User) String() string {
	return u.Tag() + "{" + u.ID.String() + "}"
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
	if (changes&userOAvatar) > 0 && j.Avatar != nil {
		u.Avatar = *j.Avatar
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
	if (changes & userOPremiumType) > 0 {
		u.PremiumType = *j.PremiumType
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
	user.Avatar = u.Avatar
	user.overwritten = u.overwritten

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

	if !u.ID.IsZero() {
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

var _ DeepCopier = (*UserPresence)(nil)
var _ fmt.Stringer = (*UserPresence)(nil)

func (p *UserPresence) String() string {
	return p.Status
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

var _ DeepCopier = (*UserConnection)(nil)

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

func ratelimitUsers() string {
	return "u"
}

// GetCurrentUserGuildsParams JSON params for func GetCurrentUserGuilds
type GetCurrentUserGuildsParams struct {
	Before Snowflake `urlparam:"before,omitempty"`
	After  Snowflake `urlparam:"after,omitempty"`
	Limit  int       `urlparam:"limit,omitempty"`
}

var _ URLQueryStringer = (*GetCurrentUserGuildsParams)(nil)

// GetCurrentUser [REST] Returns the user object of the requester's account. For OAuth2, this requires the identify
// scope, which will return the object without an email, and optionally the email scope, which returns the object
// with an email.
//  Method                  GET
//  Endpoint                /users/@me
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#get-current-user
//  Reviewed                2019-02-23
//  Comment                 -
func (c *Client) GetCurrentUser(flags ...Flag) (user *User, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.UserMe(),
	}, flags)
	r.CacheRegistry = UserCache
	r.ID = c.myID
	r.pool = c.pool.user
	r.factory = userFactory

	if user, err = getUser(r.Execute); err == nil {
		c.myID = user.ID
	}
	return user, err
}

// GetUser [REST] Returns a user object for a given user Snowflake.
//  Method                  GET
//  Endpoint                /users/{user.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#get-user
//  Reviewed                2018-06-10
//  Comment                 -
func (c *Client) GetUser(id Snowflake, flags ...Flag) (*User, error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.User(id),
	}, flags)
	r.CacheRegistry = UserCache
	r.ID = id
	r.pool = c.pool.user
	r.factory = userFactory

	return getUser(r.Execute)
}

// UpdateCurrentUser [REST] Modify the requester's user account settings. Returns a user object on success.
//  Method                  PATCH
//  Endpoint                /users/@me
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#modify-current-user
//  Reviewed                2019-02-18
//  Comment                 -
func (c *Client) UpdateCurrentUser(flags ...Flag) (builder *updateCurrentUserBuilder) {
	builder = &updateCurrentUserBuilder{}
	builder.r.itemFactory = userFactory // TODO: peak cached user
	builder.r.flags = flags
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      httd.MethodPatch,
		Endpoint:    endpoint.UserMe(),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	// TODO: cache changes?
	return builder
}

// GetCurrentUserGuilds [REST] Returns a list of partial guild objects the current user is a member of.
// Requires the guilds OAuth2 scope.
//  Method                  GET
//  Endpoint                /users/@me/guilds
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#get-current-user-guilds
//  Reviewed                2019-02-18
//  Comment                 This endpoint. returns 100 guilds by default, which is the maximum number of
//                          guilds a non-bot user can join. Therefore, pagination is not needed for
//                          integrations that need to get a list of users' guilds.
func (c *Client) GetCurrentUserGuilds(params *GetCurrentUserGuildsParams, flags ...Flag) (ret []*PartialGuild, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.UserMeGuilds(),
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*PartialGuild, 0)
		return &tmp
	}

	var vs interface{}
	if vs, err = r.Execute(); err != nil {
		return nil, err
	}

	if guilds, ok := vs.(*[]*PartialGuild); ok {
		return *guilds, nil
	}
	return nil, errors.New("unable to cast guild slice")
}

// LeaveGuild [REST] Leave a guild. Returns a 204 empty response on success.
//  Method                  DELETE
//  Endpoint                /users/@me/guilds/{guild.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#leave-guild
//  Reviewed                2019-02-18
//  Comment                 -
func (c *Client) LeaveGuild(id Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.UserMeGuild(id),
	}, flags)
	r.expectsStatusCode = http.StatusNoContent
	r.CacheRegistry = GuildCache
	r.ID = id
	r.updateCache = func(registry cacheRegistry, id Snowflake, x interface{}) (err error) {
		c.cache.DeleteGuild(id)
		return nil
	}

	_, err = r.Execute()
	return
}

// GetUserDMs [REST] Returns a list of DM channel objects.
//  Method                  GET
//  Endpoint                /users/@me/channels
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#get-user-dms
//  Reviewed                2019-02-19
//  Comment                 Apparently Discord removed support for this in 2016 and updated their docs 2 years after..
//							https://github.com/discordapp/discord-api-docs/issues/184
//							For now I'll just leave this here, until I can do a cache lookup. Making this cache
//							dependent.
// Deprecated: Needs cache checking to get the actual list of channels
func (c *Client) GetUserDMs(flags ...Flag) (ret []*Channel, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.UserMeChannels(),
	}, flags)
	r.CacheRegistry = ChannelCache
	r.factory = func() interface{} {
		tmp := make([]*Channel, 0) // TODO: use channel pool to get enough channels
		return &tmp
	}

	var vs interface{}
	if vs, err = r.Execute(); err != nil {
		return nil, err
	}

	if chans, ok := vs.(*[]*Channel); ok {
		return *chans, nil
	}
	return nil, errors.New("unable to cast guild slice")
}

// BodyUserCreateDM JSON param for func CreateDM
type BodyUserCreateDM struct {
	RecipientID Snowflake `json:"recipient_id"`
}

// CreateDM [REST] Create a new DM channel with a user. Returns a DM channel object.
//  Method                  POST
//  Endpoint                /users/@me/channels
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#create-dm
//  Reviewed                2019-02-23
//  Comment                 -
func (c *Client) CreateDM(recipientID Snowflake, flags ...Flag) (ret *Channel, err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Endpoint:    endpoint.UserMeChannels(),
		Body:        &BodyUserCreateDM{recipientID},
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.CacheRegistry = ChannelCache
	r.factory = func() interface{} {
		return &Channel{}
	}

	return getChannel(r.Execute)
}

// CreateGroupDMParams required JSON params for func CreateGroupDM
// https://discordapp.com/developers/docs/resources/user#create-group-dm
type CreateGroupDMParams struct {
	// AccessTokens access tokens of users that have granted your app the gdm.join scope
	AccessTokens []string `json:"access_tokens"`

	// map[userID] = nickname
	Nicks map[Snowflake]string `json:"nicks"`
}

// CreateGroupDM [REST] Create a new group DM channel with multiple users. Returns a DM channel object.
// This endpoint was intended to be used with the now-deprecated GameBridge SDK. DMs created with this
// endpoint will not be shown in the Discord Client
//  Method                  POST
//  Endpoint                /users/@me/channels
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#create-group-dm
//  Reviewed                2019-02-19
//  Comment                 -
func (c *Client) CreateGroupDM(params *CreateGroupDMParams, flags ...Flag) (ret *Channel, err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Endpoint:    endpoint.UserMeChannels(),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.CacheRegistry = ChannelCache
	r.factory = func() interface{} {
		return &Channel{}
	}

	// TODO: go generate casting func: return getChannel(r.Execute)
	return getChannel(r.Execute)
}

// GetUserConnections [REST] Returns a list of connection objects. Requires the connections OAuth2 scope.
//  Method                  GET
//  Endpoint                /users/@me/connections
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#get-user-connections
//  Reviewed                2019-02-19
//  Comment                 -
func (c *Client) GetUserConnections(flags ...Flag) (connections []*UserConnection, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.UserMeConnections(),
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*UserConnection, 0)
		return &tmp
	}

	var vs interface{}
	if vs, err = r.Execute(); err != nil {
		return nil, err
	}

	if cons, ok := vs.(*[]*UserConnection); ok {
		return *cons, nil
	}
	return nil, errors.New("unable to cast guild slice")
}

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

func userFactory() interface{} {
	return &User{}
}

func newUserRESTBuilder(userID Snowflake) *getUserBuilder {
	builder := &getUserBuilder{}
	builder.r.cacheRegistry = UserCache
	builder.r.cacheItemID = userID
	builder.r.itemFactory = userFactory

	return builder
}

// getUserBuilder ...
type getUserBuilder struct {
	r RESTBuilder
	c *Client
}

func (b *getUserBuilder) Execute() (user *User, err error) {
	var v interface{}
	if v, err = b.r.execute(); err != nil {
		return nil, err
	}

	return v.(*User), nil
}

// updateCurrentUserBuilder ...
//generate-rest-params: username:string, avatar:string,
//generate-rest-basic-execute: user:*User,
type updateCurrentUserBuilder struct {
	r RESTBuilder
}

// TODO: params should be url-params. But it works since we're using GET.
//generate-rest-params: before:Snowflake, after:Snowflake, limit:int,
//generate-rest-basic-execute: guilds:[]*Guild,
type getCurrentUserGuildsBuilder struct {
	r RESTBuilder
}

func (b *getCurrentUserGuildsBuilder) SetDefaultLimit() *getCurrentUserGuildsBuilder {
	delete(b.r.urlParams, "limit")
	return b
}

//generate-rest-basic-execute: cons:[]*UserConnection,
type getUserConnectionsBuilder struct {
	r RESTBuilder
}

//generate-rest-basic-execute: channel:*Channel,
type createDMBuilder struct {
	r RESTBuilder
}

//generate-rest-basic-execute: channels:[]*Channel,
type getUserDMsBuilder struct {
	r RESTBuilder
}

//generate-rest-basic-execute: channel:*Channel,
type createGroupDMBuilder struct {
	r RESTBuilder
}
