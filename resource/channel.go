package resource

import (
	"errors"

	"sync"

	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/request"
	"github.com/andersfylling/snowflake"
)

const (
	// Channel types
	// https://discordapp.com/developers/docs/resources/channel#channel-object-channel-types
	ChannelTypeGuildText uint = iota
	ChannelTypeDM
	ChannelTypeGuildVoice
	ChannelTypeGroupDM
	ChannelTypeGuildCategory
)

// ChannelMessager Methods required to create a new DM (or use an existing one) and send a DM.
type ChannelMessager interface {
	CreateMessage(*Message) error // TODO: check cache for `SEND_MESSAGES` and `SEND_TTS_MESSAGES` permissions before sending.
}

func NewChannel() *Channel {
	return &Channel{}
}

// Channel
type Channel struct {
	ID                   snowflake.ID                 `json:"id"`
	Type                 uint                         `json:"type"`
	GuildID              snowflake.ID                 `json:"guild_id,omitempty"`              // ?|
	Position             uint                         `json:"position,omitempty"`              // ?|
	PermissionOverwrites []ChannelPermissionOverwrite `json:"permission_overwrites,omitempty"` // ?|
	Name                 string                       `json:"name,omitempty"`                  // ?|
	Topic                string                       `json:"topic,omitempty"`                 // ?|
	NSFW                 bool                         `json:"nsfw,omitempty"`                  // ?|
	LastMessageID        snowflake.ID                 `json:"last_message_id,omitempty"`       // ?|?, pointer
	Bitrate              uint                         `json:"bitrate,omitempty"`               // ?|
	UserLimit            uint                         `json:"user_limit,omitempty"`            // ?|
	Recipients           []*User                      `json:"recipient,omitempty"`             // ?| , empty if not DM
	Icon                 string                       `json:"icon,omitempty"`                  // ?|?, pointer
	OwnerID              snowflake.ID                 `json:"owner_id,omitempty"`              // ?|
	ApplicationID        snowflake.ID                 `json:"applicaiton_id,omitempty"`        // ?|
	ParentID             snowflake.ID                 `json:"parent_id,omitempty"`             // ?|?, pointer
	LastPingTimestamp    discord.Timestamp            `json:"last_ping_timestamp,omitempty"`   // ?|

	mu sync.RWMutex `json:"-"`
}
type PartialChannel = Channel

func (c *Channel) Mention() string {
	return "<#" + c.ID.String() + ">"
}

func (c *Channel) Compare(other *Channel) bool {
	// eh
	return (c == nil && other == nil) || (other != nil && c.ID == other.ID)
}

func (c *Channel) Replicate(channel *Channel, recipients []*User) {
	// TODO: mutex is copied
	*c = *channel

	// WARNING: DM channels holds users. These should be fetched from cache.
	if recipients != nil && len(recipients) > 0 {
		c.Recipients = recipients
	} else {
		c.Recipients = []*User{}
	}
}

func (c *Channel) DeepCopy() *Channel {
	channel := NewChannel()

	c.mu.RLock()

	channel.ID = c.ID
	channel.Type = c.Type
	channel.GuildID = c.GuildID
	channel.Position = c.Position
	channel.PermissionOverwrites = c.PermissionOverwrites
	channel.Name = c.Name
	channel.Topic = c.Topic
	channel.NSFW = c.NSFW
	channel.LastMessageID = c.LastMessageID
	channel.Bitrate = c.Bitrate
	channel.UserLimit = c.UserLimit
	channel.Icon = c.Icon
	channel.OwnerID = c.OwnerID
	channel.ApplicationID = c.ApplicationID
	channel.ParentID = c.ParentID
	channel.LastPingTimestamp = c.LastPingTimestamp

	// add recipients if it's a DM
	if c.Type == ChannelTypeDM || c.Type == ChannelTypeGroupDM {
		for _, recipient := range c.Recipients {
			channel.Recipients = append(channel.Recipients, recipient.DeepCopy())
		}
	}

	c.mu.RUnlock()

	return channel
}

func (c *Channel) Clear() {
	// TODO
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

// ReqGetChannel [GET] 	   Get a channel by ID. Returns a channel object.
// Endpoint				   /channels/{channel.id}
// Rate limiter [MAJOR]	   /channels/{channel.id}
// Discord documentation   https://discordapp.com/developers/docs/resources/channel#get-channel
// Reviewed				   2018-06-07
// Comment				   -
func ReqGetChannel(requester request.DiscordGetter, id snowflake.ID) (*Channel, error) {
	if id.Empty() {
		return nil, errors.New("not a valid snowflake")
	}

	uri := "/channels/" + id.String()
	content := &Channel{}
	_, err := requester.Get(uri, uri, content)
	return content, err
}

// ModifyChannelParams https://discordapp.com/developers/docs/resources/channel#modify-channel-json-params
type ModifyChannelParams = Channel

// ReqModifyChannel [PUT/PATCH] Update a channels settings. Requires the 'MANAGE_CHANNELS' permission for the guild.
// 								Returns a channel on success, and a 400 BAD REQUEST on invalid parameters. Fires a
// 								Channel Update Gateway event. If modifying a category, individual Channel Update
// 								events will fire for each child channel that also changes. For the PATCH method,
// 								all the JSON Params are optional.
// Endpoint				   		/channels/{channel.id}
// Rate limiter [MAJOR]	   		/channels/{channel.id}
// Discord documentation   		https://discordapp.com/developers/docs/resources/channel#modify-channel
// Reviewed				   		2018-06-07
// Comment				   		-
func ReqModifyChannelPatch(client request.DiscordPatcher, changes *ModifyChannelParams) (*Channel, error) {
	if changes.ID.Empty() {
		return nil, errors.New("not a valid snowflake")
	}

	//uri := "/channels/" + changes.ID.String()
	//data, err := json.Marshal(changes)
	//if err != nil {
	//	return nil, err
	//}
	//err := client.Request("PUT", uri, bytes.NewBuffer(data)) // TODO implement "PATCH" logic
	return nil, nil
}

// ReqModifyChannelUpdate see ReqModifyChannelPatch
func ReqModifyChannelUpdate(client request.DiscordPutter, changes *ModifyChannelParams) (*Channel, error) {
	if changes.ID.Empty() {
		return nil, errors.New("not a valid snowflake")
	}

	//uri := "/channels/" + changes.ID.String()
	//data, err := json.Marshal(changes)
	//if err != nil {
	//	return nil, err
	//}
	//err := client.Request("PUT", uri, bytes.NewBuffer(data)) // TODO implement "PUT" logic
	return nil, nil
}

// ReqDeleteChannel [DELETE]	Delete a channel, or close a private message. Requires the 'MANAGE_CHANNELS'
// 								permission for the guild. Deleting a category does not delete its child
// 								channels; they will have their parent_id removed and a Channel Update Gateway
// 								event will fire for each of them. Returns a channel object on success. Fires a
// 								Channel Delete Gateway event.
// Endpoint				   		/channels/{channel.id}
// Rate limiter [MAJOR]	   		/channels/{channel.id}
// Discord documentation   		https://discordapp.com/developers/docs/resources/channel#deleteclose-channel
// Reviewed				   		2018-06-07
// Comment				   		Deleting a guild channel cannot be undone. Use this with caution, as it
// 								is impossible to undo this action when performed on a guild channel. In
// 								contrast, when used with a private message, it is possible to undo the
// 								action by opening a private message with the recipient again.
func ReqDeleteChannel(client request.DiscordDeleter, id snowflake.ID) (err error) {
	if id.Empty() {
		err = errors.New("not a valid snowflake")
		return
	}

	uri := "/channels/" + id.String()
	_, err = client.Delete(uri, uri)
	return
}

// -------
// Attachment

// Attachment https://discordapp.com/developers/docs/resources/channel#attachment-object
type Attachment struct {
	ID       snowflake.ID `json:"id"`
	Filename string       `json:"filename"`
	Size     uint         `json:"size"`
	URL      string       `json:"url"`
	ProxyURL string       `json:"proxy_url"`
	Height   uint         `json:"height"`
	Width    uint         `json:"width"`
}

// ---------------

// Overwrite: https://discordapp.com/developers/docs/resources/channel#overwrite-object
type ChannelPermissionOverwrite struct {
	ID    snowflake.ID `json:"id"`    // role or user id
	Type  string       `json:"type"`  // either `role` or `member`
	Allow int          `json:"allow"` // permission bit set
	Deny  int          `json:"deny"`  // permission bit set
}

func (pmo *ChannelPermissionOverwrite) Clear() {}