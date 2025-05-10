package main

import (
	"log"
	"net/http"

	"userauth/db"
	"userauth/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome"))
}

func main() {
	// Initialize database connection
	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.DisconnectDB()

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/", Hello)
	r.Post("/signup", handlers.Signup)
	r.Post("/signin", handlers.Signin)

	log.Println("Server is running on http://localhost:3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
