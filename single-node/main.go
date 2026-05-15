package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var startedAt time.Time

var totalRequests int

var mutex sync.Mutex

type HealthRespone struct {
	Status string "json:status"
	Uptime string "json:uptime"
}

type StatsResponse struct {
	TotalRequests int "json:total_requests"
}

func incrementRequestCount() {
	mutex.Lock()
	//critical section requires locking for concurrent requests behaviour
	totalRequests++

	mutex.Unlock()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	incrementRequestCount()

	hostname, err := os.Hostname()

	if err != nil {
		hostname = "unknown"
	}

	fmt.Fprintf(w, "Hello from GO HTTP Server running on %s\n", hostname)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {

	incrementRequestCount()

	uptime := time.Since(startedAt)

	response := HealthRespone{
		Status: "ok",
		Uptime: uptime.String(),
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {

	incrementRequestCount()
	//locking before sharing same reading state
	mutex.Lock()

	response := StatsResponse{
		TotalRequests: totalRequests,
	}

	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

}

func main() {

	startedAt = time.Now()

	//they are automatically concurrent go handles them
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/stats", statsHandler)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr: ":" + port,
	}

	go func() {
		fmt.Println("Server starting on port", port)

		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Server Failed:", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	<-quit

	fmt.Println("\n Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	fmt.Println("Shutting down server gracefully...")

	err := server.Shutdown(ctx)

	if err != nil {

		fmt.Println("Graceful shutdown failed:", err)
		return
	}

	fmt.Println("Server exited cleanly")
}
