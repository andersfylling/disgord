# Design decisions

## Response Header 

#### Rate Limit Bucket - Reset timestamp
Discord provices a bunch of header fields to state when a bucket is reset. I have no clue why there are so many, so the function `CorrectDiscordHeader` makes sure that the header field `X-RateLimit-Reset` is populated and displays milliseconds not seconds; when `Retry-After` or `X-RateLimit-Reset-After` or a json body with the field `retry_after` is set and holds content.

Visually the following fields wrapped in ~~ should be ignored. The `X-RateLimit-Reset` is the one of interest.
```markdown
> GET /api/v6/some-endpoint
> X-RateLimit-Precision: millisecond

< HTTP/1.1 429 TOO MANY REQUESTS
< Content-Type: application/json
~~< Retry-After: 6457~~
< X-RateLimit-Limit: 10
< X-RateLimit-Remaining: 0
< X-RateLimit-Reset: 1470173023000
~~< X-RateLimit-Reset-After: 7~~
< X-RateLimit-Bucket: abcd1234
{
  "message": "You are being rate limited.",
  ~~"retry_after": 6457,~~
  "global": false
}
```

## Rate Limit and endpoint relationships
Linking a given endpoint to a bucket is a serious hassle as Discord does not return the bucket hash on each response. This is really unfortunate as we can't establish relationships between buckets and endpoints before Discord decides to randomly send them (...). This also regards the HTTP methods, even tho the documentation states otherwise. [The docs are useless for insight on this matter.](https://github.com/discord/discord-api-docs/issues/1135)

To tackle this, every hashed endpoint is assumed to have its own bucket. To hash an endpoint before it can be linked to a bucket, the snowflakes, except major snowflakes, must be replaced with the string `{id}`.  Major snowflakes are the first snowflakes in any endpoint with the prefixes ["/channels", "/guilds", "/webhooks"].

``` 
// examples of before and after hashed endpoints (without http methods)
/test => /test
/test/4234 => /test/{id}
/channels => /channels 
/channels?limit=12 => /channels 
/channels/895349573 => /channels/895349573
/guilds/35347862384/roles/23489723 => /guilds/35347862384/roles/{id}
```

This will of course cause some hashed endpoints to point to the same bucket, once Discord gives us insight on a per request basis. While endpoint A and B might use the same bucket we won't know unless a request for both A and B returns the hash. 

The list of buckets can then be consolidated, memory wise, if two or more hashed endpoints uses the same bucket key/hash. The most recent bucket should overwrite the older bucket and all related endpoints should then point to the same bucket.

#### Concurrent requests
~~Due to the above, we can only send one request (per bucket) at the time until Discord returns rate limit information for the given bucket. Once the bucket is reset after Discord Disgord must return to sending request through the same bucket sequentially to reduce potential rate limits.~~

~~For a normal header, without any rate limit information, we can only send sequential requests for the bucket.~~
```markdown
< HTTP/1.1 200 OK
```

~~When Discord gives us some more bucket info we can use this to send concurrent requests. In this example we can send up to 7 concurrent requests until the time hits the Reset unix. Then we're back to sequential bucket requests.~~
```markdown
< HTTP/1.1 200 OK
< X-RateLimit-Limit: 10
< X-RateLimit-Remaining: 7
< X-RateLimit-Reset: 1470173023000
< X-RateLimit-Bucket: abcd1234
```

Discord is messed up and each bucket runs requests in a sequential fashion. For now this is the best that can be done with Discord.