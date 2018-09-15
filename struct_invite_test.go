package disgord

import "testing"

func TestInvite_InterfaceImplementations(t *testing.T) {
	t.Run("InviteMetadata", func(t *testing.T) {
		var u interface{} = &InviteMetadata{}
		t.Run("DeepCopier", func(t *testing.T) {
			if _, ok := u.(DeepCopier); !ok {
				t.Error("does not implement DeepCopier")
			}
		})

		t.Run("Copier", func(t *testing.T) {
			if _, ok := u.(Copier); !ok {
				t.Error("does not implement Copier")
			}
		})
	})
	t.Run("Invite", func(t *testing.T) {
		var u interface{} = &Invite{}
		t.Run("DeepCopier", func(t *testing.T) {
			if _, ok := u.(DeepCopier); !ok {
				t.Error("does not implement DeepCopier")
			}
		})

		t.Run("Copier", func(t *testing.T) {
			if _, ok := u.(Copier); !ok {
				t.Error("does not implement Copier")
			}
		})
	})
}
