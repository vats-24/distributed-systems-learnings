package main

import (
	"fmt"
	"net/http"

	"consistent-hash-ambassador/internal/proxy"
	"consistent-hash-ambassador/internal/ring"
)

func main() {

	r := ring.New(10)

	r.AddNode("http://localhost:9001")
	r.AddNode("http://localhost:9002")
	r.AddNode("http://localhost:9003")

	ambassador := proxy.New(r)

	mux := http.NewServeMux()

	mux.Handle("/", ambassador)

	fmt.Println("Ambassador listening on :8080")

	err := http.ListenAndServe(":8080", mux)

	if err != nil {

		fmt.Println("Server failed:", err)
	}
}
