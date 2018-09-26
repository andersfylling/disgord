In the event that the Discord connections need to be routed through a proxy, you can do so by using this approach.
In this example we will be using SOCKS5, but any custom implementation can be used.

```go
package main

import (
	"net/http"
	"net/url"
	"time"

	"github.com/andersfylling/disgord"
)

func main() {
	discord := disgord.NewSessionMustCompile(&disgord.Config{
		Token: "my.very.secret.bot.token",
		HTTPClient: &http.Client{
			Timeout: time.Second * 10, // important, otherwise the timeout will be infinite
			Transport: &http.Transport{
				Proxy: func(request *http.Request) (*url.URL, error) {
					return url.Parse("socks5://localhost:8080")
				},
				// You can also set the Dial and DialContext functions instead
			},
		},
	})

	if err := discord.Connect(); err != nil {
		panic(err)
	}

	discord.DisconnectOnInterrupt()
}
```