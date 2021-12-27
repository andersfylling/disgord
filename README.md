<div align='center'>
  <img src="/.github/disgord-draft-8.jpeg" alt='Build Status' />
  <p>
    <a href="https://codecov.io/gh/andersfylling/disgord">
      <img src="https://codecov.io/gh/andersfylling/disgord/branch/develop/graph/badge.svg" />
    </a>
    <a href='https://goreportcard.com/report/github.com/andersfylling/disgord'>
      <img src='https://goreportcard.com/badge/github.com/andersfylling/disgord' alt='Code coverage' />
    </a>
    <a href='https://pkg.go.dev/github.com/andersfylling/disgord'>
      <img src="https://pkg.go.dev/badge/andersfylling/disgord" alt="PkgGoDev">
    </a>
  </p>
  <p>
    <a href='https://discord.gg/fQgmBg'>
      <img src='https://img.shields.io/badge/Discord%20Gophers-%23disgord-blue.svg' alt='Discord Gophers' />
    </a>
    <a href='https://discord.gg/HBTHbme'>
      <img src='https://img.shields.io/badge/Discord%20API-%23disgord-blue.svg' alt='Discord API' />
    </a>
  </p>
</div>

## About
Go module with context support that handles some of the difficulties from interacting with Discord's bot interface for you; websocket sharding, auto-scaling of websocket connections, advanced caching (cache replacement strategies to restrict memory usage), helper functions, middlewares and lifetime controllers for event handlers, etc.

This package is intented to be used with the gateway to keep the cache up to date. You should treat data as read only, since they simply represent the discord state. To change the discord state you can use the REST methods and the gateway commands, which will eventually update your local state as well.

If you want a more lightweight experience you can disable/reject events that you do not need or want to keep track of. Be careful as this might break certain things.

## Tips
 - Use disgord.Snowflake, not snowflake.Snowflake.
 - Use disgord.Time, not time.Time when dealing with Discord timestamps.

By default DM capabilities are disabled. If you want to activate these, or some, specify their related intent.
```go
client := disgord.New(disgord.Config{
    DMIntents: disgord.IntentDirectMessages | disgord.IntentDirectMessageReactions | disgord.IntentDirectMessageTyping,
})
```

## Starter guide
> This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) for dealing with dependencies, remember to activate module support in your IDE

> Examples can be found in [examples](examples) and some open source projects Disgord projects in the [wiki](https://pkg.go.dev/github.com/andersfylling/disgord?tab=importedby)

I highly suggest reading the [Discord API documentation](https://discord.com/developers/docs/intro) and the [Disgord go doc](https://pkg.go.dev/github.com/andersfylling/disgord?tab=doc).

Simply use [this github template](https://github.com/andersfylling/disgord-starter) to create your first new bot!


## API / Interface

> In short Disgord uses the builder pattern by respecting resources

The `Client` or `Session` holds are the relevant methods for interacting with Discord. The API is split by resource, such that Guild related information is found in `Client.Guild(guild_id)`, while user related info is found in `Client.User(user_id)`, gateway interaction is found in `Client.Gateway()`, the same for Channel, CurrentUser, Emoji, AuditLog, etc.

Cancellation is supported by calling `.WithContext(context.Context` before the final REST call (.Get(), .Update(), etc.).

### Events

> every event goes through the cache layer!

For Events, Disgord uses the [reactor pattern](https://dzone.com/articles/understanding-reactor-pattern-thread-based-and-eve). This supports both channels and functions. You chose your preference.

### REST
If the request is a standard GET request, the cache is always checked first to reduce delay, network traffic and load on the Discord servers. And on responses, regardless of the http method, the data is copied into the cache.

```go
// bypasses local cache
client.CurrentUser().Get(disgord.IgnoreCache)
client.Guild(guildID).GetMembers(disgord.IgnoreCache)

// always checks the local cache first
client.CurrentUser().Get()
client.Guild(guildID).GetMembers()

// with cancellation
deadline, _ := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
client.CurrentUser().WithContext(deadline).Get()
```

### Voice
Whenever you want the bot to join a voice channel, a websocket and UDP connection is established. So if your bot is currently in 5 voice channels, then you have 5 websocket connections and 5 udp connections open to handle the voice traffic.

### Cache
The cache tries to represent the Discord state as accurate as it can. Because of this, the cache is immutable by default. Meaning the does not allow you to reference any cached objects directly, and every incoming and outgoing data of the cache is deep copied.

## Contributing
> Please see the [CONTRIBUTING.md file](CONTRIBUTING.md) (Note that it can be useful to read this regardless if you have the time)

You can contribute with pull requests, issues, wiki updates and helping out in the discord servers mentioned above.

To notify about bugs or suggesting enhancements, simply create a issue. The more the better. But be detailed enough that it can be reproduced and please provide logs.

To contribute with code, always create an issue before you open a pull request. This allows automating change logs and releases.

Remember to have stringer installed to run go generate:
`go get -u golang.org/x/tools/cmd/stringer`

## Sponsors
> [JetBrains](https://www.jetbrains.com/?from=github.com/andersfylling/disgord)

A Special thanks to the following companies for sponsoring this project!


<div align='left'>
  <a href="https://www.jetbrains.com/?from=github.com/andersfylling/disgord">
    <img src="/.github/jetbrains-variant-4.svg" alt="JetBrains" width="200px" />
  </a>
</div>

#### Software used

<div align='left'>
  <a href="https://www.jetbrains.com/go/?from=github.com/andersfylling/disgord">
    <img src="/.github/icon-goland.svg" alt="GoLand" width="150px" />
  </a>
</div>

## Q&A
> **NOTE:** To see more examples go to the [examples folder](examples). See the GoDoc for a in-depth introduction on the various topics.

```Markdown
1. How do I find my bot token and/or add my bot to a server?

Tutorial here: https://github.com/andersfylling/disgord/wiki/Get-bot-token-and-add-it-to-a-server
```

```Markdown
2. Is there an alternative Go package?

Yes, it's called DiscordGo (https://github.com/bwmarrin/discordgo). Its purpose is to provide a 
minimalistic API wrapper for Discord, it does not handle multiple websocket sharding, scaling, etc. 
behind the scenes such as Disgord does. Currently I do not have a comparison chart of Disgord and 
DiscordGo. But I do want to create one in the future, for now the biggest difference is that 
Disgord does not support self bots.
```

```Markdown
3. Why make another Discord lib in Go?

I'm trying to take over the world and then become a intergalactic war lord. Have to start somewhere.
```

```Markdown
4. Will Disgord support self bots?

No. Self bots are againts ToS and could result in account termination (see
https://support.discord.com/hc/en-us/articles/115002192352-Automated-user-accounts-self-bots-). 
In addition, self bots aren't a part of the official Discord API, meaning support could change at
any time and Disgord could break unexpectedly if this feature were to be added.
```

