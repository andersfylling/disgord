package resource

import (
	"errors"

	"net/http"

	"time"

	"encoding/json"
	"sync"

	"github.com/andersfylling/disgord/discord"
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

func NewChannel() *Channel {
	return &Channel{}
}

type Channel struct {
	ID                   snowflake.ID                 `json:"id"`
	Type                 uint                         `json:"type"`
	GuildID              snowflake.ID                 `json:"guild_id,omitempty"`              // ?|
	Position             uint                         `json:"position,omitempty"`              // ?|
	PermissionOverwrites []ChannelPermissionOverwrite `json:"permission_overwrites,omitempty"` // ?|
	Name                 string                       `json:"name,omitempty"`                  // ?|
	Topic                string                       `json:"topic,omitempty"`                 // ?|
	NSFW                 bool                         `json:"nsfw,omitempty"`                  // ?|
	LastMessageID        snowflake.ID                 `json:"last_message_id,omitempty"`       // ?|?, pointer
	Bitrate              uint                         `json:"bitrate,omitempty"`               // ?|
	UserLimit            uint                         `json:"user_limit,omitempty"`            // ?|
	Recipients           []*User                      `json:"recipient,omitempty"`             // ?| , empty if not DM
	Icon                 string                       `json:"icon,omitempty"`                  // ?|?, pointer
	OwnerID              snowflake.ID                 `json:"owner_id,omitempty"`              // ?|
	ApplicationID        snowflake.ID                 `json:"applicaiton_id,omitempty"`        // ?|
	ParentID             snowflake.ID                 `json:"parent_id,omitempty"`             // ?|?, pointer
	LastPingTimestamp    discord.Timestamp            `json:"last_ping_timestamp,omitempty"`   // ?|

	mu sync.RWMutex `json:"-"`
}

func (c *Channel) Mention() string {
	return "<#" + c.ID.String() + ">"
}

func (c *Channel) Compare(other *Channel) bool {
	// eh
	return (c == nil && other == nil) || (other != nil && c.ID == other.ID)
}

func (c *Channel) Replicate(channel *Channel, recipients []*User) {
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

// DISCORD HTTP API
// /channels/*
//
type DiscordAPIRequester interface {
	Request(method string, uri string, content interface{}) error
}

// GetChannel Get a channel by ID
func GetChannel(client DiscordAPIRequester, id snowflake.ID) (*Channel, error) {
	if id.Empty() {
		return nil, errors.New("not a valid snowflake")
	}

	uri := "/channels/" + id.String()
	content := &Channel{}
	err := client.Request(http.MethodGet, uri, content)
	return content, err
}

func UpdateChannel(client DiscordAPIRequester, changes *Channel) (*Channel, error) {
	if changes.ID.Empty() {
		return nil, errors.New("not a valid snowflake")
	}

	//uri := "/channels/" + changes.ID.String()
	//data, err := json.Marshal(changes)
	//if err != nil {
	//	return nil, err
	//}
	//err := client.Request("PUT", uri, bytes.NewBuffer(data)) // TODO implement "PUT" logic
	return nil, nil
}

func DeleteChannel(client DiscordAPIRequester, id snowflake.ID) (*Channel, error) {
	if id.Empty() {
		return nil, errors.New("not a valid snowflake")
	}

	uri := "/channels/" + id.String()
	content := &Channel{}
	err := client.Request("DELETE", uri, content)
	return content, err
}

// ---------
// ChannelEmbed ...
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

type ChannelEmbedFooter struct{}
type ChannelEmbedImage struct{}
type ChannelEmbedThumbnail struct{}
type ChannelEmbedVideo struct{}
type ChannelEmbedProvider struct{}
type ChannelEmbedAuthor struct{}
type ChannelEmbedField struct{}

// -------

const (
	_ int = iota
	MessageActivityTypeJoin
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	MessageActivityTypeJoinRequest
)

func NewMessage() *Message {
	return &Message{}
}

func NewDeletedMessage() *DeletedMessage {
	return &DeletedMessage{}
}

type Attachment struct {
	ID       snowflake.ID `json:"id"`
	Filename string       `json:"filename"`
	Size     uint         `json:"size"`
	URL      string       `json:"url"`
	ProxyURL string       `json:"proxy_url"`
	Height   uint         `json:"height"`
	Width    uint         `json:"width"`
}

type DeletedMessage struct {
	ID        snowflake.ID `json:"id"`
	ChannelID snowflake.ID `json:"channel_id"`
}

type MessageActivity struct {
	Type    int    `json:"type"`
	PartyID string `json:"party_id"`
}

type MessageApplication struct {
	ID          snowflake.ID `json:"id"`
	CoverImage  string       `json:"cover_image"`
	Description string       `json:"description"`
	Icon        string       `json:"icon"`
	Name        string       `json:"name"`
}

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

	return json.Marshal(Message(*m))
}

func (m *Message) Delete() {}
func (m *Message) Update() {}
func (m *Message) Send()   {}

func (m *Message) AddReaction(reaction *Reaction) {}
func (m *Message) RemoveReaction(id snowflake.ID) {}

// GET, based on ID? 0.o

// func (m *Message) UnmarshalJSON(data []byte) error {
// 	return json.Unmarshal(data, &m.messageJSON)
// }

func GetMessages() {}

// ---------------

type ChannelPermissionOverwrite struct {
	ID    snowflake.ID `json:"id"`    // role or user id
	Type  string       `json:"type"`  // either `role` or `member`
	Deny  int          `json:"deny"`  // permission bit set
	Allow int          `json:"allow"` // permission bit set
}

func (pmo *ChannelPermissionOverwrite) Clear() {}

// -----------

type Reaction struct {
	Count uint   `json:"count"`
	Me    bool   `json:"me"`
	Emoji *Emoji `json:"Emoji"`
}
