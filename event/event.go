package event

import "github.com/sirupsen/logrus"

// EventDispatcher is an application-level type for handling discord requests.
// All callbacks are optional, and whether they are defined or not
// is used to determine whether the EventDispatcher will send events to them.
type DispatcherInterface interface {
	// current EventHook fields here

	// OnEvent is called for all events.
	// Handlers must typecast the event type manually, and ensure
	// that it can handle receiving the same event twice if a type-specific
	// callback also exists.
	//OnEvent func(ctx *Context, ev event.DiscordEvent) error

	// OnMessageEvent is called for every message-related event.
	//OnMessageEvent func(ctx *Context, ev event.MessageEvent) error

	// OnConnectionEvent ...
	//OnUserEvent(eventName string, listener func(user *user.User))
	//OnMemberEvent(eventName string, listener func(member *guild.Member))
	// OnChannelEvent ...
	// OnGuildEvent ...

	OnEvent(eventName Type, listener interface{})
	Trigger(eventName Type, params ...interface{})
}

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{
		listeners: make(map[Type]CallbackStackInterface),
	}

	// add all the callback stacks
	// socket
	d.listeners[Ready] = &ReadyCallbackStack{}
	d.listeners[Resumed] = &ResumeCallbackStack{}

	// channel
	d.listeners[ChannelCreate] = &ChannelCreateCallbackStack{}
	d.listeners[ChannelUpdate] = &ChannelUpdateCallbackStack{}
	d.listeners[ChannelDelete] = &ChannelDeleteCallbackStack{}
	d.listeners[ChannelPinsUpdate] = &ChannelPinsUpdateCallbackStack{}

	// Guild in general
	d.listeners[GuildCreate] = &GuildCreateCallbackStack{}
	d.listeners[GuildUpdate] = &GuildUpdateCallbackStack{}
	d.listeners[GuildDelete] = &GuildDeleteCallbackStack{}
	d.listeners[GuildBanAdd] = &GuildBanAddCallbackStack{}
	d.listeners[GuildBanRemove] = &GuildBanRemoveCallbackStack{}
	d.listeners[GuildEmojisUpdate] = &GuildEmojisUpdateCallbackStack{}
	d.listeners[GuildIntegrationsUpdate] = &GuildIntegrationsUpdateCallbackStack{}

	// Guild Member
	d.listeners[GuildMemberAdd] = &GuildMemberAddCallbackStack{}
	d.listeners[GuildMemberRemove] = &GuildMemberRemoveCallbackStack{}
	d.listeners[GuildMemberUpdate] = &GuildMemberUpdateCallbackStack{}
	d.listeners[GuildMemberChunk] = &GuildMemberChunkCallbackStack{}

	// Guild role
	d.listeners[GuildRoleCreate] = &GuildRoleCreateCallbackStack{}
	d.listeners[GuildRoleUpdate] = &GuildRoleUpdateCallbackStack{}
	d.listeners[GuildRoleDelete] = &GuildRoleDeleteCallbackStack{}

	// message
	d.listeners[MessageCreate] = &MessageCreateCallbackStack{}
	d.listeners[MessageUpdate] = &MessageUpdateCallbackStack{}
	d.listeners[MessageDelete] = &MessageDeleteCallbackStack{}
	d.listeners[MessageDeleteBulk] = &MessageDeleteBulkCallbackStack{}

	// message reaction
	d.listeners[MessageReactionAdd] = &MessageReactionAddCallbackStack{}
	d.listeners[MessageReactionRemove] = &MessageReactionRemoveCallbackStack{}
	d.listeners[MessageReactionRemoveAll] = &MessageReactionRemoveAllCallbackStack{}

	// presence
	d.listeners[PresenceUpdate] = &PresenceUpdateCallbackStack{}

	// typing start
	d.listeners[TypingStart] = &TypingStartCallbackStack{}

	// user update
	d.listeners[UserUpdate] = &UserUpdateCallbackStack{}

	// voice
	d.listeners[VoiceStateUpdate] = &VoiceStateUpdateCallbackStack{}
	d.listeners[VoiceServerUpdate] = &VoiceServerUpdateCallbackStack{}

	// webhook
	d.listeners[WebhooksUpdate] = &WebhooksUpdateCallbackStack{}

	return d
}

type Dispatcher struct {
	listeners map[Type]CallbackStackInterface
}

func (d *Dispatcher) OnEvent(eventName Type, listener interface{}) {
	if listeners, ok := d.listeners[eventName]; ok {
		listeners.Add(listener)
	} else {
		logrus.Errorf("no callback interface registered for `%s`", eventName)
	}
}

func (d *Dispatcher) Trigger(eventName Type, params ...interface{}) {
	if listeners, ok := d.listeners[eventName]; ok {
		listeners.Trigger(params)
	}
}
