package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"metrics-adapter/internal/metrics"
	"metrics-adapter/internal/scraper"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	defer cancel()

	fmt.Println("Metrics adapter listening on :9090")

	metrics.Register()

	go scraper.Start(ctx)

	http.Handle("/metrics", promhttp.Handler())

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {

		sig := <-signalChan

		fmt.Println(
			"Received shutdown signal:",
			sig,
		)

		cancel()
	}()

	err := http.ListenAndServe(":9090", nil)

	if err != nil {

		fmt.Println("Server failed:", err)
	}
}
