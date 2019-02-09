package logger

// Logger super basic logging interface
type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
}

type Empty struct{}

func (Empty) Debug(v ...interface{}) {}
func (Empty) Info(v ...interface{})  {}
func (Empty) Error(v ...interface{}) {}

var _ Logger = (*Empty)(nil)
