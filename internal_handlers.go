package disgord

// handlerGuildDelete update internal state when joining or creating a guild
func (c *Client) handlerAddToConnectedGuilds(s Session, evt *GuildCreate) {
	var shard *WSShard
	if shard = c.shardManager.GetShard(evt.Guild.ID); shard == nil {
		// helps with writing unit tests
		// TODO: remove
		c.logErr("got a guild event from a unknown shard. Please notify the devs immediately")
		return
	}

	shard.Lock()
	defer shard.Unlock()

	// don't add an entry if there already is one
	for i := range shard.guilds {
		if shard.guilds[i] == evt.Guild.ID {
			return
		}
	}
	shard.guilds = append(shard.guilds, evt.Guild.ID)
}

// handlerGuildDelete update internal state when deleting or leaving a guild
func (c *Client) handlerRemoveFromConnectedGuilds(s Session, evt *GuildDelete) {
	var shard *WSShard
	if shard = c.shardManager.GetShard(evt.UnavailableGuild.ID); shard == nil {
		// helps with writing unit tests
		// TODO: remove
		c.logErr("got a guild event from a unknown shard. Please notify the devs immediately")
		return
	}

	shard.Lock()
	defer shard.Unlock()

	for i := range shard.guilds {
		if shard.guilds[i] != evt.UnavailableGuild.ID {
			continue
		}
		shard.guilds[i] = shard.guilds[len(shard.guilds)-1]
		shard.guilds = shard.guilds[:len(shard.guilds)-1]
	}
}

func (c *Client) handlerSetSelfBotID(session Session, rdy *Ready) {
	c.myID = rdy.User.ID
}
func (c *Client) handlerUpdateSelfBot(session Session, update *UserUpdate) {
	_ = session.Cache().Update(UserCache, update.User)
}
