package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDs carries correlation IDs across calls and assigns each request a local ID.
func RequestIDs() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.NewString()
		}
		c.Set("request_id", requestID)
		c.Set("correlation_id", correlationID)
		c.Header("X-Request-ID", requestID)
		c.Header("X-Correlation-ID", correlationID)
		c.Next()
	}
}
