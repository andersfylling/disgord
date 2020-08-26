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

// DeepCopy see interface at struct.go#DeepCopier
func (l *AuditLog) DeepCopy() (copy interface{}) {
	copy = &AuditLog{}
	l.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (l *AuditLog) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var log *AuditLog
	if log, ok = other.(*AuditLog); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *AuditLog")
		return
	}

	for _, webhook := range l.Webhooks {
		log.Webhooks = append(log.Webhooks, webhook.DeepCopy().(*Webhook))
	}
	for _, user := range l.Users {
		log.Users = append(log.Users, user.DeepCopy().(*User))
	}
	for _, entry := range l.AuditLogEntries {
		log.AuditLogEntries = append(log.AuditLogEntries, entry.DeepCopy().(*AuditLogEntry))
	}
	return
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

// DeepCopy see interface at struct.go#DeepCopier
func (l *AuditLogEntry) DeepCopy() (copy interface{}) {
	copy = &AuditLogEntry{}
	l.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (l *AuditLogEntry) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var log *AuditLogEntry
	if log, ok = other.(*AuditLogEntry); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *AuditLogEntry")
		return
	}

	log.TargetID = l.TargetID
	log.UserID = l.UserID
	log.ID = l.ID
	log.Event = l.Event
	log.Reason = l.Reason

	for _, change := range l.Changes {
		log.Changes = append(log.Changes, change.DeepCopy().(*AuditLogChanges))
	}

	if l.Options != nil {
		log.Options = l.Options.DeepCopy().(*AuditLogOption)
	}
	return
}

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

// DeepCopy see interface at struct.go#DeepCopier
func (l *AuditLogOption) DeepCopy() (copy interface{}) {
	copy = &AuditLogOption{}
	l.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (l *AuditLogOption) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var log *AuditLogOption
	if log, ok = other.(*AuditLogOption); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *AuditLogOption")
		return
	}

	log.DeleteMemberDays = l.DeleteMemberDays
	log.MembersRemoved = l.MembersRemoved
	log.ChannelID = l.ChannelID
	log.Count = l.Count
	log.ID = l.ID
	log.Type = l.Type
	log.RoleName = l.RoleName
	return
}

// AuditLogChanges ...
type AuditLogChanges struct {
	NewValue interface{} `json:"new_value,omitempty"`
	OldValue interface{} `json:"old_value,omitempty"`
	Key      string      `json:"key"`
}

// DeepCopy see interface at struct.go#DeepCopier
func (l *AuditLogChanges) DeepCopy() (copy interface{}) {
	copy = &AuditLogChanges{}
	l.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (l *AuditLogChanges) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var log *AuditLogChanges
	if log, ok = other.(*AuditLogChanges); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *AuditLogChanges")
		return
	}

	log.NewValue = l.NewValue
	log.OldValue = l.OldValue
	log.Key = l.Key

	return
}

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
