package resource

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/snowflake"
)

const (
	EndpointUser              = "/users/"
	EndpointUserMyself        = EndpointUser + "@me"
	EndpointUserMyGuilds      = EndpointUserMyself + "/guilds"
	EndpointUserMyChannels    = EndpointUserMyself + "/channels"
	EndpointUserMyConnections = EndpointUserMyself + "/connections"
)


type UserInterface interface {
	Mention() string
	MentionNickname() string
	String() string
}

// TODO: https://discordapp.com/developers/docs/topics/gateway#activity-object-activity-structure
type UserActivity struct{}

// ---------

// TODO: should a user object always have a ID?
func NewUser() *User {
	return &User{}
}

type User struct {
	ID            snowflake.ID `json:"id,omitempty"`
	Username      string       `json:"username,omitempty"`
	Discriminator string       `json:"discriminator,omitempty"`
	Email         string       `json:"email,omitempty"`
	Avatar        *string      `json:"avatar"` // data:image/jpeg;base64,BASE64_ENCODED_JPEG_IMAGE_DATA
	Token         string       `json:"token,omitempty"`
	Verified      bool         `json:"verified,omitempty"`
	MFAEnabled    bool         `json:"mfa_enabled,omitempty"`
	Bot           bool         `json:"bot,omitempty"`

	sync.RWMutex `json:"-"`
}

func (u *User) Mention() string {
	return "<@" + u.ID.String() + ">"
}

func (u *User) MentionNickname() string {
	return "<@!" + u.ID.String() + ">"
}

func (u *User) String() string {
	return u.Username + "#" + u.Discriminator + "{" + u.ID.String() + "}"
}

// Partial check if this is not a complete user object
// Assumption: has a snowflake.
func (u *User) Partial() bool {
	return (u.Username + u.Discriminator) == ""
}

func (u *User) MarshalJSON() ([]byte, error) {
	if u.ID.Empty() {
		return []byte("{}"), nil
	}

	return json.Marshal(User(*u))
}

// func (u *User) UnmarshalJSON(data []byte) error {
// 	return json.Unmarshal(data, &u.userJSON)
// }

func (u *User) Clear() {
	//u.d.Avatar = nil
}

func (u *User) SendMessage(requester httd.Requester, msg *Message) (channelID snowflake.ID, messageID snowflake.ID, err error) {
	return snowflake.NewID(0), snowflake.NewID(0), errors.New("not implemented")
}

func (u *User) SendMessageStr(requester httd.Requester, msg string) (channelID snowflake.ID, messageID snowflake.ID, err error) {
	return snowflake.NewID(0), snowflake.NewID(0), errors.New("not implemented")
}

func (u *User) DeepCopy() *User {
	user := NewUser()

	u.RLock()

	user.ID = u.ID
	user.Username = u.Username
	user.Discriminator = u.Discriminator
	user.Email = u.Email
	user.Token = u.Token
	user.Verified = u.Verified
	user.MFAEnabled = u.MFAEnabled
	user.Bot = u.Bot

	if u.Avatar != nil {
		avatar := *u.Avatar
		user.Avatar = &avatar
	}

	u.RUnlock()

	return user
}

// -------

type UserPresence struct {
	User    *User          `json:"user"`
	Roles   []snowflake.ID `json:"roles"`
	Game    *UserActivity  `json:"activity"`
	GuildID snowflake.ID   `json:"guild_id"`
	Nick    string         `json:"nick"`
	Status  string         `json:"status"`
}

func NewUserPresence() *UserPresence {
	return &UserPresence{}
}

func (p *UserPresence) Update(status string) {
	// Update the presence.
	// talk to the discord api
}

func (p *UserPresence) String() string {
	return p.Status
}

func (p *UserPresence) Clear() {
	p.Game = nil
}

// ----------

// ReqGetUser [GET]         Returns the user object of the requester's account. For OAuth2, this requires
//                          the identify scope, which will return the object without an email, and optionally
//                          the email scope, which returns the object with an email.
// Endpoint                 /users/@me
// Rate limiter             /users
// Discord documentation    https://discordapp.com/developers/docs/resources/user#get-current-user
// Reviewed                 2018-06-10
// Comment                  -
func ReqGetCurrentUser(client httd.Getter) (ret *User, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint: EndpointUserMyself,
	}
	resp, err := client.Get(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

// ReqGetUser [GET]         Returns a user object for a given user ID.
// Endpoint                 /users/{user.id}
// Rate limiter             /users
// Discord documentation    https://discordapp.com/developers/docs/resources/user#get-user
// Reviewed                 2018-06-10
// Comment                  -
func ReqGetUser(client httd.Getter, userID snowflake.ID) (ret *User, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint: EndpointUser + userID.String(),
	}
	resp, err := client.Get(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

type ReqModifyCurrentUserParams struct {
	Username string `json:"username,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// ReqModifyCurrentUser [PATCH]  Modify the requester's user account settings. Returns a user object on success.
// Endpoint                     /users/@me
// Rate limiter                 /users
// Discord documentation        https://discordapp.com/developers/docs/resources/user#modify-current-user
// Reviewed                     2018-06-10
// Comment                      -
func ReqModifyCurrentUser(client httd.Getter, params *ReqModifyCurrentUserParams) (ret *User, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint: EndpointUserMyself,
		JSONParams: params,
	}
	resp, err := client.Get(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

type ReqGetCurrentUserGuildsParams struct {
	Before snowflake.ID `urlparam:"before,omitempty"`
	After  snowflake.ID `urlparam:"after,omitempty"`
	Limit  int          `urlparam:"limit,omitempty"`
}

// getQueryString this ins't really pretty, but it works.
func (params *ReqGetCurrentUserGuildsParams) getQueryString() string {
	seperator := "?"
	query := ""

	if !params.Before.Empty() {
		query += seperator + params.Before.String()
		seperator = "&"
	}

	if !params.After.Empty() {
		query += seperator + params.After.String()
		seperator = "&"
	}

	if params.Limit > 0 {
		query += seperator + strconv.Itoa(params.Limit)
	}

	return query
}

// ReqGetCurrentUserGuilds [GET]  Returns a list of partial guild objects the current user is a member of.
//                                Requires the guilds OAuth2 scope.
// Endpoint                       /users/@me/guilds
// Rate limiter                   /users
// Discord documentation          https://discordapp.com/developers/docs/resources/user#get-current-user-guilds
// Reviewed                       2018-06-10
// Comment                        This endpoint returns 100 guilds by default, which is the maximum number of
//                                guilds a non-bot user can join. Therefore, pagination is not needed for
//                                integrations that need to get a list of users' guilds.
func ReqGetCurrentUserGuilds(client httd.Getter, params *ReqGetCurrentUserGuildsParams) (ret []*Guild, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint: EndpointUserMyGuilds + params.getQueryString(),
	}
	resp, err := client.Get(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

// ReqLeaveGuild [DELETE] Leave a guild. Returns a 204 empty response on success.
// Endpoint               /users/@me/guilds/{guild.id}
// Rate limiter           /users
// Discord documentation  https://discordapp.com/developers/docs/resources/user#leave-guild
// Reviewed               2018-06-10
// Comment                -
func ReqLeaveGuild(client httd.Deleter, guildID snowflake.ID) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint: EndpointUserMyGuilds + "/" + guildID.String(),
	}
	resp, err := client.Delete(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		msg := "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(http.StatusNoContent)
		err = errors.New(msg)
	}

	return
}

// ReqGetUserDMs [GET]    Returns a list of DM channel objects.
// Endpoint               /users/@me/channels
// Rate limiter           /users
// Discord documentation  https://discordapp.com/developers/docs/resources/user#get-user-dms
// Reviewed               2018-06-10
// Comment                -
func ReqGetUserDMs(client httd.Getter) (ret []*Channel, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint: EndpointUserMyChannels,
	}
	resp, err := client.Get(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

type BodyUserCreateDM struct {
	RecipientID snowflake.ID `json:"recipient_id"`
}

// ReqGetUserDMs [POST]   Create a new DM channel with a user. Returns a DM channel object.
// Endpoint               /users/@me/channels
// Rate limiter           /users
// Discord documentation  https://discordapp.com/developers/docs/resources/user#create-dm
// Reviewed               2018-06-10
// Comment                -
func ReqCreateDM(client httd.Poster, recipientID snowflake.ID) (ret *Channel, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint: EndpointUserMyChannels,
		JSONParams: &BodyUserCreateDM{recipientID},
	}
	resp, err := client.Post(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

// BodyUserCreateGroupDM
// https://discordapp.com/developers/docs/resources/user#create-group-dm
type ReqCreateGroupDMParams struct {
	AccessTokens []string                `json:"access_tokens"` // access tokens of users that have granted your app the gdm.join scope
	Nicks        map[snowflake.ID]string `json:"nicks"`         // userID => nickname
}

// ReqCreateGroupDM [POST]  Create a new group DM channel with multiple users. Returns a DM channel object.
// Endpoint                 /users/@me/channels
// Rate limiter             /users
// Discord documentation    https://discordapp.com/developers/docs/resources/user#create-group-dm
// Reviewed                 2018-06-10
// Comment                  -
func ReqCreateGroupDM(client httd.Poster, params *ReqCreateGroupDMParams) (ret *Channel, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint: EndpointUserMyChannels,
		JSONParams: params,
	}
	resp, err := client.Post(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

// TODO: is this a partial Integration object?
type IntegrationAccount struct {
	ID string `json:"id"` // id of the account
	Name string `json:"name"` // name of the account
}

type UserConnection struct {
	ID string `json:"id"` // id of the connection account
	Name string `json:"name"` // the username of the connection account
	Type string `json:"type"` // the service of the connection (twitch, youtube)
	Revoked bool `json:"revoked"` // whether the connection is revoked
	Integrations []*IntegrationAccount `json:"integrations"` // an array of partial server integrations
}

// ReqCreateGroupDM [GET] Returns a list of connection objects. Requires the connections OAuth2 scope.
// Endpoint               /users/@me/connections
// Rate limiter           /users
// Discord documentation  https://discordapp.com/developers/docs/resources/user#get-user-connections
// Reviewed               2018-06-10
// Comment                -
func ReqGetUserConnections(client httd.Getter) (ret []*UserConnection, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint: EndpointUserMyConnections,
	}
	resp, err := client.Get(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}
