package ratelimit

import (
	"sync"
	"testing"
	"time"
)

const (
	qps      = 1000
	multiple = 2
)

func TestRateLimitTry(t *testing.T) {
	r := New(qps, multiple * qps)
	r.Start()

	count := 0
	go func() {
		for i := 0; i < multiple * qps; i++ {
			time.Sleep(time.Second / time.Duration(multiple * qps))
			if r.Try() {
				count++
			}
		}
	}()
	time.Sleep(time.Second)
	if count > qps + qps/10  || count < qps - qps/10 {
		t.Errorf("expected count: %d, got: %d", qps, count)
	}

}

func TestRateLimitWait(t *testing.T) {
	r := New(qps, multiple * qps)
	r.Start()

	wg := sync.WaitGroup{}
	wg.Add(multiple * qps)
	before := time.Now()
	for i := 0; i < multiple * qps; i++ {
		go func() {
			r.Wait()
			wg.Done()
		}()
	}
	wg.Wait()
	spendSeconds := int(time.Now().Sub(before).Seconds())
	if spendSeconds != multiple {
		t.Errorf("expected: %ds, takes: %ds", multiple * qps, spendSeconds)
	}
}