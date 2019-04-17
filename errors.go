package disgord

import "github.com/pkg/errors"

var errCacheUnsupportedEvt = errors.New("given event key is unsupported for this cache repository")
var errCacheJSONObjectTooSmall = errors.New("json object does not have enough keys to be parsed")
