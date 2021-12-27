package disgord

// ThreadMember https://discord.com/developers/docs/resources/channel#thread-member-object
type ThreadMember struct {
	ThreadID            Snowflake `json:"id,omitempty"`
	UserID        Snowflake `json:"user_id,omitempty"`
	JoinTimestamp Time      `json:"join_timestamp"`
	Flags         Flag      `json:"flags"`
}

// ThreadMetadata https://discord.com/developers/docs/resources/channel#thread-metadata-object
type ThreadMetadata struct {
	Archived            bool `json:"archived"`
	AutoArchiveDuration int  `json:"auto_archive_duration"`
	ArchiveTimestamp    Time `json:"archive_timestamp"`
	Locked              bool `json:"locked"`
	Invitable           bool `json:"inviteable,omitempty"`
}

type AutoArchiveDurationTime int

const (
	AutoArchiveThreadMinute AutoArchiveDurationTime = 60
	AutoArchiveThreadDay    AutoArchiveDurationTime = 1440
	// guild must be boosted to use the below auto archive durations.
	// ref: https://discord.com/developers/docs/resources/channel#start-thread-with-message-json-params
	AutoArchiveThreadThreeDays AutoArchiveDurationTime = 4320
	AutoArchiveThreadWeek      AutoArchiveDurationTime = 10080
)
