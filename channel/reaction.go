package channel

import "github.com/andersfylling/disgord/emoji"

type Reaction struct {
	Count uint         `json:"count"`
	Me    bool         `json:"me"`
	Emoji *emoji.Emoji `json:"Emoji"`
}
