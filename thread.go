package disgord

// ref https://discord.com/developers/docs/resources/channel#thread-member-object
type ThreadMember struct {
	ID            Snowflake `json:"id,omitempty"`
	UserID        Snowflake `json:"user_id,omitempty"`
	JoinTimestamp Time      `json:"join_timestamp"`
	Flags         Flag      `json:"flags"`
}

// ref https://discord.com/developers/docs/resources/channel#thread-metadata-object
type ThreadMetadata struct {
	Archived            bool `json:"archived"`
	AutoArchiveDuration int  `json:"auto_archive_duration"`
	ArchiveTimestamp    Time `json:"archive_timestamp"`
	Locked              bool `json:"locked"`
	Invitable           bool `json:"inviteable,omitempty"`
}

// ref https://discord.com/developers/docs/resources/channel#start-thread-with-message-json-params
type CreateThreadParams struct {
	Name string `json:"name"`
	AutoArchiveDuration int `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser int `json:"rate_limit_per_user,omitempty"`
	
	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string `json:"-"`
}

// ref https://discord.com/developers/docs/resources/channel#start-thread-without-message-json-params
type CreateThreadParamsNoMessage struct {
	Name string `json:"name"`
	AutoArchiveDuration int `json:"auto_archive_duration,omitempty"`
	// In API v9, type defaults to PRIVATE_THREAD in order to match the behavior when
	// thread documentation was first published. In API v10 this will be changed to be a required field, with no default.
	Type int `json:"type,omitempty"`
	Invitable bool `json:"invitable,omitempty"`
	RateLimitPerUser int `json:"rate_limit_per_user,omitempty"`
	
	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string `json:"-"`
}
