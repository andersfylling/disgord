package disgord

import (
	"context"
	"errors"
	"fmt"
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

	//	image string // base 64 string, with prefix and everything
	guildID Snowflake
}

var _ Reseter = (*Emoji)(nil)
var _ DeepCopier = (*Emoji)(nil)
var _ Copier = (*Emoji)(nil)
var _ discordDeleter = (*Emoji)(nil)
var _ Mentioner = (*Emoji)(nil)

// var _ discordSaver = (*Emoji)(nil) // TODO
var _ fmt.Stringer = (*Emoji)(nil)

func (e *Emoji) String() string {
	return "emoji{name:" + e.Name + ", id:" + e.ID.String() + "}"
}

// PartialEmoji see Emoji
type PartialEmoji = Emoji

// SetBase64Image use this before creating the emoji for the first time
//func (e *Emoji) SetBase64Image(img string) {
//	e.image = img
//}

// Mention mentions an emoji. Adds the animation prefix, if animated
func (e *Emoji) Mention() string {
	prefix := ""
	if e.Animated {
		prefix = "a:"
	}

	return "<:" + prefix + e.Name + ":" + e.ID.String() + ">"
}

func (e *Emoji) LinkToGuild(guildID Snowflake) {
	e.guildID = guildID
}

// DeepCopy see interface at struct.go#DeepCopier
func (e *Emoji) DeepCopy() (copy interface{}) {
	copy = &Emoji{}
	e.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (e *Emoji) CopyOverTo(other interface{}) (err error) {
	var emoji *Emoji
	var ok bool
	if emoji, ok = other.(*Emoji); !ok {
		err = newErrorUnsupportedType("given type is not *Emoji")
		return
	}

	emoji.ID = e.ID
	emoji.Name = e.Name
	emoji.Roles = e.Roles
	emoji.RequireColons = e.RequireColons
	emoji.Managed = e.Managed
	emoji.Animated = e.Animated
	emoji.guildID = e.guildID

	if e.User != nil {
		emoji.User = e.User.DeepCopy().(*User)
	}
	return
}

// Missing GuildID...
//func (e *Emoji) saveToDiscord(s Session) (err error) {
//	session.Emoji
//}
func (e *Emoji) deleteFromDiscord(ctx context.Context, s Session, flags ...Flag) (err error) {
	if e.guildID.IsZero() {
		err = errors.New("missing guild ID, call Emoji.LinkToGuild")
		return
	}
	if e.ID.IsZero() {
		err = errors.New("missing emoji ID, cannot delete a not identified emoji")
		return
	}

	return s.DeleteGuildEmoji(ctx, e.guildID, e.ID, flags...)
}

//func (e *Emoji) createGuildEmoji(session Session) (err error) {
//	params := &CreateGuildEmojiParams{
//		Name:  e.Name,
//		Image: e.image,
//		Roles: e.Roles,
//	}
//
//	var creation *Emoji
//	creation, err = session.CreateGuildEmoji(, params)
//	if err != nil {
//		return
//	}
//
//	creation.CopyOverTo(e)
//	return
//}
//func (e *Emoji) modifyGuildEmoji(session Session) (err error) {
//	params := &ModifyGuildEmojiParams{
//		Name:  e.Name,
//		Roles: e.Roles,
//	}
//}

// func (e *Emoji) Clear() {
// 	// obviously don't delete the user ...
// }

// ----------------------
// CACHE

func cacheEmoji_EventGuildEmojisUpdate(cache Cacher, evt *GuildEmojisUpdate) error {
	return cacheEmoji_SetAll(cache, evt.GuildID, evt.Emojis)
}

func cacheEmoji_SetAll(cache Cacher, guildID Snowflake, emojis []*Emoji) error {
	cache.SetGuildEmojis(guildID, emojis)
	return nil
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

// GetGuildEmoji [REST] Returns an emoji object for the given guild and emoji IDs.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/emojis/{emoji.id}
//  Discord documentation   https://discord.com/developers/docs/resources/emoji#get-guild-emoji
//  Reviewed                2019-02-20
//  Comment                 -
func (c *Client) GetGuildEmoji(ctx context.Context, guildID, emojiID Snowflake, flags ...Flag) (*Emoji, error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildEmoji(guildID, emojiID),
		Ctx:      ctx,
	}, flags)
	r.pool = c.pool.emoji
	r.CacheRegistry = GuildEmojiCache
	r.preUpdateCache = func(x interface{}) {
		x.(*Emoji).guildID = guildID
	}
	r.factory = func() interface{} {
		return &Emoji{}
	}

	return getEmoji(r.Execute)
}

// GetGuildEmojis [REST] Returns a list of emoji objects for the given guild.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/emojis
//  Discord documentation   https://discord.com/developers/docs/resources/emoji#list-guild-emojis
//  Reviewed                2018-06-10
//  Comment                 -
func (c *Client) GetGuildEmojis(ctx context.Context, guildID Snowflake, flags ...Flag) (emojis []*Emoji, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildEmojis(guildID),
		Ctx:      ctx,
	}, flags)
	r.CacheRegistry = GuildEmojiCache
	r.checkCache = func() (v interface{}, err error) {
		if r.flags.Ignorecache() {
			return nil, nil
		}

		return c.cache.GetGuildEmojis(guildID)
	}
	r.factory = func() interface{} {
		tmp := make([]*Emoji, 0)
		return &tmp
	}
	r.preUpdateCache = func(x interface{}) {
		es := *x.(*[]*Emoji)
		for i := range es {
			es[i].guildID = guildID
		}
	}

	var vs interface{}
	if vs, err = r.Execute(); err != nil {
		return nil, err
	}

	if ems, ok := vs.(*[]*Emoji); ok {
		return *ems, nil
	}
	return vs.([]*Emoji), nil
}

// CreateGuildEmojiParams JSON params for func CreateGuildEmoji
type CreateGuildEmojiParams struct {
	Name  string      `json:"name"`  // required
	Image string      `json:"image"` // required
	Roles []Snowflake `json:"roles"` // optional

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string `json:"-"`
}

// CreateGuildEmoji [REST] Create a new emoji for the guild. Requires the 'MANAGE_EMOJIS' permission.
// Returns the new emoji object on success. Fires a Guild Emojis Update Gateway event.
//  Method                  POST
//  Endpoint                /guilds/{guild.id}/emojis
//  Discord documentation   https://discord.com/developers/docs/resources/emoji#create-guild-emoji
//  Reviewed                2019-02-20
//  Comment                 Emojis and animated emojis have a maximum file size of 256kb. Attempting to upload
//                          an emoji larger than this limit will fail and return 400 Bad Request and an
//                          error message, but not a JSON status code.
func (c *Client) CreateGuildEmoji(ctx context.Context, guildID Snowflake, params *CreateGuildEmojiParams, flags ...Flag) (emoji *Emoji, err error) {
	if guildID.IsZero() {
		return nil, errors.New("guildID must be set, was " + guildID.String())
	}

	if params == nil {
		return nil, errors.New("params object can not be nil")
	}
	if !validEmojiName(params.Name) {
		return nil, errors.New("invalid emoji name")
	}
	if !validAvatarPrefix(params.Image) {
		return nil, errors.New("image string must be base64 encoded with base64 prefix")
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Ctx:         ctx,
		Endpoint:    endpoint.GuildEmojis(guildID),
		ContentType: httd.ContentTypeJSON,
		Body:        params,
		Reason:      params.Reason,
	}, flags)
	r.CacheRegistry = GuildEmojiCache
	r.pool = c.pool.emoji
	r.factory = func() interface{} {
		return &Emoji{}
	}

	return getEmoji(r.Execute)
}

// UpdateGuildEmoji [REST] Modify the given emoji. Requires the 'MANAGE_EMOJIS' permission.
// Returns the updated emoji object on success. Fires a Guild Emojis Update Gateway event.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/emojis/{emoji.id}
//  Discord documentation   https://discord.com/developers/docs/resources/emoji#modify-guild-emoji
//  Reviewed                2019-02-20
//  Comment                 -
func (c *Client) UpdateGuildEmoji(ctx context.Context, guildID, emojiID Snowflake, flags ...Flag) (builder *updateGuildEmojiBuilder) {
	//if !validEmojiName(params.Name) {
	//	err = errors.New("emoji name contains illegal characters. Did not send request")
	//	return
	//}
	builder = &updateGuildEmojiBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Emoji{guildID: guildID}
	}
	builder.r.flags = flags
	builder.r.cacheRegistry = GuildEmojiCache
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         ctx,
		Endpoint:    endpoint.GuildEmoji(guildID, emojiID),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// DeleteGuildEmoji [REST] Delete the given emoji. Requires the 'MANAGE_EMOJIS' permission. Returns 204 No Content on
// success. Fires a Guild Emojis Update Gateway event.
//  Method                  DELETE
//  Endpoint                /guilds/{guild.id}/emojis/{emoji.id}
//  Discord documentation   https://discord.com/developers/docs/resources/emoji#delete-guild-emoji
//  Reviewed                2018-06-10
//  Comment                 -
func (c *Client) DeleteGuildEmoji(ctx context.Context, guildID, emojiID Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.GuildEmoji(guildID, emojiID),
		Ctx:      ctx,
	}, flags)
	r.updateCache = func(registry cacheRegistry, id Snowflake, x interface{}) (err error) {
		c.cache.DeleteGuildEmoji(guildID, emojiID)
		return nil
	}

	_, err = r.Execute()
	return
}

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

//generate-rest-params: name:string, roles:[]Snowflake,
//generate-rest-basic-execute: emoji:*Emoji,
type updateGuildEmojiBuilder struct {
	r RESTBuilder
}

//generate-rest-params: roles:[]Snowflake,
//generate-rest-basic-execute: emoji:*Emoji,
type createGuildEmojiBuilder struct {
	r RESTBuilder
}
