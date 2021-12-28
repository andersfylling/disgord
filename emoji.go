package disgord

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
)

func validEmojiName(name string) bool {
	if name == "" {
		return false
	}
	// TODO: what is the allowed format?
	// a test showed that using "-" caused regex issues
	return !strings.Contains(name, "-")
}

// Emoji ...
type Emoji struct {
	ID            Snowflake   `json:"id"`
	Name          string      `json:"name"`
	Roles         []Snowflake `json:"roles,omitempty"`
	User          *User       `json:"user,omitempty"` // the user who created the emoji
	RequireColons bool        `json:"require_colons,omitempty"`
	Managed       bool        `json:"managed,omitempty"`
	Animated      bool        `json:"animated,omitempty"`
	Available     bool        `json:"available,omitempty"`
}

var _ Reseter = (*Emoji)(nil)
var _ DeepCopier = (*Emoji)(nil)
var _ Copier = (*Emoji)(nil)
var _ Mentioner = (*Emoji)(nil)

// var _ discordSaver = (*Emoji)(nil) // TODO
var _ fmt.Stringer = (*Emoji)(nil)

func (e *Emoji) String() string {
	return "emoji{name:" + e.Name + ", id:" + e.ID.String() + "}"
}

// PartialEmoji see Emoji
type PartialEmoji = Emoji

// Mention mentions an emoji. Adds the animation prefix, if animated
func (e *Emoji) Mention() string {
	prefix := ":"
	if e.Animated {
		prefix = "a:"
	}

	return "<" + prefix + e.Name + ":" + e.ID.String() + ">"
}

//////////////////////////////////////////////////////
//
// REST Methods
//
// https://discord.com/developers/docs/resources/emoji#emoji-resource
// Routes for controlling emojis do not follow the normal rate limit conventions.
// These routes are specifically limited on a per-guild basis to prevent abuse.
// This means that the quota returned by our APIs may be inaccurate,
// and you may encounter 429s.
//
//////////////////////////////////////////////////////

type GuildEmojiQueryBuilder interface {
	WithContext(ctx context.Context) GuildEmojiQueryBuilder
	WithFlags(flags ...Flag) GuildEmojiQueryBuilder

	Get() (*Emoji, error)
	Update(params *UpdateEmoji) (*Emoji, error)
	Delete() error

	// Deprecated: use Update
	UpdateBuilder() UpdateGuildEmojiBuilder
}

func (g guildQueryBuilder) Emoji(emojiID Snowflake) GuildEmojiQueryBuilder {
	return &guildEmojiQueryBuilder{client: g.client, gid: g.gid, emojiID: emojiID}
}

type guildEmojiQueryBuilder struct {
	ctx     context.Context
	flags   Flag
	client  *Client
	gid     Snowflake
	emojiID Snowflake
}

func (g *guildEmojiQueryBuilder) validate() error {
	if g.client == nil {
		return MissingClientInstanceErr
	}
	if g.gid.IsZero() {
		return MissingGuildIDErr
	}
	if g.emojiID.IsZero() {
		return MissingEmojiIDErr
	}
	return nil
}

func (g guildEmojiQueryBuilder) WithContext(ctx context.Context) GuildEmojiQueryBuilder {
	g.ctx = ctx
	return &g
}

func (g guildEmojiQueryBuilder) WithFlags(flags ...Flag) GuildEmojiQueryBuilder {
	g.flags = mergeFlags(flags)
	return &g
}

func (g guildEmojiQueryBuilder) Get() (*Emoji, error) {
	if !ignoreCache(g.flags) {
		if emoji, _ := g.client.cache.GetGuildEmoji(g.gid, g.emojiID); emoji != nil {
			return emoji, nil
		}
	}

	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildEmoji(g.gid, g.emojiID),
		Ctx:      g.ctx,
	}, g.flags)
	r.pool = g.client.pool.emoji
	r.factory = func() interface{} {
		return &Emoji{}
	}

	return getEmoji(r.Execute)
}

// Update Modify the given emoji. Requires the 'MANAGE_EMOJIS' permission.
// Returns the updated emoji object on success. Fires a Guild Emojis Update Gateway event.
func (g guildEmojiQueryBuilder) Update(params *UpdateEmoji) (*Emoji, error) {
	if params == nil {
		return nil, MissingRESTParamsErr
	}
	if err := g.validate(); err != nil {
		return nil, err
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildEmoji(g.gid, g.emojiID),
		ContentType: httd.ContentTypeJSON,
		Body:        params,
		Reason:      params.AuditLogReason,
	}, g.flags)
	r.pool = g.client.pool.emoji
	r.factory = func() interface{} {
		return &Emoji{}
	}

	return getEmoji(r.Execute)
}

type UpdateEmoji struct {
	Name  *string      `json:"name,omitempty"`
	Roles *[]Snowflake `json:"roles,omitempty"`

	AuditLogReason string `json:"-"`
}

// Delete deletes the given emoji. Requires the 'MANAGE_EMOJIS' permission. Returns 204 No Content on
// success. Fires a Guild Emojis Update Gateway event.
func (g guildEmojiQueryBuilder) Delete() (err error) {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   http.MethodDelete,
		Endpoint: endpoint.GuildEmoji(g.gid, g.emojiID),
		Ctx:      g.ctx,
	}, g.flags)

	_, err = r.Execute()
	return
}
