# Disgord

## Health
| Branch       | Build status  | Maintainability | Test Coverage | Comment Coverage |
| ------------ |:-------------:|:---------------:|:-------------:|:----------------:|
| master       | [![CircleCI](https://circleci.com/gh/andersfylling/disgord/tree/master.svg?style=shield)](https://circleci.com/gh/andersfylling/disgord/tree/master) | [![Maintainability](https://api.codeclimate.com/v1/badges/687d02ca069eba704af9/maintainability)](https://codeclimate.com/github/andersfylling/disgord/maintainability) | [![Test Coverage](https://api.codeclimate.com/v1/badges/687d02ca069eba704af9/test_coverage)](https://codeclimate.com/github/andersfylling/disgord/test_coverage) | [![Coverage Status](https://coveralls.io/repos/github/andersfylling/disgord/badge.svg)](https://coveralls.io/github/andersfylling/disgord) |

## About
This library is split into three parts: caching, requests and events.
Disgord is currently under heavy development and should not be used for production. Contributions are welcome.

Objects will have methods to simplify interaction: `User.sendMsgStr(...)`, or `Channel.SendMsgStr(...)`.


#### Caching
The cache, of discord objects, aims to reflect the same state as of the discord servers. Therefore incoming data is deep copied, as well as return values from the cache.
This lib handles caching for you, so whenever you send a request to the REST API or receive a discord event. The contents are cached auto-magically to a separate memory space.

As a structure is sent into the cache module, everything is deep copied as mentioned, however if the object hold discord objects consistent of a snowflake, it does not do a deep copy. It converts given field to a nil, and stores the snowflake only in a separate struct/map. This makes sure that there will only exist one version of an object. Making updating fairly easy.
When the object goes out of cache, a copy is created and every sub object containing a snowflake is deep copied from the cache as well, to create a wholesome object.

#### Requests
For every REST API request (which is the only way to get objects from the discord interface, without waiting for changes as events) the request is rate limited auto-magically by the library (caching coming soon for resource funcs).
The functions in `resource` pkg are blocking, and should be used with care. For async requests, use the methods found at the `Session` interface, such as:
`Session.User(userID)` which returns a channel. The channel will get content from the REST API, if not found in the cache.

#### Events
The reactor pattern with goroutines, or a pro-actor pattern is used. This will always be the default behavior, synchronous triggering of listeners might be implemented in the future.
Incoming events from the discord servers are parsed into respective structs and dispatched to either a) callbacks, or b) through channels. Both are dispatched from the same place, and the arguments share the same memory space. So it doesn't matter which one you pick, chose your preference.

## Quick example

```go
package main

import (
    "os"
    "os/signal"
    "syscall"

    "github.com/andersfylling/disgord"
    "github.com/sirupsen/logrus"
)

func main() {
    termSignal := make(chan os.Signal, 1)
    signal.Notify(termSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

    sess, err := disgord.NewSession(&disgord.Config{
        Token: os.Getenv("DISGORD_TOKEN"),
    })
    if err != nil {
        panic(err)
    }

    // add a event listener
    sess.AddListener(event.GuildCreateKey, func(session Session, box *event.GuildCreateBox) {
        guild := box.Guild
        // do something with guild
    })

    // or use a channel to listen for events
    go func() {
        for {
            select {
            case box, alive := <- sess.Evt().GuildCreateChan():
                if !alive {
                    fmt.Println("channel is dead")
                    break
                }

                guild := box.Guild
                // do something with guild
            }
        }
    }()

    // connect to the discord gateway to receive events
    err = sess.Connect()
    if err != nil {
        panic(err)
    }

    // eg. retrieve a specific user from the Discord servers
    userID := snowflake.NewID(228846961774559232)
    user := <- sess.User(userID) // sends a request to discord
    user2 := <- sess.User(userID) // does a cache look up, to prevent rate limiting/banning

    // keep the app alive until terminated
    <-termSignal
    sess.Disconnect()
}
```

## WARNING
All the REST endpoints are implemented, but may not exist on the interface yet. Create a Disgord session/client and use the REST functions found in the rest package directly (for now).
See [using the rest functions directly](#using-the-rest-functions-directly).

## Code flow

The main design takes in a discord event and dispatches the event to a channel/callback suited for the event type. The channel can be retrieved and the callbacks set by the Session interface: `Session.Event.ChannelDeleteChan()`, `Session.Event.AddHandler(event.GuildCreateKey, func(...){})`

Note that callbacks and channels are fired from the same place, to avoid overhead. However, channels are fired before the callbacks.

## Q&A

```Markdown
1. Reason for making another Discord lib in GoLang?

I'm trying to take over the world and then become a intergalactic war lord. Have to start somewhere.
```

## Using the REST functions directly
The goal is to have all the functions in the main interface (session) to handle caching and concurrency. However, not every method is implemented in the session interface and the caching is still under development. So it might be of interest to use the functions residing within the rest package directly.

Each REST function does require a client which is found in the session interface:
```GoLang
Session.Req() // request client for REST
```

The typical design for each REST function is to take the request client, some url parameters and optionally a query object and/or a json object:
```GoLang
func GetGuildEmoji(client httd.Getter, guildID, emojiID snowflake.ID) (ret *Emoji, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/emojis/" + emojiID.String(),
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}

// -----
// usage
emoji, err := rest.GetGuildEmoji(session.Req(), guildID, emojiID)
if err != nil {
    return err
}

// tada, you now have queried the discord API for an emoji within a guild.
```







## Thanks to
* [github.com/s1kx](https://github.com/s1kx) for different design suggestions.

