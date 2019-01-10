Disgord supports the use of sharding for as explained here: [discordapp.com/.../gateway#sharding](https://discordapp.com/developers/docs/topics/gateway#sharding)

DisGord uses an internal shard manager to handle this for you. However, you are able to specify a maximum number of shards, an offset on shards IDs and the connection URL.

# Enforce N number of shards
```go
client := disgord.NewClient(&disgord.Config{
    ShardManagerConfig: &disgord.ShardManagerConfig{
        // If you have another instance running with shards 0-3, and want this instance to use the range 4-8
        // you can specify the number of the first shard this instance should have. Otherwise there is no
        // reason for you to tweak this.
        FirstID: 4, //offset. 
		
        // no less and no more than 4 shards. 
        // Setting this to 0 will allow DisGord to ask Discord for the recommented 
        // amount in respect to how many guilds your bot exists in.
        ShardLimit: 4, // number of shards.
        
        URL: "", // if not set, DisGord will contact Discord to get one.
    },
})
```

The entire ShardManagerConfig is optional as your bot will always use sharding by default, and automatically decide on how many shards you need. There is also no need to communicate with individual shards in this design, so when you interact with the DisGord interface you won't even realise that sharding exists after the client has been configured.

```go
// this client is also using shards
client := disgord.NewClient(&disgord.Config{
    BotToken: "random token",
})
```