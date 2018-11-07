package disgord

import (
	"github.com/andersfylling/disgord/httd"
	"io/ioutil"
	"testing"
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
