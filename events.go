package disgord

import (
	"context"
	"sync"

	"github.com/andersfylling/disgord/event"
)

type eventBox interface {
	registerContext(ctx context.Context)
}

// keys that does not fit within one of the existing files goes here
const EventAllEvents = "*"

// is triggered on every event type
type EventCallback = func(session Session, box interface{})

// ---------------------------

// Holds and array of presence update objects
const EventPresencesReplace = event.PresencesReplace

// holds the event content
type PresencesReplace struct {
	Presnces []*PresenceUpdate `json:"presences_replace"` // TODO: verify json tag
	Ctx      context.Context   `json:"-"`
}

func (p *PresencesReplace) registerContext(ctx context.Context) { p.Ctx = ctx }

// callback for EventPresencesReplace
type PresencesReplaceCallback = func(session Session, pr *PresencesReplace)

// ---------------------------

// The ready event is dispatched when a client has completed the initial handshake with the gateway (for new sessions).
// The ready event can be the largest and most complex event the gateway will send, as it contains all the state
// required for a client to begin interacting with the rest of the platform.
//  Fields:
//  - V int
//  - User *User
//  - PrivateChannels []*Channel
//  - Guilds []*GuildUnavailable
//  - SessionID string
//  - Trace []string
const EventReady = event.Ready

// contains the initial state information
type Ready struct {
	APIVersion int                 `json:"v"`
	User       *User               `json:"user"`
	Guilds     []*GuildUnavailable `json:"guilds"`

	// not really needed, as it is handled on the socket layer.
	SessionID string   `json:"session_id"`
	Trace     []string `json:"_trace"`

	// private_channels will be an empty array. As bots receive private messages,
	// they will be notified via Channel Create events.
	//PrivateChannels []*channel.Channel `json:"private_channels"`

	// bot can't have presences
	//Presences []*Presence         `json:"presences"`

	// bot cant have relationships
	//Relationships []interface{} `son:"relationships"`

	// bot can't have user settings
	// UserSettings interface{}        `json:"user_settings"`

	sync.RWMutex `json:"-"`
	Ctx          context.Context `json:"-"`
}

func (obj *Ready) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on READY events
type ReadyCallback = func(session Session, r *Ready)

// ---------------------------

// The resumed event is dispatched when a client has sent a resume payload to the gateway
// (for resuming existing sessions).
//  Fields:
//  - Trace []string
const EventResumed = event.Resumed

// response to Resume
type Resumed struct {
	Trace []string        `json:"_trace"`
	Ctx   context.Context `json:"-"`
}

func (obj *Resumed) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on RESUME events
type ResumedCallback = func(session Session, r *Resumed)

// ---------------------------

// failure response to Identify or Resume or invalid active session
type InvalidSession struct {
	Ctx context.Context `json:"-"`
}

func (obj *InvalidSession) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on INVALID_SESSION events
type InvalidSessionCallback = func(session Session, is *InvalidSession)

// ---------------------------

// Sent when a new channel is created, relevant to the current user. The inner payload is a DM channel or
// guild channel object.
const EventChannelCreate = event.ChannelCreate

// new channel created
type ChannelCreate struct {
	Channel *Channel        `json:"channel"`
	Ctx     context.Context `json:"-"`
}

func (obj *ChannelCreate) registerContext(ctx context.Context) { obj.Ctx = ctx }
func (obj *ChannelCreate) UnmarshalJSON(data []byte) error {
	obj.Channel = &Channel{}
	return unmarshal(data, obj.Channel)
}

// triggered on CHANNEL_CREATE events
type ChannelCreateCallback = func(session Session, cc *ChannelCreate)

// ---------------------------

// Sent when a channel is updated. The inner payload is a guild channel object.
const EventChannelUpdate = event.ChannelUpdate

// channel was updated
type ChannelUpdate struct {
	Channel *Channel        `json:"channel"`
	Ctx     context.Context `json:"-"`
}

func (obj *ChannelUpdate) registerContext(ctx context.Context) { obj.Ctx = ctx }
func (obj *ChannelUpdate) UnmarshalJSON(data []byte) error {
	obj.Channel = &Channel{}
	return unmarshal(data, obj.Channel)
}

// triggered on CHANNEL_UPDATE events
type ChannelUpdateCallback = func(session Session, cu *ChannelUpdate)

// ---------------------------

// Sent when a channel relevant to the current user is deleted. The inner payload is a DM or Guild channel object.
const EventChannelDelete = event.ChannelDelete

// channel was deleted
type ChannelDelete struct {
	Channel *Channel        `json:"channel"`
	Ctx     context.Context `json:"-"`
}

func (obj *ChannelDelete) registerContext(ctx context.Context) { obj.Ctx = ctx }
func (obj *ChannelDelete) UnmarshalJSON(data []byte) error {
	obj.Channel = &Channel{}
	return unmarshal(data, obj.Channel)
}

// triggered on CHANNEL_DELETE events
type ChannelDeleteCallback = func(session Session, cd *ChannelDelete)

// ---------------------------

// Sent when a message is pinned or unpinned in a text channel. This is not sent when a pinned message is deleted.
//  Fields:
//  - ChannelID int64 or Snowflake
//  - LastPinTimestamp time.Now().UTC().Format(time.RFC3339)
// TODO fix.
const EventChannelPinsUpdate = event.ChannelPinsUpdate

// message was pinned or unpinned
type ChannelPinsUpdate struct {
	// ChannelID snowflake	the id of the channel
	ChannelID Snowflake `json:"channel_id"`

	// LastPinTimestamp	ISO8601 timestamp	the time at which the most recent pinned message was pinned
	LastPinTimestamp Timestamp       `json:"last_pin_timestamp,omitempty"` // ?|
	Ctx              context.Context `json:"-"`
}

func (obj *ChannelPinsUpdate) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on CHANNEL_PINS_UPDATE events
type ChannelPinsUpdateCallback = func(session Session, cpu *ChannelPinsUpdate)

// ---------------------------

// Sent when a user starts typing in a channel.
//  Fields:
//  - ChannelID     Snowflake
//  - UserID        Snowflake
//  - TimestampUnix int
const EventTypingStart = event.TypingStart

// user started typing in a channel
type TypingStart struct {
	ChannelID     Snowflake       `json:"channel_id"`
	UserID        Snowflake       `json:"user_id"`
	TimestampUnix int             `json:"timestamp"`
	Ctx           context.Context `json:"-"`
}

func (obj *TypingStart) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on TYPING_START events
type TypingStartCallback = func(session Session, ts *TypingStart)

// ---------------------------

// Sent when a message is created. The inner payload is a message object.
const EventMessageCreate = event.MessageCreate

// message was created
type MessageCreate struct {
	Message *Message
	Ctx     context.Context `json:"-"`
}

func (obj *MessageCreate) registerContext(ctx context.Context) { obj.Ctx = ctx }
func (obj *MessageCreate) UnmarshalJSON(data []byte) error {
	obj.Message = &Message{}
	return unmarshal(data, obj.Message)
}

// triggered on MESSAGE_CREATE events
type MessageCreateCallback = func(session Session, mc *MessageCreate)

// ---------------------------

// Sent when a message is updated. The inner payload is a message object.
//
// NOTE! Has _at_least_ the GuildID and ChannelID fields.
const EventMessageUpdate = event.MessageUpdate

// message was edited
type MessageUpdate struct {
	Message *Message
	Ctx     context.Context `json:"-"`
}

func (obj *MessageUpdate) registerContext(ctx context.Context) { obj.Ctx = ctx }
func (obj *MessageUpdate) UnmarshalJSON(data []byte) error {
	obj.Message = &Message{}
	return unmarshal(data, obj.Message)
}

// triggered on MESSAGE_UPDATE events
type MessageUpdateCallback = func(session Session, mu *MessageUpdate)

// ---------------------------

// Sent when a message is deleted.
//  Fields:
//  - ID        Snowflake
//  - ChannelID Snowflake
const EventMessageDelete = event.MessageDelete

// message was deleted
type MessageDelete struct {
	MessageID Snowflake       `json:"id"`
	ChannelID Snowflake       `json:"channel_id"`
	Ctx       context.Context `json:"-"`
}

func (obj *MessageDelete) registerContext(ctx context.Context) { obj.Ctx = ctx }
func (obj *MessageDelete) UnmarshalJSON(data []byte) (err error) {
	msg := &Message{}
	err = unmarshal(data, msg)
	if err != nil {
		return
	}

	obj.MessageID = msg.ID
	obj.ChannelID = msg.ChannelID
	return
}

// triggered on MESSAGE_DELETE events
type MessageDeleteCallback = func(session Session, md *MessageDelete)

// ---------------------------

// Sent when multiple messages are deleted at once.
//  Fields:
//  - IDs       []Snowflake
//  - ChannelID Snowflake
const EventMessageDeleteBulk = event.MessageDeleteBulk

// multiple messages were deleted at once
type MessageDeleteBulk struct {
	MessageIDs []Snowflake     `json:"ids"`
	ChannelID  Snowflake       `json:"channel_id"`
	Ctx        context.Context `json:"-"`
}

func (obj *MessageDeleteBulk) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on MESSAGE_DELETE_BULK events
type MessageDeleteBulkCallback = func(session Session, mdb *MessageDeleteBulk)

// ---------------------------

// Sent when a user adds a reaction to a message.
//  Fields:
//  - UserID     Snowflake
//  - ChannelID  Snowflake
//  - MessageID  Snowflake
//  - Emoji      *Emoji
const EventMessageReactionAdd = event.MessageReactionAdd

// MessageReactionAdd	user reacted to a message
type MessageReactionAdd struct {
	UserID    Snowflake `json:"user_id"`
	ChannelID Snowflake `json:"channel_id"`
	MessageID Snowflake `json:"message_id"`
	// PartialEmoji id and name. id might be nil
	PartialEmoji *Emoji          `json:"emoji"`
	Ctx          context.Context `json:"-"`
}

func (obj *MessageReactionAdd) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on MESSAGE_REACTION_ADD events
type MessageReactionAddCallback = func(session Session, mra *MessageReactionAdd)

// ---------------------------

// Sent when a user removes a reaction from a message.
//  Fields:
//  - UserID     Snowflake
//  - ChannelID  Snowflake
//  - MessageID  Snowflake
//  - Emoji      *Emoji
const EventMessageReactionRemove = event.MessageReactionRemove

// user removed a reaction from a message
type MessageReactionRemove struct {
	UserID    Snowflake `json:"user_id"`
	ChannelID Snowflake `json:"channel_id"`
	MessageID Snowflake `json:"message_id"`
	// PartialEmoji id and name. id might be nil
	PartialEmoji *Emoji          `json:"emoji"`
	Ctx          context.Context `json:"-"`
}

func (obj *MessageReactionRemove) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on MESSAGE_REACTION_REMOVE events
type MessageReactionRemoveCallback = func(session Session, mrr *MessageReactionRemove)

// ---------------------------

// Sent when a user explicitly removes all reactions from a message.
//  Fields:
//  - ChannelID Snowflake
//  - MessageID Snowflake
const EventMessageReactionRemoveAll = event.MessageReactionRemoveAll

// all reactions were explicitly removed from a message
type MessageReactionRemoveAll struct {
	ChannelID Snowflake       `json:"channel_id"`
	MessageID Snowflake       `json:"id"`
	Ctx       context.Context `json:"-"`
}

func (obj *MessageReactionRemoveAll) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on MESSAGE_REACTION_REMOVE_ALL events
type MessageReactionRemoveAllCallback = func(session Session, mrra *MessageReactionRemoveAll)

// ---------------------------

// Sent when a guild's emojis have been updated.
//  Fields:
//  - GuildID Snowflake
//  - Emojis []*Emoji
const EventGuildEmojisUpdate = event.GuildEmojisUpdate

// guild emojis were updated
type GuildEmojisUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Emojis  []*Emoji        `json:"emojis"`
	Ctx     context.Context `json:"-"`
}

func (obj *GuildEmojisUpdate) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on GUILD_EMOJIS_UPDATE events
type GuildEmojisUpdateCallback = func(session Session, geu *GuildEmojisUpdate)

// ---------------------------

// This event can be sent in three different scenarios:
//  1. When a user is initially connecting, to lazily load and backfill information for all unavailable guilds
//     sent in the Ready event.
//	2. When a Guild becomes available again to the client.
// 	3. When the current user joins a new Guild.
const EventGuildCreate = event.GuildCreate

// This event can be sent in three different scenarios:
//  1. When a user is initially connecting, to lazily load and backfill information for all unavailable guilds
//     sent in the Ready event.
//	2. When a Guild becomes available again to the client.
// 	3. When the current user joins a new Guild.
type GuildCreate struct {
	Guild *Guild          `json:"guild"`
	Ctx   context.Context `json:"-"`
}

func (obj *GuildCreate) registerContext(ctx context.Context) { obj.Ctx = ctx }
func (obj *GuildCreate) UnmarshalJSON(data []byte) error {
	obj.Guild = &Guild{}
	return unmarshal(data, obj.Guild)
}

// triggered on GUILD_CREATE events
type GuildCreateCallback = func(session Session, gc *GuildCreate)

// ---------------------------

// Sent when a guild is updated. The inner payload is a guild object.
const EventGuildUpdate = event.GuildUpdate

// guild was updated
type GuildUpdate struct {
	Guild *Guild          `json:"guild"`
	Ctx   context.Context `json:"-"`
}

func (obj *GuildUpdate) registerContext(ctx context.Context) { obj.Ctx = ctx }
func (obj *GuildUpdate) UnmarshalJSON(data []byte) error {
	obj.Guild = &Guild{}
	return unmarshal(data, obj.Guild)
}

// triggered on GUILD_UPDATE events
type GuildUpdateCallback = func(session Session, gu *GuildUpdate)

// ---------------------------

// Sent when a guild becomes unavailable during a guild outage, or when the user leaves or is removed from a guild.
// The inner payload is an unavailable guild object. If the unavailable field is not set, the user was removed
// from the guild.
const EventGuildDelete = event.GuildDelete

// guild became unavailable, or user left/was removed from a guild
type GuildDelete struct {
	UnavailableGuild *GuildUnavailable `json:"guild_unavailable"`
	Ctx              context.Context   `json:"-"`
}

func (obj *GuildDelete) UserWasRemoved() bool {
	return obj.UnavailableGuild.Unavailable == false
}

func (obj *GuildDelete) registerContext(ctx context.Context) { obj.Ctx = ctx }
func (obj *GuildDelete) UnmarshalJSON(data []byte) error {
	obj.UnavailableGuild = &GuildUnavailable{}
	return unmarshal(data, obj.UnavailableGuild)
}

// triggered on GUILD_DELETE events
type GuildDeleteCallback = func(session Session, gd *GuildDelete)

// ---------------------------

// Sent when a user is banned from a guild. The inner payload is a user object, with an extra guild_id key.
const EventGuildBanAdd = event.GuildBanAdd

// user was banned from a guild
type GuildBanAdd struct {
	GuildID Snowflake       `json:"guild_id"`
	User    *User           `json:"user"`
	Ctx     context.Context `json:"-"`
}

func (obj *GuildBanAdd) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on GUILD_BAN_ADD events
type GuildBanAddCallback = func(session Session, gba *GuildBanAdd)

// ---------------------------

// Sent when a user is unbanned from a guild. The inner payload is a user object, with an extra guild_id key.
const EventGuildBanRemove = event.GuildBanRemove

// user was unbanned from a guild
type GuildBanRemove struct {
	GuildID Snowflake       `json:"guild_id"`
	User    *User           `json:"user"`
	Ctx     context.Context `json:"-"`
}

func (obj *GuildBanRemove) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on GUILD_BAN_REMOVE events
type GuildBanRemoveCallback = func(session Session, gbr *GuildBanRemove)

// ---------------------------

// Sent when a guild integration is updated.
//  Fields:
//  - GuildID Snowflake
const EventGuildIntegrationsUpdate = event.GuildIntegrationsUpdate

// guild integration was updated
type GuildIntegrationsUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Ctx     context.Context `json:"-"`
}

func (obj *GuildIntegrationsUpdate) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on GUILD_INTEGRATIONS_UPDATE events
type GuildIntegrationsUpdateCallback = func(session Session, giu *GuildIntegrationsUpdate)

// ---------------------------

// Sent when a new user joins a guild. The inner payload is a guild member object with these extra fields:
//  - GuildID Snowflake
//
//  Fields:
//  - Member *Member
const EventGuildMemberAdd = event.GuildMemberAdd

// new user joined a guild
type GuildMemberAdd struct {
	Member *Member         `json:"member"`
	Ctx    context.Context `json:"-"`
}

func (obj *GuildMemberAdd) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on GUILD_MEMBER_ADD events
type GuildMemberAddCallback = func(session Session, gma *GuildMemberAdd)

// ---------------------------

// Sent when a user is removed from a guild (leave/kick/ban).
//  Fields:
//  - GuildID   Snowflake
//  - User      *User
const EventGuildMemberRemove = event.GuildMemberRemove

// user was removed from a guild
type GuildMemberRemove struct {
	GuildID Snowflake       `json:"guild_id"`
	User    *User           `json:"user"`
	Ctx     context.Context `json:"-"`
}

func (obj *GuildMemberRemove) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on GUILD_MEMBER_REMOVE events
type GuildMemberRemoveCallback = func(session Session, gmr *GuildMemberRemove)

// ---------------------------

// Sent when a guild member is updated.
//  Fields:
//  - GuildID   Snowflake
//  - Roles     []Snowflake
//  - User      *User
//  - Nick      string
const EventGuildMemberUpdate = event.GuildMemberUpdate

// guild member was updated
type GuildMemberUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Roles   []Snowflake     `json:"roles"`
	User    *User           `json:"user"`
	Nick    string          `json:"nick"`
	Ctx     context.Context `json:"-"`
}

func (obj *GuildMemberUpdate) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on GUILD_MEMBER_UPDATE events
type GuildMemberUpdateCallback = func(session Session, gmu *GuildMemberUpdate)

// ---------------------------

// Sent in response to Gateway Request Guild Members.
//  Fields:
//  - GuildID Snowflake
//  - Members []*Member
const EventGuildMembersChunk = event.GuildMembersChunk

// response to Request Guild Members
type GuildMembersChunk struct {
	GuildID Snowflake       `json:"guild_id"`
	Members []*Member       `json:"members"`
	Ctx     context.Context `json:"-"`
}

func (obj *GuildMembersChunk) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on GUILD_MEMBERS_CHUNK events
type GuildMembersChunkCallback = func(session Session, gmc *GuildMembersChunk)

// ---------------------------

// Sent when a guild role is created.
//  Fields:
//  - GuildID   Snowflake
//  - Role      *Role
const EventGuildRoleCreate = event.GuildRoleCreate

// guild role was created
type GuildRoleCreate struct {
	GuildID Snowflake       `json:"guild_id"`
	Role    *Role           `json:"role"`
	Ctx     context.Context `json:"-"`
}

func (obj *GuildRoleCreate) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on GUILD_ROLE_CREATE events
type GuildRoleCreateCallback = func(session Session, grc *GuildRoleCreate)

// ---------------------------

// Sent when a guild role is created.
//  Fields:
//  - GuildID Snowflake
//  - Role    *Role
const EventGuildRoleUpdate = event.GuildRoleUpdate

// guild role was updated
type GuildRoleUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Role    *Role           `json:"role"`
	Ctx     context.Context `json:"-"`
}

func (obj *GuildRoleUpdate) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on GUILD_ROLE_UPDATE events
type GuildRoleUpdateCallback = func(session Session, gru *GuildRoleUpdate)

// ---------------------------

// Sent when a guild role is created.
//  Fields:
//  - GuildID Snowflake
//  - RoleID  Snowflake
const EventGuildRoleDelete = event.GuildRoleDelete

// guild role was deleted
type GuildRoleDelete struct {
	GuildID Snowflake       `json:"guild_id"`
	RoleID  Snowflake       `json:"role_id"`
	Ctx     context.Context `json:"-"`
}

func (obj *GuildRoleDelete) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on GUILD_ROLE_DELETE events
type GuildRoleDeleteCallback = func(session Session, grd *GuildRoleDelete)

// ---------------------------

// A user's presence is their current state on a guild. This event is sent when a user's presence is updated for a guild.
//  Fields:
//  - User    *User
//  - Roles   []Snowflake
//  - Game    *Activity
//  - GuildID Snowflake
//  - Status  string
const EventPresenceUpdate = event.PresenceUpdate

// user's presence was updated in a guild
type PresenceUpdate struct {
	User    *User       `json:"user"`
	RoleIDs []Snowflake `json:"roles"`
	Game    *Activity   `json:"game"`
	GuildID Snowflake   `json:"guild_id"`

	// Status either "idle", "dnd", "online", or "offline"
	// TODO: constants somewhere..
	Status string          `json:"status"`
	Ctx    context.Context `json:"-"`
}

func (obj *PresenceUpdate) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on PRESENCE_UPDATE events
type PresenceUpdateCallback = func(session Session, pu *PresenceUpdate)

// ---------------------------

// Sent when properties about the user change. Inner payload is a user object.
const EventUserUpdate = event.UserUpdate

// properties about a user changed
type UserUpdate struct {
	User *User           `json:"user"`
	Ctx  context.Context `json:"-"`
}

func (obj *UserUpdate) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on USER_UPDATE events
type UserUpdateCallback = func(session Session, uu *UserUpdate)

// ---------------------------

// Sent when someone joins/leaves/moves voice channels. Inner payload is a voice state object.
const EventVoiceStateUpdate = event.VoiceStateUpdate

// someone joined, left, or moved a voice channel
type VoiceStateUpdate struct {
	VoiceState *VoiceState     `json:"voice_state"`
	Ctx        context.Context `json:"-"`
}

func (obj *VoiceStateUpdate) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on VOICE_STATE_UPDATE events
type VoiceStateUpdateCallback = func(session Session, vsu *VoiceStateUpdate)

// ---------------------------

// Sent when a guild's voice server is updated. This is sent when initially connecting to voice, and when the current
// voice instance fails over to a new server.
//  Fields:
//  - Token     string
//  - ChannelID Snowflake
//  - Endpoint  string
const EventVoiceServerUpdate = event.VoiceServerUpdate

// guild's voice server was updated. Sent when a guild's voice server is updated. This is sent when initially
// connecting to voice, and when the current voice instance fails over to a new server.
type VoiceServerUpdate struct {
	Token    string          `json:"token"`
	GuildID  Snowflake       `json:"guild_id"`
	Endpoint string          `json:"endpoint"`
	Ctx      context.Context `json:"-"`
}

func (obj *VoiceServerUpdate) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on VOICE_SERVER_UPDATE events
type VoiceServerUpdateCallback = func(session Session, vsu *VoiceServerUpdate)

// ---------------------------

// Sent when a guild channel's webhook is created, updated, or deleted.
//  Fields:
//  - GuildID   Snowflake
//  - ChannelID Snowflake
const EventWebhooksUpdate = event.WebhooksUpdate

// guild channel webhook was created, update, or deleted
type WebhooksUpdate struct {
	GuildID   Snowflake       `json:"guild_id"`
	ChannelID Snowflake       `json:"channel_id"`
	Ctx       context.Context `json:"-"`
}

func (obj *WebhooksUpdate) registerContext(ctx context.Context) { obj.Ctx = ctx }

// triggered on WEBHOOK_UPDATE events
type WebhooksUpdateCallback = func(session Session, wu *WebhooksUpdate)
