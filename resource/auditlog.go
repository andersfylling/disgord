package resource

import . "github.com/andersfylling/snowflake"

type AuditLog struct {
	Webhooks        []*Webhook       `json:"webhooks"`
	Users           []*User          `json:"users"`
	AuditLogEntries []*AuditLogEntry `json:"audit_log_entries"`
}

type AuditLogEntry struct {
	TargetID   Snowflake         `json:"target_id"`
	Changes    []*AuditLogChange `json:"changes,omitempty"`
	UserID     Snowflake         `json:"user_id"`
	ID         Snowflake         `json:"id"`
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
	DeleteMemberDays string    `json:"delete_member_days"`
	MembersRemoved   string    `json:"members_removed"`
	ChannelID        Snowflake `json:"channel_id"`
	Count            string    `json:"count"`
	ID               Snowflake `json:"id"`
	Type             string    `json:"type"` // type of overwritten entity ("member" or "role")
	RoleName         string    `json:"role_name"`
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
