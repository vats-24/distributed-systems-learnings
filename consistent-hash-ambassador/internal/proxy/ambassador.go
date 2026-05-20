package proxy

import "container/ring"

type Ambassador struct {
	ring *ring.Ring
}

func New(r *ring.Ring) *Ambassador {
	return &Ambassador{
		ring: r,
	}
}
