package testutil

import (
	"testing"
)

func Check(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}
