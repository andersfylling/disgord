package disgord

import (
	"fmt"
	"io"
)

// #########################################
// BELOW THIS LINE YOU WILL SEE COMMON AND POPULAR LOGGER INTERFACES.
// THESE EXIST HERE JUST TO SEE WHAT CAN BE MOST SANE ONE FOR DISGORD.
//

// _logInterface_apex https://github.com/apex/log/blob/master/interface.go
// note that this has removed struct types
type _logInterface_apex interface {
	//WithFields(fields Fielder) *Entry
	//WithField(key string, value interface{}) *Entry
	//WithError(err error) *Entry
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	Debugf(msg string, v ...interface{})
	Infof(msg string, v ...interface{})
	Warnf(msg string, v ...interface{})
	Errorf(msg string, v ...interface{})
	Fatalf(msg string, v ...interface{})
	//Trace(msg string) *Entry
}

// _logInterface_apsdehal https://github.com/apsdehal/go-logger/blob/master/logger.go
// note that some fields are missing
type _logInterface_apsdehal interface {
	Fatal(message string)
	FatalF(format string, a ...interface{})
	Fatalf(format string, a ...interface{})
	Panic(message string)
	PanicF(format string, a ...interface{})
	Panicf(format string, a ...interface{})
	Critical(message string)
	CriticalF(format string, a ...interface{})
	Criticalf(format string, a ...interface{})
	Error(message string)
	ErrorF(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Warning(message string)
	WarningF(format string, a ...interface{})
	Warningf(format string, a ...interface{})
	Notice(message string)
	NoticeF(format string, a ...interface{})
	Noticef(format string, a ...interface{})
	Info(message string)
	InfoF(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Debug(message string)
	DebugF(format string, a ...interface{})
	Debugf(format string, a ...interface{})
	StackAsError(message string)
	StackAsCritical(message string)
}

// _logInterface_logxi https://github.com/mgutz/logxi/blob/master/v1/logger.go
// note some fields might be missing
type _logInterface_logxi interface {
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{}) error
	Error(msg string, args ...interface{}) error
	Fatal(msg string, args ...interface{})
	Log(level int, msg string, args []interface{})
}

// _logInterface_onelog https://github.com/francoispqt/onelog/blob/master/logger.go
// note some fields might be missing
type _logInterface_onelog interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

// _logInterface_seelog https://github.com/cihub/seelog/blob/master/logger.go
// note some fields might be missing
type _logInterface_seelog interface {
	Tracef(format string, params ...interface{})
	Debugf(format string, params ...interface{})
	Infof(format string, params ...interface{})
	Warnf(format string, params ...interface{}) error
	Errorf(format string, params ...interface{}) error
	Criticalf(format string, params ...interface{}) error
	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{}) error
	Error(v ...interface{}) error
	Critical(v ...interface{}) error
	traceWithCallDepth(callDepth int, message fmt.Stringer)
	debugWithCallDepth(callDepth int, message fmt.Stringer)
	infoWithCallDepth(callDepth int, message fmt.Stringer)
	warnWithCallDepth(callDepth int, message fmt.Stringer)
	errorWithCallDepth(callDepth int, message fmt.Stringer)
	criticalWithCallDepth(callDepth int, message fmt.Stringer)
	SetContext(context interface{})
}

// _logInterface_spew https://github.com/davecgh/go-spew/blob/master/spew/spew.go
// note some fields might be missing
type _logInterface_spew interface {
	Errorf(format string, a ...interface{}) (err error)
	Fprint(w io.Writer, a ...interface{}) (n int, err error)
	Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
	Fprintln(w io.Writer, a ...interface{}) (n int, err error)
	Print(a ...interface{}) (n int, err error)
	Printf(format string, a ...interface{}) (n int, err error)
}

// _logInterface_tail https://github.com/hpcloud/tail/blob/master/tail.go
// note some fields might be missing
type _logInterface_tail interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}
