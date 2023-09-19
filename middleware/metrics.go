package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func MetricsMiddlewares(c *gin.Context) {
	c.Next()
	customMetrics.histogram.WithLabelValues(
		c.Request.Method,
		c.Request.URL.String(),
		fmt.Sprint(c.Writer.Status()))

}

var customMetrics = newMetrics()

func newMetrics() *metrics {
	return &metrics{
		histogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: "http_request"},
			[]string{"http_method", "http_route", "http_status_code"},
		),
	}
}

type metrics struct{ histogram *prometheus.HistogramVec }

func init() { prometheus.MustRegister(customMetrics.histogram) }
