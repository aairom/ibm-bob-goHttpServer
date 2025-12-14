package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	startTime = time.Now()
	version   = "1.0.0"
)

// Response structures
type HealthResponse struct {
	Status  string `json:"status"`
	Uptime  string `json:"uptime"`
	Version string `json:"version"`
}

type InfoResponse struct {
	Version   string    `json:"version"`
	Hostname  string    `json:"hostname"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type EchoResponse struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type DataRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DataResponse struct {
	Success   bool        `json:"success"`
	Data      DataRequest `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type ErrorResponse struct {
	Error     string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}

// Middleware for logging requests
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		log.Printf("Completed in %v", time.Since(start))
	}
}

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Handler functions
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message":   "Welcome to Go HTTP Server!",
		"version":   version,
		"endpoints": "/health, /api/info, /api/echo?message=<text>, /api/data (POST)",
	}
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uptime := time.Since(startTime)

	response := HealthResponse{
		Status:  "healthy",
		Uptime:  uptime.String(),
		Version: version,
	}

	json.NewEncoder(w).Encode(response)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	response := InfoResponse{
		Version:   version,
		Hostname:  hostname,
		Timestamp: time.Now(),
		Message:   "Server information retrieved successfully",
	}

	json.NewEncoder(w).Encode(response)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	message := r.URL.Query().Get("message")
	if message == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:     "Missing 'message' query parameter",
			Timestamp: time.Now(),
		})
		return
	}

	response := EchoResponse{
		Message:   message,
		Timestamp: time.Now(),
	}

	json.NewEncoder(w).Encode(response)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:     "Method not allowed. Use POST",
			Timestamp: time.Now(),
		})
		return
	}

	var req DataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:     "Invalid JSON payload",
			Timestamp: time.Now(),
		})
		return
	}

	response := DataResponse{
		Success:   true,
		Data:      req,
		Timestamp: time.Now(),
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Setup routes with middleware
	http.HandleFunc("/", corsMiddleware(loggingMiddleware(homeHandler)))
	http.HandleFunc("/health", corsMiddleware(loggingMiddleware(healthHandler)))
	http.HandleFunc("/api/info", corsMiddleware(loggingMiddleware(infoHandler)))
	http.HandleFunc("/api/echo", corsMiddleware(loggingMiddleware(echoHandler)))
	http.HandleFunc("/api/data", corsMiddleware(loggingMiddleware(dataHandler)))

	// Create server
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s...", port)
		log.Printf("Server version: %s", version)
		log.Printf("Available endpoints:")
		log.Printf("  GET  /")
		log.Printf("  GET  /health")
		log.Printf("  GET  /api/info")
		log.Printf("  GET  /api/echo?message=<text>")
		log.Printf("  POST /api/data")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}

// Made with Bob
