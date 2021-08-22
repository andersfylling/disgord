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
	WithFlags(flags ...Flag) GuildMemberQueryBuilder

	Get() (*Member, error)
	UpdateBuilder() UpdateGuildMemberBuilder
	AddRole(roleID Snowflake) error
	RemoveRole(roleID Snowflake) error
	Kick(reason string) error
	Ban(params *BanMemberParams) error
	GetPermissions() (PermissionBit, error)
}

func (g guildQueryBuilder) Member(userID Snowflake) GuildMemberQueryBuilder {
	return &guildMemberQueryBuilder{client: g.client, gid: g.gid, uid: userID}
}

type guildMemberQueryBuilder struct {
	ctx    context.Context
	flags  Flag
	client *Client
	gid    Snowflake
	uid    Snowflake
}

func (g guildMemberQueryBuilder) WithContext(ctx context.Context) GuildMemberQueryBuilder {
	g.ctx = ctx
	return &g
}

func (g guildMemberQueryBuilder) WithFlags(flags ...Flag) GuildMemberQueryBuilder {
	g.flags = mergeFlags(flags)
	return &g
}

// Get Returns a guild member object for the specified user.
func (g guildMemberQueryBuilder) Get() (*Member, error) {
	if !ignoreCache(g.flags) {
		if member, _ := g.client.cache.GetMember(g.gid, g.uid); member != nil {
			return member, nil
		}
	}

	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildMember(g.gid, g.uid),
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = func() interface{} {
		return &Member{
			GuildID: g.gid,
			UserID:  g.uid,
		}
	}

	member, err := getMember(r.Execute)
	if err != nil {
		return nil, err
	}
	member.GuildID = g.gid
	return member, nil
}

// UpdateBuilder is used to create a builder to update a guild member.
func (g guildMemberQueryBuilder) UpdateBuilder() UpdateGuildMemberBuilder {
	builder := &updateGuildMemberBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Member{
			GuildID: g.gid,
			UserID:  g.uid,
		}
	}
	builder.r.flags = g.flags
	builder.r.setup(g.client.req, &httd.Request{
		Method:      http.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildMember(g.gid, g.uid),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	// TODO: cache member changes
	return builder
}

// AddRole adds a role to a guild member. Requires the 'MANAGE_ROLES' permission.
// Returns a 204 empty response on success. Fires a Guild Member Update Gateway event.
func (g guildMemberQueryBuilder) AddRole(roleID Snowflake) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   http.MethodPut,
		Endpoint: endpoint.GuildMemberRole(g.gid, g.uid, roleID),
		Ctx:      g.ctx,
	}, g.flags)

	_, err := r.Execute()
	return err
}

// RemoveRole removes a role from a guild member. Requires the 'MANAGE_ROLES' permission.
// Returns a 204 empty response on success. Fires a Guild Member Update Gateway event.
func (g guildMemberQueryBuilder) RemoveRole(roleID Snowflake) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   http.MethodDelete,
		Endpoint: endpoint.GuildMemberRole(g.gid, g.uid, roleID),
		Ctx:      g.ctx,
	}, g.flags)

	_, err := r.Execute()
	return err
}

// Kick kicks a member from a guild. Requires 'KICK_MEMBERS' permission.
// Returns a 204 empty response on success. Fires a Guild Member Remove Gateway event.
func (g guildMemberQueryBuilder) Kick(reason string) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   http.MethodDelete,
		Endpoint: endpoint.GuildMember(g.gid, g.uid),
		Reason:   reason,
		Ctx:      g.ctx,
	}, g.flags)

	_, err := r.Execute()
	return err
}

// Ban Create a guild ban, and optionally delete previous messages sent by the banned user. Requires
// the 'BAN_MEMBERS' permission. Returns a 204 empty response on success. Fires a Guild Ban Add Gateway event.
func (g guildMemberQueryBuilder) Ban(params *BanMemberParams) (err error) {
	if params == nil {
		return errors.New("params was nil")
	}
	if err = params.FindErrors(); err != nil {
		return err
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:   http.MethodPut,
		Endpoint: endpoint.GuildBan(g.gid, g.uid) + params.URLQueryString(),
		Ctx:      g.ctx,
		Reason:   params.Reason,
	}, g.flags)

	_, err = r.Execute()
	return
}

// GetPermissions is used to return the members permissions.
func (g guildMemberQueryBuilder) GetPermissions() (PermissionBit, error) {
	member, err := g.Get()
	if err != nil {
		return 0, err
	}
	return member.GetPermissions(g.ctx, g.client)
}
