package event

import "testing"

func generateError(name string) string {
	return name + " does not implement interface `CallbackStackInterface`"
}

func TestCallbackStackInterfaceImplementation(t *testing.T) {
	s := make(map[string]interface{})

	// socket
	s["ReadyCallbackStack"] = &ReadyCallbackStack{}
	s["ResumeCallbackStack"] = &ResumeCallbackStack{}

	// channel
	s["ChannelCreateCallbackStack"] = &ChannelCreateCallbackStack{}
	s["ChannelUpdateCallbackStack"] = &ChannelUpdateCallbackStack{}
	s["ChannelDeleteCallbackStack"] = &ChannelDeleteCallbackStack{}
	s["ChannelPinsUpdateCallbackStack"] = &ChannelPinsUpdateCallbackStack{}

	// Guild in general
	s["GuildCreateCallbackStack"] = &GuildCreateCallbackStack{}
	s["GuildUpdateCallbackStack"] = &GuildUpdateCallbackStack{}
	s["GuildDeleteCallbackStack"] = &GuildDeleteCallbackStack{}
	s["GuildBanAddCallbackStack"] = &GuildBanAddCallbackStack{}
	s["GuildBanRemoveCallbackStack"] = &GuildBanRemoveCallbackStack{}
	s["GuildEmojisUpdateCallbackStack"] = &GuildEmojisUpdateCallbackStack{}
	s["GuildIntegrationsUpdateCallbackStack"] = &GuildIntegrationsUpdateCallbackStack{}

	// Guild Member
	s["GuildMemberAddCallbackStack"] = &GuildMemberAddCallbackStack{}
	s["GuildMemberRemoveCallbackStack"] = &GuildMemberRemoveCallbackStack{}
	s["GuildMemberUpdateCallbackStack"] = &GuildMemberUpdateCallbackStack{}
	s["GuildMemberChunkCallbackStack"] = &GuildMemberChunkCallbackStack{}

	// Guild role
	s["GuildRoleCreateCallbackStack"] = &GuildRoleCreateCallbackStack{}
	s["GuildRoleUpdateCallbackStack"] = &GuildRoleUpdateCallbackStack{}
	s["GuildRoleDeleteCallbackStack"] = &GuildRoleDeleteCallbackStack{}

	// message
	s["MessageCreateCallbackStack"] = &MessageCreateCallbackStack{}
	s["MessageUpdateCallbackStack"] = &MessageUpdateCallbackStack{}
	s["MessageDeleteCallbackStack"] = &MessageDeleteCallbackStack{}
	s["MessageDeleteBulkCallbackStack"] = &MessageDeleteBulkCallbackStack{}

	// message reaction
	s["MessageReactionAddCallbackStack"] = &MessageReactionAddCallbackStack{}
	s["MessageReactionRemoveCallbackStack"] = &MessageReactionRemoveCallbackStack{}
	s["MessageReactionRemoveAllCallbackStack"] = &MessageReactionRemoveAllCallbackStack{}

	// presence
	s["PresenceUpdateCallbackStack"] = &PresenceUpdateCallbackStack{}

	// typing start
	s["TypingStartCallbackStack"] = &TypingStartCallbackStack{}

	// user update
	s["UserUpdateCallbackStack"] = &UserUpdateCallbackStack{}

	// voice
	s["VoiceStateUpdateCallbackStack"] = &VoiceStateUpdateCallbackStack{}
	s["VoiceServerUpdateCallbackStack"] = &VoiceServerUpdateCallbackStack{}

	// webhook
	s["WebhooksUpdateCallbackStack"] = &WebhooksUpdateCallbackStack{}

	for name, v := range s {
		if _, ok := v.(CallbackStackInterface); !ok {
			t.Error(name + " does not implement interface `CallbackStackInterface`")
		}
	}
}
