package metrics

import "github.com/prometheus/client_golang/prometheus"

var HttpRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of HTTP requests.",
	},
	[]string{"code", "operation"},
)

var HttpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_request_duration_seconds",
	Help: "Duration of HTTP requests.",
	Buckets: []float64{0.2, 0.5, 1, 2, 5},
}, []string{"operation"})

func init() {
	prometheus.MustRegister(HttpRequests, HttpDuration)
}