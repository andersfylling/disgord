package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	. "github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/disgord/rest/httd"
	"github.com/andersfylling/snowflake"
)

// CreateReaction [PUT]     Create a reaction for the message. This endpoint requires the 'READ_MESSAGE_HISTORY'
//                          permission to be present on the current user. Additionally, if nobody else has
//                          reacted to the message using this emoji, this endpoint requires the 'ADD_REACTIONS'
//                          permission to be present on the current user. Returns a 204 empty response on success.
//                          The maximum request size when sending a message is 8MB.
// Endpoint                 /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
// Rate limiter [MAJOR]     /channels/{channel.id}/messages TODO: I have no idea what the key is
// Discord documentation    https://discordapp.com/developers/docs/resources/channel#create-reaction
// Reviewed                 2018-06-07
// Comment                  -
// emoji either unicode (string) or *Emoji with an snowflake ID if it's custom
func CreateReaction(client httd.Puter, channelID, messageID snowflake.ID, emoji interface{}) (ret *Reaction, err error) {
	if channelID.Empty() {
		err = errors.New("channelID must be set to target the correct channel")
		return
	}
	if messageID.Empty() {
		err = errors.New("messageID must be set to target the specific channel message")
		return
	}
	if emoji == nil {
		err = errors.New("emoji must be set in order to create a message reaction")
		return
	}

	emojiCode := ""
	if _, ok := emoji.(*Emoji); ok {
		emojiCode = emoji.(*Emoji).ID.String()
	} else if _, ok := emoji.(string); ok {
		emojiCode = emoji.(string) // unicode
	} else {
		err = errors.New("emoji type can only be a unicode string or a *Emoji struct")
		return
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannelMessages(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/messages/" + messageID.String() + "/reactions/" + emojiCode + "/@me",
	}
	_, body, err := client.Put(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}

// DeleteOwnReaction [DELETE]   Delete a reaction the current user has made for the message. Returns a 204
//                              empty response on success.
// Endpoint                     /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
// Rate limiter [MAJOR]         /channels/{channel.id}/messages [DELETE] TODO: I have no idea what the key is
// Discord documentation        https://discordapp.com/developers/docs/resources/channel#delete-own-reaction
// Reviewed                     2018-06-07
// Comment                      emoji either unicode (string) or *Emoji with an snowflake ID if it's custom
func DeleteOwnReaction(client httd.Deleter, channelID, messageID snowflake.ID, emoji interface{}) (err error) {
	if channelID.Empty() {
		err = errors.New("channelID must be set to target the correct channel")
		return
	}
	if messageID.Empty() {
		err = errors.New("messageID must be set to target the specific channel message")
		return
	}
	if emoji == nil {
		err = errors.New("emoji must be set in order to create a message reaction")
		return
	}

	emojiCode := ""
	if _, ok := emoji.(*Emoji); ok {
		emojiCode = emoji.(*Emoji).ID.String()
	} else if _, ok := emoji.(string); ok {
		emojiCode = emoji.(string) // unicode
	} else {
		return errors.New("emoji type can only be a unicode string or a *Emoji struct")
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannelMessagesDelete(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/messages/" + messageID.String() + "/reactions/" + emojiCode + "/@me",
	}
	resp, _, err := client.Delete(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(http.StatusNoContent)
		err = errors.New(msg)
	}
	return
}

// CreateReaction [DELETE]  Deletes another user's reaction. This endpoint requires the 'MANAGE_MESSAGES'
//                          permission to be present on the current user. Returns a 204 empty response on success.
// Endpoint                 /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
// Rate limiter [MAJOR]     /channels/{channel.id}/messages [DELETE] TODO: I have no idea if this is the correct key
// Discord documentation    https://discordapp.com/developers/docs/resources/channel#delete-user-reaction
// Reviewed                 2018-06-07
// Comment                  emoji either unicode (string) or *Emoji with an snowflake ID if it's custom
func DeleteUserReaction(client httd.Deleter, channelID, messageID, userID snowflake.ID, emoji interface{}) (err error) {
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

	emojiCode := ""
	if _, ok := emoji.(*Emoji); ok {
		emojiCode = emoji.(*Emoji).ID.String()
	} else if _, ok := emoji.(string); ok {
		emojiCode = emoji.(string) // unicode
	} else {
		return errors.New("emoji type can only be a unicode string or a *Emoji struct")
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannelMessagesDelete(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/messages/" + messageID.String() + "/reactions/" + emojiCode + "/" + userID.String(),
	}
	resp, _, err := client.Delete(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(http.StatusNoContent)
		err = errors.New(msg)
	}
	return
}

// ReqGetReactionParams https://discordapp.com/developers/docs/resources/channel#get-reactions-query-string-params
type GetReactionParams struct {
	Before snowflake.ID `urlparam:"before,omitempty"` // get users before this user ID
	After  snowflake.ID `urlparam:"after,omitempty"`  // get users after this user ID
	Limit  int          `urlparam:"limit,omitempty"`  // max number of users to return (1-100)
}

// getQueryString this ins't really pretty, but it works.
func (params *GetReactionParams) getQueryString() string {
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

// ReqGetReaction [GET]   Get a list of users that reacted with this emoji. Returns an array of user
//                        objects on success.
// Endpoint               /channels/{channel.id}/messages/{message.id}/reactions/{emoji}
// Rate limiter [MAJOR]   /channels/{channel.id}/messages TODO: I have no idea if this is the correct key
// Discord documentation  https://discordapp.com/developers/docs/resources/channel#get-reactions
// Reviewed               2018-06-07
// Comment                -
// emoji either unicode (string) or *Emoji with an snowflake ID if it's custom
func GetReaction(client httd.Getter, channelID, messageID snowflake.ID, emoji interface{}, params *GetReactionParams) (ret []*User, err error) {
	if channelID.Empty() {
		err = errors.New("channelID must be set to target the correct channel")
		return
	}
	if messageID.Empty() {
		err = errors.New("messageID must be set to target the specific channel message")
		return
	}
	if emoji == nil {
		err = errors.New("emoji must be set in order to create a message reaction")
		return
	}

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

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannelMessages(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/messages/" + messageID.String() + "/reactions/" + emojiCode + query,
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}

// ReqCreateReaction [DELETE] Deletes all reactions on a message. This endpoint requires the 'MANAGE_MESSAGES'
//                            permission to be present on the current user.
// Endpoint                   /channels/{channel.id}/messages/{message.id}/reactions
// Rate limiter [MAJOR]       /channels/{channel.id}/messages [DELETE] TODO: I have no idea if this is the correct key
// Discord documentation      https://discordapp.com/developers/docs/resources/channel#delete-all-reactions
// Reviewed                   2018-06-07
// Comment                    -
// emoji either unicode (string) or *Emoji with an snowflake ID if it's custom
func DeleteAllReactions(client httd.Deleter, channelID, messageID snowflake.ID) (err error) {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if messageID.Empty() {
		return errors.New("messageID must be set to target the specific channel message")
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitChannelMessagesDelete(channelID),
		Endpoint:    "/channels/" + channelID.String() + "/messages/" + messageID.String() + "/reactions",
	}
	resp, _, err := client.Delete(details)
	if err != nil {
		return
	}

	// TODO: what is the response on a successful execution?
	if false && resp.StatusCode != http.StatusNoContent {
		msg := "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(http.StatusNoContent)
		err = errors.New(msg)
	}
	return
}
