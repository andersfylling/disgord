package disgord

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
)

type ClientStatus struct {
	Desktop string `json:"desktop"`
	Mobile  string `json:"mobile"`
	Web     string `json:"web"`
}

// ActivityParty ...
type ActivityParty struct {
	ID   string `json:"id,omitempty"`   // the id of the party
	Size []int  `json:"size,omitempty"` // used to show the party's current and maximum size
}

var _ Copier = (*ActivityParty)(nil)
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
	LargeImage string `json:"large_image,omitempty"` // the id for a large asset of the activity, usually a snowflake
	LargeText  string `json:"large_text,omitempty"`  //text displayed when hovering over the large image of the activity
	SmallImage string `json:"small_image,omitempty"` // the id for a small asset of the activity, usually a snowflake
	SmallText  string `json:"small_text,omitempty"`  //	text displayed when hovering over the small image of the activity
}

var _ Copier = (*ActivityAssets)(nil)
var _ DeepCopier = (*ActivityAssets)(nil)

// ActivitySecrets ...
type ActivitySecrets struct {
	Join     string `json:"join,omitempty"`     // the secret for joining a party
	Spectate string `json:"spectate,omitempty"` // the secret for spectating a game
	Match    string `json:"match,omitempty"`    // the secret for a specific instanced match
}

var _ Copier = (*ActivitySecrets)(nil)
var _ DeepCopier = (*ActivitySecrets)(nil)

// ActivityEmoji ...
type ActivityEmoji struct {
	Name     string    `json:"name"`
	ID       Snowflake `json:"id,omitempty"`
	Animated bool      `json:"animated,omitempty"`
}

var _ Copier = (*ActivityEmoji)(nil)
var _ DeepCopier = (*ActivityEmoji)(nil)

// ActivityTimestamp ...
type ActivityTimestamp struct {
	Start int `json:"start,omitempty"` // unix time (in milliseconds) of when the activity started
	End   int `json:"end,omitempty"`   // unix time (in milliseconds) of when the activity ends
}

var _ Copier = (*ActivityTimestamp)(nil)
var _ DeepCopier = (*ActivityTimestamp)(nil)

// ######################
// ##
// ## Activity
// ##
// ######################

// ActivityType https://discord.com/developers/docs/topics/gateway#activity-object-activity-types
type ActivityType uint

const (
	ActivityTypeGame ActivityType = iota
	ActivityTypeStreaming
	ActivityTypeListening
	_
	ActivityTypeCustom
	ActivityTypeCompeting
)

// ActivityFlag https://discord.com/developers/docs/topics/gateway#activity-object-activity-flags
type ActivityFlag uint

// flags for the Activity object to signify the type of action taken place
const (
	ActivityFlagInstance ActivityFlag = 1 << iota
	ActivityFlagJoin
	ActivityFlagSpectate
	ActivityFlagJoinRequest
	ActivityFlagSync
	ActivityFlagPlay
)

// Activity https://discord.com/developers/docs/topics/gateway#activity-object-activity-structure
type Activity struct {
	Name          string             `json:"name"`
	Type          ActivityType       `json:"type"`
	URL           string             `json:"url,omitempty"`
	CreatedAt     int                `json:"created_at"`
	Timestamps    *ActivityTimestamp `json:"timestamps,omitempty"`
	ApplicationID Snowflake          `json:"application_id,omitempty"`
	Details       string             `json:"details,omitempty"`
	State         string             `json:"state,omitempty"`
	Emoji         *ActivityEmoji     `json:"emoji,omitempty"`
	Party         *ActivityParty     `json:"party,omitempty"`
	Assets        *ActivityAssets    `json:"assets,omitempty"`
	Secrets       *ActivitySecrets   `json:"secrets,omitempty"`
	Instance      bool               `json:"instance,omitempty"`
	Flags         ActivityFlag       `json:"flags,omitempty"`
}

var _ Reseter = (*Activity)(nil)
var _ Copier = (*Activity)(nil)
var _ DeepCopier = (*Activity)(nil)

// ---------

type UserFlag uint64

const (
	UserFlagNone            UserFlag = 0
	UserFlagDiscordEmployee UserFlag = 1 << iota
	UserFlagDiscordPartner
	UserFlagHypeSquadEvents
	UserFlagBugHunterLevel1
	_
	_
	UserFlagHouseBravery
	UserFlagHouseBrilliance
	UserFlagHouseBalance
	UserFlagEarlySupporter
	UserFlagTeamUser
	_
	UserFlagSystem
	_
	UserFlagBugHunterLevel2
	_
	UserFlagVerifiedBot
	UserFlagVerifiedBotDeveloper
)

type PremiumType int

func (p PremiumType) String() (t string) {
	switch p {
	case PremiumTypeNone:
		t = "None"
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
	PremiumTypeNone PremiumType = iota
	PremiumTypeNitroClassic
	PremiumTypeNitro
)

// User the Discord user object which is reused in most other data structures.
type User struct {
	ID            Snowflake     `json:"id,omitempty"`
	Username      string        `json:"username,omitempty"`
	Discriminator Discriminator `json:"discriminator,omitempty"`
	Avatar        string        `json:"avatar"` // data:image/jpeg;base64,BASE64_ENCODED_JPEG_IMAGE_DATA
	Bot           bool          `json:"bot,omitempty"`
	System        bool          `json:"system,omitempty"`
	MFAEnabled    bool          `json:"mfa_enabled,omitempty"`
	Locale        string        `json:"locale,omitempty"`
	Verified      bool          `json:"verified,omitempty"`
	Email         string        `json:"email,omitempty"`
	Flags         UserFlag      `json:"flag,omitempty"`
	PremiumType   PremiumType   `json:"premium_type,omitempty"`
	PublicFlags   UserFlag      `json:"public_flag,omitempty"`
	PartialMember *Member       `json:"member"` // may be populated by Message
}

var _ Reseter = (*User)(nil)
var _ DeepCopier = (*User)(nil)
var _ Copier = (*User)(nil)
var _ Mentioner = (*User)(nil)

// Mention returns the a string that Discord clients can format into a valid Discord mention
func (u *User) Mention() string {
	return "<@" + u.ID.String() + ">"
}

// AvatarURL returns a link to the Users avatar with the given size.
func (u *User) AvatarURL(size int, preferGIF bool) (url string, err error) {
	if size > 2048 || size < 16 || (size&(size-1)) > 0 {
		return "", errors.New("image size can be any power of two between 16 and 2048")
	}

	if u.Avatar == "" {
		url = fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%d.png?size=%d", u.Discriminator%5, size)
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

// SendMsg send a message to a user where you utilize a Message object instead of a string
func (u *User) SendMsg(ctx context.Context, session Session, message *Message) (channel *Channel, msg *Message, err error) {
	channel, err = session.User(u.ID).WithContext(ctx).CreateDM()
	if err != nil {
		return
	}

	msg, err = session.WithContext(ctx).SendMsg(channel.ID, message)
	return
}

// SendMsgString send a message to given user where the message is in the form of a string.
func (u *User) SendMsgString(ctx context.Context, session Session, content string) (channel *Channel, msg *Message, err error) {
	channel, msg, err = u.SendMsg(ctx, session, &Message{
		Content: content,
	})
	return
}

// Valid ensure the user object has enough required information to be used in Discord interactions
func (u *User) Valid() bool {
	return u.ID > 0
}

// -------

// UserPresence presence info for a guild member or friend/user in a DM
type UserPresence struct {
	User    *User       `json:"user"`
	Roles   []Snowflake `json:"roles"`
	Game    *Activity   `json:"activity"`
	GuildID Snowflake   `json:"guild_id"`
	Nick    string      `json:"nick"`
	Status  string      `json:"status"`
}

var _ Copier = (*UserPresence)(nil)
var _ DeepCopier = (*UserPresence)(nil)

func (p *UserPresence) String() string {
	return p.Status
}

// UserConnection ...
type UserConnection struct {
	ID           string                `json:"id"`           // id of the connection account
	Name         string                `json:"name"`         // the username of the connection account
	Type         string                `json:"type"`         // the service of the connection (twitch, youtube)
	Revoked      bool                  `json:"revoked"`      // whether the connection is revoked
	Integrations []*IntegrationAccount `json:"integrations"` // an array of partial server integrations
}

var _ Copier = (*UserConnection)(nil)
var _ DeepCopier = (*UserConnection)(nil)

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

// RESTUser REST interface for all user endpoints
type UserQueryBuilder interface {
	WithContext(ctx context.Context) UserQueryBuilder

	// GetUser Returns a user object for a given user Snowflake.
	Get(flags ...Flag) (*User, error)

	// CreateDM Create a new DM channel with a user. Returns a DM channel object.
	CreateDM(flags ...Flag) (ret *Channel, err error)
}

// Guild is used to create a guild query builder.
func (c clientQueryBuilder) User(id Snowflake) UserQueryBuilder {
	return &userQueryBuilder{client: c.client, uid: id}
}

// The default guild query builder.
type userQueryBuilder struct {
	ctx    context.Context
	client *Client
	uid    Snowflake
}

func (c userQueryBuilder) WithContext(ctx context.Context) UserQueryBuilder {
	c.ctx = ctx
	return c
}

// GetUser [REST] Returns a user object for a given user Snowflake.
//  Method                  GET
//  Endpoint                /users/{user.id}
//  Discord documentation   https://discord.com/developers/docs/resources/user#get-user
//  Reviewed                2018-06-10
//  Comment                 -
func (c userQueryBuilder) Get(flags ...Flag) (*User, error) {
	if !ignoreCache(flags...) {
		if usr, _ := c.client.cache.GetUser(c.uid); usr != nil {
			return usr, nil
		}
	}

	r := c.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.User(c.uid),
		Ctx:      c.ctx,
	}, flags)
	r.pool = c.client.pool.user
	r.factory = userFactory

	return getUser(r.Execute)
}

// CreateDM [REST] Create a new DM channel with a user. Returns a DM channel object.
//  Method                  POST
//  Endpoint                /users/@me/channels
//  Discord documentation   https://discord.com/developers/docs/resources/user#create-dm
//  Reviewed                2019-02-23
//  Comment                 -
func (c userQueryBuilder) CreateDM(flags ...Flag) (ret *Channel, err error) {
	r := c.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodPost,
		Ctx:      c.ctx,
		Endpoint: endpoint.UserMeChannels(),
		Body: &struct {
			RecipientID Snowflake `json:"recipient_id"`
		}{c.uid},
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.factory = func() interface{} {
		return &Channel{}
	}

	return getChannel(r.Execute)
}

type CurrentUserQueryBuilder interface {
	WithContext(ctx context.Context) CurrentUserQueryBuilder

	// GetCurrentUser Returns the user object of the requester's account. For OAuth2, this requires the identify
	// scope, which will return the object without an email, and optionally the email scope, which returns the object
	// with an email.
	Get(flags ...Flag) (*User, error)

	// UpdateCurrentUser Modify the requester's user account settings. Returns a user object on success.
	UpdateBuilder(flags ...Flag) UpdateCurrentUserBuilder

	// Deprecated: use UpdateBuilder
	Update(flags ...Flag) UpdateCurrentUserBuilder

	// GetCurrentUserGuilds Returns a list of partial guild objects the current user is a member of.
	// Requires the Guilds OAuth2 scope.
	GetGuilds(params *GetCurrentUserGuildsParams, flags ...Flag) (ret []*Guild, err error)

	// LeaveGuild Leave a guild. Returns a 204 empty response on success.
	LeaveGuild(id Snowflake, flags ...Flag) (err error)

	// CreateGroupDM Create a new group DM channel with multiple Users. Returns a DM channel object.
	// This endpoint was intended to be used with the now-deprecated GameBridge SDK. DMs created with this
	// endpoint will not be shown in the Discord Client
	CreateGroupDM(params *CreateGroupDMParams, flags ...Flag) (ret *Channel, err error)

	// GetUserConnections Returns a list of connection objects. Requires the connections OAuth2 scope.
	GetUserConnections(flags ...Flag) (ret []*UserConnection, err error)
}

// Guild is used to create a guild query builder.
func (c clientQueryBuilder) CurrentUser() CurrentUserQueryBuilder {
	return &currentUserQueryBuilder{client: c.client}
}

// The default guild query builder.
type currentUserQueryBuilder struct {
	ctx    context.Context
	client *Client
}

var _ CurrentUserQueryBuilder = (*currentUserQueryBuilder)(nil)

func (c currentUserQueryBuilder) WithContext(ctx context.Context) CurrentUserQueryBuilder {
	c.ctx = ctx
	return &c
}

// GetCurrentUser [REST] Returns the user object of the requester's account. For OAuth2, this requires the identify
// scope, which will return the object without an email, and optionally the email scope, which returns the object
// with an email.
//  Method                  GET
//  Endpoint                /users/@me
//  Discord documentation   https://discord.com/developers/docs/resources/user#get-current-user
//  Reviewed                2019-02-23
//  Comment                 -
func (c currentUserQueryBuilder) Get(flags ...Flag) (user *User, err error) {
	if !ignoreCache(flags...) {
		if usr, err := c.client.cache.GetCurrentUser(); err != nil && usr != nil {
			return usr, nil
		}
	}

	r := c.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.UserMe(),
		Ctx:      c.ctx,
	}, flags)
	r.pool = c.client.pool.user
	r.factory = userFactory

	return getUser(r.Execute)
}

// UpdateCurrentUser [REST] Modify the requester's user account settings. Returns a user object on success.
//  Method                  PATCH
//  Endpoint                /users/@me
//  Discord documentation   https://discord.com/developers/docs/resources/user#modify-current-user
//  Reviewed                2019-02-18
//  Comment                 -
func (c currentUserQueryBuilder) UpdateBuilder(flags ...Flag) UpdateCurrentUserBuilder {
	builder := &updateCurrentUserBuilder{}
	builder.r.itemFactory = userFactory // TODO: peak cached user
	builder.r.flags = flags
	builder.r.setup(c.client.req, &httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         c.ctx,
		Endpoint:    endpoint.UserMe(),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	// TODO: cache changes?
	return builder
}

// GetCurrentUserGuildsParams JSON params for func GetCurrentUserGuilds
type GetCurrentUserGuildsParams struct {
	Before Snowflake `urlparam:"before,omitempty"`
	After  Snowflake `urlparam:"after,omitempty"`
	Limit  int       `urlparam:"limit,omitempty"`
}

var _ URLQueryStringer = (*GetCurrentUserGuildsParams)(nil)

// GetCurrentUserGuilds [REST] Returns a list of partial guild objects the current user is a member of.
// Requires the Guilds OAuth2 scope.
//  Method                  GET
//  Endpoint                /users/@me/guilds
//  Discord documentation   https://discord.com/developers/docs/resources/user#get-current-user-guilds
//  Reviewed                2019-02-18
//  Comment                 This endpoint. returns 100 Guilds by default, which is the maximum number of
//                          Guilds a non-bot user can join. Therefore, pagination is not needed for
//                          integrations that need to get a list of Users' Guilds.
func (c currentUserQueryBuilder) GetGuilds(params *GetCurrentUserGuildsParams, flags ...Flag) (ret []*Guild, err error) {
	r := c.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.UserMeGuilds(),
		Ctx:      c.ctx,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Guild, 0)
		return &tmp
	}

	var vs interface{}
	if vs, err = r.Execute(); err != nil {
		return nil, err
	}

	if guilds, ok := vs.(*[]*Guild); ok {
		return *guilds, nil
	}
	return nil, errors.New("unable to cast guild slice")
}

// CreateGroupDMParams required JSON params for func CreateGroupDM
// https://discord.com/developers/docs/resources/user#create-group-dm
type CreateGroupDMParams struct {
	// AccessTokens access tokens of Users that have granted your app the gdm.join scope
	AccessTokens []string `json:"access_tokens"`

	// map[UserID] = nickname
	Nicks map[Snowflake]string `json:"nicks"`
}

// LeaveGuild [REST] Leave a guild. Returns a 204 empty response on success.
//  Method                  DELETE
//  Endpoint                /users/@me/guilds/{guild.id}
//  Discord documentation   https://discord.com/developers/docs/resources/user#leave-guild
//  Reviewed                2019-02-18
//  Comment                 -
func (c currentUserQueryBuilder) LeaveGuild(id Snowflake, flags ...Flag) (err error) {
	r := c.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.UserMeGuild(id),
		Ctx:      c.ctx,
	}, flags)

	_, err = r.Execute()
	return
}

// CreateGroupDM [REST] Create a new group DM channel with multiple Users. Returns a DM channel object.
// This endpoint was intended to be used with the now-deprecated GameBridge SDK. DMs created with this
// endpoint will not be shown in the Discord Client
//  Method                  POST
//  Endpoint                /users/@me/channels
//  Discord documentation   https://discord.com/developers/docs/resources/user#create-group-dm
//  Reviewed                2019-02-19
//  Comment                 -
func (c currentUserQueryBuilder) CreateGroupDM(params *CreateGroupDMParams, flags ...Flag) (ret *Channel, err error) {
	r := c.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Ctx:         c.ctx,
		Endpoint:    endpoint.UserMeChannels(),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.factory = func() interface{} {
		return &Channel{}
	}

	// TODO: go generate casting func: return getChannel(r.Execute)
	return getChannel(r.Execute)
}

// GetUserConnections [REST] Returns a list of connection objects. Requires the connections OAuth2 scope.
//  Method                  GET
//  Endpoint                /users/@me/connections
//  Discord documentation   https://discord.com/developers/docs/resources/user#get-user-connections
//  Reviewed                2019-02-19
//  Comment                 -
func (c currentUserQueryBuilder) GetUserConnections(flags ...Flag) (connections []*UserConnection, err error) {
	r := c.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.UserMeConnections(),
		Ctx:      c.ctx,
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
