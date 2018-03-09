package state

import (
	"errors"

	"github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/snowflake"
)

// Cacher used by the session interface, so the user cannot access methods to ruin the cache state
type Cacher interface {
	User(ID snowflake.ID) (*resource.User, error)
	Channel(ID snowflake.ID) (*resource.Channel, error)
	Guild(ID snowflake.ID) (*resource.Guild, error)
	Myself() *resource.User

	// clear all the cached objects
	Clear()
}

// TODO: private interface for Disgord?
// TODO: guilds has copies of channels, that exists in the channels cacher
// TODO: channels has copies of users, that exists in the users cacher
// TODO: on caching, make sure only one of the objects exists in memory
// 		 channel.Users[0].ID == 1234 should point to the users cacher where user.ID == 1234.
// 		 Not hold a copy.
//		 Discordgo solves this by having arrays of string snowflakes for relevant objects
//		 But I feel this causes an extra step for the users, to extract a user object from say a guild/channel.

func NewCache() *Cache {
	st := &Cache{}
	st.Users = NewUserCache()
	st.Channels = NewChannelCache(st.Users)
	st.Guilds = NewGuildCache(st.Users, st.Channels)

	return st
}

type Cache struct {
	Users    UserCacher
	Channels ChannelCacher
	Guilds   GuildCacher // TODO: implement

	mySelf *resource.User
}

// -----
// Users
//
func (st *Cache) User(ID snowflake.ID) (*resource.User, error) {
	if st.Users == nil {
		return nil, errors.New("user caching has not been activated/implemented")
	}
	return st.Users.User(ID)
}

func (st *Cache) ProcessUser(details *UserDetail) {
	if st.Users != nil {
		st.Users.Process(details)
	}
}

// --------
// Channels
//

func (st *Cache) Channel(ID snowflake.ID) (*resource.Channel, error) {
	if st.Channels == nil {
		return nil, errors.New("channel caching has not been activated/implemented")
	}
	return st.Channels.Channel(ID)
}

func (st *Cache) ProcessChannel(details *ChannelDetail) {
	if st.Channels != nil {
		st.Channels.Process(details)
	}
}

// ------
// Guilds
//
func (st *Cache) Guild(ID snowflake.ID) (*resource.Guild, error) {
	if st.Guilds == nil {
		return nil, errors.New("guild caching has not been activated/implemented")
	}
	return st.Guilds.Guild(ID)
}

func (st *Cache) ProcessGuild(gd *GuildDetail) {
	if st.Guilds != nil {
		st.Guilds.Process(gd)
	}
}

// Channel listeners for object updates
//

// https://golang.org/pkg/io/#Closer
func (st *Cache) Close() (err error) {
	err = st.Users.Close()
	if err != nil {
		return
	}
	err = st.Channels.Close()
	if err != nil {
		return
	}
	err = st.Guilds.Close()
	if err != nil {
		return
	}

	return err
}

func (st *Cache) Clear() {
	st.Users.Clear()
	st.Channels.Clear()
	st.Guilds.Clear()
}

func (s *Cache) Myself() *resource.User {
	return s.mySelf
}
