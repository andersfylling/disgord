package state

import (
	"errors"

	"runtime"

	"sync"

	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/snowflake"
)

type ChannelUserCacher interface {
	Process(ud *UserDetail)
}

type ChannelCacher interface {
	Process(cd *ChannelDetail)
	Chan() chan<- *ChannelDetail
	Channel(ID snowflake.ID) (*resource.Channel, error)
	Size() int
	Clear()
	Close() error
}

func NewChannelCache(userCacher ChannelUserCacher) *ChannelCache {
	cacher := &ChannelCache{
		channel:    make(chan *ChannelDetail),
		channels:   make(map[snowflake.ID]*resource.Channel),
		userCacher: userCacher,
	}
	go cacher.channelCacher()

	return cacher
}

type ChannelCache struct {
	channel  chan *ChannelDetail
	channels map[snowflake.ID]*resource.Channel

	// for saving users to cache
	userCacher ChannelUserCacher

	mu sync.RWMutex
}

type ChannelDetail struct {
	Channel *resource.Channel
	Dirty   bool
	Action  string //event.* for specific behavior, such as delete
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
				st.channels[channel.ID] = channel.DeepCopy()
			} else if dirty {
				*st.channels[channel.ID] = *channel.DeepCopy()
			}
		}
		st.mu.Unlock()
	}
}

func (st *ChannelCache) Process(cd *ChannelDetail) {
	st.channel <- cd
}

func (st *ChannelCache) Channel(ID snowflake.ID) (*resource.Channel, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if cachedChannel, ok := st.channels[ID]; ok {
		channel := cachedChannel.DeepCopy()
		return channel, nil
	}

	return nil, errors.New("channel with ID{" + ID.String() + "} does not exist in cache")
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
	st.channels = make(map[snowflake.ID]*resource.Channel)
	runtime.GC() // Blocks thread
	st.mu.Unlock()
}

func (st *ChannelCache) Close() (err error) {
	close(st.channel)
	st.Clear()

	return
}
