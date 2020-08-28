package disgord

import (
	"context"
)

// TODO: remove this and use the http.Client for testing

// GuildQueryBuilderNop for testing
type GuildQueryBuilderNop struct{}

var _ GuildQueryBuilder = (*GuildQueryBuilderNop)(nil)

func (GuildQueryBuilderNop) Get(ctx context.Context, flags ...Flag) (guild *Guild, err error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetChannels(ctx context.Context, flags ...Flag) ([]*Channel, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetMembers(ctx context.Context, params *GetMembersParams, flags ...Flag) ([]*Member, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) Update(ctx context.Context, flags ...Flag) UpdateGuildBuilder {
	return nil
}
func (GuildQueryBuilderNop) Delete(ctx context.Context, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) CreateChannel(ctx context.Context, name string, params *CreateGuildChannelParams, flags ...Flag) (*Channel, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) UpdateChannelPositions(ctx context.Context, params []UpdateGuildChannelPositionsParams, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) GetMember(ctx context.Context, userID Snowflake, flags ...Flag) (*Member, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) AddMember(ctx context.Context, userID Snowflake, accessToken string, params *AddGuildMemberParams, flags ...Flag) (*Member, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) UpdateMember(ctx context.Context, userID Snowflake, flags ...Flag) UpdateGuildMemberBuilder {
	return nil
}
func (GuildQueryBuilderNop) SetCurrentUserNick(ctx context.Context, nick string, flags ...Flag) (newNick string, err error) {
	return "", nil
}
func (GuildQueryBuilderNop) AddMemberRole(ctx context.Context, userID, roleID Snowflake, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) KickMember(ctx context.Context, userID Snowflake, reason string, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) KickVoiceParticipant(ctx context.Context, userID Snowflake) error {
	return nil
}
func (GuildQueryBuilderNop) GetBans(ctx context.Context, flags ...Flag) ([]*Ban, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetBan(ctx context.Context, userID Snowflake, flags ...Flag) (*Ban, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) BanMember(ctx context.Context, userID Snowflake, params *BanMemberParams, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) UnbanMember(ctx context.Context, userID Snowflake, reason string, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) GetRoles(ctx context.Context, flags ...Flag) ([]*Role, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetMemberPermissions(ctx context.Context, userID Snowflake, flags ...Flag) (permissions PermissionBit, err error) {
	return 0, nil
}
func (GuildQueryBuilderNop) CreateRole(ctx context.Context, params *CreateGuildRoleParams, flags ...Flag) (*Role, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) UpdateRolePositions(ctx context.Context, guildID Snowflake, params []UpdateGuildRolePositionsParams, flags ...Flag) ([]*Role, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) UpdateRole(ctx context.Context, roleID Snowflake, flags ...Flag) (builder UpdateGuildRoleBuilder) {
	return nil
}
func (GuildQueryBuilderNop) DeleteRole(ctx context.Context, roleID Snowflake, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) EstimatePruneMembersCount(ctx context.Context, days int, flags ...Flag) (estimate int, err error) {
	return 0, nil
}
func (GuildQueryBuilderNop) PruneMembers(ctx context.Context, days int, reason string, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) GetVoiceRegions(ctx context.Context, flags ...Flag) ([]*VoiceRegion, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetInvites(ctx context.Context, flags ...Flag) ([]*Invite, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetIntegrations(ctx context.Context, flags ...Flag) ([]*Integration, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) CreateIntegration(ctx context.Context, params *CreateGuildIntegrationParams, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) UpdateIntegration(ctx context.Context, integrationID Snowflake, params *UpdateGuildIntegrationParams, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) DeleteIntegration(ctx context.Context, integrationID Snowflake, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) SyncIntegration(ctx context.Context, integrationID Snowflake, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) GetEmbed(ctx context.Context, flags ...Flag) (*GuildEmbed, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) UpdateEmbed(ctx context.Context, flags ...Flag) UpdateGuildEmbedBuilder {
	return nil
}
func (GuildQueryBuilderNop) GetVanityURL(ctx context.Context, flags ...Flag) (*PartialInvite, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetAuditLogs(ctx context.Context, flags ...Flag) GuildAuditLogsBuilder {
	return nil
}
func (GuildQueryBuilderNop) VoiceConnect(channelID Snowflake) (ret VoiceConnection, err error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetEmoji(ctx context.Context, emojiID Snowflake, flags ...Flag) (*Emoji, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetEmojis(ctx context.Context, flags ...Flag) ([]*Emoji, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) CreateEmoji(ctx context.Context, params *CreateGuildEmojiParams, flags ...Flag) (*Emoji, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) UpdateEmoji(ctx context.Context, emojiID Snowflake, flags ...Flag) UpdateGuildEmojiBuilder {
	return nil
}
func (GuildQueryBuilderNop) DeleteEmoji(ctx context.Context, emojiID Snowflake, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) GetWebhooks(ctx context.Context, flags ...Flag) (ret []*Webhook, err error) {
	return nil, nil
}

// CurrentUserQueryBuilderNop for testing
type CurrentUserQueryBuilderNop struct{}

var _ CurrentUserQueryBuilder = (*CurrentUserQueryBuilderNop)(nil)

func (CurrentUserQueryBuilderNop) Get(_ context.Context, _ ...Flag) (*User, error) {
	return nil, nil
}
func (CurrentUserQueryBuilderNop) Update(_ context.Context, _ ...Flag) UpdateCurrentUserBuilder {
	return nil
}
func (CurrentUserQueryBuilderNop) GetGuilds(_ context.Context, _ *GetCurrentUserGuildsParams, _ ...Flag) ([]*PartialGuild, error) {
	return nil, nil
}
func (CurrentUserQueryBuilderNop) GetDMChannels(_ context.Context, _ ...Flag) ([]*Channel, error) {
	return nil, nil
}
func (CurrentUserQueryBuilderNop) LeaveGuild(_ context.Context, _ Snowflake, _ ...Flag) error {
	return nil
}
func (CurrentUserQueryBuilderNop) CreateGroupDM(_ context.Context, _ *CreateGroupDMParams, _ ...Flag) (*Channel, error) {
	return nil, nil
}
func (CurrentUserQueryBuilderNop) GetUserConnections(_ context.Context, _ ...Flag) ([]*UserConnection, error) {
	return nil, nil
}

// UserQueryBuilderNop for testing
type UserQueryBuilderNop struct{}

var _ UserQueryBuilder = (*UserQueryBuilderNop)(nil)

func (UserQueryBuilderNop) Get(_ context.Context, _ ...Flag) (*User, error) {
	return nil, nil
}
func (UserQueryBuilderNop) CreateDM(_ context.Context, _ ...Flag) (*Channel, error) {
	return nil, nil
}
