package disgord

import (
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

func (gm *reqMocker) Do(req *httd.Request) (*http.Response, []byte, error) {
	gm.req = req
	return gm.resp, gm.body, gm.err
}

var _ httd.Getter = (*reqMocker)(nil)
var _ httd.Poster = (*reqMocker)(nil)
var _ httd.Puter = (*reqMocker)(nil)
var _ httd.Patcher = (*reqMocker)(nil)
var _ httd.Deleter = (*reqMocker)(nil)
var _ httd.Requester = (*reqMocker)(nil)

func TestParamHolder_URLQueryString(t *testing.T) {
	params := urlQuery{}
	params["a"] = 45
	verifyQueryString(t, params, "?a=45")

	params = urlQuery{}
	verifyQueryString(t, params, "")
}
