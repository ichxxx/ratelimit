package ratelimit

import (
	"context"
	"time"
)

type rateLimit struct {
	ctx     context.Context
	rate    time.Duration
	ticker  *time.Ticker
	tokens  chan struct{}
}

func New(qps, maxConcurrency int) *rateLimit {
	return &rateLimit{
		rate:   time.Second / time.Duration(qps),
		tokens: make(chan struct{}, maxConcurrency),
	}
}

func (r *rateLimit) generateToken() {
	for {
		select {
		case <-r.ticker.C:
			r.tokens <- struct{}{}
		case <-r.ctx.Done():
			r.ticker.Stop()
			return
		}
	}
}

func (r *rateLimit) WithContext(ctx context.Context) {
	r.ctx = ctx
}

func (r *rateLimit) Start() {
	if r.ctx == nil {
		r.ctx = context.Background()
	}
	r.ticker = time.NewTicker(r.rate)
	go r.generateToken()
}

func (r *rateLimit) SetRate(qps int) {
	r.rate = time.Second / time.Duration(qps)
	if r.ticker != nil {
		r.ticker.Reset(r.rate)
	}
}

func (r *rateLimit) Try() bool {
	select {
	case <-r.tokens:
		return true
	default:
		return false
	}
}

func (r *rateLimit) Wait() {
	<-r.tokens
}

func (r *rateLimit) WaitMulti(count int) {
	for count > 0 {
		<-r.tokens
		count--
	}
}

func (r *rateLimit) WaitTimeout(timeout time.Duration) bool {
	select {
	case <-r.tokens:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (r *rateLimit) WaitMultiTimeout(count int, timeout time.Duration) int {
	n := count
	t := time.After(timeout)
	for n > 0 {
		select {
		case <-r.tokens:
			n--
		case <-t:
			return count - n
		}
	}
	return count - n
}