package event

import (
	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/guild"
	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/disgord/voice"
	"github.com/andersfylling/disgord/webhook"
)

// socket
type ReadyCallback = func()
type ResumedCallback = func()

// channel
type ChannelCreateCallback = func(*channel.Channel)
type ChannelUpdateCallback = func(*channel.Channel)
type ChannelDeleteCallback = func(*channel.Channel)
type ChannelPinsUpdateCallback = func(*channel.Channel)

// Guild in general
type GuildCreateCallback = func(*guild.Guild)
type GuildUpdateCallback = func(*guild.Guild)
type GuildDeleteCallback = func(*guild.Guild)
type GuildBanAddCallback = func(*guild.Guild)
type GuildBanRemoveCallback = func(*guild.Guild)
type GuildEmojisUpdateCallback = func(*guild.Guild)
type GuildIntegrationsUpdateCallback = func(*guild.Guild)

// Guild Member
type GuildMemberAddCallback = func(*guild.Guild)
type GuildMemberRemoveCallback = func(*guild.Guild)
type GuildMemberUpdateCallback = func(*guild.Guild)
type GuildMemberChunkCallback = func(*guild.Guild)

// Guild role
type GuildRoleCreateCallback = func(*guild.Guild)
type GuildRoleUpdateCallback = func(*guild.Guild)
type GuildRoleDeleteCallback = func(*guild.Guild)

// message
type MessageCreateCallback = func(*channel.Message)
type MessageUpdateCallback = func(*channel.Message)
type MessageDeleteCallback = func(*channel.Message)
type MessageDeleteBulkCallback = func(*channel.Message)

// message reaction
type MessageReactionAddCallback = func(*channel.Message)
type MessageReactionRemoveCallback = func(*channel.Message)
type MessageReactionRemoveAllCallback = func(*channel.Message)

// presence
type PresenceUpdateCallback = func(*discord.Presence)

// typing start
type TypingStartCallback = func(*user.User, *channel.Channel)

// user update
type UserUpdateCallback = func(*user.User)

// voice
type VoiceStateUpdateCallback = func(*voice.State)
type VoiceServerUpdateCallback = func(*voice.State)

// webhook
type WebhooksUpdateCallback = func(*webhook.Webhook)
