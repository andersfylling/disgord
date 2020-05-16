// Package event is a universal discord package that holds all the event types one can receive (currently only bot events).
package event

// Ready The ready event is dispatched when a client has completed the initial handshake with the gateway (for new sessions).
//// The ready event can be the largest and most complex event the gateway will send, as it contains all the state
//// required for a client to begin interacting with the rest of the platform.
////  Fields:
////  - V int
////  - User *User
////  - PrivateChannels []*Channel
////  - Guilds []*GuildUnavailable
////  - SessionID string
////  - Trace []string
const Ready = "READY"

// Resumed The resumed event is dispatched when a client has sent a resume payload to the gateway
// (for resuming existing sessions).
//  Fields:
//  - Trace []string
const Resumed = "RESUMED"

// ChannelCreate Sent when a new channel is created, relevant to the current user. The inner payload is a DM channel or
// guild channel object.
const ChannelCreate = "CHANNEL_CREATE"

// ChannelUpdate Sent when a channel is updated. The inner payload is a guild channel object.
const ChannelUpdate = "CHANNEL_UPDATE"

// ChannelDelete Sent when a channel relevant to the current user is deleted. The inner payload is a DM or Guild channel object.
const ChannelDelete = "CHANNEL_DELETE"

// ChannelPinsUpdate Sent when a message is pinned or unpinned in a text channel. This is not sent when a pinned message is deleted.
//  Fields:
//  - ChannelID int64 or Snowflake
//  - LastPinTimestamp time.Now().UTC().Format(time.RFC3339)
// TODO fix.
const ChannelPinsUpdate = "CHANNEL_PINS_UPDATE"

// TypingStart Sent when a user starts typing in a channel.
//  Fields:
//  - ChannelID     Snowflake
//  - UserID        Snowflake
//  - TimestampUnix int
const TypingStart = "TYPING_START"

// InviteDelete Sent when an invite is deleted.
const InviteDelete = "INVITE_DELETE"

// MessageCreate Sent when a message is created. The inner payload is a message object.
const MessageCreate = "MESSAGE_CREATE"

// MessageUpdate Sent when a message is updated. The inner payload is a message object.
//
// NOTE! Has _at_least_ the GuildID and ChannelID fields.
const MessageUpdate = "MESSAGE_UPDATE"

// MessageDelete Sent when a message is deleted.
//  Fields:
//  - ID        Snowflake
//  - ChannelID Snowflake
const MessageDelete = "MESSAGE_DELETE"

// MessageDeleteBulk Sent when multiple messages are deleted at once.
//  Fields:
//  - IDs       []Snowflake
//  - ChannelID Snowflake
const MessageDeleteBulk = "MESSAGE_DELETE_BULK"

// MessageReactionAdd Sent when a user adds a reaction to a message.
//  Fields:
//  - UserID     Snowflake
//  - ChannelID  Snowflake
//  - MessageID  Snowflake
//  - Emoji      *Emoji
const MessageReactionAdd = "MESSAGE_REACTION_ADD"

// MessageReactionRemove Sent when a user removes a reaction from a message.
//  Fields:
//  - UserID     Snowflake
//  - ChannelID  Snowflake
//  - MessageID  Snowflake
//  - Emoji      *Emoji
const MessageReactionRemove = "MESSAGE_REACTION_REMOVE"

// MessageReactionRemoveAll Sent when a user explicitly removes all reactions from a message.
//  Fields:
//  - ChannelID Snowflake
//  - MessageID Snowflake
const MessageReactionRemoveAll = "MESSAGE_REACTION_REMOVE_ALL"

// GuildEmojisUpdate Sent when a guild's emojis have been updated.
//  Fields:
//  - GuildID Snowflake
//  - Emojis []*Emoji
const GuildEmojisUpdate = "GUILD_EMOJIS_UPDATE"

// GuildCreate This event can be sent in three different scenarios:
//  1. When a user is initially connecting, to lazily load and backfill information for all unavailable guilds
//     sent in the Ready event.
//	2. When a Guild becomes available again to the client.
// 	3. When the current user joins a new Guild.
const GuildCreate = "GUILD_CREATE"

// GuildUpdate Sent when a guild is updated. The inner payload is a guild object.
const GuildUpdate = "GUILD_UPDATE"

// GuildDelete Sent when a guild becomes unavailable during a guild outage, or when the user leaves or is removed from a guild.
// The inner payload is an unavailable guild object. If the unavailable field is not set, the user was removed
// from the guild.
const GuildDelete = "GUILD_DELETE"

// GuildBanAdd Sent when a user is banned from a guild. The inner payload is a user object, with an extra guild_id key.
const GuildBanAdd = "GUILD_BAN_ADD"

// GuildBanRemove Sent when a user is unbanned from a guild. The inner payload is a user object, with an extra guild_id key.
const GuildBanRemove = "GUILD_BAN_REMOVE"

// GuildIntegrationsUpdate Sent when a guild integration is updated.
//  Fields:
//  - GuildID Snowflake
const GuildIntegrationsUpdate = "GUILD_INTEGRATIONS_UPDATE"

// GuildMemberAdd Sent when a new user joins a guild. The inner payload is a guild member object with these extra fields:
//  - GuildID Snowflake
//
//  Fields:
//  - Member *Member
const GuildMemberAdd = "GUILD_MEMBER_ADD"

// GuildMemberRemove Sent when a user is removed from a guild (leave/kick/ban).
//  Fields:
//  - GuildID   Snowflake
//  - User      *User
const GuildMemberRemove = "GUILD_MEMBER_REMOVE"

// GuildMemberUpdate Sent when a guild member is updated.
//  Fields:
//  - GuildID   Snowflake
//  - Roles     []Snowflake
//  - User      *User
//  - Nick      string
const GuildMemberUpdate = "GUILD_MEMBER_UPDATE"

// GuildMembersChunk Sent in response to Gateway Request Guild Members.
//  Fields:
//  - GuildID Snowflake
//  - Members []*Member
const GuildMembersChunk = "GUILD_MEMBERS_CHUNK"

// GuildRoleCreate Sent when a guild role is created.
//  Fields:
//  - GuildID   Snowflake
//  - Role      *Role
const GuildRoleCreate = "GUILD_ROLE_CREATE"

// GuildRoleUpdate Sent when a guild role is created.
//  Fields:
//  - GuildID Snowflake
//  - Role    *Role
const GuildRoleUpdate = "GUILD_ROLE_UPDATE"

// GuildRoleDelete Sent when a guild role is created.
//  Fields:
//  - GuildID Snowflake
//  - RoleID  Snowflake
const GuildRoleDelete = "GUILD_ROLE_DELETE"

// PresenceUpdate A user's presence is their current state on a guild. This event is sent when a user's presence is updated for a guild.
//  Fields:
//  - User    *User
//  - Roles   []Snowflake
//  - Game    *Activity
//  - GuildID Snowflake
//  - Status  string
const PresenceUpdate = "PRESENCE_UPDATE"

// UserUpdate Sent when properties about the user change. Inner payload is a user object.
const UserUpdate = "USER_UPDATE"

// VoiceStateUpdate Sent when someone joins/leaves/moves voice channels. Inner payload is a voice state object.
const VoiceStateUpdate = "VOICE_STATE_UPDATE"

// VoiceServerUpdate Sent when a guild's voice server is updated. This is sent when initially connecting to voice, and when the current
// voice instance fails over to a new server.
//  Fields:
//  - Token     string
//  - ChannelID Snowflake
//  - Endpoint  string
const VoiceServerUpdate = "VOICE_SERVER_UPDATE"

// WebhooksUpdate Sent when a guild channel's WebHook is created, updated, or deleted.
//  Fields:
//  - GuildID   Snowflake
//  - ChannelID Snowflake
const WebhooksUpdate = "WEBHOOKS_UPDATE"

// InviteCreate Sent when a guild's invite is created.
//  Fields:
//  - Code String
//  - GuildID   Snowflake
//  - ChannelID Snowflake
//  - Inviter *User
//  - Inviter *User
//  - Target *User
//  - TargetType int
//  - CreatedAt Time
//  - MaxAge int
//  - MaxUses int
//  - Temporary bool
//  - Uses int
//  - Revoked bool
//  - Unique bool
//  - ApproximatePresenceCount int
//  - ApproximateMemberCount int
const InviteCreate = "INVITE_CREATE"
