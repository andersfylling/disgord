package main

import (
	"context"
	"net"
	"net/http"
	"os"

	"golang.org/x/net/proxy"

	"github.com/andersfylling/disgord"
)

// In the event that the Discord connections need to be routed
// through a proxy, you can do so by using this approach. In
// this example we will be using SOCKS5, but any custom
// implementation can be used. You just configure your own http client.
//
// For REST methods the only Do method is required. So any configuration, libraries, whatever that 
// implements the Do method is good enough.
//
// For websocket connection you must specify the WebsocketHttpClient config option. Currently there is a issue
// when specifying http.Client timeouts for websocket, which is why you have the option to specify both.
// When a WebsocketHttpClient is not specified, a default config is utilised.
func main() {
	p, err := proxy.SOCKS5("tcp", "localhost:8080", nil, proxy.Direct)
	if err != nil {
		panic(err)
	}
	
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return p.Dial(network, addr)
			},
		},
	}

	client := disgord.New(disgord.Config{
		BotToken:   os.Getenv("DISCORD_TOKEN"),
		HttpClient: httpClient, // REST requests with proxy support
		WebsocketHttpClient: httpClient, // Websocket setup with proxy support
	})
	defer client.Gateway().StayConnectedUntilInterrupted()
}
