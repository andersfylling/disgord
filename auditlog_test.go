// +build !integration

package disgord

import (
	"context"
	"errors"
	"github.com/andersfylling/disgord/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
)

func TestAuditLogConvertAuditLogParamsToStr(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/auditlog/auditlog1.json")
	check(err, t)

	v := AuditLog{}
	err = json.Unmarshal(data, &v)
	executeInternalUpdater(v)
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
		var u interface{} = &AuditLogChanges{}
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
	params.r.setup(nil, nil, nil)
	var wants string

	wants = ""
	verifyQueryString(t, params.r.urlParams, wants)

	s := "438543957"
	ss, _ := GetSnowflake(s)
	params.SetUserID(ss)
	wants = "?user_id=" + s
	verifyQueryString(t, params.r.urlParams, wants)

	params.SetActionType(6)
	wants += "&action_type=6"
	wantsAlternative := "?action_type=6&user_id=" + s
	got := params.r.urlParams.URLQueryString()
	if !(wants == got || wantsAlternative == got) {
		t.Errorf("incorrect query param string. Got '%s', wants '%s' or '%s'", params.r.urlParams.URLQueryString(), wants, wantsAlternative)
	}

	params.SetActionType(0)
	wants = "?user_id=" + s + "&action_type=0"
	wantsAlternative = "?action_type=0&user_id=" + s
	got = params.r.urlParams.URLQueryString()
	if !(wants == got || wantsAlternative == got) {
		t.Errorf("incorrect query param string. Got '%s', wants '%s' or '%s'", params.r.urlParams.URLQueryString(), wants, wantsAlternative)
	}
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
		builder.r.itemFactory = auditLogFactory
		builder.r.IgnoreCache().setup(client, &httd.Request{
			Method:   httd.MethodGet,
			Endpoint: endpoint.GuildAuditLogs(Snowflake(7)),
			Ctx:      context.Background(),
		}, nil)

		_, err := builder.Execute()
		if err != nil {
			t.Error(err)
		}

		if client.req.Endpoint != "/guilds/7/audit-logs" {
			t.Error("incorrect endpoint: ", client.req.Endpoint)
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
		builder.r.itemFactory = auditLogFactory
		builder.r.IgnoreCache().setup(client, &httd.Request{
			Method:   httd.MethodGet,
			Endpoint: endpoint.GuildAuditLogs(Snowflake(7)),
			Ctx:      context.Background(),
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
			t.Errorf("expected 4 Users, got %d", len(logs.Users))
		}
		if len(logs.Webhooks) != 0 {
			t.Errorf("expected 0 webhooks, got %d", len(logs.Webhooks))
		}
	})
	t.Run("missing-permission", func(t *testing.T) {
		errorMsg := "missing permission flag?"
		client := &reqMocker{
			body: []byte(`{"code":403,"message":"` + errorMsg + `"}`),
			resp: &http.Response{
				StatusCode: 403,
			},
			err: errors.New("permission issue"),
		}

		builder := &guildAuditLogsBuilder{}
		builder.r.IgnoreCache().setup(client, &httd.Request{
			Method:   httd.MethodGet,
			Endpoint: endpoint.GuildAuditLogs(Snowflake(7)),
			Ctx:      context.Background(),
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
		builder.r.itemFactory = auditLogFactory
		builder.r.IgnoreCache().setup(client, &httd.Request{
			Method:   httd.MethodGet,
			Endpoint: endpoint.GuildAuditLogs(Snowflake(7)),
			Ctx:      context.Background(),
		}, nil)

		_, err := builder.Execute()
		if err != nil {
			t.Fatal("unexpected error: " + err.Error())
		}

		// TODO: implement ErrREST check
	})
}

func TestAuditlog_Unmarshal(t *testing.T) {
	data := []byte(`{
      "target_id": "547614326257877003",
      "changes": [
        {
          "new_value": "andreesdsdsd",
          "old_value": "test",
          "key": "name"
        }
      ],
      "user_id": "486832262592069632",
      "id": "547614855067205678",
      "action_type": 61
    }`)
	var v2 *AuditLogEntry
	if err := json.Unmarshal(data, &v2); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(v2)

	data, err := ioutil.ReadFile("testdata/auditlog/logs-limit-10.json")
	if err != nil {
		t.Fatal("missing test data")
		return
	}

	var v *AuditLog
	if err := json.Unmarshal(data, &v); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(v)

	if v.Bans() != nil {
		t.Error("these logs contains no bans")
	}
}
