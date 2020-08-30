package disgord

import (
	"context"
	"errors"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
	"github.com/andersfylling/disgord/internal/logger"
)

// GetGuilds Deprecated use CurrentUser().GetGuilds() instead
func (c *Client) GetGuilds(ctx context.Context, params *GetCurrentUserGuildsParams, flags ...Flag) ([]*Guild, error) {
	panic("removed")
}

// Deprecated
func DefaultLogger(debug bool) logger.Logger {
	panic("this has been removed, please see examples/docs/logging-* for more information")
}

// Deprecated
func DefaultLoggerWithInstance(log logger.Logger) logger.Logger {
	return DefaultLogger(true)
}

// GetDMChannels [REST] Returns a list of DM channel objects.
//  Method                  GET
//  Endpoint                /users/@me/channels
//  Discord documentation   https://discord.com/developers/docs/resources/user#get-user-dms
//  Reviewed                2019-02-19
//  Comment                 Apparently Discord removed support for this in 2016 and updated their docs 2 years after..
//							https://github.com/discord/discord-api-docs/issues/184
//							For now I'll just leave this here, until I can do a cache lookup. Making this cache
//							dependent.
// Deprecated: Needs cache checking to get the actual list of Channels
func (c currentUserQueryBuilder) GetDMChannels(flags ...Flag) (ret []*Channel, err error) {
	r := c.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.UserMeChannels(),
		Ctx:      c.ctx,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Channel, 0) // TODO: use channel pool to get enough Channels
		return &tmp
	}

	var vs interface{}
	if vs, err = r.Execute(); err != nil {
		return nil, err
	}

	if chans, ok := vs.(*[]*Channel); ok {
		return *chans, nil
	}
	return nil, errors.New("unable to cast guild slice")
}
