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
	Type         OptionType                        `json:"type"`
	Name         string                            `json:"name"`
	Description  string                            `json:"description"`
	Required     bool                              `json:"required"`
	Choices      []*ApplicationCommandOptionChoice `json:"choices"`
	Options      []*ApplicationCommandOption       `json:"options"`
	ChannelTypes []ChannelType                     `json:"channel_types"`
	MinValue     float64                           `json:"min_value"`
	MaxValue     float64                           `json:"max_value"`
	Autocomplete bool                              `json:"autocomplete"`
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
	ID                       Snowflake                   `json:"id"`
	Type                     ApplicationCommandType      `json:"type,omitempty"`
	ApplicationID            Snowflake                   `json:"application_id"`
	GuildID                  Snowflake                   `json:"guild_id,omitempty"`
	Name                     string                      `json:"name"`
	NameLocalizations        map[string]string           `json:"name_localization,omitempty"`
	Description              string                      `json:"description"`
	DescriptionLocalizations map[string]string           `json:"description_localization,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *PermissionBit              `json:"default_member_permissions,omitempty"`
	DMPermission             *bool                       `json:"dm_permission,omitempty"` // Global only.
	DefaultPermission        bool                        `json:"default_permission,omitempty"`
	NSFW                     bool                        `json:"nsfw,omitempty"`
	Version                  Snowflake                   `json:"version,omitempty"`
}

type CreateApplicationCommand struct {
	Name                     string                      `json:"name"`
	NameLocalizations        map[string]string           `json:"name_localization,omitempty"`
	Description              string                      `json:"description,omitempty"`
	DescriptionLocalizations map[string]string           `json:"description_localization,omitempty"`
	Type                     ApplicationCommandType      `json:"type,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DMPermission             *bool                       `json:"dm_permission,omitempty"`
	DefaultMemberPermissions *PermissionBit              `json:"default_member_permissions,omitempty"`
	DefaultPermission        bool                        `json:"default_permission,omitempty"`
	NSFW                     *bool                       `json:"nsfw,omitempty"`
}

type UpdateApplicationCommand struct {
	Name                     *string                      `json:"name,omitempty"`
	NameLocalizations        *map[string]string           `json:"name_localization,omitempty"`
	Description              *string                      `json:"description,omitempty"`
	DescriptionLocalizations *map[string]string           `json:"description_localization,omitempty"`
	Type                     *ApplicationCommandType      `json:"type,omitempty"`
	Options                  *[]*ApplicationCommandOption `json:"options,omitempty"`
	DMPermission             *bool                        `json:"dm_permission,omitempty"` // Global only.
	DefaultMemberPermissions *PermissionBit               `json:"default_member_permissions,omitempty"`
	DefaultPermission        *bool                        `json:"default_permission,omitempty"`
	NSFW                     *bool                        `json:"nsfw,omitempty"`
}

type ApplicationCommandQueryBuilder interface {
	WithContext(ctx context.Context) ApplicationCommandQueryBuilder
	Global() ApplicationCommandFunctions
	Guild(guildID Snowflake) ApplicationCommandFunctions
}

type ApplicationCommandFunctions interface {
	Get(commandId Snowflake) (*ApplicationCommand, error)
	Commands() ([]*ApplicationCommand, error)
	Delete(commandId Snowflake) error
	Create(command *CreateApplicationCommand) error
	Update(commandId Snowflake, command *UpdateApplicationCommand) error
	BulkOverwrite(commands []*CreateApplicationCommand) error
}

type applicationCommandFunctions struct {
	appID   Snowflake
	flags   Flag
	client  *Client
	guildID Snowflake
	ctx     context.Context
}

func (c *applicationCommandFunctions) applicationID() Snowflake {
	appID := c.appID
	if appID.IsZero() {
		c.client.mu.Lock()
		appID = c.client.applicationID
		c.client.mu.Unlock()
	}

	return appID
}

func applicationCommandFactory() interface{} {
	return &ApplicationCommand{}
}

func (c *applicationCommandFunctions) Get(commandID Snowflake) (*ApplicationCommand, error) {
	var endpoint string
	if c.guildID == 0 {
		endpoint = fmt.Sprintf("/applications/%d/commands/%d", c.applicationID(), commandID)
	} else {
		endpoint = fmt.Sprintf("/applications/%d/guilds/%d/commands/%d", c.applicationID(), c.guildID, commandID)
	}

	req := &httd.Request{
		Endpoint:    endpoint,
		Method:      http.MethodGet,
		Ctx:         c.ctx,
		ContentType: httd.ContentTypeJSON,
	}

	r := c.client.newRESTRequest(req, c.flags)
	r.factory = applicationCommandFactory

	return getApplicationCommand(r.Execute)
}

func (c *applicationCommandFunctions) Commands() ([]*ApplicationCommand, error) {
	var endpoint string
	if c.guildID == 0 {
		endpoint = fmt.Sprintf("/applications/%d/commands", c.applicationID())
	} else {
		endpoint = fmt.Sprintf("/applications/%d/guilds/%d/commands", c.applicationID(), c.guildID)
	}

	req := &httd.Request{
		Endpoint:    endpoint,
		Method:      http.MethodGet,
		Ctx:         c.ctx,
		ContentType: httd.ContentTypeJSON,
	}

	r := c.client.newRESTRequest(req, c.flags)
	r.factory = func() interface{} {
		tmp := make([]*ApplicationCommand, 0)
		return &tmp
	}

	return getApplicationCommands(r.Execute)
}

func (c *applicationCommandFunctions) Create(command *CreateApplicationCommand) error {
	var endpoint string
	if c.guildID == 0 {
		endpoint = fmt.Sprintf("/applications/%d/commands", c.applicationID())
	} else {
		endpoint = fmt.Sprintf("/applications/%d/guilds/%d/commands", c.applicationID(), c.guildID)
	}

	req := &httd.Request{
		Endpoint:    endpoint,
		Method:      http.MethodPost,
		Body:        command,
		Ctx:         c.ctx,
		ContentType: httd.ContentTypeJSON,
	}

	r := c.client.newRESTRequest(req, c.flags)
	r.factory = applicationCommandFactory
	_, err := r.Execute()
	return err
}

func (c *applicationCommandFunctions) Update(commandID Snowflake, command *UpdateApplicationCommand) error {
	var endpoint string

	if c.guildID == 0 {
		endpoint = fmt.Sprintf("/applications/%d/commands/%d", c.applicationID(), commandID)
	} else {
		endpoint = fmt.Sprintf("/applications/%d/guilds/%d/commands/%d", c.applicationID(), c.guildID, commandID)
	}

	req := &httd.Request{
		Endpoint:    endpoint,
		Method:      http.MethodPatch,
		Body:        command,
		Ctx:         c.ctx,
		ContentType: httd.ContentTypeJSON,
	}

	r := c.client.newRESTRequest(req, c.flags)
	r.factory = applicationCommandFactory
	_, err := r.Execute()
	return err
}

func (c *applicationCommandFunctions) Delete(commandID Snowflake) error {
	var endpoint string

	if c.guildID == 0 {
		endpoint = fmt.Sprintf("/applications/%d/commands/%d", c.applicationID(), commandID)
	} else {
		endpoint = fmt.Sprintf("/applications/%d/guilds/%d/commands/%d", c.applicationID(), c.guildID, commandID)
	}

	req := &httd.Request{
		Endpoint:    endpoint,
		Method:      http.MethodDelete,
		Ctx:         c.ctx,
		ContentType: httd.ContentTypeJSON,
	}

	r := c.client.newRESTRequest(req, c.flags)
	r.factory = applicationCommandFactory
	_, err := r.Execute()
	return err
}

func (c *applicationCommandFunctions) BulkOverwrite(commands []*CreateApplicationCommand) error {
	var endpoint string
	if c.guildID == 0 {
		endpoint = fmt.Sprintf("/applications/%d/commands", c.applicationID())
	} else {
		endpoint = fmt.Sprintf("/applications/%d/guilds/%d/commands", c.applicationID(), c.guildID)
	}

	req := &httd.Request{
		Endpoint:    endpoint,
		Method:      http.MethodPut,
		Body:        commands,
		Ctx:         c.ctx,
		ContentType: httd.ContentTypeJSON,
	}

	r := c.client.newRESTRequest(req, c.flags)
	r.factory = func() interface{} {
		tmp := make([]*ApplicationCommand, 0)
		return &tmp
	}
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
