package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {

	backendURL, err := url.Parse("http://localhost:8080")

	if err != nil {

		fmt.Println("Failed to parse backend URL:", err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	fmt.Println("TLS reverse proxy starting on :8443")
	fmt.Println("Forwarding traffic to:", backendURL)

	err = http.ListenAndServeTLS(":8443", "certs/cert.pem", "certs/key.pem", proxy)

	if err != nil {
		fmt.Println("HTTPs server failed:", err)
	}
}
