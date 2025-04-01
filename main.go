package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/souravendra/distributed-rate-limiter-go-crispy-tribble/middleware"
)

// --- Main server ---
func main() {
	http.Handle("/", middleware.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Request allowed: %s\n", time.Now().Format(time.RFC3339))
	})))

	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
