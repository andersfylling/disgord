# Snowflake for Disgord
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)

Does not hold functionality to connect a snowflake service, but rather parsing the snowflakes for Discord only(!).

Usage:

```go
import . "github.com/andersfylling/disgord/snowflake"

type DiscordRole struct {
    ID          Snowflake    `json:"id"`
    Name        string       `json:"name"`
    Managed     bool         `json:"managed"`
    Mentionable bool         `json:"mentionable"`
    Hoist       bool         `json:"hoist"`
    Color       int          `json:"color"`
    Position    int          `json:"position"`
    Permissions uint64       `json:"permissions"`
}
```

If you're creating an API that sends JSON to a multiple different language clients, some might not be able to process uint64, such as javascript. To support both uint64 and string use the JSON struct included:

```go
import . "github.com/andersfylling/disgord/snowflake"

type DiscordRole struct {
    *SnowflakeJSON           `json:"snowflake"`
    Name        string       `json:"name"`
    Managed     bool         `json:"managed"`
    Mentionable bool         `json:"mentionable"`
    Hoist       bool         `json:"hoist"`
    Color       int          `json:"color"`
    Position    int          `json:"position"`
    Permissions uint64       `json:"permissions"`
}
```

This adds two fields: `ID` and `IDStr`. Where the first is of a snowflake.ID(uint64), and the second is a string. This creates the JSON format (IDs only. Where the dots represents the remaining DiscordRole fields):

```json
{
    "snowflake": {
          "id": 74895735435643,
          "id_str": "74895735435643",
    },
    ...
}
```

Now an alternative is to send only the string version by adding `,string` to the json tag. Which I would recommend instead:

```go
import . "github.com/andersfylling/disgord/snowflake"

type DiscordRole struct {
    ID          Snowflake    `json:"id,string"`
    Name        string       `json:"name"`
    Managed     bool         `json:"managed"`
    Mentionable bool         `json:"mentionable"`
    Hoist       bool         `json:"hoist"`
    Color       int          `json:"color"`
    Position    int          `json:"position"`
    Permissions uint64       `json:"permissions"`
}
```

This does fulfill the twitter snowflake use case described here: <https://developer.twitter.com/en/docs/basics/twitter-ids>

Remember that Discord has a different epoch. So when using the Date function, this will only function for Discord applications.
