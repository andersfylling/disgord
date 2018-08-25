package state

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/andersfylling/disgord/resource"
	. "github.com/andersfylling/snowflake"
)

func TestUsers_implementsUserCacher(t *testing.T) {
	uc := &UserCache{}
	if _, implemented := interface{}(uc).(UserCacher); !implemented {
		t.Error("UserCache does not implement interface UserCacher")
	}
}

func TestUsers_cacheSize(t *testing.T) {
	// incoming user object
	newUser := resource.NewUser()

	newUser.ID = Snowflake(11111111111111)
	newUser.Username = "new object from disgord"

	// check if it exists in cache
	cache := NewUserCache()
	cache.StartListener()
	u1, _ := cache.User(newUser.ID)
	if u1 != nil {
		t.Error("User was found in cache, even tho it was not saved in cache")
	}

	// add the user to cache
	cache.Process(&UserDetail{
		User: newUser,
	})
	time.Sleep(50 * time.Millisecond) // haxxor

	// check if the cache grew
	if cache.Size() == 0 {
		t.Error("Cache size has not grown after adding a user")
	}
	if cache.Size() != 1 {
		t.Error("expected cache to have size of 1, but got " + strconv.Itoa(cache.Size()))
	}

	// clear the cache
	cache.Close()
	if cache.Size() != 0 {
		t.Error("cache was cleared, but size is not 0. Size: " + strconv.Itoa(cache.Size()))
	}
}

func TestUsers_cacheClear(t *testing.T) {
	// generate a significant amount of random users,
	// add to cache, and clear it
	// compare memstat before and after
	N := 10000000 // 5000000 =< 1.5G

	var users []*resource.User
	// gen random partial users
	for i := 0; i < N; i++ {
		avatar := "sdfkijsdljflsdjfjsdlfjlksdjf"
		users = append(users, &resource.User{
			ID:            Snowflake(652342343 + uint64(i)),
			Username:      "iufhhsuaifuhs",
			Discriminator: "34234",
			Email:         "andersfylling@adnersfylling.internet",
			Avatar:        avatar,
		})
	}

	cache := NewUserCache()
	cache.StartListener()

	// store current mem stat
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// cache all users
	start := time.Now()
	for _, user := range users {
		cache.Process(&UserDetail{
			User:  user,
			Dirty: false,
		})
	}
	elapsed := time.Since(start)
	if false {
		fmt.Printf("copying to cache took %s", elapsed)
	}

	// store current mem stat
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	// clear cache and compare results
	cache.Close()
	var m3 runtime.MemStats
	runtime.ReadMemStats(&m3)

	// make sure the added memory from caching users was "freed"
	// TODO: not sure if this is the right way to check this..
	memUsedByCache := m2.Alloc - m1.Alloc
	memAfterClearingCache := m2.Alloc - m3.Alloc
	if memUsedByCache > memAfterClearingCache {
		t.Error("clearing cache did not 'free' memory use")
	}
}

func TestUserCache_Save(t *testing.T) {
	// make sure the data is copied correctly
	// that no race condition takes place once the cache is updated
	newUser := resource.NewUser()

	newUser.ID = Snowflake(11111111111111)
	newUser.Username = "new object from disgord"

	// add to cache
	cache := NewUserCache()
	cache.StartListener()
	cache.Process(&UserDetail{User: newUser})
	time.Sleep(50 * time.Millisecond) // haxor

	// update the cache, and make sure the newUser memory space isn't affected
	user1 := resource.NewUser()
	user1.ID = newUser.ID
	user1.Username = "different username"
	cache.Process(&UserDetail{
		User:  user1,
		Dirty: true,
	})
	time.Sleep(50 * time.Millisecond) // haxor

	// cache should not have been updated yet.
	cachedUser1, _ := cache.User(user1.ID)
	if cachedUser1.Username != "different username" {
		t.Errorf("the cached object does not hold the correct username: `%s`", cachedUser1.Username)
	}

	if cache.Size() != 1 {
		t.Error("size of user cache expected to be 1, have: " + strconv.Itoa(cache.Size()))
	}

	if newUser.Username != "new object from disgord" {
		t.Errorf("saving user to cache, altered the mem space of the original object: `%s`", newUser.Username)
	}
	// wait so the cache can update
	// kinda shitty way to ensure the cache update
	time.Sleep(time.Millisecond * 100)
	cachedUser, err := cache.User(user1.ID)
	if err != nil {
		t.Error(err)
	}

	if cachedUser.Username != "different username" {
		t.Errorf("the cached object does not hold the correct username: `%s`", cachedUser.Username)
	}

	// since cache update, verify size of cache map
	if cache.Size() != 1 {
		t.Error("size of user cache expected to be 1, have: " + strconv.Itoa(cache.Size()))
	}

	// make sure the oldest object still have not been affected by this update
	if newUser.Username != "new object from disgord" {
		t.Errorf("saving user to cache, altered the mem space of the original object: `%s`", newUser.Username)
	}

	// clearing the cache should not delete local variables
	cache.Close()
	if newUser == nil {
		t.Error("local var deleted, once the cache was cleared: newUser")
	}
	if cachedUser == nil {
		t.Error("local var deleted, once the cache was cleared: cachedUser")
	}
}
