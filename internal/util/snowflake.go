package util

import (
	"github.com/andersfylling/snowflake/v5"
)

type Snowflake = snowflake.Snowflake

// GetSnowflake see snowflake.GetSnowflake
func GetSnowflake(v interface{}) (Snowflake, error) {
	return snowflake.GetSnowflake(v)
}

// ParseSnowflakeString see snowflake.ParseSnowflakeString
func ParseSnowflakeString(v string) Snowflake {
	return snowflake.ParseSnowflakeString(v)
}
