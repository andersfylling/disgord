package endpoint

import "fmt"

func GuildAuditLogs(guildID fmt.Stringer) string {
	return guilds + "/" + guildID.String() + auditlogs
}
