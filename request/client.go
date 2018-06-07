package request

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
	BaseURL string = "https://discordapp.com/api"

	// Header
	AuthorizationFormat string = "Bot %s"
	UserAgentFormat     string = "DiscordBot (%s, %s) %s"
)

type DiscordRequester interface {
	Request(method, ratelimiter, endpoint string, target interface{}, jsonParams interface{}) (timeout int64, err error)
	Get(ratelimiter, endpoint string, target interface{}) (timeout int64, err error)
	Post(ratelimiter, endpoint string, target interface{}, jsonParams interface{}) (timeout int64, err error)
	Put(ratelimiter, endpoint string, target interface{}, jsonParams interface{}) (timeout int64, err error)
	Patch(ratelimiter, endpoint string, target interface{}, jsonParams interface{}) (timeout int64, err error)
	Delete(ratelimiter, endpoint string) (timeout int64, err error)
}

type DiscordGetter interface {
	Get(ratelimiter, endpoint string, target interface{}) (timeout int64, err error)
}

type DiscordPoster interface {
	Post(ratelimiter, endpoint string, target interface{}, jsonParams interface{}) (timeout int64, err error)
}

type DiscordPutter interface {
	Put(ratelimiter, endpoint string, target interface{}, jsonParams interface{}) (timeout int64, err error)
}

type DiscordPatcher interface {
	Patch(ratelimiter, endpoint string, target interface{}, jsonParams interface{}) (timeout int64, err error)
}

type DiscordDeleter interface {
	Delete(ratelimiter, endpoint string) (timeout int64, err error)
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

type Client struct {
	url                          string // base url with API version
	rateLimit                    *RateLimit
	reqHeader                    http.Header
	httpClient                   *http.Client
	cancelRequestWhenRateLimited bool
}

func (c *Client) Request(method, ratelimiter, endpoint string, target interface{}, jsonParams interface{}) (timeout int64, err error) {

	var jsonParamsReader io.Reader
	if jsonParams != nil {
		jsonParamsReader, err = convertStructToIOReader(jsonParams)
		if err != nil {
			return
		}
	}

	req, err := http.NewRequest(method, c.url+endpoint, jsonParamsReader)
	if err != nil {
		return
	}

	//check if rate limited
	// discord specifies this in seconds, however it is converted to milliseconds before stored in memory.
	timeout = c.rateLimit.RateLimitTimeout(ratelimiter)
	if timeout > 0 {
		// wait until rate limit is over.
		// exception; if the rate limit timeout exceeds the http client timeout, return error.
		//
		// if cancelRequestWhenRateLimited, is activated
		if c.cancelRequestWhenRateLimited || (c.httpClient.Timeout <= time.Millisecond*time.Duration(timeout)) {
			return timeout, errors.New("rate limited")
		}

		time.Sleep(time.Millisecond * time.Duration(timeout))
	}

	req.Header = c.reqHeader
	res, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	// update rate limits
	c.rateLimit.HandleResponse(ratelimiter, res)

	// successful deletes return 204. TODO: confirm
	if method == http.MethodDelete {
		if res.Status == http.MethodDelete {
			err = nil
			return
		}

		err = errors.New("Unable to delete resource at " + endpoint)
		return
	}

	// if a target has been provided for un-marshalling
	err = nil
	if target != nil {
		err = json.NewDecoder(res.Body).Decode(target)
	}

	return
}

func (c *Client) Get(ratelimiter, endpoint string, target interface{}) (timeout int64, err error) {
	return c.Request(http.MethodGet, ratelimiter, endpoint, target, nil)
}

func (c *Client) Post(ratelimiter, endpoint string, target interface{}, jsonParams interface{}) (timeout int64, err error) {
	return c.Request(http.MethodPost, ratelimiter, endpoint, target, jsonParams)
}

func (c *Client) Put(ratelimiter, endpoint string, target interface{}, jsonParams interface{}) (timeout int64, err error) {
	return c.Request(http.MethodPut, ratelimiter, endpoint, target, jsonParams)
}

func (c *Client) Patch(ratelimiter, endpoint string, target interface{}, jsonParams interface{}) (timeout int64, err error) {
	return c.Request(http.MethodPatch, ratelimiter, endpoint, target, jsonParams)
}

func (c *Client) Delete(ratelimiter, endpoint string) (timeout int64, err error) {
	return c.Request(http.MethodDelete, ratelimiter, endpoint, nil, nil)
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
