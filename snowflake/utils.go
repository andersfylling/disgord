package snowflake

import "strconv"

// ParseSnowflakeString interprets a string with a decimal number.
//         Note that in contrast to ParseUint, this function assumes the given string is
//         always valid and thus will panic rather than return an error.
//         This should only be used on checks that can be done at compile time,
//         unless you want to trust other modules to returns valid data.
func ParseSnowflakeString(v string) Snowflake {
	id, err := ParseSnowflakeUint(v, 10)
	if err != nil {
		panic(err) // TODO
	}
	return id
}

// ParseUint converts a string and given base to a Snowflake
func ParseSnowflakeUint(v string, base int) (Snowflake, error) {
	if v == "" {
		return Snowflake(0), nil
	}

	id, err := strconv.ParseUint(v, base, 64)
	return Snowflake(id), err
}
