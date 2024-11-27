package middleware

import (
	"apisecurityplatform/pkg/observability"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "not_found"
		}
		method := c.Request.Method

		// Track request size (new)
		if c.Request.ContentLength > 0 {
			observability.HTTPServerRequestSize.WithLabelValues(
				method,
				path,
			).Observe(float64(c.Request.ContentLength))
		}

		// Track active requests (new)
		observability.HTTPServerActiveRequests.WithLabelValues(
			method,
			path,
		).Inc()

		// Process request
		c.Next()

		// Decrease active requests (new)
		observability.HTTPServerActiveRequests.WithLabelValues(
			method,
			path,
		).Dec()

		// Record metrics after request is processed
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// Existing metrics
		observability.RequestCounter.WithLabelValues(method, path, status).Inc()
		observability.ResponseTime.WithLabelValues(method, path).Observe(duration)

		if c.Writer.Status() >= 400 {
			errorType := "client_error"
			if c.Writer.Status() >= 500 {
				errorType = "server_error"
			}
			observability.ErrorCounter.WithLabelValues(method, path, errorType).Inc()
		}

		// New OpenTelemetry-aligned metrics
		observability.HTTPServerDuration.WithLabelValues(
			method,
			path,
			status,
		).Observe(duration)

		// Track response size (new)
		responseSize := float64(c.Writer.Size())
		if responseSize > 0 {
			observability.HTTPServerResponseSize.WithLabelValues(
				method,
				path,
			).Observe(responseSize)
		}
	}
}
