package disgord

import (
	"errors"
	"sync"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/guild"
	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/snowflake"
)

type StateCacher interface {
	AddGuild(g *guild.Guild) *guild.Guild
	UpdateGuild(g *guild.Guild) (*guild.Guild, error)
	DeleteGuild(g guild.GuildIDer)
	DeleteGuildByID(ID snowflake.ID)
	Guild(ID snowflake.ID) (*guild.Guild, error)

	AddChannel(c *channel.Channel)
	UpdateChannel(c *channel.Channel)
	DeleteChannel(c *channel.Channel)
	DeleteChannelByID(ID snowflake.ID)

	AddUser(g user.UserInterface) user.UserInterface
	UpdateUser(g user.UserInterface) (user.UserInterface, error)
	DeleteUser(g user.UserInterface)
	DeleteUserByID(ID snowflake.ID)
	User(ID snowflake.ID) (user.UserInterface, error)
}

func NewStateCache() *StateCache {
	return &StateCache{
		guilds:   make(map[snowflake.ID]*guild.Guild),
		users:    make(map[snowflake.ID]user.UserInterface),
		channels: make(map[snowflake.ID]*channel.Channel),
	}
}

type StateCache struct {
	guilds   map[snowflake.ID]*guild.Guild
	users    map[snowflake.ID]user.UserInterface
	channels map[snowflake.ID]*channel.Channel

	guildsUpdateMutex sync.Mutex // update + delete
	guildsAddMutex    sync.Mutex // creation

	usersUpdateMutex sync.Mutex // update + delete
	usersAddMutex    sync.Mutex // creation
}

// guilds
//

// AddGuild and return reference
func (s *StateCache) AddGuild(g *guild.Guild) *guild.Guild {
	s.guildsAddMutex.Lock()
	defer s.guildsAddMutex.Unlock()

	if _, exists := s.guilds[g.ID()]; exists {
		gg, _ := s.UpdateGuild(g)
		return gg
	}
	s.guilds[g.ID()] = g
	return g
}

// UpdateGuild and return the reference stored in cache
func (s *StateCache) UpdateGuild(new *guild.Guild) (*guild.Guild, error) {
	s.guildsUpdateMutex.Lock()
	defer s.guildsUpdateMutex.Unlock()

	if _, exists := s.guilds[new.ID()]; !exists {
		return nil, errors.New("cannot update guild none-existant guild in cache")
	}

	old := s.guilds[new.ID()]

	old.Update(new)
	return old, nil
}

func (s *StateCache) DeleteGuild(g guild.GuildIDer) {
	s.DeleteGuildByID(g.ID())
}

func (s *StateCache) DeleteGuildByID(ID snowflake.ID) {
	if g, ok := s.guilds[ID]; ok {
		g.Clear()
		delete(s.guilds, ID) // TODO: how good is the golang garbage collector?
	}
}

func (s *StateCache) Guild(ID snowflake.ID) (*guild.Guild, error) {
	if g, ok := s.guilds[ID]; ok {
		return g, nil
	}

	return nil, errors.New("guild with ID{" + ID.String() + "} does not exist in cache")
}

// channels
//
// TODO: store guild channels in guild, DM in root, and voice in guild

func (s *StateCache) AddChannel(c *channel.Channel) {
	s.channels[c.ID] = c
}

func (s *StateCache) UpdateChannel(c *channel.Channel) {
	s.channels[c.ID] = c
}

func (s *StateCache) DeleteChannel(c *channel.Channel) {
	s.DeleteChannelByID(c.ID)
}

func (s *StateCache) DeleteChannelByID(ID snowflake.ID) {
	if _, ok := s.channels[ID]; ok {
		delete(s.channels, ID)
	}
}

// users
//

// AddUser and return reference
func (s *StateCache) AddUser(u user.UserInterface) (updated user.UserInterface) {
	s.usersAddMutex.Lock()
	defer s.usersAddMutex.Unlock()

	if _, exists := s.users[u.ID()]; exists {
		updated, _ = s.UpdateUser(u)
		return
	}
	s.users[u.ID()] = u
	updated = u
	return
}

// UpdateUser and return the reference stored in cache
func (s *StateCache) UpdateUser(new user.UserInterface) (user.UserInterface, error) {
	s.usersUpdateMutex.Lock()
	defer s.usersUpdateMutex.Unlock()

	if _, exists := s.users[new.ID()]; !exists {
		return nil, errors.New("cannot update guild none-existant user in cache")
	}

	old := s.users[new.ID()]

	old.Update(new)
	return old, nil
}

func (s *StateCache) DeleteUser(g user.UserInterface) {
	s.DeleteUserByID(g.ID())
}

func (s *StateCache) DeleteUserByID(ID snowflake.ID) {
	s.usersUpdateMutex.Lock()
	defer s.usersUpdateMutex.Unlock()
	if u, ok := s.users[ID]; ok {
		u.Clear()
		delete(s.users, ID) // TODO: how good is the golang garbage collector?
	}
}

func (s *StateCache) User(ID snowflake.ID) (user.UserInterface, error) {
	s.usersUpdateMutex.Lock()
	s.usersAddMutex.Lock()
	defer s.usersUpdateMutex.Unlock()
	defer s.usersAddMutex.Unlock()
	if g, ok := s.users[ID]; ok {
		return g, nil
	}

	return nil, errors.New("guild with ID{" + ID.String() + "} does not exist in cache")
}
