package disgord

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
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

// ChannelMessager Methods required to create a new DM (or use an existing one) and send a DM.
// type ChannelMessager interface {CreateMessage(*Message) error}

// ChannelFetcher holds the single method for fetching a channel from the Discord REST API
type ChannelFetcher interface {
	GetChannel(id Snowflake) (ret *Channel, err error)
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
var _ discordDeleter = (*Channel)(nil)
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
func (c *Channel) GetPermissions(ctx context.Context, s PermissionFetching, member *Member, flags ...Flag) (permissions PermissionBit, err error) {
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

func (c *Channel) deleteFromDiscord(ctx context.Context, s Session, flags ...Flag) (err error) {
	id := c.ID

	if id.IsZero() {
		err = newErrorMissingSnowflake("channel id/snowflake is empty or missing")
		return
	}
	var deleted *Channel
	if deleted, err = s.DeleteChannel(ctx, id, flags...); err != nil {
		return
	}

	_ = deleted.CopyOverTo(c)
	return
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

func (c *Channel) copyOverToCache(other interface{}) (err error) {
	return c.CopyOverTo(other)
}

//func (c *Channel) Clear() {
//	// TODO
//}

// Fetch check if there are any updates to the channel values
//func (c *Channel) Fetch(Client ChannelFetcher) (err error) {
//	if c.ID.IsZero() {
//		err = errors.New("missing channel ID")
//		return
//	}
//
//	Client.GetChannel(c.ID)
//}

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

// GetChannel [REST] Get a channel by Snowflake. Returns a channel object.
//  Method                  GET
//  Endpoint                /channels/{channel.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#get-channel
//  Reviewed                2018-06-07
//  Comment                 -
func (c *Client) GetChannel(ctx context.Context, channelID Snowflake, flags ...Flag) (ret *Channel, err error) {
	if channelID.IsZero() {
		return nil, errors.New("not a valid snowflake")
	}

	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.Channel(channelID),
		Ctx:      ctx,
	}, flags)
	r.CacheRegistry = ChannelCache
	r.ID = channelID
	r.pool = c.pool.channel
	r.factory = func() interface{} {
		return &Channel{}
	}

	return getChannel(r.Execute)
}

// UpdateChannel [REST] Update a channels settings. Requires the 'MANAGE_CHANNELS' permission for the guild. Returns
// a channel on success, and a 400 BAD REQUEST on invalid parameters. Fires a Channel Update Gateway event. If
// modifying a category, individual Channel Update events will fire for each child channel that also changes.
// For the PATCH method, all the JSON Params are optional.
//  Method                  PUT/PATCH
//  Endpoint                /channels/{channel.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#modify-channel
//  Reviewed                2018-06-07
//  Comment                 andersfylling: only implemented the patch method, as its parameters are optional.
func (c *Client) UpdateChannel(ctx context.Context, channelID Snowflake, flags ...Flag) (builder *updateChannelBuilder) {
	builder = &updateChannelBuilder{}
	builder.r.itemFactory = func() interface{} {
		return c.pool.channel.Get()
	}
	builder.r.flags = flags
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         ctx,
		Endpoint:    endpoint.Channel(channelID),
		ContentType: httd.ContentTypeJSON,
	}, nil)
	builder.r.cacheRegistry = ChannelCache
	builder.r.cacheItemID = channelID

	return builder
}

// DeleteChannel [REST] Delete a channel, or close a private message. Requires the 'MANAGE_CHANNELS' permission for
// the guild. Deleting a category does not delete its child channels; they will have their parent_id removed and a
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
func (c *Client) DeleteChannel(ctx context.Context, channelID Snowflake, flags ...Flag) (channel *Channel, err error) {
	if channelID.IsZero() {
		err = errors.New("not a valid snowflake")
		return
	}

	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.Channel(channelID),
		Ctx:      context.Background(),
	}, flags)
	r.expectsStatusCode = http.StatusOK
	r.updateCache = func(registry cacheRegistry, id Snowflake, x interface{}) (err error) {
		c.cache.DeleteChannel(id)
		return nil
	}
	r.factory = func() interface{} {
		return &Channel{}
	}

	return getChannel(r.Execute)
}

// UpdateChannelPermissionsParams https://discord.com/developers/docs/resources/channel#edit-channel-permissions-json-params
type UpdateChannelPermissionsParams struct {
	Allow PermissionBit `json:"allow"` // the bitwise value of all allowed permissions
	Deny  PermissionBit `json:"deny"`  // the bitwise value of all disallowed permissions
	Type  string        `json:"type"`  // "member" for a user or "role" for a role
}

// EditChannelPermissions [REST] Edit the channel permission overwrites for a user or role in a channel. Only usable
// for guild channels. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty response on success.
// For more information about permissions, see permissions.
//  Method                  PUT
//  Endpoint                /channels/{channel.id}/permissions/{overwrite.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#edit-channel-permissions
//  Reviewed                2018-06-07
//  Comment                 -
func (c *Client) UpdateChannelPermissions(ctx context.Context, channelID, overwriteID Snowflake, params *UpdateChannelPermissionsParams, flags ...Flag) (err error) {
	if channelID.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if overwriteID.IsZero() {
		return errors.New("overwriteID must be set to target the specific channel permissions")
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      httd.MethodPut,
		Ctx:         ctx,
		Endpoint:    endpoint.ChannelPermission(channelID, overwriteID),
		ContentType: httd.ContentTypeJSON,
		Body:        params,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent
	r.updateCache = func(registry cacheRegistry, id Snowflake, x interface{}) (err error) {
		// TODO-cache: update cache
		return nil
	}

	_, err = r.Execute()
	return err
}

// GetChannelInvites [REST] Returns a list of invite objects (with invite metadata) for the channel. Only usable for
// guild channels. Requires the 'MANAGE_CHANNELS' permission.
//  Method                  GET
//  Endpoint                /channels/{channel.id}/invites
//  Discord documentation   https://discord.com/developers/docs/resources/channel#get-channel-invites
//  Reviewed                2018-06-07
//  Comment                 -
func (c *Client) GetChannelInvites(ctx context.Context, channelID Snowflake, flags ...Flag) (invites []*Invite, err error) {
	if channelID.IsZero() {
		err = errors.New("channelID must be set to target the correct channel")
		return
	}

	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.ChannelInvites(channelID),
		Ctx:      ctx,
	}, flags)
	r.CacheRegistry = ChannelCache
	r.factory = func() interface{} {
		tmp := make([]*Invite, 0)
		return &tmp
	}

	return getInvites(r.Execute)
}

// CreateChannelInvitesParams https://discord.com/developers/docs/resources/channel#create-channel-invite-json-params
type CreateChannelInvitesParams struct {
	MaxAge    int  `json:"max_age,omitempty"`   // duration of invite in seconds before expiry, or 0 for never. default 86400 (24 hours)
	MaxUses   int  `json:"max_uses,omitempty"`  // max number of uses or 0 for unlimited. default 0
	Temporary bool `json:"temporary,omitempty"` // whether this invite only grants temporary membership. default false
	Unique    bool `json:"unique,omitempty"`    // if true, don't try to reuse a similar invite (useful for creating many unique one time use invites). default false

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string `json:"-"`
}

// CreateChannelInvites [REST] Create a new invite object for the channel. Only usable for guild channels. Requires
// the CREATE_INSTANT_INVITE permission. All JSON parameters for this route are optional, however the request body is
// not. If you are not sending any fields, you still have to send an empty JSON object ({}). Returns an invite object.
//  Method                  POST
//  Endpoint                /channels/{channel.id}/invites
//  Discord documentation   https://discord.com/developers/docs/resources/channel#create-channel-invite
//  Reviewed                2018-06-07
//  Comment                 -
func (c *Client) CreateChannelInvites(ctx context.Context, channelID Snowflake, params *CreateChannelInvitesParams, flags ...Flag) (ret *Invite, err error) {
	if channelID.IsZero() {
		err = errors.New("channelID must be set to target the correct channel")
		return nil, err
	}
	if params == nil {
		params = &CreateChannelInvitesParams{} // have to send an empty JSON object ({}). maybe just struct{}?
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Ctx:         ctx,
		Endpoint:    endpoint.ChannelInvites(channelID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
		Reason:      params.Reason,
	}, flags)
	r.factory = func() interface{} {
		return &Invite{}
	}

	return getInvite(r.Execute)
}

// DeleteChannelPermission [REST] Delete a channel permission overwrite for a user or role in a channel. Only usable
// for guild channels. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty response on success. For more
// information about permissions, see permissions: https://discord.com/developers/docs/topics/permissions#permissions
//  Method                  DELETE
//  Endpoint                /channels/{channel.id}/permissions/{overwrite.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#delete-channel-permission
//  Reviewed                2018-06-07
//  Comment                 -
func (c *Client) DeleteChannelPermission(ctx context.Context, channelID, overwriteID Snowflake, flags ...Flag) (err error) {
	if channelID.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if overwriteID.IsZero() {
		return errors.New("overwriteID must be set to target the specific channel permissions")
	}

	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.ChannelPermission(channelID, overwriteID),
		Ctx:      ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent
	r.updateCache = func(registry cacheRegistry, id Snowflake, x interface{}) (err error) {
		_ = c.cache.DeleteChannelPermissionOverwrite(channelID, overwriteID)
		return nil
	}

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
func (c *Client) AddDMParticipant(ctx context.Context, channelID Snowflake, participant *GroupDMParticipant, flags ...Flag) error {
	if channelID.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if participant == nil {
		return errors.New("params can not be nil")
	}
	if err := participant.FindErrors(); err != nil {
		return err
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      httd.MethodPut,
		Ctx:         ctx,
		Endpoint:    endpoint.ChannelRecipient(channelID, participant.UserID),
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
func (c *Client) KickParticipant(ctx context.Context, channelID, userID Snowflake, flags ...Flag) (err error) {
	if channelID.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if userID.IsZero() {
		return errors.New("UserID must be set to target the specific recipient")
	}

	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.ChannelRecipient(channelID, userID),
		Ctx:      ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

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
