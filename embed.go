package disgord

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

var _ DeepCopier = (*Embed)(nil)

// DeepCopy see interface at struct.go#DeepCopier
func (e *Embed) DeepCopy() (copy interface{}) {
	copy = &Embed{}
	_ = e.CopyOverTo(copy)

	return
}

// EmbedThumbnail https://discordapp.com/developers/docs/resources/channel#embed-object-embed-thumbnail-structure
type EmbedThumbnail struct {
	Lockable `json:"-"`

	URL      string `json:"url,omitempty"`       // ?| , source url of image (only supports http(s) and attachments)
	ProxyURL string `json:"proxy_url,omitempty"` // ?| , a proxied url of the image
	Height   int    `json:"height,omitempty"`    // ?| , height of image
	Width    int    `json:"width,omitempty"`     // ?| , width of image
}

var _ DeepCopier = (*EmbedThumbnail)(nil)

// DeepCopy see interface at struct.go#DeepCopier
func (e *EmbedThumbnail) DeepCopy() (copy interface{}) {
	copy = &EmbedThumbnail{}
	_ = e.CopyOverTo(copy)

	return
}

// EmbedVideo https://discordapp.com/developers/docs/resources/channel#embed-object-embed-video-structure
type EmbedVideo struct {
	Lockable `json:"-"`

	URL    string `json:"url,omitempty"`    // ?| , source url of video
	Height int    `json:"height,omitempty"` // ?| , height of video
	Width  int    `json:"width,omitempty"`  // ?| , width of video
}

var _ DeepCopier = (*EmbedVideo)(nil)

// DeepCopy see interface at struct.go#DeepCopier
func (e *EmbedVideo) DeepCopy() (copy interface{}) {
	copy = &EmbedVideo{}
	_ = e.CopyOverTo(copy)

	return
}

// EmbedImage https://discordapp.com/developers/docs/resources/channel#embed-object-embed-image-structure
type EmbedImage struct {
	Lockable `json:"-"`

	URL      string `json:"url,omitempty"`       // ?| , source url of image (only supports http(s) and attachments)
	ProxyURL string `json:"proxy_url,omitempty"` // ?| , a proxied url of the image
	Height   int    `json:"height,omitempty"`    // ?| , height of image
	Width    int    `json:"width,omitempty"`     // ?| , width of image
}

var _ DeepCopier = (*EmbedImage)(nil)

// DeepCopy see interface at struct.go#DeepCopier
func (e *EmbedImage) DeepCopy() (copy interface{}) {
	copy = &EmbedImage{}
	_ = e.CopyOverTo(copy)

	return
}

// EmbedProvider https://discordapp.com/developers/docs/resources/channel#embed-object-embed-provider-structure
type EmbedProvider struct {
	Lockable `json:"-"`

	Name string `json:"name,omitempty"` // ?| , name of provider
	URL  string `json:"url,omitempty"`  // ?| , url of provider
}

var _ DeepCopier = (*EmbedProvider)(nil)

// DeepCopy see interface at struct.go#DeepCopier
func (e *EmbedProvider) DeepCopy() (copy interface{}) {
	copy = &EmbedProvider{}
	_ = e.CopyOverTo(copy)

	return
}

// EmbedAuthor https://discordapp.com/developers/docs/resources/channel#embed-object-embed-author-structure
type EmbedAuthor struct {
	Lockable `json:"-"`

	Name         string `json:"name,omitempty"`           // ?| , name of author
	URL          string `json:"url,omitempty"`            // ?| , url of author
	IconURL      string `json:"icon_url,omitempty"`       // ?| , url of author icon (only supports http(s) and attachments)
	ProxyIconURL string `json:"proxy_icon_url,omitempty"` // ?| , a proxied url of author icon
}

var _ DeepCopier = (*EmbedAuthor)(nil)

// DeepCopy see interface at struct.go#DeepCopier
func (e *EmbedAuthor) DeepCopy() (copy interface{}) {
	copy = &EmbedAuthor{}
	_ = e.CopyOverTo(copy)

	return
}

// EmbedFooter https://discordapp.com/developers/docs/resources/channel#embed-object-embed-footer-structure
type EmbedFooter struct {
	Lockable `json:"-"`

	Text         string `json:"text"`                     //  | , url of author
	IconURL      string `json:"icon_url,omitempty"`       // ?| , url of footer icon (only supports http(s) and attachments)
	ProxyIconURL string `json:"proxy_icon_url,omitempty"` // ?| , a proxied url of footer icon
}

var _ DeepCopier = (*EmbedFooter)(nil)

// DeepCopy see interface at struct.go#DeepCopier
func (e *EmbedFooter) DeepCopy() (copy interface{}) {
	copy = &EmbedFooter{}
	_ = e.CopyOverTo(copy)

	return
}

// EmbedField https://discordapp.com/developers/docs/resources/channel#embed-object-embed-field-structure
type EmbedField struct {
	Lockable `json:"-"`

	Name   string `json:"name"`             //  | , name of the field
	Value  string `json:"value"`            //  | , value of the field
	Inline bool   `json:"inline,omitempty"` // ?| , whether or not this field should display inline
}

var _ DeepCopier = (*EmbedField)(nil)

// DeepCopy see interface at struct.go#DeepCopier
func (e *EmbedField) DeepCopy() (copy interface{}) {
	copy = &EmbedField{}
	_ = e.CopyOverTo(copy)

	return
}
