# Disgord
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/contains-technical-debt.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/for-you.svg)](https://forthebadge.com)

## Health
| Branch       | Build status  | Code climate | Go Report Card | Codacy |
| ------------ |:-------------:|:---------------:|:-------------:|:----------------:|
| master       | [![CircleCI](https://circleci.com/gh/andersfylling/disgord/tree/master.svg?style=shield)](https://circleci.com/gh/andersfylling/disgord/tree/master) | [![Maintainability](https://api.codeclimate.com/v1/badges/687d02ca069eba704af9/maintainability)](https://codeclimate.com/github/andersfylling/disgord/maintainability) | [![Go Report Card](https://goreportcard.com/badge/github.com/andersfylling/disgord)](https://goreportcard.com/report/github.com/andersfylling/disgord) | [![Codacy Badge](https://api.codacy.com/project/badge/Grade/a8b2edae3c114dadb7946afdc4105a51)](https://www.codacy.com/project/andersfylling/disgord/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=andersfylling/disgord&amp;utm_campaign=Badge_Grade_Dashboard) |
| develop     | [![CircleCI](https://circleci.com/gh/andersfylling/disgord/tree/develop.svg?style=shield)](https://circleci.com/gh/andersfylling/disgord/tree/develop) | - | - | [![Codacy Badge](https://api.codacy.com/project/badge/Grade/a8b2edae3c114dadb7946afdc4105a51)](https://www.codacy.com/project/andersfylling/disgord/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=andersfylling/disgord&amp;utm_campaign=Badge_Grade_Dashboard) |



## About
GoLang library for interacting with the Discord API. Supports socketing and REST functions. Compared to Discordgo this library will be a little more heavy and support helper functions on objects such as Channel, Message, etc.
 eg. `Channel.SendMsg(...)`, `Message.Respond(...)`.

To get started see the examples in [docs](docs/examples)

## package structure
```Markdown
github.com/andersfylling/disgord
└──discordws    :Deals with the discord socket connection
└──docs         :Examples, templates, (documentation)
└──event        :Data structures, callbacks, event types
└──resource     :All the Discord data structures (same setup as the Discord docs)
└──rest         :All the endpoints found in the documentation (same as resource)
└──state/cache  :Logic for caching incoming Discord information
```

## Contributing
Please see the [CONTRIBUTING.md file](CONTRIBUTING.md)

## The Wiki
Yes, the wiki might hold some information. But I believe everything will be placed within the "docs" package in the end.

## Mental model

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
    userID := NewSnowflake(228846961774559232)
    userResponse := <- sess.User(userID) // sends a request to discord
    userResponse2 := <- sess.User(userID) // does a cache look up, to prevent rate limiting/banning

    // check if there was an issue (eg. rate limited or not found)
    if userResponse.Err != nil {
        panic(userResponse.Err)
    }

    // check if this is retrieved from the cache
    if userResponse.Cache {
        // ...
    }

    // get the user info
    user := userResponse.User

    // keep the app alive until terminated
    <-termSignal
    sess.Disconnect()
}
```

## WARNING
All the REST endpoints are implemented, but may not exist on the interface yet. Create a Disgord session/client and use the REST functions found in the rest package directly (for now). See the examples in docs for using the functions in the rest package directly.

## Q&A

```Markdown
1. Reason for making another Discord lib in GoLang?

I'm trying to take over the world and then become a intergalactic war lord. Have to start somewhere.
```







## Thanks to
* [github.com/s1kx](https://github.com/s1kx) for different design suggestions.

