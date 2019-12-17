package httpc

import (
	"time"
)

// NewLimiter creates a channel that produces the provided
// count of values on every delay. The limiter does not
// compensate for time spent blocked and only waits and
// resets the count when an entire batch has been produced
// and consumed.
func newLimiter(count int, delay time.Duration) chan bool {
	limiter := make(chan bool)
	go func() {
		for {
			for i := 0; i < count; i++ {
				limiter <- true
			}
			time.Sleep(delay)
		}
	}()
	return limiter
}
