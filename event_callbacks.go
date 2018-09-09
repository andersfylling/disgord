package disgord

import (
	. "github.com/andersfylling/disgord/event"
)

// EventCallback is triggered on every event type
type EventCallback = func(session Session, box interface{})

// HelloCallback triggered in hello events
type HelloCallback = func(session Session, h *Hello)

// ReadyCallback triggered on READY events
type ReadyCallback = func(session Session, r *Ready)

// ResumedCallback triggered on RESUME events
type ResumedCallback = func(session Session, r *Resumed)

// InvalidSessionCallback triggered on INVALID_SESSION events
type InvalidSessionCallback = func(session Session, is *InvalidSession)

// ChannelCreateCallback triggered on CHANNEL_CREATE events
type ChannelCreateCallback = func(session Session, cc *ChannelCreate)

// ChannelUpdateCallback triggered on CHANNEL_UPDATE events
type ChannelUpdateCallback = func(session Session, cu *ChannelUpdate)

// ChannelDeleteCallback triggered on CHANNEL_DELETE events
type ChannelDeleteCallback = func(session Session, cd *ChannelDelete)

// ChannelPinsUpdateCallback triggered on CHANNEL_PINS_UPDATE events
type ChannelPinsUpdateCallback = func(session Session, cpu *ChannelPinsUpdate)

// GuildCreateCallback triggered on GUILD_CREATE events
type GuildCreateCallback = func(session Session, gc *GuildCreate)

// GuildUpdateCallback triggered on GUILD_UPDATE events
type GuildUpdateCallback = func(session Session, gu *GuildUpdate)

// GuildDeleteCallback triggered on GUILD_DELETE events
type GuildDeleteCallback = func(session Session, gd *GuildDelete)

// GuildBanAddCallback triggered on GUILD_BAN_ADD events
type GuildBanAddCallback = func(session Session, gba *GuildBanAdd)

// GuildBanRemoveCallback triggered on GUILD_BAN_REMOVE events
type GuildBanRemoveCallback = func(session Session, gbr *GuildBanRemove)

// GuildEmojisUpdateCallback triggered on GUILD_EMOJIS_UPDATE events
type GuildEmojisUpdateCallback = func(session Session, geu *GuildEmojisUpdate)

// GuildIntegrationsUpdateCallback triggered on GUILD_INTEGRATIONS_UPDATE events
type GuildIntegrationsUpdateCallback = func(session Session, giu *GuildIntegrationsUpdate)

// GuildMemberAddCallback triggered on GUILD_MEMBER_ADD events
type GuildMemberAddCallback = func(session Session, gma *GuildMemberAdd)

// GuildMemberRemoveCallback triggered on GUILD_MEMBER_REMOVE events
type GuildMemberRemoveCallback = func(session Session, gmr *GuildMemberRemove)

// GuildMemberUpdateCallback triggered on GUILD_MEMBER_UPDATE events
type GuildMemberUpdateCallback = func(session Session, gmu *GuildMemberUpdate)

// GuildMembersChunkCallback triggered on GUILD_MEMBERS_CHUNK events
type GuildMembersChunkCallback = func(session Session, gmc *GuildMembersChunk)

// GuildRoleCreateCallback triggered on GUILD_ROLE_CREATE events
type GuildRoleCreateCallback = func(session Session, grc *GuildRoleCreate)

// GuildRoleUpdateCallback triggered on GUILD_ROLE_UPDATE events
type GuildRoleUpdateCallback = func(session Session, gru *GuildRoleUpdate)

// GuildRoleDeleteCallback triggered on GUILD_ROLE_DELETE events
type GuildRoleDeleteCallback = func(session Session, grd *GuildRoleDelete)

// MessageCreateCallback triggered on MESSAGE_CREATE events
type MessageCreateCallback = func(session Session, mc *MessageCreate)

// MessageUpdateCallback triggered on MESSAGE_UPDATE events
type MessageUpdateCallback = func(session Session, mu *MessageUpdate)

// MessageDeleteCallback triggered on MESSAGE_DELETE events
type MessageDeleteCallback = func(session Session, md *MessageDelete)

// MessageDeleteBulkCallback triggered on MESSAGE_DELETE_BULK events
type MessageDeleteBulkCallback = func(session Session, mdb *MessageDeleteBulk)

// MessageReactionAddCallback triggered on MESSAGE_REACTION_ADD events
type MessageReactionAddCallback = func(session Session, mra *MessageReactionAdd)

// MessageReactionRemoveCallback triggered on MESSAGE_REACTION_REMOVE events
type MessageReactionRemoveCallback = func(session Session, mrr *MessageReactionRemove)

// MessageReactionRemoveAllCallback triggered on MESSAGE_REACTION_REMOVE_ALL events
type MessageReactionRemoveAllCallback = func(session Session, mrra *MessageReactionRemoveAll)

// PresenceUpdateCallback triggered on PRESENCE_UPDATE events
type PresenceUpdateCallback = func(session Session, pu *PresenceUpdate)

// TypingStartCallback triggered on TYPING_START events
type TypingStartCallback = func(session Session, ts *TypingStart)

// UserUpdateCallback triggerd on USER_UPDATE events
type UserUpdateCallback = func(session Session, uu *UserUpdate)

// VoiceStateUpdateCallback triggered on VOICE_STATE_UPDATE events
type VoiceStateUpdateCallback = func(session Session, vsu *VoiceStateUpdate)

// VoiceServerUpdateCallback triggered on VOICE_SERVER_UPDATE events
type VoiceServerUpdateCallback = func(session Session, vsu *VoiceServerUpdate)

// WebhooksUpdateCallback triggered on WEBHOOK_UPDATE events
type WebhooksUpdateCallback = func(session Session, wu *WebhooksUpdate)
