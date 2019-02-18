package disgord

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/snowflake/v3"
)

type ErrRest = httd.ErrREST

// URLQueryStringer converts a struct of values to a valid URL query string
type URLQueryStringer interface {
	URLQueryString() string
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

func newRESTBuilder(cache *Cache, client httd.Requester, config *httd.Request, middleware fRESTRequestMiddleware) *RESTBuilder {
	builder := &RESTBuilder{}
	builder.setup(cache, client, config, middleware)

	return builder
}

type paramHolder map[string]interface{}

var _ URLQueryStringer = (*paramHolder)(nil)

func (p paramHolder) URLQueryString() string {
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

type astrixREST func(repo cacheRegistry, ids []Snowflake, req *httd.Request, flags ...Flag) (v interface{}, err error)

type restStepCheckCache func() (v interface{}, err error)
type restStepDoRequest func() (resp *http.Response, body []byte, err error)
type restStepUpdateCache func(registry cacheRegistry, id Snowflake, x interface{}) (err error)
type rest struct {
	c     *client
	flags Flag // merge flags

	// caching
	ID            Snowflake
	CacheRegistry cacheRegistry
	httpMethod    string

	// item creation
	// pool is prioritzed over factory
	pool    Pool
	factory func() Reseter

	conf *httd.Request

	// steps
	checkCache  restStepCheckCache
	doRequest   restStepDoRequest
	updateCache restStepUpdateCache
}

func (r *rest) Put(x Reseter) {
	if r.pool != nil {
		r.pool.Put(x)
	} else {
		x = nil // GC
	}
}

func (r *rest) Get() (x Reseter) {
	if r.pool != nil {
		return r.pool.Get()
	}

	return r.factory()
}

var _ Pool = (*rest)(nil)

func (r *rest) init() {
	if r.conf != nil {
		if r.conf.Method == "" {
			r.conf.Method = http.MethodGet
		}
		if r.conf.Method != "" {
			r.httpMethod = r.conf.Method
		}
	}
	if r.httpMethod == "" {
		r.httpMethod = http.MethodGet
	}

	r.checkCache = r.stepCheckCache
	r.doRequest = r.stepDoRequest
	r.updateCache = r.stepUpdateCache
}

func (r *rest) bindParams(params interface{}) {
	if params == nil {
		return
	}
}

func (r *rest) stepCheckCache() (v interface{}, err error) {
	if r.httpMethod != http.MethodGet {
		return nil, nil
	}

	if r.CacheRegistry == NoCacheSpecified {
		return nil, nil
	}

	if r.ID.Empty() {
		return nil, nil
	}

	if r.flags.Ignorecache() {
		return nil, nil
	}

	return r.c.cache.Get(r.CacheRegistry, r.ID)
}

func (r *rest) stepDoRequest() (resp *http.Response, body []byte, err error) {
	if r.conf == nil {
		err = errors.New("missing httd.Request configuration")
		return
	}

	resp, body, err = r.c.req.Request(r.conf)
	return
}

// stepUpdateCache id is only used when deleting an object
func (r *rest) stepUpdateCache(registry cacheRegistry, id Snowflake, x interface{}) (err error) {
	if r.CacheRegistry == NoCacheSpecified {
		return nil
	}

	if x == nil {
		return nil
	}

	return r.c.cache.Update(registry, x)
}

func (r *rest) processContent(body []byte) (v interface{}, err error) {
	if len(body) == 0 {
		return nil, nil // nothing more to do
	}

	// next we create the object and save it to cache
	if r.pool == nil && r.factory == nil {
		return nil, errors.New("missing factory method for either pool or factory")
	}

	obj := r.Get()
	if err = httd.Unmarshal(body, obj); err != nil {
		r.Put(obj)
		return nil, err
	}

	// update the object
	if updater, implements := obj.(internalUpdater); implements {
		updater.updateInternals()
	}
	if updater, implements := obj.(internalClientUpdater); implements {
		updater.updateInternalsWithClient(r.c)
	}

	return obj, nil
}

func (r *rest) Execute() (v interface{}, err error) {
	if v, err = r.checkCache(); err == nil && v != nil {
		return v, err
	}

	var resp *http.Response
	var body []byte
	if resp, body, err = r.doRequest(); err != nil {
		return nil, err
	}

	if resp.Request.Method == http.MethodDelete && resp.StatusCode != http.StatusNoContent {
		msg := "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(http.StatusNoContent)
		err = errors.New(msg)
		return nil, err
	}

	var obj interface{}
	if obj, err = r.processContent(body); err != nil {
		return nil, err
	}

	// save it to cache / update the cache
	if err = r.updateCache(r.CacheRegistry, r.ID, obj); err != nil {
		r.c.log.Error(err)
	}

	return obj, nil
}

type fRESTRequestMiddleware func(resp *http.Response, body []byte, err error) error
type fRESTCacheMiddleware func(resp *http.Response, v interface{}, err error) error
type fRESTItemFactory func() interface{}

//go:generate go run generate/restbuilders/main.go

type RESTBuilder struct {
	middleware fRESTRequestMiddleware
	config     *httd.Request
	client     httd.Requester

	itemFactory fRESTItemFactory

	cache           *Cache
	cacheRegistry   cacheRegistry
	cacheMiddleware fRESTCacheMiddleware
	cacheItemID     snowflake.ID

	body              map[string]interface{}
	urlParams         paramHolder
	ignoreCache       bool
	cancelOnRatelimit bool
}

func (b *RESTBuilder) setup(cache *Cache, client httd.Requester, config *httd.Request, middleware fRESTRequestMiddleware) {
	b.body = make(map[string]interface{})
	b.urlParams = make(map[string]interface{})
	b.cache = cache
	b.client = client
	b.config = config
	b.middleware = middleware
}

func (b *RESTBuilder) cacheLink(registry cacheRegistry, middleware fRESTCacheMiddleware) {
	b.cacheRegistry = registry
	b.cacheMiddleware = middleware
}

func (b *RESTBuilder) prepare() {
	// update the config
	if b.config.ContentType != "" {
		b.config.Body = b.body
	}
	b.config.Endpoint += b.urlParams.URLQueryString()
}

// execute ... v must be a nil pointer.
func (b *RESTBuilder) execute() (v interface{}, err error) {
	if !b.ignoreCache && b.config.Method == http.MethodGet && !b.cacheItemID.Empty() {
		// cacheLink lookup. return on cacheLink hit
		v, err = b.cache.Get(b.cacheRegistry, b.cacheItemID)
		if v != nil && err == nil {
			return v, nil
		}
		// otherwise we perform the request
	}

	b.prepare()

	var resp *http.Response
	var body []byte
	resp, body, err = b.client.Request(b.config)
	if err != nil {
		return nil, err
	}

	if b.middleware != nil {
		if err = b.middleware(resp, body, err); err != nil {
			return nil, err
		}
	}

	if len(body) > 1 {
		v = b.itemFactory()
		if err = httd.Unmarshal(body, v); err != nil {
			return nil, err
		}

		if b.cacheRegistry == NoCacheSpecified {
			return v, err
		}

		if b.cacheMiddleware != nil {
			b.cacheMiddleware(resp, v, err)
		}

		b.cache.Update(b.cacheRegistry, v)
	}
	return v, nil
}

type restReqBuilderAsync struct {
	Data interface{}
	Err  error
	// FromCache bool // TODO
}

func (b *RESTBuilder) async() <-chan *restReqBuilderAsync {
	A := make(chan *restReqBuilderAsync)
	go func() {
		resp := &restReqBuilderAsync{}
		resp.Data, resp.Err = b.execute()

		A <- resp
		close(A)
	}()

	return A
}

func (b *RESTBuilder) param(name string, v interface{}) *RESTBuilder {
	if b.config.Method == http.MethodGet {
		// RFC says you can not send a body in a GET request
		b.queryParam(name, v)
	} else {
		b.body[name] = v
	}
	return b
}

func (b *RESTBuilder) queryParam(name string, v interface{}) *RESTBuilder {
	b.urlParams[name] = v
	return b
}

func (b *RESTBuilder) IgnoreCache() *RESTBuilder {
	b.ignoreCache = true
	return b
}

func (b *RESTBuilder) CancelOnRatelimit() *RESTBuilder {
	b.cancelOnRatelimit = true
	return b
}

// GetGateway [REST] Returns an object with a single valid WSS URL, which the client can use for Connecting.
// Clients should cacheLink this value and only call this endpoint to retrieve a new URL if they are unable to
// properly establish a connection using the cached version of the URL.
//  Method                  GET
//  Endpoint                /gateway
//  Rate limiter            /gateway
//  Discord documentation   https://discordapp.com/developers/docs/topics/gateway#get-gateway
//  Reviewed                2018-10-12
//  Comment                 This endpoint does not require authentication.
func GetGateway(client httd.Getter) (gateway *Gateway, err error) {
	var body []byte
	_, body, err = client.Get(&httd.Request{
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
	var body []byte
	_, body, err = client.Get(&httd.Request{
		Ratelimiter: "/gateway/bot",
		Endpoint:    "/gateway/bot",
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &gateway)
	return
}

func getChannel(f func() (interface{}, error)) (channel *Channel, err error) {
	var v interface{}
	if v, err = f(); err != nil {
		return nil, err
	}

	if v == nil {
		return nil, errors.New("object was nil")
	}

	return v.(*Channel), nil
}

func getUser(f func() (interface{}, error)) (user *User, err error) {
	var v interface{}
	if v, err = f(); err != nil {
		return nil, err
	}

	if v == nil {
		return nil, errors.New("object was nil")
	}

	return v.(*User), nil
}
