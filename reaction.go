package disgord

import (
	"context"
	"errors"
	"net/http"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
)

// Reaction ...
// https://discord.com/developers/docs/resources/channel#reaction-object
type Reaction struct {
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

	reaction.Count = r.Count
	reaction.Me = r.Me

	if r.Emoji != nil {
		reaction.Emoji = r.Emoji.DeepCopy().(*Emoji)
	}
	return
}

func unwrapEmoji(e string) string {
	l := len(e)
	if l >= 2 && e[0] == e[l-1] && e[0] == ':' {
		// :emoji: => emoji
		e = e[1 : l-1]
	}
	return e
}

// CreateReaction [REST] Create a reaction for the message. This endpoint requires the 'READ_MESSAGE_HISTORY'
// permission to be present on the current user. Additionally, if nobody else has reacted to the message using this
// emoji, this endpoint requires the 'ADD_REACTIONS' permission to be present on the current user. Returns a 204 empty
// response on success. The maximum request size when sending a message is 8MB.
//  Method                  PUT
//  Endpoint                /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
//  Discord documentation   https://discord.com/developers/docs/resources/channel#create-reaction
//  Reviewed                2019-01-30
//  Comment                 emoji either unicode (string) or *Emoji with an snowflake Snowflake if it's custom
func (c *Client) CreateReaction(ctx context.Context, channelID, messageID Snowflake, emoji interface{}, flags ...Flag) (err error) {
	if channelID.IsZero() {
		err = errors.New("channelID must be set to target the correct channel")
		return
	}
	if messageID.IsZero() {
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
		emojiCode = unwrapEmoji(emojiCode)
	} else {
		err = errors.New("emoji type can only be a unicode string or a *Emoji struct")
		return
	}

	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodPut,
		Endpoint: endpoint.ChannelMessageReactionMe(channelID, messageID, emojiCode),
		Ctx:      ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// DeleteOwnReaction [REST] Delete a reaction the current user has made for the message.
// Returns a 204 empty response on success.
//  Method                  DELETE
//  Endpoint                /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
//  Discord documentation   https://discord.com/developers/docs/resources/channel#delete-own-reaction
//  Reviewed                2019-01-28
//  Comment                 emoji either unicode (string) or *Emoji with an snowflake Snowflake if it's custom
func (c *Client) DeleteOwnReaction(ctx context.Context, channelID, messageID Snowflake, emoji interface{}, flags ...Flag) (err error) {
	if channelID.IsZero() {
		err = errors.New("channelID must be set to target the correct channel")
		return
	}
	if messageID.IsZero() {
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
		emojiCode = unwrapEmoji(emojiCode)
	} else {
		return errors.New("emoji type can only be a unicode string or a *Emoji struct")
	}

	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.ChannelMessageReactionMe(channelID, messageID, emojiCode),
		Ctx:      ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// DeleteUserReaction [REST] Deletes another user's reaction. This endpoint requires the 'MANAGE_MESSAGES' permission
// to be present on the current user. Returns a 204 empty response on success.
//  Method                  DELETE
//  Endpoint                /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
//  Discord documentation   https://discord.com/developers/docs/resources/channel#delete-user-reaction
//  Reviewed                2019-01-28
//  Comment                 emoji either unicode (string) or *Emoji with an snowflake Snowflake if it's custom
func (c *Client) DeleteUserReaction(ctx context.Context, channelID, messageID, userID Snowflake, emoji interface{}, flags ...Flag) (err error) {
	if channelID.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if messageID.IsZero() {
		return errors.New("messageID must be set to target the specific channel message")
	}
	if emoji == nil {
		return errors.New("emoji must be set in order to create a message reaction")
	}
	if userID.IsZero() {
		return errors.New("userID must be set to target the specific user reaction")
	}

	emojiCode := ""
	if e, ok := emoji.(*Emoji); ok {
		emojiCode = e.Name + ":" + e.ID.String()
	} else if _, ok := emoji.(string); ok {
		emojiCode = emoji.(string) // unicode
		emojiCode = unwrapEmoji(emojiCode)
	} else {
		return errors.New("emoji type can only be a unicode string or a *Emoji struct")
	}

	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.ChannelMessageReactionUser(channelID, messageID, emojiCode, userID),
		Ctx:      ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// GetReactionURLParams https://discord.com/developers/docs/resources/channel#get-reactions-query-string-params
type GetReactionURLParams struct {
	Before Snowflake `urlparam:"before,omitempty"` // get users before this user Snowflake
	After  Snowflake `urlparam:"after,omitempty"`  // get users after this user Snowflake
	Limit  int       `urlparam:"limit,omitempty"`  // max number of users to return (1-100)
}

var _ URLQueryStringer = (*GetReactionURLParams)(nil)

// GetReaction [REST] Get a list of users that reacted with this emoji. Returns an array of user objects on success.
//  Method                  GET
//  Endpoint                /channels/{channel.id}/messages/{message.id}/reactions/{emoji}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#get-reactions
//  Reviewed                2019-01-28
//  Comment                 emoji either unicode (string) or *Emoji with an snowflake Snowflake if it's custom
func (c *Client) GetReaction(ctx context.Context, channelID, messageID Snowflake, emoji interface{}, params URLQueryStringer, flags ...Flag) (ret []*User, err error) {
	if channelID.IsZero() {
		err = errors.New("channelID must be set to target the correct channel")
		return
	}
	if messageID.IsZero() {
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
		emojiCode = unwrapEmoji(emojiCode)
	} else {
		return nil, errors.New("emoji type can only be a unicode string or a *Emoji struct")
	}

	query := ""
	if params != nil {
		query += params.URLQueryString()
	}

	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.ChannelMessageReaction(channelID, messageID, emojiCode) + query,
		Ctx:      ctx,
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
//  Discord documentation   https://discord.com/developers/docs/resources/channel#delete-all-reactions
//  Reviewed                2019-01-28
//  Comment                 emoji either unicode (string) or *Emoji with an snowflake Snowflake if it's custom
func (c *Client) DeleteAllReactions(ctx context.Context, channelID, messageID Snowflake, flags ...Flag) (err error) {
	if channelID.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if messageID.IsZero() {
		return errors.New("messageID must be set to target the specific channel message")
	}

	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.ChannelMessageReactions(channelID, messageID),
		Ctx:      ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}
