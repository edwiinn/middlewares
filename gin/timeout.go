package gin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(ctx *gin.Context) {
	timeoutCh := make(chan bool, 1)
	// Use a Goroutine to handle the request
	go func() {
		ctx.Next()
		timeoutCh <- false
	}()

	timeout := time.Second * 2
	// Wait for either the request to complete or the timeout to occur
	select {
	case <-timeoutCh:
		// Request completed, do nothing
	case <-time.After(timeout):
		fmt.Println("Timeout")
		// Timeout occurred, return a timeout response
		ctx.JSON(http.StatusRequestTimeout, map[string]interface{}{
			"code":    http.StatusRequestTimeout,
			"message": "timeout",
		})
		return
	}
}
