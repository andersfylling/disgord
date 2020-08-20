# JSON dependencies

There might be interest in changing the json encoder from "encoding/json", the standard implementation, for different
reasons. Some projects offer different build tags; Gin allows you to utilise `-tags=jsoniter` to swap out their 
internal use of "encoding/json" with jsoniter. However, this introduces N dependencies and their respective 
implementation to work as a drop in replacement.

A simpler manner can be dep injection; by introducing interfaces/function pointers in the Disgord config:
```go
client := disgord.New(disgord.Config{ JSONUnmarshaler: json.Unmarshal })
```

This may seem familiar with people of OOP background. The issue is that you need to throw this reference into every 
implementation that touches json encoding. It fails as soon as you need custom unmarshal/marshal methods.

github.com/diamondburned/arikawa, another Discord lib, treats this issue by exporting json related variables inside a json package. Allowing 
devs to directly overwrite the default value. I like it, so I'm going with a variation of that.

To change the json dependency, you simply import the disgord/json package and adjust the exported variables as needed.


Here the standard json implementation is swapped out with jsoniter:
```go
import (
    "github.com/andersfylling/disgord/json"
    jsoniter "github.com/json-iterator/go"
)

var j = jsoniter.ConfigCompatibleWithStandardLibrary

json.Marshal = j.Marshal
json.Unmarshal = j.Unmarshal
```