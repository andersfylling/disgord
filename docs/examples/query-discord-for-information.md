## Query Discord for information
In this article it is detailed how to fetch information about a user. You do also use the REST API to edit information, send messages, etc.

Imagine that you want to fetch information about a Discord user. Since Disgord supports every REST endpoint found in the default Discord API, you do not have to worry about handling the response format, compression, rate limits, etc. yourself. In stead you are given two solutions to solve query the Discord API; using the REST functions in:
 A: the rest package. AKA as barebones as it gets (not recommended).
 B: the Session interface. Which returns a channel for more flexible code.

The REST functions are located in the rest package and are blocking operations. It's important to note that these functions do not update the state (or the cache). They do however implement rate limiting (which you should not bypass). A reason to use these is when you want to force a request, as the functionality in the Session interface always checks the cache first to avoid asking for information that you already have locally. See [Query Discord using the REST functions](query-discord-using-the-rest-functions).

The Session interface is the recommended way to query objects from the Discord API. It supports caching, and even handles concurrency by returning a channel. By returning a channel you can decide the query to be either blocking or handle the response later (see [Query Discord using the session interface](#query-discord-using-the-session-interface)).



### Query Discord using the REST functions
> **Note:** It is assumed you understand how to create a session. You do not need to use Session.Connect and Session.Disconnect for _most_ REST queries.
```GoLang
// The user id of this repository's owner
userID := snowflake.NewID(228846961774559232)

// send a GET request to Discord to retrieve user information.
// the response is not stored in cache.
user, err := rest.GetUser(session.Req(), userID)
if err != nil {
    ...
}

// Even though you just sent this request, you force a new request without caring that you already have this information.
user2, err := rest.GetUser(session.Req(), userID)
if err != nil {
    ...
}
```


### Query Discord using the session interface
> **Note:** It is assumed you understand how to create a session. You do not need to use Session.Connect and Session.Disconnect for _most_ REST queries.
```GoLang
// The user id of this repository's owner
userID := snowflake.NewID(228846961774559232)

user := <- session.User(userID) // sends a request to discord
user2 := <- session.User(userID) // discovers that the information is already in the cache, and does not query Discord

userChan := session.User(userID) // you can ofc just get the channel and listen for it later on
// do some stuff here
user3 := <-userChan
```