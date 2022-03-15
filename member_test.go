//go:build !integration
// +build !integration

package disgord

import (
	"testing"
	"time"

	"github.com/andersfylling/disgord/json"
)

func TestUpdateMemberParams(t *testing.T) {
	ts := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC) // 1st jan 2022
	trueBool := true

	testCases := map[string]struct {
		encoded string
		body    *UpdateMember
	}{
		"encoded future timeout time": {
			encoded: `{"mute":true,"communication_disabled_until":"2022-01-01T00:00:00.000000+00:00"}`,
			body: &UpdateMember{
				Mute:                       &trueBool,
				CommunicationDisabledUntil: &Time{Time: ts},
			},
		},
		"no timeout struct provided": {
			encoded: `{"mute":true,"deaf":true}`,
			body: &UpdateMember{
				Mute: &trueBool,
				Deaf: &trueBool,
			},
		},
		"nil timeout value": {
			encoded: `{"mute":true}`,
			body: &UpdateMember{
				Mute:                       &trueBool,
				CommunicationDisabledUntil: nil,
			},
		},
		"remove timeout": {
			encoded: `{"mute":true,"communication_disabled_until":""}`,
			body: &UpdateMember{
				Mute:                       &trueBool,
				CommunicationDisabledUntil: &Time{},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase // move to scope
		t.Run(name, func(t *testing.T) {
			bytes, err := json.Marshal(testCase.body)
			if err != nil {
				t.Fatal(err)
			}

			if resp := string(bytes); resp != testCase.encoded {
				t.Errorf("wanted %s - got %s", testCase.encoded, resp)
			}
		})
	}
}
