package red

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTPRequests - returns prometheus counter for total HTTP request
var HTTPRequests = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	},
	[]string{"status_code", "method", "route"},
)

// Register - Registers the Prometheus metrics
func Register() {
	prometheus.Register(HTTPRequests)
}

func InstrumentationMiddleware(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n.ServeHTTP(w, r)
		HTTPRequests.WithLabelValues("200", r.Method, r.RequestURI).Inc()
	})
}
