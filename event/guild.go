package event

import (
	"context"

	"github.com/andersfylling/disgord/resource"
	. "github.com/andersfylling/snowflake"
)

// KeyGuildCreate This event can be sent in three different scenarios:
//             1. When a user is initially connecting, to lazily load and
//                backfill information for all unavailable guilds sent in the
//                Ready event.
//             2. When a Guild becomes available again to the client.
//             3. When the current user joins a new Guild.
//             The inner payload is a guild object, with all the extra fields
//             specified.
const KeyGuildCreate = "GUILD_CREATE"

// GuildCreate	This event can be sent in three different scenarios:
//								1. When a user is initially connecting, to lazily load and backfill information for
//									 all unavailable guilds sent in the Ready event.
//								2. When a Guild becomes available again to the client.
// 								3. When the current user joins a new Guild.
type GuildCreate struct {
	Guild *resource.Guild `json:"guild"`
	Ctx   context.Context `json:"-"`
}

// KeyGuildUpdate Sent when a guild is updated. The inner payload is a guild
//             object.
const KeyGuildUpdate = "GUILD_UPDATE"

// GuildUpdate	guild was updated
type GuildUpdate struct {
	Guild *resource.Guild `json:"guild"`
	Ctx   context.Context `json:"-"`
}

// KeyGuildDelete Sent when a guild becomes unavailable during a guild outage,
//             or when the user leaves or is removed from a guild. The inner
//             payload is an unavailable guild object. If the unavailable
//             field is not set, the user was removed from the guild.
const KeyGuildDelete = "GUILD_DELETE"

// GuildDelete	guild became unavailable, or user left/was removed from a guild
type GuildDelete struct {
	UnavailableGuild *resource.GuildUnavailable `json:"guild_unavailable"`
	Ctx              context.Context            `json:"-"`
}

// KeyGuildBanAdd Sent when a user is banned from a guild. The inner payload is
//             a user object, with an extra guild_id key.
const KeyGuildBanAdd = "GUILD_BAN_ADD"

// GuildBanAdd	user was banned from a guild
type GuildBanAdd struct {
	User *resource.User  `json:"user"`
	Ctx  context.Context `json:"-"`
}

// KeyGuildBanRemove Sent when a user is unbanned from a guild. The inner
//                payload is a user object, with an extra guild_id key.
const KeyGuildBanRemove = "GUILD_BAN_REMOVE"

// GuildBanRemove	user was unbanned from a guild
type GuildBanRemove struct {
	User *resource.User  `json:"user"`
	Ctx  context.Context `json:"-"`
}

// KeyGuildIntegrationsUpdate Sent when a guild integration is updated.
//                        Fields:
//                        * GuildID int64 or discord.Snowflake
//                        * Emojis []*discord.emoji.Emoji
const KeyGuildIntegrationsUpdate = "GUILD_INTEGRATIONS_UPDATE"

// GuildIntegrationsUpdate	guild integration was updated
type GuildIntegrationsUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Ctx     context.Context `json:"-"`
}

// KeyGuildMemberAdd Sent when a new user joins a guild. The inner payload is a
//                guild member object with these extra fields:
//                * GuildID int64 or discord.Snowflake
const KeyGuildMemberAdd = "GUILD_MEMBER_ADD"

// GuildMemberAdd	new user joined a guild
type GuildMemberAdd struct {
	Member *resource.Member `json:"member"`
	Ctx    context.Context  `json:"-"`
}

// KeyGuildMemberRemove Sent when a user is removed from a guild
//                   (leave/kick/ban).
//                   Fields:
//                   * GuildID int64 or discord.Snowflake
//                   * User *discord.user.User
const KeyGuildMemberRemove = "GUILD_MEMBER_REMOVE"

// GuildMemberRemove	user was removed from a guild
type GuildMemberRemove struct {
	GuildID Snowflake       `json:"guild_id"`
	User    *resource.User  `json:"user"`
	Ctx     context.Context `json:"-"`
}

// KeyGuildMemberUpdate Sent when a guild member is updated.
//                   Fields:
//                   * GuildID int64 or discord.Snowflake
//                   * Roles []int64 or []discord.Snowflake
//                   * User *discord.user.User
//                   * Nick string
const KeyGuildMemberUpdate = "GUILD_MEMBER_UPDATE"

// GuildMemberUpdate	guild member was updated
type GuildMemberUpdate struct {
	GuildID Snowflake        `json:"guild_id"`
	Roles   []*resource.Role `json:"roles"`
	User    *resource.User   `json:"user"`
	Nick    string           `json:"nick"`
	Ctx     context.Context  `json:"-"`
}

// KeyGuildMemberChunk Sent in response to Gateway Request Guild Members.
//                  Fields:
//                  * GuildID int64 or discord.Snowflake
//                  * Members []*discord.member.Member
const KeyGuildMembersChunk = "GUILD_MEMBER_CHUNK"

// GuildMembersChunk	response to Request Guild Members
type GuildMembersChunk struct {
	GuildID Snowflake          `json:"guild_id"`
	Members []*resource.Member `json:"members"`
	Ctx     context.Context    `json:"-"`
}

// KeyGuildRoleCreate Sent when a guild role is created.
//                 Fields:
//                 * GuildID int64 or discord.Snowflake
//                 * Role *discord.role.Role
const KeyGuildRoleCreate = "GUILD_ROLE_CREATE"

// GuildRoleCreate	guild role was created
type GuildRoleCreate struct {
	GuildID Snowflake       `json:"guild_id"`
	Role    *resource.Role  `json:"role"`
	Ctx     context.Context `json:"-"`
}

// KeyGuildRoleUpdate Sent when a guild role is created.
//                 Fields:
//                 * GuildID int64 or discord.Snowflake
//                 * Role    *discord.role.Role
const KeyGuildRoleUpdate = "GUILD_ROLE_UPDATE"

// GuildRoleUpdate	guild role was updated
type GuildRoleUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Role    *resource.Role  `json:"role"`
	Ctx     context.Context `json:"-"`
}

// KeyGuildRoleDelete Sent when a guild role is created.
//                 Fields:
//                 * GuildID int64 or discord.Snowflake
//                 * RoleID  int64 or discord.Snowflake
const KeyGuildRoleDelete = "GUILD_ROLE_DELETE"

// GuildRoleDelete	guild role was deleted
type GuildRoleDelete struct {
	GuildID Snowflake       `json:"guild_id"`
	RoleID  Snowflake       `json:"role_id"`
	Ctx     context.Context `json:"-"`
}

// KeyPresenceUpdate A user's presence is their current state on a guild.
//                This event is sent when a user's presence is updated
//                for a guild.
//                Fields:
//                User    *discord.user.User
//                Roles   []*discord.role.Role
//                Game    *discord.game.Game
//                GuildID int64 or discord.Snowflake
//                Status  *string or *discord.presence.Status
const KeyPresenceUpdate = "PRESENCE_UPDATE"

// PresenceUpdate	user's presence was updated in a guild
type PresenceUpdate struct {
	User    *resource.User         `json:"user"`
	RoleIDs []Snowflake            `json:"roles"`
	Game    *resource.UserActivity `json:"game"`
	GuildID Snowflake              `json:"guild_id"`

	// Status either "idle", "dnd", "online", or "offline"
	// TODO: constants somewhere..
	Status string          `json:"status"`
	Ctx    context.Context `json:"-"`
}
