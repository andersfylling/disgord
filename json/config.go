package json

import "encoding/json"

var (
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)

type RawMessage = json.RawMessage

type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}
