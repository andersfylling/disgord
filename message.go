package disgord

import (
	"time"

	"github.com/andersfylling/snowflake"
)

type Message struct {
	ID              snowflake.ID   `json:"id"`
	ChannelID       snowflake.ID   `json:"channel_id"`
	Author          *User          `json:"author"`
	Content         string         `json:"content"`
	Timestamp       time.Time      `json:"timestamp"`
	EditedTimestamp time.Time      `json:"edited_timestamp"` // ?
	Tts             bool           `json:"tts"`
	MentionEveryone bool           `json:"mention_everyone"`
	Mentions        []*User        `json:"mentions"`
	MentionRoles    []snowflake.ID `json:"mention_roles"`
	Attachments     []*Attachment  `json:"attachments"`
	Embeds          []*Embed       `json:"embeds"`
	Reactions       []*Reaction    `json:"reactions"` // ?
	Nonce           snowflake.ID   `json:"nonce"`     // ?, used for validating a message was sent
	Pinned          bool           `json:"pinned"`
	WebhookID       snowflake.ID   `json:"webhook_id"` // ?
	Type            uint           `json:"type"`
}
