Disgord supports the use of sharding for as explained here: [discordapp.com/.../gateway#sharding](https://discordapp.com/developers/docs/topics/gateway#sharding)

Currently there are no shard manager so you have to write a little bit of extra code to get this working. The first way will allow you decide the number of shards on your own, while the second uses the bot gateway to let Discord decide the number of shards.

### Creating N shards
```go
// create a channel to listen for termination signals (graceful shutdown)
termSignal := make(chan os.Signal, 1)
signal.Notify(termSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

config := &disgord.Config{
    Token: os.Getenv("DISGORD_TOKEN"),
}
RESTClient := disgord.NewRESTClient(config)

// does not actually require authentication...
gateway, err := disgord.GetGateway(RESTClient)
if err != nil {
    return
}

// setup 2 shards
shardCount := uint(2)
shards := make([]disgord.Session, shardCount)
for i := uint(0); i < shardCount; i++ {
    shards[i], err = disgord.NewSession(&disgord.Config{
        Token: config.Token,
        ShardID: i,
        TotalShards: shardCount,
        WebsocketURL: gateway.URL,
    })

    // register a listener
    shards[i].On(disgord.EventMessageCreate, func(session disgord.Session, evt *disgord.MessageCreate) {
        fmt.Println("shard[" + session.ShardIDString() + "]" + evt.Message.Content)
    })
}

// connect
for i := uint(0); i < shardCount; i++ {
    shards[i].Connect()
}

// disconnect
<-termSignal
for i := uint(0); i < shardCount; i++ {
    shards[i].Disconnect()
}
```

### Letting Discord decide the number of shards (recommended)
```go
// create a channel to listen for termination signals (graceful shutdown)
termSignal := make(chan os.Signal, 1)
signal.Notify(termSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

config := &disgord.Config{
    Token: os.Getenv("DISGORD_TOKEN"),
}
RESTClient := disgord.NewRESTClient(config)

gateway, err := disgord.GetGatewayBot(RESTClient)
if err != nil {
    return
}

// setup shards
fmt.Println("shards: " + strconv.Itoa(int(gateway.Shards)))
shards := make([]disgord.Session, gateway.Shards)
for i := uint(0); i < gateway.Shards; i++ {
    shards[i], err = disgord.NewSession(&disgord.Config{
        Token: config.Token,
        ShardID: i,
        TotalShards: gateway.Shards,
        WebsocketURL: gateway.URL,
    })

    // register a listener
    shards[i].On(disgord.EventMessageCreate, func(session disgord.Session, evt *disgord.MessageCreate) {
        fmt.Println("shard[" + session.ShardIDString() + "]" + evt.Message.Content)
    })
}

// connect
for i := uint(0); i < gateway.Shards; i++ {
    shards[i].Connect()
}

// disconnect
<-termSignal
for i := uint(0); i < gateway.Shards; i++ {
    shards[i].Disconnect()
}
```


### Get guild data from a shard
```go
shardID := disgord.GetShardForGuildID(guildID, shardCount)
session := shards[shardID]
```