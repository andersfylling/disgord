package event

import (
	"context"

	"github.com/andersfylling/disgord/resource"
	. "github.com/andersfylling/snowflake"
)

// KeyChannelCreate Sent when a new channel is created, relevant to the current
//               user. The inner payload is a DM channel or guild channel
//               object.
const KeyChannelCreate = "CHANNEL_CREATE"

// ChannelCreateBox	new channel created
type ChannelCreate struct {
	Channel *resource.Channel `json:"channel"`
	Ctx     context.Context   `json:"-"`
}

// KeyChannelUpdate Sent when a channel is updated. The inner payload is a guild
//               channel object.
const KeyChannelUpdate = "CHANNEL_UPDATE"

// ChannelUpdateBox	channel was updated
type ChannelUpdate struct {
	Channel *resource.Channel `json:"channel"`
	Ctx     context.Context   `json:"-"`
}

// KeyChannelDelete Sent when a channel relevant to the current user is deleted.
//               The inner payload is a DM or Guild channel object.
const KeyChannelDelete = "CHANNEL_DELETE"

// ChannelDelete	channel was deleted
type ChannelDelete struct {
	Channel *resource.Channel `json:"channel"`
	Ctx     context.Context   `json:"-"`
}

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
	LastPinTimestamp resource.Timestamp `json:"last_pin_timestamp,omitempty"` // ?|
	Ctx              context.Context    `json:"-"`
}

// KeyTypingStart Sent when a user starts typing in a channel.
//             Fields: TODO
const KeyTypingStart = "TYPING_START"

// TypingStart	user started typing in a channel
type TypingStart struct {
	ChannelID     Snowflake       `json:"channel_id"`
	UserID        Snowflake       `json:"user_id"`
	TimestampUnix int             `json:"timestamp"`
	Ctx           context.Context `json:"-"`
}

// KeyMessageCreate Sent when a message is created. The inner payload is a
//               message object.
const KeyMessageCreate = "MESSAGE_CREATE"

// MessageCreate	message was created
type MessageCreate struct {
	Message *resource.Message
	Ctx     context.Context `json:"-"`
}

// KeyMessageUpdate Sent when a message is updated. The inner payload is a
//               message object.
//               NOTE! Has _at_least_ the GuildID and ChannelID fields.
const KeyMessageUpdate = "MESSAGE_UPDATE"

// MessageUpdate	message was edited
type MessageUpdate struct {
	Message *resource.Message
	Ctx     context.Context `json:"-"`
}

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
	PartialEmoji *resource.Emoji `json:"emoji"`
	Ctx          context.Context `json:"-"`
}

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
	PartialEmoji *resource.Emoji `json:"emoji"`
	Ctx          context.Context `json:"-"`
}

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
