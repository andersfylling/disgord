package snowflake

import (
	"strconv"
)

// NewID creates a new Snowflake Snowflake from a uint64.
func NewSnowflake(id uint64) Snowflake {
	return Snowflake(id)
}

// Snowflake Snowflake Snowflake created by twitter
type Snowflake uint64

// JSON can be useful when sending the snowflake Snowflake by a json API
type SnowflakeJSON struct {
	ID    Snowflake `json:"id"`
	IDStr string    `json:"id_str"`
}

// Empty since snowflake exists of several parts, including a timestamp,
//       I assume a valid snowflake Snowflake is never 0.
func (s Snowflake) Empty() bool {
	return uint64(s) == 0
}

// JSONStruct returns a struct that can be embedded in other structs.
//            This is useful if you have a API server, since js can't parse uint64.
//            Therefore there must a snowflake Snowflake string.
func (s Snowflake) JSONStruct() *SnowflakeJSON {
	return &SnowflakeJSON{
		ID:    s,
		IDStr: s.String(),
	}
}

// String returns the decimal representation of the snowflake Snowflake.
func (s Snowflake) String() string {
	return strconv.FormatUint(uint64(s), 10)
}

// HexString converts the Snowflake into a hexadecimal string
func (s Snowflake) HexString() string {
	return strconv.FormatUint(uint64(s), 16)
}

// HexPrettyString converts the Snowflake into a hexadecimal string with the hex prefix 0x
func (s Snowflake) HexPrettyString() string {
	return "0x" + strconv.FormatUint(uint64(s), 16)
}

// MarshalBinary create a binary literal representation as a string
func (s Snowflake) MarshalBinary() (data []byte, err error) {
	return []byte(strconv.FormatUint(uint64(s), 2)), nil
}

func (s *Snowflake) UnmarshalBinary(text []byte) (err error) {
	id, err := strconv.ParseUint(string(text), 2, 64)
	if err != nil {
		return
	}
	*s = Snowflake(id)

	return
}

func (s Snowflake) MarshalText() (text []byte, err error) {
	text = []byte(s.String())
	err = nil

	return
}

func (s *Snowflake) UnmarshalText(text []byte) (err error) {
	id, err := strconv.ParseUint(string(text), 10, 64)
	if err != nil {
		return
	}
	*s = Snowflake(id)

	return
}
