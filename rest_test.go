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

// reqMocker for testing rest endpoint configurations with having to test the request logic in httd as well
type reqMocker struct {
	err  error
	body []byte
	resp *http.Response
	req  *httd.Request
}

func (gm *reqMocker) Get(req *httd.Request) (*http.Response, []byte, error) {
	gm.req = req
	return gm.resp, gm.body, gm.err
}

func (gm *reqMocker) Post(req *httd.Request) (*http.Response, []byte, error) {
	gm.req = req
	return gm.resp, gm.body, gm.err
}

func (gm *reqMocker) Patch(req *httd.Request) (*http.Response, []byte, error) {
	gm.req = req
	return gm.resp, gm.body, gm.err
}

func (gm *reqMocker) Put(req *httd.Request) (*http.Response, []byte, error) {
	gm.req = req
	return gm.resp, gm.body, gm.err
}

func (gm *reqMocker) Delete(req *httd.Request) (*http.Response, []byte, error) {
	gm.req = req
	return gm.resp, gm.body, gm.err
}

func (gm *reqMocker) Request(req *httd.Request) (*http.Response, []byte, error) {
	gm.req = req
	return gm.resp, gm.body, gm.err
}

var _ httd.Getter = (*reqMocker)(nil)
var _ httd.Poster = (*reqMocker)(nil)
var _ httd.Puter = (*reqMocker)(nil)
var _ httd.Patcher = (*reqMocker)(nil)
var _ httd.Deleter = (*reqMocker)(nil)
var _ httd.Requester = (*reqMocker)(nil)
