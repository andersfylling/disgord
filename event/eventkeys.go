package event

// KeyType Deprecated
//type KeyType string

// The different discord Event Keys
const (
	AllEventsKey string = "GOD_DAMN_EVERYTHING"

	// Gateway events

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
	ReadyKey = "READY"
	// Resumed The resumed event is dispatched when a client has sent a resume
	//         payload to the gateway (for resuming existing sessions).
	//         Fields:
	//         * Trace []string
	ResumedKey = "RESUMED"

	// Channel events

	// ChannelCreate Sent when a new channel is created, relevant to the current
	//               user. The inner payload is a DM channel or guild channel
	//               object.
	ChannelCreateKey = "CHANNEL_CREATE"
	// ChannelUpdate Sent when a channel is updated. The inner payload is a guild
	//               channel object.
	ChannelUpdateKey = "CHANNEL_UPDATE"
	// ChannelDelete Sent when a channel relevant to the current user is deleted.
	//               The inner payload is a DM or Guild channel object.
	ChannelDeleteKey = "CHANNEL_DELETE"
	// ChannelPinsUpdate Sent when a message is pinned or unpinned in a text
	//                   channel. This is not sent when a pinned message is
	//                   deleted.
	//                   Fields:
	//                   * ChannelID int64 or discord.Snowflake
	//                   * LastPinTimestamp time.Now().UTC().Format(time.RFC3339)
	// TODO fix.
	ChannelPinsUpdateKey = "CHANNEL_PINS_UPDATE"

	// GUILD events

	// GuildCreate This event can be sent in three different scenarios:
	//             1. When a user is initially connecting, to lazily load and
	//                backfill information for all unavailable guilds sent in the
	//                Ready event.
	//             2. When a Guild becomes available again to the client.
	//             3. When the current user joins a new Guild.
	//             The inner payload is a guild object, with all the extra fields
	//             specified.
	GuildCreateKey = "GUILD_CREATE"
	// GuildUpdate Sent when a guild is updated. The inner payload is a guild
	//             object.
	GuildUpdateKey = "GUILD_UPDATE"
	// GuildDelete Sent when a guild becomes unavailable during a guild outage,
	//             or when the user leaves or is removed from a guild. The inner
	//             payload is an unavailable guild object. If the unavailable
	//             field is not set, the user was removed from the guild.
	GuildDeleteKey = "GUILD_DELETE"
	// GuildBanAdd Sent when a user is banned from a guild. The inner payload is
	//             a user object, with an extra guild_id key.
	GuildBanAddKey = "GUILD_BAN_ADD"
	// GuildBanRemove Sent when a user is unbanned from a guild. The inner
	//                payload is a user object, with an extra guild_id key.
	GuildBanRemoveKey = "GUILD_BAN_REMOVE"
	// GuildEmojisUpdate Sent when a guild's emojis have been updated.
	//                   Fields:
	//                   * GuildID int64 or discord.Snowflake
	GuildEmojisUpdateKey = "GUILD_EMOJI_UPDATE"
	//GuildIntegrationsUpdate Sent when a guild integration is updated.
	//                        Fields:
	//                        * GuildID int64 or discord.Snowflake
	//                        * Emojis []*discord.emoji.Emoji
	GuildIntegrationsUpdateKey = "GUILD_INTEGRATIONS_UPDATE"
	// GuildMemberAdd Sent when a new user joins a guild. The inner payload is a
	//                guild member object with these extra fields:
	//                * GuildID int64 or discord.Snowflake
	GuildMemberAddKey = "GUILD_MEMBER_ADD"
	// GuildMemberRemove Sent when a user is removed from a guild
	//                   (leave/kick/ban).
	//                   Fields:
	//                   * GuildID int64 or discord.Snowflake
	//                   * User *discord.user.User
	GuildMemberRemoveKey = "GUILD_MEMBER_REMOVE"
	// GuildMemberUpdate Sent when a guild member is updated.
	//                   Fields:
	//                   * GuildID int64 or discord.Snowflake
	//                   * Roles []int64 or []discord.Snowflake
	//                   * User *discord.user.User
	//                   * Nick string
	GuildMemberUpdateKey = "GUILD_MEMBER_UPDATE"
	// GuildMemberChunk Sent in response to Gateway Request Guild Members.
	//                  Fields:
	//                  * GuildID int64 or discord.Snowflake
	//                  * Members []*discord.member.Member
	GuildMembersChunkKey = "GUILD_MEMBER_CHUNK"
	// GuildRoleCreate Sent when a guild role is created.
	//                 Fields:
	//                 * GuildID int64 or discord.Snowflake
	//                 * Role *discord.role.Role
	GuildRoleCreateKey = "GUILD_ROLE_CREATE"
	// GuildRoleUpdate Sent when a guild role is created.
	//                 Fields:
	//                 * GuildID int64 or discord.Snowflake
	//                 * Role    *discord.role.Role
	GuildRoleUpdateKey = "GUILD_ROLE_UPDATE"
	// GuildRoleDelete Sent when a guild role is created.
	//                 Fields:
	//                 * GuildID int64 or discord.Snowflake
	//                 * RoleID  int64 or discord.Snowflake
	GuildRoleDeleteKey = "GUILD_ROLE_DELETE"
	// MessageCreate Sent when a message is created. The inner payload is a
	//               message object.
	MessageCreateKey = "MESSAGE_CREATE"
	// MessageUpdate Sent when a message is updated. The inner payload is a
	//               message object.
	//               NOTE! Has _at_least_ the GuildID and ChannelID fields.
	MessageUpdateKey = "MESSAGE_UPDATE"
	// MessageDelete Sent when a message is deleted.
	//               Fields:
	//               * ID        int64 or discord.Snowflake
	//               * ChannelID int64 or discord.Snowflake
	MessageDeleteKey = "MESSAGE_DELETE"
	// MessageDeleteBulk Sent when multiple messages are deleted at once.
	//                   Fields:
	//                   * IDs       []int64 or []discord.Snowflake
	//                   * ChannelID int64 or discord.Snowflake
	MessageDeleteBulkKey = "MESSAGE_DELETE_BULK"
	// MessageReactionAdd Sent when a user adds a reaction to a message.
	//                    Fields:
	//                    * UserID     int64 or discord.Snowflake
	//                    * ChannelID  int64 or discord.Snowflake
	//                    * MessageID  int64 or discord.Snowflake
	//                    * Emoji      *discord.emoji.Emoji
	MessageReactionAddKey = "MESSAGE_REACTION_ADD"
	// MessageReactionRemove Sent when a user removes a reaction from a message.
	//                       Fields:
	//                       * UserID     int64 or discord.Snowflake
	//                       * ChannelID  int64 or discord.Snowflake
	//                       * MessageID  int64 or discord.Snowflake
	//                       * Emoji      *discord.emoji.Emoji
	MessageReactionRemoveKey = "MESSAGE_REACTION_REMOVE"
	// MessageReactionRemoveAll Sent when a user explicitly removes all reactions
	//                          from a message.
	//                          Fields:
	//                          * ChannelID  int64 or discord.Snowflake
	//                          * MessageID  int64 or discord.Snowflake
	MessageReactionRemoveAllKey = "MESSAGE_REACTION_REMOVE_ALL"
	// PresenceUpdate A user's presence is their current state on a guild.
	//                This event is sent when a user's presence is updated
	//                for a guild.
	//                Fields:
	//                User    *discord.user.User
	//                Roles   []*discord.role.Role
	//                Game    *discord.game.Game
	//                GuildID int64 or discord.Snowflake
	//                Status  *string or *discord.presence.Status
	PresenceUpdateKey = "PRESENCE_UPDATE"
	// TypingStart Sent when a user starts typing in a channel.
	//             Fields: TODO
	TypingStartKey = "TYPING_START"
	// UserUpdate Sent when properties about the user change. Inner payload is a
	//            user object.
	UserUpdateKey = "USER_UPDATE"
	// VoiceStateUpdate Sent when someone joins/leaves/moves voice channels.
	//                  Inner payload is a voice state object.
	VoiceStateUpdateKey = "VOICE_STATE_UPDATE"
	// VoiceServerUpdate Sent when a guild's voice server is updated. This is
	//                   sent when initially connecting to voice, and when the
	//                   current voice instance fails over to a new server.
	//                   Fields:
	//                   * Token     string or discord.Token
	//                   * ChannelID int64 or discord.Snowflake
	//                   * Endpoint  string or discord.Endpoint
	VoiceServerUpdateKey = "VOICE_SERVER_UPDATE"
	// WebhooksUpdate Sent when a guild channel's webhook is created, updated, or
	//                deleted.
	//                Fields:
	//                * GuildID   int64 or discord.Snowflake
	//                * ChannelID int64 or discord.Snowflake
	WebhooksUpdateKey = "WEBHOOK_UPDATE"
)
