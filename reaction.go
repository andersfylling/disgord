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

func emojiReference(i interface{}) (string, error) {
	var emojiCode string
	if e, ok := i.(*Emoji); ok {
		if e.ID.IsZero() {
			emojiCode = e.Name
		} else {
			emojiCode = e.Name + ":" + e.ID.String()
		}
	} else if _, ok := i.(string); ok {
		emojiCode = i.(string) // unicode
		emojiCode = unwrapEmoji(emojiCode)
	} else {
		return "", errors.New("emoji type can only be a unicode string or a *Emoji struct")
	}
	return emojiCode, nil
}

func unwrapEmoji(e string) string {
	l := len(e)
	if l >= 2 && e[0] == e[l-1] && e[0] == ':' {
		// :emoji: => emoji
		e = e[1 : l-1]
	}
	return e
}

type ReactionQueryBuilder interface {
	WithContext(ctx context.Context) ReactionQueryBuilder

	// CreateReaction Create a reaction for the message. This endpoint requires the 'READ_MESSAGE_HISTORY'
	// permission to be present on the current user. Additionally, if nobody else has reacted to the message using this
	// emoji, this endpoint requires the 'ADD_REACTIONS' permission to be present on the current user. Returns a 204
	// empty response on success. The maximum request size when sending a message is 8MB.
	Create(flags ...Flag) (err error)

	// GetReaction Get a list of Users that reacted with this emoji. Returns an array of user objects on success.
	Get(params URLQueryStringer, flags ...Flag) (reactors []*User, err error)

	// DeleteOwnReaction Delete a reaction the current user has made for the message.
	// Returns a 204 empty response on success.
	DeleteOwn(flags ...Flag) (err error)

	// DeleteUserReaction Deletes another user's reaction. This endpoint requires the 'MANAGE_MESSAGES' permission
	// to be present on the current user. Returns a 204 empty response on success.
	DeleteUser(userID Snowflake, flags ...Flag) (err error)
}

func (m messageQueryBuilder) Reaction(emoji interface{}) ReactionQueryBuilder {
	return &reactionQueryBuilder{client: m.client, cid: m.cid, mid: m.mid, emoji: emoji}
}

type reactionQueryBuilder struct {
	ctx    context.Context
	client *Client
	cid    Snowflake
	mid    Snowflake
	emoji  interface{}
}

func (r reactionQueryBuilder) WithContext(ctx context.Context) ReactionQueryBuilder {
	r.ctx = ctx
	return &r
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
func (r reactionQueryBuilder) Create(flags ...Flag) error {
	if r.cid.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if r.mid.IsZero() {
		return errors.New("messageID must be set to target the specific channel message")
	}
	if r.emoji == nil {
		return errors.New("emoji must be set in order to create a message reaction")
	}

	emojiCode, err := emojiReference(r.emoji)
	if err != nil {
		return err
	}

	req := r.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodPut,
		Endpoint: endpoint.ChannelMessageReactionMe(r.cid, r.mid, emojiCode),
		Ctx:      r.ctx,
	}, flags)
	req.expectsStatusCode = http.StatusNoContent

	_, err = req.Execute()
	return err
}

// DeleteOwnReaction [REST] Delete a reaction the current user has made for the message.
// Returns a 204 empty response on success.
//  Method                  DELETE
//  Endpoint                /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
//  Discord documentation   https://discord.com/developers/docs/resources/channel#delete-own-reaction
//  Reviewed                2019-01-28
//  Comment                 emoji either unicode (string) or *Emoji with an snowflake Snowflake if it's custom
func (r reactionQueryBuilder) DeleteOwn(flags ...Flag) error {
	if r.cid.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if r.mid.IsZero() {
		return errors.New("messageID must be set to target the specific channel message")
	}
	if r.emoji == nil {
		return errors.New("emoji must be set in order to create a message reaction")
	}

	emojiCode, err := emojiReference(r.emoji)
	if err != nil {
		return err
	}

	req := r.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.ChannelMessageReactionMe(r.cid, r.mid, emojiCode),
		Ctx:      r.ctx,
	}, flags)
	req.expectsStatusCode = http.StatusNoContent

	_, err = req.Execute()
	return err
}

// DeleteUserReaction [REST] Deletes another user's reaction. This endpoint requires the 'MANAGE_MESSAGES' permission
// to be present on the current user. Returns a 204 empty response on success.
//  Method                  DELETE
//  Endpoint                /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
//  Discord documentation   https://discord.com/developers/docs/resources/channel#delete-user-reaction
//  Reviewed                2019-01-28
//  Comment                 emoji either unicode (string) or *Emoji with an snowflake Snowflake if it's custom
func (r reactionQueryBuilder) DeleteUser(userID Snowflake, flags ...Flag) error {
	if r.cid.IsZero() {
		return errors.New("channelID must be set to target the correct channel")
	}
	if r.mid.IsZero() {
		return errors.New("messageID must be set to target the specific channel message")
	}
	if r.emoji == nil {
		return errors.New("emoji must be set in order to create a message reaction")
	}
	if userID.IsZero() {
		return errors.New("UserID must be set to target the specific user reaction")
	}

	emojiCode, err := emojiReference(r.emoji)
	if err != nil {
		return err
	}

	req := r.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.ChannelMessageReactionUser(r.cid, r.mid, emojiCode, userID),
		Ctx:      r.ctx,
	}, flags)
	req.expectsStatusCode = http.StatusNoContent

	_, err = req.Execute()
	return err
}

// GetReactionURLParams https://discord.com/developers/docs/resources/channel#get-reactions-query-string-params
type GetReactionURLParams struct {
	Before Snowflake `urlparam:"before,omitempty"` // get Users before this user Snowflake
	After  Snowflake `urlparam:"after,omitempty"`  // get Users after this user Snowflake
	Limit  int       `urlparam:"limit,omitempty"`  // max number of Users to return (1-100)
}

var _ URLQueryStringer = (*GetReactionURLParams)(nil)

// GetReaction [REST] Get a list of Users that reacted with this emoji. Returns an array of user objects on success.
//  Method                  GET
//  Endpoint                /channels/{channel.id}/messages/{message.id}/reactions/{emoji}
//  Discord documentation   https://discord.com/developers/docs/resources/channel#get-reactions
//  Reviewed                2019-01-28
//  Comment                 emoji either unicode (string) or *Emoji with an snowflake Snowflake if it's custom
func (r reactionQueryBuilder) Get(params URLQueryStringer, flags ...Flag) (ret []*User, err error) {
	if r.cid.IsZero() {
		return nil, errors.New("channelID must be set to target the correct channel")
	}
	if r.mid.IsZero() {
		return nil, errors.New("messageID must be set to target the specific channel message")
	}
	if r.emoji == nil {
		return nil, errors.New("emoji must be set in order to create a message reaction")
	}

	emojiCode, err := emojiReference(r.emoji)
	if err != nil {
		return nil, err
	}

	query := ""
	if params != nil {
		query += params.URLQueryString()
	}

	req := r.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.ChannelMessageReaction(r.cid, r.mid, emojiCode) + query,
		Ctx:      r.ctx,
	}, flags)
	req.factory = func() interface{} {
		tmp := make([]*User, 0)
		return &tmp
	}

	return getUsers(req.Execute)
}
