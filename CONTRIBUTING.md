# Contributing

The following is a set of guidelines for contributing to Disgord.
These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

> **Note:** This CONTRIBUTIONS guideline is heavily inspired by the one created by the [Atom Organization](https://github.com/atom) on GitHub.

#### Table Of Contents

[Code of Conduct](#code-of-conduct)

[I don't want to read this whole thing, I just have a question!!!](#i-dont-want-to-read-this-whole-thing-i-just-have-a-question)

[What should I know before I get started?](#what-should-i-know-before-i-get-started)
  * [Introduction to Disgord](#introduction)
  * [Design Decisions](#design-decisions)
  * [Running Unit Tests](#running-unit-tests)

[How Can I Contribute?](#how-can-i-contribute)
  * [Reporting Bugs](#reporting-bugs)
  * [Suggesting Enhancements](#suggesting-enhancements)
  * [Your First Code Contribution](#your-first-code-contribution)
  * [Unit Tests](#unit-tests)
  * [Pull Requests](#pull-requests)

[Styleguides](#styleguides)
  * [Git Commit Messages](#git-commit-messages)

[Additional Notes](#additional-notes)
  * [Issue and Pull Request Labels](#issue-and-pull-request-labels)

## Code of Conduct
Use the GoLang formatter tool. Regarding commenting on REST functionality, do follow the commenting guide, which should include:
 1. The date it was reviewed/created/updated (just write reviewed)
 2. A description of what the endpoint does (copy paste from the discord docs if possible)
 3. The complete endpoint
 4. The rate limiter
 5. optionally, comments. (Comment#1, Comment#2, etc.)

Example (use spaces):
```GoLang
// CreateGuild [POST]       Create a new guild. Returns a guild object on success. Fires a Guild Create
//                          Gateway event.
// Endpoint                 /guilds
// Rate limiter             /guilds
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#create-guild
// Reviewed                 2018-08-16
// Comment                  This endpoint can be used only by bots in less than 10 guilds. Creating channel
//                          categories from this endpoint is not supported.
```

#### Functions and param interfaces
When creating a function, make sure that it has one purpose. It shouldn't hold any hidden functionality, and it could be wise to take dependencies as a parameter:
```GoLang
func CreateGuild(client httd.Poster, params *CreateGuildParams) (ret *Guild, err error) {
    ...
}
```

When you use a dependency as a parameter, either use an existing interface which holds all the methods you require. Or create a new interface for the functions needs (again, keep it purposeful such that other can reuse the interface).

#### Mutex
Struct's that utlises mutex must handle locking in their public methods. However, avoid locking in private methods such that they can be reused. eg. by public methods. The mutex should also be embedded (public accessible), such that developers can lock the object by their own needs.

#### Singletons
I won't accept pull requests where the author has created a singleton structure. I do not want package singletons either, as I'm worried it might cause technical debt. If you disagree you are welcome to create a discussion (not about the pattern, but why your implementation requires a singleton).

## I don't want to read this whole thing I just have a question!!!

> **Note:** While you are free to ask questions, given that you add the [help] prefix. You'll get faster results by using the resources below.

This repository has it's own discord guild/server: [Discord Gophers#Disgord](https://discord.gg/qBVmnq9)
Using the live chat application will most likely give you a faster result.

## What Should I Know Before I Get Started?
Technologies:
 * Event driven architecture
 * Reactor pattern
 * build constraints
 * golang modules
 * layout of discord docs
 * concurrency and channels

### Introduction
Compared to DiscordGo, Disgord does not focus on having the minimal implementation. Disgord hopes to simplify development and create a very configurable system. The goal is to support everything that DiscordGo does, and ontop of that; helper functions, methods, cache replacement algorithms, event channels, etc.

### Design Decisions
Utilises the Reactor pattern for handling incoming events.

#### Handlers
> Also known as listeners/callbacks. The event driven architecture of Disgord is uses the react pattern and as such the handlers are triggered in sequence.

Disgord gives the option to register multiple handlers per event type. But will not run handlers in parallel. All handlers are run in sequence and that will not change.

```GoLang
Session.On(event.MessageCreate, func(session disgord.Session, evt *disgord.MessageCreate) {
    // ...
})
```

#### Channels
> An alternative way to listen for events

While handlers are the common approach to handle events, Disgord also support channels. It is important to mark that having multiple observers on one event type channel will cause only one of them to receive the event (this is just how channels work).

#### Mutex

Every Discord object will hold a read/write mutex. This might be removed in the future. You can also disabled the mutexes by build constraints. See the main README file.


#### Go Generate

If you during your contribution have changed either `events.go` or `event/events.go`, you must run `go generate` in the root folder of the project before pushing.  
This command will ensure that all generated files have been updated accordingly.  
If this command gives you warnings, you must correct them before pushing.

All files written by Go Generate will be suffixed with `_gen.go`, they should **NOT** be edited manually as they will be overwritten by Go Generate.  
Instead, edit the templates in `generate/` or the files they're based on (see previous paragraph). 

### Running Unit Tests
> WARNING! Please do not run the unit tests for the endpoints as these are verified to work before being pushed, and rechecking every time is just spamming the Discord API for useless information.

You can run unit tests without the need of a bot token. However, if you want to properly test all the implementations you need to provide a bot token under the environment variable: "DISGORD_TEST_BOT". Any integration tests depending on this token is skipped whenever it is missing. The following environment variables must exist in order to properly execute a complete integration test (see constant package for information):
 1. DISGORD_TEST_BOT
 2. DISGORD_TEST_GUILD_ADMIN
 3. DISGORD_TEST_GUILD_DEFAULT (must have one custom emoji)
 3. DISGORD_TEST_GUILD_DEFAULT_EMOJI_SNOWFLAKE (snowflake id of emoji)


Editing the event handlers or rest package requires that you test with a bot token to verify success. (Note that tests aren't complete and is considered a work in progress).

> **Note:** The module DisgordWS (dependency) utilises a different environment variable for the bot token, so tests regarding integration testing there are skipped. There should be no need to test the DisgordWS module while developing on the Disgord module.

## How Can I Contribute?

### Reporting Bugs
Reporting a bug should help the community improving Disgord. We need you to be specific and give enough information such that it can be reproduced by others. You must use the Bug template which can be found here: [TEMPLATE_BUG.md](TEMPLATE_BUG.md).

### Suggesting Enhancements
We don't currently have a template for this. Provide benchmarks or demonstrations why your suggestion is an improvement or how it can help benefit this project is of great appreciation.

### Your First Code Contribution

### Unit Tests

### Pull Requests


## Styleguides

### Git Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line
* When only changing documentation, include `[ci skip]` in the commit title
* Consider starting the commit message with an applicable emoji:
    * :art: `:art:` when improving the format/structure of the code
    * :racehorse: `:racehorse:` when improving performance
    * :non-potable_water: `:non-potable_water:` when plugging memory leaks
    * :memo: `:memo:` when writing docs
    * :penguin: `:penguin:` when fixing something on Linux
    * :apple: `:apple:` when fixing something on macOS
    * :checkered_flag: `:checkered_flag:` when fixing something on Windows
    * :bug: `:bug:` when fixing a bug
    * :fire: `:fire:` when removing code or files
    * :green_heart: `:green_heart:` when fixing the CI build
    * :white_check_mark: `:white_check_mark:` when adding tests
    * :lock: `:lock:` when dealing with security
    * :arrow_up: `:arrow_up:` when upgrading dependencies
    * :arrow_down: `:arrow_down:` when downgrading dependencies
    * :shirt: `:shirt:` when removing linter warnings


## Additional notes

Add the following prefixes to your issues to help categorize them better:
* [help] when asking a question about functionality.
* [discussion] when starting a discussion about wanted or existing functionality
