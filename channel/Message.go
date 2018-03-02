package channel

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/snowflake"
)

const (
	_ int = iota
	MessageActivityTypeJoin
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	MessageActivityTypeJoinRequest
)

func NewMessage() *Message {
	return &Message{}
}

func NewDeletedMessage() *DeletedMessage {
	return &DeletedMessage{}
}

type Attachment struct {
	ID       snowflake.ID `json:"id"`
	Filename string       `json:"filename"`
	Size     uint         `json:"size"`
	URL      string       `json:"url"`
	ProxyURL string       `json:"proxy_url"`
	Height   uint         `json:"height"`
	Width    uint         `json:"width"`
}

type DeletedMessage struct {
	ID        snowflake.ID `json:"id"`
	ChannelID snowflake.ID `json:"channel_id"`
}

type MessageActivity struct {
	Type    int    `json:"type"`
	PartyID string `json:"party_id"`
}

type MessageApplication struct {
	ID          snowflake.ID `json:"id"`
	CoverImage  string       `json:"cover_image"`
	Description string       `json:"description"`
	Icon        string       `json:"icon"`
	Name        string       `json:"name"`
}

type Message struct {
	ID              snowflake.ID       `json:"id"`
	ChannelID       snowflake.ID       `json:"channel_id"`
	Author          *user.User         `json:"author"`
	Content         string             `json:"content"`
	Timestamp       time.Time          `json:"timestamp"`
	EditedTimestamp time.Time          `json:"edited_timestamp"` // ?
	Tts             bool               `json:"tts"`
	MentionEveryone bool               `json:"mention_everyone"`
	Mentions        []*user.User       `json:"mentions"`
	MentionRoles    []snowflake.ID     `json:"mention_roles"`
	Attachments     []*Attachment      `json:"attachments"`
	Embeds          []*Embed           `json:"embeds"`
	Reactions       []*Reaction        `json:"reactions"` // ?
	Nonce           snowflake.ID       `json:"nonce"`     // ?, used for validating a message was sent
	Pinned          bool               `json:"pinned"`
	WebhookID       snowflake.ID       `json:"webhook_id"` // ?
	Type            uint               `json:"type"`
	Activity        MessageActivity    `json:"activity"`
	Application     MessageApplication `json:"application"`

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

func GetMessages() {}
