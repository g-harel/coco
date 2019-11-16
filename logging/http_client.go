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
	if err == nil {
		fmt.Printf("%v %v %v\n", res.StatusCode, req.URL.String(), time.Since(start).Truncate(time.Millisecond))
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
