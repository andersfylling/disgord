package disgord

// handlerGuildDelete update internal state when joining or creating a guild
func (c *Client) handlerAddToConnectedGuilds(s Session, evt *GuildCreate) {
	c.connectedGuildsMutex.Lock()
	defer c.connectedGuildsMutex.Unlock()

	// don't add an entry if there already is one
	for i := range c.connectedGuilds {
		if c.connectedGuilds[i] == evt.Guild.ID {
			return
		}
	}
	c.connectedGuilds = append(c.connectedGuilds, evt.Guild.ID)
}

// handlerGuildDelete update internal state when deleting or leaving a guild
func (c *Client) handlerRemoveFromConnectedGuilds(s Session, evt *GuildDelete) {
	c.connectedGuildsMutex.Lock()
	for i := range c.connectedGuilds {
		if c.connectedGuilds[i] != evt.UnavailableGuild.ID {
			continue
		}
		c.connectedGuilds[i] = c.connectedGuilds[len(c.connectedGuilds)-1]
		c.connectedGuilds = c.connectedGuilds[:len(c.connectedGuilds)-1]
		break
	}
	c.connectedGuildsMutex.Unlock()
}

func (c *Client) handlerSetSelfBotID(session Session, rdy *Ready) {
	c.myID = rdy.User.ID
}
func (c *Client) handlerUpdateSelfBot(session Session, update *UserUpdate) {
	_ = session.Cache().Update(UserCache, update.User)
}
