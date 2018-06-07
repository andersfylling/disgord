package state

import "testing"

func missingImplError(t *testing.T, structName string) {
	t.Error(structName + " does not implement its designated handler interface")
}

func TestCache_ImplementsCacher(t *testing.T) {
	if _, implemented := interface{}(&Cache{}).(Cacher); !implemented {
		missingImplError(t, "Cache")
	}
}
