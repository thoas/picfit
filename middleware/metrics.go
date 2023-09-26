package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type RouteKey struct{}

func Route(route string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), RouteKey{}, route))
		c.Next()
	}
}

func Metrics(c *gin.Context) {
	c.Next()
	route, ok := c.Request.Context().Value(RouteKey{}).(string)
	if ok {
		customMetrics.histogram.WithLabelValues(
			c.Request.Method,
			route,
			fmt.Sprint(c.Writer.Status()))
	}
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
