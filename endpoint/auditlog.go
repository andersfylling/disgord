package endpoint

import "fmt"

// GuildAuditLogs ...
func GuildAuditLogs(guildID fmt.Stringer) string {
	return guilds + "/" + guildID.String() + auditlogs
}
