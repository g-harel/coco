package httpc

import (
	"time"
)

type limiter struct {
	max    int
	ticker chan bool
}

func newLimiter(max int) *limiter {
	ticker := make(chan bool)
	go func() {
		for {
			time.Sleep(time.Second)
			for i := 0; i < max; i++ {
				ticker <- true
			}
		}
	}()
	return &limiter{max, ticker}
}

func (l *limiter) Wait() {
	<-l.ticker
	return
}
