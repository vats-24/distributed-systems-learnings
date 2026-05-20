package main

import (
	"fmt"

	"consistent-hash-ambassador/internal/ring"
)

func main() {

	r := ring.New(5)

	r.AddNode("backend-A")
	r.AddNode("backend-B")
	r.AddNode("backend-C")

	keys := []string{
		"user-1",
		"user-2",
		"user-3",
		"user-4",
	}

	for _, key := range keys {

		node := r.GetNode(key)

		fmt.Printf("%s -> %s\n", key, node)
	}
}
