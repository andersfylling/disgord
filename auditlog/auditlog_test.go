package auditlog_test

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/auditlog"
	"github.com/andersfylling/disgord/testutil"
)

func TestConvertAuditLogParamsToStr(t *testing.T) {
	data, err := ioutil.ReadFile("examples/auditlog1.json")
	testutil.Check(err, t)

	v := auditlog.AuditLog{}
	err = testutil.ValidateJSONMarshalling(data, &v)
	testutil.Check(err, t)
}
