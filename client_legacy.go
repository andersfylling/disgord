package disgord

//////////////////////////////////////////////////////
//
// Deprecated / Legacy supported REST methods
//
// As I want to keep the method names simple, I
// do not want to make it difficult to use
// DisGord side by side with the documentation.
//
// Below I've added all the REST methods with
// their names, as in, the discord documentation.
//
//////////////////////////////////////////////////////

// Deprecated: use UpdateChannel
func (c *client) ModifyChannel(id Snowflake, flags ...Flag) *updateChannelBuilder {
	return c.UpdateChannel(id, flags...)
}

// Deprecated: use DeleteChannel
func (c *client) CloseChannel(id Snowflake, flags ...Flag) (*Channel, error) {
	return c.DeleteChannel(id, flags...)
}

// Deprecated: use DeleteMessages
func (c *client) BulkDeleteMessages(id Snowflake, params *DeleteMessagesParams, flags ...Flag) error {
	return c.DeleteMessages(id, params, flags...)
}

// Deprecated: use UpdateMessage
func (c *client) EditMessage(chanID, msgID Snowflake, flags ...Flag) *updateMessageBuilder {
	return c.UpdateMessage(chanID, msgID, flags...)
}

// Deprecated: use UpdateChannelPermissions
func (c *client) EditChannelPermissions(channelID, overwriteID Snowflake, params *UpdateChannelPermissionsParams, flags ...Flag) error {
	return c.UpdateChannelPermissions(channelID, overwriteID, params, flags...)
}

// Deprecated: use PinMessage or PinMessageID
func (c *client) AddPinnedChannelMessage(channelID, messageID Snowflake, flags ...Flag) (err error) {
	return c.PinMessageID(channelID, messageID, flags...)
}

// Deprecated: use UnpinMessage or UnpinMessageID
func (c *client) DeletePinnedChannelMessage(channelID, messageID Snowflake, flags ...Flag) (err error) {
	return c.UnpinMessageID(channelID, messageID, flags...)
}

// Deprecated: use AddDMParticipant
func (c *client) GroupDMAddRecipient(channelID Snowflake, recipient *GroupDMParticipant, flags ...Flag) (err error) {
	return c.AddDMParticipant(channelID, recipient, flags...)
}

// Deprecated: use KickParticipant
func (c *client) GroupDMRemoveRecipient(channelID, userID Snowflake, flags ...Flag) error {
	return c.KickParticipant(channelID, userID, flags...)
}

// Deprecated: use GetGuildEmojis
func (c *client) ListGuildEmojis(guildID Snowflake, flags ...Flag) ([]*Emoji, error) {
	return c.GetGuildEmojis(guildID, flags...)
}

// Deprecated: use UpdateGuildEmoji
func (c *client) ModifyGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) *updateGuildEmojiBuilder {
	return c.UpdateGuildEmoji(guildID, emojiID, flags...)
}

// Deprecated: use UpdateGuild
func (c *client) ModifyGuild(id Snowflake, flags ...Flag) *updateGuildBuilder {
	return c.UpdateGuild(id, flags...)
}

// Deprecated: use UpdateGuildChannelPositions
func (c *client) ModifyGuildChannelPositions(id Snowflake, params []UpdateGuildChannelPositionsParams, flags ...Flag) error {
	return c.UpdateGuildChannelPositions(id, params, flags...)
}

// Deprecated: use GetGuildMembers
func (c *client) ListGuildMembers(id, after Snowflake, limit int, flags ...Flag) ([]*Member, error) {
	return c.GetMembers(id, &GetMembersParams{
		After: after,
		Limit: uint32(limit),
	}, flags...)
}

// TODO: AddGuildMember => CreateGuildMember

// Deprecated: use UpdateGuildMember
func (c *client) ModifyGuildMember(guildID, userID Snowflake, flags ...Flag) *updateGuildMemberBuilder {
	return c.UpdateGuildMember(guildID, userID, flags...)
}

// Deprecated: use SetCurrentUserNick
func (c *client) ModifyCurrentUserNick(guildID Snowflake, nick string, flags ...Flag) (string, error) {
	return c.SetCurrentUserNick(guildID, nick, flags...)
}

// TODO: AddGuildMemberRole => UpdateGuildMember
// TODO: RemoveGuildMemberRole => UpdateGuildMember

// Deprecated: use KickMember
func (c *client) RemoveGuildMember(guildID, userID Snowflake, flags ...Flag) error {
	return c.KickMember(guildID, userID, flags...)
}

// Deprecated: use UnbanMember
func (c *client) RemoveGuildBan(guildID, userID Snowflake, flags ...Flag) error {
	return c.UnbanMember(guildID, userID, flags...)
}

// Deprecated: use UpdateGuildRolePositions
func (c *client) ModifyGuildRolePositions(guildID Snowflake, params []UpdateGuildRolePositionsParams, flags ...Flag) (ret []*Role, err error) {
	return c.UpdateGuildRolePositions(guildID, params, flags...)
}

// Deprecated: use DeleteGuildRole
func (c *client) RemoveGuildRole(guildID, roleID Snowflake, flags ...Flag) error {
	return c.DeleteGuildRole(guildID, roleID, flags...)
}

// Deprecated: use PruneMembers
func (c *client) BeginGuildPrune(guildID Snowflake, nrOfDays int, flags ...Flag) error {
	return c.PruneMembers(guildID, nrOfDays, flags...)
}

// Deprecated: use EstimatePruneMembersCount
func (c *client) GetGuildPruneCount(guildID Snowflake, nrOfDays int, flags ...Flag) (int, error) {
	return c.EstimatePruneMembersCount(guildID, nrOfDays, flags...)
}

// Deprecated: use UpdateGuildIntegration
func (c *client) ModifyGuildIntegration(guildID, integrationID Snowflake, params *UpdateGuildIntegrationParams, flags ...Flag) error {
	return c.UpdateGuildIntegration(guildID, integrationID, params, flags...)
}

// Deprecated: use UpdateGuildEmbed
func (c *client) ModifyGuildEmbed(guildID Snowflake, flags ...Flag) *updateGuildEmbedBuilder {
	return c.UpdateGuildEmbed(guildID, flags...)
}

// Deprecated: use UpdateCurrentUser
func (c *client) ModifyCurrentUser(flags ...Flag) *updateCurrentUserBuilder {
	return c.UpdateCurrentUser(flags...)
}

// Deprecated: use LeaveGuild
func (c *client) ListVoiceRegions(flags ...Flag) ([]*VoiceRegion, error) {
	return c.GetVoiceRegions(flags...)
}

// Deprecated: use UpdateWebhook
func (c *client) ModifyWebhook(id Snowflake, flags ...Flag) *updateWebhookBuilder {
	return c.UpdateWebhook(id, flags...)
}

// Deprecated: use UpdateWebhookWithToken
func (c *client) ModifyWebhookWithToken(newWebhook *Webhook, flags ...Flag) *updateWebhookBuilder {
	return c.UpdateWebhookWithToken(newWebhook.ID, newWebhook.Token, flags...)
}
