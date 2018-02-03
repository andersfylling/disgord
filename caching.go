package disgord

import (
	"os/user"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/guild"
	"github.com/andersfylling/snowflake"
)

type Stater interface {
	AddGuild(g *guild.Guild)
	UpdateGuild(g *guild.Guild)
	DeleteGuild(g *guild.Guild)
	DeleteGuildByID(ID snowflake.ID)

	AddChannel(c *channel.Channel)
	UpdateChannel(c *channel.Channel)
	DeleteChannel(c *channel.Channel)
	DeleteChannelByID(ID snowflake.ID)
}

type StateCache struct {
	guilds   map[snowflake.ID]*guild.Guild
	users    map[snowflake.ID]*user.User
	channels map[snowflake.ID]*channel.Channel
}

// guilds
//

func (s *StateCache) AddGuild(g *guild.Guild) {
	s.guilds[g.ID] = g
}

func (s *StateCache) UpdateGuild(g *guild.Guild) {
	s.guilds[g.ID] = g
}

func (s *StateCache) DeleteGuild(g *guild.Guild) {
	s.DeleteGuildByID(g.ID)
}

func (s *StateCache) DeleteGuildByID(ID snowflake.ID) {
	if _, ok := s.guilds[ID]; ok {
		delete(s.guilds, ID)
	}
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
