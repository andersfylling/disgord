package disgorderr

import (
	"fmt"
)

const (
	HandlerSpecErrCodeUnexpectedMiddleware uint8 = iota
	HandlerSpecErrCodeMissingHandler
	HandlerSpecErrCodeUnexpectedHandler
	HandlerSpecErrCodeUnexpectedCtrl
	HandlerSpecErrCodeNotHandlerCtrlImpl
	HandlerSpecErrCodeUnknownHandlerSignature
)

func NewHandlerSpecErr(code uint8, info string) error {
	return &HandlerSpecErr{info, code}
}

type HandlerSpecErr struct {
	info string
	code uint8
}

var _ error = (*HandlerSpecErr)(nil)

func (e *HandlerSpecErr) Error() string {
	return fmt.Sprintf("[%d] %s", e.code, e.info)
}

func (e *HandlerSpecErr) Code() uint8 {
	return e.code
}
