package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"consistent-hash-ambassador/internal/ring"
)

type Ambassador struct {
	ring *ring.Ring
}

func New(r *ring.Ring) *Ambassador {
	return &Ambassador{
		ring: r,
	}
}

func (a *Ambassador) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	shardKey := r.Header.Get("X-Shard-Key")

	if shardKey == "" {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "missing shard key header",
		})

		return
	}

	// node := a.ring.GetNode(shardKey)
	nodes := a.ring.GetNodes(shardKey)

	for _, node := range nodes {
		targetURL, err := url.Parse(node)

		if err != nil {
			fmt.Println("Invalid backend url", node)

			continue
		}

		success := true
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			fmt.Println("Backend failed:", node)
			fmt.Println("Error:", err)

			success = false
		}

		proxy.ServeHTTP(w, r)

		if success {

			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusBadGateway)

	json.NewEncoder(w).Encode(map[string]string{
		"error": "all backends unavailable",
	})
}
