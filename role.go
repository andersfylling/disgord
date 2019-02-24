package disgord

import (
	"encoding/json"
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

//func (r *Role) Clear() {
//}

// GetGuildRoles [REST] Returns a list of role objects for the guild.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/roles
//  Rate limiter            /guilds/{guild.id}/roles
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-roles
//  Reviewed                2018-08-18
//  Comment                 -
func GetGuildRoles(client httd.Getter, guildID Snowflake) (ret []*Role, err error) {
	var body []byte
	_, body, err = client.Get(&httd.Request{
		Ratelimiter: ratelimitGuildRoles(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/roles",
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	if err != nil {
		return
	}

	// add guild id to roles
	for _, role := range ret {
		role.guildID = guildID
	}
	return
}

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
func CreateGuildRole(client httd.Poster, id Snowflake, params *CreateGuildRoleParams) (ret *Role, err error) {
	var body []byte
	_, body, err = client.Post(&httd.Request{
		Ratelimiter: ratelimitGuildRoles(id),
		Endpoint:    endpoint.GuildRoles(id),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return nil, err
	}

	if err = unmarshal(body, &ret); err != nil {
		return nil, err
	}

	// add guild id to roles
	ret.guildID = id
	return ret, nil
}

// UpdateGuildRolePositionsParams ...
// https://discordapp.com/developers/docs/resources/guild#modify-guild-role-positions-json-params
type UpdateGuildRolePositionsParams struct {
	ID       Snowflake `json:"id"`
	Position uint      `json:"position"`
}

// UpdateGuildRolePositions [REST] Modify the positions of a set of role objects for the guild.
// Requires the 'MANAGE_ROLES' permission. Returns a list of all of the guild's role objects on success.
// Fires multiple Guild Role Update Gateway events.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/roles
//  Rate limiter            /guilds/{guild.id}/roles
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#modify-guild-role-positions
//  Reviewed                2018-08-18
//  Comment                 -
func UpdateGuildRolePositions(client httd.Patcher, guildID Snowflake, params []UpdateGuildRolePositionsParams) (ret []*Role, err error) {
	var body []byte
	_, body, err = client.Patch(&httd.Request{
		Ratelimiter: ratelimitGuildRoles(guildID),
		Endpoint:    endpoint.GuildRoles(guildID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return nil, err
	}

	if err = unmarshal(body, &ret); err != nil {
		return nil, err
	}

	// add guild id to roles
	for _, role := range ret {
		role.guildID = guildID
	}
	return ret, nil
}

// UpdateGuildRoleParams JSON params for func ModifyGuildRole
type UpdateGuildRoleParams struct {
	data map[string]interface{}
}

func (p *UpdateGuildRoleParams) init() {
	if p.data != nil {
		return
	}

	p.data = map[string]interface{}{}
}

func (p *UpdateGuildRoleParams) SetName(name string) {
	p.init()
	p.data["name"] = name
}

func (p *UpdateGuildRoleParams) SetPermissions(permissions uint64) {
	p.init()
	p.data["permissions"] = permissions
}

func (p *UpdateGuildRoleParams) SetColor(color uint) {
	p.init()
	p.data["color"] = color
}

func (p *UpdateGuildRoleParams) SetHoist(hoist bool) {
	p.init()
	p.data["hoist"] = hoist
}

func (p *UpdateGuildRoleParams) SetMentionable(mentionable bool) {
	p.init()
	p.data["mentionable"] = mentionable
}

func (p *UpdateGuildRoleParams) MarshalJSON() ([]byte, error) {
	if len(p.data) == 0 {
		return []byte(`{}`), nil
	}

	return httd.Marshal(p.data)
}

var _ json.Marshaler = (*UpdateGuildRoleParams)(nil)

// ModifyGuildRole [REST] Modify a guild role. Requires the 'MANAGE_ROLES' permission.
// Returns the updated role on success. Fires a Guild Role Update Gateway event.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/roles/{role.id}
//  Rate limiter            /guilds/{guild.id}/roles
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#modify-guild-role
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) UpdateGuildRole(guildID, roleID Snowflake, flags ...Flag) (builder *modifyGuildRoleBuilder) {
	builder = &modifyGuildRoleBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Role{}
	}
	builder.r.IgnoreCache().setup(c.cache, c.req, &httd.Request{
		Method:      http.MethodGet,
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

// modifyGuildRoleBuilder ...
//generate-rest-params: name:string, permissions:uint64, color:uint, hoist:bool, mentionable:bool,
type modifyGuildRoleBuilder struct {
	r RESTBuilder
}

func (b *modifyGuildRoleBuilder) Execute() (role *Role, err error) {
	var v interface{}
	if v, err = b.r.execute(); err != nil {
		return
	}

	role = v.(*Role)
	return
}

// DeleteGuildRole [REST] Delete a guild role. Requires the 'MANAGE_ROLES' permission.
// Returns a 204 empty response on success. Fires a Guild Role Delete Gateway event.
//  Method                  DELETE
//  Endpoint                /guilds/{guild.id}/roles/{role.id}
//  Rate limiter            /guilds/{guild.id}/roles
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#delete-guild-role
//  Reviewed                2018-08-18
//  Comment                 -
func DeleteGuildRole(client httd.Deleter, guildID, roleID Snowflake) (err error) {
	var resp *http.Response
	resp, _, err = client.Delete(&httd.Request{
		Ratelimiter: ratelimitGuildRoles(guildID),
		Endpoint:    endpoint.GuildRole(guildID, roleID),
	})
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not remove role. Do you have the MANAGE_ROLES permission?"
		err = errors.New(msg)
	}

	return
}
