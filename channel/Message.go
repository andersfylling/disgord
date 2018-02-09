package channel

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/snowflake"
)

func NewMessage() *Message {
	return &Message{}
}

func NewDeletedMessage() *DeletedMessage {
	return &DeletedMessage{}
}

type DeletedMessage struct {
	ID        snowflake.ID `json:"id"`
	ChannelID snowflake.ID `json:"channel_id"`
}

type Message struct {
	ID              snowflake.ID   `json:"id"`
	ChannelID       snowflake.ID   `json:"channel_id"`
	Author          *user.User     `json:"author"`
	Content         string         `json:"content"`
	Timestamp       time.Time      `json:"timestamp"`
	EditedTimestamp time.Time      `json:"edited_timestamp"` // ?
	Tts             bool           `json:"tts"`
	MentionEveryone bool           `json:"mention_everyone"`
	Mentions        []*user.User   `json:"mentions"`
	MentionRoles    []snowflake.ID `json:"mention_roles"`
	Attachments     []*Attachment  `json:"attachments"`
	Embeds          []*Embed       `json:"embeds"`
	Reactions       []*Reaction    `json:"reactions"` // ?
	Nonce           snowflake.ID   `json:"nonce"`     // ?, used for validating a message was sent
	Pinned          bool           `json:"pinned"`
	WebhookID       snowflake.ID   `json:"webhook_id"` // ?
	Type            uint           `json:"type"`

	sync.RWMutex `json:"-"`
}

func (m *Message) MarshalJSON() ([]byte, error) {
	if m.ID.Empty() {
		return []byte("{}"), nil
	}

	return json.Marshal(Message(*m))
}

// func (m *Message) UnmarshalJSON(data []byte) error {
// 	return json.Unmarshal(data, &m.messageJSON)
// }

type TypingStart struct {
	UserID    snowflake.ID `json:"user_id"`
	Timestamp int          `json:"timestamp"` // unix timestamp
	ChannelID snowflake.ID `json:"channel_id"`
}
