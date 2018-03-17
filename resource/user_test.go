package resource

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

func verifyUserMashaller(t *testing.T, file string) {
	data, err := ioutil.ReadFile(file)
	testutil.Check(err, t)

	user := User{}
	err = testutil.ValidateJSONMarshalling(data, &user)
	testutil.Check(err, t)
}

func TestUserMarshalling(t *testing.T) {
	files := []string{
		"testdata/user/user1.json",
		"testdata/user/user2.json",
	}

	for _, file := range files {
		verifyUserMashaller(t, file)
	}
}
