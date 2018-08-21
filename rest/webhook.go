package rest

import (
	"errors"
	"net/http"

	. "github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/disgord/rest/httd"
	"github.com/andersfylling/snowflake"
)

const (
	EndpointWebhooks      = "/webhooks"
	EndpointSlackWebhook  = "/slack"
	EndpointGitHubWebhook = "/github"
)

// CreateWebhookParams json params for the create webhook rest request
// avatar string: https://discordapp.com/developers/docs/resources/user#avatar-data
type CreateWebhookParams struct {
	Name   string `json:"name"`   // name of the webhook (2-32 characters)
	Avatar string `json:"avatar"` // avatar data uri scheme, image for the default webhook avatar
}

// CreateWebhook [POST]     Create a new webhook. Requires the 'MANAGE_WEBHOOKS' permission.
//                          Returns a webhook object on success.
// Endpoint                 /channels/{channel.id}/webhooks
// Rate limiter             /channels/{channel.id}/webhooks
// Discord documentation    https://discordapp.com/developers/docs/resources/webhook#create-webhook
// Reviewed                 2018-08-14
// Comment                  -
func CreateWebhook(client httd.Poster, channelID snowflake.ID, params *CreateWebhookParams) (ret *Webhook, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannelWebhooks(channelID),
		Endpoint:    EndpointChannels + "/" + channelID.String() + EndpointWebhooks,
		JSONParams:  params,
	}
	_, body, err := client.Post(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetChannelWebhooks [GET] Returns a list of channel webhook objects. Requires the 'MANAGE_WEBHOOKS' permission.
// Endpoint                 /channels/{channel.id}/webhooks
// Rate limiter             /channels/{channel.id}/webhooks
// Discord documentation    https://discordapp.com/developers/docs/resources/webhook#get-channel-webhooks
// Reviewed                 2018-08-14
// Comment                  -
func GetChannelWebhooks(client httd.Getter, channelID snowflake.ID) (ret []*Webhook, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannelWebhooks(channelID),
		Endpoint:    EndpointChannels + "/" + channelID.String() + EndpointWebhooks,
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetGuildWebhooks [GET]   Returns a list of guild webhook objects. Requires the 'MANAGE_WEBHOOKS' permission.
// Endpoint                 /guilds/{guild.id}/webhooks
// Rate limiter             /guilds/{guild.id}/webhooks
// Discord documentation    https://discordapp.com/developers/docs/resources/webhook#get-guild-webhooks
// Reviewed                 2018-08-14
// Comment                  -
func GetGuildWebhooks(client httd.Getter, guildID snowflake.ID) (ret []*Webhook, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannelWebhooks(guildID),
		Endpoint:    EndpointChannels + "/" + guildID.String() + EndpointWebhooks,
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetWebhook [GET]         Returns the new webhook object for the given id.
// Endpoint                 /webhooks/{webhook.id}
// Rate limiter             /webhooks/{webhook.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/webhook#get-webhook
// Reviewed                 2018-08-14
// Comment                  -
func GetWebhook(client httd.Getter, webhookID snowflake.ID) (ret *Webhook, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitWebhook(webhookID),
		Endpoint:    EndpointWebhooks + "/" + webhookID.String(),
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetWebhookWithToken [GET]    Same as GetWebhook, except this call does not require authentication and returns
//                              no user in the webhook object.
// Endpoint                     /webhooks/{webhook.id}/{webhook.token}
// Rate limiter                 /webhooks/{webhook.id}
// Discord documentation        https://discordapp.com/developers/docs/resources/webhook#get-webhook-with-token
// Reviewed                     2018-08-14
// Comment                      -
func GetWebhookWithToken(client httd.Getter, webhookID snowflake.ID, token string) (ret *Webhook, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitWebhook(webhookID),
		Endpoint:    EndpointWebhooks + "/" + webhookID.String() + "/" + token,
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// ModifyWebhook [PATCH]    Modify a webhook. Requires the 'MANAGE_WEBHOOKS' permission.
//                          Returns the updated webhook object on success.
// Endpoint                 /webhooks/{webhook.id}
// Rate limiter             /webhooks/{webhook.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/webhook#modify-webhook
// Reviewed                 2018-08-14
// Comment                  All parameters to this endpoint are optional. Not tested:extra json fields might
//                          cause an issue. Consider writing a json params object.
func ModifyWebhook(client httd.Patcher, newWebhook *Webhook) (ret *Webhook, err error) {
	if newWebhook.ID.Empty() {
		err = errors.New("not a valid snowflake")
		return
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitWebhook(newWebhook.ID),
		Endpoint:    EndpointWebhooks + "/" + newWebhook.ID.String(),
	}
	_, body, err := client.Patch(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// ModifyWebhookWithToken [PATCH]   Same as ModifyWebhook, except this call does not require authentication,
//                                  does not accept a channel_id parameter in the body, and does not return
//                                  a user in the webhook object.
// Endpoint                         /webhooks/{webhook.id}/{webhook.token}
// Rate limiter                     /webhooks/{webhook.id}
// Discord documentation            https://discordapp.com/developers/docs/resources/webhook#modify-webhook-with-token
// Reviewed                         2018-08-14
// Comment                          All parameters to this endpoint are optional. Not tested:extra json fields might
//                                  cause an issue. Consider writing a json params object.
func ModifyWebhookWithToken(client httd.Patcher, newWebhook *Webhook) (ret *Webhook, err error) {
	if newWebhook.ID.Empty() {
		err = errors.New("not a valid snowflake")
		return
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitWebhook(newWebhook.ID),
		Endpoint:    EndpointWebhooks + "/" + newWebhook.ID.String() + "/" + newWebhook.Token,
	}
	_, body, err := client.Patch(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// DeleteWebhook [DELETE]   Delete a webhook permanently. User must be owner. Returns a 204 NO CONTENT response
//                          on success.
// Endpoint                 /webhooks/{webhook.id}
// Rate limiter             /webhooks/{webhook.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/webhook#delete-webhook
// Reviewed                 2018-08-14
// Comment                  -
func DeleteWebhook(client httd.Deleter, webhookID snowflake.ID) (err error) {
	return DeleteWebhookWithToken(client, webhookID, "")
}

// DeleteWebhookWithToken [DELETE]  Same as DeleteWebhook, except this call does not require authentication.
// Endpoint                         /webhooks/{webhook.id}/{webhook.token}
// Rate limiter                     /webhooks/{webhook.id}
// Discord documentation            https://discordapp.com/developers/docs/resources/webhook#delete-webhook-with-token
// Reviewed                         2018-08-14
// Comment                          -
func DeleteWebhookWithToken(client httd.Deleter, webhookID snowflake.ID, token string) (err error) {
	endpoint := EndpointWebhooks + "/" + webhookID.String()
	if token != "" {
		endpoint += "/" + token
	}
	details := &httd.Request{
		Ratelimiter: httd.RatelimitWebhook(webhookID),
		Endpoint:    endpoint,
	}
	resp, _, err := client.Delete(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(http.StatusNoContent)
		err = errors.New(msg)
	}
	return
}

func NewExecuteWebhookParams(id snowflake.ID, token string) (ret *ExecuteWebhookParams, err error) {
	return &ExecuteWebhookParams{
		WebhookID: id,
		Token:     token,
	}, nil
}

type ExecuteWebhookParams struct {
	WebhookID snowflake.ID `json:"-"`
	Token     string       `json:"-"`

	Content   string          `json:"content"`
	Username  string          `json:"username"`
	AvatarURL string          `json:"avatar_url"`
	TTS       bool            `json:"tts"`
	File      interface{}     `json:"file"`
	Embeds    []*ChannelEmbed `json:"embeds"`
}

// ExecuteWebhook [POST]    Trigger a webhook in Discord.
// Endpoint                 /webhooks/{webhook.id}/{webhook.token}
// Rate limiter             /webhooks/{webhook.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/webhook#execute-webhook
// Reviewed                 2018-08-14
// Comment                  This endpoint supports both JSON and form data bodies. It does require
//                          multipart/form-data requests instead of the normal JSON request type when
//                          uploading files. Make sure you set your Content-Type to multipart/form-data if
//                          you're doing that. Note that in that case, the embeds field cannot be used, but
//                          you can pass an url-encoded JSON body as a form value for payload_json.
// Comment#2                For the webhook embed objects, you can set every field except type (it will be
//                          rich regardless of if you try to set it), provider, video, and any height, width,
//                          or proxy_url values for images.
func ExecuteWebhook(client httd.Poster, params *ExecuteWebhookParams, wait bool, URLSuffix string) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitWebhook(params.WebhookID),
		Endpoint:    EndpointWebhooks + "/" + params.WebhookID.String() + "/" + params.Token + URLSuffix,
	}
	_, _, err = client.Post(details)
	if err != nil {
		return
	}

	//err = unmarshal(body, ret) // TODO: how to verify success?
	return
}

// ExecuteSlackWebhook [POST]   Trigger a webhook in Discord from the Slack app.
// Endpoint                     /webhooks/{webhook.id}/{webhook.token}
// Rate limiter                 /webhooks
// Discord documentation        https://discordapp.com/developers/docs/resources/webhook#execute-slackcompatible-webhook
// Reviewed                     2018-08-14
// Comment                      Refer to Slack's documentation for more information. We do not support Slack's channel,
//                              icon_emoji, mrkdwn, or mrkdwn_in properties.
func ExecuteSlackWebhook(client httd.Poster, params *ExecuteWebhookParams, wait bool) (err error) {
	return ExecuteWebhook(client, params, wait, EndpointSlackWebhook)
}

// ExecuteGitHubWebhook [POST]  Trigger a webhook in Discord from the GitHub app.
// Endpoint                     /webhooks/{webhook.id}/{webhook.token}
// Rate limiter                 /webhooks
// Discord documentation        https://discordapp.com/developers/docs/resources/webhook#execute-githubcompatible-webhook
// Reviewed                     2018-08-14
// Comment                      Add a new webhook to your GitHub repo (in the repo's settings), and use this endpoint
//                              as the "Payload URL." You can choose what events your Discord channel receives by
//                              choosing the "Let me select individual events" option and selecting individual
//                              events for the new webhook you're configuring.
func ExecuteGitHubWebhook(client httd.Poster, params *ExecuteWebhookParams, wait bool) (err error) {
	return ExecuteWebhook(client, params, wait, EndpointGitHubWebhook)
}
