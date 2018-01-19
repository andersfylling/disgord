package channel

import (
	"time"
)

// Embed ...
type Embed struct {
	Title       string          `json:"title"`       // title of embed
	Type        string          `json:"type"`        // type of embed (always "rich" for webhook embeds)
	Description string          `json:"description"` // description of embed
	URL         string          `json:"url"`         // url of embed
	Timestamp   time.Time       `json:"timestamp"`   // timestamp	timestamp of embed content
	Color       int             `json:"color"`       // color code of the embed
	Footer      *EmbedFooter    `json:"footer"`      // embed footer object	footer information
	Image       *EmbedImage     `json:"image"`       // embed image object	image information
	Thumbnail   *EmbedThumbnail `json:"thumbnail"`   // embed thumbnail object	thumbnail information
	Video       *EmbedVideo     `json:"video"`       // embed video object	video information
	Provider    *EmbedProvider  `json:"provider"`    // embed provider object	provider information
	Author      *EmbedAuthor    `json:"author"`      // embed author object	author information
	Fields      []*EmbedField   `json:"fields"`      //	array of embed field objects	fields information
}

type EmbedFooter struct{}
type EmbedImage struct{}
type EmbedThumbnail struct{}
type EmbedVideo struct{}
type EmbedProvider struct{}
type EmbedAuthor struct{}
type EmbedField struct{}
