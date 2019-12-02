package internal

import (
	"fmt"
	"net/http"
	"time"
)

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
		"%v %4vms  %v\u001b[0m\n",
		res.StatusCode,
		time.Since(start).Truncate(time.Millisecond).Nanoseconds()/1e6,
		req.URL.String(),
	)
	if res.StatusCode == 200 {
		fmt.Print(message)
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

func LogError(format string, a ...interface{}) {
	err := fmt.Sprintf(format, a...)
	fmt.Printf("\u001b[31m%v\u001b[0m", err)
}
