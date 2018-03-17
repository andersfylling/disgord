package resource

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/testutil"
)

func TestAuditLogConvertAuditLogParamsToStr(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/auditlog/auditlog1.json")
	testutil.Check(err, t)

	v := AuditLog{}
	err = testutil.ValidateJSONMarshalling(data, &v)
	testutil.Check(err, t)
}
