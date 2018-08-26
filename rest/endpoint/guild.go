package endpoint

import . "github.com/andersfylling/snowflake"

// Guilds /guilds
func Guilds() string {
	return guilds
}

// Guild /guilds/{guild.id}
func Guild(id Snowflake) string {
	return guilds + "/" + id.String()
}

// Guild /guilds/{guild.id}/channels
func GuildChannels(id Snowflake) string {
	return Guild(id) + channels
}

// Guild /guilds/{guild.id}/channels/{channel.id}
func GuildChannel(guildID, channelID Snowflake) string {
	return Guild(guildID) + Channel(channelID)
}

// Guild /guilds/{guild.id}/members
func GuildMembers(id Snowflake) string {
	return Guild(id) + members
}

// Guild /guilds/{guild.id}/members/{user.id}
func GuildMember(guildID, userID Snowflake) string {
	return GuildMembers(guildID) + "/" + userID.String()
}

// GuildMemberMeNick /guilds/{guild.id}/members/@me/nick
func GuildMembersMeNick(guildID Snowflake) string {
	return GuildMembers(guildID) + me + nick
}

// GuildMemberRole /guilds/{guild.id}/members/{user.id}/roles/{role.id}
func GuildMemberRole(guildID, userID, roleID Snowflake) string {
	return GuildMember(guildID, userID) + roles + "/" +roleID.String()
}

// GuildBans /builds/{guild.id}/bans
func GuildBans(id Snowflake) string {
	return Guild(id) + bans
}

// GuildBans /builds/{guild.id}/bans/{user.id}
func GuildBan(guildID, userID Snowflake) string {
	return Guild(guildID) + bans + "/" + userID.String()
}

// GuildRoles /guilds/{guild.id}/roles
func GuildRoles(id Snowflake) string {
	return Guild(id) + roles
}

// GuildRole /guilds/{guild.id}/roles/{role.id}
func GuildRole(guildID, roleID Snowflake) string {
	return GuildRoles(guildID) + "/" + roleID.String()
}

// GuildPrune /guilds/{guild.id}/prune
func GuildPrune(id Snowflake) string {
	return Guild(id) + prune
}

// GuildRegions /guilds/{guild.id}/regions
func GuildRegions(id Snowflake) string {
	return Guild(id) + regions
}

// GuildInvites /guilds/{guild.id}/invites
func GuildInvites(id Snowflake) string {
	return Guild(id) + invites
}

// GuildIntegrations /guilds/{guild.id}/integrations
func GuildIntegrations(id Snowflake) string {
	return Guild(id) + integrations
}

// GuildIntegration /guilds/{guild.id}/integrations/{integration.id}
func GuildIntegration(guildID, integrationID Snowflake) string {
	return GuildIntegrations(guildID) + "/" + integrationID.String()
}

// GuildIntegration /guilds/{guild.id}/integrations/{integration.id}/sync
func GuildIntegrationSync(guildID, integrationID Snowflake) string {
	return GuildIntegration(guildID, integrationID) + sync
}

// GuildEmbed /guilds/{guild.id}/embed
func GuildEmbed(id Snowflake) string {
	return Guild(id) + embed
}

// GuildVanityURL /guilds/{guild.id}/vanity-url
func GuildVanityURL(id Snowflake) string {
	return Guild(id) + vanityURL
}