package disgord

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/httd"
)

func verifyUserMashaller(t *testing.T, file string) {
	data, err := ioutil.ReadFile(file)
	check(err, t)

	user := User{}
	err = httd.Unmarshal(data, &user)
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

func equals(a, b *User) bool {
	return a.ID == b.ID && a.Username == b.Username && a.Discriminator == b.Discriminator && a.Email == b.Email && a.Token == b.Token && a.Verified == b.Verified && a.MFAEnabled == b.MFAEnabled && a.Bot == b.Bot && b.Avatar == a.Avatar
}

type userCopyOverToCacheTestSets struct {
	user  User
	bytes []byte
}

func TestUser_copyOverToCache(t *testing.T) {
	datas := make([]userCopyOverToCacheTestSets, 6)

	user1 := User{}
	unmarshal([]byte(`{"id": "133218433276182528","username":"alak","discriminator":"1149","avatar":"38d04eba240fa3cad581a947025644ad","verified":true}`), &user1)
	datas[0] = userCopyOverToCacheTestSets{user1, []byte(`{"id": "133218433276182528","username":"alak","discriminator":"1149","avatar":"38d04eba240fa3cad581a947025644ad","verified":true}`)}

	user2 := User(user1)
	user2.Bot = true
	datas[1] = userCopyOverToCacheTestSets{user2, []byte(`{"id": "133218433276182528","username":"alak","discriminator":"1149","avatar":"38d04eba240fa3cad581a947025644ad","bot":true}`)}

	user3 := User(user2)
	user3.Discriminator = Discriminator(1849)
	datas[2] = userCopyOverToCacheTestSets{user3, []byte(`{"id": "133218433276182528","username":"alak","discriminator":"1849","avatar":"38d04eba240fa3cad581a947025644ad"}`)}

	user4 := User(user3)
	user4.Avatar = ""
	datas[3] = userCopyOverToCacheTestSets{user4, []byte(`{"id": "133218433276182528","username":"alak","discriminator":"1849","avatar":null}`)}

	user5 := User(user4)
	datas[4] = userCopyOverToCacheTestSets{user5, []byte(`{"id": "133218433276182528"}`)}

	user6 := User(user5)
	user6.Username = "sdfsd"
	a := "aaaaaaaaa"
	user6.Avatar = a
	user6.Discriminator = Discriminator(1249)
	user6.Verified = false
	datas[5] = userCopyOverToCacheTestSets{user6, []byte(`{"id": "133218433276182528","username":"sdfsd","discriminator":"1249","avatar":"aaaaaaaaa","verified":false}`)}

	var cache User
	var err error
	for i := range datas {
		bytes := datas[i].bytes
		expected := datas[i].user
		var user User
		err = unmarshal(bytes, &user)
		if err != nil {
			t.Error(err)
			return
		}

		user.copyOverToCache(&cache)

		//fmt.Printf("##: %+v\n", cacheLink)

		if !equals(&cache, &expected) {
			t.Errorf("different users. \nGot \t%+v, \nWants \t%+v", cache, expected)
		}
	}
}

func TestGetCurrentUserGuildsParams(t *testing.T) {
	params := &getCurrentUserGuildsBuilder{}
	params.r.setup(nil, nil, nil, nil)
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
