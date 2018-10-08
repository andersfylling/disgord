// +build json_std

package httd

import "encoding/json"

// Unmarshal is the json unmarshalling implementation that is defined by the used build tags.
func Unmarshal(data []byte, v interface{}) error {
	if j, has := v.(json.Unmarshaler); has {
		return j.UnmarshalJSON(data)
	}
	return json.Unmarshal(data, v)
}

func Marshal(v interface{}) (data []byte, err error) {
	if j, has := v.(json.Marshaler); has {
		return j.MarshalJSON()
	}
	return json.Marshal(v)
}
