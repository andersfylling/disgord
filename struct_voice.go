package disgord

import "sync"

// VoiceState Voice State structure
// https://discordapp.com/developers/docs/resources/voice#voice-state-object
type VoiceState struct {
	sync.RWMutex `json:"-"`

	// GuildID the guild id this voice state is for
	GuildID Snowflake `json:"guild_id,omitempty"` // ? |

	// ChannelID the channel id this user is connected to
	ChannelID Snowflake `json:"channel_id"` // |

	// UserID the user id this voice state is for
	UserID Snowflake `json:"user_id"` // |

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

//func (vst *VoiceState) Clear() {
//
//}

func (v *VoiceState) DeepCopy() (copy interface{}) {
	copy = &VoiceState{}
	v.CopyOverTo(copy)

	return
}

func (v *VoiceState) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var voiceState *VoiceState
	if voiceState, ok = other.(*VoiceState); !ok {
		err = NewErrorUnsupportedType("given interface{} was not of type *VoiceState")
		return
	}

	v.RLock()
	voiceState.Lock()

	voiceState.GuildID = v.GuildID
	voiceState.ChannelID = v.ChannelID
	voiceState.UserID = v.UserID
	voiceState.SessionID = v.SessionID
	voiceState.Deaf = v.Deaf
	voiceState.Mute = v.Mute
	voiceState.SelfDeaf = v.SelfDeaf
	voiceState.SelfMute = v.SelfMute
	voiceState.Suppress = v.Suppress

	v.RUnlock()
	voiceState.Unlock()

	return
}

// VoiceRegion voice region structure
// https://discordapp.com/developers/docs/resources/voice#voice-region
type VoiceRegion struct {
	// Snowflake unique Snowflake for the region
	ID string `json:"id"`

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
