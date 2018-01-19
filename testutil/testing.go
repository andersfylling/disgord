package testutil

import (
	"testing"
)

func Check(err error, t *testing.T) {
	// Hide function from stacktrace
	t.Helper()
	// Assert error
	if err != nil {
		t.Error(err)
	}
}
