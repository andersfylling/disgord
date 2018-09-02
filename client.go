package disgord

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"errors"

	"github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/disgord/rest"
	"github.com/andersfylling/disgord/rest/httd"
	"github.com/andersfylling/disgord/state"
	"github.com/andersfylling/disgordws"
	. "github.com/andersfylling/snowflake"
	"github.com/sirupsen/logrus"
)

// Session the discord api is split in two. socket for keeping the client up to date, and http api for requests.
type Session interface {
	// main modules
	//

	// Request For interacting with Discord. Sending messages, creating channels, guilds, etc.
	// To read object state such as guilds, State() should be used in stead. However some data
	// might not exist in the state. If so it should be requested. Note that this only holds http
	// CRUD operation and not the actual rest endpoints for discord (See Rest()).
	Req() httd.Requester

	// todo
	//Rest()

	// Event let's developers listen for specific events, event groups, or every event as one listener.
	// Supports both channels and callbacks
	Evt() EvtDispatcher

	// State reflects the latest changes received from Discord gateway.
	// Should be used instead of requesting objects.
	State() state.Cacher

	// RateLimiter the ratelimiter for the discord REST API
	RateLimiter() httd.RateLimiter

	// Discord Gateway, web socket
	//
	Connect() error
	Disconnect() error

	// module wrappers
	//

	// event callbacks
	AddListener(evtName string, callback interface{})
	AddListenerOnce(evtName string, callback interface{})

	// state/caching module
	// checks the cache first, otherwise do a http request

	// all discord REST functions
	// TODO: support caching for each
	// Audit-log
	GetGuildAuditLogs(guildID Snowflake, params *rest.AuditLogParams) (log *resource.AuditLog, err error)
	// Channel
	GetChannel(id Snowflake) (ret *resource.Channel, err error)
	ModifyChannel(changes *rest.ModifyChannelParams) (ret *resource.Channel, err error)
	DeleteChannel(id Snowflake) (err error)
	EditChannelPermissions(chanID, overwriteID Snowflake, params *rest.EditChannelPermissionsParams) (err error)
	GetChannelInvites(id Snowflake) (ret []*resource.Invite, err error)
	CreateChannelInvites(id Snowflake, params *rest.CreateChannelInvitesParams) (ret *resource.Invite, err error)
	DeleteChannelPermission(channelID, overwriteID Snowflake) (err error)
	TriggerTypingIndicator(channelID Snowflake) (err error)
	GetPinnedMessages(channelID Snowflake) (ret []*resource.Message, err error)
	AddPinnedChannelMessage(channelID, msgID Snowflake) (err error)
	DeletePinnedChannelMessage(channelID, msgID Snowflake) (err error)
	GroupDMAddRecipient(channelID, userID Snowflake, params *rest.GroupDMAddRecipientParams) (err error)
	GroupDMRemoveRecipient(channelID, userID Snowflake) (err error)
	GetChannelMessages(channelID Snowflake, params rest.URLParameters) (ret []*resource.Message, err error)
	GetChannelMessage(channelID, messageID Snowflake) (ret *resource.Message, err error)
	CreateChannelMessage(channelID Snowflake, params *rest.CreateChannelMessageParams) (ret *resource.Message, err error)
	EditMessage(chanID, msgID Snowflake, params *rest.EditMessageParams) (ret *resource.Message, err error)
	DeleteMessage(channelID, msgID Snowflake) (err error)
	BulkDeleteMessages(chanID Snowflake, params *rest.BulkDeleteMessagesParams) (err error)
	CreateReaction(channelID, messageID Snowflake, emoji interface{}) (ret *resource.Reaction, err error)
	DeleteOwnReaction(channelID, messageID Snowflake, emoji interface{}) (err error)
	DeleteUserReaction(channelID, messageID, userID Snowflake, emoji interface{}) (err error)
	GetReaction(channelID, messageID Snowflake, emoji interface{}, params rest.URLParameters) (ret []*resource.User, err error)
	DeleteAllReactions(channelID, messageID Snowflake) (err error)
	// Emoji
	GetGuildEmojis(id Snowflake) (ret []*resource.Emoji, err error)
	GetGuildEmoji(guildID, emojiID Snowflake) (ret *resource.Emoji, err error)
	CreateGuildEmoji(guildID Snowflake) (ret *resource.Emoji, err error)
	ModifyGuildEmoji(guildID, emojiID Snowflake) (ret *resource.Emoji, err error)
	DeleteGuildEmoji(guildID, emojiID Snowflake) (err error)
	// Guild
	CreateGuild(params *rest.CreateGuildParams) (ret *resource.Guild, err error)
	GetGuild(id Snowflake) (ret *resource.Guild, err error)
	ModifyGuild(id Snowflake, params *rest.ModifyGuildParams) (ret *resource.Guild, err error)
	DeleteGuild(id Snowflake) (err error)
	GetGuildChannels(id Snowflake) (ret []*resource.Channel, err error)
	CreateGuildChannel(id Snowflake, params *rest.CreateGuildChannelParams) (ret *resource.Channel, err error)
	GetGuildMember(guildID, userID Snowflake) (ret *resource.Member, err error)
	GetGuildMembers(guildID, after Snowflake, limit int) (ret []*resource.Member, err error)
	AddGuildMember(guildID, userID Snowflake, params *rest.AddGuildMemberParams) (ret *resource.Member, err error)
	ModifyGuildMember(guildID, userID Snowflake, params *rest.ModifyGuildMemberParams) (err error)
	ModifyCurrentUserNick(id Snowflake, params *rest.ModifyCurrentUserNickParams) (nick string, err error)
	AddGuildMemberRole(guildID, userID, roleID Snowflake) (err error)
	RemoveGuildMemberRole(guildID, userID, roleID Snowflake) (err error)
	RemoveGuildMember(guildID, userID Snowflake) (err error)
	GetGuildBans(id Snowflake) (ret []*resource.Ban, err error)
	GetGuildBan(guildID, userID Snowflake) (ret *resource.Ban, err error)
	CreateGuildBan(guildID, userID Snowflake, params *rest.CreateGuildBanParams) (err error)
	RemoveGuildBan(guildID, userID Snowflake) (err error)
	GetGuildRoles(guildID Snowflake) (ret []*resource.Role, err error)
	CreateGuildRole(id Snowflake, params *rest.CreateGuildRoleParams) (ret *resource.Role, err error)
	ModifyGuildRolePositions(guildID Snowflake, params *rest.ModifyGuildRolePositionsParams) (ret []*resource.Role, err error)
	ModifyGuildRole(guildID, roleID Snowflake, params *rest.ModifyGuildRoleParams) (ret []*resource.Role, err error)
	DeleteGuildRole(guildID, roleID Snowflake) (err error)
	GetGuildPruneCount(id Snowflake, params *rest.GuildPruneParams) (ret *resource.GuildPruneCount, err error)
	BeginGuildPrune(id Snowflake, params *rest.GuildPruneParams) (ret *resource.GuildPruneCount, err error)
	GetGuildVoiceRegions(id Snowflake) (ret []*resource.VoiceRegion, err error)
	GetGuildInvites(id Snowflake) (ret []*resource.Invite, err error)
	GetGuildIntegrations(id Snowflake) (ret []*resource.Integration, err error)
	CreateGuildIntegration(guildID Snowflake, params *rest.CreateGuildIntegrationParams) (err error)
	ModifyGuildIntegration(guildID, integrationID Snowflake, params *rest.ModifyGuildIntegrationParams) (err error)
	DeleteGuildIntegration(guildID, integrationID Snowflake) (err error)
	SyncGuildIntegration(guildID, integrationID Snowflake) (err error)
	GetGuildEmbed(guildID Snowflake) (ret *resource.GuildEmbed, err error)
	ModifyGuildEmbed(guildID Snowflake, params *resource.GuildEmbed) (ret *resource.GuildEmbed, err error)
	GetGuildVanityURL(guildID Snowflake) (ret *resource.PartialInvite, err error)
	// Invite
	GetInvite(inviteCode string, withCounts bool) (invite *resource.Invite, err error)
	DeleteInvite(inviteCode string) (invite *resource.Invite, err error)
	// User
	GetCurrentUser() (ret *resource.User, err error)
	GetUser(id Snowflake) (ret *resource.User, err error)
	ModifyCurrentUser(params *rest.ModifyCurrentUserParams) (ret *resource.User, err error)
	GetCurrentUserGuilds(params *rest.GetCurrentUserGuildsParams) (ret []*resource.Guild, err error)
	LeaveGuild(id Snowflake) (err error)
	GetUserDMs() (ret []*resource.Channel, err error)
	CreateDM(recipientID Snowflake) (ret *resource.Channel, err error)
	CreateGroupDM(params *rest.CreateGroupDMParams) (ret *resource.Channel, err error)
	GetUserConnections() (ret []*resource.UserConnection, err error)
	// Voice
	GetVoiceRegions() (ret []*resource.VoiceRegion, err error)
	// Webhook
	CreateWebhook(channelID Snowflake, params *rest.CreateWebhookParams) (ret *resource.Webhook, err error)
	GetChannelWebhooks(channelID Snowflake) (ret []*resource.Webhook, err error)
	GetGuildWebhooks(guildID Snowflake) (ret []*resource.Webhook, err error)
	GetWebhook(id Snowflake) (ret *resource.Webhook, err error)
	GetWebhookWithToken(id Snowflake, token string) (ret *resource.Webhook, err error)
	ModifyWebhook(newWebhook *resource.Webhook) (ret *resource.Webhook, err error)
	ModifyWebhookWithToken(newWebhook *resource.Webhook) (ret *resource.Webhook, err error)
	DeleteWebhook(webhookID Snowflake) (err error)
	DeleteWebhookWithToken(id Snowflake, token string) (err error)
	ExecuteWebhook(params *rest.ExecuteWebhookParams, wait bool, URLSuffix string) (err error)
	ExecuteSlackWebhook(params *rest.ExecuteWebhookParams, wait bool) (err error)
	ExecuteGitHubWebhook(params *rest.ExecuteWebhookParams, wait bool) (err error)
	// Custom
	SendMsg(channelID Snowflake, message *resource.Message) (msg *resource.Message, err error)
	SendMsgString(channelID Snowflake, content string) (msg *resource.Message, err error)

	// same as above. Except these returns a channel
	// WARNING: none below should be assumed to be working.
	GuildChan(guildID Snowflake) <-chan *resource.Guild
	ChannelChan(channelID Snowflake) <-chan *resource.Channel
	ChannelsChan(guildID Snowflake) <-chan map[Snowflake]*resource.Channel
	MsgChan(msgID Snowflake) <-chan *resource.Message
	UserChan(userID Snowflake) <-chan *UserChan
	MemberChan(guildID, userID Snowflake) <-chan *resource.Member
	MembersChan(guildID Snowflake) <-chan map[Snowflake]*resource.Member
}

type Config struct {
	Token      string
	HTTPClient *http.Client

	APIVersion  int    // eg. version 6. 0 defaults to lowest supported api version
	APIEncoding string // eg. json, use const. defaults to json

	CancelRequestWhenRateLimited bool

	LoadAllMembers   bool
	LoadAllChannels  bool
	LoadAllRoles     bool
	LoadAllPresences bool

	Debug bool
}

// NewClient creates a new default disgord instance
func NewClient(conf *Config) (*Client, error) {

	// ensure valid api version
	if conf.APIVersion == 0 {
		conf.APIVersion = 6 // the current discord API, for now v6
	}
	switch conf.APIVersion { // todo: simplify
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		fallthrough
	case 4:
		fallthrough
	case 5:
		return nil, errors.New("outdated API version")
	case 6: // supported
	default:
		return nil, errors.New("Discord API version is not yet supported")
	}

	if conf.HTTPClient == nil {
		// http client configuration
		conf.HTTPClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	// Use disgordws to keep the socket connection going
	// default communication encoding to json
	if conf.APIEncoding == "" {
		conf.APIEncoding = JSONEncoding
	}
	dws, err := disgordws.NewClient(&disgordws.Config{
		// user settings
		Token:      conf.Token,
		HTTPClient: conf.HTTPClient,
		Debug:      conf.Debug,

		// lib specific
		DAPIVersion:  conf.APIVersion,
		DAPIEncoding: conf.APIEncoding,
	})
	if err != nil {
		return nil, err
	}

	// request client
	reqConf := &httd.Config{
		APIVersion:                   conf.APIVersion,
		BotToken:                     conf.Token,
		UserAgentSourceURL:           GitHubURL,
		UserAgentVersion:             Version,
		HTTPClient:                   conf.HTTPClient,
		CancelRequestWhenRateLimited: conf.CancelRequestWhenRateLimited,
	}
	reqClient := httd.NewClient(reqConf)

	// event dispatcher
	evtDispatcher := NewDispatch()

	// create a disgord client/instance/session
	c := &Client{
		httpClient:    conf.HTTPClient,
		ws:            dws,
		socketEvtChan: dws.GetEventChannel(),
		token:         conf.Token,
		evtDispatch:   evtDispatcher,
		state:         state.NewCache(),
		req:           reqClient,
	}

	return c, nil
}

func NewClientMustCompile(conf *Config) *Client {
	client, err := NewClient(conf)
	if err != nil {
		panic(err)
	}

	return client
}

func NewSession(conf *Config) (Session, error) {
	return NewClient(conf)
}

func NewSessionMustCompile(conf *Config) Session {
	return NewClientMustCompile(conf)
}

type Client struct {
	sync.RWMutex

	token string

	ws            *disgordws.Client
	socketEvtChan <-chan disgordws.EventInterface

	// register listeners for events
	evtDispatch *Dispatch

	// cancelRequestWhenRateLimited by default the client waits until either the HTTPClient.timeout or
	// the rate limit ends before closing a request channel. If activated, in stead, requests will
	// instantly be denied, and the channel closed.
	cancelRequestWhenRateLimited bool

	// discord http api
	req *httd.Client

	httpClient *http.Client

	// cache
	state *state.Cache
}

func (c *Client) logInfo(msg string) {
	logrus.WithFields(logrus.Fields{
		"lib": c.ws.String(),
	}).Info(msg)
}

func (c *Client) logErr(msg string) {
	logrus.WithFields(logrus.Fields{
		"lib": c.ws.String(),
	}).Error(msg)
}

func (c *Client) String() string {
	return c.ws.String()
}

func (c *Client) RateLimiter() httd.RateLimiter {
	return c.req.RateLimiter()
}

// Connect establishes a websocket connection to the discord API
func (c *Client) Connect() (err error) {
	c.logInfo("Connecting to discord Gateway")
	err = c.ws.Connect()
	if err != nil {
		c.logErr(err.Error())
		return
	}
	c.logInfo("Connected")

	// setup event observer
	go c.eventHandler()

	return nil
}

// Disconnect closes the discord websocket connection
func (c *Client) Disconnect() (err error) {
	fmt.Println()
	c.logInfo("Closing Discord gateway connection")
	err = c.ws.Disconnect()
	if err != nil {
		c.logErr(err.Error())
		return
	}
	c.logInfo("Disconnected")

	return nil
}

func (c *Client) Req() httd.Requester {
	return c.req
}

func (c *Client) Evt() EvtDispatcher {
	return c.evtDispatch
}

func (c *Client) State() state.Cacher {
	return c.state
}

func (c *Client) AddListener(evtName string, listener interface{}) {
	c.evtDispatch.AddHandler(evtName, listener)
}

// AddListenerOnce not implemented. Do not use.
func (c *Client) AddListenerOnce(evtName string, listener interface{}) {
	c.evtDispatch.AddHandlerOnce(evtName, listener)
}

// Audit-log
func (c *Client) GetGuildAuditLogs(guildID Snowflake, params *rest.AuditLogParams) (log *resource.AuditLog, err error) {
	log, err = rest.GuildAuditLogs(c.req, guildID, params)
	return
}

// Channel
func (c *Client) GetChannel(id Snowflake) (ret *resource.Channel, err error) {
	ret, err = rest.GetChannel(c.req, id)
	return
}
func (c *Client) ModifyChannel(changes *rest.ModifyChannelParams) (ret *resource.Channel, err error) {
	ret, err = rest.ModifyChannel(c.req, changes)
	return
}
func (c *Client) DeleteChannel(id Snowflake) (err error) {
	err = rest.DeleteChannel(c.req, id)
	return
}
func (c *Client) EditChannelPermissions(chanID, overwriteID Snowflake, params *rest.EditChannelPermissionsParams) (err error) {
	err = rest.EditChannelPermissions(c.req, chanID, overwriteID, params)
	return
}
func (c *Client) GetChannelInvites(id Snowflake) (ret []*resource.Invite, err error) {
	ret, err = rest.GetChannelInvites(c.req, id)
	return
}
func (c *Client) CreateChannelInvites(id Snowflake, params *rest.CreateChannelInvitesParams) (ret *resource.Invite, err error) {
	ret, err = rest.CreateChannelInvites(c.req, id, params)
	return
}
func (c *Client) DeleteChannelPermission(channelID, overwriteID Snowflake) (err error) {
	err = rest.DeleteChannelPermission(c.req, channelID, overwriteID)
	return
}
func (c *Client) TriggerTypingIndicator(channelID Snowflake) (err error) {
	err = rest.TriggerTypingIndicator(c.req, channelID)
	return
}
func (c *Client) GetPinnedMessages(channelID Snowflake) (ret []*resource.Message, err error) {
	ret, err = rest.GetPinnedMessages(c.req, channelID)
	return
}
func (c *Client) AddPinnedChannelMessage(channelID, msgID Snowflake) (err error) {
	err = rest.AddPinnedChannelMessage(c.req, channelID, msgID)
	return
}
func (c *Client) DeletePinnedChannelMessage(channelID, msgID Snowflake) (err error) {
	err = rest.DeletePinnedChannelMessage(c.req, channelID, msgID)
	return
}
func (c *Client) GroupDMAddRecipient(channelID, userID Snowflake, params *rest.GroupDMAddRecipientParams) (err error) {
	err = rest.GroupDMAddRecipient(c.req, channelID, userID, params)
	return
}
func (c *Client) GroupDMRemoveRecipient(channelID, userID Snowflake) (err error) {
	err = rest.GroupDMRemoveRecipient(c.req, channelID, userID)
	return
}
func (c *Client) GetChannelMessages(channelID Snowflake, params rest.URLParameters) (ret []*resource.Message, err error) {
	ret, err = rest.GetChannelMessages(c.req, channelID, params)
	return
}
func (c *Client) GetChannelMessage(channelID, messageID Snowflake) (ret *resource.Message, err error) {
	ret, err = rest.GetChannelMessage(c.req, channelID, messageID)
	return
}
func (c *Client) CreateChannelMessage(channelID Snowflake, params *rest.CreateChannelMessageParams) (ret *resource.Message, err error) {
	ret, err = rest.CreateChannelMessage(c.req, channelID, params)
	return
}
func (c *Client) EditMessage(chanID, msgID Snowflake, params *rest.EditMessageParams) (ret *resource.Message, err error) {
	ret, err = rest.EditMessage(c.req, chanID, msgID, params)
	return
}
func (c *Client) DeleteMessage(channelID, msgID Snowflake) (err error) {
	err = rest.DeleteMessage(c.req, channelID, msgID)
	return
}
func (c *Client) BulkDeleteMessages(chanID Snowflake, params *rest.BulkDeleteMessagesParams) (err error) {
	err = rest.BulkDeleteMessages(c.req, chanID, params)
	return
}
func (c *Client) CreateReaction(channelID, messageID Snowflake, emoji interface{}) (ret *resource.Reaction, err error) {
	ret, err = rest.CreateReaction(c.req, channelID, messageID, emoji)
	return
}
func (c *Client) DeleteOwnReaction(channelID, messageID Snowflake, emoji interface{}) (err error) {
	err = rest.DeleteOwnReaction(c.req, channelID, messageID, emoji)
	return
}
func (c *Client) DeleteUserReaction(channelID, messageID, userID Snowflake, emoji interface{}) (err error) {
	err = rest.DeleteUserReaction(c.req, channelID, messageID, userID, emoji)
	return
}
func (c *Client) GetReaction(channelID, messageID Snowflake, emoji interface{}, params rest.URLParameters) (ret []*resource.User, err error) {
	ret, err = rest.GetReaction(c.req, channelID, messageID, emoji, params)
	return
}
func (c *Client) DeleteAllReactions(channelID, messageID Snowflake) (err error) {
	err = rest.DeleteAllReactions(c.req, channelID, messageID)
	return
}

// Emoji
func (c *Client) GetGuildEmojis(id Snowflake) (ret []*resource.Emoji, err error) {
	ret, err = rest.ListGuildEmojis(c.req, id)
	return
}
func (c *Client) GetGuildEmoji(guildID, emojiID Snowflake) (ret *resource.Emoji, err error) {
	ret, err = rest.GetGuildEmoji(c.req, guildID, emojiID)
	return
}
func (c *Client) CreateGuildEmoji(guildID Snowflake) (ret *resource.Emoji, err error) {
	ret, err = rest.CreateGuildEmoji(c.req, guildID)
	return
}
func (c *Client) ModifyGuildEmoji(guildID, emojiID Snowflake) (ret *resource.Emoji, err error) {
	ret, err = rest.ModifyGuildEmoji(c.req, guildID, emojiID)
	return
}
func (c *Client) DeleteGuildEmoji(guildID, emojiID Snowflake) (err error) {
	err = rest.DeleteGuildEmoji(c.req, guildID, emojiID)
	return
}

// Guild
func (c *Client) CreateGuild(params *rest.CreateGuildParams) (ret *resource.Guild, err error) {
	ret, err = rest.CreateGuild(c.req, params)
	return
}
func (c *Client) GetGuild(id Snowflake) (ret *resource.Guild, err error) {
	ret, err = rest.GetGuild(c.req, id)
	return
}
func (c *Client) ModifyGuild(id Snowflake, params *rest.ModifyGuildParams) (ret *resource.Guild, err error) {
	ret, err = rest.ModifyGuild(c.req, id, params)
	return
}
func (c *Client) DeleteGuild(id Snowflake) (err error) {
	err = rest.DeleteGuild(c.req, id)
	return
}
func (c *Client) GetGuildChannels(id Snowflake) (ret []*resource.Channel, err error) {
	ret, err = rest.GetGuildChannels(c.req, id)
	return
}
func (c *Client) CreateGuildChannel(id Snowflake, params *rest.CreateGuildChannelParams) (ret *resource.Channel, err error) {
	ret, err = rest.CreateGuildChannel(c.req, id, params)
	return
}
func (c *Client) GetGuildMember(guildID, userID Snowflake) (ret *resource.Member, err error) {
	ret, err = rest.GetGuildMember(c.req, guildID, userID)
	return
}
func (c *Client) GetGuildMembers(guildID, after Snowflake, limit int) (ret []*resource.Member, err error) {
	ret, err = rest.GetGuildMembers(c.req, guildID, after, limit)
	return
}
func (c *Client) AddGuildMember(guildID, userID Snowflake, params *rest.AddGuildMemberParams) (ret *resource.Member, err error) {
	ret, err = rest.AddGuildMember(c.req, guildID, userID, params)
	return
}
func (c *Client) ModifyGuildMember(guildID, userID Snowflake, params *rest.ModifyGuildMemberParams) (err error) {
	err = rest.ModifyGuildMember(c.req, guildID, userID, params)
	return
}
func (c *Client) ModifyCurrentUserNick(id Snowflake, params *rest.ModifyCurrentUserNickParams) (nick string, err error) {
	nick, err = rest.ModifyCurrentUserNick(c.req, id, params)
	return
}
func (c *Client) AddGuildMemberRole(guildID, userID, roleID Snowflake) (err error) {
	err = rest.AddGuildMemberRole(c.req, guildID, userID, roleID)
	return
}
func (c *Client) RemoveGuildMemberRole(guildID, userID, roleID Snowflake) (err error) {
	err = rest.RemoveGuildMemberRole(c.req, guildID, userID, roleID)
	return
}
func (c *Client) RemoveGuildMember(guildID, userID Snowflake) (err error) {
	err = rest.RemoveGuildMember(c.req, guildID, userID)
	return
}
func (c *Client) GetGuildBans(id Snowflake) (ret []*resource.Ban, err error) {
	ret, err = rest.GetGuildBans(c.req, id)
	return
}
func (c *Client) GetGuildBan(guildID, userID Snowflake) (ret *resource.Ban, err error) {
	ret, err = rest.GetGuildBan(c.req, guildID, userID)
	return
}
func (c *Client) CreateGuildBan(guildID, userID Snowflake, params *rest.CreateGuildBanParams) (err error) {
	err = rest.CreateGuildBan(c.req, guildID, userID, params)
	return
}
func (c *Client) RemoveGuildBan(guildID, userID Snowflake) (err error) {
	err = rest.RemoveGuildBan(c.req, guildID, userID)
	return
}
func (c *Client) GetGuildRoles(guildID Snowflake) (ret []*resource.Role, err error) {
	ret, err = rest.GetGuildRoles(c.req, guildID)
	return
}
func (c *Client) CreateGuildRole(id Snowflake, params *rest.CreateGuildRoleParams) (ret *resource.Role, err error) {
	ret, err = rest.CreateGuildRole(c.req, id, params)
	return
}
func (c *Client) ModifyGuildRolePositions(guildID Snowflake, params *rest.ModifyGuildRolePositionsParams) (ret []*resource.Role, err error) {
	ret, err = rest.ModifyGuildRolePositions(c.req, guildID, params)
	return
}
func (c *Client) ModifyGuildRole(guildID, roleID Snowflake, params *rest.ModifyGuildRoleParams) (ret []*resource.Role, err error) {
	ret, err = rest.ModifyGuildRole(c.req, guildID, roleID, params)
	return
}
func (c *Client) DeleteGuildRole(guildID, roleID Snowflake) (err error) {
	err = rest.DeleteGuildRole(c.req, guildID, roleID)
	return
}
func (c *Client) GetGuildPruneCount(id Snowflake, params *rest.GuildPruneParams) (ret *resource.GuildPruneCount, err error) {
	ret, err = rest.GetGuildPruneCount(c.req, id, params)
	return
}
func (c *Client) BeginGuildPrune(id Snowflake, params *rest.GuildPruneParams) (ret *resource.GuildPruneCount, err error) {
	ret, err = rest.BeginGuildPrune(c.req, id, params)
	return
}
func (c *Client) GetGuildVoiceRegions(id Snowflake) (ret []*resource.VoiceRegion, err error) {
	ret, err = rest.GetGuildVoiceRegions(c.req, id)
	return
}
func (c *Client) GetGuildInvites(id Snowflake) (ret []*resource.Invite, err error) {
	ret, err = rest.GetGuildInvites(c.req, id)
	return
}
func (c *Client) GetGuildIntegrations(id Snowflake) (ret []*resource.Integration, err error) {
	ret, err = rest.GetGuildIntegrations(c.req, id)
	return
}
func (c *Client) CreateGuildIntegration(guildID Snowflake, params *rest.CreateGuildIntegrationParams) (err error) {
	err = rest.CreateGuildIntegration(c.req, guildID, params)
	return
}
func (c *Client) ModifyGuildIntegration(guildID, integrationID Snowflake, params *rest.ModifyGuildIntegrationParams) (err error) {
	err = rest.ModifyGuildIntegration(c.req, guildID, integrationID, params)
	return
}
func (c *Client) DeleteGuildIntegration(guildID, integrationID Snowflake) (err error) {
	err = rest.DeleteGuildIntegration(c.req, guildID, integrationID)
	return
}
func (c *Client) SyncGuildIntegration(guildID, integrationID Snowflake) (err error) {
	err = rest.SyncGuildIntegration(c.req, guildID, integrationID)
	return
}
func (c *Client) GetGuildEmbed(guildID Snowflake) (ret *resource.GuildEmbed, err error) {
	ret, err = rest.GetGuildEmbed(c.req, guildID)
	return
}
func (c *Client) ModifyGuildEmbed(guildID Snowflake, params *resource.GuildEmbed) (ret *resource.GuildEmbed, err error) {
	ret, err = rest.ModifyGuildEmbed(c.req, guildID, params)
	return
}
func (c *Client) GetGuildVanityURL(guildID Snowflake) (ret *resource.PartialInvite, err error) {
	ret, err = rest.GetGuildVanityURL(c.req, guildID)
	return
}

// Invite
func (c *Client) GetInvite(inviteCode string, withCounts bool) (invite *resource.Invite, err error) {
	invite, err = rest.GetInvite(c.req, inviteCode, withCounts)
	return
}
func (c *Client) DeleteInvite(inviteCode string) (invite *resource.Invite, err error) {
	invite, err = rest.DeleteInvite(c.req, inviteCode)
	return
}

// User
func (c *Client) GetCurrentUser() (ret *resource.User, err error) {
	ret, err = rest.GetCurrentUser(c.req)
	return
}
func (c *Client) GetUser(id Snowflake) (ret *resource.User, err error) {
	ret, err = rest.GetUser(c.req, id)
	return
}
func (c *Client) ModifyCurrentUser(params *rest.ModifyCurrentUserParams) (ret *resource.User, err error) {
	ret, err = rest.ModifyCurrentUser(c.req, params)
	return
}
func (c *Client) GetCurrentUserGuilds(params *rest.GetCurrentUserGuildsParams) (ret []*resource.Guild, err error) {
	ret, err = rest.GetCurrentUserGuilds(c.req, params)
	return
}
func (c *Client) LeaveGuild(id Snowflake) (err error) {
	err = rest.LeaveGuild(c.req, id)
	return
}
func (c *Client) GetUserDMs() (ret []*resource.Channel, err error) {
	ret, err = rest.GetUserDMs(c.req)
	return
}
func (c *Client) CreateDM(recipientID Snowflake) (ret *resource.Channel, err error) {
	ret, err = rest.CreateDM(c.req, recipientID)
	return
}
func (c *Client) CreateGroupDM(params *rest.CreateGroupDMParams) (ret *resource.Channel, err error) {
	ret, err = rest.CreateGroupDM(c.req, params)
	return
}
func (c *Client) GetUserConnections() (ret []*resource.UserConnection, err error) {
	ret, err = rest.GetUserConnections(c.req)
	return
}

// Voice
func (c *Client) GetVoiceRegions() (ret []*resource.VoiceRegion, err error) {
	ret, err = rest.ListVoiceRegions(c.req)
	return
}

// Webhook
func (c *Client) CreateWebhook(channelID Snowflake, params *rest.CreateWebhookParams) (ret *resource.Webhook, err error) {
	ret, err = rest.CreateWebhook(c.req, channelID, params)
	return
}
func (c *Client) GetChannelWebhooks(channelID Snowflake) (ret []*resource.Webhook, err error) {
	ret, err = rest.GetChannelWebhooks(c.req, channelID)
	return
}
func (c *Client) GetGuildWebhooks(guildID Snowflake) (ret []*resource.Webhook, err error) {
	ret, err = rest.GetGuildWebhooks(c.req, guildID)
	return
}
func (c *Client) GetWebhook(id Snowflake) (ret *resource.Webhook, err error) {
	ret, err = rest.GetWebhook(c.req, id)
	return
}
func (c *Client) GetWebhookWithToken(id Snowflake, token string) (ret *resource.Webhook, err error) {
	ret, err = rest.GetWebhookWithToken(c.req, id, token)
	return
}
func (c *Client) ModifyWebhook(newWebhook *resource.Webhook) (ret *resource.Webhook, err error) {
	ret, err = rest.ModifyWebhook(c.req, newWebhook)
	return
}
func (c *Client) ModifyWebhookWithToken(newWebhook *resource.Webhook) (ret *resource.Webhook, err error) {
	ret, err = rest.ModifyWebhookWithToken(c.req, newWebhook)
	return
}
func (c *Client) DeleteWebhook(webhookID Snowflake) (err error) {
	err = rest.DeleteWebhook(c.req, webhookID)
	return
}
func (c *Client) DeleteWebhookWithToken(id Snowflake, token string) (err error) {
	err = rest.DeleteWebhookWithToken(c.req, id, token)
	return
}
func (c *Client) ExecuteWebhook(params *rest.ExecuteWebhookParams, wait bool, URLSuffix string) (err error) {
	err = rest.ExecuteWebhook(c.req, params, wait, URLSuffix)
	return
}
func (c *Client) ExecuteSlackWebhook(params *rest.ExecuteWebhookParams, wait bool) (err error) {
	err = rest.ExecuteSlackWebhook(c.req, params, wait)
	return
}
func (c *Client) ExecuteGitHubWebhook(params *rest.ExecuteWebhookParams, wait bool) (err error) {
	err = rest.ExecuteGitHubWebhook(c.req, params, wait)
	return
}

// Custom
func (c *Client) SendMsg(channelID Snowflake, message *resource.Message) (msg *resource.Message, err error) {
	message.RLock()
	params := &rest.CreateChannelMessageParams{
		Content: message.Content,
		Nonce:   message.Nonce,
		Tts:     message.Tts,
		// File: ...
		// Embed: ...
	}
	if len(message.Embeds) > 0 {
		params.Embed = message.Embeds[0]
	}
	message.RUnlock()

	return c.CreateChannelMessage(channelID, params)
}
func (c *Client) SendMsgString(channelID Snowflake, content string) (msg *resource.Message, err error) {
	params := &rest.CreateChannelMessageParams{
		Content: content,
	}

	msg, err = c.CreateChannelMessage(channelID, params)
	return
}

func (c *Client) ChannelChan(channelID Snowflake) <-chan *resource.Channel {
	ch := make(chan *resource.Channel)

	go func(receiver chan<- *resource.Channel, storage *state.Cache) {
		result := &resource.Channel{}
		cached := true

		// check cache

		// do http request if none found
		if result == nil {
			cached = false
		}

		// return result
		receiver <- result

		// update cache with new result, if not found
		if !cached {
			//storage.MemberChan <- result
		}

		// kill the channel
		close(ch)
	}(ch, c.state)

	return ch
}

func (c *Client) ChannelsChan(GuildID Snowflake) <-chan map[Snowflake]*resource.Channel {
	ch := make(chan map[Snowflake]*resource.Channel)

	go func(receiver chan<- map[Snowflake]*resource.Channel, storage *state.Cache) {
		result := make(map[Snowflake]*resource.Channel)
		cached := true

		// check cache

		// do http request if none found
		if result == nil {
			cached = false
		}

		// return result
		receiver <- result

		// update cache with new result, if not found
		if !cached {
			//storage.MemberChan <- result
		}

		// kill the channel
		close(ch)
	}(ch, c.state)

	return ch
}

// state/caching module
func (c *Client) GuildChan(guildID Snowflake) <-chan *resource.Guild {
	ch := make(chan *resource.Guild)

	go func(receiver chan<- *resource.Guild, storage *state.Cache) {
		result := &resource.Guild{}
		cached := true

		// check cache

		// do http request if none found
		if result == nil {
			cached = false
		}

		// return result
		receiver <- result

		// update cache with new result, if not found
		if !cached {
			//storage.MemberChan <- result
		}

		// kill the channel
		close(ch)
	}(ch, c.state)

	return ch
}
func (c *Client) MsgChan(msgID Snowflake) <-chan *resource.Message {
	ch := make(chan *resource.Message)

	go func(receiver chan<- *resource.Message, storage *state.Cache) {
		result := &resource.Message{}
		cached := true

		// check cache

		// do http request if none found
		if result == nil {
			cached = false
		}

		// return result
		receiver <- result

		// update cache with new result, if not found
		if !cached {
			//storage.MemberChan <- result
		}

		// kill the channel
		close(ch)
	}(ch, c.state)

	return ch
}

type UserChan struct {
	User  *resource.User
	Err   error
	Cache bool
}

func (c *Client) UserChan(userID Snowflake) <-chan *UserChan {
	ch := make(chan *UserChan)

	go func(userID Snowflake, receiver chan<- *UserChan, storage *state.Cache) {
		response := &UserChan{
			Cache: true,
		}

		// check cache
		response.User, response.Err = storage.User(userID)
		if response.Err != nil {
			response.Cache = false
			response.Err = nil
			response.User, response.Err = rest.GetUser(c.req, userID)
		}

		// TODO: cache dead objects, to avoid http requesting the same none existent object?
		// will this ever be a problem

		// return result
		receiver <- response

		// update cache with new result, if not found
		if !response.Cache && response.User != nil {
			storage.ProcessUser(&state.UserDetail{
				User:  response.User,
				Dirty: false,
			})
		}

		// kill the channel
		close(receiver)
	}(userID, ch, c.state)

	return ch
}
func (c *Client) MemberChan(guildID, userID Snowflake) <-chan *resource.Member {
	ch := make(chan *resource.Member)

	go func(receiver chan<- *resource.Member, storage *state.Cache) {
		result := &resource.Member{}
		cached := true

		// check cache

		// do http request if none found
		if result == nil {
			cached = false
		}

		// return result
		receiver <- result

		// update cache with new result, if not found
		if !cached {
			//storage.MemberChan <- result
		}

		// kill the channel
		close(ch)
	}(ch, c.state)

	return ch
}
func (c *Client) MembersChan(guildID Snowflake) <-chan map[Snowflake]*resource.Member {
	ch := make(chan map[Snowflake]*resource.Member)

	go func(receiver chan<- map[Snowflake]*resource.Member, storage *state.Cache) {
		result := make(map[Snowflake]*resource.Member)
		cached := true

		// check cache

		// do http request if none found
		if result == nil {
			cached = false
		}

		// return result
		receiver <- result

		// update cache with new result, if not found
		if !cached {
			//storage.MemberChan <- result
		}

		// kill the channel
		close(ch)
	}(ch, c.state)

	return ch
}
