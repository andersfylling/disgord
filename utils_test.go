// +build !integration

package disgord

import (
	"errors"
	"testing"

	"github.com/andersfylling/disgord/internal/disgorderr"
)

func TestValidateHandlerInputs(t *testing.T) {
	var testHandler Handler = func() {}
	var testMiddleware Middleware = func(i interface{}) interface{} { return nil }
	var testCtrl HandlerCtrl = &Ctrl{}
	var e *disgorderr.HandlerSpecErr

	t.Run("valid", func(t *testing.T) {
		t.Run("handler", func(t *testing.T) {
			err := ValidateHandlerInputs(testHandler)
			if err != nil {
				t.Error(err)
			}
		})
		t.Run("handler ctrl", func(t *testing.T) {
			err := ValidateHandlerInputs(testHandler, testCtrl)
			if err != nil {
				t.Error(err)
			}
		})
		t.Run("middleware handler", func(t *testing.T) {
			err := ValidateHandlerInputs(testMiddleware, testHandler)
			if err != nil {
				t.Error(err)
			}
		})
		t.Run("all", func(t *testing.T) {
			err := ValidateHandlerInputs(testMiddleware, testHandler, testCtrl)
			if err != nil {
				t.Error(err)
			}
		})
		t.Run("all multiple", func(t *testing.T) {
			err := ValidateHandlerInputs(testMiddleware, testMiddleware, testMiddleware, testHandler, testHandler, testHandler, testCtrl)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("missingHandler", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			err := ValidateHandlerInputs()
			if err != nil && errors.As(err, &e) {
				if e.Code() == disgorderr.HandlerSpecErrCodeMissingHandler {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("middleware", func(t *testing.T) {
			err := ValidateHandlerInputs(testMiddleware)
			if err != nil && errors.As(err, &e) {
				if e.Code() == disgorderr.HandlerSpecErrCodeMissingHandler {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("ctrl", func(t *testing.T) {
			err := ValidateHandlerInputs(testCtrl)
			if err != nil && errors.As(err, &e) {
				if e.Code() == disgorderr.HandlerSpecErrCodeMissingHandler {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("middleware ctrl", func(t *testing.T) {
			err := ValidateHandlerInputs(testMiddleware, testCtrl)
			if err != nil && errors.As(err, &e) {
				if e.Code() == disgorderr.HandlerSpecErrCodeMissingHandler {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
	})

	t.Run("middleware", func(t *testing.T) {
		t.Run("ctrl first", func(t *testing.T) {
			err := ValidateHandlerInputs(testCtrl, testMiddleware)
			if err != nil && errors.As(err, &e) {
				if e.Code() == disgorderr.HandlerSpecErrCodeUnexpectedMiddleware {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("handler first", func(t *testing.T) {
			err := ValidateHandlerInputs(testHandler, testMiddleware)
			if err != nil && errors.As(err, &e) {
				if e.Code() == disgorderr.HandlerSpecErrCodeUnexpectedMiddleware {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("middleware handler middleware", func(t *testing.T) {
			err := ValidateHandlerInputs(testMiddleware, testHandler, testMiddleware)
			if err != nil && errors.As(err, &e) {
				if e.Code() == disgorderr.HandlerSpecErrCodeUnexpectedMiddleware {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
	})

	t.Run("ctrl", func(t *testing.T) {
		t.Run("multiple ctrl", func(t *testing.T) {
			err := ValidateHandlerInputs(testMiddleware, testHandler, testCtrl, testCtrl)
			if err != nil && errors.As(err, &e) {
				if e.Code() == disgorderr.HandlerSpecErrCodeUnexpectedCtrl {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("ctrl handler ctrl", func(t *testing.T) {
			err := ValidateHandlerInputs(testCtrl, testHandler, testCtrl)
			if err != nil && errors.As(err, &e) {
				if e.Code() == disgorderr.HandlerSpecErrCodeUnexpectedCtrl {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
	})

	t.Run("invalid handler", func(t *testing.T) {
		testInvalidHandler := func(s Session, mc MessageCreate) {}

		t.Run("invalidHandler", func(t *testing.T) {
			err := ValidateHandlerInputs(testInvalidHandler)
			if err != nil && errors.As(err, &e) {
				if e.Code() == disgorderr.HandlerSpecErrCodeUnknownHandlerSignature {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("handler invalidHandler", func(t *testing.T) {
			err := ValidateHandlerInputs(testHandler, testInvalidHandler)
			if err != nil && errors.As(err, &e) {
				if e.Code() == disgorderr.HandlerSpecErrCodeUnknownHandlerSignature {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
	})

	t.Run("invalid ctrl", func(t *testing.T) {
		testInvalidCtrl := Ctrl{}

		t.Run("invalidCtrl", func(t *testing.T) {
			err := ValidateHandlerInputs(testInvalidCtrl)
			if err != nil && errors.As(err, &e) {
				if e.Code() == disgorderr.HandlerSpecErrCodeNotHandlerCtrlImpl {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("handler invalidCtrl", func(t *testing.T) {
			err := ValidateHandlerInputs(testHandler, testInvalidCtrl)
			if err != nil && errors.As(err, &e) {
				if e.Code() == disgorderr.HandlerSpecErrCodeNotHandlerCtrlImpl {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
	})
}
