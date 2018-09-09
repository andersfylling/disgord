package endpoint

import "fmt"

func Channels() string {
	return channels
}

func Channel(id fmt.Stringer) string {
	return channels + "/" + id.String()
}

func ChannelTyping(id fmt.Stringer) string {
	return Channel(id) + typing
}

func ChannelInvites(id fmt.Stringer) string {
	return Channel(id) + invites
}

func ChannelInvite(channelID, inviteID fmt.Stringer) string {
	return ChannelInvites(channelID) + "/" + inviteID.String()
}

func ChannelRecipients(channelID fmt.Stringer) string {
	return Channel(channelID) + recipients
}

func ChannelRecipient(channelID, recipientID fmt.Stringer) string {
	return ChannelRecipients(channelID) + "/" + recipientID.String()
}

// ChannelPermissions /channels/{channel.id}/permissions
func ChannelPermissions(channelID fmt.Stringer) string {
	return Channel(channelID) + permissions
}

// ChannelPermission /channels/{channel.id}/permissions/{overwrite.id}
func ChannelPermission(channelID, overwriteID fmt.Stringer) string {
	return ChannelPermissions(channelID) + "/" + overwriteID.String()
}

func ChannelPins(channelID fmt.Stringer) string {
	return Channel(channelID) + pins
}

func ChannelPin(channelID, messageID fmt.Stringer) string {
	return ChannelPins(channelID) + "/" + messageID.String()
}

func ChannelMessages(channelID fmt.Stringer) string {
	return Channel(channelID) + messages
}

func ChannelMessagesBulkDelete(channelID fmt.Stringer) string {
	return ChannelMessages(channelID) + bulkDelete
}

func ChannelMessage(channelID, messageID fmt.Stringer) string {
	return ChannelMessages(channelID) + "/" + messageID.String()
}

func ChannelMessageReactions(channelID, messageID fmt.Stringer) string {
	return ChannelMessage(channelID, messageID) + reactions
}

func ChannelMessageReaction(channelID, messageID fmt.Stringer, emoji string) string {
	return ChannelMessage(channelID, messageID) + reactions + "/" + emoji
}

func ChannelMessageReactionMe(channelID, messageID fmt.Stringer, emoji string) string {
	return ChannelMessage(channelID, messageID) + reactions + "/" + emoji + me
}

func ChannelMessageReactionUser(channelID, messageID fmt.Stringer, emoji string, userID fmt.Stringer) string {
	return ChannelMessage(channelID, messageID) + reactions + "/" + emoji + "/" + userID.String()
}
