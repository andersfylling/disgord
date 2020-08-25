package disgord

type CacheImmutable struct {
	Cache
}

func (c *CacheImmutable) ChannelCreate(data []byte) (evt *ChannelCreate, err error) {
	if evt, _ = c.ChannelCreate(data); evt != nil {
		evt = evt.DeepCopy().(*ChannelCreate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) ChannelDelete(data []byte) (evt *ChannelDelete, err error) {
	if evt, _ = c.ChannelDelete(data); evt != nil {
		evt = evt.DeepCopy().(*ChannelDelete)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) ChannelPinsUpdate(data []byte) (evt *ChannelPinsUpdate, err error) {
	if evt, _ = c.ChannelPinsUpdate(data); evt != nil {
		evt = evt.DeepCopy().(*ChannelPinsUpdate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) ChannelUpdate(data []byte) (evt *ChannelUpdate, err error) {
	if evt, _ = c.ChannelUpdate(data); evt != nil {
		evt = evt.DeepCopy().(*ChannelUpdate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) GuildBanAdd(data []byte) (evt *GuildBanAdd, err error) {
	if evt, _ = c.GuildBanAdd(data); evt != nil {
		evt = evt.DeepCopy().(*GuildBanAdd)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) GuildBanRemove(data []byte) (evt *GuildBanRemove, err error) {
	if evt, _ = c.GuildBanRemove(data); evt != nil {
		evt = evt.DeepCopy().(*GuildBanRemove)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) GuildCreate(data []byte) (evt *GuildCreate, err error) {
	if evt, _ = c.GuildCreate(data); evt != nil {
		evt = evt.DeepCopy().(*GuildCreate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) GuildDelete(data []byte) (evt *GuildDelete, err error) {
	if evt, _ = c.GuildDelete(data); evt != nil {
		evt = evt.DeepCopy().(*GuildDelete)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) GuildEmojisUpdate(data []byte) (evt *GuildEmojisUpdate, err error) {
	if evt, _ = c.GuildEmojisUpdate(data); evt != nil {
		evt = evt.DeepCopy().(*GuildEmojisUpdate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) GuildIntegrationsUpdate(data []byte) (evt *GuildIntegrationsUpdate, err error) {
	if evt, _ = c.GuildIntegrationsUpdate(data); evt != nil {
		evt = evt.DeepCopy().(*GuildIntegrationsUpdate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) GuildMemberAdd(data []byte) (evt *GuildMemberAdd, err error) {
	if evt, _ = c.GuildMemberAdd(data); evt != nil {
		evt = evt.DeepCopy().(*GuildMemberAdd)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) GuildMemberRemove(data []byte) (evt *GuildMemberRemove, err error) {
	if evt, _ = c.GuildMemberRemove(data); evt != nil {
		evt = evt.DeepCopy().(*GuildMemberRemove)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) GuildMemberUpdate(data []byte) (evt *GuildMemberUpdate, err error) {
	if evt, _ = c.GuildMemberUpdate(data); evt != nil {
		evt = evt.DeepCopy().(*GuildMemberUpdate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) GuildMembersChunk(data []byte) (evt *GuildMembersChunk, err error) {
	if evt, _ = c.GuildMembersChunk(data); evt != nil {
		evt = evt.DeepCopy().(*GuildMembersChunk)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) GuildRoleCreate(data []byte) (evt *GuildRoleCreate, err error) {
	if evt, _ = c.GuildRoleCreate(data); evt != nil {
		evt = evt.DeepCopy().(*GuildRoleCreate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) GuildRoleDelete(data []byte) (evt *GuildRoleDelete, err error) {
	if evt, _ = c.GuildRoleDelete(data); evt != nil {
		evt = evt.DeepCopy().(*GuildRoleDelete)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) GuildRoleUpdate(data []byte) (evt *GuildRoleUpdate, err error) {
	if evt, _ = c.GuildRoleUpdate(data); evt != nil {
		evt = evt.DeepCopy().(*GuildRoleUpdate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) GuildUpdate(data []byte) (evt *GuildUpdate, err error) {
	if evt, _ = c.GuildUpdate(data); evt != nil {
		evt = evt.DeepCopy().(*GuildUpdate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) InviteCreate(data []byte) (evt *InviteCreate, err error) {
	if evt, _ = c.InviteCreate(data); evt != nil {
		evt = evt.DeepCopy().(*InviteCreate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) InviteDelete(data []byte) (evt *InviteDelete, err error) {
	if evt, _ = c.InviteDelete(data); evt != nil {
		evt = evt.DeepCopy().(*InviteDelete)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) MessageCreate(data []byte) (evt *MessageCreate, err error) {
	if evt, _ = c.MessageCreate(data); evt != nil {
		evt = evt.DeepCopy().(*MessageCreate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) MessageDelete(data []byte) (evt *MessageDelete, err error) {
	if evt, _ = c.MessageDelete(data); evt != nil {
		evt = evt.DeepCopy().(*MessageDelete)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) MessageDeleteBulk(data []byte) (evt *MessageDeleteBulk, err error) {
	if evt, _ = c.MessageDeleteBulk(data); evt != nil {
		evt = evt.DeepCopy().(*MessageDeleteBulk)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) MessageReactionAdd(data []byte) (evt *MessageReactionAdd, err error) {
	if evt, _ = c.MessageReactionAdd(data); evt != nil {
		evt = evt.DeepCopy().(*MessageReactionAdd)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) MessageReactionRemove(data []byte) (evt *MessageReactionRemove, err error) {
	if evt, _ = c.MessageReactionRemove(data); evt != nil {
		evt = evt.DeepCopy().(*MessageReactionRemove)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) MessageReactionRemoveAll(data []byte) (evt *MessageReactionRemoveAll, err error) {
	if evt, _ = c.MessageReactionRemoveAll(data); evt != nil {
		evt = evt.DeepCopy().(*MessageReactionRemoveAll)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) MessageUpdate(data []byte) (evt *MessageUpdate, err error) {
	if evt, _ = c.MessageUpdate(data); evt != nil {
		evt = evt.DeepCopy().(*MessageUpdate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) PresenceUpdate(data []byte) (evt *PresenceUpdate, err error) {
	if evt, _ = c.PresenceUpdate(data); evt != nil {
		evt = evt.DeepCopy().(*PresenceUpdate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) Ready(data []byte) (evt *Ready, err error) {
	if evt, _ = c.Ready(data); evt != nil {
		evt = evt.DeepCopy().(*Ready)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) Resumed(data []byte) (evt *Resumed, err error) {
	if evt, _ = c.Resumed(data); evt != nil {
		evt = evt.DeepCopy().(*Resumed)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) TypingStart(data []byte) (evt *TypingStart, err error) {
	if evt, _ = c.TypingStart(data); evt != nil {
		evt = evt.DeepCopy().(*TypingStart)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) UserUpdate(data []byte) (evt *UserUpdate, err error) {
	if evt, _ = c.UserUpdate(data); evt != nil {
		evt = evt.DeepCopy().(*UserUpdate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) VoiceServerUpdate(data []byte) (evt *VoiceServerUpdate, err error) {
	if evt, _ = c.VoiceServerUpdate(data); evt != nil {
		evt = evt.DeepCopy().(*VoiceServerUpdate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) VoiceStateUpdate(data []byte) (evt *VoiceStateUpdate, err error) {
	if evt, _ = c.VoiceStateUpdate(data); evt != nil {
		evt = evt.DeepCopy().(*VoiceStateUpdate)
		return evt, nil
	}
	return nil, nil
}
func (c *CacheImmutable) WebhooksUpdate(data []byte) (evt *WebhooksUpdate, err error) {
	if evt, _ = c.WebhooksUpdate(data); evt != nil {
		evt = evt.DeepCopy().(*WebhooksUpdate)
		return evt, nil
	}
	return nil, nil
}

// REST lookup
func (c *CacheImmutable) GetMessage(channelID, messageID Snowflake) (*Message, error) {
	return nil, nil
}
func (c *CacheImmutable) GetChannel(id Snowflake) (*Channel, error)                { return nil, nil }
func (c *CacheImmutable) GetGuildEmoji(guildID, emojiID Snowflake) (*Emoji, error) { return nil, nil }
func (c *CacheImmutable) GetGuildEmojis(id Snowflake) ([]*Emoji, error)            { return nil, nil }
func (c *CacheImmutable) GetGuild(id Snowflake) (*Guild, error)                    { return nil, nil }
func (c *CacheImmutable) GetGuildChannels(id Snowflake) ([]*Channel, error)        { return nil, nil }
func (c *CacheImmutable) GetMember(guildID, userID Snowflake) (*Member, error)     { return nil, nil }
func (c *CacheImmutable) GetGuildRoles(guildID Snowflake) ([]*Role, error)         { return nil, nil }
func (c *CacheImmutable) GetCurrentUser() (*User, error)                           { return nil, nil }
func (c *CacheImmutable) GetUser(id Snowflake) (*User, error)                      { return nil, nil }
func (c *CacheImmutable) GetCurrentUserGuilds(p *GetCurrentUserGuildsParams) ([]*PartialGuild, error) {
	return nil, nil
}
func (c *CacheImmutable) GetMessages(channel Snowflake, p *GetMessagesParams) ([]*Message, error) {
	return nil, nil
}
func (c *CacheImmutable) GetMembers(guildID Snowflake, p *GetMembersParams) ([]*Member, error) {
	return nil, nil
}
