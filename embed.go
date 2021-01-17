package disgord

// limitations: https://discord.com/developers/docs/resources/channel#embed-limits
// TODO: implement NewEmbedX functions that ensures limitations

type EmbedType string

const (
	EmbedTypeRich    EmbedType = "rich"
	EmbedTypeImage   EmbedType = "image"
	EmbedTypeVideo   EmbedType = "video"
	EmbedTypeGIFV    EmbedType = "gifv"
	EmbedTypeArticle EmbedType = "article"
	EmbedTypeLink    EmbedType = "link"
)

// Embed https://discord.com/developers/docs/resources/channel#embed-object
type Embed struct {
	Title       string          `json:"title,omitempty"`       // title of embed
	Type        EmbedType       `json:"type,omitempty"`        // type of embed (always "rich" for webhook embeds)
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

var _ Copier = (*Embed)(nil)
var _ DeepCopier = (*Embed)(nil)

// EmbedThumbnail https://discord.com/developers/docs/resources/channel#embed-object-embed-thumbnail-structure
type EmbedThumbnail struct {
	URL      string `json:"url,omitempty"`       // ?| , source url of image (only supports http(s) and attachments)
	ProxyURL string `json:"proxy_url,omitempty"` // ?| , a proxied url of the image
	Height   int    `json:"height,omitempty"`    // ?| , height of image
	Width    int    `json:"width,omitempty"`     // ?| , width of image
}

var _ Copier = (*EmbedThumbnail)(nil)
var _ DeepCopier = (*EmbedThumbnail)(nil)

// EmbedVideo https://discord.com/developers/docs/resources/channel#embed-object-embed-video-structure
type EmbedVideo struct {
	URL    string `json:"url,omitempty"`    // ?| , source url of video
	Height int    `json:"height,omitempty"` // ?| , height of video
	Width  int    `json:"width,omitempty"`  // ?| , width of video
}

var _ Copier = (*EmbedVideo)(nil)
var _ DeepCopier = (*EmbedVideo)(nil)

// EmbedImage https://discord.com/developers/docs/resources/channel#embed-object-embed-image-structure
type EmbedImage struct {
	URL      string `json:"url,omitempty"`       // ?| , source url of image (only supports http(s) and attachments)
	ProxyURL string `json:"proxy_url,omitempty"` // ?| , a proxied url of the image
	Height   int    `json:"height,omitempty"`    // ?| , height of image
	Width    int    `json:"width,omitempty"`     // ?| , width of image
}

var _ Copier = (*EmbedImage)(nil)
var _ DeepCopier = (*EmbedImage)(nil)

// EmbedProvider https://discord.com/developers/docs/resources/channel#embed-object-embed-provider-structure
type EmbedProvider struct {
	Name string `json:"name,omitempty"` // ?| , name of provider
	URL  string `json:"url,omitempty"`  // ?| , url of provider
}

var _ Copier = (*EmbedProvider)(nil)
var _ DeepCopier = (*EmbedProvider)(nil)

// EmbedAuthor https://discord.com/developers/docs/resources/channel#embed-object-embed-author-structure
type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`           // ?| , name of author
	URL          string `json:"url,omitempty"`            // ?| , url of author
	IconURL      string `json:"icon_url,omitempty"`       // ?| , url of author icon (only supports http(s) and attachments)
	ProxyIconURL string `json:"proxy_icon_url,omitempty"` // ?| , a proxied url of author icon
}

var _ Copier = (*EmbedAuthor)(nil)
var _ DeepCopier = (*EmbedAuthor)(nil)

// EmbedFooter https://discord.com/developers/docs/resources/channel#embed-object-embed-footer-structure
type EmbedFooter struct {
	Text         string `json:"text"`                     //  | , url of author
	IconURL      string `json:"icon_url,omitempty"`       // ?| , url of footer icon (only supports http(s) and attachments)
	ProxyIconURL string `json:"proxy_icon_url,omitempty"` // ?| , a proxied url of footer icon
}

var _ Copier = (*EmbedFooter)(nil)
var _ DeepCopier = (*EmbedFooter)(nil)

// EmbedField https://discord.com/developers/docs/resources/channel#embed-object-embed-field-structure
type EmbedField struct {
	Name   string `json:"name"`             //  | , name of the field
	Value  string `json:"value"`            //  | , value of the field
	Inline bool   `json:"inline,omitempty"` // ?| , whether or not this field should display inline
}

var _ Copier = (*EmbedField)(nil)
var _ DeepCopier = (*EmbedField)(nil)
