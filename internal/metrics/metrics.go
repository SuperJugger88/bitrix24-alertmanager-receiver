package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "request_duration_seconds",
			Help: "Duration of requests in seconds",
		},
		[]string{"method", "endpoint"},
	)
)
