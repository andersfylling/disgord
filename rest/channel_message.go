package rest

import (
	"errors"
	"net/http"
	"strconv"
	"sync"

	. "github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/disgord/rest/httd"
	"github.com/andersfylling/snowflake"
)

// GetChannelMessagesParams https://discordapp.com/developers/docs/resources/channel#get-channel-messages-query-string-params
// TODO: ensure limits
type GetChannelMessagesParams struct {
	Around snowflake.ID `urlparam:"around,omitempty"`
	Before snowflake.ID `urlparam:"before,omitempty"`
	After  snowflake.ID `urlparam:"after,omitempty"`
	Limit  int          `urlparam:"limit,omitempty"`
}

// getQueryString this ins't really pretty, but it works.
func (params *GetChannelMessagesParams) getQueryString() string {
	separator := "?"
	query := ""

	if !params.Around.Empty() {
		query += separator + params.Around.String()
		separator = "&"
	}

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

// GetChannelMessages [GET] Returns the messages for a channel. If operating on a guild channel, this
//                          endpoint requires the 'VIEW_CHANNEL' permission to be present on the current
//                          user. If the current user is missing the 'READ_MESSAGE_HISTORY' permission
//                          in the channel then this will return no messages (since they cannot read
//                          the message history). Returns an array of message objects on success.
// Endpoint                 /channels/{channel.id}/messages
// Rate limiter [MAJOR]     /channels/{channel.id}/messages
// Discord documentation    https://discordapp.com/developers/docs/resources/channel#get-channel-messages
// Reviewed                 2018-06-10
// Comment                  The before, after, and around keys are mutually exclusive, only one may
//                          be passed at a time. see ReqGetChannelMessagesParams.
func GetChannelMessages(client httd.Getter, channelID snowflake.ID, params *GetChannelMessagesParams) (ret []*Message, err error) {
	if channelID.Empty() {
		err = errors.New("channelID must be set to get channel messages")
		return
	}
	query := ""
	if params != nil {
		query += params.getQueryString()
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannelMessages(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/messages" + query,
		JSONParams:  params,
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetChannelMessage [GET] Returns a specific message in the channel. If operating on a guild channel,
//                         this endpoints requires the 'READ_MESSAGE_HISTORY' permission to be present
//                         on the current user. Returns a message object on success.
// Endpoint                /channels/{channel.id}/messages/{message.id}
// Rate limiter [MAJOR]    /channels/{channel.id}/messages
// Discord documentation   https://discordapp.com/developers/docs/resources/channel#get-channel-message
// Reviewed                2018-06-10
// Comment                 -
func GetChannelMessage(client httd.Getter, channelID, messageID snowflake.ID) (ret *Message, err error) {
	if channelID.Empty() {
		err = errors.New("channelID must be set to get channel messages")
		return
	}
	if messageID.Empty() {
		err = errors.New("messageID must be set to get a specific message from a channel")
		return
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannelMessages(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/messages/" + messageID.String(),
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

func NewCreateMessageByString(content string) *CreateMessageParams {
	return &CreateMessageParams{
		Content: content,
	}
}

type CreateMessageParams struct {
	Content     string        `json:"content"`
	Nonce       snowflake.ID  `json:"nonce,omitempty"`
	Tts         bool          `json:"tts,omitempty"`
	File        interface{}   `json:"file,omitempty"`  // TODO: what is this supposed to be?
	Embed       *ChannelEmbed `json:"embed,omitempty"` // embedded rich content
	PayloadJSON string        `json:"payload_json,omitempty"`
}

// CreateChannelMessage [POST] Post a message to a guild text or DM channel. If operating on a guild channel,
//                             this endpoint requires the 'SEND_MESSAGES' permission to be present on the
//                             current user. If the tts field is set to true, the SEND_TTS_MESSAGES permission
//                             is required for the message to be spoken. Returns a message object. Fires a
//                             Message Create Gateway event. See message formatting for more information on
//                             how to properly format messages.
//                             The maximum request size when sending a message is 8MB.
// Endpoint                    /channels/{channel.id}/messages
// Rate limiter [MAJOR]        /channels/{channel.id}/messages
// Discord documentation       https://discordapp.com/developers/docs/resources/channel#create-message
// Reviewed                    2018-06-10
// Comment                     Before using this endpoint, you must connect to and identify with a gateway
//                             at least once. This endpoint supports both JSON and form data bodies. It does
//                             require multipart/form-data requests instead of the normal JSON request type
//                             when uploading files. Make sure you set your Content-Type to multipart/form-data
//                             if you're doing that. Note that in that case, the embed field cannot be used,
//                             but you can pass an url-encoded JSON body as a form value for payload_json.
func CreateChannelMessage(client httd.Poster, channelID snowflake.ID, params *CreateMessageParams) (ret *Message, err error) {
	if channelID.Empty() {
		err = errors.New("channelID must be set to get channel messages")
		return
	}
	if params == nil {
		err = errors.New("message must be set")
		return
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannelMessages(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/messages",
		JSONParams:  params,
	}
	_, body, err := client.Post(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// EditMessageParams https://discordapp.com/developers/docs/resources/channel#edit-message-json-params
type EditMessageParams struct {
	Content string        `json:"content,omitempty"`
	Embed   *ChannelEmbed `json:"embed,omitempty"` // embedded rich content
}

// EditMessage [PATCH]      Edit a previously sent message. You can only edit messages that have been sent by
//                          the current user. Returns a message object. Fires a Message Update Gateway event.
// Endpoint                 /channels/{channel.id}/messages/{message.id}
// Rate limiter [MAJOR]     /channels/{channel.id}/messages
// Discord documentation    https://discordapp.com/developers/docs/resources/channel#edit-message
// Reviewed                 2018-06-10
// Comment                  All parameters to this endpoint are optional.
func EditMessage(client httd.Patcher, chanID, msgID snowflake.ID, params *EditMessageParams) (ret *Message, err error) {
	if chanID.Empty() {
		err = errors.New("channelID must be set to get channel messages")
		return
	}
	if msgID.Empty() {
		err = errors.New("msgID must be set to edit the message")
		return
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannelMessages(chanID),
		Endpoint:    "/channels/" + chanID.String() + "/messages/" + msgID.String(),
		JSONParams:  params,
	}
	_, body, err := client.Patch(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// DeleteMessage [DELETE]   Delete a message. If operating on a guild channel and trying to delete a
//                          message that was not sent by the current user, this endpoint requires the
//                          'MANAGE_MESSAGES' permission. Returns a 204 empty response on success.
//                          Fires a Message Delete Gateway event.
// Endpoint                 /channels/{channel.id}/messages/{message.id}
// Rate limiter [MAJOR]     /channels/{channel.id}/messages [DELETE]
// Discord documentation    https://discordapp.com/developers/docs/resources/channel#delete-message
// Reviewed                 2018-06-10
// Comment                  -
func DeleteMessage(client httd.Deleter, channelID, msgID snowflake.ID) (err error) {
	if channelID.Empty() {
		err = errors.New("channelID must be set to get channel messages")
		return
	}
	if msgID.Empty() {
		err = errors.New("msgID must be set to delete the message")
		return
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannelMessagesDelete(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/messages/" + msgID.String(),
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

// BulkDeleteMessagesParams https://discordapp.com/developers/docs/resources/channel#bulk-delete-messages-json-params
type BulkDeleteMessagesParams struct {
	Messages []snowflake.ID `json:"messages"`
	m        sync.RWMutex
}

func (p *BulkDeleteMessagesParams) tooMany(messages int) (err error) {
	if messages > 100 {
		err = errors.New("must be 100 or less messages to delete")
	}

	return
}

func (p *BulkDeleteMessagesParams) tooFew(messages int) (err error) {
	if messages < 2 {
		err = errors.New("must be at least two messages to delete")
	}

	return
}

func (p *BulkDeleteMessagesParams) Valid() (err error) {
	p.m.RLock()
	defer p.m.RUnlock()

	messages := len(p.Messages)
	err = p.tooMany(messages)
	if err != nil {
		return
	}
	err = p.tooFew(messages)
	return
}

func (p *BulkDeleteMessagesParams) AddMessage(msg *Message) (err error) {
	p.m.Lock()
	defer p.m.Unlock()

	err = p.tooMany(len(p.Messages) + 1)
	if err != nil {
		return
	}

	// TODO: check for duplicates as those are counted only once

	p.Messages = append(p.Messages, msg.ID)
	return
}

// BulkDeleteMessages [POST]    Delete multiple messages in a single request. This endpoint can only be used
//                              on guild channels and requires the 'MANAGE_MESSAGES' permission. Returns a 204
//                              empty response on success. Fires multiple Message Delete Gateway events.Any message
//                              IDs given that do not exist or are invalid will count towards the minimum and
//                              maximum message count (currently 2 and 100 respectively). Additionally,
//                              duplicated IDs will only be counted once.
// Endpoint                     /channels/{channel.id}/messages/bulk-delete
// Rate limiter [MAJOR]         /channels/{channel.id}/messages [DELETE] TODO: is this limiter key incorrect?
// Discord documentation        https://discordapp.com/developers/docs/resources/channel#delete-message
// Reviewed                     2018-06-10
// Comment                      This endpoint will not delete messages older than 2 weeks, and will fail if
//                              any message provided is older than that.
func BulkDeleteMessages(client httd.Poster, chanID snowflake.ID, params *BulkDeleteMessagesParams) (err error) {
	if chanID.Empty() {
		err = errors.New("channelID must be set to get channel messages")
		return
	}
	err = params.Valid()
	if err != nil {
		return
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannelMessagesDelete(chanID),
		Endpoint:    "/channels/" + chanID.String() + "/messages/bulk-delete",
	}
	resp, _, err := client.Post(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(http.StatusNoContent)
		err = errors.New(msg)
	}
	return
}
