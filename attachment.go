package disgord

import "github.com/andersfylling/snowflake"

type Attachment struct {
	ID       snowflake.ID `json:"id,string"`
	Filename string       `json:"filename"`
	Size     uint         `json:"size"`
	URL      string       `json:"url"`
	ProxyURL string       `json:"proxy_url"`
	Height   uint         `json:"height"`
	Width    uint         `json:"width"`
}
