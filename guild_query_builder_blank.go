package disgord

import "context"

// This struct exists to allow for unit testing with the guild query builder.
type guildQueryBuilderBlank struct{}

// Defines the blank methods.
func (guildQueryBuilderBlank) Get(ctx context.Context, flags ...Flag) (guild *Guild, err error) {
	return nil, nil
}
func (guildQueryBuilderBlank) GetChannels(ctx context.Context, flags ...Flag) ([]*Channel, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) GetMembers(ctx context.Context, params *GetMembersParams, flags ...Flag) ([]*Member, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) Update(ctx context.Context, flags ...Flag) UpdateGuildBuilder {
	return nil
}
func (guildQueryBuilderBlank) Delete(ctx context.Context, flags ...Flag) error { return nil }
func (guildQueryBuilderBlank) CreateChannel(ctx context.Context, name string, params *CreateGuildChannelParams, flags ...Flag) (*Channel, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) UpdateChannelPositions(ctx context.Context, params []UpdateGuildChannelPositionsParams, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderBlank) GetMember(ctx context.Context, userID Snowflake, flags ...Flag) (*Member, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) AddMember(ctx context.Context, userID Snowflake, accessToken string, params *AddGuildMemberParams, flags ...Flag) (*Member, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) UpdateMember(ctx context.Context, userID Snowflake, flags ...Flag) UpdateGuildMemberBuilder {
	return nil
}
func (guildQueryBuilderBlank) SetCurrentUserNick(ctx context.Context, nick string, flags ...Flag) (newNick string, err error) {
	return "", nil
}
func (guildQueryBuilderBlank) AddMemberRole(ctx context.Context, userID, roleID Snowflake, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderBlank) KickMember(ctx context.Context, userID Snowflake, reason string, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderBlank) KickVoiceParticipant(ctx context.Context, userID Snowflake) error {
	return nil
}
func (guildQueryBuilderBlank) GetBans(ctx context.Context, flags ...Flag) ([]*Ban, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) GetBan(ctx context.Context, userID Snowflake, flags ...Flag) (*Ban, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) BanMember(ctx context.Context, userID Snowflake, params *BanMemberParams, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderBlank) UnbanMember(ctx context.Context, userID Snowflake, reason string, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderBlank) GetRoles(ctx context.Context, flags ...Flag) ([]*Role, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) GetMemberPermissions(ctx context.Context, userID Snowflake, flags ...Flag) (permissions PermissionBit, err error) {
	return 0, nil
}
func (guildQueryBuilderBlank) CreateRole(ctx context.Context, params *CreateGuildRoleParams, flags ...Flag) (*Role, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) UpdateRolePositions(ctx context.Context, guildID Snowflake, params []UpdateGuildRolePositionsParams, flags ...Flag) ([]*Role, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) UpdateRole(ctx context.Context, roleID Snowflake, flags ...Flag) (builder UpdateGuildRoleBuilder) {
	return nil
}
func (guildQueryBuilderBlank) DeleteRole(ctx context.Context, roleID Snowflake, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderBlank) EstimatePruneMembersCount(ctx context.Context, days int, flags ...Flag) (estimate int, err error) {
	return 0, nil
}
func (guildQueryBuilderBlank) PruneMembers(ctx context.Context, days int, reason string, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderBlank) GetVoiceRegions(ctx context.Context, flags ...Flag) ([]*VoiceRegion, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) GetInvites(ctx context.Context, flags ...Flag) ([]*Invite, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) GetIntegrations(ctx context.Context, flags ...Flag) ([]*Integration, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) CreateIntegration(ctx context.Context, params *CreateGuildIntegrationParams, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderBlank) UpdateIntegration(ctx context.Context, integrationID Snowflake, params *UpdateGuildIntegrationParams, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderBlank) DeleteIntegration(ctx context.Context, integrationID Snowflake, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderBlank) SyncIntegration(ctx context.Context, integrationID Snowflake, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderBlank) GetEmbed(ctx context.Context, flags ...Flag) (*GuildEmbed, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) UpdateEmbed(ctx context.Context, flags ...Flag) UpdateGuildEmbedBuilder {
	return nil
}
func (guildQueryBuilderBlank) GetVanityURL(ctx context.Context, flags ...Flag) (*PartialInvite, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) GetAuditLogs(ctx context.Context, flags ...Flag) GuildAuditLogsBuilder {
	return nil
}
func (guildQueryBuilderBlank) VoiceConnect(channelID Snowflake) (ret VoiceConnection, err error) {
	return nil, nil
}
func (guildQueryBuilderBlank) GetEmoji(ctx context.Context, emojiID Snowflake, flags ...Flag) (*Emoji, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) GetEmojis(ctx context.Context, flags ...Flag) ([]*Emoji, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) CreateEmoji(ctx context.Context, params *CreateGuildEmojiParams, flags ...Flag) (*Emoji, error) {
	return nil, nil
}
func (guildQueryBuilderBlank) UpdateEmoji(ctx context.Context, emojiID Snowflake, flags ...Flag) UpdateGuildEmojiBuilder {
	return nil
}
func (guildQueryBuilderBlank) DeleteEmoji(ctx context.Context, emojiID Snowflake, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderBlank) GetWebhooks(ctx context.Context, flags ...Flag) (ret []*Webhook, err error) {
	return nil, nil
}
