// Code generated by generate/interfaces; DO NOT EDIT.

package disgord

import (
	"context"
	"github.com/andersfylling/disgord/internal/gateway"
	"net/url"
)

type channelQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ ChannelQueryBuilder = &channelQueryBuilderNop{}

func (c channelQueryBuilderNop) WithContext(ctx context.Context) ChannelQueryBuilder {
	c.Ctx = ctx
	return &c
}

func (c channelQueryBuilderNop) WithFlags(flags ...Flag) ChannelQueryBuilder {
	c.Flags = mergeFlags(flags)
	return &c
}

func (c *channelQueryBuilderNop) AddDMParticipant(_ *GroupDMParticipant) error {
	return nil
}

func (c *channelQueryBuilderNop) CreateInvite() CreateChannelInviteBuilder {
	return nil
}

func (c *channelQueryBuilderNop) CreateMessage(_ *CreateMessageParams) (*Message, error) {
	return nil, nil
}

func (c *channelQueryBuilderNop) CreateWebhook(_ *CreateWebhookParams) (*Webhook, error) {
	return nil, nil
}

func (c *channelQueryBuilderNop) Delete() (*Channel, error) {
	return nil, nil
}

func (c *channelQueryBuilderNop) DeleteMessages(_ *DeleteMessagesParams) error {
	return nil
}

func (c *channelQueryBuilderNop) DeletePermission(_ Snowflake) error {
	return nil
}

func (c *channelQueryBuilderNop) Get() (*Channel, error) {
	return nil, nil
}

func (c *channelQueryBuilderNop) GetInvites() ([]*Invite, error) {
	return nil, nil
}

func (c *channelQueryBuilderNop) GetMessages(_ *GetMessagesParams) ([]*Message, error) {
	return nil, nil
}

func (c *channelQueryBuilderNop) GetPinnedMessages() ([]*Message, error) {
	return nil, nil
}

func (c *channelQueryBuilderNop) GetWebhooks() ([]*Webhook, error) {
	return nil, nil
}

func (c *channelQueryBuilderNop) KickParticipant(_ Snowflake) error {
	return nil
}

func (c *channelQueryBuilderNop) Message(_ Snowflake) MessageQueryBuilder {
	return nil
}

func (c *channelQueryBuilderNop) TriggerTypingIndicator() error {
	return nil
}

func (c *channelQueryBuilderNop) UpdateBuilder() UpdateChannelBuilder {
	return nil
}

func (c *channelQueryBuilderNop) UpdatePermissions(_ Snowflake, _ *UpdateChannelPermissionsParams) error {
	return nil
}

type clientQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ ClientQueryBuilder = &clientQueryBuilderNop{}

func (c clientQueryBuilderNop) WithContext(ctx context.Context) ClientQueryBuilderExecutables {
	c.Ctx = ctx
	return &c
}

func (c clientQueryBuilderNop) WithFlags(flags ...Flag) ClientQueryBuilderExecutables {
	c.Flags = mergeFlags(flags)
	return &c
}

func (c *clientQueryBuilderNop) BotAuthorizeURL() (*url.URL, error) {
	return nil, nil
}

func (c *clientQueryBuilderNop) Channel(_ Snowflake) ChannelQueryBuilder {
	return nil
}

func (c *clientQueryBuilderNop) CreateGuild(_ string, _ *CreateGuildParams) (*Guild, error) {
	return nil, nil
}

func (c *clientQueryBuilderNop) CurrentUser() CurrentUserQueryBuilder {
	return nil
}

func (c *clientQueryBuilderNop) Gateway() GatewayQueryBuilder {
	return nil
}

func (c *clientQueryBuilderNop) GetVoiceRegions() ([]*VoiceRegion, error) {
	return nil, nil
}

func (c *clientQueryBuilderNop) Guild(_ Snowflake) GuildQueryBuilder {
	return nil
}

func (c *clientQueryBuilderNop) Invite(_ string) InviteQueryBuilder {
	return nil
}

func (c *clientQueryBuilderNop) SendMsg(_ Snowflake, _ []interface{}) (*Message, error) {
	return nil, nil
}

func (c *clientQueryBuilderNop) User(_ Snowflake) UserQueryBuilder {
	return nil
}

type currentUserQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ CurrentUserQueryBuilder = &currentUserQueryBuilderNop{}

func (c currentUserQueryBuilderNop) WithContext(ctx context.Context) CurrentUserQueryBuilder {
	c.Ctx = ctx
	return &c
}

func (c currentUserQueryBuilderNop) WithFlags(flags ...Flag) CurrentUserQueryBuilder {
	c.Flags = mergeFlags(flags)
	return &c
}

func (c *currentUserQueryBuilderNop) CreateGroupDM(_ *CreateGroupDMParams) (*Channel, error) {
	return nil, nil
}

func (c *currentUserQueryBuilderNop) Get() (*User, error) {
	return nil, nil
}

func (c *currentUserQueryBuilderNop) GetGuilds(_ *GetCurrentUserGuildsParams) ([]*Guild, error) {
	return nil, nil
}

func (c *currentUserQueryBuilderNop) GetUserConnections() ([]*UserConnection, error) {
	return nil, nil
}

func (c *currentUserQueryBuilderNop) LeaveGuild(_ Snowflake) error {
	return nil
}

func (c *currentUserQueryBuilderNop) UpdateBuilder() UpdateCurrentUserBuilder {
	return nil
}

type gatewayQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ GatewayQueryBuilder = &gatewayQueryBuilderNop{}

func (g gatewayQueryBuilderNop) WithContext(ctx context.Context) GatewayQueryBuilder {
	g.Ctx = ctx
	return &g
}

func (g *gatewayQueryBuilderNop) BotGuildsReady(_ func()) {
	return
}

func (g *gatewayQueryBuilderNop) BotReady(_ func()) {
	return
}

func (g *gatewayQueryBuilderNop) ChannelCreate(_ ChannelCreate, _ ChannelCreate) {
	return
}

func (g *gatewayQueryBuilderNop) ChannelCreateChan(_ chan *ChannelCreate, _ []chan *ChannelCreate) {
	return
}

func (g *gatewayQueryBuilderNop) ChannelDelete(_ ChannelDelete, _ ChannelDelete) {
	return
}

func (g *gatewayQueryBuilderNop) ChannelDeleteChan(_ chan *ChannelDelete, _ []chan *ChannelDelete) {
	return
}

func (g *gatewayQueryBuilderNop) ChannelPinsUpdate(_ ChannelPinsUpdate, _ ChannelPinsUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) ChannelPinsUpdateChan(_ chan *ChannelPinsUpdate, _ []chan *ChannelPinsUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) ChannelUpdate(_ ChannelUpdate, _ ChannelUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) ChannelUpdateChan(_ chan *ChannelUpdate, _ []chan *ChannelUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) Connect() error {
	return nil
}

func (g *gatewayQueryBuilderNop) Disconnect() error {
	return nil
}

func (g *gatewayQueryBuilderNop) DisconnectOnInterrupt() error {
	return nil
}

func (g *gatewayQueryBuilderNop) Dispatch(_ gatewayCmdName, _ CmdPayload) (Snowflake, error) {
	return 0, nil
}

func (g *gatewayQueryBuilderNop) Get() (*gateway.Gateway, error) {
	return nil, nil
}

func (g *gatewayQueryBuilderNop) GetBot() (*gateway.GatewayBot, error) {
	return nil, nil
}

func (g *gatewayQueryBuilderNop) GuildBanAdd(_ GuildBanAdd, _ GuildBanAdd) {
	return
}

func (g *gatewayQueryBuilderNop) GuildBanAddChan(_ chan *GuildBanAdd, _ []chan *GuildBanAdd) {
	return
}

func (g *gatewayQueryBuilderNop) GuildBanRemove(_ GuildBanRemove, _ GuildBanRemove) {
	return
}

func (g *gatewayQueryBuilderNop) GuildBanRemoveChan(_ chan *GuildBanRemove, _ []chan *GuildBanRemove) {
	return
}

func (g *gatewayQueryBuilderNop) GuildCreate(_ GuildCreate, _ GuildCreate) {
	return
}

func (g *gatewayQueryBuilderNop) GuildCreateChan(_ chan *GuildCreate, _ []chan *GuildCreate) {
	return
}

func (g *gatewayQueryBuilderNop) GuildDelete(_ GuildDelete, _ GuildDelete) {
	return
}

func (g *gatewayQueryBuilderNop) GuildDeleteChan(_ chan *GuildDelete, _ []chan *GuildDelete) {
	return
}

func (g *gatewayQueryBuilderNop) GuildEmojisUpdate(_ GuildEmojisUpdate, _ GuildEmojisUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) GuildEmojisUpdateChan(_ chan *GuildEmojisUpdate, _ []chan *GuildEmojisUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) GuildIntegrationsUpdate(_ GuildIntegrationsUpdate, _ GuildIntegrationsUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) GuildIntegrationsUpdateChan(_ chan *GuildIntegrationsUpdate, _ []chan *GuildIntegrationsUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) GuildMemberAdd(_ GuildMemberAdd, _ GuildMemberAdd) {
	return
}

func (g *gatewayQueryBuilderNop) GuildMemberAddChan(_ chan *GuildMemberAdd, _ []chan *GuildMemberAdd) {
	return
}

func (g *gatewayQueryBuilderNop) GuildMemberRemove(_ GuildMemberRemove, _ GuildMemberRemove) {
	return
}

func (g *gatewayQueryBuilderNop) GuildMemberRemoveChan(_ chan *GuildMemberRemove, _ []chan *GuildMemberRemove) {
	return
}

func (g *gatewayQueryBuilderNop) GuildMembersChunk(_ GuildMembersChunk, _ GuildMembersChunk) {
	return
}

func (g *gatewayQueryBuilderNop) GuildMembersChunkChan(_ chan *GuildMembersChunk, _ []chan *GuildMembersChunk) {
	return
}

func (g *gatewayQueryBuilderNop) GuildMemberUpdate(_ GuildMemberUpdate, _ GuildMemberUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) GuildMemberUpdateChan(_ chan *GuildMemberUpdate, _ []chan *GuildMemberUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) GuildRoleCreate(_ GuildRoleCreate, _ GuildRoleCreate) {
	return
}

func (g *gatewayQueryBuilderNop) GuildRoleCreateChan(_ chan *GuildRoleCreate, _ []chan *GuildRoleCreate) {
	return
}

func (g *gatewayQueryBuilderNop) GuildRoleDelete(_ GuildRoleDelete, _ GuildRoleDelete) {
	return
}

func (g *gatewayQueryBuilderNop) GuildRoleDeleteChan(_ chan *GuildRoleDelete, _ []chan *GuildRoleDelete) {
	return
}

func (g *gatewayQueryBuilderNop) GuildRoleUpdate(_ GuildRoleUpdate, _ GuildRoleUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) GuildRoleUpdateChan(_ chan *GuildRoleUpdate, _ []chan *GuildRoleUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) GuildUpdate(_ GuildUpdate, _ GuildUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) GuildUpdateChan(_ chan *GuildUpdate, _ []chan *GuildUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) InteractionCreate(_ InteractionCreate, _ InteractionCreate) {
	return
}

func (g *gatewayQueryBuilderNop) InteractionCreateChan(_ chan *InteractionCreate, _ []chan *InteractionCreate) {
	return
}

func (g *gatewayQueryBuilderNop) InviteCreate(_ InviteCreate, _ InviteCreate) {
	return
}

func (g *gatewayQueryBuilderNop) InviteCreateChan(_ chan *InviteCreate, _ []chan *InviteCreate) {
	return
}

func (g *gatewayQueryBuilderNop) InviteDelete(_ InviteDelete, _ InviteDelete) {
	return
}

func (g *gatewayQueryBuilderNop) InviteDeleteChan(_ chan *InviteDelete, _ []chan *InviteDelete) {
	return
}

func (g *gatewayQueryBuilderNop) MessageCreate(_ MessageCreate, _ MessageCreate) {
	return
}

func (g *gatewayQueryBuilderNop) MessageCreateChan(_ chan *MessageCreate, _ []chan *MessageCreate) {
	return
}

func (g *gatewayQueryBuilderNop) MessageDelete(_ MessageDelete, _ MessageDelete) {
	return
}

func (g *gatewayQueryBuilderNop) MessageDeleteBulk(_ MessageDeleteBulk, _ MessageDeleteBulk) {
	return
}

func (g *gatewayQueryBuilderNop) MessageDeleteBulkChan(_ chan *MessageDeleteBulk, _ []chan *MessageDeleteBulk) {
	return
}

func (g *gatewayQueryBuilderNop) MessageDeleteChan(_ chan *MessageDelete, _ []chan *MessageDelete) {
	return
}

func (g *gatewayQueryBuilderNop) MessageReactionAdd(_ MessageReactionAdd, _ MessageReactionAdd) {
	return
}

func (g *gatewayQueryBuilderNop) MessageReactionAddChan(_ chan *MessageReactionAdd, _ []chan *MessageReactionAdd) {
	return
}

func (g *gatewayQueryBuilderNop) MessageReactionRemove(_ MessageReactionRemove, _ MessageReactionRemove) {
	return
}

func (g *gatewayQueryBuilderNop) MessageReactionRemoveAll(_ MessageReactionRemoveAll, _ MessageReactionRemoveAll) {
	return
}

func (g *gatewayQueryBuilderNop) MessageReactionRemoveAllChan(_ chan *MessageReactionRemoveAll, _ []chan *MessageReactionRemoveAll) {
	return
}

func (g *gatewayQueryBuilderNop) MessageReactionRemoveChan(_ chan *MessageReactionRemove, _ []chan *MessageReactionRemove) {
	return
}

func (g *gatewayQueryBuilderNop) MessageReactionRemoveEmoji(_ MessageReactionRemoveEmoji, _ MessageReactionRemoveEmoji) {
	return
}

func (g *gatewayQueryBuilderNop) MessageReactionRemoveEmojiChan(_ chan *MessageReactionRemoveEmoji, _ []chan *MessageReactionRemoveEmoji) {
	return
}

func (g *gatewayQueryBuilderNop) MessageUpdate(_ MessageUpdate, _ MessageUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) MessageUpdateChan(_ chan *MessageUpdate, _ []chan *MessageUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) PresenceUpdate(_ PresenceUpdate, _ PresenceUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) PresenceUpdateChan(_ chan *PresenceUpdate, _ []chan *PresenceUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) Ready(_ Ready, _ Ready) {
	return
}

func (g *gatewayQueryBuilderNop) ReadyChan(_ chan *Ready, _ []chan *Ready) {
	return
}

func (g *gatewayQueryBuilderNop) Resumed(_ Resumed, _ Resumed) {
	return
}

func (g *gatewayQueryBuilderNop) ResumedChan(_ chan *Resumed, _ []chan *Resumed) {
	return
}

func (g *gatewayQueryBuilderNop) StayConnectedUntilInterrupted() error {
	return nil
}

func (g *gatewayQueryBuilderNop) TypingStart(_ TypingStart, _ TypingStart) {
	return
}

func (g *gatewayQueryBuilderNop) TypingStartChan(_ chan *TypingStart, _ []chan *TypingStart) {
	return
}

func (g *gatewayQueryBuilderNop) UserUpdate(_ UserUpdate, _ UserUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) UserUpdateChan(_ chan *UserUpdate, _ []chan *UserUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) VoiceServerUpdate(_ VoiceServerUpdate, _ VoiceServerUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) VoiceServerUpdateChan(_ chan *VoiceServerUpdate, _ []chan *VoiceServerUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) VoiceStateUpdate(_ VoiceStateUpdate, _ VoiceStateUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) VoiceStateUpdateChan(_ chan *VoiceStateUpdate, _ []chan *VoiceStateUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) WebhooksUpdate(_ WebhooksUpdate, _ WebhooksUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) WebhooksUpdateChan(_ chan *WebhooksUpdate, _ []chan *WebhooksUpdate) {
	return
}

func (g *gatewayQueryBuilderNop) WithCtrl(_ HandlerCtrl) SocketHandlerRegistrator {
	return nil
}

func (g *gatewayQueryBuilderNop) WithMiddleware(_ func(interface{}) interface{}, _ []func(interface{}) interface{}) SocketHandlerRegistrator {
	return nil
}

type guildEmojiQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ GuildEmojiQueryBuilder = &guildEmojiQueryBuilderNop{}

func (g guildEmojiQueryBuilderNop) WithContext(ctx context.Context) GuildEmojiQueryBuilder {
	g.Ctx = ctx
	return &g
}

func (g guildEmojiQueryBuilderNop) WithFlags(flags ...Flag) GuildEmojiQueryBuilder {
	g.Flags = mergeFlags(flags)
	return &g
}

func (g *guildEmojiQueryBuilderNop) Delete() error {
	return nil
}

func (g *guildEmojiQueryBuilderNop) Get() (*Emoji, error) {
	return nil, nil
}

func (g *guildEmojiQueryBuilderNop) UpdateBuilder() UpdateGuildEmojiBuilder {
	return nil
}

type guildMemberQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ GuildMemberQueryBuilder = &guildMemberQueryBuilderNop{}

func (g guildMemberQueryBuilderNop) WithContext(ctx context.Context) GuildMemberQueryBuilder {
	g.Ctx = ctx
	return &g
}

func (g guildMemberQueryBuilderNop) WithFlags(flags ...Flag) GuildMemberQueryBuilder {
	g.Flags = mergeFlags(flags)
	return &g
}

func (g *guildMemberQueryBuilderNop) AddRole(_ Snowflake) error {
	return nil
}

func (g *guildMemberQueryBuilderNop) Ban(_ *BanMemberParams) error {
	return nil
}

func (g *guildMemberQueryBuilderNop) Get() (*Member, error) {
	return nil, nil
}

func (g *guildMemberQueryBuilderNop) GetPermissions() (PermissionBit, error) {
	return 0, nil
}

func (g *guildMemberQueryBuilderNop) Kick(_ string) error {
	return nil
}

func (g *guildMemberQueryBuilderNop) RemoveRole(_ Snowflake) error {
	return nil
}

func (g *guildMemberQueryBuilderNop) UpdateBuilder() UpdateGuildMemberBuilder {
	return nil
}

type guildQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ GuildQueryBuilder = &guildQueryBuilderNop{}

func (g guildQueryBuilderNop) WithContext(ctx context.Context) GuildQueryBuilder {
	g.Ctx = ctx
	return &g
}

func (g guildQueryBuilderNop) WithFlags(flags ...Flag) GuildQueryBuilder {
	g.Flags = mergeFlags(flags)
	return &g
}

func (g *guildQueryBuilderNop) CreateChannel(_ string, _ *CreateGuildChannelParams) (*Channel, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) CreateEmoji(_ *CreateGuildEmojiParams) (*Emoji, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) CreateIntegration(_ *CreateGuildIntegrationParams) error {
	return nil
}

func (g *guildQueryBuilderNop) CreateMember(_ Snowflake, _ string, _ *AddGuildMemberParams) (*Member, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) CreateRole(_ *CreateGuildRoleParams) (*Role, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) Delete() error {
	return nil
}

func (g *guildQueryBuilderNop) DeleteIntegration(_ Snowflake) error {
	return nil
}

func (g *guildQueryBuilderNop) Emoji(_ Snowflake) GuildEmojiQueryBuilder {
	return nil
}

func (g *guildQueryBuilderNop) EstimatePruneMembersCount(_ int) (int, error) {
	return 0, nil
}

func (g *guildQueryBuilderNop) Get() (*Guild, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) GetAuditLogs() GuildAuditLogsBuilder {
	return nil
}

func (g *guildQueryBuilderNop) GetBan(_ Snowflake) (*Ban, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) GetBans() ([]*Ban, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) GetChannels() ([]*Channel, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) GetEmbed() (*GuildEmbed, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) GetEmojis() ([]*Emoji, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) GetIntegrations() ([]*Integration, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) GetInvites() ([]*Invite, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) GetMembers(_ *GetMembersParams) ([]*Member, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) GetRoles() ([]*Role, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) GetVanityURL() (*Invite, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) GetVoiceRegions() ([]*VoiceRegion, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) GetWebhooks() ([]*Webhook, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) KickVoiceParticipant(_ Snowflake) error {
	return nil
}

func (g *guildQueryBuilderNop) Member(_ Snowflake) GuildMemberQueryBuilder {
	return nil
}

func (g *guildQueryBuilderNop) PruneMembers(_ int, _ string) error {
	return nil
}

func (g *guildQueryBuilderNop) Role(_ Snowflake) GuildRoleQueryBuilder {
	return nil
}

func (g *guildQueryBuilderNop) SetCurrentUserNick(_ string) (string, error) {
	return "", nil
}

func (g *guildQueryBuilderNop) SyncIntegration(_ Snowflake) error {
	return nil
}

func (g *guildQueryBuilderNop) UnbanUser(_ Snowflake, _ string) error {
	return nil
}

func (g *guildQueryBuilderNop) UpdateBuilder() UpdateGuildBuilder {
	return nil
}

func (g *guildQueryBuilderNop) UpdateChannelPositions(_ []UpdateGuildChannelPositionsParams) error {
	return nil
}

func (g *guildQueryBuilderNop) UpdateEmbedBuilder() UpdateGuildEmbedBuilder {
	return nil
}

func (g *guildQueryBuilderNop) UpdateIntegration(_ Snowflake, _ *UpdateGuildIntegrationParams) error {
	return nil
}

func (g *guildQueryBuilderNop) UpdateRolePositions(_ []UpdateGuildRolePositionsParams) ([]*Role, error) {
	return nil, nil
}

func (g *guildQueryBuilderNop) VoiceChannel(_ Snowflake) VoiceChannelQueryBuilder {
	return nil
}

type guildRoleQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ GuildRoleQueryBuilder = &guildRoleQueryBuilderNop{}

func (g guildRoleQueryBuilderNop) WithContext(ctx context.Context) GuildRoleQueryBuilder {
	g.Ctx = ctx
	return &g
}

func (g guildRoleQueryBuilderNop) WithFlags(flags ...Flag) GuildRoleQueryBuilder {
	g.Flags = mergeFlags(flags)
	return &g
}

func (g *guildRoleQueryBuilderNop) Delete() error {
	return nil
}

func (g *guildRoleQueryBuilderNop) UpdateBuilder() UpdateGuildRoleBuilder {
	return nil
}

type inviteQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ InviteQueryBuilder = &inviteQueryBuilderNop{}

func (i inviteQueryBuilderNop) WithContext(ctx context.Context) InviteQueryBuilder {
	i.Ctx = ctx
	return &i
}

func (i inviteQueryBuilderNop) WithFlags(flags ...Flag) InviteQueryBuilder {
	i.Flags = mergeFlags(flags)
	return &i
}

func (i *inviteQueryBuilderNop) Delete() (*Invite, error) {
	return nil, nil
}

func (i *inviteQueryBuilderNop) Get(_ bool) (*Invite, error) {
	return nil, nil
}

type messageQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ MessageQueryBuilder = &messageQueryBuilderNop{}

func (m messageQueryBuilderNop) WithContext(ctx context.Context) MessageQueryBuilder {
	m.Ctx = ctx
	return &m
}

func (m messageQueryBuilderNop) WithFlags(flags ...Flag) MessageQueryBuilder {
	m.Flags = mergeFlags(flags)
	return &m
}

func (m *messageQueryBuilderNop) CrossPost() (*Message, error) {
	return nil, nil
}

func (m *messageQueryBuilderNop) Delete() error {
	return nil
}

func (m *messageQueryBuilderNop) DeleteAllReactions() error {
	return nil
}

func (m *messageQueryBuilderNop) Get() (*Message, error) {
	return nil, nil
}

func (m *messageQueryBuilderNop) Pin() error {
	return nil
}

func (m *messageQueryBuilderNop) Reaction(_ interface{}) ReactionQueryBuilder {
	return nil
}

func (m *messageQueryBuilderNop) SetContent(_ string) (*Message, error) {
	return nil, nil
}

func (m *messageQueryBuilderNop) SetEmbed(_ *Embed) (*Message, error) {
	return nil, nil
}

func (m *messageQueryBuilderNop) Unpin() error {
	return nil
}

func (m *messageQueryBuilderNop) UpdateBuilder() UpdateMessageBuilder {
	return nil
}

type reactionQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ ReactionQueryBuilder = &reactionQueryBuilderNop{}

func (r reactionQueryBuilderNop) WithContext(ctx context.Context) ReactionQueryBuilder {
	r.Ctx = ctx
	return &r
}

func (r reactionQueryBuilderNop) WithFlags(flags ...Flag) ReactionQueryBuilder {
	r.Flags = mergeFlags(flags)
	return &r
}

func (r *reactionQueryBuilderNop) Create() error {
	return nil
}

func (r *reactionQueryBuilderNop) DeleteOwn() error {
	return nil
}

func (r *reactionQueryBuilderNop) DeleteUser(_ Snowflake) error {
	return nil
}

func (r *reactionQueryBuilderNop) Get(_ URLQueryStringer) ([]*User, error) {
	return nil, nil
}

type userQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ UserQueryBuilder = &userQueryBuilderNop{}

func (u userQueryBuilderNop) WithContext(ctx context.Context) UserQueryBuilder {
	u.Ctx = ctx
	return &u
}

func (u userQueryBuilderNop) WithFlags(flags ...Flag) UserQueryBuilder {
	u.Flags = mergeFlags(flags)
	return &u
}

func (u *userQueryBuilderNop) CreateDM() (*Channel, error) {
	return nil, nil
}

func (u *userQueryBuilderNop) Get() (*User, error) {
	return nil, nil
}

type voiceChannelQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ VoiceChannelQueryBuilder = &voiceChannelQueryBuilderNop{}

func (v voiceChannelQueryBuilderNop) WithContext(ctx context.Context) ChannelQueryBuilder {
	v.Ctx = ctx
	return &v
}

func (v voiceChannelQueryBuilderNop) WithFlags(flags ...Flag) ChannelQueryBuilder {
	v.Flags = mergeFlags(flags)
	return &v
}

func (v *voiceChannelQueryBuilderNop) Connect(_ bool, _ bool) (VoiceConnection, error) {
	return nil, nil
}

func (v *voiceChannelQueryBuilderNop) CreateInvite() CreateChannelInviteBuilder {
	return nil
}

func (v *voiceChannelQueryBuilderNop) Delete() (*Channel, error) {
	return nil, nil
}

func (v *voiceChannelQueryBuilderNop) DeletePermission(_ Snowflake) error {
	return nil
}

func (v *voiceChannelQueryBuilderNop) Get() (*Channel, error) {
	return nil, nil
}

func (v *voiceChannelQueryBuilderNop) GetInvites() ([]*Invite, error) {
	return nil, nil
}

func (v *voiceChannelQueryBuilderNop) JoinManual(_ bool, _ bool) (*VoiceStateUpdate, *VoiceServerUpdate, error) {
	return nil, nil, nil
}

func (v *voiceChannelQueryBuilderNop) UpdateBuilder() UpdateChannelBuilder {
	return nil
}

func (v *voiceChannelQueryBuilderNop) UpdatePermissions(_ Snowflake, _ *UpdateChannelPermissionsParams) error {
	return nil
}

type webhookQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ WebhookQueryBuilder = &webhookQueryBuilderNop{}

func (w webhookQueryBuilderNop) WithContext(ctx context.Context) WebhookQueryBuilder {
	w.Ctx = ctx
	return &w
}

func (w webhookQueryBuilderNop) WithFlags(flags ...Flag) WebhookQueryBuilder {
	w.Flags = mergeFlags(flags)
	return &w
}

func (w *webhookQueryBuilderNop) Delete() error {
	return nil
}

func (w *webhookQueryBuilderNop) Execute(_ *ExecuteWebhookParams, _ bool, _ string) (*Message, error) {
	return nil, nil
}

func (w *webhookQueryBuilderNop) ExecuteGitHubWebhook(_ *ExecuteWebhookParams, _ bool) (*Message, error) {
	return nil, nil
}

func (w *webhookQueryBuilderNop) ExecuteSlackWebhook(_ *ExecuteWebhookParams, _ bool) (*Message, error) {
	return nil, nil
}

func (w *webhookQueryBuilderNop) Get() (*Webhook, error) {
	return nil, nil
}

func (w *webhookQueryBuilderNop) UpdateBuilder() UpdateWebhookBuilder {
	return nil
}

func (w *webhookQueryBuilderNop) WithToken(_ string) WebhookWithTokenQueryBuilder {
	return nil
}

type webhookWithTokenQueryBuilderNop struct {
	Ctx       context.Context
	Flags     Flag
	ChannelID Snowflake
	GuildID   Snowflake
	UserID    Snowflake
}

var _ WebhookWithTokenQueryBuilder = &webhookWithTokenQueryBuilderNop{}

func (w webhookWithTokenQueryBuilderNop) WithContext(ctx context.Context) WebhookWithTokenQueryBuilder {
	w.Ctx = ctx
	return &w
}

func (w webhookWithTokenQueryBuilderNop) WithFlags(flags ...Flag) WebhookWithTokenQueryBuilder {
	w.Flags = mergeFlags(flags)
	return &w
}

func (w *webhookWithTokenQueryBuilderNop) Delete() error {
	return nil
}

func (w *webhookWithTokenQueryBuilderNop) Execute(_ *ExecuteWebhookParams, _ bool, _ string) (*Message, error) {
	return nil, nil
}

func (w *webhookWithTokenQueryBuilderNop) Get() (*Webhook, error) {
	return nil, nil
}

func (w *webhookWithTokenQueryBuilderNop) UpdateBuilder() UpdateWebhookBuilder {
	return nil
}
