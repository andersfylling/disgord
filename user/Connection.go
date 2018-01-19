package user

import (
	"encoding/json"

	"github.com/andersfylling/snowflake"
)

// Connection The connection object that the user has attached.
// https://discordapp.com/developers/docs/resources/user#avatar-data
// WARNING! Due to dependency issues, the Integrations (array) refers to Integration IDs only!
//          It breaks the lib pattern, but there's nothing I can do. To retrieve the Integration
//          Object, use *disgord.Client.Integration(id) (*Integration, error)
type Connection struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Revoked bool   `json:"revoked"`

	// Since this does not hold real guild.Integration objects we need to empty
	// the integration slice, it's misguiding. But at least the output is "correct".
	// the UnmarshalJSON method is a hack for input
	Integrations []snowflake.ID `json:"-"`
}

func (conn *Connection) UnmarshalJSON(b []byte) error {
	mock := struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		Revoked      bool   `json:"revoked"`
		Integrations []struct {
			ID snowflake.ID `json:"id"`
		} `json:"integrations"`
	}{}

	err := json.Unmarshal(b, &mock)
	if err != nil {
		return err
	}

	conn.ID = mock.ID
	conn.Name = mock.Name
	conn.Type = mock.Type
	conn.Revoked = mock.Revoked

	// empty integration slice
	//conn.Integrations = conn.Integrations[:0]

	// set new slice size
	conn.Integrations = make([]snowflake.ID, len(mock.Integrations))

	// add new data
	for index, id := range mock.Integrations {
		conn.Integrations[index] = id.ID
	}

	return nil
}
