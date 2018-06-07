package resource

import (
	"strconv"

	"github.com/andersfylling/disgord/request"
	"github.com/andersfylling/snowflake"
)

type AuditLog struct {
	Webhooks        []*Webhook       `json:"webhooks"`
	Users           []*User          `json:"users"`
	AuditLogEntries []*AuditLogEntry `json:"audit_log_entries"`
}

type AuditLogEntry struct {
	TargetID   snowflake.ID      `json:"target_id"`
	Changes    []*AuditLogChange `json:"changes,omitempty"`
	UserID     snowflake.ID      `json:"user_id"`
	ID         snowflake.ID      `json:"id"`
	ActionType uint              `json:"action_type"`
	Options    []*AuditLogOption `json:"options,omitempty"`
	Reason     string            `json:"reason,omitempty"`
}

const (
	AuditLogEvtGuildUpdate      = 1
	AuditLogEvtChannelCreate    = 10
	AuditLogEvtChannelUpdate    = 11
	AuditLogEvtChannelDelete    = 12
	AuditLogEvtOverwriteCreate  = 13
	AuditLogEvtOverwriteUpdate  = 14
	AuditLogEvtOverwriteDelete  = 15
	AuditLogEvtMemberKick       = 20
	AuditLogEvtMemberPrune      = 21
	AuditLogEvtMemberBanAdd     = 22
	AuditLogEvtMemberBanRemove  = 23
	AuditLogEvtMemberUpdate     = 24
	AuditLogEvtMemberRoleUpdate = 25
	AuditLogEvtRoleCreate       = 30
	AuditLogEvtRoleUpdate       = 31
	AuditLogEvtRoleDelete       = 32
	AuditLogEvtInviteCreate     = 40
	AuditLogEvtInviteUpdate     = 41
	AuditLogEvtInviteDelete     = 42
	AuditLogEvtWebhookCreate    = 50
	AuditLogEvtWebhookUpdate    = 51
	AuditLogEvtWebhookDelete    = 52
	AuditLogEvtEmojiCreate      = 60
	AuditLogEvtEmojiUpdate      = 61
	AuditLogEvtEmojiDelete      = 62
	AuditLogEvtMessageDelete    = 72
)

type AuditLogOption struct {
	DeleteMemberDays string       `json:"delete_member_days"`
	MembersRemoved   string       `json:"members_removed"`
	ChannelID        snowflake.ID `json:"channel_id"`
	Count            string       `json:"count"`
	ID               snowflake.ID `json:"id"`
	Type             string       `json:"type"` // type of overwritten entity ("member" or "role")
	RoleName         string       `json:"role_name"`
}

type AuditLogChange struct {
	NewValue interface{} `json:"new_value,omitempty"`
	OldValue interface{} `json:"old_value,omitempty"`
	Key      string      `json:"key"`
}

const (
	// key name,								          identifier                       changed, type,   description
	AuditLogChangeKeyName                        = "name"                          // guild	string	name changed
	AuditLogChangeKeyIconHash                    = "icon_hash"                     // guild	string	icon changed
	AuditLogChangeKeySplashHash                  = "splash_hash"                   // guild	string	invite splash page artwork changed
	AuditLogChangeKeyOwnerID                     = "owner_id"                      // guild	snowflake	owner changed
	AuditLogChangeKeyRegion                      = "region"                        // guild	string	region changed
	AuditLogChangeKeyAFKChannelID                = "afk_channel_id"                // guild	snowflake	afk channel changed
	AuditLogChangeKeyAFKTimeout                  = "afk_timeout"                   // guild	integer	afk timeout duration changed
	AuditLogChangeKeyMFALevel                    = "mfa_level"                     // guild	integer	two-factor auth requirement changed
	AuditLogChangeKeyVerificationLevel           = "verification_level"            // guild	integer	required verification level changed
	AuditLogChangeKeyExplicitContentFilter       = "explicit_content_filter"       // guild	integer	change in whose messages are scanned and deleted for explicit content in the server
	AuditLogChangeKeyDefaultMessageNotifications = "default_message_notifications" // guild	integer	default message notification level changed
	AuditLogChangeKeyVanityURLCode               = "vanity_url_code"               // guild	string	guild invite vanity url changed
	AuditLogChangeKeyAdd                         = "$add"                          // add	guild	array of role objects	new role added
	AuditLogChangeKeyRemove                      = "$remove"                       // remove	guild	array of role objects	role removed
	AuditLogChangeKeyPruneDeleteDays             = "prune_delete_days"             // guild	integer	change in number of days after which inactive and role-unassigned members are kicked
	AuditLogChangeKeyWidgetEnabled               = "widget_enabled"                // guild	bool	server widget enabled/disable
	AuditLogChangeKeyWidgetChannelID             = "widget_channel_id"             // guild	snowflake	channel id of the server widget changed
	AuditLogChangeKeyPosition                    = "position"                      // channel	integer	text or voice channel position changed
	AuditLogChangeKeyTopic                       = "topic"                         // channel	string	text channel topic changed
	AuditLogChangeKeyBitrate                     = "bitrate"                       // channel	integer	voice channel bitrate changed
	AuditLogChangeKeyPermissionOverwrites        = "permission_overwrites"         // channel	array of channel overwrite objects	permissions on a channel changed
	AuditLogChangeKeyNSFW                        = "nsfw"                          // channel	bool	channel nsfw restriction changed
	AuditLogChangeKeyApplicationID               = "application_id"                // channel	snowflake	application id of the added or removed webhook or bot
	AuditLogChangeKeyPermissions                 = "permissions"                   // role	integer	permissions for a role changed
	AuditLogChangeKeyColor                       = "color"                         // role	integer	role color changed
	AuditLogChangeKeyHoist                       = "hoist"                         // role	bool	role is now displayed/no longer displayed separate from online users
	AuditLogChangeKeyMentionable                 = "mentionable"                   // role	bool	role is now mentionable/unmentionable
	AuditLogChangeKeyAllow                       = "allow"                         // role	integer	a permission on a text or voice channel was allowed for a role
	AuditLogChangeKeyDeny                        = "deny"                          // role	integer	a permission on a text or voice channel was denied for a role
	AuditLogChangeKeyCode                        = "code"                          // invite	string	invite code changed
	AuditLogChangeKeyChannelID                   = "channel_id"                    // invite	snowflake	channel for invite code changed
	AuditLogChangeKeyInviterID                   = "inviter_id"                    // invite	snowflake	person who created invite code changed
	AuditLogChangeKeyMaxUses                     = "max_uses"                      // invite	integer	change to max number of times invite code can be used
	AuditLogChangeKeyUses                        = "uses"                          // invite	integer	number of times invite code used changed
	AuditLogChangeKeyMaxAge                      = "max_age"                       // invite	integer	how long invite code lasts changed
	AuditLogChangeKeyTemporary                   = "temporary"                     // invite	bool	invite code is temporary/never expires
	AuditLogChangeKeyDeaf                        = "deaf"                          // user	bool	user server deafened/undeafened
	AuditLogChangeKeyMute                        = "mute"                          // user	bool	user server muted/unmuteds
	AuditLogChangeKeyNick                        = "nick"                          // user	string	user nickname changed
	AuditLogChangeKeyAvatarHash                  = "avatar_hash"                   // user	string	user avatar changed
	AuditLogChangeKeyID                          = "id"                            // any	snowflake	the id of the changed entity - sometimes used in conjunction with other keys
	AuditLogChangeKeyType                        = "type"                          // any	integer (channel type) or string	type of entity created
)

// AuditLogParams set params used in endpoint request
// https://discordapp.com/developers/docs/resources/audit-log#get-guild-audit-log-query-string-parameters
type AuditLogParams struct {
	UserID     snowflake.ID `urlparam:"user_id,omitempty"`     // filter the log for a user id
	ActionType uint         `urlparam:"action_type,omitempty"` // the type of audit log event
	Before     snowflake.ID `urlparam:"before,omitempty"`      // filter the log before a certain entry id
	Limit      int          `urlparam:"limit,omitempty"`       // how many entries are returned (default 50, minimum 1, maximum 100)
}

// getQueryString this ins't really pretty, but it works.
func (params *AuditLogParams) getQueryString() string {
	seperator := "?"
	query := ""

	if !params.UserID.Empty() {
		query += seperator + params.UserID.String()
		seperator = "&"
	}

	if params.ActionType > 0 {
		query += seperator + strconv.FormatUint(uint64(params.ActionType), 10)
		seperator = "&"
	}

	if !params.Before.Empty() {
		query += seperator + params.Before.String()
		seperator = "&"
	}

	if params.Limit > 0 {
		query += seperator + strconv.Itoa(params.Limit)
	}

	return query
}

// ReqGuildAuditLogs [GET] Returns an audit log object for the guild.
// 						   Requires the 'VIEW_AUDIT_LOG' permission.
// Endpoint				   /guilds/{guild.id}/audit-logs
// Rate limiter [MAJOR]	   /guilds/{guild.id}
// Discord documentation   https://discordapp.com/developers/docs/resources/audit-log#get-guild-audit-log
// Reviewed				   2018-06-05
// Comment				   -
func ReqGuildAuditLogs(requester request.DiscordGetter, guildID string, params *AuditLogParams) (*AuditLog, error) {
	endpoint := EndpointGuild + "/" + guildID
	path := endpoint + "audit-logs" + params.getQueryString()

	logs := &AuditLog{}
	_, err := requester.Get(endpoint, path, logs)

	return logs, err
}
