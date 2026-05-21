package main

import (
	"fmt"
	"metrics-adapter/internal/metrics"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	fmt.Println("Metrics adapter listening on :9090")

	metrics.Register()

	// go scraper.Start()

	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":9090", nil)

	if err != nil {

		fmt.Println("Server failed:", err)
	}
}
