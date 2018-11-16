package disgord

import (
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/ratelimit"
	"github.com/andersfylling/snowflake/v3"
)

func TestAuditLogConvertAuditLogParamsToStr(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/auditlog/auditlog1.json")
	check(err, t)

	v := AuditLog{}
	err = httd.Unmarshal(data, &v)
	check(err, t)
}

func TestAuditLog_InterfaceImplementations(t *testing.T) {
	t.Run("AuditLog", func(t *testing.T) {
		var u interface{} = &AuditLog{}
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
	t.Run("AuditLogEntry", func(t *testing.T) {
		var u interface{} = &AuditLogEntry{}
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
	t.Run("AuditLogOption", func(t *testing.T) {
		var u interface{} = &AuditLogOption{}
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
	t.Run("AuditLogChange", func(t *testing.T) {
		var u interface{} = &AuditLogChange{}
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

func TestAuditLogParams(t *testing.T) {
	params := &guildAuditLogsBuilder{}
	params.setup(nil, nil, nil, nil)
	var wants string

	wants = ""
	verifyQueryString(t, params.urlParams, wants)

	s := "438543957"
	ss, _ := snowflake.GetSnowflake(s)
	params.UserID(ss)
	wants = "?user_id=" + s
	verifyQueryString(t, params.urlParams, wants)

	params.ActionType(6)
	wants += "&action_type=6"
	verifyQueryString(t, params.urlParams, wants)

	params.ActionType(0)
	wants = "?user_id=" + s + "&action_type=0"
	verifyQueryString(t, params.urlParams, wants)
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

		builder := &guildAuditLogsBuilder{}
		builder.itemFactory = auditLogFactory
		builder.IgnoreCache().setup(nil, client, &httd.Request{
			Method:      http.MethodGet,
			Ratelimiter: ratelimit.GuildAuditLogs(7),
			Endpoint:    endpoint.GuildAuditLogs(snowflake.ID(7)),
		}, nil)

		_, err := builder.Execute()
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

		builder := &guildAuditLogsBuilder{}
		builder.itemFactory = auditLogFactory
		builder.IgnoreCache().setup(nil, client, &httd.Request{
			Method:      http.MethodGet,
			Ratelimiter: ratelimit.GuildAuditLogs(7),
			Endpoint:    endpoint.GuildAuditLogs(snowflake.ID(7)),
		}, nil)

		logs, err := builder.Execute()
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
			err: errors.New("permissing issue"),
		}

		builder := &guildAuditLogsBuilder{}
		builder.IgnoreCache().setup(nil, client, &httd.Request{
			Method:      http.MethodGet,
			Ratelimiter: ratelimit.GuildAuditLogs(7),
			Endpoint:    endpoint.GuildAuditLogs(snowflake.ID(7)),
		}, nil)

		logs, err := builder.Execute()
		if logs != nil {
			t.Error("expected logs to be nil (not initiated)")
		}
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

		builder := &guildAuditLogsBuilder{}
		builder.itemFactory = auditLogFactory
		builder.IgnoreCache().setup(nil, client, &httd.Request{
			Method:      http.MethodGet,
			Ratelimiter: ratelimit.GuildAuditLogs(7),
			Endpoint:    endpoint.GuildAuditLogs(snowflake.ID(7)),
		}, nil)

		_, err := builder.Execute()
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

	builder := &guildAuditLogsBuilder{}
	builder.IgnoreCache().setup(nil, client, &httd.Request{
		Method:      http.MethodGet,
		Ratelimiter: ratelimit.GuildAuditLogs(keys.GuildAdmin),
		Endpoint:    endpoint.GuildAuditLogs(snowflake.ID(keys.GuildAdmin)),
	}, nil)

	log, err := builder.Execute()
	if err != nil {
		t.Error(err)
	}

	if log == nil {
		t.Error("did not get a datastructure from rest.GuildAuditLogs()")
	}
}
