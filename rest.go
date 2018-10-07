package disgord

import (
	"github.com/andersfylling/disgord/httd"
)

// URLParameters converts a struct of values to a valid URL query string
type URLParameters interface {
	GetQueryString() string
}

func unmarshal(data []byte, v interface{}) error {
	return httd.Unmarshal(data, v)
}

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