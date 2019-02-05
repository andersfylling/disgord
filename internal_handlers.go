package disgord

import "fmt"

// handlerGuildDelete update internal state when joining or creating a guild
func (c *Client) handlerAddToConnectedGuilds(s Session, evt *GuildCreate) {
	// NOTE: during unit tests, you must remember that shards are usually added dynamically at runtime
	//  meaning, you might have to add your own shards if you get a panic here
	shard, _ := c.shardManager.GetShard(evt.Guild.ID)
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
	// NOTE: during unit tests, you must remember that shards are usually added dynamically at runtime
	//  meaning, you might have to add your own shards if you get a panic here
	shard, _ := c.shardManager.GetShard(evt.UnavailableGuild.ID)
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
	fmt.Println("Got READY in selfbot setter")
	c.myID = rdy.User.ID
}
func (c *Client) handlerUpdateSelfBot(session Session, update *UserUpdate) {
	_ = session.Cache().Update(UserCache, update.User)
}
