package event

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
		GuildMemberChunkEvent:  NewGuildMembersChunkCallbackStack(),

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
	GuildMemberChunkEvent  *GuildMembersChunkCallbackStack

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
func (d *Dispatcher) On(eventName Type, listener interface{}) {
}

// Trigger listeners based on the event type
func (d *Dispatcher) Trigger(eventName Type, params ...*interface{}) {
}
