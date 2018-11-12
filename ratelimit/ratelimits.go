package ratelimit

import (
	"github.com/andersfylling/snowflake/v3"
)

func GetB() *b {
	return &b{}
}

type b struct{}

func (a *b) Test() bool {
	return true
}

// endpoints/paths
const (
	discordAPI   = "https://discordapp.com/api"
	auditlogs    = "/audit-logs"
	channels     = "/channels"
	messages     = "/messages"
	bulkDelete   = "/bulk-delete"
	recipients   = "/recipients"
	pins         = "/pins"
	typing       = "/typing"
	permissions  = "/permissions"
	invites      = "/invites"
	reactions    = "/reactions"
	me           = "/@me"
	emojis       = "/emojis"
	guilds       = "/guilds"
	users        = "/users"
	connections  = "/connections"
	voice        = "/voice"
	regions      = "/regions"
	webhooks     = "/webhooks"
	slack        = "/slack"
	github       = "/github"
	members      = "/members"
	nick         = "/nick"
	roles        = "/roles"
	bans         = "/bans"
	prune        = "/prune"
	integrations = "/integrations"
	sync         = "/sync"
	embed        = "/embed"
	vanityURL    = "/vanity-url"
	gateway      = "/gateway"
	version      = "/v"
)

// --------------------
// Audit Log

func GuildAuditLogs(id snowflake.ID) string {
	return Guild(id) + ":a-l"
}

// --------------------
// Guild

func Guild(id snowflake.ID) string {
	return "g:" + id.String()
}
func GuildEmbed(id snowflake.ID) string {
	return Guild(id) + ":e"
}
func GuildVanityURL(id snowflake.ID) string {
	return Guild(id) + ":vurl"
}
func GuildChannels(id snowflake.ID) string {
	return Guild(id) + ":c"
}
func GuildMembers(id snowflake.ID) string {
	return Guild(id) + ":m"
}
func GuildBans(id snowflake.ID) string {
	return Guild(id) + ":b"
}
func GuildRoles(id snowflake.ID) string {
	return Guild(id) + ":r"
}
func GuildRegions(id snowflake.ID) string {
	return Guild(id) + ":regions"
}
func GuildIntegrations(id snowflake.ID) string {
	return Guild(id) + ":i"
}
func GuildInvites(id snowflake.ID) string {
	return Guild(id) + ":inv"
}
func GuildPrune(id snowflake.ID) string {
	return Guild(id) + ":p"
}
func GuildWebhooks(id snowflake.ID) string {
	return Guild(id) + ":w"
}

// --------------------
// Guild

// --------------------
// Guild

// --------------------
// Guild

// --------------------
// Guild

// --------------------
// Guild

// --------------------
// Guild

// --------------------
// Guild

// --------------------
// Guild

// --------------------
// Invite

// Invites /invites
func Invites() string {
	return invites
}

// --------------------
// Voice

// VoiceRegions /voice/regions
func VoiceRegions() string {
	return voice + regions
}

// --------------------
// Guild
