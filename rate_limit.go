package ratelimit

import (
	"time"
)

type rateLimit struct {
	enable  bool
	rate    time.Duration
	ticker  *time.Ticker
	tokens  chan struct{}
}

func New(qps, maxConcurrency int) *rateLimit {
	return &rateLimit{
		rate:    getRate(qps),
		tokens:  make(chan struct{}, maxConcurrency),
	}
}

func (r *rateLimit) generateToken() {
	for r.enable {
		<- r.ticker.C
		r.tokens <- struct{}{}
	}
}

func (r *rateLimit) Start() {
	if r.enable {
		return
	}
	r.enable = true
	r.ticker = time.NewTicker(r.rate)
	go r.generateToken()
}

func (r *rateLimit) Stop() {
	if r.enable {
		r.ticker.Stop()
		r.enable = false
	}
}

func (r *rateLimit) SetRate(qps int) {
	r.rate = getRate(qps)
	if r.enable {
		r.ticker.Reset(r.rate)
	}
}

func (r *rateLimit) Try() bool {
	if r.enable {
		select {
		case <-r.tokens:
			return true
		default:
			return false
		}
	}
	return true
}

func (r *rateLimit) Wait() {
	if r.enable {
		<-r.tokens
	}
}

func (r *rateLimit) WaitTimeout(timeout time.Duration) bool {
	if r.enable {
		select {
		case <-r.tokens:
			return true
		case <-time.After(timeout):
			return false
		}
	}
	return true
}

func getRate(qps int) time.Duration {
	return time.Second / time.Duration(qps)
}