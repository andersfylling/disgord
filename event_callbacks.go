package disgord

import (
	. "github.com/andersfylling/disgord/event"
)

type EventCallback = func(session Session, box interface{})

// socket
type HelloCallback = func(session Session, h *HelloBox)
type ReadyCallback = func(session Session, r *ReadyBox)
type ResumedCallback = func(session Session, r *ResumedBox)
type InvalidSessionCallback = func(session Session, is *InvalidSessionBox)

// channel
type ChannelCreateCallback = func(session Session, cc *ChannelCreateBox)
type ChannelUpdateCallback = func(session Session, cu *ChannelUpdateBox)
type ChannelDeleteCallback = func(session Session, cd *ChannelDeleteBox)
type ChannelPinsUpdateCallback = func(session Session, cpu *ChannelPinsUpdateBox)

// Guild in general
type GuildCreateCallback = func(session Session, gc *GuildCreateBox)
type GuildUpdateCallback = func(session Session, gu *GuildUpdateBox)
type GuildDeleteCallback = func(session Session, gd *GuildDeleteBox)
type GuildBanAddCallback = func(session Session, gba *GuildBanAddBox)
type GuildBanRemoveCallback = func(session Session, gbr *GuildBanRemoveBox)
type GuildEmojisUpdateCallback = func(session Session, geu *GuildEmojisUpdateBox)
type GuildIntegrationsUpdateCallback = func(session Session, giu *GuildIntegrationsUpdateBox)

// Guild Member
type GuildMemberAddCallback = func(session Session, gma *GuildMemberAddBox)
type GuildMemberRemoveCallback = func(session Session, gmr *GuildMemberRemoveBox)
type GuildMemberUpdateCallback = func(session Session, gmu *GuildMemberUpdateBox)
type GuildMembersChunkCallback = func(session Session, gmc *GuildMembersChunkBox)

// Guild role
type GuildRoleCreateCallback = func(session Session, grc *GuildRoleCreateBox)
type GuildRoleUpdateCallback = func(session Session, gru *GuildRoleUpdateBox)
type GuildRoleDeleteCallback = func(session Session, grd *GuildRoleDeleteBox)

// message
type MessageCreateCallback = func(session Session, mc *MessageCreateBox)
type MessageUpdateCallback = func(session Session, mu *MessageUpdateBox)
type MessageDeleteCallback = func(session Session, md *MessageDeleteBox)
type MessageDeleteBulkCallback = func(session Session, mdb *MessageDeleteBulkBox)

// message reaction
type MessageReactionAddCallback = func(session Session, mra *MessageReactionAddBox)
type MessageReactionRemoveCallback = func(session Session, mrr *MessageReactionRemoveBox)
type MessageReactionRemoveAllCallback = func(session Session, mrra *MessageReactionRemoveAllBox)

// presence
type PresenceUpdateCallback = func(session Session, pu *PresenceUpdateBox)

// typing start
type TypingStartCallback = func(session Session, ts *TypingStartBox)

// user update
type UserUpdateCallback = func(session Session, uu *UserUpdateBox)

// voice
type VoiceStateUpdateCallback = func(session Session, vsu *VoiceStateUpdateBox)
type VoiceServerUpdateCallback = func(session Session, vsu *VoiceServerUpdateBox)

// webhook
type WebhooksUpdateCallback = func(session Session, wu *WebhooksUpdateBox)
