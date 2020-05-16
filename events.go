package disgord

//go:generate go run generate/events/main.go

// This file contains resource objects for the event reactor

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/andersfylling/disgord/internal/util"
)

// Resource represents a discord event.
// This is used internally for readability only.
type resource = interface{}

func cacheEvent(cache Cacher, event string, v interface{}, data json.RawMessage) (err error) {
	// updates holds key and object to be cached
	updates := map[cacheRegistry]([]interface{}){}

	switch event {
	case EvtReady:
		ready := v.(*Ready)
		updates[UserCache] = append(updates[UserCache], ready.User)

		for _, guild := range ready.Guilds {
			updates[GuildCache] = append(updates[GuildCache], guild)
		}
	case EvtVoiceStateUpdate:
		update := v.(*VoiceStateUpdate)
		updates[VoiceStateCache] = append(updates[VoiceStateCache], update.VoiceState)
	case EvtChannelCreate, EvtChannelUpdate:
		var channel *Channel
		if event == EvtChannelCreate {
			channel = (v.(*ChannelCreate)).Channel
			if !channel.GuildID.IsZero() {
				cache.AddGuildChannel(channel.GuildID, channel.ID)
			}
		} else if event == EvtChannelUpdate {
			channel = (v.(*ChannelUpdate)).Channel
		}
		if len(channel.Recipients) > 0 {
			for i := range channel.Recipients {
				updates[UserCache] = append(updates[UserCache], channel.Recipients[i])
			}
		}

		updates[ChannelCache] = append(updates[ChannelCache], channel)
	case EvtChannelDelete:
		channel := (v.(*ChannelDelete)).Channel
		cache.DeleteChannel(channel.ID)
		cache.DeleteGuildChannel(channel.GuildID, channel.ID)
	case EvtChannelPinsUpdate:
		evt := v.(*ChannelPinsUpdate)
		cache.UpdateChannelPin(evt.ChannelID, evt.LastPinTimestamp)
	case EvtGuildCreate, EvtGuildUpdate:
		var guild *Guild
		if event == EvtGuildCreate {
			guild = (v.(*GuildCreate)).Guild
		} else if event == EvtGuildUpdate {
			guild = (v.(*GuildUpdate)).Guild
		}
		updates[GuildCache] = append(updates[GuildCache], guild)

		// update all users
		if len(guild.Members) > 0 {
			updates[UserCache] = make([]interface{}, len(guild.Members))
			for i := range guild.Members {
				updates[UserCache][i] = guild.Members[i].User
			}
		}
		// update all channels
		if len(guild.Channels) > 0 {
			updates[ChannelCache] = make([]interface{}, len(guild.Channels))
			for i := range guild.Channels {
				updates[ChannelCache][i] = guild.Channels[i]
			}
		}
	case EvtGuildDelete:
		uguild := (v.(*GuildDelete)).UnavailableGuild
		cache.DeleteGuild(uguild.ID)
	case EvtGuildRoleDelete:
		evt := v.(*GuildRoleDelete)
		cache.DeleteGuildRole(evt.GuildID, evt.RoleID)
	case EvtGuildEmojisUpdate:
		err = cacheEmoji_EventGuildEmojisUpdate(cache, v.(*GuildEmojisUpdate))
	case EvtUserUpdate:
		usr := v.(*UserUpdate).User
		updates[UserCache] = append(updates[UserCache], usr)
	case EvtMessageCreate:
		// TODO: performance issues?
		msg := (v.(*MessageCreate)).Message
		cache.UpdateChannelLastMessageID(msg.ChannelID, msg.ID)
	case EvtGuildMembersChunk:
		evt := v.(*GuildMembersChunk)
		updates[GuildMembersCache] = append(updates[GuildMembersCache], evt)

		// update all users
		if len(evt.Members) > 0 {
			updates[UserCache] = make([]interface{}, len(evt.Members))
			for i := range evt.Members {
				updates[UserCache][i] = evt.Members[i].User
			}
		}
	case EvtGuildMemberUpdate:
		evt := v.(*GuildMemberUpdate)
		cache.UpdateMemberAndUser(evt.GuildID, evt.User.ID, data)
	case EvtGuildMemberAdd:
		evt := v.(*GuildMemberAdd)
		cache.AddGuildMember(evt.Member.GuildID, evt.Member)
		updates[UserCache] = append(updates[UserCache], evt.Member.User)
	case EvtGuildMemberRemove:
		evt := v.(*GuildMemberRemove)
		cache.RemoveGuildMember(evt.GuildID, evt.User.ID)
	// TODO: mark user as free from guild...
	case EvtGuildRoleCreate:
		evt := v.(*GuildRoleCreate)
		cache.AddGuildRole(evt.GuildID, evt.Role)
	case EvtGuildRoleUpdate:
		evt := v.(*GuildRoleUpdate)
		if updated := cache.UpdateGuildRole(evt.GuildID, evt.Role, data); !updated {
			cache.AddGuildRole(evt.GuildID, evt.Role)
		}
	default:
		//case EventResumed:
		//case EventGuildBanAdd:
		//case EventGuildBanRemove:
		//case EventGuildIntegrationsUpdate:
		//case EventMessageUpdate:
		//case EventMessageDelete:
		//case EventMessageDeleteBulk:
		//case EventMessageReactionAdd:
		//case EventMessageReactionRemove:
		//case EventMessageReactionRemoveAll:
		//case EventPresenceUpdate:
		//case EventTypingStart:
		//case EventVoiceServerUpdate:
		//case EventWebhooksUpdate:
	}

	for key, structs := range updates {
		if err = cache.Updates(key, structs); err != nil {
			// TODO: logging? or append all errs to the return statement?
		}
	}
	return nil
}

// ---------------------------

type evtResource interface {
	registerContext(ctx context.Context)
	setShardID(id uint)
}

// ---------------------------

// Ready contains the initial state information
type Ready struct {
	APIVersion int                 `json:"v"`
	User       *User               `json:"user"`
	Guilds     []*GuildUnavailable `json:"guilds"`

	// not really needed, as it is handled on the socket layer.
	SessionID string `json:"session_id"`

	// private_channels will be an empty array. As bots receive private messages,
	// they will be notified via Channel Create events.
	//PrivateChannels []*channel.Channel `json:"private_channels"`

	// bot can't have presences
	//Presences []*Presence         `json:"presences"`

	// bot cant have relationships
	//Relationships []interface{} `son:"relationships"`

	// bot can't have user settings
	// UserSettings interface{}        `json:"user_settings"`

	sync.RWMutex `json:"-"`
	Ctx          context.Context `json:"-"`
	ShardID      uint            `json:"-"`
}

// ---------------------------

// Resumed response to Resume
type Resumed struct {
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

// ---------------------------

// ChannelCreate new channel created
type ChannelCreate struct {
	Channel *Channel        `json:"channel"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

// UnmarshalJSON ...
func (obj *ChannelCreate) UnmarshalJSON(data []byte) error {
	obj.Channel = &Channel{}
	return unmarshal(data, obj.Channel)
}

// ---------------------------

// ChannelUpdate channel was updated
type ChannelUpdate struct {
	Channel *Channel        `json:"channel"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

// UnmarshalJSON ...
func (obj *ChannelUpdate) UnmarshalJSON(data []byte) error {
	obj.Channel = &Channel{}
	return unmarshal(data, obj.Channel)
}

// ---------------------------

// ChannelDelete channel was deleted
type ChannelDelete struct {
	Channel *Channel        `json:"channel"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

// UnmarshalJSON ...
func (obj *ChannelDelete) UnmarshalJSON(data []byte) error {
	obj.Channel = &Channel{}
	return unmarshal(data, obj.Channel)
}

// ---------------------------

// ChannelPinsUpdate message was pinned or unpinned
type ChannelPinsUpdate struct {
	// ChannelID snowflake	the id of the channel
	ChannelID Snowflake `json:"channel_id"`

	// LastPinTimestamp	ISO8601 timestamp	the time at which the most recent pinned message was pinned
	LastPinTimestamp Time            `json:"last_pin_timestamp,omitempty"` // ?|
	Ctx              context.Context `json:"-"`
	ShardID          uint            `json:"-"`
}

// ---------------------------

// TypingStart user started typing in a channel
type TypingStart struct {
	ChannelID     Snowflake       `json:"channel_id"`
	UserID        Snowflake       `json:"user_id"`
	TimestampUnix int             `json:"timestamp"`
	Ctx           context.Context `json:"-"`
	ShardID       uint            `json:"-"`
}

// ---------------------------

// InviteDelete Sent when an invite is deleted.
type InviteDelete struct {
	ChannelID Snowflake       `json:"channel_id"`
	GuildID   Snowflake       `json:"guild_id"`
	Code      string          `json:"code"`
	Ctx       context.Context `json:"-"`
	ShardID   uint            `json:"-"`
}

// ---------------------------

// MessageCreate message was created
type MessageCreate struct {
	Message *Message
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

var _ Reseter = (*MessageCreate)(nil)
var _ internalUpdater = (*MessageCreate)(nil)

func (obj *MessageCreate) updateInternals() {
	obj.Message.updateInternals()
}

// UnmarshalJSON ...
func (obj *MessageCreate) UnmarshalJSON(data []byte) error {
	obj.Message = &Message{}
	return unmarshal(data, obj.Message)
}

// ---------------------------

// MessageUpdate message was edited
type MessageUpdate struct {
	Message *Message
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

var _ internalUpdater = (*MessageUpdate)(nil)

func (obj *MessageUpdate) updateInternals() {
	obj.Message.updateInternals()
}

// UnmarshalJSON ...
func (obj *MessageUpdate) UnmarshalJSON(data []byte) error {
	obj.Message = &Message{}
	return unmarshal(data, obj.Message)
}

// ---------------------------

// MessageDelete message was deleted
type MessageDelete struct {
	MessageID Snowflake       `json:"id"`
	ChannelID Snowflake       `json:"channel_id"`
	GuildID   Snowflake       `json:"guild_id,omitempty"`
	Ctx       context.Context `json:"-"`
	ShardID   uint            `json:"-"`
}

// ---------------------------

// MessageDeleteBulk multiple messages were deleted at once
type MessageDeleteBulk struct {
	MessageIDs []Snowflake     `json:"ids"`
	ChannelID  Snowflake       `json:"channel_id"`
	Ctx        context.Context `json:"-"`
	ShardID    uint            `json:"-"`
}

// ---------------------------

// MessageReactionAdd user reacted to a message
// Note! do not cache emoji, unless it's updated with guildID
// TODO: find guildID when given userID, ChannelID and MessageID
type MessageReactionAdd struct {
	UserID    Snowflake `json:"user_id"`
	ChannelID Snowflake `json:"channel_id"`
	MessageID Snowflake `json:"message_id"`
	// PartialEmoji id and name. id might be nil
	PartialEmoji *Emoji          `json:"emoji"`
	Ctx          context.Context `json:"-"`
	ShardID      uint            `json:"-"`
}

// ---------------------------

// MessageReactionRemove user removed a reaction from a message
// Note! do not cache emoji, unless it's updated with guildID
// TODO: find guildID when given userID, ChannelID and MessageID
type MessageReactionRemove struct {
	UserID    Snowflake `json:"user_id"`
	ChannelID Snowflake `json:"channel_id"`
	MessageID Snowflake `json:"message_id"`
	// PartialEmoji id and name. id might be nil
	PartialEmoji *Emoji          `json:"emoji"`
	Ctx          context.Context `json:"-"`
	ShardID      uint            `json:"-"`
}

// ---------------------------

// MessageReactionRemoveAll all reactions were explicitly removed from a message
type MessageReactionRemoveAll struct {
	ChannelID Snowflake       `json:"channel_id"`
	MessageID Snowflake       `json:"message_id"`
	Ctx       context.Context `json:"-"`
	ShardID   uint            `json:"-"`
}

// ---------------------------

// GuildEmojisUpdate guild emojis were updated
type GuildEmojisUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Emojis  []*Emoji        `json:"emojis"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

var _ internalUpdater = (*GuildEmojisUpdate)(nil)

func (g *GuildEmojisUpdate) updateInternals() {
	for i := range g.Emojis {
		g.Emojis[i].guildID = g.GuildID
	}
}

// ---------------------------

// GuildCreate This event can be sent in three different scenarios:
//  1. When a user is initially connecting, to lazily load and backfill information for all unavailable guilds
//     sent in the Ready event.
//	2. When a Guild becomes available again to the Client.
// 	3. When the current user joins a new Guild.
type GuildCreate struct {
	Guild   *Guild          `json:"guild"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

var _ internalUpdater = (*GuildCreate)(nil)

func (g *GuildCreate) updateInternals() {
	g.Guild.updateInternals()
}

// UnmarshalJSON ...
func (obj *GuildCreate) UnmarshalJSON(data []byte) error {
	obj.Guild = &Guild{}
	return unmarshal(data, obj.Guild)
}

// ---------------------------

// GuildUpdate guild was updated
type GuildUpdate struct {
	Guild   *Guild          `json:"guild"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

var _ internalUpdater = (*GuildUpdate)(nil)

func (g *GuildUpdate) updateInternals() {
	g.Guild.updateInternals()
}

// UnmarshalJSON ...
func (obj *GuildUpdate) UnmarshalJSON(data []byte) error {
	obj.Guild = &Guild{}
	return unmarshal(data, obj.Guild)
}

// ---------------------------

// GuildDelete guild became unavailable, or user left/was removed from a guild
type GuildDelete struct {
	UnavailableGuild *GuildUnavailable `json:"guild_unavailable"`
	Ctx              context.Context   `json:"-"`
	ShardID          uint              `json:"-"`
}

// UserWasRemoved ... TODO
func (obj *GuildDelete) UserWasRemoved() bool {
	return obj.UnavailableGuild.Unavailable == false
}

// UnmarshalJSON ...
func (obj *GuildDelete) UnmarshalJSON(data []byte) error {
	obj.UnavailableGuild = &GuildUnavailable{}
	return unmarshal(data, obj.UnavailableGuild)
}

// ---------------------------

// GuildBanAdd user was banned from a guild
type GuildBanAdd struct {
	GuildID Snowflake       `json:"guild_id"`
	User    *User           `json:"user"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

// ---------------------------

// GuildBanRemove user was unbanned from a guild
type GuildBanRemove struct {
	GuildID Snowflake       `json:"guild_id"`
	User    *User           `json:"user"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

// ---------------------------

// GuildIntegrationsUpdate guild integration was updated
type GuildIntegrationsUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

// ---------------------------

// GuildMemberAdd new user joined a guild
type GuildMemberAdd struct {
	Member  *Member         `json:"member"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

var _ internalUpdater = (*GuildMemberAdd)(nil)

func (g *GuildMemberAdd) updateInternals() {
	g.Member.updateInternals()
}

// UnmarshalJSON ...
func (obj *GuildMemberAdd) UnmarshalJSON(data []byte) error {
	obj.Member = &Member{}
	return util.Unmarshal(data, obj.Member)
}

// ---------------------------

// GuildMemberRemove user was removed from a guild
type GuildMemberRemove struct {
	GuildID Snowflake       `json:"guild_id"`
	User    *User           `json:"user"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

// ---------------------------

// GuildMemberUpdate guild member was updated
type GuildMemberUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Roles   []Snowflake     `json:"roles"`
	User    *User           `json:"user"`
	Nick    string          `json:"nick"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

// ---------------------------

// GuildMembersChunk response to Request Guild Members
type GuildMembersChunk struct {
	GuildID Snowflake       `json:"guild_id"`
	Members []*Member       `json:"members"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

var _ internalUpdater = (*GuildMembersChunk)(nil)

func (g *GuildMembersChunk) updateInternals() {
	for i := range g.Members {
		g.Members[i].updateInternals()
	}
}

// ---------------------------

// GuildRoleCreate guild role was created
type GuildRoleCreate struct {
	GuildID Snowflake       `json:"guild_id"`
	Role    *Role           `json:"role"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

var _ internalUpdater = (*GuildRoleCreate)(nil)

func (g *GuildRoleCreate) updateInternals() {
	g.Role.guildID = g.GuildID
}

// ---------------------------

// GuildRoleUpdate guild role was updated
type GuildRoleUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Role    *Role           `json:"role"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

var _ internalUpdater = (*GuildRoleUpdate)(nil)

func (g *GuildRoleUpdate) updateInternals() {
	g.Role.guildID = g.GuildID
}

// ---------------------------

// GuildRoleDelete a guild role was deleted
type GuildRoleDelete struct {
	GuildID Snowflake       `json:"guild_id"`
	RoleID  Snowflake       `json:"role_id"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

// ---------------------------

// PresenceUpdate user's presence was updated in a guild
type PresenceUpdate struct {
	User    *User       `json:"user"`
	RoleIDs []Snowflake `json:"roles"`
	Game    *Activity   `json:"game"`
	GuildID Snowflake   `json:"guild_id"`

	// Status either "idle", "dnd", "online", or "offline"
	// TODO: constants somewhere..
	Status  string          `json:"status"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

// ---------------------------

// UserUpdate properties about a user changed
type UserUpdate struct {
	User    *User           `json:"user"`
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

// ---------------------------

// VoiceStateUpdate someone joined, left, or moved a voice channel
type VoiceStateUpdate struct {
	*VoiceState
	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}

// ---------------------------

// VoiceServerUpdate guild's voice server was updated. Sent when a guild's voice server is updated. This is sent when initially
// connecting to voice, and when the current voice instance fails over to a new server.
type VoiceServerUpdate struct {
	Token    string          `json:"token"`
	GuildID  Snowflake       `json:"guild_id"`
	Endpoint string          `json:"endpoint"`
	Ctx      context.Context `json:"-"`
	ShardID  uint            `json:"-"`
}

// ---------------------------

// WebhooksUpdate guild channel webhook was created, update, or deleted
type WebhooksUpdate struct {
	GuildID   Snowflake       `json:"guild_id"`
	ChannelID Snowflake       `json:"channel_id"`
	Ctx       context.Context `json:"-"`
	ShardID   uint            `json:"-"`
}

// InviteCreate guild invite was created
type InviteCreate struct {
	// Code the invite code (unique Snowflake)
	Code string `json:"code"`

	// Guild the guild this invite is for
	Guild *PartialGuild `json:"guild"`

	// Channel the channel this invite is for
	Channel *PartialChannel `json:"channel"`

	// Inviter the user that created the invite
	Inviter *User `json:"inviter"`

	// CreatedAt the time at which the invite was created
	CreatedAt Time `json:"created_at"`

	// MaxAge how long the invite is valid for (in seconds)
	MaxAge int `json:"max_age"`

	// MaxUses the maximum number of times the invite can be used
	MaxUses int `json:"max_uses"`

	// Temporary whether or not the invite is temporary (invited users will be kicked on disconnect unless they're assigned a role)
	Temporary bool `json:"temporary"`

	// Uses how many times the invite has been used (always will be 0)
	Uses int `json:"uses"`

	Revoked bool `json:"revoked"`
	Unique  bool `json:"unique"`

	// ApproximatePresenceCount approximate count of online members
	ApproximatePresenceCount int `json:"approximate_presence_count,omitempty"`

	// ApproximatePresenceCount approximate count of total members
	ApproximateMemberCount int `json:"approximate_member_count,omitempty"`

	Ctx     context.Context `json:"-"`
	ShardID uint            `json:"-"`
}
