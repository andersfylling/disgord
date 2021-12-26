package disgord

import (
	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
	"net/http"
)

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

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

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

// UpdateWebhookParams https://discord.com/developers/docs/resources/webhook#modify-webhook-json-params
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

//////////////////////////////////////////////////////
//
// REST Wrappers
//
//////////////////////////////////////////////////////

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
