package gateway

import (
	"github.com/andersfylling/disgord/internal/event"
)

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

func EventToIntent(evt string, direct bool) Intent {
	var intent Intent

	if direct {
		switch evt {
		case event.MessageCreate:
			intent = IntentDirectMessages
		case event.MessageUpdate:
			intent = IntentDirectMessages
		case event.MessageDelete:
			intent = IntentDirectMessages
		case event.MessageDeleteBulk:
			intent = IntentDirectMessages
		case event.MessageReactionAdd:
			intent = IntentDirectMessageReactions
		case event.MessageReactionRemove:
			intent = IntentDirectMessageReactions
		case event.MessageReactionRemoveAll:
			intent = IntentDirectMessageReactions
		// case event.MessageReactionRemoveEmoji:
		// 	intent = IntentDirectMessageReactions
		case event.TypingStart:
			intent = IntentDirectMessageTyping
		}
	} else {
		switch evt {
		case event.GuildCreate:
			intent = IntentGuilds
		case event.GuildUpdate:
			intent = IntentGuilds
		case event.GuildDelete:
			intent = IntentGuilds
		case event.GuildRoleCreate:
			intent = IntentGuilds
		case event.GuildRoleUpdate:
			intent = IntentGuilds
		case event.GuildRoleDelete:
			intent = IntentGuilds
		case event.ChannelCreate:
			intent = IntentGuilds
		case event.ChannelUpdate:
			intent = IntentGuilds
		case event.ChannelDelete:
			intent = IntentGuilds
		case event.ChannelPinsUpdate:
			intent = IntentGuilds
		case event.GuildMemberAdd:
			intent = IntentGuildMembers
		case event.GuildMemberUpdate:
			intent = IntentGuildMembers
		case event.GuildMemberRemove:
			intent = IntentGuildMembers
		case event.GuildBanAdd:
			intent = IntentGuildBans
		case event.GuildBanRemove:
			intent = IntentGuildBans
		case event.GuildEmojisUpdate:
			intent = IntentGuildEmojis
		case event.GuildIntegrationsUpdate:
			intent = IntentGuildIntegrations
		case event.WebhooksUpdate:
			intent = IntentGuildWebhooks
		case event.InviteCreate:
			intent = IntentGuildInvites
		case event.InviteDelete:
			intent = IntentGuildInvites
		case event.VoiceStateUpdate:
			intent = IntentGuildVoiceStates
		case event.PresenceUpdate:
			intent = IntentGuildPresences
		case event.MessageCreate:
			intent = IntentGuildMessages
		case event.MessageUpdate:
			intent = IntentGuildMessages
		case event.MessageDelete:
			intent = IntentGuildMessages
		case event.MessageDeleteBulk:
			intent = IntentGuildMessages
		case event.MessageReactionAdd:
			intent = IntentGuildMessageReactions
		case event.MessageReactionRemove:
			intent = IntentGuildMessageReactions
		case event.MessageReactionRemoveAll:
			intent = IntentGuildMessageReactions
		// case event.MessageReactionRemoveEmoji:
		// 	intent = IntentGuildMessageReactions
		case event.TypingStart:
			intent = IntentGuildMessageTyping
		}
	}

	return intent
}
