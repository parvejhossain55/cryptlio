package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		logFields := map[string]interface{}{
			"method":     method,
			"status":     status,
			"latency":    latency.String(),
			"path":       path,
			"client_ip":  clientIP,
			"user_agent": c.Request.UserAgent(),
		}

		if query != "" {
			logFields["query"] = query
		}

		if userID, exists := c.Get("user_id"); exists {
			logFields["user_id"] = userID
		}

		if status >= 500 {
			log.Error().Fields(logFields).Msg("request failed")
		} else if status >= 400 {
			log.Warn().Fields(logFields).Msg("request error")
		} else {
			log.Info().Fields(logFields).Msg("request completed")
		}
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(500, gin.H{"error": err})
		}
		c.AbortWithStatus(500)
	})
}
