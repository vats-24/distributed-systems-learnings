package proxy

import (
	"encoding/json"
	"net/http"

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

	node := a.ring.GetNode(shardKey)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"shard_key": shardKey,
		"backend":   node,
	})
}
