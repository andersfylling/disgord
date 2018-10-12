// +build !json_std

package httd

import (
	"encoding/json"

	"github.com/json-iterator/go"
)

// Unmarshal is the json unmarshaler implementation that is defined by the used build tags.
func Unmarshal(data []byte, v interface{}) error {
	if j, has := v.(json.Unmarshaler); has {
		return j.UnmarshalJSON(data)
	}
	return jsoniter.Unmarshal(data, v)
}

// Marshal is the json marshaler implementation for jsoniter, depending on the build tags.
func Marshal(v interface{}) (data []byte, err error) {
	if j, has := v.(json.Marshaler); has {
		return j.MarshalJSON()
	}
	return jsoniter.Marshal(v)
}
