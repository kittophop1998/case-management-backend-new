package monitoring

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	HTTPRequestCounter *prometheus.CounterVec
	HTTPRequestLatency *prometheus.HistogramVec
}

func NewPrometheus(namespace string) *Prometheus {
	return &Prometheus{
		HTTPRequestCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests.",
			},
			[]string{"method", "path", "status_code"},
		),
		HTTPRequestLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_latency_seconds",
				Help:      "Latency of HTTP requests.",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}
}
func (p *Prometheus) Handler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}
