package scraper

import (
	"encoding/json"
	"fmt"
	"metrics-adapter/internal/metrics"
	"net/http"
	"time"
)

type StatsResponse struct {
	Requests int `json:"requests"`
}

func Start() {
	ticker := time.NewTicker(10 * time.Second)

	for {
		<-ticker.C

		fmt.Println("Scrapping foundation service...")

		resp, err := http.Get(
			"http://localhost:8080/stats",
		)

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
