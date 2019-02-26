package disgord

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"

	"github.com/andersfylling/disgord/constant"
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

// PermissionOverwrite https://discordapp.com/developers/docs/resources/channel#overwrite-object
type PermissionOverwrite struct {
	ID    Snowflake `json:"id"`    // role or user id
	Type  string    `json:"type"`  // either `role` or `member`
	Allow int       `json:"allow"` // permission bit set
	Deny  int       `json:"deny"`  // permission bit set
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
	Lockable `json:"-"`
	ID       Snowflake `json:"id"`
	Name     string    `json:"name"`
	Type     uint      `json:"type"`
}

// Channel ...
type Channel struct {
	Lockable             `json:"-"`
	ID                   Snowflake             `json:"id"`
	Type                 uint                  `json:"type"`
	GuildID              Snowflake             `json:"guild_id,omitempty"`              // ?|
	Position             uint                  `json:"position,omitempty"`              // ?|
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"` // ?|
	Name                 string                `json:"name,omitempty"`                  // ?|
	Topic                string                `json:"topic,omitempty"`                 // ?|?
	NSFW                 bool                  `json:"nsfw,omitempty"`                  // ?|
	LastMessageID        Snowflake             `json:"last_message_id,omitempty"`       // ?|?
	Bitrate              uint                  `json:"bitrate,omitempty"`               // ?|
	UserLimit            uint                  `json:"user_limit,omitempty"`            // ?|
	RateLimitPerUser     uint                  `json:"rate_limit_per_user,omitempty"`   // ?|
	Recipients           []*User               `json:"recipient,omitempty"`             // ?| , empty if not DM/GroupDM
	Icon                 *string               `json:"icon,omitempty"`                  // ?|?
	OwnerID              Snowflake             `json:"owner_id,omitempty"`              // ?|
	ApplicationID        Snowflake             `json:"application_id,omitempty"`        // ?|
	ParentID             Snowflake             `json:"parent_id,omitempty"`             // ?|?
	LastPinTimestamp     Timestamp             `json:"last_pin_timestamp,omitempty"`    // ?|

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

func (c *Channel) String() string {
	return "channel{" + c.ID.String() + "}"
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

// Mention creates a channel mention string. Mention format is according the Discord protocol.
func (c *Channel) Mention() string {
	return "<#" + c.ID.String() + ">"
}

// Compare checks if channel A is the same as channel B
func (c *Channel) Compare(other *Channel) bool {
	// eh
	return (c == nil && other == nil) || (other != nil && c.ID == other.ID)
}

func (c *Channel) saveToDiscord(session Session, changes discordSaver) (err error) {
	var updated *Channel

	// verify discord request
	defer func() {
		if err != nil {
			return
		}

		*c = *updated
	}()

	if c.ID.Empty() {
		if c.Type != ChannelTypeDM && c.Type != ChannelTypeGroupDM {
			// create
			if c.Name == "" {
				err = newErrorEmptyValue("must have a channel name before creating channel")
			}
			params := CreateGuildChannelParams{
				Name:                 c.Name,
				PermissionOverwrites: c.PermissionOverwrites,
			}

			// specific
			if c.Type == ChannelTypeGuildText {
				params.NSFW = c.NSFW
				params.Topic = c.Topic
				params.RateLimitPerUser = c.RateLimitPerUser
			} else if c.Type == ChannelTypeGuildVoice {
				params.Bitrate = c.Bitrate
				params.UserLimit = c.UserLimit
			}

			// shared
			if c.Type == ChannelTypeGuildVoice || c.Type == ChannelTypeGuildText {
				params.ParentID = c.ParentID
			}

			updated, err = session.CreateGuildChannel(c.GuildID, &params)
		} else if c.Type == ChannelTypeDM {
			if len(c.Recipients) != 1 {
				err = errors.New("must have only one recipient in Channel.Recipient (with ID) for creating a DM. Got " + strconv.Itoa(len(c.Recipients)))
				return err
			}
			updated, err = session.CreateDM(c.Recipients[0].ID)
		} else if c.Type == ChannelTypeGroupDM {
			err = errors.New("creating group DM using SaveToDiscord has not been implemented")
			//if len(c.Recipients) == 0 {
			//	err = errors.New("must have at least one recipient in Channel.Recipient (with access token) for creating a group DM. Got 0")
			//	return
			//}
			//total := len(c.Recipients)
			//params := CreateGroupDMParams{}
			//params.AccessTokens = make([]string, total)
			//params.Nicks = make(map[Snowflake]string, total)
			//
			//for i := 0; i < total; i++ {
			//	params.AccessTokens[i] = c.Recipients[i].
			//}
			//
			//updated, err = session.CreateGroupDM()
		} else {
			err = errors.New("cannot save to discord. Does not recognise what needs to be saved")
		}
		return err
	}

	return errors.New("updating discord objects are not yet implemented - only saving new ones")
	//
	//
	//// modify / update channel
	//changes := ModifyChannelParams{}
	//
	//// specific
	//if c.Type == ChannelTypeDM {
	//	// nothing to change
	//} else if c.Type == ChannelTypeGroupDM {
	//	// nothing to change
	//} else if c.Type == ChannelTypeGuildText {
	//	changes.SetNSFW(c.NSFW)
	//	changes.SetTopic(c.Topic)
	//	changes.SetRateLimitPerUser(c.RateLimitPerUser)
	//} else if c.Type == ChannelTypeGuildVoice {
	//	changes.SetBitrate(c.Bitrate)
	//	changes.SetUserLimit(c.UserLimit)
	//}
	//
	//// shared
	//if c.Type == ChannelTypeGuildVoice || c.Type == ChannelTypeGuildText {
	//	if c.ParentID.Empty() {
	//		changes.RemoveParentID()
	//	} else {
	//		changes.SetParentID(c.ParentID)
	//	}
	//}
	//
	//// for all
	//changes.SetName(c.Name)
	//changes.SetPosition(c.Position)
	//changes.SetPermissionOverwrites(c.PermissionOverwrites)
	//
	//updated, err = session.ModifyChannel(c.ID, &changes)
	//return
}

func (c *Channel) deleteFromDiscord(session Session) (err error) {
	if c.ID.Empty() {
		err = newErrorMissingSnowflake("channel id/snowflake is empty or missing")
		return
	}
	var deleted *Channel
	if deleted, err = session.DeleteChannel(c.ID); err != nil {
		return
	}

	*c = *deleted
	return
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *Channel) DeepCopy() (copy interface{}) {
	copy = NewChannel()
	c.CopyOverTo(copy)

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

	if constant.LockedMethods {
		c.RWMutex.RLock()
		channel.RWMutex.Lock()
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
	for _, recipient := range c.Recipients {
		channel.Recipients = append(channel.Recipients, recipient.DeepCopy().(*User))
	}

	if constant.LockedMethods {
		c.RWMutex.RUnlock()
		channel.RWMutex.Unlock()
	}

	return
}

func (c *Channel) copyOverToCache(other interface{}) (err error) {
	channel := other.(*Channel)

	if constant.LockedMethods {
		channel.Lock()
		c.RLock()
	}

	channel.ID = c.ID
	channel.Type = c.Type

	if c.Type == ChannelTypeGroupDM || c.Type == ChannelTypeDM {
		if c.Type == ChannelTypeGroupDM {
			channel.Icon = c.Icon
			channel.OwnerID = c.OwnerID
			channel.Name = c.Name
			channel.LastPinTimestamp = c.LastPinTimestamp
		}
		channel.LastMessageID = c.LastMessageID

		if len(c.recipientsIDs) == len(c.Recipients) {
			channel.recipientsIDs = c.recipientsIDs
		} else {
			channel.recipientsIDs = make([]Snowflake, len(c.Recipients))
			for i := range c.Recipients {
				channel.recipientsIDs[i] = c.Recipients[i].ID
			}
		}
	} else if c.Type == ChannelTypeGuildText {
		channel.NSFW = c.NSFW
		channel.Name = c.Name
		channel.Position = c.Position
		channel.PermissionOverwrites = c.PermissionOverwrites
		channel.Topic = c.Topic
		channel.LastMessageID = c.LastMessageID
		channel.RateLimitPerUser = c.RateLimitPerUser
		channel.LastPinTimestamp = c.LastPinTimestamp
		channel.ParentID = c.ParentID
		channel.GuildID = c.GuildID
	} else if c.Type == ChannelTypeGuildVoice {
		channel.Name = c.Name
		channel.Position = c.Position
		channel.PermissionOverwrites = c.PermissionOverwrites
		channel.ParentID = c.ParentID
		channel.Bitrate = c.Bitrate
		channel.UserLimit = c.UserLimit
		channel.GuildID = c.GuildID
	}

	// TODO: evaluate
	channel.ApplicationID = c.ApplicationID

	if constant.LockedMethods {
		channel.Unlock()
		c.RUnlock()
	}
	return nil
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

// SendMsgString same as SendMsg, however this only takes the message content (string) as a argument for the message
func (c *Channel) SendMsgString(client MessageSender, content string) (msg *Message, err error) {
	if c.ID.Empty() {
		err = newErrorMissingSnowflake("snowflake ID not set for channel")
		return
	}
	params := &CreateMessageParams{
		Content: content,
	}

	msg, err = client.CreateMessage(c.ID, params)
	return
}

// SendMsg sends a message to a channel
func (c *Channel) SendMsg(client MessageSender, message *Message) (msg *Message, err error) {
	if c.ID.Empty() {
		err = newErrorMissingSnowflake("snowflake ID not set for channel")
		return
	}
	message.RLock()
	params := &CreateMessageParams{
		Content: message.Content,
		Nonce:   message.Nonce,
		Tts:     message.Tts,
		// File: ...
		// Embed: ...
	}
	if len(message.Embeds) > 0 {
		params.Embed = message.Embeds[0]
	}
	message.RUnlock()

	msg, err = client.CreateMessage(c.ID, params)
	return
}

// -----------------
// Embed

// limitations: https://discordapp.com/developers/docs/resources/channel#embed-limits
// TODO: implement NewEmbedX functions that ensures limitations

// ChannelEmbed https://discordapp.com/developers/docs/resources/channel#embed-object
type ChannelEmbed struct {
	Lockable `json:"-"`

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

// DeepCopy see interface at struct.go#DeepCopier
func (c *ChannelEmbed) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbed{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *ChannelEmbed) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbed
	var valid bool
	if embed, valid = other.(*ChannelEmbed); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *ChannelEmbed")
		return
	}

	if constant.LockedMethods {
		c.RLock()
		embed.Lock()
	}

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

	if constant.LockedMethods {
		c.RUnlock()
		embed.Unlock()
	}
	return nil
}

// ChannelEmbedThumbnail https://discordapp.com/developers/docs/resources/channel#embed-object-embed-thumbnail-structure
type ChannelEmbedThumbnail struct {
	Lockable `json:"-"`

	URL      string `json:"url,omitempty"`       // ?| , source url of image (only supports http(s) and attachments)
	ProxyURL string `json:"proxy_url,omitempty"` // ?| , a proxied url of the image
	Height   int    `json:"height,omitempty"`    // ?| , height of image
	Width    int    `json:"width,omitempty"`     // ?| , width of image
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *ChannelEmbedThumbnail) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbedThumbnail{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *ChannelEmbedThumbnail) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbedThumbnail
	var valid bool
	if embed, valid = other.(*ChannelEmbedThumbnail); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *ChannelEmbedThumbnail")
		return
	}

	if constant.LockedMethods {
		c.RLock()
		embed.Lock()
	}

	embed.URL = c.URL
	embed.ProxyURL = c.ProxyURL
	embed.Height = c.Height
	embed.Width = c.Width

	if constant.LockedMethods {
		c.RUnlock()
		embed.Unlock()
	}
	return
}

// ChannelEmbedVideo https://discordapp.com/developers/docs/resources/channel#embed-object-embed-video-structure
type ChannelEmbedVideo struct {
	Lockable `json:"-"`

	URL    string `json:"url,omitempty"`    // ?| , source url of video
	Height int    `json:"height,omitempty"` // ?| , height of video
	Width  int    `json:"width,omitempty"`  // ?| , width of video
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *ChannelEmbedVideo) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbedVideo{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *ChannelEmbedVideo) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbedVideo
	var valid bool
	if embed, valid = other.(*ChannelEmbedVideo); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *ChannelEmbedVideo")
		return
	}

	if constant.LockedMethods {
		c.RLock()
		embed.Lock()
	}

	embed.URL = c.URL
	embed.Height = c.Height
	embed.Width = c.Width

	if constant.LockedMethods {
		c.RUnlock()
		embed.Unlock()
	}
	return nil
}

// ChannelEmbedImage https://discordapp.com/developers/docs/resources/channel#embed-object-embed-image-structure
type ChannelEmbedImage struct {
	Lockable `json:"-"`

	URL      string `json:"url,omitempty"`       // ?| , source url of image (only supports http(s) and attachments)
	ProxyURL string `json:"proxy_url,omitempty"` // ?| , a proxied url of the image
	Height   int    `json:"height,omitempty"`    // ?| , height of image
	Width    int    `json:"width,omitempty"`     // ?| , width of image
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *ChannelEmbedImage) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbedImage{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *ChannelEmbedImage) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbedImage
	var valid bool
	if embed, valid = other.(*ChannelEmbedImage); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *ChannelEmbedImage")
		return
	}

	if constant.LockedMethods {
		c.RLock()
		embed.Lock()
	}

	embed.URL = c.URL
	embed.ProxyURL = c.ProxyURL
	embed.Height = c.Height
	embed.Width = c.Width

	if constant.LockedMethods {
		c.RUnlock()
		embed.Unlock()
	}
	return nil
}

// ChannelEmbedProvider https://discordapp.com/developers/docs/resources/channel#embed-object-embed-provider-structure
type ChannelEmbedProvider struct {
	Lockable `json:"-"`

	Name string `json:"name,omitempty"` // ?| , name of provider
	URL  string `json:"url,omitempty"`  // ?| , url of provider
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *ChannelEmbedProvider) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbedProvider{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *ChannelEmbedProvider) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbedProvider
	var valid bool
	if embed, valid = other.(*ChannelEmbedProvider); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *ChannelEmbedProvider")
		return
	}

	if constant.LockedMethods {
		c.RLock()
		embed.Lock()
	}

	embed.URL = c.URL
	embed.Name = c.Name

	if constant.LockedMethods {
		c.RUnlock()
		embed.Unlock()
	}
	return nil
}

// ChannelEmbedAuthor https://discordapp.com/developers/docs/resources/channel#embed-object-embed-author-structure
type ChannelEmbedAuthor struct {
	Lockable `json:"-"`

	Name         string `json:"name,omitempty"`           // ?| , name of author
	URL          string `json:"url,omitempty"`            // ?| , url of author
	IconURL      string `json:"icon_url,omitempty"`       // ?| , url of author icon (only supports http(s) and attachments)
	ProxyIconURL string `json:"proxy_icon_url,omitempty"` // ?| , a proxied url of author icon
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *ChannelEmbedAuthor) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbedAuthor{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *ChannelEmbedAuthor) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbedAuthor
	var valid bool
	if embed, valid = other.(*ChannelEmbedAuthor); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *ChannelEmbedAuthor")
		return
	}

	if constant.LockedMethods {
		c.RLock()
		embed.Lock()
	}

	embed.Name = c.Name
	embed.URL = c.URL
	embed.IconURL = c.IconURL
	embed.ProxyIconURL = c.ProxyIconURL

	if constant.LockedMethods {
		c.RUnlock()
		embed.Unlock()
	}
	return nil
}

// ChannelEmbedFooter https://discordapp.com/developers/docs/resources/channel#embed-object-embed-footer-structure
type ChannelEmbedFooter struct {
	Lockable `json:"-"`

	Text         string `json:"text"`                     //  | , url of author
	IconURL      string `json:"icon_url,omitempty"`       // ?| , url of footer icon (only supports http(s) and attachments)
	ProxyIconURL string `json:"proxy_icon_url,omitempty"` // ?| , a proxied url of footer icon
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *ChannelEmbedFooter) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbedFooter{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *ChannelEmbedFooter) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbedFooter
	var valid bool
	if embed, valid = other.(*ChannelEmbedFooter); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *ChannelEmbedFooter")
		return
	}

	if constant.LockedMethods {
		c.RLock()
		embed.Lock()
	}

	embed.Text = c.Text
	embed.IconURL = c.IconURL
	embed.ProxyIconURL = c.ProxyIconURL

	if constant.LockedMethods {
		c.RUnlock()
		embed.Unlock()
	}
	return nil
}

// ChannelEmbedField https://discordapp.com/developers/docs/resources/channel#embed-object-embed-field-structure
type ChannelEmbedField struct {
	Lockable `json:"-"`

	Name   string `json:"name"`           //  | , name of the field
	Value  string `json:"value"`          //  | , value of the field
	Inline bool   `json:"bool,omitempty"` // ?| , whether or not this field should display inline
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *ChannelEmbedField) DeepCopy() (copy interface{}) {
	copy = &ChannelEmbedField{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *ChannelEmbedField) CopyOverTo(other interface{}) (err error) {
	var embed *ChannelEmbedField
	var valid bool
	if embed, valid = other.(*ChannelEmbedField); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *ChannelEmbedField")
		return
	}

	if constant.LockedMethods {
		c.RLock()
		embed.Lock()
	}

	embed.Name = c.Name
	embed.Value = c.Value
	embed.Inline = c.Inline

	if constant.LockedMethods {
		c.RUnlock()
		embed.Unlock()
	}
	return nil
}

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

func ratelimitChannel(id Snowflake) string {
	return "c:" + id.String()
}
func ratelimitChannelPermissions(id Snowflake) string {
	return ratelimitChannel(id) + ":perm"
}
func ratelimitChannelInvites(id Snowflake) string {
	return ratelimitChannel(id) + ":i"
}
func ratelimitChannelTyping(id Snowflake) string {
	return ratelimitChannel(id) + ":t"
}
func ratelimitChannelPins(id Snowflake) string {
	return ratelimitChannel(id) + ":pins"
}
func ratelimitChannelRecipients(id Snowflake) string {
	return ratelimitChannel(id) + ":r"
}
func ratelimitChannelMessages(id Snowflake) string {
	return ratelimitChannel(id) + ":m"
}
func ratelimitChannelMessagesDelete(id Snowflake) string {
	return ratelimitChannelMessages(id) + "_"
}
func ratelimitChannelWebhooks(id Snowflake) string {
	return ratelimitChannel(id) + ":w"
}

// GetChannel [REST] Get a channel by Snowflake. Returns a channel object.
//  Method                  GET
//  Endpoint                /channels/{channel.id}
//  Rate limiter [MAJOR]    /channels/{channel.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/channel#get-channel
//  Reviewed                2018-06-07
//  Comment                 -
func (c *client) GetChannel(id Snowflake, flags ...Flag) (ret *Channel, err error) {
	if id.Empty() {
		return nil, errors.New("not a valid snowflake")
	}

	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitChannel(id),
		Endpoint:    endpoint.Channel(id),
	}, flags)
	r.CacheRegistry = ChannelCache
	r.ID = id
	r.pool = c.pool.channel

	return getChannel(r.Execute)
}

// UpdateChannel [REST] Update a channels settings. Requires the 'MANAGE_CHANNELS' permission for the guild. Returns
// a channel on success, and a 400 BAD REQUEST on invalid parameters. Fires a Channel Update Gateway event. If
// modifying a category, individual Channel Update events will fire for each child channel that also changes.
// For the PATCH method, all the JSON Params are optional.
//  Method                  PUT/PATCH
//  Endpoint                /channels/{channel.id}
//  Rate limiter [MAJOR]    /channels/{channel.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/channel#modify-channel
//  Reviewed                2018-06-07
//  Comment                 andersfylling: only implemented the patch method, as its parameters are optional.
func (c *client) UpdateChannel(id Snowflake, flags ...Flag) (builder *updateChannelBuilder) {
	builder = &updateChannelBuilder{}
	builder.r.itemFactory = func() interface{} {
		return c.pool.channel.Get()
	}
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      http.MethodPatch,
		Ratelimiter: ratelimitChannel(id),
		Endpoint:    endpoint.Channel(id),
		ContentType: httd.ContentTypeJSON,
	}, nil)
	builder.r.cacheRegistry = ChannelCache
	builder.r.cacheItemID = id

	return builder
}

// updateChannelBuilder https://discordapp.com/developers/docs/resources/channel#modify-channel-json-params
//generate-rest-params: parent_id:Snowflake, permission_overwrites:[]PermissionOverwrite, user_limit:int, bitrate:uint, rate_limit_per_user:uint, nsfw:bool, topic:string, position:uint, name:string,
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

// DeleteChannel [REST] Delete a channel, or close a private message. Requires the 'MANAGE_CHANNELS' permission for
// the guild. Deleting a category does not delete its child channels; they will have their parent_id removed and a
// Channel Update Gateway event will fire for each of them. Returns a channel object on success.
// Fires a Channel Delete Gateway event.
//  Method                  Delete
//  Endpoint                /channels/{channel.id}
//  Rate limiter [MAJOR]    /channels/{channel.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/channel#deleteclose-channel
//  Reviewed                2018-10-09
//  Comment                 Deleting a guild channel cannot be undone. Use this with caution, as it
//                          is impossible to undo this action when performed on a guild channel. In
//                          contrast, when used with a private message, it is possible to undo the
//                          action by opening a private message with the recipient again.
func (c *client) DeleteChannel(id Snowflake, flags ...Flag) (channel *Channel, err error) {
	if id.Empty() {
		err = errors.New("not a valid snowflake")
		return
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodDelete,
		Ratelimiter: ratelimitChannel(id),
		Endpoint:    endpoint.Channel(id),
	}, flags)
	r.expectsStatusCode = http.StatusOK
	r.updateCache = func(registry cacheRegistry, id Snowflake, x interface{}) (err error) {
		c.cache.DeleteChannel(id)
		return nil
	}

	return getChannel(r.Execute)
}

// UpdateChannelPermissionsParams https://discordapp.com/developers/docs/resources/channel#edit-channel-permissions-json-params
type UpdateChannelPermissionsParams struct {
	Allow int    `json:"allow"` // the bitwise value of all allowed permissions
	Deny  int    `json:"deny"`  // the bitwise value of all disallowed permissions
	Type  string `json:"type"`  // "member" for a user or "role" for a role
}

// EditChannelPermissions [REST] Edit the channel permission overwrites for a user or role in a channel. Only usable
// for guild channels. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty response on success.
// For more information about permissions, see permissions.
//  Method                  PUT
//  Endpoint                /channels/{channel.id}/permissions/{overwrite.id}
//  Rate limiter [MAJOR]    /channels/{channel.id}/permissions
//  Discord documentation   https://discordapp.com/developers/docs/resources/channel#edit-channel-permissions
//  Reviewed                2018-06-07
//  Comment                 -
func (c *client) UpdateChannelPermissions(chanID, overwriteID Snowflake, params *UpdateChannelPermissionsParams, flags ...Flag) (err error) {
	if chanID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if overwriteID.Empty() {
		return errors.New("overwriteID must be set to target the specific channel permissions")
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPut,
		Ratelimiter: ratelimitChannelPermissions(chanID),
		Endpoint:    endpoint.ChannelPermission(chanID, overwriteID),
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
//  Rate limiter [MAJOR]    /channels/{channel.id}/invites
//  Discord documentation   https://discordapp.com/developers/docs/resources/channel#get-channel-invites
//  Reviewed                2018-06-07
//  Comment                 -
func (c *client) GetChannelInvites(id Snowflake, flags ...Flag) (invites []*Invite, err error) {
	if id.Empty() {
		err = errors.New("channelID must be set to target the correct channel")
		return
	}

	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitChannelInvites(id),
		Endpoint:    endpoint.ChannelInvites(id),
	}, flags)
	r.CacheRegistry = ChannelCache
	r.factory = func() interface{} {
		tmp := make([]*Invite, 0)
		return &tmp
	}

	var vs interface{}
	if vs, err = r.Execute(); err != nil {
		return nil, err
	}

	if tmp, ok := vs.(*[]*Invite); ok {
		return *tmp, nil
	}
	return nil, errors.New("unable to cast Invite slice")
}

// CreateChannelInvitesParams https://discordapp.com/developers/docs/resources/channel#create-channel-invite-json-params
type CreateChannelInvitesParams struct {
	MaxAge    int  `json:"max_age,omitempty"`   // duration of invite in seconds before expiry, or 0 for never. default 86400 (24 hours)
	MaxUses   int  `json:"max_uses,omitempty"`  // max number of uses or 0 for unlimited. default 0
	Temporary bool `json:"temporary,omitempty"` // whether this invite only grants temporary membership. default false
	Unique    bool `json:"unique,omitempty"`    // if true, don't try to reuse a similar invite (useful for creating many unique one time use invites). default false
}

// CreateChannelInvites [REST] Create a new invite object for the channel. Only usable for guild channels. Requires
// the CREATE_INSTANT_INVITE permission. All JSON parameters for this route are optional, however the request body is
// not. If you are not sending any fields, you still have to send an empty JSON object ({}). Returns an invite object.
//  Method                  POST
//  Endpoint                /channels/{channel.id}/invites
//  Rate limiter [MAJOR]    /channels/{channel.id}/invites
//  Discord documentation   https://discordapp.com/developers/docs/resources/channel#create-channel-invite
//  Reviewed                2018-06-07
//  Comment                 -
func (c *client) CreateChannelInvites(id Snowflake, params *CreateChannelInvitesParams, flags ...Flag) (ret *Invite, err error) {
	if id.Empty() {
		err = errors.New("channelID must be set to target the correct channel")
		return nil, err
	}
	if params == nil {
		params = &CreateChannelInvitesParams{} // have to send an empty JSON object ({}). maybe just struct{}?
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPost,
		Ratelimiter: ratelimitChannelInvites(id),
		Endpoint:    endpoint.ChannelInvites(id),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.factory = func() interface{} {
		return &Invite{}
	}

	return getInvite(r.Execute)
}

// DeleteChannelPermission [REST] Delete a channel permission overwrite for a user or role in a channel. Only usable
// for guild channels. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty response on success. For more
// information about permissions, see permissions: https://discordapp.com/developers/docs/topics/permissions#permissions
//  Method                  DELETE
//  Endpoint                /channels/{channel.id}/permissions/{overwrite.id}
//  Rate limiter [MAJOR]    /channels/{channel.id}/permissions
//  Discord documentation   https://discordapp.com/developers/docs/resources/channel#delete-channel-permission
//  Reviewed                2018-06-07
//  Comment                 -
func (c *client) DeleteChannelPermission(channelID, overwriteID Snowflake, flags ...Flag) (err error) {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if overwriteID.Empty() {
		return errors.New("overwriteID must be set to target the specific channel permissions")
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodDelete,
		Ratelimiter: ratelimitChannelPermissions(channelID),
		Endpoint:    endpoint.ChannelPermission(channelID, overwriteID),
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
	if g.UserID.Empty() {
		return errors.New("missing userID")
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
//  Rate limiter [MAJOR]    /channels/{channel.id}/recipients
//  Discord documentation   https://discordapp.com/developers/docs/resources/channel#group-dm-add-recipient
//  Reviewed                2018-06-10
//  Comment                 -
func (c *client) AddDMParticipant(channelID Snowflake, participant *GroupDMParticipant, flags ...Flag) error {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if participant == nil {
		return errors.New("params can not be nil")
	}
	if err := participant.FindErrors(); err != nil {
		return err
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPut,
		Ratelimiter: ratelimitChannelRecipients(channelID),
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
//  Rate limiter [MAJOR]    /channels/{channel.id}/recipients
//  Discord documentation   https://discordapp.com/developers/docs/resources/channel#group-dm-remove-recipient
//  Reviewed                2018-06-10
//  Comment                 -
func (c *client) KickParticipant(channelID, userID Snowflake, flags ...Flag) (err error) {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if userID.Empty() {
		return errors.New("userID must be set to target the specific recipient")
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodDelete,
		Ratelimiter: ratelimitChannelRecipients(channelID),
		Endpoint:    endpoint.ChannelRecipient(channelID, userID),
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}
