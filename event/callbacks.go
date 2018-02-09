package event

import (
	"context"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/guild"
	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/disgord/voice"
	"github.com/andersfylling/disgord/webhook"
)

// socket
type ReadyCallback = func(ctx context.Context, r discord.Ready)
type ResumedCallback = func(ctx context.Context, resumed discord.Resumed)

// channel
type ChannelCreateCallback = func(ctx context.Context, c channel.Channel)
type ChannelUpdateCallback = func(ctx context.Context, c channel.Channel)
type ChannelDeleteCallback = func(ctx context.Context, c channel.Channel)
type ChannelPinsUpdateCallback = func(ctx context.Context, c channel.Channel)

// Guild in general
type GuildCreateCallback = func(ctx context.Context, g guild.Guild)
type GuildUpdateCallback = func(ctx context.Context, g guild.Guild)
type GuildDeleteCallback = func(ctx context.Context, g guild.Guild)
type GuildBanAddCallback = func(ctx context.Context, g guild.Guild)
type GuildBanRemoveCallback = func(ctx context.Context, g guild.Guild)
type GuildEmojisUpdateCallback = func(ctx context.Context, g guild.Guild)
type GuildIntegrationsUpdateCallback = func(ctx context.Context, g guild.Guild)

// Guild Member
type GuildMemberAddCallback = func(ctx context.Context, member guild.Member)
type GuildMemberRemoveCallback = func(ctx context.Context, member guild.Member)
type GuildMemberUpdateCallback = func(ctx context.Context, member guild.Member)
type GuildMemberChunkCallback = func(ctx context.Context, members []guild.Member)

// Guild role
type GuildRoleCreateCallback = func(ctx context.Context, role discord.RoleEvent)
type GuildRoleUpdateCallback = func(ctx context.Context, role discord.RoleEvent)
type GuildRoleDeleteCallback = func(ctx context.Context, role discord.RoleDeleteEvent)

// message
type MessageCreateCallback = func(ctx context.Context, msg channel.Message)
type MessageUpdateCallback = func(ctx context.Context, msg channel.Message)
type MessageDeleteCallback = func(ctx context.Context, msg channel.DeletedMessage)
type MessageDeleteBulkCallback = func(ctx context.Context, msgs []channel.Message)

// message reaction
type MessageReactionAddCallback = func(ctx context.Context, msg channel.Message)
type MessageReactionRemoveCallback = func(ctx context.Context, msg channel.Message)
type MessageReactionRemoveAllCallback = func(ctx context.Context, msgs []channel.Message)

// presence
type PresenceUpdateCallback = func(ctx context.Context, presence discord.Presence)

// typing start
type TypingStartCallback = func(ctx context.Context, ts channel.TypingStart)

// user update
type UserUpdateCallback = func(ctx context.Context, user user.User)

// voice
type VoiceStateUpdateCallback = func(ctx context.Context, voiceState voice.State)
type VoiceServerUpdateCallback = func(ctx context.Context, voiceState voice.State)

// webhook
type WebhooksUpdateCallback = func(ctx context.Context, webhook webhook.Webhook)
