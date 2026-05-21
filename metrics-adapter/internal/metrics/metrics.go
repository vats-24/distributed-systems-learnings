package metrics

import "github.com/prometheus/client_golang/prometheus"

var AdapterScrapesTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "adapter_scrapes_total",
		Help: "Total unsuccessful adapter scrapes",
	},
	[]string{"service"},
)

var FoundationRequests = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "foundation_requests",
		Help: "Current request count from foundation service",
	},
	[]string{"service"},
)

func Register() {
	prometheus.MustRegister(AdapterScrapesTotal)

	prometheus.MustRegister(FoundationRequests)
}
