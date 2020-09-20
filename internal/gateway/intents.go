package gateway

type Intent uint64

const (
	// IntentGuilds
	// - GUILD_CREATE
	// - GUILD_UPDATE
	// - GUILD_DELETE
	// - GUILD_ROLE_CREATE
	// - GUILD_ROLE_UPDATE
	// - GUILD_ROLE_DELETE
	// - CHANNEL_CREATE
	// - CHANNEL_UPDATE
	// - CHANNEL_DELETE
	// - CHANNEL_PINS_UPDATE
	IntentGuilds Intent = 1 << iota

	// IntentGuildMembers
	// - GUILD_MEMBER_ADD
	// - GUILD_MEMBER_UPDATE
	// - GUILD_MEMBER_REMOVE
	IntentGuildMembers

	// IntentGuildBans
	// - GUILD_BAN_ADD
	// - GUILD_BAN_REMOVE
	IntentGuildBans

	// IntentGuildEmojis
	// - GUILD_EMOJIS_UPDATE
	IntentGuildEmojis

	// IntentGuildIntegrations
	// - GUILD_INTEGRATIONS_UPDATE
	IntentGuildIntegrations

	IntentGuildWebhooks
	IntentGuildInvites
	IntentGuildVoiceStates
	IntentGuildPresences
	IntentGuildMessages
	IntentGuildMessageReactions
	IntentGuildMessageTyping
	IntentDirectMessages
	IntentDirectMessageReactions
	IntentDirectMessageTyping
)
