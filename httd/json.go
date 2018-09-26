package httd

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"

	jsoniter "github.com/json-iterator/go"
)

// unmarshalJSONIterator https://github.com/json-iterator/go
func unmarshalJSONIterator(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

// unmarshalSTD standard GoLang implementation
func unmarshalSTD(data []byte, v interface{}) (err error) {
	return json.Unmarshal(data, v)
}

func Unmarshal(data []byte, v interface{}) error {
	if j, has := v.(json.Unmarshaler); has {
		return j.UnmarshalJSON(data)
	}
	return unmarshalJSONIterator(data, v)
}

func BinaryToText(packet []byte) (text []byte, err error) {
	b := bytes.NewReader(packet)
	var r io.ReadCloser

	r, err = zlib.NewReader(b)
	if err != nil {
		return
	}
	defer r.Close()

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(r)
	if err != nil {
		return
	}

	text = buffer.Bytes()
	return
}
