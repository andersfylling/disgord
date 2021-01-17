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
func (guildQueryBuilderNop) VoiceChannel(channelID Snowflake) VoiceChannelQueryBuilder {
	return nil
}
func (guildQueryBuilderNop) Get(flags ...Flag) (guild *Guild, err error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetChannels(flags ...Flag) ([]*Channel, error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetMembers(params *GetMembersParams, flags ...Flag) ([]*Member, error) {
	return nil, nil
}
func (guildQueryBuilderNop) UpdateBuilder(flags ...Flag) UpdateGuildBuilder {
	return nil
}
func (guildQueryBuilderNop) Delete(flags ...Flag) error {
	return nil
}
func (guildQueryBuilderNop) CreateChannel(name string, params *CreateGuildChannelParams, flags ...Flag) (*Channel, error) {
	return nil, nil
}
func (guildQueryBuilderNop) UpdateChannelPositions(params []UpdateGuildChannelPositionsParams, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderNop) CreateMember(userID Snowflake, accessToken string, params *AddGuildMemberParams, flags ...Flag) (*Member, error) {
	return nil, nil
}
func (guildQueryBuilderNop) SetCurrentUserNick(nick string, flags ...Flag) (newNick string, err error) {
	return "", nil
}
func (guildQueryBuilderNop) KickVoiceParticipant(userID Snowflake) error {
	return nil
}
func (guildQueryBuilderNop) GetBans(flags ...Flag) ([]*Ban, error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetBan(userID Snowflake, flags ...Flag) (*Ban, error) {
	return nil, nil
}
func (guildQueryBuilderNop) UnbanUser(userID Snowflake, reason string, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderNop) GetRoles(flags ...Flag) ([]*Role, error) {
	return nil, nil
}
func (guildQueryBuilderNop) CreateRole(params *CreateGuildRoleParams, flags ...Flag) (*Role, error) {
	return nil, nil
}
func (guildQueryBuilderNop) UpdateRolePositions(params []UpdateGuildRolePositionsParams, flags ...Flag) ([]*Role, error) {
	return nil, nil
}
func (guildQueryBuilderNop) EstimatePruneMembersCount(days int, flags ...Flag) (estimate int, err error) {
	return 0, nil
}
func (guildQueryBuilderNop) PruneMembers(days int, reason string, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderNop) GetVoiceRegions(flags ...Flag) ([]*VoiceRegion, error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetInvites(flags ...Flag) ([]*Invite, error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetIntegrations(flags ...Flag) ([]*Integration, error) {
	return nil, nil
}
func (guildQueryBuilderNop) CreateIntegration(params *CreateGuildIntegrationParams, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderNop) UpdateIntegration(integrationID Snowflake, params *UpdateGuildIntegrationParams, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderNop) DeleteIntegration(integrationID Snowflake, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderNop) SyncIntegration(integrationID Snowflake, flags ...Flag) error {
	return nil
}
func (guildQueryBuilderNop) GetEmbed(flags ...Flag) (*GuildEmbed, error) {
	return nil, nil
}
func (guildQueryBuilderNop) UpdateEmbedBuilder(flags ...Flag) UpdateGuildEmbedBuilder {
	return nil
}
func (guildQueryBuilderNop) GetVanityURL(flags ...Flag) (*PartialInvite, error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetAuditLogs(flags ...Flag) GuildAuditLogsBuilder {
	return nil
}
func (guildQueryBuilderNop) VoiceConnect(channelID Snowflake) (ret VoiceConnection, err error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetEmojis(flags ...Flag) ([]*Emoji, error) {
	return nil, nil
}
func (guildQueryBuilderNop) CreateEmoji(params *CreateGuildEmojiParams, flags ...Flag) (*Emoji, error) {
	return nil, nil
}
func (guildQueryBuilderNop) GetWebhooks(flags ...Flag) (ret []*Webhook, err error) {
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
func (currentUserQueryBuilderNop) Get(_ ...Flag) (*User, error) {
	return nil, nil
}
func (currentUserQueryBuilderNop) UpdateBuilder(_ ...Flag) UpdateCurrentUserBuilder {
	return nil
}
func (currentUserQueryBuilderNop) GetGuilds(_ *GetCurrentUserGuildsParams, _ ...Flag) ([]*Guild, error) {
	return nil, nil
}
func (currentUserQueryBuilderNop) LeaveGuild(_ Snowflake, _ ...Flag) error {
	return nil
}
func (currentUserQueryBuilderNop) CreateGroupDM(_ *CreateGroupDMParams, _ ...Flag) (*Channel, error) {
	return nil, nil
}
func (currentUserQueryBuilderNop) GetUserConnections(_ ...Flag) ([]*UserConnection, error) {
	return nil, nil
}

// userQueryBuilderNop for testing
type userQueryBuilderNop struct{}

var _ UserQueryBuilder = (*userQueryBuilderNop)(nil)

func (u userQueryBuilderNop) WithContext(_ context.Context) UserQueryBuilder {
	return u
}
func (userQueryBuilderNop) Get(_ ...Flag) (*User, error) {
	return nil, nil
}
func (userQueryBuilderNop) CreateDM(_ ...Flag) (*Channel, error) {
	return nil, nil
}
