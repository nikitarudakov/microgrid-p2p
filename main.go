package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// RequestBody defines the structure of the expected JSON request
type RequestBody struct {
	Seller   string `json:"seller"`
	Capacity int    `json:"capacity"`
}

// Middleware function to add CORS headers to the response
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow specific methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// Handle preflight request for CORS
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Ensure content-type is application/json
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Decode JSON request body
	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Log the received data
	log.Printf("Received: Seller=%s, Capacity=%d", reqBody.Seller, reqBody.Capacity)

	// Send a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Request received successfully"})
}

func main() {
	port := "8080"
	fmt.Println("Server running on port " + port)

	// Wrap the handler with CORS support and register it
	http.Handle("/", enableCORS(http.HandlerFunc(handler)))

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
