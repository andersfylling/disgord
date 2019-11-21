// +build disgord_legacy

package disgord

import "context"

//////////////////////////////////////////////////////
//
// Deprecated / Legacy supported REST methods
//
// The main REST method names are renamed by a standard
// operation (ModifyChannel => UpdateChannel,
//            EditMessage => UpdateMessage,
//            etc.).
//
// However, to avoid confusion when developing
// with the discord docs as aid; this file holds
// the original names for the REST methods.
// (This might be deleted in the future)
//
//////////////////////////////////////////////////////

// Deprecated: use UpdateChannel
func (c *Client) ModifyChannel(ctx context.Context, id Snowflake, flags ...Flag) *updateChannelBuilder {
	return c.UpdateChannel(ctx, id, flags...)
}

// Deprecated: use DeleteChannel
func (c *Client) CloseChannel(ctx context.Context, id Snowflake, flags ...Flag) (*Channel, error) {
	return c.DeleteChannel(ctx, id, flags...)
}

// Deprecated: use DeleteMessages
func (c *Client) BulkDeleteMessages(ctx context.Context, id Snowflake, params *DeleteMessagesParams, flags ...Flag) error {
	return c.DeleteMessages(ctx, id, params, flags...)
}

// Deprecated: use UpdateMessage
func (c *Client) EditMessage(ctx context.Context, chanID, msgID Snowflake, flags ...Flag) *updateMessageBuilder {
	return c.UpdateMessage(ctx, chanID, msgID, flags...)
}

// Deprecated: use UpdateChannelPermissions
func (c *Client) EditChannelPermissions(ctx context.Context, channelID, overwriteID Snowflake, params *UpdateChannelPermissionsParams, flags ...Flag) error {
	return c.UpdateChannelPermissions(ctx, channelID, overwriteID, params, flags...)
}

// Deprecated: use PinMessage or PinMessageID
func (c *Client) AddPinnedChannelMessage(ctx context.Context, channelID, messageID Snowflake, flags ...Flag) (err error) {
	return c.PinMessageID(ctx, channelID, messageID, flags...)
}

// Deprecated: use UnpinMessage or UnpinMessageID
func (c *Client) DeletePinnedChannelMessage(ctx context.Context, channelID, messageID Snowflake, flags ...Flag) (err error) {
	return c.UnpinMessageID(ctx, channelID, messageID, flags...)
}

// Deprecated: use AddDMParticipant
func (c *Client) GroupDMAddRecipient(ctx context.Context, channelID Snowflake, recipient *GroupDMParticipant, flags ...Flag) (err error) {
	return c.AddDMParticipant(ctx, channelID, recipient, flags...)
}

// Deprecated: use KickParticipant
func (c *Client) GroupDMRemoveRecipient(ctx context.Context, channelID, userID Snowflake, flags ...Flag) error {
	return c.KickParticipant(ctx, channelID, userID, flags...)
}

// Deprecated: use GetGuildEmojis
func (c *Client) ListGuildEmojis(ctx context.Context, guildID Snowflake, flags ...Flag) ([]*Emoji, error) {
	return c.GetGuildEmojis(ctx, guildID, flags...)
}

// Deprecated: use UpdateGuildEmoji
func (c *Client) ModifyGuildEmoji(ctx context.Context, guildID, emojiID Snowflake, flags ...Flag) *updateGuildEmojiBuilder {
	return c.UpdateGuildEmoji(ctx, guildID, emojiID, flags...)
}

// Deprecated: use UpdateGuild
func (c *Client) ModifyGuild(ctx context.Context, id Snowflake, flags ...Flag) *updateGuildBuilder {
	return c.UpdateGuild(ctx, id, flags...)
}

// Deprecated: use UpdateGuildChannelPositions
func (c *Client) ModifyGuildChannelPositions(ctx context.Context, id Snowflake, params []UpdateGuildChannelPositionsParams, flags ...Flag) error {
	return c.UpdateGuildChannelPositions(ctx, id, params, flags...)
}

// Deprecated: use GetGuildMembers
func (c *Client) ListGuildMembers(ctx context.Context, id, after Snowflake, limit int, flags ...Flag) ([]*Member, error) {
	return c.GetMembers(ctx, id, &GetMembersParams{
		After: after,
		Limit: uint32(limit),
	}, flags...)
}

// TODO: AddGuildMember => CreateGuildMember

// Deprecated: use UpdateGuildMember
func (c *Client) ModifyGuildMember(ctx context.Context, guildID, userID Snowflake, flags ...Flag) *updateGuildMemberBuilder {
	return c.UpdateGuildMember(ctx, guildID, userID, flags...)
}

// Deprecated: use SetCurrentUserNick
func (c *Client) ModifyCurrentUserNick(ctx context.Context, guildID Snowflake, nick string, flags ...Flag) (string, error) {
	return c.SetCurrentUserNick(ctx, guildID, nick, flags...)
}

// TODO: AddGuildMemberRole => UpdateGuildMember
// TODO: RemoveGuildMemberRole => UpdateGuildMember

// Deprecated: use KickMember
func (c *Client) RemoveGuildMember(ctx context.Context, guildID, userID Snowflake, reason string, flags ...Flag) error {
	return c.KickMember(guildID, userID, reason, flags...)
}

// Deprecated: use UnbanMember
func (c *Client) RemoveGuildBan(ctx context.Context, guildID, userID Snowflake, reason string, flags ...Flag) error {
	return c.UnbanMember(guildID, userID, reason, flags...)
}

// Deprecated: use UpdateGuildRolePositions
func (c *Client) ModifyGuildRolePositions(ctx context.Context, guildID Snowflake, params []UpdateGuildRolePositionsParams, flags ...Flag) (ret []*Role, err error) {
	return c.UpdateGuildRolePositions(ctx, guildID, params, flags...)
}

// Deprecated: use DeleteGuildRole
func (c *Client) RemoveGuildRole(ctx context.Context, guildID, roleID Snowflake, flags ...Flag) error {
	return c.DeleteGuildRole(ctx, guildID, roleID, flags...)
}

// Deprecated: use PruneMembers
func (c *Client) BeginGuildPrune(ctx context.Context, guildID Snowflake, nrOfDays int, reason string, flags ...Flag) error {
	return c.PruneMembers(guildID, nrOfDays, reason, flags...)
}

// Deprecated: use EstimatePruneMembersCount
func (c *Client) GetGuildPruneCount(ctx context.Context, guildID Snowflake, nrOfDays int, flags ...Flag) (int, error) {
	return c.EstimatePruneMembersCount(ctx, guildID, nrOfDays, flags...)
}

// Deprecated: use UpdateGuildIntegration
func (c *Client) ModifyGuildIntegration(ctx context.Context, guildID, integrationID Snowflake, params *UpdateGuildIntegrationParams, flags ...Flag) error {
	return c.UpdateGuildIntegration(ctx, guildID, integrationID, params, flags...)
}

// Deprecated: use UpdateGuildEmbed
func (c *Client) ModifyGuildEmbed(ctx context.Context, guildID Snowflake, flags ...Flag) *updateGuildEmbedBuilder {
	return c.UpdateGuildEmbed(ctx, guildID, flags...)
}

// Deprecated: use UpdateCurrentUser
func (c *Client) ModifyCurrentUser(ctx context.Context, flags ...Flag) *updateCurrentUserBuilder {
	return c.UpdateCurrentUser(ctx, flags...)
}

// Deprecated: use LeaveGuild
func (c *Client) ListVoiceRegions(ctx context.Context, flags ...Flag) ([]*VoiceRegion, error) {
	return c.GetVoiceRegions(ctx, flags...)
}

// Deprecated: use UpdateWebhook
func (c *Client) ModifyWebhook(ctx context.Context, id Snowflake, flags ...Flag) *updateWebhookBuilder {
	return c.UpdateWebhook(ctx, id, flags...)
}

// Deprecated: use UpdateWebhookWithToken
func (c *Client) ModifyWebhookWithToken(ctx context.Context, newWebhook *Webhook, flags ...Flag) *updateWebhookBuilder {
	return c.UpdateWebhookWithToken(ctx, newWebhook.ID, newWebhook.Token, flags...)
}
