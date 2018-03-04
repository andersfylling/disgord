package disgord

import (
	"sync"

	"errors"

	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/schema"
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
	User(ID snowflake.ID) (*schema.User, error)
	//
	//UpdateMySelf(*user.User)
	GetMySelf() *schema.User

	// channels to receive changes
	UserChan() chan<- *schema.User
	MemberChan() chan<- *schema.Member
	MessageChan() chan<- *schema.Message

	// Closer interface
	Close() error
}

func NewStateCache(evtDispatcher EvtDispatcher) *StateCache {
	st := &StateCache{
		guilds:   make(map[snowflake.ID]*schema.Guild),
		users:    make(map[snowflake.ID]*schema.User),
		channels: make(map[snowflake.ID]*schema.Channel),
		mySelf:   &schema.User{},

		userChan:   make(chan *schema.User),
		memberChan: make(chan *schema.Member),
		msgChan:    make(chan *schema.Message),
		guildChan:  make(chan *schema.Guild),
	}

	// listen for changes, and update the cache
	//go st.updaterGuild(evtDispatcher)
	go st.updaterUser(evtDispatcher)

	return st
}

type StateCache struct {
	guilds   map[snowflake.ID]*schema.Guild
	channels map[snowflake.ID]*schema.Channel // DM, one-one, or groups
	users    map[snowflake.ID]*schema.User
	mySelf   *schema.User

	guildsUpdateMutex sync.Mutex // update + delete
	guildsAddMutex    sync.Mutex // creation

	usersMutex sync.Mutex

	// channels
	userChan   chan *schema.User
	memberChan chan *schema.Member
	msgChan    chan *schema.Message
	guildChan  chan *schema.Guild
}

// Channel listeners for object updates
//

func (st *StateCache) updaterGuild(evtDispatcher EvtDispatcher) {
	for {
		var guild *schema.Guild
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
			guild = schema.NewGuildFromUnavailable(box.UnavailableGuild)
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
			st.guilds[guild.ID] = &schema.Guild{}
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

func (st *StateCache) updaterUser(evtDispatcher EvtDispatcher) {
	for {
		var user *schema.User
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
			st.users[user.ID] = &schema.User{}
			newUser = true
		}

		// false: the user exists, but the incoming user object hasn't changed. It's just cached cause of activity
		if triggeredByChange || newUser {
			st.users[user.ID].Replicate(user)
		}

		st.usersMutex.Unlock()
	}
}

func (st *StateCache) UserChan() chan<- *schema.User {
	return st.userChan
}
func (st *StateCache) MemberChan() chan<- *schema.Member {
	return st.memberChan
}
func (st *StateCache) MessageChan() chan<- *schema.Message {
	return st.msgChan
}
func (st *StateCache) GuildChan() chan<- *schema.Guild {
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
//// users
////
//
//// AddUser and return reference
//func (s *StateCache) AddUser(u *user.User) (updated *user.User) {
//	s.usersAddMutex.Lock()
//	defer s.usersAddMutex.Unlock()
//
//	if _, exists := s.users[u.ID]; exists {
//		updated, _ = s.UpdateUser(u)
//		return
//	}
//	s.users[u.ID] = u
//	updated = u
//	return
//}
//
//// UpdateUser and return the reference stored in cache
//func (s *StateCache) UpdateUser(new *user.User) (*user.User, error) {
//	s.usersUpdateMutex.Lock()
//	defer s.usersUpdateMutex.Unlock()
//
//	if _, exists := s.users[new.ID]; !exists {
//		return nil, errors.New("cannot update guild none-existant user in cache")
//	}
//
//	old := s.users[new.ID]
//
//	old.Update(new)
//	return old, nil
//}
//
//func (s *StateCache) DeleteUser(u *user.User) {
//	s.DeleteUserByID(u.ID)
//}
//
//func (s *StateCache) DeleteUserByID(ID snowflake.ID) {
//	s.usersUpdateMutex.Lock()
//	defer s.usersUpdateMutex.Unlock()
//	if u, ok := s.users[ID]; ok {
//		u.Clear()
//		delete(s.users, ID) // TODO: how good is the golang garbage collector?
//	}
//}

// User get a copy from the cache, which can be safely distributed without ruining the up to date discord cache.
// See st.updaterUser(...) for more information why it's a copy only.
func (st *StateCache) User(ID snowflake.ID) (*schema.User, error) {
	st.usersMutex.Lock()
	defer st.usersMutex.Unlock()

	if u, ok := st.users[ID]; ok {
		return u, nil
	}

	return nil, errors.New("guild with ID{" + ID.String() + "} does not exist in cache")
}

//func (s *StateCache) UpdateMySelf(new *user.User) {
//	s.mySelf.Update(new)
//}
func (s *StateCache) GetMySelf() *schema.User {
	return s.mySelf
}
