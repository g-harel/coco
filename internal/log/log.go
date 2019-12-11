package log

import (
	"fmt"
	"os"

	"github.com/g-harel/coco/internal/flags"
)

func Info(format string, a ...interface{}) {
	if *flags.LogInfo {
		msg := fmt.Sprintf(format, a...)
		fmt.Printf("\u001b[38;5;244m%v\u001b[0m", msg)
	}
}

func Error(format string, a ...interface{}) {
	if *flags.LogErrors {
		err := fmt.Sprintf(format, a...)
		fmt.Fprintf(os.Stderr, "\u001b[31m%v\u001b[0m", err)
	}
}
