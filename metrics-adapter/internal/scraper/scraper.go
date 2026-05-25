package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"metrics-adapter/internal/metrics"
	"net/http"
	"time"
)

type StatsResponse struct {
	Requests int `json:"requests"`
}

var client = &http.Client{

	Transport: &http.Transport{
		MaxIdleConns: 10,

		MaxIdleConnsPerHost: 5,

		IdleConnTimeout: 30 * time.Second,
	},
}

func Start(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)

	for {

		select {

		case <-ctx.Done():

			fmt.Println("Scrapper shutting down")

			return

		case <-ticker.C:

			fmt.Println("Scrapping foundation service...")

			ctx, cancel := context.WithTimeout(
				context.Background(),
				3*time.Second,
			)

			defer cancel()

			req, err := http.NewRequestWithContext(
				ctx,
				http.MethodGet,
				"http://localhost:8080/stats",
				nil,
			)

			if err != nil {

				fmt.Println("Request creation failed:", err)

				continue
			}

			if err != nil {

				fmt.Println("Scrape failed:", err)

				continue
			}

			resp, err := client.Do(req)

			if err != nil {

				fmt.Println("Scrape failed:", err)

				continue
			}

			var stats StatsResponse

			err = json.NewDecoder(resp.Body).Decode(&stats)

			resp.Body.Close()

			if err != nil {

				fmt.Println("JSON decode failed:", err)

				continue
			}

			metrics.FoundationRequests.WithLabelValues("foundation").Set(float64(stats.Requests))

			metrics.AdapterScrapesTotal.WithLabelValues("foundation").Inc()

			fmt.Println("Update metrics:", stats.Requests)

		}
	}
}
