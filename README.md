#### ratelimit
A token bucket rate limiter in less than 100 lines of code.

#### Usage
```go
qps := 100
maxConcurrency := 200
r := ratelimit.New(qps, maxConcurrency)
r.Start()

r.Try()
r.Wait()
r.WaitTimeout(time.Second)

r.Stop()
```
