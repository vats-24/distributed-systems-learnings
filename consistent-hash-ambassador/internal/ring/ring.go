package ring

import (
	"fmt"
	"hash/fnv"
	"sort"
	"sync"
)

type Ring struct {
	mu sync.RWMutex

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

	r.mu.Lock()

	defer r.mu.Unlock()

	for i := 0; i < r.virtualNodes; i++ {
		virtualKey := fmt.Sprintf("%s#%d", node, i)

		position := hash(virtualKey)

		fmt.Println(position)

		r.positions = append(r.positions, position)

		r.nodes[position] = node

	}
	sort.Ints(r.positions)
}

func (r *Ring) GetNode(key string) string {

	r.mu.RLock()

	defer r.mu.RUnlock()

	if len(r.positions) == 0 {
		return ""
	}

	hashValue := hash(key)

	index := sort.SearchInts(r.positions, hashValue)

	if index == len(r.positions) {
		index = 0
	}

	position := r.positions[index]

	return r.nodes[position]
}

func (r *Ring) GetNodes(key string) []string {
	r.mu.Lock()

	defer r.mu.Unlock()

	if len(r.positions) == 0 {
		return nil
	}

	hashValue := hash(key)

	index := sort.SearchInts(r.positions, hashValue)

	if index == len(r.positions) {
		index = 0
	}

	seen := make(map[string]bool)

	var result []string

	for i := 0; i < len(r.positions); i++ {

		position := r.positions[(index+i)%len(r.positions)]

		node := r.nodes[position]

		if seen[node] {
			continue
		}

		seen[node] = true

		result = append(result, node)

	}
	return result
}
