package disgord

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/websocket"
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

type urlQuery map[string]interface{}

var _ URLQueryStringer = (*urlQuery)(nil)
var _ fmt.Stringer = (*urlQuery)(nil)

func (p urlQuery) String() string {
	return p.URLQueryString()
}

func (p urlQuery) URLQueryString() string {
	if len(p) == 0 {
		return ""
	}

	query := url.Values{}
	for k, v := range p {
		query.Add(k, fmt.Sprint(v))
	}

	if len(query) > 0 {
		return "?" + query.Encode()
	}
	return ""
}

type restStepCheckCache func() (v interface{}, err error)
type restStepDoRequest func() (resp *http.Response, body []byte, err error)
type restStepUpdateCache func(registry cacheRegistry, id Snowflake, x interface{}) (err error)
type rest struct {
	c     *Client
	flags Flag // merge flags

	// caching
	ID            Snowflake
	CacheRegistry cacheRegistry
	httpMethod    string

	// item creation
	// pool is prioritized over factory
	pool    Pool
	factory func() interface{}

	conf              *httd.Request
	expectsStatusCode int

	// steps
	checkCache     restStepCheckCache
	doRequest      restStepDoRequest
	preUpdateCache func(x interface{})
	updateContent  func(x interface{})
	updateCache    restStepUpdateCache
}

func (r *rest) Put(x interface{}) {
	if r.pool != nil {
		r.pool.Put(x.(Reseter))
	} else {
		x = nil // GC
	}
}

func (r *rest) Get() (x interface{}) {
	if r.pool != nil {
		return r.pool.Get()
	}

	return r.factory()
}

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
	r.updateContent = r.stepUpdateContent
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

	if r.preUpdateCache != nil {
		r.preUpdateCache(x)
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
	r.updateContent(obj)

	return obj, nil
}

func (r *rest) stepUpdateContent(x interface{}) {
	if x == nil {
		return
	}
	executeInternalUpdater(x)
	executeInternalClientUpdater(r.c, x)
}

func (r *rest) Execute() (v interface{}, err error) {
	if !r.flags.Ignorecache() {
		if v, err = r.checkCache(); err == nil && v != nil {
			return v, err
		}
	}

	var resp *http.Response
	var body []byte
	if resp, body, err = r.doRequest(); err != nil {
		return nil, err
	}

	if r.expectsStatusCode > 0 && resp.StatusCode != r.expectsStatusCode {
		err = &httd.ErrREST{
			HTTPCode: resp.StatusCode,
			Msg:      "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(r.expectsStatusCode),
		}
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

	if r.flags.Sort() {
		Sort(obj, r.flags)
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

	flags []Flag // TODO: checking

	prerequisites []string // error msg

	itemFactory fRESTItemFactory

	cache           *Cache
	cacheRegistry   cacheRegistry
	cacheMiddleware fRESTCacheMiddleware
	cacheItemID     Snowflake

	body              map[string]interface{}
	urlParams         urlQuery
	ignoreCache       bool
	cancelOnRatelimit bool
}

// addPrereq the naming here is kinda reversed..
// just think that each resembles a normal; if true => error
func (b *RESTBuilder) addPrereq(condition bool, errorMsg string) {
	if condition == false {
		return
	}
	b.prerequisites = append(b.prerequisites, errorMsg)
}

func (b *RESTBuilder) setup(cache *Cache, client httd.Requester, config *httd.Request, middleware fRESTRequestMiddleware) {
	b.body = make(map[string]interface{})
	b.urlParams = make(map[string]interface{})
	b.cache = cache
	b.client = client
	b.config = config
	b.middleware = middleware

	if b.config == nil {
		b.config = &httd.Request{
			Method: http.MethodGet,
		}
	}
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

	if b.cache == nil {
		b.IgnoreCache()
	}

	flags := mergeFlags(b.flags)
	if flags.Ignorecache() {
		b.IgnoreCache()
	}
}

// execute ... v must be a nil pointer.
func (b *RESTBuilder) execute() (v interface{}, err error) {
	for i := range b.prerequisites {
		return nil, errors.New(b.prerequisites[i])
	}

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

	if len(body) > 1 && b.itemFactory != nil {
		v = b.itemFactory()
		if err = httd.Unmarshal(body, v); err != nil {
			return nil, err
		}

		executeInternalUpdater(v)
		// executeInternalClientUpdater(disgord.Client, v)

		if b.cacheRegistry == NoCacheSpecified {
			return v, err
		}

		if b.cacheMiddleware != nil {
			b.cacheMiddleware(resp, v, err)
		}

		b.cache.Update(b.cacheRegistry, v)
	}
	if mergeFlags(b.flags).Sort() {
		Sort(v, b.flags...)
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

//generate-rest-basic-execute: err:error,
type basicBuilder struct {
	r RESTBuilder
}

// GetGateway [REST] Returns an object with a single valid WSS URL, which the Client can use for Connecting.
// Clients should cacheLink this value and only call this endpoint to retrieve a new URL if they are unable to
// properly establish a connection using the cached version of the URL.
//  Method                  GET
//  Endpoint                /gateway
//  Rate limiter            /gateway
//  Discord documentation   https://discordapp.com/developers/docs/topics/gateway#get-gateway
//  Reviewed                2018-10-12
//  Comment                 This endpoint does not require authentication.
func (c *Client) GetGateway(client httd.Getter) (gateway *websocket.Gateway, err error) {
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
func (c *Client) GetGatewayBot() (gateway *websocket.GatewayBot, err error) {
	var body []byte
	_, body, err = c.req.Get(&httd.Request{
		Ratelimiter: "/gateway/bot",
		Endpoint:    "/gateway/bot",
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &gateway)
	return
}

func exec(f func() (interface{}, error), flags ...Flag) (v interface{}, err error) {
	if v, err = f(); err != nil {
		return nil, err
	}

	if v == nil {
		return nil, errors.New("object was nil")
	}

	return v, nil
}

// TODO: auto generate
func getChannel(f func() (interface{}, error), flags ...Flag) (channel *Channel, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	return v.(*Channel), nil
}

// TODO: auto generate
func getChannels(f func() (interface{}, error), flags ...Flag) (channels []*Channel, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	if list, ok := v.(*[]*Channel); ok {
		return *list, nil
	} else if list, ok := v.([]*Channel); ok {
		return list, nil
	}
	panic("v was not assumed type. Got " + fmt.Sprint(v))
}

// TODO: auto generate
func getRole(f func() (interface{}, error), flags ...Flag) (role *Role, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	return v.(*Role), nil
}

// TODO: auto generate
func getRoles(f func() (interface{}, error), flags ...Flag) (roles []*Role, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	if list, ok := v.(*[]*Role); ok {
		return *list, nil
	} else if list, ok := v.([]*Role); ok {
		return list, nil
	}
	panic("v was not assumed type. Got " + fmt.Sprint(v))
}

// TODO: auto generate
func getMember(f func() (interface{}, error), flags ...Flag) (member *Member, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	return v.(*Member), nil
}

// TODO: auto generate
func getMembers(f func() (interface{}, error), flags ...Flag) (members []*Member, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	if list, ok := v.(*[]*Member); ok {
		return *list, nil
	} else if list, ok := v.([]*Member); ok {
		return list, nil
	}
	panic("v was not assumed type. Got " + fmt.Sprint(v))
}

// TODO: auto generate
func getWebhook(f func() (interface{}, error), flags ...Flag) (wh *Webhook, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	return v.(*Webhook), nil
}

// TODO: auto generate
func getWebhooks(f func() (interface{}, error), flags ...Flag) (whs []*Webhook, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	if list, ok := v.(*[]*Webhook); ok {
		return *list, nil
	} else if list, ok := v.([]*Webhook); ok {
		return list, nil
	}
	panic("v was not assumed type. Got " + fmt.Sprint(v))
}

// TODO: auto generate
func getMessage(f func() (interface{}, error), flags ...Flag) (msg *Message, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	return v.(*Message), nil
}

// TODO: auto generate
func getMessages(f func() (interface{}, error), flags ...Flag) (msgs []*Message, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	if list, ok := v.(*[]*Message); ok {
		return *list, nil
	} else if list, ok := v.([]*Message); ok {
		return list, nil
	}
	panic("v was not assumed type. Got " + fmt.Sprint(v))
}

// TODO: auto generate
func getUser(f func() (interface{}, error), flags ...Flag) (user *User, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	return v.(*User), nil
}

// TODO: auto generate
func getUsers(f func() (interface{}, error), flags ...Flag) (users []*User, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	if list, ok := v.(*[]*User); ok {
		return *list, nil
	} else if list, ok := v.([]*User); ok {
		return list, nil
	}
	panic("v was not assumed type. Got " + fmt.Sprint(v))
}

// TODO: auto generate
func getNickName(f func() (interface{}, error), flags ...Flag) (nick string, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return "", err
	}
	return v.(*nickNameResponse).Nickname, nil
}

// TODO: auto generate
func getBan(f func() (interface{}, error), flags ...Flag) (ban *Ban, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	return v.(*Ban), nil
}

// TODO: auto generate
func getEmoji(f func() (interface{}, error), flags ...Flag) (emoji *Emoji, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	return v.(*Emoji), nil
}

// TODO: auto generate
func getInvite(f func() (interface{}, error), flags ...Flag) (invite *Invite, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	return v.(*Invite), nil
}

// TODO: auto generate
func getInvites(f func() (interface{}, error), flags ...Flag) (invite []*Invite, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	if list, ok := v.(*[]*Invite); ok {
		return *list, nil
	} else if list, ok := v.([]*Invite); ok {
		return list, nil
	}
	panic("v was not assumed type. Got " + fmt.Sprint(v))
}

// TODO: auto generate
func getGuild(f func() (interface{}, error), flags ...Flag) (guild *Guild, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	return v.(*Guild), nil
}

// TODO: auto generate
func getIntegrations(f func() (interface{}, error), flags ...Flag) (integrations []*Integration, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	if list, ok := v.(*[]*Integration); ok {
		return *list, nil
	} else if list, ok := v.([]*Integration); ok {
		return list, nil
	}
	panic("v was not assumed type. Got " + fmt.Sprint(v))
}

// TODO: auto generate
func getVoiceRegions(f func() (interface{}, error), flags ...Flag) (regions []*VoiceRegion, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	if list, ok := v.(*[]*VoiceRegion); ok {
		return *list, nil
	} else if list, ok := v.([]*VoiceRegion); ok {
		return list, nil
	}
	panic("v was not assumed type. Got " + fmt.Sprint(v))
}

// TODO: auto generate
func getVoiceRegion(f func() (interface{}, error), flags ...Flag) (region *VoiceRegion, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	return v.(*VoiceRegion), nil
}

// TODO: auto generate
func getPartialInvite(f func() (interface{}, error), flags ...Flag) (invite *PartialInvite, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	return v.(*PartialInvite), nil
}

// TODO: auto generate
func getGuildEmbed(f func() (interface{}, error), flags ...Flag) (embed *GuildEmbed, err error) {
	var v interface{}
	if v, err = exec(f, flags...); err != nil {
		return nil, err
	}
	return v.(*GuildEmbed), nil
}
