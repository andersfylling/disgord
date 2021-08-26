package disgord

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andersfylling/disgord/internal/httd"
)

type ApplicationCommandType int

const (
	_ ApplicationCommandType = iota
	ApplicationCommandChatInput
	ApplicationCommandUser
	ApplicationCommandMessage
)

type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type ApplicationCommandOption struct {
	Type        OptionType                        `json:"type"`
	Name        string                            `json:"name"`
	Description string                            `json:"description"`
	Required    bool                              `json:"required"`
	Choices     []*ApplicationCommandOptionChoice `json:"choices"`
	Options     []*ApplicationCommandOption       `json:"options"`
}

type ApplicationCommandDataOption struct {
	Name    string                          `json:"name"`
	Type    OptionType                      `json:"type"`
	Value   interface{}                     `json:"value"`
	Options []*ApplicationCommandDataOption `json:"options"`
}

type ApplicationCommandPermissionType int

const (
	_ ApplicationCommandPermissionType = iota
	ApplicationCommandPermissionRole
	ApplicationCommandPermissionUser
)

type ApplicationCommandPermissions struct {
	ID         Snowflake                        `json:"id"`
	Type       ApplicationCommandPermissionType `json:"type"`
	Permission bool                             `json:"permission"`
}
type GuildApplicationCommandPermissions struct {
	ID            Snowflake                        `json:"id"`
	ApplicationID Snowflake                        `json:"application_id"`
	GuildID       Snowflake                        `json:"guild_id"`
	Permissions   []*ApplicationCommandPermissions `json:"permissions"`
}

type ApplicationCommand struct {
	ID                Snowflake                   `json:"id"`
	Type              ApplicationCommandType      `json:"type"`
	ApplicationID     Snowflake                   `json:"application_id"`
	GuildID           Snowflake                   `json:"guild_id"`
	Name              string                      `json:"name"`
	Description       string                      `json:"description"`
	Options           []*ApplicationCommandOption `json:"options"`
	DefaultPermission bool                        `json:"default_permission,omitempty"`
}

type UpdateApplicationCommand struct {
	Name              string                      `json:"name,omitempty"`
	DefaultPermission bool                        `json:"default_permission,omitempty"`
	Description       string                      `json:"description,omitempty"`
	Options           []*ApplicationCommandOption `json:"options,omitempty"`
}

type ApplicationCommandQueryBuilder interface {
	WithContext(ctx context.Context) ApplicationCommandQueryBuilder
	Global() ApplicationCommandFunctions
	Guild(guildID Snowflake) ApplicationCommandFunctions
}

type ApplicationCommandFunctions interface {
	Delete(commandId Snowflake) error
	Create(command *ApplicationCommand) error
	Update(commandId Snowflake, command *UpdateApplicationCommand) error
}

type applicationCommandFunctions struct {
	appID   Snowflake
	flags   Flag
	client  *Client
	guildID Snowflake
	ctx     context.Context
}

func applicationCommandFactory() interface{} {
	return &ApplicationCommand{}
}

func (c *applicationCommandFunctions) Create(command *ApplicationCommand) error {
	var endpoint string
	if c.guildID == 0 {
		endpoint = fmt.Sprintf("/applications/%d/commands", c.appID)
	} else {
		endpoint = fmt.Sprintf("/applications/%d/guilds/%d/commands", c.appID, c.guildID)
	}
	ctx := c.ctx
	req := &httd.Request{
		Endpoint:    endpoint,
		Method:      http.MethodPost,
		Body:        command,
		Ctx:         ctx,
		ContentType: httd.ContentTypeJSON,
	}
	r := c.client.newRESTRequest(req, c.flags)
	r.factory = applicationCommandFactory
	_, err := r.Execute()
	return err
}

func (c *applicationCommandFunctions) Update(commandID Snowflake, command *UpdateApplicationCommand) error {
	var endpoint string
	ctx := c.ctx

	if c.guildID == 0 {
		endpoint = fmt.Sprintf("applications/%d/commands/%d", c.appID, commandID)
	} else {
		endpoint = fmt.Sprintf("applications/%d/guilds/%d/commands/%d", c.appID, c.guildID, commandID)
	}
	req := &httd.Request{
		Endpoint:    endpoint,
		Method:      http.MethodPatch,
		Body:        command,
		Ctx:         ctx,
		ContentType: httd.ContentTypeJSON,
	}
	r := c.client.newRESTRequest(req, c.flags)
	r.factory = applicationCommandFactory
	_, err := r.Execute()
	return err
}

func (c *applicationCommandFunctions) Delete(commandID Snowflake) error {
	var endpoint string
	ctx := c.ctx

	if c.guildID == 0 {
		endpoint = fmt.Sprintf("applications/%d/commands/%d", c.appID, commandID)
	} else {
		endpoint = fmt.Sprintf("applications/%d/guilds/%d/commands/%d", c.appID, c.guildID, commandID)
	}
	req := &httd.Request{
		Endpoint:    endpoint,
		Method:      http.MethodDelete,
		Ctx:         ctx,
		ContentType: httd.ContentTypeJSON,
	}
	r := c.client.newRESTRequest(req, c.flags)
	r.factory = applicationCommandFactory
	_, err := r.Execute()
	return err
}

type applicationCommandQueryBuilder struct {
	ctx    context.Context
	client *Client
	appID  Snowflake
	flags  Flag
}

func (q applicationCommandQueryBuilder) Guild(guildId Snowflake) ApplicationCommandFunctions {
	return &applicationCommandFunctions{appID: q.appID, ctx: q.ctx, guildID: guildId, client: q.client, flags: q.flags}
}
func (q applicationCommandQueryBuilder) Global() ApplicationCommandFunctions {
	return &applicationCommandFunctions{appID: q.appID, ctx: q.ctx, client: q.client, flags: q.flags}
}

func (q applicationCommandQueryBuilder) WithFlags(flags ...Flag) ApplicationCommandQueryBuilder {
	q.flags = mergeFlags(flags)
	return &q
}

func (q applicationCommandQueryBuilder) WithContext(ctx context.Context) ApplicationCommandQueryBuilder {
	q.ctx = ctx
	return &q
}

func (c clientQueryBuilder) ApplicationCommand(id Snowflake) ApplicationCommandQueryBuilder {
	return &applicationCommandQueryBuilder{client: c.client, appID: id}
}
