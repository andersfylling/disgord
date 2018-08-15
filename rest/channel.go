package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andersfylling/disgord/rest/httd"
	. "github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/snowflake"
)

const (
	EndpointChannels = "/channels"
)

// GetChannel [GET]         Get a channel by ID. Returns a channel object.
// Endpoint                 /channels/{channel.id}
// Rate limiter [MAJOR]     /channels/{channel.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/channel#get-channel
// Reviewed                 2018-06-07
// Comment                  -
func GetChannel(requester httd.Getter, channelID snowflake.ID) (ret *Channel, err error) {
	if channelID.Empty() {
		return nil, errors.New("not a valid snowflake")
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannel(channelID),
		Endpoint:    "/channels/" + channelID.String(),
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

// ModifyChannel [PUT/PATCH]    Update a channels settings. Requires the 'MANAGE_CHANNELS' permission for the guild.
//                              Returns a channel on success, and a 400 BAD REQUEST on invalid parameters. Fires a
//                              Channel Update Gateway event. If modifying a category, individual Channel Update
//                              events will fire for each child channel that also changes. For the PATCH method,
//                              all the JSON Params are optional.
// Endpoint                     /channels/{channel.id}
// Rate limiter [MAJOR]         /channels/{channel.id}
// Discord documentation        https://discordapp.com/developers/docs/resources/channel#modify-channel
// Reviewed                     2018-06-07
// Comment                      andersfylling: only implemented the patch method, as its parameters are optional.
func ModifyChannel(client httd.Patcher, changes *ModifyChannelParams) (ret *Channel, err error) {
	if changes.ID.Empty() {
		err = errors.New("not a valid snowflake")
		return
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannel(changes.ID),
		Endpoint:    "/channels/" + changes.ID.String(),
	}
	resp, err := client.Patch(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

// DeleteChannel [DELETE]   Delete a channel, or close a private message. Requires the 'MANAGE_CHANNELS'
//                          permission for the guild. Deleting a category does not delete its child
//                          channels; they will have their parent_id removed and a Channel Update Gateway
//                          event will fire for each of them. Returns a channel object on success. Fires a
//                          Channel Delete Gateway event.
// Endpoint                 /channels/{channel.id}
// Rate limiter [MAJOR]     /channels/{channel.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/channel#deleteclose-channel
// Reviewed                 2018-06-07
// Comment                  Deleting a guild channel cannot be undone. Use this with caution, as it
//                          is impossible to undo this action when performed on a guild channel. In
//                          contrast, when used with a private message, it is possible to undo the
//                          action by opening a private message with the recipient again.
func DeleteChannel(client httd.Deleter, channelID snowflake.ID) (err error) {
	if channelID.Empty() {
		err = errors.New("not a valid snowflake")
		return
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannel(channelID),
		Endpoint:    "/channels/" + channelID.String(),
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

// EditChannelPermissionsParams https://discordapp.com/developers/docs/resources/channel#edit-channel-permissions-json-params
type EditChannelPermissionsParams struct {
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
func EditChannelPermissions(client httd.Puter, chanID, overwriteID snowflake.ID, params *EditChannelPermissionsParams) (err error) {
	if chanID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if overwriteID.Empty() {
		return errors.New("overwriteID must be set to target the specific channel permissions")
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannel(chanID),
		Endpoint:    "/channels/" + chanID.String() + "/permissions/" + overwriteID.String(),
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
func GetChannelInvites(client httd.Getter, channelID snowflake.ID) (ret []*Invite, err error) {
	if channelID.Empty() {
		err = errors.New("channelID must be set to target the correct channel")
		return
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannel(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/invites",
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
type CreateChannelInvitesParams struct {
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
func CreateChannelInvites(client httd.Poster, channelID snowflake.ID, params *CreateChannelInvitesParams) (ret *Invite, err error) {
	if channelID.Empty() {
		err = errors.New("channelID must be set to target the correct channel")
		return
	}
	if params == nil {
		params = &CreateChannelInvitesParams{} // have to send an empty JSON object ({})
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannel(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/invites",
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
func DeleteChannelPermission(client httd.Deleter, channelID, overwriteID snowflake.ID) (err error) {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if overwriteID.Empty() {
		return errors.New("overwriteID must be set to target the specific channel permissions")
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannel(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/permissions/" + overwriteID.String(),
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
func TriggerTypingIndicator(client httd.Poster, channelID snowflake.ID) (err error) {

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannel(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/typing",
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
func GetPinnedMessages(client httd.Getter, channelID snowflake.ID) (ret []*Message, err error) {

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannel(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/pins",
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
func AddPinnedChannelMessage(client httd.Puter, channelID, msgID snowflake.ID) (err error) {

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannel(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/pints/" + msgID.String(),
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
func DeletePinnedChannelMessage(client httd.Deleter, channelID, msgID snowflake.ID) (err error) {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if msgID.Empty() {
		return errors.New("messageID must be set to target the specific channel message")
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannel(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/pins/" + msgID.String(),
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

type GroupDMAddRecipientParams struct {
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
func GroupDMAddRecipient(client httd.Puter, channelID, userID snowflake.ID, params *GroupDMAddRecipientParams) (err error) {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if userID.Empty() {
		return errors.New("userID must be set to target the specific recipient")
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannel(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/recipients/" + userID.String(),
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
func GroupDMRemoveRecipient(client httd.Deleter, channelID, userID snowflake.ID) (err error) {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if userID.Empty() {
		return errors.New("userID must be set to target the specific recipient")
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannel(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/recipients/" + userID.String(),
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
