// Code generated by generate/interfaces; DO NOT EDIT.

package disgordutil

import (
	"context"
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/internal/gateway"
	"net/url"
)

func mergeFlags(flags []disgord.Flag) (f disgord.Flag) {
	for i := range flags {
		f |= flags[i]
	}

	return f
}

type ApplicationCommandQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.ApplicationCommandQueryBuilder = &ApplicationCommandQueryBuilderNop{}

func (a ApplicationCommandQueryBuilderNop) WithContext(ctx context.Context) disgord.ApplicationCommandQueryBuilder {
	a.Ctx = ctx
	return &a
}

func (a *ApplicationCommandQueryBuilderNop) Global() disgord.ApplicationCommandFunctions {
	return nil
}

func (a *ApplicationCommandQueryBuilderNop) Guild(_ disgord.Snowflake) disgord.ApplicationCommandFunctions {
	return nil
}

type ChannelQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.ChannelQueryBuilder = &ChannelQueryBuilderNop{}

func (c ChannelQueryBuilderNop) WithContext(ctx context.Context) disgord.ChannelQueryBuilder {
	c.Ctx = ctx
	return &c
}

func (c ChannelQueryBuilderNop) WithFlags(flags ...disgord.Flag) disgord.ChannelQueryBuilder {
	c.Flags = mergeFlags(flags)
	return &c
}

func (c *ChannelQueryBuilderNop) AddDMParticipant(_ *disgord.GroupDMParticipant) error {
	return nil
}

func (c *ChannelQueryBuilderNop) AddThreadMember(_ disgord.Snowflake) error {
	return nil
}

func (c *ChannelQueryBuilderNop) CreateInvite() disgord.CreateChannelInviteBuilder {
	return nil
}

func (c *ChannelQueryBuilderNop) CreateMessage(_ *disgord.CreateMessage) (*disgord.Message, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) CreateThread(_ *disgord.CreateThreadWithoutMessage) (*disgord.Channel, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) CreateWebhook(_ *disgord.CreateWebhook) (*disgord.Webhook, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) Delete() (*disgord.Channel, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) DeleteMessages(_ *disgord.DeleteMessages) error {
	return nil
}

func (c *ChannelQueryBuilderNop) DeletePermission(_ disgord.Snowflake) error {
	return nil
}

func (c *ChannelQueryBuilderNop) Get() (*disgord.Channel, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) GetInvites() ([]*disgord.Invite, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) GetJoinedPrivateArchivedThreads(_ *disgord.GetArchivedThreads) (*disgord.ArchivedThreads, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) GetMessages(_ *disgord.GetMessages) ([]*disgord.Message, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) GetPinnedMessages() ([]*disgord.Message, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) GetPrivateArchivedThreads(_ *disgord.GetArchivedThreads) (*disgord.ArchivedThreads, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) GetPublicArchivedThreads(_ *disgord.GetArchivedThreads) (*disgord.ArchivedThreads, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) GetThreadMember(_ disgord.Snowflake) (*disgord.ThreadMember, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) GetThreadMembers() ([]*disgord.ThreadMember, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) GetWebhooks() ([]*disgord.Webhook, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) JoinThread() error {
	return nil
}

func (c *ChannelQueryBuilderNop) KickParticipant(_ disgord.Snowflake) error {
	return nil
}

func (c *ChannelQueryBuilderNop) LeaveThread() error {
	return nil
}

func (c *ChannelQueryBuilderNop) Message(_ disgord.Snowflake) disgord.MessageQueryBuilder {
	return nil
}

func (c *ChannelQueryBuilderNop) RemoveThreadMember(_ disgord.Snowflake) error {
	return nil
}

func (c *ChannelQueryBuilderNop) TriggerTypingIndicator() error {
	return nil
}

func (c *ChannelQueryBuilderNop) Update(_ *disgord.UpdateChannel, _ string) (*disgord.Channel, error) {
	return nil, nil
}

func (c *ChannelQueryBuilderNop) UpdateBuilder() disgord.UpdateChannelBuilder {
	return nil
}

func (c *ChannelQueryBuilderNop) UpdatePermissions(_ disgord.Snowflake, _ *disgord.UpdateChannelPermissions) error {
	return nil
}

type ClientQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.ClientQueryBuilder = &ClientQueryBuilderNop{}

func (c ClientQueryBuilderNop) WithContext(ctx context.Context) disgord.ClientQueryBuilderExecutables {
	c.Ctx = ctx
	return &c
}

func (c ClientQueryBuilderNop) WithFlags(flags ...disgord.Flag) disgord.ClientQueryBuilderExecutables {
	c.Flags = mergeFlags(flags)
	return &c
}

func (c *ClientQueryBuilderNop) ApplicationCommand(_ disgord.Snowflake) disgord.ApplicationCommandQueryBuilder {
	return nil
}

func (c *ClientQueryBuilderNop) BotAuthorizeURL() (*url.URL, error) {
	return nil, nil
}

func (c *ClientQueryBuilderNop) Channel(_ disgord.Snowflake) disgord.ChannelQueryBuilder {
	return nil
}

func (c *ClientQueryBuilderNop) CreateGuild(_ string, _ *disgord.CreateGuild) (*disgord.Guild, error) {
	return nil, nil
}

func (c *ClientQueryBuilderNop) CurrentUser() disgord.CurrentUserQueryBuilder {
	return nil
}

func (c *ClientQueryBuilderNop) Gateway() disgord.GatewayQueryBuilder {
	return nil
}

func (c *ClientQueryBuilderNop) GetVoiceRegions() ([]*disgord.VoiceRegion, error) {
	return nil, nil
}

func (c *ClientQueryBuilderNop) Guild(_ disgord.Snowflake) disgord.GuildQueryBuilder {
	return nil
}

func (c *ClientQueryBuilderNop) Invite(_ string) disgord.InviteQueryBuilder {
	return nil
}

func (c *ClientQueryBuilderNop) SendMsg(_ disgord.Snowflake, _ ...interface{}) (*disgord.Message, error) {
	return nil, nil
}

func (c *ClientQueryBuilderNop) User(_ disgord.Snowflake) disgord.UserQueryBuilder {
	return nil
}

type CurrentUserQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.CurrentUserQueryBuilder = &CurrentUserQueryBuilderNop{}

func (c CurrentUserQueryBuilderNop) WithContext(ctx context.Context) disgord.CurrentUserQueryBuilder {
	c.Ctx = ctx
	return &c
}

func (c CurrentUserQueryBuilderNop) WithFlags(flags ...disgord.Flag) disgord.CurrentUserQueryBuilder {
	c.Flags = mergeFlags(flags)
	return &c
}

func (c *CurrentUserQueryBuilderNop) CreateGroupDM(_ *disgord.CreateGroupDM) (*disgord.Channel, error) {
	return nil, nil
}

func (c *CurrentUserQueryBuilderNop) Get() (*disgord.User, error) {
	return nil, nil
}

func (c *CurrentUserQueryBuilderNop) GetConnections() ([]*disgord.UserConnection, error) {
	return nil, nil
}

func (c *CurrentUserQueryBuilderNop) GetGuilds(_ *disgord.GetCurrentUserGuilds) ([]*disgord.Guild, error) {
	return nil, nil
}

func (c *CurrentUserQueryBuilderNop) GetUserConnections() ([]*disgord.UserConnection, error) {
	return nil, nil
}

func (c *CurrentUserQueryBuilderNop) LeaveGuild(_ disgord.Snowflake) error {
	return nil
}

func (c *CurrentUserQueryBuilderNop) Update(_ *disgord.UpdateUser) (*disgord.User, error) {
	return nil, nil
}

func (c *CurrentUserQueryBuilderNop) UpdateBuilder() disgord.UpdateCurrentUserBuilder {
	return nil
}

type GatewayQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.GatewayQueryBuilder = &GatewayQueryBuilderNop{}

func (g GatewayQueryBuilderNop) WithContext(ctx context.Context) disgord.GatewayQueryBuilder {
	g.Ctx = ctx
	return &g
}

func (g *GatewayQueryBuilderNop) BotGuildsReady(_ func()) {
	return
}

func (g *GatewayQueryBuilderNop) BotReady(_ func()) {
	return
}

func (g *GatewayQueryBuilderNop) ChannelCreate(_ func(disgord.Session, *disgord.ChannelCreate), _ ...func(disgord.Session, *disgord.ChannelCreate)) {
	return
}

func (g *GatewayQueryBuilderNop) ChannelCreateChan(_ chan *disgord.ChannelCreate, _ ...chan *disgord.ChannelCreate) {
	return
}

func (g *GatewayQueryBuilderNop) ChannelDelete(_ func(disgord.Session, *disgord.ChannelDelete), _ ...func(disgord.Session, *disgord.ChannelDelete)) {
	return
}

func (g *GatewayQueryBuilderNop) ChannelDeleteChan(_ chan *disgord.ChannelDelete, _ ...chan *disgord.ChannelDelete) {
	return
}

func (g *GatewayQueryBuilderNop) ChannelPinsUpdate(_ func(disgord.Session, *disgord.ChannelPinsUpdate), _ ...func(disgord.Session, *disgord.ChannelPinsUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) ChannelPinsUpdateChan(_ chan *disgord.ChannelPinsUpdate, _ ...chan *disgord.ChannelPinsUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) ChannelUpdate(_ func(disgord.Session, *disgord.ChannelUpdate), _ ...func(disgord.Session, *disgord.ChannelUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) ChannelUpdateChan(_ chan *disgord.ChannelUpdate, _ ...chan *disgord.ChannelUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) Connect() error {
	return nil
}

func (g *GatewayQueryBuilderNop) Disconnect() error {
	return nil
}

func (g *GatewayQueryBuilderNop) DisconnectOnInterrupt() error {
	return nil
}

func (g *GatewayQueryBuilderNop) Dispatch(_ disgord.GatewayCmdName, _ gateway.CmdPayload) ([]disgord.Snowflake, error) {
	return nil, nil
}

func (g *GatewayQueryBuilderNop) Get() (*gateway.Gateway, error) {
	return nil, nil
}

func (g *GatewayQueryBuilderNop) GetBot() (*gateway.GatewayBot, error) {
	return nil, nil
}

func (g *GatewayQueryBuilderNop) GuildBanAdd(_ func(disgord.Session, *disgord.GuildBanAdd), _ ...func(disgord.Session, *disgord.GuildBanAdd)) {
	return
}

func (g *GatewayQueryBuilderNop) GuildBanAddChan(_ chan *disgord.GuildBanAdd, _ ...chan *disgord.GuildBanAdd) {
	return
}

func (g *GatewayQueryBuilderNop) GuildBanRemove(_ func(disgord.Session, *disgord.GuildBanRemove), _ ...func(disgord.Session, *disgord.GuildBanRemove)) {
	return
}

func (g *GatewayQueryBuilderNop) GuildBanRemoveChan(_ chan *disgord.GuildBanRemove, _ ...chan *disgord.GuildBanRemove) {
	return
}

func (g *GatewayQueryBuilderNop) GuildCreate(_ func(disgord.Session, *disgord.GuildCreate), _ ...func(disgord.Session, *disgord.GuildCreate)) {
	return
}

func (g *GatewayQueryBuilderNop) GuildCreateChan(_ chan *disgord.GuildCreate, _ ...chan *disgord.GuildCreate) {
	return
}

func (g *GatewayQueryBuilderNop) GuildDelete(_ func(disgord.Session, *disgord.GuildDelete), _ ...func(disgord.Session, *disgord.GuildDelete)) {
	return
}

func (g *GatewayQueryBuilderNop) GuildDeleteChan(_ chan *disgord.GuildDelete, _ ...chan *disgord.GuildDelete) {
	return
}

func (g *GatewayQueryBuilderNop) GuildEmojisUpdate(_ func(disgord.Session, *disgord.GuildEmojisUpdate), _ ...func(disgord.Session, *disgord.GuildEmojisUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) GuildEmojisUpdateChan(_ chan *disgord.GuildEmojisUpdate, _ ...chan *disgord.GuildEmojisUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) GuildIntegrationsUpdate(_ func(disgord.Session, *disgord.GuildIntegrationsUpdate), _ ...func(disgord.Session, *disgord.GuildIntegrationsUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) GuildIntegrationsUpdateChan(_ chan *disgord.GuildIntegrationsUpdate, _ ...chan *disgord.GuildIntegrationsUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) GuildMemberAdd(_ func(disgord.Session, *disgord.GuildMemberAdd), _ ...func(disgord.Session, *disgord.GuildMemberAdd)) {
	return
}

func (g *GatewayQueryBuilderNop) GuildMemberAddChan(_ chan *disgord.GuildMemberAdd, _ ...chan *disgord.GuildMemberAdd) {
	return
}

func (g *GatewayQueryBuilderNop) GuildMemberRemove(_ func(disgord.Session, *disgord.GuildMemberRemove), _ ...func(disgord.Session, *disgord.GuildMemberRemove)) {
	return
}

func (g *GatewayQueryBuilderNop) GuildMemberRemoveChan(_ chan *disgord.GuildMemberRemove, _ ...chan *disgord.GuildMemberRemove) {
	return
}

func (g *GatewayQueryBuilderNop) GuildMembersChunk(_ func(disgord.Session, *disgord.GuildMembersChunk), _ ...func(disgord.Session, *disgord.GuildMembersChunk)) {
	return
}

func (g *GatewayQueryBuilderNop) GuildMembersChunkChan(_ chan *disgord.GuildMembersChunk, _ ...chan *disgord.GuildMembersChunk) {
	return
}

func (g *GatewayQueryBuilderNop) GuildMemberUpdate(_ func(disgord.Session, *disgord.GuildMemberUpdate), _ ...func(disgord.Session, *disgord.GuildMemberUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) GuildMemberUpdateChan(_ chan *disgord.GuildMemberUpdate, _ ...chan *disgord.GuildMemberUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) GuildRoleCreate(_ func(disgord.Session, *disgord.GuildRoleCreate), _ ...func(disgord.Session, *disgord.GuildRoleCreate)) {
	return
}

func (g *GatewayQueryBuilderNop) GuildRoleCreateChan(_ chan *disgord.GuildRoleCreate, _ ...chan *disgord.GuildRoleCreate) {
	return
}

func (g *GatewayQueryBuilderNop) GuildRoleDelete(_ func(disgord.Session, *disgord.GuildRoleDelete), _ ...func(disgord.Session, *disgord.GuildRoleDelete)) {
	return
}

func (g *GatewayQueryBuilderNop) GuildRoleDeleteChan(_ chan *disgord.GuildRoleDelete, _ ...chan *disgord.GuildRoleDelete) {
	return
}

func (g *GatewayQueryBuilderNop) GuildRoleUpdate(_ func(disgord.Session, *disgord.GuildRoleUpdate), _ ...func(disgord.Session, *disgord.GuildRoleUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) GuildRoleUpdateChan(_ chan *disgord.GuildRoleUpdate, _ ...chan *disgord.GuildRoleUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) GuildUpdate(_ func(disgord.Session, *disgord.GuildUpdate), _ ...func(disgord.Session, *disgord.GuildUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) GuildUpdateChan(_ chan *disgord.GuildUpdate, _ ...chan *disgord.GuildUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) InteractionCreate(_ func(disgord.Session, *disgord.InteractionCreate), _ ...func(disgord.Session, *disgord.InteractionCreate)) {
	return
}

func (g *GatewayQueryBuilderNop) InteractionCreateChan(_ chan *disgord.InteractionCreate, _ ...chan *disgord.InteractionCreate) {
	return
}

func (g *GatewayQueryBuilderNop) InviteCreate(_ func(disgord.Session, *disgord.InviteCreate), _ ...func(disgord.Session, *disgord.InviteCreate)) {
	return
}

func (g *GatewayQueryBuilderNop) InviteCreateChan(_ chan *disgord.InviteCreate, _ ...chan *disgord.InviteCreate) {
	return
}

func (g *GatewayQueryBuilderNop) InviteDelete(_ func(disgord.Session, *disgord.InviteDelete), _ ...func(disgord.Session, *disgord.InviteDelete)) {
	return
}

func (g *GatewayQueryBuilderNop) InviteDeleteChan(_ chan *disgord.InviteDelete, _ ...chan *disgord.InviteDelete) {
	return
}

func (g *GatewayQueryBuilderNop) MessageCreate(_ func(disgord.Session, *disgord.MessageCreate), _ ...func(disgord.Session, *disgord.MessageCreate)) {
	return
}

func (g *GatewayQueryBuilderNop) MessageCreateChan(_ chan *disgord.MessageCreate, _ ...chan *disgord.MessageCreate) {
	return
}

func (g *GatewayQueryBuilderNop) MessageDelete(_ func(disgord.Session, *disgord.MessageDelete), _ ...func(disgord.Session, *disgord.MessageDelete)) {
	return
}

func (g *GatewayQueryBuilderNop) MessageDeleteBulk(_ func(disgord.Session, *disgord.MessageDeleteBulk), _ ...func(disgord.Session, *disgord.MessageDeleteBulk)) {
	return
}

func (g *GatewayQueryBuilderNop) MessageDeleteBulkChan(_ chan *disgord.MessageDeleteBulk, _ ...chan *disgord.MessageDeleteBulk) {
	return
}

func (g *GatewayQueryBuilderNop) MessageDeleteChan(_ chan *disgord.MessageDelete, _ ...chan *disgord.MessageDelete) {
	return
}

func (g *GatewayQueryBuilderNop) MessageReactionAdd(_ func(disgord.Session, *disgord.MessageReactionAdd), _ ...func(disgord.Session, *disgord.MessageReactionAdd)) {
	return
}

func (g *GatewayQueryBuilderNop) MessageReactionAddChan(_ chan *disgord.MessageReactionAdd, _ ...chan *disgord.MessageReactionAdd) {
	return
}

func (g *GatewayQueryBuilderNop) MessageReactionRemove(_ func(disgord.Session, *disgord.MessageReactionRemove), _ ...func(disgord.Session, *disgord.MessageReactionRemove)) {
	return
}

func (g *GatewayQueryBuilderNop) MessageReactionRemoveAll(_ func(disgord.Session, *disgord.MessageReactionRemoveAll), _ ...func(disgord.Session, *disgord.MessageReactionRemoveAll)) {
	return
}

func (g *GatewayQueryBuilderNop) MessageReactionRemoveAllChan(_ chan *disgord.MessageReactionRemoveAll, _ ...chan *disgord.MessageReactionRemoveAll) {
	return
}

func (g *GatewayQueryBuilderNop) MessageReactionRemoveChan(_ chan *disgord.MessageReactionRemove, _ ...chan *disgord.MessageReactionRemove) {
	return
}

func (g *GatewayQueryBuilderNop) MessageReactionRemoveEmoji(_ func(disgord.Session, *disgord.MessageReactionRemoveEmoji), _ ...func(disgord.Session, *disgord.MessageReactionRemoveEmoji)) {
	return
}

func (g *GatewayQueryBuilderNop) MessageReactionRemoveEmojiChan(_ chan *disgord.MessageReactionRemoveEmoji, _ ...chan *disgord.MessageReactionRemoveEmoji) {
	return
}

func (g *GatewayQueryBuilderNop) MessageUpdate(_ func(disgord.Session, *disgord.MessageUpdate), _ ...func(disgord.Session, *disgord.MessageUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) MessageUpdateChan(_ chan *disgord.MessageUpdate, _ ...chan *disgord.MessageUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) PresenceUpdate(_ func(disgord.Session, *disgord.PresenceUpdate), _ ...func(disgord.Session, *disgord.PresenceUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) PresenceUpdateChan(_ chan *disgord.PresenceUpdate, _ ...chan *disgord.PresenceUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) Ready(_ func(disgord.Session, *disgord.Ready), _ ...func(disgord.Session, *disgord.Ready)) {
	return
}

func (g *GatewayQueryBuilderNop) ReadyChan(_ chan *disgord.Ready, _ ...chan *disgord.Ready) {
	return
}

func (g *GatewayQueryBuilderNop) Resumed(_ func(disgord.Session, *disgord.Resumed), _ ...func(disgord.Session, *disgord.Resumed)) {
	return
}

func (g *GatewayQueryBuilderNop) ResumedChan(_ chan *disgord.Resumed, _ ...chan *disgord.Resumed) {
	return
}

func (g *GatewayQueryBuilderNop) StayConnectedUntilInterrupted() error {
	return nil
}

func (g *GatewayQueryBuilderNop) ThreadCreate(_ func(disgord.Session, *disgord.ThreadCreate), _ ...func(disgord.Session, *disgord.ThreadCreate)) {
	return
}

func (g *GatewayQueryBuilderNop) ThreadCreateChan(_ chan *disgord.ThreadCreate, _ ...chan *disgord.ThreadCreate) {
	return
}

func (g *GatewayQueryBuilderNop) ThreadDelete(_ func(disgord.Session, *disgord.ThreadDelete), _ ...func(disgord.Session, *disgord.ThreadDelete)) {
	return
}

func (g *GatewayQueryBuilderNop) ThreadDeleteChan(_ chan *disgord.ThreadDelete, _ ...chan *disgord.ThreadDelete) {
	return
}

func (g *GatewayQueryBuilderNop) ThreadListSync(_ func(disgord.Session, *disgord.ThreadListSync), _ ...func(disgord.Session, *disgord.ThreadListSync)) {
	return
}

func (g *GatewayQueryBuilderNop) ThreadListSyncChan(_ chan *disgord.ThreadListSync, _ ...chan *disgord.ThreadListSync) {
	return
}

func (g *GatewayQueryBuilderNop) ThreadMembersUpdate(_ func(disgord.Session, *disgord.ThreadMembersUpdate), _ ...func(disgord.Session, *disgord.ThreadMembersUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) ThreadMembersUpdateChan(_ chan *disgord.ThreadMembersUpdate, _ ...chan *disgord.ThreadMembersUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) ThreadMemberUpdate(_ func(disgord.Session, *disgord.ThreadMemberUpdate), _ ...func(disgord.Session, *disgord.ThreadMemberUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) ThreadMemberUpdateChan(_ chan *disgord.ThreadMemberUpdate, _ ...chan *disgord.ThreadMemberUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) ThreadUpdate(_ func(disgord.Session, *disgord.ThreadUpdate), _ ...func(disgord.Session, *disgord.ThreadUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) ThreadUpdateChan(_ chan *disgord.ThreadUpdate, _ ...chan *disgord.ThreadUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) TypingStart(_ func(disgord.Session, *disgord.TypingStart), _ ...func(disgord.Session, *disgord.TypingStart)) {
	return
}

func (g *GatewayQueryBuilderNop) TypingStartChan(_ chan *disgord.TypingStart, _ ...chan *disgord.TypingStart) {
	return
}

func (g *GatewayQueryBuilderNop) UserUpdate(_ func(disgord.Session, *disgord.UserUpdate), _ ...func(disgord.Session, *disgord.UserUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) UserUpdateChan(_ chan *disgord.UserUpdate, _ ...chan *disgord.UserUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) VoiceServerUpdate(_ func(disgord.Session, *disgord.VoiceServerUpdate), _ ...func(disgord.Session, *disgord.VoiceServerUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) VoiceServerUpdateChan(_ chan *disgord.VoiceServerUpdate, _ ...chan *disgord.VoiceServerUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) VoiceStateUpdate(_ func(disgord.Session, *disgord.VoiceStateUpdate), _ ...func(disgord.Session, *disgord.VoiceStateUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) VoiceStateUpdateChan(_ chan *disgord.VoiceStateUpdate, _ ...chan *disgord.VoiceStateUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) WebhooksUpdate(_ func(disgord.Session, *disgord.WebhooksUpdate), _ ...func(disgord.Session, *disgord.WebhooksUpdate)) {
	return
}

func (g *GatewayQueryBuilderNop) WebhooksUpdateChan(_ chan *disgord.WebhooksUpdate, _ ...chan *disgord.WebhooksUpdate) {
	return
}

func (g *GatewayQueryBuilderNop) WithCtrl(_ disgord.HandlerCtrl) disgord.SocketHandlerRegistrator {
	return nil
}

func (g *GatewayQueryBuilderNop) WithMiddleware(_ func(interface{}) interface{}, _ ...func(interface{}) interface{}) disgord.SocketHandlerRegistrator {
	return nil
}

type GuildEmojiQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.GuildEmojiQueryBuilder = &GuildEmojiQueryBuilderNop{}

func (g GuildEmojiQueryBuilderNop) WithContext(ctx context.Context) disgord.GuildEmojiQueryBuilder {
	g.Ctx = ctx
	return &g
}

func (g GuildEmojiQueryBuilderNop) WithFlags(flags ...disgord.Flag) disgord.GuildEmojiQueryBuilder {
	g.Flags = mergeFlags(flags)
	return &g
}

func (g *GuildEmojiQueryBuilderNop) Delete() error {
	return nil
}

func (g *GuildEmojiQueryBuilderNop) Get() (*disgord.Emoji, error) {
	return nil, nil
}

func (g *GuildEmojiQueryBuilderNop) Update(_ *disgord.UpdateEmoji) (*disgord.Emoji, error) {
	return nil, nil
}

func (g *GuildEmojiQueryBuilderNop) UpdateBuilder() disgord.UpdateGuildEmojiBuilder {
	return nil
}

type GuildMemberQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.GuildMemberQueryBuilder = &GuildMemberQueryBuilderNop{}

func (g GuildMemberQueryBuilderNop) WithContext(ctx context.Context) disgord.GuildMemberQueryBuilder {
	g.Ctx = ctx
	return &g
}

func (g GuildMemberQueryBuilderNop) WithFlags(flags ...disgord.Flag) disgord.GuildMemberQueryBuilder {
	g.Flags = mergeFlags(flags)
	return &g
}

func (g *GuildMemberQueryBuilderNop) AddRole(_ disgord.Snowflake) error {
	return nil
}

func (g *GuildMemberQueryBuilderNop) Ban(_ *disgord.BanMember) error {
	return nil
}

func (g *GuildMemberQueryBuilderNop) Get() (*disgord.Member, error) {
	return nil, nil
}

func (g *GuildMemberQueryBuilderNop) GetPermissions() (disgord.PermissionBit, error) {
	return 0, nil
}

func (g *GuildMemberQueryBuilderNop) Kick(_ string) error {
	return nil
}

func (g *GuildMemberQueryBuilderNop) RemoveRole(_ disgord.Snowflake) error {
	return nil
}

func (g *GuildMemberQueryBuilderNop) Update(_ *disgord.UpdateMember) (*disgord.Member, error) {
	return nil, nil
}

func (g *GuildMemberQueryBuilderNop) UpdateBuilder() disgord.UpdateGuildMemberBuilder {
	return nil
}

type GuildQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.GuildQueryBuilder = &GuildQueryBuilderNop{}

func (g GuildQueryBuilderNop) WithContext(ctx context.Context) disgord.GuildQueryBuilder {
	g.Ctx = ctx
	return &g
}

func (g GuildQueryBuilderNop) WithFlags(flags ...disgord.Flag) disgord.GuildQueryBuilder {
	g.Flags = mergeFlags(flags)
	return &g
}

func (g *GuildQueryBuilderNop) CreateChannel(_ string, _ *disgord.CreateGuildChannel) (*disgord.Channel, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) CreateEmoji(_ *disgord.CreateGuildEmoji) (*disgord.Emoji, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) CreateIntegration(_ *disgord.CreateGuildIntegration) error {
	return nil
}

func (g *GuildQueryBuilderNop) CreateMember(_ disgord.Snowflake, _ string, _ *disgord.AddGuildMember) (*disgord.Member, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) CreateRole(_ *disgord.CreateGuildRole) (*disgord.Role, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) Delete() error {
	return nil
}

func (g *GuildQueryBuilderNop) DeleteIntegration(_ disgord.Snowflake) error {
	return nil
}

func (g *GuildQueryBuilderNop) DisconnectVoiceParticipant(_ disgord.Snowflake) error {
	return nil
}

func (g *GuildQueryBuilderNop) Emoji(_ disgord.Snowflake) disgord.GuildEmojiQueryBuilder {
	return nil
}

func (g *GuildQueryBuilderNop) EstimatePruneMembersCount(_ int) (int, error) {
	return 0, nil
}

func (g *GuildQueryBuilderNop) Get() (*disgord.Guild, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetActiveThreads() (*disgord.ActiveGuildThreads, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetAuditLogs(_ *disgord.GetAuditLogs) (*disgord.AuditLog, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetBan(_ disgord.Snowflake) (*disgord.Ban, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetBans() ([]*disgord.Ban, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetChannels() ([]*disgord.Channel, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetEmbed() (*disgord.GuildWidget, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetEmojis() ([]*disgord.Emoji, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetIntegrations() ([]*disgord.Integration, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetInvites() ([]*disgord.Invite, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetMembers(_ *disgord.GetMembers) ([]*disgord.Member, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetRoles() ([]*disgord.Role, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetVanityURL() (*disgord.Invite, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetVoiceRegions() ([]*disgord.VoiceRegion, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetWebhooks() ([]*disgord.Webhook, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) GetWidget() (*disgord.GuildWidget, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) KickVoiceParticipant(_ disgord.Snowflake) error {
	return nil
}

func (g *GuildQueryBuilderNop) Leave() error {
	return nil
}

func (g *GuildQueryBuilderNop) Member(_ disgord.Snowflake) disgord.GuildMemberQueryBuilder {
	return nil
}

func (g *GuildQueryBuilderNop) PruneMembers(_ int, _ string) error {
	return nil
}

func (g *GuildQueryBuilderNop) Role(_ disgord.Snowflake) disgord.GuildRoleQueryBuilder {
	return nil
}

func (g *GuildQueryBuilderNop) SetCurrentUserNick(_ string) (string, error) {
	return "", nil
}

func (g *GuildQueryBuilderNop) SyncIntegration(_ disgord.Snowflake) error {
	return nil
}

func (g *GuildQueryBuilderNop) UnbanUser(_ disgord.Snowflake, _ string) error {
	return nil
}

func (g *GuildQueryBuilderNop) Update(_ *disgord.UpdateGuild) (*disgord.Guild, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) UpdateBuilder() disgord.UpdateGuildBuilder {
	return nil
}

func (g *GuildQueryBuilderNop) UpdateChannelPositions(_ []disgord.UpdateGuildChannelPositions) error {
	return nil
}

func (g *GuildQueryBuilderNop) UpdateEmbedBuilder() disgord.UpdateGuildEmbedBuilder {
	return nil
}

func (g *GuildQueryBuilderNop) UpdateIntegration(_ disgord.Snowflake, _ *disgord.UpdateGuildIntegration) error {
	return nil
}

func (g *GuildQueryBuilderNop) UpdateRolePositions(_ []disgord.UpdateGuildRolePositions) ([]*disgord.Role, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) UpdateWidget(_ *disgord.UpdateGuildWidget) (*disgord.GuildWidget, error) {
	return nil, nil
}

func (g *GuildQueryBuilderNop) VoiceChannel(_ disgord.Snowflake) disgord.VoiceChannelQueryBuilder {
	return nil
}

type GuildRoleQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.GuildRoleQueryBuilder = &GuildRoleQueryBuilderNop{}

func (g GuildRoleQueryBuilderNop) WithContext(ctx context.Context) disgord.GuildRoleQueryBuilder {
	g.Ctx = ctx
	return &g
}

func (g GuildRoleQueryBuilderNop) WithFlags(flags ...disgord.Flag) disgord.GuildRoleQueryBuilder {
	g.Flags = mergeFlags(flags)
	return &g
}

func (g *GuildRoleQueryBuilderNop) Delete() error {
	return nil
}

func (g *GuildRoleQueryBuilderNop) Update(_ *disgord.UpdateRole, _ string) (*disgord.Role, error) {
	return nil, nil
}

func (g *GuildRoleQueryBuilderNop) UpdateBuilder() disgord.UpdateGuildRoleBuilder {
	return nil
}

type InviteQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.InviteQueryBuilder = &InviteQueryBuilderNop{}

func (i InviteQueryBuilderNop) WithContext(ctx context.Context) disgord.InviteQueryBuilder {
	i.Ctx = ctx
	return &i
}

func (i InviteQueryBuilderNop) WithFlags(flags ...disgord.Flag) disgord.InviteQueryBuilder {
	i.Flags = mergeFlags(flags)
	return &i
}

func (i *InviteQueryBuilderNop) Delete() (*disgord.Invite, error) {
	return nil, nil
}

func (i *InviteQueryBuilderNop) Get(_ bool) (*disgord.Invite, error) {
	return nil, nil
}

type MessageQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.MessageQueryBuilder = &MessageQueryBuilderNop{}

func (m MessageQueryBuilderNop) WithContext(ctx context.Context) disgord.MessageQueryBuilder {
	m.Ctx = ctx
	return &m
}

func (m MessageQueryBuilderNop) WithFlags(flags ...disgord.Flag) disgord.MessageQueryBuilder {
	m.Flags = mergeFlags(flags)
	return &m
}

func (m *MessageQueryBuilderNop) CreateThread(_ *disgord.CreateThread) (*disgord.Channel, error) {
	return nil, nil
}

func (m *MessageQueryBuilderNop) CrossPost() (*disgord.Message, error) {
	return nil, nil
}

func (m *MessageQueryBuilderNop) Delete() error {
	return nil
}

func (m *MessageQueryBuilderNop) DeleteAllReactions() error {
	return nil
}

func (m *MessageQueryBuilderNop) Get() (*disgord.Message, error) {
	return nil, nil
}

func (m *MessageQueryBuilderNop) Pin() error {
	return nil
}

func (m *MessageQueryBuilderNop) Reaction(_ interface{}) disgord.ReactionQueryBuilder {
	return nil
}

func (m *MessageQueryBuilderNop) SetContent(_ string) (*disgord.Message, error) {
	return nil, nil
}

func (m *MessageQueryBuilderNop) SetEmbed(_ *disgord.Embed) (*disgord.Message, error) {
	return nil, nil
}

func (m *MessageQueryBuilderNop) Unpin() error {
	return nil
}

func (m *MessageQueryBuilderNop) Update(_ *disgord.UpdateMessage) (*disgord.Message, error) {
	return nil, nil
}

func (m *MessageQueryBuilderNop) UpdateBuilder() disgord.UpdateMessageBuilder {
	return nil
}

type ReactionQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.ReactionQueryBuilder = &ReactionQueryBuilderNop{}

func (r ReactionQueryBuilderNop) WithContext(ctx context.Context) disgord.ReactionQueryBuilder {
	r.Ctx = ctx
	return &r
}

func (r ReactionQueryBuilderNop) WithFlags(flags ...disgord.Flag) disgord.ReactionQueryBuilder {
	r.Flags = mergeFlags(flags)
	return &r
}

func (r *ReactionQueryBuilderNop) Create() error {
	return nil
}

func (r *ReactionQueryBuilderNop) DeleteOwn() error {
	return nil
}

func (r *ReactionQueryBuilderNop) DeleteUser(_ disgord.Snowflake) error {
	return nil
}

func (r *ReactionQueryBuilderNop) Get(_ disgord.URLQueryStringer) ([]*disgord.User, error) {
	return nil, nil
}

type UserQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.UserQueryBuilder = &UserQueryBuilderNop{}

func (u UserQueryBuilderNop) WithContext(ctx context.Context) disgord.UserQueryBuilder {
	u.Ctx = ctx
	return &u
}

func (u UserQueryBuilderNop) WithFlags(flags ...disgord.Flag) disgord.UserQueryBuilder {
	u.Flags = mergeFlags(flags)
	return &u
}

func (u *UserQueryBuilderNop) CreateDM() (*disgord.Channel, error) {
	return nil, nil
}

func (u *UserQueryBuilderNop) Get() (*disgord.User, error) {
	return nil, nil
}

type VoiceChannelQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.VoiceChannelQueryBuilder = &VoiceChannelQueryBuilderNop{}

func (v VoiceChannelQueryBuilderNop) WithContext(ctx context.Context) disgord.VoiceChannelQueryBuilder {
	v.Ctx = ctx
	return &v
}

func (v VoiceChannelQueryBuilderNop) WithFlags(flags ...disgord.Flag) disgord.VoiceChannelQueryBuilder {
	v.Flags = mergeFlags(flags)
	return &v
}

func (v *VoiceChannelQueryBuilderNop) Connect(_ bool, _ bool) (disgord.VoiceConnection, error) {
	return nil, nil
}

func (v *VoiceChannelQueryBuilderNop) CreateInvite() disgord.CreateChannelInviteBuilder {
	return nil
}

func (v *VoiceChannelQueryBuilderNop) Delete() (*disgord.Channel, error) {
	return nil, nil
}

func (v *VoiceChannelQueryBuilderNop) DeletePermission(_ disgord.Snowflake) error {
	return nil
}

func (v *VoiceChannelQueryBuilderNop) Get() (*disgord.Channel, error) {
	return nil, nil
}

func (v *VoiceChannelQueryBuilderNop) GetInvites() ([]*disgord.Invite, error) {
	return nil, nil
}

func (v *VoiceChannelQueryBuilderNop) JoinManual(_ bool, _ bool) (*disgord.VoiceStateUpdate, *disgord.VoiceServerUpdate, error) {
	return nil, nil, nil
}

func (v *VoiceChannelQueryBuilderNop) Update(_ *disgord.UpdateChannel, _ string) (*disgord.Channel, error) {
	return nil, nil
}

func (v *VoiceChannelQueryBuilderNop) UpdateBuilder() disgord.UpdateChannelBuilder {
	return nil
}

func (v *VoiceChannelQueryBuilderNop) UpdatePermissions(_ disgord.Snowflake, _ *disgord.UpdateChannelPermissions) error {
	return nil
}

type WebhookQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.WebhookQueryBuilder = &WebhookQueryBuilderNop{}

func (w WebhookQueryBuilderNop) WithContext(ctx context.Context) disgord.WebhookQueryBuilder {
	w.Ctx = ctx
	return &w
}

func (w WebhookQueryBuilderNop) WithFlags(flags ...disgord.Flag) disgord.WebhookQueryBuilder {
	w.Flags = mergeFlags(flags)
	return &w
}

func (w *WebhookQueryBuilderNop) Delete() error {
	return nil
}

func (w *WebhookQueryBuilderNop) Execute(_ *disgord.ExecuteWebhook, _ bool, _ disgord.Snowflake, _ string) (*disgord.Message, error) {
	return nil, nil
}

func (w *WebhookQueryBuilderNop) ExecuteGitHubWebhook(_ *disgord.ExecuteWebhook, _ bool, _ disgord.Snowflake) (*disgord.Message, error) {
	return nil, nil
}

func (w *WebhookQueryBuilderNop) ExecuteSlackWebhook(_ *disgord.ExecuteWebhook, _ bool, _ disgord.Snowflake) (*disgord.Message, error) {
	return nil, nil
}

func (w *WebhookQueryBuilderNop) Get() (*disgord.Webhook, error) {
	return nil, nil
}

func (w *WebhookQueryBuilderNop) Update(_ *disgord.UpdateWebhook) (*disgord.Webhook, error) {
	return nil, nil
}

func (w *WebhookQueryBuilderNop) UpdateBuilder() disgord.UpdateWebhookBuilder {
	return nil
}

func (w *WebhookQueryBuilderNop) WithToken(_ string) disgord.WebhookWithTokenQueryBuilder {
	return nil
}

type WebhookWithTokenQueryBuilderNop struct {
	Ctx       context.Context
	Flags     disgord.Flag
	ChannelID disgord.Snowflake
	GuildID   disgord.Snowflake
	UserID    disgord.Snowflake
}

var _ disgord.WebhookWithTokenQueryBuilder = &WebhookWithTokenQueryBuilderNop{}

func (w WebhookWithTokenQueryBuilderNop) WithContext(ctx context.Context) disgord.WebhookWithTokenQueryBuilder {
	w.Ctx = ctx
	return &w
}

func (w WebhookWithTokenQueryBuilderNop) WithFlags(flags ...disgord.Flag) disgord.WebhookWithTokenQueryBuilder {
	w.Flags = mergeFlags(flags)
	return &w
}

func (w *WebhookWithTokenQueryBuilderNop) Delete() error {
	return nil
}

func (w *WebhookWithTokenQueryBuilderNop) Execute(_ *disgord.ExecuteWebhook, _ bool, _ disgord.Snowflake, _ string) (*disgord.Message, error) {
	return nil, nil
}

func (w *WebhookWithTokenQueryBuilderNop) Get() (*disgord.Webhook, error) {
	return nil, nil
}

func (w *WebhookWithTokenQueryBuilderNop) Update(_ *disgord.UpdateWebhook) (*disgord.Webhook, error) {
	return nil, nil
}

func (w *WebhookWithTokenQueryBuilderNop) UpdateBuilder() disgord.UpdateWebhookBuilder {
	return nil
}
