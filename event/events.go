package event

// Event the discord api event type
type Event int8

const (

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
	Ready Event = iota
	// Resumed The resumed event is dispatched when a client has sent a resume
	//         payload to the gateway (for resuming existing sessions).
	//         Fields:
	//         * Trace []string
	Resumed Event = iota

	// Channel events

	// ChannelCreate Sent when a new channel is created, relevant to the current
	//               user. The inner payload is a DM channel or guild channel
	//               object.
	ChannelCreate Event = iota
	// ChannelUpdate Sent when a channel is updated. The inner payload is a guild
	//               channel object.
	ChannelUpdate Event = iota
	// ChannelDelete Sent when a channel relevant to the current user is deleted.
	//               The inner payload is a DM or Guild channel object.
	ChannelDelete Event = iota
	// ChannelPinsUpdate Sent when a message is pinned or unpinned in a text
	//                   channel. This is not sent when a pinned message is
	//                   deleted.
	//                   Fields:
	//                   * ChannelID int64 or discord.Snowflake
	//                   * LastPinTimestamp time.Now().UTC().Format(time.RFC3339)
	// TODO fix.
	ChannelPinsUpdate Event = iota

	// GUILD events

	// GuildCreate This event can be sent in three different scenarios:
	//             1. When a user is initially connecting, to lazily load and
	//                backfill information for all unavailable guilds sent in the
	//                Ready event.
	//             2. When a Guild becomes available again to the client.
	//             3. When the current user joins a new Guild.
	//             The inner payload is a guild object, with all the extra fields
	//             specified.
	GuildCreate Event = iota
	// GuildUpdate Sent when a guild is updated. The inner payload is a guild
	//             object.
	GuildUpdate Event = iota
	// GuildDelete Sent when a guild becomes unavailable during a guild outage,
	//             or when the user leaves or is removed from a guild. The inner
	//             payload is an unavailable guild object. If the unavailable
	//             field is not set, the user was removed from the guild.
	GuildDelete Event = iota
	// GuildBanAdd Sent when a user is banned from a guild. The inner payload is
	//             a user object, with an extra guild_id key.
	GuildBanAdd Event = iota
	// GuildBanRemove Sent when a user is unbanned from a guild. The inner
	//                payload is a user object, with an extra guild_id key.
	GuildBanRemove Event = iota
	// GuildEmojisUpdate Sent when a guild's emojis have been updated.
	//                   Fields:
	//                   * GuildID int64 or discord.Snowflake
	GuildEmojisUpdate Event = iota
	//GuildIntegrationsUpdate Sent when a guild integration is updated.
	//                        Fields:
	//                        * GuildID int64 or discord.Snowflake
	//                        * Emojis []*discord.emoji.Emoji
	GuildIntegrationsUpdate Event = iota
	// GuildMemberAdd Sent when a new user joins a guild. The inner payload is a
	//                guild member object with these extra fields:
	//                * GuildID int64 or discord.Snowflake
	GuildMemberAdd Event = iota
	// GuildMemberRemove Sent when a user is removed from a guild
	//                   (leave/kick/ban).
	//                   Fields:
	//                   * GuildID int64 or discord.Snowflake
	//                   * User *discord.user.User
	GuildMemberRemove Event = iota
	// GuildMemberUpdate Sent when a guild member is updated.
	//                   Fields:
	//                   * GuildID int64 or discord.Snowflake
	//                   * Roles []int64 or []discord.Snowflake
	//                   * User *discord.user.User
	//                   * Nick string
	GuildMemberUpdate Event = iota
	// GuildMemberChunk Sent in response to Gateway Request Guild Members.
	//                  Fields:
	//                  * GuildID int64 or discord.Snowflake
	//                  * Members []*discord.member.Member
	GuildMemberChunk Event = iota
	// GuildRoleCreate Sent when a guild role is created.
	//                 Fields:
	//                 * GuildID int64 or discord.Snowflake
	//                 * Role *discord.role.Role
	GuildRoleCreate Event = iota
	// GuildRoleUpdate Sent when a guild role is created.
	//                 Fields:
	//                 * GuildID int64 or discord.Snowflake
	//                 * Role    *discord.role.Role
	GuildRoleUpdate Event = iota
	// GuildRoleDelete Sent when a guild role is created.
	//                 Fields:
	//                 * GuildID int64 or discord.Snowflake
	//                 * RoleID  int64 or discord.Snowflake
	GuildRoleDelete Event = iota
	// MessageCreate Sent when a message is created. The inner payload is a
	//               message object.
	MessageCreate Event = iota
	// MessageUpdate Sent when a message is updated. The inner payload is a
	//               message object.
	//               NOTE! Has _at_least_ the GuildID and ChannelID fields.
	MessageUpdate Event = iota
	// MessageDelete Sent when a message is deleted.
	//               Fields:
	//               * ID        int64 or discord.Snowflake
	//               * ChannelID int64 or discord.Snowflake
	MessageDelete Event = iota
	// MessageDeleteBulk Sent when multiple messages are deleted at once.
	//                   Fields:
	//                   * IDs       []int64 or []discord.Snowflake
	//                   * ChannelID int64 or discord.Snowflake
	MessageDeleteBulk Event = iota
	// MessageReactionAdd Sent when a user adds a reaction to a message.
	//                    Fields:
	//                    * UserID     int64 or discord.Snowflake
	//                    * ChannelID  int64 or discord.Snowflake
	//                    * MessageID  int64 or discord.Snowflake
	//                    * Emoji      *discord.emoji.Emoji
	MessageReactionAdd Event = iota
	// MessageReactionRemove Sent when a user removes a reaction from a message.
	//                       Fields:
	//                       * UserID     int64 or discord.Snowflake
	//                       * ChannelID  int64 or discord.Snowflake
	//                       * MessageID  int64 or discord.Snowflake
	//                       * Emoji      *discord.emoji.Emoji
	MessageReactionRemove Event = iota
	// MessageReactionRemoveAll Sent when a user explicitly removes all reactions
	//                          from a message.
	//                          Fields:
	//                          * ChannelID  int64 or discord.Snowflake
	//                          * MessageID  int64 or discord.Snowflake
	MessageReactionRemoveAll Event = iota
	// PresenceUpdate A user's presence is their current state on a guild.
	//                This event is sent when a user's presence is updated
	//                for a guild.
	//                Fields:
	//                User    *discord.user.User
	//                Roles   []*discord.role.Role
	//                Game    *discord.game.Game
	//                GuildID int64 or discord.Snowflake
	//                Status  *string or *discord.presence.Status
	PresenceUpdate Event = iota
	// TypingStart Sent when a user starts typing in a channel.
	//             Fields: TODO
	TypingStart Event = iota
	// UserUpdate Sent when properties about the user change. Inner payload is a
	//            user object.
	UserUpdate Event = iota
	// VoiceStateUpdate Sent when someone joins/leaves/moves voice channels.
	//                  Inner payload is a voice state object.
	VoiceStateUpdate Event = iota
	// VoiceServerUpdate Sent when a guild's voice server is updated. This is
	//                   sent when initially connecting to voice, and when the
	//                   current voice instance fails over to a new server.
	//                   Fields:
	//                   * Token     string or discord.Token
	//                   * ChannelID int64 or discord.Snowflake
	//                   * Endpoint  string or discord.Endpoint
	VoiceServerUpdate Event = iota
	// WebhooksUpdate Sent when a guild channel's webhook is created, updated, or
	//                deleted.
	//                Fields:
	//                * GuildID   int64 or discord.Snowflake
	//                * ChannelID int64 or discord.Snowflake
	WebhooksUpdate Event = iota
)
