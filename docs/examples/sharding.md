> NOTE! outdated: sharding is handled automatically now. Use disgord.ShardConfig.

DisGord supports the use of sharding for as explained here: [discordapp.com/.../gateway#sharding](https://discordapp.com/developers/docs/topics/gateway#sharding)

DisGord uses an internal shard manager to handle this for you. However, you can customize this and should have enough control to handle sharding across N instances of disgord (see the godoc).

# Enforce N number of shards
```go
client := disgord.New(disgord.Config{
    ShardConfig: disgord.ShardConfig{
        ShardIDs: []uint{0, 1, 2, 3}, // must be valid shard ids
    },
    BotToken: "secret token",
    Logger: disgord.DefaultLogger(false), // optional logging, debug=false
})
```

The entire shard config is optional as your bot will always use sharding by default, and automatically decide on how many shards you need. There is also no need to communicate with individual shards in this design, so there is no difference when you interact with the DisGord interface regardless of how many shards are being used.

```go
// this client is also using shards
client := disgord.New(disgord.Config{
    BotToken: "secret token",
    Logger: disgord.DefaultLogger(false), // optional logging, debug=false
})
```
