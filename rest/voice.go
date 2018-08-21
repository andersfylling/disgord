package rest

import (
	"encoding/json"

	. "github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/disgord/rest/httd"
)

// EndpointVoiceRegions List Voice Regions
// https://discordapp.com/developers/docs/resources/voice#list-voice-regions
const EndpointVoiceRegions = "/voice/regions"

// ListVoiceRegions [GET]   Returns an array of voice region objects that can be used when creating servers.
// Endpoint                 /voice/regions
// Rate limiter             /voice/regions
// Discord documentation    https://discordapp.com/developers/docs/resources/voice#list-voice-regions
// Reviewed                 2018-08-21
// Comment                  -
func ListVoiceRegions(client httd.Getter) (ret []*VoiceRegion, err error) {
	details := &httd.Request{
		Ratelimiter: EndpointVoiceRegions,
		Endpoint:    EndpointVoiceRegions,
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}
