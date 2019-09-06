package red

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var RequestTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	},
	[]string{"status_code", "method", "route"},
)

var RequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Time (in seconds) spent serving HTTP requests.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"status_code", "method", "route"},
)

func Register() {
	prometheus.Register(RequestDuration)
	prometheus.Register(RequestTotal)
}

type statusWriter struct {
	http.ResponseWriter
	Status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.Status == 0 {
		w.Status = 200
	}
	return w.ResponseWriter.Write(b)
}

func InstrumentationMiddleware(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := statusWriter{ResponseWriter: w}
		startTime := time.Now()
		n.ServeHTTP(&sw, r)
		duration := time.Now().Sub(startTime)
		RequestTotal.WithLabelValues(strconv.Itoa(sw.Status), r.Method, r.RequestURI).Inc()
		RequestDuration.WithLabelValues(strconv.Itoa(sw.Status), r.Method, r.RequestURI).Observe(duration.Seconds())
	})
}
