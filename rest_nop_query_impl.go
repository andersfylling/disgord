package disgord

import (
	"context"
)

// TODO: remove this and use the http.Client for testing

// GuildQueryBuilderNop for testing
type GuildQueryBuilderNop struct{}

var _ GuildQueryBuilder = (*GuildQueryBuilderNop)(nil)

func (g GuildQueryBuilderNop) WithContext(_ context.Context) GuildQueryBuilder {
	return g
}
func (GuildQueryBuilderNop) Get(flags ...Flag) (guild *Guild, err error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetChannels(flags ...Flag) ([]*Channel, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetMembers(params *GetMembersParams, flags ...Flag) ([]*Member, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) Update(flags ...Flag) UpdateGuildBuilder {
	return nil
}
func (GuildQueryBuilderNop) Delete(flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) CreateChannel(name string, params *CreateGuildChannelParams, flags ...Flag) (*Channel, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) UpdateChannelPositions(params []UpdateGuildChannelPositionsParams, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) CreateMember(userID Snowflake, accessToken string, params *AddGuildMemberParams, flags ...Flag) (*Member, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) SetCurrentUserNick(nick string, flags ...Flag) (newNick string, err error) {
	return "", nil
}
func (GuildQueryBuilderNop) KickVoiceParticipant(userID Snowflake) error {
	return nil
}
func (GuildQueryBuilderNop) GetBans(flags ...Flag) ([]*Ban, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetBan(userID Snowflake, flags ...Flag) (*Ban, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) UnbanUser(userID Snowflake, reason string, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) GetRoles(flags ...Flag) ([]*Role, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetMemberPermissions(userID Snowflake, flags ...Flag) (permissions PermissionBit, err error) {
	return 0, nil
}
func (GuildQueryBuilderNop) CreateRole(params *CreateGuildRoleParams, flags ...Flag) (*Role, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) UpdateRolePositions(params []UpdateGuildRolePositionsParams, flags ...Flag) ([]*Role, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) EstimatePruneMembersCount(days int, flags ...Flag) (estimate int, err error) {
	return 0, nil
}
func (GuildQueryBuilderNop) PruneMembers(days int, reason string, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) GetVoiceRegions(flags ...Flag) ([]*VoiceRegion, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetInvites(flags ...Flag) ([]*Invite, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetIntegrations(flags ...Flag) ([]*Integration, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) CreateIntegration(params *CreateGuildIntegrationParams, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) UpdateIntegration(integrationID Snowflake, params *UpdateGuildIntegrationParams, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) DeleteIntegration(integrationID Snowflake, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) SyncIntegration(integrationID Snowflake, flags ...Flag) error {
	return nil
}
func (GuildQueryBuilderNop) GetEmbed(flags ...Flag) (*GuildEmbed, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) UpdateEmbed(flags ...Flag) UpdateGuildEmbedBuilder {
	return nil
}
func (GuildQueryBuilderNop) GetVanityURL(flags ...Flag) (*PartialInvite, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetAuditLogs(flags ...Flag) GuildAuditLogsBuilder {
	return nil
}
func (GuildQueryBuilderNop) VoiceConnect(channelID Snowflake) (ret VoiceConnection, err error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetEmojis(flags ...Flag) ([]*Emoji, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) CreateEmoji(params *CreateGuildEmojiParams, flags ...Flag) (*Emoji, error) {
	return nil, nil
}
func (GuildQueryBuilderNop) GetWebhooks(flags ...Flag) (ret []*Webhook, err error) {
	return nil, nil
}
func (GuildQueryBuilderNop) Member(userID Snowflake) GuildMemberQueryBuilder {
	return nil
}
func (GuildQueryBuilderNop) Role(roleID Snowflake) GuildRoleQueryBuilder {
	return nil
}
func (GuildQueryBuilderNop) Emoji(emojiID Snowflake) GuildEmojiQueryBuilder {
	return nil
}

// CurrentUserQueryBuilderNop for testing
type CurrentUserQueryBuilderNop struct{}

var _ CurrentUserQueryBuilder = (*CurrentUserQueryBuilderNop)(nil)

func (c CurrentUserQueryBuilderNop) WithContext(_ context.Context) CurrentUserQueryBuilder {
	return &c
}
func (CurrentUserQueryBuilderNop) Get(_ ...Flag) (*User, error) {
	return nil, nil
}
func (CurrentUserQueryBuilderNop) Update(_ ...Flag) UpdateCurrentUserBuilder {
	return nil
}
func (CurrentUserQueryBuilderNop) GetGuilds(_ *GetCurrentUserGuildsParams, _ ...Flag) ([]*PartialGuild, error) {
	return nil, nil
}
func (CurrentUserQueryBuilderNop) GetDMChannels(_ ...Flag) ([]*Channel, error) {
	return nil, nil
}
func (CurrentUserQueryBuilderNop) LeaveGuild(_ Snowflake, _ ...Flag) error {
	return nil
}
func (CurrentUserQueryBuilderNop) CreateGroupDM(_ *CreateGroupDMParams, _ ...Flag) (*Channel, error) {
	return nil, nil
}
func (CurrentUserQueryBuilderNop) GetUserConnections(_ ...Flag) ([]*UserConnection, error) {
	return nil, nil
}

// UserQueryBuilderNop for testing
type UserQueryBuilderNop struct{}

var _ UserQueryBuilder = (*UserQueryBuilderNop)(nil)

func (u UserQueryBuilderNop) WithContext(_ context.Context) UserQueryBuilder {
	return u
}
func (UserQueryBuilderNop) Get(_ ...Flag) (*User, error) {
	return nil, nil
}
func (UserQueryBuilderNop) CreateDM(_ ...Flag) (*Channel, error) {
	return nil, nil
}
