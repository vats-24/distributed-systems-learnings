package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var startedAt time.Time

func homeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Hello from GO HTTP Server")
}

func main() {

	http.HandleFunc("/", homeHandler)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("Server starting at port", port)

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		fmt.Println("Server failed", err)
	}
}
