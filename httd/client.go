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

// defaults and string format's for Discord interaction
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

// Requester holds all the sub-request interface for Discord interaction
type Requester interface {
	Request(req *Request) (resp *http.Response, body []byte, err error)
	Getter
	Poster
	Puter
	Patcher
	Deleter
}

// Getter interface which holds the Get method for sending get requests to Discord
type Getter interface {
	Get(req *Request) (resp *http.Response, body []byte, err error)
}

// Poster interface which holds the Post method for sending post requests to Discord
type Poster interface {
	Post(req *Request) (resp *http.Response, body []byte, err error)
}

// Puter interface which holds the Put method for sending put requests to Discord
type Puter interface {
	Put(req *Request) (resp *http.Response, body []byte, err error)
}

// Patcher interface which holds the Patch method for sending patch requests to Discord
type Patcher interface {
	Patch(req *Request) (resp *http.Response, body []byte, err error)
}

// Deleter interface which holds the Delete method for sending delete requests to Discord
type Deleter interface {
	Delete(req *Request) (resp *http.Response, body []byte, err error)
}

type ErrREST struct {
	Code       int    `json:"code"`
	Msg        string `json:"message"`
	Suggestion string `json:"-"`
	HTTPCode   int    `json:"-"`
}

func (e *ErrREST) Error() string {
	return e.Msg
}

// Client is the httd client for handling Discord requests
type Client struct {
	url                          string // base url with API version
	rateLimit                    *RateLimit
	reqHeader                    http.Header
	httpClient                   *http.Client // TODO: decouple to allow better unit testing of REST requests
	cancelRequestWhenRateLimited bool
}

// Get handles Discord get requests
func (c *Client) Get(req *Request) (resp *http.Response, body []byte, err error) {
	req.Method = http.MethodGet
	return c.Request(req)
}

// Post handles Discord post requests
func (c *Client) Post(req *Request) (resp *http.Response, body []byte, err error) {
	req.Method = http.MethodPost
	return c.Request(req)
}

// Put handles Discord put requests
func (c *Client) Put(req *Request) (resp *http.Response, body []byte, err error) {
	req.Method = http.MethodPut
	return c.Request(req)
}

// Patch handles Discord patch requests
func (c *Client) Patch(req *Request) (resp *http.Response, body []byte, err error) {
	req.Method = http.MethodPatch
	return c.Request(req)
}

// Delete handles Discord delete requests
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

// NewClient ...
func NewClient(conf *Config) (*Client, error) {
	if !SupportsDiscordAPIVersion(conf.APIVersion) {
		return nil, errors.New(fmt.Sprintf("Discord API version %d is not supported", conf.APIVersion))
	}

	if conf.BotToken == "" {
		return nil, errors.New("no Discord Bot Token was provided")
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
		return nil, errors.New("both a source(url) and a version must be present for sending requests to the Discord REST API")
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
	}, nil
}

// Config is the configuration options for the httd.Client structure. Essentially the behaviour of all requests
// sent to Discord.
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

// Details ...
type Details struct {
	Ratelimiter     string
	Endpoint        string // always as a suffix to Ratelimiter(!)
	ResponseStruct  interface{}
	SuccessHTTPCode int
}

// RateLimitAdjuster acts as a middleware when receiving ratelimits from Discord.
type RateLimitAdjuster func(timeout time.Duration) time.Duration

// Request is populated before executing a Discord request to correctly generate a http request
type Request struct {
	Method            string
	Ratelimiter       string
	RateLimitAdjuster RateLimitAdjuster
	Endpoint          string
	Body              interface{} // will automatically marshal to JSON if the ContentType is httd.ContentTypeJSON
	ContentType       string
}

func (c *Client) decodeResponseBody(resp *http.Response) (body []byte, err error) {
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.Header.Get(ContentEncoding) {
	case GZIPCompression:
		b := bytes.NewBuffer(buffer)

		var r io.Reader
		r, err = gzip.NewReader(b)
		if err != nil {
			return nil, err
		}

		var resB bytes.Buffer
		_, err = resB.ReadFrom(r)
		if err != nil {
			return nil, err
		}

		body = resB.Bytes()
	default:
		body = buffer
	}

	return body, nil
}

func makeRateLimitCompliant(key string, c *Client) (err error) {
	var timeout time.Duration
	remaining := c.httpClient.Timeout.Nanoseconds()
	for {
		timeout, err = c.rateLimit.RequestPermit(key)
		if err != nil {
			return err // no way around the rate limiter in this case
		}
		tns := timeout.Nanoseconds()

		if tns > 0 && c.cancelRequestWhenRateLimited {
			return errors.New("rate limited")
		}

		if tns > 0 {
			if remaining < tns {
				return errors.New("rate limit timeout is higher than http.Client.Timeout, cannot wait")
			}
			remaining -= tns

			// if we are given a timeout, we weren't allowed to send a request so we need to retry
			<-time.After(timeout)
			continue
		}

		return nil
	}
}

// Request execute a Discord request
func (c *Client) Request(r *Request) (resp *http.Response, body []byte, err error) {
	err = makeRateLimitCompliant(r.Ratelimiter, c)
	if err != nil {
		return nil, nil, err
	}

	// prepare request body
	var bodyReader io.Reader
	if r.Body != nil {
		switch b := r.Body.(type) { // Determine the type of the passed body so we can treat it differently
		case io.Reader:
			bodyReader = b
		default:
			// If the type is unknown, possibly Marshal it as JSON
			if r.ContentType != ContentTypeJSON {
				return nil, nil, errors.New("unknown request body types and only be used in conjunction with httd.ContentTypeJSON")
			}

			bodyReader, err = convertStructToIOReader(r.Body)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	// create request
	req, err := http.NewRequest(r.Method, c.url+r.Endpoint, bodyReader)
	if err != nil {
		return nil, nil, err
	}
	req.Header = c.reqHeader
	req.Header.Set(ContentType, r.ContentType) // unique for each request

	// send request
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err = c.decodeResponseBody(resp)

	// update rate limits
	c.rateLimit.UpdateRegisters(r.Ratelimiter, r.RateLimitAdjuster, resp, body)

	// check if request was successful
	noDiff := resp.StatusCode == http.StatusNotModified
	withinSuccessScope := 200 <= resp.StatusCode && resp.StatusCode < 300
	if !(noDiff || withinSuccessScope) {
		// not within successful http range
		msg := "response was not within the successful http code range [200, 300). code: "
		msg += strconv.Itoa(resp.StatusCode)

		err = &ErrREST{
			Suggestion: msg,
			HTTPCode:   resp.StatusCode,
		}

		// store the Discord error if it exists
		if len(body) > 0 {
			_ = Unmarshal(body, err)
		}
		return nil, nil, err
	}

	return resp, body, nil
}

// RateLimiter get the rate limit manager
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
