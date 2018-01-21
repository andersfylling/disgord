package guild

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/testutil"
)

func TestBanObject(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/ban1.json")
	testutil.Check(err, t)

	ban := Ban{}
	err = testutil.ValidateJSONMarshalling(data, &ban)
	testutil.Check(err, t)
}
