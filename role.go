package disgord

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
)

type roles []*Role

var _ sort.Interface = (roles)(nil)

func (r roles) Len() int {
	return len(r)
}

// Less is reversed due to the visual ordering in Discord.
func (r roles) Less(i, j int) bool {
	a := r[i]
	b := r[j]

	if a.Position == b.Position {
		return a.ID < b.ID
	}

	return a.Position > b.Position
}

func (r roles) Swap(i, j int) {
	tmp := r[i]
	r[i] = r[j]
	r[j] = tmp
}

// Role https://discord.com/developers/docs/topics/permissions#role-object
type Role struct {
	ID          Snowflake     `json:"id"`
	Name        string        `json:"name"`
	Color       uint          `json:"color"`
	Hoist       bool          `json:"hoist"`
	Position    int           `json:"position"` // can be -1
	Permissions PermissionBit `json:"permissions"`
	Managed     bool          `json:"managed"`
	Mentionable bool          `json:"mentionable"`
	guildID     Snowflake
}

var _ Mentioner = (*Role)(nil)
var _ Reseter = (*Role)(nil)
var _ DeepCopier = (*Role)(nil)
var _ Copier = (*Role)(nil)
var _ fmt.Stringer = (*Role)(nil)

func (r *Role) String() string {
	return r.Name
}

// Mention gives a formatted version of the role such that it can be parsed by Discord clients
func (r *Role) Mention() string {
	return "<@&" + r.ID.String() + ">"
}

// SetGuildID link role to a guild before running session.SaveToDiscord(*Role)
func (r *Role) SetGuildID(id Snowflake) {
	r.guildID = id
}

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

type GuildRoleQueryBuilder interface {
	WithContext(ctx context.Context) GuildRoleQueryBuilder
	WithFlags(flags ...Flag) GuildRoleQueryBuilder
	Update(params *UpdateRole) (*Role, error)
	Delete() error
	Get(ID Snowflake) (*Role, error)
	// Deprecated: use Update
	UpdateBuilder() (builder UpdateGuildRoleBuilder)
}

func (g guildQueryBuilder) Role(id Snowflake) GuildRoleQueryBuilder {
	return &guildRoleQueryBuilder{client: g.client, gid: g.gid, roleID: id}
}

type guildRoleQueryBuilder struct {
	ctx    context.Context
	flags  Flag
	client *Client
	gid    Snowflake
	roleID Snowflake
}

func (g *guildRoleQueryBuilder) validate() error {
	if g.client == nil {
		return ErrMissingClientInstance
	}
	if g.gid.IsZero() {
		return ErrMissingGuildID
	}
	if g.roleID.IsZero() {
		return ErrMissingRoleID
	}
	return nil
}

func (g guildRoleQueryBuilder) WithContext(ctx context.Context) GuildRoleQueryBuilder {
	g.ctx = ctx
	return &g
}

func (g guildRoleQueryBuilder) WithFlags(flags ...Flag) GuildRoleQueryBuilder {
	g.flags = mergeFlags(flags)
	return &g
}

// Delete Deletes a guild role. Requires the 'MANAGE_ROLES' permission.
// Returns a 204 empty response on success. Fires a Guild Role Delete Gateway event.
func (g guildRoleQueryBuilder) Delete() error {
	if err := g.validate(); err != nil {
		return err
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:   http.MethodDelete,
		Endpoint: endpoint.GuildRole(g.gid, g.roleID),
		Ctx:      g.ctx,
	}, g.flags)

	_, err := r.Execute()
	return err
}

// Update update a role
func (g guildRoleQueryBuilder) Update(params *UpdateRole) (*Role, error) {
	if params == nil {
		return nil, MissingRESTParamsErr
	}
	if err := g.validate(); err != nil {
		return nil, err
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildRole(g.gid, g.roleID),
		ContentType: httd.ContentTypeJSON,
		Body:        params,
		Reason:      params.AuditLogReason,
	}, g.flags)
	r.factory = func() interface{} {
		return &Role{}
	}

	return getRole(r.Execute)
}

type UpdateRole struct {
	Name        *string        `json:"name,omitempty"`
	Permissions *PermissionBit `json:"permissions,omitempty"`
	Color       *int           `json:"color,omitempty"`
	Hoist       *bool          `json:"hoist,omitempty"`
	Mentionable *bool          `json:"mentionable,omitempty"`

	AuditLogReason string `json:"-"`
}
