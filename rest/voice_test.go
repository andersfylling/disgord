package rest

import (
	"testing"
)

func TestListVoiceRegions(t *testing.T) {
	client, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	list, err := ListVoiceRegions(client)
	if err != nil {
		t.Error(err)
	}

	if len(list) == 0 {
		t.Error("expected at least one voice region")
	}
}
