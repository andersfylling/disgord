package resource

// Invite Represents a code that when used, adds a user to a guild.
// https://discordapp.com/developers/docs/resources/invite#invite-object
// Reviewed: 2018-06-10
type Invite struct {
	// Code the invite code (unique Snowflake)
	Code string `json:"code"`

	// Guild the guild this invite is for
	Guild *PartialGuild `json:"guild"`

	// Channel the channel this invite is for
	Channel *PartialChannel `json:"channel"`

	// ApproximatePresenceCount approximate count of online members
	ApproximatePresenceCount int `json:"approximate_presence_count,omitempty"`

	// ApproximatePresenceCount approximate count of total members
	ApproximateMemberCount int `json:"approximate_member_count,omitempty"`
}

// InviteMetadata Object
// https://discordapp.com/developers/docs/resources/invite#invite-metadata-object
// Reviewed: 2018-06-10
type InviteMetadata struct {
	// Inviter user who created the invite
	Inviter *User `json:"inviter"`

	// Uses number of times this invite has been used
	Uses int `json:"uses"`

	// MaxUses max number of times this invite can be used
	MaxUses int `json:"max_uses"`

	// MaxAge duration (in seconds) after which the invite expires
	MaxAge int `json:"max_age"`

	// Temporary whether this invite only grants temporary membership
	Temporary bool `json:"temporary"`

	// CreatedAt when this invite was created
	CreatedAt Timestamp `json:"created_at"`

	// Revoked whether this invite is revoked
	Revoked bool `json:"revoked"`
}

// PartialInvite
// {
//    "code": "abc"
// }
type PartialInvite = Invite
