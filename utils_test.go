package disgord

import (
	"testing"
)

func TestValidateHandlerInputs(t *testing.T) {
	var testHandler Handler = func() {}
	var testMiddleware Middleware = func(i interface{}) interface{} { return nil }
	var testCtrl HandlerCtrl = &Ctrl{Runs: 1}

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
			if err != nil {
				if err.Error() == "missing handler(s)" {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("middleware", func(t *testing.T) {
			err := ValidateHandlerInputs(testMiddleware)
			if err != nil {
				if err.Error() == "missing handler(s)" {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("ctrl", func(t *testing.T) {
			err := ValidateHandlerInputs(testCtrl)
			if err != nil {
				if err.Error() == "missing handler(s)" {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("middleware ctrl", func(t *testing.T) {
			err := ValidateHandlerInputs(testMiddleware, testCtrl)
			if err != nil {
				if err.Error() == "missing handler(s)" {
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
			if err != nil {
				if err.Error() == "middlewares can only be in the beginning. Grouped together" {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("handler first", func(t *testing.T) {
			err := ValidateHandlerInputs(testHandler, testMiddleware)
			if err != nil {
				if err.Error() == "middlewares can only be in the beginning. Grouped together" {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("middleware handler middleware", func(t *testing.T) {
			err := ValidateHandlerInputs(testMiddleware, testHandler, testMiddleware)
			if err != nil {
				if err.Error() == "middlewares can only be in the beginning. Grouped together" {
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
			if err != nil {
				if err.Error() == "a handlerCtrl's can only be at the end of the definition and only one" {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("ctrl handler", func(t *testing.T) {
			err := ValidateHandlerInputs(testCtrl, testHandler, testCtrl)
			if err != nil {
				if err.Error() == "a handlerCtrl's can only be at the end of the definition and only one" {
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
			if err != nil {
				if err.Error() == "invalid handler signature. General tip: no handlers can use the param type `*disgord.Session`, try `disgord.Session` instead" {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
		t.Run("handler invalidHandler", func(t *testing.T) {
			err := ValidateHandlerInputs(testHandler, testInvalidHandler)
			if err != nil {
				if err.Error() == "invalid handler signature. General tip: no handlers can use the param type `*disgord.Session`, try `disgord.Session` instead" {
					return
				}
				t.Error(err)
			}
			t.Fail()
		})
	})
}
