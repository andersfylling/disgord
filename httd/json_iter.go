// +build !json_std

package httd

import (
	"encoding/json"
	"io"

	"github.com/json-iterator/go"
)

var jiter = jsoniter.Config{
	EscapeHTML:                    false,
	MarshalFloatWith6Digits:       true, // will lose precession
	ObjectFieldMustBeSimpleString: true, // do not unescape object field
	CaseSensitive:                 false,
	ValidateJsonRawMessage:        false,
	SortMapKeys:                   false,
}.Froze()

// Unmarshal is the json unmarshaler implementation that is defined by the used build tags.
func Unmarshal(data []byte, v interface{}) error {
	if j, has := v.(json.Unmarshaler); has {
		return j.UnmarshalJSON(data)
	}
	return jiter.Unmarshal(data, v)
}

func JSONEncode(w io.WriteCloser, v interface{}) error {
	err1 := jsoniter.NewEncoder(w).Encode(v)
	err2 := w.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

// Marshal is the json marshaler implementation for jsoniter, depending on the build tags.
func Marshal(v interface{}) (data []byte, err error) {
	if j, has := v.(json.Marshaler); has {
		return j.MarshalJSON()
	}
	return jiter.Marshal(v)
}
