package disgord

import (
	"testing"
)

func TestGetCurrentUser(t *testing.T) {
	client, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	_, err = GetCurrentUser(client)
	if err != nil {
		t.Error(err)
	}
}
func TestGetUser(t *testing.T) {
	client, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	userID := NewSnowflake(140413331470024704)
	user, err := GetUser(client, userID)
	if err != nil {
		t.Error(err)
	}

	if user.ID != userID {
		t.Error("user ID missmatch")
	}
}
func TestModifyCurrentUser(t *testing.T) {
	client, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	// this has been verified to work
	// however, you cannot change username often so this is
	// can give an error

	var originalUsername string
	t.Run("getting original username", func(t *testing.T) {
		user, err := GetCurrentUser(client)
		if err != nil {
			t.Error(err)
		}

		originalUsername = user.Username
	})

	t.Run("changing username", func(t *testing.T) {
		if originalUsername == "" {
			t.Skip()
			return
		}
		randomName := "sldfhksghs"
		params := &ModifyCurrentUserParams{
			Username: randomName,
		}
		_, err := ModifyCurrentUser(client, params)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("resetting username", func(t *testing.T) {
		if originalUsername == "" {
			t.Skip()
			return
		}
		params := &ModifyCurrentUserParams{
			Username: originalUsername,
		}
		_, err := ModifyCurrentUser(client, params)
		if err != nil {
			t.Error(err)
		}
	})
}
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
func TestLeaveGuild(t *testing.T) {
	// Nope. Not gonna automate this.
}
func TestUserDMs(t *testing.T) {
	// TODO
}
func TestCreateDM(t *testing.T) {
	// TODO
}
func TestCreateGroupDM(t *testing.T) {
	// TODO
}
func TestGetUserConnections(t *testing.T) {
	// Missing OAuth2
}
