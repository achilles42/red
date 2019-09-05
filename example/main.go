package main

import (
	"log"
	"net/http"
	"os"

	red "github.com/achilles42/red"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func HandleHelloRequest() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
}

func main() {
	addr := os.Getenv("ADDR")

	mux := http.NewServeMux()
	mux.Handle("/v1/hello", red.InstrumentationMiddleware(HandleHelloRequest()))
	mux.Handle("/metrics", promhttp.Handler())

	log.Printf("server is listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
