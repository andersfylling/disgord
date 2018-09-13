package disgord

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestUser_InterfaceImplementations(t *testing.T) {
	var u interface{} = &User{}

	t.Run("UserInterface", func(t *testing.T) {
		if _, ok := u.(UserInterface); !ok {
			t.Error("User does not implement UserInterface")
		}
	})

	t.Run("DeepCopier", func(t *testing.T) {
		if _, ok := u.(DeepCopier); !ok {
			t.Error("User does not implement DeepCopier")
		}
	})

	t.Run("Copier", func(t *testing.T) {
		if _, ok := u.(Copier); !ok {
			t.Error("User does not implement Copier")
		}
	})

	t.Run("DiscordSaver", func(t *testing.T) {
		if _, ok := u.(discordSaver); !ok {
			t.Error("User does not implement discordSaver")
		}
	})
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

	// TODO
	// t.Run("DiscordSaver", func(t *testing.T) {
	// 	if _, ok := u.(discordSaver); !ok {
	// 		t.Error("UserPresence does not implement discordSaver")
	// 	}
	// })
}
