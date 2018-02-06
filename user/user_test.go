package user

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/testutil"
)

func TestImplementsUserInterface(t *testing.T) {
	var u interface{} = &User{}
	if _, ok := u.(UserInterface); !ok {
		t.Error("User does not implement UserInterface")
	}
}

func TestUserMarshalling(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/user1.json")
	testutil.Check(err, t)

	user := User{}
	err = testutil.ValidateJSONMarshalling(data, &user)
	testutil.Check(err, t)
}
