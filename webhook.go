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

var _ Copier = (*Webhook)(nil)
var _ DeepCopier = (*Webhook)(nil)

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

type WebhookQueryBuilder interface {
	WithContext(ctx context.Context) WebhookQueryBuilder
	WithFlags(flags ...Flag) WebhookQueryBuilder

	// Get Returns the new webhook object for the given id.
	Get() (*Webhook, error)

	// Update Modify a webhook. Requires the 'MANAGE_WEBHOOKS' permission.
	// Returns the updated webhook object on success.
	Update(*UpdateWebhook) (*Webhook, error)

	// Deprecated: use Update instead
	UpdateBuilder() UpdateWebhookBuilder

	// Delete Deletes a webhook permanently. User must be owner. Returns a 204 NO CONTENT response on success.
	Delete() error

	WithToken(token string) WebhookWithTokenQueryBuilder
}

func (c clientQueryBuilder) Webhook(id Snowflake) WebhookQueryBuilder {
	return &webhookQueryBuilder{client: c.client, webhookID: id}
}

type webhookQueryBuilder struct {
	ctx       context.Context
	flags     Flag
	client    *Client
	cid       Snowflake
	webhookID Snowflake
}

func (w *webhookQueryBuilder) validate() error {
	if w.client == nil {
		return ErrMissingClientInstance
	}
	if w.cid.IsZero() {
		return ErrMissingChannelID
	}
	if w.webhookID.IsZero() {
		return ErrMissingWebhookID
	}
	return nil
}

func (w webhookQueryBuilder) WithContext(ctx context.Context) WebhookQueryBuilder {
	w.ctx = ctx
	return &w
}

func (w webhookQueryBuilder) WithFlags(flags ...Flag) WebhookQueryBuilder {
	w.flags = mergeFlags(flags)
	return &w
}

// Get [REST] Returns the new webhook object for the given id.
//  Method                  GET
//  Endpoint                /webhooks/{webhook.id}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#get-webhook
//  Reviewed                2018-08-14
//  Comment                 -
func (w webhookQueryBuilder) Get() (ret *Webhook, err error) {
	r := w.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.Webhook(w.webhookID),
		Ctx:      w.ctx,
	}, w.flags)
	r.factory = func() interface{} {
		return &Webhook{}
	}

	return getWebhook(r.Execute)
}

// Update [REST] Modify a webhook. Requires the 'MANAGE_WEBHOOKS' permission.
// Returns the updated webhook object on success.
//  Method                  PATCH
//  Endpoint                /webhooks/{webhook.id}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#modify-webhook
//  Reviewed                2018-08-14
//  Comment                 All parameters to this endpoint.
func (w webhookQueryBuilder) Update(params *UpdateWebhook) (*Webhook, error) {
	if params == nil {
		return nil, ErrMissingRESTParams
	}
	if err := w.validate(); err != nil {
		return nil, err
	}

	r := w.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ctx:         w.ctx,
		Endpoint:    endpoint.Webhook(w.webhookID),
		ContentType: httd.ContentTypeJSON,
		Body:        params,
	}, w.flags)
	r.factory = func() interface{} {
		return &Webhook{}
	}

	return getWebhook(r.Execute)
}

// Delete [REST] Delete a webhook permanently. User must be owner. Returns a 204 NO CONTENT response on success.
//  Method                  DELETE
//  Endpoint                /webhooks/{webhook.id}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#delete-webhook
//  Reviewed                2018-08-14
//  Comment                 -
func (w webhookQueryBuilder) Delete() (err error) {
	return w.WithToken("").WithFlags(w.flags).WithContext(w.ctx).Delete()
}

// ExecuteWebhook JSON params for func ExecuteWebhook
type ExecuteWebhook struct {
	Content   string      `json:"content"`
	Username  string      `json:"username"`
	AvatarURL string      `json:"avatar_url"`
	TTS       bool        `json:"tts"`
	File      interface{} `json:"file"`
	Embeds    []*Embed    `json:"embeds"`
}

type execWebhook struct {
	Wait     *bool      `urlparam:"wait"`
	ThreadID *Snowflake `urlparam:"thread_id"`
}

var _ URLQueryStringer = (*execWebhook)(nil)

type WebhookWithTokenQueryBuilder interface {
	WithContext(ctx context.Context) WebhookWithTokenQueryBuilder
	WithFlags(flags ...Flag) WebhookWithTokenQueryBuilder

	// Get Same as GetWebhook, except this call does not require authentication and
	// returns no user in the webhook object.
	Get() (*Webhook, error)

	// Update Same as UpdateWebhook, except this call does not require authentication,
	// does _not_ accept a channel_id parameter in the body, and does not return a user in the webhook object.
	Update(*UpdateWebhook) (*Webhook, error)

	// Deprecated: use Update instead
	UpdateBuilder() UpdateWebhookBuilder

	// Delete Same as DeleteWebhook, except this call does not require authentication.
	Delete() error

	// Execute Trigger a webhook in Discord.
	Execute(params *ExecuteWebhook, wait *bool, threadID *Snowflake, URLSuffix string) (*Message, error)

	// ExecuteSlackWebhook Trigger a webhook in Discord from the Slack app.
	ExecuteSlackWebhook(params *ExecuteWebhook, wait *bool, threadID *Snowflake) (*Message, error)

	// ExecuteGitHubWebhook Trigger a webhook in Discord from the GitHub app.
	ExecuteGitHubWebhook(params *ExecuteWebhook, wait *bool, threadID *Snowflake) (*Message, error)

	// GetMessage gets a previously sent message from this webhook.
	GetMessage(messageId Snowflake, threadID *Snowflake) (*Message, error)

	// EditMessage edits a previously sent message from this webhook.
	EditMessage(params *ExecuteWebhook, messageId Snowflake, threadID *Snowflake) (*Message, error)

	// DeleteMessage Deletes a previously sent message from this webhook.
	DeleteMessage(messageId Snowflake, threadID *Snowflake) error
}

func (w webhookQueryBuilder) WithToken(token string) WebhookWithTokenQueryBuilder {
	return &webhookWithTokenQueryBuilder{client: w.client, webhookID: w.webhookID, token: token}
}

type webhookWithTokenQueryBuilder struct {
	ctx       context.Context
	flags     Flag
	client    *Client
	cid       Snowflake
	webhookID Snowflake
	token     string
}

func (w *webhookWithTokenQueryBuilder) validate() error {
	if w.client == nil {
		return ErrMissingClientInstance
	}
	if w.cid.IsZero() {
		return ErrMissingChannelID
	}
	if w.webhookID.IsZero() {
		return ErrMissingWebhookID
	}
	if w.token == "" {
		return ErrMissingWebhookToken
	}
	return nil
}

func (w webhookWithTokenQueryBuilder) WithContext(ctx context.Context) WebhookWithTokenQueryBuilder {
	w.ctx = ctx
	return &w
}

func (w webhookWithTokenQueryBuilder) WithFlags(flags ...Flag) WebhookWithTokenQueryBuilder {
	w.flags = mergeFlags(flags)
	return &w
}

// Get [REST] Same as GetWebhook, except this call does not require authentication and
// returns no user in the webhook object.
//  Method                  GET
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#get-webhook-with-token
//  Reviewed                2018-08-14
//  Comment                 -
func (w webhookWithTokenQueryBuilder) Get() (*Webhook, error) {
	r := w.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.WebhookToken(w.webhookID, w.token),
		Ctx:      w.ctx,
	}, w.flags)
	r.factory = func() interface{} {
		return &Webhook{}
	}

	return getWebhook(r.Execute)
}

// Update [REST] Same as UpdateWebhook, except this call does not require authentication,
// does _not_ accept a channel_id parameter in the body, and does not return a user in the webhook object.
//  Method                  PATCH
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#modify-webhook-with-token
//  Reviewed                2018-08-14
//  Comment                 All parameters to this endpoint. are optional.
func (w webhookWithTokenQueryBuilder) Update(params *UpdateWebhook) (*Webhook, error) {
	if params == nil {
		return nil, ErrMissingRESTParams
	}
	if err := w.validate(); err != nil {
		return nil, err
	}

	r := w.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ctx:         w.ctx,
		Endpoint:    endpoint.WebhookToken(w.webhookID, w.token),
		ContentType: httd.ContentTypeJSON,
		Body:        params,
	}, w.flags)
	r.factory = func() interface{} {
		return &Webhook{}
	}

	return getWebhook(r.Execute)
}

type UpdateWebhook struct {
	Name      *string    `json:"name,omitempty"`
	Avatar    *string    `json:"avatar,omitempty"`
	ChannelID *Snowflake `json:"channel_id,omitempty"`
}

// Delete [REST] Same as DeleteWebhook, except this call does not require authentication.
//  Method                  DELETE
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#delete-webhook-with-token
//  Reviewed                2018-08-14
//  Comment                 -
func (w webhookWithTokenQueryBuilder) Delete() error {
	var e string
	if w.token != "" {
		e = endpoint.WebhookToken(w.webhookID, w.token)
	} else {
		e = endpoint.Webhook(w.webhookID)
	}

	r := w.client.newRESTRequest(&httd.Request{
		Method:   http.MethodDelete,
		Endpoint: e,
		Ctx:      w.ctx,
	}, w.flags)

	_, err := r.Execute()
	return err
}

// Execute [REST] Trigger a webhook in Discord.
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
func (w webhookWithTokenQueryBuilder) Execute(params *ExecuteWebhook, wait *bool, threadID *Snowflake, URLSuffix string) (message *Message, err error) {
	if params == nil {
		return nil, errors.New("params can not be nil")
	}

	if w.webhookID.IsZero() {
		return nil, ErrMissingWebhookID
	}
	if w.token == "" {
		return nil, ErrMissingWebhookToken
	}

	var contentType string
	if params.File == nil {
		contentType = httd.ContentTypeJSON
	} else {
		contentType = "multipart/form-data"
	}

	urlparams := &execWebhook{
		Wait:     wait,
		ThreadID: threadID,
	}
	r := w.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPost,
		Ctx:         w.ctx,
		Endpoint:    endpoint.WebhookToken(w.webhookID, w.token) + URLSuffix + urlparams.URLQueryString(),
		Body:        params,
		ContentType: contentType,
	}, w.flags)
	// Discord only returns the message when wait=true.
	if wait != nil && *wait == true {
		r.pool = w.client.pool.message
		return getMessage(r.Execute)
	}
	_, err = r.Execute()
	return nil, err
}

// ExecuteSlackWebhook [REST] Trigger a webhook in Discord from the Slack app.
//  Method                  POST
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#execute-slackcompatible-webhook
//  Reviewed                2020-05-21
//  Comment                 Refer to Slack's documentation for more information. We do not support Slack's channel,
//                          icon_emoji, mrkdwn, or mrkdwn_in properties.
func (w *webhookWithTokenQueryBuilder) ExecuteSlackWebhook(params *ExecuteWebhook, wait *bool, threadID *Snowflake) (*Message, error) {
	return w.Execute(params, wait, threadID, endpoint.Slack())
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
func (w *webhookWithTokenQueryBuilder) ExecuteGitHubWebhook(params *ExecuteWebhook, wait *bool, threadID *Snowflake) (*Message, error) {
	return w.Execute(params, wait, threadID, endpoint.GitHub())
}

// GetMessage [REST] Gets a Message sent from this Webhook.
//  Method                  GET
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#get-webhook-message
//  Reviewed                2021-02-16
//  Comment                 Returns a previously-sent webhook message from the same token. Returns a message object on success.
func (w *webhookWithTokenQueryBuilder) GetMessage(messageId Snowflake, threadID *Snowflake) (*Message, error) {
	urlparams := &execWebhook{
		ThreadID: threadID,
	}

	r := w.client.newRESTRequest(&httd.Request{
		Method:      http.MethodGet,
		Ctx:         w.ctx,
		Endpoint:    endpoint.WebhookMessage(w.webhookID, w.token, messageId) + urlparams.URLQueryString(),
	}, w.flags)
	r.pool = w.client.pool.message
	return getMessage(r.Execute)
}


// EditMessage [REST] Edits a previously-sent Webhook Message from the same token. Returns a Message object on success.
//  Method                  PATCH
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#execute-webhook
//  Reviewed                2020-05-21
//  Comment                 When the content field is edited, the mentions array in the message object will be reconstructed from scratch based on the new content.
// 							The allowed_mentions field of the edit request controls how this happens.
//							If there is no explicit allowed_mentions in the edit request, the content will be parsed with default allowances, that is, without regard to whether or not an allowed_mentions was present in the request that originally created the message.
// 							Refer to Uploading Files for details on attachments and multipart/form-data requests.
// 							Any provided files will be appended to the message.
//							To remove or replace files you will have to supply the attachments field which specifies the files to retain on the message after edit.
// Comment #2				All parameters to this endpoint are optional and nullable.
func (w *webhookWithTokenQueryBuilder) EditMessage(params *ExecuteWebhook, messageId Snowflake, threadID *Snowflake) (*Message, error) {
	if params == nil {
		return nil, errors.New("params can not be nil")
	}

	if w.webhookID.IsZero() {
		return nil, ErrMissingWebhookID
	}
	if w.token == "" {
		return nil, ErrMissingWebhookToken
	}

	var contentType string
	if params.File == nil {
		contentType = httd.ContentTypeJSON
	} else {
		contentType = "multipart/form-data"
	}

	urlparams := &execWebhook{
		ThreadID: threadID,
	}

	r := w.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ctx:         w.ctx,
		Endpoint:    endpoint.WebhookMessage(w.webhookID, w.token, messageId) + urlparams.URLQueryString(),
		Body:        params,
		ContentType: contentType,
	}, w.flags)

	r.pool = w.client.pool.message
	return getMessage(r.Execute)
}

// DeleteMessage [REST] Deletes a Message sent from this Webhook.
//  Method                  DELETE
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
//  Discord documentation   https://discord.com/developers/docs/resources/webhook#delete-webhook-message
//  Reviewed                2021-02-16
//  Comment                 Deletes a message that was created by the webhook. Returns a 204 No Content response on success.
func (w *webhookWithTokenQueryBuilder) DeleteMessage(messageId Snowflake, threadID *Snowflake) error {
	urlparams := &execWebhook{
		ThreadID: threadID,
	}

	r := w.client.newRESTRequest(&httd.Request{
		Method:      http.MethodDelete,
		Ctx:         w.ctx,
		Endpoint:    endpoint.WebhookMessage(w.webhookID, w.token, messageId) + urlparams.URLQueryString(),
	}, w.flags)
	r.pool = w.client.pool.message
	_, err := r.Execute()
	if err != nil {
		return err
	}
	return nil
}