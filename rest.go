package disgord

import (
	"encoding/json"
	"github.com/andersfylling/disgord/httd"
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

//
//func newRESTBuilder(client httd.Requester, config *httd.Request, middleware RESTRequestMiddleware) *RESTRequestBuilder {
//	builder := &RESTRequestBuilder{}
//	builder.setup(client, config, middleware)
//
//	return builder
//}
//
//type RESTRequestMiddleware func(resp *http.Response, body []byte, err error) error
//
//type RESTRequestBuilder struct {
//	middleware RESTRequestMiddleware
//	config     *httd.Request
//	client     httd.Requester
//
//	cacheID uint
//	cacheActionUpdate bool
//
//	params            map[string]interface{}
//	ignoreCache       bool
//	cancelOnRatelimit bool
//}
//
//func (b *RESTRequestBuilder) setup(client httd.Requester, config *httd.Request, middleware RESTRequestMiddleware) {
//	b.client = client
//	b.config = config
//	b.middleware = middleware
//}
//
//func (b *RESTRequestBuilder) cache(id uint, update bool) {
//	b.cacheID = id
//	b.cacheActionUpdate = update
//}
//
//func (b *RESTRequestBuilder) Param(name string, v interface{}) *RESTRequestBuilder {
//	b.params[name] = v
//	return b
//}
//
//func (b *RESTRequestBuilder) IgnoreCache() *RESTRequestBuilder {
//	b.ignoreCache = true
//	return b
//}
//
//func (b *RESTRequestBuilder) CancelOnRatelimit() *RESTRequestBuilder {
//	b.cancelOnRatelimit = true
//	return b
//}
//
//func (b *RESTRequestBuilder) ExecuteSimple(v interface{}) (err error) {
//	if !b.ignoreCache && b.config.Method == http.MethodGet {
//
//	}
//
//	var resp *http.Response
//	var body []byte
//	resp, body, err = b.client.Request(b.config)
//	if err != nil {
//		return
//	}
//
//	if b.middleware != nil {
//		err = b.middleware(resp, body, err)
//		if err != nil {
//			return
//		}
//	}
//
//	if !b.ignoreCache {
//
//	}
//
//	err = httd.Unmarshal(body, v)
//	return
//}
//
//func (c *Client) GettUser(id Snowflake) *RESTRequestBuilder {
//	return newRESTBuilder(c.req, &httd.Request{
//		Method:      http.MethodGet,
//		Ratelimiter: ratelimitUsers(),
//		Endpoint:    endpoint.User(id),
//	}, nil)
//}

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
