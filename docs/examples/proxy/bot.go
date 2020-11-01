package main

import (
	"os"

	"golang.org/x/net/proxy"

	"github.com/andersfylling/disgord"
)

// In the event that the Discord connections need to be routed
// through a proxy, you can do so by using this approach. In
// this example we will be using SOCKS5, but any custom
// implementation can be used as long as they satisfy the
// proxy.Dialer interface.
func main() {
	p, err := proxy.SOCKS5("tcp", "localhost:8080", nil, proxy.Direct)
	if err != nil {
		panic(err)
	}

	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISCORD_TOKEN"),
		Proxy:    p, // Anything satisfying the proxy.Dialer interface will work
	})
	defer client.Gateway().StayConnectedUntilInterrupted()
}
