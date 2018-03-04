package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	BaseURL string = "https://discordapp.com/api"

	// Header
	AuthorizationFormat string = "Bot %s"
	UserAgentFormat     string = "DiscordBot (%s, %s) %s"
)

type Requester interface {
	Request(method, uri string, target interface{}, jsonParams io.Reader) error
	Get(uri string, target interface{}) error
	Post(uri string, target interface{}, jsonParams interface{}) error
	Put(uri string, target interface{}, jsonParams interface{}) error
	Patch(uri string, target interface{}, jsonParams interface{}) error
	Delete(uri string) error
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
		conf.HTTPClient = &http.Client{}
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

	// WaitIfRateLimited if discord has rate limited the lib,
	// this will add any request to a que and complete them ASAP.
	WaitIfRateLimited bool

	// Header field: `User-Agent: DiscordBot ({Source}, {Version}) {Extra}`
	UserAgentVersion   string
	UserAgentSourceURL string
	UserAgentExtra     string
}

type Client struct {
	url        string // base url with API version
	rateLimit  *RateLimit
	reqHeader  http.Header
	httpClient *http.Client
}

func (c *Client) Request(method, uri string, target interface{}, jsonParams io.Reader) error {
	req, err := http.NewRequest(method, c.url+uri, jsonParams)
	if err != nil {
		return err
	}

	req.Header = c.reqHeader
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// successful deletes return 204
	if method == http.MethodDelete {
		if res.Status == http.MethodDelete {
			return nil
		}

		return errors.New("Unable to delete resource at " + uri)
	}

	// if a target has been provided for un-marshalling
	if target != nil {

		return json.NewDecoder(res.Body).Decode(target)
	}

	return nil
}

func (c *Client) Get(uri string, target interface{}) error {
	return c.Request(http.MethodGet, uri, target, nil)
}

func (c *Client) Post(uri string, target interface{}, jsonParams interface{}) error {
	jsonParamsBytes, err := json.Marshal(jsonParams)
	if err != nil {
		return err
	}
	jsonParamsReader := bytes.NewReader(jsonParamsBytes)

	return c.Request(http.MethodPost, uri, target, jsonParamsReader)
}

func (c *Client) Put(uri string, target interface{}, jsonParams interface{}) error {
	jsonParamsReader, err := convertStructToIOReader(jsonParams)
	if err != nil {
		return err
	}

	return c.Request(http.MethodPut, uri, target, jsonParamsReader)
}

func (c *Client) Patch(uri string, target interface{}, jsonParams interface{}) error {
	jsonParamsReader, err := convertStructToIOReader(jsonParams)
	if err != nil {
		return err
	}

	return c.Request(http.MethodPatch, uri, target, jsonParamsReader)
}

func (c *Client) Delete(uri string) error {
	return c.Request(http.MethodDelete, uri, nil, nil)
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
