package resource

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/snowflake"
)

const (
	// Channel types
	// https://discordapp.com/developers/docs/resources/channel#channel-object-channel-types
	ChannelTypeGuildText uint = iota
	ChannelTypeDM
	ChannelTypeGuildVoice
	ChannelTypeGroupDM
	ChannelTypeGuildCategory
)

// ChannelMessager Methods required to create a new DM (or use an existing one) and send a DM.
type ChannelMessager interface {
	CreateMessage(*Message) error // TODO: check cache for `SEND_MESSAGES` and `SEND_TTS_MESSAGES` permissions before sending.
}

// Attachment https://discordapp.com/developers/docs/resources/channel#attachment-object
type Attachment struct {
	ID       snowflake.ID `json:"id"`
	Filename string       `json:"filename"`
	Size     uint         `json:"size"`
	URL      string       `json:"url"`
	ProxyURL string       `json:"proxy_url"`
	Height   uint         `json:"height"`
	Width    uint         `json:"width"`
}

// Overwrite: https://discordapp.com/developers/docs/resources/channel#overwrite-object
type PermissionOverwrite struct {
	ID    snowflake.ID `json:"id"`    // role or user id
	Type  string       `json:"type"`  // either `role` or `member`
	Allow int          `json:"allow"` // permission bit set
	Deny  int          `json:"deny"`  // permission bit set
}

func (pmo *PermissionOverwrite) Clear() {}

func NewChannel() *Channel {
	return &Channel{}
}

// Channel
type Channel struct {
	ID                   snowflake.ID          `json:"id"`
	Type                 uint                  `json:"type"`
	GuildID              snowflake.ID          `json:"guild_id,omitempty"`              // ?|
	Position             uint                  `json:"position,omitempty"`              // ?|
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"` // ?|
	Name                 string                `json:"name,omitempty"`                  // ?|
	Topic                string                `json:"topic,omitempty"`                 // ?|
	NSFW                 bool                  `json:"nsfw,omitempty"`                  // ?|
	LastMessageID        snowflake.ID          `json:"last_message_id,omitempty"`       // ?|?, pointer
	Bitrate              uint                  `json:"bitrate,omitempty"`               // ?|
	UserLimit            uint                  `json:"user_limit,omitempty"`            // ?|
	Recipients           []*User               `json:"recipient,omitempty"`             // ?| , empty if not DM
	Icon                 string                `json:"icon,omitempty"`                  // ?|?, pointer
	OwnerID              snowflake.ID          `json:"owner_id,omitempty"`              // ?|
	ApplicationID        snowflake.ID          `json:"applicaiton_id,omitempty"`        // ?|
	ParentID             snowflake.ID          `json:"parent_id,omitempty"`             // ?|?, pointer
	LastPingTimestamp    discord.Timestamp     `json:"last_ping_timestamp,omitempty"`   // ?|

	mu sync.RWMutex `json:"-"`
}
type PartialChannel = Channel

func (c *Channel) Mention() string {
	return "<#" + c.ID.String() + ">"
}

func (c *Channel) Compare(other *Channel) bool {
	// eh
	return (c == nil && other == nil) || (other != nil && c.ID == other.ID)
}

func (c *Channel) Replicate(channel *Channel, recipients []*User) {
	// TODO: mutex is copied
	*c = *channel

	// WARNING: DM channels holds users. These should be fetched from cache.
	if recipients != nil && len(recipients) > 0 {
		c.Recipients = recipients
	} else {
		c.Recipients = []*User{}
	}
}

func (c *Channel) DeepCopy() *Channel {
	channel := NewChannel()

	c.mu.RLock()

	channel.ID = c.ID
	channel.Type = c.Type
	channel.GuildID = c.GuildID
	channel.Position = c.Position
	channel.PermissionOverwrites = c.PermissionOverwrites
	channel.Name = c.Name
	channel.Topic = c.Topic
	channel.NSFW = c.NSFW
	channel.LastMessageID = c.LastMessageID
	channel.Bitrate = c.Bitrate
	channel.UserLimit = c.UserLimit
	channel.Icon = c.Icon
	channel.OwnerID = c.OwnerID
	channel.ApplicationID = c.ApplicationID
	channel.ParentID = c.ParentID
	channel.LastPingTimestamp = c.LastPingTimestamp

	// add recipients if it's a DM
	if c.Type == ChannelTypeDM || c.Type == ChannelTypeGroupDM {
		for _, recipient := range c.Recipients {
			channel.Recipients = append(channel.Recipients, recipient.DeepCopy())
		}
	}

	c.mu.RUnlock()

	return channel
}

func (c *Channel) Clear() {
	// TODO
}

func (c *Channel) Update() {

}

func (c *Channel) Delete() {

}

func (c *Channel) Create() {
	// check if channel already exists.
}

func (c *Channel) SendMsgStr(client ChannelMessager, msgStr string) (msg *Message, err error) {
	return &Message{}, errors.New("not implemented")
}

func (c *Channel) SendMsg(client ChannelMessager, msg *Message) (err error) {
	return errors.New("not implemented")
}

// ReqGetChannel [GET]    Get a channel by ID. Returns a channel object.
// Endpoint               /channels/{channel.id}
// Rate limiter [MAJOR]   /channels/{channel.id}
// Discord documentation  https://discordapp.com/developers/docs/resources/channel#get-channel
// Reviewed               2018-06-07
// Comment                -
func ReqGetChannel(requester httd.Getter, channelID snowflake.ID) (ret *Channel, err error) {
	if channelID.Empty() {
		return nil, errors.New("not a valid snowflake")
	}

	details := &httd.Request{
		Ratelimiter: "/channels/" + channelID.String(),
	}
	resp, err := requester.Get(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

// ModifyChannelParams https://discordapp.com/developers/docs/resources/channel#modify-channel-json-params
type ModifyChannelParams = Channel

// ReqModifyChannel [PUT/PATCH] Update a channels settings. Requires the 'MANAGE_CHANNELS' permission for the guild.
//                              Returns a channel on success, and a 400 BAD REQUEST on invalid parameters. Fires a
//                              Channel Update Gateway event. If modifying a category, individual Channel Update
//                              events will fire for each child channel that also changes. For the PATCH method,
//                              all the JSON Params are optional.
// Endpoint                     /channels/{channel.id}
// Rate limiter [MAJOR]         /channels/{channel.id}
// Discord documentation        https://discordapp.com/developers/docs/resources/channel#modify-channel
// Reviewed                     2018-06-07
// Comment                      andersfylling: only implemented the patch method, as its parameters are optional.
func ReqModifyChannel(client httd.Patcher, changes *ModifyChannelParams) (ret *Channel, err error) {
	if changes.ID.Empty() {
		err = errors.New("not a valid snowflake")
		return
	}

	details := &httd.Request{
		Ratelimiter: "/channels/" + changes.ID.String(),
	}
	resp, err := client.Patch(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

// ReqDeleteChannel [DELETE]  Delete a channel, or close a private message. Requires the 'MANAGE_CHANNELS'
//                            permission for the guild. Deleting a category does not delete its child
//                            channels; they will have their parent_id removed and a Channel Update Gateway
//                            event will fire for each of them. Returns a channel object on success. Fires a
//                            Channel Delete Gateway event.
// Endpoint                   /channels/{channel.id}
// Rate limiter [MAJOR]       /channels/{channel.id}
// Discord documentation      https://discordapp.com/developers/docs/resources/channel#deleteclose-channel
// Reviewed                   2018-06-07
// Comment                    Deleting a guild channel cannot be undone. Use this with caution, as it
//                            is impossible to undo this action when performed on a guild channel. In
//                            contrast, when used with a private message, it is possible to undo the
//                            action by opening a private message with the recipient again.
func ReqDeleteChannel(client httd.Deleter, channelID snowflake.ID) (err error) {
	if channelID.Empty() {
		err = errors.New("not a valid snowflake")
		return
	}

	details := &httd.Request{
		Ratelimiter: "/channels/" + channelID.String(),
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

// ReqEditChannelPermissionsParams https://discordapp.com/developers/docs/resources/channel#edit-channel-permissions-json-params
type ReqEditChannelPermissionsParams struct {
	Allow int    `json:"allow"` // the bitwise value of all allowed permissions
	Deny  int    `json:"deny"`  // the bitwise value of all disallowed permissions
	Type  string `json:"type"`  // "member" for a user or "role" for a role
}

// ReqEditChannelPermissions [PUT]  Edit the channel permission overwrites for a user or role in a channel.
//                                  Only usable for guild channels. Requires the 'MANAGE_ROLES' permission.
//                                  Returns a 204 empty response on success. For more information about
//                                  permissions, see permissions.
// Endpoint                         /channels/{channel.id}/permissions/{overwrite.id}
// Rate limiter [MAJOR]             /channels/{channel.id}
// Discord documentation            https://discordapp.com/developers/docs/resources/channel#edit-channel-permissions
// Reviewed                         2018-06-07
// Comment                          -
func ReqEditChannelPermissions(client httd.Puter, chanID, overwriteID snowflake.ID, params *ReqEditChannelPermissionsParams) (err error) {
	if chanID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if overwriteID.Empty() {
		return errors.New("overwriteID must be set to target the specific channel permissions")
	}

	details := &httd.Request{
		Ratelimiter:     "/channels/" + chanID.String(),
		Endpoint:        "/permissions/" + overwriteID.String(),
	}
	resp, err := client.Put(details)
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

// ReqGetChannelInvites [GET] Returns a list of invite objects (with invite metadata) for the channel.
//                            Only usable for guild channels. Requires the 'MANAGE_CHANNELS' permission.
// Endpoint                   /channels/{channel.id}/invites
// Rate limiter [MAJOR]       /channels/{channel.id}
// Discord documentation      https://discordapp.com/developers/docs/resources/channel#get-channel-invites
// Reviewed                   2018-06-07
// Comment                    -
func ReqGetChannelInvites(client httd.Getter, channelID snowflake.ID) (ret []*Invite, err error) {
	if channelID.Empty() {
		err = errors.New("channelID must be set to target the correct channel")
		return
	}

	details := &httd.Request{
		Ratelimiter:     "/channels/" + channelID.String(),
		Endpoint:        "/invites",
	}
	resp, err := client.Get(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

// ReqCreateChannelInvitesParams https://discordapp.com/developers/docs/resources/channel#create-channel-invite-json-params
type ReqCreateChannelInvitesParams struct {
	MaxAge    int  `json:"max_age,omitempty"`   // duration of invite in seconds before expiry, or 0 for never. default 86400 (24 hours)
	MaxUses   int  `json:"max_uses,omitempty"`  // max number of uses or 0 for unlimited. default 0
	Temporary bool `json:"temporary,omitempty"` // whether this invite only grants temporary membership. default false
	Unique    bool `json:"unique,omitempty"`    // if true, don't try to reuse a similar invite (useful for creating many unique one time use invites). default false
}

// ReqCreateChannelInvites [POST] Create a new invite object for the channel. Only usable for guild channels.
//                                Requires the CREATE_INSTANT_INVITE permission. All JSON paramaters for this
//                                route are optional, however the request body is not. If you are not sending
//                                any fields, you still have to send an empty JSON object ({}).
//                                Returns an invite object.
// Endpoint                       /channels/{channel.id}/invites
// Rate limiter [MAJOR]           /channels/{channel.id}
// Discord documentation          https://discordapp.com/developers/docs/resources/channel#create-channel-invite
// Reviewed                       2018-06-07
// Comment                        -
func ReqCreateChannelInvites(client httd.Poster, channelID snowflake.ID, params *ReqCreateChannelInvitesParams) (ret *Invite, err error) {
	if channelID.Empty() {
		err = errors.New("channelID must be set to target the correct channel")
		return
	}
	if params == nil {
		params = &ReqCreateChannelInvitesParams{} // have to send an empty JSON object ({})
	}

	details := &httd.Request{
		Ratelimiter:     "/channels/" + channelID.String(),
		Endpoint:        "/invites",
	}
	resp, err := client.Post(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

// ReqDeleteChannelPermission [DELETE]  Delete a channel permission overwrite for a user or role in a channel.
//                                      Only usable for guild channels. Requires the 'MANAGE_ROLES' permission.
//                                      Returns a 204 empty response on success. For more information about
//                                      permissions, see permissions:
//                                      https://discordapp.com/developers/docs/topics/permissions#permissions
// Endpoint                             /channels/{channel.id}/permissions/{overwrite.id}
// Rate limiter [MAJOR]                 /channels/{channel.id}
// Discord documentation                https://discordapp.com/developers/docs/resources/channel#delete-channel-permission
// Reviewed                             2018-06-07
// Comment                              -
func ReqDeleteChannelPermission(client httd.Deleter, channelID, overwriteID snowflake.ID) (err error) {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if overwriteID.Empty() {
		return errors.New("overwriteID must be set to target the specific channel permissions")
	}

	details := &httd.Request{
		Ratelimiter:     "/channels/" + channelID.String(),
		Endpoint:        "/permissions/" + overwriteID.String(),
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

// ReqTriggerTypingIndicator [POST] Post a typing indicator for the specified channel. Generally bots should
//                                  not implement this route. However, if a bot is responding to a command and
//                                  expects the computation to take a few seconds, this endpoint may be called
//                                  to let the user know that the bot is processing their message. Returns a 204
//                                  empty response on success. Fires a Typing Start Gateway event.
// Endpoint                         /channels/{channel.id}/typing
// Rate limiter [MAJOR]             /channels/{channel.id}
// Discord documentation            https://discordapp.com/developers/docs/resources/channel#trigger-typing-indicator
// Reviewed                         2018-06-10
// Comment                          -
func ReqTriggerTypingIndicator(client httd.Poster, channelID snowflake.ID) (err error) {

	details := &httd.Request{
		Ratelimiter: "/channels/" + channelID.String(),
		Endpoint: "/typing",
	}
	resp, err := client.Post(details)
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

// ReqGetPinnedMessages [GET] Returns all pinned messages in the channel as an array of message objects.
// Endpoint                   /channels/{channel.id}/pins
// Rate limiter [MAJOR]       /channels/{channel.id}
// Discord documentation      https://discordapp.com/developers/docs/resources/channel#get-pinned-messages
// Reviewed                   2018-06-10
// Comment                    -
func ReqGetPinnedMessages(client httd.Getter, channelID snowflake.ID) (ret []*Message, err error) {

	details := &httd.Request{
		Ratelimiter: "/channels/" + channelID.String(),
		Endpoint: "/pins",
	}
	resp, err := client.Get(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

// ReqAddPinnedChannelMessage [GET] Pin a message in a channel. Requires the 'MANAGE_MESSAGES' permission.
//                                  Returns a 204 empty response on success.
// Endpoint                         /channels/{channel.id}/pins/{message.id}
// Rate limiter [MAJOR]             /channels/{channel.id}
// Discord documentation            https://discordapp.com/developers/docs/resources/channel#add-pinned-channel-message
// Reviewed                         2018-06-10
// Comment                          -
func ReqAddPinnedChannelMessage(client httd.Puter, channelID, msgID snowflake.ID) (err error) {

	details := &httd.Request{
		Ratelimiter: "/channels/" + channelID.String(),
		Endpoint: "/pints/" + msgID.String(),
	}
	resp, err := client.Put(details)
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

// ReqDeletePinnedChannelMessage [DELETE] Delete a pinned message in a channel. Requires the 'MANAGE_MESSAGES'
//                                        permission. Returns a 204 empty response on success.
//                                        Returns a 204 empty response on success.
// Endpoint                               /channels/{channel.id}/pins/{message.id}
// Rate limiter [MAJOR]                   /channels/{channel.id}
// Discord documentation                  https://discordapp.com/developers/docs/resources/channel#delete-pinned-channel-message
// Reviewed                               2018-06-10
// Comment                                -
func ReqDeletePinnedChannelMessage(client httd.Deleter, channelID, msgID snowflake.ID) (err error) {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if msgID.Empty() {
		return errors.New("messageID must be set to target the specific channel message")
	}

	details := &httd.Request{
		Ratelimiter:     "/channels/" + channelID.String(),
		Endpoint:        "/pins/" + msgID.String(),
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

type ReqGroupDMAddRecipientParams struct {
	AccessToken string `json:"access_token"` // access token of a user that has granted your app the gdm.join scope
	Nickname    string `json:"nick"`         // nickname of the user being added
}

// ReqGroupDMAddRecipient [PUT] Adds a recipient to a Group DM using their access token.
//                              Returns a 204 empty response on success.
// Endpoint                     /channels/{channel.id}/recipients/{user.id}
// Rate limiter [MAJOR]         /channels/{channel.id}
// Discord documentation        https://discordapp.com/developers/docs/resources/channel#group-dm-add-recipient
// Reviewed                     2018-06-10
// Comment                      -
func ReqGroupDMAddRecipient(client httd.Puter, channelID, userID snowflake.ID, params *ReqGroupDMAddRecipientParams) (err error) {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if userID.Empty() {
		return errors.New("userID must be set to target the specific recipient")
	}

	details := &httd.Request{
		Ratelimiter:     "/channels/" + channelID.String(),
		Endpoint:        "/recipients/" + userID.String(),
	}
	resp, err := client.Put(details)
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

// ReqGroupDMRemoveRecipient [DELETE] Removes a recipient from a Group DM.
//                                    Returns a 204 empty response on success.
// Endpoint                           /channels/{channel.id}/recipients/{user.id}
// Rate limiter [MAJOR]               /channels/{channel.id}
// Discord documentation              https://discordapp.com/developers/docs/resources/channel#group-dm-remove-recipient
// Reviewed                           2018-06-10
// Comment                            -
func ReqGroupDMRemoveRecipient(client httd.Deleter, channelID, userID snowflake.ID) (err error) {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if userID.Empty() {
		return errors.New("userID must be set to target the specific recipient")
	}

	details := &httd.Request{
		Ratelimiter: "/channels/" + channelID.String(),
		Endpoint:    "/recipients/" + userID.String(),
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
