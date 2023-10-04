package ratelimiter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// In the token bucket algorithm, tokens are added to a bucket at a constant rate.
// Requests or actions are allowed if there are enough tokens in the bucket to cover them.
// If there are no tokens, the request is rejected.
// This algorithm allows bursts of requests up to the bucket size.
type tokenBucket struct {
	capacity        int64
	tokens          int64
	rate            time.Duration
	lastUpdatedTime time.Time
}

func NewTokenBucket(capacity int64, rate time.Duration) (bucket tokenBucket) {
	return tokenBucket{
		capacity:        capacity,
		tokens:          capacity,
		lastUpdatedTime: time.Now(),
		rate:            rate,
	}
}

func (t *tokenBucket) TokenBucketRatelimiterMiddleware(ctx *gin.Context) {
	elapsedTime := time.Since(t.lastUpdatedTime)
	t.tokens += int64(elapsedTime / t.rate)
	if t.tokens > t.capacity {
		t.tokens = t.capacity
	}
	if t.tokens >= 1 {
		ctx.Next()
	} else {
		fmt.Println("Rate limit exceeded")
		// Timeout occurred, return a timeout response
		ctx.JSON(http.StatusTooManyRequests, map[string]interface{}{
			"code":    http.StatusTooManyRequests,
			"message": "Rate limit exceeded. Please try again later.",
		})
		return
	}
}
