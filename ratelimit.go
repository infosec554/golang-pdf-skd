package pdfsdk

import (
	"sync"
	"time"
)

type RateLimiter struct {
	tokens     chan struct{}
	interval   time.Duration
	maxTokens  int
	mu         sync.Mutex
	stopCh     chan struct{}
	refillOnce sync.Once
}

func NewRateLimiter(maxOps int, interval time.Duration) *RateLimiter {
	if maxOps <= 0 {
		maxOps = 100
	}
	if interval <= 0 {
		interval = time.Second
	}

	rl := &RateLimiter{
		tokens:    make(chan struct{}, maxOps),
		interval:  interval,
		maxTokens: maxOps,
		stopCh:    make(chan struct{}),
	}

	for i := 0; i < maxOps; i++ {
		rl.tokens <- struct{}{}
	}

	rl.startRefill()

	return rl
}

func (rl *RateLimiter) startRefill() {
	rl.refillOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(rl.interval)
			defer ticker.Stop()

			for {
				select {
				case <-rl.stopCh:
					return
				case <-ticker.C:
					rl.refill()
				}
			}
		}()
	})
}

func (rl *RateLimiter) refill() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for len(rl.tokens) < rl.maxTokens {
		select {
		case rl.tokens <- struct{}{}:
		default:
			return
		}
	}
}

func (rl *RateLimiter) Acquire() {
	<-rl.tokens
}

func (rl *RateLimiter) TryAcquire() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

func (rl *RateLimiter) Available() int {
	return len(rl.tokens)
}

func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
}

type RateLimitedSDK struct {
	*SDK
	limiter *RateLimiter
}

func (sdk *SDK) WithRateLimiter(maxOps int, interval time.Duration) *RateLimitedSDK {
	return &RateLimitedSDK{
		SDK:     sdk,
		limiter: NewRateLimiter(maxOps, interval),
	}
}

func (rls *RateLimitedSDK) Close() {
	rls.limiter.Stop()
	rls.SDK.Close()
}

func (rls *RateLimitedSDK) RateLimiter() *RateLimiter {
	return rls.limiter
}
