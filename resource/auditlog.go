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
	AuditLogChangeKeyName                        string = "name"                          // guild	string	name changed
	AuditLogChangeKeyIconHash                    string = "icon_hash"                     // guild	string	icon changed
	AuditLogChangeKeySplashHash                  string = "splash_hash"                   // guild	string	invite splash page artwork changed
	AuditLogChangeKeyOwnerID                     string = "owner_id"                      // guild	snowflake	owner changed
	AuditLogChangeKeyRegion                      string = "region"                        // guild	string	region changed
	AuditLogChangeKeyAFKChannelID                string = "afk_channel_id"                // guild	snowflake	afk channel changed
	AuditLogChangeKeyAFKTimeout                  string = "afk_timeout"                   // guild	integer	afk timeout duration changed
	AuditLogChangeKeyMFALevel                    string = "mfa_level"                     // guild	integer	two-factor auth requirement changed
	AuditLogChangeKeyVerificationLevel           string = "verification_level"            // guild	integer	required verification level changed
	AuditLogChangeKeyExplicitContentFilter       string = "explicit_content_filter"       // guild	integer	change in whose messages are scanned and deleted for explicit content in the server
	AuditLogChangeKeyDefaultMessageNotifications string = "default_message_notifications" // guild	integer	default message notification level changed
	AuditLogChangeKeyVanityURLCode               string = "vanity_url_code"               // guild	string	guild invite vanity url changed
	AuditLogChangeKeyAdd                         string = "$add"                          // add	guild	array of role objects	new role added
	AuditLogChangeKeyRemove                      string = "$remove"                       // remove	guild	array of role objects	role removed
	AuditLogChangeKeyPruneDeleteDays             string = "prune_delete_days"             // guild	integer	change in number of days after which inactive and role-unassigned members are kicked
	AuditLogChangeKeyWidgetEnabled               string = "widget_enabled"                // guild	bool	server widget enabled/disable
	AuditLogChangeKeyWidgetChannelID             string = "widget_channel_id"             // guild	snowflake	channel id of the server widget changed
	AuditLogChangeKeyPosition                    string = "position"                      // channel	integer	text or voice channel position changed
	AuditLogChangeKeyTopic                       string = "topic"                         // channel	string	text channel topic changed
	AuditLogChangeKeyBitrate                     string = "bitrate"                       // channel	integer	voice channel bitrate changed
	AuditLogChangeKeyPermissionOverwrites        string = "permission_overwrites"         // channel	array of channel overwrite objects	permissions on a channel changed
	AuditLogChangeKeyNSFW                        string = "nsfw"                          // channel	bool	channel nsfw restriction changed
	AuditLogChangeKeyApplicationID               string = "application_id"                // channel	snowflake	application id of the added or removed webhook or bot
	AuditLogChangeKeyPermissions                 string = "permissions"                   // role	integer	permissions for a role changed
	AuditLogChangeKeyColor                       string = "color"                         // role	integer	role color changed
	AuditLogChangeKeyHoist                       string = "hoist"                         // role	bool	role is now displayed/no longer displayed separate from online users
	AuditLogChangeKeyMentionable                 string = "mentionable"                   // role	bool	role is now mentionable/unmentionable
	AuditLogChangeKeyAllow                       string = "allow"                         // role	integer	a permission on a text or voice channel was allowed for a role
	AuditLogChangeKeyDeny                        string = "deny"                          // role	integer	a permission on a text or voice channel was denied for a role
	AuditLogChangeKeyCode                        string = "code"                          // invite	string	invite code changed
	AuditLogChangeKeyChannelID                   string = "channel_id"                    // invite	snowflake	channel for invite code changed
	AuditLogChangeKeyInviterID                   string = "inviter_id"                    // invite	snowflake	person who created invite code changed
	AuditLogChangeKeyMaxUses                     string = "max_uses"                      // invite	integer	change to max number of times invite code can be used
	AuditLogChangeKeyUses                        string = "uses"                          // invite	integer	number of times invite code used changed
	AuditLogChangeKeyMaxAge                      string = "max_age"                       // invite	integer	how long invite code lasts changed
	AuditLogChangeKeyTemporary                   string = "temporary"                     // invite	bool	invite code is temporary/never expires
	AuditLogChangeKeyDeaf                        string = "deaf"                          // user	bool	user server deafened/undeafened
	AuditLogChangeKeyMute                        string = "mute"                          // user	bool	user server muted/unmuteds
	AuditLogChangeKeyNick                        string = "nick"                          // user	string	user nickname changed
	AuditLogChangeKeyAvatarHash                  string = "avatar_hash"                   // user	string	user avatar changed
	AuditLogChangeKeyID                          string = "id"                            // any	snowflake	the id of the changed entity - sometimes used in conjunction with other keys
	AuditLogChangeKeyType                        string = "type"                          // any	integer (channel type) or string	type of entity created
)

// AuditLogParams set params used in endpoint request
type AuditLogParams struct {
	UserID     snowflake.ID `urlparam:"user_id,omitempty"`
	ActionType uint         `urlparam:"action_type,omitempty"`
	Before     snowflake.ID `urlparam:"before,omitempty"`
	Limit      int          `urlparam:"limit,omitempty"`
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

func ReqGuildAuditLogs(requester request.DiscordGetter, guildID string, params *AuditLogParams) (*AuditLog, error) {
	endpoint := EndpointGuild + "/" + guildID
	path := endpoint + "audit-logs" + params.getQueryString()

	logs := &AuditLog{}
	_, err := requester.Get(endpoint, path, logs)

	return logs, err
}
