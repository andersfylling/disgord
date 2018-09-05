package rest

import (
	"testing"
)

func TestGetInvite(t *testing.T) {
	client, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	// TODO: invite codes....
	return // will only return 404 without an invite code

	inviteCode := ""
	_, err = GetInvite(client, inviteCode, false)
	if err != nil {
		t.Error(err)
	}

	_, err = GetInvite(client, inviteCode, true)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteInvite(t *testing.T) {
	client, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	// TODO: invite codes....
	return // will only return 404 without an invite code

	inviteCode := ""
	_, err = DeleteInvite(client, inviteCode)
	if err != nil {
		t.Error(err)
	}
}
