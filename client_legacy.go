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
func (c *Client) ModifyChannel(id Snowflake, flags ...Flag) *updateChannelBuilder {
	return c.UpdateChannel(id, flags...)
}

// Deprecated: use DeleteChannel
func (c *Client) CloseChannel(id Snowflake, flags ...Flag) (*Channel, error) {
	return c.DeleteChannel(id, flags...)
}

// Deprecated: use DeleteMessages
func (c *Client) BulkDeleteMessages(id Snowflake, params *DeleteMessagesParams, flags ...Flag) error {
	return c.DeleteMessages(id, params, flags...)
}

// Deprecated: use UpdateMessage
func (c *Client) EditMessage(chanID, msgID Snowflake, flags ...Flag) *updateMessageBuilder {
	return c.UpdateMessage(chanID, msgID, flags...)
}

// Deprecated: use UpdateChannelPermissions
func (c *Client) EditChannelPermissions(channelID, overwriteID Snowflake, params *UpdateChannelPermissionsParams, flags ...Flag) error {
	return c.UpdateChannelPermissions(channelID, overwriteID, params, flags...)
}

// Deprecated: use PinMessage or PinMessageID
func (c *Client) AddPinnedChannelMessage(channelID, messageID Snowflake, flags ...Flag) (err error) {
	return c.PinMessageID(channelID, messageID, flags...)
}

// Deprecated: use UnpinMessage or UnpinMessageID
func (c *Client) DeletePinnedChannelMessage(channelID, messageID Snowflake, flags ...Flag) (err error) {
	return c.UnpinMessageID(channelID, messageID, flags...)
}

// Deprecated: use AddDMParticipant
func (c *Client) GroupDMAddRecipient(channelID Snowflake, recipient *GroupDMParticipant, flags ...Flag) (err error) {
	return c.AddDMParticipant(channelID, recipient, flags...)
}

// Deprecated: use KickParticipant
func (c *Client) GroupDMRemoveRecipient(channelID, userID Snowflake, flags ...Flag) error {
	return c.KickParticipant(channelID, userID, flags...)
}

// Deprecated: use GetGuildEmojis
func (c *Client) ListGuildEmojis(guildID Snowflake, flags ...Flag) ([]*Emoji, error) {
	return c.GetGuildEmojis(guildID, flags...)
}

// Deprecated: use UpdateGuildEmoji
func (c *Client) ModifyGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) *updateGuildEmojiBuilder {
	return c.UpdateGuildEmoji(guildID, emojiID, flags...)
}

// Deprecated: use UpdateGuild
func (c *Client) ModifyGuild(id Snowflake, flags ...Flag) *updateGuildBuilder {
	return c.UpdateGuild(id, flags...)
}

// Deprecated: use UpdateGuildChannelPositions
func (c *Client) ModifyGuildChannelPositions(id Snowflake, params []UpdateGuildChannelPositionsParams, flags ...Flag) error {
	return c.UpdateGuildChannelPositions(id, params, flags...)
}

// Deprecated: use GetGuildMembers
func (c *Client) ListGuildMembers(id, after Snowflake, limit int, flags ...Flag) ([]*Member, error) {
	return c.GetMembers(id, &GetMembersParams{
		After: after,
		Limit: uint32(limit),
	}, flags...)
}

// TODO: AddGuildMember => CreateGuildMember

// Deprecated: use UpdateGuildMember
func (c *Client) ModifyGuildMember(guildID, userID Snowflake, flags ...Flag) *updateGuildMemberBuilder {
	return c.UpdateGuildMember(guildID, userID, flags...)
}

// Deprecated: use SetCurrentUserNick
func (c *Client) ModifyCurrentUserNick(guildID Snowflake, nick string, flags ...Flag) (string, error) {
	return c.SetCurrentUserNick(guildID, nick, flags...)
}

// TODO: AddGuildMemberRole => UpdateGuildMember
// TODO: RemoveGuildMemberRole => UpdateGuildMember

// Deprecated: use KickMember
func (c *Client) RemoveGuildMember(guildID, userID Snowflake, flags ...Flag) error {
	return c.KickMember(guildID, userID, flags...)
}

// Deprecated: use UnbanMember
func (c *Client) RemoveGuildBan(guildID, userID Snowflake, flags ...Flag) error {
	return c.UnbanMember(guildID, userID, flags...)
}

// Deprecated: use UpdateGuildRolePositions
func (c *Client) ModifyGuildRolePositions(guildID Snowflake, params []UpdateGuildRolePositionsParams, flags ...Flag) (ret []*Role, err error) {
	return c.UpdateGuildRolePositions(guildID, params, flags...)
}

// Deprecated: use DeleteGuildRole
func (c *Client) RemoveGuildRole(guildID, roleID Snowflake, flags ...Flag) error {
	return c.DeleteGuildRole(guildID, roleID, flags...)
}

// Deprecated: use PruneMembers
func (c *Client) BeginGuildPrune(guildID Snowflake, nrOfDays int, flags ...Flag) error {
	return c.PruneMembers(guildID, nrOfDays, flags...)
}

// Deprecated: use EstimatePruneMembersCount
func (c *Client) GetGuildPruneCount(guildID Snowflake, nrOfDays int, flags ...Flag) (int, error) {
	return c.EstimatePruneMembersCount(guildID, nrOfDays, flags...)
}

// Deprecated: use UpdateGuildIntegration
func (c *Client) ModifyGuildIntegration(guildID, integrationID Snowflake, params *UpdateGuildIntegrationParams, flags ...Flag) error {
	return c.UpdateGuildIntegration(guildID, integrationID, params, flags...)
}

// Deprecated: use UpdateGuildEmbed
func (c *Client) ModifyGuildEmbed(guildID Snowflake, flags ...Flag) *updateGuildEmbedBuilder {
	return c.UpdateGuildEmbed(guildID, flags...)
}

// Deprecated: use UpdateCurrentUser
func (c *Client) ModifyCurrentUser(flags ...Flag) *updateCurrentUserBuilder {
	return c.UpdateCurrentUser(flags...)
}

// Deprecated: use LeaveGuild
func (c *Client) ListVoiceRegions(flags ...Flag) ([]*VoiceRegion, error) {
	return c.GetVoiceRegions(flags...)
}

// Deprecated: use UpdateWebhook
func (c *Client) ModifyWebhook(id Snowflake, flags ...Flag) *updateWebhookBuilder {
	return c.UpdateWebhook(id, flags...)
}

// Deprecated: use UpdateWebhookWithToken
func (c *Client) ModifyWebhookWithToken(newWebhook *Webhook, flags ...Flag) *updateWebhookBuilder {
	return c.UpdateWebhookWithToken(newWebhook.ID, newWebhook.Token, flags...)
}
