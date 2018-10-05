package disgord

import (
	"github.com/andersfylling/disgord/cache/interfaces"
)

func createGuildCacher(conf *CacheConfig) (cacher interfaces.CacheAlger, err error) {
	if conf.DisableGuildCaching {
		return nil, nil
	}

	const channelWeight = 1 // MiB. TODO: what is the actual max size?
	limit := conf.ChannelCacheLimitMiB / channelWeight

	cacher, err = constructSpecificCacher(conf.ChannelCacheAlgorithm, limit, conf.ChannelCacheLifetime)
	return
}

type guildCacheItem struct {
	guild    *Guild
	channels []Snowflake
}

func (g *guildCacheItem) process(guild *Guild, immutable bool) {
	if immutable {
		g.guild = guild.DeepCopy().(*Guild)

		for _, member := range g.guild.Members {
			member.userID = member.User.ID
			member.User = nil
		}
	} else {
		g.guild = guild

		for _, member := range g.guild.Members {
			member.userID = member.User.ID
		}
	}

}

func (g *guildCacheItem) build(cache *Cache) (guild *Guild) {
	var err error

	if cache.immutable {
		guild = g.guild.DeepCopy().(*Guild)
		guild.Channels = make([]*Channel, len(g.channels))
		for i := range g.channels {
			guild.Channels[i], err = cache.GetChannel(g.channels[i])
			if err != nil {
				guild.Channels[i] = &Channel{
					ID: g.channels[i],
				}
			}
		}
		// TODO: voice state
	} else {
		guild = g.guild

		channels := make([]*Channel, len(g.channels))
		for i := range g.channels {
			channels[i], err = cache.GetChannel(g.channels[i])
			if err != nil {
				channels[i] = &Channel{
					ID: g.channels[i],
				}
			}
		}
		guild.Channels = channels
	}

	return
}

func (g *guildCacheItem) update(fresh *Guild, immutable bool) {
	if immutable {
		fresh.copyOverToCache(g.guild)
		// roles
		if len(fresh.Roles) > 0 {
			g.guild.Roles = make([]*Role, len(fresh.Roles))
		}
		for i, role := range fresh.Roles {
			if role == nil {
				continue
			}
			g.guild.Roles[i] = role.DeepCopy().(*Role)
		}
		// emojis
		if len(fresh.Emojis) > 0 {
			g.guild.Emojis = make([]*Emoji, len(fresh.Emojis))
		}
		for i, emoji := range fresh.Emojis {
			if emoji == nil {
				continue
			}
			g.guild.Emojis[i] = emoji.DeepCopy().(*Emoji)
		}
		// voice states
		if len(fresh.VoiceStates) > 0 {
			g.guild.VoiceStates = make([]*VoiceState, len(fresh.VoiceStates))
		}
		for i, state := range fresh.VoiceStates {
			if state == nil {
				continue
			}
			g.guild.VoiceStates[i] = state.DeepCopy().(*VoiceState)
		}
		// members
		if len(fresh.Members) > 0 {
			g.guild.Members = make([]*Member, len(fresh.Members))
		}
		for i, m := range fresh.Members {
			if m == nil {
				continue
			}
			m.userID = m.User.ID
			m.User = nil
			g.guild.Members[i] = m.DeepCopy().(*Member)
		}
		// presences
		if len(fresh.Presences) > 0 {
			g.guild.Presences = make([]*UserPresence, len(fresh.Presences))
		}
		for i, p := range fresh.Presences {
			if p == nil {
				continue
			}
			g.guild.Presences[i] = p.DeepCopy().(*UserPresence)
		}
		// channels
		if len(fresh.Channels) > 0 {
			g.channels = make([]Snowflake, len(fresh.Channels))
		}
		for i, c := range fresh.Channels {
			if c == nil {
				continue
			}
			g.channels[i] = c.ID
		}
	} else {
		if len(fresh.Roles) == 0 && len(g.guild.Roles) > 0 {
			fresh.Roles = g.guild.Roles
		}
		if len(fresh.Emojis) == 0 && len(g.guild.Emojis) > 0 {
			fresh.Emojis = g.guild.Emojis
		}
		if len(fresh.VoiceStates) == 0 && len(g.guild.VoiceStates) > 0 {
			fresh.VoiceStates = g.guild.VoiceStates
		}
		if len(fresh.Members) == 0 && len(g.guild.Members) > 0 {
			fresh.Members = g.guild.Members
		}
		if len(fresh.Channels) == 0 && len(g.guild.Channels) > 0 {
			fresh.Channels = g.guild.Channels
		}
		if len(fresh.Presences) == 0 && len(g.guild.Presences) > 0 {
			fresh.Presences = g.guild.Presences
		}
		g.guild = fresh
	}
}

func (g *guildCacheItem) deleteChannel(id Snowflake) {
	for i := range g.channels {
		if g.channels[i] != id {
			continue
		}

		g.channels[i] = g.channels[len(g.channels)-1]
		g.channels = g.channels[:len(g.channels)-1]
		g.guild.DeleteChannelByID(id)
		return
	}
}

func (c *Cache) SetGuild(guild *Guild) {
	if c.guilds == nil || guild == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()
	if item, exists := c.guilds.Get(guild.ID); exists {
		item.Object().(*guildCacheItem).update(guild, c.immutable)
		c.guilds.RefreshAfterDiscordUpdate(item)
	} else {
		content := &guildCacheItem{}
		content.process(guild, c.immutable)
		c.guilds.Set(guild.ID, c.guilds.CreateCacheableItem(content))
	}
}

func (c *Cache) SetGuildEmojis(guildID Snowflake, emojis []*Emoji) {
	if c.guilds == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()
	if item, exists := c.guilds.Get(guildID); exists {
		guild := item.Object().(*guildCacheItem).guild
		if c.immutable {
			emojisCopy := make([]*Emoji, len(emojis))
			for i := range emojis {
				emojisCopy[i] = emojis[i].DeepCopy().(*Emoji)
			}
			// TODO-racecondition
			guild.Emojis = emojisCopy
		} else {
			// code smell => try to update only affected emoji's
			guild.Emojis = emojis
		}
		c.guilds.RefreshAfterDiscordUpdate(item)
	} else {
		content := &guildCacheItem{}
		content.process(&Guild{
			ID:     guildID,
			Emojis: emojis,
		}, c.immutable)
		c.guilds.Set(guildID, c.guilds.CreateCacheableItem(content))
	}
}

func (c *Cache) GetGuild(id Snowflake) (guild *Guild, err error) {
	if c.guilds == nil {
		err = NewErrorUsingDeactivatedCache("guilds")
		return
	}

	c.guilds.RLock()
	defer c.guilds.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.guilds.Get(id); !exists {
		err = NewErrorCacheItemNotFound(id)
		return
	}

	guild = result.Object().(*guildCacheItem).build(c)
	return
}

func (c *Cache) DeleteGuild(id Snowflake) {
	if c.guilds == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()

	c.guilds.Delete(id)
}

func (c *Cache) DeleteGuildChannel(guildID, channelID Snowflake) {
	if c.guilds == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()
	if item, exists := c.guilds.Get(guildID); exists {
		item.Object().(*guildCacheItem).deleteChannel(channelID)
		c.guilds.RefreshAfterDiscordUpdate(item)
	}
}

func (c *Cache) DeleteGuildRole(guildID, roleID Snowflake) {
	if c.guilds == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()
	if item, exists := c.guilds.Get(guildID); exists {
		item.Object().(*guildCacheItem).guild.DeleteRoleByID(roleID)
		c.guilds.RefreshAfterDiscordUpdate(item)
	}
}
