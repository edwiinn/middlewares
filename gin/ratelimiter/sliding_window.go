package ratelimiter

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type slidingWindow struct {
	windowDuration time.Duration
	limit          int64
	windows        []time.Time
	mutex          sync.Mutex
}

func NewSlidingWindow(limit int64, windowDuration time.Duration) slidingWindow {
	return slidingWindow{
		limit:          limit,
		windowDuration: windowDuration,
		windows:        make([]time.Time, 0),
	}
}

func (s *slidingWindow) allow() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	for i := 0; i < len(s.windows); i++ {
		// If window time
		if now.Sub(s.windows[i]) > s.windowDuration {
			fmt.Println("here", i, now.Sub(s.windows[i]), s.windowDuration, len(s.windows))
			s.windows = s.windows[i+1:]
			i--
		} else {
			fmt.Println("1")
			break
		}
	}

	if int64(len(s.windows)) < s.limit {
		s.windows = append(s.windows, now)
		return true
	}
	return false
}

func (s *slidingWindow) SlidingWindowRatelimiterMiddleware(ctx *gin.Context) {
	if s.allow() {
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
