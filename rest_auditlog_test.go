package disgord

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAuditLogParams(t *testing.T) {
	params := &GuildAuditLogsParams{}
	var wants string

	wants = ""
	verifyQueryString(t, params, wants)

	s := "438543957"
	params.UserID, _ = GetSnowflake(s)
	wants = "?user_id=" + s
	verifyQueryString(t, params, wants)

	params.ActionType = 6
	wants += "&action_type=6"
	verifyQueryString(t, params, wants)

	params.ActionType = 0
	wants = "?user_id=" + s
	verifyQueryString(t, params, wants)
}

func TestGuildAuditLogs(t *testing.T) {
	t.Run("configuration", func(t *testing.T) {
		// successfull response
		client := &reqMocker{
			body: []byte(`{}`),
			resp: &http.Response{
				StatusCode: 200,
			},
		}

		_, err := GuildAuditLogs(client, 7, &GuildAuditLogsParams{})
		if err != nil {
			t.Error(err)
		}

		if client.req.Endpoint != "/guilds/7/audit-logs" {
			t.Error("incorrect endpoint")
		}

		if client.req.Ratelimiter != "g:7:a-l" { // why even test this?
			t.Error("incorrect rate limit key")
		}
	})
	t.Run("success", func(t *testing.T) {
		body, err := ioutil.ReadFile("testdata/auditlog/auditlog1.json")
		check(err, t)

		// successful response
		client := &reqMocker{
			body: body,
			resp: &http.Response{
				StatusCode: 200,
			},
		}

		params := &GuildAuditLogsParams{}
		logs, err := GuildAuditLogs(client, 7, params)
		if err != nil {
			t.Error(err)
		}

		if logs == nil {
			t.Fatal("logs was nil but expected content")
		}

		if len(logs.AuditLogEntries) != 10 {
			t.Errorf("expected 10 log entries, got %d", len(logs.AuditLogEntries))
		}
		if len(logs.Users) != 4 {
			t.Errorf("expected 4 users, got %d", len(logs.Users))
		}
		if len(logs.Webhooks) != 0 {
			t.Errorf("expected 0 webhooks, got %d", len(logs.Webhooks))
		}
	})
	t.Run("missing-permission", func(t *testing.T) {
		errorMsg := "missing permissiong flag?"
		client := &reqMocker{
			body: []byte(`{"code":403,"message":"` + errorMsg + `"}`),
			resp: &http.Response{
				StatusCode: 403,
			},
		}

		params := &GuildAuditLogsParams{}
		_, err := GuildAuditLogs(client, 7, params)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}

		// TODO: implement ErrREST check
	})
	t.Run("no-content", func(t *testing.T) {
		client := &reqMocker{
			body: []byte(`{}`),
			resp: &http.Response{
				StatusCode: 204,
			},
		}

		params := &GuildAuditLogsParams{}
		_, err := GuildAuditLogs(client, 7, params)
		if err != nil {
			t.Fatal("unexpected error: " + err.Error())
		}

		// TODO: implement ErrREST check
	})
}

func TestLive_GuildAuditLogs(t *testing.T) {
	client, keys, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	params := &GuildAuditLogsParams{}
	log, err := GuildAuditLogs(client, keys.GuildAdmin, params)
	if err != nil {
		t.Error(err)
	}

	if log == nil {
		t.Error("did not get a datastructure from rest.GuildAuditLogs()")
	}
}
