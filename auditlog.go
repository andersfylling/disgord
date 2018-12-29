package disgord

import (
	"net/http"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/ratelimit"
	"github.com/andersfylling/snowflake/v3"
)

// Audit-log event types
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

// all the different keys for an audit log change
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

// AuditLog ...
type AuditLog struct {
	Lockable `json:"-"`

	Webhooks        []*Webhook       `json:"webhooks"`
	Users           []*User          `json:"users"`
	AuditLogEntries []*AuditLogEntry `json:"audit_log_entries"`
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

	if constant.LockedMethods {
		l.RLock()
		log.Lock()
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

	if constant.LockedMethods {
		l.RUnlock()
		log.Unlock()
	}
	return
}

// AuditLogEntry ...
type AuditLogEntry struct {
	Lockable `json:"-"`

	TargetID   Snowflake         `json:"target_id"`
	Changes    []*AuditLogChange `json:"changes,omitempty"`
	UserID     Snowflake         `json:"user_id"`
	ID         Snowflake         `json:"id"`
	ActionType uint              `json:"action_type"`
	Options    []*AuditLogOption `json:"options,omitempty"`
	Reason     string            `json:"reason,omitempty"`
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

	if constant.LockedMethods {
		l.RLock()
		log.Lock()
	}

	log.TargetID = l.TargetID
	log.UserID = l.UserID
	log.ID = l.ID
	log.ActionType = l.ActionType
	log.Reason = l.Reason

	for _, change := range l.Changes {
		log.Changes = append(log.Changes, change.DeepCopy().(*AuditLogChange))
	}

	for _, option := range l.Options {
		log.Options = append(log.Options, option.DeepCopy().(*AuditLogOption))
	}

	if constant.LockedMethods {
		l.RUnlock()
		log.Unlock()
	}
	return
}

// AuditLogOption ...
type AuditLogOption struct {
	Lockable `json:"-"`

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

	if constant.LockedMethods {
		l.RLock()
		log.Lock()
	}

	log.DeleteMemberDays = l.DeleteMemberDays
	log.MembersRemoved = l.MembersRemoved
	log.ChannelID = l.ChannelID
	log.Count = l.Count
	log.ID = l.ID
	log.Type = l.Type
	log.RoleName = l.RoleName

	if constant.LockedMethods {
		l.RUnlock()
		log.Unlock()
	}
	return
}

// AuditLogChange ...
type AuditLogChange struct {
	Lockable `json:"-"`

	NewValue interface{} `json:"new_value,omitempty"`
	OldValue interface{} `json:"old_value,omitempty"`
	Key      string      `json:"key"`
}

// DeepCopy see interface at struct.go#DeepCopier
func (l *AuditLogChange) DeepCopy() (copy interface{}) {
	copy = &AuditLogChange{}
	l.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (l *AuditLogChange) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var log *AuditLogChange
	if log, ok = other.(*AuditLogChange); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *AuditLogChange")
		return
	}

	if constant.LockedMethods {
		l.RLock()
		log.Lock()
	}

	log.NewValue = l.NewValue
	log.OldValue = l.OldValue
	log.Key = l.Key

	if constant.LockedMethods {
		l.RUnlock()
		log.Unlock()
	}

	return
}

// auditLogFactory temporary until flyweight is implemented
func auditLogFactory() interface{} {
	return &AuditLog{}
}

// GetGuildAuditLogs [REST] Returns an audit log object for the guild. Requires the 'VIEW_AUDIT_LOG' permission.
// Note that this request will _always_ send a REST request, regardless of you calling IgnoreCache or not.
//  Method                   GET
//  Endpoint                 /guilds/{guild.id}/audit-logs
//  Rate limiter [MAJOR]     /guilds/{guild.id}/audit-logs
//  Discord documentation    https://discordapp.com/developers/docs/resources/audit-log#get-guild-audit-log
//  Reviewed                 2018-06-05
//  Comment                  -
//  Note                     Check the last entry in the cacheLink, to avoid fetching data we already got
func (c *Client) GetGuildAuditLogs(guildID snowflake.ID) (builder *guildAuditLogsBuilder) {
	builder = &guildAuditLogsBuilder{}
	builder.r.itemFactory = auditLogFactory
	builder.r.IgnoreCache().setup(c.cache, c.req, &httd.Request{
		Method:      http.MethodGet,
		Ratelimiter: ratelimit.GuildAuditLogs(guildID),
		Endpoint:    endpoint.GuildAuditLogs(guildID), // body are added automatically
	}, nil)

	return builder
}

// guildAuditLogsBuilder for building the GetGuildAuditLogs request
type guildAuditLogsBuilder struct {
	r RESTRequestBuilder
}

// UserID filter the log for a user id
func (b *guildAuditLogsBuilder) UserID(id snowflake.ID) *guildAuditLogsBuilder {
	b.r.queryParam("user_id", id)
	return b
}

// ActionType the type of audit log event
func (b *guildAuditLogsBuilder) ActionType(action uint) *guildAuditLogsBuilder {
	b.r.queryParam("action_type", action)
	return b
}

// Before filter the log before a certain entry id
func (b *guildAuditLogsBuilder) Before(id snowflake.ID) *guildAuditLogsBuilder {
	b.r.queryParam("before", id)
	return b
}

// Before filter the log before a certain entry id
func (b *guildAuditLogsBuilder) Limit(limit int) *guildAuditLogsBuilder {
	b.r.queryParam("limit", limit)
	return b
}

func (b *guildAuditLogsBuilder) IgnoreCache() *guildAuditLogsBuilder {
	b.r.IgnoreCache()
	return b
}

func (b *guildAuditLogsBuilder) CancelOnRatelimit() *guildAuditLogsBuilder {
	b.r.CancelOnRatelimit()
	return b
}

func (b *guildAuditLogsBuilder) Execute() (log *AuditLog, err error) {
	// TODO: support caching of audit log entries. So we only fetch those we don't have.
	var v interface{}
	v, err = b.r.execute()
	if err != nil {
		return
	}

	log = v.(*AuditLog)
	return
}
