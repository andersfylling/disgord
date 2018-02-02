package event

import "testing"

func generateError(name string) string {
	return name + " does not implement interface `CallbackStackInterface`"
}

func TestCallbackStackInterfaceImplementation(t *testing.T) {
}
