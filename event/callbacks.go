package event

import (
	"context"
)

// socket
type HelloCallback = func(context.Context, *HelloBox)
type ReadyCallback = func(context.Context, *ReadyBox)
type ResumedCallback = func(context.Context, *ResumedBox)
type InvalidSessionCallback = func(context.Context, *InvalidSessionBox)

// channel
type ChannelCreateCallback = func(context.Context, *ChannelCreateBox)
type ChannelUpdateCallback = func(context.Context, *ChannelUpdateBox)
type ChannelDeleteCallback = func(context.Context, *ChannelDeleteBox)
type ChannelPinsUpdateCallback = func(context.Context, *ChannelPinsUpdateBox)

// Guild in general
type GuildCreateCallback = func(context.Context, *GuildCreateBox)
type GuildUpdateCallback = func(context.Context, *GuildUpdateBox)
type GuildDeleteCallback = func(context.Context, *GuildDeleteBox)
type GuildBanAddCallback = func(context.Context, *GuildBanAddBox)
type GuildBanRemoveCallback = func(context.Context, *GuildBanRemoveBox)
type GuildEmojisUpdateCallback = func(context.Context, *GuildEmojisUpdateBox)
type GuildIntegrationsUpdateCallback = func(context.Context, *GuildIntegrationsUpdateBox)

// Guild Member
type GuildMemberAddCallback = func(context.Context, *GuildMemberAddBox)
type GuildMemberRemoveCallback = func(context.Context, *GuildMemberRemoveBox)
type GuildMemberUpdateCallback = func(context.Context, *GuildMemberUpdateBox)
type GuildMembersChunkCallback = func(context.Context, *GuildMembersChunkBox)

// Guild role
type GuildRoleCreateCallback = func(context.Context, *GuildRoleCreateBox)
type GuildRoleUpdateCallback = func(context.Context, *GuildRoleUpdateBox)
type GuildRoleDeleteCallback = func(context.Context, *GuildRoleDeleteBox)

// message
type MessageCreateCallback = func(context.Context, *MessageCreateBox)
type MessageUpdateCallback = func(context.Context, *MessageUpdateBox)
type MessageDeleteCallback = func(context.Context, *MessageDeleteBox)
type MessageDeleteBulkCallback = func(context.Context, *MessageDeleteBulkBox)

// message reaction
type MessageReactionAddCallback = func(context.Context, *MessageReactionAddBox)
type MessageReactionRemoveCallback = func(context.Context, *MessageReactionRemoveBox)
type MessageReactionRemoveAllCallback = func(context.Context, *MessageReactionRemoveAllBox)

// presence
type PresenceUpdateCallback = func(context.Context, *PresenceUpdateBox)

// typing start
type TypingStartCallback = func(context.Context, *TypingStartBox)

// user update
type UserUpdateCallback = func(context.Context, *UserUpdateBox)

// voice
type VoiceStateUpdateCallback = func(context.Context, *VoiceStateUpdateBox)
type VoiceServerUpdateCallback = func(context.Context, *VoiceServerUpdateBox)

// webhook
type WebhooksUpdateCallback = func(context.Context, *WebhooksUpdateBox)
