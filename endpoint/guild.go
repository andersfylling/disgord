package endpoint

import "fmt"

// Guilds /guilds
func Guilds() string {
	return guilds
}

// Guild /guilds/{guild.id}
func Guild(id fmt.Stringer) string {
	return guilds + "/" + id.String()
}

// GuildChannels /guilds/{guild.id}/channels
func GuildChannels(id fmt.Stringer) string {
	return Guild(id) + channels
}

// GuildChannel /guilds/{guild.id}/channels/{channel.id}
func GuildChannel(guildID, channelID fmt.Stringer) string {
	return Guild(guildID) + Channel(channelID)
}

// GuildMembers /guilds/{guild.id}/members
func GuildMembers(id fmt.Stringer) string {
	return Guild(id) + members
}

// GuildMember /guilds/{guild.id}/members/{user.id}
func GuildMember(guildID, userID fmt.Stringer) string {
	return GuildMembers(guildID) + "/" + userID.String()
}

// GuildMembersMeNick /guilds/{guild.id}/members/@me/nick
func GuildMembersMeNick(guildID fmt.Stringer) string {
	return GuildMembers(guildID) + me + nick
}

// GuildMemberRole /guilds/{guild.id}/members/{user.id}/roles/{role.id}
func GuildMemberRole(guildID, userID, roleID fmt.Stringer) string {
	return GuildMember(guildID, userID) + roles + "/" + roleID.String()
}

// GuildBans /guilds/{guild.id}/bans
func GuildBans(id fmt.Stringer) string {
	return Guild(id) + bans
}

// GuildBan /guilds/{guild.id}/bans/{user.id}
func GuildBan(guildID, userID fmt.Stringer) string {
	return Guild(guildID) + bans + "/" + userID.String()
}

// GuildRoles /guilds/{guild.id}/roles
func GuildRoles(id fmt.Stringer) string {
	return Guild(id) + roles
}

// GuildRole /guilds/{guild.id}/roles/{role.id}
func GuildRole(guildID, roleID fmt.Stringer) string {
	return GuildRoles(guildID) + "/" + roleID.String()
}

// GuildPrune /guilds/{guild.id}/prune
func GuildPrune(id fmt.Stringer) string {
	return Guild(id) + prune
}

// GuildRegions /guilds/{guild.id}/regions
func GuildRegions(id fmt.Stringer) string {
	return Guild(id) + regions
}

// GuildInvites /guilds/{guild.id}/invites
func GuildInvites(id fmt.Stringer) string {
	return Guild(id) + invites
}

// GuildIntegrations /guilds/{guild.id}/integrations
func GuildIntegrations(id fmt.Stringer) string {
	return Guild(id) + integrations
}

// GuildIntegration /guilds/{guild.id}/integrations/{integration.id}
func GuildIntegration(guildID, integrationID fmt.Stringer) string {
	return GuildIntegrations(guildID) + "/" + integrationID.String()
}

// GuildIntegrationSync /guilds/{guild.id}/integrations/{integration.id}/sync
func GuildIntegrationSync(guildID, integrationID fmt.Stringer) string {
	return GuildIntegration(guildID, integrationID) + sync
}

// GuildEmbed /guilds/{guild.id}/embed
func GuildEmbed(id fmt.Stringer) string {
	return Guild(id) + embed
}

// GuildVanityURL /guilds/{guild.id}/vanity-url
func GuildVanityURL(id fmt.Stringer) string {
	return Guild(id) + vanityURL
}
