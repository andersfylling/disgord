package disgord

type GuildScheduledEvent struct {
	ID                 Snowflake   `json:"id"`
	GuildID            Snowflake   `json:"guild_id"`
	ChannelID          Snowflake   `json:"channel_id"`
	CreatorID          Snowflake   `json:"creator_id"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	ScheduledStartTime Time        `json:"scheduled_start_time"`
	ScheduledEndTime   Time        `json:"scheduled_end_time"`
	PrivacyLevel       int         `json:"privacy_level"`
	EventStatus        int         `json:"event_status"`
	EntityType         int         `json:"entity_type"`
	EntityMetadata     interface{} `json:"entity_metadata"`
	Creator            *User       `json:"creator"`
	UserCount          int         `json:"user_count"`
}
