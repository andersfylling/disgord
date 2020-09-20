package logger

import "fmt"

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

type FmtPrinter struct{}

func (FmtPrinter) Debug(v ...interface{}) {
	fmt.Print("debug -- ")
	fmt.Println(v...)
}
func (FmtPrinter) Info(v ...interface{}) {
	fmt.Print("info -- ")
	fmt.Println(v...)
}
func (FmtPrinter) Error(v ...interface{}) {
	fmt.Print("error -- ")
	fmt.Println(v...)
}

var _ Logger = (*FmtPrinter)(nil)
