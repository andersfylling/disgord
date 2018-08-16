package rest

import (
	"encoding/json"

	"github.com/andersfylling/disgord/rest/httd"
	. "github.com/andersfylling/disgord/resource"
)

// EndpointVoiceRegions List Voice Regions
// https://discordapp.com/developers/docs/resources/voice#list-voice-regions
const EndpointVoiceRegions = "/voice/regions"

func ReqVoiceRegions(client httd.Getter) (ret []*VoiceRegion, err error) {
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
