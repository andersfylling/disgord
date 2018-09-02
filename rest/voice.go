package rest

import (
	. "github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/disgord/rest/endpoint"
	"github.com/andersfylling/disgord/rest/httd"
)

// ListVoiceRegions [GET]   Returns an array of voice region objects that can be used when creating servers.
// Endpoint                 /voice/regions
// Rate limiter             /voice/regions
// Discord documentation    https://discordapp.com/developers/docs/resources/voice#list-voice-regions
// Reviewed                 2018-08-21
// Comment                  -
func ListVoiceRegions(client httd.Getter) (ret []*VoiceRegion, err error) {
	details := &httd.Request{
		Ratelimiter: endpoint.VoiceRegions(),
		Endpoint:    endpoint.VoiceRegions(),
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}
