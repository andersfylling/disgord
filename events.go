package disgord

//go:generate go run generate/events/main.go

import (
	"context"
	"sync"
)

func prepareBox(evtName string, box interface{}) {
	switch evtName {
	case EventGuildCreate:
		guild := (box.(*GuildCreate)).Guild
		for _, role := range guild.Roles {
			role.guildID = guild.ID
		}
	case EventGuildUpdate:
		guild := (box.(*GuildUpdate)).Guild
		for _, role := range guild.Roles {
			role.guildID = guild.ID
		}
	case EventGuildRoleCreate:
		(box.(*GuildRoleCreate)).Role.guildID = (box.(*GuildRoleCreate)).GuildID
	case EventGuildRoleUpdate:
		(box.(*GuildRoleUpdate)).Role.guildID = (box.(*GuildRoleUpdate)).GuildID
	case EventGuildEmojisUpdate:
		evt := box.(*GuildEmojisUpdate)
		for i := range evt.Emojis {
			evt.Emojis[i].guildID = evt.GuildID
		}
	}
}

func cacheEvent(cache Cacher, event string, v interface{}) (err error) {
	// updates holds key and object to be cached
	updates := map[cacheRegistry]([]interface{}){}

	switch event {
	case EventReady:
		ready := v.(*Ready)
		updates[UserCache] = append(updates[UserCache], ready.User)

		for _, guild := range ready.Guilds {
			updates[GuildCache] = append(updates[GuildCache], guild)
		}
	case EventVoiceStateUpdate:
		update := v.(*VoiceStateUpdate)
		updates[VoiceStateCache] = append(updates[VoiceStateCache], update.VoiceState)
	case EventChannelCreate, EventChannelUpdate:
		var channel *Channel
		if event == EventChannelCreate {
			channel = (v.(*ChannelCreate)).Channel
			if !channel.GuildID.Empty() {
				cache.AddGuildChannel(channel.GuildID, channel.ID)
			}
		} else if event == EventChannelUpdate {
			channel = (v.(*ChannelUpdate)).Channel
		}
		if len(channel.Recipients) > 0 {
			for i := range channel.Recipients {
				updates[UserCache] = append(updates[UserCache], channel.Recipients[i])
			}
		}

		updates[ChannelCache] = append(updates[ChannelCache], channel)
	case EventChannelDelete:
		channel := (v.(*ChannelDelete)).Channel
		cache.DeleteChannel(channel.ID)
		cache.DeleteGuildChannel(channel.GuildID, channel.ID)
	case EventChannelPinsUpdate:
		evt := v.(*ChannelPinsUpdate)
		cache.UpdateChannelPin(evt.ChannelID, evt.LastPinTimestamp)
	case EventGuildCreate, EventGuildUpdate:
		var guild *Guild
		if event == EventGuildCreate {
			guild = (v.(*GuildCreate)).Guild
		} else if event == EventGuildUpdate {
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
	case EventGuildDelete:
		uguild := (v.(*GuildDelete)).UnavailableGuild
		cache.DeleteGuild(uguild.ID)
	case EventGuildRoleDelete:
		evt := v.(*GuildRoleDelete)
		cache.DeleteGuildRole(evt.GuildID, evt.RoleID)
	case EventGuildEmojisUpdate:
		err = cacheEmoji_EventGuildEmojisUpdate(cache, v.(*GuildEmojisUpdate))
	case EventUserUpdate:
		usr := v.(*UserUpdate).User
		updates[UserCache] = append(updates[UserCache], usr)
	case EventMessageCreate:
		// TODO: performance issues?
		msg := (v.(*MessageCreate)).Message
		cache.UpdateChannelLastMessageID(msg.ChannelID, msg.ID)
	case EventGuildMembersChunk:
		evt := v.(*GuildMembersChunk)
		updates[GuildMembersCache] = append(updates[GuildMembersCache], evt)

		// update all users
		if len(evt.Members) > 0 {
			updates[UserCache] = make([]interface{}, len(evt.Members))
			for i := range evt.Members {
				updates[UserCache][i] = evt.Members[i].User
			}
		}
	default:
		//case EventResumed:
		//case EventGuildBanAdd:
		//case EventGuildBanRemove:
		//case EventGuildIntegrationsUpdate:
		//case EventGuildMemberAdd:
		//case EventGuildMemberRemove:
		//case EventGuildMemberUpdate:
		//case EventGuildRoleCreate:
		//case EventGuildRoleUpdate:
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
			// TODO: logging
		}
	}
	return
}

// ---------------------------

type eventBox interface {
	registerContext(ctx context.Context)
}

// ---------------------------

// PresencesReplace holds the event content
type PresencesReplace struct {
	Presnces []*PresenceUpdate `json:"presences_replace"` // TODO: verify json tag
	Ctx      context.Context   `json:"-"`
}

// ---------------------------

// Ready contains the initial state information
type Ready struct {
	APIVersion int                 `json:"v"`
	User       *User               `json:"user"`
	Guilds     []*GuildUnavailable `json:"guilds"`

	// not really needed, as it is handled on the socket layer.
	SessionID string   `json:"session_id"`
	Trace     []string `json:"_trace"`

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
}

// ---------------------------

// Resumed response to Resume
type Resumed struct {
	Trace []string        `json:"_trace"`
	Ctx   context.Context `json:"-"`
}

// ---------------------------

// ChannelCreate new channel created
type ChannelCreate struct {
	Channel *Channel        `json:"channel"`
	Ctx     context.Context `json:"-"`
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
	LastPinTimestamp Timestamp       `json:"last_pin_timestamp,omitempty"` // ?|
	Ctx              context.Context `json:"-"`
}

// ---------------------------

// TypingStart user started typing in a channel
type TypingStart struct {
	ChannelID     Snowflake       `json:"channel_id"`
	UserID        Snowflake       `json:"user_id"`
	TimestampUnix int             `json:"timestamp"`
	Ctx           context.Context `json:"-"`
}

// ---------------------------

// MessageCreate message was created
type MessageCreate struct {
	Message *Message
	Ctx     context.Context `json:"-"`
}

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
}

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
}

// ---------------------------

// MessageDeleteBulk multiple messages were deleted at once
type MessageDeleteBulk struct {
	MessageIDs []Snowflake     `json:"ids"`
	ChannelID  Snowflake       `json:"channel_id"`
	Ctx        context.Context `json:"-"`
}

// ---------------------------

// MessageReactionAdd user reacted to a message
type MessageReactionAdd struct {
	UserID    Snowflake `json:"user_id"`
	ChannelID Snowflake `json:"channel_id"`
	MessageID Snowflake `json:"message_id"`
	// PartialEmoji id and name. id might be nil
	PartialEmoji *Emoji          `json:"emoji"`
	Ctx          context.Context `json:"-"`
}

// ---------------------------

// MessageReactionRemove user removed a reaction from a message
type MessageReactionRemove struct {
	UserID    Snowflake `json:"user_id"`
	ChannelID Snowflake `json:"channel_id"`
	MessageID Snowflake `json:"message_id"`
	// PartialEmoji id and name. id might be nil
	PartialEmoji *Emoji          `json:"emoji"`
	Ctx          context.Context `json:"-"`
}

// ---------------------------

// MessageReactionRemoveAll all reactions were explicitly removed from a message
type MessageReactionRemoveAll struct {
	ChannelID Snowflake       `json:"channel_id"`
	MessageID Snowflake       `json:"id"`
	Ctx       context.Context `json:"-"`
}

// ---------------------------

// GuildEmojisUpdate guild emojis were updated
type GuildEmojisUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Emojis  []*Emoji        `json:"emojis"`
	Ctx     context.Context `json:"-"`
}

// ---------------------------

// GuildCreate This event can be sent in three different scenarios:
//  1. When a user is initially connecting, to lazily load and backfill information for all unavailable guilds
//     sent in the Ready event.
//	2. When a Guild becomes available again to the client.
// 	3. When the current user joins a new Guild.
type GuildCreate struct {
	Guild *Guild          `json:"guild"`
	Ctx   context.Context `json:"-"`
}

// UnmarshalJSON ...
func (obj *GuildCreate) UnmarshalJSON(data []byte) error {
	obj.Guild = &Guild{}
	return unmarshal(data, obj.Guild)
}

// ---------------------------

// GuildUpdate guild was updated
type GuildUpdate struct {
	Guild *Guild          `json:"guild"`
	Ctx   context.Context `json:"-"`
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
}

// ---------------------------

// GuildBanRemove user was unbanned from a guild
type GuildBanRemove struct {
	GuildID Snowflake       `json:"guild_id"`
	User    *User           `json:"user"`
	Ctx     context.Context `json:"-"`
}

// ---------------------------

// GuildIntegrationsUpdate guild integration was updated
type GuildIntegrationsUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Ctx     context.Context `json:"-"`
}

// ---------------------------

// GuildMemberAdd new user joined a guild
type GuildMemberAdd struct {
	Member *Member         `json:"member"`
	Ctx    context.Context `json:"-"`
}

// ---------------------------

// GuildMemberRemove user was removed from a guild
type GuildMemberRemove struct {
	GuildID Snowflake       `json:"guild_id"`
	User    *User           `json:"user"`
	Ctx     context.Context `json:"-"`
}

// ---------------------------

// GuildMemberUpdate guild member was updated
type GuildMemberUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Roles   []Snowflake     `json:"roles"`
	User    *User           `json:"user"`
	Nick    string          `json:"nick"`
	Ctx     context.Context `json:"-"`
}

// ---------------------------

// GuildMembersChunk response to Request Guild Members
type GuildMembersChunk struct {
	GuildID Snowflake       `json:"guild_id"`
	Members []*Member       `json:"members"`
	Ctx     context.Context `json:"-"`
}

// ---------------------------

// GuildRoleCreate guild role was created
type GuildRoleCreate struct {
	GuildID Snowflake       `json:"guild_id"`
	Role    *Role           `json:"role"`
	Ctx     context.Context `json:"-"`
}

// ---------------------------

// GuildRoleUpdate guild role was updated
type GuildRoleUpdate struct {
	GuildID Snowflake       `json:"guild_id"`
	Role    *Role           `json:"role"`
	Ctx     context.Context `json:"-"`
}

// ---------------------------

// GuildRoleDelete a guild role was deleted
type GuildRoleDelete struct {
	GuildID Snowflake       `json:"guild_id"`
	RoleID  Snowflake       `json:"role_id"`
	Ctx     context.Context `json:"-"`
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
	Status string          `json:"status"`
	Ctx    context.Context `json:"-"`
}

// ---------------------------

// UserUpdate properties about a user changed
type UserUpdate struct {
	User *User           `json:"user"`
	Ctx  context.Context `json:"-"`
}

// ---------------------------

// VoiceStateUpdate someone joined, left, or moved a voice channel
type VoiceStateUpdate struct {
	*VoiceState
	Ctx context.Context `json:"-"`
}

func (evt *VoiceStateUpdate) updateInternalsWithClient(c *Client) {
	c.voiceRepository.onVoiceStateUpdate(evt)
}

// ---------------------------

// VoiceServerUpdate guild's voice server was updated. Sent when a guild's voice server is updated. This is sent when initially
// connecting to voice, and when the current voice instance fails over to a new server.
type VoiceServerUpdate struct {
	Token    string          `json:"token"`
	GuildID  Snowflake       `json:"guild_id"`
	Endpoint string          `json:"endpoint"`
	Ctx      context.Context `json:"-"`
}

func (evt *VoiceServerUpdate) updateInternalsWithClient(c *Client) {
	c.voiceRepository.onVoiceServerUpdate(evt)
}

// ---------------------------

// WebhooksUpdate guild channel webhook was created, update, or deleted
type WebhooksUpdate struct {
	GuildID   Snowflake       `json:"guild_id"`
	ChannelID Snowflake       `json:"channel_id"`
	Ctx       context.Context `json:"-"`
}
