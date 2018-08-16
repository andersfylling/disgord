package httd

import (
	"net/http"
)

type Requester interface {
	Request(req *Request) (resp *http.Response, body []byte, err error)
	Get(req *Request) (resp *http.Response, body []byte, err error)
	Post(req *Request) (resp *http.Response, body []byte, err error)
	Put(req *Request) (resp *http.Response, body []byte, err error)
	Patch(req *Request) (resp *http.Response, body []byte, err error)
	Delete(req *Request) (resp *http.Response, body []byte, err error)
}

type Getter interface {
	Get(req *Request) (resp *http.Response, body []byte, err error)
}

type Poster interface {
	Post(req *Request) (resp *http.Response, body []byte, err error)
}

type Puter interface {
	Put(req *Request) (resp *http.Response, body []byte, err error)
}

type Patcher interface {
	Patch(req *Request) (resp *http.Response, body []byte, err error)
}

type Deleter interface {
	Delete(req *Request) (resp *http.Response, body []byte, err error)
}

func (c *Client) Get(req *Request) (resp *http.Response, body []byte, err error) {
	req.Method = http.MethodGet
	return c.Request(req)
}

func (c *Client) Post(req *Request) (resp *http.Response, body []byte, err error) {
	req.Method = http.MethodPost
	return c.Request(req)
}

func (c *Client) Put(req *Request) (resp *http.Response, body []byte, err error) {
	req.Method = http.MethodPut
	return c.Request(req)
}

func (c *Client) Patch(req *Request) (resp *http.Response, body []byte, err error) {
	req.Method = http.MethodPatch
	return c.Request(req)
}

func (c *Client) Delete(req *Request) (resp *http.Response, body []byte, err error) {
	req.Method = http.MethodDelete
	return c.Request(req)
}

/*
type Response struct {
	StatusCode int
	Body       io.ReadCloser
}

func Get(endpoint string, bucket *Bucket) (resp *Response, err error) {

}
func Post()   {}
func Put()    {}
func Patch()  {}
func Delete() {}
*/
