package util

import (
	"github.com/andersfylling/snowflake/v4"
)

type Snowflake = snowflake.Snowflake

// GetSnowflake see snowflake.GetSnowflake
func GetSnowflake(v interface{}) (Snowflake, error) {
	return snowflake.GetSnowflake(v)
}

// NewSnowflake see snowflake.NewSnowflake
func NewSnowflake(id uint64) Snowflake {
	return snowflake.NewSnowflake(id)
}

// ParseSnowflakeString see snowflake.ParseSnowflakeString
func ParseSnowflakeString(v string) Snowflake {
	return snowflake.ParseSnowflakeString(v)
}
