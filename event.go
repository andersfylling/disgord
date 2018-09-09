package disgord

import (
	"context"
	"sync"
)

// KeyAllEvents keys that does not fit within one of the existing files goes here
const KeyAllEvents = "GOD_DAMN_EVERYTHING"

// EventCallback is triggered on every event type
type EventCallback = func(session Session, box interface{})

// ---------------------------

// Hello defines the heartbeat interval
type Hello struct {
	HeartbeatInterval uint            `json:"heartbeat_interval"`
	Trace             []string        `json:"_trace"`
	Ctx               context.Context `json:"-"`
}

// HelloCallback triggered in hello events
type HelloCallback = func(session Session, h *Hello)

// ---------------------------

// KeyReady The ready event is dispatched when a client has completed the
//       initial handshake with the gateway (for new sessions). The ready
//       event can be the largest and most complex event the gateway will
//       send, as it contains all the state required for a client to begin
//       interacting with the rest of the platform.
//       Fields:
//       * V uint8
//       * User *discord.user.User
//       * PrivateChannels []*discord.channel.Private
//       * Guilds []*discord.guild.Unavailable
//       * SessionID string
//       * Trace []string
const KeyReady = "READY"

// Ready	contains the initial state information
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

// ReadyCallback triggered on READY events
type ReadyCallback = func(session Session, r *Ready)

// ---------------------------

// KeyResumed The resumed event is dispatched when a client has sent a resume
//         payload to the gateway (for resuming existing sessions).
//         Fields:
//         * Trace []string
const KeyResumed = "RESUMED"

// Resumed	response to Resume
type Resumed struct {
	Trace []string        `json:"_trace"`
	Ctx   context.Context `json:"-"`
}

// ResumedCallback triggered on RESUME events
type ResumedCallback = func(session Session, r *Resumed)

// ---------------------------

// InvalidSession	failure response to Identify or Resume or invalid active session
type InvalidSession struct {
	Ctx context.Context `json:"-"`
}

// InvalidSessionCallback triggered on INVALID_SESSION events
type InvalidSessionCallback = func(session Session, is *InvalidSession)

// ---------------------------

// KeyChannelCreate Sent when a new channel is created, relevant to the current
//               user. The inner payload is a DM channel or guild channel
//               object.
const KeyChannelCreate = "CHANNEL_CREATE"

// ChannelCreateBox	new channel created
type ChannelCreate struct {
	Channel *Channel        `json:"channel"`
	Ctx     context.Context `json:"-"`
}

// ChannelCreateCallback triggered on CHANNEL_CREATE events
type ChannelCreateCallback = func(session Session, cc *ChannelCreate)

// ---------------------------

// KeyChannelUpdate Sent when a channel is updated. The inner payload is a guild
//               channel object.
const KeyChannelUpdate = "CHANNEL_UPDATE"

// ChannelUpdateBox	channel was updated
type ChannelUpdate struct {
	Channel *Channel        `json:"channel"`
	Ctx     context.Context `json:"-"`
}

// ChannelUpdateCallback triggered on CHANNEL_UPDATE events
type ChannelUpdateCallback = func(session Session, cu *ChannelUpdate)

// ---------------------------

// KeyChannelDelete Sent when a channel relevant to the current user is deleted.
//               The inner payload is a DM or Guild channel object.
const KeyChannelDelete = "CHANNEL_DELETE"

// ChannelDelete	channel was deleted
type ChannelDelete struct {
	Channel *Channel        `json:"channel"`
	Ctx     context.Context `json:"-"`
}

// ChannelDeleteCallback triggered on CHANNEL_DELETE events
type ChannelDeleteCallback = func(session Session, cd *ChannelDelete)

// ---------------------------

// KeyChannelPinsUpdate Sent when a message is pinned or unpinned in a text
//                   channel. This is not sent when a pinned message is
//                   deleted.
//                   Fields:
//                   * ChannelID int64 or discord.Snowflake
//                   * LastPinTimestamp time.Now().UTC().Format(time.RFC3339)
// TODO fix.
const KeyChannelPinsUpdate = "CHANNEL_PINS_UPDATE"

// ChannelPinsUpdate	message was pinned or unpinned
type ChannelPinsUpdate struct {
	// ChannelID snowflake	the id of the channel
	ChannelID Snowflake `json:"channel_id"`

	// LastPinTimestamp	ISO8601 timestamp	the time at which the most recent pinned message was pinned
	LastPinTimestamp Timestamp       `json:"last_pin_timestamp,omitempty"` // ?|
	Ctx              context.Context `json:"-"`
}

// ChannelPinsUpdateCallback triggered on CHANNEL_PINS_UPDATE events
type ChannelPinsUpdateCallback = func(session Session, cpu *ChannelPinsUpdate)

// ---------------------------

// KeyTypingStart Sent when a user starts typing in a channel.
//             Fields: TODO
const KeyTypingStart = "TYPING_START"

// TypingStart user started typing in a channel
type TypingStart struct {
	ChannelID     Snowflake       `json:"channel_id"`
	UserID        Snowflake       `json:"user_id"`
	TimestampUnix int             `json:"timestamp"`
	Ctx           context.Context `json:"-"`
}

// TypingStartCallback triggered on TYPING_START events
type TypingStartCallback = func(session Session, ts *TypingStart)

// ---------------------------

// KeyMessageCreate Sent when a message is created. The inner payload is a
//               message object.
const KeyMessageCreate = "MESSAGE_CREATE"

// MessageCreate	message was created
type MessageCreate struct {
	Message *Message
	Ctx     context.Context `json:"-"`
}

// MessageCreateCallback triggered on MESSAGE_CREATE events
type MessageCreateCallback = func(session Session, mc *MessageCreate)

// ---------------------------

// KeyMessageUpdate Sent when a message is updated. The inner payload is a
//               message object.
//               NOTE! Has _at_least_ the GuildID and ChannelID fields.
const KeyMessageUpdate = "MESSAGE_UPDATE"

// MessageUpdate	message was edited
type MessageUpdate struct {
	Message *Message
	Ctx     context.Context `json:"-"`
}

// MessageUpdateCallback triggered on MESSAGE_UPDATE events
type MessageUpdateCallback = func(session Session, mu *MessageUpdate)

// ---------------------------

// KeyMessageDelete Sent when a message is deleted.
//               Fields:
//               * ID        int64 or discord.Snowflake
//               * ChannelID int64 or discord.Snowflake
const KeyMessageDelete = "MESSAGE_DELETE"

// MessageDelete	message was deleted
type MessageDelete struct {
	MessageID Snowflake       `json:"id"`
	ChannelID Snowflake       `json:"channel_id"`
	Ctx       context.Context `json:"-"`
}

// MessageDeleteCallback triggered on MESSAGE_DELETE events
type MessageDeleteCallback = func(session Session, md *MessageDelete)

// ---------------------------

// KeyMessageDeleteBulk Sent when multiple messages are deleted at once.
//                   Fields:
//                   * IDs       []int64 or []discord.Snowflake
//                   * ChannelID int64 or discord.Snowflake
const KeyMessageDeleteBulk = "MESSAGE_DELETE_BULK"

// MessageDeleteBulk	multiple messages were deleted at once
type MessageDeleteBulk struct {
	MessageIDs []Snowflake     `json:"ids"`
	ChannelID  Snowflake       `json:"channel_id"`
	Ctx        context.Context `json:"-"`
}

// MessageDeleteBulkCallback triggered on MESSAGE_DELETE_BULK events
type MessageDeleteBulkCallback = func(session Session, mdb *MessageDeleteBulk)

// ---------------------------

// KeyMessageReactionAdd Sent when a user adds a reaction to a message.
//                    Fields:
//                    * UserID     int64 or discord.Snowflake
//                    * ChannelID  int64 or discord.Snowflake
//                    * MessageID  int64 or discord.Snowflake
//                    * Emoji      *discord.emoji.Emoji
const KeyMessageReactionAdd = "MESSAGE_REACTION_ADD"

// MessageReactionAdd	user reacted to a message
type MessageReactionAdd struct {
	UserID    Snowflake `json:"user_id"`
	ChannelID Snowflake `json:"channel_id"`
	MessageID Snowflake `json:"message_id"`
	// PartialEmoji id and name. id might be nil
	PartialEmoji *Emoji          `json:"emoji"`
	Ctx          context.Context `json:"-"`
}

// MessageReactionAddCallback triggered on MESSAGE_REACTION_ADD events
type MessageReactionAddCallback = func(session Session, mra *MessageReactionAdd)

// ---------------------------

// KeyMessageReactionRemove Sent when a user removes a reaction from a message.
//                       Fields:
//                       * UserID     int64 or discord.Snowflake
//                       * ChannelID  int64 or discord.Snowflake
//                       * MessageID  int64 or discord.Snowflake
//                       * Emoji      *discord.emoji.Emoji
const KeyMessageReactionRemove = "MESSAGE_REACTION_REMOVE"

// MessageReactionRemove	user removed a reaction from a message
type MessageReactionRemove struct {
	UserID    Snowflake `json:"user_id"`
	ChannelID Snowflake `json:"channel_id"`
	MessageID Snowflake `json:"message_id"`
	// PartialEmoji id and name. id might be nil
	PartialEmoji *Emoji          `json:"emoji"`
	Ctx          context.Context `json:"-"`
}

// MessageReactionRemoveCallback triggered on MESSAGE_REACTION_REMOVE events
type MessageReactionRemoveCallback = func(session Session, mrr *MessageReactionRemove)

// ---------------------------

// KeyMessageReactionRemoveAll Sent when a user explicitly removes all reactions
//                          from a message.
//                          Fields:
//                          * ChannelID  int64 or discord.Snowflake
//                          * MessageID  int64 or discord.Snowflake
const KeyMessageReactionRemoveAll = "MESSAGE_REACTION_REMOVE_ALL"

// MessageReactionRemoveAll	all reactions were explicitly removed from a message
type MessageReactionRemoveAll struct {
	ChannelID Snowflake       `json:"channel_id"`
	MessageID Snowflake       `json:"id"`
	Ctx       context.Context `json:"-"`
}

// MessageReactionRemoveAllCallback triggered on MESSAGE_REACTION_REMOVE_ALL events
type MessageReactionRemoveAllCallback = func(session Session, mrra *MessageReactionRemoveAll)

// ---------------------------

// KeyGuildEmojisUpdate Sent when a guild's emojis have been updated.
//                   Fields:
//                   * GuildID int64 or discord.Snowflake
const KeyGuildEmojisUpdate = "GUILD_EMOJI_UPDATE"

// GuildEmojisUpdate	guild emojis were updated
type GuildEmojisUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Emojis  []*Emoji        `json:"emojis"`
	Ctx     context.Context `json:"-"`
}

// GuildEmojisUpdateCallback triggered on GUILD_EMOJIS_UPDATE events
type GuildEmojisUpdateCallback = func(session Session, geu *GuildEmojisUpdate)

// ---------------------------

// KeyGuildCreate This event can be sent in three different scenarios:
//             1. When a user is initially connecting, to lazily load and
//                backfill information for all unavailable guilds sent in the
//                Ready event.
//             2. When a Guild becomes available again to the client.
//             3. When the current user joins a new Guild.
//             The inner payload is a guild object, with all the extra fields
//             specified.
const KeyGuildCreate = "GUILD_CREATE"

// GuildCreate	This event can be sent in three different scenarios:
//								1. When a user is initially connecting, to lazily load and backfill information for
//									 all unavailable guilds sent in the Ready event.
//								2. When a Guild becomes available again to the client.
// 								3. When the current user joins a new Guild.
type GuildCreate struct {
	Guild *Guild          `json:"guild"`
	Ctx   context.Context `json:"-"`
}

// GuildCreateCallback triggered on GUILD_CREATE events
type GuildCreateCallback = func(session Session, gc *GuildCreate)

// ---------------------------

// KeyGuildUpdate Sent when a guild is updated. The inner payload is a guild
//             object.
const KeyGuildUpdate = "GUILD_UPDATE"

// GuildUpdate	guild was updated
type GuildUpdate struct {
	Guild *Guild          `json:"guild"`
	Ctx   context.Context `json:"-"`
}

// GuildUpdateCallback triggered on GUILD_UPDATE events
type GuildUpdateCallback = func(session Session, gu *GuildUpdate)

// ---------------------------

// KeyGuildDelete Sent when a guild becomes unavailable during a guild outage,
//             or when the user leaves or is removed from a guild. The inner
//             payload is an unavailable guild object. If the unavailable
//             field is not set, the user was removed from the guild.
const KeyGuildDelete = "GUILD_DELETE"

// GuildDelete	guild became unavailable, or user left/was removed from a guild
type GuildDelete struct {
	UnavailableGuild *GuildUnavailable `json:"guild_unavailable"`
	Ctx              context.Context   `json:"-"`
}

// GuildDeleteCallback triggered on GUILD_DELETE events
type GuildDeleteCallback = func(session Session, gd *GuildDelete)

// ---------------------------

// KeyGuildBanAdd Sent when a user is banned from a guild. The inner payload is
//             a user object, with an extra guild_id key.
const KeyGuildBanAdd = "GUILD_BAN_ADD"

// GuildBanAdd	user was banned from a guild
type GuildBanAdd struct {
	User *User           `json:"user"`
	Ctx  context.Context `json:"-"`
}

// GuildBanAddCallback triggered on GUILD_BAN_ADD events
type GuildBanAddCallback = func(session Session, gba *GuildBanAdd)

// ---------------------------

// KeyGuildBanRemove Sent when a user is unbanned from a guild. The inner
//                payload is a user object, with an extra guild_id key.
const KeyGuildBanRemove = "GUILD_BAN_REMOVE"

// GuildBanRemove	user was unbanned from a guild
type GuildBanRemove struct {
	User *User           `json:"user"`
	Ctx  context.Context `json:"-"`
}

// GuildBanRemoveCallback triggered on GUILD_BAN_REMOVE events
type GuildBanRemoveCallback = func(session Session, gbr *GuildBanRemove)

// ---------------------------

// KeyGuildIntegrationsUpdate Sent when a guild integration is updated.
//                        Fields:
//                        * GuildID int64 or discord.Snowflake
//                        * Emojis []*discord.emoji.Emoji
const KeyGuildIntegrationsUpdate = "GUILD_INTEGRATIONS_UPDATE"

// GuildIntegrationsUpdate	guild integration was updated
type GuildIntegrationsUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Ctx     context.Context `json:"-"`
}

// GuildIntegrationsUpdateCallback triggered on GUILD_INTEGRATIONS_UPDATE events
type GuildIntegrationsUpdateCallback = func(session Session, giu *GuildIntegrationsUpdate)

// ---------------------------

// KeyGuildMemberAdd Sent when a new user joins a guild. The inner payload is a
//                guild member object with these extra fields:
//                * GuildID int64 or discord.Snowflake
const KeyGuildMemberAdd = "GUILD_MEMBER_ADD"

// GuildMemberAdd	new user joined a guild
type GuildMemberAdd struct {
	Member *Member         `json:"member"`
	Ctx    context.Context `json:"-"`
}

// GuildMemberAddCallback triggered on GUILD_MEMBER_ADD events
type GuildMemberAddCallback = func(session Session, gma *GuildMemberAdd)

// ---------------------------

// KeyGuildMemberRemove Sent when a user is removed from a guild
//                   (leave/kick/ban).
//                   Fields:
//                   * GuildID int64 or discord.Snowflake
//                   * User *discord.user.User
const KeyGuildMemberRemove = "GUILD_MEMBER_REMOVE"

// GuildMemberRemove	user was removed from a guild
type GuildMemberRemove struct {
	GuildID Snowflake       `json:"guild_id"`
	User    *User           `json:"user"`
	Ctx     context.Context `json:"-"`
}

// GuildMemberRemoveCallback triggered on GUILD_MEMBER_REMOVE events
type GuildMemberRemoveCallback = func(session Session, gmr *GuildMemberRemove)

// ---------------------------

// KeyGuildMemberUpdate Sent when a guild member is updated.
//                   Fields:
//                   * GuildID int64 or discord.Snowflake
//                   * Roles []int64 or []discord.Snowflake
//                   * User *discord.user.User
//                   * Nick string
const KeyGuildMemberUpdate = "GUILD_MEMBER_UPDATE"

// GuildMemberUpdate	guild member was updated
type GuildMemberUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Roles   []*Role         `json:"roles"`
	User    *User           `json:"user"`
	Nick    string          `json:"nick"`
	Ctx     context.Context `json:"-"`
}

// GuildMemberUpdateCallback triggered on GUILD_MEMBER_UPDATE events
type GuildMemberUpdateCallback = func(session Session, gmu *GuildMemberUpdate)

// ---------------------------

// KeyGuildMemberChunk Sent in response to Gateway Request Guild Members.
//                  Fields:
//                  * GuildID int64 or discord.Snowflake
//                  * Members []*discord.member.Member
const KeyGuildMembersChunk = "GUILD_MEMBER_CHUNK"

// GuildMembersChunk	response to Request Guild Members
type GuildMembersChunk struct {
	GuildID Snowflake       `json:"guild_id"`
	Members []*Member       `json:"members"`
	Ctx     context.Context `json:"-"`
}

// GuildMembersChunkCallback triggered on GUILD_MEMBERS_CHUNK events
type GuildMembersChunkCallback = func(session Session, gmc *GuildMembersChunk)

// ---------------------------

// KeyGuildRoleCreate Sent when a guild role is created.
//                 Fields:
//                 * GuildID int64 or discord.Snowflake
//                 * Role *discord.role.Role
const KeyGuildRoleCreate = "GUILD_ROLE_CREATE"

// GuildRoleCreate	guild role was created
type GuildRoleCreate struct {
	GuildID Snowflake       `json:"guild_id"`
	Role    *Role           `json:"role"`
	Ctx     context.Context `json:"-"`
}

// GuildRoleCreateCallback triggered on GUILD_ROLE_CREATE events
type GuildRoleCreateCallback = func(session Session, grc *GuildRoleCreate)

// ---------------------------

// KeyGuildRoleUpdate Sent when a guild role is created.
//                 Fields:
//                 * GuildID int64 or discord.Snowflake
//                 * Role    *discord.role.Role
const KeyGuildRoleUpdate = "GUILD_ROLE_UPDATE"

// GuildRoleUpdate	guild role was updated
type GuildRoleUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Role    *Role           `json:"role"`
	Ctx     context.Context `json:"-"`
}

// GuildRoleUpdateCallback triggered on GUILD_ROLE_UPDATE events
type GuildRoleUpdateCallback = func(session Session, gru *GuildRoleUpdate)

// ---------------------------

// KeyGuildRoleDelete Sent when a guild role is created.
//                 Fields:
//                 * GuildID int64 or discord.Snowflake
//                 * RoleID  int64 or discord.Snowflake
const KeyGuildRoleDelete = "GUILD_ROLE_DELETE"

// GuildRoleDelete	guild role was deleted
type GuildRoleDelete struct {
	GuildID Snowflake       `json:"guild_id"`
	RoleID  Snowflake       `json:"role_id"`
	Ctx     context.Context `json:"-"`
}

// GuildRoleDeleteCallback triggered on GUILD_ROLE_DELETE events
type GuildRoleDeleteCallback = func(session Session, grd *GuildRoleDelete)

// ---------------------------

// KeyPresenceUpdate A user's presence is their current state on a guild.
//                This event is sent when a user's presence is updated
//                for a guild.
//                Fields:
//                User    *discord.user.User
//                Roles   []*discord.role.Role
//                Game    *discord.game.Game
//                GuildID int64 or discord.Snowflake
//                Status  *string or *discord.presence.Status
const KeyPresenceUpdate = "PRESENCE_UPDATE"

// PresenceUpdate	user's presence was updated in a guild
type PresenceUpdate struct {
	User    *User         `json:"user"`
	RoleIDs []Snowflake   `json:"roles"`
	Game    *UserActivity `json:"game"`
	GuildID Snowflake     `json:"guild_id"`

	// Status either "idle", "dnd", "online", or "offline"
	// TODO: constants somewhere..
	Status string          `json:"status"`
	Ctx    context.Context `json:"-"`
}

// PresenceUpdateCallback triggered on PRESENCE_UPDATE events
type PresenceUpdateCallback = func(session Session, pu *PresenceUpdate)

// ---------------------------

// KeyUserUpdate Sent when properties about the user change. Inner payload is a
//            user object.
const KeyUserUpdate = "USER_UPDATE"

// UserUpdate	properties about a user changed
type UserUpdate struct {
	User *User           `json:"user"`
	Ctx  context.Context `json:"-"`
}

// UserUpdateCallback triggerd on USER_UPDATE events
type UserUpdateCallback = func(session Session, uu *UserUpdate)

// ---------------------------

// KeyVoiceStateUpdate Sent when someone joins/leaves/moves voice channels.
//                  Inner payload is a voice state object.
const KeyVoiceStateUpdate = "VOICE_STATE_UPDATE"

// VoiceStateUpdate	someone joined, left, or moved a voice channel
type VoiceStateUpdate struct {
	VoiceState *VoiceState     `json:"voice_state"`
	Ctx        context.Context `json:"-"`
}

// VoiceStateUpdateCallback triggered on VOICE_STATE_UPDATE events
type VoiceStateUpdateCallback = func(session Session, vsu *VoiceStateUpdate)

// ---------------------------

// KeyVoiceServerUpdate Sent when a guild's voice server is updated. This is
//                   sent when initially connecting to voice, and when the
//                   current voice instance fails over to a new server.
//                   Fields:
//                   * Token     string or discord.Token
//                   * ChannelID int64 or discord.Snowflake
//                   * Endpoint  string or discord.Endpoint
const KeyVoiceServerUpdate = "VOICE_SERVER_UPDATE"

// VoiceServerUpdate	guild's voice server was updated
// Sent when a guild's voice server is updated.
// This is sent when initially connecting to voice,
// and when the current voice instance fails over to a new server.
type VoiceServerUpdate struct {
	Token    string          `json:"token"`
	GuildID  Snowflake       `json:"guild_id"`
	Endpoint string          `json:"endpoint"`
	Ctx      context.Context `json:"-"`
}

// VoiceServerUpdateCallback triggered on VOICE_SERVER_UPDATE events
type VoiceServerUpdateCallback = func(session Session, vsu *VoiceServerUpdate)

// ---------------------------

// KeyWebhooksUpdate Sent when a guild channel's webhook is created, updated, or
//                deleted.
//                Fields:
//                * GuildID   int64 or discord.Snowflake
//                * ChannelID int64 or discord.Snowflake
const KeyWebhooksUpdate = "WEBHOOK_UPDATE"

// WebhooksUpdate guild channel webhook was created, update, or deleted
type WebhooksUpdate struct {
	GuildID   Snowflake       `json:"guild_id"`
	ChannelID Snowflake       `json:"channel_id"`
	Ctx       context.Context `json:"-"`
}

// WebhooksUpdateCallback triggered on WEBHOOK_UPDATE events
type WebhooksUpdateCallback = func(session Session, wu *WebhooksUpdate)
