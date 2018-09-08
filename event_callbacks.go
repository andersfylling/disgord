package disgord

import (
	. "github.com/andersfylling/disgord/event"
)

type EventCallback = func(session Session, box interface{})

// socket
type HelloCallback = func(session Session, h *Hello)
type ReadyCallback = func(session Session, r *Ready)
type ResumedCallback = func(session Session, r *Resumed)
type InvalidSessionCallback = func(session Session, is *InvalidSession)

// channel
type ChannelCreateCallback = func(session Session, cc *ChannelCreate)
type ChannelUpdateCallback = func(session Session, cu *ChannelUpdate)
type ChannelDeleteCallback = func(session Session, cd *ChannelDelete)
type ChannelPinsUpdateCallback = func(session Session, cpu *ChannelPinsUpdate)

// Guild in general
type GuildCreateCallback = func(session Session, gc *GuildCreate)
type GuildUpdateCallback = func(session Session, gu *GuildUpdate)
type GuildDeleteCallback = func(session Session, gd *GuildDelete)
type GuildBanAddCallback = func(session Session, gba *GuildBanAdd)
type GuildBanRemoveCallback = func(session Session, gbr *GuildBanRemove)
type GuildEmojisUpdateCallback = func(session Session, geu *GuildEmojisUpdate)
type GuildIntegrationsUpdateCallback = func(session Session, giu *GuildIntegrationsUpdate)

// Guild Member
type GuildMemberAddCallback = func(session Session, gma *GuildMemberAdd)
type GuildMemberRemoveCallback = func(session Session, gmr *GuildMemberRemove)
type GuildMemberUpdateCallback = func(session Session, gmu *GuildMemberUpdate)
type GuildMembersChunkCallback = func(session Session, gmc *GuildMembersChunk)

// Guild role
type GuildRoleCreateCallback = func(session Session, grc *GuildRoleCreate)
type GuildRoleUpdateCallback = func(session Session, gru *GuildRoleUpdate)
type GuildRoleDeleteCallback = func(session Session, grd *GuildRoleDelete)

// message
type MessageCreateCallback = func(session Session, mc *MessageCreate)
type MessageUpdateCallback = func(session Session, mu *MessageUpdate)
type MessageDeleteCallback = func(session Session, md *MessageDelete)
type MessageDeleteBulkCallback = func(session Session, mdb *MessageDeleteBulk)

// message reaction
type MessageReactionAddCallback = func(session Session, mra *MessageReactionAdd)
type MessageReactionRemoveCallback = func(session Session, mrr *MessageReactionRemove)
type MessageReactionRemoveAllCallback = func(session Session, mrra *MessageReactionRemoveAll)

// presence
type PresenceUpdateCallback = func(session Session, pu *PresenceUpdate)

// typing start
type TypingStartCallback = func(session Session, ts *TypingStart)

// user update
type UserUpdateCallback = func(session Session, uu *UserUpdate)

// voice
type VoiceStateUpdateCallback = func(session Session, vsu *VoiceStateUpdate)
type VoiceServerUpdateCallback = func(session Session, vsu *VoiceServerUpdate)

// webhook
type WebhooksUpdateCallback = func(session Session, wu *WebhooksUpdate)
