package disgord

import (
	"github.com/andersfylling/disgord/httd"
)

// URLParameters converts a struct of values to a valid URL query string
type URLParameters interface {
	GetQueryString() string
}

func unmarshal(data []byte, v interface{}) error {
	return httd.Unmarshal(data, v)
}
