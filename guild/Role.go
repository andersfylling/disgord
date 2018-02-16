package guild

import "github.com/andersfylling/snowflake"

type Role struct {
	ID          snowflake.ID `json:"id"`
	Name        string       `json:"name"`
	Managed     bool         `json:"managed"`
	Mentionable bool         `json:"mentionable"`
	Hoist       bool         `json:"hoist"`
	Color       int          `json:"color"`
	Position    int          `json:"position"`
	Permissions uint64       `json:"permissions"`
}

func NewRole() *Role {
	return &Role{}
}

func (r *Role) Mention() string {
	return "<@&" + r.ID.String() + ">"
}

func (r *Role) Clear() {

}
