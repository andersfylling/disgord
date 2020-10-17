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

// SortRoles sorts a slice of roles such that the first element is the top one in the Discord Guild Settings UI.
func SortRoles(rs []*Role) {
	sort.Sort(roles(rs))
}

// NewRole ...
func NewRole() *Role {
	return &Role{}
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

// DeepCopy see interface at struct.go#DeepCopier
func (r *Role) DeepCopy() (copy interface{}) {
	copy = NewRole()
	_ = r.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (r *Role) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var role *Role
	if role, ok = other.(*Role); !ok {
		return newErrorUnsupportedType("given interface{} was not a *Role")
	}

	role.ID = r.ID
	role.Name = r.Name
	role.Color = r.Color
	role.Hoist = r.Hoist
	role.Position = r.Position
	role.Permissions = r.Permissions
	role.Managed = r.Managed
	role.Mentionable = r.Mentionable
	role.guildID = r.guildID
	return
}

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

type GuildRoleQueryBuilder interface {
	WithContext(ctx context.Context) GuildRoleQueryBuilder

	Update(flags ...Flag) (builder UpdateGuildRoleBuilder)
	Delete(flags ...Flag) error
}

func (g guildQueryBuilder) Role(id Snowflake) GuildRoleQueryBuilder {
	return &guildRoleQueryBuilder{client: g.client, gid: g.gid, roleID: id}
}

type guildRoleQueryBuilder struct {
	ctx    context.Context
	client *Client
	gid    Snowflake
	roleID Snowflake
}

func (g guildRoleQueryBuilder) WithContext(ctx context.Context) GuildRoleQueryBuilder {
	g.ctx = ctx
	return &g
}

// UpdateRole Modify a guild role. Requires the 'MANAGE_ROLES' permission.
// Returns the updated role on success. Fires a Guild Role Update Gateway event.
func (g guildRoleQueryBuilder) Update(flags ...Flag) UpdateGuildRoleBuilder {
	builder := &updateGuildRoleBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Role{}
	}
	builder.r.flags = flags
	builder.r.IgnoreCache().setup(g.client.req, &httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildRole(g.gid, g.roleID),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// DeleteRole Delete a guild role. Requires the 'MANAGE_ROLES' permission.
// Returns a 204 empty response on success. Fires a Guild Role Delete Gateway event.
func (g guildRoleQueryBuilder) Delete(flags ...Flag) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.GuildRole(g.gid, g.roleID),
		Ctx:      g.ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err := r.Execute()
	return err
}

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

// updateGuildRoleBuilder ...
//generate-rest-basic-execute: role:*Role,
//generate-rest-params: name:string, permissions:PermissionBit, color:uint, hoist:bool, mentionable:bool,
type updateGuildRoleBuilder struct {
	r RESTBuilder
}
