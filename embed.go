package disgord

// limitations: https://discord.com/developers/docs/resources/channel#embed-limits
// TODO: implement NewEmbedX functions that ensures limitations

// Embed https://discord.com/developers/docs/resources/channel#embed-object
type Embed struct {
	Title       string          `json:"title,omitempty"`       // title of embed
	Type        string          `json:"type,omitempty"`        // type of embed (always "rich" for webhook embeds)
	Description string          `json:"description,omitempty"` // description of embed
	URL         string          `json:"url,omitempty"`         // url of embed
	Timestamp   Time            `json:"timestamp,omitempty"`   // timestamp	timestamp of embed content
	Color       int             `json:"color,omitempty"`       // color code of the embed
	Footer      *EmbedFooter    `json:"footer,omitempty"`      // embed footer object	footer information
	Image       *EmbedImage     `json:"image,omitempty"`       // embed image object	image information
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`   // embed thumbnail object	thumbnail information
	Video       *EmbedVideo     `json:"video,omitempty"`       // embed video object	video information
	Provider    *EmbedProvider  `json:"provider,omitempty"`    // embed provider object	provider information
	Author      *EmbedAuthor    `json:"author,omitempty"`      // embed author object	author information
	Fields      []*EmbedField   `json:"fields,omitempty"`      //	array of embed field objects	fields information
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *Embed) deepCopy() interface{} {
	cp := &Embed{}
	_ = DeepCopyOver(cp, c)
	return cp
}

// CopyOverTo see interface at struct.go#Copier
func (c *Embed) copyOverTo(other interface{}) (err error) {
	var embed *Embed
	var valid bool
	if embed, valid = other.(*Embed); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *Embed")
		return
	}

	embed.Title = c.Title
	embed.Type = c.Type
	embed.Description = c.Description
	embed.URL = c.URL
	embed.Timestamp = c.Timestamp
	embed.Color = c.Color

	if c.Footer != nil {
		embed.Footer = DeepCopy(c.Footer).(*EmbedFooter)
	}
	if c.Image != nil {
		embed.Image = DeepCopy(c.Image).(*EmbedImage)
	}
	if c.Thumbnail != nil {
		embed.Thumbnail = DeepCopy(c.Thumbnail).(*EmbedThumbnail)
	}
	if c.Video != nil {
		embed.Video = DeepCopy(c.Video).(*EmbedVideo)
	}
	if c.Provider != nil {
		embed.Provider = DeepCopy(c.Provider).(*EmbedProvider)
	}
	if c.Author != nil {
		embed.Author = DeepCopy(c.Author).(*EmbedAuthor)
	}

	embed.Fields = make([]*EmbedField, len(c.Fields))
	for i, field := range c.Fields {
		embed.Fields[i] = DeepCopy(field).(*EmbedField)
	}
	return nil
}

// EmbedThumbnail https://discord.com/developers/docs/resources/channel#embed-object-embed-thumbnail-structure
type EmbedThumbnail struct {
	URL      string `json:"url,omitempty"`       // ?| , source url of image (only supports http(s) and attachments)
	ProxyURL string `json:"proxy_url,omitempty"` // ?| , a proxied url of the image
	Height   int    `json:"height,omitempty"`    // ?| , height of image
	Width    int    `json:"width,omitempty"`     // ?| , width of image
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *EmbedThumbnail) deepCopy() interface{} {
	cp := &EmbedThumbnail{}
	_ = DeepCopyOver(cp, c)
	return cp
}

// CopyOverTo see interface at struct.go#Copier
func (c *EmbedThumbnail) copyOverTo(other interface{}) (err error) {
	var embed *EmbedThumbnail
	var valid bool
	if embed, valid = other.(*EmbedThumbnail); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *EmbedThumbnail")
		return
	}

	embed.URL = c.URL
	embed.ProxyURL = c.ProxyURL
	embed.Height = c.Height
	embed.Width = c.Width
	return
}

// EmbedVideo https://discord.com/developers/docs/resources/channel#embed-object-embed-video-structure
type EmbedVideo struct {
	URL    string `json:"url,omitempty"`    // ?| , source url of video
	Height int    `json:"height,omitempty"` // ?| , height of video
	Width  int    `json:"width,omitempty"`  // ?| , width of video
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *EmbedVideo) deepCopy() interface{} {
	cp := &EmbedVideo{}
	_ = DeepCopyOver(cp, c)
	return cp
}

// CopyOverTo see interface at struct.go#Copier
func (c *EmbedVideo) copyOverTo(other interface{}) (err error) {
	var embed *EmbedVideo
	var valid bool
	if embed, valid = other.(*EmbedVideo); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *EmbedVideo")
		return
	}

	embed.URL = c.URL
	embed.Height = c.Height
	embed.Width = c.Width
	return nil
}

// EmbedImage https://discord.com/developers/docs/resources/channel#embed-object-embed-image-structure
type EmbedImage struct {
	URL      string `json:"url,omitempty"`       // ?| , source url of image (only supports http(s) and attachments)
	ProxyURL string `json:"proxy_url,omitempty"` // ?| , a proxied url of the image
	Height   int    `json:"height,omitempty"`    // ?| , height of image
	Width    int    `json:"width,omitempty"`     // ?| , width of image
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *EmbedImage) deepCopy() interface{} {
	cp := &EmbedImage{}
	_ = DeepCopyOver(cp, c)
	return cp
}

// CopyOverTo see interface at struct.go#Copier
func (c *EmbedImage) copyOverTo(other interface{}) (err error) {
	var embed *EmbedImage
	var valid bool
	if embed, valid = other.(*EmbedImage); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *EmbedImage")
		return
	}

	embed.URL = c.URL
	embed.ProxyURL = c.ProxyURL
	embed.Height = c.Height
	embed.Width = c.Width
	return nil
}

// EmbedProvider https://discord.com/developers/docs/resources/channel#embed-object-embed-provider-structure
type EmbedProvider struct {
	Name string `json:"name,omitempty"` // ?| , name of provider
	URL  string `json:"url,omitempty"`  // ?| , url of provider
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *EmbedProvider) deepCopy() interface{} {
	cp := &EmbedProvider{}
	_ = DeepCopyOver(cp, c)
	return cp
}

// CopyOverTo see interface at struct.go#Copier
func (c *EmbedProvider) copyOverTo(other interface{}) (err error) {
	var embed *EmbedProvider
	var valid bool
	if embed, valid = other.(*EmbedProvider); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *EmbedProvider")
		return
	}

	embed.URL = c.URL
	embed.Name = c.Name
	return nil
}

// EmbedAuthor https://discord.com/developers/docs/resources/channel#embed-object-embed-author-structure
type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`           // ?| , name of author
	URL          string `json:"url,omitempty"`            // ?| , url of author
	IconURL      string `json:"icon_url,omitempty"`       // ?| , url of author icon (only supports http(s) and attachments)
	ProxyIconURL string `json:"proxy_icon_url,omitempty"` // ?| , a proxied url of author icon
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *EmbedAuthor) deepCopy() interface{} {
	cp := &EmbedAuthor{}
	_ = DeepCopyOver(cp, c)
	return cp
}

// CopyOverTo see interface at struct.go#Copier
func (c *EmbedAuthor) copyOverTo(other interface{}) (err error) {
	var embed *EmbedAuthor
	var valid bool
	if embed, valid = other.(*EmbedAuthor); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *EmbedAuthor")
		return
	}

	embed.Name = c.Name
	embed.URL = c.URL
	embed.IconURL = c.IconURL
	embed.ProxyIconURL = c.ProxyIconURL
	return nil
}

// EmbedFooter https://discord.com/developers/docs/resources/channel#embed-object-embed-footer-structure
type EmbedFooter struct {
	Text         string `json:"text"`                     //  | , url of author
	IconURL      string `json:"icon_url,omitempty"`       // ?| , url of footer icon (only supports http(s) and attachments)
	ProxyIconURL string `json:"proxy_icon_url,omitempty"` // ?| , a proxied url of footer icon
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *EmbedFooter) deepCopy() interface{} {
	cp := &EmbedFooter{}
	_ = DeepCopyOver(cp, c)
	return cp
}

// CopyOverTo see interface at struct.go#Copier
func (c *EmbedFooter) copyOverTo(other interface{}) (err error) {
	var embed *EmbedFooter
	var valid bool
	if embed, valid = other.(*EmbedFooter); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *EmbedFooter")
		return
	}

	embed.Text = c.Text
	embed.IconURL = c.IconURL
	embed.ProxyIconURL = c.ProxyIconURL
	return nil
}

// EmbedField https://discord.com/developers/docs/resources/channel#embed-object-embed-field-structure
type EmbedField struct {
	Name   string `json:"name"`             //  | , name of the field
	Value  string `json:"value"`            //  | , value of the field
	Inline bool   `json:"inline,omitempty"` // ?| , whether or not this field should display inline
}

// DeepCopy see interface at struct.go#DeepCopier
func (c *EmbedField) deepCopy() interface{} {
	cp := &EmbedField{}
	_ = DeepCopyOver(cp, c)
	return cp
}

// CopyOverTo see interface at struct.go#Copier
func (c *EmbedField) copyOverTo(other interface{}) (err error) {
	var embed *EmbedField
	var valid bool
	if embed, valid = other.(*EmbedField); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *EmbedField")
		return
	}

	embed.Name = c.Name
	embed.Value = c.Value
	embed.Inline = c.Inline
	return nil
}
