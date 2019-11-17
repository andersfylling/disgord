// +build json_std

package util

import (
	"encoding/json"
	"io"
)

// Unmarshal is the json unmarshalling implementation that is defined by the used build tags.
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func JSONEncode(w io.WriteCloser, v interface{}) error {
	err1 := json.NewEncoder(w).Encode(v)
	err2 := w.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

func Marshal(v interface{}) (data []byte, err error) {
	return json.Marshal(v)
}
