package middleware

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/penonpaper/event-api/internal/metrics"
)

func PrometheusMiddleware(m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		m.InFlightRequests.Inc()
		defer m.InFlightRequests.Dec()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = "unknown"
		}

		m.HttpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
		m.HttpRequestDuration.WithLabelValues(method, endpoint, status).Observe(duration)
	}
}

func LogginMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[Started] method=%s path=%s", c.Request.Method, c.Request.URL.Path)

		c.Next()

		log.Printf("[Completed] method=%s path=%s status=%d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	}
}
