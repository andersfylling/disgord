package disgord

import (
	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
	"github.com/andersfylling/disgord/json"
)

// VoiceState Voice State structure
// https://discord.com/developers/docs/resources/voice#voice-state-object
// reviewed 2018-09-29
type VoiceState struct {
	// GuildID the guild id this voice state is for
	GuildID Snowflake `json:"guild_id,omitempty"` // ? |

	// ChannelID the channel id this user is connected to
	ChannelID Snowflake `json:"channel_id"` // | ?

	// UserID the user id this voice state is for
	UserID Snowflake `json:"user_id"` // |

	// the guild member this voice state is for
	Member *Member `json:"member,omitempty"`

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

var _ Reseter = (*VoiceState)(nil)
var _ Copier = (*VoiceState)(nil)
var _ DeepCopier = (*VoiceState)(nil)

//func (vst *VoiceState) Clear() {
//
//}

// DeepCopy see interface at struct.go#DeepCopier
func (v *VoiceState) DeepCopy() (copy interface{}) {
	copy = &VoiceState{}
	v.CopyOverTo(copy)

	return
}

// UnmarshalJSON is used to unmarshal Discord's JSON.
func (v *VoiceState) UnmarshalJSON(data []byte) error {
	type s2 VoiceState
	if err := json.Unmarshal(data, (*s2)(v)); err != nil {
		return err
	}
	if v.Member != nil {
		v.Member.GuildID = v.GuildID
		v.Member.UserID = v.UserID
	}
	return nil
}

// CopyOverTo see interface at struct.go#Copier
func (v *VoiceState) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var voiceState *VoiceState
	if voiceState, ok = other.(*VoiceState); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *VoiceState")
		return
	}

	voiceState.GuildID = v.GuildID
	voiceState.ChannelID = v.ChannelID
	voiceState.UserID = v.UserID
	voiceState.SessionID = v.SessionID
	voiceState.Deaf = v.Deaf
	voiceState.Mute = v.Mute
	voiceState.SelfDeaf = v.SelfDeaf
	voiceState.SelfMute = v.SelfMute
	voiceState.Suppress = v.Suppress

	return
}

// VoiceRegion voice region structure
// https://discord.com/developers/docs/resources/voice#voice-region
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

	// Optimal true for a single server that is closest to the current user's Client
	Optimal bool `json:"optimal"`

	// Deprecated 	whether this is a deprecated voice region (avoid switching to these)
	Deprecated bool `json:"deprecated"`

	// Custom whether this is a custom voice region (used for events/etc)
	Custom bool `json:"custom"`
}

var _ Reseter = (*VoiceRegion)(nil)
var _ Copier = (*VoiceRegion)(nil)
var _ DeepCopier = (*VoiceRegion)(nil)

// DeepCopy see interface at struct.go#DeepCopier
func (v *VoiceRegion) DeepCopy() (copy interface{}) {
	copy = &VoiceRegion{}
	v.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (v *VoiceRegion) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var voice *VoiceRegion
	if voice, ok = other.(*VoiceRegion); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *VoiceRegion")
		return
	}

	voice.ID = v.ID
	voice.Name = v.Name
	voice.SampleHostname = v.SampleHostname
	voice.SamplePort = v.SamplePort
	voice.VIP = v.VIP
	voice.Optimal = v.Optimal
	voice.Deprecated = v.Deprecated
	voice.Custom = v.Custom

	return
}

// GetVoiceRegionsBuilder [REST] Returns an array of voice region objects that can be used when creating servers.
//  Method                  GET
//  Endpoint                /voice/regions
//  Discord documentation   https://discord.com/developers/docs/resources/voice#list-voice-regions
//  Reviewed                2018-08-21
//  Comment                 -
func (c clientQueryBuilder) GetVoiceRegions(flags ...Flag) (regions []*VoiceRegion, err error) {
	r := c.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.VoiceRegions(),
		Ctx:      c.ctx,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*VoiceRegion, 0)
		return &tmp
	}

	var vs interface{}
	if vs, err = r.Execute(); err != nil {
		return nil, err
	}

	if ems, ok := vs.(*[]*VoiceRegion); ok {
		return *ems, nil
	}
	return vs.([]*VoiceRegion), nil
}
