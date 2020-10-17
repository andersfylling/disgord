package disgord

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
	"github.com/andersfylling/disgord/json"
)

// Channel types
// https://discord.com/developers/docs/resources/channel#channel-object-channel-types
const (
	ChannelTypeGuildText uint = iota
	ChannelTypeDM
	ChannelTypeGuildVoice
	ChannelTypeGroupDM
	ChannelTypeGuildCategory
	ChannelTypeGuildNews
	ChannelTypeGuildStore
)

// Attachment https://discord.com/developers/docs/resources/channel#attachment-object
type Attachment struct {
	ID       Snowflake `json:"id"`
	Filename string    `json:"filename"`
	Size     uint      `json:"size"`
	URL      string    `json:"url"`
	ProxyURL string    `json:"proxy_url"`
	Height   uint      `json:"height"`
	Width    uint      `json:"width"`

	SpoilerTag bool `json:"-"`
}

var _ internalUpdater = (*Attachment)(nil)

func (a *Attachment) updateInternals() {
	a.SpoilerTag = strings.HasPrefix(a.Filename, AttachmentSpoilerPrefix)
}

// DeepCopy see interface at struct.go#DeepCopier
func (a *Attachment) DeepCopy() (copy interface{}) {
	copy = &Attachment{
		ID:       a.ID,
		Filename: a.Filename,
		Size:     a.Size,
		URL:      a.URL,
		ProxyURL: a.ProxyURL,
		Height:   a.Height,
		Width:    a.Width,
	}

	return
}

// PermissionOverwrite https://discord.com/developers/docs/resources/channel#overwrite-object
type PermissionOverwrite struct {
	ID    Snowflake     `json:"id"`    // role or user id
	Type  string        `json:"type"`  // either `role` or `member`
	Allow PermissionBit `json:"allow"` // permission bit set
	Deny  PermissionBit `json:"deny"`  // permission bit set
}

// NewChannel ...
func NewChannel() *Channel {
	return &Channel{}
}

// type ChannelDeleter interface { DeleteChannel(id Snowflake) (err error) }
// type ChannelUpdater interface {}

// PartialChannel ...
// example of partial channel
// // "channel": {
// //   "id": "165176875973476352",
// //   "name": "illuminati",
// //   "type": 0
// // }
type PartialChannel struct {
	ID   Snowflake `json:"id"`
	Name string    `json:"name"`
	Type uint      `json:"type"`
}

// Channel ...
type Channel struct {
	ID                   Snowflake             `json:"id"`
	Type                 uint                  `json:"type"`
	GuildID              Snowflake             `json:"guild_id,omitempty"`              // ?|
	Position             int                   `json:"position,omitempty"`              // ?| can be less than 0
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"` // ?|
	Name                 string                `json:"name,omitempty"`                  // ?|
	Topic                string                `json:"topic,omitempty"`                 // ?|?
	NSFW                 bool                  `json:"nsfw,omitempty"`                  // ?|
	LastMessageID        Snowflake             `json:"last_message_id,omitempty"`       // ?|?
	Bitrate              uint                  `json:"bitrate,omitempty"`               // ?|
	UserLimit            uint                  `json:"user_limit,omitempty"`            // ?|
	RateLimitPerUser     uint                  `json:"rate_limit_per_user,omitempty"`   // ?|
	Recipients           []*User               `json:"recipient,omitempty"`             // ?| , empty if not DM/GroupDM
	Icon                 string                `json:"icon,omitempty"`                  // ?|?
	OwnerID              Snowflake             `json:"owner_id,omitempty"`              // ?|
	ApplicationID        Snowflake             `json:"application_id,omitempty"`        // ?|
	ParentID             Snowflake             `json:"parent_id,omitempty"`             // ?|?
	LastPinTimestamp     Time                  `json:"last_pin_timestamp,omitempty"`    // ?|

	// set to true when the object is not incomplete. Used in situations
	// like cacheLink to avoid overwriting correct information.
	// A partial or incomplete channel can be
	//  "channel": {
	//    "id": "165176875973476352",
	//    "name": "illuminati",
	//    "type": 0
	//  }
	complete      bool
	recipientsIDs []Snowflake
}

var _ Reseter = (*Channel)(nil)
var _ fmt.Stringer = (*Channel)(nil)
var _ Copier = (*Channel)(nil)
var _ DeepCopier = (*Channel)(nil)
var _ Mentioner = (*Channel)(nil)

func (c *Channel) String() string {
	return "channel{name:'" + c.Name + "', id:" + c.ID.String() + "}"
}

func (c *Channel) valid() bool {
	if c.RateLimitPerUser > 120 {
		return false
	}

	if len(c.Topic) > 1024 {
		return false
	}

	if c.Name != "" && (len(c.Name) > 100 || len(c.Name) < 2) {
		return false
	}

	return true
}

// GetPermissions is used to get a members permissions in a channel.
func (c *Channel) GetPermissions(ctx context.Context, s GuildQueryBuilderCaller, member *Member, flags ...Flag) (permissions PermissionBit, err error) {
	// Get the guild permissions.
	permissions, err = member.GetPermissions(ctx, s, flags...)
	if err != nil {
		return 0, err
	}

	// Handle permission overwrites.
	apply := func(o PermissionOverwrite) {
		permissions |= o.Allow
		permissions &= (-o.Deny) - 1
	}
	for _, overwrite := range c.PermissionOverwrites {
		if overwrite.Type == "member" {
			// This is a member. Is it me?
			if overwrite.ID == member.UserID {
				// It is! Time to apply the overwrites.
				apply(overwrite)
			}
			continue
		}

		for _, role := range member.Roles {
			if role == overwrite.ID {
				apply(overwrite)
				break
			}
		}
	}

	// Return the result.
	return
}

// Mention creates a channel mention string. Mention format is according the Discord protocol.
func (c *Channel) Mention() string {
	return "<#" + c.ID.String() + ">"
}

// Compare checks if channel A is the same as channel B
func (c *Channel) Compare(other *Channel) bool {
	// eh
	return (c == nil && other == nil) || (other != nil && c.ID == other.ID)
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *Channel) DeepCopy() (copy interface{}) {
	copy = NewChannel()
	_ = c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *Channel) CopyOverTo(other interface{}) (err error) {
	var channel *Channel
	var valid bool
	if channel, valid = other.(*Channel); !valid {
		err = newErrorUnsupportedType("argument given is not a *Channel type")
		return
	}

	channel.ID = c.ID
	channel.Type = c.Type
	channel.GuildID = c.GuildID
	channel.Position = c.Position
	channel.PermissionOverwrites = c.PermissionOverwrites // TODO: check for pointer
	channel.Name = c.Name
	channel.Topic = c.Topic
	channel.NSFW = c.NSFW
	channel.LastMessageID = c.LastMessageID
	channel.Bitrate = c.Bitrate
	channel.UserLimit = c.UserLimit
	channel.RateLimitPerUser = c.RateLimitPerUser
	channel.Icon = c.Icon
	channel.OwnerID = c.OwnerID
	channel.ApplicationID = c.ApplicationID
	channel.ParentID = c.ParentID
	channel.LastPinTimestamp = c.LastPinTimestamp
	channel.LastMessageID = c.LastMessageID

	// add recipients if it's a DM
	channel.Recipients = make([]*User, 0, len(c.Recipients))
	for _, recipient := range c.Recipients {
		channel.Recipients = append(channel.Recipients, recipient.DeepCopy().(*User))
	}

	return
}

// SendMsgString same as SendMsg, however this only takes the message content (string) as a argument for the message
func (c *Channel) SendMsgString(ctx context.Context, client MessageSender, content string) (msg *Message, err error) {
	if c.ID.IsZero() {
		err = newErrorMissingSnowflake("snowflake ID not set for channel")
		return
	}
	params := &CreateMessageParams{
		Content: content,
	}

	msg, err = client.CreateMessage(ctx, c.ID, params)
	return
}

// SendMsg sends a message to a channel
func (c *Channel) SendMsg(ctx context.Context, client MessageSender, message *Message) (msg *Message, err error) {
	if c.ID.IsZero() {
		err = newErrorMissingSnowflake("snowflake ID not set for channel")
		return
	}
	nonce := fmt.Sprint(message.Nonce)
	if len(nonce) > 25 {
		return nil, errors.New("nonce can not be longer than 25 characters")
	}

	params := &CreateMessageParams{
		Content: message.Content,
		Nonce:   nonce, // THIS IS A STRING. NOT A SNOWFLAKE! DONT TOUCH!
		Tts:     message.Tts,
		// File: ...
		// Embed: ...
	}
	if len(message.Embeds) > 0 {
		params.Embed = message.Embeds[0]
	}

	msg, err = client.CreateMessage(ctx, c.ID, params)
	return
}

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

func (c clientQueryBuilder) Channel(id Snowflake) ChannelQueryBuilder {
	return &channelQueryBuilder{client: c.client, cid: id}
}

// ChannelQueryBuilder REST interface for all Channel endpoints
type ChannelQueryBuilder interface {
	WithContext(ctx context.Context) ChannelQueryBuilder

	// TriggerTypingIndicator Post a typing indicator for the specified channel. Generally bots should not implement
	// this route. However, if a bot is responding to a command and expects the computation to take a few seconds, this
	// endpoint may be called to let the user know that the bot is processing their message. Returns a 204 empty response
	// on success. Fires a Typing Start Gateway event.
	TriggerTypingIndicator(flags ...Flag) error

	// GetChannel Get a channel by Snowflake. Returns a channel object.
	Get(flags ...Flag) (*Channel, error)

	// UpdateChannel Update a Channels settings. Requires the 'MANAGE_CHANNELS' permission for the guild. Returns
	// a channel on success, and a 400 BAD REQUEST on invalid parameters. Fires a Channel Update Gateway event. If
	// modifying a category, individual Channel Update events will fire for each child channel that also changes.
	// For the PATCH method, all the JSON Params are optional.
	Update(flags ...Flag) *updateChannelBuilder

	// DeleteChannel Delete a channel, or close a private message. Requires the 'MANAGE_CHANNELS' permission for
	// the guild. Deleting a category does not delete its child Channels; they will have their parent_id removed and a
	// Channel Update Gateway event will fire for each of them. Returns a channel object on success.
	// Fires a Channel Delete Gateway event.
	Delete(flags ...Flag) (*Channel, error)

	// EditChannelPermissions Edit the channel permission overwrites for a user or role in a channel. Only usable
	// for guild Channels. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty response on success.
	// For more information about permissions, see permissions.
	UpdatePermissions(overwriteID Snowflake, params *UpdateChannelPermissionsParams, flags ...Flag) error

	// GetChannelInvites Returns a list of invite objects (with invite metadata) for the channel. Only usable for
	// guild Channels. Requires the 'MANAGE_CHANNELS' permission.
	GetInvites(flags ...Flag) ([]*Invite, error)

	// CreateChannelInvite Create a new invite object for the channel. Only usable for guild Channels. Requires
	// the CREATE_INSTANT_INVITE permission. All JSON parameters for this route are optional, however the request
	// body is not. If you are not sending any fields, you still have to send an empty JSON object ({}).
	// Returns an invite object.
	CreateInvite(flags ...Flag) *createChannelInviteBuilder

	// DeleteChannelPermission Delete a channel permission overwrite for a user or role in a channel. Only usable
	// for guild Channels. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty response on success. For more
	// information about permissions,
	// see permissions: https://discord.com/developers/docs/topics/permissions#permissions
	DeletePermission(overwriteID Snowflake, flags ...Flag) error

	// AddDMParticipant Adds a recipient to a Group DM using their access token. Returns a 204 empty response
	// on success.
	AddDMParticipant(participant *GroupDMParticipant, flags ...Flag) error

	// KickParticipant Removes a recipient from a Group DM. Returns a 204 empty response on success.
	KickParticipant(userID Snowflake, flags ...Flag) error

	// GetPinnedMessages Returns all pinned messages in the channel as an array of message objects.
	GetPinnedMessages(flags ...Flag) ([]*Message, error)

	// DeleteMessages Delete multiple messages in a single request. This endpoint can only be used on guild
	// Channels and requires the 'MANAGE_MESSAGES' permission. Returns a 204 empty response on success. Fires multiple
	// Message Delete Gateway events.Any message IDs given that do not exist or are invalid will count towards
	// the minimum and maximum message count (currently 2 and 100 respectively). Additionally, duplicated IDs
	// will only be counted once.
	DeleteMessages(params *DeleteMessagesParams, flags ...Flag) error

	// GetMessages Returns the messages for a channel. If operating on a guild channel, this endpoint requires
	// the 'VIEW_CHANNEL' permission to be present on the current user. If the current user is missing
	// the 'READ_MESSAGE_HISTORY' permission in the channel then this will return no messages
	// (since they cannot read the message history). Returns an array of message objects on success.
	GetMessages(params *GetMessagesParams, flags ...Flag) ([]*Message, error)

	// CreateMessage Post a message to a guild text or DM channel. If operating on a guild channel, this
	// endpoint requires the 'SEND_MESSAGES' permission to be present on the current user. If the tts field is set to true,
	// the SEND_TTS_MESSAGES permission is required for the message to be spoken. Returns a message object. Fires a
	// Message Create Gateway event. See message formatting for more information on how to properly format messages.
	// The maximum request size when sending a message is 8MB.
	CreateMessage(params *CreateMessageParams, flags ...Flag) (*Message, error)

	// CreateWebhook Create a new webhook. Requires the 'MANAGE_WEBHOOKS' permission.
	// Returns a webhook object on success.
	CreateWebhook(params *CreateWebhookParams, flags ...Flag) (ret *Webhook, err error)

	// GetChannelWebhooks Returns a list of channel webhook objects. Requires the 'MANAGE_WEBHOOKS' permission.
	GetWebhooks(flags ...Flag) (ret []*Webhook, err error)

	Message(id Snowflake) MessageQueryBuilder
}

type channelQueryBuilder struct {
	ctx    context.Context
	client *Client
	cid    Snowflake
}

var _ ChannelQueryBuilder = (*channelQueryBuilder)(nil)

func (c channelQueryBuilder) WithContext(ctx context.Context) ChannelQueryBuilder {
	c.ctx = ctx
	return &c
}

// GetChannel [REST] Get a channel by Snowflake. Returns a channel object.
//  Method                  GET
//  Endpoint                /channels/{channel.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#get-channel
//  Reviewed                2018-06-07
//  Comment                 -
func (c channelQueryBuilder) Get(flags ...Flag) (ret *Channel, err error) {
	if c.cid.IsZero() {
		return nil, errors.New("not a valid snowflake")
	}

	channel, _ := c.client.cache.GetChannel(c.cid)
	if channel != nil {
		return channel, nil
	}

	r := c.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.Channel(c.cid),
		Ctx:      c.ctx,
	}, flags)
	r.pool = c.client.pool.channel
	r.factory = func() interface{} {
		return &Channel{}
	}

	return getChannel(r.Execute)
}

// UpdateChannel [REST] Update a Channels settings. Requires the 'MANAGE_CHANNELS' permission for the guild. Returns
// a channel on success, and a 400 BAD REQUEST on invalid parameters. Fires a Channel Update Gateway event. If
// modifying a category, individual Channel Update events will fire for each child channel that also changes.
// For the PATCH method, all the JSON Params are optional.
//  Method                  PUT/PATCH
//  Endpoint                /channels/{channel.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#modify-channel
//  Reviewed                2018-06-07
//  Comment                 andersfylling: only implemented the patch method, as its parameters are optional.
func (c channelQueryBuilder) Update(flags ...Flag) (builder *updateChannelBuilder) {
	builder = &updateChannelBuilder{}
	builder.r.itemFactory = func() interface{} {
		return c.client.pool.channel.Get()
	}
	builder.r.flags = flags
	builder.r.setup(c.client.req, &httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         c.ctx,
		Endpoint:    endpoint.Channel(c.cid),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// DeleteChannel [REST] Delete a channel, or close a private message. Requires the 'MANAGE_CHANNELS' permission for
// the guild. Deleting a category does not delete its child Channels; they will have their parent_id removed and a
// Channel Update Gateway event will fire for each of them. Returns a channel object on success.
// Fires a Channel Delete Gateway event.
//  Method                  Delete
//  Endpoint                /channels/{channel.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#deleteclose-channel
//  Reviewed                2018-10-09
//  Comment                 Deleting a guild channel cannot be undone. Use this with caution, as it
//                          is impossible to undo this action when performed on a guild channel. In
//                          contrast, when used with a private message, it is possible to undo the
//                          action by opening a private message with the recipient again.
func (c channelQueryBuilder) Delete(flags ...Flag) (channel *Channel, err error) {
	if c.cid.IsZero() {
		err = errors.New("not a valid snowflake")
		return
	}

	r := c.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.Channel(c.cid),
		Ctx:      context.Background(),
	}, flags)
	r.expectsStatusCode = http.StatusOK
	r.factory = func() interface{} {
		return &Channel{}
	}

	return getChannel(r.Execute)
}

// TriggerTypingIndicator [REST] Post a typing indicator for the specified channel. Generally bots should not implement
// this route. However, if a bot is responding to a command and expects the computation to take a few seconds, this
// endpoint may be called to let the user know that the bot is processing their message. Returns a 204 empty response
// on success. Fires a Typing Start Gateway event.
//  Method                  POST
//  Endpoint                /channels/{channel.id}/typing
//  Discord documentation   https://discord.com/developers/docs/resources/channel#trigger-typing-indicator
//  Reviewed                2018-06-10
//  Comment                 -
func (c channelQueryBuilder) TriggerTypingIndicator(flags ...Flag) (err error) {
	r := c.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodPost,
		Endpoint: endpoint.ChannelTyping(c.cid),
		Ctx:      c.ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// UpdateChannelPermissionsParams https://discord.com/developers/docs/resources/channel#edit-channel-permissions-json-params
type UpdateChannelPermissionsParams struct {
	Allow PermissionBit `json:"allow"` // the bitwise value of all allowed permissions
	Deny  PermissionBit `json:"deny"`  // the bitwise value of all disallowed permissions
	Type  string        `json:"type"`  // "member" for a user or "role" for a role
}

// EditChannelPermissions [REST] Edit the channel permission overwrites for a user or role in a channel. Only usable
// for guild Channels. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty response on success.
// For more information about permissions, see permissions.
//  Method                  PUT
//  Endpoint                /channels/{channel.id}/permissions/{overwrite.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#edit-channel-permissions
//  Reviewed                2018-06-07
//  Comment                 -
func (c channelQueryBuilder) UpdatePermissions(overwriteID Snowflake, params *UpdateChannelPermissionsParams, flags ...Flag) (err error) {
	if c.cid.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if overwriteID.IsZero() {
		return errors.New("overwriteID must be set to target the specific channel permissions")
	}

	r := c.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPut,
		Ctx:         c.ctx,
		Endpoint:    endpoint.ChannelPermission(c.cid, overwriteID),
		ContentType: httd.ContentTypeJSON,
		Body:        params,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// GetChannelInvites [REST] Returns a list of invite objects (with invite metadata) for the channel. Only usable for
// guild Channels. Requires the 'MANAGE_CHANNELS' permission.
//  Method                  GET
//  Endpoint                /channels/{channel.id}/invites
//  Discord documentation   https://discord.com/developers/docs/resources/channel#get-channel-invites
//  Reviewed                2018-06-07
//  Comment                 -
func (c channelQueryBuilder) GetInvites(flags ...Flag) (invites []*Invite, err error) {
	if c.cid.IsZero() {
		err = errors.New("channelID must be set to target the correct channel")
		return
	}

	r := c.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.ChannelInvites(c.cid),
		Ctx:      c.ctx,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Invite, 0)
		return &tmp
	}

	return getInvites(r.Execute)
}

// CreateChannelInvite [REST] Create a new invite object for the channel. Only usable for guild Channels. Requires
// the CREATE_INSTANT_INVITE permission. All JSON parameters for this route are optional, however the request body is
// not. If you are not sending any fields, you still have to send an empty JSON object ({}). Returns an invite object.
//  Method                  POST
//  Endpoint                /channels/{channel.id}/invites
//  Discord documentation   https://discord.com/developers/docs/resources/channel#create-channel-invite
//  Reviewed                2018-06-07
//  Comment                 -
func (c channelQueryBuilder) CreateInvite(flags ...Flag) (builder *createChannelInviteBuilder) {
	builder = &createChannelInviteBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Invite{}
	}
	builder.r.flags = flags
	builder.r.setup(c.client.req, &httd.Request{
		Method:      httd.MethodPost,
		Ctx:         c.ctx,
		Endpoint:    endpoint.ChannelInvites(c.cid),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// DeleteChannelPermission [REST] Delete a channel permission overwrite for a user or role in a channel. Only usable
// for guild Channels. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty response on success. For more
// information about permissions, see permissions: https://discord.com/developers/docs/topics/permissions#permissions
//  Method                  DELETE
//  Endpoint                /channels/{channel.id}/permissions/{overwrite.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#delete-channel-permission
//  Reviewed                2018-06-07
//  Comment                 -
func (c channelQueryBuilder) DeletePermission(overwriteID Snowflake, flags ...Flag) (err error) {
	if c.cid.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if overwriteID.IsZero() {
		return errors.New("overwriteID must be set to target the specific channel permissions")
	}

	r := c.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.ChannelPermission(c.cid, overwriteID),
		Ctx:      c.ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// GroupDMParticipant Information needed to add a recipient to a group chat
type GroupDMParticipant struct {
	AccessToken string    `json:"access_token"`   // access token of a user that has granted your app the gdm.join scope
	Nickname    string    `json:"nick,omitempty"` // nickname of the user being added
	UserID      Snowflake `json:"-"`
}

func (g *GroupDMParticipant) FindErrors() error {
	if g.UserID.IsZero() {
		return errors.New("missing UserID")
	}
	if g.AccessToken == "" {
		return errors.New("missing access token")
	}
	if err := ValidateUsername(g.Nickname); err != nil && g.Nickname != "" {
		return err
	}

	return nil
}

// AddDMParticipant [REST] Adds a recipient to a Group DM using their access token. Returns a 204 empty response
// on success.
//  Method                  PUT
//  Endpoint                /channels/{channel.id}/recipients/{user.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#group-dm-add-recipient
//  Reviewed                2018-06-10
//  Comment                 -
func (c channelQueryBuilder) AddDMParticipant(participant *GroupDMParticipant, flags ...Flag) error {
	if c.cid.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if participant == nil {
		return errors.New("params can not be nil")
	}
	if err := participant.FindErrors(); err != nil {
		return err
	}

	r := c.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPut,
		Ctx:         c.ctx,
		Endpoint:    endpoint.ChannelRecipient(c.cid, participant.UserID),
		Body:        participant,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err := r.Execute()
	return err
}

// KickParticipant [REST] Removes a recipient from a Group DM. Returns a 204 empty response on success.
//  Method                  DELETE
//  Endpoint                /channels/{channel.id}/recipients/{user.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#group-dm-remove-recipient
//  Reviewed                2018-06-10
//  Comment                 -
func (c channelQueryBuilder) KickParticipant(userID Snowflake, flags ...Flag) (err error) {
	if c.cid.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if userID.IsZero() {
		return errors.New("UserID must be set to target the specific recipient")
	}

	r := c.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.ChannelRecipient(c.cid, userID),
		Ctx:      c.ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// GetChannelMessagesParams https://discord.com/developers/docs/resources/channel#get-channel-messages-query-string-params
// TODO: ensure limits
type GetMessagesParams struct {
	Around Snowflake `urlparam:"around,omitempty"`
	Before Snowflake `urlparam:"before,omitempty"`
	After  Snowflake `urlparam:"after,omitempty"`
	Limit  uint      `urlparam:"limit,omitempty"`
}

func (g *GetMessagesParams) Validate() error {
	var mutuallyExclusives int
	if !g.Around.IsZero() {
		mutuallyExclusives++
	}
	if !g.Before.IsZero() {
		mutuallyExclusives++
	}
	if !g.After.IsZero() {
		mutuallyExclusives++
	}

	if mutuallyExclusives > 1 {
		return errors.New(`only one of the keys "around", "before" and "after" can be set at the time`)
	}
	return nil
}

var _ URLQueryStringer = (*GetMessagesParams)(nil)

// getMessages [REST] Returns the messages for a channel. If operating on a guild channel, this endpoint requires
// the 'VIEW_CHANNEL' permission to be present on the current user. If the current user is missing
// the 'READ_MESSAGE_HISTORY' permission in the channel then this will return no messages
// (since they cannot read the message history). Returns an array of message objects on success.
//  Method                  GET
//  Endpoint                /channels/{channel.id}/messages
//  Discord documentation   https://discord.com/developers/docs/resources/channel#get-channel-messages
//  Reviewed                2018-06-10
//  Comment                 The before, after, and around keys are mutually exclusive, only one may
//                          be passed at a time. see ReqGetChannelMessagesParams.
func (c channelQueryBuilder) getMessages(params URLQueryStringer, flags ...Flag) (ret []*Message, err error) {
	if c.cid.IsZero() {
		err = errors.New("channelID must be set to get channel messages")
		return
	}

	var query string
	if params != nil {
		query += params.URLQueryString()
	}

	r := c.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.ChannelMessages(c.cid) + query,
		Ctx:      c.ctx,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Message, 0)
		return &tmp
	}

	return getMessages(r.Execute)
}

// GetMessages bypasses discord limitations and iteratively fetches messages until the set filters are met.
func (c channelQueryBuilder) GetMessages(filter *GetMessagesParams, flags ...Flag) (messages []*Message, err error) {
	// discord values
	const filterLimit = 100
	const filterDefault = 50

	if err = filter.Validate(); err != nil {
		return nil, err
	}

	if filter.Limit == 0 {
		filter.Limit = filterDefault
		// we hardcode it here in case discord goes dumb and decided to randomly change it.
		// This avoids that the bot do not experience a new, random, behaviour on API changes
	}

	if filter.Limit <= filterLimit {
		return c.getMessages(filter, flags...)
	}

	latestSnowflake := func(msgs []*Message) (latest Snowflake) {
		for i := range msgs {
			// if msgs[i].ID.Date().After(latest.Date()) {
			if msgs[i].ID > latest {
				latest = msgs[i].ID
			}
		}
		return
	}
	earliestSnowflake := func(msgs []*Message) (earliest Snowflake) {
		for i := range msgs {
			// if msgs[i].ID.Date().Before(earliest.Date()) {
			if msgs[i].ID < earliest {
				earliest = msgs[i].ID
			}
		}
		return
	}

	// scenario#1: filter.Around is not 0 AND filter.Limit is above 100
	//  divide the limit by half and use .Before and .After tags on each quotient limit.
	//  Use the .After on potential remainder.
	//  Note! This method can be used recursively
	if !filter.Around.IsZero() {
		beforeParams := *filter
		beforeParams.Before = beforeParams.Around
		beforeParams.Around = 0
		beforeParams.Limit = filter.Limit / 2
		befores, err := c.GetMessages(&beforeParams, flags...)
		if err != nil {
			return nil, err
		}
		messages = append(messages, befores...)

		afterParams := *filter
		afterParams.After = afterParams.Around
		afterParams.Around = 0
		afterParams.Limit = filter.Limit / 2
		afters, err := c.GetMessages(&afterParams, flags...)
		if err != nil {
			return nil, err
		}
		messages = append(messages, afters...)

		// filter.Around includes the given ID, so should .Before and .After iterations do as well
		if msg, _ := c.Message(filter.Around).Get(c.ctx, flags...); msg != nil {
			// assumption: error here can be caused by the message ID not actually being a real message
			//             and that it was used to get messages in the vicinity. Therefore the err is ignored.
			// TODO: const discord errors.
			messages = append(messages, msg)
		}
	} else {
		// scenario#3: filter.After or filter.Before is set.
		// note that none might be set, which will cause filter.Before to be set after the first 100 messages.
		//
		for {
			if filter.Limit <= 0 {
				break
			}

			f := *filter
			if f.Limit > 100 {
				f.Limit = 100
			}
			filter.Limit -= f.Limit
			msgs, err := c.getMessages(&f, flags...)
			if err != nil {
				return nil, err
			}
			messages = append(messages, msgs...)
			if !filter.After.IsZero() {
				filter.After = latestSnowflake(msgs)
			} else {
				// no snowflake or filter.Before
				filter.Before = earliestSnowflake(msgs)
			}
		}
	}

	// duplicates should not exist as we use snowflakes to fetch unique segments in time
	return messages, nil
}

// DeleteMessagesParams https://discord.com/developers/docs/resources/channel#bulk-delete-messages-json-params
type DeleteMessagesParams struct {
	Messages []Snowflake `json:"messages"`
	m        sync.RWMutex
}

func (p *DeleteMessagesParams) tooMany(messages int) (err error) {
	if messages > 100 {
		err = errors.New("must be 100 or less messages to delete")
	}

	return
}

func (p *DeleteMessagesParams) tooFew(messages int) (err error) {
	if messages < 2 {
		err = errors.New("must be at least two messages to delete")
	}

	return
}

// Valid validates the DeleteMessagesParams data
func (p *DeleteMessagesParams) Valid() (err error) {
	p.m.RLock()
	defer p.m.RUnlock()

	messages := len(p.Messages)
	if err = p.tooMany(messages); err != nil {
		return
	}
	err = p.tooFew(messages)
	return
}

// AddMessage Adds a message to be deleted
func (p *DeleteMessagesParams) AddMessage(msg *Message) (err error) {
	p.m.Lock()
	defer p.m.Unlock()

	if err = p.tooMany(len(p.Messages) + 1); err != nil {
		return
	}

	// TODO: check for duplicates as those are counted only once

	p.Messages = append(p.Messages, msg.ID)
	return
}

// DeleteMessages [REST] Delete multiple messages in a single request. This endpoint can only be used on guild
// Channels and requires the 'MANAGE_MESSAGES' permission. Returns a 204 empty response on success. Fires multiple
// Message Delete Gateway events.Any message IDs given that do not exist or are invalid will count towards
// the minimum and maximum message count (currently 2 and 100 respectively). Additionally, duplicated IDs
// will only be counted once.
//  Method                  POST
//  Endpoint                /channels/{channel.id}/messages/bulk-delete
//  Discord documentation   https://discord.com/developers/docs/resources/channel#delete-message
//  Reviewed                2018-06-10
//  Comment                 This endpoint will not delete messages older than 2 weeks, and will fail if any message
//                          provided is older than that.
func (c channelQueryBuilder) DeleteMessages(params *DeleteMessagesParams, flags ...Flag) (err error) {
	if c.cid.IsZero() {
		err = errors.New("channelID must be set to get channel messages")
		return err
	}
	if err = params.Valid(); err != nil {
		return err
	}

	r := c.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Ctx:         c.ctx,
		Endpoint:    endpoint.ChannelMessagesBulkDelete(c.cid),
		ContentType: httd.ContentTypeJSON,
		Body:        params,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// AllowedMentions allows finer control over mentions in a message, see
// https://discord.com/developers/docs/resources/channel#allowed-mentions-object for more info.
// Any strings in the Parse value must be any from ["everyone", "users", "roles"].
type AllowedMentions struct {
	Parse []string `json:"parse"` // this is purposefully not marked as omitempty as to allow `parse: []` which blocks all mentions.

	Roles []Snowflake `json:"roles,omitempty"`
	Users []Snowflake `json:"users,omitempty"`
}

// CreateMessageFileParams contains the information needed to upload a file to Discord, it is part of the
// CreateMessageParams struct.
type CreateMessageFileParams struct {
	Reader   io.Reader `json:"-"` // always omit as we don't want this as part of the JSON payload
	FileName string    `json:"-"`

	// SpoilerTag lets discord know that this image should be blurred out.
	// Current Discord behaviour is that whenever a message with one or more images is marked as
	// spoiler tag, all the images in that message are blurred out. (independent of msg.Content)
	SpoilerTag bool `json:"-"`
}

// write helper for file uploading in messages
func (f *CreateMessageFileParams) write(i int, mp *multipart.Writer) error {
	var filename string
	if f.SpoilerTag {
		filename = AttachmentSpoilerPrefix + f.FileName
	} else {
		filename = f.FileName
	}
	w, err := mp.CreateFormFile("file"+strconv.FormatInt(int64(i), 10), filename)
	if err != nil {
		return err
	}

	if _, err = io.Copy(w, f.Reader); err != nil {
		return err
	}

	return nil
}

// CreateMessageParams JSON params for CreateChannelMessage
type CreateMessageParams struct {
	Content string `json:"content"`
	Nonce   string `json:"nonce,omitempty"` // THIS IS A STRING. NOT A SNOWFLAKE! DONT TOUCH!
	Tts     bool   `json:"tts,omitempty"`
	Embed   *Embed `json:"embed,omitempty"` // embedded rich content

	Files []CreateMessageFileParams `json:"-"` // Always omit as this is included in multipart, not JSON payload

	SpoilerTagContent        bool `json:"-"`
	SpoilerTagAllAttachments bool `json:"-"`

	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"` // The allowed mentions object for the message.
}

func (p *CreateMessageParams) prepare() (postBody interface{}, contentType string, err error) {
	// spoiler tag
	if p.SpoilerTagContent && len(p.Content) > 0 {
		p.Content = "|| " + p.Content + " ||"
	}

	if len(p.Files) == 0 {
		postBody = p
		contentType = httd.ContentTypeJSON
		return
	}

	if p.SpoilerTagAllAttachments {
		for i := range p.Files {
			p.Files[i].SpoilerTag = true
		}
	}

	if p.Embed != nil {
		// check for spoilers
		for i := range p.Files {
			if p.Files[i].SpoilerTag && strings.Contains(p.Embed.Image.URL, p.Files[i].FileName) {
				s := strings.Split(p.Embed.Image.URL, p.Files[i].FileName)
				if len(s) > 0 {
					s[0] += AttachmentSpoilerPrefix + p.Files[i].FileName
					p.Embed.Image.URL = strings.Join(s, "")
				}
			}
		}
	}

	// Set up a new multipart writer, as we'll be using this for the POST body instead
	buf := new(bytes.Buffer)
	mp := multipart.NewWriter(buf)

	// Write the existing JSON payload
	var payload []byte
	payload, err = json.Marshal(p)
	if err != nil {
		return
	}
	if err = mp.WriteField("payload_json", string(payload)); err != nil {
		return
	}

	// Iterate through all the files and write them to the multipart blob
	for i, file := range p.Files {
		if err = file.write(i, mp); err != nil {
			return
		}
	}

	mp.Close()

	postBody = buf
	contentType = mp.FormDataContentType()

	return
}

// CreateMessage [REST] Post a message to a guild text or DM channel. If operating on a guild channel, this
// endpoint requires the 'SEND_MESSAGES' permission to be present on the current user. If the tts field is set to true,
// the SEND_TTS_MESSAGES permission is required for the message to be spoken. Returns a message object. Fires a
// Message Create Gateway event. See message formatting for more information on how to properly format messages.
// The maximum request size when sending a message is 8MB.
//  Method                  POST
//  Endpoint                /channels/{channel.id}/messages
//  Discord documentation   https://discord.com/developers/docs/resources/channel#create-message
//  Reviewed                2018-06-10
//  Comment                 Before using this endpoint, you must connect to and identify with a gateway at least once.
func (c channelQueryBuilder) CreateMessage(params *CreateMessageParams, flags ...Flag) (ret *Message, err error) {
	if c.cid.IsZero() {
		err = errors.New("channelID must be set to get channel messages")
		return nil, err
	}
	if params == nil {
		err = errors.New("message must be set")
		return nil, err
	}

	var (
		postBody    interface{}
		contentType string
	)

	if postBody, contentType, err = params.prepare(); err != nil {
		return nil, err
	}

	r := c.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Ctx:         c.ctx,
		Endpoint:    "/channels/" + c.cid.String() + "/messages",
		Body:        postBody,
		ContentType: contentType,
	}, flags)
	r.pool = c.client.pool.message
	r.factory = func() interface{} {
		return &Message{}
	}

	return getMessage(r.Execute)
}

// GetPinnedMessages [REST] Returns all pinned messages in the channel as an array of message objects.
//  Method                  GET
//  Endpoint                /channels/{channel.id}/pins
//  Discord documentation   https://discord.com/developers/docs/resources/channel#get-pinned-messages
//  Reviewed                2018-06-10
//  Comment                 -
func (c channelQueryBuilder) GetPinnedMessages(flags ...Flag) (ret []*Message, err error) {
	r := c.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.ChannelPins(c.cid),
		Ctx:      c.ctx,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Message, 0)
		return &tmp
	}

	return getMessages(r.Execute)
}

// CreateWebhookParams json params for the create webhook rest request avatar string
// https://discord.com/developers/docs/resources/user#avatar-data
type CreateWebhookParams struct {
	Name   string `json:"name"`   // name of the webhook (2-32 characters)
	Avatar string `json:"avatar"` // avatar data uri scheme, image for the default webhook avatar

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string `json:"-"`
}

func (c *CreateWebhookParams) FindErrors() error {
	if c.Name == "" {
		return errors.New("webhook must have a name")
	}
	if !(2 <= len(c.Name) && len(c.Name) <= 32) {
		return errors.New("webhook name must be 2 to 32 characters long")
	}
	return nil
}

// CreateWebhook [REST] Create a new webhook. Requires the 'MANAGE_WEBHOOKS' permission.
// Returns a webhook object on success.
//  Method                  POST
//  Endpoint                /channels/{channel.id}/webhooks
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#create-webhook
//  Reviewed                2018-08-14
//  Comment                 -
func (c channelQueryBuilder) CreateWebhook(params *CreateWebhookParams, flags ...Flag) (ret *Webhook, err error) {
	if params == nil {
		return nil, errors.New("params was nil")
	}
	if err = params.FindErrors(); err != nil {
		return nil, err
	}

	r := c.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Ctx:         c.ctx,
		Endpoint:    endpoint.ChannelWebhooks(c.cid),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
		Reason:      params.Reason,
	}, flags)
	r.factory = func() interface{} {
		return &Webhook{}
	}

	return getWebhook(r.Execute)
}

// GetChannelWebhooks [REST] Returns a list of channel webhook objects. Requires the 'MANAGE_WEBHOOKS' permission.
//  Method                  POST
//  Endpoint                /channels/{channel.id}/webhooks
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#get-channel-webhooks
//  Reviewed                2018-08-14
//  Comment                 -
func (c channelQueryBuilder) GetWebhooks(flags ...Flag) (ret []*Webhook, err error) {
	r := c.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.ChannelWebhooks(c.cid),
		Ctx:      c.ctx,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Webhook, 0)
		return &tmp
	}

	return getWebhooks(r.Execute)
}

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

//generate-rest-params: max_age:int, max_uses:int, temporary:bool, unique:bool,
//generate-rest-basic-execute: invite:*Invite,
type createChannelInviteBuilder struct {
	r RESTBuilder
}

func (b *createChannelInviteBuilder) WithReason(reason string) *createChannelInviteBuilder {
	b.r.headerReason = reason
	return b
}

// updateChannelBuilder https://discord.com/developers/docs/resources/channel#modify-channel-json-params
//generate-rest-params: parent_id:Snowflake, permission_overwrites:[]PermissionOverwrite, user_limit:uint, bitrate:uint, rate_limit_per_user:uint, nsfw:bool, topic:string, position:int, name:string,
//generate-rest-basic-execute: channel:*Channel,
type updateChannelBuilder struct {
	r RESTBuilder
}

func (b *updateChannelBuilder) AddPermissionOverwrite(permission PermissionOverwrite) *updateChannelBuilder {
	if _, exists := b.r.body["permission_overwrites"]; !exists {
		b.SetPermissionOverwrites([]PermissionOverwrite{permission})
	} else {
		s := b.r.body["permission_overwrites"].([]PermissionOverwrite)
		s = append(s, permission)
		b.SetPermissionOverwrites(s)
	}
	return b
}
func (b *updateChannelBuilder) AddPermissionOverwrites(permissions []PermissionOverwrite) *updateChannelBuilder {
	for i := range permissions {
		b.AddPermissionOverwrite(permissions[i])
	}
	return b
}

func (b *updateChannelBuilder) RemoveParentID() *updateChannelBuilder {
	b.r.param("parent_id", nil)
	return b
}
