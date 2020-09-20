package disgord

import (
	"github.com/andersfylling/disgord/internal/disgorderr"
)

// TODO: go generate from internal/errors/*
type Err = disgorderr.Err
type CloseConnectionErr = disgorderr.ClosedConnectionErr
type HandlerSpecErr = disgorderr.HandlerSpecErr
