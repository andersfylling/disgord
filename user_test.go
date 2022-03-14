//go:build !integration
// +build !integration

package disgord

import (
	"fmt"
	"testing"
)

func TestUserPresence_InterfaceImplementations(t *testing.T) {
	var u interface{} = &UserPresence{}

	t.Run("Stringer", func(t *testing.T) {
		if _, ok := u.(fmt.Stringer); !ok {
			t.Error("UserPresence does not implement fmt.Stringer")
		}
	})

	t.Run("DeepCopier", func(t *testing.T) {
		if _, ok := u.(DeepCopier); !ok {
			t.Error("UserPresence does not implement DeepCopier")
		}
	})

	t.Run("Copier", func(t *testing.T) {
		if _, ok := u.(Copier); !ok {
			t.Error("UserPresence does not implement Copier")
		}
	})
}
