package endpoint

import . "github.com/andersfylling/snowflake"

// UserMe /users/@me
func UserMe() string {
	return users + me
}

func User(id Snowflake) string {
	return users + "/" + id.String()
}

func UserMeGuilds() string {
	return UserMe() + guilds
}
func UserMeGuild(id Snowflake) string {
	return UserMe() + guilds
}

func UserMeChannels() string {
	return UserMe() + channels
}

func UserMeChannel(id Snowflake) string {
	return UserMe() + channels
}

func UserMeConnections() string {
	return UserMe() + connections
}
