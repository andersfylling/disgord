package disgord

import (
	"net/http"
	"testing"

	"github.com/andersfylling/disgord/httd"
)

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
