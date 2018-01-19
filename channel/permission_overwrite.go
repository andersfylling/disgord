package channel

import "github.com/andersfylling/snowflake"

type PermissionOverwrite struct {
	ID    snowflake.ID `json:"id"`    // role or user id
	Type  string       `json:"type"`  // either `role` or `member`
	Deny  int          `json:"deny"`  // permission bit set
	Allow int          `json:"allow"` // permission bit set
}
