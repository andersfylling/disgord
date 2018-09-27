// +build !json_std

package httd

import (
	"encoding/json"

	"github.com/json-iterator/go"
)

// Unmarshal is the json unmarshalling implementation that is defined by the used build tags.
func Unmarshal(data []byte, v interface{}) error {
	if j, has := v.(json.Unmarshaler); has {
		return j.UnmarshalJSON(data)
	}
	return jsoniter.Unmarshal(data, v)
}
