package resource

import (
	"errors"
	"strconv"

	"github.com/andersfylling/disgord/request"
	"github.com/andersfylling/snowflake"
)

// https://discordapp.com/developers/docs/resources/channel#reaction-object
type Reaction struct {
	Count uint          `json:"count"`
	Me    bool          `json:"me"`
	Emoji *PartialEmoji `json:"Emoji"`
}

// ReqCreateReaction [PUT]	Create a reaction for the message. This endpoint requires the 'READ_MESSAGE_HISTORY'
// 							permission to be present on the current user. Additionally, if nobody else has
// 							reacted to the message using this emoji, this endpoint requires the 'ADD_REACTIONS'
// 							permission to be present on the current user. Returns a 204 empty response on success.
// 							The maximum request size when sending a message is 8MB.
// Endpoint				   	/channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
// Rate limiter [MAJOR]	   	/channels/{channel.id}
// Discord documentation   	https://discordapp.com/developers/docs/resources/channel#create-reaction
// Reviewed				   	2018-06-07
// Comment				   	-
// emoji either unicode (string) or *Emoji with an snowflake ID if it's custom
func ReqCreateReaction(client request.DiscordPutter, channelID, messageID snowflake.ID, emoji interface{}) (*Reaction, error) {
	if channelID.Empty() {
		return nil, errors.New("channelID must be set to target the correct channel")
	}
	if messageID.Empty() {
		return nil, errors.New("messageID must be set to target the specific channel message")
	}
	if emoji == nil {
		return nil, errors.New("emoji must be set in order to create a message reaction")
	}

	ratelimiter := "/channels/" + channelID.String()
	emojiCode := ""
	if _, ok := emoji.(*Emoji); ok {
		emojiCode = emoji.(*Emoji).ID.String()
	} else if _, ok := emoji.(string); ok {
		emojiCode = emoji.(string) // unicode
	} else {
		return nil, errors.New("emoji type can only be a unicode string or a *Emoji struct")
	}

	endpoint := ratelimiter + "/messages/" + messageID.String() + "/reactions/" + emojiCode + "/@me"
	_, err := client.Put(ratelimiter, endpoint, nil, nil)
	return nil, err
}

// ReqDeleteOwnReaction [DELETE]	Delete a reaction the current user has made for the message. Returns a 204
// 									empty response on success.
// Endpoint				   			/channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
// Rate limiter [MAJOR]	   			/channels/{channel.id}
// Discord documentation   			https://discordapp.com/developers/docs/resources/channel#delete-own-reaction
// Reviewed				   			2018-06-07
// Comment				   			-
// emoji either unicode (string) or *Emoji with an snowflake ID if it's custom
func ReqDeleteOwnReaction(client request.DiscordDeleter, channelID, messageID snowflake.ID, emoji interface{}) error {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if messageID.Empty() {
		return errors.New("messageID must be set to target the specific channel message")
	}
	if emoji == nil {
		return errors.New("emoji must be set in order to create a message reaction")
	}

	ratelimiter := "/channels/" + channelID.String()
	emojiCode := ""
	if _, ok := emoji.(*Emoji); ok {
		emojiCode = emoji.(*Emoji).ID.String()
	} else if _, ok := emoji.(string); ok {
		emojiCode = emoji.(string) // unicode
	} else {
		return errors.New("emoji type can only be a unicode string or a *Emoji struct")
	}

	endpoint := ratelimiter + "/messages/" + messageID.String() + "/reactions/" + emojiCode + "/@me"
	_, err := client.Delete(ratelimiter, endpoint)
	return err
}

// ReqCreateReaction [DELETE]	Deletes another user's reaction. This endpoint requires the 'MANAGE_MESSAGES'
// 								permission to be present on the current user. Returns a 204 empty response on success.
// Endpoint				   		/channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
// Rate limiter [MAJOR]	   		/channels/{channel.id}
// Discord documentation   		https://discordapp.com/developers/docs/resources/channel#delete-user-reaction
// Reviewed				   		2018-06-07
// Comment				   		-
// emoji either unicode (string) or *Emoji with an snowflake ID if it's custom
func ReqDeleteUserReaction(client request.DiscordDeleter, channelID, messageID, userID snowflake.ID, emoji interface{}) error {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if messageID.Empty() {
		return errors.New("messageID must be set to target the specific channel message")
	}
	if emoji == nil {
		return errors.New("emoji must be set in order to create a message reaction")
	}
	if userID.Empty() {
		return errors.New("userID must be set to target the specific user reaction")
	}

	ratelimiter := "/channels/" + channelID.String()
	emojiCode := ""
	if _, ok := emoji.(*Emoji); ok {
		emojiCode = emoji.(*Emoji).ID.String()
	} else if _, ok := emoji.(string); ok {
		emojiCode = emoji.(string) // unicode
	} else {
		return errors.New("emoji type can only be a unicode string or a *Emoji struct")
	}

	endpoint := ratelimiter + "/messages/" + messageID.String() + "/reactions/" + emojiCode + "/" + userID.String()
	_, err := client.Delete(ratelimiter, endpoint)
	return err
}

// ReqGetReactionParams https://discordapp.com/developers/docs/resources/channel#get-reactions-query-string-params
type ReqGetReactionParams struct {
	Before snowflake.ID `urlparam:"before,omitempty"` // get users before this user ID
	After  snowflake.ID `urlparam:"after,omitempty"`  // get users after this user ID
	Limit  int          `urlparam:"limit,omitempty"`  // max number of users to return (1-100)
}

// getQueryString this ins't really pretty, but it works.
func (params *ReqGetReactionParams) getQueryString() string {
	seperator := "?"
	query := ""

	if !params.Before.Empty() {
		query += seperator + params.Before.String()
		seperator = "&"
	}

	if !params.After.Empty() {
		query += seperator + params.After.String()
		seperator = "&"
	}

	if params.Limit > 0 {
		query += seperator + strconv.Itoa(params.Limit)
	}

	return query
}

// ReqGetReaction [GET]		Get a list of users that reacted with this emoji. Returns an array of user
// 							objects on success.
// Endpoint				   	/channels/{channel.id}/messages/{message.id}/reactions/{emoji}
// Rate limiter [MAJOR]	   	/channels/{channel.id}
// Discord documentation   	https://discordapp.com/developers/docs/resources/channel#get-reactions
// Reviewed				   	2018-06-07
// Comment				   	-
// emoji either unicode (string) or *Emoji with an snowflake ID if it's custom
func ReqGetReaction(client request.DiscordGetter, channelID, messageID snowflake.ID, emoji interface{}, params *ReqGetReactionParams) ([]*User, error) {
	if channelID.Empty() {
		return nil, errors.New("channelID must be set to target the correct channel")
	}
	if messageID.Empty() {
		return nil, errors.New("messageID must be set to target the specific channel message")
	}
	if emoji == nil {
		return nil, errors.New("emoji must be set in order to create a message reaction")
	}

	ratelimiter := "/channels/" + channelID.String()
	emojiCode := ""
	if _, ok := emoji.(*Emoji); ok {
		emojiCode = emoji.(*Emoji).ID.String()
	} else if _, ok := emoji.(string); ok {
		emojiCode = emoji.(string) // unicode
	} else {
		return nil, errors.New("emoji type can only be a unicode string or a *Emoji struct")
	}

	query := ""
	if params != nil {
		query += params.getQueryString()
	}
	endpoint := ratelimiter + "/messages/" + messageID.String() + "/reactions/" + emojiCode + query
	var users []*User
	_, err := client.Get(ratelimiter, endpoint, users)
	return users, err
}

// ReqCreateReaction [DELETE]	Deletes all reactions on a message. This endpoint requires the 'MANAGE_MESSAGES'
// 								permission to be present on the current user.
// Endpoint				   		/channels/{channel.id}/messages/{message.id}/reactions
// Rate limiter [MAJOR]	   		/channels/{channel.id}
// Discord documentation   		https://discordapp.com/developers/docs/resources/channel#delete-all-reactions
// Reviewed				   		2018-06-07
// Comment				   		-
// emoji either unicode (string) or *Emoji with an snowflake ID if it's custom
func ReqDeleteAllReactions(client request.DiscordDeleter, channelID, messageID snowflake.ID) error {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if messageID.Empty() {
		return errors.New("messageID must be set to target the specific channel message")
	}

	ratelimiter := "/channels/" + channelID.String()
	endpoint := ratelimiter + "/messages/" + messageID.String() + "/reactions"
	_, err := client.Delete(ratelimiter, endpoint)
	return err
}
