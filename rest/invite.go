package rest

import (
	. "github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/disgord/rest/httd"
)

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
func GetInvite(client httd.Getter, inviteCode string, withCounts bool) (invite *Invite, err error) {
	query := ""
	if withCounts {
		query += "?with_counts=true"
	}

	details := &httd.Request{
		Ratelimiter: EndpointInvite,
		Endpoint:    EndpointInvite + "/" + inviteCode + query,
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &invite)
	return
}

// ReqDeleteInvite [DELETE] Delete an invite. Requires the MANAGE_CHANNELS permission. Returns an invite
//                          object on success.
// Endpoint                 /invites/{invite.code}
// Rate limiter             /invites/{invite.code}
// Discord documentation    https://discordapp.com/developers/docs/resources/invite#delete-invite
// Reviewed                 2018-06-10
// Comment                  -
func DeleteInvite(client httd.Deleter, inviteCode string) (invite *Invite, err error) {

	details := &httd.Request{
		Ratelimiter: EndpointInvite,
		Endpoint:    EndpointInvite + "/" + inviteCode,
	}
	_, body, err := client.Delete(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &invite)
	return
}
