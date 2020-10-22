// +build !integration

package disgord

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
)

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
