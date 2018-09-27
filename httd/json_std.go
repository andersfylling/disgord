// +build json_std

package httd

import "encoding/json"

func Unmarshal(data []byte, v interface{}) error {
	if j, has := v.(json.Unmarshaler); has {
		return j.UnmarshalJSON(data)
	}
	return json.Unmarshal(data, v)
}
