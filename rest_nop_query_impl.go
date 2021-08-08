package disgord

import (
	"context"
)

// TODO: remove this and use the http.Client for testing

// guildQueryBuilderNop for testing
type guildQueryBuilderNop struct{}

var _ GuildQueryBuilder = (*guildQueryBuilderNop)(nil)

func (g guildQueryBuilderNop) WithContext(_ context.Context) GuildQueryBuilder {
	return g
}
func (g guildQueryBuilderNop) WithFlags(_ ...Flag) GuildQueryBuilder {
	return g
}
func (guildQueryBuilderNop) VoiceChannel(channelID Snowflake) VoiceChannelQueryBuilder {
	return nil
}
func (guildQueryBuilderNop) Get() (guild *Guild, err error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetChannels() ([]*Channel, error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetMembers(params *GetMembersParams) ([]*Member, error) {
	return nil, nil
}
func (guildQueryBuilderNop) UpdateBuilder() UpdateGuildBuilder {
	return nil
}
func (guildQueryBuilderNop) Delete() error {
	return nil
}
func (guildQueryBuilderNop) CreateChannel(name string, params *CreateGuildChannelParams) (*Channel, error) {
	return nil, nil
}
func (guildQueryBuilderNop) UpdateChannelPositions(params []UpdateGuildChannelPositionsParams) error {
	return nil
}
func (guildQueryBuilderNop) CreateMember(userID Snowflake, accessToken string, params *AddGuildMemberParams) (*Member, error) {
	return nil, nil
}
func (guildQueryBuilderNop) SetCurrentUserNick(nick string) (newNick string, err error) {
	return "", nil
}
func (guildQueryBuilderNop) KickVoiceParticipant(userID Snowflake) error {
	return nil
}
func (guildQueryBuilderNop) GetBans() ([]*Ban, error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetBan(userID Snowflake) (*Ban, error) {
	return nil, nil
}
func (guildQueryBuilderNop) UnbanUser(userID Snowflake, reason string) error {
	return nil
}
func (guildQueryBuilderNop) GetRoles() ([]*Role, error) {
	return nil, nil
}
func (guildQueryBuilderNop) CreateRole(params *CreateGuildRoleParams) (*Role, error) {
	return nil, nil
}
func (guildQueryBuilderNop) UpdateRolePositions(params []UpdateGuildRolePositionsParams) ([]*Role, error) {
	return nil, nil
}
func (guildQueryBuilderNop) EstimatePruneMembersCount(days int) (estimate int, err error) {
	return 0, nil
}
func (guildQueryBuilderNop) PruneMembers(days int, reason string) error {
	return nil
}
func (guildQueryBuilderNop) GetVoiceRegions() ([]*VoiceRegion, error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetInvites() ([]*Invite, error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetIntegrations() ([]*Integration, error) {
	return nil, nil
}
func (guildQueryBuilderNop) CreateIntegration(params *CreateGuildIntegrationParams) error {
	return nil
}
func (guildQueryBuilderNop) UpdateIntegration(integrationID Snowflake, params *UpdateGuildIntegrationParams) error {
	return nil
}
func (guildQueryBuilderNop) DeleteIntegration(integrationID Snowflake) error {
	return nil
}
func (guildQueryBuilderNop) SyncIntegration(integrationID Snowflake) error {
	return nil
}
func (guildQueryBuilderNop) GetEmbed() (*GuildEmbed, error) {
	return nil, nil
}
func (guildQueryBuilderNop) UpdateEmbedBuilder() UpdateGuildEmbedBuilder {
	return nil
}
func (guildQueryBuilderNop) GetVanityURL() (*PartialInvite, error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetAuditLogs() GuildAuditLogsBuilder {
	return nil
}
func (guildQueryBuilderNop) VoiceConnect(channelID Snowflake) (ret VoiceConnection, err error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetEmojis() ([]*Emoji, error) {
	return nil, nil
}
func (guildQueryBuilderNop) CreateEmoji(params *CreateGuildEmojiParams) (*Emoji, error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetWebhooks() (ret []*Webhook, err error) {
	return nil, nil
}
func (guildQueryBuilderNop) Member(userID Snowflake) GuildMemberQueryBuilder {
	return nil
}
func (guildQueryBuilderNop) Role(roleID Snowflake) GuildRoleQueryBuilder {
	return nil
}
func (guildQueryBuilderNop) Emoji(emojiID Snowflake) GuildEmojiQueryBuilder {
	return nil
}

// currentUserQueryBuilderNop for testing
type currentUserQueryBuilderNop struct{}

var _ CurrentUserQueryBuilder = (*currentUserQueryBuilderNop)(nil)

func (c currentUserQueryBuilderNop) WithContext(_ context.Context) CurrentUserQueryBuilder {
	return &c
}
func (g currentUserQueryBuilderNop) WithFlags(_ ...Flag) CurrentUserQueryBuilder {
	return g
}
func (currentUserQueryBuilderNop) Get() (*User, error) {
	return nil, nil
}
func (currentUserQueryBuilderNop) UpdateBuilder() UpdateCurrentUserBuilder {
	return nil
}
func (currentUserQueryBuilderNop) GetGuilds(_ *GetCurrentUserGuildsParams) ([]*Guild, error) {
	return nil, nil
}
func (currentUserQueryBuilderNop) LeaveGuild(_ Snowflake) error {
	return nil
}
func (currentUserQueryBuilderNop) CreateGroupDM(_ *CreateGroupDMParams) (*Channel, error) {
	return nil, nil
}
func (currentUserQueryBuilderNop) GetUserConnections() ([]*UserConnection, error) {
	return nil, nil
}

// userQueryBuilderNop for testing
type userQueryBuilderNop struct{}

var _ UserQueryBuilder = (*userQueryBuilderNop)(nil)

func (u userQueryBuilderNop) WithContext(_ context.Context) UserQueryBuilder {
	return u
}
func (g userQueryBuilderNop) WithFlags(_ ...Flag) UserQueryBuilder {
	return g
}
func (userQueryBuilderNop) Get() (*User, error) {
	return nil, nil
}
func (userQueryBuilderNop) CreateDM() (*Channel, error) {
	return nil, nil
}
