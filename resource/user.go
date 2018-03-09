package resource

import (
	"encoding/json"
	"errors"
	"sync"

	"fmt"

	"github.com/andersfylling/disgord/request"
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

// ---------

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

func (u *User) SendMessage(requester request.DiscordRequester, msg *Message) (channelID snowflake.ID, messageID snowflake.ID, err error) {
	return snowflake.NewID(0), snowflake.NewID(0), errors.New("not implemented")
}

func (u *User) SendMessageStr(requester request.DiscordRequester, msg string) (channelID snowflake.ID, messageID snowflake.ID, err error) {
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

// https://discordapp.com/developers/docs/resources/user#connection-object
// TODO
type UserConnection struct {
}

// ----------

// GetUser [GET] Returns a user object for a given user ID.
func ReqUser(requester request.DiscordGetter, id snowflake.ID) (*User, error) {
	var endpoint string = EndpointUser
	path := EndpointUser + id.String()

	result := NewUser()
	_, err := requester.Get(endpoint, path, result)

	return result, err
}

func ReqMyself(requester request.DiscordGetter) (*User, error) {
	endpoint := EndpointUser
	path := EndpointUserMyself

	result := NewUser()
	_, err := requester.Get(endpoint, path, result)

	return result, err
}

// RequestMyGuilds [GET] Returns a list of partial guild objects the current user is a member of.
//                       Requires the guilds OAuth2 scope.
func ReqMyGuilds(requester request.DiscordGetter) ([]*Guild, error) {
	endpoint := EndpointUser
	path := EndpointUserMyGuilds

	var result []*Guild
	_, err := requester.Get(endpoint, path, result)

	return result, err
}

// ReqMyDMs [GET] Returns a list of DM channel objects.
func ReqMyDMs(requester request.DiscordGetter) ([]*Channel, error) {
	endpoint := EndpointUser
	path := EndpointUserMyChannels

	var result []*Channel
	_, err := requester.Get(endpoint, path, result)

	return result, err
}

// ReqLeaveGuild [DELETE] Leave a guild.
// 						  Returns a 204 empty response on success.
func ReqLeaveGuild(requester request.DiscordDeleter, id snowflake.ID) error {
	endpoint := EndpointUser
	path := EndpointUserMyGuilds + "/" + id.String()

	_, err := requester.Delete(endpoint, path)

	return err
}

type BodyUserCreateDM struct {
	RecipientID snowflake.ID `json:"recipient_id"`
}

// ReqCreateDM [POST, JSON] Create a new DM channel with a user. Returns a DM channel object.
func ReqCreateDM(requester request.DiscordPoster, user *User) (*Channel, error) {
	endpoint := EndpointUser
	path := EndpointUserMyChannels
	params := BodyUserCreateDM{
		RecipientID: user.ID,
	}

	var result *Channel
	_, err := requester.Post(endpoint, path, result, &params)

	return result, err
}

// BodyUserCreateGroupDM
// https://discordapp.com/developers/docs/resources/user#create-group-dm
type BodyUserCreateGroupDM struct {
	AccessTokens []string                `json:"access_tokens"` // access tokens of users that have granted your app the gdm.join scope
	Nicks        map[snowflake.ID]string `json:"nicks"`         // userID => nickname
}

// ReqCreateGroupDM [POST, JSON] Create a new group DM channel with multiple users. Returns a DM channel object.
func ReqCreateGroupDM(requester request.DiscordPoster, user *User) (*Channel, error) {
	fmt.Println("ReqCreateGroupDM HAS NOT YET BEEN IMPLEMENTED!")
	return nil, errors.New("not implemented")
	endpoint := EndpointUser
	path := EndpointUserMyChannels
	params := BodyUserCreateGroupDM{}

	var result *Channel
	_, err := requester.Post(endpoint, path, result, &params)

	return result, err
}

// ReqMyConnections [GET] Returns a list of connection objects. Requires the connections OAuth2 scope.
func ReqMyConnections(requester request.DiscordGetter) ([]*UserConnection, error) {
	endpoint := EndpointUser
	path := EndpointUserMyConnections

	var result []*UserConnection
	_, err := requester.Get(endpoint, path, result)

	return result, err
}

// --------

// TODO: https://discordapp.com/developers/docs/topics/gateway#activity-object-activity-structure
type UserActivity struct{}
