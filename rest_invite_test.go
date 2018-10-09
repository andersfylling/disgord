package disgord

import (
	"testing"
)

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
