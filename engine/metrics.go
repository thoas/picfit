package engine

import (
	"github.com/prometheus/client_golang/prometheus"
)

var defaultMetrics = newMetrics()

type metrics struct {
	operationSeconds *prometheus.HistogramVec
}

func newMetrics() *metrics {
	return &metrics{
		operationSeconds: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "picfit_operation_seconds",
				Help: "Duration of an image operation by backend",
			},
			[]string{"operation", "backend", "content_type"},
		),
	}
}

func init() {
	prometheus.MustRegister(defaultMetrics.operationSeconds)
}
