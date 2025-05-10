package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hannanaarif/crudapi/db"
	"github.com/hannanaarif/crudapi/handlers"
)

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
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Todo API"))
	})

	// Todo routes
	r.Route("/todos", func(r chi.Router) {
		r.Get("/", handlers.GetTodos)
		r.Post("/", handlers.CreateTodo)
		r.Get("/{id}", handlers.GetTodo)
		r.Put("/{id}", handlers.UpdateTodo)
		r.Delete("/{id}", handlers.DeleteTodo)
	})

	port := ":3000"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
