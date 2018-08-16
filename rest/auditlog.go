package rest

import (
	"encoding/json"
	"strconv"

	"github.com/andersfylling/disgord/rest/httd"
	. "github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/snowflake"
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

// ReqGuildAuditLogs [GET]  Returns an audit log object for the guild.
//                          Requires the 'VIEW_AUDIT_LOG' permission.
// Endpoint                 /guilds/{guild.id}/audit-logs
// Rate limiter [MAJOR]     /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/audit-log#get-guild-audit-log
// Reviewed                 2018-06-05
// Comment                  -
func GuildAuditLogs(client httd.Getter, guildID snowflake.ID, params *AuditLogParams) (log *AuditLog, err error) {

	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    EndpointGuild + guildID.String() + "/audit-logs" + params.getQueryString(),
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &log)
	return
}
