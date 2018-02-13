package event

type Type string

// Event types
const (

	// Gateway events

	// Hello defines the heartbeat interval
	Hello Type = "HELLO"
	// Ready The ready event is dispatched when a client has completed the
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
	Ready = "READY"
	// Resumed The resumed event is dispatched when a client has sent a resume
	//         payload to the gateway (for resuming existing sessions).
	//         Fields:
	//         * Trace []string
	Resumed = "RESUMED"
	// InvalidSession failure response to Identify or Resume or invalid active session
	InvalidSession = "INVALID_SESSION"

	// Channel events

	// ChannelCreate Sent when a new channel is created, relevant to the current
	//               user. The inner payload is a DM channel or guild channel
	//               object.
	ChannelCreate = "CHANNEL_CREATE"
	// ChannelUpdate Sent when a channel is updated. The inner payload is a guild
	//               channel object.
	ChannelUpdate = "CHANNEL_UPDATE"
	// ChannelDelete Sent when a channel relevant to the current user is deleted.
	//               The inner payload is a DM or Guild channel object.
	ChannelDelete = "CHANNEL_DELETE"
	// ChannelPinsUpdate Sent when a message is pinned or unpinned in a text
	//                   channel. This is not sent when a pinned message is
	//                   deleted.
	//                   Fields:
	//                   * ChannelID int64 or discord.Snowflake
	//                   * LastPinTimestamp time.Now().UTC().Format(time.RFC3339)
	// TODO fix.
	ChannelPinsUpdate = "CHANNEL_PINS_UPDATE"

	// GUILD events

	// GuildCreate This event can be sent in three different scenarios:
	//             1. When a user is initially connecting, to lazily load and
	//                backfill information for all unavailable guilds sent in the
	//                Ready event.
	//             2. When a Guild becomes available again to the client.
	//             3. When the current user joins a new Guild.
	//             The inner payload is a guild object, with all the extra fields
	//             specified.
	GuildCreate = "GUILD_CREATE"
	// GuildUpdate Sent when a guild is updated. The inner payload is a guild
	//             object.
	GuildUpdate = "GUILD_UPDATE"
	// GuildDelete Sent when a guild becomes unavailable during a guild outage,
	//             or when the user leaves or is removed from a guild. The inner
	//             payload is an unavailable guild object. If the unavailable
	//             field is not set, the user was removed from the guild.
	GuildDelete = "GUILD_DELETE"
	// GuildBanAdd Sent when a user is banned from a guild. The inner payload is
	//             a user object, with an extra guild_id key.
	GuildBanAdd = "GUILD_BAN_ADD"
	// GuildBanRemove Sent when a user is unbanned from a guild. The inner
	//                payload is a user object, with an extra guild_id key.
	GuildBanRemove = "GUILD_BAN_REMOVE"
	// GuildEmojisUpdate Sent when a guild's emojis have been updated.
	//                   Fields:
	//                   * GuildID int64 or discord.Snowflake
	GuildEmojisUpdate = "GUILD_EMOJI_UPDATE"
	//GuildIntegrationsUpdate Sent when a guild integration is updated.
	//                        Fields:
	//                        * GuildID int64 or discord.Snowflake
	//                        * Emojis []*discord.emoji.Emoji
	GuildIntegrationsUpdate = "GUILD_INTEGRATIONS_UPDATE"
	// GuildMemberAdd Sent when a new user joins a guild. The inner payload is a
	//                guild member object with these extra fields:
	//                * GuildID int64 or discord.Snowflake
	GuildMemberAdd = "GUILD_MEMBER_ADD"
	// GuildMemberRemove Sent when a user is removed from a guild
	//                   (leave/kick/ban).
	//                   Fields:
	//                   * GuildID int64 or discord.Snowflake
	//                   * User *discord.user.User
	GuildMemberRemove = "GUILD_MEMBER_REMOVE"
	// GuildMemberUpdate Sent when a guild member is updated.
	//                   Fields:
	//                   * GuildID int64 or discord.Snowflake
	//                   * Roles []int64 or []discord.Snowflake
	//                   * User *discord.user.User
	//                   * Nick string
	GuildMemberUpdate = "GUILD_MEMBER_UPDATE"
	// GuildMemberChunk Sent in response to Gateway Request Guild Members.
	//                  Fields:
	//                  * GuildID int64 or discord.Snowflake
	//                  * Members []*discord.member.Member
	GuildMemberChunk = "GUILD_MEMBER_CHUNK"
	// GuildRoleCreate Sent when a guild role is created.
	//                 Fields:
	//                 * GuildID int64 or discord.Snowflake
	//                 * Role *discord.role.Role
	GuildRoleCreate = "GUILD_ROLE_CREATE"
	// GuildRoleUpdate Sent when a guild role is created.
	//                 Fields:
	//                 * GuildID int64 or discord.Snowflake
	//                 * Role    *discord.role.Role
	GuildRoleUpdate = "GUILD_ROLE_UPDATE"
	// GuildRoleDelete Sent when a guild role is created.
	//                 Fields:
	//                 * GuildID int64 or discord.Snowflake
	//                 * RoleID  int64 or discord.Snowflake
	GuildRoleDelete = "GUILD_ROLE_DELETE"
	// MessageCreate Sent when a message is created. The inner payload is a
	//               message object.
	MessageCreate = "MESSAGE_CREATE"
	// MessageUpdate Sent when a message is updated. The inner payload is a
	//               message object.
	//               NOTE! Has _at_least_ the GuildID and ChannelID fields.
	MessageUpdate = "MESSAGE_UPDATE"
	// MessageDelete Sent when a message is deleted.
	//               Fields:
	//               * ID        int64 or discord.Snowflake
	//               * ChannelID int64 or discord.Snowflake
	MessageDelete = "MESSAGE_DELETE"
	// MessageDeleteBulk Sent when multiple messages are deleted at once.
	//                   Fields:
	//                   * IDs       []int64 or []discord.Snowflake
	//                   * ChannelID int64 or discord.Snowflake
	MessageDeleteBulk = "MESSAGE_DELETE_BULK"
	// MessageReactionAdd Sent when a user adds a reaction to a message.
	//                    Fields:
	//                    * UserID     int64 or discord.Snowflake
	//                    * ChannelID  int64 or discord.Snowflake
	//                    * MessageID  int64 or discord.Snowflake
	//                    * Emoji      *discord.emoji.Emoji
	MessageReactionAdd = "MESSAGE_REACTION_ADD"
	// MessageReactionRemove Sent when a user removes a reaction from a message.
	//                       Fields:
	//                       * UserID     int64 or discord.Snowflake
	//                       * ChannelID  int64 or discord.Snowflake
	//                       * MessageID  int64 or discord.Snowflake
	//                       * Emoji      *discord.emoji.Emoji
	MessageReactionRemove = "MESSAGE_REACTION_REMOVE"
	// MessageReactionRemoveAll Sent when a user explicitly removes all reactions
	//                          from a message.
	//                          Fields:
	//                          * ChannelID  int64 or discord.Snowflake
	//                          * MessageID  int64 or discord.Snowflake
	MessageReactionRemoveAll = "MESSAGE_REACTION_REMOVE_ALL"
	// PresenceUpdate A user's presence is their current state on a guild.
	//                This event is sent when a user's presence is updated
	//                for a guild.
	//                Fields:
	//                User    *discord.user.User
	//                Roles   []*discord.role.Role
	//                Game    *discord.game.Game
	//                GuildID int64 or discord.Snowflake
	//                Status  *string or *discord.presence.Status
	PresenceUpdate = "PRESENCE_UPDATE"
	// TypingStart Sent when a user starts typing in a channel.
	//             Fields: TODO
	TypingStart = "TYPING_START"
	// UserUpdate Sent when properties about the user change. Inner payload is a
	//            user object.
	UserUpdate = "USER_UPDATE"
	// VoiceStateUpdate Sent when someone joins/leaves/moves voice channels.
	//                  Inner payload is a voice state object.
	VoiceStateUpdate = "VOICE_STATE_UPDATE"
	// VoiceServerUpdate Sent when a guild's voice server is updated. This is
	//                   sent when initially connecting to voice, and when the
	//                   current voice instance fails over to a new server.
	//                   Fields:
	//                   * Token     string or discord.Token
	//                   * ChannelID int64 or discord.Snowflake
	//                   * Endpoint  string or discord.Endpoint
	VoiceServerUpdate = "VOICE_SERVER_UPDATE"
	// WebhooksUpdate Sent when a guild channel's webhook is created, updated, or
	//                deleted.
	//                Fields:
	//                * GuildID   int64 or discord.Snowflake
	//                * ChannelID int64 or discord.Snowflake
	WebhooksUpdate = "WEBHOOK_UPDATE"
)
