package rest

import (
	"errors"
	"net/http"
	"strconv"

	. "github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/disgord/rest/httd"
	"github.com/andersfylling/snowflake"
)

const (
	EndpointUser              = "/users/"
	EndpointUserMyself        = EndpointUser + "@me"
	EndpointUserMyGuilds      = EndpointUserMyself + "/guilds"
	EndpointUserMyChannels    = EndpointUserMyself + "/channels"
	EndpointUserMyConnections = EndpointUserMyself + "/connections"
)

// ----------

// ReqGetUser [GET]         Returns the user object of the requester's account. For OAuth2, this requires
//                          the identify scope, which will return the object without an email, and optionally
//                          the email scope, which returns the object with an email.
// Endpoint                 /users/@me
// Rate limiter             /users
// Discord documentation    https://discordapp.com/developers/docs/resources/user#get-current-user
// Reviewed                 2018-06-10
// Comment                  -
func GetCurrentUser(client httd.Getter) (ret *User, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint:    EndpointUserMyself,
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// ReqGetUser [GET]         Returns a user object for a given user ID.
// Endpoint                 /users/{user.id}
// Rate limiter             /users
// Discord documentation    https://discordapp.com/developers/docs/resources/user#get-user
// Reviewed                 2018-06-10
// Comment                  -
func GetUser(client httd.Getter, userID snowflake.ID) (ret *User, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint:    EndpointUser + userID.String(),
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

type ModifyCurrentUserParams struct {
	Username string `json:"username,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

// ModifyCurrentUser [PATCH]    Modify the requester's user account settings. Returns a user object on success.
// Endpoint                     /users/@me
// Rate limiter                 /users
// Discord documentation        https://discordapp.com/developers/docs/resources/user#modify-current-user
// Reviewed                     2018-06-10
// Comment                      -
func ModifyCurrentUser(client httd.Getter, params *ModifyCurrentUserParams) (ret *User, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint:    EndpointUserMyself,
		JSONParams:  params,
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

type GetCurrentUserGuildsParams struct {
	Before snowflake.ID `urlparam:"before,omitempty"`
	After  snowflake.ID `urlparam:"after,omitempty"`
	Limit  int          `urlparam:"limit,omitempty"`
}

// getQueryString this ins't really pretty, but it works.
func (params *GetCurrentUserGuildsParams) getQueryString() string {
	separator := "?"
	query := ""

	if !params.Before.Empty() {
		query += separator + params.Before.String()
		separator = "&"
	}

	if !params.After.Empty() {
		query += separator + params.After.String()
		separator = "&"
	}

	if params.Limit > 0 {
		query += separator + strconv.Itoa(params.Limit)
	}

	return query
}

// ReqGetCurrentUserGuilds [GET]  Returns a list of partial guild objects the current user is a member of.
//                                Requires the guilds OAuth2 scope.
// Endpoint                       /users/@me/guilds
// Rate limiter                   /users TODO: is this correct?
// Discord documentation          https://discordapp.com/developers/docs/resources/user#get-current-user-guilds
// Reviewed                       2018-06-10
// Comment                        This endpoint returns 100 guilds by default, which is the maximum number of
//                                guilds a non-bot user can join. Therefore, pagination is not needed for
//                                integrations that need to get a list of users' guilds.
func GetCurrentUserGuilds(client httd.Getter, params *GetCurrentUserGuildsParams) (ret []*Guild, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint:    EndpointUserMyGuilds + params.getQueryString(),
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// ReqLeaveGuild [DELETE] Leave a guild. Returns a 204 empty response on success.
// Endpoint               /users/@me/guilds/{guild.id}
// Rate limiter           /users TODO: is this correct?
// Discord documentation  https://discordapp.com/developers/docs/resources/user#leave-guild
// Reviewed               2018-06-10
// Comment                -
func LeaveGuild(client httd.Deleter, guildID snowflake.ID) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint:    EndpointUserMyGuilds + "/" + guildID.String(),
	}
	resp, _, err := client.Delete(details)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusNoContent {
		msg := "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(http.StatusNoContent)
		err = errors.New(msg)
	}

	return
}

// ReqGetUserDMs [GET]    Returns a list of DM channel objects.
// Endpoint               /users/@me/channels
// Rate limiter           /users TODO: is this correct?
// Discord documentation  https://discordapp.com/developers/docs/resources/user#get-user-dms
// Reviewed               2018-06-10
// Comment                -
func GetUserDMs(client httd.Getter) (ret []*Channel, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint:    EndpointUserMyChannels,
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

type BodyUserCreateDM struct {
	RecipientID snowflake.ID `json:"recipient_id"`
}

// ReqGetUserDMs [POST]   Create a new DM channel with a user. Returns a DM channel object.
// Endpoint               /users/@me/channels
// Rate limiter           /users TODO: is this correct?
// Discord documentation  https://discordapp.com/developers/docs/resources/user#create-dm
// Reviewed               2018-06-10
// Comment                -
func CreateDM(client httd.Poster, recipientID snowflake.ID) (ret *Channel, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint:    EndpointUserMyChannels,
		JSONParams:  &BodyUserCreateDM{recipientID},
	}
	_, body, err := client.Post(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// BodyUserCreateGroupDM
// https://discordapp.com/developers/docs/resources/user#create-group-dm
type CreateGroupDMParams struct {
	AccessTokens []string                `json:"access_tokens"` // access tokens of users that have granted your app the gdm.join scope
	Nicks        map[snowflake.ID]string `json:"nicks"`         // userID => nickname
}

// ReqCreateGroupDM [POST]  Create a new group DM channel with multiple users. Returns a DM channel object.
// Endpoint                 /users/@me/channels
// Rate limiter             /users TODO: is this correct?
// Discord documentation    https://discordapp.com/developers/docs/resources/user#create-group-dm
// Reviewed                 2018-06-10
// Comment                  -
func CreateGroupDM(client httd.Poster, params *CreateGroupDMParams) (ret *Channel, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint:    EndpointUserMyChannels,
		JSONParams:  params,
	}
	_, body, err := client.Post(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// ReqCreateGroupDM [GET] Returns a list of connection objects. Requires the connections OAuth2 scope.
// Endpoint               /users/@me/connections
// Rate limiter           /users TODO: is this correct?
// Discord documentation  https://discordapp.com/developers/docs/resources/user#get-user-connections
// Reviewed               2018-06-10
// Comment                -
func GetUserConnections(client httd.Getter) (ret []*UserConnection, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitUsers(),
		Endpoint:    EndpointUserMyConnections,
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}
