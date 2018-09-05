package endpoint

import . "github.com/andersfylling/snowflake"

func GuildAuditLogs(guildID Snowflake) string {
	return guilds + "/" + guildID.String() + auditlogs
}
