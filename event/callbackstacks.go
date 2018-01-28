package event

import (
	"errors"
)

type CallbackStackInterface interface {
	Add(interface{}) error
	Trigger(...interface{}) error // TODO: the param should be a specific event holder type
}

type ReadyCallbackStack struct {
	listeners []ReadyCallback
}

func (rds *ReadyCallbackStack) Add(i interface{}) (err error) {
	cb, ok := (i).(ReadyCallback)
	if !ok {
		return errors.New("cannot convert interface to *ReadyCallback")
	}

	if rds.listeners == nil {
		rds.listeners = []ReadyCallback{}
	}

	rds.listeners = append(rds.listeners, cb)

	return nil
}

func (rds *ReadyCallbackStack) Trigger(is ...interface{}) (err error) {

	for _, listener := range rds.listeners {
		listener()
	}

	return nil
}

type ResumeCallbackStack struct{}

// channel
type ChannelCreateCallbackStack struct{}
type ChannelUpdateCallbackStack struct{}
type ChannelDeleteCallbackStack struct{}
type ChannelPinsUpdateCallbackStack struct{}

// Guild in general
type GuildCreateCallbackStack struct{}
type GuildUpdateCallbackStack struct{}
type GuildDeleteCallbackStack struct{}
type GuildBanAddCallbackStack struct{}
type GuildBanRemoveCallbackStack struct{}
type GuildEmojisUpdateCallbackStack struct{}
type GuildIntegrationsUpdateCallbackStack struct{}

// Guild Member
type GuildMemberAddCallbackStack struct{}
type GuildMemberRemoveCallbackStack struct{}
type GuildMemberUpdateCallbackStack struct{}
type GuildMemberChunkCallbackStack struct{}

// Guild role
type GuildRoleCreateCallbackStack struct{}
type GuildRoleUpdateCallbackStack struct{}
type GuildRoleDeleteCallbackStack struct{}

// message
type MessageCreateCallbackStack struct{}
type MessageUpdateCallbackStack struct{}
type MessageDeleteCallbackStack struct{}
type MessageDeleteBulkCallbackStack struct{}

// message reaction
type MessageReactionAddCallbackStack struct{}
type MessageReactionRemoveCallbackStack struct{}
type MessageReactionRemoveAllCallbackStack struct{}

// presence
type PresenceUpdateCallbackStack struct{}

// typing start
type TypingStartCallbackStack struct{}

// user update
type UserUpdateCallbackStack struct{}

// voice
type VoiceStateUpdateCallbackStack struct{}
type VoiceServerUpdateCallbackStack struct{}

// webhook
type WebhooksUpdateCallbackStack struct{}
