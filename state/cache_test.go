package state

import "testing"

func missingImplError(t *testing.T, structName string) {
	t.Error(structName + " does not implement its designated handler interface")
}

func TestCache_ImplementsCacher(t *testing.T) {
	cache := &Cache{}
	if _, implemented := interface{}(cache).(Cacher); !implemented {
		missingImplError(t, "Cache")
	}
}
