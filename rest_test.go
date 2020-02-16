// +build !integration

package disgord

import (
	"context"
	"net/http"
	"testing"

	"github.com/andersfylling/disgord/internal/httd"
)

func verifyQueryString(t *testing.T, params URLQueryStringer, wants string) {
	got := params.URLQueryString()
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

func (gm *reqMocker) Do(ctx context.Context, req *httd.Request) (*http.Response, []byte, error) {
	gm.req = req
	gm.req.Ctx = ctx
	return gm.resp, gm.body, gm.err
}

var _ httd.Requester = (*reqMocker)(nil)

func TestParamHolder_URLQueryString(t *testing.T) {
	params := urlQuery{}
	params["a"] = 45
	verifyQueryString(t, params, "?a=45")

	params = urlQuery{}
	verifyQueryString(t, params, "")
}
