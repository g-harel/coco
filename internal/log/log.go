package log

import (
	"fmt"
	"os"

	"github.com/g-harel/coco/internal/flags"
)

// Info prints muted text to the console.
// It can be enabled or disabled using the "log-info" flag.
func Info(format string, a ...interface{}) {
	if *flags.LogInfo {
		msg := fmt.Sprintf(format, a...)
		fmt.Printf("\u001b[38;5;244m%v\u001b[0m", msg)
	}
}

// Error prints red error messages to stderr.
// It can be enabled or disabled using the "log-error" flag.
func Error(format string, a ...interface{}) {
	if *flags.LogErrors {
		err := fmt.Sprintf(format, a...)
		fmt.Fprintf(os.Stderr, "\u001b[31m%v\u001b[0m", err)
	}
}

// Output prints program output.
func Output(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}
