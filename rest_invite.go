package disgord

import (
	"fmt"
	"net/http"

	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
)

// GetInvite [GET]          Returns an invite object for the given code.
// Endpoint                 /invites/{invite.code}
// Rate limiter             /invites
// Discord documentation    https://discordapp.com/developers/docs/resources/invite#get-invite
// Reviewed                 2018-06-10
// Comment                  -
//
// withCounts whether the invite should contain approximate member counts
func GetInvite(client httd.Getter, inviteCode string, withCounts bool) (invite *Invite, err error) {
	query := ""
	if withCounts {
		query += "?with_counts=true"
	}

	resp, body, err := client.Get(&httd.Request{
		Ratelimiter: endpoint.Invites(),
		Endpoint:    endpoint.Invite(inviteCode) + query,
	})
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
	}

	err = unmarshal(body, &invite)
	return
}

// DeleteInvite [DELETE]    Delete an invite. Requires the MANAGE_CHANNELS permission. Returns an invite
//                          object on success.
// Endpoint                 /invites/{invite.code}
// Rate limiter             /invites
// Discord documentation    https://discordapp.com/developers/docs/resources/invite#delete-invite
// Reviewed                 2018-06-10
// Comment                  -
func DeleteInvite(client httd.Deleter, inviteCode string) (invite *Invite, err error) {
	_, body, err := client.Delete(&httd.Request{
		Ratelimiter: endpoint.Invites(),
		Endpoint:    endpoint.Invite(inviteCode),
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &invite)
	return
}
