package voice

import "github.com/andersfylling/snowflake"

// State Voice State structure
// https://discordapp.com/developers/docs/resources/voice#voice-state-object
type State struct {
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

func (vst *State) Clear() {

}
