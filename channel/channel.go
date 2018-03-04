package channel

import (
	"errors"

	"net/http"

	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/snowflake"
)

// ChannelMessager Methods required to create a new DM (or use an existing one) and send a DM.
type ChannelMessager interface {
	CreateMessage(*Message) error // TODO: check cache for `SEND_MESSAGES` and `SEND_TTS_MESSAGES` permissions before sending.
}

type Channel struct {
	ID                   snowflake.ID           `json:"id"`
	Type                 uint                   `json:"type"`
	GuildID              snowflake.ID           `json:"guild_id,omitempty"`
	Position             uint                   `json:"position,omitempty"`
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Name                 string                 `json:"name,omitempty"`
	Topic                string                 `json:"topic,omitempty"`
	NSFW                 bool                   `json:"nsfw,omitempty"`
	LastMessageID        *snowflake.ID          `json:"last_message_id,omitempty"`
	Bitrate              uint                   `json:"bitrate,omitempty"`
	UserLimit            uint                   `json:"user_limit,omitempty"`
	Recipients           []*user.User           `json:"recipient,omitempty"`
	Icon                 *string                `json:"icon,omitempty"`
	OwnerID              snowflake.ID           `json:"owner_id,omitempty"`
	ApplicationID        snowflake.ID           `json:"applicaiton_id,omitempty"`
	ParentID             snowflake.ID           `json:"parent_id,omitempty"`
	LastPingTimestamp    discord.Timestamp      `json:"last_ping_timestamp,omitempty"`

	// Messages used for caching only. is always empty when fresh from the discord API
	Messages []*Message `json:"-"` // should prolly set a cache limit of 100
}

func NewChannel() *Channel {
	return &Channel{}
}

func (c *Channel) Mention() string {
	return "<#" + c.ID.String() + ">"
}

func (c *Channel) Compare(other *Channel) bool {
	// eh
	return (c == nil && other == nil) || (other != nil && c.ID == other.ID)
}

func (c *Channel) Clear() {
	c.LastMessageID = nil
	// c.Icon = nil // Do I really want to clear this?
	for _, pmo := range c.PermissionOverwrites {
		pmo.Clear()
		pmo = nil
	}
	c.PermissionOverwrites = nil

	//for _,
}

func (c *Channel) Update() {

}

func (c *Channel) Delete() {

}

func (c *Channel) Create() {
	// check if channel already exists.
}

func (c *Channel) SendMsgStr(client ChannelMessager, msgStr string) (msg *Message, err error) {
	return &Message{}, errors.New("not implemented")
}

func (c *Channel) SendMsg(client ChannelMessager, msg *Message) (err error) {
	return errors.New("not implemented")
}

// DISCORD HTTP API
// /channels/*
//
type DiscordAPIRequester interface {
	Request(method string, uri string, content interface{}) error
}

// GetChannel Get a channel by ID
func GetChannel(client DiscordAPIRequester, id snowflake.ID) (*Channel, error) {
	if id.Empty() {
		return nil, errors.New("Not a valid snowflake")
	}

	uri := "/channels/" + id.String()
	content := &Channel{}
	err := client.Request(http.MethodGet, uri, content)
	return content, err
}

func UpdateChannel(client DiscordAPIRequester, changes *Channel) (*Channel, error) {
	if changes.ID.Empty() {
		return nil, errors.New("Not a valid snowflake")
	}

	//uri := "/channels/" + changes.ID.String()
	//data, err := json.Marshal(changes)
	//if err != nil {
	//	return nil, err
	//}
	//err := client.Request("PUT", uri, bytes.NewBuffer(data)) // TODO implement "PUT" logic
	return nil, nil
}

func DeleteChannel(client DiscordAPIRequester, id snowflake.ID) (*Channel, error) {
	if id.Empty() {
		return nil, errors.New("Not a valid snowflake")
	}

	uri := "/channels/" + id.String()
	content := &Channel{}
	err := client.Request("DELETE", uri, content)
	return content, err
}
