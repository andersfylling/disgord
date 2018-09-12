package disgord

import "testing"

func TestError_InterfaceImplementations(t *testing.T) {
	var u interface{} = &ErrorUnsupportedType{}

	t.Run("error", func(t *testing.T) {
		if _, ok := u.(error); !ok {
			t.Error("ErrorUnsupportedType does not implement error")
		}
	})
}
