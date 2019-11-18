package logging

import (
	"fmt"
	"net/http"
	"time"
)

var DefaultClient = Wrap(&http.Client{
	Transport: http.DefaultTransport,
})

type WrappedRoundTripper struct {
	original http.RoundTripper
}

func (w *WrappedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	res, err := w.original.RoundTrip(req)

	timeMillisecond := time.Since(start).Truncate(time.Millisecond).Nanoseconds() / 1e6
	colorCode := "" // none
	if timeMillisecond > 1000 {
		colorCode = "\u001b[33m" // yellow
	}
	if res.StatusCode != 200 {
		colorCode = "\u001b[31m" // red
	}
	if err == nil {
		fmt.Printf(
			"%v%v %4vms %v\u001b[0m\n",
			colorCode,
			res.StatusCode,
			timeMillisecond,
			req.URL.String(),
		)
	}
	return res, err
}

func Wrap(original *http.Client) *http.Client {
	return &http.Client{
		CheckRedirect: original.CheckRedirect,
		Jar:           original.Jar,
		Timeout:       original.Timeout,
		Transport:     &WrappedRoundTripper{original: original.Transport},
	}
}
