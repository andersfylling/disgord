package disgord

import (
	"net/http"
	"sort"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
)

type roles []*Role

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

var _ discordSaver = (*roles)(nil)
var _ sort.Interface = (roles)(nil)

func (rp *roles) saveToDiscord(s Session, flags ...Flag) error {
	r := *rp
	var guildID Snowflake
	for i := range r {
		if !r[i].guildID.Empty() {
			guildID = r[i].guildID
			break
		}
	}
	updated, err := s.UpdateGuildRolePositions(guildID, NewUpdateGuildRolePositionsParams(r), flags...)
	if err != nil {
		return err
	}

	// Since the updating guild role positions _requires_ you to send all the roles
	// you should be given the exact same roles in return. We sort them such that we only need to iterate
	// with a O(N) instead of a O(N*M). However, since I don't trust Discord (...) I keep the option open
	// that more than the local number of roles might be returned.
	SortRoles(r)
	SortRoles(updated)
	var newRoles []*Role
	for j := range updated {
		var handled bool
		for i := range r {
			if r[i].ID != updated[j].ID {
				continue
			}

			_ = updated[j].CopyOverTo(r[i])
			updated[j] = nil
			updated[j] = updated[len(updated)-1]
			updated = updated[:len(updated)-1]
			handled = true
			break
		}

		if !handled {
			newRoles = append(newRoles, updated[j])
		}
	}
	*rp = append(r, newRoles...)
	SortRoles(*rp)

	return err
}

// SortRoles sorts a slice of roles such that the first element is the top one in the Discord Guild Settings UI.
func SortRoles(rs []*Role) {
	sort.Sort(roles(rs))
}

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

func (r *Role) saveToDiscord(s Session, flags ...Flag) (err error) {
	if constant.LockedMethods {
		r.RLock()
	}
	guildID := r.guildID
	id := r.ID
	//pos := r.Position
	if constant.LockedMethods {
		r.RUnlock()
	}

	if guildID.Empty() {
		err = newErrorMissingSnowflake("role has no guildID. Use Role.SetGuildID(..)")
		return
	}

	var role *Role
	if id.Empty() {
		// create role
		if constant.LockedMethods {
			r.RLock()
		}
		params := CreateGuildRoleParams{
			Name:        r.Name,
			Permissions: r.Permissions,
			Color:       r.Color,
			Hoist:       r.Hoist,
			Mentionable: r.Mentionable,
		}
		if constant.LockedMethods {
			r.RUnlock()
		}
		role, err = s.CreateGuildRole(guildID, &params, flags...)
	} else {
		if constant.LockedMethods {
			r.RLock()
		}
		builder := s.UpdateGuildRole(guildID, id, flags...).
			SetName(r.Name).
			SetColor(r.Color).
			SetHoist(r.Hoist).
			SetMentionable(r.Mentionable).
			SetPermissions(r.Permissions)
		if constant.LockedMethods {
			r.RUnlock()
		}
		role, err = builder.Execute()
		if err == nil {
			// TODO: handle role position
			//  this is a little tricky as a user might not want to change the positions

			//if role != nil && role.Position != pos {
			//
			////var roles []*Role
			////roles, err = s.GetGuildRoles(guildID, flags...)
			////if err != nil {
			////	err = errors.New("unable to update role position: " + err.Error())
			////} else {
			////	params := NewUpdateGuildRolePositionsParams(roles)
			////	_, err = s.UpdateGuildRolePositions(guildID, params, flags...)
			////	if err != nil {
			////		err = errors.New("unable to update role position: " + err.Error())
			////	} else {
			////		err = role.CopyOverTo(r)
			////	}
			////}
			//}

		}
	}

	if err != nil {
		return err
	}

	err = role.CopyOverTo(r)
	return err
}

func (r *Role) deleteFromDiscord(s Session, flags ...Flag) (err error) {
	if constant.LockedMethods {
		r.RLock()
	}
	guildID := r.guildID
	id := r.ID
	if constant.LockedMethods {
		r.RUnlock()
	}

	if id.Empty() {
		return newErrorMissingSnowflake("role has no ID")
	}
	if guildID.Empty() {
		return newErrorMissingSnowflake("role has no guildID")
	}

	err = s.DeleteGuildRole(guildID, id, flags...)
	return err
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
