// +build !integration

package disgord

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/json"
)

func verifyUserMashaller(t *testing.T, file string) {
	data, err := ioutil.ReadFile(file)
	check(err, t)

	user := &User{}
	err = json.Unmarshal(data, user)
	check(err, t)
}

func TestUserUpdateUnmarshal(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/user/user_update.json")
	check(err, t)

	u := &UserUpdate{}
	err = json.Unmarshal(data, u)
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

func TestUserPresence_InterfaceImplementations(t *testing.T) {
	var u interface{} = NewUserPresence()

	t.Run("Stringer", func(t *testing.T) {
		if _, ok := u.(fmt.Stringer); !ok {
			t.Error("UserPresence does not implement fmt.Stringer")
		}
	})

	t.Run("DeepCopier", func(t *testing.T) {
		if _, ok := u.(DeepCopier); !ok {
			t.Error("UserPresence does not implement DeepCopier")
		}
	})

	t.Run("Copier", func(t *testing.T) {
		if _, ok := u.(Copier); !ok {
			t.Error("UserPresence does not implement Copier")
		}
	})
}

func TestGetCurrentUserGuildsParams(t *testing.T) {
	params := &getCurrentUserGuildsBuilder{}
	params.r.setup(nil, nil, nil)
	var wants string

	wants = ""
	verifyQueryString(t, params.r.urlParams, wants)

	wants = "?before=438543957"
	params.SetBefore(438543957)
	verifyQueryString(t, params.r.urlParams, wants)

	wants += "&limit=6"
	params.SetLimit(6)
	verifyQueryString(t, params.r.urlParams, wants)

	wants = "?before=438543957"
	params.SetDefaultLimit()
	verifyQueryString(t, params.r.urlParams, wants)
}
