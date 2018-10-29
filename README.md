# Disgord [![Documentation](https://godoc.org/github.com/andersfylling/disgord?status.svg)](http://godoc.org/github.com/andersfylling/disgord)
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/contains-technical-debt.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/for-you.svg)](https://forthebadge.com)

## Health
| Branch       | Build status  | Code climate | Go Report Card | Codacy |
| ------------ |:-------------:|:---------------:|:-------------:|:----------------:|
| develop     | [![CircleCI](https://circleci.com/gh/andersfylling/disgord/tree/develop.svg?style=shield)](https://circleci.com/gh/andersfylling/disgord/tree/develop) | [![Maintainability](https://api.codeclimate.com/v1/badges/687d02ca069eba704af9/maintainability)](https://codeclimate.com/github/andersfylling/disgord/maintainability) | [![Go Report Card](https://goreportcard.com/badge/github.com/andersfylling/disgord)](https://goreportcard.com/report/github.com/andersfylling/disgord) | [![Codacy Badge](https://api.codacy.com/project/badge/Grade/a8b2edae3c114dadb7946afdc4105a51)](https://www.codacy.com/project/andersfylling/disgord/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=andersfylling/disgord&amp;utm_campaign=Badge_Grade_Dashboard) |

## About
GoLang module for interacting with the Discord API. Supports socketing and REST functionality. Discord object will also have implemented helper functions such as `Message.RespondString(session, "hello")`, or `Session.SaveToDiscord(&Emoji)` for simplicity/readability.

Disgord has complete implementation for Discord's documented REST API. It lacks comprehensive testing, although unit-tests have been created for several of the Disgord REST implementations. The socketing is not complete, but does support all event types that are documented (using both channels and callbacks).

Disgord does not utilize reflection, except in unit tests and unmarshalling/marshalling of JSON. But does return custom error messages for some functions which can be type checked in a switch for a more readable error handling as well potentially giving access to more information.

To get started see the examples in [docs](docs/examples)

Alternative GoLang package for Discord: [DiscordGo](https://github.com/bwmarrin/discordgo)

Discord channel/server: [Discord Gophers#Disgord](https://discord.gg/qBVmnq9)

## Issues and behavior you must be aware of
Channels are not correctly implemented (see [issue #78](https://github.com/andersfylling/disgord/issues/78)). So information is lost at random. For now, only use handlers.

The cache is not complete (see [issue #65](https://github.com/andersfylling/disgord/issues/65)). It might also have bugs, and require more unit tests and stress testing would be great.
The idea behind it is to make it as configurable as possible such that you as a developer can experiment and tweak your cache setup to greater performance.

Once we have enough insight/discussion on what works best, we will use build constraints to let people choose between the very configurable cache (current), and then a more performance focused cache (future default).

## Package structure
None of the sub-packages should be used outside the library. If there exists a requirement for that, please create an issue or pull request.
```Markdown
github.com/andersfylling/disgord
└──.circleci    :CircleCI configuration
└──cache        :Different cache replacement algorithms
└──constant     :Constants such as version, GitHub URL, etc.
└──docs         :Examples, templates, (documentation)
└──endpoint     :All the REST endpoints of Discord
└──event        :All the Discord event identifiers
└──generate     :All go generate scripts for "generic" code
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

### Build constraints
> For **advanced users only**.

If you do not wish to use json-iterator, you can pass `-tags=json-std` to switch to `"encoding/json"`.
However, json-iterator is the recommended default for this library.

Disgord has the option to use mutexes (sync.RWMutex) on Discord objects. By default, methods of Discord objects are not locked as this is
not needed in our event driven architecture unless you create a parallel computing environment.
If you want the internal methods to deal with read-write locks on their own, you can pass `-tags=parallelism`, which will activate built-in locking.
Making all methods thread safe.

If you want to remove the extra memory used by mutexes, or you just want to completely avoid potential deadlocks by disabling
mutexes you can pass `-tags=removeDiscordMutex` which will replace the RWMutex with an empty struct, causing mutexes (in Discord objects only) to be removed at compile time.
This cannot be the default behaviour, as it creates confusion whether or not a mutex exists and leads to more error prone code. The developer has to be aware themselves whether or not
their code can be run without the need of mutexes. This option is not affected by the `parallelism` tag.

## Setup / installation guide
As this is a go module, it is expected that your project utilises the module concept (minimum Go version: 1.11). If you do not, then there is no guarantee that using this will work. To get this, simply use go get: `go get github.com/andersfylling/disgord`. I have been using this project in none module projects, so it might function for you as well. But official, this is not supported.

Read more about modules here: [https://github.com/golang/go/wiki/Modules](https://github.com/golang/go/wiki/Modules)

### Creating a fresh project using Disgord
So if you haven't used modules before and you just want to create a Bot using Disgord, this is how it's done (Linux):
 1. Create a folder with your project name: `mkdir my-bot && cd my-bot` (outside the go path!)
 2. Create a main.go file, and add the following:
    ```go
    package main

    import "github.com/andersfylling/disgord"
    import "fmt"

    func main() {
        session, err := disgord.NewSession(&disgord.Config{
            Token: "DISGORD_TOKEN",
        })
        if err != nil {
            panic(err)
        }

        myself, err := session.GetCurrentUser()
        if err != nil {
            panic(err)
        }

        fmt.Printf("Hello, %s!\n", myself.String())
    }
    ```
 3. Make sure you have activated go modules: `export GO111MODULE=auto`
 4. Initiate the project as a module: `go mod init my-bot` (you should now see a `go.mod` file)
 5. Start building, this will find all your dependencies and store them in the go.mod file: `go build .`
 6. You can now start the bot, and see the greeting: `go run .`

If you experience any issues with this guide, please create a issue.

## Contributing
Please see the [CONTRIBUTING.md file](CONTRIBUTING.md) (Note that it can be useful to read this regardless if you have the time)

## Git branching model
The branch:develop holds the most recent changes, as it name implies. There is no master branch as there will never be a "stable latest branch" except the git tags (or releases).

## Mental model
#### Caching
The cache can be either immutable (recommended) or mutable. When the cache is mutable you will share the memory space with the cache, such that if you change your data structure you might also change the cache directly. However, by using the immutable option all incoming data is deep copied to the cache and you will not be able to directly access the memory space, this should keep your code less error-prone and allow for concurrent cache access in case you want to use channels or other long-running tasks/processes.

#### Requests
For every REST API request the request is rate limited and cached auto-magically by Disgord. This means that when you utilize the Session interface you won't have to worry about rate limits and data is cached to improve performance. See the GoDoc for how to bypass the caching.

#### Events
The reactor pattern is used. This will always be the default behavior, however channels will ofcourse work more as a pro-actor system as you deal with the data parallel to other functions.
Incoming events from the discord servers are parsed into respective structs and dispatched to either a) handlers, or b) through channels. Both are dispatched from the same place, and the arguments share the same memory space. Pick handlers (register them using Session.On method) simplicity as they run in sequence, while channels are executed in a parallel setting (it's expected you understand how channels work so I won't go in-depth here).

## Quick example
> **NOTE:** To see more examples go visit the docs/examples folder.
See the GoDoc for a in-depth introduction on the various topics (or disgord.go package comment). Below is an example of the traditional ping-pong bot.
```go
// create a Disgord session
session, err := disgord.NewSession(&disgord.Config{
    Token: os.Getenv("DISGORD_TOKEN"),
})
if err != nil {
    panic(err)
}

// create a handler and bind it to new message events
session.On(disgord.EventMessageCreate, func(session disgord.Session, data *disgord.MessageCreate) {
    msg := data.Message

    if msg.Content == "ping" {
        msg.RespondString(session, "pong")
    }
})

// connect to the discord gateway to receive events
err = session.Connect()
if err != nil {
    panic(err)
}

// Keep the socket connection alive, until you terminate the application
session.DisconnectOnInterrupt()
```

## Q&A

```Markdown
1. Reason for making another Discord lib in GoLang?

I'm trying to take over the world and then become a intergalactic war lord. Have to start somewhere.
```
