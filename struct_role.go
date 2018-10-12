package disgord

import "sync"

// NewRole ...
func NewRole() *Role {
	return &Role{}
}

// Role https://discordapp.com/developers/docs/topics/permissions#role-object
type Role struct {
	sync.RWMutex `json:"-"`

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

	r.RLock()
	role.Lock()

	role.ID = r.ID
	role.Name = r.Name
	role.Color = r.Color
	role.Hoist = r.Hoist
	role.Position = r.Position
	role.Permissions = r.Permissions
	role.Managed = r.Managed
	role.Mentionable = r.Mentionable
	role.guildID = r.guildID

	r.RUnlock()
	role.Unlock()

	return
}

func (r *Role) saveToDiscord(session Session) (err error) {
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
	} else {
		// modify/update role
		params := ModifyGuildRoleParams{
			Name:        r.Name,
			Permissions: r.Permissions,
			Color:       r.Color,
			Hoist:       r.Hoist,
			Mentionable: r.Mentionable,
		}
		role, err = session.ModifyGuildRole(r.guildID, r.ID, &params)
		if err != nil {
			return
		}
		if role.Position != r.Position {
			// update the position
			params := ModifyGuildRolePositionsParams{
				ID:       r.ID,
				Position: r.Position,
			}
			_, err = session.ModifyGuildRolePositions(r.guildID, &params)
			if err != nil {
				return
			}
			role.Position = r.Position
		}
	}

	return
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
