package rest

import "testing"
import . "github.com/andersfylling/snowflake"

func TestGetCurrentUser(t *testing.T)    {}
func TestGetUser(t *testing.T)           {}
func TestModifyCurrentUser(t *testing.T) {}
func TestGetCurrentUserGuildsParams(t *testing.T) {
	params := &GetCurrentUserGuildsParams{}
	var wants string

	wants = ""
	verifyQueryString(t, params, wants)

	s := "438543957"
	params.Before, _ = GetSnowflake(s)
	wants = "?before=" + s
	verifyQueryString(t, params, wants)

	params.Limit = 6
	wants += "&limit=6"
	verifyQueryString(t, params, wants)

	params.Limit = 0
	wants = "?before=" + s
	verifyQueryString(t, params, wants)
}
func TestLeaveGuild(t *testing.T)         {}
func TestUserDMs(t *testing.T)            {}
func TestCreateDM(t *testing.T)           {}
func TestCreateGroupDM(t *testing.T)      {}
func TestGetUserConnections(t *testing.T) {}
