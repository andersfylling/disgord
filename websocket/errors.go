package websocket

// NewErrorUnsupportedEventName ...
func NewErrorUnsupportedEventName(event string) *ErrorUnsupportedEventName {
	return &ErrorUnsupportedEventName{
		info: "unsupported event name '" + event + "' was given",
	}
}

// ErrorUnsupportedEventName is an error to identity unsupported event types request by the user
type ErrorUnsupportedEventName struct {
	info string
}

func (e *ErrorUnsupportedEventName) Error() string {
	return e.info
}
