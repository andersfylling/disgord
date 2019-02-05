package constant

// I know this isn't optimal...

// Logger super basic logging interface
type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
}
