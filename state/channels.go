package state

import (
	"errors"
	"runtime"
	"sync"

	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/resource"
	. "github.com/andersfylling/disgord/snowflake"
)

type ChannelUserCacher interface {
	Process(ud *UserDetail)
	User(ID Snowflake) (*resource.User, error)
}

type ChannelCacher interface {
	Process(cd *ChannelDetail)
	Chan() chan<- *ChannelDetail
	Channel(ID Snowflake) (*resource.Channel, error)
	Size() int
	Clear()
	Close() error
}

func NewChannelCache(userCacher ChannelUserCacher) *ChannelCache {
	cacher := &ChannelCache{
		channel:    make(chan *ChannelDetail),
		channels:   make(map[Snowflake]*channelHolder),
		userCacher: userCacher,
	}
	go cacher.channelCacher()

	return cacher
}

type ChannelCache struct {
	channel  chan *ChannelDetail
	channels map[Snowflake]*channelHolder

	// for saving users to cache
	userCacher ChannelUserCacher

	mu sync.RWMutex
}

// ChannelDetail holds information about the "incoming" channel object.
// how it was created, should it be considered dirty, etc.
type ChannelDetail struct {
	Channel *resource.Channel
	Dirty   bool
	Action  string //event.* for specific behavior, such as delete
}

// channelHolder has the channel struct along with partials for any snowflake Snowflake'able objects.
// as an example Channel.Recipients is set to nil, after all the user snowflakes have been
// moved to a map/slice. This saves memory, and makes it easy to lookup users later.
type channelHolder struct {
	// the channel object itself
	Channel *resource.Channel
	// Recipients, snowflake only
	// Since discord has limited the size to 10, we use a slice as that is faster.
	Recipients []Snowflake
}

func newChannelHolder(channel *resource.Channel) *channelHolder {
	holder := &channelHolder{
		Channel:    channel.DeepCopy(),
		Recipients: make([]Snowflake, len(channel.Recipients)),
	}

	for index, recipient := range holder.Channel.Recipients {
		holder.Recipients[index] = recipient.ID
	}

	holder.Channel.Recipients = nil

	return holder
}

func (st *ChannelCache) channelCacher() {
	for {
		var channelDetail *ChannelDetail
		var chanOpen bool

		// listen for channel changes
		select {
		case channelDetail, chanOpen = <-st.channel:
			if !chanOpen {
				break
			}
		}

		st.mu.Lock()
		channel := channelDetail.Channel
		action := channelDetail.Action
		dirty := channelDetail.Dirty
		switch action {
		case event.ChannelDeleteKey:
			if _, exists := st.channels[channel.ID]; exists {
				delete(st.channels, channel.ID)
			}
		default:
			if channel.Type == resource.ChannelTypeDM || channel.Type == resource.ChannelTypeGroupDM {
				// ensure the users for channel are stored as well
				// TODO-issue: blocking if high rate of user caching
				for _, recipient := range channel.Recipients {
					st.userCacher.Process(&UserDetail{
						User: recipient,
					})
				}
			}
			if _, exists := st.channels[channel.ID]; !exists {
				st.channels[channel.ID] = newChannelHolder(channel)
			} else if dirty {
				*st.channels[channel.ID] = *newChannelHolder(channel)
			}
		}
		st.mu.Unlock()
	}
}

func (st *ChannelCache) Process(cd *ChannelDetail) {
	st.channel <- cd
}

func (st *ChannelCache) Channel(ID Snowflake) (*resource.Channel, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if cachedChannelHolder, ok := st.channels[ID]; ok {
		channel := cachedChannelHolder.Channel.DeepCopy()
		channel.Recipients = []*resource.User{}
		// if DM, supply users
		if channel.Type == resource.ChannelTypeDM || channel.Type == resource.ChannelTypeGroupDM {
			for _, id := range cachedChannelHolder.Recipients {
				user, err := st.userCacher.User(id)
				if err != nil {
					// user not in cache
					continue
				}
				channel.Recipients = append(channel.Recipients, user.DeepCopy())
			}
		}

		return channel, nil
	}

	return nil, errors.New("channel with Snowflake{" + ID.String() + "} does not exist in cache")
}

func (st *ChannelCache) Chan() chan<- *ChannelDetail {
	return st.channel
}

func (st *ChannelCache) Size() int {
	return len(st.channels)
}

// Clear empty the cache
func (st *ChannelCache) Clear() {
	st.mu.Lock()
	st.channels = make(map[Snowflake]*channelHolder)
	runtime.GC() // Blocks thread
	st.mu.Unlock()
}

func (st *ChannelCache) Close() (err error) {
	close(st.channel)
	st.Clear()

	return
}
