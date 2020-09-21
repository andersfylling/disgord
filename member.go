package disgord

import (
	"context"
	"errors"
	"net/http"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
)

type GuildMemberQueryBuilder interface {
	WithContext(ctx context.Context) GuildMemberQueryBuilder

	Get(flags ...Flag) (*Member, error)
	Update(flags ...Flag) UpdateGuildMemberBuilder
	AddRole(roleID Snowflake, flags ...Flag) error
	RemoveRole(roleID Snowflake, flags ...Flag) error
	Kick(reason string, flags ...Flag) error
	Ban(params *BanMemberParams, flags ...Flag) error
}

func (g guildQueryBuilder) Member(userID Snowflake) GuildMemberQueryBuilder {
	return &guildMemberQueryBuilder{client: g.client, gid: g.gid, uid: userID}
}

type guildMemberQueryBuilder struct {
	client *Client
	gid    Snowflake
	uid    Snowflake
	ctx    context.Context
}

func (g guildMemberQueryBuilder) WithContext(ctx context.Context) GuildMemberQueryBuilder {
	g.ctx = ctx
	return &g
}

// GetMember Returns a guild member object for the specified user.
func (g guildMemberQueryBuilder) Get(flags ...Flag) (member *Member, err error) {
	if member, _ = g.client.cache.GetMember(g.gid, g.uid); member != nil {
		return member, nil
	}

	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildMember(g.gid, g.uid),
		Ctx:      g.ctx,
	}, flags)
	r.factory = func() interface{} {
		return &Member{
			GuildID: g.gid,
			UserID:  g.uid,
		}
	}

	member, err = getMember(r.Execute)
	if err != nil {
		return
	}
	member.GuildID = g.gid
	return
}

// UpdateMember is used to create a builder to update a guild member.
func (g guildMemberQueryBuilder) Update(flags ...Flag) UpdateGuildMemberBuilder {
	builder := &updateGuildMemberBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Member{
			GuildID: g.gid,
			UserID:  g.uid,
		}
	}
	builder.r.flags = flags
	builder.r.setup(g.client.req, &httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildMember(g.gid, g.uid),
		ContentType: httd.ContentTypeJSON,
	}, func(resp *http.Response, body []byte, err error) error {
		if resp.StatusCode != http.StatusNoContent {
			msg := "could not change attributes of member. Does the member exist, and do you have permissions?"
			return errors.New(msg)
		}
		return nil
	})

	// TODO: cache member changes
	return builder
}

// AddGuildMemberRole adds a role to a guild member. Requires the 'MANAGE_ROLES' permission.
// Returns a 204 empty response on success. Fires a Guild Member Update Gateway event.
func (g guildMemberQueryBuilder) AddRole(roleID Snowflake, flags ...Flag) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodPut,
		Endpoint: endpoint.GuildMemberRole(g.gid, g.uid, roleID),
		Ctx:      g.ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err := r.Execute()
	return err
}

// RemoveMemberRole removes a role from a guild member. Requires the 'MANAGE_ROLES' permission.
// Returns a 204 empty response on success. Fires a Guild Member Update Gateway event.
func (g guildMemberQueryBuilder) RemoveRole(roleID Snowflake, flags ...Flag) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.GuildMemberRole(g.gid, g.uid, roleID),
		Ctx:      g.ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err := r.Execute()
	return err
}

// KickMember kicks a member from a guild. Requires 'KICK_MEMBERS' permission.
// Returns a 204 empty response on success. Fires a Guild Member Remove Gateway event.
func (g guildMemberQueryBuilder) Kick(reason string, flags ...Flag) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.GuildMember(g.gid, g.uid),
		Reason:   reason,
		Ctx:      g.ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err := r.Execute()
	return err
}

// BanMember Create a guild ban, and optionally delete previous messages sent by the banned user. Requires
// the 'BAN_MEMBERS' permission. Returns a 204 empty response on success. Fires a Guild Ban Add Gateway event.
func (g guildMemberQueryBuilder) Ban(params *BanMemberParams, flags ...Flag) (err error) {
	if params == nil {
		return errors.New("params was nil")
	}
	if err = params.FindErrors(); err != nil {
		return err
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodPut,
		Endpoint: endpoint.GuildBan(g.gid, g.uid) + params.URLQueryString(),
		Ctx:      g.ctx,
		Reason:   params.Reason,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return
}
