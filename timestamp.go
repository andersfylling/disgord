package disgord

import (
	"encoding/json"
	"fmt"
	"time"
)

// discordTimeFormat to be able to correctly convert timestamps back into json,
// we need the micro timestamp with an addition at the ending.
// time.RFC3331 does not yield an output similar to the discord timestamp input, the date is however correct.
const discordTimeFormat string = "2006-01-02T15:04:05.000000+00:00"

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalJSON(data []byte) error
}

type DiscordTimestamp time.Time

// error: https://stackoverflow.com/questions/28464711/go-strange-json-hyphen-unmarshall-error
func (t DiscordTimestamp) MarshalJSON() ([]byte, error) {
	// wrap in double qoutes for valid json parsing
	jsonReady := fmt.Sprintf("\"%s\"", t.String())

	return []byte(jsonReady), nil
}

func (t *DiscordTimestamp) UnmarshalJSON(data []byte) error {
	var ts time.Time
	err := json.Unmarshal(data, &ts)
	if err != nil {
		return err
	}

	*t = DiscordTimestamp(ts)
	return nil
}

// String converts the timestamp into a discord formatted timestamp. time.RFC3331 does not suffice
func (t DiscordTimestamp) String() string {
	return t.Time().Format(discordTimeFormat)
}

// Time converts the DiscordTimestamp into a time.Time type.......
func (t DiscordTimestamp) Time() time.Time {
	return time.Time(t)
}
