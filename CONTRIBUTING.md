# Contributing

The following is a set of guidelines for contributing to DisGord.
These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

> This was inspired by the [Atom Organization](https://github.com/atom) GitHub guideline.

#### Table Of Contents

[Code of Conduct](#code-of-conduct)

[I don't want to read this whole thing, I just have a question!!!](#i-dont-want-to-read-this-whole-thing-i-just-have-a-question)

[What should I know before I get started?](#what-should-i-know-before-i-get-started)
  * [Introduction to DisGord](#introduction)
  * [Design Decisions](#design-decisions)
  * [Running Unit Tests](#running-unit-tests)

[How Can I Contribute?](#how-can-i-contribute)
  * [Reporting Bugs](#reporting-bugs)
  * [Suggesting Enhancements](#suggesting-enhancements)
  * [Your First Code Contribution](#your-first-code-contribution)
  * [Unit Tests](#unit-tests)
  * [Pull Requests](#pull-requests)

[Styleguide](#styleguide)
  * [Git Commit Messages](#git-commit-messages)

[Additional Notes](#additional-notes)

## Code of Conduct
Use the GoLang formatter tool. Regarding commenting on REST functionality, do follow the commenting guide, which should include:
 1. The date it was reviewed/created/updated (just write reviewed)
 2. A description of what the endpoint does (copy paste from the discord docs if possible)
 3. The complete endpoint
 4. The rate limiter
 5. optionally, comments. (Comment#1, Comment#2, etc.)

Example (use spaces):
```go
// CreateGuild [POST]       Create a new guild. Returns a guild object on success. Fires a Guild Create
//                          Gateway event.
// Endpoint                 /guilds
// Rate limiter             /guilds
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#create-guild
// Reviewed                 2018-08-16
// Comment                  This endpoint can be used only by bots in less than 10 guilds. Creating channel
//                          categories from this endpoint is not supported.
```

#### Mutex
For discord specific structs (Message, User, etc.) use Lockable. This is to allow deactivating/activating the mutex in public methods.

If the mutex does not need to be publicly accessible, then use the `mu` prefix.

#### Singletons
Nope. But discussions are welcome.

## I don't want to read this whole thing I just have a question!!!

> **Note:** While you are free to ask questions, given that you add the [help] prefix. You'll get faster results by using the resources below.

You can find a support channel for DisGord in Discord. We exist in both the Gopher server and the Discord API server:
 - [Discord Gophers](https://discord.gg/qBVmnq9)
 - [Discord API](https://discord.gg/HBTHbme)

Using the live chat application will most likely give you a result quicker.

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
Compared to DiscordGo, DisGord does not focus on having a minimalistic implementation. DisGord hopes to simplify development and give developers a very configurable system. The goal is to support everything that DiscordGo does, and on top of that; helper functions, methods, cache replacement algorithms, event channels, etc.

### Design Decisions
DisGord should handle events, REST, voice, caching; these can be split into separate logical parts. Because of this DisGord must have an event driven architecture to support events and voice. REST methods should be written idiomatic, reusing code for readability is acceptable: I want these methods to stay flexible for future changes, and there might be requirements to directly change the json data. Lastly, caching should be done behind the scenes. Any REST calls, and incoming events should go through the cache before the dev/user gets access to the data.

#### Caching
Caching in discord is happens behind the scenes and unless a user explicitly tells disgord to not check the cache before sending a request, disgord should always do so.

The cache is split into repositories; users, channels, guilds, presences. This is to reduce wait time from locking and complex look ups when fetching data. If necessary, more repositories can be created if reports states the need.

Each repository has their own demultiplexer for events and each handler is a unique method in the naming scheme of `onEventName(data []byte) (updated interface{}, err error)`. The handlers should do copying on their own to support immutability, this is to stop interacting with the object in cache as soon as possible such that new events can be handled. Each handler should also be able to run in parallel with other repositories for the same event and handle flags that lets the developer decide if heap allocations are necessary.

Events should always be read from cache, not the actual raw data from discord. This is to ensure that the user has complete objects which gives them a better experience.

#### Event Handlers (functions and channels)
> Also known as listeners/callbacks, but are named handlers to stay close to the reactor pattern naming conventions.

> Handlers are both functions and channels in DisGord.

DisGord gives the option to register multiple handlers per event type. But will not run handlers in parallel. All handlers are run in sequence.

> Note! The handlers run in sequence per event. But events execute concurrently.

```go
session.On(event.MessageCreate, func(session disgord.Session, evt *disgord.MessageCreate) {
    // ...
})
```

```go
session.On(event.MessageCreate, messageChan)
```

#### Event Middlewares
Middlewares are executed in sequence and can manipulate the event content directly. Once all have executed, and none returns nil, the handler(s) are executed. Middlewares only applies on a par-registration basis, meaning they only apply to the handlers/middlewares that are arguments in the same registration as them.

```go
session.On(event.MessageCreate, middleware1, middleware2, messageChan)
```

It's an alternative way of doing fail-fast are specifying requirements that can be reused. Or directly manipulate the incoming events before a handler process it, to ensure certain values are added/specified.

#### Event Handlers lifetime
DisGord allows a controller that dictates the lifetime of the handlers. Such that the handler(s) run only once, five time, or only within five minutes, or whatever kind of behaviour is desired. These are optional and are injected at the end of the registration function.

```go
session.On(event.MessageCreate, messageChan, &disgord.Ctrl{Deadline:5*time.Second})
```

#### Mutex
Every Discord object will hold a read/write mutex. You can also disabled the mutexes by build constraints. See the main README file.

#### Go Generate

If you during your contribution have changed either `events.go` or `event/events.go`, you must run `go generate` in the root folder of the project before pushing.  
This command will ensure that all generated files have been updated accordingly.  
If this command gives you warnings, you must correct them before pushing.

All files written by Go Generate will be suffixed with `_gen.go`, they should **NOT** be edited manually as they will be overwritten by Go Generate.  
Instead, edit the templates in `generate/` or the files they're based on (see previous paragraph). 

### Running Unit Tests
> WARNING! Please do not run the unit tests for the endpoints as these are verified to work before being pushed, and rechecking every time is just spamming the Discord API for useless information.

In DisGord you will see both local unit tests and unit tests that verify directly against the Discord API. Note that the "live" tests (that depends on the Discord API) needs to be activated through environment variables. However, these should only be run when some breaking changes to the REST implementation changes. You will most likely never have to worry about it.

But if you want to properly test all the implementations you need to provide a bot token under the environment variable: "DISGORD_TEST_BOT". Any integration tests depending on this token is skipped whenever it is missing. The following environment variables must exist in order to properly execute a complete integration test (see constant package for information):
 1. DISGORD_TEST_BOT
 2. DISGORD_TEST_GUILD_ADMIN
 3. DISGORD_TEST_GUILD_DEFAULT (must have one custom emoji)
 3. DISGORD_TEST_GUILD_DEFAULT_EMOJI_SNOWFLAKE (snowflake id of emoji)


Editing the event handlers or rest package requires that you test with a bot token to verify success. (Note that tests aren't complete and is considered a work in progress).

For the local tests (the main tests) DisGord tries to decouple.. well.. everything. The websocket connection is decoupled to allow mocking input/output socket communication with the Discord API. And the REST methods are decoupled by `Getter`, `Poster`, etc. interfaces. However, if anyone creates a solution for decoupling the `http/Client` instead of the `disgord/httd/Client` that would be a great improvement, as we could do integration tests of the entire `httd` package as well, to verify more accurate DisGord behavior.

## How Can I Contribute?

### Reporting Bugs
Reporting a bug should help the community improving DisGord. We need you to be specific and give enough information such that it can be reproduced by others. You must use the Bug template which can be found here: [TEMPLATE_BUG.md](docs/TEMPLATE_BUG.md).

### Suggesting Enhancements
We don't currently have a template for this. Provide benchmarks or demonstrations why your suggestion is an improvement or how it can help benefit this project is of great appreciation.

### Your First Code Contribution
Remember to run go fmt to properly format your code and add unit tests if you provide new content. Benchmarks are also welcome as we can use these in future decisions (!).

### Unit Tests
Make them readable. And don't have external dependencies (eg. Discord API).

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
