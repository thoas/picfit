package picfit

import (
	"github.com/prometheus/client_golang/prometheus"
)

var defaultMetrics = newMetrics()

type metrics struct {
	histogram *prometheus.HistogramVec
}

func newMetrics() *metrics {
	return &metrics{
		histogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: "picfit_action_seconds"},
			[]string{"picfit_method", "picfit_image_type"},
		),
	}
}

func init() {
	prometheus.MustRegister(defaultMetrics.histogram)
}
