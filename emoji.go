package disgord

import (
	"errors"
	"net/http"
	"strings"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/ratelimit"
	"github.com/andersfylling/snowflake/v3"
)

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
	err = session.DeleteGuildEmoji(e.guildID, e.ID)
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

// ----------------------
// REST

// GetGuildEmojis [REST] Returns a list of emoji objects for the given guild.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/emojis
//  Rate limiter [MAJOR]    /guilds/{guild.id} // TODO: no idea if this is correct
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#list-guild-emojis
//  Reviewed                2018-06-10
//  Comment                 -
func (c *client) GetGuildEmojis(guildID snowflake.ID, flags ...Flag) (builder *listGuildEmojisBuilder) {
	builder = &listGuildEmojisBuilder{
		guildID: guildID,
	}
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      http.MethodGet,
		Ratelimiter: ratelimit.Guild(guildID),
		Endpoint:    endpoint.GuildEmojis(guildID),
	}, nil)

	return builder
}

type listGuildEmojisBuilder struct {
	r       RESTBuilder
	guildID snowflake.ID
}

func (b *listGuildEmojisBuilder) Execute() (emojis []*Emoji, err error) {
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

// GetGuildEmoji .
func (c *client) GetGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) (ret *Emoji, err error) {
	// TODO place emojis in their own cacheLink system
	var guild *Guild
	guild, err = c.cache.GetGuild(guildID)
	if err != nil {
		ret, err = GetGuildEmoji(c.req, guildID, emojiID)
		// TODO: cacheLink
		return
	}
	ret, err = guild.Emoji(emojiID)
	if err != nil {
		ret, err = GetGuildEmoji(c.req, guildID, emojiID)
		// TODO: cacheLink
		return
	}
	return
}

// CreateGuildEmoji .
func (c *client) CreateGuildEmoji(guildID Snowflake, params *CreateGuildEmojiParams, flags ...Flag) (ret *Emoji, err error) {
	ret, err = CreateGuildEmoji(c.req, guildID, params)
	return
}

// ModifyGuildEmoji .
func (c *client) ModifyGuildEmoji(guildID, emojiID Snowflake, params *ModifyGuildEmojiParams, flags ...Flag) (ret *Emoji, err error) {
	ret, err = ModifyGuildEmoji(c.req, guildID, emojiID, params)
	return
}

// DeleteGuildEmoji .
func (c *client) DeleteGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) (err error) {
	err = DeleteGuildEmoji(c.req, guildID, emojiID)
	return
}

// endpoints
//
// https://discordapp.com/developers/docs/resources/emoji#emoji-resource
// Routes for controlling emojis do not follow the normal rate limit conventions.
// These routes are specifically limited on a per-guild basis to prevent abuse.
// This means that the quota returned by our APIs may be inaccurate,
// and you may encounter 429s.

// GetGuildEmoji [REST] Returns an emoji object for the given guild and emoji IDs.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/emojis/{emoji.id}
//  Rate limiter [MAJOR]    /guilds/{guild.id} // TODO: no idea if this is correct
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#get-guild-emoji
//  Reviewed                2018-06-10
//  Comment                 -
func GetGuildEmoji(client httd.Getter, guildID, emojiID Snowflake) (ret *Emoji, err error) {
	var body []byte
	_, body, err = client.Get(&httd.Request{
		Ratelimiter: ratelimitGuild(guildID),
		Endpoint:    endpoint.GuildEmoji(guildID, emojiID),
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

func validEmojiName(name string) bool {
	// TODO: what is the allowed format?
	// a test showed that using "-" caused regex issues
	return !strings.Contains(name, "-")
}

// CreateGuildEmojiParams JSON params for func CreateGuildEmoji
type CreateGuildEmojiParams struct {
	Name  string      `json:"name"`
	Image string      `json:"image"`
	Roles []Snowflake `json:"roles"`
}

// CreateGuildEmoji [REST] Create a new emoji for the guild. Requires the 'MANAGE_EMOJIS' permission.
// Returns the new emoji object on success. Fires a Guild Emojis Update Gateway event.
//  Method                  POST
//  Endpoint                /guilds/{guild.id}/emojis
//  Rate limiter [MAJOR]    /guilds/{guild.id} // TODO: no idea if this is correct
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#create-guild-emoji
//  Reviewed                2018-06-10
//  Comment                 Emojis and animated emojis have a maximum file size of 256kb. Attempting to upload
//                          an emoji larger than this limit will fail and return 400 Bad Request and an
//                          error message, but not a JSON status code.
func CreateGuildEmoji(client httd.Poster, guildID Snowflake, params *CreateGuildEmojiParams) (ret *Emoji, err error) {
	if !validEmojiName(params.Name) {
		err = errors.New("emoji name contains illegal characters. Did not send request")
		return
	}
	var body []byte
	_, body, err = client.Post(&httd.Request{
		Ratelimiter: ratelimitGuild(guildID),
		Endpoint:    endpoint.GuildEmojis(guildID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// ModifyGuildEmojiParams JSON params for func ModifyGuildEmoji
type ModifyGuildEmojiParams struct {
	Name  string      `json:"name"`
	Roles []Snowflake `json:"roles"`
}

// ModifyGuildEmoji [REST] Modify the given emoji. Requires the 'MANAGE_EMOJIS' permission.
// Returns the updated emoji object on success. Fires a Guild Emojis Update Gateway event.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/emojis/{emoji.id}
//  Rate limiter [MAJOR]    /guilds/{guild.id} // TODO: no idea if this is correct
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#modify-guild-emoji
//  Reviewed                2018-06-10
//  Comment                 -
func ModifyGuildEmoji(client httd.Patcher, guildID, emojiID Snowflake, params *ModifyGuildEmojiParams) (ret *Emoji, err error) {
	if !validEmojiName(params.Name) {
		err = errors.New("emoji name contains illegal characters. Did not send request")
		return
	}
	var body []byte
	_, body, err = client.Patch(&httd.Request{
		Ratelimiter: ratelimitGuild(guildID),
		Endpoint:    endpoint.GuildEmoji(guildID, emojiID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// DeleteGuildEmoji [REST] Delete the given emoji. Requires the 'MANAGE_EMOJIS' permission. Returns 204 No Content on
// success. Fires a Guild Emojis Update Gateway event.
//  Method                  DELETE
//  Endpoint                /guilds/{guild.id}/emojis/{emoji.id}
//  Rate limiter [MAJOR]    /guilds/{guild.id} // TODO: no idea if this is correct
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#delete-guild-emoji
//  Reviewed                2018-06-10
//  Comment                 -
func DeleteGuildEmoji(client httd.Deleter, guildID, emojiID Snowflake) (err error) {
	var resp *http.Response
	resp, _, err = client.Delete(&httd.Request{
		Ratelimiter: ratelimitGuild(guildID),
		Endpoint:    endpoint.GuildEmoji(guildID, emojiID),
	})
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(http.StatusNoContent)
		err = errors.New(msg)
	}
	return
}
