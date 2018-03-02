package request

type Requester interface {
	Request(method, uri string)
	Get(uri string)
	Post(uri string)
	Put(uri string)
	Patch(uri string)
	Delete(uri string)
}

func NewClient(conf *Config) *Client {
	return &Client{}
}

type Config struct{}

type Client struct {
	url string // base url with API version
}
