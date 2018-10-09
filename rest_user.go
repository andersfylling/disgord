package disgord

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
)

func ratelimitUsers() string {
	return "u"
}

// [REST] Returns the user object of the requester's account. For OAuth2, this requires the identify scope, which
// will return the object without an email, and optionally the email scope, which returns the object with an email.
//  Method                  GET
//  Endpoint                /users/@me
//  Rate limiter            /users
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#get-current-user
//  Reviewed                2018-06-10
//  Comment                 -
func GetCurrentUser(client httd.Getter) (ret *User, err error) {
	_, body, err := client.Get(&httd.Request{
		Ratelimiter: ratelimitUsers(),
		Endpoint:    endpoint.UserMe(),
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// [REST] Returns a user object for a given user Snowflake.
//  Method                  GET
//  Endpoint                /users/{user.id}
//  Rate limiter            /users
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#get-user
//  Reviewed                2018-06-10
//  Comment                 -
func GetUser(client httd.Getter, id Snowflake) (ret *User, err error) {
	_, body, err := client.Get(&httd.Request{
		Ratelimiter: ratelimitUsers(),
		Endpoint:    endpoint.User(id),
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// JSON params for func ModifyCurrentUser
type ModifyCurrentUserParams struct {
	Username string `json:"username,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

// [REST] Modify the requester's user account settings. Returns a user object on success.
//  Method                  PATCH
//  Endpoint                /users/@me
//  Rate limiter            /users
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#modify-current-user
//  Reviewed                2018-06-10
//  Comment                 -
func ModifyCurrentUser(client httd.Patcher, params *ModifyCurrentUserParams) (ret *User, err error) {
	_, body, err := client.Patch(&httd.Request{
		Ratelimiter: ratelimitUsers(),
		Endpoint:    endpoint.UserMe(),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// JSON params for func GetCurrentUserGuilds
type GetCurrentUserGuildsParams struct {
	Before Snowflake `urlparam:"before,omitempty"`
	After  Snowflake `urlparam:"after,omitempty"`
	Limit  int       `urlparam:"limit,omitempty"`
}

// GetQueryString ...
func (params *GetCurrentUserGuildsParams) GetQueryString() string {
	separator := "?"
	query := ""

	if !params.Before.Empty() {
		query += separator + "before=" + params.Before.String()
		separator = "&"
	}

	if !params.After.Empty() {
		query += separator + "after=" + params.After.String()
		separator = "&"
	}

	if params.Limit > 0 {
		query += separator + "limit=" + strconv.Itoa(params.Limit)
	}

	return query
}

// [REST] Returns a list of partial guild objects the current user is a member of. Requires the guilds OAuth2 scope.
//  Method                  GET
//  Endpoint                /users/@me/guilds
//  Rate limiter            /users TODO: is this correct?
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#get-current-user-guilds
//  Reviewed                2018-06-10
//  Comment                 This endpoint. returns 100 guilds by default, which is the maximum number of
//                          guilds a non-bot user can join. Therefore, pagination is not needed for
//                          integrations that need to get a list of users' guilds.
func GetCurrentUserGuilds(client httd.Getter, params *GetCurrentUserGuildsParams) (ret []*Guild, err error) {
	_, body, err := client.Get(&httd.Request{
		Ratelimiter: ratelimitUsers(),
		Endpoint:    endpoint.UserMeGuilds() + params.GetQueryString(),
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// [REST] Leave a guild. Returns a 204 empty response on success.
//  Method                  DELETE
//  Endpoint                /users/@me/guilds/{guild.id}
//  Rate limiter            /users TODO: is this correct?
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#leave-guild
//  Reviewed                2018-06-10
//  Comment                 -
func LeaveGuild(client httd.Deleter, id Snowflake) (err error) {
	resp, _, err := client.Delete(&httd.Request{
		Ratelimiter: ratelimitUsers(),
		Endpoint:    endpoint.UserMeGuild(id),
	})
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusNoContent {
		msg := "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(http.StatusNoContent)
		err = errors.New(msg)
	}

	return
}

// [REST] Returns a list of DM channel objects.
//  Method                  GET
//  Endpoint                /users/@me/channels
//  Rate limiter            /users TODO: is this correct?
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#get-user-dms
//  Reviewed                2018-06-10
//  Comment                 -
func GetUserDMs(client httd.Getter) (ret []*Channel, err error) {
	_, body, err := client.Get(&httd.Request{
		Ratelimiter: ratelimitUsers(),
		Endpoint:    endpoint.UserMeChannels(),
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// JSON param for func CreateDM
type BodyUserCreateDM struct {
	RecipientID Snowflake `json:"recipient_id"`
}

// [REST] Create a new DM channel with a user. Returns a DM channel object.
//  Method                  POST
//  Endpoint                /users/@me/channels
//  Rate limiter            /users TODO: is this correct?
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#create-dm
//  Reviewed                2018-06-10
//  Comment                 -
// TODO: review
func CreateDM(client httd.Poster, recipientID Snowflake) (ret *Channel, err error) {
	_, body, err := client.Post(&httd.Request{
		Ratelimiter: ratelimitUsers(),
		Endpoint:    endpoint.UserMeChannels(),
		Body:        &BodyUserCreateDM{recipientID},
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// JSON params for func CreateGroupDM https://discordapp.com/developers/docs/resources/user#create-group-dm
type CreateGroupDMParams struct {
	AccessTokens []string             `json:"access_tokens"` // access tokens of users that have granted your app the gdm.join scope
	Nicks        map[Snowflake]string `json:"nicks"`         // map[userID] = nickname
}

// [REST] Create a new group DM channel with multiple users. Returns a DM channel object.
//  Method                  POST
//  Endpoint                /users/@me/channels
//  Rate limiter            /users TODO: is this correct?
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#create-group-dm
//  Reviewed                2018-06-10
//  Comment                 -
func CreateGroupDM(client httd.Poster, params *CreateGroupDMParams) (ret *Channel, err error) {
	_, body, err := client.Post(&httd.Request{
		Ratelimiter: ratelimitUsers(),
		Endpoint:    endpoint.UserMeChannels(),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// [REST] Returns a list of connection objects. Requires the connections OAuth2 scope.
//  Method                  GET
//  Endpoint                /users/@me/connections
//  Rate limiter            /users TODO: is this correct?
//  Discord documentation   https://discordapp.com/developers/docs/resources/user#get-user-connections
//  Reviewed                2018-06-10
//  Comment                 -
func GetUserConnections(client httd.Getter) (ret []*UserConnection, err error) {
	_, body, err := client.Get(&httd.Request{
		Ratelimiter: ratelimitUsers(),
		Endpoint:    endpoint.UserMeConnections(),
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}
