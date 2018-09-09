package disgord

import (
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/httd"
)

// testing
func createTestRequester() (*httd.Client, error) {
	if os.Getenv(constant.DisgordTestLive) != "true" {
		return nil, errors.New("live testing is deactivated")
	}

	reqConf := &httd.Config{
		APIVersion:         6,
		BotToken:           os.Getenv(constant.DisgordTestBot),
		UserAgentSourceURL: constant.GitHubURL,
		UserAgentVersion:   constant.Version,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		CancelRequestWhenRateLimited: false,
	}
	if reqConf.BotToken == "" {
		return nil, errors.New("missing bot token")
	}
	return httd.NewClient(reqConf), nil
}

func verifyQueryString(t *testing.T, params URLParameters, wants string) {
	got := params.GetQueryString()
	if got != wants {
		t.Errorf("incorrect query param string. Got '%s', wants '%s'", got, wants)
	}
}
