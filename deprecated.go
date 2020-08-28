package disgord

import (
	"context"
)

// GetGuilds Deprecated use CurrentUser().GetGuilds() instead
func (c *Client) GetGuilds(ctx context.Context, params *GetCurrentUserGuildsParams, flags ...Flag) ([]*Guild, error) {
	panic("removed")
}
