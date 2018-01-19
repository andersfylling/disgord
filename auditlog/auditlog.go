package auditlog

import (
	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/disgord/webhook"
	"github.com/andersfylling/snowflake"
)

type AuditLog struct {
	Webhooks        []*webhook.Webhook `json:"webhooks"`
	Users           []*user.User       `json:"users"`
	AuditLogEntries []*AuditLogEntry   `json:"audit_log_entries"`
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

type AuditLogChange struct {
	// will this even work? TODO, NOTE
	NewValue interface{} `json:"new_value,omitempty"`
	OldValue interface{} `json:"old_value,omitempty"`
	Key      string      `json:"key"`
}

type AuditLogOption struct {
	DeleteMemberDays string       `json:"delete_member_days"`
	MembersRemoved   string       `json:"members_removed"`
	ChannelID        snowflake.ID `json:"channel_id"`
	Count            string       `json:"count"`
	ID               snowflake.ID `json:"id,string"`
	Type             string       `json:"type"` // type of overwritten entity ("member" or "role")
	RoleName         string       `json:"role_name"`
}

// AuditLogParams set params used in endpoint request
type AuditLogParams struct {
	UserID     snowflake.ID `urlparam:"user_id,omitempty"`
	ActionType uint         `urlparam:"action_type,omitempty"`
	Before     snowflake.ID `urlparam:"before,omitempty"`
	Limit      int          `urlparam:"limit,omitempty"`
}

//
// func convertAuditLogParamsToStr(params *AuditLogParams) string {
// 	var getParams string
//
// 	v := reflect.ValueOf(*params)
// 	t := reflect.TypeOf(*params)
// 	// Iterate over all available fields and read the tag value
// 	for i := 0; i < t.Elem().NumField(); i++ {
// 		// Get the field, returns https://golang.org/pkg/reflect/#StructField
// 		field := t.Field(i)
//
// 		// Get the field tag value
// 		tag := field.Tag.Get("urlparam")
//
// 		// check if it's omitempty
// 		tags := strings.Split(tag, ",")
// 		if len(tags) > 1 {
// 			var skip bool
// 			for _, tagDetail := range tags {
// 				if tagDetail == "omitempty" && reflect.DeepEqual(field, reflect.Zero(field.Type).Interface()) {
// 					skip = true
// 				}
// 			}
// 			if skip {
// 				continue
// 			}
// 		}
//
// 		getParams += "&" + tags[0] + "=" + v.Field(i).Interface().(string)
// 	}
//
// 	urlParams := ""
// 	if getParams != "" {
// 		urlParams = "?" + getParams[1:len(getParams)]
// 	}
//
// 	return urlParams
// }
