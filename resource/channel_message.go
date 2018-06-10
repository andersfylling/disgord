package resource

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/snowflake"
)

const (
	_ = iota
	MessageActivityTypeJoin
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	MessageActivityTypeJoinRequest
)
const (
	MessageTypeDefault = iota
	MessageTypeRecipientAdd
	MessageTypeRecipientRemove
	MessageTypeCall
	MessageTypeChannelNameChange
	MessageTypeChannelIconChange
	MessageTypeChannelPinnedMessage
	MessageTypeGuildMemberJoin
)

func NewMessage() *Message {
	return &Message{}
}

func NewDeletedMessage() *DeletedMessage {
	return &DeletedMessage{}
}

type DeletedMessage struct {
	ID        snowflake.ID `json:"id"`
	ChannelID snowflake.ID `json:"channel_id"`
}

// https://discordapp.com/developers/docs/resources/channel#message-object-message-activity-structure
type MessageActivity struct {
	Type    int    `json:"type"`
	PartyID string `json:"party_id"`
}

// https://discordapp.com/developers/docs/resources/channel#message-object-message-application-structure
type MessageApplication struct {
	ID          snowflake.ID `json:"id"`
	CoverImage  string       `json:"cover_image"`
	Description string       `json:"description"`
	Icon        string       `json:"icon"`
	Name        string       `json:"name"`
}

// Message https://discordapp.com/developers/docs/resources/channel#message-object-message-structure
type Message struct {
	ID              snowflake.ID       `json:"id"`
	ChannelID       snowflake.ID       `json:"channel_id"`
	Author          *User              `json:"author"`
	Content         string             `json:"content"`
	Timestamp       time.Time          `json:"timestamp"`
	EditedTimestamp time.Time          `json:"edited_timestamp"` // ?
	Tts             bool               `json:"tts"`
	MentionEveryone bool               `json:"mention_everyone"`
	Mentions        []*User            `json:"mentions"`
	MentionRoles    []snowflake.ID     `json:"mention_roles"`
	Attachments     []*Attachment      `json:"attachments"`
	Embeds          []*ChannelEmbed    `json:"embeds"`
	Reactions       []*Reaction        `json:"reactions"` // ?
	Nonce           snowflake.ID       `json:"nonce"`     // ?, used for validating a message was sent
	Pinned          bool               `json:"pinned"`
	WebhookID       snowflake.ID       `json:"webhook_id"` // ?
	Type            uint               `json:"type"`
	Activity        MessageActivity    `json:"activity"`
	Application     MessageApplication `json:"application"`

	sync.RWMutex `json:"-"`
}

func (m *Message) MarshalJSON() ([]byte, error) {
	if m.ID.Empty() {
		return []byte("{}"), nil
	}

	//TODO: remove copying of mutex
	return json.Marshal(Message(*m))
}

func (m *Message) Delete() {}
func (m *Message) Update() {}
func (m *Message) Send()   {}

func (m *Message) AddReaction(reaction *Reaction) {}
func (m *Message) RemoveReaction(id snowflake.ID) {}

// ReqGetChannelMessagesParams https://discordapp.com/developers/docs/resources/channel#get-channel-messages-query-string-params
// TODO: ensure limits
type ReqGetChannelMessagesParams struct {
	Around snowflake.ID `urlparam:"around,omitempty"`
	Before snowflake.ID `urlparam:"before,omitempty"`
	After  snowflake.ID `urlparam:"after,omitempty"`
	Limit  int          `urlparam:"limit,omitempty"`
}

// getQueryString this ins't really pretty, but it works.
func (params *ReqGetChannelMessagesParams) getQueryString() string {
	seperator := "?"
	query := ""

	if !params.Around.Empty() {
		query += seperator + params.Around.String()
		seperator = "&"
	}

	if !params.Before.Empty() {
		query += seperator + params.Before.String()
		seperator = "&"
	}

	if !params.After.Empty() {
		query += seperator + params.After.String()
		seperator = "&"
	}

	if params.Limit > 0 {
		query += seperator + strconv.Itoa(params.Limit)
	}

	return query
}

// ReqGetChannelMessages [GET]  Returns the messages for a channel. If operating on a guild channel, this
//                              endpoint requires the 'VIEW_CHANNEL' permission to be present on the current
//                              user. If the current user is missing the 'READ_MESSAGE_HISTORY' permission
//                              in the channel then this will return no messages (since they cannot read
//                              the message history). Returns an array of message objects on success.
// Endpoint                     /channels/{channel.id}/messages
// Rate limiter [MAJOR]         /channels/{channel.id}
// Discord documentation        https://discordapp.com/developers/docs/resources/channel#get-channel-messages
// Reviewed                     2018-06-10
// Comment                      The before, after, and around keys are mutually exclusive, only one may
//                              be passed at a time. see ReqGetChannelMessagesParams.
func ReqGetChannelMessages(client httd.Getter, channelID snowflake.ID, params *ReqGetChannelMessagesParams) (ret []*Message, err error) {
	if channelID.Empty() {
		err = errors.New("channelID must be set to get channel messages")
		return
	}
	query := ""
	if params != nil {
		query += params.getQueryString()
	}

	details := &httd.Request{
		Ratelimiter: "/channels/" + channelID.String(),
		Endpoint: "/messages" + query,
		JSONParams: params,
	}
	resp, err := client.Get(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

// ReqGetChannelMessage [GET] Returns a specific message in the channel. If operating on a guild channel,
//                            this endpoints requires the 'READ_MESSAGE_HISTORY' permission to be present
//                            on the current user. Returns a message object on success.
// Endpoint                   /channels/{channel.id}/message/{message.id}
// Rate limiter [MAJOR]       /channels/{channel.id}
// Discord documentation      https://discordapp.com/developers/docs/resources/channel#get-channel-message
// Reviewed                   2018-06-10
// Comment                    -
func ReqGetChannelMessage(client httd.Getter, channelID, messageID snowflake.ID) (ret *Message, err error) {
	if channelID.Empty() {
		err = errors.New("channelID must be set to get channel messages")
		return
	}
	if messageID.Empty() {
		err = errors.New("messageID must be set to get a specific message from a channel")
		return
	}

	details := &httd.Request{
		Ratelimiter: "/channels/" + channelID.String(),
		Endpoint: "/messages/" + messageID.String(),
	}
	resp, err := client.Get(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

type ReqCreateMessageParams struct {
	Content     string        `json:"content"`
	Nonce       snowflake.ID  `json:"nonce,omitempty"`
	Tts         bool          `json:"tts,omitempty"`
	File        interface{}   `json:"file,omitempty"`  // TODO: what is this supposed to be?
	Embed       *ChannelEmbed `json:"embed,omitempty"` // embedded rich content
	PayloadJSON string        `json:"payload_json,omitempty"`
}

// ReqCreateChannelMessage [POST] Post a message to a guild text or DM channel. If operating on a guild channel,
//                                this endpoint requires the 'SEND_MESSAGES' permission to be present on the
//                                current user. If the tts field is set to true, the SEND_TTS_MESSAGES permission
//                                is required for the message to be spoken. Returns a message object. Fires a
//                                Message Create Gateway event. See message formatting for more information on
//                                how to properly format messages.
//                                The maximum request size when sending a message is 8MB.
// Endpoint                       /channels/{channel.id}/messages
// Rate limiter [MAJOR]           /channels/{channel.id}
// Discord documentation          https://discordapp.com/developers/docs/resources/channel#create-message
// Reviewed                       2018-06-10
// Comment                        Before using this endpoint, you must connect to and identify with a gateway
//                                at least once. This endpoint supports both JSON and form data bodies. It does
//                                require multipart/form-data requests instead of the normal JSON request type
//                                when uploading files. Make sure you set your Content-Type to multipart/form-data
//                                if you're doing that. Note that in that case, the embed field cannot be used,
//                                but you can pass an url-encoded JSON body as a form value for payload_json.
func ReqCreateChannelMessage(client httd.Poster, channelID snowflake.ID, params *ReqCreateMessageParams) (ret *Message, err error) {
	if channelID.Empty() {
		err = errors.New("channelID must be set to get channel messages")
		return
	}
	if params == nil {
		err = errors.New("message must be set")
		return
	}

	details := &httd.Request{
		Ratelimiter: "/channels/" + channelID.String(),
		Endpoint: "/messages",
		JSONParams: params,
	}
	resp, err := client.Post(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

// ReqEditMessageParams https://discordapp.com/developers/docs/resources/channel#edit-message-json-params
type ReqEditMessageParams struct {
	Content string        `json:"content,omitempty"`
	Embed   *ChannelEmbed `json:"embed,omitempty"` // embedded rich content
}

// ReqEditMessage [PATCH] Edit a previously sent message. You can only edit messages that have been sent by
//                        the current user. Returns a message object. Fires a Message Update Gateway event.
// Endpoint               /channels/{channel.id}/messages/{message.id}
// Rate limiter [MAJOR]   /channels/{channel.id}
// Discord documentation  https://discordapp.com/developers/docs/resources/channel#edit-message
// Reviewed               2018-06-10
// Comment                All parameters to this endpoint are optional.
func ReqEditMessage(client httd.Patcher, chanID, msgID snowflake.ID, params *ReqEditMessageParams) (ret *Message, err error) {
	if chanID.Empty() {
		err = errors.New("channelID must be set to get channel messages")
		return
	}
	if msgID.Empty() {
		err = errors.New("msgID must be set to edit the message")
		return
	}

	details := &httd.Request{
		Ratelimiter: "/channels/" + chanID.String(),
		Endpoint: "/messages/" + msgID.String(),
		JSONParams: params,
	}
	resp, err := client.Patch(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	return
}

// ReqDeleteMessage [DELETE]  Delete a message. If operating on a guild channel and trying to delete a
//                            message that was not sent by the current user, this endpoint requires the
//                            'MANAGE_MESSAGES' permission. Returns a 204 empty response on success.
//                            Fires a Message Delete Gateway event.
// Endpoint                   /channels/{channel.id}/messages/{message.id}
// Rate limiter [MAJOR]       /channels/{channel.id}
// Discord documentation      https://discordapp.com/developers/docs/resources/channel#delete-message
// Reviewed                   2018-06-10
// Comment                    -
func ReqDeleteMessage(client httd.Deleter, chanID, msgID snowflake.ID) (err error) {
	if chanID.Empty() {
		err = errors.New("channelID must be set to get channel messages")
		return
	}
	if msgID.Empty() {
		err = errors.New("msgID must be set to delete the message")
		return
	}

	details := &httd.Request{
		Ratelimiter: "/channels/" + chanID.String(),
		Endpoint: "/messages/" + msgID.String(),
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

// ReqBulkDeleteMessagesParams https://discordapp.com/developers/docs/resources/channel#bulk-delete-messages-json-params
type ReqBulkDeleteMessagesParams struct {
	Messages []snowflake.ID `json:"messages"`
	m        sync.RWMutex   `json:"-"`
}

func (p *ReqBulkDeleteMessagesParams) tooMany(messages int) (err error) {
	if messages > 100 {
		err = errors.New("must be 100 or less messages to delete")
	}

	return
}

func (p *ReqBulkDeleteMessagesParams) tooFew(messages int) (err error) {
	if messages < 2 {
		err = errors.New("must be at least two messages to delete")
	}

	return
}

func (p *ReqBulkDeleteMessagesParams) Valid() (err error) {
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

func (p *ReqBulkDeleteMessagesParams) AddMessage(msg *Message) (err error) {
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

// ReqBulkDeleteMessages [POST] Delete multiple messages in a single request. This endpoint can only be used
//                              on guild channels and requires the 'MANAGE_MESSAGES' permission. Returns a 204
//                              empty response on success. Fires multiple Message Delete Gateway events.Any message
//                              IDs given that do not exist or are invalid will count towards the minimum and
//                              maximum message count (currently 2 and 100 respectively). Additionally,
//                              duplicated IDs will only be counted once.
// Endpoint                     /channels/{channel.id}/messages/bulk-delete
// Rate limiter [MAJOR]         /channels/{channel.id}
// Discord documentation        https://discordapp.com/developers/docs/resources/channel#delete-message
// Reviewed                     2018-06-10
// Comment                      This endpoint will not delete messages older than 2 weeks, and will fail if
//                              any message provided is older than that.
func ReqBulkDeleteMessages(client httd.Poster, chanID snowflake.ID, params *ReqBulkDeleteMessagesParams) (err error) {
	if chanID.Empty() {
		err = errors.New("channelID must be set to get channel messages")
		return
	}
	err = params.Valid()
	if err != nil {
		return
	}

	details := &httd.Request{
		Ratelimiter: "/channels/" + chanID.String(),
		Endpoint: "/messages/bulk-delete",
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
