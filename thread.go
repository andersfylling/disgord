package disgord

// ref https://discord.com/developers/docs/resources/channel#thread-member-object
type ThreadMember struct {
	ID Snowflake
	UserID Snowflake
	JoinTimestamp Time
	Flags Flag
}

// ref https://discord.com/developers/docs/resources/channel#thread-metadata-object
type ThreadMetadata struct {
	Archived bool
	AutoArchiveDuration	int
	ArchiveTimestamp Time
	Locked bool
	Invitable bool
}
