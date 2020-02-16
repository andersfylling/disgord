// +build !integration

package disgord

import (
	"testing"
)

func TestSort(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		sameOrder := func(a, b []*Role) bool {
			for i := range a {
				if a[i].Name != b[i].Name {
					return false
				}
			}
			return true
		}
		roles := []*Role{
			{Name: "b"},
			{Name: "k"},
			{Name: "j"},
			{Name: "a"},
			{Name: "e"},
			{Name: "g"},
		}
		rolesAsc := []*Role{
			{Name: "a"},
			{Name: "b"},
			{Name: "e"},
			{Name: "g"},
			{Name: "j"},
			{Name: "k"},
		}
		rolesDesc := []*Role{
			{Name: "k"},
			{Name: "j"},
			{Name: "g"},
			{Name: "e"},
			{Name: "b"},
			{Name: "a"},
		}

		t.Run("ascending", func(t *testing.T) {
			Sort(roles, SortByName)
			if !sameOrder(roles, rolesAsc) {
				t.Error("roles were not sorted into ascending order")
			}
		})

		t.Run("descending", func(t *testing.T) {
			Sort(roles, SortByName, OrderDescending)
			if !sameOrder(roles, rolesDesc) {
				t.Error("roles were not sorted into descending order")
			}
		})
	})
}
