package disgord

import (
	"errors"
	"net/http"
	"time"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
)

// Reaction ...
// https://discordapp.com/developers/docs/resources/channel#reaction-object
type Reaction struct {
	Lockable `json:"-"`

	Count uint          `json:"count"`
	Me    bool          `json:"me"`
	Emoji *PartialEmoji `json:"Emoji"`
}

var _ Reseter = (*Reaction)(nil)

// DeepCopy see interface at struct.go#DeepCopier
func (r *Reaction) DeepCopy() (copy interface{}) {
	copy = &Reaction{}
	r.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (r *Reaction) CopyOverTo(other interface{}) (err error) {
	var reaction *Reaction
	var valid bool
	if reaction, valid = other.(*Reaction); !valid {
		err = newErrorUnsupportedType("given interface{} is not of type *Reaction")
		return
	}

	if constant.LockedMethods {
		r.RLock()
		reaction.Lock()
	}

	reaction.Count = r.Count
	reaction.Me = r.Me

	if r.Emoji != nil {
		reaction.Emoji = r.Emoji.DeepCopy().(*Emoji)
	}

	if constant.LockedMethods {
		r.RUnlock()
		reaction.Unlock()
	}
	return
}

func reactionEndpointRLAdjuster(d time.Duration) time.Duration {
	if d.Seconds() <= 2 { // the time diff is not accurate at all.. might be 1s or 2s.
		d = time.Duration(250) * time.Millisecond // 1/250ms
	}
	return d
}

// CreateReaction [REST] Create a reaction for the message. This endpoint requires the 'READ_MESSAGE_HISTORY'
// permission to be present on the current user. Additionally, if nobody else has reacted to the message using this
// emoji, this endpoint requires the 'ADD_REACTIONS' permission to be present on the current user. Returns a 204 empty
// response on success. The maximum request size when sending a message is 8MB.
//  Method                  PUT
//  Endpoint                /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
//  Rate limiter [MAJOR]    /channels/{channel.id}/messages/reactions
//  Discord documentation   https://discordapp.com/developers/docs/resources/channel#create-reaction
//  Reviewed                2019-01-30
//  Comment                 emoji either unicode (string) or *Emoji with an snowflake Snowflake if it's custom
func (c *Client) CreateReaction(channelID, messageID Snowflake, emoji interface{}, flags ...Flag) (err error) {
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
	if e, ok := emoji.(*Emoji); ok {
		emojiCode = e.Name + ":" + e.ID.String()
	} else if _, ok := emoji.(string); ok {
		emojiCode = emoji.(string) // unicode
	} else {
		err = errors.New("emoji type can only be a unicode string or a *Emoji struct")
		return
	}

	r := c.newRESTRequest(&httd.Request{
		Method:            http.MethodPut,
		Ratelimiter:       ratelimitChannelMessages(channelID) + "/reactions",
		Endpoint:          endpoint.ChannelMessageReactionMe(channelID, messageID, emojiCode),
		RateLimitAdjuster: reactionEndpointRLAdjuster,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// DeleteOwnReaction [REST] Delete a reaction the current user has made for the message.
// Returns a 204 empty response on success.
//  Method                  DELETE
//  Endpoint                /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
//  Rate limiter [MAJOR]    /channels/{channel.id}/messages/reactions
//  Discord documentation   https://discordapp.com/developers/docs/resources/channel#delete-own-reaction
//  Reviewed                2019-01-28
//  Comment                 emoji either unicode (string) or *Emoji with an snowflake Snowflake if it's custom
func (c *Client) DeleteOwnReaction(channelID, messageID Snowflake, emoji interface{}, flags ...Flag) (err error) {
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
	if e, ok := emoji.(*Emoji); ok {
		emojiCode = e.Name + ":" + e.ID.String()
	} else if _, ok := emoji.(string); ok {
		emojiCode = emoji.(string) // unicode
	} else {
		return errors.New("emoji type can only be a unicode string or a *Emoji struct")
	}

	r := c.newRESTRequest(&httd.Request{
		Method:            http.MethodDelete,
		Ratelimiter:       ratelimitChannelMessages(channelID) + "/reactions",
		Endpoint:          endpoint.ChannelMessageReactionMe(channelID, messageID, emojiCode),
		RateLimitAdjuster: reactionEndpointRLAdjuster,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// DeleteUserReaction [REST] Deletes another user's reaction. This endpoint requires the 'MANAGE_MESSAGES' permission
// to be present on the current user. Returns a 204 empty response on success.
//  Method                  DELETE
//  Endpoint                /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
//  Rate limiter [MAJOR]    /channels/{channel.id}/messages/reactions
//  Discord documentation   https://discordapp.com/developers/docs/resources/channel#delete-user-reaction
//  Reviewed                2019-01-28
//  Comment                 emoji either unicode (string) or *Emoji with an snowflake Snowflake if it's custom
func (c *Client) DeleteUserReaction(channelID, messageID, userID Snowflake, emoji interface{}, flags ...Flag) (err error) {
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
	if e, ok := emoji.(*Emoji); ok {
		emojiCode = e.Name + ":" + e.ID.String()
	} else if _, ok := emoji.(string); ok {
		emojiCode = emoji.(string) // unicode
	} else {
		return errors.New("emoji type can only be a unicode string or a *Emoji struct")
	}

	r := c.newRESTRequest(&httd.Request{
		Method:            http.MethodDelete,
		Ratelimiter:       ratelimitChannelMessages(channelID) + "/reactions",
		Endpoint:          endpoint.ChannelMessageReactionUser(channelID, messageID, emojiCode, userID),
		RateLimitAdjuster: reactionEndpointRLAdjuster,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// GetReactionURLParams https://discordapp.com/developers/docs/resources/channel#get-reactions-query-string-params
type GetReactionURLParams struct {
	Before Snowflake `urlparam:"before,omitempty"` // get users before this user Snowflake
	After  Snowflake `urlparam:"after,omitempty"`  // get users after this user Snowflake
	Limit  int       `urlparam:"limit,omitempty"`  // max number of users to return (1-100)
}

var _ URLQueryStringer = (*GetReactionURLParams)(nil)

// GetReaction [REST] Get a list of users that reacted with this emoji. Returns an array of user objects on success.
//  Method                  GET
//  Endpoint                /channels/{channel.id}/messages/{message.id}/reactions/{emoji}
//  Rate limiter [MAJOR]    /channels/{channel.id}/messages/reactions
//  Discord documentation   https://discordapp.com/developers/docs/resources/channel#get-reactions
//  Reviewed                2019-01-28
//  Comment                 emoji either unicode (string) or *Emoji with an snowflake Snowflake if it's custom
func (c *Client) GetReaction(channelID, messageID Snowflake, emoji interface{}, params URLQueryStringer, flags ...Flag) (ret []*User, err error) {
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
	if e, ok := emoji.(*Emoji); ok {
		emojiCode = e.Name + ":" + e.ID.String()
	} else if _, ok := emoji.(string); ok {
		emojiCode = emoji.(string) // unicode
	} else {
		return nil, errors.New("emoji type can only be a unicode string or a *Emoji struct")
	}

	query := ""
	if params != nil {
		query += params.URLQueryString()
	}

	r := c.newRESTRequest(&httd.Request{
		Ratelimiter:       ratelimitChannelMessages(channelID) + "/reactions",
		Endpoint:          endpoint.ChannelMessageReaction(channelID, messageID, emojiCode) + query,
		RateLimitAdjuster: reactionEndpointRLAdjuster,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*User, 0)
		return &tmp
	}

	return getUsers(r.Execute)
}

// DeleteAllReactions [REST] Deletes all reactions on a message. This endpoint requires the 'MANAGE_MESSAGES'
// permission to be present on the current user.
//  Method                  DELETE
//  Endpoint                /channels/{channel.id}/messages/{message.id}/reactions
//  Rate limiter [MAJOR]    /channels/{channel.id}/messages/reactions
//  Discord documentation   https://discordapp.com/developers/docs/resources/channel#delete-all-reactions
//  Reviewed                2019-01-28
//  Comment                 emoji either unicode (string) or *Emoji with an snowflake Snowflake if it's custom
func (c *Client) DeleteAllReactions(channelID, messageID Snowflake, flags ...Flag) (err error) {
	if channelID.Empty() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if messageID.Empty() {
		return errors.New("messageID must be set to target the specific channel message")
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodDelete,
		Ratelimiter: ratelimitChannelMessages(channelID) + "/reactions",
		Endpoint:    endpoint.ChannelMessageReactions(channelID, messageID),
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}
