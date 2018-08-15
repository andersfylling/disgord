package resource

import (
	"io/ioutil"
	"testing"
)

func TestAuditLogConvertAuditLogParamsToStr(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/auditlog/auditlog1.json")
	check(err, t)

	v := AuditLog{}
	err = validateJSONMarshalling(data, &v)
	check(err, t)
}
