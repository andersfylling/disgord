package httd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
)

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
		"Authorization": {authorization},
		"User-Agent":    {userAgent},
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
}

type Client struct {
	url                          string // base url with API version
	rateLimit                    *RateLimit
	reqHeader                    http.Header
	httpClient                   *http.Client
	cancelRequestWhenRateLimited bool
}

func (c *Client) Request(r *Request) (resp *http.Response, err error) {
	var jsonParamsReader io.Reader
	if r.JSONParams != nil {
		jsonParamsReader, err = convertStructToIOReader(r.JSONParams)
		if err != nil {
			return
		}
	}

	req, err := http.NewRequest(r.Method, c.url+r.Endpoint, jsonParamsReader)
	if err != nil {
		return
	}

	// check if rate limited
	timeout := int64(0)

	// check global rate limit
	if c.rateLimit.global.remaining == 0 {
		timeout = c.rateLimit.global.timeout()
	}

	// route specific rate limit
	if r.Ratelimiter != "" && timeout == 0 {
		timeout = c.rateLimit.RateLimitTimeout(r.Ratelimiter)
	}

	// discord specifies this in seconds, however it is converted to milliseconds before stored in memory.
	if timeout > 0 {
		// wait until rate limit is over.
		// exception; if the rate limit timeout exceeds the http client timeout, return error.
		//
		// if cancelRequestWhenRateLimited, is activated
		deadtime := time.Millisecond*time.Duration(timeout)
		if c.cancelRequestWhenRateLimited || (c.httpClient.Timeout <= deadtime) {
			err = errors.New("rate limited")
			// TODO: add the timeout to the return
			return
		}

		time.Sleep(deadtime)
	}

	req.Header = c.reqHeader
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return
	}

	// update rate limits
	c.rateLimit.HandleResponse(r.Ratelimiter, resp)

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
