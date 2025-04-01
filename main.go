package main

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)    // Logs all requests
	r.Use(middleware.Recoverer) // Graceful panic recovery
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to my Go API!"))
	})

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte(`{"message": "Hello from Go!"}`))
		json.NewEncoder(w).Encode(map[string]string{"message": "Hello from Go!"})
	})

	http.ListenAndServe(":8080", r)
}
