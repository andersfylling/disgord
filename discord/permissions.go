package discord

// copied from: https://github.com/bwmarrin/discordgo/blob/master/structs.go
// but altered

// Constants for the different bit offsets of text channel permissions
const (
	ReadMessagesPermission = 1 << (iota + 10)
	SendMessagesPermission
	SendTTSMessagesPermission
	ManageMessagesPermission
	EmbedLinksPermission
	AttachFilesPermission
	ReadMessageHistoryPermission
	MentionEveryonePermission
	UseExternalEmojisPermission
)

// Constants for the different bit offsets of voice permissions
const (
	VoiceConnectPermission = 1 << (iota + 20)
	VoiceSpeakPermission
	VoiceMuteMembersPermission
	VoiceDeafenMembersPermission
	VoiceMoveMembersPermission
	VoiceUseVADPermission
)

// Constants for general management.
const (
	ChangeNicknamePermission = 1 << (iota + 26)
	ManageNicknamesPermission
	ManageRolesPermission
	ManageWebhooksPermission
	ManageEmojisPermission
)

// Constants for the different bit offsets of general permissions
const (
	CreateInstantInvitePermission = 1 << iota
	KickMembersPermission
	BanMembersPermission
	AdministratorPermission
	ManageChannelsPermission
	ManageServerPermission
	AddReactionsPermission
	ViewAuditLogsPermission

	AllTextPermission = ReadMessagesPermission |
		SendMessagesPermission |
		SendTTSMessagesPermission |
		ManageMessagesPermission |
		EmbedLinksPermission |
		AttachFilesPermission |
		ReadMessageHistoryPermission |
		MentionEveryonePermission
	AllVoicePermission = VoiceConnectPermission |
		VoiceSpeakPermission |
		VoiceMuteMembersPermission |
		VoiceDeafenMembersPermission |
		VoiceMoveMembersPermission |
		VoiceUseVADPermission
	AllChannelPermission = AllTextPermission |
		AllVoicePermission |
		CreateInstantInvitePermission |
		ManageRolesPermission |
		ManageChannelsPermission |
		AddReactionsPermission |
		ViewAuditLogsPermission
	AllPermission = AllChannelPermission |
		KickMembersPermission |
		BanMembersPermission |
		ManageServerPermission |
		AdministratorPermission
)
