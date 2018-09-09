// Package disgord GoLang module for interacting with the Discord API
package disgord

import (
	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/snowflake"
)

const (
	// JSONEncoding const for JSON encoding type
	JSONEncoding = "JSON"

	// APIVersion desired API version to use
	APIVersion = 6 // February 5, 2018
	// DefaultAPIVersion the default Discord API version
	DefaultAPIVersion = 6

	// GitHubURL complete url for this project
	GitHubURL = "https://github.com/andersfylling/disgord"
)

// LibraryInfo returns name + version
func LibraryInfo() string {
	return "Disgord " + constant.Version
}

// Wrapper for github.com/andersfylling/snowflake
// ------------------

// Snowflake twitter snowflake identification for Discord
type Snowflake = snowflake.Snowflake

func GetSnowflake(v interface{}) (Snowflake, error) {
	s, err := snowflake.GetSnowflake(v)
	return Snowflake(s), err
}

func NewSnowflake(id uint64) Snowflake {
	return Snowflake(snowflake.NewSnowflake(id))
}

func ParseSnowflakeString(v string) Snowflake {
	return Snowflake(snowflake.ParseSnowflakeString(v))
}
