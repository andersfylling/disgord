package event

import (
	"sync"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/emoji"
	"github.com/andersfylling/disgord/guild"
	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/disgord/voice"
	"github.com/andersfylling/snowflake"
)

// event.Box is a container for a given event type which holds different data structures.

// HelloBox defines the heartbeat interval
type HelloBox struct {
	HeartbeatInterval uint     `json:"heartbeat_interval"`
	Trace             []string `json:"_trace"`
}

// ReadyBox	contains the initial state information
type ReadyBox struct {
	APIVersion int                         `json:"v"`
	User       *user.User                  `json:"user"`
	Guilds     []*discord.GuildUnavailable `json:"guilds"`

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
}

// ResumedBox	response to Resume
type ResumedBox struct {
	Trace []string `json:"_trace"`
}

// InvalidSessionBox	failure response to Identify or Resume or invalid active session
type InvalidSessionBox struct{}

// ChannelCreateBox	new channel created
type ChannelCreateBox struct {
	Channel *channel.Channel `json:"channel"`
}

// ChannelUpdateBox	channel was updated
type ChannelUpdateBox struct {
	Channel *channel.Channel `json:"channel"`
}

// ChannelDeleteBox	channel was deleted
type ChannelDeleteBox struct {
	Channel *channel.Channel `json:"channel"`
}

// ChannelPinsUpdateBox	message was pinned or unpinned
type ChannelPinsUpdateBox struct {
	// ChannelID snowflake	the id of the channel
	ChannelID snowflake.ID `json:"channel_id"`

	// LastPinTimestamp	ISO8601 timestamp	the time at which the most recent pinned message was pinned
	LastPinTimestamp discord.Timestamp `json:"last_pin_timestamp,omitempty"` // ?|
}

// GuildCreateBox	This event can be sent in three different scenarios:
//								1. When a user is initially connecting, to lazily load and backfill information for
//									 all unavailable guilds sent in the Ready event.
//								2. When a Guild becomes available again to the client.
// 								3. When the current user joins a new Guild.
type GuildCreateBox struct {
	Guild *guild.Guild `json:"guild"`
}

// GuildUpdateBox	guild was updated
type GuildUpdateBox struct {
	Guild *guild.Guild `json:"guild"`
}

// GuildDeleteBox	guild became unavailable, or user left/was removed from a guild
type GuildDeleteBox struct {
	UnavailableGuild *discord.GuildUnavailable `json:"guild_unavailable"`
}

// GuildBanAddBox	user was banned from a guild
type GuildBanAddBox struct {
	User *user.User `json:"user"`
}

// GuildBanRemoveBox	user was unbanned from a guild
type GuildBanRemoveBox struct {
	User *user.User `json:"user"`
}

// GuildEmojisUpdateBox	guild emojis were updated
type GuildEmojisUpdateBox struct {
	GuildID snowflake.ID   `json:"guild_id"`
	Emojis  []*emoji.Emoji `json:"emojis"`
}

// GuildIntegrationsUpdateBox	guild integration was updated
type GuildIntegrationsUpdateBox struct {
	GuildID snowflake.ID `json:"guild_id"`
}

// GuildMemberAddBox	new user joined a guild
type GuildMemberAddBox struct {
	Member *guild.Member `json:"member"`
}

// GuildMemberRemoveBox	user was removed from a guild
type GuildMemberRemoveBox struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    *user.User   `json:"user"`
}

// GuildMemberUpdateBox	guild member was updated
type GuildMemberUpdateBox struct {
	GuildID snowflake.ID    `json:"guild_id"`
	Roles   []*discord.Role `json:"roles"`
	User    *user.User      `json:"user"`
	Nick    string          `json:"nick"`
}

// GuildMembersChunkBox	response to Request Guild Members
type GuildMembersChunkBox struct {
	GuildID snowflake.ID    `json:"guild_id"`
	Members []*guild.Member `json:"members"`
}

// GuildRoleCreateBox	guild role was created
type GuildRoleCreateBox struct {
	GuildID snowflake.ID  `json:"guild_id"`
	Role    *discord.Role `json:"role"`
}

// GuildRoleUpdateBox	guild role was updated
type GuildRoleUpdateBox struct {
	GuildID snowflake.ID  `json:"guild_id"`
	Role    *discord.Role `json:"role"`
}

// GuildRoleDeleteBox	guild role was deleted
type GuildRoleDeleteBox struct {
	GuildID snowflake.ID `json:"guild_id"`
	RoleID  snowflake.ID `json:"role_id"`
}

// MessageCreateBox	message was created
type MessageCreateBox struct {
	Message *channel.Message
}

// MessageUpdateBox	message was edited
type MessageUpdateBox struct {
	Message *channel.Message
}

// MessageDeleteBox	message was deleted
type MessageDeleteBox struct {
	MessageID snowflake.ID `json:"id"`
	ChannelID snowflake.ID `json:"channel_id"`
}

// MessageDeleteBulkBox	multiple messages were deleted at once
type MessageDeleteBulkBox struct {
	MessageIDs []snowflake.ID `json:"ids"`
	ChannelID  snowflake.ID   `json:"channel_id"`
}

// MessageReactionAddBox	user reacted to a message
type MessageReactionAddBox struct {
	UserID    snowflake.ID `json:"user_id"`
	ChannelID snowflake.ID `json:"channel_id"`
	MessageID snowflake.ID `json:"message_id"`
	// PartialEmoji id and name. id might be nil
	PartialEmoji *emoji.Emoji `json:"emoji"`
}

// MessageReactionRemoveBox	user removed a reaction from a message
type MessageReactionRemoveBox struct {
	UserID    snowflake.ID `json:"user_id"`
	ChannelID snowflake.ID `json:"channel_id"`
	MessageID snowflake.ID `json:"message_id"`
	// PartialEmoji id and name. id might be nil
	PartialEmoji *emoji.Emoji `json:"emoji"`
}

// MessageReactionRemoveAllBox	all reactions were explicitly removed from a message
type MessageReactionRemoveAllBox struct {
	ChannelID snowflake.ID `json:"channel_id"`
	MessageID snowflake.ID `json:"id"`
}

// PresenceUpdateBox	user's presence was updated in a guild
type PresenceUpdateBox struct {
	User    *user.User        `json:"user"`
	RoleIDs []snowflake.ID    `json:"roles"`
	Game    *discord.Activity `json:"game"`
	GuildID snowflake.ID      `json:"guild_id"`

	// Status either "idle", "dnd", "online", or "offline"
	// TODO: constants somewhere..
	Status string `json:"status"`
}

// TypingStartBox	user started typing in a channel
type TypingStartBox struct {
	ChannelID     snowflake.ID `json:"channel_id"`
	UserID        snowflake.ID `json:"user_id"`
	TimestampUnix int          `json:"timestamp"`
}

// UserUpdateBox	properties about a user changed
type UserUpdateBox struct {
	User *user.User `json:"user"`
}

// VoiceStateUpdateBox	someone joined, left, or moved a voice channel
type VoiceStateUpdateBox struct {
	VoiceState *voice.State `json:"voice_state"`
}

// VoiceServerUpdateBox	guild's voice server was updated
// Sent when a guild's voice server is updated.
// This is sent when initially connecting to voice,
// and when the current voice instance fails over to a new server.
type VoiceServerUpdateBox struct {
	Token    string       `json:"token"`
	GuildID  snowflake.ID `json:"guild_id"`
	Endpoint string       `json:"endpoint"`
}

// WebhooksUpdateBox guild channel webhook was created, update, or deleted
type WebhooksUpdateBox struct {
	GuildID   snowflake.ID `json:"guild_id"`
	ChannelID snowflake.ID `json:"channel_id"`
}
