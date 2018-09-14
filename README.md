# Disgord
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/contains-technical-debt.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/for-you.svg)](https://forthebadge.com)

## Health
| Branch       | Build status  | Code climate | Go Report Card | Codacy |
| ------------ |:-------------:|:---------------:|:-------------:|:----------------:|
| develop     | [![CircleCI](https://circleci.com/gh/andersfylling/disgord/tree/develop.svg?style=shield)](https://circleci.com/gh/andersfylling/disgord/tree/develop) | [![Maintainability](https://api.codeclimate.com/v1/badges/687d02ca069eba704af9/maintainability)](https://codeclimate.com/github/andersfylling/disgord/maintainability) | [![Go Report Card](https://goreportcard.com/badge/github.com/andersfylling/disgord)](https://goreportcard.com/report/github.com/andersfylling/disgord) | [![Codacy Badge](https://api.codacy.com/project/badge/Grade/a8b2edae3c114dadb7946afdc4105a51)](https://www.codacy.com/project/andersfylling/disgord/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=andersfylling/disgord&amp;utm_campaign=Badge_Grade_Dashboard) |

## WARNING
Missing caching. It's under development, but in the current state it does not exist.

## About
GoLang module for interacting with the Discord API. Supports socketing and REST functionality. Discord object will also have implemented helper functions such as `Message.RespondString(session, "hello")`, or `Channel.SendMsg(session, &Message{...})` for simplicity/readability.

Disgord has complete implementation for Discord's documented REST API. It lacks comprehensive testing, although unit-tests have been created for several of the disgord REST implementations. The socketing is not complete, but does support all event types (using both channels and callbacks).

Note that caching is yet to be implemented. Versions from v0.5.1 and below, had caching to some degree, but was scrapped once a complete rework of the project structure was done.

Disgord does not utilize reflection, except in unit tests and unmarshalling/marshalling of JSON.

To get started see the examples in [docs](docs/examples)

Alternative GoLang package for Discord: [DiscordGo](https://github.com/bwmartin/discordgo)

## Package structure
None of the sub-packages should be used outside the library. If there exists a requirement for that, please create an issue or pull request.
```Markdown
github.com/andersfylling/disgord
└──.circleci    :CircleCI configuration
└──constant     :Constants such as version, GitHub URL, etc.
└──docs         :Examples, templates, (documentation)
└──endpoint     :All the REST endpoints of Discord
└──httd         :Deals with rate limits and http calls
└──testdata     :Holds all test data for unit tests (typically JSON files)
└──websocket    :Discord Websocket logic (reconnect, resume, etc.)
```

### Dependencies
```Markdown
github.com/andersfylling/disgord
└──github.com/andersfylling/snowflake  :The snowflake ID designed for Discord
└──github.com/json-iterator/go         :For faster JSON decoding/encoding
└──github.com/sergi/go-diff            :Unit testing for checking JSON encoding/decoding of structs
└──github.com/sirupsen/logrus          :Logging (will be replaced with a simplified interface for DI)
```

## Contributing
Please see the [CONTRIBUTING.md file](CONTRIBUTING.md)

## Git branching model
The branch:develop holds the most recent changes, as it name implies. There is no master branch as there will never be a "stable latest branch" except the git tags (or releases).

## Mental model
#### Caching
The cache, of discord objects, aims to reflect the same state as of the discord servers. Therefore incoming data is deep copied, as well as return values from the cache. This lib handles caching for you, so whenever you send a request to the REST API or receive a discord event. The contents are cached auto-magically to a separate memory space.

As a structure is sent into the cache module, everything is deep copied as mentioned, however if the object hold discord objects consistent of a snowflake, it does not do a deep copy. It converts given field to a nil, and stores only the snowflake in a separate struct/map. This makes sure that there will only exist one version of an object. Making updating fairly easy.
When the object goes out of cache, a copy is created and every sub object containing a snowflake is deep copied from the cache as well, to return a wholesome object.

#### Requests
For every REST API request the request is rate limited auto-magically by disgord. The functions in `resource` pkg are blocking, and should be used with care. In the future there might be implementations of channel methods if requested.

#### Events
The reactor pattern with goroutines, or a pro-actor pattern is used. This will always be the default behavior, synchronous triggering of listeners might be implemented in the future as an option.
Incoming events from the discord servers are parsed into respective structs and dispatched to either a) callbacks, or b) through channels. Both are dispatched from the same place, and the arguments share the same memory space. So it doesn't matter which one you pick, chose your preference.

## Quick example
> **NOTE:** To see more examples go visit the docs/examples folder.

The following example is used as a prerequisite for the coming examples later on in this README file.
```go
var err error
termSignal := make(chan os.Signal, 1)
signal.Notify(termSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

sess, err := disgord.NewSession(&disgord.Config{
  Token: os.Getenv("DISGORD_TOKEN"),
})
if err != nil {
    panic(err)
}
```

Listening for events can be done in two ways. Firstly, the reactor pattern and secondly, a GoLang channel:
```GoLang
// add a event listener
sess.AddListener(event.KeyGuildCreate, func(session Session, data *event.GuildCreate) {
  guild := data.Guild
  // do something with guild
})

// or use a channel to listen for events
go func() {
    for {
        select {
        case data, alive := <- sess.Evt().GuildCreateChan():
            if !alive {
                fmt.Println("channel is dead")
                break
            }

            guild := data.Guild
            // do something with guild
        }
    }
}()

// connect to the discord gateway to receive events
err = sess.Connect()
if err != nil {
    panic(err)
}
```

Remember that when you call Session.Connect() it is recommended to call Session.Disconnect for a graceful shutdown and closing channels and Goroutines.
```GoLang
// keep the app alive until terminated
<-termSignal
sess.Disconnect()
```

To retrieve information from the Discord REST API you utilize the Session interface as it will in the future implement features such as caching, control checks, etc
```GoLang
// retrieve a specific user from the Discord API
var user *resource.User
userID := NewSnowflake(228846961774559232)
user, err = session.GetUser(userID) // will do a cache lookup in the future
if err != nil {
   panic(err)
}
```
However, if you think the Session interface is incorrect (outdated cache, or another issue) you can bypass the interface and call the REST method directly while you wait for a patch:
```GoLang
// bypassing the cache
user, err = rest.GetUser(session.Req(), userID)
if err != nil {
   panic(err)
}
```

There's also another way to retrieve content: channels. These methods will return a GoLang channel to help with concurrency. However, their implementation is down prioritized and I recommend using the normal REST methods for now. (Currently only the Session.User method is working, and might be temporary deprecated).
```GoLang
// eg. retrieve a specific user from the Discord API using GoLang channels
userResponse := <- sess.UserChan(userID) // sends a request to discord
userResponse2 := <- sess.UserChan(userID) // does a cache look up, to prevent rate limiting/banning

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
```

## Q&A

```Markdown
1. Reason for making another Discord lib in GoLang?

I'm trying to take over the world and then become a intergalactic war lord. Have to start somewhere.
```







## Thanks to
* [github.com/s1kx](https://github.com/s1kx) for different design suggestions.
