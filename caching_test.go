package disgord

import "testing"

func missingImplError(t *testing.T, structName string) {
	t.Error(structName + " does not implement its designated handler interface")
}

func TestStateCache_ImplementsStateCacher(t *testing.T) {
	if _, implemented := interface{}(&StateCache{}).(StateCacher); !implemented {
		missingImplError(t, "StateCache")
	}
}
