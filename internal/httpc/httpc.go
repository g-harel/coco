package httpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/g-harel/coco/internal/flags"
	"github.com/g-harel/coco/internal/log"
)

var zero = time.Now().Truncate(time.Millisecond).UnixNano() / 1e6
var rateLimiter = newLimiter(*flags.RateLimit, time.Second)

// LogHTTP logs a formatted version of the request with:
// - response status code
// - program time at which request was sent (in ms)
// - round trip time (in ms)
// - request URL
func logHTTP(url *url.URL, res *http.Response, start time.Time) {
	message := fmt.Sprintf(
		"%v %v+%vms %v\n",
		res.StatusCode,
		start.Truncate(time.Millisecond).UnixNano()/1e6-zero,
		time.Since(start).Truncate(time.Millisecond).Nanoseconds()/1e6,
		url.String(),
	)
	if res.StatusCode == 200 {
		log.Info(message)
	} else {
		log.Error(message)
	}
}

// Get is a wrapper to make simple GET requests.
// It abstracts away logic around rate limiting, logging,
// response decoding and error handling.
func Get(rawURL string, headers http.Header, responseBody interface{}) (*http.Header, error) {
	<-rateLimiter

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	res, err := http.DefaultClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: headers,
	})
	if err != nil {
		return nil, err
	}
	logHTTP(u, res, start)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %v", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(responseBody)
	if err != nil {
		return nil, err
	}

	return &res.Header, nil
}
