# Contributing

The following is a set of guidelines for contributing to Disgord.
These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

> This was inspired by the [Atom Organization](https://github.com/atom) GitHub guideline.

#### Table Of Contents

[Code of Conduct](#code-of-conduct)

[I don't want to read this whole thing, I just have a question!!!](#i-dont-want-to-read-this-whole-thing-i-just-have-a-question)

[Development environment](#development-environment)

[What should I know before I get started?](#what-should-i-know-before-i-get-started)
  * [Introduction to Disgord](#introduction)
  * [Design Decisions](#design-decisions)
  * [Running Unit Tests](#running-unit-tests)

[How Can I Contribute?](#how-can-i-contribute)
  * [Reporting Bugs](#reporting-bugs)
  * [Suggesting Enhancements](#suggesting-enhancements)
  * [Your First Code Contribution](#your-first-code-contribution)
  * [Tests](#tests)
  * [Pull Requests](#pull-requests)

[Styleguide](#styleguide)
  * [Git Commit Messages](#git-commit-messages)

[Additional Notes](#additional-notes)

## Code of Conduct
Use the GoLang formatter tool. Regarding commenting on REST functionality, do follow the commenting guide, which should include:
 1. The date it was reviewed/created/updated (just write reviewed)
 2. A description of what the endpoint does (copy paste from the discord docs if possible)
 3. The complete endpoint
 4. optionally, comments. (Comment#1, Comment#2, etc.)

Example (use spaces):
```go
// CreateGuild [POST]       Create a new guild. Returns a guild object on success. Fires a Guild Create
//                          Gateway event.
// Endpoint                 /guilds
// Discord documentation    https://discord.com/developers/docs/resources/guild#create-guild
// Reviewed                 2018-08-16
// Comment                  This endpoint can be used only by bots in less than 10 guilds. Creating channel
//                          categories from this endpoint is not supported.
```

## I don't want to read this whole thing I just have a question!!!

> **Note:** While you are free to ask questions, given that you add the [help] prefix. You'll get faster results by using the resources below.

You can find a support channel for Disgord in Discord. We exist in both the Gopher server and the Discord API server:
 - [Discord Gophers](https://discord.gg/qBVmnq9)
 - [Discord API](https://discord.gg/HBTHbme)

Using the live chat application will most likely give you a result quicker.

## Development environment

I personally use pre-commit, a python script that triggers a bunch of go tools on `git commit`; [https://pre-commit.com/](https://pre-commit.com/).

The project is already configured so you can simply run `pre-commit install`.

### testing with a live bot
Generally you simply need to create a fork and clone it to start development. However, sometimes you might want to test the changes you've made to disgord with an actual bot.

The following text assumes the directory layout:
```go
+-+dev
|
+-+disgord (your fork)
|
+-+bot (some bot project)
  +-go.mod
  +-main.go
```

Here we see that your disgord folder and bot folder lays under the parent directory "dev". The namings are not important, it is just to illustrate the relative path between "bot" and "disgord".

your `bot/go.mod` file will include at least disgord as a dependency. However, this links to the github repository found online. To change this to your local fork of disgord, a "replace directive" can be added to the `bot/go.mod` file.

Before:
```go
module github.com/myusername/bot

go 1.15

require (
	github.com/andersfylling/disgord v0.24.0
)
```

After:
```go
module github.com/myusername/bot

go 1.15

require (
	github.com/andersfylling/disgord v0.24.0
)

replace github.com/andersfylling/disgord => ../disgord
```

We have now stated that the `github.com/andersfylling/disgord` dependency can be found at the relative path `../disgord`. Such that any time you make changes to your disgord folder and run/build your bot, it will fetch the disgord code from `../disgord`. This allows you to test your changes immediately.


## What Should I Know Before I Get Started?
Depending on what you want to contribute to, here's a few:
 * Event driven architecture
 * Reactor pattern
 * build constraints
 * golang modules
 * discord docs
 * concurrency and channels
 * sharding
 * caching

### Introduction
Compared to DiscordGo, Disgord does not focus on having a minimalistic implementation that should represent the discord docs. Mostly because this isn't possible (eg. setting default values in REST requests, you'll have to do something hacky to get away with that or use the builder pattern). Disgord hopes to simplify development and give developers a very configurable system. The goal is to support everything that DiscordGo does, and ontop of that; helper functions, methods, event channels, etc.


## Architecture & Behavior
Discord provide communication in different forms. Disgord tackles the main ones, events (ws), voice (udp + ws), and REST calls.

You can think of Disgord as layered, in which case it will look something like:
![Simple way to think about Disgord architecture from a layered perspective](docs/disgord-layered-version.png)

### Design Decisions
Disgord should handle events, REST, voice, caching; these can be split into separate logical parts. Because of this Disgord must have an event driven architecture to support events and voice. Caching should be done behind the scenes. 

#### Code flow / design
Prefer procedural when possible. Note that disgord.Snowflake and disgord.Time, should be treated as OOP. Especially their .IsZero() implementation to avoid any potential zero checking.

You may see OOP code that could easily be procedural, feel to rectify this. But remember that this, is after all, Go.

#### JSON encoding
For now, Disgord will utilise JSON for Discord communication. ETF is on hold. For changing the unmarshal/marshal logic, see the disgord/json pkg.

#### REST requests
All GET REST calls, and incoming events should go through the cache before the dev/user gets access to the data. Also, the calls will most likely utilise a dedicated data structures with a "Params" suffix.

All REST methods are resource based. Meaning you will see Channel(id).Update()..., Guild(id).Member(id).Kick(). This is horrible to mock, which is why you should inject your own http.RoundTripper if you have to test anything.

> see ClientQueryBuilder

context.Context is optional, and can injected into every resource using `.WithContext(context.Context)`. Note that the context.Context is only relevant for the depth you inject it. It does not continue to the next level. And the state is not mutable to avoid confusions.

```go
guildResource := client.Guild(id)
guildResourceWithCtx := guildResource.WithContext(ctx) // returns a copy with ctx
memberResource := guildResourceWithCtx.Member(uid) // does not contain the ctx
memberResourceWithCtx := memberResource.WithContext(ctx) // now contains the ctx
```

#### Cache
In Disgord, the caching layer is both a creational layer and a caching layer. Meaning your cache layer is responsible to initialising the incoming data structures from events and GET requests.

There exists a Nop implementation - that still holds logic - but since it does not alter a local state/cache, it's given the "Nop" suffix.

When creating a custom cache, remember to embed disgord.CacheNop to avoid having to implement methods you don't care about / don't need.

#### Interface implementation
Empty implementations, to e.g. satisfy some interface, should use "Nop" in it's name to signify so. 


#### Events
 
##### Handlers (functions and channels)
> Also known as listeners/callbacks, but are named handlers to stay close to the reactor pattern naming scheme.

> Handlers are both functions and channels in Disgord.

When you register a handler, you create a handler-specification. This isolated unit contains the event, N middlewares and M handlers, all of which are executed sequentially by default. However, handler-specifications are run concurrently.

```go
client.On(disgord.EvtMessageCreate, func(session disgord.Session, evt *disgord.MessageCreate) {
    // ...
})
```

```go
messageChan := make(chan *disgord.MessageCreate)
client.On(event.MessageCreate, messageChan)
```

##### Event Middlewares
Middlewares are executed in sequence and can manipulate the event content directly. Once all have executed, and none returns nil, the handler(s) are executed. Middlewares only applies on a par-registration basis, meaning they only apply to the handlers/middlewares that are arguments in the same registration as them.

```go
client.On(event.MessageCreate, middleware1, middleware2, messageChan)
```

It's an alternative way of doing fail-fast are specifying requirements that can be reused. Or directly manipulate the incoming events before a handler process it, to ensure certain values are added/specified.

##### Handler-specification lifetime
Disgord allows a controller that dictates the lifetime of the handler-specification. Such that the handler(s) run only once, five time, or only within five minutes, or whatever kind of behaviour. These are optional and are injected at the end of the registration function.

```go
client.On(event.MessageCreate, messageChan, &disgord.Ctrl{Deadline:5*time.Second})
```

##### Registration with compile time polymorphism
Until there are generics, we will have to use the builder pattern. Call the .Event() method from session or client to gain access to compile time constrained handler-specification registration.

```go
client.Event().GuildUpdate(func(s Session, evt *GuildUpdate) {
})
client.Event().WithMdlw(excludeBots).MessageCreate(func(s Session, evt *MessageCreate) {
})
client.Event().WithCtrl(&Ctrl{Runs: 3}).MessageCreate(func(s Session, evt *MessageCreate) {
})
client.Event().WithCtrl(&Ctrl{Runs: 3}).WithMdlw(excludeBots).MessageCreate(func(s Session, evt *MessageCreate) {
})
```

#### Go Generate

Please run `go generate` before every commit. I recommend using a git hook.

All files written by Go Generate will be suffixed with `_gen.go`, they should **NOT** be edited manually as they will be overwritten by Go Generate.
Instead, edit the templates in `generate/` or the files they're based on (see previous paragraph). 

### Unit Tests

`go test ./... -race`

Do not overdo it. Hardening something that interacts with the Discord API is forced to be rewritten, and possibly tests deleted as discord can suddenly change something. Public functions should however be well tested.

In Disgord you will see both local unit tests and unit tests that verify directly against the Discord API. Note that the "live" tests (that depends on the Discord API) needs to be activated through environment variables. However, these should only be run when some breaking changes to the REST implementation changes. You will most likely never have to worry about it.

But if you want to properly test all the implementations you need to provide a bot token under the environment variable: "DISGORD_TEST_BOT". Any integration tests depending on this token is skipped whenever it is missing. The following environment variables must exist in order to properly execute a complete integration test (see constant package for information):
 1. DISGORD_TEST_BOT
 2. DISGORD_TEST_GUILD_ADMIN
 3. DISGORD_TEST_GUILD_DEFAULT (must have one custom emoji)
 3. DISGORD_TEST_GUILD_DEFAULT_EMOJI_SNOWFLAKE (snowflake id of emoji)


Editing the event handlers or rest package requires that you test with a bot token to verify success. (Note that tests aren't complete and is considered a work in progress).

For the local tests (the main tests) Disgord tries to decouple.. well.. everything. The websocket connection is decoupled to allow mocking input/output socket communication with the Discord API. And the REST methods are decoupled by `Getter`, `Poster`, etc. interfaces. However, if anyone creates a solution for decoupling the `http/Client` instead of the `disgord/httd/Client` that would be a great improvement, as we could do integration tests of the entire `httd` package as well, to verify more accurate Disgord behavior.

## How Can I Contribute?

### Reporting Bugs
Reporting a bug should help the community improving Disgord. We need you to be specific and give enough information such that it can be reproduced by others. You must use the Bug template which can be found here: [TEMPLATE_BUG.md](docs/TEMPLATE_BUG.md).

### Suggesting Enhancements
We don't currently have a template for this. Provide benchmarks or demonstrations why your suggestion is an improvement or how it can help benefit this project is of great appreciation.

### Your First Code Contribution
Remember to run go fmt to properly format your code and add unit tests if you provide new content. Benchmarks are also welcome as we can use these in future decisions (!).

### Tests
Make them readable. Tests that is for the public interface of Disgord, should be placed in a file with _integration_test suffix (also remember the integration build tag!).

### Pull Requests
If your PR is not ready yet, make it a Draft.

Deadlines:
 - If you do not fix the required changes within 30 days your PR will be closed. 
 - If you have created a PR before that was closed due to rule #1, the deadline is reduced from 30 days to 20 days (this only applies when you have 2 or more PR in a row that are a victim to rule #1).
 - Rule #1 and Rule #2 does not apply if you mark your PR as a draft. Such that if you forget about it for 29 days and mark it as a draft, it will not be closed on day 30 and the counter is reset.
 - PR drafts have no deadline.

## Styleguide
`go fmt ./...`

### Git Commit Messages
* Use the present tense ("Adds feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line
* When only changing documentation, include `[ci skip]` in the commit title

## Additional notes
Add the following prefixes to your issues to help categorize them better:
* [help] when asking a question about functionality.
* [discussion] when starting a discussion about wanted or existing functionality
* [proposal] when you have written a comprehensive suggestion with examples
