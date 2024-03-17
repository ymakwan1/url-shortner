package middleware

import (
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	burst      int           // Max no. of tokens in bucket
	refill     time.Duration //Interval at which tokens are added to the bucket
	tokens     int           //Current number of toeksn in the bucket
	lastRefill time.Time     //Last time the bucket was refilled
	mutex      sync.Mutex    //Mutex lock for thread safe
}

func NewTokenBucket(burst int, refill time.Duration) *TokenBucket {
	return &TokenBucket{
		burst:      burst,
		refill:     refill,
		tokens:     burst, //Intitial bucket is full
		lastRefill: time.Now(),
		mutex:      sync.Mutex{},
	}
}

func (tb *TokenBucket) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tb.mutex.Lock()
		defer tb.mutex.Unlock()

		tb.refillTokens()

		if tb.tokens > 0 {
			tb.tokens--
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Rate Limit Exceeded", http.StatusTooManyRequests)
		}
	})
}

func (tb *TokenBucket) refillTokens() {
	now := time.Now()
	timeElapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int(timeElapsed / tb.refill)
	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd

		if tb.tokens > tb.burst {
			tb.tokens = tb.burst
		}
		tb.lastRefill = now
	}

}
