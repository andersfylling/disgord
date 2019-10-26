package disgorderr

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	// TODO: stack?
	return &Err{
		cause: err,
		msg:   message,
	}
}

type Err struct {
	msg   string
	cause error
}

var _ error = (*Err)(nil)

func (e *Err) Error() string {
	return e.msg + ": " + e.cause.Error()
}
