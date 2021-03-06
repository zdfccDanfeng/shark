package helper

import (
	"fmt"
	"runtime"
	"time"
)

// loggingOn is a simple flag to turn logging on or off.
var loggingOn = true

// TurnLoggingOff sets the logging flag to off.
func TurnLoggingOff() {
	loggingOn = false
}

// WriteStdout is used to write message directly stdout.
func WriteStdout(goRoutine string, functionName string, message string) {
	if loggingOn == true {
		fmt.Printf("%s : %s : %s : %s\n", time.Now().Format("2006-01-02T15:04:05.000"), goRoutine, functionName, message)
	}
}

// WriteStdoutf is used to write a formatted message directly stdout.
func WriteStdoutf(goRoutine string, functionName string, format string, a ...interface{}) {
	WriteStdout(goRoutine, functionName, fmt.Sprintf(format, a...))
}

// CatchPanic is used to catch and display panics
func CatchPanic(err *error, goRoutine string, function string) {
	if r := recover(); r != nil {
		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		WriteStdoutf(goRoutine, function, "PANIC Defered [%v] : Stack Trace : %v", r, string(buf))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}
