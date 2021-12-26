package disgord

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/andersfylling/disgord/internal/constant"
	"github.com/andersfylling/disgord/json"

	"github.com/andersfylling/disgord/internal/httd"
)

type ErrRest = httd.ErrREST

// URLQueryStringer converts a struct of values to a valid URL query string
type URLQueryStringer interface {
	URLQueryString() string
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

type restStepDoRequest func() (resp *http.Response, body []byte, err error)
type rest struct {
	c          *Client
	flags      Flag // merge flags
	httpMethod string

	// item creation
	// pool is prioritized over factory
	pool    Pool
	factory func() interface{}

	conf *httd.Request

	// steps
	doRequest restStepDoRequest
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
		r.conf.PopulateMissing()
	}
	r.httpMethod = r.conf.Method
	r.doRequest = r.stepDoRequest
}

func (r *rest) bindParams(params interface{}) {
	if params == nil {
		return
	}
}

func (r *rest) stepDoRequest() (resp *http.Response, body []byte, err error) {
	if r.conf == nil {
		err = errors.New("missing httd.Request configuration")
		return
	}

	resp, body, err = r.c.req.Do(r.conf.Ctx, r.conf)
	return
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
	if err = json.Unmarshal(body, obj); err != nil {
		r.Put(obj)
		return nil, err
	}
	executeInternalUpdater(obj)

	return obj, nil
}

func (r *rest) Execute() (v interface{}, err error) {
	var resp *http.Response
	var body []byte
	if resp, body, err = r.doRequest(); err != nil {
		return nil, err
	}

	successful := func(code int) bool {
		return code >= http.StatusOK && code < 300
	}

	if !successful(resp.StatusCode) {
		err = &httd.ErrREST{
			HTTPCode: resp.StatusCode,
			Msg:      "unexpected http response code. Got " + resp.Status,
		}
		return nil, err
	}

	var obj interface{}
	if obj, err = r.processContent(body); err != nil {
		return nil, err
	}

	if r.flags.Sort() {
		Sort(obj, r.flags)
	}

	return obj, nil
}

type fRESTRequestMiddleware func(resp *http.Response, body []byte, err error) error
type fRESTCacheMiddleware func(resp *http.Response, v interface{}, err error) error
type fRESTItemFactory func() interface{}

//go:generate go run internal/generate/restbuilders/main.go

type RESTBuilder struct {
	middleware fRESTRequestMiddleware
	config     *httd.Request
	client     httd.Requester

	flags Flag // TODO: checking

	prerequisites []string // error msg

	itemFactory fRESTItemFactory

	body              map[string]interface{}
	urlParams         urlQuery
	ignoreCache       bool
	cancelOnRatelimit bool

	headerReason string
}

// addPrereq the naming here is kinda reversed..
// just think that each resembles a normal; if true => error
func (b *RESTBuilder) addPrereq(condition bool, errorMsg string) {
	if condition == false {
		return
	}
	b.prerequisites = append(b.prerequisites, errorMsg)
}

func (b *RESTBuilder) setup(client httd.Requester, config *httd.Request, middleware fRESTRequestMiddleware) {
	b.body = make(map[string]interface{})
	b.urlParams = make(map[string]interface{})
	b.client = client
	b.config = config
	b.middleware = middleware

	if b.config == nil {
		b.config = &httd.Request{
			Ctx:    context.Background(),
			Method: http.MethodGet,
		}
	}
}

func (b *RESTBuilder) prepare() {
	// update the config
	if b.config.ContentType != "" {
		b.config.Body = b.body
	}
	b.config.Endpoint += b.urlParams.URLQueryString()

	if b.flags.Ignorecache() {
		b.IgnoreCache()
	}
	if b.config.Ctx == nil {
		b.config.Ctx = context.Background()
	}
}

// execute ... v must be a nil pointer.
func (b *RESTBuilder) execute() (v interface{}, err error) {
	for i := range b.prerequisites {
		return nil, errors.New(b.prerequisites[i])
	}
	b.prepare()

	if b.headerReason != "" {
		b.config.Reason = b.headerReason
	}

	var resp *http.Response
	var body []byte
	resp, body, err = b.client.Do(b.config.Ctx, b.config)
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
		if err = json.Unmarshal(body, v); err != nil {
			return nil, err
		}
		executeInternalUpdater(v)
	}
	if b.flags.Sort() {
		Sort(v, b.flags)
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

var MissingRESTParamsErr = errors.New("this method requires REST parameters, but none were given")

type ClientQueryBuilderExecutables interface {
	// CreateGuild Create a new guild. Returns a guild object on success. Fires a Guild Create Gateway event.
	CreateGuild(guildName string, params *CreateGuildParams) (*Guild, error)

	// GetVoiceRegions Returns an array of voice region objects that can be used when creating servers.
	GetVoiceRegions() ([]*VoiceRegion, error)

	BotAuthorizeURL() (*url.URL, error)
	SendMsg(channelID Snowflake, data ...interface{}) (*Message, error)
}

type ClientQueryBuilder interface {
	WithContext(ctx context.Context) ClientQueryBuilderExecutables
	WithFlags(flags ...Flag) ClientQueryBuilderExecutables

	ClientQueryBuilderExecutables

	Invite(code string) InviteQueryBuilder
	Channel(cid Snowflake) ChannelQueryBuilder
	User(uid Snowflake) UserQueryBuilder
	CurrentUser() CurrentUserQueryBuilder
	Guild(id Snowflake) GuildQueryBuilder
	Gateway() GatewayQueryBuilder
	ApplicationCommand(appID Snowflake) ApplicationCommandQueryBuilder
}

type clientQueryBuilder struct {
	ctx    context.Context
	flags  Flag
	client *Client
}

func (c clientQueryBuilder) WithContext(ctx context.Context) ClientQueryBuilderExecutables {
	c.ctx = ctx
	return &c
}

func (c clientQueryBuilder) WithFlags(flags ...Flag) ClientQueryBuilderExecutables {
	c.flags = mergeFlags(flags)
	return &c
}

func (c clientQueryBuilder) WithContextAndFlags(ctx context.Context, flags ...Flag) ClientQueryBuilderExecutables {
	c.ctx = ctx
	c.flags = mergeFlags(flags)
	return &c
}

// SendMsg should convert all inputs into a single message. If you supply a object with an ID
// such as a channel, message, role, etc. It will become a reference.  If say the Message provided
// does not have an ID, the Message will populate a CreateMessage with it's fields.
//
// If you want to affect the actual message data besides .Content; provide a
// MessageCreateParams. The reply message will be updated by the last one provided.
func (c clientQueryBuilder) SendMsg(channelID Snowflake, data ...interface{}) (msg *Message, err error) {
	params := &CreateMessageParams{}
	addEmbed := func(e *Embed) error {
		if params.Embed != nil {
			return errors.New("can only send one embed")
		}
		params.Embed = e
		return nil
	}
	msgToParams := func(m *Message) (s string, err error) {
		if s, err = m.DiscordURL(); err != nil {
			// try to reference the message, otherwise use it to
			// populate the params
			if len(m.Embeds) > 1 {
				return "", errors.New("can only create a message with a single embed")
			} else if len(m.Embeds) > 0 {
				params.Embed = m.Embeds[0]
			}

			params.Content = m.Content
			params.Components = m.Components
			params.MessageReference = m.MessageReference
			params.SpoilerTagAllAttachments = m.SpoilerTagAllAttachments
			params.SpoilerTagContent = m.SpoilerTagContent
			params.Tts = m.Tts
			return "", nil
		}
		return s, nil
	}
	for i := range data {
		if data[i] == nil {
			continue
		}

		var s string
		switch t := data[i].(type) {
		case *CreateMessageParams:
			*params = *t
		case CreateMessageParams:
			*params = t
		case CreateMessageFileParams:
			params.Files = append(params.Files, t)
		case *CreateMessageFileParams:
			params.Files = append(params.Files, *t)
		case Embed:
			if err = addEmbed(&t); err != nil {
				return nil, err
			}
		case *Embed:
			if err = addEmbed(t); err != nil {
				return nil, err
			}
		case *os.File:
			return nil, errors.New("can not handle *os.File, use a CreateMessageFileParams instead")
		case string:
			s = t
		case Message:
			if s, err = msgToParams(&t); err != nil {
				return nil, err
			}
		case *Message:
			if s, err = msgToParams(t); err != nil {
				return nil, err
			}
		case AllowedMentions:
			params.AllowedMentions = &t
		case *AllowedMentions:
			params.AllowedMentions = t
		default:
			var mentioned bool
			if mentionable, ok := t.(Mentioner); ok {
				if s = mentionable.Mention(); len(s) > 5 {
					mentioned = true
				}
			}

			if !mentioned {
				if str, ok := t.(fmt.Stringer); ok {
					s = str.String()
				} else {
					s = fmt.Sprint(t)
				}
			}
		}

		if s != "" {
			params.Content += " " + s
		}
	}

	// wtf?
	if data == nil {
		if c.flags.IgnoreEmptyParams() {
			params.Content = ""
		} else {
			return nil, errors.New("params were nil")
		}
	}

	return c.Channel(channelID).WithContext(c.ctx).CreateMessage(params)
}

// BotAuthorizeURL creates a URL that can be used to invite this bot to a guild/server.
// Note that it depends on the bot ID to be after the Discord update where the Client ID
// is the same as the Bot ID.
//
// By default the permissions will be 0, as in none. If you want to add/set the minimum required permissions
// for your bot to run successfully, you should utilise
//  Client.
func (c clientQueryBuilder) BotAuthorizeURL() (*url.URL, error) {
	format := "https://discord.com/oauth2/authorize?scope=bot&client_id=%s&permissions=%d"
	u := fmt.Sprintf(format, c.client.botID.String(), c.client.permissions)
	return url.Parse(u)
}

func ensureDiscordGatewayURLHasQueryParams(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return urlString, err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return urlString, err
	}
	q.Add("encoding", constant.Encoding)
	q.Add("v", strconv.FormatUint(uint64(constant.DiscordVersion), 10))
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func exec(f func() (interface{}, error)) (v interface{}, err error) {
	if v, err = f(); err != nil {
		return nil, err
	}

	if v == nil {
		return nil, errors.New("object was nil")
	}

	return v, nil
}

// TODO: auto generate
func getChannel(f func() (interface{}, error)) (channel *Channel, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*Channel), nil
}

// TODO: auto generate
func getChannels(f func() (interface{}, error)) (channels []*Channel, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
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
func getRole(f func() (interface{}, error)) (role *Role, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*Role), nil
}

// TODO: auto generate
func getRoles(f func() (interface{}, error)) (roles []*Role, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
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
func getMember(f func() (interface{}, error)) (member *Member, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*Member), nil
}

// TODO: auto generate
func getMembers(f func() (interface{}, error)) (members []*Member, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
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
func getWebhook(f func() (interface{}, error)) (wh *Webhook, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*Webhook), nil
}

// TODO: auto generate
func getWebhooks(f func() (interface{}, error)) (whs []*Webhook, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
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
func getMessage(f func() (interface{}, error)) (msg *Message, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*Message), nil
}

// TODO: auto generate
func getMessages(f func() (interface{}, error)) (msgs []*Message, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
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
func getUser(f func() (interface{}, error)) (user *User, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*User), nil
}

// TODO: auto generate
func getUsers(f func() (interface{}, error)) (users []*User, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
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
func getNickName(f func() (interface{}, error)) (nick string, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return "", err
	}
	return v.(*nickNameResponse).Nickname, nil
}

// TODO: auto generate
func getBan(f func() (interface{}, error)) (ban *Ban, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*Ban), nil
}

// TODO: auto generate
func getEmoji(f func() (interface{}, error)) (emoji *Emoji, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*Emoji), nil
}

// TODO: auto generate
func getInvite(f func() (interface{}, error)) (invite *Invite, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*Invite), nil
}

// TODO: auto generate
func getInvites(f func() (interface{}, error)) (invite []*Invite, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
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
func getGuild(f func() (interface{}, error)) (guild *Guild, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*Guild), nil
}

// TODO: auto generate
func getIntegrations(f func() (interface{}, error)) (integrations []*Integration, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
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
func getVoiceRegions(f func() (interface{}, error)) (regions []*VoiceRegion, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
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
func getVoiceRegion(f func() (interface{}, error)) (region *VoiceRegion, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*VoiceRegion), nil
}

// TODO: auto generate
func getPartialInvite(f func() (interface{}, error)) (invite *PartialInvite, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*PartialInvite), nil
}

// TODO: auto generate
func getGuildEmbed(f func() (interface{}, error)) (embed *GuildEmbed, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*GuildEmbed), nil
}

// TODO: auto generate
func getThreadMember(f func() (interface{}, error)) (threadMember *ThreadMember, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*ThreadMember), nil
}

// TODO: auto generate
func getThreadMembers(f func() (interface{}, error)) (threadMembers []*ThreadMember, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	if list, ok := v.(*[]*ThreadMember); ok {
		return *list, nil
	} else if list, ok := v.([]*ThreadMember); ok {
		return list, nil
	}
	panic("v was not assumed type. Got " + fmt.Sprint(v))
}

func getResponseBodyThreads(f func() (interface{}, error)) (concreteBody *ResponseBodyThreads, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*ResponseBodyThreads), nil
}

func getResponseBodyGuildThreads(f func() (interface{}, error)) (concreteBody *ResponseBodyGuildThreads, err error) {
	var v interface{}
	if v, err = exec(f); err != nil {
		return nil, err
	}
	return v.(*ResponseBodyGuildThreads), nil
}
