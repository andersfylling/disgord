package endpoint

import . "github.com/andersfylling/snowflake"


func Channels() string {
	return channels
}

func Channel(id Snowflake) string {
	return channels + "/" + id.String()
}

func ChannelTyping(id Snowflake) string {
	return Channel(id) + typing
}

func ChannelInvites(id Snowflake) string {
	return Channel(id) + invites
}

func ChannelInvite(channelID, inviteID Snowflake) string {
	return ChannelInvites(channelID) + "/" + inviteID.String()
}

func ChannelRecipients(channelID Snowflake) string {
	return Channel(channelID) + recipients
}

func ChannelRecipient(channelID, recipientID Snowflake) string {
	return ChannelRecipients(channelID) + "/" + recipientID.String()
}

// ChannelPermissions /channels/{channel.id}/permissions
func ChannelPermissions(channelID Snowflake) string {
	return Channel(channelID) + permissions
}

// ChannelPermission /channels/{channel.id}/permissions/{overwrite.id}
func ChannelPermission(channelID, overwriteID Snowflake) string {
	return ChannelPermissions(channelID) + "/" + overwriteID.String()
}

func ChannelPins(channelID Snowflake) string {
	return Channel(channelID) + pins
}

func ChannelPin(channelID, messageID Snowflake) string {
	return ChannelPins(channelID) + "/" + messageID.String()
}

func ChannelMessages(channelID Snowflake) string {
	return Channel(channelID) + messages
}

func ChannelMessagesBulkDelete(channelID Snowflake) string {
	return ChannelMessages(channelID) + bulkDelete
}

func ChannelMessage(channelID, messageID Snowflake) string {
	return ChannelMessages(channelID) + "/" + messageID.String()
}

func ChannelMessageReactions(channelID, messageID Snowflake) string {
	return ChannelMessage(channelID, messageID) + reactions
}

func ChannelMessageReaction(channelID, messageID Snowflake, emoji string) string {
	return ChannelMessage(channelID, messageID) + reactions + "/" + emoji
}

func ChannelMessageReactionMe(channelID, messageID Snowflake, emoji string) string {
	return ChannelMessage(channelID, messageID) + reactions + "/" + emoji + me
}

func ChannelMessageReactionUser(channelID, messageID Snowflake, emoji string, userID Snowflake) string {
	return ChannelMessage(channelID, messageID) + reactions + "/" + emoji + "/" + userID.String()
}