package user

import "github.com/andersfylling/snowflake"
import "github.com/sirupsen/logrus"

type EndpointInterface interface {
	GetCurrentUser() string
	GetUser(id snowflake.ID) string
	ModifyCurrentUser() string
	GetCurrentUserGuilds(params interface{}) string
	LeaveGuild(id snowflake.ID) string
	GetUserDMs() string
	CreateDM() string
	CreateGroupDM() string
	GetUserConnections() string
}

type Endpoint struct{}

// GetCurrentUser [GET] Returns the user object of the requester's account. For OAuth2, this requires
//                      the identify scope, which will return the object without an email, and optionally
//                      the email scope, which returns the object with an email.
func (e *Endpoint) GetCurrentUser() string {
	return "/user/@me"
}

// GetUser [GET] Returns a user object for a given user ID.
func (e *Endpoint) GetUser(id snowflake.ID) string {
	return "/users/" + id.String()
}

// ModifyCurrentUser [PATCH, JSON] Modify the requester's user account settings. Returns a user object on success.
func (e *Endpoint) ModifyCurrentUser() string {
	return e.GetCurrentUser()
}

// GetCurrentUserGuilds [GET] Returns a list of partial guild objects the current user is a member of.
//                            Requires the guilds OAuth2 scope.
func (e *Endpoint) GetCurrentUserGuilds(params interface{}) string {
	logrus.WithFields(
		logrus.Fields{
			"package": "disgord.user",
			"Func":    "GetCurrentUserGuilds(params interface{}) string",
		}).Warnln("Params not parsed!")

	return "/users/@me/guilds"
}

// LeaveGuild [DELETE] Leave a guild. Returns a 204 empty response on success.
func (e *Endpoint) LeaveGuild(id snowflake.ID) string {
	return "/users/@me/guilds/" + id.String()
}

// GetUserDMs [GET] Returns a list of DM channel objects.
func (e *Endpoint) GetUserDMs() string {
	return "/users/@me/channels"
}

// CreateDM [POST, JSON] Create a new DM channel with a user. Returns a DM channel object.
func (e *Endpoint) CreateDM() string {
	return e.GetUserDMs()
}

// CreateGroupDM [POST, JSON] Create a new group DM channel with multiple users. Returns a DM channel object.
func (e *Endpoint) CreateGroupDM() string {
	return e.CreateDM()
}

// GetUserConnections [GET] Returns a list of connection objects. Requires the connections OAuth2 scope.
func (e *Endpoint) GetUserConnections() string {
	return "/users/@me/connections"
}

var _ EndpointInterface = (*Endpoint)(nil)
