package disgord

import (
	"encoding/json"
	"errors"
	"sync"

	"time"
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

// Attachment https://discordapp.com/developers/docs/resources/channel#attachment-object
type Attachment struct {
	ID       Snowflake `json:"id"`
	Filename string    `json:"filename"`
	Size     uint      `json:"size"`
	URL      string    `json:"url"`
	ProxyURL string    `json:"proxy_url"`
	Height   uint      `json:"height"`
	Width    uint      `json:"width"`
}

// Overwrite: https://discordapp.com/developers/docs/resources/channel#overwrite-object
type PermissionOverwrite struct {
	ID    Snowflake `json:"id"`    // role or user id
	Type  string    `json:"type"`  // either `role` or `member`
	Allow int       `json:"allow"` // permission bit set
	Deny  int       `json:"deny"`  // permission bit set
}

func (pmo *PermissionOverwrite) Clear() {}

func NewChannel() *Channel {
	return &Channel{}
}

// ChannelMessager Methods required to create a new DM (or use an existing one) and send a DM.
type ChannelMessager interface {
	CreateMessage(*Message) error
}
type ChannelFetcher interface {
	GetChannel(id Snowflake) (ret *Channel, err error)
}
type ChannelDeleter interface {
	DeleteChannel(id Snowflake) (err error)
}
type ChannelUpdater interface {
}

// Channel
type Channel struct {
	ID                   Snowflake             `json:"id"`
	Type                 uint                  `json:"type"`
	GuildID              Snowflake             `json:"guild_id,omitempty"`              // ?|
	Position             uint                  `json:"position,omitempty"`              // ?|
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"` // ?|
	Name                 string                `json:"name,omitempty"`                  // ?|
	Topic                string                `json:"topic,omitempty"`                 // ?|
	NSFW                 bool                  `json:"nsfw,omitempty"`                  // ?|
	LastMessageID        Snowflake             `json:"last_message_id,omitempty"`       // ?|?, pointer
	Bitrate              uint                  `json:"bitrate,omitempty"`               // ?|
	UserLimit            uint                  `json:"user_limit,omitempty"`            // ?|
	Recipients           []*User               `json:"recipient,omitempty"`             // ?| , empty if not DM
	Icon                 string                `json:"icon,omitempty"`                  // ?|?, pointer
	OwnerID              Snowflake             `json:"owner_id,omitempty"`              // ?|
	ApplicationID        Snowflake             `json:"application_id,omitempty"`        // ?|
	ParentID             Snowflake             `json:"parent_id,omitempty"`             // ?|?, pointer
	LastPingTimestamp    Timestamp             `json:"last_ping_timestamp,omitempty"`   // ?|

	sync.RWMutex
}
type PartialChannel = Channel

func (c *Channel) Mention() string {
	return "<#" + c.ID.String() + ">"
}

func (c *Channel) Compare(other *Channel) bool {
	// eh
	return (c == nil && other == nil) || (other != nil && c.ID == other.ID)
}

func (c *Channel) SaveToDiscord(session Session) (err error) {
	if c.GuildID.Empty() {
		err = NewErrorMissingSnowflake("guild id/snowflake is empty or missing")
		return
	}
	if c.Name == "" {
		err = NewErrorEmptyValue("must have a channel name before creating channel")
	}
	params := &CreateGuildChannelParams{
		Name:                 c.Name,
		Type:                 c.Type,
		Topic:                c.Topic,
		Bitrate:              c.Bitrate,
		UserLimit:            c.UserLimit,
		PermissionOverwrites: c.PermissionOverwrites,
		ParentID:             c.ParentID,
		NSFW:                 c.NSFW,
	}
	var creation *Channel
	creation, err = session.CreateGuildChannel(c.GuildID, params)
	if err != nil {
		return
	}

	// update current channel object
	creation.CopyOverTo(c)
	return
}

func (c *Channel) DeleteFromDiscord(session Session) (err error) {
	if c.ID.Empty() {
		err = NewErrorMissingSnowflake("channel id/snowflake is empty or missing")
		return
	}
	err = session.DeleteChannel(c.ID)
	return
}

func (c *Channel) DeepCopy() (copy interface{}) {
	copy = NewChannel()
	c.CopyOverTo(copy)

	return
}

func (c *Channel) CopyOverTo(other interface{}) (err error) {
	var channel *Channel
	var valid bool
	if channel, valid = other.(*Channel); !valid {
		err = NewErrorUnsupportedType("argument given is not a *Channel type")
		return
	}

	c.RWMutex.RLock()
	channel.RWMutex.Lock()

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
	channel.Icon = c.Icon
	channel.OwnerID = c.OwnerID
	channel.ApplicationID = c.ApplicationID
	channel.ParentID = c.ParentID
	channel.LastPingTimestamp = c.LastPingTimestamp

	// add recipients if it's a DM
	if c.Type == ChannelTypeDM || c.Type == ChannelTypeGroupDM {
		for _, recipient := range c.Recipients {
			channel.Recipients = append(channel.Recipients, recipient.DeepCopy().(*User))
		}
	}

	c.RWMutex.RUnlock()
	channel.RWMutex.Unlock()

	return
}

func (c *Channel) Clear() {
	// TODO
}

// Update send channel changes to the Discord API
func (c *Channel) Update(client ChannelUpdater) {

}

// Delete sends a Discord REST request to delete the related channel
func (c *Channel) Delete(client ChannelDeleter) (err error) {
	err = client.DeleteChannel(c.ID)
	return
}

func (c *Channel) Create() {
	// check if channel already exists.
}

// Fetch check if there are any updates to the channel values
//func (c *Channel) Fetch(client ChannelFetcher) (err error) {
//	if c.ID.Empty() {
//		err = errors.New("missing channel ID")
//		return
//	}
//
//	client.GetChannel(c.ID)
//}

func (c *Channel) SendMsgString(client MessageSender, content string) (msg *Message, err error) {
	msg, err = client.SendMsgString(c.ID, content)
	return
}

func (c *Channel) SendMsg(client MessageSender, message *Message) (msg *Message, err error) {
	msg, err = client.SendMsg(c.ID, message)
	return
}

// -----------------------------
// Message

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
	ID        Snowflake `json:"id"`
	ChannelID Snowflake `json:"channel_id"`
}

// https://discordapp.com/developers/docs/resources/channel#message-object-message-activity-structure
type MessageActivity struct {
	Type    int    `json:"type"`
	PartyID string `json:"party_id"`
}

// https://discordapp.com/developers/docs/resources/channel#message-object-message-application-structure
type MessageApplication struct {
	ID          Snowflake `json:"id"`
	CoverImage  string    `json:"cover_image"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Name        string    `json:"name"`
}

// Message https://discordapp.com/developers/docs/resources/channel#message-object-message-structure
type Message struct {
	ID              Snowflake          `json:"id"`
	ChannelID       Snowflake          `json:"channel_id"`
	Author          *User              `json:"author"`
	Content         string             `json:"content"`
	Timestamp       time.Time          `json:"timestamp"`
	EditedTimestamp time.Time          `json:"edited_timestamp"` // ?
	Tts             bool               `json:"tts"`
	MentionEveryone bool               `json:"mention_everyone"`
	Mentions        []*User            `json:"mentions"`
	MentionRoles    []Snowflake        `json:"mention_roles"`
	Attachments     []*Attachment      `json:"attachments"`
	Embeds          []*ChannelEmbed    `json:"embeds"`
	Reactions       []*Reaction        `json:"reactions"` // ?
	Nonce           Snowflake          `json:"nonce"`     // ?, used for validating a message was sent
	Pinned          bool               `json:"pinned"`
	WebhookID       Snowflake          `json:"webhook_id"` // ?
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

type MessageDeleter interface {
	DeleteMessage(channelID, msgID Snowflake) (err error)
}

// Delete sends a delete request to Discord for the related message
func (m *Message) Delete(client MessageDeleter) (err error) {
	if m.ID.Empty() {
		err = errors.New("message is missing snowflake")
		return
	}

	err = client.DeleteMessage(m.ChannelID, m.ID)
	return
}

type MessageUpdater interface {
	UpdateMessage(message *Message) (msg *Message, err error)
}

// Update after changing the message object, call update to notify Discord
//        about any changes made
func (m *Message) Update(client MessageUpdater) (msg *Message, err error) {
	msg, err = client.UpdateMessage(m)
	return
}

type MessageSender interface {
	SendMsg(channelID Snowflake, message *Message) (msg *Message, err error)
	SendMsgString(channelID Snowflake, content string) (msg *Message, err error)
}

func (m *Message) Send(client MessageSender) (msg *Message, err error) {
	msg, err = client.SendMsg(m.ChannelID, m)
	return
}
func (m *Message) Respond(client MessageSender, message *Message) (msg *Message, err error) {
	message.ChannelID = m.ChannelID
	msg, err = message.Send(client)
	return
}
func (m *Message) RespondString(client MessageSender, content string) (msg *Message, err error) {
	msg, err = client.SendMsgString(m.ChannelID, content)
	return
}

func (m *Message) AddReaction(reaction *Reaction) {}
func (m *Message) RemoveReaction(id Snowflake)    {}

// ----------------
// Reaction

// https://discordapp.com/developers/docs/resources/channel#reaction-object
type Reaction struct {
	Count uint          `json:"count"`
	Me    bool          `json:"me"`
	Emoji *PartialEmoji `json:"Emoji"`
}

// -----------------
// Embed

// limitations: https://discordapp.com/developers/docs/resources/channel#embed-limits
// TODO: implement NewEmbedX functions that ensures limitations

// ChannelEmbed https://discordapp.com/developers/docs/resources/channel#embed-object
type ChannelEmbed struct {
	Title       string                 `json:"title"`       // title of embed
	Type        string                 `json:"type"`        // type of embed (always "rich" for webhook embeds)
	Description string                 `json:"description"` // description of embed
	URL         string                 `json:"url"`         // url of embed
	Timestamp   time.Time              `json:"timestamp"`   // timestamp	timestamp of embed content
	Color       int                    `json:"color"`       // color code of the embed
	Footer      *ChannelEmbedFooter    `json:"footer"`      // embed footer object	footer information
	Image       *ChannelEmbedImage     `json:"image"`       // embed image object	image information
	Thumbnail   *ChannelEmbedThumbnail `json:"thumbnail"`   // embed thumbnail object	thumbnail information
	Video       *ChannelEmbedVideo     `json:"video"`       // embed video object	video information
	Provider    *ChannelEmbedProvider  `json:"provider"`    // embed provider object	provider information
	Author      *ChannelEmbedAuthor    `json:"author"`      // embed author object	author information
	Fields      []*ChannelEmbedField   `json:"fields"`      //	array of embed field objects	fields information
}

// ChannelEmbedThumbnail https://discordapp.com/developers/docs/resources/channel#embed-object-embed-thumbnail-structure
type ChannelEmbedThumbnail struct {
	Url      string `json:"url,omitempty"`       // ?| , source url of image (only supports http(s) and attachments)
	ProxyUrl string `json:"proxy_url,omitempty"` // ?| , a proxied url of the image
	Height   int    `json:"height,omitempty"`    // ?| , height of image
	Width    int    `json:"width,omitempty"`     // ?| , width of image
}

// ChannelEmbedVideo https://discordapp.com/developers/docs/resources/channel#embed-object-embed-video-structure
type ChannelEmbedVideo struct {
	Url    string `json:"url,omitempty"`    // ?| , source url of video
	Height int    `json:"height,omitempty"` // ?| , height of video
	Width  int    `json:"width,omitempty"`  // ?| , width of video
}

// ChannelEmbedImage https://discordapp.com/developers/docs/resources/channel#embed-object-embed-image-structure
type ChannelEmbedImage struct {
	Url      string `json:"url,omitempty"`       // ?| , source url of image (only supports http(s) and attachments)
	ProxyUrl string `json:"proxy_url,omitempty"` // ?| , a proxied url of the image
	Height   int    `json:"height,omitempty"`    // ?| , height of image
	Width    int    `json:"width,omitempty"`     // ?| , width of image
}

// ChannelEmbedProvider https://discordapp.com/developers/docs/resources/channel#embed-object-embed-provider-structure
type ChannelEmbedProvider struct {
	Name string `json:"name,omitempty"` // ?| , name of provider
	Url  string `json:"url,omitempty"`  // ?| , url of provider
}

// ChannelEmbedAuthor https://discordapp.com/developers/docs/resources/channel#embed-object-embed-author-structure
type ChannelEmbedAuthor struct {
	Name         string `json:"name,omitempty"`           // ?| , name of author
	Url          string `json:"url,omitempty"`            // ?| , url of author
	IconUrl      string `json:"icon_url,omitempty"`       // ?| , url of author icon (only supports http(s) and attachments)
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"` // ?| , a proxied url of author icon
}

// ChannelEmbedFooter https://discordapp.com/developers/docs/resources/channel#embed-object-embed-footer-structure
type ChannelEmbedFooter struct {
	Text         string `json:"text"`                     //  | , url of author
	IconUrl      string `json:"icon_url,omitempty"`       // ?| , url of footer icon (only supports http(s) and attachments)
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"` // ?| , a proxied url of footer icon
}

// ChannelEmbedField https://discordapp.com/developers/docs/resources/channel#embed-object-embed-field-structure
type ChannelEmbedField struct {
	Name   string `json:"name"`           //  | , name of the field
	Value  string `json:"value"`          //  | , value of the field
	Inline bool   `json:"bool,omitempty"` // ?| , whether or not this field should display inline
}
