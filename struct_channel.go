package disgord

import (
	"sync"

	"time"
)

// Channel types
// https://discordapp.com/developers/docs/resources/channel#channel-object-channel-types
const (
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

func NewPartialChannel(id Snowflake, name string, t uint) *PartialChannel {
	return &PartialChannel{
		ID:      id,
		Name:    name,
		Type:    t,
		partial: true,
	}
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

// PartialChannel ...
// example of partial channel
// // "channel": {
// //   "id": "165176875973476352",
// //   "name": "illuminati",
// //   "type": 0
// // }
type PartialChannel = Channel

type channelJSON struct {
	ID                   Snowflake             `json:"id"`
	Type                 uint                  `json:"type"`
	GuildID              Snowflake             `json:"guild_id,omitempty"`              // ?|
	Position             uint                  `json:"position,omitempty"`              // ?|
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"` // ?|
	Name                 string                `json:"name,omitempty"`                  // ?|
	Topic                *string               `json:"topic,omitempty"`                 // ?|?, pointer
	NSFW                 bool                  `json:"nsfw,omitempty"`                  // ?|
	LastMessageID        *Snowflake            `json:"last_message_id,omitempty"`       // ?|?, pointer
	Bitrate              uint                  `json:"bitrate,omitempty"`               // ?|
	UserLimit            uint                  `json:"user_limit,omitempty"`            // ?|
	RateLimitPerUser     uint                  `json:"rate_limit_per_user,omitempty"`   // ?|
	Recipients           []*User               `json:"recipient,omitempty"`             // ?| , empty if not DM
	Icon                 *string               `json:"icon,omitempty"`                  // ?|?, pointer
	OwnerID              Snowflake             `json:"owner_id,omitempty"`              // ?|
	ApplicationID        Snowflake             `json:"application_id,omitempty"`        // ?|
	ParentID             *Snowflake            `json:"parent_id,omitempty"`             // ?|?, pointer
	LastPinTimestamp     Timestamp             `json:"last_pin_timestamp,omitempty"`    // ?|
}

// Channel ...
type Channel struct {
	sync.RWMutex         `json:"-"`
	ID                   Snowflake             `json:"id"`
	Type                 uint                  `json:"type"`
	GuildID              Snowflake             `json:"guild_id,omitempty"`              // ?|
	Position             uint                  `json:"position,omitempty"`              // ?|
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"` // ?|
	Name                 string                `json:"name,omitempty"`                  // ?|
	Topic                string                `json:"topic,omitempty"`                 // ?|?, pointer
	NSFW                 bool                  `json:"nsfw,omitempty"`                  // ?|
	LastMessageID        Snowflake             `json:"last_message_id,omitempty"`       // ?|?, pointer
	Bitrate              uint                  `json:"bitrate,omitempty"`               // ?|
	UserLimit            uint                  `json:"user_limit,omitempty"`            // ?|
	RateLimitPerUser     uint                  `json:"rate_limit_per_user,omitempty"`   // ?|
	Recipients           []*User               `json:"recipient,omitempty"`             // ?| , empty if not DM
	Icon                 string                `json:"icon,omitempty"`                  // ?|?, pointer
	OwnerID              Snowflake             `json:"owner_id,omitempty"`              // ?|
	ApplicationID        Snowflake             `json:"application_id,omitempty"`        // ?|
	ParentID             Snowflake             `json:"parent_id,omitempty"`             // ?|?, pointer
	LastPinTimestamp     Timestamp             `json:"last_pin_timestamp,omitempty"`    // ?|

	// set to true when the object is incomplete. Used in situations
	// like cache to avoid overwriting correct information.
	// A partial channel is assumed to be
	//  "channel": {
	//    "id": "165176875973476352",
	//    "name": "illuminati",
	//    "type": 0
	//  }
	partial bool
}

func (c *Channel) UnmarshalJSON(data []byte) (err error) {
	channel := channelJSON{}
	err = unmarshal(data, &channel)
	if err != nil {
		return
	}

	// copy over to c
	c.Lock()

	c.ID = channel.ID
	c.Type = channel.Type
	c.GuildID = channel.GuildID
	c.Position = channel.Position
	c.PermissionOverwrites = channel.PermissionOverwrites // TODO: check for pointer
	c.Name = channel.Name
	c.NSFW = channel.NSFW
	c.Bitrate = channel.Bitrate
	c.UserLimit = channel.UserLimit
	c.RateLimitPerUser = channel.RateLimitPerUser
	c.Recipients = channel.Recipients
	c.OwnerID = channel.OwnerID
	c.ApplicationID = channel.ApplicationID
	c.LastPinTimestamp = channel.LastPinTimestamp

	if channel.Topic != nil {
		c.Topic = *channel.Topic
	}
	if channel.Icon != nil {
		c.Icon = *channel.Icon
	}
	if channel.ParentID != nil {
		c.ParentID = *channel.ParentID
	}
	if channel.LastMessageID != nil {
		c.LastMessageID = *channel.LastMessageID
	}

	c.Unlock()
	return
}

func (c *Channel) MarshalJSON() (data []byte, err error) {
	// copy over to channel
	channel := channelJSON{}
	c.RLock()

	channel.ID = c.ID
	channel.Type = c.Type
	channel.GuildID = c.GuildID
	channel.Position = c.Position
	channel.PermissionOverwrites = c.PermissionOverwrites // TODO: check for pointer
	channel.Name = c.Name
	channel.NSFW = c.NSFW
	channel.Bitrate = c.Bitrate
	channel.UserLimit = c.UserLimit
	channel.RateLimitPerUser = c.RateLimitPerUser
	channel.Recipients = c.Recipients
	channel.OwnerID = c.OwnerID
	channel.ApplicationID = c.ApplicationID
	channel.LastPinTimestamp = c.LastPinTimestamp

	if c.Topic != "" {
		channel.Topic = &c.Topic
	}
	if c.Icon != "" {
		channel.Icon = &c.Icon
	}
	if !c.ParentID.Empty() {
		channel.ParentID = &c.ParentID
	}
	if !c.LastMessageID.Empty() {
		channel.LastMessageID = &c.LastMessageID
	}

	c.RUnlock()

	data, err = marshal(&channel)
	return
}

func (c *Channel) valid() bool {
	if c.RateLimitPerUser > 120 {
		return false
	}

	if len(c.Topic) > 1024 {
		return false
	}

	if len(c.Name) > 100 || len(c.Name) < 2 {
		return false
	}

	return true
}

func (c *Channel) Mention() string {
	return "<#" + c.ID.String() + ">"
}

func (c *Channel) Compare(other *Channel) bool {
	// eh
	return (c == nil && other == nil) || (other != nil && c.ID == other.ID)
}

func (c *Channel) saveToDiscord(session Session) (err error) {
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
		RateLimitPerUser:     c.RateLimitPerUser,
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

func (c *Channel) deleteFromDiscord(session Session) (err error) {
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
	channel.RateLimitPerUser = c.RateLimitPerUser
	channel.Icon = c.Icon
	channel.OwnerID = c.OwnerID
	channel.ApplicationID = c.ApplicationID
	channel.ParentID = c.ParentID
	channel.LastPinTimestamp = c.LastPinTimestamp
	channel.LastMessageID = c.LastMessageID

	// add recipients if it's a DM
	for _, recipient := range c.Recipients {
		channel.Recipients = append(channel.Recipients, recipient.DeepCopy().(*User))
	}

	c.RWMutex.RUnlock()
	channel.RWMutex.Unlock()

	return
}

func (c *Channel) copyOverToCache(other interface{}) (err error) {
	channel := other.(*Channel)

	channel.Lock()
	c.RLock()

	if !c.ID.Empty() {
		channel.ID = c.ID
	}
	if channel.Type == 0 && c.Type > 0 {
		// if channel type is not set(?) then it can be overwritten
		channel.Type = c.Type
	}
	if !c.GuildID.Empty() {
		channel.GuildID = c.GuildID
	}
	channel.Position = c.Position // TODO: how to avoid an partial channel to overwrite this?
	channel.PermissionOverwrites = c.PermissionOverwrites
	channel.Name = c.Name
	channel.Topic = c.Topic
	channel.NSFW = c.NSFW
	channel.LastMessageID = c.LastMessageID
	channel.Bitrate = c.Bitrate
	channel.UserLimit = c.UserLimit
	channel.RateLimitPerUser = c.RateLimitPerUser
	channel.Icon = c.Icon
	channel.OwnerID = c.OwnerID
	if !c.ApplicationID.Empty() {
		channel.ApplicationID = c.ApplicationID
	}
	if !c.ParentID.Empty() {
		channel.ParentID = c.ParentID
	}
	if !c.LastPinTimestamp.Empty() {
		channel.LastPinTimestamp = c.LastPinTimestamp
	}
	if !c.LastMessageID.Empty() {
		channel.LastMessageID = c.LastMessageID
	}

	channel.Unlock()
	c.RUnlock()
	return
}

//func (c *Channel) Clear() {
//	// TODO
//}

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
	if c.ID.Empty() {
		err = NewErrorMissingSnowflake("snowflake ID not set for channel")
		return
	}
	params := &CreateChannelMessageParams{
		Content: content,
	}

	msg, err = client.CreateChannelMessage(c.ID, params)
	return
}

func (c *Channel) SendMsg(client MessageSender, message *Message) (msg *Message, err error) {
	if c.ID.Empty() {
		err = NewErrorMissingSnowflake("snowflake ID not set for channel")
		return
	}
	message.RLock()
	params := &CreateChannelMessageParams{
		Content: message.Content,
		Nonce:   *message.Nonce,
		Tts:     message.Tts,
		// File: ...
		// Embed: ...
	}
	if len(message.Embeds) > 0 {
		params.Embed = message.Embeds[0]
	}
	message.RUnlock()

	msg, err = client.CreateChannelMessage(c.ID, params)
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
	Reactions       []*Reaction        `json:"reactions"`       // ?
	Nonce           *Snowflake         `json:"nonce,omitempty"` // ?, used for validating a message was sent
	Pinned          bool               `json:"pinned"`
	WebhookID       Snowflake          `json:"webhook_id"` // ?
	Type            uint               `json:"type"`
	Activity        MessageActivity    `json:"activity"`
	Application     MessageApplication `json:"application"`

	sync.RWMutex `json:"-"`
}

// TODO: why is this method needed?
//func (m *Message) MarshalJSON() ([]byte, error) {
//	if m.ID.Empty() {
//		return []byte("{}"), nil
//	}
//
//	//TODO: remove copying of mutex
//	return json.Marshal(Message(*m))
//}

// TODO: await for caching
//func (m *Message) DirectMessage(session Session) bool {
//	return m.Type == ChannelTypeDM
//}

type messageDeleter interface {
	DeleteMessage(channelID, msgID Snowflake) (err error)
}

func (m *Message) CopyOverTo(other interface{}) (err error) {
	var message *Message
	var valid bool
	if message, valid = other.(*Message); !valid {
		err = NewErrorUnsupportedType("argument given is not a *Message type")
		return
	}

	m.RLock()
	message.Lock()

	message.ID = m.ID
	message.ChannelID = m.ChannelID
	message.Content = m.Content
	message.Timestamp = m.Timestamp
	message.EditedTimestamp = m.EditedTimestamp
	message.Tts = m.Tts
	message.MentionEveryone = m.MentionEveryone
	message.MentionRoles = m.MentionRoles
	message.Pinned = m.Pinned
	message.WebhookID = m.WebhookID
	message.Type = m.Type
	message.Activity = m.Activity
	message.Application = m.Application

	if m.Author != nil {
		message.Author = m.Author.DeepCopy().(*User)
	}

	if m.Nonce != nil {
		s := *m.Nonce
		message.Nonce = &s
	}

	for _, mention := range m.Mentions {
		message.Mentions = append(message.Mentions, mention.DeepCopy().(*User))
	}

	for _, attachment := range m.Attachments {
		message.Attachments = append(message.Attachments, attachment.DeepCopy().(*Attachment))
	}

	for _, embed := range m.Embeds {
		message.Embeds = append(message.Embeds, embed.DeepCopy().(*ChannelEmbed))
	}

	for _, reaction := range m.Reactions {
		message.Reactions = append(message.Reactions, reaction.DeepCopy().(*Reaction))
	}

	m.RUnlock()
	message.Unlock()

	return
}

func (m *Message) deleteFromDiscord(session Session) (err error) {
	if m.ID.Empty() {
		err = NewErrorMissingSnowflake("message is missing snowflake")
		return
	}

	err = session.DeleteMessage(m.ChannelID, m.ID)
	return
}
func (m *Message) saveToDiscord(session Session) (err error) {
	var message *Message
	if m.ID.Empty() {
		message, err = m.Send(session)
	} else {
		message, err = m.update(session)
	}

	message.CopyOverTo(m)
	return
}

type MessageUpdater interface {
	UpdateMessage(message *Message) (msg *Message, err error)
}

// Update after changing the message object, call update to notify Discord
//        about any changes made
func (m *Message) update(client MessageUpdater) (msg *Message, err error) {
	msg, err = client.UpdateMessage(m)
	return
}

type MessageSender interface {
	CreateChannelMessage(channelID Snowflake, params *CreateChannelMessageParams) (ret *Message, err error)
}

func (m *Message) Send(client MessageSender) (msg *Message, err error) {
	m.RLock()
	params := &CreateChannelMessageParams{
		Content: m.Content,
		Tts:     m.Tts,
		// File: ...
		// Embed: ...
	}
	if m.Nonce != nil {
		params.Nonce = *m.Nonce
	}
	if len(m.Embeds) > 0 {
		params.Embed = m.Embeds[0]
	}
	channelID := m.ChannelID
	m.RUnlock()

	msg, err = client.CreateChannelMessage(channelID, params)
	return
}
func (m *Message) Respond(client MessageSender, message *Message) (msg *Message, err error) {
	m.RLock()
	message.Lock()
	message.ChannelID = m.ChannelID
	m.RUnlock()
	msg, err = message.Send(client)
	message.Unlock()
	return
}
func (m *Message) RespondString(client MessageSender, content string) (msg *Message, err error) {
	params := &CreateChannelMessageParams{
		Content: content,
	}

	m.RLock()
	defer m.RUnlock()
	msg, err = client.CreateChannelMessage(m.ChannelID, params)
	return
}

func (m *Message) AddReaction(reaction *Reaction) {}
func (m *Message) RemoveReaction(id Snowflake)    {}

// ----------------
// Reaction

// https://discordapp.com/developers/docs/resources/channel#reaction-object
type Reaction struct {
	sync.RWMutex `json:"-"`

	Count uint          `json:"count"`
	Me    bool          `json:"me"`
	Emoji *PartialEmoji `json:"Emoji"`
}

func (r *Reaction) DeepCopy() (copy interface{}) {
	copy = &Reaction{}
	r.CopyOverTo(copy)

	return
}

func (r *Reaction) CopyOverTo(other interface{}) (err error) {
	var reaction *Reaction
	var valid bool
	if reaction, valid = other.(*Reaction); !valid {
		err = NewErrorUnsupportedType("given interface{} is not of type *Reaction")
		return
	}

	r.RLock()
	reaction.Lock()

	reaction.Count = r.Count
	reaction.Me = r.Me

	if r.Emoji != nil {
		reaction.Emoji = r.Emoji.DeepCopy().(*Emoji)
	}

	r.RUnlock()
	reaction.Unlock()
	return
}

// -----------------
// Embed

// limitations: https://discordapp.com/developers/docs/resources/channel#embed-limits
// TODO: implement NewEmbedX functions that ensures limitations

// ChannelEmbed https://discordapp.com/developers/docs/resources/channel#embed-object
type ChannelEmbed struct {
	sync.RWMutex `json:"-"`

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

func (c *ChannelEmbed) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbed{}
	c.CopyOverTo(copy)

	return
}

func (c *ChannelEmbed) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbed
	var valid bool
	if embed, valid = other.(*ChannelEmbed); !valid {
		err = NewErrorUnsupportedType("given interface{} is not of type *ChannelEmbed")
		return
	}

	c.RLock()
	embed.Lock()

	embed.Title = c.Title
	embed.Type = c.Type
	embed.Description = c.Description
	embed.URL = c.URL
	embed.Timestamp = c.Timestamp
	embed.Color = c.Color

	if c.Footer != nil {
		embed.Footer = c.Footer.DeepCopy().(*ChannelEmbedFooter)
	}
	if c.Image != nil {
		embed.Image = c.Image.DeepCopy().(*ChannelEmbedImage)
	}
	if c.Thumbnail != nil {
		embed.Thumbnail = c.Thumbnail.DeepCopy().(*ChannelEmbedThumbnail)
	}
	if c.Video != nil {
		embed.Video = c.Video.DeepCopy().(*ChannelEmbedVideo)
	}
	if c.Provider != nil {
		embed.Provider = c.Provider.DeepCopy().(*ChannelEmbedProvider)
	}
	if c.Author != nil {
		embed.Author = c.Author.DeepCopy().(*ChannelEmbedAuthor)
	}

	embed.Fields = make([]*ChannelEmbedField, len(c.Fields))
	for i, field := range c.Fields {
		embed.Fields[i] = field.DeepCopy().(*ChannelEmbedField)
	}

	c.RUnlock()
	embed.Unlock()
	return
}

// ChannelEmbedThumbnail https://discordapp.com/developers/docs/resources/channel#embed-object-embed-thumbnail-structure
type ChannelEmbedThumbnail struct {
	sync.RWMutex `json:"-"`

	URL      string `json:"url,omitempty"`       // ?| , source url of image (only supports http(s) and attachments)
	ProxyURL string `json:"proxy_url,omitempty"` // ?| , a proxied url of the image
	Height   int    `json:"height,omitempty"`    // ?| , height of image
	Width    int    `json:"width,omitempty"`     // ?| , width of image
}

func (c *ChannelEmbedThumbnail) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbedThumbnail{}
	c.CopyOverTo(copy)

	return
}

func (c *ChannelEmbedThumbnail) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbedThumbnail
	var valid bool
	if embed, valid = other.(*ChannelEmbedThumbnail); !valid {
		err = NewErrorUnsupportedType("given interface{} is not of type *ChannelEmbedThumbnail")
		return
	}

	c.RLock()
	embed.Lock()

	embed.URL = c.URL
	embed.ProxyURL = c.ProxyURL
	embed.Height = c.Height
	embed.Width = c.Width

	c.RUnlock()
	embed.Unlock()
	return
}

// ChannelEmbedVideo https://discordapp.com/developers/docs/resources/channel#embed-object-embed-video-structure
type ChannelEmbedVideo struct {
	sync.RWMutex `json:"-"`

	URL    string `json:"url,omitempty"`    // ?| , source url of video
	Height int    `json:"height,omitempty"` // ?| , height of video
	Width  int    `json:"width,omitempty"`  // ?| , width of video
}

func (c *ChannelEmbedVideo) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbedVideo{}
	c.CopyOverTo(copy)

	return
}

func (c *ChannelEmbedVideo) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbedVideo
	var valid bool
	if embed, valid = other.(*ChannelEmbedVideo); !valid {
		err = NewErrorUnsupportedType("given interface{} is not of type *ChannelEmbedVideo")
		return
	}

	c.RLock()
	embed.Lock()

	embed.URL = c.URL
	embed.Height = c.Height
	embed.Width = c.Width

	c.RUnlock()
	embed.Unlock()
	return
}

// ChannelEmbedImage https://discordapp.com/developers/docs/resources/channel#embed-object-embed-image-structure
type ChannelEmbedImage struct {
	sync.RWMutex `json:"-"`

	URL      string `json:"url,omitempty"`       // ?| , source url of image (only supports http(s) and attachments)
	ProxyURL string `json:"proxy_url,omitempty"` // ?| , a proxied url of the image
	Height   int    `json:"height,omitempty"`    // ?| , height of image
	Width    int    `json:"width,omitempty"`     // ?| , width of image
}

func (c *ChannelEmbedImage) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbedImage{}
	c.CopyOverTo(copy)

	return
}

func (c *ChannelEmbedImage) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbedImage
	var valid bool
	if embed, valid = other.(*ChannelEmbedImage); !valid {
		err = NewErrorUnsupportedType("given interface{} is not of type *ChannelEmbedImage")
		return
	}

	c.RLock()
	embed.Lock()

	embed.URL = c.URL
	embed.ProxyURL = c.ProxyURL
	embed.Height = c.Height
	embed.Width = c.Width

	c.RUnlock()
	embed.Unlock()
	return
}

// ChannelEmbedProvider https://discordapp.com/developers/docs/resources/channel#embed-object-embed-provider-structure
type ChannelEmbedProvider struct {
	sync.RWMutex `json:"-"`

	Name string `json:"name,omitempty"` // ?| , name of provider
	URL  string `json:"url,omitempty"`  // ?| , url of provider
}

func (c *ChannelEmbedProvider) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbedProvider{}
	c.CopyOverTo(copy)

	return
}

func (c *ChannelEmbedProvider) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbedProvider
	var valid bool
	if embed, valid = other.(*ChannelEmbedProvider); !valid {
		err = NewErrorUnsupportedType("given interface{} is not of type *ChannelEmbedProvider")
		return
	}

	c.RLock()
	embed.Lock()

	embed.URL = c.URL
	embed.Name = c.Name

	c.RUnlock()
	embed.Unlock()
	return
}

// ChannelEmbedAuthor https://discordapp.com/developers/docs/resources/channel#embed-object-embed-author-structure
type ChannelEmbedAuthor struct {
	sync.RWMutex `json:"-"`

	Name         string `json:"name,omitempty"`           // ?| , name of author
	URL          string `json:"url,omitempty"`            // ?| , url of author
	IconURL      string `json:"icon_url,omitempty"`       // ?| , url of author icon (only supports http(s) and attachments)
	ProxyIconURL string `json:"proxy_icon_url,omitempty"` // ?| , a proxied url of author icon
}

func (c *ChannelEmbedAuthor) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbedAuthor{}
	c.CopyOverTo(copy)

	return
}

func (c *ChannelEmbedAuthor) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbedAuthor
	var valid bool
	if embed, valid = other.(*ChannelEmbedAuthor); !valid {
		err = NewErrorUnsupportedType("given interface{} is not of type *ChannelEmbedAuthor")
		return
	}

	c.RLock()
	embed.Lock()

	embed.Name = c.Name
	embed.URL = c.URL
	embed.IconURL = c.IconURL
	embed.ProxyIconURL = c.ProxyIconURL

	c.RUnlock()
	embed.Unlock()
	return
}

// ChannelEmbedFooter https://discordapp.com/developers/docs/resources/channel#embed-object-embed-footer-structure
type ChannelEmbedFooter struct {
	sync.RWMutex `json:"-"`

	Text         string `json:"text"`                     //  | , url of author
	IconURL      string `json:"icon_url,omitempty"`       // ?| , url of footer icon (only supports http(s) and attachments)
	ProxyIconURL string `json:"proxy_icon_url,omitempty"` // ?| , a proxied url of footer icon
}

func (c *ChannelEmbedFooter) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbedFooter{}
	c.CopyOverTo(copy)

	return
}

func (c *ChannelEmbedFooter) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbedFooter
	var valid bool
	if embed, valid = other.(*ChannelEmbedFooter); !valid {
		err = NewErrorUnsupportedType("given interface{} is not of type *ChannelEmbedFooter")
		return
	}

	c.RLock()
	embed.Lock()

	embed.Text = c.Text
	embed.IconURL = c.IconURL
	embed.ProxyIconURL = c.ProxyIconURL

	c.RUnlock()
	embed.Unlock()
	return
}

// ChannelEmbedField https://discordapp.com/developers/docs/resources/channel#embed-object-embed-field-structure
type ChannelEmbedField struct {
	sync.RWMutex `json:"-"`

	Name   string `json:"name"`           //  | , name of the field
	Value  string `json:"value"`          //  | , value of the field
	Inline bool   `json:"bool,omitempty"` // ?| , whether or not this field should display inline
}

func (c *ChannelEmbedField) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbedField{}
	c.CopyOverTo(copy)

	return
}

func (c *ChannelEmbedField) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbedField
	var valid bool
	if embed, valid = other.(*ChannelEmbedField); !valid {
		err = NewErrorUnsupportedType("given interface{} is not of type *ChannelEmbedField")
		return
	}

	c.RLock()
	embed.Lock()

	embed.Name = c.Name
	embed.Value = c.Value
	embed.Inline = c.Inline

	c.RUnlock()
	embed.Unlock()
	return
}
