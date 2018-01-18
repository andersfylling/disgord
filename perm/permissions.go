package perm

// copied from: https://github.com/bwmarrin/discordgo/blob/master/structs.go
// but altered

// Constants for the different bit offsets of text channel permissions
const (
	ReadMessages = 1 << (iota + 10)
	SendMessages
	SendTTSMessages
	ManageMessages
	EmbedLinks
	AttachFiles
	ReadMessageHistory
	MentionEveryone
	UseExternalEmojis
)

// Constants for the different bit offsets of voice permissions
const (
	VoiceConnect = 1 << (iota + 20)
	VoiceSpeak
	VoiceMuteMembers
	VoiceDeafenMembers
	VoiceMoveMembers
	VoiceUseVAD
)

// Constants for general management.
const (
	ChangeNickname = 1 << (iota + 26)
	ManageNicknames
	ManageRoles
	ManageWebhooks
	ManageEmojis
)

// Constants for the different bit offsets of general permissions
const (
	CreateInstantInvite = 1 << iota
	KickMembers
	BanMembers
	Administrator
	ManageChannels
	ManageServer
	AddReactions
	ViewAuditLogs

	AllText = ReadMessages |
		SendMessages |
		SendTTSMessages |
		ManageMessages |
		EmbedLinks |
		AttachFiles |
		ReadMessageHistory |
		MentionEveryone
	AllVoice = VoiceConnect |
		VoiceSpeak |
		VoiceMuteMembers |
		VoiceDeafenMembers |
		VoiceMoveMembers |
		VoiceUseVAD
	AllChannel = AllText |
		AllVoice |
		CreateInstantInvite |
		ManageRoles |
		ManageChannels |
		AddReactions |
		ViewAuditLogs
	All = AllChannel |
		KickMembers |
		BanMembers |
		ManageServer |
		Administrator
)
