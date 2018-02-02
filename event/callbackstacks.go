package event

import (
	"context"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/guild"
	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/disgord/voice"
	"github.com/andersfylling/disgord/webhook"
)

type Handler interface {
	Add(interface{}) error
	Trigger(...*interface{}) error // TODO: the param should be a specific event holder type
}

// socket
//

// ReadyCallbackStack ***************
type ReadyHandler interface {
	Add(cb ReadyCallback)
	Trigger(ctx context.Context)
}
type ReadyCallbackStack struct {
	listeners []ReadyCallback
}

func (rds *ReadyCallbackStack) Add(cb ReadyCallback) (err error) {
	if rds.listeners == nil {
		rds.listeners = []ReadyCallback{}
	}

	rds.listeners = append(rds.listeners, cb)
	return nil
}

func (rds *ReadyCallbackStack) Trigger(ctx context.Context) (err error) {
	for _, listener := range rds.listeners {
		listener(ctx)
	}

	return nil
}

// ResumedCallbackStack **********
type ResumedHandler interface {
	Add(cb ReadyCallback)
	Trigger(ctx context.Context)
}
type ResumedCallbackStack struct {
	listeners []ResumedCallback
}

func (stack *ResumedCallbackStack) Add(cb ResumedCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []ResumedCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *ResumedCallbackStack) Trigger(ctx context.Context) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx)
	}

	return nil
}

// channel
//

// ChannelCreateCallbackStack **************
type ChannelCreateHandler interface {
	Add(ChannelCreateCallback)
	Trigger(context.Context, *channel.Channel)
}
type ChannelCreateCallbackStack struct {
	listeners []ChannelCreateCallback
}

func (stack *ChannelCreateCallbackStack) Add(cb ChannelCreateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []ChannelCreateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *ChannelCreateCallbackStack) Trigger(ctx context.Context, c *channel.Channel) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, c)
	}

	return nil
}

// ChannelUpdateCallbackStack ************
type ChannelUpdateHandler interface {
	Add(ChannelUpdateCallback)
	Trigger(context.Context, *channel.Channel)
}
type ChannelUpdateCallbackStack struct {
	listeners []ChannelUpdateCallback
}

func (stack *ChannelUpdateCallbackStack) Add(cb ChannelUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []ChannelUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *ChannelUpdateCallbackStack) Trigger(ctx context.Context, c *channel.Channel) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, c)
	}

	return nil
}

// ChannelDeleteCallbackStack ***********
type ChannelDeleteHandler interface {
	Add(ChannelDeleteCallback)
	Trigger(context.Context, *channel.Channel)
}
type ChannelDeleteCallbackStack struct {
	listeners []ChannelDeleteCallback
}

func (stack *ChannelDeleteCallbackStack) Add(cb ChannelDeleteCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []ChannelDeleteCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *ChannelDeleteCallbackStack) Trigger(ctx context.Context, c *channel.Channel) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, c)
	}

	return nil
}

// ChannelPinsUpdateCallbackStack **********
type ChannelPinsUpdateHandler interface {
	Add(ChannelPinsUpdateCallback)
	Trigger(context.Context, *channel.Channel)
}
type ChannelPinsUpdateCallbackStack struct {
	listeners []ChannelPinsUpdateCallback
}

func (stack *ChannelPinsUpdateCallbackStack) Add(cb ChannelPinsUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []ChannelPinsUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *ChannelPinsUpdateCallbackStack) Trigger(ctx context.Context, c *channel.Channel) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, c)
	}

	return nil
}

// Guild in general
//

// GuildCreateCallbackStack **********
type GuildCreateHandler interface {
	Add(GuildCreateCallback)
	Trigger(context.Context, *guild.Guild)
}
type GuildCreateCallbackStack struct {
	listeners []GuildCreateCallback
}

func (stack *GuildCreateCallbackStack) Add(cb GuildCreateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildCreateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildCreateCallbackStack) Trigger(ctx context.Context, g *guild.Guild) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, g)
	}

	return nil
}

// GuildUpdateCallbackStack .....
type GuildUpdateHandler interface {
	Add(GuildUpdateCallback)
	Trigger(context.Context, *guild.Guild)
}
type GuildUpdateCallbackStack struct {
	listeners []GuildUpdateCallback
}

func (stack *GuildUpdateCallbackStack) Add(cb GuildUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildUpdateCallbackStack) Trigger(ctx context.Context, g *guild.Guild) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, g)
	}

	return nil
}

// GuildDeleteCallbackStack *********
type GuildDeleteHandler interface {
	Add(GuildDeleteCallback)
	Trigger(context.Context, *guild.Guild)
}
type GuildDeleteCallbackStack struct {
	listeners []GuildDeleteCallback
}

func (stack *GuildDeleteCallbackStack) Add(cb GuildDeleteCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildDeleteCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildDeleteCallbackStack) Trigger(ctx context.Context, g *guild.Guild) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, g)
	}

	return nil
}

// GuildBanAddCallbackStack **************
type GuildBanAddHandler interface {
	Add(GuildBanAddCallback)
	Trigger(context.Context, *guild.Guild)
}
type GuildBanAddCallbackStack struct {
	listeners []GuildBanAddCallback
}

func (stack *GuildBanAddCallbackStack) Add(cb GuildBanAddCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildBanAddCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildBanAddCallbackStack) Trigger(ctx context.Context, g *guild.Guild) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, g)
	}

	return nil
}

// GuildBanRemoveCallbackStack *********
type GuildBanRemoveHandler interface {
	Add(GuildBanRemoveCallback)
	Trigger(context.Context, *guild.Guild)
}
type GuildBanRemoveCallbackStack struct {
	listeners []GuildBanRemoveCallback
}

func (stack *GuildBanRemoveCallbackStack) Add(cb GuildBanRemoveCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildBanRemoveCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildBanRemoveCallbackStack) Trigger(ctx context.Context, g *guild.Guild) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, g)
	}

	return nil
}

// GuildEmojisUpdateCallbackStack ***********
type GuildEmojisUpdateHandler interface {
	Add(GuildEmojisUpdateCallback)
	Trigger(context.Context, *guild.Guild)
}
type GuildEmojisUpdateCallbackStack struct {
	listeners []GuildEmojisUpdateCallback
}

func (stack *GuildEmojisUpdateCallbackStack) Add(cb GuildEmojisUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildEmojisUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildEmojisUpdateCallbackStack) Trigger(ctx context.Context, g *guild.Guild) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, g)
	}

	return nil
}

// GuildIntegrationsUpdateCallbackStack *******************
type GuildIntegrationsUpdateHandler interface {
	Add(GuildIntegrationsUpdateCallback)
	Trigger(context.Context, *guild.Guild)
}
type GuildIntegrationsUpdateCallbackStack struct {
	listeners []GuildIntegrationsUpdateCallback
}

func (stack *GuildIntegrationsUpdateCallbackStack) Add(cb GuildIntegrationsUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildIntegrationsUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildIntegrationsUpdateCallbackStack) Trigger(ctx context.Context, g *guild.Guild) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, g)
	}

	return nil
}

// Guild Member
//

// GuildMemberAddCallbackStack ***********************
type GuildMemberAddHandler interface {
	Add(GuildMemberAddCallback)
	Trigger(context.Context, *guild.Member)
}
type GuildMemberAddCallbackStack struct {
	listeners []GuildMemberAddCallback
}

func (stack *GuildMemberAddCallbackStack) Add(cb GuildMemberAddCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildMemberAddCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildMemberAddCallbackStack) Trigger(ctx context.Context, member *guild.Member) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, member)
	}

	return nil
}

// GuildMemberRemoveCallbackStack *******************
type GuildMemberRemoveHandler interface {
	Add(GuildMemberRemoveCallback)
	Trigger(context.Context, *guild.Member)
}
type GuildMemberRemoveCallbackStack struct {
	listeners []GuildMemberRemoveCallback
}

func (stack *GuildMemberRemoveCallbackStack) Add(cb GuildMemberRemoveCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildMemberRemoveCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildMemberRemoveCallbackStack) Trigger(ctx context.Context, member *guild.Member) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, member)
	}

	return nil
}

// GuildMemberUpdateCallbackStack **************
type GuildMemberUpdateHandler interface {
	Add(GuildMemberUpdateCallback)
	Trigger(context.Context, *guild.Member)
}
type GuildMemberUpdateCallbackStack struct {
	listeners []GuildMemberUpdateCallback
}

func (stack *GuildMemberUpdateCallbackStack) Add(cb GuildMemberUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildMemberUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildMemberUpdateCallbackStack) Trigger(ctx context.Context, member *guild.Member) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, member)
	}

	return nil
}

// GuildMemberChunkCallbackStack **************
type GuildMemberChunkHandler interface {
	Add(GuildMemberChunkCallback)
	Trigger(context.Context, []*guild.Member)
}
type GuildMemberChunkCallbackStack struct {
	listeners []GuildMemberChunkCallback
}

func (stack *GuildMemberChunkCallbackStack) Add(cb GuildMemberChunkCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildMemberChunkCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildMemberChunkCallbackStack) Trigger(ctx context.Context, members []*guild.Member) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, members)
	}

	return nil
}

// Guild role
//

// GuildRoleCreateCallbackStack *************
type GuildRoleCreateHandler interface {
	Add(GuildRoleCreateCallback)
	Trigger(ctx context.Context, role *discord.Role)
}
type GuildRoleCreateCallbackStack struct {
	listeners []GuildRoleCreateCallback
}

func (stack *GuildRoleCreateCallbackStack) Add(cb GuildRoleCreateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildRoleCreateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildRoleCreateCallbackStack) Trigger(ctx context.Context, role *discord.Role) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, role)
	}

	return nil
}

// GuildRoleUpdateCallbackStack ***************
type GuildRoleUpdateHandler interface {
	Add(GuildRoleUpdateCallback)
	Trigger(ctx context.Context, role *discord.Role)
}
type GuildRoleUpdateCallbackStack struct {
	listeners []GuildRoleUpdateCallback
}

func (stack *GuildRoleUpdateCallbackStack) Add(cb GuildRoleUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildRoleUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildRoleUpdateCallbackStack) Trigger(ctx context.Context, role *discord.Role) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, role)
	}

	return nil
}

// GuildRoleDeleteCallbackStack **************
type GuildRoleDeleteHandler interface {
	Add(GuildRoleDeleteCallback)
	Trigger(ctx context.Context, role *discord.Role)
}
type GuildRoleDeleteCallbackStack struct {
	listeners []GuildRoleDeleteCallback
}

func (stack *GuildRoleDeleteCallbackStack) Add(cb GuildRoleDeleteCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildRoleDeleteCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildRoleDeleteCallbackStack) Trigger(ctx context.Context, role *discord.Role) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, role)
	}

	return nil
}

// message
//

// MessageCreateCallbackStack ********************
type MessageCreateHandler interface {
	Add(MessageCreateCallback)
	Trigger(context.Context, *channel.Message)
}
type MessageCreateCallbackStack struct {
	listeners []MessageCreateCallback
}

func (stack *MessageCreateCallbackStack) Add(cb MessageCreateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []MessageCreateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *MessageCreateCallbackStack) Trigger(ctx context.Context, msg *channel.Message) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, msg)
	}

	return nil
}

// MessageUpdateCallbackStack ****************
type MessageUpdateHandler interface {
	Add(MessageUpdateCallback)
	Trigger(context.Context, *channel.Message)
}
type MessageUpdateCallbackStack struct {
	listeners []MessageUpdateCallback
}

func (stack *MessageUpdateCallbackStack) Add(cb MessageUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []MessageUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *MessageUpdateCallbackStack) Trigger(ctx context.Context, msg *channel.Message) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, msg)
	}

	return nil
}

// MessageDeleteCallbackStack ***************
type MessageDeleteHandler interface {
	Add(MessageDeleteCallback)
	Trigger(context.Context, *channel.Message)
}
type MessageDeleteCallbackStack struct {
	listeners []MessageDeleteCallback
}

func (stack *MessageDeleteCallbackStack) Add(cb MessageDeleteCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []MessageDeleteCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *MessageDeleteCallbackStack) Trigger(ctx context.Context, msg *channel.Message) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, msg)
	}

	return nil
}

// MessageDeleteBulkCallbackStack ****************
type MessageDeleteBulkHandler interface {
	Add(MessageDeleteBulkCallback)
	Trigger(context.Context, []*channel.Message)
}
type MessageDeleteBulkCallbackStack struct {
	listeners []MessageDeleteBulkCallback
}

func (stack *MessageDeleteBulkCallbackStack) Add(cb MessageDeleteBulkCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []MessageDeleteBulkCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *MessageDeleteBulkCallbackStack) Trigger(ctx context.Context, msgs []*channel.Message) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, msgs)
	}

	return nil
}

// message reaction
//

// MessageReactionAddCallbackStack ************
type MessageReactionAddHandler interface {
	Add(MessageReactionAddCallback)
	Trigger(context.Context, *channel.Message)
}
type MessageReactionAddCallbackStack struct {
	listeners []MessageReactionAddCallback
}

func (stack *MessageReactionAddCallbackStack) Add(cb MessageReactionAddCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []MessageReactionAddCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *MessageReactionAddCallbackStack) Trigger(ctx context.Context, msg *channel.Message) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, msg)
	}

	return nil
}

// MessageReactionRemoveCallbackStack *********
type MessageReactionRemoveHandler interface {
	Add(MessageReactionRemoveCallback)
	Trigger(context.Context, *channel.Message)
}
type MessageReactionRemoveCallbackStack struct {
	listeners []MessageReactionRemoveCallback
}

func (stack *MessageReactionRemoveCallbackStack) Add(cb MessageReactionRemoveCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []MessageReactionRemoveCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *MessageReactionRemoveCallbackStack) Trigger(ctx context.Context, msg *channel.Message) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, msg)
	}

	return nil
}

// MessageReactionRemoveAllCallbackStack *********
type MessageReactionRemoveAllHandler interface {
	Add(MessageReactionRemoveAllCallback)
	Trigger(context.Context, *channel.Message)
}
type MessageReactionRemoveAllCallbackStack struct {
	listeners []MessageReactionRemoveAllCallback
}

func (stack *MessageReactionRemoveAllCallbackStack) Add(cb MessageReactionRemoveAllCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []MessageReactionRemoveAllCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *MessageReactionRemoveAllCallbackStack) Trigger(ctx context.Context, msgs []*channel.Message) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, msgs)
	}

	return nil
}

// presence
//

// PresenceUpdateCallbackStack *************
type PresenceUpdateHandler interface {
	Add(PresenceUpdateCallback)
	Trigger(context.Context, *discord.Presence)
}
type PresenceUpdateCallbackStack struct {
	listeners []PresenceUpdateCallback
}

func (stack *PresenceUpdateCallbackStack) Add(cb PresenceUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []PresenceUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *PresenceUpdateCallbackStack) Trigger(ctx context.Context, p *discord.Presence) (err error) {

	for _, listener := range stack.listeners {
		listener(ctx, p)
	}

	return nil
}

// typing start
//

// TypingStartCallbackStack ******************
type TypingStartHandler interface {
	Add(TypingStartCallback)
	Trigger(context.Context, *user.User, *channel.Channel)
}
type TypingStartCallbackStack struct {
	listeners []TypingStartCallback
}

func (stack *TypingStartCallbackStack) Add(cb TypingStartCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []TypingStartCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *TypingStartCallbackStack) Trigger(ctx context.Context, u *user.User, c *channel.Channel) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, u, c)
	}

	return nil
}

// user update
type UserUpdateHandler interface {
	Add(UserUpdateCallback)
	Trigger(context.Context, *user.User)
}
type UserUpdateCallbackStack struct {
	listeners []UserUpdateCallback
}

func (stack *UserUpdateCallbackStack) Add(cb UserUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []UserUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *UserUpdateCallbackStack) Trigger(ctx context.Context, u *user.User) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, u)
	}

	return nil
}

// voice
//

// VoiceStateUpdateCallbackStack *************************
type VoiceStateUpdateHandler interface {
	Add(VoiceStateUpdateCallback)
	Trigger(context.Context, *voice.State)
}
type VoiceStateUpdateCallbackStack struct {
	listeners []VoiceStateUpdateCallback
}

func (stack *VoiceStateUpdateCallbackStack) Add(cb VoiceStateUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []VoiceStateUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *VoiceStateUpdateCallbackStack) Trigger(ctx context.Context, vst *voice.State) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, vst)
	}

	return nil
}

// VoiceServerUpdateCallbackStack ***********************
type VoiceServerUpdateHandler interface {
	Add(VoiceServerUpdateCallback)
	Trigger(context.Context, *voice.State)
}
type VoiceServerUpdateCallbackStack struct {
	listeners []VoiceServerUpdateCallback
}

func (stack *VoiceServerUpdateCallbackStack) Add(cb VoiceServerUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []VoiceServerUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *VoiceServerUpdateCallbackStack) Trigger(ctx context.Context, vst *voice.State) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, vst)
	}

	return nil
}

// WebhooksUpdateCallbackStack *******************
type WebhooksUpdateHandler interface {
	Add(cb WebhooksUpdateCallback)
	Trigger(ctx context.Context, wb *webhook.Webhook)
}
type WebhooksUpdateCallbackStack struct {
	listeners []WebhooksUpdateCallback
}

func (stack *WebhooksUpdateCallbackStack) Add(cb WebhooksUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []WebhooksUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *WebhooksUpdateCallbackStack) Trigger(ctx context.Context, wb *webhook.Webhook) (err error) {
	for _, listener := range stack.listeners {
		listener(ctx, wb)
	}

	return nil
}
