package disgord

import (
	"errors"
	"net/http"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
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

func ratelimitWebhook(id Snowflake) string {
	return "wh:" + id.String()
}

func NewCreateWebhookParams(name string, avatar *string) *CreateWebhookParams {
	return &CreateWebhookParams{
		Name:   name,
		Avatar: avatar,
	}
}

// CreateWebhookParams json params for the create webhook rest request avatar string
// https://discordapp.com/developers/docs/resources/user#avatar-data
type CreateWebhookParams struct {
	Name   string  `json:"name"`   // name of the webhook (2-32 characters)
	Avatar *string `json:"avatar"` // avatar data uri scheme, image for the default webhook avatar
}

// CreateWebhook [REST] Create a new webhook. Requires the 'MANAGE_WEBHOOKS' permission.
// Returns a webhook object on success.
//  Method                  POST
//  Endpoint                /channels/{channel.id}/webhooks
//  Rate limiter            /channels/{channel.id}/webhooks
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#create-webhook
//  Reviewed                2018-08-14
//  Comment                 -
func CreateWebhook(client httd.Poster, channelID Snowflake, params *CreateWebhookParams) (ret *Webhook, err error) {
	var body []byte
	_, body, err = client.Post(&httd.Request{
		Ratelimiter: ratelimitChannelWebhooks(channelID),
		Endpoint:    endpoint.ChannelWebhooks(channelID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetChannelWebhooks [REST] Returns a list of channel webhook objects. Requires the 'MANAGE_WEBHOOKS' permission.
//  Method                  POST
//  Endpoint                /channels/{channel.id}/webhooks
//  Rate limiter            /channels/{channel.id}/webhooks
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#get-channel-webhooks
//  Reviewed                2018-08-14
//  Comment                 -
func GetChannelWebhooks(client httd.Getter, channelID Snowflake) (ret []*Webhook, err error) {
	var body []byte
	_, body, err = client.Get(&httd.Request{
		Ratelimiter: ratelimitChannelWebhooks(channelID),
		Endpoint:    endpoint.ChannelWebhooks(channelID),
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetGuildWebhooks [REST] Returns a list of guild webhook objects. Requires the 'MANAGE_WEBHOOKS' permission.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/webhooks
//  Rate limiter            /guilds/{guild.id}/webhooks
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#get-guild-webhooks
//  Reviewed                2018-08-14
//  Comment                 -
func GetGuildWebhooks(client httd.Getter, guildID Snowflake) (ret []*Webhook, err error) {
	var body []byte
	_, body, err = client.Get(&httd.Request{
		Ratelimiter: ratelimitGuildWebhooks(guildID),
		Endpoint:    endpoint.GuildWebhooks(guildID),
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetWebhook [REST] Returns the new webhook object for the given id.
//  Method                  GET
//  Endpoint                /webhooks/{webhook.id}
//  Rate limiter            /webhooks/{webhook.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#get-webhook
//  Reviewed                2018-08-14
//  Comment                 -
func GetWebhook(client httd.Getter, id Snowflake) (ret *Webhook, err error) {
	var body []byte
	_, body, err = client.Get(&httd.Request{
		Ratelimiter: ratelimitWebhook(id),
		Endpoint:    endpoint.Webhook(id),
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetWebhookWithToken [REST] Same as GetWebhook, except this call does not require authentication and
// returns no user in the webhook object.
//  Method                  GET
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Rate limiter            /webhooks/{webhook.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#get-webhook-with-token
//  Reviewed                2018-08-14
//  Comment                 -
func GetWebhookWithToken(client httd.Getter, id Snowflake, token string) (ret *Webhook, err error) {
	var body []byte
	_, body, err = client.Get(&httd.Request{
		Ratelimiter: ratelimitWebhook(id),
		Endpoint:    endpoint.WebhookToken(id, token),
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// UpdateWebhookParams https://discordapp.com/developers/docs/resources/webhook#modify-webhook-json-params
// Allows changing the name of the webhook, avatar and moving it to another channel. It also allows to resetting the
// avatar by providing a nil to SetAvatar.
//  params = &UpdateWebhookParams{}
//  params.UseDefaultAvatar() // will reset any image data, if present
type UpdateWebhookParams struct {
	avatarIsSet bool
	name        string
	avatar      string
	channelID   Snowflake
}

var _ AvatarParamHolder = (*UpdateWebhookParams)(nil)

func NewUpdateWebhookParams() *UpdateWebhookParams {
	return &UpdateWebhookParams{}
}

func (m *UpdateWebhookParams) Empty() bool {
	return m.name == "" && m.channelID.Empty() && !m.avatarIsSet
}

func (m *UpdateWebhookParams) SetName(name string) {
	m.name = name
}

// SetAvatar updates the avatar image. Must be abase64 encoded string.
// provide a nil to reset the avatar.
func (m *UpdateWebhookParams) SetAvatar(avatar string) {
	m.avatar = avatar
	m.avatarIsSet = avatar != ""
}

// UseDefaultAvatar sets the avatar param to null, and let's Discord assign a default avatar image.
// Note that the avatar value will never hold content, as default avatars only works on null values.
//
// Use this to reset an avatar image.
func (m *UpdateWebhookParams) UseDefaultAvatar() {
	m.avatar = ""
	m.avatarIsSet = true
}
func (m *UpdateWebhookParams) SetChannelID(channelID Snowflake) {
	m.channelID = channelID
}

func (m *UpdateWebhookParams) MarshalJSON() ([]byte, error) {
	var v interface{}
	if m.avatarIsSet {
		p := &modifyWebhookParamsWithAvatar{
			Name:      m.name,
			ChannelID: m.channelID,
		}
		if m.avatar != "" {
			p.Avatar = &m.avatar
		}

		v = p
	} else {
		v = &modifyWebhookParamsWithoutAvatar{
			Name:      m.name,
			ChannelID: m.channelID,
		}
	}

	return marshal(v)
}

type modifyWebhookParamsWithoutAvatar struct {
	Name      string    `json:"name,omitempty"` // name of the webhook (2-32 characters)
	ChannelID Snowflake `json:"channel_id,omitempty"`
}

type modifyWebhookParamsWithAvatar struct {
	Name      string    `json:"name,omitempty"` // name of the webhook (2-32 characters)
	Avatar    *string   `json:"avatar"`         // avatar data uri scheme, image for the default webhook avatar
	ChannelID Snowflake `json:"channel_id,omitempty"`
}

// UpdateWebhook [REST] Modify a webhook. Requires the 'MANAGE_WEBHOOKS' permission.
// Returns the updated webhook object on success.
//  Method                  PATCH
//  Endpoint                /webhooks/{webhook.id}
//  Rate limiter            /webhooks/{webhook.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#modify-webhook
//  Reviewed                2018-08-14
//  Comment                 All parameters to this endpoint.
func UpdateWebhook(client httd.Patcher, id Snowflake, params *UpdateWebhookParams) (ret *Webhook, err error) {
	if id.Empty() {
		err = errors.New("not a valid snowflake")
		return
	}

	var body []byte
	_, body, err = client.Patch(&httd.Request{
		Ratelimiter: ratelimitWebhook(id),
		Endpoint:    endpoint.Webhook(id),
		ContentType: httd.ContentTypeJSON,
		Body:        params,
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// UpdateWebhookWithToken [REST] Same as ModifyWebhook, except this call does not require authentication,
// does not accept a channel_id parameter in the body, and does not return a user in the webhook object.
//  Method                  PATCH
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Rate limiter            /webhooks/{webhook.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#modify-webhook-with-token
//  Reviewed                2018-08-14
//  Comment                 All parameters to this endpoint. are optional. Not tested:extra json fields might cause
//                          an issue. Consider writing a json params object.
func UpdateWebhookWithToken(client httd.Patcher, newWebhook *Webhook) (ret *Webhook, err error) {
	id := newWebhook.ID
	if id.Empty() {
		err = errors.New("not a valid snowflake")
		return
	}

	var body []byte
	_, body, err = client.Patch(&httd.Request{
		Ratelimiter: ratelimitWebhook(id),
		Endpoint:    endpoint.WebhookToken(id, newWebhook.Token),
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// DeleteWebhook [REST] Delete a webhook permanently. User must be owner. Returns a 204 NO CONTENT response on success.
//  Method                  DELETE
//  Endpoint                /webhooks/{webhook.id}
//  Rate limiter            /webhooks/{webhook.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#delete-webhook
//  Reviewed                2018-08-14
//  Comment                 -
func DeleteWebhook(client httd.Deleter, webhookID Snowflake) (err error) {
	return DeleteWebhookWithToken(client, webhookID, "")
}

// DeleteWebhookWithToken [REST] Same as DeleteWebhook, except this call does not require authentication.
//  Method                  DELETE
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Rate limiter            /webhooks/{webhook.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#delete-webhook-with-token
//  Reviewed                2018-08-14
//  Comment                 -
func DeleteWebhookWithToken(client httd.Deleter, id Snowflake, token string) (err error) {
	var e string
	if token != "" {
		e = endpoint.WebhookToken(id, token)
	} else {
		e = endpoint.Webhook(id)
	}

	var resp *http.Response
	resp, _, err = client.Delete(&httd.Request{
		Ratelimiter: ratelimitWebhook(id),
		Endpoint:    e,
	})
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(http.StatusNoContent)
		err = errors.New(msg)
	}
	return
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

	Content   string          `json:"content"`
	Username  string          `json:"username"`
	AvatarURL string          `json:"avatar_url"`
	TTS       bool            `json:"tts"`
	File      interface{}     `json:"file"`
	Embeds    []*ChannelEmbed `json:"embeds"`
}

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
// TODO
func ExecuteWebhook(client httd.Poster, params *ExecuteWebhookParams, wait bool, URLSuffix string) (err error) {
	_, _, err = client.Post(&httd.Request{
		Ratelimiter: ratelimitWebhook(params.WebhookID),
		Endpoint:    endpoint.WebhookToken(params.WebhookID, params.Token) + URLSuffix,
		ContentType: httd.ContentTypeJSON,
	})
	return // TODO: how to verify success?
}

// ExecuteSlackWebhook [REST] Trigger a webhook in Discord from the Slack app.
//  Method                  POST
//  Endpoint                /webhooks/{webhook.id}/{webhook.token}
//  Rate limiter            /webhooks
//  Discord documentation   https://discordapp.com/developers/docs/resources/webhook#execute-slackcompatible-webhook
//  Reviewed                2018-08-14
//  Comment                 Refer to Slack's documentation for more information. We do not support Slack's channel,
//                          icon_emoji, mrkdwn, or mrkdwn_in properties.
func ExecuteSlackWebhook(client httd.Poster, params *ExecuteWebhookParams, wait bool) (err error) {
	return ExecuteWebhook(client, params, wait, endpoint.Slack())
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
func ExecuteGitHubWebhook(client httd.Poster, params *ExecuteWebhookParams, wait bool) (err error) {
	return ExecuteWebhook(client, params, wait, endpoint.GitHub())
}

// CreateWebhook .
func (c *client) CreateWebhook(channelID Snowflake, params *CreateWebhookParams, flags ...Flag) (ret *Webhook, err error) {
	ret, err = CreateWebhook(c.req, channelID, params)
	return
}

// GetChannelWebhooks .
func (c *client) GetChannelWebhooks(channelID Snowflake, flags ...Flag) (ret []*Webhook, err error) {
	ret, err = GetChannelWebhooks(c.req, channelID)
	return
}

// GetGuildWebhooks .
func (c *client) GetGuildWebhooks(guildID Snowflake, flags ...Flag) (ret []*Webhook, err error) {
	ret, err = GetGuildWebhooks(c.req, guildID)
	return
}

// GetWebhook .
func (c *client) GetWebhook(id Snowflake, flags ...Flag) (ret *Webhook, err error) {
	ret, err = GetWebhook(c.req, id)
	return
}

// GetWebhookWithToken .
func (c *client) GetWebhookWithToken(id Snowflake, token string, flags ...Flag) (ret *Webhook, err error) {
	ret, err = GetWebhookWithToken(c.req, id, token)
	return
}

// ModifyWebhook .
func (c *client) UpdateWebhook(id Snowflake, params *UpdateWebhookParams, flags ...Flag) (ret *Webhook, err error) {
	if id.Empty() {
		err = errors.New("given webhook ID was not set, there is nothing to modify")
		return
	}
	if params == nil {
		err = errors.New("given param object was nil, there is nothing to modify")
		return
	}
	if params.Empty() {
		err = errors.New("given param object was empty, there is nothing to modify")
		return
	}

	// verify avatar string prefix
	if params.avatarIsSet && params.avatar != "" && !validAvatarPrefix(params.avatar) {
		err = errors.New("given avatar string is invalid. Must specify data encoding. Eg. 'data:image/jpeg;base64,'")
		return
	}

	// TODO: check if user has permission to modify webhook
	ret, err = UpdateWebhook(c.req, id, params)
	return
}

// ModifyWebhookWithToken .
func (c *client) UpdateWebhookWithToken(newWebhook *Webhook, flags ...Flag) (ret *Webhook, err error) {
	ret, err = UpdateWebhookWithToken(c.req, newWebhook)
	return
}

// DeleteWebhook .
func (c *client) DeleteWebhook(webhookID Snowflake, flags ...Flag) (err error) {
	err = DeleteWebhook(c.req, webhookID)
	return
}

// DeleteWebhookWithToken .
func (c *client) DeleteWebhookWithToken(id Snowflake, token string, flags ...Flag) (err error) {
	err = DeleteWebhookWithToken(c.req, id, token)
	return
}

// ExecuteWebhook .
func (c *client) ExecuteWebhook(params *ExecuteWebhookParams, wait bool, URLSuffix string, flags ...Flag) (err error) {
	err = ExecuteWebhook(c.req, params, wait, URLSuffix)
	return
}

// ExecuteSlackWebhook .
func (c *client) ExecuteSlackWebhook(params *ExecuteWebhookParams, wait bool, flags ...Flag) (err error) {
	err = ExecuteSlackWebhook(c.req, params, wait)
	return
}

// ExecuteGitHubWebhook .
func (c *client) ExecuteGitHubWebhook(params *ExecuteWebhookParams, wait bool, flags ...Flag) (err error) {
	err = ExecuteGitHubWebhook(c.req, params, wait)
	return
}
