package picfit

import (
	"github.com/prometheus/client_golang/prometheus"
)

var defaultMetrics = newMetrics()

type metrics struct {
	histogram    *prometheus.HistogramVec
	sourcePixels *prometheus.HistogramVec
}

func newMetrics() *metrics {
	return &metrics{
		histogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: "picfit_action_seconds"},
			[]string{"picfit_method", "picfit_image_type"},
		),
		// decoded size is roughly 4 bytes per pixel: 0.1 MP (400 KB) up to 1.6 GP (6.4 GB)
		sourcePixels: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "picfit_source_image_pixels",
				Help:    "Number of pixels (width x height) of source images to process",
				Buckets: prometheus.ExponentialBuckets(1e5, 4, 8),
			},
			[]string{"format"},
		),
	}
}

func init() {
	prometheus.MustRegister(defaultMetrics.histogram, defaultMetrics.sourcePixels)
}
