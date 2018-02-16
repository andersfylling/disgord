# Progress

## Logic

- [ ] Sharding
- [ ] Sharding for large bots (+100,000 guilds)
- [ ] Rate limiting
- [x] Guild availability
- [ ] Socketing

  - [x] Connecting
  - [x] Reconnect/resume (Needs revision)
  - [x] Handling invalid connection
  - [x] Sequence tracking
  - [x] JSON support (Won't support ETF)
  - [ ] Transport compression
  - [x] heartbeat
  - [x] Identify

- [ ] OAuth2

- [ ] RPC

- [ ] Voice

## Data structs for events:

_Every event data struct has a `Box` suffix to clarify they're a container for the event data._ _While event keys, event identifiers, has the `Key` suffix._

- ~~[ ] Hello~~
- [x] Ready
- [x] Resumed
- ~~[ ] InvalidSession~~
- [x] ChannelCreate
- [x] ChannelUpdate
- [x] ChannelDelete
- [x] ChannelPinsUpdate
- [x] GuildCreate
- [x] GuildUpdate
- [x] GuildDelete
- [x] GuildBanAdd
- [x] GuildBanRemove
- [x] GuildEmojisUpdate
- [x] GuildIntegrationsUpdate
- [x] GuildMemberAdd
- [x] GuildMemberRemove
- [x] GuildMemberUpdate
- [x] GuildMemberChunk
- [x] GuildRoleCreate
- [x] GuildRoleUpdate
- [x] GuildRoleDelete
- [x] MessageCreate
- [x] MessageUpdate
- [x] MessageDelete
- [x] MessageDeleteBulk
- [x] MessageReactionAdd
- [x] MessageReactionRemove
- [x] MessageReactionRemoveAll
- [x] PresenceUpdate (TODO: review)
- [x] TypingStart
- [x] UserUpdate
- [x] VoiceStateUpdate
- [x] VoiceServerUpdate
- [x] WebhooksUpdate

## Event dispatchers:

- ~~[ ] Hello~~
- [x] Ready
- [x] Resumed
- ~~[ ] InvalidSession~~
- [x] ChannelCreate
- [x] ChannelUpdate
- [x] ChannelDelete
- [x] ChannelPinsUpdate
- [x] GuildCreate
- [x] GuildUpdate
- [x] GuildDelete
- [x] GuildBanAdd
- [x] GuildBanRemove
- [x] GuildEmojisUpdate
- [x] GuildIntegrationsUpdate
- [x] GuildMemberAdd
- [x] GuildMemberRemove
- [x] GuildMemberUpdate
- [x] GuildMemberChunk
- [x] GuildRoleCreate
- [x] GuildRoleUpdate
- [x] GuildRoleDelete
- [x] MessageCreate
- [x] MessageUpdate
- [x] MessageDelete
- [x] MessageDeleteBulk
- [x] MessageReactionAdd
- [x] MessageReactionRemove
- [x] MessageReactionRemoveAll
- [x] PresenceUpdate
- [x] TypingStart
- [x] UserUpdate
- [x] VoiceStateUpdate
- [x] VoiceServerUpdate
- [x] WebhooksUpdate

## Caching:

- ~~[ ] Hello~~
- ~~[ ] Ready~~
- ~~[ ] Resumed~~
- ~~[ ] InvalidSession~~
- [ ] ChannelCreate
- [ ] ChannelUpdate
- [ ] ChannelDelete
- [ ] ChannelPinsUpdate
- [ ] GuildCreate
- [ ] GuildUpdate
- [ ] GuildDelete
- [ ] GuildBanAdd
- [ ] GuildBanRemove
- [ ] GuildEmojisUpdate
- [ ] GuildIntegrationsUpdate
- [ ] GuildMemberAdd
- [ ] GuildMemberRemove
- [ ] GuildMemberUpdate
- [ ] GuildMemberChunk
- [ ] GuildRoleCreate
- [ ] GuildRoleUpdate
- [ ] GuildRoleDelete
- [ ] MessageCreate
- [ ] MessageUpdate
- [ ] MessageDelete
- [ ] MessageDeleteBulk
- [ ] MessageReactionAdd
- [ ] MessageReactionRemove
- [ ] MessageReactionRemoveAll
- [ ] PresenceUpdate
- ~~[ ] TypingStart~~
- [ ] UserUpdate
- [ ] VoiceStateUpdate
- [ ] VoiceServerUpdate
- [ ] WebhooksUpdate
