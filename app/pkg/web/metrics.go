package web

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func newMetricsServer(address string) *http.Server {
	mux := httprouter.New()
	mux.Handler("GET", "/metrics", promhttp.Handler())

	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         address,
		Handler:      mux,
	}
}
