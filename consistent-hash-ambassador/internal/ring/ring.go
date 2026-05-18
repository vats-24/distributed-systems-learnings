package ring

import (
	"fmt"
	"hash/fnv"
	"sort"
)

type Ring struct {
	virtualNodes int

	positions []int

	nodes map[int]string
}

func New(virtualNodes int) *Ring {

	return &Ring{
		virtualNodes: virtualNodes,
		positions:    []int{},
		nodes:        make(map[int]string),
	}
}

func hash(key string) int {
	h := fnv.New32a()

	h.Write([]byte(key))

	return int(h.Sum32())
}

func (r *Ring) AddNode(node string) {
	for i := 0; i < r.virtualNodes; i++ {
		virtualKey := fmt.Sprintf("%s#%d", node, i)

		position := hash(virtualKey)

		r.positions = append(r.positions, position)

		r.nodes[position] = node

	}

	sort.Ints(r.positions)

}
