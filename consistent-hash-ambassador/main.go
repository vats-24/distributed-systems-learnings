package main

import (
	"fmt"
	"net/http"

	"consistent-hash-ambassador/internal/proxy"
	"consistent-hash-ambassador/internal/ring"
)

func main() {

	r := ring.New(10)

	r.AddNode("backend-A")
	r.AddNode("backend-B")
	r.AddNode("backend-C")

	ambassador := proxy.New(r)

	mux := http.NewServeMux()

	mux.Handle("/", ambassador)

	fmt.Println("Ambassador listening on :8080")

	err := http.ListenAndServe(":8080", mux)

	if err != nil {

		fmt.Println("Server failed:", err)
	}
}
