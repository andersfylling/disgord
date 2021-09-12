package disgord

//go:generate go run internal/generate/events/main.go

// This file contains resource objects for the event reactor

import (
	"errors"

	"github.com/andersfylling/disgord/json"
)

// Resource represents a discord event.
// This is used internally for readability only.
type resource = interface{}

// ---------------------------

type EventType interface {
	evtResource
}

type evtResource interface {
	setShardID(id uint)
}

// ---------------------------

// Ready contains the initial state information
type Ready struct {
	APIVersion int                 `json:"v"`
	User       *User               `json:"user"`
	Guilds     []*GuildUnavailable `json:"guilds"`

	// not really needed, as it is handled on the socket layer.
	SessionID string `json:"session_id"`

	// private_channels will be an empty array. As bots receive private messages,
	// they will be notified via Channel Create events.
	//PrivateChannels []*channel.Channel `json:"private_channels"`

	// bot can't have presences
	//Presences []*Presence         `json:"presences"`

	// bot cant have relationships
	//Relationships []interface{} `son:"relationships"`

	// bot can't have user settings
	// UserSettings interface{}        `json:"user_settings"`

	ShardID uint `json:"-"`
}

// ---------------------------

// Resumed response to Resume
type Resumed struct {
	ShardID uint `json:"-"`
}

// ---------------------------

// ChannelCreate new channel created
type ChannelCreate struct {
	Channel *Channel `json:"channel"`
	ShardID uint     `json:"-"`
}

// UnmarshalJSON ...
func (c *ChannelCreate) UnmarshalJSON(data []byte) error {
	c.Channel = &Channel{}
	return json.Unmarshal(data, c.Channel)
}

// ---------------------------

// ChannelUpdate channel was updated
type ChannelUpdate struct {
	Channel *Channel `json:"channel"`
	ShardID uint     `json:"-"`
}

// UnmarshalJSON ...
func (obj *ChannelUpdate) UnmarshalJSON(data []byte) error {
	obj.Channel = &Channel{}
	return json.Unmarshal(data, obj.Channel)
}

// ---------------------------

// ChannelDelete channel was deleted
type ChannelDelete struct {
	Channel *Channel `json:"channel"`
	ShardID uint     `json:"-"`
}

// UnmarshalJSON ...
func (obj *ChannelDelete) UnmarshalJSON(data []byte) error {
	obj.Channel = &Channel{}
	return json.Unmarshal(data, obj.Channel)
}

// ---------------------------

// ChannelPinsUpdate message was pinned or unpinned. Not sent when a message is deleted.
type ChannelPinsUpdate struct {
	ChannelID        Snowflake `json:"channel_id"`
	GuildID          Snowflake `json:"guild_id,omitempty"`
	LastPinTimestamp Time      `json:"last_pin_timestamp,omitempty"`
	ShardID          uint      `json:"-"`
}

// ---------------------------

// TypingStart user started typing in a channel
type TypingStart struct {
	ChannelID     Snowflake `json:"channel_id"`
	GuildID       Snowflake `json:"guild_id"`
	UserID        Snowflake `json:"user_id"`
	Member        *Member   `json:"member"`
	TimestampUnix int       `json:"timestamp"`
	ShardID       uint      `json:"-"`
}

// ---------------------------

// InviteDelete Sent when an invite is deleted.
type InviteDelete struct {
	ChannelID Snowflake `json:"channel_id"`
	GuildID   Snowflake `json:"guild_id"`
	Code      string    `json:"code"`
	ShardID   uint      `json:"-"`
}

// ---------------------------

// MessageCreate message was created
type MessageCreate struct {
	Message *Message
	ShardID uint `json:"-"`
}

var _ Reseter = (*MessageCreate)(nil)
var _ internalUpdater = (*MessageCreate)(nil)

func (obj *MessageCreate) updateInternals() {
	obj.Message.updateInternals()
}

// UnmarshalJSON ...
func (obj *MessageCreate) UnmarshalJSON(data []byte) error {
	obj.Message = &Message{}
	if err := json.Unmarshal(data, obj.Message); err != nil {
		return err
	}
	if obj.Message.Member != nil {
		obj.Message.Member.GuildID = obj.Message.GuildID
	}
	return nil
}

// ---------------------------

// MessageUpdate message was edited
type MessageUpdate struct {
	Message *Message
	ShardID uint `json:"-"`
}

var _ internalUpdater = (*MessageUpdate)(nil)

func (obj *MessageUpdate) updateInternals() {
	obj.Message.updateInternals()
}

// UnmarshalJSON ...
func (obj *MessageUpdate) UnmarshalJSON(data []byte) error {
	obj.Message = &Message{}
	if err := json.Unmarshal(data, obj.Message); err != nil {
		return err
	}
	if obj.Message.Member != nil {
		obj.Message.Member.GuildID = obj.Message.GuildID
	}
	return nil
}

// ---------------------------

// MessageDelete message was deleted
type MessageDelete struct {
	MessageID Snowflake `json:"id"`
	ChannelID Snowflake `json:"channel_id"`
	GuildID   Snowflake `json:"guild_id,omitempty"`
	ShardID   uint      `json:"-"`
}

// ---------------------------

// MessageDeleteBulk multiple messages were deleted at once
type MessageDeleteBulk struct {
	MessageIDs []Snowflake `json:"ids"`
	ChannelID  Snowflake   `json:"channel_id"`
	ShardID    uint        `json:"-"`
}

// ---------------------------

// MessageReactionAdd user reacted to a message
// Note! do not cache emoji, unless it's updated with guildID
// TODO: find guildID when given UserID, ChannelID and MessageID
type MessageReactionAdd struct {
	UserID    Snowflake `json:"user_id"`
	ChannelID Snowflake `json:"channel_id"`
	MessageID Snowflake `json:"message_id"`
	// PartialEmoji id and name. id might be nil
	PartialEmoji *Emoji `json:"emoji"`
	ShardID      uint   `json:"-"`
}

// ---------------------------

type InteractionCreate struct {
	ID            Snowflake                          `json:"id"`
	ApplicationID Snowflake                          `json:"application_id"`
	Type          InteractionType                    `json:"type"`
	Data          *ApplicationCommandInteractionData `json:"data"`
	GuildID       Snowflake                          `json:"guild_id"`
	ChannelID     Snowflake                          `json:"channel_id"`
	Member        *Member                            `json:"member"`
	User          *User                              `json:"user"`
	Token         string                             `json:"token"`
	Version       int                                `json:"version"`
	Message       *Message                           `json:"message"`
	ShardID       uint                               `json:"-"`
}

// ---------------------------

// MessageReactionRemove user removed a reaction from a message
// Note! do not cache emoji, unless it's updated with guildID
// TODO: find guildID when given UserID, ChannelID and MessageID
type MessageReactionRemove struct {
	UserID    Snowflake `json:"user_id"`
	ChannelID Snowflake `json:"channel_id"`
	MessageID Snowflake `json:"message_id"`
	// PartialEmoji id and name. id might be nil
	PartialEmoji *Emoji `json:"emoji"`
	ShardID      uint   `json:"-"`
}

// ---------------------------

// MessageReactionRemoveAll all reactions were explicitly removed from a message
type MessageReactionRemoveAll struct {
	ChannelID Snowflake `json:"channel_id"`
	MessageID Snowflake `json:"message_id"`
	ShardID   uint      `json:"-"`
}

// ---------------------------

// MessageReactionRemoveEmoji Sent when a bot removes all instances of a given emoji from the reactions of a message
type MessageReactionRemoveEmoji struct {
	ChannelID Snowflake `json:"channel_id"`
	GuildID   Snowflake `json:"guild_id"`
	MessageID Snowflake `json:"message_id"`
	Emoji     *Emoji    `json:"emoji"`
	ShardID   uint      `json:"-"`
}

// ---------------------------

// GuildEmojisUpdate guild emojis were updated
type GuildEmojisUpdate struct {
	GuildID Snowflake `json:"guild_id"`
	Emojis  []*Emoji  `json:"emojis"`
	ShardID uint      `json:"-"`
}

var _ internalUpdater = (*GuildEmojisUpdate)(nil)

func (g *GuildEmojisUpdate) updateInternals() {
}

// ---------------------------

// GuildCreate This event can be sent in three different scenarios:
//  1. When a user is initially connecting, to lazily load and backfill information for all unavailable Guilds
//     sent in the Ready event.
//	2. When a Guild becomes available again to the Client.
// 	3. When the current user joins a new Guild.
type GuildCreate struct {
	Guild   *Guild `json:"guild"`
	ShardID uint   `json:"-"`
}

var _ internalUpdater = (*GuildCreate)(nil)

func (g *GuildCreate) updateInternals() {
	g.Guild.updateInternals()
}

// UnmarshalJSON ...
func (obj *GuildCreate) UnmarshalJSON(data []byte) error {
	obj.Guild = &Guild{}
	if err := json.Unmarshal(data, obj.Guild); err != nil {
		return err
	}
	for _, v := range obj.Guild.Members {
		v.GuildID = obj.Guild.ID
	}
	return nil
}

// ---------------------------

// GuildUpdate guild was updated
type GuildUpdate struct {
	Guild   *Guild `json:"guild"`
	ShardID uint   `json:"-"`
}

var _ internalUpdater = (*GuildUpdate)(nil)

func (g *GuildUpdate) updateInternals() {
	g.Guild.updateInternals()
}

// UnmarshalJSON ...
func (obj *GuildUpdate) UnmarshalJSON(data []byte) error {
	obj.Guild = &Guild{}
	return json.Unmarshal(data, obj.Guild)
}

// ---------------------------

// GuildDelete guild became unavailable, or user left/was removed from a guild
type GuildDelete struct {
	UnavailableGuild *GuildUnavailable `json:"guild_unavailable"`
	ShardID          uint              `json:"-"`
}

// UserWasRemoved ... TODO
func (obj *GuildDelete) UserWasRemoved() bool {
	return obj.UnavailableGuild.Unavailable == false
}

// UnmarshalJSON ...
func (obj *GuildDelete) UnmarshalJSON(data []byte) error {
	obj.UnavailableGuild = &GuildUnavailable{}
	return json.Unmarshal(data, obj.UnavailableGuild)
}

// ---------------------------

// GuildBanAdd user was banned from a guild
type GuildBanAdd struct {
	GuildID Snowflake `json:"guild_id"`
	User    *User     `json:"user"`
	ShardID uint      `json:"-"`
}

// ---------------------------

// GuildBanRemove user was unbanned from a guild
type GuildBanRemove struct {
	GuildID Snowflake `json:"guild_id"`
	User    *User     `json:"user"`
	ShardID uint      `json:"-"`
}

// ---------------------------

// GuildIntegrationsUpdate guild integration was updated
type GuildIntegrationsUpdate struct {
	GuildID Snowflake `json:"guild_id"`
	ShardID uint      `json:"-"`
}

// ---------------------------

// GuildMemberAdd new user joined a guild
type GuildMemberAdd struct {
	Member  *Member `json:"member"`
	ShardID uint    `json:"-"`
}

var _ internalUpdater = (*GuildMemberAdd)(nil)

func (g *GuildMemberAdd) updateInternals() {
	g.Member.updateInternals()
}

// UnmarshalJSON ...
func (obj *GuildMemberAdd) UnmarshalJSON(data []byte) error {
	obj.Member = &Member{}
	return json.Unmarshal(data, obj.Member)
}

// ---------------------------

// GuildMemberRemove user was removed from a guild
type GuildMemberRemove struct {
	GuildID Snowflake `json:"guild_id"`
	User    *User     `json:"user"`
	ShardID uint      `json:"-"`
}

// ---------------------------

// GuildMemberUpdate guild member was updated
type GuildMemberUpdate struct {
	*Member
	ShardID uint `json:"-"`
}

// ---------------------------

// GuildMembersChunk response to Request Guild Members
type GuildMembersChunk struct {
	GuildID    Snowflake         `json:"guild_id"`
	Members    []*Member         `json:"members"`
	ChunkIndex uint              `json:"chunk_index"`
	ChunkCount uint              `json:"chunk_count"`
	NotFound   []interface{}     `json:"not_found"`
	Presences  []*PresenceUpdate `json:"presences"`
	Nonce      string            `json:"nonce"`
	ShardID    uint              `json:"-"`
}

var _ internalUpdater = (*GuildMembersChunk)(nil)

func (g *GuildMembersChunk) updateInternals() {
	for i := range g.Members {
		g.Members[i].GuildID = g.GuildID
		g.Members[i].updateInternals()
	}
	for i := range g.Presences {
		executeInternalUpdater(g.Presences[i])
	}
}

// ---------------------------

// GuildRoleCreate guild role was created
type GuildRoleCreate struct {
	GuildID Snowflake `json:"guild_id"`
	Role    *Role     `json:"role"`
	ShardID uint      `json:"-"`
}

var _ internalUpdater = (*GuildRoleCreate)(nil)

func (g *GuildRoleCreate) updateInternals() {
	g.Role.guildID = g.GuildID
}

// ---------------------------

// GuildRoleUpdate guild role was updated
type GuildRoleUpdate struct {
	GuildID Snowflake `json:"guild_id"`
	Role    *Role     `json:"role"`
	ShardID uint      `json:"-"`
}

var _ internalUpdater = (*GuildRoleUpdate)(nil)

func (g *GuildRoleUpdate) updateInternals() {
	g.Role.guildID = g.GuildID
}

// ---------------------------

// GuildRoleDelete a guild role was deleted
type GuildRoleDelete struct {
	GuildID Snowflake `json:"guild_id"`
	RoleID  Snowflake `json:"role_id"`
	ShardID uint      `json:"-"`
}

// ---------------------------

// PresenceUpdate user's presence was updated in a guild
type PresenceUpdate struct {
	User         *User        `json:"user"`
	GuildID      Snowflake    `json:"guild_id"`
	Status       string       `json:"status"`
	Activities   []*Activity  `json:"activities"`
	ClientStatus ClientStatus `json:"client_status"`
	ShardID      uint         `json:"-"`
}

func (h *PresenceUpdate) Game() (*Activity, error) {
	if len(h.Activities) > 0 {
		return h.Activities[0], nil
	} else {
		return nil, errors.New("no activities")
	}
}

// ---------------------------

// UserUpdate properties about a user changed
type UserUpdate struct {
	*User
	ShardID uint `json:"-"`
}

// UnmarshalJSON ...
func (obj *UserUpdate) UnmarshalJSON(data []byte) error {
	obj.User = &User{}
	return json.Unmarshal(data, obj.User)
}

// ---------------------------

// VoiceStateUpdate someone joined, left, or moved a voice channel
type VoiceStateUpdate struct {
	*VoiceState
	ShardID uint `json:"-"`
}

// UnmarshalJSON ...
func (h *VoiceStateUpdate) UnmarshalJSON(data []byte) error {
	h.VoiceState = &VoiceState{}
	return json.Unmarshal(data, h.VoiceState)
}

// ---------------------------

// VoiceServerUpdate guild's voice server was updated. Sent when a guild's voice server is updated. This is sent when initially
// connecting to voice, and when the current voice instance fails over to a new server.
type VoiceServerUpdate struct {
	Token    string    `json:"token"`
	GuildID  Snowflake `json:"guild_id"`
	Endpoint string    `json:"endpoint"`
	ShardID  uint      `json:"-"`
}

// ---------------------------

// WebhooksUpdate guild channel webhook was created, update, or deleted
type WebhooksUpdate struct {
	GuildID   Snowflake `json:"guild_id"`
	ChannelID Snowflake `json:"channel_id"`
	ShardID   uint      `json:"-"`
}

// InviteCreate guild invite was created
type InviteCreate struct {
	ChannelID  Snowflake `json:"channel_id"`
	Code       string    `json:"code"`
	CreatedAt  Time      `json:"created_at"`
	GuildID    Snowflake `json:"guild_id,omitempty"`
	Inviter    *User     `json:"inviter"`
	MaxAge     int       `json:"max_age"`
	MaxUses    int       `json:"max_uses"`
	Target     *User     `json:"target_user,omitempty"`
	TargetType int       `json:"target_user_type"`
	Temporary  bool      `json:"temporary"`
	Uses       int       `json:"uses"`

	ShardID uint `json:"-"`
}
