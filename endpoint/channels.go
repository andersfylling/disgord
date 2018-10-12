package endpoint

import "fmt"

// Channels ...
func Channels() string {
	return channels
}

// Channel ...
func Channel(id fmt.Stringer) string {
	return channels + "/" + id.String()
}

// ChannelTyping ...
func ChannelTyping(id fmt.Stringer) string {
	return Channel(id) + typing
}

// ChannelInvites ...
func ChannelInvites(id fmt.Stringer) string {
	return Channel(id) + invites
}

// ChannelInvite ...
func ChannelInvite(channelID, inviteID fmt.Stringer) string {
	return ChannelInvites(channelID) + "/" + inviteID.String()
}

// ChannelRecipients ...
func ChannelRecipients(channelID fmt.Stringer) string {
	return Channel(channelID) + recipients
}

// ChannelRecipient ...
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

// ChannelPins ...
func ChannelPins(channelID fmt.Stringer) string {
	return Channel(channelID) + pins
}

// ChannelPin ...
func ChannelPin(channelID, messageID fmt.Stringer) string {
	return ChannelPins(channelID) + "/" + messageID.String()
}

// ChannelMessages ...
func ChannelMessages(channelID fmt.Stringer) string {
	return Channel(channelID) + messages
}

// ChannelMessagesBulkDelete ...
func ChannelMessagesBulkDelete(channelID fmt.Stringer) string {
	return ChannelMessages(channelID) + bulkDelete
}

// ChannelMessage ...
func ChannelMessage(channelID, messageID fmt.Stringer) string {
	return ChannelMessages(channelID) + "/" + messageID.String()
}

// ChannelMessageReactions ...
func ChannelMessageReactions(channelID, messageID fmt.Stringer) string {
	return ChannelMessage(channelID, messageID) + reactions
}

// ChannelMessageReaction ...
func ChannelMessageReaction(channelID, messageID fmt.Stringer, emoji string) string {
	return ChannelMessage(channelID, messageID) + reactions + "/" + emoji
}

// ChannelMessageReactionMe ...
func ChannelMessageReactionMe(channelID, messageID fmt.Stringer, emoji string) string {
	return ChannelMessage(channelID, messageID) + reactions + "/" + emoji + me
}

// ChannelMessageReactionUser ...
func ChannelMessageReactionUser(channelID, messageID fmt.Stringer, emoji string, userID fmt.Stringer) string {
	return ChannelMessage(channelID, messageID) + reactions + "/" + emoji + "/" + userID.String()
}
