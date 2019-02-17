// Package shortevent is a short version of the event pkg. Please import this pkg as event.evt. This package should be
// use with case, as Discord might add a new event that can cause the names here to break or not make sense anymore.
// I really do not recommend using this pkg instead of the event pkg.
package shortevent

import "github.com/andersfylling/disgord/event"

const Rdy = event.Ready
const Resumed = event.Resumed
const Channel = event.ChannelCreate
const CUpdate = event.ChannelUpdate
const CDelete = event.ChannelDelete
const CPinsUpdate = event.ChannelPinsUpdate
const Typing = event.TypingStart
const Msg = event.MessageCreate
const MsgUpdate = event.MessageUpdate
const MsgDelete = event.MessageDelete
const MsgDeleteBulk = event.MessageDeleteBulk
const MsgReaction = event.MessageReactionAdd
const MsgReactionRm = event.MessageReactionRemove
const MsgReactionRmAll = event.MessageReactionRemoveAll
const EmojisUpdate = event.GuildEmojisUpdate
const Guild = event.GuildCreate
const GUpdate = event.GuildUpdate
const GDelete = event.GuildDelete
const Ban = event.GuildBanAdd
const UnBan = event.GuildBanRemove
const IntegrationsUpdate = event.GuildIntegrationsUpdate
const Member = event.GuildMemberAdd
const MemberRm = event.GuildMemberRemove
const MemberUpdate = event.GuildMemberUpdate
const MembersChunk = event.GuildMembersChunk
const Role = event.GuildRoleCreate
const RoleUpdate = event.GuildRoleUpdate
const RoleDelete = event.GuildRoleDelete
const PresenceUpdate = event.PresenceUpdate
const PresencesReplace = event.PresencesReplace
const UserUpdate = event.UserUpdate
const VoiceStateUpdate = event.VoiceStateUpdate
const VoiceServerUpdate = event.VoiceServerUpdate
const WebhooksUpdate = event.WebhooksUpdate
