package rest

import (
	"encoding/json"
	"github.com/json-iterator/go"
)

type URLParameters interface {
	GetQueryString() string
}

// unmarshalJSONIterator https://github.com/json-iterator/go
func unmarshalJSONIterator(data []byte, v interface{}) (err error) {
	err = jsoniter.Unmarshal(data, v)
	return
}

// unmarshalSTD standard GoLang implementation
func unmarshalSTD(data []byte, v interface{}) (err error) {
	err = json.Unmarshal(data, v)
	return
}

func unmarshal(data []byte, v interface{}) error {
	return unmarshalJSONIterator(data, v)
}
