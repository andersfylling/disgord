Disgord allows you to inject a logger. The interface is fairly simple:

```go
type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
}
```

Both sirupsen/logrus and op/go-logging works out of the box without any need for wrapping.