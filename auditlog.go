package disgord

type AuditLogEvt uint

// Audit-log event types
const (
	AuditLogEvtGuildUpdate AuditLogEvt = 1
)
const (
	AuditLogEvtChannelCreate AuditLogEvt = 10 + iota
	AuditLogEvtChannelUpdate
	AuditLogEvtChannelDelete
	AuditLogEvtOverwriteCreate
	AuditLogEvtOverwriteUpdate
	AuditLogEvtOverwriteDelete
)
const (
	AuditLogEvtMemberKick AuditLogEvt = 20 + iota
	AuditLogEvtMemberPrune
	AuditLogEvtMemberBanAdd
	AuditLogEvtMemberBanRemove
	AuditLogEvtMemberUpdate
	AuditLogEvtMemberRoleUpdate
	AuditLogEvtMemberMove
	AuditLogEvtMemberDisconnect
	AuditLogEvtBotAdd
)
const (
	AuditLogEvtRoleCreate AuditLogEvt = 30 + iota
	AuditLogEvtRoleUpdate
	AuditLogEvtRoleDelete
)
const (
	AuditLogEvtInviteCreate AuditLogEvt = 40
	AuditLogEvtInviteUpdate
	AuditLogEvtInviteDelete
)
const (
	AuditLogEvtWebhookCreate AuditLogEvt = 50 + iota
	AuditLogEvtWebhookUpdate
	AuditLogEvtWebhookDelete
)
const (
	AuditLogEvtEmojiCreate AuditLogEvt = 60 + iota
	AuditLogEvtEmojiUpdate
	AuditLogEvtEmojiDelete
)
const (
	AuditLogEvtMessageDelete AuditLogEvt = 72
)

type AuditLogChange string

// all the different keys for an audit log change
const (
	// key name,								          identifier                       changed, type,   description
	AuditLogChangeName                        AuditLogChange = "name"                          // guild	string	name changed
	AuditLogChangeIconHash                    AuditLogChange = "icon_hash"                     // guild	string	icon changed
	AuditLogChangeSplashHash                  AuditLogChange = "splash_hash"                   // guild	string	invite splash page artwork changed
	AuditLogChangeOwnerID                     AuditLogChange = "owner_id"                      // guild	snowflake	owner changed
	AuditLogChangeRegion                      AuditLogChange = "region"                        // guild	string	region changed
	AuditLogChangeAFKChannelID                AuditLogChange = "afk_channel_id"                // guild	snowflake	afk channel changed
	AuditLogChangeAFKTimeout                  AuditLogChange = "afk_timeout"                   // guild	integer	afk timeout duration changed
	AuditLogChangeMFALevel                    AuditLogChange = "mfa_level"                     // guild	integer	two-factor auth requirement changed
	AuditLogChangeVerificationLevel           AuditLogChange = "verification_level"            // guild	integer	required verification level changed
	AuditLogChangeExplicitContentFilter       AuditLogChange = "explicit_content_filter"       // guild	integer	change in whose messages are scanned and deleted for explicit content in the server
	AuditLogChangeDefaultMessageNotifications AuditLogChange = "default_message_notifications" // guild	integer	default message notification level changed
	AuditLogChangeVanityURLCode               AuditLogChange = "vanity_url_code"               // guild	string	guild invite vanity url changed
	AuditLogChangeAdd                         AuditLogChange = "$add"                          // add	guild	array of role objects	new role added
	AuditLogChangeRemove                      AuditLogChange = "$remove"                       // remove	guild	array of role objects	role removed
	AuditLogChangePruneDeleteDays             AuditLogChange = "prune_delete_days"             // guild	integer	change in number of days after which inactive and role-unassigned members are kicked
	AuditLogChangeWidgetEnabled               AuditLogChange = "widget_enabled"                // guild	bool	server widget enabled/disable
	AuditLogChangeWidgetChannelID             AuditLogChange = "widget_channel_id"             // guild	snowflake	channel id of the server widget changed
	AuditLogChangePosition                    AuditLogChange = "position"                      // channel	integer	text or voice channel position changed
	AuditLogChangeTopic                       AuditLogChange = "topic"                         // channel	string	text channel topic changed
	AuditLogChangeBitrate                     AuditLogChange = "bitrate"                       // channel	integer	voice channel bitrate changed
	AuditLogChangePermissionOverwrites        AuditLogChange = "permission_overwrites"         // channel	array of channel overwrite objects	permissions on a channel changed
	AuditLogChangeNSFW                        AuditLogChange = "nsfw"                          // channel	bool	channel nsfw restriction changed
	AuditLogChangeApplicationID               AuditLogChange = "application_id"                // channel	snowflake	application id of the added or removed webhook or bot
	AuditLogChangePermissions                 AuditLogChange = "permissions"                   // role	integer	permissions for a role changed
	AuditLogChangeColor                       AuditLogChange = "color"                         // role	integer	role color changed
	AuditLogChangeHoist                       AuditLogChange = "hoist"                         // role	bool	role is now displayed/no longer displayed separate from online users
	AuditLogChangeMentionable                 AuditLogChange = "mentionable"                   // role	bool	role is now mentionable/unmentionable
	AuditLogChangeAllow                       AuditLogChange = "allow"                         // role	integer	a permission on a text or voice channel was allowed for a role
	AuditLogChangeDeny                        AuditLogChange = "deny"                          // role	integer	a permission on a text or voice channel was denied for a role
	AuditLogChangeCode                        AuditLogChange = "code"                          // invite	string	invite code changed
	AuditLogChangeChannelID                   AuditLogChange = "channel_id"                    // invite	snowflake	channel for invite code changed
	AuditLogChangeInviterID                   AuditLogChange = "inviter_id"                    // invite	snowflake	person who created invite code changed
	AuditLogChangeMaxUses                     AuditLogChange = "max_uses"                      // invite	integer	change to max number of times invite code can be used
	AuditLogChangeUses                        AuditLogChange = "uses"                          // invite	integer	number of times invite code used changed
	AuditLogChangeMaxAge                      AuditLogChange = "max_age"                       // invite	integer	how long invite code lasts changed
	AuditLogChangeTemporary                   AuditLogChange = "temporary"                     // invite	bool	invite code is temporary/never expires
	AuditLogChangeDeaf                        AuditLogChange = "deaf"                          // user	bool	user server deafened/undeafened
	AuditLogChangeMute                        AuditLogChange = "mute"                          // user	bool	user server muted/unmuteds
	AuditLogChangeNick                        AuditLogChange = "nick"                          // user	string	user nickname changed
	AuditLogChangeAvatarHash                  AuditLogChange = "avatar_hash"                   // user	string	user avatar changed
	AuditLogChangeID                          AuditLogChange = "id"                            // any	snowflake	the id of the changed entity - sometimes used in conjunction with other keys
	AuditLogChangeType                        AuditLogChange = "type"                          // any	integer (channel type) or string	type of entity created
)

// AuditLog ...
type AuditLog struct {
	Webhooks        []*Webhook       `json:"webhooks"`
	Users           []*User          `json:"users"`
	AuditLogEntries []*AuditLogEntry `json:"audit_log_entries"`
}

var _ Copier = (*AuditLog)(nil)
var _ DeepCopier = (*AuditLog)(nil)

func (l *AuditLog) Bans() (bans []*PartialBan) {
	for _, e := range l.AuditLogEntries {
		if e.Event == AuditLogEvtMemberBanAdd {
			bans = append(bans, &PartialBan{
				Reason:                 e.Reason,
				ModeratorResponsibleID: e.UserID,
				BannedUserID:           e.TargetID,
			})
		}
	}
	return bans
}

// AuditLogEntry ...
type AuditLogEntry struct {
	TargetID Snowflake          `json:"target_id"`
	Changes  []*AuditLogChanges `json:"changes,omitempty"`
	UserID   Snowflake          `json:"user_id"`
	ID       Snowflake          `json:"id"`
	Event    AuditLogEvt        `json:"action_type"`
	Options  *AuditLogOption    `json:"options,omitempty"`
	Reason   string             `json:"reason,omitempty"`
}

var _ Copier = (*AuditLogEntry)(nil)
var _ DeepCopier = (*AuditLogEntry)(nil)

// AuditLogOption ...
type AuditLogOption struct {
	DeleteMemberDays string    `json:"delete_member_days"`
	MembersRemoved   string    `json:"members_removed"`
	ChannelID        Snowflake `json:"channel_id"`
	Count            string    `json:"count"`
	ID               Snowflake `json:"id"`
	Type             string    `json:"type"` // type of overwritten entity ("member" or "role")
	RoleName         string    `json:"role_name"`
}

var _ Copier = (*AuditLogOption)(nil)
var _ DeepCopier = (*AuditLogOption)(nil)

// AuditLogChanges ...
type AuditLogChanges struct {
	NewValue interface{} `json:"new_value,omitempty"`
	OldValue interface{} `json:"old_value,omitempty"`
	Key      string      `json:"key"`
}

var _ Copier = (*AuditLogChanges)(nil)
var _ DeepCopier = (*AuditLogChanges)(nil)

// auditLogFactory temporary until flyweight is implemented
func auditLogFactory() interface{} {
	return &AuditLog{}
}

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

// guildAuditLogsBuilder for building the GetGuildAuditLogs request.
// TODO: support caching of audit log entries. So we only fetch those we don't have.
//generate-rest-params: user_id:Snowflake, action_type:uint, before:Snowflake, limit:int,
//generate-rest-basic-execute: log:*AuditLog,
type guildAuditLogsBuilder struct {
	r RESTBuilder
}
