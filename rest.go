package disgord

import (
	"encoding/json"
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/snowflake/v3"
	"net/http"
	"strconv"
)

type ErrRest = httd.ErrREST

// URLParameters converts a struct of values to a valid URL query string
type URLParameters interface {
	GetQueryString() string
}

func unmarshal(data []byte, v interface{}) error {
	return httd.Unmarshal(data, v)
}

func marshal(v interface{}) ([]byte, error) {
	return httd.Marshal(v)
}

// AvatarParamHolder is used when handling avatar related REST structs.
// since a Avatar can be reset by using nil, it causes some extra issues as omit empty cannot be used
// to get around this, the struct requires an internal state and must also handle custom marshalling
type AvatarParamHolder interface {
	json.Marshaler
	Empty() bool
	SetAvatar(avatar string)
	UseDefaultAvatar()
}

func newRESTBuilder(client httd.Requester, config *httd.Request, middleware RESTRequestMiddleware) *RESTRequestBuilder {
	builder := &RESTRequestBuilder{}
	builder.setup(client, config, middleware)

	return builder
}

type paramHolder map[string]interface{}

func (p paramHolder) GetQueryString() string {
	if len(p) == 0 {
		return ""
	}

	var params string
	seperator := "?"
	for k, v := range p {
		var str string
		var uintHolder uint64
		switch t := v.(type) {
		case snowflake.ID:
			str = t.String()
		case int8:
			uintHolder = uint64(t)
		case int16:
			uintHolder = uint64(t)
		case int32:
			uintHolder = uint64(t)
		case int64:
			uintHolder = uint64(t)
		case int:
			uintHolder = uint64(t)
		case uint8:
			uintHolder = uint64(t)
		case uint16:
			uintHolder = uint64(t)
		case uint32:
			uintHolder = uint64(t)
		case uint:
			uintHolder = uint64(t)
		case uint64:
			uintHolder = t
		case bool:
			if t {
				str = "true"
			} else {
				str = "false"
			}
		case string:
			str = t
		}
		// TODO float
		if str == "" {
			str = strconv.FormatUint(uintHolder, 10)
		}

		params += seperator + k + "=" + str
		seperator = "&"
	}

	return params
}

var _ URLParameters = (*paramHolder)(nil)

type RESTRequestMiddleware func(resp *http.Response, body []byte, err error) error

type RESTRequestBuilder struct {
	middleware RESTRequestMiddleware
	config     *httd.Request
	client     httd.Requester

	cacheID           uint
	cacheActionUpdate bool

	body              map[string]interface{}
	urlParams         paramHolder
	ignoreCache       bool
	cancelOnRatelimit bool
}

func (b *RESTRequestBuilder) setup(client httd.Requester, config *httd.Request, middleware RESTRequestMiddleware) {
	b.body = make(map[string]interface{})
	b.urlParams = make(map[string]interface{})
	b.client = client
	b.config = config
	b.middleware = middleware
}

func (b *RESTRequestBuilder) cache(id uint, update bool) {
	b.cacheID = id
	b.cacheActionUpdate = update
}

// execute ... v must be a nil pointer.
func (b *RESTRequestBuilder) execute(v interface{}) (err error) {
	if !b.ignoreCache && b.config.Method == http.MethodGet {
		// cache lookup. return on cache hit

	}

	// update the config
	if b.config.ContentType != "" {
		b.config.Body = b.body
	}
	b.config.Endpoint += b.urlParams.GetQueryString()

	var resp *http.Response
	var body []byte
	resp, body, err = b.client.Request(b.config)
	if err != nil {
		return
	}

	if b.middleware != nil {
		err = b.middleware(resp, body, err)
		if err != nil {
			return
		}
	}

	if !b.ignoreCache {

	}

	if len(body) > 1 {
		err = httd.Unmarshal(body, v)
	}
	return
}

func (b *RESTRequestBuilder) Param(name string, v interface{}) *RESTRequestBuilder {
	b.body[name] = v
	return b
}

func (b *RESTRequestBuilder) IgnoreCache() *RESTRequestBuilder {
	b.ignoreCache = true
	return b
}

func (b *RESTRequestBuilder) CancelOnRatelimit() *RESTRequestBuilder {
	b.cancelOnRatelimit = true
	return b
}

// GetGateway [REST] Returns an object with a single valid WSS URL, which the client can use for Connecting.
// Clients should cache this value and only call this endpoint to retrieve a new URL if they are unable to
// properly establish a connection using the cached version of the URL.
//  Method                  GET
//  Endpoint                /gateway
//  Rate limiter            /gateway
//  Discord documentation   https://discordapp.com/developers/docs/topics/gateway#get-gateway
//  Reviewed                2018-10-12
//  Comment                 This endpoint does not require authentication.
func GetGateway(client httd.Getter) (gateway *Gateway, err error) {
	_, body, err := client.Get(&httd.Request{
		Ratelimiter: "/gateway",
		Endpoint:    "/gateway",
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &gateway)
	return
}

// GetGatewayBot [REST] Returns an object based on the information in Get Gateway, plus additional metadata
// that can help during the operation of large or sharded bots. Unlike the Get Gateway, this route should not
// be cached for extended periods of time as the value is not guaranteed to be the same per-call, and
// changes as the bot joins/leaves guilds.
//  Method                  GET
//  Endpoint                /gateway/bot
//  Rate limiter            /gateway/bot
//  Discord documentation   https://discordapp.com/developers/docs/topics/gateway#get-gateway-bot
//  Reviewed                2018-10-12
//  Comment                 This endpoint requires authentication using a valid bot token.
func GetGatewayBot(client httd.Getter) (gateway *GatewayBot, err error) {
	_, body, err := client.Get(&httd.Request{
		Ratelimiter: "/gateway/bot",
		Endpoint:    "/gateway/bot",
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &gateway)
	return
}
