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
    // should now be connected

    // add a event listener
    dg.ReadyEvent.Add(func(session disgordctx.Context, box *event.ReadyBox) {
      // fmt.Printf("\n----\n:%s:\n%+v\n-------\n", event.ReadyKey, box)
    })

    // add a event listener using the abstract method
    dg.OnEvent(event.ReadyKey, func(session disgordctx.Context, box *event.ReadyBox) {
      fmt.Printf("\n----\n:%s:\n%+v\n-------\n", event.ReadyKey, box)
      // Tip: for blocking/long running tasks check box.Ctx.Done()
    })


    <-termSignal
    dg.Disconnect()
}
```

Output:

```
╰─ go build && ./disgordtest
INFO[2018-02-16 19:05:47] Connecting to discord Gateway                 lib="Disgord v0.0.0"
INFO[2018-02-16 19:05:48] Connected                                     lib="Disgord v0.0.0"

----
:READY:
&{APIVersion:6 User:disgord#2355{40472951282397185} Guilds:[0xc4203922d0] SessionID:4dc1bab8ff8fgfg234f7e0997d7d28a Trace:[gateway-prd-main-1t4gc discord-sessions-prd-2-6] RWMutex:{w:{state:0 sema:0} writerSem:0 readerSem:0 readerCount:0 readerWait:0}}
-------
```

Then on a system interrupt, here pressing `Ctrl+C`, you will see the following:

```
^C
INFO[2018-02-16 19:07:28] Closing Discord gateway connection            lib="Disgord v0.0.0"
INFO[2018-02-16 19:07:30] Disconnected                                  lib="Disgord v0.0.0"
```

## Progress

[Progression for different Discord implementations](PROGRESS.md)
