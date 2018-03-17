package resource

import (
	"github.com/andersfylling/snowflake"
)

// Webhook Used to represent a webhook
// https://discordapp.com/developers/docs/resources/webhook#webhook-object
type Webhook struct {
	ID        snowflake.ID `json:"id"`                 //  |
	GuildID   snowflake.ID `json:"guild_id,omitempty"` //  |?
	ChannelID snowflake.ID `json:"channel_id"`         //  |
	User      *User        `json:"user,omitempty"`     // ?|
	Name      string       `json:"name"`               //  |?
	Avatar    string       `json:"avatar"`             //  |?
	Token     string       `json:"token"`              //  |
}

func ReqCreateWebhook()                  {}
func ReqGetChannelWebhooks()             {}
func ReqGetGuildWebhooks()               {}
func ReqGetWebhook()                     {}
func ReqGetWebhookWithToken()            {}
func ReqModifyWebhook()                  {}
func ReqModifyWebhookWithToken()         {}
func ReqDeleteWebhook()                  {}
func ReqDeleteWebhookWithToken()         {}
func ReqExecuteWebhook()                 {}
func ReqExecuteSlackCompatibleWebhook()  {}
func ReqExecuteGitHubCompatibleWebgook() {}
