# Disgord

[![CircleCI](https://circleci.com/gh/andersfylling/disgord/tree/master.svg?style=shield)](https://circleci.com/gh/andersfylling/disgord/tree/master)[![Maintainability](https://api.codeclimate.com/v1/badges/687d02ca069eba704af9/maintainability)](https://codeclimate.com/github/andersfylling/disgord/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/687d02ca069eba704af9/test_coverage)](https://codeclimate.com/github/andersfylling/disgord/test_coverage) [![Coverage Status](https://coveralls.io/repos/github/andersfylling/disgord/badge.svg)](https://coveralls.io/github/andersfylling/disgord)

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

    dg, err := disgord.NewDisgord(&disgord.Config{
        Token: os.Getenv("DISGORD_TOKEN"),
        //Debug: true,
    })
    if err != nil {
        panic(err)
    }

    err = dg.Connect()
    if err != nil {
        panic(err)
        return
    }
    <-termSignal
    dg.Disconnect()
}
```

Gives a output similar to (note that it only gives out data to the terminal atm):

```go
INFO[2018-01-25 04:14:15] Connecting to discord Gateway                 lib="Disgord v0.0.0"
INFO[2018-01-25 04:14:17] Connected                                     lib="Disgord v0.0.0"
INFO[2018-01-25 04:14:17] Event{READY}
{"v":6,"user_settings":{},"user":{"verified":true,"username":"disgord" ......
```

And the disconnect method provides graceful shutdown (buggy on reconnect periodes):

```go

^C
INFO[2018-01-25 04:30:09] Closing Discord gateway connection            lib="Disgord v0.0.0"
INFO[2018-01-25 04:30:11] Disconnected                                  lib="Disgord v0.0.0"
```

TODO:

Event dispatchers + caching(ish):

- [x] Ready
- [ ] Resumed
- [ ] ChannelCreate
- [ ] ChannelUpdate
- [ ] ChannelDelete
- [ ] ChannelPinsUpdate
- [x] GuildCreate
- [x] GuildUpdate
- [x] GuildDelete
- [ ] GuildBanAdd
- [ ] GuildBanRemove
- [ ] GuildEmojisUpdate
- [ ] GuildIntegrationsUpdate
- [ ] GuildMemberAdd
- [ ] GuildMemberRemove
- [ ] GuildMemberUpdate
- [ ] GuildMemberChunk
- [ ] GuildRoleCreate
- [ ] GuildRoleUpdate
- [ ] GuildRoleDelete
- [x] MessageCreate
- [x] MessageUpdate
- [x] MessageDelete
- [ ] MessageDeleteBulk
- [ ] MessageReactionAdd
- [ ] MessageReactionRemove
- [ ] MessageReactionRemoveAll
- [ ] PresenceUpdate
- [ ] TypingStart
- [x] UserUpdate
- [ ] VoiceStateUpdate
- [ ] VoiceServerUpdate
- [ ] WebhooksUpdate
