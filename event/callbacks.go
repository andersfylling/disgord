package event

import "github.com/andersfylling/disgord/disgordctx"

// socket
type HelloCallback = func(ctx disgordctx.Context, h *HelloBox)
type ReadyCallback = func(ctx disgordctx.Context, r *ReadyBox)
type ResumedCallback = func(ctx disgordctx.Context, r *ResumedBox)
type InvalidSessionCallback = func(ctx disgordctx.Context, is *InvalidSessionBox)

// channel
type ChannelCreateCallback = func(ctx disgordctx.Context, cc *ChannelCreateBox)
type ChannelUpdateCallback = func(ctx disgordctx.Context, cu *ChannelUpdateBox)
type ChannelDeleteCallback = func(ctx disgordctx.Context, cd *ChannelDeleteBox)
type ChannelPinsUpdateCallback = func(ctx disgordctx.Context, cpu *ChannelPinsUpdateBox)

// Guild in general
type GuildCreateCallback = func(ctx disgordctx.Context, gc *GuildCreateBox)
type GuildUpdateCallback = func(ctx disgordctx.Context, gu *GuildUpdateBox)
type GuildDeleteCallback = func(ctx disgordctx.Context, gd *GuildDeleteBox)
type GuildBanAddCallback = func(ctx disgordctx.Context, gba *GuildBanAddBox)
type GuildBanRemoveCallback = func(ctx disgordctx.Context, gbr *GuildBanRemoveBox)
type GuildEmojisUpdateCallback = func(ctx disgordctx.Context, geu *GuildEmojisUpdateBox)
type GuildIntegrationsUpdateCallback = func(ctx disgordctx.Context, giu *GuildIntegrationsUpdateBox)

// Guild Member
type GuildMemberAddCallback = func(ctx disgordctx.Context, gma *GuildMemberAddBox)
type GuildMemberRemoveCallback = func(ctx disgordctx.Context, gmr *GuildMemberRemoveBox)
type GuildMemberUpdateCallback = func(ctx disgordctx.Context, gmu *GuildMemberUpdateBox)
type GuildMembersChunkCallback = func(ctx disgordctx.Context, gmc *GuildMembersChunkBox)

// Guild role
type GuildRoleCreateCallback = func(ctx disgordctx.Context, grc *GuildRoleCreateBox)
type GuildRoleUpdateCallback = func(ctx disgordctx.Context, gru *GuildRoleUpdateBox)
type GuildRoleDeleteCallback = func(ctx disgordctx.Context, grd *GuildRoleDeleteBox)

// message
type MessageCreateCallback = func(ctx disgordctx.Context, mc *MessageCreateBox)
type MessageUpdateCallback = func(ctx disgordctx.Context, mu *MessageUpdateBox)
type MessageDeleteCallback = func(ctx disgordctx.Context, md *MessageDeleteBox)
type MessageDeleteBulkCallback = func(ctx disgordctx.Context, mdb *MessageDeleteBulkBox)

// message reaction
type MessageReactionAddCallback = func(ctx disgordctx.Context, mra *MessageReactionAddBox)
type MessageReactionRemoveCallback = func(ctx disgordctx.Context, mrr *MessageReactionRemoveBox)
type MessageReactionRemoveAllCallback = func(ctx disgordctx.Context, mrra *MessageReactionRemoveAllBox)

// presence
type PresenceUpdateCallback = func(ctx disgordctx.Context, pu *PresenceUpdateBox)

// typing start
type TypingStartCallback = func(ctx disgordctx.Context, ts *TypingStartBox)

// user update
type UserUpdateCallback = func(ctx disgordctx.Context, uu *UserUpdateBox)

// voice
type VoiceStateUpdateCallback = func(ctx disgordctx.Context, vsu *VoiceStateUpdateBox)
type VoiceServerUpdateCallback = func(ctx disgordctx.Context, vsu *VoiceServerUpdateBox)

// webhook
type WebhooksUpdateCallback = func(ctx disgordctx.Context, wu *WebhooksUpdateBox)
