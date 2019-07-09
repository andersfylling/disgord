package disgord

import (
	"github.com/andersfylling/disgord/constant"
)

// limitations: https://discordapp.com/developers/docs/resources/channel#embed-limits
// TODO: implement NewEmbedX functions that ensures limitations

// Embed https://discordapp.com/developers/docs/resources/channel#embed-object
type Embed struct {
	Lockable `json:"-"`

	Title       string          `json:"title,omitempty"`       // title of embed
	Type        string          `json:"type,omitempty"`        // type of embed (always "rich" for webhook embeds)
	Description string          `json:"description,omitempty"` // description of embed
	URL         string          `json:"url,omitempty"`         // url of embed
	Timestamp   Time            `json:"timestamp,omitempty"`   // timestamp	timestamp of embed content
	Color       int             `json:"color"`                 // color code of the embed
	Footer      *EmbedFooter    `json:"footer,omitempty"`      // embed footer object	footer information
	Image       *EmbedImage     `json:"image,omitempty"`       // embed image object	image information
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`   // embed thumbnail object	thumbnail information
	Video       *EmbedVideo     `json:"video,omitempty"`       // embed video object	video information
	Provider    *EmbedProvider  `json:"provider,omitempty"`    // embed provider object	provider information
	Author      *EmbedAuthor    `json:"author,omitempty"`      // embed author object	author information
	Fields      []*EmbedField   `json:"fields,omitempty"`      //	array of embed field objects	fields information
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *Embed) DeepCopy() (copy interface{}) {
	copy = &Embed{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *Embed) CopyOverTo(other interface{}) (err error) {
	var embed *Embed
	var valid bool
	if embed, valid = other.(*Embed); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *Embed")
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
		embed.Footer = c.Footer.DeepCopy().(*EmbedFooter)
	}
	if c.Image != nil {
		embed.Image = c.Image.DeepCopy().(*EmbedImage)
	}
	if c.Thumbnail != nil {
		embed.Thumbnail = c.Thumbnail.DeepCopy().(*EmbedThumbnail)
	}
	if c.Video != nil {
		embed.Video = c.Video.DeepCopy().(*EmbedVideo)
	}
	if c.Provider != nil {
		embed.Provider = c.Provider.DeepCopy().(*EmbedProvider)
	}
	if c.Author != nil {
		embed.Author = c.Author.DeepCopy().(*EmbedAuthor)
	}

	embed.Fields = make([]*EmbedField, len(c.Fields))
	for i, field := range c.Fields {
		embed.Fields[i] = field.DeepCopy().(*EmbedField)
	}

	if constant.LockedMethods {
		c.RUnlock()
		embed.Unlock()
	}
	return nil
}

// EmbedThumbnail https://discordapp.com/developers/docs/resources/channel#embed-object-embed-thumbnail-structure
type EmbedThumbnail struct {
	Lockable `json:"-"`

	URL      string `json:"url,omitempty"`       // ?| , source url of image (only supports http(s) and attachments)
	ProxyURL string `json:"proxy_url,omitempty"` // ?| , a proxied url of the image
	Height   int    `json:"height,omitempty"`    // ?| , height of image
	Width    int    `json:"width,omitempty"`     // ?| , width of image
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *EmbedThumbnail) DeepCopy() (copy interface{}) {
	copy = &EmbedThumbnail{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *EmbedThumbnail) CopyOverTo(other interface{}) (err error) {
	var embed *EmbedThumbnail
	var valid bool
	if embed, valid = other.(*EmbedThumbnail); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *EmbedThumbnail")
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

// EmbedVideo https://discordapp.com/developers/docs/resources/channel#embed-object-embed-video-structure
type EmbedVideo struct {
	Lockable `json:"-"`

	URL    string `json:"url,omitempty"`    // ?| , source url of video
	Height int    `json:"height,omitempty"` // ?| , height of video
	Width  int    `json:"width,omitempty"`  // ?| , width of video
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *EmbedVideo) DeepCopy() (copy interface{}) {
	copy = &EmbedVideo{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *EmbedVideo) CopyOverTo(other interface{}) (err error) {
	var embed *EmbedVideo
	var valid bool
	if embed, valid = other.(*EmbedVideo); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *EmbedVideo")
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

// EmbedImage https://discordapp.com/developers/docs/resources/channel#embed-object-embed-image-structure
type EmbedImage struct {
	Lockable `json:"-"`

	URL      string `json:"url,omitempty"`       // ?| , source url of image (only supports http(s) and attachments)
	ProxyURL string `json:"proxy_url,omitempty"` // ?| , a proxied url of the image
	Height   int    `json:"height,omitempty"`    // ?| , height of image
	Width    int    `json:"width,omitempty"`     // ?| , width of image
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *EmbedImage) DeepCopy() (copy interface{}) {
	copy = &EmbedImage{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *EmbedImage) CopyOverTo(other interface{}) (err error) {
	var embed *EmbedImage
	var valid bool
	if embed, valid = other.(*EmbedImage); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *EmbedImage")
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

// EmbedProvider https://discordapp.com/developers/docs/resources/channel#embed-object-embed-provider-structure
type EmbedProvider struct {
	Lockable `json:"-"`

	Name string `json:"name,omitempty"` // ?| , name of provider
	URL  string `json:"url,omitempty"`  // ?| , url of provider
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *EmbedProvider) DeepCopy() (copy interface{}) {
	copy = &EmbedProvider{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *EmbedProvider) CopyOverTo(other interface{}) (err error) {
	var embed *EmbedProvider
	var valid bool
	if embed, valid = other.(*EmbedProvider); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *EmbedProvider")
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

// EmbedAuthor https://discordapp.com/developers/docs/resources/channel#embed-object-embed-author-structure
type EmbedAuthor struct {
	Lockable `json:"-"`

	Name         string `json:"name,omitempty"`           // ?| , name of author
	URL          string `json:"url,omitempty"`            // ?| , url of author
	IconURL      string `json:"icon_url,omitempty"`       // ?| , url of author icon (only supports http(s) and attachments)
	ProxyIconURL string `json:"proxy_icon_url,omitempty"` // ?| , a proxied url of author icon
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *EmbedAuthor) DeepCopy() (copy interface{}) {
	copy = &EmbedAuthor{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *EmbedAuthor) CopyOverTo(other interface{}) (err error) {
	var embed *EmbedAuthor
	var valid bool
	if embed, valid = other.(*EmbedAuthor); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *EmbedAuthor")
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

// EmbedFooter https://discordapp.com/developers/docs/resources/channel#embed-object-embed-footer-structure
type EmbedFooter struct {
	Lockable `json:"-"`

	Text         string `json:"text"`                     //  | , url of author
	IconURL      string `json:"icon_url,omitempty"`       // ?| , url of footer icon (only supports http(s) and attachments)
	ProxyIconURL string `json:"proxy_icon_url,omitempty"` // ?| , a proxied url of footer icon
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *EmbedFooter) DeepCopy() (copy interface{}) {
	copy = &EmbedFooter{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *EmbedFooter) CopyOverTo(other interface{}) (err error) {
	var embed *EmbedFooter
	var valid bool
	if embed, valid = other.(*EmbedFooter); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *EmbedFooter")
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

// EmbedField https://discordapp.com/developers/docs/resources/channel#embed-object-embed-field-structure
type EmbedField struct {
	Lockable `json:"-"`

	Name   string `json:"name"`           //  | , name of the field
	Value  string `json:"value"`          //  | , value of the field
	Inline bool   `json:"inline,omitempty"` // ?| , whether or not this field should display inline
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *EmbedField) DeepCopy() (copy interface{}) {
	copy = &EmbedField{}
	c.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (c *EmbedField) CopyOverTo(other interface{}) (err error) {
	var embed *EmbedField
	var valid bool
	if embed, valid = other.(*EmbedField); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *EmbedField")
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
