package resource

import "github.com/andersfylling/snowflake"

func NewRole() *Role {
	return &Role{}
}

// Role https://discordapp.com/developers/docs/topics/permissions#role-object
type Role struct {
	ID          snowflake.ID `json:"id"`
	Name        string       `json:"name"`
	Color       int          `json:"color"`
	Hoist       bool         `json:"hoist"`
	Position    int          `json:"position"`
	Permissions uint64       `json:"permissions"`
	Managed     bool         `json:"managed"`
	Mentionable bool         `json:"mentionable"`
}

func (r *Role) Mention() string {
	return "<@&" + r.ID.String() + ">"
}

func (r *Role) Clear() {
}
