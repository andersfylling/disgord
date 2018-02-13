package event

import (
	"context"
)

// socket
type HelloCallback = func(ctx context.Context, h *HelloBox)
type ReadyCallback = func(ctx context.Context, r *ReadyBox)
type ResumedCallback = func(ctx context.Context, r *ResumedBox)
type InvalidSessionCallback = func(ctx context.Context, is *InvalidSessionBox)

// channel
type ChannelCreateCallback = func(ctx context.Context, cc *ChannelCreateBox)
type ChannelUpdateCallback = func(ctx context.Context, cu *ChannelUpdateBox)
type ChannelDeleteCallback = func(ctx context.Context, cd *ChannelDeleteBox)
type ChannelPinsUpdateCallback = func(ctx context.Context, cpu *ChannelPinsUpdateBox)

// Guild in general
type GuildCreateCallback = func(ctx context.Context, gc *GuildCreateBox)
type GuildUpdateCallback = func(ctx context.Context, gu *GuildUpdateBox)
type GuildDeleteCallback = func(ctx context.Context, gd *GuildDeleteBox)
type GuildBanAddCallback = func(ctx context.Context, gba *GuildBanAddBox)
type GuildBanRemoveCallback = func(ctx context.Context, gbr *GuildBanRemoveBox)
type GuildEmojisUpdateCallback = func(ctx context.Context, geu *GuildEmojisUpdateBox)
type GuildIntegrationsUpdateCallback = func(ctx context.Context, giu *GuildIntegrationsUpdateBox)

// Guild Member
type GuildMemberAddCallback = func(ctx context.Context, gma *GuildMemberAddBox)
type GuildMemberRemoveCallback = func(ctx context.Context, gmr *GuildMemberRemoveBox)
type GuildMemberUpdateCallback = func(ctx context.Context, gmu *GuildMemberUpdateBox)
type GuildMembersChunkCallback = func(ctx context.Context, gmc *GuildMembersChunkBox)

// Guild role
type GuildRoleCreateCallback = func(ctx context.Context, grc *GuildRoleCreateBox)
type GuildRoleUpdateCallback = func(ctx context.Context, gru *GuildRoleUpdateBox)
type GuildRoleDeleteCallback = func(ctx context.Context, grd *GuildRoleDeleteBox)

// message
type MessageCreateCallback = func(ctx context.Context, mc *MessageCreateBox)
type MessageUpdateCallback = func(ctx context.Context, mu *MessageUpdateBox)
type MessageDeleteCallback = func(ctx context.Context, md *MessageDeleteBox)
type MessageDeleteBulkCallback = func(ctx context.Context, mdb *MessageDeleteBulkBox)

// message reaction
type MessageReactionAddCallback = func(ctx context.Context, mra *MessageReactionAddBox)
type MessageReactionRemoveCallback = func(ctx context.Context, mrr *MessageReactionRemoveBox)
type MessageReactionRemoveAllCallback = func(ctx context.Context, mrra *MessageReactionRemoveAllBox)

// presence
type PresenceUpdateCallback = func(ctx context.Context, pu *PresenceUpdateBox)

// typing start
type TypingStartCallback = func(ctx context.Context, ts *TypingStartBox)

// user update
type UserUpdateCallback = func(ctx context.Context, uu *UserUpdateBox)

// voice
type VoiceStateUpdateCallback = func(ctx context.Context, vsu *VoiceStateUpdateBox)
type VoiceServerUpdateCallback = func(ctx context.Context, vsu *VoiceServerUpdateBox)

// webhook
type WebhooksUpdateCallback = func(ctx context.Context, wu *WebhooksUpdateBox)
