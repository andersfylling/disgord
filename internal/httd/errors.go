package httd

import "time"

type Error struct {
	message string
	reset   time.Time
}

var _ error = (*Error)(nil)

func (e *Error) Error() string {
	return e.message
}

var (
	ErrRateLimited error = &Error{"rate limited", time.Unix(0, 0)}
)
