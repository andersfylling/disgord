// +build !integration

package disgord

import (
	"fmt"
	"testing"
)

func TestUserPresence_InterfaceImplementations(t *testing.T) {
	var u interface{} = &UserPresence{}

	t.Run("Stringer", func(t *testing.T) {
		if _, ok := u.(fmt.Stringer); !ok {
			t.Error("UserPresence does not implement fmt.Stringer")
		}
	})

	t.Run("DeepCopier", func(t *testing.T) {
		if _, ok := u.(DeepCopier); !ok {
			t.Error("UserPresence does not implement DeepCopier")
		}
	})

	t.Run("Copier", func(t *testing.T) {
		if _, ok := u.(Copier); !ok {
			t.Error("UserPresence does not implement Copier")
		}
	})
}

func TestGetCurrentUserGuildsParams(t *testing.T) {
	params := &getCurrentUserGuildsBuilder{}
	params.r.setup(nil, nil, nil)
	var wants string

	wants = ""
	verifyQueryString(t, params.r.urlParams, wants)

	wants = "?before=438543957"
	params.SetBefore(438543957)
	verifyQueryString(t, params.r.urlParams, wants)

	wants += "&limit=6"
	params.SetLimit(6)
	verifyQueryString(t, params.r.urlParams, wants)

	wants = "?before=438543957"
	params.SetDefaultLimit()
	verifyQueryString(t, params.r.urlParams, wants)
}
