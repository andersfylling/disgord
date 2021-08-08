package disgord

import (
	"context"
	"fmt"

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

// GetVoiceRegions [REST] Returns an array of voice region objects that can be used when creating servers.
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

// VoiceChannel is used to handle making a voice connection.
func (g guildQueryBuilder) VoiceChannel(channelID Snowflake) VoiceChannelQueryBuilder {
	vc := &voiceChannelQueryBuilder{}
	vc.gid = g.gid
	vc.cid = channelID
	vc.client = g.client
	vc.ctx = context.Background()
	return vc
}

type VoiceChannelQueryBuilder interface {
	WithContext(ctx context.Context) ChannelQueryBuilder

	// Get Get a channel by Snowflake. Returns a channel object.
	Get(flags ...Flag) (*Channel, error)

	// UpdateBuilder Update a Channels settings. Requires the 'MANAGE_CHANNELS' permission for the guild. Returns
	// a channel on success, and a 400 BAD REQUEST on invalid parameters. Fires a Channel Update Gateway event. If
	// modifying a category, individual Channel Update events will fire for each child channel that also changes.
	// For the PATCH method, all the JSON Params are optional.
	UpdateBuilder(flags ...Flag) *updateChannelBuilder

	// Delete Delete a channel, or close a private message. Requires the 'MANAGE_CHANNELS' permission for
	// the guild. Deleting a category does not delete its child Channels; they will have their parent_id removed and a
	// Channel Update Gateway event will fire for each of them. Returns a channel object on success.
	// Fires a Channel Delete Gateway event.
	Delete(flags ...Flag) (*Channel, error)

	// UpdatePermissions Edit the channel permission overwrites for a user or role in a channel. Only usable
	// for guild Channels. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty response on success.
	// For more information about permissions, see permissions.
	UpdatePermissions(overwriteID Snowflake, params *UpdateChannelPermissionsParams, flags ...Flag) error

	// GetInvites Returns a list of invite objects (with invite metadata) for the channel. Only usable for
	// guild Channels. Requires the 'MANAGE_CHANNELS' permission.
	GetInvites(flags ...Flag) ([]*Invite, error)

	// CreateInvite Create a new invite object for the channel. Only usable for guild Channels. Requires
	// the CREATE_INSTANT_INVITE permission. All JSON parameters for this route are optional, however the request
	// body is not. If you are not sending any fields, you still have to send an empty JSON object ({}).
	// Returns an invite object.
	CreateInvite(flags ...Flag) *createChannelInviteBuilder

	// DeletePermission Delete a channel permission overwrite for a user or role in a channel. Only usable
	// for guild Channels. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty response on success. For more
	// information about permissions,
	// see permissions: https://discord.com/developers/docs/topics/permissions#permissions
	DeletePermission(overwriteID Snowflake, flags ...Flag) error

	// Connect param{deaf} is deprecated
	Connect(mute, deaf bool) (VoiceConnection, error)

	JoinManual(mute, deaf bool) (*VoiceStateUpdate, *VoiceServerUpdate, error)
}

type voiceChannelQueryBuilder struct {
	channelQueryBuilder
	gid Snowflake
}

// Connect is used to handle making a voice connection.
func (v voiceChannelQueryBuilder) Connect(mute, deaf bool) (VoiceConnection, error) {
	return v.client.voiceConnectOptions(v.gid, v.cid, deaf, mute)
}

// JoinManual Tells Discord we want to start a new voice connection. However, it is the caller that
// is responsible for handling the voice connection compared to Connect(..).
//
// Useful for cases where third party libs, like Lavalink, is the preferred way to handle voice connections.
func (v voiceChannelQueryBuilder) JoinManual(mute, deaf bool) (*VoiceStateUpdate, *VoiceServerUpdate, error) {
	state := make(chan *VoiceStateUpdate, 2)
	stateCtrl := &Ctrl{Channel: state}
	defer stateCtrl.CloseChannel()

	server := make(chan *VoiceServerUpdate, 2)
	serverCtrl := &Ctrl{Channel: server}
	defer serverCtrl.CloseChannel()

	v.client.Gateway().WithCtrl(stateCtrl).VoiceStateUpdate(func(_ Session, evt *VoiceStateUpdate) {
		if evt.GuildID != v.gid {
			return
		}
		if evt.ChannelID != v.cid {
			return
		}
		state <- evt
	})

	v.client.Gateway().WithCtrl(serverCtrl).VoiceServerUpdate(func(_ Session, evt *VoiceServerUpdate) {
		if evt.GuildID != v.gid {
			return
		}
		server <- evt
	})

	_, err := v.client.Gateway().Dispatch(UpdateVoiceState, &UpdateVoiceStatePayload{
		GuildID:   v.gid,
		ChannelID: v.cid,
		SelfDeaf:  deaf,
		SelfMute:  mute,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send op 4 to discord: %w", err)
	}

	var (
		stateData  *VoiceStateUpdate
		serverData *VoiceServerUpdate
	)

	select {
	case serverData = <-server:
	case <-v.ctx.Done():
		return nil, nil, v.ctx.Err()
	}

	select {
	case stateData = <-state:
	case <-v.ctx.Done():
		return nil, serverData, v.ctx.Err()
	}

	return stateData, serverData, nil
}
