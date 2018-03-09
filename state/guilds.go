package state

import (
	"sync"

	"runtime"

	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/snowflake"
)

type GuildUserCacher interface {
	Process(ud *UserDetail)
}

type GuildChannelCacher interface {
	Process(ud *ChannelDetail)
}

type GuildCacher interface {
	Process(cd *GuildDetail)
	Chan() chan<- *GuildDetail
	Guild(ID snowflake.ID) (*resource.Guild, error)
	Clear()
	Close() error
}

// NewGuildCache creates a new guild cacher, and starts listening for inputs
func NewGuildCache(userCacher GuildUserCacher, channelCacher GuildChannelCacher) *GuildCache {
	cacher := &GuildCache{
		guilds:  make(map[snowflake.ID]*resource.Guild),
		channel: make(chan *GuildDetail),

		userCacher:    userCacher,
		channelCacher: channelCacher,
	}
	go cacher.guildCacher()

	return cacher
}

// GuildCache handles guild caching
type GuildCache struct {
	guilds  map[snowflake.ID]*resource.Guild
	channel chan *GuildDetail

	// saving channels and users
	channelCacher GuildChannelCacher
	userCacher    GuildUserCacher

	mu sync.RWMutex
}

type GuildDetail struct {
	Guild  *resource.Guild
	Dirty  bool
	Action string // event type
}

func (st *GuildCache) guildCacher() {
	for {
		var guildDetail *GuildDetail
		var openChan bool

		// listen for guild changes
		select {
		case guildDetail, openChan = <-st.channel:
			if !openChan {
				break
			}
		}

		st.mu.Lock()
		guild := guildDetail.Guild
		action := guildDetail.Action
		dirty := guildDetail.Dirty
		switch action {
		case event.GuildDeleteKey:
			if _, exists := st.guilds[guild.ID]; exists {
				delete(st.guilds, guild.ID)
			}
		default:
			if _, exists := st.guilds[guild.ID]; !exists {
				st.guilds[guild.ID] = guild.DeepCopy()
			} else if dirty {
				*st.guilds[guild.ID] = *guild.DeepCopy()
			}
		}
		st.mu.Unlock()
	}
}

func (st *GuildCache) Process(gd *GuildDetail) {
	st.channel <- gd
}

func (st *GuildCache) Chan() chan<- *GuildDetail {
	return st.channel
}

func (st *GuildCache) Guild(ID snowflake.ID) (*resource.Guild, error) {
	return nil, nil
}

// Clear empty the cache
func (st *GuildCache) Clear() {
	st.mu.Lock()
	// remove every guild channel from memory
	for _, guild := range st.guilds {
		for _, channel := range guild.Channels {
			st.channelCacher.Process(&ChannelDetail{
				Channel: channel,
				Action:  event.ChannelDeleteKey,
			})
		}

		guild.Channels = nil
	}

	st.guilds = make(map[snowflake.ID]*resource.Guild)
	runtime.GC() // Blocks thread
	st.mu.Unlock()
}

func (st *GuildCache) Close() (err error) {
	close(st.channel)
	st.Clear()

	return
}
