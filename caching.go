package disgord

import (
	"sync"

	"errors"

	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/snowflake"
)

type StateCacher interface {
	//AddGuild(g *guild.Guild) *guild.Guild
	//UpdateGuild(g *guild.Guild) (*guild.Guild, error)
	//DeleteGuild(g *guild.Guild)
	//DeleteGuildByID(ID snowflake.ID)
	//Guild(ID snowflake.ID) (*guild.Guild, error)
	//
	//AddChannel(c *channel.Channel)
	//UpdateChannel(c *channel.Channel)
	//DeleteChannel(c *channel.Channel)
	//DeleteChannelByID(ID snowflake.ID)
	//
	//AddUser(*user.User) *user.User
	//UpdateUser(*user.User) (*user.User, error)
	//DeleteUser(*user.User)
	//DeleteUserByID(ID snowflake.ID)
	User(ID snowflake.ID) (*resource.User, error)
	Channel(ID snowflake.ID) (*resource.Channel, error)
	//
	//UpdateMySelf(*user.User)
	GetMySelf() *resource.User

	// channels to receive changes
	UserChan() chan<- *resource.User
	MemberChan() chan<- *resource.Member
	MessageChan() chan<- *resource.Message

	// Closer interface
	Close() error
}

func NewStateCache(evtDispatcher EvtDispatcher) *StateCache {
	st := &StateCache{
		guilds:   make(map[snowflake.ID]*resource.Guild),
		users:    make(map[snowflake.ID]*resource.User),
		channels: make(map[snowflake.ID]*resource.Channel),
		mySelf:   &resource.User{},

		userChan:   make(chan *resource.User),
		memberChan: make(chan *resource.Member),
		msgChan:    make(chan *resource.Message),
		guildChan:  make(chan *resource.Guild),
	}

	// listen for changes, and update the cache
	//go st.updaterGuild(evtDispatcher)
	go st.updaterUser(evtDispatcher)
	go st.updaterChannel(evtDispatcher)

	return st
}

type StateCache struct {
	guilds   map[snowflake.ID]*resource.Guild
	channels map[snowflake.ID]*resource.Channel // DM, one-one, or groups
	users    map[snowflake.ID]*resource.User
	mySelf   *resource.User

	usersMutex    sync.RWMutex
	channelsMutex sync.RWMutex
	guildsMutex   sync.RWMutex

	// channels
	userChan   chan *resource.User
	memberChan chan *resource.Member
	msgChan    chan *resource.Message
	guildChan  chan *resource.Guild
}

// Channel listeners for object updates
//

func (st *StateCache) updaterGuild(evtDispatcher EvtDispatcher) {
	for {
		var guild *resource.Guild
		var action string

		// listen for guild changes
		select {
		case box, alive := <-evtDispatcher.GuildCreateChan():
			if !alive {
				continue
			}
			guild = box.Guild
			action = event.GuildCreateKey
		case box, alive := <-evtDispatcher.GuildUpdateChan():
			if !alive {
				continue
			}
			guild = box.Guild
			action = event.GuildUpdateKey
		case box, alive := <-evtDispatcher.GuildDeleteChan():
			if !alive {
				continue
			}
			guild = resource.NewGuildFromUnavailable(box.UnavailableGuild)
			action = event.GuildDeleteKey
		case g, alive := <-st.guildChan:
			if !alive {
				continue
			}
			guild = g
			if g.Unavailable {
				action = event.GuildDeleteKey
			} else {
				// never GuildCreate as the bot cannot read a guild, it is not a member of
				action = event.GuildUpdateKey
			}
		}

		switch action {
		case event.GuildCreateKey:
			// Make sure changes to the cache, doesn't ruin the reactor pattern.
			st.guilds[guild.ID] = &resource.Guild{}
			*(st.guilds[guild.ID]) = *guild // don't alter the pointer, but merely data at the mem location.
		case event.GuildUpdateKey:
			//TODO: store cached arrays, delete, set new guild, and update respective arrays
		case event.GuildDeleteKey:
			if _, exists := st.guilds[guild.ID]; exists {
				delete(st.guilds, guild.ID)
			}
			// TODO: delete content in arrays as well, but not public data such as users
		}
	}
}

func (st *StateCache) cacheChannel(channel *resource.Channel) {
	if _, exists := st.channels[channel.ID]; !exists {
		st.channels[channel.ID] = &resource.Channel{}
	}

	var recipients []*resource.User
	// DM will holds user objects
	if channel.Type == resource.ChannelTypeDM || channel.Type == resource.ChannelTypeGroupDM {
		recipients = make([]*resource.User, len(channel.Recipients))
		// TODO: predefined length vs append speed. since the user objects will most likely exist.
		for index, recipient := range channel.Recipients {
			var user *resource.User
			user, _ = st.User(recipient.ID)

			// if the user is not in cache, he should be cached.
			// with this users can suddenly not exist, and causes unnecessary requests
			// TODO
			if user == nil {
				user = resource.NewUser()
				user.Replicate(recipient)
			}

			recipients[index] = user
		}
	}

	st.channels[channel.ID].Replicate(channel, recipients)
}

func (st *StateCache) updaterChannel(evtDispatcher EvtDispatcher) {
	for {
		var channel *resource.Channel
		var guild *resource.Guild
		var action string

		// listen for channel changes
		select {
		case box, alive := <-evtDispatcher.ChannelCreateChan():
			if !alive {
				continue
			}
			channel = box.Channel
			action = event.ChannelCreateKey
		case box, alive := <-evtDispatcher.ChannelUpdateChan():
			if !alive {
				continue
			}
			channel = box.Channel
			action = event.ChannelUpdateKey
		case box, alive := <-evtDispatcher.ChannelDeleteChan():
			if !alive {
				continue
			}
			channel = box.Channel
			action = event.ChannelDeleteKey
		case box, alive := <-evtDispatcher.GuildCreateChan():
			if !alive {
				continue
			}
			guild = box.Guild
			action = event.GuildCreateKey
		}

		st.channelsMutex.Lock()
		if action == event.ChannelDeleteKey {
			if _, exists := st.channels[channel.ID]; exists {
				delete(st.channels, channel.ID)
			}
		} else if action == event.GuildCreateKey {
			for _, ch := range guild.Channels {
				st.cacheChannel(ch)
			}
		} else {
			st.cacheChannel(channel)
		}
		st.channelsMutex.Unlock()
	}
}

func (st *StateCache) updaterUser(evtDispatcher EvtDispatcher) {
	for {
		var user *resource.User
		var triggeredByChange bool

		// listen for guild changes
		select {
		case box, alive := <-evtDispatcher.UserUpdateChan():
			if !alive {
				continue
			}
			user = box.User
			triggeredByChange = true
		case box, alive := <-evtDispatcher.MessageCreateChan():
			if !alive {
				continue
			}
			user = box.Message.Author
		case member, alive := <-st.memberChan:
			if !alive {
				continue
			}
			user = member.User
		case u, alive := <-st.userChan:
			if !alive {
				continue
			}
			user = u
		}

		// the users doesn't hold any pointers, and can be safely swapped out without the need to update
		// sub values. See st.updaterGuild for a scenario where this does not apply.

		// Keeping behavior stable
		// When a object is put into the cache, it can be updated at any time, so any variable that points to a
		// object in the cache will be auto updated without knowing it.
		// When an incoming user object updates an existing user object it will only alter data where the cached obj
		// points to. not changing the pointer address itself. This means that the incoming pointer, wont be updated
		// if the cache is updated.
		//
		// To keep this behavior on newly generated user objects, we need to initiate a new pointer before assigning
		// the data.
		//
		// It should therefore be noted that, referencing a cached object, will always hold the latest data, and can
		// change at any time. But when the object is retrieved from a socket event or a request, it will never be
		// altered on cache changes.
		//
		// This does cause inconsistent behavior when retrieving an object from a abstract method; where cached is
		// checked, and a http request is performed to get the content if missing. If it references the cache, it can
		// change, it the value comes from a http request it won't.
		//
		// TODO: should I therefore create a copy when data is requested from the cache?
		//		 It will make sure the cache is always valid, and reflects the discord server as much as possible,
		// 		 since devs cannot change the cached objects. But Mere copies.
		// 		 Issue:
		// 				For long running tasks, it can be of interest to always have the latest up to date version.
		//				This does however cause a possibility of the cache being altered, and not correctly reflecting
		// 				the Discord state.
		//		 Solution1:
		//				Have a wrapper that always asks for the latest user object whenever a operation is done.
		//				It's slow, but will reflect the latest change without writing changes to the cache.
		st.usersMutex.Lock()
		var newUser bool
		if _, exists := st.users[user.ID]; !exists {
			// new user object
			st.users[user.ID] = &resource.User{}
			newUser = true
		}

		// TODO: method for saving/updating user object

		// false: the user exists, but the incoming user object hasn't changed. It's just cached cause of activity
		if triggeredByChange || newUser {
			st.users[user.ID].Replicate(user)
		}

		st.usersMutex.Unlock()
	}
}

func (st *StateCache) UserChan() chan<- *resource.User {
	return st.userChan
}
func (st *StateCache) MemberChan() chan<- *resource.Member {
	return st.memberChan
}
func (st *StateCache) MessageChan() chan<- *resource.Message {
	return st.msgChan
}
func (st *StateCache) GuildChan() chan<- *resource.Guild {
	return st.guildChan
}

// https://golang.org/pkg/io/#Closer
func (st *StateCache) Close() error {
	// destroy channels
	close(st.guildChan)
	close(st.msgChan)
	close(st.memberChan)

	// clear cache

	return nil
}

// guilds
//
//
//// AddGuild and return reference
//func (s *StateCache) AddGuild(g *guild.Guild) *guild.Guild {
//	s.guildsAddMutex.Lock()
//	defer s.guildsAddMutex.Unlock()
//
//	if _, exists := s.guilds[g.ID]; exists {
//		gg, _ := s.UpdateGuild(g)
//		return gg
//	}
//	s.guilds[g.ID] = g
//	return g
//}
//
//// UpdateGuild and return the reference stored in cache
//func (s *StateCache) UpdateGuild(new *guild.Guild) (*guild.Guild, error) {
//	s.guildsUpdateMutex.Lock()
//	defer s.guildsUpdateMutex.Unlock()
//
//	if _, exists := s.guilds[new.ID]; !exists {
//		return nil, errors.New("cannot update guild none-existant guild in cache")
//	}
//
//	old := s.guilds[new.ID]
//
//	old.Update(new)
//	return old, nil
//}
//
//func (s *StateCache) DeleteGuild(g *guild.Guild) {
//	s.DeleteGuildByID(g.ID)
//}
//
//func (s *StateCache) DeleteGuildByID(ID snowflake.ID) {
//	if g, ok := s.guilds[ID]; ok {
//		g.Clear()
//		delete(s.guilds, ID) // TODO: how good is the golang garbage collector?
//	}
//}
//
//func (s *StateCache) Guild(ID snowflake.ID) (*guild.Guild, error) {
//	if g, ok := s.guilds[ID]; ok {
//		return g, nil
//	}
//
//	return nil, errors.New("guild with ID{" + ID.String() + "} does not exist in cache")
//}
//
//// channels
////
//// TODO: store guild channels in guild, DM in root, and voice in guild
//
//func (s *StateCache) AddChannel(c *channel.Channel) {
//	s.channels[c.ID] = c
//}
//
//func (s *StateCache) UpdateChannel(c *channel.Channel) {
//	s.channels[c.ID] = c
//}
//
//func (s *StateCache) DeleteChannel(c *channel.Channel) {
//	s.DeleteChannelByID(c.ID)
//}
//
//func (s *StateCache) DeleteChannelByID(ID snowflake.ID) {
//	if _, ok := s.channels[ID]; ok {
//		delete(s.channels, ID)
//	}
//}
//

func (st *StateCache) Channel(ID snowflake.ID) (*resource.Channel, error) {
	st.channelsMutex.Lock()
	defer st.channelsMutex.Unlock()

	if cachedChannel, ok := st.channels[ID]; ok {
		channel := resource.NewChannel()

		var recipients []*resource.User
		// TODO: this duplicates code from st.cacheChannel
		// DM will holds user objects
		if cachedChannel.Type == resource.ChannelTypeDM || cachedChannel.Type == resource.ChannelTypeGroupDM {
			recipients = make([]*resource.User, len(cachedChannel.Recipients))
			// TODO: predefined length vs append speed. since the user objects will most likely exist.
			for index, recipient := range cachedChannel.Recipients {
				var user *resource.User
				user, _ = st.User(recipient.ID)

				// if the user is not in cache, he should be cached.
				// with this users can suddenly not exist, and causes unnecessary requests
				// TODO
				if user == nil {
					user = resource.NewUser()
					user.Replicate(recipient)
				}

				recipients[index] = user
			}
		}
		channel.Replicate(cachedChannel, recipients)

		return channel, nil
	}

	return nil, errors.New("channel with ID{" + ID.String() + "} does not exist in cache")
}

// User get a copy from the cache, which can be safely distributed without ruining the up to date discord cache.
// See st.updaterUser(...) for more information why it's a copy only.
func (st *StateCache) User(ID snowflake.ID) (*resource.User, error) {
	st.usersMutex.Lock()
	defer st.usersMutex.Unlock()

	if cachedUser, ok := st.users[ID]; ok {
		user := resource.NewUser()
		user.Replicate(cachedUser)

		return user, nil
	}

	return nil, errors.New("user with ID{" + ID.String() + "} does not exist in cache")
}

func (s *StateCache) GetMySelf() *resource.User {
	return s.mySelf
}
