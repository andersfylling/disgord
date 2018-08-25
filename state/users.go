package state

import (
	"errors"
	"runtime"
	"sync"

	"github.com/andersfylling/disgord/resource"
	. "github.com/andersfylling/snowflake"
)

type UserCacher interface {
	Process(ud *UserDetail)
	Chan() chan<- *UserDetail
	User(ID Snowflake) (*resource.User, error)
	Size() int
	Clear()
	StartListener() error
	Close() error
}

// NewUserCache creates a new user cacher, and starts listening for inputs
func NewUserCache() *UserCache {
	cacher := &UserCache{
		users:   make(map[Snowflake]*resource.User),
		channel: make(chan *UserDetail, 100),
	}

	return cacher
}

// UserCache handles user caching
type UserCache struct {
	users   map[Snowflake]*resource.User
	channel chan *UserDetail
	mu      sync.RWMutex
}

// UserCache
// made this a struct, in case I want to add details/data in the future for optimizing caching.
type UserDetail struct {
	User  *resource.User
	Dirty bool // if the user is a part of another struct, like Member, we only need to check that it exists. otherwise dirty.
}

// userCacher handles incoming user objects, and copies them to the cache
func (st *UserCache) userCacher() {
	for {
		var userDetail *UserDetail
		var openChan bool

		select {
		case userDetail, openChan = <-st.channel:
		}

		if !openChan {
			break
		}

		// make sure it has a legal snowflake
		if !userDetail.User.Valid() {
			continue
		}

		st.mu.Lock()
		if _, exists := st.users[userDetail.User.ID]; !exists {
			// new user object
			st.users[userDetail.User.ID] = userDetail.User.DeepCopy()
		} else {
			if userDetail.Dirty {
				*(st.users[userDetail.User.ID]) = *(userDetail.User.DeepCopy())
			}
		}
		st.mu.Unlock()
	}
}

func (st *UserCache) Process(uc *UserDetail) {
	st.channel <- uc
}

func (st *UserCache) Chan() chan<- *UserDetail {
	return st.channel
}

// User get a copy from the cache, which can be safely distributed without ruining the up to date discord cache.
// See st.updaterUser(...) for more information why it's a copy only.
func (st *UserCache) User(ID Snowflake) (*resource.User, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if cachedUser, ok := st.users[ID]; ok {
		user := cachedUser.DeepCopy()
		return user, nil
	}

	return nil, errors.New("user with Snowflake{" + ID.String() + "} does not exist in cache")
}

func (st *UserCache) Size() int {
	return len(st.users)
}

// Clear empty the cache
func (st *UserCache) Clear() {
	st.mu.Lock()
	st.users = make(map[Snowflake]*resource.User)
	runtime.GC() // Blocks thread
	st.mu.Unlock()
}

func (st *UserCache) StartListener() (err error) {
	if st.users == nil || st.channel == nil {
		err = errors.New("users map and/or channel have not been instantiated")
		return
	}

	go st.userCacher()
	return nil
}

func (st *UserCache) Close() (err error) {
	close(st.channel)
	st.Clear()

	return
}
