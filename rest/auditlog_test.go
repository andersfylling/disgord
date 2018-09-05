package rest

import (
	"os"
	"testing"

	"github.com/andersfylling/disgord/constant"
	. "github.com/andersfylling/snowflake"
)

func verifyQueryString(t *testing.T, params URLParameters, wants string) {
	got := params.GetQueryString()
	if got != wants {
		t.Errorf("incorrect query param string. Got '%s', wants '%s'", got, wants)
	}
}

func TestAuditLogParams(t *testing.T) {
	params := &GuildAuditLogsParams{}
	var wants string

	wants = ""
	verifyQueryString(t, params, wants)

	s := "438543957"
	params.UserID, _ = GetSnowflake(s)
	wants = "?user_id=" + s

	params.ActionType = 6
	wants += "&action_type=6"
	verifyQueryString(t, params, wants)

	params.ActionType = 0
	wants = "?user_id=" + s
	verifyQueryString(t, params, wants)
}

func TestGuildAuditLogs(t *testing.T) {
	client, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	s, err := GetSnowflake(os.Getenv(constant.DisgordTestGuildAdmin))
	if err != nil {
		t.Skip()
		return
	}

	params := &GuildAuditLogsParams{}
	log, err := GuildAuditLogs(client, s, params)
	if err != nil {
		t.Error(err)
	}

	if log == nil {
		t.Error("did not get a datastructure from rest.GuildAuditLogs()")
	}
}
