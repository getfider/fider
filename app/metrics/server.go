package metrics

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartHTTPServer(address string) {
	mux := httprouter.New()
	mux.Handler("GET", "/metrics", promhttp.Handler())
	err := http.ListenAndServe(address, mux)
	if err != nil {
		panic(err)
	}
}
