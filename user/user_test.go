package user

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/testutil"
)

func TestUserMarshalling(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/user1.json")
	testutil.Check(err, t)

	user := User{}
	err = testutil.ValidateJSONMarshalling(data, &user)
	testutil.Check(err, t)
}
