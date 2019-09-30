### Removed
 - Client.SaveToDiscord / Session.SaveToDiscord: The behaviour was too ambiguous and introduced just another version of creating/updating structures.
 - Event channels (you now register channels instead, same as handlers)
 - support for Go versions below 1.13
 - removed a bunch of string pointers in the discord structs
 
### New
 - You can register channels to receive events (with controller, Ctrl, support)
 - Activity type consts
 - User.Tag(): return username#discriminator
 - User.AvatarURL(...): returns a valid URL to a user avatar, with GIF support.
 - internal loop for GetMessages so you don't have to worry about the discord limits
 - more functionality to the std pkg
 - allows LFU cache algorithm only (in preparations to weighted time aware LFU caching)
 - specify events to ignore instead of accepting
 - better support for distributed disgord instances with respect to ws sharding
 
### Internal changes
#### websocket
 - switched from gorilla/websocket to nhooyr/websocket (use disgord_websoket_gorilla build tag to revert)
 - simpler shard sync logic (might have fixed the possible deadlock)
 - stricter conditions for websocket communications (on errors, regardless, reconnect)
 - message queue, such that you do not lose outgoing websocket messages (resent on reconnects if required)
 - auto scaling of shards when discord requires you to scale up
 - you can now ignore presence and typing updates on a websocket level
 
#### Other
 - more accurate rate limiter support for REST requests (second => millisecond)
 - upgrade to snowflake v4
