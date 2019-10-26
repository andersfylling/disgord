package disgord

import (
	"errors"

	"github.com/andersfylling/disgord/internal/disgorderr"
)

var errCacheUnsupportedEvt = errors.New("given event key is unsupported for this cache repository")
var errCacheJSONObjectTooSmall = errors.New("json object does not have enough keys to be parsed")

// TODO: go generate from internal/errors/*
type Err = disgorderr.Err
