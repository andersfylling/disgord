package channel

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/testutil"
)

func TestTypingStartMarshaler(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/typing_start.json")
	testutil.Check(err, t)

	v := &TypingStart{}
	err = testutil.ValidateJSONMarshalling(data, &v)
	testutil.Check(err, t)
}
