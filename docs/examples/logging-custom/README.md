Disgord allows you to inject a logger. The interface is fairly simple:

```go
type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
}
```

In the file example-custom.go we create the basic requirement for a wrapper around the zap library. To improve readability, remember to skip one caller level as shown in the example.
