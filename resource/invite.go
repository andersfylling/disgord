package resource

import (
	"encoding/json"

	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/httd"
)

// Invite Represents a code that when used, adds a user to a guild.
// https://discordapp.com/developers/docs/resources/invite#invite-object
// Reviewed: 2018-06-10
type Invite struct {
	// Code the invite code (unique ID)
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
	CreatedAt discord.Timestamp `json:"created_at"`

	// Revoked whether this invite is revoked
	Revoked bool `json:"revoked"`
}

const (
	EndpointInvite = "/invites"
)

// ReqGetInvite [GET]     Returns an invite object for the given code.
// Endpoint               /invites/{invite.code}
// Rate limiter           /invites/{invite.code}
// Discord documentation  https://discordapp.com/developers/docs/resources/invite#get-invite
// Reviewed               2018-06-10
// Comment                -
//
// withCounts whether the invite should contain approximate member counts
func ReqGetInvite(requester httd.Getter, inviteCode string, withCounts bool) (invite *Invite, err error) {
	query := ""
	if withCounts {
		query += "?with_counts=true"
	}

	details := &httd.Request{
		Ratelimiter: EndpointInvite + "/" + inviteCode,
		Endpoint:    query,
	}
	resp, err := requester.Get(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(invite)
	return
}

// ReqDeleteInvite [DELETE] Delete an invite. Requires the MANAGE_CHANNELS permission. Returns an invite
//                          object on success.
// Endpoint                 /invites/{invite.code}
// Rate limiter             /invites/{invite.code}
// Discord documentation    https://discordapp.com/developers/docs/resources/invite#delete-invite
// Reviewed                 2018-06-10
// Comment                  -
func ReqDeleteInvite(requester httd.Deleter, inviteCode string) (invite *Invite, err error) {

	details := &httd.Request{
		Ratelimiter: EndpointInvite + "/" + inviteCode,
	}
	resp, err := requester.Delete(details)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(invite)
	return
}
