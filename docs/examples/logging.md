# Logging
You should have seen numerous examples of this already; but you can inject a logger of your choice into the disgord instance. If you don't inject a instance, nothing will be logged.

To inject a logger, the logger of your choice must implement the disgord.Logger interface:
```go
// See internal/logger/logger.go
type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
}
```

Out of the box you can use the default logger, a really bad wrapper of zap, where you can enable or disable debugging:
```go
client := disgord.New(disgord.Config{
    BotToken: "secret",
    Logger:   disgord.DefaultLogger(false), // debugging disabled
})
```

eg. sirupsen/logrus and op/go-logging also implements this interface, so you can inject them directly:
```go
client := disgord.New(disgord.Config{
    BotToken: "secret",
    Logger:   logrus.New(),
})
client2 := disgord.New(disgord.Config{
    BotToken: "secret",
    Logger:   logging.MustGetLogger("example"),
})
```

If you want to use another logger that does not implement the disgord.Logger interface by default, just write a wrapper:
```go
// logger.go
type MyLogger struct {
  ...
}

var _ disgord.Logger = (*MyLogger)(nil)

func (l *MyLogger) Debug(v ...interface{}) {...}
func (l *MyLogger) Info(v ...interface{}) {...}
func (l *MyLogger) Error(v ...interface{}) {...}

// main.go
client := disgord.New(disgord.Config{
    BotToken: "secret",
    Logger:   &MyLogger{},
})
```