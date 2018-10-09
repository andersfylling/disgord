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

type keys struct {
	GuildAdmin   Snowflake
	GuildDefault Snowflake
}

// testing
func createTestRequester() (*httd.Client, *keys, error) {
	if os.Getenv(constant.DisgordTestLive) != "true" {
		return nil, nil, errors.New("live testing is deactivated")
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
		return nil, nil, errors.New("missing bot token")
	}

	keys := &keys{}

	str1 := os.Getenv(constant.DisgordTestGuildDefault)
	g1, err := GetSnowflake(str1)
	if err != nil {
		return nil, nil, errors.New("missing default guild id")
	}
	keys.GuildDefault = g1

	str2 := os.Getenv(constant.DisgordTestGuildAdmin)
	g2, err := GetSnowflake(str2)
	if err != nil {
		return nil, nil, errors.New("missing admin guild id")
	}
	keys.GuildAdmin = g2

	return httd.NewClient(reqConf), keys, nil
}

func verifyQueryString(t *testing.T, params URLParameters, wants string) {
	got := params.GetQueryString()
	if got != wants {
		t.Errorf("incorrect query param string. Got '%s', wants '%s'", got, wants)
	}
}
