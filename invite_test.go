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

func TestGetInvite(t *testing.T) {
	_, _, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	// TODO: invite codes....
	// will only return 404 without an invite code
	//
	// inviteCode := ""
	// _, err = GetInvite(client, inviteCode, false)
	// if err != nil {
	// 	t.Error(err)
	// }
	//
	// _, err = GetInvite(client, inviteCode, true)
	// if err != nil {
	// 	t.Error(err)
	// }
}

func TestDeleteInvite(t *testing.T) {
	_, _, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	// TODO: invite codes....
	// will only return 404 without an invite code
	//
	// inviteCode := ""
	// _, err = DeleteInvite(client, inviteCode)
	// if err != nil {
	// 	t.Error(err)
	// }
}
