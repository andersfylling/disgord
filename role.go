package disgord

import (
	"context"
	"fmt"
	"sort"
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
	ID          Snowflake `json:"id"`
	Name        string    `json:"name"`
	Color       uint      `json:"color"`
	Hoist       bool      `json:"hoist"`
	Position    int       `json:"position"` // can be -1
	Permissions uint64    `json:"permissions"`
	Managed     bool      `json:"managed"`
	Mentionable bool      `json:"mentionable"`

	guildID Snowflake
}

var _ Mentioner = (*Role)(nil)
var _ Reseter = (*Role)(nil)
var _ DeepCopier = (*Role)(nil)
var _ Copier = (*Role)(nil)
var _ discordDeleter = (*Role)(nil)
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
	r.CopyOverTo(copy)

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

func (r *Role) deleteFromDiscord(ctx context.Context, s Session, flags ...Flag) (err error) {
	guildID := r.guildID
	id := r.ID

	if id.IsZero() {
		return newErrorMissingSnowflake("role has no ID")
	}
	if guildID.IsZero() {
		return newErrorMissingSnowflake("role has no guildID")
	}

	err = s.Guild(guildID).DeleteRole(ctx, id, flags...)
	return err
}

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

// CreateGuildRoleParams ...
// https://discord.com/developers/docs/resources/guild#create-guild-role-json-params
type CreateGuildRoleParams struct {
	Name        string `json:"name,omitempty"`
	Permissions uint64 `json:"permissions,omitempty"`
	Color       uint   `json:"color,omitempty"`
	Hoist       bool   `json:"hoist,omitempty"`
	Mentionable bool   `json:"mentionable,omitempty"`

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string `json:"-"`
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
