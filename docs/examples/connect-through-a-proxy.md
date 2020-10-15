In the event that the Discord connections need to be routed through a proxy, you can do so by using this approach.
In this example we will be using SOCKS5, but any custom implementation can be used as long as they satisfy the `proxy.Dialer` interface.

```go
package main

import (
	"context"
	"os"

	"github.com/andersfylling/disgord"
	"golang.org/x/net/proxy"
)

func main() {
	p, err := proxy.SOCKS5("tcp", "localhost:8080", nil, proxy.Direct)
	if err != nil {
		panic(err)
	}
	
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
		Proxy:    p, // Anything satisfying the proxy.Dialer interface will work
	})

	if err := client.Connect(context.Background()); err != nil {
		panic(err)
	}

	client.DisconnectOnInterrupt()
}
```
