package httd

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	BaseURL = "https://discordapp.com/api"

	// Header
	AuthorizationFormat = "Bot %s"
	UserAgentFormat     = "DiscordBot (%s, %s) %s"

	HTTPCodeRateLimit int = 429

	ContentEncoding = "Content-Encoding"
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
	GZIPCompression = "gzip"
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

// SupportsDiscordAPIVersion check if a given discord api version is supported by this package.
func SupportsDiscordAPIVersion(version int) bool {
	supports := []int{
		6,
	}

	var supported bool
	for _, supportedVersion := range supports {
		if supportedVersion == version {
			supported = true
			break
		}
	}

	return supported
}

func NewClient(conf *Config) *Client {
	if !SupportsDiscordAPIVersion(conf.APIVersion) {
		panic(fmt.Sprintf("Discord API version %d is not supported", conf.APIVersion))
	}

	if conf.BotToken == "" {
		panic("No Discord Bot Token was provided")
	}

	// if no http client was provided, create a new one
	if conf.HTTPClient == nil {
		conf.HTTPClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	// Clients using the HTTP API must provide a valid User Agent which specifies
	// information about the client library and version in the following format:
	//	User-Agent: DiscordBot ($url, $versionNumber)
	if conf.UserAgentSourceURL == "" || conf.UserAgentVersion == "" {
		panic("Both a source(url) and a version must be present for sending requests to the Discord REST API")
	}

	// setup the required http request header fields
	authorization := fmt.Sprintf(AuthorizationFormat, conf.BotToken)
	userAgent := fmt.Sprintf(UserAgentFormat, conf.UserAgentSourceURL, conf.UserAgentVersion, conf.UserAgentExtra)
	header := map[string][]string{
		"Authorization":   {authorization},
		"User-Agent":      {userAgent},
		"Accept-Encoding": {"gzip"},
	}

	return &Client{
		url:        BaseURL + "/v" + strconv.Itoa(conf.APIVersion),
		reqHeader:  header,
		httpClient: conf.HTTPClient,
		rateLimit:  NewRateLimit(),
	}
}

type Config struct {
	APIVersion int
	BotToken   string

	HTTPClient *http.Client

	CancelRequestWhenRateLimited bool

	// Header field: `User-Agent: DiscordBot ({Source}, {Version}) {Extra}`
	UserAgentVersion   string
	UserAgentSourceURL string
	UserAgentExtra     string
}

type Details struct {
	Ratelimiter     string
	Endpoint        string // always as a suffix to Ratelimiter(!)
	ResponseStruct  interface{}
	SuccessHttpCode int
}

type Request struct {
	Method      string
	Ratelimiter string
	Endpoint    string
	JSONParams  interface{}
	ContentType string
}

type Client struct {
	url                          string // base url with API version
	rateLimit                    *RateLimit
	reqHeader                    http.Header
	httpClient                   *http.Client
	cancelRequestWhenRateLimited bool
}

func (c *Client) decodeResponseBody(resp *http.Response) (body []byte, err error) {
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	switch resp.Header.Get(ContentEncoding) {
	case GZIPCompression:
		b := bytes.NewBuffer(buffer)

		var r io.Reader
		r, err = gzip.NewReader(b)
		if err != nil {
			return
		}

		var resB bytes.Buffer
		_, err = resB.ReadFrom(r)
		if err != nil {
			return
		}

		body = resB.Bytes()
	default:
		body = buffer
	}

	return
}

func WaitIfRateLimited(c *Client, r *Request) (waited bool, err error) {
	deadtime := c.RateLimiter().WaitTime(r)
	if deadtime.Nanoseconds() > 0 {
		if c.cancelRequestWhenRateLimited {
			err = errors.New("rate limited")
			return
		}
		waitTimeLongerThanHTTPTimeout := (c.httpClient.Timeout.Nanoseconds() <= deadtime.Nanoseconds())
		if waitTimeLongerThanHTTPTimeout {
			err = errors.New("rate limit timeout is higher than http.Client.Timeout, cannot wait")
			return
		}

		<-time.After(deadtime)
	}

	waited = true
	return
}

func (c *Client) Request(r *Request) (resp *http.Response, body []byte, err error) {
	var jsonParamsReader io.Reader
	if r.JSONParams != nil {
		jsonParamsReader, err = convertStructToIOReader(r.JSONParams)
		if err != nil {
			return
		}
	}

	// check the rate limiter for how long we must wait before sending the request
	_, err = WaitIfRateLimited(c, r)
	if err != nil {
		return
	}

	// create request
	req, err := http.NewRequest(r.Method, c.url+r.Endpoint, jsonParamsReader)
	if err != nil {
		return
	}
	req.Header = c.reqHeader
	req.Header.Set(ContentType, r.ContentType) // unique for each request

	// send request
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = c.decodeResponseBody(resp)

	// update rate limits
	c.RateLimiter().UpdateRegisters(r.Ratelimiter, resp, body)

	// check if request was successful
	noDiff := resp.StatusCode == http.StatusNotModified
	withinSuccessScope := 200 <= resp.StatusCode && resp.StatusCode < 300
	if !(noDiff || withinSuccessScope) {
		// not within successful http range
		// TODO: redirects?
		msg := "response was not within the successful http code range [200, 300). code: "
		msg += strconv.Itoa(resp.StatusCode) + ", response: "
		msg += string(body)
		err = errors.New(msg)
	}

	return
}

func (c *Client) RateLimiter() RateLimiter {
	return c.rateLimit
}

// helper functions
func convertStructToIOReader(v interface{}) (io.Reader, error) {
	jsonParamsBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	jsonParamsReader := bytes.NewReader(jsonParamsBytes)

	return jsonParamsReader, nil
}
