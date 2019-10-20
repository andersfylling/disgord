package disgord

import (
	"errors"
	"net/http"

	"github.com/andersfylling/disgord/internal/constant"
	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
)

// Webhook Used to represent a webhook
// https://discordapp.com/developers/docs/resources/webhook#webhook-object
type Webhook struct {
	Lockable `json:"-"`

	ID        Snowflake `json:"id"`                 //  |
	GuildID   Snowflake `json:"guild_id,omitempty"` //  |?
	ChannelID Snowflake `json:"channel_id"`         //  |
	User      *User     `json:"user,omitempty"`     // ?|
	Name      string    `json:"name"`               //  |?
	Avatar    string    `json:"avatar"`             //  |?
	Token     string    `json:"token"`              //  |
}

// DeepCopy see interface at struct.go#DeepCopier
func (w *Webhook) DeepCopy() (copy interface{}) {
	copy = &Webhook{}
	w.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (w *Webhook) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var hook *Webhook
	if hook, ok = other.(*Webhook); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *Webhook")
		return
	}

	if constant.LockedMethods {
		w.RLock()
		hook.Lock()
	}

	hook.ID = w.ID
	hook.GuildID = w.GuildID
	hook.ChannelID = w.ChannelID
	hook.User = w.User.DeepCopy().(*User)
	hook.Name = w.Name
	hook.Avatar = w.Avatar
	hook.Token = w.Token

	if constant.LockedMethods {
		w.RUnlock()
		hook.Unlock()
	}
	return
}

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

func ratelimitWebhook(id Snowflake) string {
	return "wh:" + id.String()
}

// CreateWebhookParams json params for the create webhook rest request avatar string
// https://discordapp.com/developers/docs/resources/user#avatar-data
type CreateWebhookParams struct {
	Name   string `json:"name"`   // name of the webhook (2-32 characters)
	Avatar string `json:"avatar"` // avatar data uri scheme, image for the default webhook avatar
}

func (c *CreateWebhookParams) FindErrors() error {
	if c.Name == "" {
		return errors.New("webhook must have a name")
	}
	if !(2 <= len(c.Name) && len(c.Name) <= 32) {
		return errors.New("webhook name must be 2 to 32 characters long")
	}
	return nil
}

// CreateWebhook [REST] Create a new webhook. Requires the 'MANAGE_WEBHOOKS' permission.
// Returns a webhook object on success.
//  Method                  POST
//  Endpoint                /channels/{channel.id}/webhooks
//  Rate limiter            /channels/{channel.id}/webhooks
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#create-webhook
//  Reviewed                2018-08-14
//  Comment                 -
func (c *Client) CreateWebhook(channelID Snowflake, params *CreateWebhookParams, flags ...Flag) (ret *Webhook, err error) {
	if params == nil {
		return nil, errors.New("params was nil")
	}
	if err = params.FindErrors(); err != nil {
		return nil, err
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Endpoint:    endpoint.ChannelWebhooks(channelID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.factory = func() interface{} {
		return &Webhook{}
	}

	return getWebhook(r.Execute)
}

// GetChannelWebhooks [REST] Returns a list of channel webhook objects. Requires the 'MANAGE_WEBHOOKS' permission.
//  Method                  POST
//  Endpoint                /channels/{channel.id}/webhooks
//  Rate limiter            /channels/{channel.id}/webhooks
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#get-channel-webhooks
//  Reviewed                2018-08-14
//  Comment                 -
func (c *Client) GetChannelWebhooks(channelID Snowflake, flags ...Flag) (ret []*Webhook, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.ChannelWebhooks(channelID),
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Webhook, 0)
		return &tmp
	}

	return getWebhooks(r.Execute)
}

// GetGuildWebhooks [REST] Returns a list of guild webhook objects. Requires the 'MANAGE_WEBHOOKS' permission.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/webhooks
//  Rate limiter            /guilds/{guild.id}/webhooks
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#get-guild-webhooks
//  Reviewed                2018-08-14
//  Comment                 -
func (c *Client) GetGuildWebhooks(guildID Snowflake, flags ...Flag) (ret []*Webhook, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildWebhooks(guildID),
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Webhook, 0)
		return &tmp
	}

	return getWebhooks(r.Execute)
}

// GetWebhook [REST] Returns the new webhook object for the given id.
//  Method                  GET
//  Endpoint                /webhooks/{webhook.id}
//  Rate limiter            /webhooks/{webhook.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#get-webhook
//  Reviewed                2018-08-14
//  Comment                 -
func (c *Client) GetWebhook(id Snowflake, flags ...Flag) (ret *Webhook, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.Webhook(id),
	}, flags)
	r.factory = func() interface{} {
		return &Webhook{}
	}

	return getWebhook(r.Execute)
}

// GetWebhookWithToken [REST] Same as GetWebhook, except this call does not require authentication and
// returns no user in the webhook object.
//  Method                  GET
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Rate limiter            /webhooks/{webhook.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#get-webhook-with-token
//  Reviewed                2018-08-14
//  Comment                 -
func (c *Client) GetWebhookWithToken(id Snowflake, token string, flags ...Flag) (ret *Webhook, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.WebhookToken(id, token),
	}, flags)
	r.factory = func() interface{} {
		return &Webhook{}
	}

	return getWebhook(r.Execute)
}

// UpdateWebhook [REST] Modify a webhook. Requires the 'MANAGE_WEBHOOKS' permission.
// Returns the updated webhook object on success.
//  Method                  PATCH
//  Endpoint                /webhooks/{webhook.id}
//  Rate limiter            /webhooks/{webhook.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#modify-webhook
//  Reviewed                2018-08-14
//  Comment                 All parameters to this endpoint.
func (c *Client) UpdateWebhook(id Snowflake, flags ...Flag) (builder *updateWebhookBuilder) {
	builder = &updateWebhookBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Webhook{}
	}
	builder.r.flags = flags
	builder.r.addPrereq(id.IsZero(), "given webhook ID was not set, there is nothing to modify")
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      httd.MethodPatch,
		Endpoint:    endpoint.Webhook(id),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// UpdateWebhookWithToken [REST] Same as UpdateWebhook, except this call does not require authentication,
// does _not_ accept a channel_id parameter in the body, and does not return a user in the webhook object.
//  Method                  PATCH
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Rate limiter            /webhooks/{webhook.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#modify-webhook-with-token
//  Reviewed                2018-08-14
//  Comment                 All parameters to this endpoint. are optional.
func (c *Client) UpdateWebhookWithToken(id Snowflake, token string, flags ...Flag) (builder *updateWebhookBuilder) {
	builder = &updateWebhookBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Webhook{}
	}
	builder.r.flags = flags
	builder.r.addPrereq(id.IsZero(), "given webhook ID was not set, there is nothing to modify")
	builder.r.addPrereq(token == "", "given webhook token was not set")
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      httd.MethodPatch,
		Endpoint:    endpoint.WebhookToken(id, token),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// DeleteWebhook [REST] Delete a webhook permanently. User must be owner. Returns a 204 NO CONTENT response on success.
//  Method                  DELETE
//  Endpoint                /webhooks/{webhook.id}
//  Rate limiter            /webhooks/{webhook.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#delete-webhook
//  Reviewed                2018-08-14
//  Comment                 -
func (c *Client) DeleteWebhook(id Snowflake, flags ...Flag) (err error) {
	return c.DeleteWebhookWithToken(id, "", flags...)
}

// DeleteWebhookWithToken [REST] Same as DeleteWebhook, except this call does not require authentication.
//  Method                  DELETE
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Rate limiter            /webhooks/{webhook.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#delete-webhook-with-token
//  Reviewed                2018-08-14
//  Comment                 -
func (c *Client) DeleteWebhookWithToken(id Snowflake, token string, flags ...Flag) (err error) {
	var e string
	if token != "" {
		e = endpoint.WebhookToken(id, token)
	} else {
		e = endpoint.Webhook(id)
	}

	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: e,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// NewExecuteWebhookParams creates params for func ExecuteWebhook
func NewExecuteWebhookParams(id Snowflake, token string) (ret *ExecuteWebhookParams, err error) {
	return &ExecuteWebhookParams{
		WebhookID: id,
		Token:     token,
	}, nil
}

// ExecuteWebhookParams JSON params for func ExecuteWebhook
type ExecuteWebhookParams struct {
	WebhookID Snowflake `json:"-"`
	Token     string    `json:"-"`

	Content   string      `json:"content"`
	Username  string      `json:"username"`
	AvatarURL string      `json:"avatar_url"`
	TTS       bool        `json:"tts"`
	File      interface{} `json:"file"`
	Embeds    []*Embed    `json:"embeds"`
}

type execWebhookParams struct {
	Wait bool `urlparam:"wait"`
}

var _ URLQueryStringer = (*execWebhookParams)(nil)

// ExecuteWebhook [REST] Trigger a webhook in Discord.
//  Method                  POST
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Rate limiter            /webhooks/{webhook.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#execute-webhook
//  Reviewed                2018-08-14
//  Comment                 This endpoint. supports both JSON and form data bodies. It does require
//                          multipart/form-data requests instead of the normal JSON request type when
//                          uploading files. Make sure you set your Content-Type to multipart/form-data if
//                          you're doing that. Note that in that case, the embeds field cannot be used, but
//                          you can pass an url-encoded JSON body as a form value for payload_json.
//  Comment#2               For the webhook embed objects, you can set every field except type (it will be
//                          rich regardless of if you try to set it), provider, video, and any height, width,
//                          or proxy_url values for images.
func (c *Client) ExecuteWebhook(params *ExecuteWebhookParams, wait bool, URLSuffix string, flags ...Flag) (err error) {
	if params == nil {
		return errors.New("params can not be nil")
	}

	if params.WebhookID.IsZero() {
		return errors.New("webhook id is required")
	}
	if params.Token == "" {
		return errors.New("webhook token is required")
	}

	var contentType string
	if params.File == nil {
		contentType = httd.ContentTypeJSON
	} else {
		contentType = "multipart/form-data"
	}

	urlparams := &execWebhookParams{wait}
	r := c.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Endpoint:    endpoint.WebhookToken(params.WebhookID, params.Token) + URLSuffix + urlparams.URLQueryString(),
		Body:        params,
		ContentType: contentType,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent // TODO: verify

	_, err = r.Execute()
	return err
}

// ExecuteSlackWebhook [REST] Trigger a webhook in Discord from the Slack app.
//  Method                  POST
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Rate limiter            /webhooks
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#execute-slackcompatible-webhook
//  Reviewed                2018-08-14
//  Comment                 Refer to Slack's documentation for more information. We do not support Slack's channel,
//                          icon_emoji, mrkdwn, or mrkdwn_in properties.
func (c *Client) ExecuteSlackWebhook(params *ExecuteWebhookParams, wait bool, flags ...Flag) (err error) {
	return c.ExecuteWebhook(params, wait, endpoint.Slack(), flags...)
}

// ExecuteGitHubWebhook [REST] Trigger a webhook in Discord from the GitHub app.
//  Method                  POST
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Rate limiter            /webhooks
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#execute-githubcompatible-webhook
//  Reviewed                2018-08-14
//  Comment                 Add a new webhook to your GitHub repo (in the repo's settings), and use this endpoint.
//                          as the "Payload URL." You can choose what events your Discord channel receives by
//                          choosing the "Let me select individual events" option and selecting individual
//                          events for the new webhook you're configuring.
func (c *Client) ExecuteGitHubWebhook(params *ExecuteWebhookParams, wait bool, flags ...Flag) (err error) {
	return c.ExecuteWebhook(params, wait, endpoint.GitHub(), flags...)
}

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

// UpdateWebhookParams https://discordapp.com/developers/docs/resources/webhook#modify-webhook-json-params
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
