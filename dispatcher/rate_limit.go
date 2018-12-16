package dispatcher

// https://github.com/golang/go/wiki/RateLimiting

import "time"

// RateLimit describe a rate limiter
type RateLimit struct {
	rate     time.Duration
	throttle chan time.Time
	ticker   *time.Ticker
	stopChan chan struct{}
}

// NewRateLimit return a RateLimiter instance
func NewRateLimit(rate time.Duration, burstLimit int) *RateLimit {
	return &RateLimit{
		rate:     rate,
		throttle: make(chan time.Time, burstLimit),
		stopChan: make(chan struct{}),
	}
}

// Start starts time.Ticker to rate limit.
func (r *RateLimit) Start() {
	r.ticker = time.NewTicker(r.rate)

	go func(ticker *time.Ticker, throttle chan time.Time) {
		for t := range ticker.C {
			select {
			case throttle <- t:
			default:
			}
		} // does not exit after tick.Stop()
	}(r.ticker, r.throttle)
}

// Throttle checks the operation rate.
func (r *RateLimit) Throttle() {
	<-r.throttle
}

// Stop stops rate limitting.
func (r *RateLimit) Stop() {
	r.ticker.Stop()
}
