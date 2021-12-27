package disgord

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
)

// different message activity types
const (
	_ = iota
	MessageActivityTypeJoin
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	MessageActivityTypeJoinRequest
)

// MessageFlag https://discord.com/developers/docs/resources/channel#message-object-message-flags
type MessageFlag uint

const (
	MessageFlagCrossposted MessageFlag = 1 << iota
	MessageFlagIsCrosspost
	MessageFlagSupressEmbeds
	MessageFlagSourceMessageDeleted
	MessageFlagUrgent
	MessageFlagHasThread
	MessageFlagEphemeral
	MessageFlagLoading
)

// MessageType The different message types usually generated by Discord. eg. "a new user joined"
type MessageType uint // TODO: once auto generated, un-export this.

const (
	MessageTypeDefault MessageType = iota
	MessageTypeRecipientAdd
	MessageTypeRecipientRemove
	MessageTypeCall
	MessageTypeChannelNameChange
	MessageTypeChannelIconChange
	MessageTypeChannelPinnedMessage
	MessageTypeGuildMemberJoin
	MessageTypeUserPremiumGuildSubscription
	MessageTypeUserPremiumGuildSubscriptionTier1
	MessageTypeUserPremiumGuildSubscriptionTier2
	MessageTypeUserPremiumGuildSubscriptionTier3
	MessageTypeChannelFollowAdd
	_
	MessageTypeGuildDiscoveryDisqualified
	MessageTypeGuildDiscoveryRequalified
	_
	_
	MessageTypeThreadCreated
	MessageTypeReply
	MessageTypeApplicationCommand
	MessageTypeThreadStarterMessage
)

const (
	AttachmentSpoilerPrefix = "SPOILER_"
)

// MessageActivity https://discord.com/developers/docs/resources/channel#message-object-message-activity-structure
type MessageActivity struct {
	Type    int    `json:"type"`
	PartyID string `json:"party_id"`
}

type MentionChannel struct {
	ID      Snowflake   `json:"id"`
	GuildID Snowflake   `json:"guild_id"`
	Type    ChannelType `json:"type"`
	Name    string      `json:"name"`
}

var _ Copier = (*MentionChannel)(nil)
var _ DeepCopier = (*MentionChannel)(nil)

type MessageReference struct {
	MessageID Snowflake `json:"message_id"`
	ChannelID Snowflake `json:"channel_id"`
	GuildID   Snowflake `json:"guild_id"`
}

type MessageComponentType = int

const (
	_ MessageComponentType = iota
	MessageComponentActionRow
	MessageComponentButton
)

type ButtonStyle = int

const (
	_ ButtonStyle = iota
	Primary
	Secondary
	Success
	Danger
	Link
)

type MessageComponent struct {
	Type       MessageComponentType `json:"type"`
	Style      ButtonStyle          `json:"style"`
	Label      string               `json:"label"`
	Emoji      *Emoji               `json:"emoji"`
	CustomID   string               `json:"custom_id"`
	Url        string               `json:"url"`
	Disabled   bool                 `json:"disabled"`
	Components []*MessageComponent  `json:"components"`
}

var _ Copier = (*MessageComponent)(nil)
var _ DeepCopier = (*MessageComponent)(nil)

// MessageApplication https://discord.com/developers/docs/resources/channel#message-object-message-application-structure
type MessageApplication struct {
	ID          Snowflake `json:"id"`
	CoverImage  string    `json:"cover_image"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Name        string    `json:"name"`
}

type MessageStickerFormatType int

const (
	_ MessageStickerFormatType = iota
	MessageStickerFormatPNG
	MessageStickerFormatAPNG
	MessageStickerFormatLOTTIE
)

type StickerItem struct {
	ID         Snowflake                `json:"id"`
	Name       string                   `json:"name"`
	FormatType MessageStickerFormatType `json:"format_type"`
}

var _ Copier = (*StickerItem)(nil)
var _ DeepCopier = (*StickerItem)(nil)

type MessageSticker struct {
	ID           Snowflake                `json:"id"`
	PackID       Snowflake                `json:"pack_id"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	Tags         string                   `json:"tags"`
	Asset        string                   `json:"asset"`
	PreviewAsset string                   `json:"preview_asset"`
	FormatType   MessageStickerFormatType `json:"format_type"`
}

var _ Copier = (*MessageSticker)(nil)
var _ DeepCopier = (*MessageSticker)(nil)

// Message https://discord.com/developers/docs/resources/channel#message-object-message-structure
type Message struct {
	ID                Snowflake           `json:"id"`
	ChannelID         Snowflake           `json:"channel_id"`
	GuildID           Snowflake           `json:"guild_id"`
	Author            *User               `json:"author"`
	Member            *Member             `json:"member"`
	Content           string              `json:"content"`
	Timestamp         Time                `json:"timestamp"`
	EditedTimestamp   Time                `json:"edited_timestamp"` // ?
	Tts               bool                `json:"tts"`
	MentionEveryone   bool                `json:"mention_everyone"`
	Mentions          []*User             `json:"mentions"`
	MentionRoles      []Snowflake         `json:"mention_roles"`
	MentionChannels   []*MentionChannel   `json:"mention_channels"`
	Attachments       []*Attachment       `json:"attachments"`
	Embeds            []*Embed            `json:"embeds"`
	Reactions         []*Reaction         `json:"reactions"` // ?
	Nonce             interface{}         `json:"nonce"`     // NOT A SNOWFLAKE! DONT TOUCH!
	Pinned            bool                `json:"pinned"`
	WebhookID         Snowflake           `json:"webhook_id"` // ?
	Type              MessageType         `json:"type"`
	Activity          MessageActivity     `json:"activity"`
	Application       MessageApplication  `json:"application"`
	MessageReference  *MessageReference   `json:"message_reference"`
	ReferencedMessage *Message            `json:"referenced_message"`
	Flags             MessageFlag         `json:"flags"`
	StickerItems      []*StickerItem      `json:"sticker_items"`
	Components        []*MessageComponent `json:"components"`
	Interaction       *MessageInteraction `json:"interaction"`
	// SpoilerTagContent is only true if the entire message text is tagged as a spoiler (aka completely wrapped in ||)
	SpoilerTagContent        bool `json:"-"`
	SpoilerTagAllAttachments bool `json:"-"`
	HasSpoilerImage          bool `json:"-"`
}

var _ Reseter = (*Message)(nil)
var _ fmt.Stringer = (*Message)(nil)
var _ internalUpdater = (*Message)(nil)
var _ Copier = (*Message)(nil)
var _ DeepCopier = (*Message)(nil)

func (m *Message) String() string {
	return "message{" + m.ID.String() + "}"
}

// DiscordURL returns the Discord link to the message. This can be used to jump
// directly to a message from within the client.
//
// Example: https://discord.com/channels/319567980491046913/644376487331495967/646925626523254795
func (m *Message) DiscordURL() (string, error) {
	if m.ID.IsZero() {
		return "", errors.New("missing message ID")
	}
	if m.GuildID.IsZero() {
		return "", errors.New("missing guild ID")
	}
	if m.ChannelID.IsZero() {
		return "", errors.New("missing channel ID")
	}

	return fmt.Sprintf(
		"https://discord.com/channels/%d/%d/%d",
		m.GuildID, m.ChannelID, m.ID,
	), nil
}

func (m *Message) updateInternals() {
	if len(m.Content) >= len("||||") {
		prefix := m.Content[0:2]
		suffix := m.Content[len(m.Content)-2 : len(m.Content)]
		m.SpoilerTagContent = prefix+suffix == "||||"
	}

	m.SpoilerTagAllAttachments = len(m.Attachments) > 0
	for i := range m.Attachments {
		m.Attachments[i].updateInternals()
		if !m.Attachments[i].SpoilerTag {
			m.SpoilerTagAllAttachments = false
			break
		} else {
			m.HasSpoilerImage = true
		}
	}

	if m.Author != nil && m.Member != nil {
		m.Member.UserID = m.Author.ID
	}
}

// IsDirectMessage checks if the message is from a direct message channel.
//
// WARNING! Note that, when fetching messages using the REST API the
// guildID might be empty -> giving a false positive.
func (m *Message) IsDirectMessage() bool {
	return m.Type == MessageTypeDefault && m.GuildID.IsZero()
}

// Send sends this message to discord.
func (m *Message) Send(ctx context.Context, s Session) (msg *Message, err error) {
	nonce := fmt.Sprint(m.Nonce)
	if len(nonce) > 25 {
		return nil, errors.New("nonce can not be more than 25 characters")
	}

	// TODO: attachments
	params := &CreateMessage{
		Content:          m.Content,
		Tts:              m.Tts,
		MessageReference: m.MessageReference,
		Nonce:            nonce,
		// File: ...
		// Embed: ...
	}
	if len(m.Embeds) > 0 {
		params.Embed = &Embed{}
		_ = DeepCopyOver(params.Embed, m.Embeds[0])
	}
	channelID := m.ChannelID

	msg, err = s.Channel(channelID).WithContext(ctx).CreateMessage(params)
	return
}

// Reply input any type as an reply. int, string, an object, etc.
func (m *Message) Reply(ctx context.Context, s Session, data ...interface{}) (*Message, error) {
	return s.WithContext(ctx).SendMsg(m.ChannelID, data...)
}

func (m *Message) React(ctx context.Context, s Session, emoji interface{}) error {
	if m.ID.IsZero() {
		return errors.New("missing message ID")
	} else if m.ChannelID.IsZero() {
		return errors.New("missing channel ID")
	}

	return s.Channel(m.ChannelID).Message(m.ID).Reaction(emoji).WithContext(ctx).Create()
}

func (m *Message) Unreact(ctx context.Context, s Session, emoji interface{}) error {
	if m.ID.IsZero() {
		return errors.New("missing message ID")
	} else if m.ChannelID.IsZero() {
		return errors.New("missing channel ID")
	}

	return s.Channel(m.ChannelID).Message(m.ID).Reaction(emoji).WithContext(ctx).DeleteOwn()
}

// AddReaction adds a reaction to the message
//func (m *Message) AddReaction(reaction *Reaction) {}

// RemoveReaction removes a reaction from the message
//func (m *Message) RemoveReaction(id Snowflake)    {}

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

type MessageQueryBuilder interface {
	WithContext(ctx context.Context) MessageQueryBuilder
	WithFlags(flags ...Flag) MessageQueryBuilder

	// Pin Pin a message by its ID and channel ID. Requires the 'MANAGE_MESSAGES' permission.
	Pin() error

	// Unpin Delete a pinned message in a channel. Requires the 'MANAGE_MESSAGES' permission.
	Unpin() error

	// Get Returns a specific message in the channel. If operating on a guild channel, this endpoints
	// requires the 'READ_MESSAGE_HISTORY' permission to be present on the current user.
	// Returns a message object on success.
	Get() (*Message, error)

	// UpdateBuilder Edit a previously sent message. You can only edit messages that have been sent by the
	// current user. Returns a message object. Fires a Message Update Gateway event.
	UpdateBuilder() UpdateMessageBuilder
	SetContent(content string) (*Message, error)
	SetEmbed(embed *Embed) (*Message, error)

	CrossPost() (*Message, error)

	// Delete Delete a message. If operating on a guild channel and trying to delete a message that was not
	// sent by the current user, this endpoint requires the 'MANAGE_MESSAGES' permission. Fires a Message Delete Gateway event.
	Delete() error

	// DeleteAllReactions Deletes all reactions on a message. This endpoint requires the 'MANAGE_MESSAGES'
	// permission to be present on the current user.
	DeleteAllReactions() error

	Reaction(emoji interface{}) ReactionQueryBuilder
}

func (c channelQueryBuilder) Message(id Snowflake) MessageQueryBuilder {
	return &messageQueryBuilder{client: c.client, cid: c.cid, mid: id}
}

type messageQueryBuilder struct {
	ctx    context.Context
	flags  Flag
	client *Client
	cid    Snowflake
	mid    Snowflake
}

func (m messageQueryBuilder) WithContext(ctx context.Context) MessageQueryBuilder {
	m.ctx = ctx
	return &m
}

func (m messageQueryBuilder) WithFlags(flags ...Flag) MessageQueryBuilder {
	m.flags = mergeFlags(flags)
	return &m
}

// Get Returns a specific message in the channel. If operating on a guild channel, this endpoints
// requires the 'READ_MESSAGE_HISTORY' permission to be present on the current user.
// Returns a message object on success.
//  Method                  GET
//  Endpoint                /channels/{channel.id}/messages/{message.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#get-channel-message
//  Reviewed                2018-06-10
//  Comment                 -
func (m messageQueryBuilder) Get() (*Message, error) {
	if m.cid.IsZero() {
		err := errors.New("channelID must be set to get channel messages")
		return nil, err
	}
	if m.mid.IsZero() {
		err := errors.New("messageID must be set to get a specific message from a channel")
		return nil, err
	}

	if !ignoreCache(m.flags) {
		if msg, _ := m.client.cache.GetMessage(m.cid, m.mid); msg != nil {
			return msg, nil
		}
	}

	r := m.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.ChannelMessage(m.cid, m.mid),
		Ctx:      m.ctx,
	}, m.flags)
	r.pool = m.client.pool.message
	r.factory = func() interface{} {
		return &Message{}
	}

	return getMessage(r.Execute)
}

// UpdateBuilder Edit a previously sent message. You can only edit messages that have been sent by the
// current user. Returns a message object. Fires a Message Update Gateway event.
//  Method                  PATCH
//  Endpoint                /channels/{channel.id}/messages/{message.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#edit-message
//  Reviewed                2018-06-10
//  Comment                 All parameters to this endpoint are optional.
func (m messageQueryBuilder) UpdateBuilder() UpdateMessageBuilder {
	builder := &updateMessageBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Message{}
	}
	builder.r.flags = m.flags
	builder.r.addPrereq(m.cid.IsZero(), "channelID must be set to get channel messages")
	builder.r.addPrereq(m.mid.IsZero(), "msgID must be set to edit the message")
	builder.r.setup(m.client.req, &httd.Request{
		Method:      http.MethodPatch,
		Ctx:         m.ctx,
		Endpoint:    "/channels/" + m.cid.String() + "/messages/" + m.mid.String(),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// Delete If operating on a guild channel and trying to delete a message that was not
// sent by the current user, this endpoint requires the 'MANAGE_MESSAGES' permission. Returns a 204 empty response
// on success. Fires a Message Delete Gateway event.
//  Method                  DELETE
//  Endpoint                /channels/{channel.id}/messages/{message.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#delete-message
//  Reviewed                2018-06-10
//  Comment                 -
func (m messageQueryBuilder) Delete() (err error) {
	if m.cid.IsZero() {
		err = errors.New("channelID must be set to get channel messages")
		return
	}
	if m.mid.IsZero() {
		err = errors.New("msgID must be set to delete the message")
		return
	}

	r := m.client.newRESTRequest(&httd.Request{
		Method:   http.MethodDelete,
		Endpoint: endpoint.ChannelMessage(m.cid, m.mid),
		Ctx:      m.ctx,
	}, m.flags)

	_, err = r.Execute()
	return err
}

// Pin a message by its ID and channel ID. Requires the 'MANAGE_MESSAGES' permission.
// Returns a 204 empty response on success.
//  Method                  PUT
//  Endpoint                /channels/{channel.id}/pins/{message.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#add-pinned-channel-message
//  Reviewed                2018-06-10
//  Comment                 -
func (m messageQueryBuilder) Pin() (err error) {
	r := m.client.newRESTRequest(&httd.Request{
		Method:   http.MethodPut,
		Endpoint: endpoint.ChannelPin(m.cid, m.mid),
		Ctx:      m.ctx,
	}, m.flags)

	_, err = r.Execute()
	return err
}

// Unpin [REST] Delete a pinned message in a channel. Requires the 'MANAGE_MESSAGES' permission.
// Returns a 204 empty response on success. Returns a 204 empty response on success.
//  Method                  DELETE
//  Endpoint                /channels/{channel.id}/pins/{message.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#delete-pinned-channel-message
//  Reviewed                2018-06-10
//  Comment                 -
func (m messageQueryBuilder) Unpin() (err error) {
	if m.cid.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if m.mid.IsZero() {
		return errors.New("messageID must be set to target the specific channel message")
	}

	r := m.client.newRESTRequest(&httd.Request{
		Method:   http.MethodDelete,
		Endpoint: endpoint.ChannelPin(m.cid, m.mid),
		Ctx:      m.ctx,
	}, m.flags)

	_, err = r.Execute()
	return err
}

// CrossPost Crosspost a message in a News Channel to following channels.
//  Method                  POST
//  Endpoint                /channels/{channel.id}/messages/{message.id}/crosspost
//  Discord documentation   https://discord.com/developers/docs/resources/channel#crosspost-message
//  Reviewed                2021-04-07
//  Comment                 -
func (m messageQueryBuilder) CrossPost() (*Message, error) {
	if m.cid.IsZero() {
		return nil, errors.New("channelID must be set to target the correct channel")
	}
	if m.mid.IsZero() {
		return nil, errors.New("messageID must be set to target the specific channel message")
	}

	r := m.client.newRESTRequest(&httd.Request{
		Method:   http.MethodPost,
		Endpoint: endpoint.ChannelMessageCrossPost(m.cid, m.mid),
		Ctx:      m.ctx,
	}, m.flags)
	r.pool = m.client.pool.message
	r.factory = func() interface{} {
		return &Message{}
	}

	msg, err := r.Execute()
	if err != nil {
		return nil, err
	}
	return msg.(*Message), nil
}

// DeleteAllReactions [REST] Deletes all reactions on a message. This endpoint requires the 'MANAGE_MESSAGES'
// permission to be present on the current user.
//  Method                  DELETE
//  Endpoint                /channels/{channel.id}/messages/{message.id}/reactions
//  Discord documentation   https://discord.com/developers/docs/resources/channel#delete-all-reactions
//  Reviewed                2019-01-28
func (m messageQueryBuilder) DeleteAllReactions() error {
	if m.cid.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if m.mid.IsZero() {
		return errors.New("messageID must be set to target the specific channel message")
	}

	r := m.client.newRESTRequest(&httd.Request{
		Method:   http.MethodDelete,
		Endpoint: endpoint.ChannelMessageReactions(m.cid, m.mid),
		Ctx:      m.ctx,
	}, m.flags)

	_, err := r.Execute()
	return err
}

//////////////////////////////////////////////////////
//
// REST Wrappers
//
//////////////////////////////////////////////////////

func (m messageQueryBuilder) SetContent(content string) (*Message, error) {
	builder := m.WithContext(m.ctx).UpdateBuilder()
	return builder.
		SetContent(content).
		Execute()
}

func (m messageQueryBuilder) SetEmbed(embed *Embed) (*Message, error) {
	builder := m.WithContext(m.ctx).UpdateBuilder()
	return builder.
		SetEmbed(embed).
		Execute()
}

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

// updateMessageBuilder, params here
//  https://discord.com/developers/docs/resources/channel#edit-message-json-params
//generate-rest-params: content:string, embed:*Embed,
//generate-rest-basic-execute: message:*Message,
type updateMessageBuilder struct {
	r RESTBuilder
}

// SetAllowedMentions sets the allowed mentions for the updateMessageBuilder then returns the builder to allow chaining.
func (b *updateMessageBuilder) SetAllowedMentions(mentions *AllowedMentions) *updateMessageBuilder {
	b.r.param("allowed_mentions", mentions)
	return b
}
