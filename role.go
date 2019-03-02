package disgord

import (
	"errors"
	"net/http"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
)

// NewRole ...
func NewRole() *Role {
	return &Role{}
}

// Role https://discordapp.com/developers/docs/topics/permissions#role-object
type Role struct {
	Lockable `json:"-"`

	ID          Snowflake `json:"id"`
	Name        string    `json:"name"`
	Color       uint      `json:"color"`
	Hoist       bool      `json:"hoist"`
	Position    uint      `json:"position"`
	Permissions uint64    `json:"permissions"`
	Managed     bool      `json:"managed"`
	Mentionable bool      `json:"mentionable"`

	guildID Snowflake
}

var _ Reseter = (*Role)(nil)

// Mention gives a formatted version of the role such that it can be parsed by Discord clients
func (r *Role) Mention() string {
	return "<@&" + r.ID.String() + ">"
}

// SetGuildID link role to a guild before running session.SaveToDiscord(*Role)
func (r *Role) SetGuildID(id Snowflake) {
	r.ID = id
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

	if constant.LockedMethods {
		r.RLock()
		role.Lock()
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

	if constant.LockedMethods {
		r.RUnlock()
		role.Unlock()
	}

	return
}

func (r *Role) saveToDiscord(session Session, changes discordSaver) (err error) {
	if r.guildID.Empty() {
		err = newErrorMissingSnowflake("role has no guildID")
		return
	}

	var role *Role
	if r.ID.Empty() {
		// create role
		params := CreateGuildRoleParams{
			Name:        r.Name,
			Permissions: r.Permissions,
			Color:       r.Color,
			Hoist:       r.Hoist,
			Mentionable: r.Mentionable,
		}
		role, err = session.CreateGuildRole(r.guildID, &params)
		if err != nil {
			return
		}
		err = role.CopyOverTo(r)
		return err
	}

	return errors.New("updating discord objects are not yet implemented - only saving new ones")
	//
	//// modify/update role
	//params := UpdateGuildRoleParams{}
	//params.SetName(r.Name)
	//params.SetPermissions(r.Permissions)
	//params.SetColor(r.Color)
	//params.SetHoist(r.Hoist)
	//params.SetMentionable(r.Mentionable)
	//role, err = session.ModifyGuildRole(r.guildID, r.ID, &params)
	//if err != nil {
	//	return
	//}
	//if role.Position != r.Position {
	//	// update the position
	//	params := []UpdateGuildRolePositionsParams{{
	//		ID:       r.ID,
	//		Position: r.Position,
	//	}}
	//	_, err = session.UpdateGuildRolePositions(r.guildID, params)
	//	if err != nil {
	//		return
	//	}
	//	role.Position = r.Position
	//}
	//
	//return
}

func (r *Role) deleteFromDiscord(session Session) (err error) {
	if r.ID.Empty() {
		err = newErrorMissingSnowflake("role has no ID")
		return
	}
	if r.guildID.Empty() {
		err = newErrorMissingSnowflake("role has no guildID")
		return
	}

	err = session.DeleteGuildRole(r.guildID, r.ID)
	return
}

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

// CreateGuildRoleParams ...
// https://discordapp.com/developers/docs/resources/guild#create-guild-role-json-params
type CreateGuildRoleParams struct {
	Name        string `json:"name,omitempty"`
	Permissions uint64 `json:"permissions,omitempty"`
	Color       uint   `json:"color,omitempty"`
	Hoist       bool   `json:"hoist,omitempty"`
	Mentionable bool   `json:"mentionable,omitempty"`
}

// CreateGuildRole [REST] Create a new role for the guild. Requires the 'MANAGE_ROLES' permission.
// Returns the new role object on success. Fires a Guild Role Create Gateway event.
//  Method                  POST
//  Endpoint                /guilds/{guild.id}/roles
//  Rate limiter            /guilds/{guild.id}/roles
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#create-guild-role
//  Reviewed                2018-08-18
//  Comment                 All JSON params are optional.
func (c *client) CreateGuildRole(id Snowflake, params *CreateGuildRoleParams, flags ...Flag) (ret *Role, err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPost,
		Ratelimiter: ratelimitGuildRoles(id),
		Endpoint:    endpoint.GuildRoles(id),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.CacheRegistry = GuildRoleCache
	r.factory = func() interface{} {
		return &Role{}
	}
	r.preUpdateCache = func(x interface{}) {
		r := x.(*Role)
		r.guildID = id
	}

	return getRole(r.Execute)
}

// ModifyGuildRole [REST] Modify a guild role. Requires the 'MANAGE_ROLES' permission.
// Returns the updated role on success. Fires a Guild Role Update Gateway event.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/roles/{role.id}
//  Rate limiter            /guilds/{guild.id}/roles
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#modify-guild-role
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) UpdateGuildRole(guildID, roleID Snowflake, flags ...Flag) (builder *updateGuildRoleBuilder) {
	builder = &updateGuildRoleBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Role{}
	}
	builder.r.flags = flags
	builder.r.IgnoreCache().setup(c.cache, c.req, &httd.Request{
		Method:      http.MethodPatch,
		Ratelimiter: ratelimitGuildRoles(guildID),
		Endpoint:    endpoint.GuildRole(guildID, roleID),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	builder.r.cacheMiddleware = func(resp *http.Response, v interface{}, err error) error {
		role := v.(*Role)
		role.guildID = guildID
		return nil
	}

	return builder
}

// DeleteGuildRole [REST] Delete a guild role. Requires the 'MANAGE_ROLES' permission.
// Returns a 204 empty response on success. Fires a Guild Role Delete Gateway event.
//  Method                  DELETE
//  Endpoint                /guilds/{guild.id}/roles/{role.id}
//  Rate limiter            /guilds/{guild.id}/roles
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#delete-guild-role
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) DeleteGuildRole(guildID, roleID Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodDelete,
		Ratelimiter: ratelimitGuildRoles(guildID),
		Endpoint:    endpoint.GuildRole(guildID, roleID),
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// GetGuildRoles [REST] Returns a list of role objects for the guild.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/roles
//  Rate limiter            /guilds/{guild.id}/roles
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-roles
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) GetGuildRoles(guildID Snowflake, flags ...Flag) (ret []*Role, err error) {
	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitGuildRoles(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/roles",
	}, flags)
	r.CacheRegistry = GuildRolesCache
	r.factory = func() interface{} {
		tmp := make([]*Role, 0)
		return &tmp
	}
	r.preUpdateCache = func(x interface{}) {
		roles := *x.(*[]*Role)
		for i := range roles {
			roles[i].guildID = guildID
		}
	}

	return getRoles(r.Execute)
}

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

// updateGuildRoleBuilder ...
//generate-rest-basic-execute: role:*Role,
//generate-rest-params: name:string, permissions:uint64, color:uint, hoist:bool, mentionable:bool,
type updateGuildRoleBuilder struct {
	r RESTBuilder
}
