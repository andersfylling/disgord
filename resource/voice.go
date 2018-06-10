package resource

import (
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/snowflake"
)

// State Voice State structure
// https://discordapp.com/developers/docs/resources/voice#voice-state-object
type VoiceState struct {
	// GuildID the guild id this voice state is for
	GuildID snowflake.ID `json:"guild_id,omitempty"` // ? |

	// ChannelID the channel id this user is connected to
	ChannelID snowflake.ID `json:"channel_id"` // |

	// UserID the user id this voice state is for
	UserID snowflake.ID `json:"user_id"` // |

	// SessionID the session id for this voice state
	SessionID string `json:"session_id"` // |

	// Deaf whether this user is deafened by the server
	Deaf bool `json:"deaf"` // |

	// Mute whether this user is muted by the server
	Mute bool `json:"mute"` // |

	// SelfDeaf whether this user is locally deafened
	SelfDeaf bool `json:"self_deaf"` // |

	// SelfMute whether this user is locally muted
	SelfMute bool `json:"self_mute"` // |

	// Suppress whether this user is muted by the current user
	Suppress bool `json:"suppress"` // |
}

func (vst *VoiceState) Clear() {

}

// Region voice region structure
// https://discordapp.com/developers/docs/resources/voice#voice-region
type VoiceRegion struct {
	// ID unique ID for the region
	ID snowflake.ID `json:"id"`

	// Name name of the region
	Name string `json:"name"`

	// SampleHostname an example hostname for the region
	SampleHostname string `json:"sample_hostname"`

	// SamplePort an example port for the region
	SamplePort uint `json:"sample_port"`

	// VIP true if this is a vip-only server
	VIP bool `json:"vip"`

	// Optimal true for a single server that is closest to the current user's client
	Optimal bool `json:"optimal"`

	// Deprecated 	whether this is a deprecated voice region (avoid switching to these)
	Deprecated bool `json:"deprecated"`

	// Custom whether this is a custom voice region (used for events/etc)
	Custom bool `json:"custom"`
}

// EndpointVoiceRegions List Voice Regions
// https://discordapp.com/developers/docs/resources/voice#list-voice-regions
const EndpointVoiceRegions = "/voice/regions"

func ReqVoiceRegions(requester httd.Getter) (regions []*VoiceRegion, err error) {
	_, err = requester.Get(EndpointVoiceRegions, EndpointVoiceRegions, regions)

	return regions, err
}
