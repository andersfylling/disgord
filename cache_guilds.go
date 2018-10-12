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

func (g *guildCacheItem) updateMembers(members []*Member, immutable bool) {
	newMembers := []*Member{}

	g.guild.Lock()
	defer g.guild.Unlock()

	var userID Snowflake
	for i := range members {
		userID = members[i].User.ID
		for j := range g.guild.Members {
			if g.guild.Members[j].userID == userID {
				userID = 0
				*g.guild.Members[j] = *members[i]
				g.guild.Members[j].User = nil
				break
			}
		}

		if !userID.Empty() {
			newMembers = append(newMembers, members[i])
		}
	}

	var member *Member
	for i := range newMembers {
		member = newMembers[i]
		member.userID = member.User.ID
		member.User = nil
		g.guild.Members = append(g.guild.Members, member)
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

// SetGuild adds a new guild to cache or updates an existing one
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

// SetGuildEmojis adds a new guild to cache if no guild exist for the emojis or updates an existing guild with the new emojis
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

// SetGuildMember calls SetGuildMembers
func (c *Cache) SetGuildMember(guildID Snowflake, member *Member) {
	if c.guilds == nil || member == nil {
		return
	}

	c.SetGuildMembers(guildID, []*Member{member})
}

// SetGuildMembers adds the members to a guild or updates an existing guild
func (c *Cache) SetGuildMembers(guildID Snowflake, members []*Member) {
	if c.guilds == nil || members == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()
	if item, exists := c.guilds.Get(guildID); exists {
		item.Object().(*guildCacheItem).updateMembers(members, c.immutable)
		c.guilds.RefreshAfterDiscordUpdate(item)
	} else {
		content := &guildCacheItem{}
		content.process(&Guild{
			ID:      guildID,
			Members: members,
		}, c.immutable)
		c.guilds.Set(guildID, c.guilds.CreateCacheableItem(content))
	}
}

// SetGuildRoles creates a new guild if none is found and updates the roles for a given guild
func (c *Cache) SetGuildRoles(guildID Snowflake, roles []*Role) {
	if c.guilds == nil || roles == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()
	if item, exists := c.guilds.Get(guildID); exists {
		guild := item.Object().(*guildCacheItem).guild
		var newRoles []*Role
		if c.immutable {
			newRoles = make([]*Role, len(roles))
			for i := range roles {
				newRoles[i] = roles[i].DeepCopy().(*Role)
			}
		} else {
			newRoles = roles
		}
		guild.Roles = newRoles
		c.guilds.RefreshAfterDiscordUpdate(item)
	} else {
		content := &guildCacheItem{}
		content.process(&Guild{
			ID:    guildID,
			Roles: roles,
		}, c.immutable)
		c.guilds.Set(guildID, c.guilds.CreateCacheableItem(content))
	}
}

// GetGuild ...
func (c *Cache) GetGuild(id Snowflake) (guild *Guild, err error) {
	if c.guilds == nil {
		err = newErrorUsingDeactivatedCache("guilds")
		return
	}

	c.guilds.RLock()
	defer c.guilds.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.guilds.Get(id); !exists {
		err = newErrorCacheItemNotFound(id)
		return
	}

	guild = result.Object().(*guildCacheItem).build(c)
	return
}

// GetGuildRoles ...
func (c *Cache) GetGuildRoles(id Snowflake) (roles []*Role, err error) {
	if c.guilds == nil {
		err = newErrorUsingDeactivatedCache("guilds")
		return
	}

	c.guilds.RLock()
	defer c.guilds.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.guilds.Get(id); !exists {
		err = newErrorCacheItemNotFound(id)
		return
	}

	rolePs := result.Object().(*guildCacheItem).guild.Roles
	if c.immutable {
		roles = make([]*Role, len(rolePs))
		for i := range rolePs {
			roles[i] = rolePs[i].DeepCopy().(*Role)
		}
	} else {
		roles = rolePs
	}

	return
}

// GetGuildMember ...
func (c *Cache) GetGuildMember(guildID, userID Snowflake) (member *Member, err error) {
	if c.guilds == nil {
		err = newErrorUsingDeactivatedCache("guilds")
		return
	}

	c.guilds.RLock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.guilds.Get(guildID); !exists {
		err = newErrorCacheItemNotFound(guildID)
		return
	}

	guild := result.Object().(*guildCacheItem).guild
	for i := range guild.Members {
		if guild.Members[i].userID == userID {
			member = guild.Members[i]
			if c.immutable {
				member = member.DeepCopy().(*Member)
			}
			break
		}
	}
	c.guilds.RUnlock()

	if member == nil {
		err = newErrorCacheItemNotFound(userID)
		return
	}

	// add user object
	member.User, err = c.GetUser(userID)
	if err != nil {
		member.User = &User{
			ID: userID,
		}
	}
	return
}

// GetGuildMembersAfter ...
func (c *Cache) GetGuildMembersAfter(guildID, after Snowflake, limit int) (members []*Member, err error) {
	if c.guilds == nil {
		err = newErrorUsingDeactivatedCache("guilds")
		return
	}

	c.guilds.RLock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.guilds.Get(guildID); !exists {
		err = newErrorCacheItemNotFound(guildID)
		return
	}

	guild := result.Object().(*guildCacheItem).guild
	for i := range guild.Members {
		if guild.Members[i].userID > after && len(members) <= limit {
			member := guild.Members[i]
			if c.immutable {
				member = member.DeepCopy().(*Member)
			}
			members = append(members, member)
		} else if len(members) > limit {
			break
		}
	}
	c.guilds.RUnlock()

	for i := range members {
		// add user object
		members[i].User, err = c.GetUser(members[i].userID)
		if err != nil {
			members[i].User = &User{
				ID: members[i].userID,
			}
		}
	}
	return
}

// DeleteGuild ...
func (c *Cache) DeleteGuild(id Snowflake) {
	if c.guilds == nil {
		return
	}

	c.guilds.Lock()
	defer c.guilds.Unlock()

	c.guilds.Delete(id)
}

// DeleteGuildChannel removes a channel from a cached guild object without removing the guild
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

// DeleteGuildRole removes a role from a cached guild object without removing the guild
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
