package disgord

import (
	"errors"
	"fmt"
	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
	"net/http"
)

// Deprecated: use ActiveGuildThreads
type ResponseBodyGuildThreads = ActiveGuildThreads

// Deprecated: use GetArchivedThreads
type GetThreads = GetArchivedThreads

// Deprecated: use ArchivedThreads
type ResponseBodyThreads = ArchivedThreads

// Deprecated: use CreateThreadWithoutMessage
type CreateThreadNoMessage = CreateThreadWithoutMessage

// Deprecated: use GuildWidget
type GuildEmbed = GuildWidget

// Deprecated: use ErrMissingRequiredField
var MissingRequiredFieldErr = ErrMissingRequiredField

// Deprecated: use ErrMissingGuildID
var MissingGuildIDErr = fmt.Errorf("guild: %w", MissingIDErr)

// Deprecated: use ErrMissingID
var MissingIDErr = fmt.Errorf("id: %w", MissingRequiredFieldErr)

// Deprecated: use ErrMissingChannelID
var MissingChannelIDErr = fmt.Errorf("channel: %w", MissingIDErr)

// Deprecated: use ErrMissingUserID
var MissingUserIDErr = fmt.Errorf("user: %w", MissingIDErr)

// Deprecated: use ErrMissingMessageID
var MissingMessageIDErr = fmt.Errorf("message: %w", MissingIDErr)

// Deprecated: use ErrMissingEmojiID
var MissingEmojiIDErr = fmt.Errorf("emoji: %w", MissingIDErr)

// Deprecated: use ErrMissingRoleID
var MissingRoleIDErr = fmt.Errorf("role: %w", MissingIDErr)

// Deprecated: use ErrMissingWebhookID
var MissingWebhookIDErr = fmt.Errorf("webhook: %w", MissingIDErr)

// Deprecated: use ErrMissingPermissionOverwriteID
var MissingPermissionOverwriteIDErr = fmt.Errorf("channel permission overwrite: %w", MissingIDErr)

// Deprecated: use ErrMissingName
var MissingNameErr = fmt.Errorf("name: %w", MissingRequiredFieldErr)

// Deprecated: use ErrMissingGuildName
var MissingGuildNameErr = fmt.Errorf("guild: %w", MissingNameErr)

// Deprecated: use ErrMissingChannelName
var MissingChannelNameErr = fmt.Errorf("channel: %w", MissingNameErr)

// Deprecated: use ErrMissingWebhookName
var MissingWebhookNameErr = fmt.Errorf("webhook: %w", MissingNameErr)

// Deprecated: use ErrMissingThreadName
var MissingThreadNameErr = fmt.Errorf("thread: %w", MissingNameErr)

// Deprecated: use ErrMissingWebhookToken
var MissingWebhookTokenErr = errors.New("webhook token was not set")

// Deprecated: use ErrIllegalValue
var IllegalValueErr = errors.New("illegal value")

func (g guildQueryBuilder) KickVoiceParticipant(userID Snowflake) error {
	return g.DisconnectVoiceParticipant(userID)
}

//generate-rest-params: roles:[]Snowflake,
//generate-rest-basic-execute: emoji:*Emoji,
type createGuildEmojiBuilder struct {
	r RESTBuilder
}

// updateGuildMemberBuilder ...
// https://discord.com/developers/docs/resources/guild#modify-guild-member-json-params
//generate-rest-params: nick:string, roles:[]Snowflake, mute:bool, deaf:bool, channel_id:Snowflake,
//generate-rest-basic-execute: err:error,
type updateGuildMemberBuilder struct {
	r RESTBuilder
}

func (c currentUserQueryBuilder) LeaveGuild(id Snowflake) (err error) {
	return c.client.Guild(id).Leave()
}

// KickFromVoice kicks member out of voice channel. Assuming they are in one.
func (b *updateGuildMemberBuilder) KickFromVoice() UpdateGuildMemberBuilder {
	b.r.param("channel_id", 0)
	return b
}

// DeleteNick removes nickname for user. Requires permission MANAGE_NICKNAMES
func (b *updateGuildMemberBuilder) DeleteNick() UpdateGuildMemberBuilder {
	b.r.param("nick", "")
	return b
}

//generate-rest-params: enabled:bool, channel_id:Snowflake,
//generate-rest-basic-execute: embed:*GuildEmbed,
type updateGuildEmbedBuilder struct {
	r RESTBuilder
}

// UpdateEmbedBuilder Modify a guild embed object for the guild. All attributes may be passed in with JSON and
// modified. Requires the 'MANAGE_GUILD' permission. Returns the updated guild embed object.
func (g guildQueryBuilder) UpdateEmbedBuilder() UpdateGuildEmbedBuilder {
	builder := &updateGuildEmbedBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &GuildEmbed{}
	}
	builder.r.flags = g.flags
	builder.r.setup(g.client.req, &httd.Request{
		Method:      http.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildEmbed(g.gid),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// Deprecated: use Update instead
func (m messageQueryBuilder) UpdateBuilder() UpdateMessageBuilder {
	builder := &updateMessageBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Message{}
	}
	builder.r.flags = m.flags
	builder.r.addPrereq(m.cid.IsZero(), "channelID must be set to get channel messages")
	builder.r.addPrereq(m.mid.IsZero(), "msgID must be set to edit the message")
	builder.r.setup(m.client.req, &httd.Request{
		Method:      http.MethodPatch,
		Ctx:         m.ctx,
		Endpoint:    "/channels/" + m.cid.String() + "/messages/" + m.mid.String(),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// UpdateBuilder [REST] Update a Channels settings. Requires the 'MANAGE_CHANNELS' permission for the guild. Returns
// a channel on success, and a 400 BAD REQUEST on invalid parameters. Fires a Channel Update Gateway event. If
// modifying a category, individual Channel Update events will fire for each child channel that also changes.
// For the PATCH method, all the JSON Params are optional.
//  Method                  PUT/PATCH
//  Endpoint                /channels/{channel.id}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#modify-channel
//  Reviewed                2018-06-07
//  Comment                 andersfylling: only implemented the patch method, as its parameters are optional.
func (c channelQueryBuilder) UpdateBuilder() UpdateChannelBuilder {
	builder := &updateChannelBuilder{}
	builder.r.itemFactory = func() interface{} {
		return c.client.pool.channel.Get()
	}
	builder.r.flags = c.flags
	builder.r.setup(c.client.req, &httd.Request{
		Method:      http.MethodPatch,
		Ctx:         c.ctx,
		Endpoint:    endpoint.Channel(c.cid),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// UpdateBuilder Modify the given emoji. Requires the 'MANAGE_EMOJIS' permission.
// Returns the updated emoji object on success. Fires a Guild Emojis Update Gateway event.
func (g guildEmojiQueryBuilder) UpdateBuilder() UpdateGuildEmojiBuilder {
	builder := &updateGuildEmojiBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Emoji{}
	}
	builder.r.flags = g.flags
	builder.r.setup(g.client.req, &httd.Request{
		Method:      http.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildEmoji(g.gid, g.emojiID),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// UpdateBuilder is used to create a guild update builder.
func (g guildQueryBuilder) UpdateBuilder() UpdateGuildBuilder {
	builder := &updateGuildBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Guild{}
	}
	builder.r.setup(g.client.req, &httd.Request{
		Method:      http.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.Guild(g.gid),
		ContentType: httd.ContentTypeJSON,
	}, nil)
	builder.r.flags = g.flags

	return builder
}

// UpdateBuilder is used to create a builder to update a guild member.
func (g guildMemberQueryBuilder) UpdateBuilder() UpdateGuildMemberBuilder {
	builder := &updateGuildMemberBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Member{
			GuildID: g.gid,
			UserID:  g.uid,
		}
	}
	builder.r.flags = g.flags
	builder.r.setup(g.client.req, &httd.Request{
		Method:      http.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildMember(g.gid, g.uid),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	// TODO: cache member changes
	return builder
}

// UpdateBuilder Modify a guild role. Requires the 'MANAGE_ROLES' permission.
// Returns the updated role on success. Fires a Guild Role Update Gateway event.
func (g guildRoleQueryBuilder) UpdateBuilder() UpdateGuildRoleBuilder {
	builder := &updateGuildRoleBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Role{}
	}
	builder.r.flags = g.flags
	builder.r.IgnoreCache().setup(g.client.req, &httd.Request{
		Method:      http.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildRole(g.gid, g.roleID),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// UpdateBuilder [REST] Modify the requester's user account settings. Returns a user object on success.
//  Method                  PATCH
//  Endpoint                /users/@me
//  Discord documentation   https://discord.com/developers/docs/resources/user#modify-current-user
//  Reviewed                2019-02-18
//  Comment                 -
func (c currentUserQueryBuilder) UpdateBuilder() UpdateCurrentUserBuilder {
	builder := &updateCurrentUserBuilder{}
	builder.r.itemFactory = userFactory // TODO: peak cached user
	builder.r.flags = c.flags
	builder.r.setup(c.client.req, &httd.Request{
		Method:      http.MethodPatch,
		Ctx:         c.ctx,
		Endpoint:    endpoint.UserMe(),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	// TODO: cache changes?
	return builder
}

// UpdateBuilder [REST] Same as UpdateWebhook, except this call does not require authentication,
// does _not_ accept a channel_id parameter in the body, and does not return a user in the webhook object.
//  Method                  PATCH
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#modify-webhook-with-token
//  Reviewed                2018-08-14
//  Comment                 All parameters to this endpoint. are optional.
func (w webhookWithTokenQueryBuilder) UpdateBuilder() UpdateWebhookBuilder {
	builder := &updateWebhookBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Webhook{}
	}
	builder.r.flags = w.flags
	builder.r.addPrereq(w.webhookID.IsZero(), "given webhook ID was not set, there is nothing to modify")
	builder.r.addPrereq(w.token == "", "given webhook token was not set")
	builder.r.setup(w.client.req, &httd.Request{
		Method:      http.MethodPatch,
		Ctx:         w.ctx,
		Endpoint:    endpoint.WebhookToken(w.webhookID, w.token),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// UpdateBuilder [REST] Modify a webhook. Requires the 'MANAGE_WEBHOOKS' permission.
// Returns the updated webhook object on success.
//  Method                  PATCH
//  Endpoint                /webhooks/{webhook.id}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#modify-webhook
//  Reviewed                2018-08-14
//  Comment                 All parameters to this endpoint.
func (w webhookQueryBuilder) UpdateBuilder() UpdateWebhookBuilder {
	builder := &updateWebhookBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Webhook{}
	}
	builder.r.flags = w.flags
	builder.r.addPrereq(w.webhookID.IsZero(), "given webhook ID was not set, there is nothing to modify")
	builder.r.setup(w.client.req, &httd.Request{
		Method:      http.MethodPatch,
		Ctx:         w.ctx,
		Endpoint:    endpoint.Webhook(w.webhookID),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

func (g guildQueryBuilder) EstimatePruneMembersCount(days int) (estimate int, err error) {
	return g.GetPruneMembersCount(&GetPruneMembersCount{
		Days: &days,
	})
}

func (g guildQueryBuilder) GetEmbed() (*GuildWidget, error) {
	return g.GetWidget()
}

// Deprecated: use ErrCacheMiss
var CacheMissErr = ErrCacheMiss

// Deprecated: use ErrCacheEntryAlreadyExists
var CacheEntryAlreadyExistsErr = ErrCacheEntryAlreadyExists

// Deprecated: use ErrMissingClientInstance
var MissingClientInstanceErr = ErrMissingClientInstance

// Deprecated: use ErrMissingRESTParams
var MissingRESTParamsErr = ErrMissingRESTParams

const (
	// Deprecated: ...
	SUB_COMMAND = OptionTypeSubCommand
	// Deprecated: ...
	SUB_COMMAND_GROUP = OptionTypeSubCommandGroup
	// Deprecated: ...
	STRING = OptionTypeString
	// Deprecated: ...
	INTEGER = OptionTypeInteger
	// Deprecated: ...
	BOOLEAN = OptionTypeBoolean
	// Deprecated: ...
	USER = OptionTypeUser
	// Deprecated: ...
	CHANNEL = OptionTypeChannel
	// Deprecated: ...
	ROLE = OptionTypeRole
	// Deprecated: ...
	MENTIONABLE = OptionTypeMentionable
	// Deprecated: ...
	NUMBER = OptionTypeNumber
)

// updateMessageBuilder, params here
//  https://discord.com/developers/docs/resources/channel#edit-message-json-params
//generate-rest-params: content:string, embed:*Embed,
//generate-rest-basic-execute: message:*Message,
type updateMessageBuilder struct {
	r RESTBuilder
}

// SetAllowedMentions sets the allowed mentions for the updateMessageBuilder then returns the builder to allow chaining.
func (b *updateMessageBuilder) SetAllowedMentions(mentions *AllowedMentions) *updateMessageBuilder {
	b.r.param("allowed_mentions", mentions)
	return b
}

// updateChannelBuilder https://discord.com/developers/docs/resources/channel#modify-channel-json-params
//generate-rest-params: parent_id:Snowflake, permission_overwrites:[]PermissionOverwrite, user_limit:uint, bitrate:uint, rate_limit_per_user:uint, nsfw:bool, topic:string, position:int, name:string,
//generate-rest-basic-execute: channel:*Channel,
type updateChannelBuilder struct {
	r RESTBuilder
}

func (b *updateChannelBuilder) AddPermissionOverwrite(permission PermissionOverwrite) *updateChannelBuilder {
	if _, exists := b.r.body["permission_overwrites"]; !exists {
		b.SetPermissionOverwrites([]PermissionOverwrite{permission})
	} else {
		s := b.r.body["permission_overwrites"].([]PermissionOverwrite)
		s = append(s, permission)
		b.SetPermissionOverwrites(s)
	}
	return b
}
func (b *updateChannelBuilder) AddPermissionOverwrites(permissions []PermissionOverwrite) *updateChannelBuilder {
	for i := range permissions {
		b.AddPermissionOverwrite(permissions[i])
	}
	return b
}

func (b *updateChannelBuilder) RemoveParentID() *updateChannelBuilder {
	b.r.param("parent_id", nil)
	return b
}

//generate-rest-params: name:string, roles:[]Snowflake,
//generate-rest-basic-execute: emoji:*Emoji,
type updateGuildEmojiBuilder struct {
	r RESTBuilder
}

// updateGuildBuilder https://discord.com/developers/docs/resources/guild#modify-guild-json-params
//generate-rest-params: name:string, region:string, verification_level:int, default_message_notifications:DefaultMessageNotificationLvl, explicit_content_filter:ExplicitContentFilterLvl, afk_channel_id:Snowflake, afk_timeout:int, icon:string, owner_id:Snowflake, splash:string, system_channel_id:Snowflake,
//generate-rest-basic-execute: guild:*Guild,
type updateGuildBuilder struct {
	r RESTBuilder
}

// updateGuildRoleBuilder ...
//generate-rest-basic-execute: role:*Role,
//generate-rest-params: name:string, permissions:PermissionBit, color:uint, hoist:bool, mentionable:bool,
type updateGuildRoleBuilder struct {
	r RESTBuilder
}

// updateCurrentUserBuilder ...
//generate-rest-params: username:string, avatar:string,
//generate-rest-basic-execute: user:*User,
type updateCurrentUserBuilder struct {
	r RESTBuilder
}

// UpdateWebhook https://discord.com/developers/docs/resources/webhook#modify-webhook-json-params
// Allows changing the name of the webhook, avatar and moving it to another channel. It also allows to resetting the
// avatar by providing a nil to SetAvatar.
//
//generate-rest-params: name:string, avatar:string, channel_id:Snowflake,
//generate-rest-basic-execute: webhook:*Webhook,
type updateWebhookBuilder struct {
	r RESTBuilder
}

// SetDefaultAvatar will reset the webhook image
func (u *updateWebhookBuilder) SetDefaultAvatar() *updateWebhookBuilder {
	u.r.param("avatar", nil)
	return u
}

// Deprecated: specify permissions when using the Client.AuthorizeBotURL method
func (c *Client) AddPermission(permission PermissionBit) (updatedPermissions PermissionBit) {
	if permission < 0 {
		permission = 0
	}

	c.permissions |= permission
	return c.GetPermissions()
}

// Deprecated: ...
func (c *Client) GetPermissions() (permissions PermissionBit) {
	return c.permissions
}

func (c currentUserQueryBuilder) GetUserConnections() (connections []*UserConnection, err error) {
	return c.GetConnections()
}

// Deprecated: use Update instead
func (m messageQueryBuilder) SetContent(content string) (*Message, error) {
	builder := m.WithContext(m.ctx).UpdateBuilder()
	return builder.
		SetContent(content).
		Execute()
}

// Deprecated: use Update instead
func (m messageQueryBuilder) SetEmbed(embed *Embed) (*Message, error) {
	builder := m.WithContext(m.ctx).UpdateBuilder()
	return builder.
		SetEmbed(embed).
		Execute()
}
