package resource

import (
	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/request"
)

type Invite struct {
	// Code the invite code (unique ID)
	Code string `json:"code"`

	// Guild the guild this invite is for
	Guild *PartialGuild

	// Channel the channel this invite is for
	Channel *PartialChannel
}

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
	EndpointInvite string = "/invites"
)

// ReqGetInvite Returns an invite object for the given code.
func ReqGetInvite(requester request.DiscordGetter, code string) (invite *Invite, err error) {
	path := EndpointInvite + "/" + code
	_, err = requester.Get(EndpointInvite, path, invite)

	return invite, err
}

// ReqDeleteInvite Delete an invite. Requires the MANAGE_CHANNELS permission. Returns an invite object on success.
func ReqDeleteInvite(requester request.DiscordGetter, code string) (invite *Invite, err error) {
	path := EndpointInvite + "/" + code
	_, err = requester.Get(EndpointInvite, path, invite)

	return invite, err
}

// @Deprecated
// func ReqAcceptInvite() {}
