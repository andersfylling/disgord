package disgord

import (
	"context"
	"errors"
	"net/http"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
)

// Webhook Used to represent a webhook
// https://discord.com/developers/docs/resources/webhook#webhook-object
type Webhook struct {
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

	hook.ID = w.ID
	hook.GuildID = w.GuildID
	hook.ChannelID = w.ChannelID
	hook.User = w.User.DeepCopy().(*User)
	hook.Name = w.Name
	hook.Avatar = w.Avatar
	hook.Token = w.Token

	return
}

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

type WebhookQueryBuilder interface {
	WithContext(ctx context.Context) WebhookQueryBuilder

	// GetWebhook Returns the new webhook object for the given id.
	Get(flags ...Flag) (*Webhook, error)

	// UpdateWebhook Modify a webhook. Requires the 'MANAGE_WEBHOOKS' permission.
	// Returns the updated webhook object on success.
	Update(flags ...Flag) *updateWebhookBuilder

	// DeleteWebhook Delete a webhook permanently. User must be owner. Returns a 204 NO CONTENT response on success.
	Delete(flags ...Flag) error

	// ExecuteWebhook Trigger a webhook in Discord.
	Execute(params *ExecuteWebhookParams, wait bool, URLSuffix string, flags ...Flag) (*Message, error)

	// ExecuteSlackWebhook Trigger a webhook in Discord from the Slack app.
	ExecuteSlackWebhook(params *ExecuteWebhookParams, wait bool, flags ...Flag) (*Message, error)

	// ExecuteGitHubWebhook Trigger a webhook in Discord from the GitHub app.
	ExecuteGitHubWebhook(params *ExecuteWebhookParams, wait bool, flags ...Flag) (*Message, error)

	WithToken(token string) WebhookWithTokenQueryBuilder
}

func (c clientQueryBuilder) Webhook(id Snowflake) WebhookQueryBuilder {
	return &webhookQueryBuilder{client: c.client, webhookID: id}
}

type webhookQueryBuilder struct {
	ctx       context.Context
	client    *Client
	cid       Snowflake
	webhookID Snowflake
}

func (w webhookQueryBuilder) WithContext(ctx context.Context) WebhookQueryBuilder {
	w.ctx = ctx
	return &w
}

// GetWebhook [REST] Returns the new webhook object for the given id.
//  Method                  GET
//  Endpoint                /webhooks/{webhook.id}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#get-webhook
//  Reviewed                2018-08-14
//  Comment                 -
func (w webhookQueryBuilder) Get(flags ...Flag) (ret *Webhook, err error) {
	r := w.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.Webhook(w.webhookID),
		Ctx:      w.ctx,
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
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#modify-webhook
//  Reviewed                2018-08-14
//  Comment                 All parameters to this endpoint.
func (w webhookQueryBuilder) Update(flags ...Flag) (builder *updateWebhookBuilder) {
	builder = &updateWebhookBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Webhook{}
	}
	builder.r.flags = flags
	builder.r.addPrereq(w.webhookID.IsZero(), "given webhook ID was not set, there is nothing to modify")
	builder.r.setup(w.client.req, &httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         w.ctx,
		Endpoint:    endpoint.Webhook(w.webhookID),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// DeleteWebhook [REST] Delete a webhook permanently. User must be owner. Returns a 204 NO CONTENT response on success.
//  Method                  DELETE
//  Endpoint                /webhooks/{webhook.id}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#delete-webhook
//  Reviewed                2018-08-14
//  Comment                 -
func (w webhookQueryBuilder) Delete(flags ...Flag) (err error) {
	return w.WithToken("").WithContext(w.ctx).Delete(flags...)
}

// ExecuteWebhookParams JSON params for func ExecuteWebhook
type ExecuteWebhookParams struct {
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

// Execute Trigger a webhook in Discord.
func (w webhookQueryBuilder) Execute(params *ExecuteWebhookParams, wait bool, URLSuffix string, flags ...Flag) (message *Message, err error) {
	return w.WithToken("").WithContext(w.ctx).Execute(params, wait, URLSuffix, flags...)
}

// ExecuteSlackWebhook [REST] Trigger a webhook in Discord from the Slack app.
//  Method                  POST
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#execute-slackcompatible-webhook
//  Reviewed                2020-05-21
//  Comment                 Refer to Slack's documentation for more information. We do not support Slack's channel,
//                          icon_emoji, mrkdwn, or mrkdwn_in properties.
func (w webhookQueryBuilder) ExecuteSlackWebhook(params *ExecuteWebhookParams, wait bool, flags ...Flag) (message *Message, err error) {
	return w.WithToken("").WithContext(w.ctx).Execute(params, wait, endpoint.Slack(), flags...)
}

// ExecuteGitHubWebhook [REST] Trigger a webhook in Discord from the GitHub app.
//  Method                  POST
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#execute-githubcompatible-webhook
//  Reviewed                2020-05-21
//  Comment                 Add a new webhook to your GitHub repo (in the repo's settings), and use this endpoint.
//                          as the "Payload URL." You can choose what events your Discord channel receives by
//                          choosing the "Let me select individual events" option and selecting individual
//                          events for the new webhook you're configuring.
func (w webhookQueryBuilder) ExecuteGitHubWebhook(params *ExecuteWebhookParams, wait bool, flags ...Flag) (message *Message, err error) {
	return w.WithToken("").WithContext(w.ctx).Execute(params, wait, endpoint.GitHub(), flags...)
}

type WebhookWithTokenQueryBuilder interface {
	WithContext(ctx context.Context) WebhookWithTokenQueryBuilder

	// GetWebhookWithToken Same as GetWebhook, except this call does not require authentication and
	// returns no user in the webhook object.
	Get(flags ...Flag) (*Webhook, error)

	// UpdateWebhookWithToken Same as UpdateWebhook, except this call does not require authentication,
	// does _not_ accept a channel_id parameter in the body, and does not return a user in the webhook object.
	Update(flags ...Flag) *updateWebhookBuilder

	// DeleteWebhookWithToken Same as DeleteWebhook, except this call does not require authentication.
	Delete(flags ...Flag) error

	Execute(params *ExecuteWebhookParams, wait bool, URLSuffix string, flags ...Flag) (*Message, error)
}

func (w webhookQueryBuilder) WithToken(token string) WebhookWithTokenQueryBuilder {
	return &webhookWithTokenQueryBuilder{client: w.client, webhookID: w.webhookID, token: token}
}

type webhookWithTokenQueryBuilder struct {
	ctx       context.Context
	client    *Client
	cid       Snowflake
	webhookID Snowflake
	token     string
}

func (w webhookWithTokenQueryBuilder) WithContext(ctx context.Context) WebhookWithTokenQueryBuilder {
	w.ctx = ctx
	return &w
}

// GetWebhookWithToken [REST] Same as GetWebhook, except this call does not require authentication and
// returns no user in the webhook object.
//  Method                  GET
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#get-webhook-with-token
//  Reviewed                2018-08-14
//  Comment                 -
func (w webhookWithTokenQueryBuilder) Get(flags ...Flag) (*Webhook, error) {
	r := w.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.WebhookToken(w.webhookID, w.token),
		Ctx:      w.ctx,
	}, flags)
	r.factory = func() interface{} {
		return &Webhook{}
	}

	return getWebhook(r.Execute)
}

// UpdateWebhookWithToken [REST] Same as UpdateWebhook, except this call does not require authentication,
// does _not_ accept a channel_id parameter in the body, and does not return a user in the webhook object.
//  Method                  PATCH
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#modify-webhook-with-token
//  Reviewed                2018-08-14
//  Comment                 All parameters to this endpoint. are optional.
func (w webhookWithTokenQueryBuilder) Update(flags ...Flag) (builder *updateWebhookBuilder) {
	builder = &updateWebhookBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Webhook{}
	}
	builder.r.flags = flags
	builder.r.addPrereq(w.webhookID.IsZero(), "given webhook ID was not set, there is nothing to modify")
	builder.r.addPrereq(w.token == "", "given webhook token was not set")
	builder.r.setup(w.client.req, &httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         w.ctx,
		Endpoint:    endpoint.WebhookToken(w.webhookID, w.token),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// DeleteWebhookWithToken [REST] Same as DeleteWebhook, except this call does not require authentication.
//  Method                  DELETE
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#delete-webhook-with-token
//  Reviewed                2018-08-14
//  Comment                 -
func (w webhookWithTokenQueryBuilder) Delete(flags ...Flag) error {
	var e string
	if w.token != "" {
		e = endpoint.WebhookToken(w.webhookID, w.token)
	} else {
		e = endpoint.Webhook(w.webhookID)
	}

	r := w.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: e,
		Ctx:      w.ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err := r.Execute()
	return err
}

// ExecuteWebhook [REST] Trigger a webhook in Discord.
//  Method                  POST
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#execute-webhook
//  Reviewed                2020-05-21
//  Comment                 This endpoint. supports both JSON and form data bodies. It does require
//                          multipart/form-data requests instead of the normal JSON request type when
//                          uploading files. Make sure you set your Content-Type to multipart/form-data if
//                          you're doing that. Note that in that case, the embeds field cannot be used, but
//                          you can pass an url-encoded JSON body as a form value for payload_json.
//  Comment#2               For the webhook embed objects, you can set every field except type (it will be
//                          rich regardless of if you try to set it), provider, video, and any height, width,
//                          or proxy_url values for images.
func (w webhookWithTokenQueryBuilder) Execute(params *ExecuteWebhookParams, wait bool, URLSuffix string, flags ...Flag) (message *Message, err error) {
	if params == nil {
		return nil, errors.New("params can not be nil")
	}

	if w.webhookID.IsZero() {
		return nil, errors.New("webhook id is required")
	}
	if w.token == "" {
		return nil, errors.New("webhook token is required")
	}

	var contentType string
	if params.File == nil {
		contentType = httd.ContentTypeJSON
	} else {
		contentType = "multipart/form-data"
	}

	urlparams := &execWebhookParams{wait}
	r := w.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Ctx:         w.ctx,
		Endpoint:    endpoint.WebhookToken(w.webhookID, w.token) + URLSuffix + urlparams.URLQueryString(),
		Body:        params,
		ContentType: contentType,
	}, flags)
	// Discord only returns the message when wait=true.
	if wait {
		r.pool = w.client.pool.message
		r.expectsStatusCode = http.StatusOK
		return getMessage(r.Execute)
	}
	r.expectsStatusCode = http.StatusNoContent
	_, err = r.Execute()
	return nil, err
}

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

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
