package event

import (
	"encoding/json"

	"github.com/andersfylling/disgord/disgordctx"
)

func NewDispatcher() *Dispatcher {

	return &Dispatcher{

		// socket
		HelloEvent:          NewHelloCallbackStack(),
		ReadyEvent:          NewReadyCallbackStack(),
		ResumedEvent:        NewResumedCallbackStack(),
		InvalidSessionEvent: NewInvalidSessionCallbackStack(),

		// channel
		ChannelCreateEvent:     NewChannelCreateCallbackStack(),
		ChannelUpdateEvent:     NewChannelUpdateCallbackStack(),
		ChannelDeleteEvent:     NewChannelDeleteCallbackStack(),
		ChannelPinsUpdateEvent: NewChannelPinsUpdateCallbackStack(),

		// Guild in general
		GuildCreateEvent:             NewGuildCreateCallbackStack(),
		GuildUpdateEvent:             NewGuildUpdateCallbackStack(),
		GuildDeleteEvent:             NewGuildDeleteCallbackStack(),
		GuildBanAddEvent:             NewGuildBanAddCallbackStack(),
		GuildBanRemoveEvent:          NewGuildBanRemoveCallbackStack(),
		GuildEmojisUpdateEvent:       NewGuildEmojisUpdateCallbackStack(),
		GuildIntegrationsUpdateEvent: NewGuildIntegrationsUpdateCallbackStack(),

		// Guild Member
		GuildMemberAddEvent:    NewGuildMemberAddCallbackStack(),
		GuildMemberRemoveEvent: NewGuildMemberRemoveCallbackStack(),
		GuildMemberUpdateEvent: NewGuildMemberUpdateCallbackStack(),
		GuildMembersChunkEvent: NewGuildMembersChunkCallbackStack(),

		// Guild role
		GuildRoleCreateEvent: NewGuildRoleCreateCallbackStack(),
		GuildRoleUpdateEvent: NewGuildRoleUpdateCallbackStack(),
		GuildRoleDeleteEvent: NewGuildRoleDeleteCallbackStack(),

		// message
		MessageCreateEvent:     NewMessageCreateCallbackStack(),
		MessageUpdateEvent:     NewMessageUpdateCallbackStack(),
		MessageDeleteEvent:     NewMessageDeleteCallbackStack(),
		MessageDeleteBulkEvent: NewMessageDeleteBulkCallbackStack(),

		// message reaction
		MessageReactionAddEvent:       NewMessageReactionAddCallbackStack(),
		MessageReactionRemoveEvent:    NewMessageReactionRemoveCallbackStack(),
		MessageReactionRemoveAllEvent: NewMessageReactionRemoveAllCallbackStack(),

		// presence
		PresenceUpdateEvent: NewPresenceUpdateCallbackStack(),

		// typing start
		TypingStartEvent: NewTypingStartCallbackStack(),

		// user update
		UserUpdateEvent: NewUserUpdateCallbackStack(),

		// voice
		VoiceStateUpdateEvent:  NewVoiceStateUpdateCallbackStack(),
		VoiceServerUpdateEvent: NewVoiceServerUpdateCallbackStack(),

		// webhook
		WebhooksUpdateEvent: NewWebhooksUpdateCallbackStack(),
	}
}

//
// type DispatcherInterface interface {
// 	// // add all the callback stacks
// 	// // socket
// 	// ReadyHandler
// 	// ResumedHandler
// 	//
// 	// // channel
// 	// ChannelCreateHandler
// 	// ChannelUpdateHandler
// 	// ChannelDeleteHandler
// 	// ChannelPinsUpdateHandler
// 	//
// 	// // Guild in general
// 	// GuildCreateHandler
// 	// GuildUpdateHandler
// 	// GuildDeleteHandler
// 	// GuildBanAddHandler
// 	// GuildBanRemoveHandler
// 	// GuildEmojisUpdateHandler
// 	// GuildIntegrationsUpdateHandler
// 	//
// 	// // Guild Member
// 	// GuildMemberAddHandler
// 	// GuildMemberRemoveHandler
// 	// GuildMemberUpdateHandler
// 	// GuildMemberChunkHandler
// 	//
// 	// // Guild role
// 	// GuildRoleCreateHandler
// 	// GuildRoleUpdateHandler
// 	// GuildRoleDeleteHandler
// 	//
// 	// // message
// 	// MessageCreateHandler
// 	// MessageUpdateHandler
// 	// MessageDeleteHandler
// 	// MessageDeleteBulkHandler
// 	//
// 	// // message reaction
// 	// MessageReactionAddHandler
// 	// MessageReactionRemoveHandler
// 	// MessageReactionRemoveAllHandler
// 	//
// 	// // presence
// 	// PresenceUpdateHandler
// 	//
// 	// // typing start
// 	// TypingStartHandler
// 	//
// 	// // user update
// 	// UserUpdateHandler
// 	//
// 	// // voice
// 	// VoiceStateUpdateHandler
// 	// VoiceServerUpdateHandler
// 	//
// 	// // webhook
// 	// WebhooksUpdateHandler
//
// 	Event(eventName Type, listener interface{})
// 	Trigger(eventName Type, params ...*interface{})
// }

type Dispatcher struct {
	// add all the callback stacks
	// socket
	HelloEvent          *HelloCallbackStack
	ReadyEvent          *ReadyCallbackStack
	ResumedEvent        *ResumedCallbackStack
	InvalidSessionEvent *InvalidSessionCallbackStack

	// channel
	ChannelCreateEvent     *ChannelCreateCallbackStack
	ChannelUpdateEvent     *ChannelUpdateCallbackStack
	ChannelDeleteEvent     *ChannelDeleteCallbackStack
	ChannelPinsUpdateEvent *ChannelPinsUpdateCallbackStack

	// Guild in general
	GuildCreateEvent             *GuildCreateCallbackStack
	GuildUpdateEvent             *GuildUpdateCallbackStack
	GuildDeleteEvent             *GuildDeleteCallbackStack
	GuildBanAddEvent             *GuildBanAddCallbackStack
	GuildBanRemoveEvent          *GuildBanRemoveCallbackStack
	GuildEmojisUpdateEvent       *GuildEmojisUpdateCallbackStack
	GuildIntegrationsUpdateEvent *GuildIntegrationsUpdateCallbackStack

	// Guild Member
	GuildMemberAddEvent    *GuildMemberAddCallbackStack
	GuildMemberRemoveEvent *GuildMemberRemoveCallbackStack
	GuildMemberUpdateEvent *GuildMemberUpdateCallbackStack
	GuildMembersChunkEvent *GuildMembersChunkCallbackStack

	// Guild role
	GuildRoleUpdateEvent *GuildRoleUpdateCallbackStack
	GuildRoleCreateEvent *GuildRoleCreateCallbackStack
	GuildRoleDeleteEvent *GuildRoleDeleteCallbackStack

	// message
	MessageCreateEvent     *MessageCreateCallbackStack
	MessageUpdateEvent     *MessageUpdateCallbackStack
	MessageDeleteEvent     *MessageDeleteCallbackStack
	MessageDeleteBulkEvent *MessageDeleteBulkCallbackStack

	// message reaction
	MessageReactionAddEvent       *MessageReactionAddCallbackStack
	MessageReactionRemoveEvent    *MessageReactionRemoveCallbackStack
	MessageReactionRemoveAllEvent *MessageReactionRemoveAllCallbackStack

	// presence
	PresenceUpdateEvent *PresenceUpdateCallbackStack

	// typing start
	TypingStartEvent *TypingStartCallbackStack

	// user update
	UserUpdateEvent *UserUpdateCallbackStack

	// voice
	VoiceStateUpdateEvent  *VoiceStateUpdateCallbackStack
	VoiceServerUpdateEvent *VoiceServerUpdateCallbackStack

	// webhook
	WebhooksUpdateEvent *WebhooksUpdateCallbackStack
}

// On places listeners into their respected stacks
func (d *Dispatcher) OnEvent(eventKey KeyType, listener interface{}) {
	d.handleStateUpdate(eventKey, listener, nil)
}

// Trigger listeners based on the event type
func (d *Dispatcher) Trigger(eventKey KeyType, box interface{}, ctx disgordctx.Context) {
	d.handleStateUpdate(eventKey, box, ctx)
}

func (d *Dispatcher) handleStateUpdate(eventKey KeyType, content interface{}, ctx disgordctx.Context) {
	switch eventKey {
	case ReadyKey:
		if ctx == nil {
			d.ReadyEvent.Add(content.(ReadyCallback))
		} else {
			d.ReadyEvent.Trigger(ctx, content.(*ReadyBox))
		}
	case ResumedKey:
		if ctx == nil {
			d.ResumedEvent.Add(content.(ResumedCallback))
		} else {
			d.ResumedEvent.Trigger(ctx, content.(*ResumedBox))
		}
	case ChannelCreateKey:
		if ctx == nil {
			d.ChannelCreateEvent.Add(content.(ChannelCreateCallback))
		} else {
			d.ChannelCreateEvent.Trigger(ctx, content.(*ChannelCreateBox))
		}
	case ChannelUpdateKey:
		if ctx == nil {
			d.ChannelUpdateEvent.Add(content.(ChannelUpdateCallback))
		} else {
			d.ChannelUpdateEvent.Trigger(ctx, content.(*ChannelUpdateBox))
		}
	case ChannelDeleteKey:
		if ctx == nil {
			d.ChannelDeleteEvent.Add(content.(ChannelDeleteCallback))
		} else {
			d.ChannelDeleteEvent.Trigger(ctx, content.(*ChannelDeleteBox))
		}
	case ChannelPinsUpdateKey:
		if ctx == nil {
			d.ChannelPinsUpdateEvent.Add(content.(ChannelPinsUpdateCallback))
		} else {
			d.ChannelPinsUpdateEvent.Trigger(ctx, content.(*ChannelPinsUpdateBox))
		}
	case GuildCreateKey:
		if ctx == nil {
			d.GuildCreateEvent.Add(content.(GuildCreateCallback))
		} else {
			d.GuildCreateEvent.Trigger(ctx, content.(*GuildCreateBox))
		}
	case GuildUpdateKey:
		if ctx == nil {
			d.GuildUpdateEvent.Add(content.(GuildUpdateCallback))
		} else {
			d.GuildUpdateEvent.Trigger(ctx, content.(*GuildUpdateBox))
		}
	case GuildDeleteKey:
		if ctx == nil {
			d.GuildDeleteEvent.Add(content.(GuildDeleteCallback))
		} else {
			d.GuildDeleteEvent.Trigger(ctx, content.(*GuildDeleteBox))
		}
	case GuildBanAddKey:
		if ctx == nil {
			d.GuildBanAddEvent.Add(content.(GuildBanAddCallback))
		} else {
			d.GuildBanAddEvent.Trigger(ctx, content.(*GuildBanAddBox))
		}
	case GuildBanRemoveKey:
		if ctx == nil {
			d.GuildBanRemoveEvent.Add(content.(GuildBanRemoveCallback))
		} else {
			d.GuildBanRemoveEvent.Trigger(ctx, content.(*GuildBanRemoveBox))
		}
	case GuildEmojisUpdateKey:
		if ctx == nil {
			d.GuildEmojisUpdateEvent.Add(content.(GuildEmojisUpdateCallback))
		} else {
			d.GuildEmojisUpdateEvent.Trigger(ctx, content.(*GuildEmojisUpdateBox))
		}
	case GuildIntegrationsUpdateKey:
		if ctx == nil {
			d.GuildIntegrationsUpdateEvent.Add(content.(GuildIntegrationsUpdateCallback))
		} else {
			d.GuildIntegrationsUpdateEvent.Trigger(ctx, content.(*GuildIntegrationsUpdateBox))
		}
	case GuildMemberAddKey:
		if ctx == nil {
			d.GuildMemberAddEvent.Add(content.(GuildMemberAddCallback))
		} else {
			d.GuildMemberAddEvent.Trigger(ctx, content.(*GuildMemberAddBox))
		}
	case GuildMemberRemoveKey:
		if ctx == nil {
			d.GuildMemberRemoveEvent.Add(content.(GuildMemberRemoveCallback))
		} else {
			d.GuildMemberRemoveEvent.Trigger(ctx, content.(*GuildMemberRemoveBox))
		}
	case GuildMemberUpdateKey:
		if ctx == nil {
			d.GuildMemberUpdateEvent.Add(content.(GuildMemberUpdateCallback))
		} else {
			d.GuildMemberUpdateEvent.Trigger(ctx, content.(*GuildMemberUpdateBox))
		}
	case GuildMembersChunkKey:
		if ctx == nil {
			d.GuildMembersChunkEvent.Add(content.(GuildMembersChunkCallback))
		} else {
			d.GuildMembersChunkEvent.Trigger(ctx, content.(*GuildMembersChunkBox))
		}
	case GuildRoleCreateKey:
		if ctx == nil {
			d.GuildRoleCreateEvent.Add(content.(GuildRoleCreateCallback))
		} else {
			d.GuildRoleCreateEvent.Trigger(ctx, content.(*GuildRoleCreateBox))
		}
	case GuildRoleUpdateKey:
		if ctx == nil {
			d.GuildRoleUpdateEvent.Add(content.(GuildRoleUpdateCallback))
		} else {
			d.GuildRoleUpdateEvent.Trigger(ctx, content.(*GuildRoleUpdateBox))
		}
	case GuildRoleDeleteKey:
		if ctx == nil {
			d.GuildRoleDeleteEvent.Add(content.(GuildRoleDeleteCallback))
		} else {
			d.GuildRoleDeleteEvent.Trigger(ctx, content.(*GuildRoleDeleteBox))
		}
	case MessageDeleteBulkKey:
		if ctx == nil {
			d.MessageDeleteBulkEvent.Add(content.(MessageDeleteBulkCallback))
		} else {
			d.MessageDeleteBulkEvent.Trigger(ctx, content.(*MessageDeleteBulkBox))
		}
	case MessageReactionAddKey:
		if ctx == nil {
			d.MessageReactionAddEvent.Add(content.(MessageReactionAddCallback))
		} else {
			d.MessageReactionAddEvent.Trigger(ctx, content.(*MessageReactionAddBox))
		}
	case MessageReactionRemoveKey:
		if ctx == nil {
			d.MessageReactionRemoveEvent.Add(content.(MessageReactionRemoveCallback))
		} else {
			d.MessageReactionRemoveEvent.Trigger(ctx, content.(*MessageReactionRemoveBox))
		}
	case MessageReactionRemoveAllKey:
		if ctx == nil {
			d.MessageReactionRemoveAllEvent.Add(content.(MessageReactionRemoveAllCallback))
		} else {
			d.MessageReactionRemoveAllEvent.Trigger(ctx, content.(*MessageReactionRemoveAllBox))
		}
	case PresenceUpdateKey:
		if ctx == nil {
			d.PresenceUpdateEvent.Add(content.(PresenceUpdateCallback))
		} else {
			d.PresenceUpdateEvent.Trigger(ctx, content.(*PresenceUpdateBox))
		}
	case TypingStartKey:
		if ctx == nil {
			d.TypingStartEvent.Add(content.(TypingStartCallback))
		} else {
			d.TypingStartEvent.Trigger(ctx, content.(*TypingStartBox))
		}
	case UserUpdateKey:
		if ctx == nil {
			d.UserUpdateEvent.Add(content.(UserUpdateCallback))
		} else {
			d.UserUpdateEvent.Trigger(ctx, content.(*UserUpdateBox))
		}
	case VoiceStateUpdateKey:
		if ctx == nil {
			d.VoiceStateUpdateEvent.Add(content.(VoiceStateUpdateCallback))
		} else {
			d.VoiceStateUpdateEvent.Trigger(ctx, content.(*VoiceStateUpdateBox))
		}
	case VoiceServerUpdateKey:
		if ctx == nil {
			d.VoiceServerUpdateEvent.Add(content.(VoiceServerUpdateCallback))
		} else {
			d.VoiceServerUpdateEvent.Trigger(ctx, content.(*VoiceServerUpdateBox))
		}
	case WebhooksUpdateKey:
		if ctx == nil {
			d.WebhooksUpdateEvent.Add(content.(WebhooksUpdateCallback))
		} else {
			d.WebhooksUpdateEvent.Trigger(ctx, content.(*WebhooksUpdateBox))
		}

	default:
		msg := "unknown event key was used to "
		if ctx == nil {
			msg = msg + " register event reactor/listener"
		} else {
			msg = msg + " trigger reactors/listeners"
		}
		panic(msg)
	}
}

func Unmarshal(data []byte, box interface{}) {
	err := json.Unmarshal(data, box)
	if err != nil {
		panic(err)
	}
}
