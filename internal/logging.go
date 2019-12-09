package internal

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/g-harel/coco/internal/flags"
)

var zero = time.Now().Truncate(time.Millisecond).UnixNano() / 1e6

var DefaultLoggingClient = NewLoggingClient(&http.Client{
	Transport: http.DefaultTransport,
})

type loggingRoundTripper struct {
	original http.RoundTripper
}

func (w *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	res, err := w.original.RoundTrip(req)
	if err != nil {
		return res, err
	}

	message := fmt.Sprintf(
		"%v %v+%vms %v\n",
		res.StatusCode,
		start.Truncate(time.Millisecond).UnixNano()/1e6-zero,
		time.Since(start).Truncate(time.Millisecond).Nanoseconds()/1e6,
		req.URL.String(),
	)
	if res.StatusCode == 200 {
		LogInfo(message)
	} else {
		LogError(message)
	}

	return res, err
}

func NewLoggingClient(original *http.Client) *http.Client {
	return &http.Client{
		CheckRedirect: original.CheckRedirect,
		Jar:           original.Jar,
		Timeout:       original.Timeout,
		Transport:     &loggingRoundTripper{original: original.Transport},
	}
}

func LogInfo(format string, a ...interface{}) {
	if *flags.LogInfo {
		msg := fmt.Sprintf(format, a...)
		fmt.Printf("\u001b[38;5;244m%v\u001b[0m", msg)
	}
}

func LogError(format string, a ...interface{}) {
	if *flags.LogErrors {
		err := fmt.Sprintf(format, a...)
		fmt.Fprintf(os.Stderr, "\u001b[31m%v\u001b[0m", err)
	}
}
