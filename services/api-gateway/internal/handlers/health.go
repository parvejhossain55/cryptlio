package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Services  map[string]string `json:"services,omitempty"`
	Uptime    string            `json:"uptime,omitempty"`
}

var startTime = time.Now()

func HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		services := make(map[string]string)

		// TODO: Add health checks for backend services
		services["wallet"] = "ok"
		services["trade_engine"] = "ok"
		services["kyc"] = "ok"
		services["notification"] = "ok"
		services["dispute"] = "ok"

		c.JSON(http.StatusOK, HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now().Format(time.RFC3339),
			Services:  services,
			Uptime:    time.Since(startTime).String(),
		})
	}
}

func LivenessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
		})
	}
}

func ReadinessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	}
}
