package disgord

// ref https://discord.com/developers/docs/resources/channel#thread-member-object
type ThreadMember struct {
	ID Snowflake `json:"id,omitempty"`
	UserID Snowflake `json:"user_id,omitempty"`
	JoinTimestamp Time `json:"join_timestamp"`
	Flags Flag `json:"flags"`
}

// ref https://discord.com/developers/docs/resources/channel#thread-metadata-object
type ThreadMetadata struct {
	Archived bool `json:"archived"`
	AutoArchiveDuration	int `json:"auto_archive_duration"`
	ArchiveTimestamp Time `json:"archive_timestamp"`
	Locked bool `json:"locked"`
	Invitable bool `json:"inviteable,omitempty"`
}
