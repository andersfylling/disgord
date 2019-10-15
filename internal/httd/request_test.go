package httd

import (
	"testing"
)

func TestRequest_RateLimitID(t *testing.T) {
	table := map[string]string{
		"/test":                            "GET:/test",
		"/test/345345":                     "GET:/test/{id}",
		"/test/345345/lol":                 "GET:/test/{id}/lol",
		"/test/345345/lol/2652354":         "GET:/test/{id}/lol/{id}",
		"/test/345345/lol/2652354?limit=4": "GET:/test/{id}/lol/{id}",
		// major
		"/guilds/345345":                 "GET:/guilds/345345",
		"/guilds/345345/sdfsdf":          "GET:/guilds/345345/sdfsdf",
		"/guilds/345345/sdfsdf/32987234": "GET:/guilds/345345/sdfsdf/{id}",
		// major
		"/channels/345345":                 "GET:/channels/345345",
		"/channels/345345/sdfsdf":          "GET:/channels/345345/sdfsdf",
		"/channels/345345/sdfsdf/32987234": "GET:/channels/345345/sdfsdf/{id}",
		// major
		"/webhooks/345345":                 "GET:/webhooks/345345",
		"/webhooks/345345/sdfsdf":          "GET:/webhooks/345345/sdfsdf",
		"/webhooks/345345/sdfsdf/32987234": "GET:/webhooks/345345/sdfsdf/{id}",
	}

	for endpoint, wants := range table {
		r := Request{Endpoint: endpoint}
		r.PopulateMissing() // calls generator

		if r.rateLimitKey != wants {
			t.Errorf("got %s, wants %s", r.rateLimitKey, wants)
		}
	}
}
