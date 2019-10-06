package endpoint

import "fmt"

// UserMe /users/@me
func UserMe() string {
	return users + me
}

// User ...
func User(id fmt.Stringer) string {
	return users + "/" + id.String()
}

// UserMeGuilds ...
func UserMeGuilds() string {
	return UserMe() + guilds
}

// UserMeGuild ...
func UserMeGuild(id fmt.Stringer) string {
	return UserMe() + guilds + "/" + id.String()
}

// UserMeChannels ...
func UserMeChannels() string {
	return UserMe() + channels
}

// UserMeChannel ...
func UserMeChannel(id fmt.Stringer) string {
	return UserMe() + channels + "/" + id.String()
}

// UserMeConnections ...
func UserMeConnections() string {
	return UserMe() + connections
}
