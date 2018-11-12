package disgord

import (
	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/ratelimit"
	"net/http"
)

// PartialInvite ...
// {
//    "code": "abc"
// }
type PartialInvite = Invite

// Invite Represents a code that when used, adds a user to a guild.
// https://discordapp.com/developers/docs/resources/invite#invite-object
// Reviewed: 2018-06-10
type Invite struct {
	Lockable `json:"-"`

	// Code the invite code (unique Snowflake)
	Code string `json:"code"`

	// Guild the guild this invite is for
	Guild *PartialGuild `json:"guild"`

	// Channel the channel this invite is for
	Channel *PartialChannel `json:"channel"`

	// ApproximatePresenceCount approximate count of online members
	ApproximatePresenceCount int `json:"approximate_presence_count,omitempty"`

	// ApproximatePresenceCount approximate count of total members
	ApproximateMemberCount int `json:"approximate_member_count,omitempty"`
}

// DeepCopy see interface at struct.go#DeepCopier
func (i *Invite) DeepCopy() (copy interface{}) {
	copy = &Invite{}
	i.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (i *Invite) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var invite *Invite
	if invite, ok = other.(*Invite); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *Invite")
		return
	}

	if constant.LockedMethods {
		i.RLock()
		invite.Lock()
	}

	invite.Code = i.Code
	invite.ApproximatePresenceCount = i.ApproximatePresenceCount
	invite.ApproximateMemberCount = i.ApproximateMemberCount

	if i.Guild != nil {
		invite.Guild = NewPartialGuild(i.Guild.ID)
	}
	if i.Channel != nil {
		c := i.Channel
		invite.Channel = &PartialChannel{
			ID:   c.ID,
			Name: c.Name,
			Type: c.Type,
		}
	}

	if constant.LockedMethods {
		i.RUnlock()
		invite.Unlock()
	}

	return
}

// InviteMetadata Object
// https://discordapp.com/developers/docs/resources/invite#invite-metadata-object
// Reviewed: 2018-06-10
type InviteMetadata struct {
	Lockable `json:"-"`

	// Inviter user who created the invite
	Inviter *User `json:"inviter"`

	// Uses number of times this invite has been used
	Uses int `json:"uses"`

	// MaxUses max number of times this invite can be used
	MaxUses int `json:"max_uses"`

	// MaxAge duration (in seconds) after which the invite expires
	MaxAge int `json:"max_age"`

	// Temporary whether this invite only grants temporary membership
	Temporary bool `json:"temporary"`

	// CreatedAt when this invite was created
	CreatedAt Timestamp `json:"created_at"`

	// Revoked whether this invite is revoked
	Revoked bool `json:"revoked"`
}

// DeepCopy see interface at struct.go#DeepCopier
func (i *InviteMetadata) DeepCopy() (copy interface{}) {
	copy = &InviteMetadata{}
	i.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (i *InviteMetadata) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var invite *InviteMetadata
	if invite, ok = other.(*InviteMetadata); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *InviteMetadata")
		return
	}

	if constant.LockedMethods {
		i.RLock()
		invite.Lock()
	}

	invite.Uses = i.Uses
	invite.MaxUses = i.MaxUses
	invite.MaxAge = i.MaxAge
	invite.Temporary = i.Temporary
	invite.CreatedAt = i.CreatedAt
	invite.Revoked = i.Revoked

	if i.Inviter != nil {
		invite.Inviter = i.Inviter.DeepCopy().(*User)
	}

	if constant.LockedMethods {
		i.RUnlock()
		invite.Unlock()
	}

	return
}

// GetInvite [REST] Returns an invite object for the given code.
//  Method                  GET
//  Endpoint                /invites/{invite.code}
//  Rate limiter            /invites
//  Discord documentation   https://discordapp.com/developers/docs/resources/invite#get-invite
//  Reviewed                2018-06-10
//  Comment                 withCounts whether the invite should contain approximate member counts
func (c *Client) GetInvite(inviteCode string) (builder *getInviteBuilder) {
	builder = &getInviteBuilder{}
	builder.setup(c.req, &httd.Request{
		Method:      http.MethodGet,
		Ratelimiter: ratelimit.Invites(),
		Endpoint:    endpoint.Invite(inviteCode),
	}, nil)

	return builder
}

type getInviteBuilder struct {
	RESTRequestBuilder
}

func (b *getInviteBuilder) AsInviteCode() *getInviteBuilder {
	b.urlParams["with_counts"] = true
	return b
}

func (b *getInviteBuilder) Execute() (invite *Invite, err error) {
	err = b.execute(invite)
	return
}

// DeleteInvite [REST] Delete an invite. Requires the MANAGE_CHANNELS permission. Returns an invite object on success.
//  Method                  DELETE
//  Endpoint                /invites/{invite.code}
//  Rate limiter            /invites
//  Discord documentation   https://discordapp.com/developers/docs/resources/invite#delete-invite
//  Reviewed                2018-06-10
//  Comment                 -
func (c *Client) DeleteInvite(inviteCode string) (builder *deleteInviteBuilder) {
	builder = &deleteInviteBuilder{}
	builder.setup(c.req, &httd.Request{
		Method:      http.MethodDelete,
		Ratelimiter: ratelimit.Invites(),
		Endpoint:    endpoint.Invite(inviteCode),
	}, nil)

	return builder
}

type deleteInviteBuilder struct {
	RESTRequestBuilder
}

func (b *deleteInviteBuilder) Execute() (invite *Invite, err error) {
	err = b.execute(invite)
	return
}
