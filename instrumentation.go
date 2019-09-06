package red

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTPRequests - returns prometheus counter for total HTTP request
var HTTPRequests = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	},
	[]string{"status_code", "method", "route", "content_length"},
)

// Register - Registers the Prometheus metrics
func Register() {
	prometheus.Register(HTTPRequests)
}

type statusWriter struct {
	http.ResponseWriter
	Status int
	Length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.Status == 0 {
		w.Status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.Length += n
	return n, err
}

func InstrumentationMiddleware(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := statusWriter{ResponseWriter: w}
		n.ServeHTTP(&sw, r)
		HTTPRequests.WithLabelValues(strconv.Itoa(sw.Status), r.Method, r.RequestURI, strconv.Itoa(sw.Length)).Inc()
	})
}
