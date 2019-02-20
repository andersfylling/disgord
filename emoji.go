package disgord

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/ratelimit"
	"github.com/andersfylling/snowflake/v3"
)

func validEmojiName(name string) bool {
	// TODO: what is the allowed format?
	// a test showed that using "-" caused regex issues
	return !strings.Contains(name, "-")
}

// Emoji ...
type Emoji struct {
	mu Lockable

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

	return "<" + prefix + e.Name + ":" + e.ID.String() + ">"
}

func (e *Emoji) LinkToGuild(guildID snowflake.ID) {
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

	if constant.LockedMethods {
		e.mu.RLock()
		emoji.mu.Lock()
	}

	emoji.ID = e.ID
	emoji.Name = e.Name
	emoji.Roles = e.Roles
	emoji.RequireColons = e.RequireColons
	emoji.Managed = e.Managed
	emoji.Animated = e.Animated
	emoji.guildID = e.guildID
	emoji.mu = Lockable{}

	if e.User != nil {
		emoji.User = e.User.DeepCopy().(*User)
	}

	if constant.LockedMethods {
		e.mu.RUnlock()
		emoji.mu.Unlock()
	}

	return
}

// Missing GuildID...
//func (e *Emoji) saveToDiscord(session Session) (err error) {
//	session.Emoji
//}
func (e *Emoji) deleteFromDiscord(session Session) (err error) {
	if e.guildID.Empty() {
		err = errors.New("missing guild ID, call Emoji.LinkToGuild")
		return
	}
	if e.ID.Empty() {
		err = errors.New("missing emoji ID, cannot delete a not identified emoji")
		return
	}
	err = session.DeleteGuildEmoji(e.guildID, e.ID).Execute()
	return
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

func cacheEmoji_SetAll(cache Cacher, guildID snowflake.ID, emojis []*Emoji) error {
	cache.SetGuildEmojis(guildID, emojis)
	return nil
}

//////////////////////////////////////////////////////
//
// REST Methods
//
// https://discordapp.com/developers/docs/resources/emoji#emoji-resource
// Routes for controlling emojis do not follow the normal rate limit conventions.
// These routes are specifically limited on a per-guild basis to prevent abuse.
// This means that the quota returned by our APIs may be inaccurate,
// and you may encounter 429s.
//
//////////////////////////////////////////////////////

// GetGuildEmoji [REST] Returns an emoji object for the given guild and emoji IDs.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/emojis/{emoji.id}
//  Rate limiter [MAJOR]    /guilds/{guild.id}/emojis
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#get-guild-emoji
//  Reviewed                2019-02-20
//  Comment                 -
func (c *client) GetGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) (builder *getGuildEmojiBuilder) {
	builder = &getGuildEmojiBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Emoji{guildID: guildID}
	}
	builder.r.cacheRegistry = GuildEmojiCache
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      http.MethodGet,
		Ratelimiter: ratelimit.GuildEmojis(guildID),
		Endpoint:    endpoint.GuildEmoji(guildID, emojiID),
	}, nil)

	return builder
}

// GetGuildEmojis [REST] Returns a list of emoji objects for the given guild.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/emojis
//  Rate limiter [MAJOR]    /guilds/{guild.id}/emojis
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#list-guild-emojis
//  Reviewed                2018-06-10
//  Comment                 -
func (c *client) GetGuildEmojis(guildID snowflake.ID, flags ...Flag) (builder *getGuildEmojisBuilder) {
	builder = &getGuildEmojisBuilder{
		guildID: guildID,
	}
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      http.MethodGet,
		Ratelimiter: ratelimit.GuildEmojis(guildID),
		Endpoint:    endpoint.GuildEmojis(guildID),
	}, nil)

	return builder
}

// CreateGuildEmoji [REST] Create a new emoji for the guild. Requires the 'MANAGE_EMOJIS' permission.
// Returns the new emoji object on success. Fires a Guild Emojis Update Gateway event.
//  Method                  POST
//  Endpoint                /guilds/{guild.id}/emojis
//  Rate limiter [MAJOR]    /guilds/{guild.id}/emojis
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#create-guild-emoji
//  Reviewed                2019-02-20
//  Comment                 Emojis and animated emojis have a maximum file size of 256kb. Attempting to upload
//                          an emoji larger than this limit will fail and return 400 Bad Request and an
//                          error message, but not a JSON status code.
func (c *client) CreateGuildEmoji(guildID Snowflake, name, image string, flags ...Flag) (builder *createGuildEmojiBuilder) {
	builder = &createGuildEmojiBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Emoji{guildID: guildID}
	}
	builder.r.cacheRegistry = GuildEmojiCache
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      http.MethodPost,
		Ratelimiter: ratelimit.GuildEmojis(guildID),
		Endpoint:    endpoint.GuildEmojis(guildID),
		ContentType: httd.ContentTypeJSON,
	}, nil)
	builder.r.param("name", name)
	builder.r.param("image", image)
	builder.r.addPrereq(!validEmojiName(name), "invalid emoji name")
	builder.r.addPrereq(!validAvatarPrefix(image), "image string must be base64 encoded with base64 prefix")

	return builder
}

// ModifyGuildEmoji [REST] Modify the given emoji. Requires the 'MANAGE_EMOJIS' permission.
// Returns the updated emoji object on success. Fires a Guild Emojis Update Gateway event.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/emojis/{emoji.id}
//  Rate limiter [MAJOR]    /guilds/{guild.id}/emojis
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#modify-guild-emoji
//  Reviewed                2019-02-20
//  Comment                 -
func (c *client) ModifyGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) (builder *modifyGuildEmojiBuilder) {
	//if !validEmojiName(params.Name) {
	//	err = errors.New("emoji name contains illegal characters. Did not send request")
	//	return
	//}
	builder = &modifyGuildEmojiBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Emoji{guildID: guildID}
	}
	builder.r.cacheRegistry = GuildEmojiCache
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      http.MethodPatch,
		Ratelimiter: ratelimit.GuildEmojis(guildID),
		Endpoint:    endpoint.GuildEmoji(guildID, emojiID),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// DeleteGuildEmoji [REST] Delete the given emoji. Requires the 'MANAGE_EMOJIS' permission. Returns 204 No Content on
// success. Fires a Guild Emojis Update Gateway event.
//  Method                  DELETE
//  Endpoint                /guilds/{guild.id}/emojis/{emoji.id}
//  Rate limiter [MAJOR]    /guilds/{guild.id}/emojis
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#delete-guild-emoji
//  Reviewed                2018-06-10
//  Comment                 -
func (c *client) DeleteGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) (builder *basicBuilder) {
	builder = &basicBuilder{}
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      http.MethodDelete,
		Ratelimiter: ratelimitGuild(guildID),
		Endpoint:    endpoint.GuildEmoji(guildID, emojiID),
	}, func(resp *http.Response, body []byte, err error) error {
		if resp.StatusCode != http.StatusNoContent {
			msg := "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(http.StatusNoContent)
			err = errors.New(msg)
		}
		c.cache.DeleteGuildEmoji(guildID, emojiID)
		return nil
	})

	return builder
}

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

//generate-rest-params: name:string, roles:[]Snowflake,
//generate-rest-basic-execute: emoji:*Emoji,
type modifyGuildEmojiBuilder struct {
	r RESTBuilder
}

//generate-rest-params: roles:[]Snowflake,
//generate-rest-basic-execute: emoji:*Emoji,
type createGuildEmojiBuilder struct {
	r RESTBuilder
}

//generate-rest-basic-execute: emoji:*Emoji,
type getGuildEmojiBuilder struct {
	r RESTBuilder
}

type getGuildEmojisBuilder struct {
	r       RESTBuilder
	guildID snowflake.ID
}

func (b *getGuildEmojisBuilder) Execute() (emojis []*Emoji, err error) {
	if !b.r.ignoreCache {
		emojis, err = b.r.cache.GetGuildEmojis(b.guildID)
		if emojis != nil && err == nil {
			return
		}
	}

	b.r.prepare()
	var body []byte
	_, body, err = b.r.client.Request(b.r.config)
	if err != nil {
		return
	}

	if len(body) > 1 {
		err = httd.Unmarshal(body, &emojis)
		if err != nil {
			return
		}

		for i := range emojis {
			emojis[i].guildID = b.guildID
		}
		b.r.cache.SetGuildEmojis(b.guildID, emojis)
	}
	return
}
