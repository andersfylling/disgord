package state

import (
	"runtime"
	"strconv"
	"testing"
	"time"

	"errors"
	"github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/snowflake"
)

func TestChannels_implementsChannelCacher(t *testing.T) {
	if _, implemented := interface{}(&ChannelCache{}).(ChannelCacher); !implemented {
		t.Error("ChannelCache does not implement interface ChannelCacher")
	}
}

type Mock_UserCacher struct {
	users map[snowflake.ID]*resource.User
}

func (muc *Mock_UserCacher) Process(ud *UserDetail) {
	muc.users[ud.User.ID] = ud.User
}
func (muc *Mock_UserCacher) User(id snowflake.ID) (*resource.User, error) {
	if _, exists := muc.users[id]; exists {
		return muc.users[id], nil
	}
	return nil, errors.New("no such user")
}

func TestChannels_cacheSize(t *testing.T) {
	// incoming user object
	newChannel := resource.NewChannel()

	newChannel.ID = snowflake.NewID(11111111111111)
	newChannel.Name = "new object from disgord"

	// check if it exists in cache
	userCacher := &Mock_UserCacher{}
	cache := NewChannelCache(userCacher)
	u1, _ := cache.Channel(newChannel.ID)
	if u1 != nil {
		t.Error("Channel was found in cache, even tho it was not saved in cache")
	}

	// add the user to cache
	cache.Process(&ChannelDetail{
		Channel: newChannel,
	})

	// check if the cache grew
	if cache.Size() == 0 {
		t.Error("Cache size has not grown after adding a user")
	}
	if cache.Size() != 1 {
		t.Error("expected cache to have size of 1, but got " + strconv.Itoa(cache.Size()))
	}

	// clear the cache
	cache.Clear()
	if cache.Size() != 0 {
		t.Error("cache was cleared, but size is not 0. Size: " + strconv.Itoa(cache.Size()))
	}
}

func TestChannels_cacheClear(t *testing.T) {
	// generate a significant amount of random users,
	// add to cache, and clear it
	// compare memstat before and after
	N := 1000000 // 0.6GiB
	userCacher := &Mock_UserCacher{}
	cache := NewChannelCache(userCacher)

	var channels []*resource.Channel
	// gen random partial users
	for i := 0; i < N; i++ {
		channel := resource.NewChannel()
		channel.Type = resource.ChannelTypeGuildCategory
		channel.ID = snowflake.NewID(234234 + uint64(i))
		channel.Name = "242sdflkjsfjlksdhfkjsdf"
		channel.Topic = "sldkflks lks jdlf j klsdjf lskdjf klljkds flksdj fjkl"

		channels = append(channels, channel)
	}

	// store current mem stat
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// cache all channels
	for _, channel := range channels {
		cache.Process(&ChannelDetail{
			Channel: channel,
		})
	}

	// store current mem stat
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	// clear cache and compare results
	cache.Clear()
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

func TestChannelCache_Save(t *testing.T) {
	// make sure the data is copied correctly
	// that no race condition takes place once the cache is updated
	newChannel := resource.NewChannel()

	newChannel.ID = snowflake.NewID(11111111111111)
	newChannel.Name = "new object from disgord"

	// add to cache
	userCacher := &Mock_UserCacher{}
	cache := NewChannelCache(userCacher)
	cache.Process(&ChannelDetail{Channel: newChannel})

	// update the cache, and make sure the newChannel memory space isn't affected
	user1 := resource.NewChannel()
	user1.ID = newChannel.ID
	user1.Name = "different username"
	cache.Process(&ChannelDetail{
		Channel: user1,
		Dirty:   true,
	})

	// cache should not have been updated yet.
	cachedChannel1, _ := cache.Channel(user1.ID)
	if cachedChannel1.Name != "new object from disgord" {
		t.Errorf("the cached object does not hold the correct username: `%s`", cachedChannel1.Name)
	}

	if cache.Size() != 1 {
		t.Error("size of user cache expected to be 1, have: " + strconv.Itoa(cache.Size()))
	}

	if newChannel.Name != "new object from disgord" {
		t.Errorf("saving user to cache, altered the mem space of the original object: `%s`", newChannel.Name)
	}
	// wait so the cache can update
	// kinda shitty way to ensure the cache update
	time.Sleep(time.Millisecond * 100)
	cachedChannel, err := cache.Channel(user1.ID)
	if err != nil {
		t.Error(err)
	}

	if cachedChannel.Name != "different username" {
		t.Errorf("the cached object does not hold the correct username: `%s`", cachedChannel.Name)
	}

	// since cache update, verify size of cache map
	if cache.Size() != 1 {
		t.Error("size of user cache expected to be 1, have: " + strconv.Itoa(cache.Size()))
	}

	// make sure the oldest object still have not been affected by this update
	if newChannel.Name != "new object from disgord" {
		t.Errorf("saving user to cache, altered the mem space of the original object: `%s`", newChannel.Name)
	}

	// clearing the cache should not delete local variables
	cache.Clear()
	if newChannel == nil {
		t.Error("local var deleted, once the cache was cleared: newChannel")
	}
	if cachedChannel == nil {
		t.Error("local var deleted, once the cache was cleared: cachedChannel")
	}
}

func TestChannelCache_inputOutput(t *testing.T) {
	// make sure that the recipients are the same after as before caching
	channel := resource.NewChannel()

	channel.ID = snowflake.NewID(11111111111111)
	channel.Name = "new object from disgord"
	channel.Type = resource.ChannelTypeGroupDM
	for i := 0; i < 10; i++ {
		user := resource.NewUser()
		user.ID = snowflake.NewID(3546345 + uint64(i+i*i))
		channel.Recipients = append(channel.Recipients, user)
	}

	// add to cache
	userCacher := &Mock_UserCacher{
		users: make(map[snowflake.ID]*resource.User),
	}
	cache := NewChannelCache(userCacher)
	cache.Process(&ChannelDetail{
		Channel: channel,
		Dirty:   true,
	})

	cachedChannel, err := cache.Channel(channel.ID)
	if err != nil {
		t.Error(err)
	}

	if len(channel.Recipients) != len(cachedChannel.Recipients) {
		t.Error("incorrect number of users in DM channel")
	}

	for index, recipient := range cachedChannel.Recipients {
		if channel.Recipients[index].ID != recipient.ID {
			t.Error("incorrect user id")
		}
	}

}
