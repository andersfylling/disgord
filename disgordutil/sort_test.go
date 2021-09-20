// +build !integration

package disgordutil

import (
	"github.com/andersfylling/disgord"
	"testing"
)

func TestSort(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		sameOrder := func(a, b []*disgord.Role) bool {
			for i := range a {
				if a[i].Name != b[i].Name {
					return false
				}
			}
			return true
		}
		roles := []*disgord.Role{
			{Name: "b"},
			{Name: "k"},
			{Name: "j"},
			{Name: "a"},
			{Name: "e"},
			{Name: "g"},
		}
		rolesAsc := []*disgord.Role{
			{Name: "a"},
			{Name: "b"},
			{Name: "e"},
			{Name: "g"},
			{Name: "j"},
			{Name: "k"},
		}
		rolesDesc := []*disgord.Role{
			{Name: "k"},
			{Name: "j"},
			{Name: "g"},
			{Name: "e"},
			{Name: "b"},
			{Name: "a"},
		}

		t.Run("ascending", func(t *testing.T) {
			data := make([]*disgord.Role, len(roles))
			copy(data, roles)

			Sort(data, SortByName, OrderNone)
			if !sameOrder(data, rolesAsc) {
				t.Error("roles were not sorted into ascending order")
			}
		})

		t.Run("ascending-explicit", func(t *testing.T) {
			data := make([]*disgord.Role, len(roles))
			copy(data, roles)

			Sort(data, SortByName, OrderAscending)
			if !sameOrder(data, rolesAsc) {
				t.Error("roles were not sorted into ascending order")
			}
		})

		t.Run("descending", func(t *testing.T) {
			data := make([]*disgord.Role, len(roles))
			copy(data, roles)

			Sort(data, SortByName, OrderDescending)
			if !sameOrder(data, rolesDesc) {
				t.Error("roles were not sorted into descending order")
			}
		})

		t.Run("descending pointer", func(t *testing.T) {
			data := make([]*disgord.Role, len(roles))
			copy(data, roles)

			Sort(&data, SortByName, OrderDescending)
			if !sameOrder(data, rolesDesc) {
				t.Error("roles were not sorted into descending order")
			}
		})
	})
}
