package endpoint

import "fmt"

// UserMe /users/@me
func UserMe() string {
	return users + me
}

func User(id fmt.Stringer) string {
	return users + "/" + id.String()
}

func UserMeGuilds() string {
	return UserMe() + guilds
}
func UserMeGuild(id fmt.Stringer) string {
	return UserMe() + guilds
}

func UserMeChannels() string {
	return UserMe() + channels
}

func UserMeChannel(id fmt.Stringer) string {
	return UserMe() + channels
}

func UserMeConnections() string {
	return UserMe() + connections
}
