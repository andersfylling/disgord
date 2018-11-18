This example showcases an easy to way to change your Discord client's avatar. Note: put this code in any place you have session defined.

```go
// load our png image, and close the body
file, _ := os.Open("./avatar.png")
defer file.Body.Close()
// initialize modify user parameters
params := &disgord.ModifyCurrentUserParams{}
params.SetAvatarImage(file)
// update our client
// s - *disgord.Session
s.ModifyCurrentUser(params)
```

>  **Note**: if you already have a base64 string - you can just use the SetAvatar method instead.
```go
params.SetAvatar(base64str)
```