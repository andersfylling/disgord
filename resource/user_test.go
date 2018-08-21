package resource

import (
	"io/ioutil"
	"testing"
)

func TestImplementsUserInterface(t *testing.T) {
	var u interface{} = &User{}
	if _, ok := u.(UserInterface); !ok {
		t.Error("User does not implement UserInterface")
	}
}

func verifyUserMashaller(t *testing.T, file string) {
	data, err := ioutil.ReadFile(file)
	check(err, t)

	user := User{}
	err = validateJSONMarshalling(data, &user)
	check(err, t)
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
