package resource

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/andersfylling/snowflake"
	"github.com/sirupsen/logrus"
)

// UserMessager Methods required to create a new DM (or use an existing one) and send a DM.
type UserMessager interface { // TODO: wtf?
	//CreateAndSendDM(recipientID snowflake.ID, msg *Message) error // hmmm...
}

type UserInterface interface {
	Mention() string
	MentionNickname() string
	String() string
	//
	//// Update internal structure
	//Update(*User) error
	//Clear()
	//
	//// Send a direct message to this user
	//SendMessage(UserMessager, string) (snowflake.ID, snowflake.ID, error)
}

type User struct {
	ID            snowflake.ID `json:"id,omitempty"`
	Username      string       `json:"username,omitempty"`
	Discriminator string       `json:"discriminator,omitempty"`
	Email         string       `json:"email,omitempty"`
	Avatar        *string      `json:"avatar"`
	Token         string       `json:"token,omitempty"`
	Verified      bool         `json:"verified,omitempty"`
	MFAEnabled    bool         `json:"mfa_enabled,omitempty"`
	Bot           bool         `json:"bot,omitempty"`

	sync.RWMutex `json:"-"`
}

func NewUser() *User {
	return &User{}
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

func (u *User) SendMessage(client UserMessager, msg string) (channelID snowflake.ID, messageID snowflake.ID, err error) {
	return snowflake.NewID(0), snowflake.NewID(0), errors.New("not implemented")
}

func (u *User) Replicate(user *User) error {
	if user == nil {
		return errors.New("cannot copy nil object")
	}

	if u == user {
		return errors.New("cannot copy itself, makes no sense")
	}

	user.RLock()
	u.Lock()

	// deep copy, without changing the avatar pointer address
	var avatarAddress *string
	if u.Avatar != nil {
		avatarAddress = u.Avatar
	} else {
		avatar := ""
		avatarAddress = &avatar
	}

	uMutex := u.RWMutex
	*u = *user                   // copy all the fields
	u.RWMutex = uMutex           // make sure the mutex isn't deleted
	u.Avatar = avatarAddress     // point to a different location than user.Avatar
	*(u.Avatar) = *(user.Avatar) // copy the Base64 image

	u.Unlock()
	user.RUnlock()

	return nil
}

type UserPresence struct {
	User    *User          `json:"user"`
	Roles   []snowflake.ID `json:"roles"`
	Game    *UserActivity  `json:"activty"`
	GuildID snowflake.ID   `json:"guild_id"`
	Nick    *string        `json:"nick"`
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

type UserEndpointInterface interface {
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

type UserEndpoint struct{}

// GetCurrentUser [GET] Returns the user object of the requester's account. For OAuth2, this requires
//                      the identify scope, which will return the object without an email, and optionally
//                      the email scope, which returns the object with an email.
func (e *UserEndpoint) GetCurrentUser() string {
	return "/user/@me"
}

// GetUser [GET] Returns a user object for a given user ID.
func (e *UserEndpoint) GetUser(id snowflake.ID) string {
	return "/users/" + id.String()
}

// ModifyCurrentUser [PATCH, JSON] Modify the requester's user account settings. Returns a user object on success.
func (e *UserEndpoint) ModifyCurrentUser() string {
	return e.GetCurrentUser()
}

// GetCurrentUserGuilds [GET] Returns a list of partial guild objects the current user is a member of.
//                            Requires the guilds OAuth2 scope.
func (e *UserEndpoint) GetCurrentUserGuilds(params interface{}) string {
	logrus.WithFields(
		logrus.Fields{
			"package": "disgord.schema",
			"Func":    "GetCurrentUserGuilds(params interface{}) string",
		}).Warnln("Params not parsed!")

	return "/users/@me/guilds"
}

// LeaveGuild [DELETE] Leave a guild. Returns a 204 empty response on success.
func (e *UserEndpoint) LeaveGuild(id snowflake.ID) string {
	return "/users/@me/guilds/" + id.String()
}

// GetUserDMs [GET] Returns a list of DM channel objects.
func (e *UserEndpoint) GetUserDMs() string {
	return "/users/@me/channels"
}

// CreateDM [POST, JSON] Create a new DM channel with a user. Returns a DM channel object.
func (e *UserEndpoint) CreateDM() string {
	return e.GetUserDMs()
}

// CreateGroupDM [POST, JSON] Create a new group DM channel with multiple users. Returns a DM channel object.
func (e *UserEndpoint) CreateGroupDM() string {
	return e.CreateDM()
}

// GetUserConnections [GET] Returns a list of connection objects. Requires the connections OAuth2 scope.
func (e *UserEndpoint) GetUserConnections() string {
	return "/users/@me/connections"
}

// Connection The connection object that the user has attached.
// https://discordapp.com/developers/docs/resources/user#avatar-data
// WARNING! Due to dependency issues, the Integrations (array) refers to Integration IDs only!
//          It breaks the lib pattern, but there's nothing I can do. To retrieve the Integration
//          Object, use *disgord.Client.Integration(id) (*Integration, error)
type UserConnection struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Revoked bool   `json:"revoked"`

	// Since this does not hold real guild.Integration objects we need to empty
	// the integration slice, it's misguiding. But at least the output is "correct".
	// the UnmarshalJSON method is a hack for input
	Integrations []snowflake.ID `json:"-"`
}

func (conn *UserConnection) UnmarshalJSON(b []byte) error {
	mock := struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		Revoked      bool   `json:"revoked"`
		Integrations []struct {
			ID snowflake.ID `json:"id"`
		} `json:"integrations"`
	}{}

	err := json.Unmarshal(b, &mock)
	if err != nil {
		return err
	}

	conn.ID = mock.ID
	conn.Name = mock.Name
	conn.Type = mock.Type
	conn.Revoked = mock.Revoked

	// empty integration slice
	//conn.Integrations = conn.Integrations[:0]

	// set new slice size
	conn.Integrations = make([]snowflake.ID, len(mock.Integrations))

	// add new data
	for index, id := range mock.Integrations {
		conn.Integrations[index] = id.ID
	}

	return nil
}

// TODO: https://discordapp.com/developers/docs/topics/gateway#activity-object-activity-structure
type UserActivity struct{}
