<div align='center'>
  <img src="/docs/disgord-draft-8.jpeg" alt='Build Status' />
  <p>
    <a href='https://circleci.com/gh/andersfylling/disgord/tree/develop'>
      <img src='https://circleci.com/gh/andersfylling/disgord/tree/develop.svg?style=shield'
           alt='Build Status' />
    </a>
    <a href='https://codeclimate.com/github/andersfylling/disgord/test_coverage'>
      <img src='https://api.codeclimate.com/v1/badges/687d02ca069eba704af9/test_coverage'
           alt='Maintainability' />
    </a>
    <a href='https://goreportcard.com/report/github.com/andersfylling/disgord'>
      <img src='https://goreportcard.com/badge/github.com/andersfylling/disgord'
           alt='Code coverage' />
    </a>
  </p>
  <p>
    <a href='http://godoc.org/github.com/andersfylling/disgord'>
      <img src='https://godoc.org/github.com/andersfylling/disgord?status.svg'
           alt='Godoc' />
    </a>
  </p>
</div>

## About
Go module for interacting with the Discord API. Supports events, REST calls and voice (sending only).

The goal is to make bot development easy and handle some nastiness internally; sharding, auto-scaling of shards, caching, provide helper functions, middlewares for events, allow concurrent use of rate limiters, etc.

DisGord has complete implementation for Discord's documented REST API. It lacks battle testing, so any bug report/feedback is greatly appreciated!

To get started see the examples in [docs](docs/examples)

Some projects using DisGord can be found [here](docs/PROJECTS.md).

Talk to us on Discord! We exist in both the Gopher server and the Discord API server:
 - [Discord Gophers](https://discord.gg/qBVmnq9)
 - [Discord API](https://discord.gg/HBTHbme)

## Warning
The develop branch is under continuous breaking changes, as the interface and exported funcs/consts are still undergoing planning. Because DisGord is under development and pushing for a satisfying interface, the SemVer logic is not according to spec. Until v1.0.0, every minor release is considered possibly breaking and patch releases might contain additional features. Please see the issue and current PR's to get an idea about coming changes before v1.

There might be bugs in the cache, or the cache processing might not exist yet for some REST methods. Bypass the cache for REST methods by supplying the flag argument `disgord.IgnoreCache`. eg. `client.GetCurrentUser(disgord.IgnoreCache)`.

Remember to read the docs/code for whatever version of disgord you are using. This README file reflects the latest state in the develop branch, or at least, I try to reflect the latest state.

## Starter guide (Linux)
> Note! this is a Go module project, and Go module support should activated to properly use DisGord. It might work using only the GOPATH. But officially this is not supported: Read more about modules here: [https://github.com/golang/go/wiki/Modules](https://github.com/golang/go/wiki/Modules) 

To create a new bot you can use the disgord.sh script to automate the boring copy/paste process. Paste the following into your terminal:

```bash
bash <(curl -s -L https://git.io/disgord-script)
``` 
> Remember to activate module support. Your IDE might require you to activate it in the settings menu.

Starter guide as a gif: https://terminalizer.com/view/469961d0695


## Architecture & Behavior
Discord provide communication in different forms. DisGord tackles the main ones, events (ws), voice (udp + ws), and REST calls.

You can think of DisGord as layered, in which case it will look something like:
![Simple way to think about DisGord architecture from a layered perspective](docs/disgord-layered-version.png)

#### Events
For Events, DisGord uses the [reactor pattern](https://dzone.com/articles/understanding-reactor-pattern-thread-based-and-eve). Every incoming event from Discord is processed and checked if any handler is registered for it, otherwise it's discarded to save time and resource use. Once a desired event is received, DisGord starts up a Go routine and runs all the related handlers in sequence; avoiding locking the need to use mutexes the handlers. 

In addition to traditional handlers, DisGord allows you to use Go channels. Note that if you use more than one channel per event, one of the channels will randomly receive the event data; this is how go channels work. It will act as a randomized load balancer.

But before either channels or handlers are triggered, the cache is updated.

#### REST
The "REST manager", or the `httd.Client`, handles rate limiting for outgoing requests, and updated the internal logic on responses. All the REST methods are defined on the `disgord.Client` and checks for issues before the request is sent out.

If the request is a standard GET request, the cache is always checked first to reduce delay, network traffic and load on the Discord servers. And on responses, regardless of the http method, the data is copied into the cache.

Some of the REST methods (updating existing data structures) will use the builder+command pattern. While the remaining will take a simple config struct. 

> Note: Methods that update a single field, like SetCurrentUserNick, does not use the builder pattern.
```go
// bypasses local cache
client.GetCurrentUser(disgord.DisableCache)
client.GetGuildMembers(guildID, disgord.DisableCache)

// always checks the local cache first
client.GetCurrentUser()
client.GetGuildMembers(guildID)
```

#### Voice
Whenever you want the bot to join a voice channel, a websocket and UDP connection is established. So if your bot is currently in 5 voice channels, then you have 5 websocket connections and 5 udp connections open to handle the voice traffic.

#### Cache
The cache tries to represent the Discord state as accurate as it can. Because of this, the cache is immutable by default. Meaning the does not allow you to reference any cached objects directly, and every incoming and outgoing data of the cache is deep copied.

### Package structure
None of the sub-packages should be used outside the library. If there exists a requirement for that, please create an issue or pull request.
```Markdown
github.com/andersfylling/disgord
└──.circleci    :CircleCI configuration
└──.githooks    :Hooks that can help speed up development for DisGord contributors
└──.github      :GitHub templates, issues, PR, etc.
└──crs          :Cache Replacement Algorithm
└──cmd          :Private content for live testing
└──constant     :Constants such as version, GitHub URL, etc.
└──docs         :Examples, templates, (documentation)
└──endpoint     :All the REST endpoints of Discord
└──event        :All the Discord event identifiers
└──generate     :go:generate logic
└──httd         :Deals with rate limits and http calls
└──logger       :Logger interface and Zap wrapper
└──ratelimit    :All the ratelimit keys for the REST endpoints
└──std          :Standard implementations/functionality that bot developers can use with DisGord logic
└──testdata     :Holds all test data for unit tests (typically JSON files)
└──websocket    :Discord Websocket logic (reconnect, resume, etc.)
```
The root pkg (disgord) holds all the data structures and the main client. Essentially all the features that should be used by the developer for creating bots. If you need access to the Snowflake type used by DisGord, then you should use `github.com/andersfylling/snowflake`.

### Dependencies
```Markdown
github.com/andersfylling/disgord
└──github.com/andersfylling/snowflake  :The snowflake ID designed for Discord
└──github.com/json-iterator/go         :For faster JSON decoding/encoding
└──github.com/gorilla/websocket        :Default websocket client
└──github.com/sergi/go-diff            :Unit testing for checking JSON encoding/decoding of structs
└──github.com/uber-go/zap              :Logging (optional)
└──go.uber.org/atomic                  
```


### Logging
DisGord requires you to inject a logger instance if you want DisGord to log internal messages (recommended). Logrus is supported out of the box, while other projects might require you to wrap them to comply with `disgord.Logger`. 

You can also use the default logger, which is a wrapped Zap instance: `disgord.DefaultLogger(false)` (see logging.go).

To use the log instance later on, call the method `Session.Logger()`.
```go
client := disgord.New(&disgord.Config{
	Logger: disgord.DefaultLogger(false),
})
client.Logger().Error("oh no")
```

### Build tags
> For **advanced users only**.

If you do not wish to use json-iterator, you can pass `-tags=json-std` to switch to `"encoding/json"`.
However, json-iterator is the recommended default for this library.

DisGord has the option to use mutexes (sync.RWMutex) on Discord objects. By default, methods of Discord objects are not locked as this is
not needed in our event driven architecture unless you create a parallel computing environment.
If you want the internal methods to deal with read-write locks on their own, you can pass `-tags=disgord_parallelism`, which will activate built-in locking.
Making all methods thread safe.

If you want to remove the extra memory used by mutexes, or you just want to completely avoid potential deadlocks by disabling
mutexes you can pass `-tags=disgord_removeDiscordMutex` which will replace the RWMutex with an empty struct, causing mutexes (in Discord objects only) to be removed at compile time.
This cannot be the default behaviour, as it creates confusion whether or not a mutex exists and leads to more error prone code. The developer has to be aware themselves whether or not
their code can be run without the need of mutexes. This option is not affected by the `disgord_parallelism` tag.


## Contributing
Please see the [CONTRIBUTING.md file](CONTRIBUTING.md) (Note that it can be useful to read this regardless if you have the time)

## Q&A
> **NOTE:** To see more examples go visit the docs/examples folder.
See the GoDoc for a in-depth introduction on the various topics (or disgord.go package comment). Below is an example of the traditional ping-pong bot and then some.

```Markdown
1. How do I find my bot token and/or add my bot to a server?

Tutorial here: https://github.com/andersfylling/disgord/wiki/Get-bot-token-and-add-it-to-a-server
```

```Markdown
2. Is there an alternative Go package?

Yes, it's called DiscordGo (https://github.com/bwmarrin/discordgo). Its purpose is to provide low 
level bindings for Discord, while DisGord wants to provide a more configurable system with more 
features (channels, build constraints, tailored unmarshal methods, etc.). 
Currently I do not have a comparison chart of DisGord and DiscordGo. But I do want to create one in the 
future, for now the biggest difference is that DisGord does not support self bots.
```

```Markdown
3. Why make another Discord lib in GoLang?

I'm trying to take over the world and then become a intergalactic war lord. Have to start somewhere.
```

```Markdown
4. Does this project re-use any code from DiscordGo?

Yes. See guild.go. The permission consts are pretty much a copy from DiscordGo.
```

```Markdown
5. Will DisGord support self bots?

No. Self bots are againts ToS and could result in account termination (see
https://support.discordapp.com/hc/en-us/articles/115002192352-Automated-user-accounts-self-bots-). 
In addition, self bots aren't a part of the official Discord API, meaning support could change at any 
time and DisGord could break unexpectedly if this feature were to be added.
```

