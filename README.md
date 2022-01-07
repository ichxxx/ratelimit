A token bucket rate limiter in less than 100 lines of code.

#### Usage
```go
r := ratelimit.New(qps, maxConcurrency).WithContext(ctx)
r.Start()

r.Try()
r.Wait()
r.WaitMulti(10)
r.WaitTimeout(time.Second)
r.WaitMultiTimeout(10, time.Second)
```
