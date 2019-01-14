package constant

// I know this isn't optimal...

// Logger super basic logging interface
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
}
