package event

import "testing"

func generateError(name string) string {
	return name + " does not implement its designated handler interface"
}

func TestCallbackStackInterfaceImplementation(t *testing.T) {
	// ... i wish i knew go::generate
	if _, ok := interface{}(&HelloCallbackStack{}).(HelloHandler); !ok {
		t.Error(generateError("HelloCallbackStack"))
	}

	if _, ok := interface{}(&ReadyCallbackStack{}).(ReadyHandler); !ok {
		t.Error(generateError("ReadyCallbackStack"))
	}

	if _, ok := interface{}(&ResumedCallbackStack{}).(ResumedHandler); !ok {
		t.Error(generateError("ResumedCallbackStack"))
	}

	if _, ok := interface{}(&InvalidSessionCallbackStack{}).(InvalidSessionHandler); !ok {
		t.Error(generateError("InvalidSessionCallbackStack"))
	}

	if _, ok := interface{}(&ChannelCreateCallbackStack{}).(ChannelCreateHandler); !ok {
		t.Error(generateError("ChannelCreateCallbackStack"))
	}

	if _, ok := interface{}(&ChannelUpdateCallbackStack{}).(ChannelUpdateHandler); !ok {
		t.Error(generateError("ChannelUpdateCallbackStack"))
	}

	if _, ok := interface{}(&ChannelDeleteCallbackStack{}).(ChannelDeleteHandler); !ok {
		t.Error(generateError("ChannelDeleteCallbackStack"))
	}

	if _, ok := interface{}(&ChannelPinsUpdateCallbackStack{}).(ChannelPinsUpdateHandler); !ok {
		t.Error(generateError("ChannelPinsUpdateCallbackStack"))
	}

	if _, ok := interface{}(&GuildCreateCallbackStack{}).(GuildCreateHandler); !ok {
		t.Error(generateError("GuildCreateCallbackStack"))
	}

	if _, ok := interface{}(&GuildUpdateCallbackStack{}).(GuildUpdateHandler); !ok {
		t.Error(generateError("GuildUpdateCallbackStack"))
	}

	if _, ok := interface{}(&GuildDeleteCallbackStack{}).(GuildDeleteHandler); !ok {
		t.Error(generateError("GuildDeleteCallbackStack"))
	}

	if _, ok := interface{}(&GuildBanAddCallbackStack{}).(GuildBanAddHandler); !ok {
		t.Error(generateError("GuildBanAddCallbackStack"))
	}

	if _, ok := interface{}(&GuildBanRemoveCallbackStack{}).(GuildBanRemoveHandler); !ok {
		t.Error(generateError("GuildBanRemoveCallbackStack"))
	}

	if _, ok := interface{}(&GuildEmojisUpdateCallbackStack{}).(GuildEmojisUpdateHandler); !ok {
		t.Error(generateError("GuildEmojisUpdateCallbackStack"))
	}

	if _, ok := interface{}(&GuildIntegrationsUpdateCallbackStack{}).(GuildIntegrationsUpdateHandler); !ok {
		t.Error(generateError("GuildIntegrationsUpdateCallbackStack"))
	}

	if _, ok := interface{}(&GuildMemberAddCallbackStack{}).(GuildMemberAddHandler); !ok {
		t.Error(generateError("GuildMemberAddCallbackStack"))
	}

	if _, ok := interface{}(&GuildMemberRemoveCallbackStack{}).(GuildMemberRemoveHandler); !ok {
		t.Error(generateError("GuildMemberRemoveCallbackStack"))
	}

	if _, ok := interface{}(&GuildMemberUpdateCallbackStack{}).(GuildMemberUpdateHandler); !ok {
		t.Error(generateError("GuildMemberUpdateCallbackStack"))
	}

	if _, ok := interface{}(&GuildMembersChunkCallbackStack{}).(GuildMembersChunkHandler); !ok {
		t.Error(generateError("GuildMembersChunkCallbackStack"))
	}

	if _, ok := interface{}(&GuildRoleCreateCallbackStack{}).(GuildRoleCreateHandler); !ok {
		t.Error(generateError("GuildRoleCreateCallbackStack"))
	}

	if _, ok := interface{}(&GuildRoleUpdateCallbackStack{}).(GuildRoleUpdateHandler); !ok {
		t.Error(generateError("GuildRoleUpdateCallbackStack"))
	}

	if _, ok := interface{}(&GuildRoleDeleteCallbackStack{}).(GuildRoleDeleteHandler); !ok {
		t.Error(generateError("GuildRoleDeleteCallbackStack"))
	}

	if _, ok := interface{}(&MessageCreateCallbackStack{}).(MessageCreateHandler); !ok {
		t.Error(generateError("MessageCreateCallbackStack"))
	}

	if _, ok := interface{}(&MessageUpdateCallbackStack{}).(MessageUpdateHandler); !ok {
		t.Error(generateError("MessageUpdateCallbackStack"))
	}

	if _, ok := interface{}(&MessageDeleteCallbackStack{}).(MessageDeleteHandler); !ok {
		t.Error(generateError("MessageDeleteCallbackStack"))
	}

	if _, ok := interface{}(&MessageDeleteBulkCallbackStack{}).(MessageDeleteBulkHandler); !ok {
		t.Error(generateError("MessageDeleteBulkCallbackStack"))
	}

	if _, ok := interface{}(&MessageReactionAddCallbackStack{}).(MessageReactionAddHandler); !ok {
		t.Error(generateError("MessageReactionAddCallbackStack"))
	}

	if _, ok := interface{}(&MessageReactionRemoveCallbackStack{}).(MessageReactionRemoveHandler); !ok {
		t.Error(generateError("MessageReactionRemoveCallbackStack"))
	}

	if _, ok := interface{}(&MessageReactionRemoveAllCallbackStack{}).(MessageReactionRemoveAllHandler); !ok {
		t.Error(generateError("MessageReactionRemoveAllCallbackStack"))
	}

	if _, ok := interface{}(&PresenceUpdateCallbackStack{}).(PresenceUpdateHandler); !ok {
		t.Error(generateError("PresenceUpdateCallbackStack"))
	}

	if _, ok := interface{}(&TypingStartCallbackStack{}).(TypingStartHandler); !ok {
		t.Error(generateError("TypingStartCallbackStack"))
	}

	if _, ok := interface{}(&UserUpdateCallbackStack{}).(UserUpdateHandler); !ok {
		t.Error(generateError("UserUpdateCallbackStack"))
	}

	if _, ok := interface{}(&VoiceStateUpdateCallbackStack{}).(VoiceStateUpdateHandler); !ok {
		t.Error(generateError("VoiceStateUpdateCallbackStack"))
	}

	if _, ok := interface{}(&VoiceServerUpdateCallbackStack{}).(VoiceServerUpdateHandler); !ok {
		t.Error(generateError("VoiceServerUpdateCallbackStack"))
	}

	if _, ok := interface{}(&WebhooksUpdateCallbackStack{}).(WebhooksUpdateHandler); !ok {
		t.Error(generateError("WebhooksUpdateCallbackStack"))
	}
}
